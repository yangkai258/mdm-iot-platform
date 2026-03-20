# 模块 PRD：设备管理 (Device Management)

**版本：** V1.4
**模块负责人：** agentcp
**编制日期：** 2026-03-20

---

## 1. 概述

设备管理是 MDM 中台的核心模块，负责建立全局设备台账，实现 M5Stack 设备从出厂烧录到用户绑定、维修、报废的全链路追踪，并提供远程指令控制能力。

**业务目标：**
- 建立设备唯一标识体系（DeviceID / MAC / SN）
- 追踪设备生命周期状态
- 支持扫码绑定/解绑
- 提供远程指令下发能力

---

## 2. 功能列表

| 功能 | 描述 | 优先级 | 触发方式 | 前端入口/按钮 |
|------|------|--------|----------|--------------|
| 设备注册 | 设备首次上电注册，生成 DeviceID | P0 | 自动 | 无按钮 |
| 设备列表 | 分页查询设备，支持状态/型号/关键词筛选 | P0 | 自动 | 无按钮 |
| 设备详情 | 查看设备完整信息及影子状态 | P0 | 自动 | 无按钮 |
| 扫码绑定 | 通过 SN Code 绑定设备到用户 | P0 | 人工 | 「绑定」按钮 |
| 解绑设备 | 解除设备与用户的绑定关系 | P1 | 人工 | 「解绑」按钮 |
| 设备状态更新 | 更新设备生命周期状态（维修/报废） | P1 | 人工 | 「编辑」按钮 |
| 远程指令下发 | 向设备下发控制指令 | P1 | 人工 | 「指令」按钮 |
| 指令历史查询 | 查询设备历史指令及执行结果 | P2 | 自动 | 无按钮 |
| 设备分组管理 | 设备归属分组，支持树形分组结构 | P1 | 人工 | 「分组」按钮 |
| 设备标签管理 | 设备打标签，支持颜色标签和类型标签 | P1 | 人工 | 「标签」按钮 |
| 批量绑定 | 批量将设备绑定到用户 | P1 | 人工 | 「批量绑定」按钮 |
| 批量解绑 | 批量解绑设备 | P1 | 人工 | 「批量解绑」按钮 |
| 批量状态更新 | 批量更新设备生命周期状态 | P1 | 人工 | 「批量状态更新」按钮 |
| 批量删除 | 批量删除设备（仅限待激活/已报废状态） | P1 | 人工 | 「批量删除」按钮 |
| 设备操作日志 | 记录设备所有操作行为 | P1 | 自动 | 无按钮 |

---

## 3. 数据模型

### 3.1 设备主表 (devices)

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| id | uint | 主键 | PK, auto |
| device_id | string | 全局唯一设备标识 | unique, not null, UUID |
| mac_address | string | 物理 MAC 地址 | unique, not null, format: XX:XX:XX:XX:XX:XX |
| sn_code | string | 产品序列号 | unique, not null |
| hardware_model | string | 硬件型号 | not null |
| firmware_version | string | 当前固件版本 | not null |
| bind_user_id | string | 绑定用户 ID | nullable, FK → sys_users |
| lifecycle_status | int | 生命周期状态 | 1=待激活 2=服役中 3=维修中 4=已挂失 5=已报废 |
| created_at | datetime | 创建时间 | auto |
| updated_at | datetime | 更新时间 | auto |

**lifecycle_status 状态机：**
```
[1:待激活] ──注册──► [2:服役中] ──解绑──► [1:待激活]
                    ├──维修──► [3:维修中] ──修好──► [2:服役中]
                    ├──挂失──► [4:已挂失]
                    └──报废──► [5:已报废]
```

### 3.2 指令历史表 (command_histories)

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| device_id | string | 设备ID |
| cmd_id | string | 指令ID (UUID) |
| cmd_type | string | 指令类型：ota/restart/config/wake |
| action | string | 具体动作描述 |
| status | string | sent/delivered/executed/failed |
| sent_at | datetime | 发送时间 |
| created_at | datetime | 创建时间 |

### 3.3 设备分组表 (device_groups)

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| id | uint | 主键 | PK, auto |
| name | string | 分组名称 | not null |
| description | string | 分组描述 | nullable |
| parent_id | uint | 父分组 ID | nullable, FK → device_groups |
| sort_order | int | 排序权重 | default 0 |
| created_at | datetime | 创建时间 | auto |
| updated_at | datetime | 更新时间 | auto |

