# OpenClaw宠物控制台

**版本：** V1.1  
**模块负责人：** agentcp  
**编制日期：** 2026-03-22  

---

## 1. 概述

OpenClaw宠物控制台是AI电子宠物的Web管理界面，为用户提供与宠物对话、状态监控、快捷指令和个性化设置的能力。控制台作为OpenClaw AI层的前端入口，连接设备管理层与AI大脑，是用户感知宠物智能的核心触点。

**业务目标：**
- 提供自然语言对话界面，支持文字和语音输入
- 实时展示宠物状态（位置/表情/情绪/电量）
- 支持快捷指令快速触发预设动作
- 支持宠物个性化设置（名字/性格/外观）
- 实现宠物状态实时推送（WebSocket）

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
| 实时状态推送 | WebSocket实时接收宠物状态变化 | P0 | 自动 | 无按钮 |
| 语音录制上传 | 录制语音并上传识别 | P0 | 人工 | 输入框旁「语音」按钮 |

---

## 3. 数据模型

### 3.1 对话记录表 (conversations)

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| id | uint | 主键 | PK, auto_increment |
| conversation_id | string | 对话会话唯一标识 | unique, not null, UUID |
| user_id | uint | 所属用户 | FK to sys_users, not null |
| device_id | string | 关联设备 | FK to devices, nullable |
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
| conversation_id | string | 所属会话 | FK to conversations, not null |
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
| device_id | string | 设备ID | FK to devices, unique, not null |
| pet_name | string | 宠物名字 | not null, default 小爪 |
| pet_type | string | 宠物类型 | default cat |
| personality | json | 性格配置 | JSON |
| appearance | json | 外观配置（颜色/配件） | JSON |
| mood | int | 当前心情值 | 0-100, default 50 |
| energy | int | 精力值 | 0-100, default 100 |
| hunger | int | 饥饿值 | 0-100, default 0 |
| position_x | float | X坐标 | default 0 |
| position_y | float | Y坐标 | default 0 |
| current_expression | string | 当前表情 | default happy |
| current_action | string | 当前动作 | nullable |
| is_online | bool | 在线状态 | default false |
| last_seen_at | datetime | 最后在线时间 | auto |
| updated_at | datetime | 更新时间 | auto |

### 3.4 宠物快捷指令表 (pet_quick_actions)

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| id | uint | 主键 | PK, auto_increment |
| device_id | string | 设备ID | FK to devices, not null |
| action_id | string | 动作ID | not null |
| action_name | string | 显示名称 | not null |
| icon | string | 图标 | nullable |
| sort_order | int | 排序 | default 0 |
| is_enabled | bool | 是否启用 | default true |
| created_at | datetime | 创建时间 | auto |
| updated_at | datetime | 更新时间 | auto |

### 3.5 宠物表情配置表 (pet_expressions)

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| id | uint | 主键 | PK, auto_increment |
| expression_id | string | 表情唯一标识 | unique, not null |
| expression_name | string | 表情名称 | not null |
| emotion_value | int | 对应心情值 | 0-100 |
| icon_emoji | string | Emoji图标 | nullable |
| animation_type | string | 动画类型 | idle/happy/sad/angry/sleep |
| created_at | datetime | 创建时间 | auto |

---

## 4. 接口定义

### 4.1 获取宠物当前状态

```
GET /api/v1/pets/{device_id}/status
```

**请求参数：**

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
    "battery_level": 85,
    "last_seen_at": "2026-03-20T10:30:00Z"
  }
}
```

### 4.2 发送消息

```
POST /api/v1/pets/{device_id}/messages
```

**请求参数：**

| 参数 | 类型 | 位置 | 必填 | 说明 |
|------|------|------|------|------|
| device_id | string | path | 是 | 设备ID |
| content | string | body | 是 | 消息内容 |
| content_type | int | body | 否 | 1=文本 2=语音，默认1 |
| sender_type | int | body | 否 | 1=用户，默认1 |
| conversation_id | string | body | 否 | 所属会话ID，不传则创建新会话 |

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
    "conversation_id": "conv-uuid-001",
    "content": "今天天气晴朗，气温20-25度，适合出门散步哦！",
    "sender_type": 2,
    "intent": "weather_query",
    "confidence": 0.95,
    "metadata": {
      "emotion": "happy",
      "action": "tail_wag",
      "tts_text": "今天天气晴朗，气温20-25度，适合出门散步哦！"
    },
    "created_at": "2026-03-20T10:31:00Z"
  }
}
```

### 4.3 获取对话历史会话列表

```
GET /api/v1/conversations
```

**请求参数：**

