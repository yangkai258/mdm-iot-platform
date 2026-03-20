# 模块 PRD：设备影子 (Device Shadow)

**版本：** V1.4
**模块负责人：** agentcp
**编制日期：** 2026-03-20

---

## 1. 概述

设备影子是 MDM 中台的"状态缓冲池"，在服务器端维护 M5Stack 设备的实时物理快照，解决网络高延迟或断网时状态不同步的问题。

**业务目标：**
- 缓存设备最新状态（在线/离线/电量/模式）
- 90秒心跳超时自动标记离线
- 支持期望状态（desired_state）下发到设备
- 为告警系统提供实时数据源

---

## 2. 功能列表

| 功能 | 描述 | 优先级 | 触发方式 | 前端入口/按钮 |
|------|------|--------|----------|--------------|
| 心跳上报处理 | 接收MQTT心跳，更新Redis影子 | P0 | 自动 | 无按钮 |
| 离线检测 | 90秒未收到心跳自动标记离线 | P0 | 自动 | 无按钮 |
| 实时状态查询 | HTTP接口查询设备影子 | P0 | 自动 | 无按钮 |
| 期望状态下发 | MQTT下发desired_config到设备 | P1 | 人工 | 「保存配置」按钮 |
| 宠物模式管理 | current_mode状态读写 | P1 | 人工 | 「保存配置」按钮 |
| 设备影子与DB同步 | 离线时同步更新PostgreSQL | P0 | 自动 | 无按钮 |

---

## 3. 数据模型

### 3.1 设备影子Redis结构

**Key:** `shadow:{device_id}`
**TTL:** 90秒（每次心跳刷新）

```json
{
  "device_id": "550e8400-e29b-41d4-a716-446655440000",
  "is_online": true,
  "battery_level": 85,
  "current_mode": "idle",
  "last_ip": "192.168.1.100",
  "rssi": -55,
  "last_heartbeat": "2026-03-20T10:30:00Z",
  "desired_config": {
    "dnd_start_time": "23:00",
    "dnd_end_time": "08:00",
    "desired_firmware": "v1.3.0"
  }
}
```

### 3.2 设备影子表 (device_shadows) — PostgreSQL

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| device_id | string | 设备唯一标识 | PK, not null |
| is_online | bool | 在线状态 | default false |
| battery_level | int | 电量0-100 | - |
| current_mode | string | 当前模式 | default 'idle' |
| last_ip | string | 最后公网IP | nullable |
| last_heartbeat | datetime | 最后心跳时间 | index |
| desired_config | jsonb | 期望配置 | 免打扰/目标固件 |

**current_mode 枚举：** sleeping=休眠模式, roaming=漫游模式, listening=倾听模式, talking=对话模式, idle=空闲

### 3.3 宠物配置表 (pet_profiles)

| 字段 | 类型 | 说明 |
|------|------|------|
| device_id | string | 设备ID, PK |
| pet_name | string | 宠物名字，默认Mimi |
| personality | string | 性格 lively/cool/angry/clingy |
| interaction_freq | string | 交互频率 high/medium/low |
| dnd_start_time | string | 免打扰开始时间 HH:MM |
| dnd_end_time | string | 免打扰结束时间 HH:MM |
| custom_rules | jsonb | 自定义规则 |

---

## 4. 接口定义

### 4.1 查询设备影子

```
GET /api/v1/devices/:device_id/shadow
```

**响应：** Redis优先，若无数据则查PostgreSQL

```json
{
  "code": 0,
  "data": {
    "device_id": "xxx",
    "is_online": true,
    "battery_level": 85,
    "current_mode": "idle",
    "last_ip": "192.168.1.100",
    "last_heartbeat": "2026-03-20T10:30:00Z",
    "desired_config": {
      "dnd_start_time": "23:00",
      "dnd_end_time": "08:00"
    }
  }
}
```

### 4.2 更新期望状态

```
PUT /api/v1/devices/:device_id/shadow/desired
```

**请求体：**
```json
{
  "desired_config": {
    "dnd_start_time": "23:00",
    "dnd_end_time": "08:00",
    "desired_firmware": "v1.3.0",
    "interaction_freq": "low"
  }
}
```

**响应：** 更新后的desired_config + mqtt_published=true

**说明：** 更新后通过MQTT下发到设备

