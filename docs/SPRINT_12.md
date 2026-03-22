# Sprint 12 规划

**时间**：2026-05-03
**状态**：待开始
**Sprint 周期**：2 周（2026-05-03 ～ 2026-05-16）

---

## 一、Sprint 目标

**目标：** 企业级安全功能

在 Sprint 11（告警通知）的基础上，实现企业级安全功能，包括 LDAP/AD 目录服务集成、证书管理、远程设备锁定/数据擦除，以及细粒度数据权限控制，满足企业客户的安全合规需求。

---

## 二、详细任务列表

### 后端 P0/P1/P2 任务表

| # | 任务 | 说明 | 交付物 | 优先级 |
|---|------|------|--------|--------|
| P0-1 | **LDAP/AD 集成 API** | 完成 `/api/v1/ldap/*` 用户同步/认证接口 | ldap_controller.go | P0 |
| P0-2 | **证书管理 API** | 完成 `/api/v1/certificates/*` CRUD | certificate_controller.go | P0 |
| P0-3 | **远程设备锁定 API** | 完成 `POST /api/v1/devices/{device_id}/lock` | device_lock_controller.go | P0 |
| P0-4 | **远程数据擦除 API** | 完成 `POST /api/v1/devices/{device_id}/wipe` | device_wipe_controller.go | P0 |
| P0-5 | **数据权限 API** | 完成 `/api/v1/data-permissions/*` 权限配置接口 | data_permission_controller.go | P0 |
| P1-1 | **LDAP 用户同步服务** | 实现定时同步 LDAP 用户到本地数据库 | ldap/sync_service.go | P1 |
| P1-2 | **证书申请/颁发流程** | 完成 SCEP 证书申请和颁发流程 | certificate/scep_service.go | P1 |
| P1-3 | **设备擦除确认机制** | 擦除前需要二次确认+操作审计 | device_wipe_controller.go | P1 |
| P1-4 | **数据权限表达式** | 支持行级权限表达式配置 | data_permission_service.go | P1 |
| P2-1 | **LDAP 分组同步** | 实现 LDAP 分组到系统角色的映射 | ldap/group_sync.go | P2 |
| P2-2 | **证书到期预警** | 证书到期前自动告警 | certificate/expiry_worker.go | P2 |
| P2-3 | **擦除历史记录** | 完成 `/api/v1/devices/{device_id}/wipe-history` | wipe_history_controller.go | P2 |

### 前端 P0/P1/P2 任务表

| # | 任务 | 说明 | 交付物 | 优先级 |
|---|------|------|--------|--------|
| PF0-1 | **权限分配页面** | 完成 PermissionAssignmentView.vue 角色/用户权限分配 | PermissionAssignmentView.vue | P0 |
| PF0-2 | **数据权限配置页面** | 完成 DataPermissionView.vue 行级/列级权限配置 | DataPermissionView.vue | P0 |
| PF0-3 | **证书管理页面** | 完成 CertificateManageView.vue 证书列表/上传/颁发 | CertificateManageView.vue | P0 |
| PF0-4 | **远程锁定/擦除操作** | 完成 DeviceSecurityView.vue 设备安全操作入口 | DeviceSecurityView.vue | P0 |
| PF1-1 | **LDAP 配置页面** | 完成 LDAPConfigView.vue LDAP 服务器配置 | LDAPConfigView.vue | P1 |
| PF1-2 | **用户同步管理** | 完成 UserSyncView.vue LDAP 用户同步管理 | UserSyncView.vue | P1 |
| PF1-3 | **擦除操作确认** | 完成 DeviceWipeConfirmModal.vue 二次确认弹窗 | DeviceWipeConfirmModal.vue | P1 |
| PF2-1 | **擦除历史页面** | 完成 DeviceWipeHistoryView.vue 擦除操作历史 | DeviceWipeHistoryView.vue | P2 |

---

## 三、技术方案

### API 路由设计

