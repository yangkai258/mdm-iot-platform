# 宠物管理 API

**控制器：** `controllers/pet_controller.go`  
**路由前缀：** `/api/v1/pets`

---

## 概述

宠物管理 API 用于管理宠物的状态、对话、动作和设置。每个设备关联一个虚拟宠物实体。

### 宠物属性说明

| 属性 | 类型 | 说明 |
|------|------|------|
| device_id | string | 关联的设备ID |
| pet_name | string | 宠物名称，默认"小爪" |
| pet_type | string | 宠物类型：`cat` / `dog` / `bird` 等 |
| mood | int | 心情值 0-100 |
| energy | int | 能量值 0-100 |
| hunger | int | 饥饿值 0-100 |
| current_expression | string | 当前表情：`happy` / `sad` / `angry` / `sleeping` 等 |
| is_online | bool | 关联设备是否在线 |
| personality | object | 个性配置 JSON |
| appearance | object | 外观配置 JSON |

---

## 1. 获取宠物状态

获取指定设备宠物的当前状态。

### 请求

```
GET /api/v1/pets/:device_id/status
```

### 路径参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| device_id | string | 是 | 设备ID |

### 响应示例

```json
{
  "code": 0,
  "data": {
    "device_id": "dev_abc123",
    "pet_name": "小爪",
    "pet_type": "cat",
    "mood": 75,
    "energy": 80,
    "hunger": 20,
    "current_expression": "happy",
    "is_online": true,
    "last_seen_at": "2026-03-22T10:00:00Z",
    "personality": {
      "trait": "curious",
      "favorite_food": "fish"
    },
    "appearance": {
      "color": "orange",
      "pattern": "solid"
    }
  }
}
```

### 响应字段说明

| 字段 | 类型 | 说明 |
|------|------|------|
| mood | int | 心情值 0-100，值越高心情越好 |
| energy | int | 能量值 0-100，值越高越精神 |
| hunger | int | 饥饿值 0-100，值越高越饥饿 |
| current_expression | string | 当前表情状态 |
| is_online | bool | 设备是否在线（通过 Redis 设备影子判断） |
| last_seen_at | datetime | 最后一次心跳时间 |

---

## 2. 发送消息

向宠物发送消息，支持文本、表情、图片等。

### 请求

```
POST /api/v1/pets/:device_id/messages
```

### 路径参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| device_id | string | 是 | 设备ID |

### 请求参数 (JSON)

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| content | string | 是 | 消息内容 |
| content_type | int | 否 | 内容类型：1=文本 2=表情 3=图片 4=语音，默认 1 |
| media_url | string | 否 | 媒体资源URL（当 content_type 非文本时） |
| metadata | object | 否 | 元数据 |

### 请求示例

```json
{
  "content": "你好，小爪！",
  "content_type": 1,
  "metadata": {
    "client": "web"
  }
}
```

### 响应示例

```json
{
  "code": 0,
  "data": {
    "message": {
      "message_id": "msg_xyz789",
      "content": "你好，小爪！",
      "content_type": 1,
      "sender_type": "user",
      "created_at": "2026-03-22T10:00:00Z"
    },
    "conversation_id": "conv_abc123"
  }
}
```

### 业务逻辑

1. 如果用户与该设备不存在对话会话，自动创建新会话
2. 消息保存到数据库 `messages` 表
3. 通过 MQTT 将消息下发到设备 `/device/{device_id}/down/action`
4. MQTT Payload 示例：

```json
{
  "message_id": "msg_xyz789",
  "content": "你好，小爪！",
  "content_type": 1,
  "timestamp": "2026-03-22T10:00:00Z"
}
```

---

## 3. 下发快捷指令

向设备下发预设的快捷动作指令。

### 请求

```
POST /api/v1/pets/:device_id/actions
```

### 路径参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| device_id | string | 是 | 设备ID |

### 请求参数 (JSON)

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| action_id | string | 是 | 动作ID（来自动作库） |
| parameters | object | 否 | 动作参数 |

### 请求示例

```json
{
  "action_id": "action_wiggle",
  "parameters": {
    "duration_ms": 500,
    "intensity": 3
  }
}
```

### 响应示例

```json
{
  "code": 0,
  "data": {
    "action_id": "action_wiggle",
    "status": "sent",
    "message": "动作已下发"
  }
}
```

### 错误码

| code | 说明 |
|------|------|
| 404 | 动作不存在 |
| 503 | MQTT 服务不可用 |