### 4.3 宠物配置读取

```
GET /api/v1/devices/:device_id/profile
```

**响应：** pet_name/personality/interaction_freq/dnd/custom_rules

### 4.4 宠物配置更新

```
PUT /api/v1/devices/:device_id/profile
```

**请求体：**
```json
{
  "pet_name": "Mimi Pro",
  "personality": "cool",
  "interaction_freq": "low",
  "dnd_start_time": "22:30",
  "dnd_end_time": "07:30",
  "custom_rules": { "max_talk_per_day": 10 }
}
```

---

## 5. 流程图

### 5.1 MQTT消息格式

#### 5.1.1 设备上报心跳 (订阅)

**Topic:** `/device/{device_id}/up/status`

**Payload:**
```json
{
  "device_id": "550e8400-e29b-41d4-a716-446655440000",
  "battery": 85,
  "mode": "idle",
  "ip": "192.168.1.100",
  "rssi": -55,
  "timestamp": 1710912600
}
```

#### 5.1.2 下发期望状态 (发布)

**Topic:** `/device/{device_id}/down/desired`

**Payload:**
```json
{
  "desired_config": {
    "dnd_start_time": "23:00",
    "dnd_end_time": "08:00",
    "interaction_freq": "low"
  },
  "timestamp": 1710912600
}
```

#### 5.1.3 下发设备指令 (发布)

**Topic:** `/device/{device_id}/down/cmd`

**Payload:**
```json
{
  "cmd_id": "cmd-uuid-456",
  "cmd_type": "restart",
  "action": "force_restart",
  "params": { "delay_seconds": 5 },
  "timestamp": 1710912600
}
```

### 5.2 心跳处理流程

```
MQTT Broker收到 /device/{id}/up/status
    │
    ▼
StatusMessageHandler()
    ├─► 解析JSON payload
    ├─► 校验device_id存在
    ▼
Redis SET shadow:{device_id} (TTL 90s)
    │  is_online=true, battery_level, current_mode, last_ip, last_heartbeat=now()
    ▼
CheckAlerts(db, device_id, payload)
    │  电池<15%? ──是──► 创建battery_low告警
    │  离线状态? ──是──► 创建offline告警
    ▼
检查desired_config是否存在
    ├─► 存在 ──MQTT下发──► /device/{id}/down/desired
    ▼
流程结束
```

### 5.3 离线检测流程

```
OTA Worker (每30s轮询)
    │
    ▼
遍历所有shadow:* keys (使用SCAN，非KEYS)
    │
    ▼
检查last_heartbeat距今是否>90s
    │
    ├─► 超时 ── shadow.is_online = false
    │       ── 写回Redis (TTL永久)
    │       ── 异步更新PostgreSQL
    │       ── 创建offline告警
    │
    └─► 未超时 ── 跳过
```

---

## 6. 模块联动

| 联动模块 | 联动方式 | 说明 |
|----------|----------|------|
| 设备管理 | 设备注册时创建影子 | 设备下线时同步DB |
| OTA升级 | 期望固件版本在desired_config | OTA Worker下发指令 |
| 告警系统 | 心跳数据触发CheckAlerts | 电量低/离线告警 |
| 宠物配置 | desired_config含DND时间 | 设备端拉取应用 |
| 数据分析 | 统计在线/离线设备数量 | Redis SCAN计数 |
| 策略管理 | 合规策略通过desired_config下发 | 策略绑定后心跳时推送 |
| 策略管理 | 心跳时触发合规检查CheckCompliance | - |
| 应用管理 | 应用分发指令通过desired_config下发 | App安装/卸载指令 |
| 内容管理 | 内容分发通过desired_config下发 | 文档推送指令 |

---

## 7. 验收标准

### P0 验收标准

| 用例 | 验收条件 | 测试方法 |
|------|----------|----------|
| 心跳上报 | 设备发送MQTT后Redis shadow正确更新 | 模拟MQTT消息，检查Redis |
| TTL刷新 | 连续心跳保持is_online=true | 持续发送心跳，检查TTL |
| 离线检测 | 90秒无心跳is_online=false | 停止发送心跳90s后检查 |
| 影子查询 | HTTP接口正确返回影子数据 | GET /api/v1/devices/:id/shadow |
| DB同步 | 离线时PostgreSQL last_heartbeat更新 | 检查DB记录 |
| CheckAlerts调用 | 电量<15%创建告警记录 | 模拟低电量心跳 |