| 参数 | 类型 | 位置 | 必填 | 说明 |
|------|------|------|------|------|
| page | int | query | 否 | 页码，默认1 |
| page_size | int | query | 否 | 每页条数，默认20，最大100 |
| device_id | string | query | 否 | 设备ID筛选 |
| status | int | query | 否 | 状态筛选 1=活跃 2=已归档 |

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {
        "conversation_id": "conv-uuid-001",
        "title": "关于天气的话题",
        "last_message": "适合出门散步哦！",
        "last_message_at": "2026-03-20T10:31:00Z",
        "message_count": 12,
        "device_id": "pet-001",
        "pet_name": "小爪",
        "status": 1
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
    "conversation_id": "conv-uuid-001",
    "items": [
      {
        "message_id": "msg-001",
        "sender_type": 1,
        "sender_name": "用户",
        "content": "今天天气怎么样？",
        "content_type": 1,
        "created_at": "2026-03-20T10:30:00Z"
      },
      {
        "message_id": "msg-002",
        "sender_type": 2,
        "sender_name": "小爪",
        "content": "今天天气晴朗，气温20-25度，适合出门散步哦！",
        "content_type": 1,
        "intent": "weather_query",
        "metadata": { "emotion": "happy", "action": "tail_wag" },
        "created_at": "2026-03-20T10:30:05Z"
      }
    ],
    "pagination": { "page": 1, "page_size": 50, "total": 12, "total_pages": 1 }
  }
}
```

### 4.5 创建新会话

```
POST /api/v1/conversations
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "conversation_id": "conv-uuid-002",
    "title": "新对话",
    "device_id": "pet-001",
    "created_at": "2026-03-20T11:00:00Z"
  }
}
```

### 4.6 删除会话

```
DELETE /api/v1/conversations/{conversation_id}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": { "conversation_id": "conv-uuid-001", "deleted_message_count": 12 }
}
```

### 4.7 更新宠物设置

```
PUT /api/v1/pets/{device_id}/settings
```

**请求示例：**
```json
{
  "pet_name": "小爪",
  "personality": { "playfulness": 80, "friendliness": 90, "curiosity": 70 },
  "appearance": { "color": "#FFB6C1", "accessories": ["bowtie", "collar"] }
}
```

### 4.8 快捷指令下发

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
    "action_name": "过来",
    "status": "queued",
    "sequence_id": "seq-uuid-001",
    "created_at": "2026-03-20T10:36:00Z"
  }
}
```

### 4.9 获取快捷指令列表

```
GET /api/v1/pets/{device_id}/quick-actions
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      { "action_id": "come_here", "action_name": "过来", "icon": "cat-come", "sort_order": 1 },
      { "action_id": "sit", "action_name": "坐下", "icon": "cat-sit", "sort_order": 2 },
      { "action_id": "spin", "action_name": "转圈", "icon": "cat-spin", "sort_order": 3 },
      { "action_id": "sleep", "action_name": "睡觉", "icon": "cat-sleep", "sort_order": 4 },
      { "action_id": "feed", "action_name": "投喂", "icon": "cat-eat", "sort_order": 5 },
      { "action_id": "play", "action_name": "玩耍", "icon": "cat-play", "sort_order": 6 }
    ]
  }
}
```

### 4.10 心情激励

```
POST /api/v1/pets/{device_id}/boost
```

**请求示例：**
```json
{ "type": "feed" }
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "device_id": "pet-001",
    "type": "feed",
    "hunger_before": 70,
    "hunger_after": 30,
    "mood_change": "+10",
    "action_executed": "eating_animation",
    "message": "小爪吃得很开心！"
  }
}
```

### 4.11 获取表情列表

```
GET /api/v1/pets/expressions
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      { "expression_id": "happy", "expression_name": "开心", "emotion_value": 80, "icon_emoji": "😸" },
      { "expression_id": "sad", "expression_name": "难过", "emotion_value": 30, "icon_emoji": "😿" },
      { "expression_id": "angry", "expression_name": "生气", "emotion_value": 20, "icon_emoji": "😾" },
      { "expression_id": "sleep", "expression_name": "睡觉", "emotion_value": 50, "icon_emoji": "😴" }
    ]
  }
}
```

### 4.12 切换宠物表情

```
PUT /api/v1/pets/{device_id}/expression
```

**请求示例：**
```json
{ "expression": "happy" }
```

### 4.13 语音上传识别

```
POST /api/v1/pets/{device_id}/voice
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "message_id": "msg-uuid-voice-001",
    "text": "今天天气怎么样",
    "confidence": 0.95,
    "duration_ms": 2500,
    "audio_url": "https://cdn.example.com/audio/msg-uuid-voice-001.wav"
  }
}
```

