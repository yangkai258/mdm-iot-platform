# 模块 PRD：策略与配置管理 (Policy Management)

**版本：** V1.3  
**模块负责人：** agentcp  
**编制日期：** 2026-03-20  

---

## 1. 概述

策略与配置管理是 MDM 中台的核心管控模块，负责统一管理设备配置策略和合规规则，实现设备的规范化管控。设备注册后自动触发合规检查，不合规设备可被隔离或远程擦除。

**业务目标：**
- 配置文件库统一管理（Wi-Fi、VPN、Email、证书、限制策略）
- 策略定义与多维度绑定（设备/用户/组/组织单元）
- 合规策略配置（加密要求、禁止越狱、密码复杂度）
- 不合规自动动作（隔离、擦除、提醒）
- 策略版本管理与回滚

---

## 2. 功能列表

| 功能 | 描述 | 优先级 | 触发方式 | 前端入口/按钮 |
|------|------|--------|----------|--------------|
| 配置文件库 | 配置文件 CRUD（Wi-Fi/VPN/Email/证书/限制） | P0 | 人工 | 「新建配置」/「编辑」/「删除」按钮 |
| 策略定义 | 策略 CRUD、版本管理 | P0 | 人工 | 「新建策略」/「编辑」/「版本历史」按钮 |
| 策略绑定 | 策略绑定到设备/用户/组/组织单元 | P0 | 人工 | 「绑定策略」按钮 |
| 合规规则配置 | 定义合规检查规则（加密/越狱/密码等） | P0 | 人工 | 「新建规则」/「编辑」按钮 |
| 不合规动作 | 不合规时执行动作（隔离/擦除/提醒） | P0 | 自动 | 无按钮 |
| 策略版本回滚 | 回滚到历史策略版本 | P1 | 人工 | 「回滚」按钮 |
| 策略冲突检测 | 检测同一设备的多策略冲突 | P1 | 自动 | 无按钮 |

---

## 3. 数据模型

### 3.1 配置文件表 (policy_configs)

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| id | uint | 主键 | PK, auto |
| name | string | 配置名称 | not null |
| config_type | string | 配置类型 wifi/vpn/email/certificate/restriction | not null |
| config_content | jsonb | 配置内容 | not null |
| hardware_model | string | 适用硬件型号，空表示通用 | nullable |
| is_active | bool | 是否启用 | default true |
| created_by | string | 创建人 | not null |
| created_at | datetime | 创建时间 | auto |
| updated_at | datetime | 更新时间 | auto |

**config_content 示例（Wi-Fi）：**
```json
{
  "ssid": "CompanyWiFi",
  "security_type": "WPA2",
  "password": "***",
  "auto_join": true,
  "proxy_config": {}
}
```

**config_type 枚举：**
- `wifi` - Wi-Fi 配置
- `vpn` - VPN 配置
- `email` - 邮件账户配置
- `certificate` - 证书配置
- `restriction` - 限制策略（如禁止截屏、禁止 App 安装）

### 3.2 策略主表 (policies)

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| id | uint | 主键 | PK, auto |
| policy_code | string | 策略编码, unique | not null |
| policy_name | string | 策略名称 | not null |
| policy_type | string | 类型 device/user/group/org_unit | not null |
| description | string | 策略描述 | nullable |
| priority | int | 优先级 1-100，数字越大优先级越高 | default 50 |
| version | int | 当前版本号 | default 1 |
| config_ids | jsonb | 关联的配置文件 ID 列表 | not null |
| compliance_rules | jsonb | 合规规则配置 | not null |
| remediation_action | string | 不合规动作 quarantine/wipe/notify | default 'notify' |
| is_active | bool | 是否启用 | default true |
| is_system | bool | 是否系统内置策略 | default false |
| created_by | string | 创建人 | not null |
| created_at | datetime | 创建时间 | auto |
| updated_at | datetime | 更新时间 | auto |

**compliance_rules 示例：**
```json
{
  "require_encryption": true,
  "min_password_length": 8,
  "require_alphanumeric": true,
  "jailbreak_detection": true,
  "geofencing": {
    "enabled": false,
    "allowed_regions": ["CN", "US"]
  },
  "max_failed_attempts": 5
}
```