### P1 验收标准

| 用例 | 验收条件 | 测试方法 |
|------|----------|----------|
| 期望状态下发 | PUT desired后MQTT消息正确 | 检查/down/desired Topic |
| 宠物配置更新 | 修改DND时间后设备收到 | 设备端日志验证 |
| 模式切换 | current_mode正确更新 | 模拟不同模式上报 |

---

## 8. UI设计指引

### 8.1 页面结构
- **设备详情页 → Tab切换**：「基本信息」/「实时状态」/「宠物配置」/「指令历史」
- **顶部区域**：设备基本信息卡片
- **中间区域**：实时状态可视化（电量进度条、在线指示灯、最后心跳时间）
- **宠物配置Tab**：a-form表单编辑

### 8.2 组件选用
| 组件 | 用途 |
|------|------|
| a-card | 设备基本信息卡片、实时状态卡片 |
| a-descriptions | 设备影子键值对展示 |
| a-progress | 电量进度条（绿>50%，黄20-50%，红<20%）|
| a-badge | 在线/离线状态指示 |
| a-form | 宠物配置编辑表单 |
| a-switch | 免打扰开关 |
| a-time-picker | DND开始/结束时间选择 |
| a-tag | 当前模式标签（idle/sleeping/roaming等）|
| a-tabs | Tab切换：基本信息/实时状态/宠物配置/指令历史 |

### 8.3 参考模板
```
┌──────────────────────────────────────────────────────────────┐
│  设备详情: 550e8400-e29b-41d4-a716-446655440000              │
├──────────────────────────────────────────────────────────────┤
│  [Tab: 基本信息 | 实时状态 | 宠物配置 | 指令历史]              │
├──────────────────────────────────────────────────────────────┤
│  【实时状态Tab】                                             │
│  ┌─────────────────────┐  ┌─────────────────────┐          │
│  │ 电量    [====85%]  │  │ 🟢 在线  已在线2h35m │          │
│  └─────────────────────┘  └─────────────────────┘          │
│  ┌─────────────────────┐  ┌─────────────────────┐          │
│  │ 当前模式  idle     │  │ 信号强度  -55dBm    │          │
│  └─────────────────────┘  └─────────────────────┘          │
│  最后心跳: 2026-03-20 10:30:00  最后IP: 192.168.1.100       │
│                                                               │
│  【宠物配置Tab】                                             │
│  ┌─────────────────────────────────────────────┐            │
│  │ 宠物名称: [Mimi       ]  性格: [活泼 ▼]    │            │
│  │ 交互频率: [中等 ▼]                          │            │
│  │ 免打扰:   [22:30 - 07:30]  [开关:ON]      │            │
│  │                    [保存配置]               │            │
│  └─────────────────────────────────────────────┘            │
└──────────────────────────────────────────────────────────────┘
```

### 8.4 交互流程
```
设备详情页 → 切换「实时状态」Tab
    ├── 页面加载 ──► 从Redis获取shadow数据 ──► 实时展示
    ├── 实时数据刷新 ──► 每30s自动轮询shadow ──► 更新页面
    └── 电量<20% ──► 进度条变红 + 告警提示

设备详情页 → 切换「宠物配置」Tab
    ├── 加载宠物配置数据 ──► 展示a-form
    ├── 修改配置 ──► a-form自动校验
    └── 点击「保存」──► PUT /devices/:id/profile ──► 成功提示
```

### 8.5 关键状态显示
- **在线状态**：a-badge，绿色=status="success"=在线，灰色=status="default"=离线
- **电量进度条**：a-progress，stroke-color按阈值变色
- **模式标签**：a-tag，显示current_mode枚举值
- **DND时间**：a-time-picker HH:mm格式，配合a-switch启用/禁用

---

## 附录 C. 修订记录

| 版本 | 日期 | 修订人 | 修订内容 |
|------|------|--------|----------|
| V1.0 | 2026-03-20 | agentcp | 初稿，基于代码调研 |
| V1.2 | 2026-03-20 | agentcp | 修订功能列表，补充触发方式和前端入口按钮列 |
| V1.4 | 2026-03-20 | agentcp | 重建文档结构，统一使用8章节格式，将MQTT消息格式并入流程图章节