---

## 4. 更新宠物设置

更新宠物的名称、类型、个性等设置。

### 请求

```
PUT /api/v1/pets/:device_id/settings
```

### 路径参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| device_id | string | 是 | 设备ID |

### 请求参数 (JSON)

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| pet_name | string | 否 | 宠物名称 |
| pet_type | string | 否 | 宠物类型 |
| personality | object | 否 | 个性配置 |
| appearance | object | 否 | 外观配置 |

### 请求示例

```json
{
  "pet_name": "小橘",
  "pet_type": "cat",
  "personality": {
    "trait": "playful",
    "favorite_food": "fish",
    "favorite_toy": "feather"
  },
  "appearance": {
    "color": "orange",
    "pattern": "solid",
    "has_collar": true
  }
}
```

### 响应示例

```json
{
  "code": 0,
  "data": {
    "device_id": "dev_abc123",
    "pet_name": "小橘",
    "pet_type": "cat",
    "personality": {
      "trait": "playful",
      "favorite_food": "fish",
      "favorite_toy": "feather"
    },
    "appearance": {
      "color": "orange",
      "pattern": "solid",
      "has_collar": true
    }
  }
}
```

---

## 5. 心情激励

给宠物发送心情激励，提升宠物心情值。

### 请求

```
POST /api/v1/pets/:device_id/boost
```

### 路径参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| device_id | string | 是 | 设备ID |

### 请求参数 (JSON)

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| boost_type | string | 是 | 激励类型：`food`（投喂）/ `play`（玩耍）/ `pet`（抚摸）/ `talk`（聊天） |
| amount | int | 否 | 激励量，默认 10 |

### 请求示例

```json
{
  "boost_type": "play",
  "amount": 15
}
```

### 响应示例

```json
{
  "code": 0,
  "data": {
    "device_id": "dev_abc123",
    "boost_type": "play",
    "mood_before": 60,
    "mood_after": 75,
    "mood_change": 15,
    "energy_before": 80,
    "energy_after": 65,
    "energy_change": -15
  }
}
```

### 激励效果说明

| 激励类型 | 心情 | 能量 | 饥饿 |
|----------|------|------|------|
| food | +5~10 | 0 | -20~30 |
| play | +10~20 | -10~20 | +10~15 |
| pet | +5~15 | 0 | 0 |
| talk | +3~8 | 0 | 0 |

---

## 6. 对话列表

获取用户的宠物对话会话列表。

### 请求

```
GET /api/v1/conversations
```

### 查询参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| page | int | 否 | 页码，默认 1 |
| page_size | int | 否 | 每页条数，默认 20 |
| device_id | string | 否 | 设备ID筛选 |

### 响应示例

```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "conversation_id": "conv_abc123",
        "device_id": "dev_abc123",
        "pet_name": "小爪",
        "title": "与小爪的对话",
        "last_message": "今天心情不错！",
        "last_message_at": "2026-03-22T10:00:00Z",
        "message_count": 42,
        "created_at": "2026-03-01T10:00:00Z"
      }
    ],
    "total": 5,
    "page": 1,
    "page_size": 20
  }
}
```

---

## 7. 消息历史

获取指定对话的消息历史。

### 请求

```
GET /api/v1/conversations/:id/messages
```

### 路径参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | string | 是 | 对话ID |

### 查询参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| page | int | 否 | 页码，默认 1 |
| page_size | int | 否 | 每页条数，默认 50 |
| before | string | 否 | 获取此时间之前的消息（用于滚动加载） |

### 响应示例

```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "message_id": "msg_001",
        "conversation_id": "conv_abc123",
        "sender_type": "user",
        "sender_id": "user_123",
        "content": "你好，小爪！",
        "content_type": 1,
        "created_at": "2026-03-22T09:00:00Z"
      },
      {
        "message_id": "msg_002",
        "conversation_id": "conv_abc123",
        "sender_type": "pet",
        "content": "你好呀！今天心情不错！",
        "content_type": 1,
        "created_at": "2026-03-22T09:00:05Z"
      }
    ],
    "total": 100,
    "page": 1,
    "page_size": 50,
    "has_more": true
  }
}
```

### 消息发送者类型

| 类型 | 说明 |
|------|------|
| user | 用户发送 |
| pet | 宠物回复 |
| system | 系统消息 |

---

## 8. 获取宠物档案