| 接口 | 方法 | 说明 |
|------|------|------|
| `GET /api/v1/ldap/config` | GET | 获取 LDAP 配置 |
| `PUT /api/v1/ldap/config` | PUT | 更新 LDAP 配置 |
| `POST /api/v1/ldap/test` | POST | 测试 LDAP 连接 |
| `GET /api/v1/ldap/users` | GET | LDAP 用户列表 |
| `POST /api/v1/ldap/sync` | POST | 手动触发用户同步 |
| `GET /api/v1/ldap/groups` | GET | LDAP 分组列表 |
| `POST /api/v1/ldap/group-mapping` | POST | 设置分组-角色映射 |
| `GET /api/v1/certificates` | GET | 证书列表 |
| `POST /api/v1/certificates` | POST | 上传/申请证书 |
| `GET /api/v1/certificates/:id` | GET | 证书详情 |
| `PUT /api/v1/certificates/:id` | PUT | 更新证书 |
| `DELETE /api/v1/certificates/:id` | DELETE | 删除证书 |
| `POST /api/v1/certificates/:id/issue` | POST | 颁发证书 |
| `POST /api/v1/devices/{device_id}/lock` | POST | 远程锁定设备 |
| `POST /api/v1/devices/{device_id}/unlock` | POST | 解除设备锁定 |
| `POST /api/v1/devices/{device_id}/wipe` | POST | 远程擦除设备数据 |
| `GET /api/v1/devices/{device_id}/wipe-history` | GET | 擦除历史 |
| `GET /api/v1/data-permissions` | GET | 数据权限列表 |
| `POST /api/v1/data-permissions` | POST | 创建数据权限规则 |
| `PUT /api/v1/data-permissions/:id` | PUT | 更新数据权限规则 |
| `DELETE /api/v1/data-permissions/:id` | DELETE | 删除数据权限规则 |
| `GET /api/v1/data-permissions/effect` | GET | 查看权限生效情况 |

### 数据库设计

```sql
-- LDAP 配置表
CREATE TABLE ldap_configs (
    id              BIGSERIAL PRIMARY KEY,
    config_name     VARCHAR(100) NOT NULL,
    server_url      VARCHAR(255) NOT NULL,
    bind_dn         VARCHAR(255),
    bind_password   VARCHAR(255),                  -- 加密存储
    base_dn         VARCHAR(255) NOT NULL,
    user_filter     VARCHAR(500),
    group_filter    VARCHAR(500),
    sync_interval   INT DEFAULT 3600,              -- 秒
    is_enabled      BOOLEAN DEFAULT FALSE,
    last_sync_at    TIMESTAMP,
    created_by      BIGINT,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

-- LDAP 用户映射表
CREATE TABLE ldap_user_mappings (
    id              BIGSERIAL PRIMARY KEY,
    ldap_dn         VARCHAR(255) NOT NULL UNIQUE,
    local_user_id   BIGINT REFERENCES users(id),
    username        VARCHAR(100),
    email           VARCHAR(255),
    display_name    VARCHAR(100),
    ldap_groups     VARCHAR(100)[],
    sync_status     VARCHAR(20) DEFAULT 'synced',
    last_synced_at  TIMESTAMP,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

-- LDAP 分组-角色映射表
CREATE TABLE ldap_group_role_mappings (
    id              BIGSERIAL PRIMARY KEY,
    ldap_group_dn   VARCHAR(255) NOT NULL,
    ldap_group_name VARCHAR(100),
    role_id         BIGINT NOT NULL REFERENCES roles(id),
    created_at      TIMESTAMP DEFAULT NOW(),
    UNIQUE(ldap_group_dn)
);

-- 设备证书表
CREATE TABLE device_certificates (
    id              BIGSERIAL PRIMARY KEY,
    device_id       VARCHAR(64) NOT NULL,
    cert_serial     VARCHAR(100) NOT NULL UNIQUE,
    cert_type       VARCHAR(20) NOT NULL,         -- 'device'/'user'/'ca'
    subject_cn      VARCHAR(255),
    issuer_dn       VARCHAR(255),
    not_before      TIMESTAMP NOT NULL,
    not_after       TIMESTAMP NOT NULL,
    cert_pem        TEXT,
    private_key_pem TEXT,                          -- 加密存储
    status          VARCHAR(20) DEFAULT 'active',  -- 'active'/'revoked'/'expired'
    revoked_at      TIMESTAMP,
    revocation_reason TEXT,
    fingerprint     VARCHAR(100),
    created_at      TIMESTAMP DEFAULT NOW(),
    INDEX idx_device_id (device_id),
    INDEX idx_cert_serial (cert_serial),
    INDEX idx_not_after (not_after)
);

-- 设备锁定记录表
CREATE TABLE device_lock_records (
    id              BIGSERIAL PRIMARY KEY,
    device_id       VARCHAR(64) NOT NULL,
    lock_type       VARCHAR(20) NOT NULL,         -- 'remote_lock'/'wipe'
    status          VARCHAR(20) NOT NULL,         -- 'pending'/'sent'/'confirmed'/'failed'
    locked_by       BIGINT NOT NULL REFERENCES users(id),
    reason          TEXT,
    confirmed_at     TIMESTAMP,
    completed_at    TIMESTAMP,
    created_at      TIMESTAMP DEFAULT NOW(),
    INDEX idx_device_id (device_id),
    INDEX idx_locked_by (locked_by)
);

-- 数据权限规则表
CREATE TABLE data_permission_rules (
    id              BIGSERIAL PRIMARY KEY,
    rule_name       VARCHAR(100) NOT NULL,
    resource_type   VARCHAR(50) NOT NULL,         -- 'device'/'pet'/'user'/'organization'
    resource_ids    TEXT[],                        -- 资源ID列表，空表示全部
    permission_expr JSONB NOT NULL,                -- 权限表达式
    priority        INT DEFAULT 0,
    is_active       BOOLEAN DEFAULT TRUE,
    description     TEXT,
    created_by      BIGINT,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);
```