### 3.4 设备分组关联表 (device_group_relations)

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| id | uint | 主键 | PK, auto |
| device_id | string | 设备 ID | not null, FK → devices |
| group_id | uint | 分组 ID | not null, FK → device_groups |
| created_at | datetime | 分配时间 | auto |

### 3.5 设备标签表 (device_tags)

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| id | uint | 主键 | PK, auto |
| name | string | 标签名称 | not null |
| color | string | 标签颜色 | not null, HEX like #FF5733 |
| tag_type | string | 标签类型 | 枚举: color(颜色标签) / type(类型标签) |
| created_at | datetime | 创建时间 | auto |
| updated_at | datetime | 更新时间 | auto |

### 3.6 设备标签关联表 (device_tag_relations)

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| id | uint | 主键 | PK, auto |
| device_id | string | 设备 ID | not null, FK → devices |
| tag_id | uint | 标签 ID | not null, FK → device_tags |
| created_at | datetime | 关联时间 | auto |

### 3.7 设备操作日志表 (device_operation_logs)

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| id | uint | 主键 | PK, auto |
| device_id | string | 设备 ID | not null, FK → devices |
| operator_id | uint | 操作人 ID | not null, FK → sys_users |
| operator_name | string | 操作人名称 | not null |
| operation_type | string | 操作类型 | 枚举: bind/unbind/cmd/status_update/delete/group_assign/tag_assign |
| operation_detail | string | 操作详情 | JSON 格式 |
| result | string | 操作结果 | success/failed |
| ip | string | IP 地址 | 操作来源 IP |
| created_at | datetime | 操作时间 | auto |

**operation_type 枚举值：** bind=设备绑定, unbind=设备解绑, cmd=指令下发, status_update=状态变更, delete=设备删除, group_assign=分组分配, tag_assign=标签分配

---

## 4. 接口定义

### 4.1 设备注册

```
POST /api/v1/devices/register
```

**请求参数：** mac_address (MAC地址), sn_code (序列号), hardware_model (硬件型号), firmware_version (固件版本)

**响应：** device_id, lifecycle_status=1, created_at

**错误码：** 4005=MAC格式无效, 5001=服务器内部错误

### 4.2 设备列表

```
GET /api/v1/devices
```

**Query:** page, page_size, status(online/offline), lifecycle_status, hardware_model, search

**响应：** 设备分页列表，含 is_online/battery_level（来自Redis）

### 4.3 设备详情

```
GET /api/v1/devices/:device_id
```

**响应：** 设备完整信息含 shadow(设备影子) + pet_profile(宠物配置)

### 4.4 扫码绑定

```
POST /api/v1/devices/bind/:sn_code
```

**请求体：** { "bind_user_id": "user-uuid-123" }

**响应：** 绑定成功，lifecycle_status: 1→2

**错误码：** 4002=设备不存在, 4003=非法设备状态

### 4.5 解绑设备

```
POST /api/v1/devices/unbind/:sn_code
```

**响应：** 解绑成功，lifecycle_status: 2→1, bind_user_id清空

### 4.6 更新设备状态

```
PUT /api/v1/devices/:device_id/status
```

**请求体：** { "lifecycle_status": 3, "remark": "设备故障" }

**响应：** 更新后的lifecycle_status

### 4.7 删除设备

```
DELETE /api/v1/devices/:device_id
```

**响应：** success

### 4.8 发送设备指令

```
POST /api/v1/devices/:device_id/commands
```

**请求体：** { "cmd_type": "restart", "action": "force_restart", "params": { "delay_seconds": 5 } }

**cmd_type枚举：** ota/restart/config/wake

**响应：** cmd_id, status=sent, sent_at

### 4.9 查询指令历史

```
GET /api/v1/devices/:device_id/commands
```

**Query:** page, page_size

**响应：** 指令历史分页列表，含cmd_id/cmd_type/action/status

### 4.10 设备分组管理接口

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/device/groups | 分组列表（树形结构） |
| POST | /api/v1/device/groups | 创建分组 |
| PUT | /api/v1/device/groups/:id | 更新分组 |
| DELETE | /api/v1/device/groups/:id | 删除分组（含级联） |
| GET | /api/v1/device/groups/:id/devices | 查询分组下设备 |
| POST | /api/v1/devices/:device_id/group | 设备分配分组 |

### 4.11 设备标签管理接口

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/device/tags | 标签列表 |
| POST | /api/v1/device/tags | 创建标签 |
| PUT | /api/v1/device/tags/:id | 更新标签 |
| DELETE | /api/v1/device/tags/:id | 删除标签 |
| POST | /api/v1/devices/:device_id/tags | 设备添加标签 |
| DELETE | /api/v1/devices/:device_id/tags/:tag_id | 设备移除标签 |

