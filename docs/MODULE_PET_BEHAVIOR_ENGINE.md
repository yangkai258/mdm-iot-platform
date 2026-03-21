# 宠物行为引擎

**版本：** V1.1  
**模块负责人：** agentcp  
**编制日期：** 2026-03-22  

---

## 1. 概述

宠物行为引擎是OpenClaw AI层的核心执行模块，负责将AI决策转化为MiniClaw设备可执行的动作序列。行为引擎通过决策树和优先级管理，协调传感器的实时输入与AI的语义理解，生成流畅、自然的宠物行为表现。

**业务目标：**
- 将AI语义指令转换为设备可执行的动作序列
- 处理传感器实时数据（防跌落/避障）
- 管理多动作并发与优先级冲突
- 确保行为安全性和响应实时性

---

## 2. 功能列表

| 功能 | 描述 | 优先级 | 触发方式 | 前端入口/按钮 |
|------|------|--------|----------|--------------|
| 动作序列规划 | 将AI指令分解为走/停/看/说/听等基础动作序列 | P0 | 自动 | 无按钮 |
| 决策树执行 | 基于规则的决策引擎，处理条件分支和状态转换 | P0 | 自动 | 无按钮 |
| 传感器数据处理 | 实时接收并处理防跌落/避障/触摸等传感器数据 | P0 | 自动 | 无按钮 |
| 动作优先级管理 | 紧急动作（如跌落检测）优先于普通动作 | P0 | 自动 | 无按钮 |
| 决策规则管理 | CRUD管理决策规则，配置触发条件和动作 | P1 | 人工 | 「决策规则」按钮 |
| 基础动作库查询 | 根据设备型号查询兼容的基础动作 | P1 | 自动 | 无按钮 |
| 动作序列调试 | 开发环境下的动作序列调试与预览 | P2 | 人工 | 「调试」按钮 |
| 自定义动作序列 | 用户创建自定义动作序列并保存 | P2 | 人工 | 「动作编排」按钮 |
| 动作回放 | 重放历史动作序列用于分析 | P2 | 人工 | 「回放」按钮 |

---

## 3. 数据模型

### 3.1 动作库 (action_library)

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| id | uint | 主键 | PK, auto_increment |
| action_id | string | 动作唯一标识 | unique, not null |
| action_name | string | 动作名称 | not null |
| action_name_en | string | 英文名称 | not null |
| category | string | 动作类别 | move/speak/look/listen/emotion |
| description | string | 动作描述 | nullable |
| duration_ms | int | 动作持续时间(毫秒) | not null |
| priority | int | 优先级 | 1-10, 10最高 |
| is_emergency | bool | 是否紧急动作 | default false |
| compatible_models | json | 兼容硬件型号列表 | not null |
| parameters | json | 动作参数定义 | nullable |
| animation_data | json | 动画数据（M5Stack屏幕） | nullable |
| motor_commands | json | 电机控制指令 | nullable |
| audio_file | string | 关联音频文件路径 | nullable |
| created_at | datetime | 创建时间 | auto |
| updated_at | datetime | 更新时间 | auto |

### 3.2 动作序列表 (action_sequences)

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| id | uint | 主键 | PK, auto_increment |
| sequence_id | string | 序列唯一标识 | unique, not null |
| sequence_name | string | 序列名称 | not null |
| trigger_type | string | 触发类型 | intent/emotion/sensor/manual |
| trigger_condition | json | 触发条件 | nullable |
| actions | json | 动作列表及参数 | not null |
| total_duration_ms | int | 总持续时间 | auto |
| is_loop | bool | 是否循环 | default false |
| loop_count | int | 循环次数 | default 1 |
| device_id | string | 关联设备（null表示通用） | FK to devices, nullable |
| created_by | uint | 创建用户 | FK to sys_users, nullable |
| is_public | bool | 是否公开 | default false |
| usage_count | int | 使用次数统计 | default 0 |
| created_at | datetime | 创建时间 | auto |
| updated_at | datetime | 更新时间 | auto |