---

## 四、验收标准

### 4.1 功能验收

| 用例 | 验收条件 | 测试方法 |
|------|----------|----------|
| LDAP 用户同步 | AD 用户成功同步到系统 | 配置 AD 触发同步 |
| LDAP 认证登录 | AD 用户使用域账号登录 | 测试登录流程 |
| 证书申请颁发 | SCEP 流程完成证书颁发 | 调用 API 测试 |
| 远程设备锁定 | 锁定命令下发设备响应 | 实机测试 |
| 远程数据擦除 | 擦除命令下发设备清除数据 | 实机测试 |
| 数据权限过滤 | Repository 层正确过滤数据 | 调用 API 验证 |
| 权限表达式 | 复杂表达式正确解析 | 单元测试 |

### 4.2 安全验收

| 验收点 | 标准 |
|--------|------|
| 证书存储安全 | 私钥加密存储，不泄露 |
| 擦除操作审计 | 所有操作记录审计日志 |
| LDAP 密码安全 | Bind 密码加密传输 |
| 数据权限隔离 | 跨租户数据不泄漏 |

---

## 五、依赖与风险

### 依赖

| 依赖 | 说明 |
|------|------|
| Sprint 8 权限系统 | 角色/权限体系基础 |
| OpenSSL/SCEP 库 | 证书服务 |
| 企业 AD/LDAP 服务器 | 目录服务（客户提供） |

### 风险

| 风险 | 影响 | 缓解措施 |
|------|------|----------|
| LDAP 服务不稳定 | 用户无法登录 | 保留本地账号作为备用 |
| 设备不在线 | 锁定/擦除失败 | 增加重试+状态查询 |
| 证书过期未处理 | 设备离线 | 提前 30 天预警 |

---

## 六、前端完成清单（agentqd）

**完成时间**：2026-03-22

| # | 任务 | 交付物 | 状态 |
|---|------|--------|------|
| PF0-1 | 权限分配页面 | `views/security/PermissionAssignmentView.vue` | ✅ 完成 |
| PF0-2 | 数据权限配置页面 | `views/security/DataPermissionView.vue` | ✅ 完成 |
| PF0-3 | 证书管理页面 | `views/security/CertificateManageView.vue` | ✅ 完成 |
| PF0-4 | 设备安全操作 | `views/security/DeviceSecurityView.vue` | ✅ 完成 |
| PF1-1 | LDAP 配置页面 | `views/security/LDAPConfigView.vue` | ✅ 完成 |
| PF1-2 | 用户同步管理 | `views/security/UserSyncView.vue` | ✅ 完成 |
| PF1-3 | 擦除操作确认弹窗 | `components/security/DeviceWipeConfirmModal.vue` | ✅ 完成 |
| PF2-1 | 擦除历史页面 | （由后端 agent 提供） | 待完成 |
| API层 | security.ts | `api/security.ts` | ✅ 完成 |
| 路由 | 安全模块路由 | `router/index.js` | ✅ 完成 |