### 4.12 批量操作接口

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/v1/devices/batch/bind | 批量绑定（最多100个） |
| POST | /api/v1/devices/batch/unbind | 批量解绑 |
| POST | /api/v1/devices/batch/status | 批量状态更新 |
| POST | /api/v1/devices/batch/delete | 批量删除（仅待激活/已报废） |

### 4.13 设备操作日志接口

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/devices/:device_id/operation-logs | 单设备操作日志 |
| GET | /api/v1/device/operation-logs | 全局操作日志查询 |

**日志记录时机：** bind记录bind_user_id, cmd记录cmd_type/action/cmd_id, status_update记录old/new_status, group_assign记录old/new_group_id, tag_assign记录tag_id/action

---

## 5. 流程图

### 5.1 设备注册流程

```
设备上电 → POST /api/v1/devices/register → 校验MAC格式 → 查询MAC是否存在
  → 已存在: 更新固件版本 → 返回device_id+status
  → 不存在: 创建新记录(生成UUID, lifecycle_status=1) → 返回device_id+lifecycle_status
```

### 5.2 设备绑定流程

```
用户扫码 → 调用 /bind/:sn_code → 查询设备
  → 设备不存在(4002) 或 状态非法(4003) → 返回错误
  → lifecycle_status in [1,2] → 更新bind_user_id, lifecycle_status=2
  → 触发OpenClaw初始化宠物记忆库 → 返回绑定成功
```

---

## 6. 模块联动

| 联动模块 | 联动方式 | 说明 |
|----------|----------|------|
| 设备影子 | 设备注册时创建Redis shadow | 设备影子随注册初始化 |
| 设备影子 | 解绑时清除Redis shadow | 设备解绑后影子失效 |
| OTA升级 | 设备注册后检测OTA任务 | CheckOTA在绑定流程中触发 |
| 会员管理 | 绑定时关联bind_user_id | 设备与会员一对一绑定 |
| 告警系统 | 心跳数据触发CheckAlerts | 设备上报时检查告警规则 |
| 宠物配置 | 绑定后初始化PetProfile | 解绑时宠物配置保留 |
| 策略管理 | 设备绑定时检查合规策略 | 合规检查触发 |
| 应用管理 | target_type=device应用分发 | 设备可安装/卸载应用 |
| 内容管理 | target_type=device内容推送 | 设备可接收文件分发 |
| 通知管理 | MQTT /device/{id}/down/notification | 设备接收推送通知 |
| 数据分析 | 设备操作日志记录 | 触发device_operation_logs |

---

## 7. 验收标准

### P0 验收标准

| 用例 | 验收条件 | 测试方法 |
|------|----------|----------|
| 设备注册 | 首次注册返回device_id，lifecycle_status=1 | POST /api/v1/devices/register |
| 设备注册 | MAC已存在时只更新固件版本，不新建记录 | 同一MAC重复注册 |
| 设备列表 | 支持按lifecycle_status/hardware_model/search筛选 | GET /api/v1/devices?lifecycle_status=2 |
| 设备列表 | 返回数据含is_online/battery_level（来自Redis） | 验证影子字段非空 |
| 设备详情 | 返回设备+影子+宠物配置完整信息 | GET /api/v1/devices/:id |
| 扫码绑定 | 正确设备状态流转1→2 | 验证DB lifecycle_status变化 |
| 解绑设备 | 状态从2→1，bind_user_id清空 | 验证DB字段 |
| MAC格式校验 | 无效格式返回4005错误码 | 提交XX-XX-XX-XX-XX-XX格式 |

### P1 验收标准

| 用例 | 验收条件 | 测试方法 |
|------|----------|----------|
| 设备状态更新 | 支持更新为3/4/5（维修/挂失/报废） | PUT /api/v1/devices/:id/status |
| 远程指令下发 | 正确生成cmd_id并通过MQTT下发 | 查看MQTT Topic消息 |
| 指令历史 | 正确记录每条指令及状态 | 多次下发后查询历史 |
| 设备分组 | 支持树形分组结构和设备分配 | 创建分组并分配设备 |
| 设备标签 | 支持创建标签和设备打标 | 创建标签并关联设备 |
| 批量绑定 | 一次最多绑定100个设备 | 调用批量绑定接口 |
| 批量解绑 | 批量解绑正确更新状态 | 调用批量解绑接口 |
| 批量状态更新 | 批量更新正确写入DB | 调用批量状态更新接口 |
| 批量删除 | 仅允许删除待激活/已报废设备 | 调用批量删除验证 |

