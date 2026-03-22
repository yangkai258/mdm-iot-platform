# 模块 PRD：系统管理 (System Management)

**版本：** V1.4
**模块负责人：** agentcp
**编制日期：** 2026-03-20

---

## 1. 概述

系统管理模块负责 MDM 中台的后台运营人员管理、权限控制和系统配置，包括用户管理、角色权限、菜单权限、字典管理和操作日志。

**业务目标：**
- 运营人员账号管理
- RBAC角色权限控制
- 精细化的菜单/按钮权限
- 操作日志审计
- 统一的字典配置

---

## 2. 功能列表

| 功能 | 描述 | 优先级 | 触发方式 | 前端入口/按钮 |
|------|------|--------|----------|--------------|
| 用户管理 | 用户CRUD、状态管理 | P0 | 人工 | 「新建用户」/「编辑」/「删除」/「禁用」/「重置密码」按钮 |
| 角色管理 | 角色CRUD、权限分配 | P0 | 人工 | 「新建角色」/「编辑」/「删除」/「权限配置」按钮 |
| 菜单管理 | 菜单树CRUD、权限标识 | P1 | 人工 | 「新增菜单」/「编辑」/「删除」按钮 |
| 权限管理 | 按钮/操作权限CRUD | P1 | 人工 | 「新建权限」/「编辑」/「删除」按钮 |
| 登录日志 | 登录历史查询 | P1 | 自动 | 无按钮 |
| 操作日志 | 操作记录查询、审计 | P1 | 自动 | 无按钮 |
| 字典管理 | 统一字典配置 | P2 | 人工 | 「新建字典」/「编辑」/「删除」按钮 |

---

## 3. 数据模型

### 3.1 用户表 (sys_users)

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| username | string | 用户名, unique |
| password | string | 密码(加密) |
| nickname | string | 昵称 |
| email | string | 邮箱 |
| phone | string | 手机号 |
| status | int | 1=正常 0=禁用 |
| role_id | uint | 角色ID |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |

### 3.2 角色表 (sys_roles)

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| name | string | 角色名称, unique |
| code | string | 角色编码, unique |
| description | string | 描述 |
| status | int | 1=正常 0=禁用 |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |

### 3.3 菜单表 (sys_menus)

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| parent_id | uint | 上级菜单, 0=顶级 |
| name | string | 菜单名称 |
| path | string | 路由路径 |
| component | string | 组件路径 |
| icon | string | 图标 |
| sort | int | 排序 |
| visible | int | 1=显示 0=隐藏 |
| permission | string | 权限标识 |
| type | int | 1=菜单 2=按钮 |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |

### 3.4 操作日志表 (sys_operation_logs)

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| user_id | uint | 用户ID |
| username | string | 用户名 |
| module | string | 模块 |
| operation | string | 操作类型 |
| method | string | 请求方法 GET/POST/PUT/DELETE |
| path | string | 请求路径 |
| ip | string | IP地址 |
| location | string | 地理位置 |
| params | text | 请求参数 |
| result | text | 响应结果 |
| status | int | 1=成功 0=失败 |
| error_msg | text | 错误信息 |
| duration | int | 耗时(ms) |
| created_at | datetime | 创建时间 |

### 3.5 登录日志表 (sys_login_logs)

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| user_id | uint | 用户ID |
| username | string | 用户名 |
| ip | string | IP地址 |
| location | string | 地理位置 |
| browser | string | 浏览器 |
| os | string | 操作系统 |
| status | int | 1=成功 0=失败 |
| msg | string | 消息 |
| login_time | datetime | 登录时间 |

### 3.6 字典表 (sys_dictionaries)

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| type | string | 字典类型, unique |
| name | string | 字典名称 |
| label | string | 字典标签 |
| value | string | 字典值 |
| sort | int | 排序 |
| status | int | 1=正常 0=禁用 |
| remark | string | 备注 |

### 3.7 设备事件日志表 (device_event_logs)

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| device_id | string | 设备ID, index |
| event_type | string | 事件类型 boot/connect/disconnect/error/update |
| event_data | jsonb | 事件数据详情 |
| severity | string | 严重程度 info/warning/error |
| reported_at | datetime | 设备上报时间 |
| created_at | datetime | 服务端记录时间 |

**event_type枚举：** boot=设备启动, connect=MQTT连接, disconnect=MQTT断开, error=设备错误, update=固件更新, app_install=App安装, app_uninstall=App卸载

