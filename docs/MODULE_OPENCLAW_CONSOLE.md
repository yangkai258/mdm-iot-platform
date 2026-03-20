# OpenClaw宠物控制台

**版本：** V1.0  
**模块负责人：** agentcp  
**编制日期：** 2026-03-20  

---

## 1. 概述

OpenClaw宠物控制台是AI电子宠物的Web管理界面，为用户提供与宠物对话、状态监控、快捷指令和个性化设置的能力。控制台作为OpenClaw AI层的前端入口，连接设备管理层与AI大脑，是用户感知宠物智能的核心触点。

**业务目标：**
- 提供自然语言对话界面，支持文字和语音输入
- 实时展示宠物状态（位置/表情/情绪/电量）
- 支持快捷指令快速触发预设动作
- 支持宠物个性化设置（名字/性格/外观）

---

## 2. 功能列表

| 功能 | 描述 | 优先级 | 触发方式 | 前端入口/按钮 |
|------|------|--------|----------|--------------|
| AI对话界面 | 聊天窗口，支持文字输入和语音录制 | P0 | 自动 | 主界面左侧聊天区 |
| 宠物状态展示 | 实时显示宠物位置/表情/情绪/电量/在线状态 | P0 | 自动 | 主界面右侧状态卡片 |
| 快捷指令面板 | 预设动作按钮（过来/坐下/转圈/睡觉等） | P0 | 人工 | 主界面底部快捷指令栏 |
| 宠物设置 | 修改宠物名字/性格/外观/声音 | P1 | 人工 | 右上角「设置」图标 |
| 历史对话查询 | 分页查看历史对话记录 | P1 | 人工 | 左侧边栏「历史记录」按钮 |
| 表情切换 | 手动切换宠物显示表情 | P2 | 人工 | 状态卡片内「表情」按钮 |
| 心情激励 | 手动触发心情恢复（饥饿/疲劳时） | P2 | 人工 | 状态卡片内「投喂/休息」按钮 |

---

## 3. 数据模型

### 3.1 对话记录表 (conversations)

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| id | uint | 主键 | PK, auto_increment |
| conversation_id | string | 对话会话唯一标识 | unique, not null, UUID |
| user_id | uint | 所属用户 | FK → sys_users, not null |
| device_id | string | 关联设备 | FK → devices, nullable |
| title | string | 对话标题（自动生成或用户自定义） | not null |
| last_message | text | 最后一条消息摘要 | nullable |
| last_message_at | datetime | 最后消息时间 | auto |
| message_count | int | 消息总数 | default 0 |
| status | int | 状态 | 1=活跃 2=已归档 |
| created_at | datetime | 创建时间 | auto |
| updated_at | datetime | 更新时间 | auto |

### 3.2 对话消息表 (messages)

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| id | uint | 主键 | PK, auto_increment |
| message_id | string | 消息唯一标识 | unique, not null, UUID |
| conversation_id | string | 所属会话 | FK → conversations, not null |
| sender_type | int | 发送者类型 | 1=用户 2=宠物(AI) 3=系统 |
| sender_id | uint | 发送者ID | nullable |
| content | text | 消息内容 | not null |
| content_type | int | 内容类型 | 1=文本 2=语音 3=图片 4=指令 |
| media_url | string | 媒体文件URL | nullable |
| intent | string | 识别意图 | nullable |
| confidence | float | 意图置信度 | nullable |
| metadata | json | 扩展字段（情绪/动作等） | nullable |
| created_at | datetime | 创建时间 | auto |