**remediation_action 枚举：**
- `quarantine` - 隔离（限制功能）
- `wipe` - 远程擦除
- `notify` - 仅提醒

### 3.3 策略绑定表 (policy_bindings)

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| id | uint | 主键 | PK, auto |
| policy_id | uint | 策略 ID, FK | not null |
| binding_type | string | 绑定类型 device/user/group/org_unit | not null |
| binding_id | string | 绑定目标 ID | not null |
| is_inherited | bool | 是否继承（上级组织单元继承） | default false |
| created_at | datetime | 创建时间 | auto |

**绑定优先级规则：**
1. 设备级策略 > 用户级策略 > 组级策略 > 组织单元策略
2. 同级别多策略时，按 policy.priority 决出最高优先级

### 3.4 合规规则表 (compliance_rules)

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| id | uint | 主键 | PK, auto |
| rule_code | string | 规则编码, unique | not null |
| rule_name | string | 规则名称 | not null |
| rule_type | string | 类型 encryption/password/jailbreak/geofence | not null |
| condition | jsonb | 检查条件 | not null |
| action | string | 不合规动作 | not null |
| severity | int | 严重程度 1-4 | default 2 |
| enabled | bool | 是否启用 | default true |
| created_at | datetime | 创建时间 | auto |

### 3.5 策略版本历史表 (policy_versions)

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| id | uint | 主键 | PK, auto |
| policy_id | uint | 策略 ID, FK | not null |
| version | int | 版本号 | not null |
| config_ids | jsonb | 配置 ID 列表快照 | not null |
| compliance_rules | jsonb | 合规规则快照 | not null |
| changed_by | string | 变更人 | not null |
| change_reason | string | 变更原因 | nullable |
| created_at | datetime | 创建时间 | auto |

### 3.6 设备合规状态表 (device_compliance_status)

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| id | uint | 主键 | PK, auto |
| device_id | string | 设备 ID, unique | not null |
| policy_id | uint | 生效策略 ID | nullable |
| is_compliant | bool | 是否合规 | default true |
| last_check_at | datetime | 最后检查时间 | nullable |
| last_check_result | jsonb | 最后检查结果详情 | nullable |
| remediation_status | string | 处置状态 none/pending/completed | default 'none' |
| updated_at | datetime | 更新时间 | auto |

---

## 4. 接口定义

### 4.1 配置文件管理