### 3.3 决策规则表 (decision_rules)

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| id | uint | 主键 | PK, auto_increment |
| rule_id | string | 规则唯一标识 | unique, not null |
| rule_name | string | 规则名称 | not null |
| rule_type | string | 规则类型 | priority/conflict/filter/override |
| conditions | json | 条件表达式 | not null |
| condition_logic | string | 条件逻辑 | AND/OR |
| actions | json | 匹配后执行的动作 | not null |
| priority | int | 规则优先级 | 1-100 |
| is_active | bool | 是否启用 | default true |
| device_id | string | 关联设备（null表示通用） | FK to devices, nullable |
| description | string | 规则描述 | nullable |
| created_at | datetime | 创建时间 | auto |
| updated_at | datetime | 更新时间 | auto |

### 3.4 传感器事件表 (sensor_events)

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| id | uint | 主键 | PK, auto_increment |
| event_id | string | 事件唯一标识 | unique, not null |
| device_id | string | 设备ID | FK to devices, not null |
| sensor_type | string | 传感器类型 | fall/collision/touch/temperature/ultrasonic |
| sensor_data | json | 传感器原始数据 | not null |
| event_level | int | 事件级别 | 1=Info 2=Warning 3=Critical |
| processed | bool | 是否已处理 | default false |
| handled_by_rule | string | 处理的规则ID | nullable |
| created_at | datetime | 创建时间 | auto |

### 3.5 动作执行记录表 (action_executions)

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| id | uint | 主键 | PK, auto_increment |
| execution_id | string | 执行唯一标识 | unique, not null |
| sequence_id | string | 关联动作序列ID | FK to action_sequences |
| device_id | string | 设备ID | FK to devices |
| status | string | 执行状态 | pending/running/completed/failed/cancelled |
| current_action_index | int | 当前动作索引 | default 0 |
| started_at | datetime | 开始时间 | auto |
| completed_at | datetime | 完成时间 | nullable |
| error_message | string | 错误信息 | nullable |
| created_at | datetime | 创建时间 | auto |

---

## 4. 接口定义

### 4.1 生成动作序列

```
POST /api/v1/behavior/generate
```

**请求参数：**

| 参数 | 类型 | 位置 | 必填 | 说明 |
|------|------|------|------|------|
| device_id | string | body | 是 | 设备ID |
| intent | string | body | 是 | AI识别的意图 |
| context | json | body | 否 | 上下文信息 |
| emotion | int | body | 否 | 当前情绪值 0-100 |
| user_id | uint | body | 否 | 用户ID |

**请求示例：**
```json
{
  "device_id": "pet-001",
  "intent": "greeting",
  "context": { "time_of_day": "morning", "owner_mood": "happy" },
  "emotion": 75,
  "user_id": 10001
}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "sequence_id": "seq-uuid-001",
    "execution_id": "exec-uuid-001",
    "actions": [
      { "action_id": "look_at_owner", "duration_ms": 500, "params": {} },
      { "action_id": "tail_wag", "duration_ms": 800, "params": { "speed": 3 } },
      { "action_id": "happy_bark", "duration_ms": 1000, "params": { "volume": 70 } },
      { "action_id": "walk_forward", "duration_ms": 1500, "params": { "distance": 30 } }
    ],
    "total_duration_ms": 3800,
    "is_loop": false
  }
}
```

### 4.2 传感器事件上报

```
POST /api/v1/behavior/sensor-event
```

**请求参数：**

| 参数 | 类型 | 位置 | 必填 | 说明 |
|------|------|------|------|------|
| device_id | string | body | 是 | 设备ID |
| sensor_type | string | body | 是 | 传感器类型 |
| sensor_data | json | body | 是 | 传感器数据 |
| event_level | int | body | 是 | 事件级别 1/2/3 |