### 3.3 宠物状态表 (pet_status)

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| id | uint | 主键 | PK, auto_increment |
| device_id | string | 设备ID | FK → devices, unique, not null |
| pet_name | string | 宠物名字 | not null, default "小爪" |
| pet_type | string | 宠物类型 | default "cat" |
| personality | string | 性格配置 | JSON |
| appearance | json | 外观配置（颜色/配件） | JSON |
| mood | int | 当前心情值 | 0-100, default 50 |
| energy | int | 精力值 | 0-100, default 100 |
| hunger | int | 饥饿值 | 0-100, default 0 |
| position_x | float | X坐标 | default 0 |
| position_y | float | Y坐标 | default 0 |
| current_expression | string | 当前表情 | default "happy" |
| current_action | string | 当前动作 | nullable |
| is_online | bool | 在线状态 | default false |
| last_seen_at | datetime | 最后在线时间 | auto |
| updated_at | datetime | 更新时间 | auto |

---

## 4. 接口定义

### 4.1 获取宠物当前状态

```
GET /api/v1/pets/{device_id}/status
```

**参数：**

| 参数 | 类型 | 位置 | 必填 | 说明 |
|------|------|------|------|------|
| device_id | string | path | 是 | 设备ID |

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "device_id": "pet-001",
    "pet_name": "小爪",
    "pet_type": "cat",
    "mood": 75,
    "energy": 80,
    "hunger": 20,
    "position": { "x": 120.5, "y": 80.3 },
    "expression": "happy",
    "action": "idle",
    "is_online": true,
    "last_seen_at": "2026-03-20T10:30:00Z"
  }
}
```

### 4.2 发送消息

```
POST /api/v1/pets/{device_id}/messages
```

**参数：**

| 参数 | 类型 | 位置 | 必填 | 说明 |
|------|------|------|------|------|
| device_id | string | path | 是 | 设备ID |
| content | string | body | 是 | 消息内容 |
| content_type | int | body | 否 | 1=文本 2=语音，默认1 |
| sender_type | int | body | 否 | 1=用户，默认1 |

**请求示例：**
```json
{
  "content": "今天天气怎么样？",
  "content_type": 1,
  "sender_type": 1
}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "message_id": "msg-uuid-001",
    "content": "今天天气晴朗，气温20-25度，适合出门散步哦！",
    "sender_type": 2,
    "created_at": "2026-03-20T10:31:00Z"
  }
}
```

### 4.3 获取对话历史

```
GET /api/v1/conversations
```

**参数：**

| 参数 | 类型 | 位置 | 必填 | 说明 |
|------|------|------|------|------|
| page | int | query | 否 | 页码，默认1 |
| page_size | int | query | 否 | 每页条数，默认20 |
| device_id | string | query | 否 | 设备ID筛选 |

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {
        "conversation_id": "conv-001",
        "title": "关于天气的话题",
        "last_message": "适合出门散步哦！",
        "last_message_at": "2026-03-20T10:31:00Z",
        "message_count": 12
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 45,
      "total_pages": 3
    }
  }
}
```

### 4.4 获取会话消息列表

```
GET /api/v1/conversations/{conversation_id}/messages
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {
        "message_id": "msg-001",
        "sender_type": 1,
        "content": "今天天气怎么样？",
        "created_at": "2026-03-20T10:30:00Z"
      },
      {
        "message_id": "msg-002",
        "sender_type": 2,
        "content": "今天天气晴朗，气温20-25度，适合出门散步哦！",
        "created_at": "2026-03-20T10:30:05Z"
      }
    ],
    "pagination": { "page": 1, "page_size": 50, "total": 12, "total_pages": 1 }
  }
}
```

### 4.5 更新宠物设置

```
PUT /api/v1/pets/{device_id}/settings
```

**请求示例：**
```json
{
  "pet_name": "小爪",
  "personality": {
    "playfulness": 80,
    "friendliness": 90,
    "curiosity": 70
  },
  "appearance": {
    "color": "#FFB6C1",
    "accessories": ["bowtie", "collar"]
  }
}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "device_id": "pet-001",
    "pet_name": "小爪",
    "updated_at": "2026-03-20T10:35:00Z"
  }
}
```

### 4.6 快捷指令下发

```
POST /api/v1/pets/{device_id}/actions
```