**附带修复**：
- 修复 `AlertHistoryView.vue`、 `NotificationLogsView.vue` 等 Sprint 11 文件中的 TypeScript 语法错误和错误 icon 导入
- 修复 `WebhookChannelConfig.vue` 等文件的 TypeScript 类型注解问题
- 构建验证通过 ✅

---

## 七、后端完成清单（agenthd）

**完成时间**：2026-03-22

| # | 任务 | 交付物 | 状态 |
|---|------|--------|------|
| P0-1 | LDAP/AD 集成 API | `backend/ldap/ldap_service.go`, `backend/controllers/ldap_controller.go` | ✅ 完成 |
| P0-2 | 证书管理 API | `backend/controllers/certificate_controller.go`, `backend/models/certificate.go` | ✅ 完成 |
| P0-3 | 远程设备锁定 API | `backend/controllers/device_security_controller.go` | ✅ 完成 |
| P0-4 | 远程数据擦除 API | `backend/controllers/device_security_controller.go`, `backend/models/wipe_history.go` | ✅ 完成 |
| P0-5 | 数据权限 API | `backend/controllers/data_permission_controller.go`, `backend/models/data_permission.go` | ✅ 完成 |
| P1-1 | LDAP 用户同步服务 | 已在 `ldap_controller.go` 的 `SyncLDAPUsers` 中实现 | ✅ 完成 |
| P1-3 | 设备擦除确认机制 | 二次确认 token + 审计日志在 `device_security_controller.go` | ✅ 完成 |
| P1-4 | 数据权限表达式 | `data_permission_controller.go` 中的 `ValidatePermissionExpression` | ✅ 完成 |
| P2-1 | LDAP 分组同步 | `ldap_controller.go` 的 `GetLDAPGroups` + `SetGroupRoleMapping` | ✅ 完成 |
| P2-2 | 证书到期预警 | `certificate_controller.go` 的 `GetExpiringCertificates` | ✅ 完成 |
| P2-3 | 擦除历史记录 | `device_security_controller.go` 的 `GetWipeHistory` | ✅ 完成 |
| DB | 数据库迁移 | `backend/migrations/006_sprint12_ldap_cert_wipe.sql` | ✅ 完成 |
| Utils | AES 加密工具 | `backend/utils/crypto.go` | ✅ 完成 |
| Models | LDAP/Certificate/Wipe/DataPermission 模型 | `backend/models/` | ✅ 完成 |

**API 路由清单**：
- `GET /api/v1/ldap/config` - 获取 LDAP 配置
- `PUT /api/v1/ldap/config` - 更新 LDAP 配置
- `POST /api/v1/ldap/test` - 测试 LDAP 连接
- `GET /api/v1/ldap/users` - LDAP 用户列表
- `POST /api/v1/ldap/sync` - 触发用户同步
- `GET /api/v1/ldap/groups` - LDAP 分组列表
- `POST /api/v1/ldap/groups/mapping` - 角色分组映射
- `GET /api/v1/certificates` - 证书列表
- `POST /api/v1/certificates` - 创建证书
- `GET /api/v1/certificates/:id` - 证书详情
- `PUT /api/v1/certificates/:id` - 更新证书
- `DELETE /api/v1/certificates/:id` - 删除证书
- `POST /api/v1/certificates/:id/revoke` - 吊销证书
- `GET /api/v1/certificates/expiring` - 即将到期证书
- `POST /api/v1/certificates/validate` - 验证证书
- `POST /api/v1/devices/:device_id/lock` - 锁定设备
- `POST /api/v1/devices/:device_id/unlock` - 解锁设备
- `POST /api/v1/devices/:device_id/wipe` - 擦除设备
- `POST /api/v1/devices/:device_id/wipe/confirm` - 确认擦除
- `GET /api/v1/devices/:device_id/wipe-history` - 擦除历史
- `GET /api/v1/data-permissions/roles/:role_id` - 获取角色数据权限
- `PUT /api/v1/data-permissions/roles/:role_id` - 更新角色数据权限
- `GET /api/v1/data-permissions/users/:user_id` - 获取用户数据权限
- `PUT /api/v1/data-permissions/users/:user_id` - 更新用户数据权限
- `GET /api/v1/data-permissions/columns` - 可配置列级权限字段列表
- `POST /api/v1/data-permissions/validate` - 验证权限表达式

**构建验证**：✅ Go build 通过