#### 4.1.1 配置文件列表
```
GET /api/v1/policy/configs
```
**Query:** config_type, hardware_model, is_active, page, page_size

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "name": "公司WiFi配置",
        "config_type": "wifi",
        "hardware_model": "",
        "is_active": true,
        "created_at": "2026-03-20T10:00:00Z"
      }
    ],
    "pagination": { "page": 1, "page_size": 20, "total": 10 }
  }
}
```

#### 4.1.2 创建配置文件
```
POST /api/v1/policy/configs
```
**请求体：**
```json
{
  "name": "公司WiFi配置",
  "config_type": "wifi",
  "config_content": {
    "ssid": "CompanyWiFi",
    "security_type": "WPA2",
    "password": "password123",
    "auto_join": true
  },
  "hardware_model": "",
  "is_active": true,
  "created_by": "admin"
}
```

#### 4.1.3 更新配置文件
```
PUT /api/v1/policy/configs/:id
```

#### 4.1.4 删除配置文件
```
DELETE /api/v1/policy/configs/:id
```

---

### 4.2 策略管理

#### 4.2.1 策略列表
```
GET /api/v1/policies
```
**Query:** policy_type, is_active, page, page_size

#### 4.2.2 策略详情
```
GET /api/v1/policies/:id
```
返回策略详情含绑定关系和版本历史

#### 4.2.3 创建策略
```
POST /api/v1/policies
```
**请求体：**
```json
{
  "policy_code": "POLICY_DEVICE_001",
  "policy_name": "标准设备策略",
  "policy_type": "device",
  "description": "所有设备默认执行的策略",
  "priority": 50,
  "config_ids": [1, 2, 3],
  "compliance_rules": {
    "require_encryption": true,
    "min_password_length": 8,
    "jailbreak_detection": true
  },
  "remediation_action": "quarantine",
  "is_active": true,
  "created_by": "admin"
}
```

#### 4.2.4 更新策略
```
PUT /api/v1/policies/:id
```
**说明：** 更新后 version++，旧版本存入 policy_versions

#### 4.2.5 删除策略
```
DELETE /api/v1/policies/:id
```
**限制：** 有绑定关系的策略不可删除

#### 4.2.6 策略版本历史
```
GET /api/v1/policies/:id/versions
```

#### 4.2.7 回滚策略版本
```
POST /api/v1/policies/:id/rollback
```
**请求体：**
```json
{
  "target_version": 2,
  "change_reason": "配置错误，需要回滚"
}
```

---

### 4.3 策略绑定

#### 4.3.1 绑定策略到目标
```
POST /api/v1/policies/:id/bindings
```
**请求体：**
```json
{
  "binding_type": "device",
  "binding_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

#### 4.3.2 解除绑定
```
DELETE /api/v1/policies/:id/bindings/:binding_id
```

#### 4.3.3 查询目标绑定的策略
```
GET /api/v1/policy/targets/:binding_type/:binding_id/policies
```

#### 4.3.4 批量绑定
```
POST /api/v1/policies/batch/bind
```
**请求体：**
```json
{
  "policy_id": 1,
  "binding_type": "group",
  "binding_ids": ["group1", "group2", "group3"]
}
```

---

### 4.4 合规检查

#### 4.4.1 手动触发设备合规检查
```
POST /api/v1/policy/devices/:device_id/compliance-check
```
**说明：** 人工触发，检查设备当前合规状态

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "device_id": "550e8400-e29b-41d4-a716-446655440000",
    "is_compliant": false,
    "violations": [
      {
        "rule_code": "JAILBREAK_001",
        "rule_name": "越狱检测",
        "message": "检测到设备已越狱",
        "severity": 4
      }
    ],
    "remediation_action": "quarantine",
    "checked_at": "2026-03-20T10:30:00Z"
  }
}
```

#### 4.4.2 设备合规状态查询
```
GET /api/v1/policy/devices/:device_id/compliance-status
```

#### 4.4.3 批量合规检查
```
POST /api/v1/policy/compliance-check/batch
```
**请求体：**
```json
{
  "device_ids": ["device1", "device2", "device3"]
}
```

---

### 4.5 合规规则

#### 4.5.1 规则列表
```
GET /api/v1/policy/compliance-rules
```

#### 4.5.2 创建规则
```
POST /api/v1/policy/compliance-rules
```
**请求体：**
```json
{
  "rule_code": "JAILBREAK_001",
  "rule_name": "越狱检测",
  "rule_type": "jailbreak",
  "condition": {
    "check_file_paths": ["/Applications/Cydia.app", "/private/var/lib/apt"],
    "check_sudo": true,
    "check_ssh": true
  },
  "action": "quarantine",
  "severity": 4,
  "enabled": true
}
```

---

## 5. 流程图

### 5.1 设备注册合规检查流程

```
设备注册 POST /api/v1/devices/register
    │
    ▼
查询设备绑定的策略
    │
    ├─→ 无绑定策略 ──► 流程结束，正常服役
    │
    ▼
查询 compliance_rules（所有启用的规则）
    │
    ▼
遍历规则进行合规检查:
    │
    ├─→ require_encryption=true ──► 检查设备存储是否加密
    │       │
    │       ├─→ 合规 ──► 继续下一规则
    │       └─→ 不合规 ──► 记录 violation
    │
    ├─→ jailbreak_detection=true ──► 设备端检测越狱
    │       │
    │       ├─→ 合规 ──► 继续下一规则
    │       └─→ 不合规 ──► 记录 violation
    │
    ├─→ min_password_length ──► 检查密码复杂度
    │       │
    │       ├─→ 合规 ──► 继续下一规则
    │       └─► 不合规 ──► 记录 violation
    │
    └─→ geofencing ──► 检查设备位置是否在允许区域内
            │
            ├─→ 合规 ──► 继续下一规则
            └─→ 不合规 ──► 记录 violation
    │
    ▼