### 4.14 WebSocket实时状态推送

```
WS /api/v1/ws/pets/{device_id}?token={jwt_token}
```

**订阅消息类型：**

| type | 说明 | payload |
|------|------|---------|
| status_update | 宠物状态变更 | 完整pet_status对象 |
| action_result | 动作执行结果 | action_id + status |
| message | 新消息(AI回复) | message对象 |
| emotion_change | 情绪变化 | mood + emotion值 |
| ota_progress | OTA升级进度 | progress_percent |

**服务端推送示例（状态更新）：**
```json
{
  "type": "status_update",
  "data": {
    "device_id": "pet-001",
    "mood": 78,
    "energy": 75,
    "expression": "happy",
    "action": "tail_wag",
    "is_online": true,
    "battery_level": 82,
    "updated_at": "2026-03-20T10:32:00Z"
  }
}
```

### 4.15 错误码定义

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 10001 | 设备不在线 |
| 10002 | 设备不存在 |
| 10003 | 会话不存在 |
| 10004 | 消息内容为空 |
| 10005 | 语音识别失败 |
| 10006 | 动作执行失败 |
| 10007 | 宠物不在线，无法发送消息 |

---

## 5. 流程图

### 5.1 用户对话完整流程

用户输入 -> 控制台前端 -> OpenClaw AI大脑
         1.输入消息         
         2.发送消息请求 -> 3.查询主人画像
                         -> 4.查询知识库
                         -> 5.生成回复+动作
         6.返回AI回复        
         7.展示回复+动画     
         8.MQTT下发动作指令  -> MiniClaw设备
         9.状态上报->记忆更新
         10.更新界面状态

### 5.2 快捷指令流程

点击快捷按钮 -> 控制台前端 -> MiniClaw设备
1.点击过来     2.MQTT下发action  3.执行动作(移动+语音)
              4.状态上报
5.更新宠物状态卡片

### 5.3 WebSocket实时状态推送流程

页面加载 -> WebSocket Service -> MQTT Broker
建立连接    订阅device_id    -> MiniClaw设备
                            状态上报
           转发状态消息
           推送前端更新 -> 前端更新UI

---

## 6. 前端页面详细设计

### 6.1 主页面布局（PetConsoleView）

整体布局分为：顶部导航栏、左侧历史会话区、中央对话区、底部快捷指令+状态卡片

**顶部导航栏：**
- Logo
- 设备选择下拉（显示宠物名-device_id）
- 搜索对话输入框
- 设置图标
- 用户头像

**左侧历史会话栏：**
- 新建对话按钮
- 会话列表（按时间倒序）
- 每条显示标题+最后消息预览+时间

**中央对话区：**
- 消息气泡区（用户右对齐橙色，AI左对齐白色）
- 输入框（支持文字+语音切换）
- 发送按钮

**底部：**
- 快捷指令面板（6+3个预设动作按钮）
- 宠物状态卡片（心情/精力/饥饿/电量/在线状态）

### 6.2 宠物设置弹窗

包含：基本信息（名字/类型）、性格配置（活泼/友善/好奇）、外观定制（颜色/配件）、声音设置（音色/音量）

### 6.3 设备选择下拉

显示所有已绑定设备，在线状态用绿色/灰色圆点标识

### 6.4 语音录制组件

按住说话样式，显示录制时长和波形动画

### 6.5 组件状态说明

| 组件 | 状态 | 说明 |
|------|------|------|
| 快捷按钮 | default | 默认可用 |
| 快捷按钮 | loading | 指令下发中，显示loading动画 |
| 快捷按钮 | disabled | 设备离线时禁用 |
| 快捷按钮 | success | 指令执行成功，绿色闪烁 |
| 快捷按钮 | failed | 指令执行失败，红色提示 |
| 消息气泡 | sending | 发送中，灰色 |
| 消息气泡 | sent | 已发送，单勾 |
| 消息气泡 | delivered | 已送达，双勾蓝色 |
| 消息气泡 | failed | 发送失败，红色可重发 |
| 状态卡片 | online | 绿色圆点 |
| 状态卡片 | offline | 灰色圆点 |

### 6.6 前端目录结构