获取宠物的完整档案信息。

### 请求

```
GET /api/v1/pets/:device_id/profile
```

### 响应示例

```json
{
  "code": 0,
  "data": {
    "device_id": "dev_abc123",
    "basic_info": {
      "pet_name": "小爪",
      "pet_type": "cat",
      "breed": "橘猫",
      "gender": "male",
      "birthday": "2024-06-01",
      "age_months": 21
    },
    "status": {
      "mood": 75,
      "energy": 80,
      "hunger": 20,
      "health": 95,
      "current_expression": "happy"
    },
    "appearance": {
      "color": "orange",
      "pattern": "solid",
      "has_collar": true,
      "collar_color": "red"
    },
    "personality": {
      "trait": "curious",
      "favorite_food": "fish",
      "favorite_toy": "feather",
      "sleep_hours": 14
    },
    "statistics": {
      "total_conversations": 150,
      "total_messages": 1200,
      "total_actions": 80,
      "online_days": 45
    },
    "created_at": "2024-06-01T10:00:00Z",
    "updated_at": "2026-03-22T10:00:00Z"
  }
}
```

---

## 9. 更新宠物档案

更新宠物档案的详细信息。

### 请求

```
PUT /api/v1/pets/:device_id/profile
```

### 请求参数 (JSON)

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| pet_name | string | 否 | 宠物名称 |
| pet_type | string | 否 | 宠物类型 |
| breed | string | 否 | 品种 |
| gender | string | 否 | 性别 |
| birthday | string | 否 | 生日 `YYYY-MM-DD` |
| appearance | object | 否 | 外观配置 |
| personality | object | 否 | 个性配置 |

### 请求示例

```json
{
  "pet_name": "小橘",
  "breed": "英短",
  "gender": "male",
  "birthday": "2023-01-15",
  "appearance": {
    "color": "gray",
    "pattern": "tabby",
    "has_collar": true,
    "collar_color": "blue"
  }
}
```

---

## MQTT 协议说明

### 宠物消息下发

**Topic:** `/device/{device_id}/down/action`

**Payload:**

```json
{
  "type": "message",
  "message_id": "msg_xyz789",
  "content": "你好，小爪！",
  "content_type": 1,
  "timestamp": "2026-03-22T10:00:00Z"
}
```

### 宠物动作下发

**Topic:** `/device/{device_id}/down/action`

**Payload:**

```json
{
  "type": "action",
  "action_id": "action_wiggle",
  "action_name": "摇尾巴",
  "duration_ms": 500,
  "priority": 1,
  "parameters": {
    "intensity": 3
  },
  "motor_commands": [
    {"motor": "tail", "angle": 30, "duration_ms": 500}
  ],
  "timestamp": "2026-03-22T10:00:00Z"
}
```

### 设备上行消息

**Topic:** `/device/{device_id}/up/message`

**Payload:**

```json
{
  "message_id": "msg_xyz789",
  "content": "你好呀！",
  "content_type": 1,
  "pet_expression": "happy",
  "timestamp": "2026-03-22T10:00:05Z"
}
```

---

## 相关数据模型

### PetStatusV2

| 字段 | 类型 | 说明 |
|------|------|------|
| device_id | string | 设备ID（主键） |
| pet_name | string | 宠物名称 |
| pet_type | string | 宠物类型 |
| mood | int | 心情值 0-100 |
| energy | int | 能量值 0-100 |
| hunger | int | 饥饿值 0-100 |
| health | int | 健康值 0-100 |
| current_expression | string | 当前表情 |
| personality | text | 个性JSON |
| appearance | text | 外观JSON |
| is_online | bool | 设备在线状态 |

### Conversation

| 字段 | 类型 | 说明 |
|------|------|------|
| conversation_id | string | 对话ID（主键） |
| user_id | int | 用户ID |
| device_id | string | 设备ID |
| title | string | 对话标题 |
| last_message | text | 最后一条消息 |
| last_message_at | datetime | 最后消息时间 |
| message_count | int | 消息数量 |

### Message

| 字段 | 类型 | 说明 |
|------|------|------|
| message_id | string | 消息ID（主键） |
| conversation_id | string | 关联对话ID |
| sender_type | string | 发送者类型：user/pet/system |
| sender_id | int | 发送者ID |
| content | text | 消息内容 |
| content_type | int | 内容类型：1=文本 2=表情 3=图片 4=语音 |
| media_url | string | 媒体URL |