所有规则检查完成
    │
    ├─→ 无 violations ──► is_compliant=true
    │       ──► 更新 device_compliance_status
    │       ──► 流程结束，正常服役
    │
    └─→ 有 violations ──► is_compliant=false
            ──► 更新 device_compliance_status
            ──► 执行 remediation_action:
                    │
                    ├─→ notify ──► 创建告警，通知运营
                    ├─→ quarantine ──► 限制设备功能，通知运营
                    └─→ wipe ──► 远程擦除，通知运营
```

### 5.2 策略绑定与生效流程

```
管理员绑定策略到设备/用户/组/组织单元
    │
    ▼
POST /api/v1/policies/:id/bindings
    │
    ▼
写入 policy_bindings 表
    │
    ▼
设备下次心跳上报时
    │
    ▼
MQTTHandler 获取设备绑定的策略列表
    │
    ▼
按优先级排序，合并配置（去重）
    │
    ▼
通过 /down/desired 下发到设备
    │
    ▼
设备端应用配置
```

### 5.3 策略版本回滚流程

```
管理员选择回滚目标版本
    │
    ▼
POST /api/v1/policies/:id/rollback
    │
    ▼
查询 policy_versions WHERE version = :target_version
    │
    ▼
更新 policies 表:
    ├─→ config_ids = 快照.config_ids
    ├─→ compliance_rules = 快照.compliance_rules
    └─→ version++
    │
    ▼
新增 policy_versions 记录（标记回滚操作）
    │
    ▼