**请求示例：**
```json
{
  "device_id": "pet-001",
  "sensor_type": "fall",
  "sensor_data": { "accelerometer": { "x": -0.2, "y": -0.8, "z": -9.6 }, "timestamp": 1710902400000 },
  "event_level": 3
}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": { 
    "event_id": "evt-001", 
    "processed": true, 
    "handled_action": "stop_all", 
    "message": "紧急停止所有动作",
    "mqtt_published": true
  }
}
```

### 4.3 查询动作库

```
GET /api/v1/behavior/actions
```

**请求参数：**

| 参数 | 类型 | 位置 | 必填 | 说明 |
|------|------|------|------|------|
| category | string | query | 否 | 动作类别筛选 |
| device_model | string | query | 否 | 设备型号筛选 |
| page | int | query | 否 | 页码，默认1 |
| page_size | int | query | 否 | 每页条数，默认50 |

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      { 
        "action_id": "walk_forward", 
        "action_name": "前进", 
        "category": "move", 
        "duration_ms": 1500, 
        "priority": 5, 
        "is_emergency": false,
        "parameters": { "distance": {"type": "int", "default": 30}, "speed": {"type": "int", "default": 3} }
      },
      { 
        "action_id": "stop_immediately", 
        "action_name": "立即停止", 
        "category": "move", 
        "duration_ms": 100, 
        "priority": 10, 
        "is_emergency": true 
      }
    ],
    "pagination": { "page": 1, "page_size": 50, "total": 120, "total_pages": 3 }
  }
}
```

### 4.4 创建自定义动作序列

```
POST /api/v1/behavior/sequences
```

**请求参数：**

| 参数 | 类型 | 位置 | 必填 | 说明 |
|------|------|------|------|------|
| sequence_name | string | body | 是 | 序列名称 |
| trigger_type | string | body | 是 | 触发类型 |
| trigger_condition | json | body | 否 | 触发条件 |
| actions | array | body | 是 | 动作列表 |
| is_loop | bool | body | 否 | 是否循环 |
| loop_count | int | body | 否 | 循环次数 |
| device_id | string | body | 否 | 关联设备 |

**请求示例：**
```json
{
  "sequence_name": "开心迎接",
  "trigger_type": "intent",
  "trigger_condition": { "intent": "owner_home" },
  "actions": [
    { "action_id": "ear_up", "duration_ms": 300, "params": {} },
    { "action_id": "tail_wag", "duration_ms": 1000, "params": { "speed": 5 } },
    { "action_id": "happy_sound", "duration_ms": 800, "params": { "volume": 80 } }
  ],
  "is_loop": false
}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "sequence_id": "seq-uuid-002",
    "sequence_name": "开心迎接",
    "total_duration_ms": 2100,
    "created_at": "2026-03-22T10:00:00Z"
  }
}
```

### 4.5 查询动作序列列表

```
GET /api/v1/behavior/sequences
```

**请求参数：**

| 参数 | 类型 | 位置 | 必填 | 说明 |
|------|------|------|------|------|
| trigger_type | string | query | 否 | 触发类型筛选 |
| device_id | string | query | 否 | 设备ID筛选 |
| page | int | query | 否 | 页码 |

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {
        "sequence_id": "seq-uuid-001",
        "sequence_name": "开心迎接",
        "trigger_type": "intent",
        "actions_count": 3,
        "total_duration_ms": 2100,
        "usage_count": 45,
        "created_at": "2026-03-20T10:00:00Z"
      }
    ],
    "pagination": { "page": 1, "page_size": 20, "total": 15, "total_pages": 1 }
  }
}
```

### 4.6 决策规则管理 - 获取规则列表

```
GET /api/v1/behavior/rules
```

**请求参数：**