### 3.8 APNs配置表 (apns_config)

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| environment | string | 环境 development/production |
| key_id | string | APNs Key ID |
| team_id | string | Apple Team ID |
| bundle_id | string | Bundle Identifier |
| key_file | string | APNs私钥文件路径 |
| is_active | bool | 是否启用 |
| updated_at | datetime | 更新时间 |

---

## 4. 接口定义

### 4.1 认证接口

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/v1/auth/login | 登录（body: username, password → token+user info） |
| POST | /api/v1/auth/logout | 登出 |
| GET | /api/v1/auth/current | 当前用户信息 |

### 4.2 用户管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/system/users | 用户列表（Query: keyword, status, role_id, page, page_size） |
| POST | /api/v1/system/users | 创建用户（username, password, nickname, email, phone, role_id） |
| PUT | /api/v1/system/users/:id | 更新用户 |
| DELETE | /api/v1/system/users/:id | 删除用户 |
| PUT | /api/v1/system/users/:id/password | 重置密码（body: password） |
| PUT | /api/v1/system/users/:id/status | 启用/禁用用户（body: status） |

### 4.3 角色管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/roles | 角色列表 |
| POST | /api/v1/roles | 创建角色（name, code, description） |
| PUT | /api/v1/roles/:id | 更新角色 |
| DELETE | /api/v1/roles/:id | 删除角色 |
| GET | /api/v1/roles/:id/perms | 获取角色权限 → permissions数组 |
| PUT | /api/v1/roles/:id/perms | 设置角色权限（body: permissions数组） |

### 4.4 菜单管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/system/menus/tree | 菜单树 |
| POST | /api/v1/system/menus | 创建菜单（parent_id, name, path, component, icon, sort, type, permission） |

### 4.5 权限管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/permissions | 权限列表 |
| POST | /api/v1/permissions | 创建权限（name, permission, type） |

### 4.6 日志管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/system/logs/operation | 操作日志（Query: user_id, module, start_time, end_time, page, page_size） |
| GET | /api/v1/system/logs/login | 登录日志（Query: username, start_time, end_time, page, page_size） |

---

## 5. 流程图

### 5.1 登录流程

```
用户提交username + password
    │
    ▼
查询sys_users WHERE username=:username
    │
    ├─► 用户不存在 ──返回401
    ├─► 用户已禁用 ──返回403
    └─► 验证密码bcrypt.Compare
            │
            ├─► 密码错误 ──记录登录日志(status=0)──► 返回401
            └─► 密码正确 ──生成JWT token(24h) ──记录登录日志(status=1) ──► 返回token+user info
```

### 5.2 权限校验流程

```
请求带着JWT token → JWT Middleware解析token
    │
    ├─► token无效/过期 ──返回401
    └─► token有效 → 提取user_id → 查询用户角色 → 查询角色权限
    │
    ▼
检查请求路径+方法是否在权限列表中
    │
    ├─► 有权限 ──继续处理请求
    └─► 无权限 ──返回403 Forbidden
```

---

## 6. 模块联动

| 联动模块 | 联动方式 | 说明 |
|----------|----------|------|
| 设备管理 | 运营人员管理设备 | 权限控制 |
| OTA升级 | 运营人员管理OTA | 权限控制 |
| 会员管理 | 运营人员管理会员 | 权限控制 |
| 告警系统 | 运营人员处理告警 | 权限控制 |
| 组织架构 | 用户关联员工/部门 | UserExt表 |

---

## 7. 验收标准

### P0 验收标准

| 用例 | 验收条件 | 测试方法 |
|------|----------|----------|
| 用户登录 | 正确账号密码返回token | POST /auth/login |
| JWT校验 | 无效token返回401 | 携带过期token请求 |
| 用户CRUD | 完整增删改查 | 调用各接口验证 |
| 角色CRUD | 完整增删改查 | 调用各接口验证 |
| 角色权限设置 | 设置后用户获得对应权限 | 设置权限后验证 |

### P1 验收标准

| 用例 | 验收条件 | 测试方法 |
|------|----------|----------|
| 菜单树 | 正确返回树形结构 | GET /menus/tree |
| 权限校验 | 无权限接口返回403 | 无权限角色调用受限接口 |
| 登录日志 | 每次登录正确记录 | 登录后查询日志 |
| 操作日志 | 请求正确记录操作 | 操作后查询日志 |

### P2 验收标准

| 用例 | 验收条件 | 测试方法 |
|------|----------|----------|
| 字典管理 | 支持多类型字典 | 创建/查询字典 |
| 用户禁用 | 禁用后无法登录 | 禁用用户后尝试登录 |

---

