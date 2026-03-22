# Sprint 8 规划

**时间**：2026-03-21  
**状态**：待开始  
**Sprint 周期**：2 周（2026-03-21 ～ 2026-04-03）  

---

## 一、Sprint 目标

**目标：** 完成数据权限维度体系建设，实现行级数据隔离能力

在 Sprint 7（角色管理 + 权限分配 UI）的基础上，完成数据权限的后端实现，确保不同角色用户只能访问其授权范围内的业务数据（设备、告警、OTA 任务等）。

---

## 二、详细任务列表

### P0（必须完成）

| # | 任务 | 说明 | 交付物 |
|---|------|------|--------|
| P0-1 | **DataScope 中间件实现** | 所有业务查询自动注入数据权限过滤条件 | `middleware/data_scope.go` |
| P0-2 | **数据权限维度配置表完善** | 完善 `data_scope_config` 表，支持多资源类型 | `models/data_scope_config.go` + migration |
| P0-3 | **用户上下文增强** | 登录时加载用户完整 data_scope 信息写入 Context | `middleware/auth.go` 改造 |
| P0-4 | **Repository 层数据过滤** | 改造所有业务 Repository，自动拼接 WHERE tenant_id + org_id | 改造 device/alert/ota 等 Repository |
| P0-5 | **部门维度数据权限** | data_scope=org / org_and_children 过滤逻辑实现 | 部门级数据隔离 |

### P1（高优先级）

| # | 任务 | 说明 | 交付物 |
|---|------|------|--------|
| P1-1 | **自定义数据权限范围** | data_scope=custom 时，支持 org_ids 自定义范围配置 | 自定义范围 API + UI 配置项 |
| P1-2 | **角色权限继承关系** | 权限合并时 data_scope 取最大范围（all > org_and_children > org > self > custom） | 权限合并逻辑 |
| P1-3 | **API 权限 + 数据权限双重验证** | 同时校验 API 操作权限和数据可见范围 | 双重验证中间件 |
| P1-4 | **设备列表数据权限** | 设备查询自动过滤，只返回权限范围内的设备 | device Repository 改造 |
| P1-5 | **告警列表数据权限** | 告警查询自动过滤，只返回权限范围内的告警 | alert Repository 改造 |

### P2（提升体验）

| # | 任务 | 说明 | 交付物 |
|---|------|------|--------|
| P2-1 | **数据权限配置预览** | 角色详情页展示该角色可访问的数据范围描述 | 角色详情页增强 |
| P2-2 | **OTA 任务数据权限** | OTA 部署任务查询过滤，按创建人/设备归属组织过滤 | ota Repository 改造 |
| P2-3 | **权限变更影响分析** | 角色 data_scope 变更时，提示受影响的数据范围 | 变更提示组件 |

---

## 三、技术方案

### 3.1 架构设计

```
请求进入 API
       │
       ▼
JWT Token 解析 → user_id
       │
       ▼
AuthMiddleware → 查询用户角色 + 权限 → 写入 Context
       │         │
       │         ├─ permissions []        （功能权限列表）
       │         ├─ data_scope_type        （all/org/org_and_children/self/custom）
       │         ├─ org_ids []            （自定义组织范围）
       │         └─ tenant_id
       │
       ▼
DataScopeMiddleware（数据权限过滤）
       │
       ├─ data_scope_type='all' ── 无过滤
       │
       ├─ data_scope_type='org' ── WHERE org_id = :user_org_id
       │
       ├─ data_scope_type='org_and_children' ── WHERE org_id IN (:user_org_id, :children_ids)
       │
       ├─ data_scope_type='self' ── WHERE created_by = :user_id
       │
       └─ data_scope_type='custom' ── WHERE org_id IN (:custom_org_ids)
       │
       ▼
业务 Handler（已自动携带 data_scope 过滤条件）
```

### 3.2 核心接口

| 接口 | 说明 |
|------|------|
| `GET /api/v1/roles/:id/data-scope` | 获取角色数据权限配置 |
| `PUT /api/v1/roles/:id/data-scope` | 更新角色数据权限配置 |
| `GET /api/v1/auth/data-scope` | 获取当前用户数据权限范围 |
| `GET /api/v1/permissions/resource-types` | 获取可配置数据权限的资源类型 |

### 3.3 数据权限配置请求体

```json
{
  "data_scope_type": "org_and_children",
  "resource_configs": [
    { "resource_type": "device", "org_ids": [] },
    { "resource_type": "alert", "org_ids": [] },
    { "resource_type": "ota_deployment", "org_ids": [] }
  ]
}
```