**请求示例：**
```json
{
  "action": "come_here",
  "params": {}
}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "action_id": "act-uuid-001",
    "action": "come_here",
    "status": "queued",
    "created_at": "2026-03-20T10:36:00Z"
  }
}
```

---

## 5. 流程图

### 5.1 用户对话完整流程

```
┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│   用户输入   │     │  控制台前端  │     │  OpenClaw    │
│  (文字/语音) │     │   (Vue.js)   │     │   AI大脑     │
└──────┬───────┘     └──────┬───────┘     └──────┬───────┘
       │                    │                    │
       │ 1.输入消息         │                    │
       │───────────────────>│                    │
       │                    │                    │
       │                    │ 2.发送消息请求      │
       │                    │───────────────────>│
       │                    │                    │
       │                    │                    │ 3.查询主人画像
       │                    │                    │ 4.查询知识库
       │                    │                    │ 5.生成回复+动作
       │                    │                    │
       │                    │ 6.返回AI回复        │
       │                    │<──────────────────│
       │                    │                    │
       │ 7.展示回复+动画     │                    │
       │<───────────────────│                    │
       │                    │                    │
       │ 8.MQTT下发动作指令  │                    │
       │                    │═══════════════════│ (同步)
       │                    │                    │
       │                    │     ┌──────────────┴───────┐
       │                    │     │      MiniClaw       │
       │                    │     │       设备          │
       │                    │     └──────────────────────┘
       │                    │                    │
       │                    │ 9.状态上报→记忆更新 │
       │                    │<───────────────────│
       │                    │                    │
       │ 10.更新界面状态     │                    │
       │<───────────────────│                    │
```

### 5.2 快捷指令流程

```
┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│  点击快捷    │     │  控制台前端  │     │   MiniClaw   │
│    按钮      │     │   (Vue.js)   │     │     设备      │
└──────┬───────┘     └──────┬───────┘     └──────┬───────┘
       │                    │                    │
       │ 1.点击"过来"        │                    │
       │───────────────────>│                    │
       │                    │                    │
       │                    │ 2.MQTT下发action   │
       │                    │───────────────────>│
       │                    │                    │
       │                    │                    │ 3.执行动作
       │                    │                    │ (移动+语音)
       │                    │                    │
       │                    │ 4.状态上报          │
       │                    │<───────────────────│
       │                    │                    │
       │ 5.更新宠物状态卡片  │                    │
       │<───────────────────│                    │
```

---

## 6. 模块联动

### 6.1 与主人画像库(OWNER_PROFILE)联动

- **触发时机：** 用户首次发送消息时
- **联动内容：** 
  - 控制台调用 `GET /api/v1/owner-profile/{user_id}` 获取主人偏好
  - 将偏好数据发送给OpenClaw对话引擎，生成个性化回复
- **数据流向：** 主人画像库 → 控制台 → 对话引擎

### 6.2 与知识库(KNOWLEDGE_BASE)联动

- **触发时机：** 用户询问天气/新闻/常识时
- **联动内容：**
  - 控制台发送消息时携带上下文
  - 知识库模块提供实时查询结果（天气/新闻）
- **数据流向：** 知识库 → OpenClaw → 控制台

### 6.3 与宠物行为引擎(PET_BEHAVIOR_ENGINE)联动

- **触发时机：** 快捷指令触发或AI回复包含动作时
- **联动内容：**
  - 控制台下发动作指令 `POST /api/v1/pets/{device_id}/actions`
  - 行为引擎生成动作序列 → MQTT下发 → MiniClaw执行
- **数据流向：** 控制台 → 行为引擎 → MiniClaw

### 6.4 与宠物记忆库(PET_MEMORY)联动

- **触发时机：** 每次对话结束后
- **联动内容：**
  - 控制台将消息元数据同步到记忆库
  - 支持对话上下文理解（短期记忆）
- **数据流向：** 控制台 → 记忆库

### 6.5 与MiniClaw通信协议(MINICLAW_PROTOCOL)联动