触发已绑定设备重新检查合规
```

---

## 6. 模块联动

| 联动模块 | 联动方式 | 说明 |
|----------|----------|------|
| 设备管理 | 设备注册时自动触发合规检查 | 设备影子提供设备信息 |
| 设备影子 | 策略配置通过 desired_config 下发 | 心跳处理中集成合规检查 |
| 告警系统 | 不合规时创建告警 | remediation_action=notify 时触发 |
| 会员管理 | 可绑定策略到用户/用户组 | 用户维度策略分发 |
| 组织架构 | 可绑定策略到组织单元 | 组织架构维度策略继承 |
| 应用管理 | 策略可包含禁止安装的 App 列表 | restriction 类型配置 |
| 内容管理 | 策略可包含内容访问权限 | 敏感内容保护 |

---

## 7. 验收标准

### P0 验收标准

| 用例 | 验收条件 | 测试方法 |
|------|----------|----------|
| 配置文件 CRUD | Wi-Fi/VPN/Email/证书/限制配置完整增删改查 | 调用各接口验证 |
| 策略 CRUD | 策略创建后正确保存，列表正确 | POST /policies 后 GET 验证 |
| 策略绑定 | 绑定后设备收到对应策略配置 | 检查设备 MQTT 下发消息 |
| 合规检查触发 | 设备注册时自动检查合规规则 | 创建设备后检查合规状态 |
| 不合规动作 | 不合规时按 remediation_action 执行 | 模拟不合规场景验证 |
| 越狱检测 | 设备越狱时创建告警并隔离 | 模拟越狱设备验证 |

### P1 验收标准

| 用例 | 验收条件 | 测试方法 |
|------|----------|----------|
| 策略版本回滚 | 回滚后配置还原到历史版本 | 回滚后检查配置内容 |
| 策略冲突检测 | 同设备多策略时按优先级决出 | 绑定多个策略验证生效结果 |
| 批量绑定 | 一次绑定多个目标 | 调用批量接口验证 |
| 手动合规检查 | 人工触发检查立即返回结果 | POST /compliance-check 验证 |

---

## 8. UI设计指引

### 页面结构
- **左侧菜单**：策略管理 → 配置文件库 / 策略管理 / 合规规则 / 合规状态
- **顶部区域**：策略统计卡片（总策略数 / 绑定中 / 合规率 / 待处置）
- **中间区域**：Tab 页签：配置文件库 / 策略管理 / 合规规则 / 合规状态
- **底部区域**：分页组件

### 组件选用
| 组件 | 用途 |
|------|------|
| a-table | 配置文件列表、策略列表、合规规则列表 |
| a-card | 顶部统计卡片，4列布局 |
| a-tabs | Tab 切换：配置文件库 / 策略管理 / 合规规则 / 合规状态 |
| a-drawer | 创建/编辑配置文件、创建/编辑策略 |
| a-modal | 删除确认、版本回滚确认 |
| a-form | 配置内容表单（不同类型配置有不同表单项）|
| a-input-search | 策略名称/编码搜索 |
| a-select | 配置类型筛选、策略类型筛选、状态筛选 |
| a-tag | 配置类型标签（Wi-Fi=蓝色，VPN=绿色，证书=橙色，限制=红色）|
| a-switch | 启用/禁用开关 |
| a-tree-select | 组织单元选择（绑定目标）|
| a-steps | 策略版本历史时间线 |
| a-badge | 合规状态指示（合规=绿色，不合规=红色）|

### 参考模板
```
┌──────────────────────────────────────────────────────────────┐
│  [统计卡片]  策略总数:20  绑定中:150  合规率:95%  待处置:5    │
├──────────────────────────────────────────────────────────────┤
│  [Tab: 配置文件库 | 策略管理 | 合规规则 | 合规状态]            │
├──────────────────────────────────────────────────────────────┤
│  【策略管理 Tab】                                             │
│  ┌──────────────────────────────────────────────────────┐   │
│  │ [+新建策略]  [筛选: 类型▼  状态▼]                      │   │
│  ├──────────────────────────────────────────────────────┤   │
│  │ 策略名称    │类型  │优先级│配置数│状态│操作            │   │
│  │ 标准设备策略│device│  50  │  3   │启用│详情绑定编辑删 │   │
│  │ 严格安全策略│device│  80  │  5   │启用│详情绑定编辑删 │   │
│  └──────────────────────────────────────────────────────┘   │
│                                                              │
│  【合规规则 Tab】                                             │
│  ┌──────────────────────────────────────────────────────┐   │
│  │ [+新建规则]                                            │   │
│  │ 规则名称    │类型      │严重│状态│操作                │   │
│  │ 越狱检测    │jailbreak│ 4  │启用│编辑禁用删除        │   │
│  │ 加密要求    │encryption│ 3  │启用│编辑禁用删除        │   │
│  │ 地理围栏    │geofence │ 3  │禁用│编辑启用删除        │   │
│  └──────────────────────────────────────────────────────┘   │
└──────────────────────────────────────────────────────────────┘
```

### 交互流程
```
配置文件库页
    │
    ├── 点击「新建配置」──► a-drawer ──► 选择类型 ──► 填写配置内容 ──► 保存
    │
    ├── 点击「编辑」──► a-drawer ──► 修改配置 ──► 保存（自动版本++）
    │
    └── 点击「删除」──► 检查是否有策略引用 ──► 有则提示不可删除

策略管理页
    │
    ├── 点击「新建策略」──► a-drawer ──► 选择类型/优先级/配置/合规规则 ──► 保存
    │
    ├── 点击「绑定」──► a-modal ──► 选择目标类型+搜索目标 ──► 批量绑定
    │
    ├── 点击「详情」──► a-drawer ──► 查看完整策略内容+版本历史
    │
    └── 点击「版本历史」──► a-steps ──► 选择版本 ──► 回滚确认

合规状态页
    │
    ├── 查看所有设备合规状态列表
    │
    ├── 点击「重新检查」──► POST 合规检查 ──► 刷新状态
    │
    └── 点击「处置」──► a-modal ──► 执行 remediation_action