| 参数 | 类型 | 位置 | 必填 | 说明 |
|------|------|------|------|------|
| rule_type | string | query | 否 | 规则类型筛选 |
| is_active | bool | query | 否 | 是否启用 |
| device_id | string | query | 否 | 设备ID筛选 |

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {
        "rule_id": "rule-001",
        "rule_name": "跌落紧急停止",
        "rule_type": "override",
        "conditions": { "sensor_type": "fall", "event_level": 3 },
        "condition_logic": "AND",
        "actions": [{ "action_id": "stop_immediately", "params": {} }],
        "priority": 100,
        "is_active": true,
        "description": "检测到跌落时立即停止所有动作",
        "created_at": "2026-03-20T10:00:00Z"
      }
    ]
  }
}
```

### 4.7 决策规则管理 - 创建规则

```
POST /api/v1/behavior/rules
```

**请求示例：**
```json
{
  "rule_name": "碰撞后退",
  "rule_type": "filter",
  "conditions": { "sensor_type": "collision", "event_level": 2 },
  "condition_logic": "AND",
  "actions": [{ "action_id": "walk_backward", "params": { "distance": 20, "speed": 2 } }],
  "priority": 80,
  "is_active": true,
  "description": "检测到碰撞时后退"
}
```

### 4.8 决策规则管理 - 更新规则

```
PUT /api/v1/behavior/rules/{rule_id}
```

### 4.9 决策规则管理 - 删除规则

```
DELETE /api/v1/behavior/rules/{rule_id}
```

### 4.10 动作执行状态查询

```
GET /api/v1/behavior/executions/{execution_id}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "execution_id": "exec-001",
    "sequence_id": "seq-001",
    "status": "running",
    "current_action_index": 2,
    "current_action": "tail_wag",
    "elapsed_ms": 1300,
    "remaining_ms": 2500
  }
}
```

### 4.11 取消动作执行

```
DELETE /api/v1/behavior/executions/{execution_id}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": { "execution_id": "exec-001", "status": "cancelled", "message": "动作序列已取消" }
}
```

### 4.12 获取传感器事件历史

```
GET /api/v1/behavior/sensor-events
```

**请求参数：**

| 参数 | 类型 | 位置 | 必填 | 说明 |
|------|------|------|------|------|
| device_id | string | query | 是 | 设备ID |
| sensor_type | string | query | 否 | 传感器类型 |
| event_level | int | query | 否 | 事件级别 |
| start_date | string | query | 否 | 开始日期 |
| end_date | string | query | 否 | 结束日期 |
| page | int | query | 否 | 页码 |

### 4.13 错误码定义

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 20001 | 设备不在线 |
| 20002 | 动作库中找不到指定动作 |
| 20003 | 动作序列执行超时 |
| 20004 | 动作执行失败 |
| 20005 | 决策规则不存在 |
| 20006 | 传感器事件处理异常 |

---

## 5. 流程图

### 5.1 动作序列生成与执行流程

```
用户意图/情绪输入 (OpenClaw)
       │
       ▼
┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐
│   主人画像       │  │   知识库        │  │   记忆库        │
│  查询偏好        │  │  上下文信息     │  │  历史行为        │
└─────────────────┘  └─────────────────┘  └─────────────────┘
       │
       └─────────────────┼─────────────────┘
                        ▼
            ┌───────────────────────┐
            │     决策树引擎        │
            │  匹配触发条件和规则   │
            └───────────┬───────────┘
                        │
       ┌────────────────┼────────────────┐
       ▼                │                ▼
┌───────────────┐       │        ┌───────────────┐
│ 动作序列生成  │       │        │   规则冲突     │
│ 规划动作步骤   │       │        │   优先级处理   │
└───────┬───────┘       │        └───────┬───────┘
        │               │                │
        └───────────────┼────────────────┘
                        ▼
            ┌───────────────────────┐
            │    动作优先级管理器    │
            │  紧急动作插入/排队     │
            └───────────┬───────────┘
                        │
    ┌───────────────────┼───────────────────┐
    ▼                   ▼                   ▼