### 3.4 Repository 改造示例

```go
// 改造前
func (r *DeviceRepository) List(filter DeviceFilter) ([]Device, error) {
    var devices []Device
    query := r.db.Where("tenant_id = ?", filter.TenantID)
    // ... 无数据权限过滤
    return devices, nil
}

// 改造后
func (r *DeviceRepository) List(filter DeviceFilter, scope *DataScope) ([]Device, error) {
    var devices []Device
    query := r.db.Where("tenant_id = ?", filter.TenantID)
    
    // 注入数据权限过滤
    if scope != nil && scope.Type != "all" {
        query = scope.ApplyToQuery(query, "org_id", "created_by")
    }
    
    // ... 原有业务过滤条件
    return devices, nil
}
```

### 3.5 关键文件

| 文件 | 说明 |
|------|------|
| `middleware/data_scope.go` | **新建** — DataScope 中间件 |
| `middleware/auth.go` | **改造** — 登录时加载完整 data_scope |
| `models/data_scope_config.go` | **改造** — 支持多资源类型 |
| `repository/device.go` | **改造** — 注入数据权限过滤 |
| `repository/alert.go` | **改造** — 注入数据权限过滤 |
| `repository/ota.go` | **改造** — 注入数据权限过滤 |
| `service/permission.go` | **改造** — data_scope 合并逻辑 |
| `api/v1/role_data_scope.go` | **新建** — 数据权限配置 API |

### 3.6 data_scope 合并规则

多角色权限叠加时，data_scope 取所有角色中范围**最大**的值：

| 角色A | 角色B | 合并后 |
|-------|-------|--------|
| all | 任意 | all |
| org_and_children | org | org_and_children |
| org | self | org |
| self | custom | self（但 custom org_ids 并集） |
| custom | custom | custom（org_ids 取并集） |

---

## 四、验收标准

### 功能验收

| 用例 | 验收条件 | 测试方法 |
|------|----------|----------|
| 部门级数据隔离 | 角色 data_scope=org，用户只能看到本部门数据 | 用该角色账号查询，验证返回数据 org_id 匹配 |
| 部门及下级隔离 | data_scope=org_and_children，可看到本部门及下级部门数据 | 创建多级部门，验证下级数据可见 |
| 仅本人数据 | data_scope=self，只能看到自己创建的数据 | 用非创建人账号查询，验证 created_by 过滤 |
| 自定义范围 | data_scope=custom + 指定 org_ids，只看到指定组织数据 | 配置特定 org_ids，验证过滤正确 |
| 全部数据 | data_scope=all，超级管理员可见所有数据 | super_admin 账号验证 |
| 权限叠加 | 用户有多角色时，data_scope 取最大范围 | 绑定 org + org_and_children 两个角色验证 |
| OTA 任务隔离 | 只能操作本部门/本人的 OTA 部署任务 | 创建其他部门任务，验证无法查看/操作 |
| 告警列表隔离 | 只能查看本部门/本人的告警 | 跨部门验证告警列表过滤 |

### 技术验收

| 验收点 | 标准 |
|--------|------|
| 中间件无侵入 | Repository 改造不影响现有业务逻辑 |
| 性能 | 数据权限过滤不导致查询性能下降 > 10% |
| 缓存 | data_scope 信息支持 Redis 缓存（TTL 5min） |
| 超级管理员 bypass | is_system=super_admin 的用户跳过所有数据权限检查 |

---

## 五、依赖与风险

### 依赖

| 依赖 | 说明 |
|------|------|
| Sprint 7 角色管理 UI | Sprint 7 的角色 CRUD + 权限分配 UI 必须先完成 |
| 组织管理模块 | 数据权限依赖 org_id，需要组织架构已完成 |
| 租户中间件 | 数据权限必须同时过滤 tenant_id |

### 风险

| 风险 | 影响 | 缓解措施 |
|------|------|----------|
| Repository 改造范围大 | 涉及多个业务表的查询改造 | 抽取公共方法，分批改造 |
| 权限叠加逻辑复杂 | 多角色 data_scope 合并可能遗漏边界情况 | 编写UT覆盖所有合并场景 |
| 性能影响 | 每次查询都做 org_id 过滤 | 合理使用索引，必要时加缓存 |

---

## 六、任务分配（初稿）

| 任务 | 负责人 |
|------|--------|
| DataScope 中间件 + Repository 改造 | agenthd |
| 数据权限 API + service 层逻辑 | agenthd |
| 角色数据权限配置 UI | agentqd |
| UT 编写 + 集成测试 | agentcs |