```

### 关键状态显示
- **配置类型**：a-tag，Wi-Fi=蓝，VPN=绿，Email=紫，证书=橙，限制=红
- **策略状态**：a-tag，启用=绿色，禁用=灰色
- **合规状态**：a-badge，合规=绿色 success，不合规=红色 error
- **严重程度**：a-tag，1=蓝，2=黄，3=橙，4=红
- **处置状态**：a-tag，none=灰，pending=黄，completed=绿

---

## 修订记录

| 版本 | 日期 | 修订人 | 修订内容 |
|------|------|--------|----------|
| V1.3 | 2026-03-20 | agentcp | 全新模块，基于新增功能需求重建 |


---

## 9. 页面布局规范

### 9.1 配置文件库页面

**布局结构：**
1. 面包屑 → 页面标题
2. 筛选区（浅灰卡片 #F2F3F5）：配置类型 / 设备型号 / 状态
3. 操作栏（新建配置靠左，其他靠右）
4. 配置文件列表表格

**按钮规范：**
- [新建配置] — 左对齐
- [编辑] [删除] — 行内右对齐

**表格列：**
| 列名 | 宽度 | 说明 |
|------|------|------|
| 序号 | 60px | - |
| 配置名称 | 200px | - |
| 配置类型 | 120px | Wi-Fi/VPN/Email/证书/限制 |
| 设备型号 | 120px | - |
| 状态 | 80px | 启用/禁用 |
| 创建时间 | 150px | - |
| 操作 | 120px | 编辑/删除 |

**分页：** 右下角，10/20/50/100 条

### 9.2 策略管理页面

**布局结构：**
1. 面包屑 → 页面标题
2. 筛选区（浅灰卡片）：策略类型 / 状态
3. 操作栏（新建策略靠左，其他靠右）
4. 策略列表表格

**按钮规范：**
- [新建策略] — 左对齐
- [详情] [绑定] [编辑] [删除] — 行内右对齐

**表格列：**
| 列名 | 宽度 | 说明 |
|------|------|------|
| 序号 | 60px | - |
| 策略名称 | 200px | - |
| 策略类型 | 120px | device/user/group/org_unit |
| 优先级 | 80px | - |
| 配置数 | 80px | - |
| 状态 | 80px | 启用/禁用 |
| 操作 | 120px | 详情/绑定/编辑/删除 |

**分页：** 右下角，10/20/50/100 条

### 9.3 合规规则页面

**布局结构：**
1. 面包屑 → 页面标题
2. 操作栏（新建规则靠左）
3. 合规规则列表表格

**按钮规范：**
- [新建规则] — 左对齐
- [编辑] [启用/禁用] [删除] — 行内右对齐

**表格列：**
| 列名 | 宽度 | 说明 |
|------|------|------|
| 序号 | 60px | - |
| 规则名称 | 200px | - |
| 规则类型 | 120px | encryption/password/jailbreak/geofence |
| 严重程度 | 100px | - |
| 处置动作 | 120px | quarantine/wipe/notify |
| 状态 | 80px | 启用/禁用 |
| 操作 | 120px | 编辑/禁用/删除 |

**分页：** 右下角，10/20/50/100 条

### 9.4 合规状态页面

**布局结构：**
1. 面包屑 → 页面标题
2. 筛选区（浅灰卡片）：合规状态 / 处置状态
3. 操作栏（重新检查靠右）
4. 设备合规状态列表表格

**按钮规范：**
- [重新检查] — 右对齐

**表格列：**
| 列名 | 宽度 | 说明 |
|------|------|------|
| 序号 | 60px | - |
| 设备ID | 150px | - |
| 生效策略 | 150px | - |
| 合规状态 | 100px | 合规/不合规 |
| 最后检查时间 | 150px | - |
| 违规项 | 200px | - |
| 处置状态 | 120px | none/pending/completed |

**分页：** 右下角，10/20/50/100 条

### 9.5 弹窗规范

| 类型 | 使用场景 |
|------|----------|
| Drawer 抽屉 | 新建/编辑配置文件、新建/编辑策略、合规详情 |
| Dialog 对话框 | 删除确认、版本回滚确认、策略绑定 |
| 全屏模态 | 暂无复杂表单场景 |