## 8. UI设计指引

### 8.1 页面结构
- **左侧菜单**：系统管理 → 用户管理 / 角色管理 / 菜单管理 / 权限管理 / 日志管理
- **顶部区域**：用户统计卡片（总用户数/在线用户/今日登录/角色数）
- **中间区域**：左侧树形菜单 + 右侧数据表格
- **底部区域**：分页组件

### 8.2 组件选用
| 组件 | 用途 |
|------|------|
| a-table | 用户列表、角色列表、菜单列表、权限列表、日志列表 |
| a-tree | 左侧菜单树（系统管理展开）|
| a-drawer | 创建/编辑用户、创建/编辑角色、菜单配置 |
| a-modal | 删除确认、重置密码确认、角色权限配置 |
| a-form | 用户表单、角色表单、菜单表单 |
| a-input-search | 用户名/手机号搜索 |
| a-select | 角色筛选、状态筛选 |
| a-switch | 用户启用/禁用、角色启用/禁用 |
| a-cascader | 部门级联选择（用户归属部门）|
| a-transfer | 角色权限配置（左侧可选权限 → 右侧已选权限）|
| a-password | 重置密码输入框 |

### 8.3 参考模板
```
┌──────────────────────────────────────────────────────────────┐
│  [统计卡片]  用户总数:50  在线:12  今日登录:8  角色数:5    │
├──────────┬───────────────────────────────────────────────────┤
│  系统管理 │                                                    │
│  ├─用户管理│  【用户列表】                                      │
│  ├─角色管理│  ┌──────────────────────────────────────────┐   │
│  ├─菜单管理│  │ [关键词搜索]  角色▼  [+新建用户]          │   │
│  ├─权限管理│  ├──────────────────────────────────────────┤   │
│  └─日志管理│  │用户名  │昵称  │角色   │状态│最后登录│操作  │   │
│            │  │ admin  │管理员│超级管理│ 正常│ 10:30 │编辑删│   │
│            │  │ op01   │运营  │运营   │ 正常│ 09:15 │编辑删│   │
│            │  └──────────────────────────────────────────┘   │
│            │                                                    │
│            │  【角色权限配置 a-modal】                          │
│            │  ┌──────────────────────────────────────────┐   │
│            │  │ 角色名称: [设备管理员      ]              │   │
│            │  │ 权限配置:                                 │   │
│            │  │  ◄ [device:list >>] [device:bind >>]   │   │
│            │  │    [device:unbind >>]  [ota:list >>]   │   │
│            │  │                    [确定]  [取消]        │   │
│            │  └──────────────────────────────────────────┘   │
└──────────┴───────────────────────────────────────────────────┘
```

### 8.4 交互流程
```
用户管理页
    ├── 点击「新建用户」──► a-drawer ──► 填写信息+选择角色+选择部门 ──► 创建
    ├── 点击「编辑」──► a-drawer ──► 修改信息 ──► 保存
    ├── 点击「禁用」──► a-modal确认 ──► 用户不可登录
    └── 点击「重置密码」──► a-modal ──► 确认 ──► 密码重置为默认值

角色管理页
    ├── 点击「新建角色」──► a-drawer ──► 填写角色信息 ──► 创建
    ├── 点击「权限配置」──► a-modal ──► a-transfer选择权限 ──► 保存
    └── 点击「删除」──► 检查是否有用户归属该角色 ──► 不可删除提示

菜单管理页
    ├── 左侧树形展示菜单层级
    ├── 点击「新增菜单」──► a-drawer ──► 选择上级+填写信息 ──► 保存
    └── 点击「编辑/删除」──► 操作后刷新树形结构

日志管理页
    ├── 操作日志 / 登录日志Tab切换
    └── 筛选条件：时间范围+用户+模块 ──► 查询结果
```

### 8.5 关键状态显示
- **用户状态**：a-switch，启用=绿色，禁用=灰色
- **角色状态**：a-tag，正常=绿色，禁用=灰色
- **登录状态**：a-tag，成功=绿色，失败=红色
- **权限树**：a-tree-select，支持多选父子联动

---

## 附录 C. 修订记录

| 版本 | 日期 | 修订人 | 修订内容 |
|------|------|--------|----------|
| V1.0 | 2026-03-20 | agentcp | 初稿，基于代码调研 |
| V1.2 | 2026-03-20 | agentcp | 修订功能列表，补充触发方式和前端入口按钮列 |
| V1.4 | 2026-03-20 | agentcp | 重建文档结构，统一使用8章节格式，将关键问题移至附录