```
frontend/src/views/pet/
├── PetConsoleView.vue          # 主控制台页面
├── components/
│   ├── ChatArea.vue            # 对话区域组件
│   ├── ChatMessage.vue         # 单条消息组件
│   ├── ChatInput.vue           # 输入框组件
│   ├── VoiceRecord.vue         # 语音录制组件
│   ├── ConversationList.vue    # 历史会话列表
│   ├── QuickActions.vue        # 快捷指令面板
│   ├── PetStatusCard.vue       # 宠物状态卡片
│   ├── ExpressionPicker.vue    # 表情选择器
│   ├── PetSettingsModal.vue    # 宠物设置弹窗
│   └── DeviceSelector.vue      # 设备选择下拉
├── composables/
│   ├── usePetStatus.ts         # 宠物状态管理
│   ├── useWebSocket.ts         # WebSocket连接
│   ├── useConversation.ts      # 对话管理
│   └── useQuickActions.ts      # 快捷指令
├── api/
│   └── pet.ts                  # 宠物控制台API
└── types/
    └── pet.ts                  # 类型定义
```

---

## 7. 模块联动

### 7.1 与主人画像库(OWNER_PROFILE)联动

- **触发时机：** 用户首次发送消息时
- **联动内容：** 获取主人偏好用于生成个性化回复
- **数据流向：** 主人画像库 -> 控制台 -> 对话引擎

### 7.2 与知识库(KNOWLEDGE_BASE)联动

- **触发时机：** 用户询问天气/新闻/常识时
- **联动内容：** 知识库提供实时查询结果
- **数据流向：** 知识库 -> OpenClaw -> 控制台

### 7.3 与宠物行为引擎(PET_BEHAVIOR_ENGINE)联动

- **触发时机：** 快捷指令触发或AI回复包含动作时
- **联动内容：** 行为引擎生成动作序列 -> MQTT下发 -> MiniClaw执行
- **数据流向：** 控制台 -> 行为引擎 -> MiniClaw

### 7.4 与宠物记忆库(PET_MEMORY)联动

- **触发时机：** 每次对话结束后
- **联动内容：** 控制台将消息元数据同步到记忆库
- **数据流向：** 控制台 -> 记忆库

### 7.5 与MiniClaw通信协议(MINICLAW_PROTOCOL)联动

- **触发时机：** 所有设备通信
- **联动内容：** 通过MQTT协议与MiniClaw设备通信
- **数据流向：** 控制台 <-> MQTT Broker <-> MiniClaw

---

## 8. 验收标准

### 8.1 功能验收

| 功能 | 验收条件 |
|------|----------|
| AI对话界面 | 支持文字输入，3秒内返回AI回复；支持语音输入（录制<=60秒） |
| 宠物状态展示 | 状态卡片实时更新（心跳间隔<=5秒），显示宠物名/心情/电量/表情 |
| 快捷指令面板 | 至少支持6个预设动作按钮，点击后设备在10秒内响应 |
| 宠物设置 | 可修改名字（2-8字符）、性格参数、外观配色 |
| 历史对话查询 | 支持分页查询，默认显示最近30天记录 |
| WebSocket推送 | 建立连接后实时接收状态更新，延迟<=1秒 |

### 8.2 性能验收

- 页面首屏加载时间 <= 2秒
- AI对话响应时间 <= 3秒（不含网络延迟）
- 状态更新延迟 <= 5秒
- 支持100个并发用户

### 8.3 兼容性验收

- 浏览器：Chrome/Firefox/Safari/Edge 最新版
- 移动端：iOS 14+ / Android 10+
- 响应式布局，支持手机/平板/PC

---

## 9. 配色方案

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


---

## 10. 页面布局规范

### 10.1 宠物控制台主页面（PetConsoleView）

**布局结构：**（非标准三段式 — 全屏沉浸式对话布局）
1. 顶部导航栏（Logo / 设备选择下拉 / 搜索 / 设置 / 用户头像）
2. 左侧边栏（历史会话列表，可折叠）
3. 中央对话区（消息气泡 + 输入框）
4. 底部快捷指令面板 + 宠物状态卡片

**非标准说明：** 本页面为对话型沉浸式界面，不适用标准三段式布局，采用全屏左右分栏布局。

**快捷指令面板按钮：**
- [过来] [坐下] [转圈] [睡觉] [投喂] [玩耍] — 底部横向排列，均等分布
- 设备离线时全部禁用（disabled）

### 10.2 宠物设置页面（Modal）

**布局结构：**
1. Modal 全屏模态框
2. Tab 页签：基本信息 / 性格配置 / 外观定制 / 声音设置

**按钮规范：**
- [取消] — 左下角
- [保存] — 右下角

### 10.3 历史会话列表（侧边栏）

**布局结构：**
1. 新建对话按钮（顶部）
2. 会话列表（按时间倒序，每条显示：标题 + 最后消息预览 + 时间）
3. 每条支持右键菜单（删除/归档）

### 10.4 弹窗规范

| 类型 | 使用场景 |
|------|----------|
| 全屏模态 | 宠物设置弹窗 |
| Dialog 对话框 | 确认删除会话 |