- **触发时机：** 所有设备通信
- **联动内容：**
  - 控制台通过MQTT协议与MiniClaw设备通信
  - 订阅设备状态Topic：`/miniclaw/{device_id}/up/status`
  - 发布动作Topic：`/miniclaw/{device_id}/down/action`
- **数据流向：** 控制台 ↔ MQTT Broker ↔ MiniClaw

---

## 7. 验收标准

### 7.1 功能验收

| 功能 | 验收条件 |
|------|----------|
| AI对话界面 | 支持文字输入，3秒内返回AI回复；支持语音输入（录制≤60秒） |
| 宠物状态展示 | 状态卡片实时更新（心跳间隔≤5秒），显示宠物名/心情/电量/表情 |
| 快捷指令面板 | 至少支持6个预设动作按钮，点击后设备在10秒内响应 |
| 宠物设置 | 可修改名字（2-8字符）、性格参数、外观配色 |
| 历史对话查询 | 支持分页查询，默认显示最近30天记录 |

### 7.2 性能验收

- 页面首屏加载时间 ≤ 2秒
- AI对话响应时间 ≤ 3秒（不含网络延迟）
- 状态更新延迟 ≤ 5秒
- 支持100个并发用户

### 7.3 兼容性验收

- 浏览器：Chrome/Firefox/Safari/Edge 最新版
- 移动端：iOS 14+ / Android 10+
- 响应式布局，支持手机/平板/PC

---

## 8. UI设计指引

### 8.1 整体布局

```
┌─────────────────────────────────────────────────────────────┐
│  Logo  │  设备选择下拉  │        搜索      │ ⚙️ │ 👤 │      │
├────────┴────────────────────────────────────────────────────┤
│          │                                                    │
│ 历史记录  │              AI对话区域                           │
│          │  ┌──────────────────────────────────────────┐    │
│ conv-01  │  │  🐱 小爪                                 │    │
│ conv-02  │  │  今天天气晴朗，适合出门散步哦！          │    │
│ conv-03  │  │                           用户 10:30 ✓✓ │    │
│          │  │  用户                                  │    │
│          │  │  今天天气怎么样？           10:30 ✓✓   │    │
│          │  └──────────────────────────────────────────┘    │
│          │                                                    │
│          │  ┌──────────────────────────────────────────────┐ │
│          │  │ 💬 输入消息...                    🎤 │ 发送 │ │
│          │  └──────────────────────────────────────────────┘ │
├──────────┴───────────────────────────────────────────────────┤
│  [过来] [坐下] [转圈] [睡觉] [投喂] [玩耍]    心情: 😊 精力: 80│
└─────────────────────────────────────────────────────────────┘
```

### 8.2 配色方案

| 用途 | 颜色 | 色值 |
|------|------|------|
| 主色 | 温暖橙 | #FF8C42 |
| 辅助色 | 柔和蓝 | #4ECDC4 |
| 背景色 | 浅米白 | #FFF8F0 |
| 文字色 | 深灰 | #2D3436 |
| 心情-开心 | 明亮黄 | #FFD93D |
| 心情-平静 | 舒适绿 | #6BCB77 |
| 心情-难过 | 忧郁蓝 | #74B9FF |
| 心情-生气 | 愤怒红 | #FF6B6B |

### 8.3 组件规范

- **聊天气泡：** 用户消息右对齐（橙色背景），AI消息左对齐（白色背景+阴影）
- **快捷按钮：** 圆角胶囊形状，高度48px，支持图标+文字
- **状态卡片：** 圆角卡片，阴影 `0 4px 12px rgba(0,0,0,0.1)`
- **心情展示：** 进度条+emoji组合，颜色随心情值渐变

### 8.4 动画规范

- **消息出现：** 淡入+上移，时长300ms，ease-out
- **宠物动作：** 与MiniClaw执行同步，设备响应后平滑过渡
- **状态更新：** 数字变化使用计数动画，防止跳跃