---

## 8. UI设计指引

### 8.1 页面结构
- **左侧菜单**：设备管理 → 设备列表 / 设备详情
- **顶部区域**：统计卡片（总数 / 在线数 / 离线数 / 今日新增）
- **中间区域**：数据表格（设备列表，含复选框批量操作）
- **底部区域**：分页组件

### 8.2 组件选用
| 组件 | 用途 |
|------|------|
| a-table | 设备列表主表格，含复选框批量操作 |
| a-card | 顶部统计数字卡片，4列布局 |
| a-drawer | 右侧抽屉：设备详情/编辑/扫码绑定/分组/标签 |
| a-modal | 解绑确认、状态变更二次确认、批量操作确认 |
| a-select | 生命周期状态筛选、硬件型号筛选 |
| a-input-search | 关键词搜索（device_id / sn_code） |
| a-tag | 设备状态标签（绿色=在线，灰色=离线，红色=告警） |
| a-badge | 表格行内在线状态指示灯 |
| a-tree | 设备分组树形结构 |
| a-checkbox | 批量操作复选框 |
| a-button | 各类操作按钮 |

### 8.3 参考模板
```
┌────────────────────────────────────────────────────────────┐
│  [统计卡片]  总设备:1,000  在线:850  离线:150  新增:12  │
├────────────────────────────────────────────────────────────┤
│  [筛选区]  状态▼  型号▼  搜索device_id/sn_code  [+注册设备] │
├────────────────────────────────────────────────────────────┤
│  a-table                                                     │
│  ┌────┬──────────────┬────────┬──────┬─────┬──────┬────┐│
│  │勾选│设备ID         │型号    │固件  │状态 │电量  │操作││
│  ├────┼──────────────┼────────┼──────┼─────┼──────┼────┤│
│  │ ☐  │550e8400-...│CoreS3  │v1.2.0│🟢在线│ 85% │详情││
│  │ ☐  │660e8400-...│StickC  │v1.1.0│⚫离线│ 45% │详情││
│  └────┴──────────────┴────────┴──────┴─────┴──────┴────┘│
│  [◀上一页]  第1/50页  [下一页▶]   每页20条              │
└────────────────────────────────────────────────────────────┘
```

### 8.4 交互流程
```
设备列表页
    ├── 点击「注册设备」──► a-modal ──► 填写MAC/SN/型号 ──► 提交
    ├── 点击「详情」──► a-drawer ──► 显示完整信息+影子+宠物配置
    ├── 点击「绑定」──► a-modal ──► 扫码或输入用户ID ──► 绑定成功
    ├── 点击「解绑」──► a-modal二次确认 ──► 确认后解绑
    ├── 点击「编辑」──► a-drawer ──► 修改信息 ──► 保存
    ├── 点击「指令」──► a-drawer ──► 选择指令类型 ──► 下发并查看结果
    ├── 点击「分组」──► a-drawer ──► 选择分组 ──► 分配成功
    ├── 点击「标签」──► a-drawer ──► 选择标签 ──► 打标成功
    ├── 勾选多个 ──► 「批量绑定」──► a-modal ──► 输入用户ID ──► 批量绑定
    ├── 勾选多个 ──► 「批量解绑」──► a-modal确认 ──► 批量解绑
    ├── 勾选多个 ──► 「批量状态更新」──► a-modal ──► 选择目标状态 ──► 批量更新
    └── 勾选多个 ──► 「批量删除」──► a-modal二次确认 ──► 批量删除
```

### 8.5 关键状态显示
- **在线状态**：a-tag/a-badge，绿色=在线，灰色=离线
- **lifecycle_status**：数字→文字映射，a-tag显示状态
- **电量**：低于20%显示红色警示
- **表格加载态**：a-spin + skeleton骨架屏

---

## 附录 B. 修订记录

| 版本 | 日期 | 修订人 | 修订内容 |
|------|------|--------|----------|
| V1.0 | 2026-03-20 | agentcp | 初稿，基于代码调研 |
| V1.1 | 2026-03-20 | agentcp | 补充设备分组管理、设备标签管理、批量操作、设备操作日志章节 |
| V1.2 | 2026-03-20 | agentcp | 修订功能列表，补充触发方式和前端入口按钮列 |
| V1.4 | 2026-03-20 | agentcp | 重建文档结构，统一使用8章节格式，合并重复的九、UI设计指引章节