┌────────────┐   ┌────────────┐   ┌────────────┐
│ 电机控制   │   │ 音频播放   │   │ 屏幕动画   │
│ 指令序列   │   │ 指令序列   │   │ 指令序列   │
└─────┬──────┘   └─────┬──────┘   └─────┬──────┘
      │                 │                 │
      └─────────────────┼─────────────────┘
                        ▼
            ┌───────────────────────┐
            │    MQTT下发          │
            │  /miniclaw/{device_id}
            │      /down/action    │
            └───────────┬───────────┘
                        │
                        ▼
            ┌───────────────────────┐
            │      MiniClaw         │
            │       设备            │
            └───────────┬───────────┘
                        │
                        ▼
            ┌───────────────────────┐
            │     传感器数据         │
            │  处理 (避障/防跌落)   │
            └───────────┬───────────┘
                        │
       ┌────────────────┼────────────────┐
       ▼                ▼                ▼
┌────────────┐   ┌────────────┐   ┌────────────┐
│  跌落检测   │   │  碰撞检测   │   │  触摸事件  │
│  → 紧急停止 │   │  → 后退    │   │  → 反馈    │
└────────────┘   └────────────┘   └────────────┘
```

### 5.2 传感器事件紧急处理流程

```
传感器数据上报 -> 行为引擎接收 -> 决策规则匹配
       │
       ├─ Critical级别 -> 立即执行紧急动作 -> MQTT下发stop_immediate
       │
       ├─ Warning级别 -> 队列优先处理 -> 动作替换/插入
       │
       └─ Info级别 -> 正常队列处理 -> 继续原计划
```

---

## 6. 模块联动

### 6.1 与OpenClaw控制台(OPENCLAW_CONSOLE)联动

- **触发时机：** 用户发送消息或点击快捷指令
- **联动内容：** 控制台发送意图(intent)到行为引擎，行为引擎返回动作序列
- **数据流向：** 控制台 → 行为引擎 → MQTT → MiniClaw

### 6.2 与主人画像库(OWNER_PROFILE)联动

- **触发时机：** 动作序列生成时
- **联动内容：** 行为引擎查询主人偏好（如动作风格：活泼/安静），根据偏好调整动作速度和幅度
- **数据流向：** 主人画像库 → 行为引擎

### 6.3 与知识库(KNOWLEDGE_BASE)联动

- **触发时机：** 需要环境上下文时
- **联动内容：** 行为引擎查询天气/时间等信息，生成与环境匹配的动作
- **数据流向：** 知识库 → 行为引擎

### 6.4 与宠物记忆库(PET_MEMORY)联动

- **触发时机：** 动作执行完成后
- **联动内容：** 行为引擎记录动作执行结果到记忆库，记忆库提供历史行为模式供决策参考
- **数据流向：** 行为引擎 ↔ 记忆库

### 6.5 与MiniClaw通信协议(MINICLAW_PROTOCOL)联动

- **触发时机：** 动作下发和状态上报
- **联动内容：** 行为引擎通过MQTT下发动作指令，接收MiniClaw传感器数据进行处理
- **数据流向：** 行为引擎 ↔ MQTT Broker ↔ MiniClaw

---

## 7. 验收标准

### 7.1 功能验收

| 功能 | 验收条件 |
|------|----------|
| 动作序列生成 | 给定意图，3秒内返回有效动作序列 |
| 传感器处理 | Critical事件100ms内响应，Warning事件500ms内响应 |
| 决策树执行 | 规则匹配准确率>=95% |
| 优先级管理 | 紧急动作可立即中断正在执行的普通动作 |
| 自定义序列 | 用户可创建包含最多20个动作的序列 |
| 规则管理 | 支持CRUD决策规则，支持启用/禁用 |

### 7.2 性能验收

- 动作序列生成延迟 <= 3秒
- 紧急停止响应延迟 <= 200ms
- 支持同时执行10个设备的动作序列
- 单设备动作队列最大长度：50个

### 7.3 兼容性验收

- 支持M5Stack Basic/Gray/Stack三种硬件型号
- 支持固件版本 >= V1.2.0

---

## 8. 动作库标准动作

### 8.1 移动类动作

| action_id | 动作名称 | 英文名称 | 默认时长 | 优先级 |
|-----------|----------|----------|----------|--------|
| walk_forward | 前进 | Walk Forward | 1500ms | 5 |
| walk_backward | 后退 | Walk Backward | 1500ms | 5 |
| turn_left | 左转 | Turn Left | 1000ms | 5 |
| turn_right | 右转 | Turn Right | 1000ms | 5 |
| stop_immediately | 立即停止 | Emergency Stop | 100ms | 10 |
| dance | 跳舞 | Dance | 3000ms | 6 |

### 8.2 情感类动作

| action_id | 动作名称 | 默认时长 | 优先级 |
|-----------|----------|----------|--------|
| tail_wag | 摇尾巴 | 800ms | 4 |
| ear_up | 耳朵竖起 | 300ms | 4 |
| happy_bark | 开心叫 | 1000ms | 4 |
| sad_whimper | 悲伤呜咽 | 1200ms | 4 |
| excited_jump | 兴奋跳跃 | 1500ms | 6 |

### 8.3 声音类动作

| action_id | 动作名称 | 默认时长 | 优先级 |
|-----------|----------|----------|--------|
| speak_greeting | 问候 | 1500ms | 3 |
| speak_bye | 告别 | 1200ms | 3 |
| speak_hungry | 叫声-饿 | 1000ms | 3 |
| speak_play | 叫声-想玩 | 1000ms | 3 |

---

## 9. 前端页面设计

### 9.1 动作调试界面

```
┌─────────────────────────────────────────────────────────────┐
│  动作调试工具                                    [返回] [保存]│
├─────────────────────────────────────────────────────────────┤
│ 设备: [pet-001 v]  │  动作序列: [开心迎接 v]  │  [添加动作]  │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌─────────────────────────┐  ┌─────────────────────────┐ │
│  │    宠物动画预览         │  │    动作列表             │ │
│  │                         │  │                         │ │
│  │    [🐱 动画区域]        │  │  1. look_at_owner  0.5s │ │
│  │                         │  │  2. tail_wag       0.8s │ │
│  │    表情: 😊            │  │  3. happy_bark     1.0s │ │
│  │    心情: 75             │  │  4. walk_forward   1.5s │ │
│  │                         │  │                         │ │
│  └─────────────────────────┘  │  [上移] [下移] [删除]   │ │
│                                └─────────────────────────┘ │
├─────────────────────────────────────────────────────────────┤
│  总时长: 3.8s  │  优先级: 5  │  [▶ 预览] [⏹ 停止]        │
└─────────────────────────────────────────────────────────────┘
```

### 9.2 决策规则管理界面

```
┌─────────────────────────────────────────────────────────────┐
│  决策规则管理                                               │
├─────────────────────────────────────────────────────────────┤
│ [新建规则]                                                  │
├─────────────────────────────────────────────────────────────┤
│ 规则名称: 跌落紧急停止                                      │
│ 规则类型: [override v]                                      │
│ 优先级: [100 v]                                             │
│ 启用: [✓]                                                   │
│                                                            │
│ 触发条件:                                                   │
│   传感器类型: [fall v]  事件级别: [3 v]                     │
│                                                            │
│ 执行动作:                                                   │
│   [stop_immediately]  参数: {}                            │
│                                                            │
│ 说明: 检测到跌落时立即停止所有动作                           │
│                                                            │
│                              [取消]  [保存]                │
└─────────────────────────────────────────────────────────────┘
```
