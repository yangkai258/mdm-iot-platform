# Sprint 9 规划

**时间**：2026-03-22  
**状态**：进行中（前端已完成）  
**Sprint 周期**：2 周（2026-03-22 ～ 2026-04-04）  

---

## 一、Sprint 目标

**目标：** 完成 OpenClaw AI层核心功能开发，打通宠物控制台与MiniClaw设备的完整交互链路

在 Sprint 8（数据权限维度体系）的基础上，完成 OpenClaw 宠物智能相关的核心功能开发，包括宠物控制台、行为引擎、记忆库、固件管理和通信协议的对接，实现用户与AI宠物之间的完整交互体验。

---

## 二、OpenClaw 功能全景图

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          用户 (Web/App)                                     │
└─────────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                    OpenClaw 宠物控制台 (MODULE_OPENCLAW_CONSOLE)            │
│  - AI对话界面    - 宠物状态展示    - 快捷指令面板    - 宠物设置            │
└─────────────────────────────────────────────────────────────────────────────┘
          │                    │                    │                    │
          ▼                    ▼                    ▼                    ▼
┌──────────────────┐ ┌──────────────────┐ ┌──────────────────┐ ┌──────────────────┐
│ 主人画像库       │ │  知识库          │ │ 宠物行为引擎     │ │ 宠物记忆库      │
│ OWNER_PROFILE    │ │ KNOWLEDGE_BASE   │ │ PET_BEHAVIOR_ENG │ │ PET_MEMORY      │
└──────────────────┘ └──────────────────┘ └──────────────────┘ └──────────────────┘
          │                    │                    │                    │
          └──────────────────┴──────────────────┴──────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                    MiniClaw 通信协议 (MODULE_MINICLAW_PROTOCOL)              │
│  - MQTT Topic定义    - 消息格式    - 心跳管理    - OTA协议                 │
└─────────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                    MiniClaw 固件管理 (MODULE_MINICLAW_FIRMWARE)              │
│  - 固件仓库    - OTA升级    - 版本管理    - 兼容性检查                     │
└─────────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                         MiniClaw 设备 (M5Stack)                             │
│  - 动作执行    - 传感器采集    - 语音识别    - TTS播报                     │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 三、详细任务列表

### P0（必须完成）

| # | 任务 | 说明 | 交付物 | 优先级 |
|---|------|------|--------|--------|
| P0-1 | **宠物状态API实现** | 完成 `GET /api/v1/pets/{device_id}/status` 接口 | pet_controller.go | P0 |
| P0-2 | **消息发送API实现** | 完成 `POST /api/v1/pets/{device_id}/messages` 接口 | pet_controller.go | P0 |
| P0-3 | **快捷指令下发API** | 完成 `POST /api/v1/pets/{device_id}/actions` 接口 | pet_controller.go | P0 |
| P0-4 | **MQTT动作指令下发** | 实现 `/miniclaw/{device_id}/down/action` MQTT下发逻辑 | mqtt/action_publisher.go | P0 |
| P0-5 | **设备状态上报处理** | 实现 `/miniclaw/{device_id}/up/status` 上报消息处理 | mqtt/status_handler.go | P0 |
| P0-6 | **宠物状态表创建** | 创建 pet_status 表及 CRUD | models/pet_status.go + migration | P0 |
| P0-7 | **对话记录表创建** | 创建 conversations + messages 表 | models/conversation.go | P0 |
| P0-8 | **WebSocket状态推送** | 实现 `WS /api/v1/ws/pets/{device_id}` 实时推送 | websocket/pet_ws.go | P0 |

### P1（高优先级）

| # | 任务 | 说明 | 交付物 | 优先级 |
|---|------|------|--------|--------|
| P1-1 | **对话历史API** | 完成 `GET /api/v1/conversations` 和消息列表接口 | pet_controller.go | P1 |
| P1-2 | **宠物设置API** | 完成 `PUT /api/v1/pets/{device_id}/settings` | pet_controller.go | P1 |
| P1-3 | **心情激励API** | 完成 `POST /api/v1/pets/{device_id}/boost` | pet_controller.go | P1 |
| P1-4 | **动作序列生成** | 完成 `POST /api/v1/behavior/generate` | behavior_controller.go | P1 |
| P1-5 | **传感器事件处理** | 完成 `POST /api/v1/behavior/sensor-event` | behavior_controller.go | P1 |
| P1-6 | **短期记忆API** | 完成 `POST /api/v1/memory/short-term` | memory_controller.go | P1 |
| P1-7 | **上下文加载API** | 完成 `GET /api/v1/memory/context/{device_id}` | memory_controller.go | P1 |

### P2（提升体验）

| # | 任务 | 说明 | 交付物 | 优先级 |
|---|------|------|--------|--------|
| P2-1 | **表情切换API** | 完成 `PUT /api/v1/pets/{device_id}/expression` | pet_controller.go | P2 |
| P2-2 | **语音识别API** | 完成 `POST /api/v1/pets/{device_id}/voice` | pet_controller.go | P2 |
| P2-3 | **动作库查询API** | 完成 `GET /api/v1/behavior/actions` | behavior_controller.go | P2 |
| P2-4 | **决策规则管理** | 完成 CRUD `/api/v1/behavior/rules` | behavior_controller.go | P2 |
| P2-5 | **长期记忆API** | 完成 `POST/GET /api/v1/memory/long-term` | memory_controller.go | P2 |
| P2-6 | **学习记录API** | 完成 `GET /api/v1/memory/learning/{device_id}` | memory_controller.go | P2 |

---

## 四、前端任务列表

### P0（必须完成）

| # | 任务 | 说明 | 交付物 |
|---|------|------|--------|
| PF0-1 | **宠物控制台主页面** | 完成 PetConsoleView.vue 主布局 | PetConsoleView.vue |
| PF0-2 | **对话区域组件** | 完成 ChatArea.vue + ChatMessage.vue | ChatArea.vue, ChatMessage.vue |
| PF0-3 | **宠物状态卡片** | 完成 PetStatusCard.vue | PetStatusCard.vue |
| PF0-4 | **快捷指令面板** | 完成 QuickActions.vue | QuickActions.vue |
| PF0-5 | **WebSocket集成** | 完成 useWebSocket.ts | useWebSocket.ts |
| PF0-6 | **API层封装** | 完成 pet.ts API调用 | api/pet.ts |

### P1（高优先级）

| # | 任务 | 说明 | 交付物 |
|---|------|------|--------|
| PF1-1 | **历史会话列表** | 完成 ConversationList.vue | ConversationList.vue |
| PF1-2 | **宠物设置弹窗** | 完成 PetSettingsModal.vue | PetSettingsModal.vue |
| PF1-3 | **设备选择下拉** | 完成 DeviceSelector.vue | DeviceSelector.vue |
| PF1-4 | **表情选择器** | 完成 ExpressionPicker.vue | ExpressionPicker.vue |

### P2（提升体验）

| # | 任务 | 说明 | 交付物 |
|---|------|------|--------|
| PF2-1 | **语音录制组件** | 完成 VoiceRecord.vue | VoiceRecord.vue |
| PF2-2 | **宠物学习记录页** | 完成 LearningRecordView.vue | views/LearningRecordView.vue |

---

## 五、技术方案

### 5.1 架构设计

```
请求进入 API
       │
       ▼
┌─────────────────┐
│ PetController   │
│ - /pets/*      │
│ - /conversations/* │
└────────┬────────┘
         │
         ▼
┌─────────────────┐     ┌─────────────────┐
│  BehaviorEngine │ ──> │    MQTT        │
│  - generate    │     │  /down/action  │
│  - sensor_event │     └────────┬────────┘
└────────┬────────┘              │
         │                       ▼
         ▼              ┌─────────────────┐
┌─────────────────┐     │   MiniClaw     │
│  MemoryService  │     │    设备        │
│  - short_term   │     └────────┬────────┘
│  - long_term    │              │
└────────┬────────┘              │
         │                       ▼
         │              ┌─────────────────┐
         └─────────────> │    MQTT        │
                          │  /up/status   │
                          └────────┬────────┘
                                   │
                                   ▼
                          ┌─────────────────┐
                          │ StatusHandler   │
                          │ 更新宠物状态     │
                          └────────┬────────┘
                                   │
                                   ▼
                          ┌─────────────────┐
                          │  WebSocket      │
                          │  推送前端       │
                          └─────────────────┘
```

### 5.2 核心目录结构

```
backend/
├── controllers/
│   ├── pet_controller.go          # 宠物控制台API
│   ├── behavior_controller.go     # 行为引擎API
│   └── memory_controller.go       # 记忆库API
├── models/
│   ├── pet_status.go              # 宠物状态模型
│   ├── conversation.go           # 对话模型
│   ├── action_library.go          # 动作库模型
│   ├── decision_rules.go          # 决策规则模型
│   ├── short_term_memory.go       # 短期记忆模型
│   └── long_term_memory.go        # 长期记忆模型
├── services/
│   ├── pet_service.go             # 宠物服务
│   ├── behavior_service.go        # 行为引擎服务
│   ├── memory_service.go          # 记忆库服务
│   └── websocket_service.go        # WebSocket服务
├── mqtt/
│   ├── client.go                  # MQTT客户端
│   ├── action_publisher.go        # 动作下发
│   ├── status_handler.go          # 状态上报处理
│   └── sensor_handler.go          # 传感器数据处理
└── websocket/
    └── pet_ws.go                  # 宠物WebSocket处理器

frontend/src/views/pet/
├── PetConsoleView.vue
├── components/
│   ├── ChatArea.vue
│   ├── ChatMessage.vue
│   ├── ChatInput.vue
│   ├── VoiceRecord.vue
│   ├── ConversationList.vue
│   ├── QuickActions.vue
│   ├── PetStatusCard.vue
│   ├── ExpressionPicker.vue
│   ├── PetSettingsModal.vue
│   └── DeviceSelector.vue
├── composables/
│   ├── usePetStatus.ts
│   ├── useWebSocket.ts
│   ├── useConversation.ts
│   └── useQuickActions.ts
└── api/
    └── pet.ts
```

### 5.3 关键API路由

| 接口 | 方法 | 说明 |
|------|------|------|
| `/api/v1/pets/{device_id}/status` | GET | 获取宠物状态 |
| `/api/v1/pets/{device_id}/messages` | POST | 发送消息 |
| `/api/v1/pets/{device_id}/actions` | POST | 快捷指令下发 |
| `/api/v1/pets/{device_id}/settings` | PUT | 更新宠物设置 |
| `/api/v1/pets/{device_id}/boost` | POST | 心情激励 |
| `/api/v1/conversations` | GET | 会话列表 |
| `/api/v1/conversations/{id}/messages` | GET | 消息列表 |
| `/api/v1/behavior/generate` | POST | 生成动作序列 |
| `/api/v1/behavior/sensor-event` | POST | 传感器事件 |
| `/api/v1/behavior/rules` | GET/POST | 决策规则 |
| `/api/v1/memory/short-term` | POST | 短期记忆 |
| `/api/v1/memory/context/{device_id}` | GET | 获取上下文 |
| `/api/v1/memory/long-term/search` | GET | 检索长期记忆 |
| `/api/v1/ws/pets/{device_id}` | WS | WebSocket |

---

## 六、数据库设计

### 6.1 新增表

```sql
-- 宠物状态表
CREATE TABLE pet_status (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    device_id VARCHAR(64) NOT NULL UNIQUE,
    pet_name VARCHAR(32) NOT NULL DEFAULT '小爪',
    pet_type VARCHAR(16) DEFAULT 'cat',
    personality JSON,
    appearance JSON,
    mood INT DEFAULT 50,
    energy INT DEFAULT 100,
    hunger INT DEFAULT 0,
    position_x FLOAT DEFAULT 0,
    position_y FLOAT DEFAULT 0,
    current_expression VARCHAR(32) DEFAULT 'happy',
    current_action VARCHAR(32),
    is_online BOOLEAN DEFAULT FALSE,
    last_seen_at DATETIME,
    updated_at DATETIME,
    INDEX idx_device_id (device_id)
);

-- 对话记录表
CREATE TABLE conversations (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    conversation_id VARCHAR(64) NOT NULL UNIQUE,
    user_id BIGINT NOT NULL,
    device_id VARCHAR(64),
    title VARCHAR(256) NOT NULL,
    last_message TEXT,
    last_message_at DATETIME,
    message_count INT DEFAULT 0,
    status TINYINT DEFAULT 1,
    created_at DATETIME,
    updated_at DATETIME,
    INDEX idx_user_id (user_id),
    INDEX idx_device_id (device_id)
);

-- 对话消息表
CREATE TABLE messages (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    message_id VARCHAR(64) NOT NULL UNIQUE,
    conversation_id VARCHAR(64) NOT NULL,
    sender_type TINYINT NOT NULL,
    sender_id BIGINT,
    content TEXT NOT NULL,
    content_type TINYINT DEFAULT 1,
    media_url VARCHAR(512),
    intent VARCHAR(64),
    confidence FLOAT,
    metadata JSON,
    created_at DATETIME,
    INDEX idx_conversation_id (conversation_id)
);

-- 动作库表
CREATE TABLE action_library (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    action_id VARCHAR(64) NOT NULL UNIQUE,
    action_name VARCHAR(64) NOT NULL,
    action_name_en VARCHAR(64),
    category VARCHAR(32) NOT NULL,
    description TEXT,
    duration_ms INT NOT NULL,
    priority INT DEFAULT 5,
    is_emergency BOOLEAN DEFAULT FALSE,
    compatible_models JSON NOT NULL,
    parameters JSON,
    animation_data JSON,
    motor_commands JSON,
    audio_file VARCHAR(256),
    created_at DATETIME,
    updated_at DATETIME
);

-- 决策规则表
CREATE TABLE decision_rules (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    rule_id VARCHAR(64) NOT NULL UNIQUE,
    rule_name VARCHAR(128) NOT NULL,
    rule_type VARCHAR(32) NOT NULL,
    conditions JSON NOT NULL,
    condition_logic VARCHAR(8) DEFAULT 'AND',
    actions JSON NOT NULL,
    priority INT DEFAULT 50,
    is_active BOOLEAN DEFAULT TRUE,
    device_id VARCHAR(64),
    description TEXT,
    created_at DATETIME,
    updated_at DATETIME,
    INDEX idx_device_id (device_id),
    INDEX idx_is_active (is_active)
);

-- 短期记忆表
CREATE TABLE short_term_memory (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    memory_id VARCHAR(64) NOT NULL UNIQUE,
    device_id VARCHAR(64) NOT NULL,
    user_id BIGINT NOT NULL,
    session_id VARCHAR(64) NOT NULL,
    message_id VARCHAR(64),
    memory_type VARCHAR(32) NOT NULL,
    content JSON NOT NULL,
    importance FLOAT DEFAULT 0.5,
    access_count INT DEFAULT 0,
    last_accessed_at DATETIME,
    expires_at DATETIME,
    created_at DATETIME,
    INDEX idx_device_id (device_id),
    INDEX idx_session_id (session_id),
    INDEX idx_expires_at (expires_at)
);

-- 长期记忆表
CREATE TABLE long_term_memory (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    memory_id VARCHAR(64) NOT NULL UNIQUE,
    device_id VARCHAR(64) NOT NULL,
    user_id BIGINT NOT NULL,
    memory_category VARCHAR(32) NOT NULL,
    content JSON NOT NULL,
    keywords JSON,
    embedding JSON,
    confidence FLOAT DEFAULT 0.8,
    reinforcement_count INT DEFAULT 1,
    last_reinforced_at DATETIME,
    decay_score FLOAT DEFAULT 1.0,
    is_locked BOOLEAN DEFAULT FALSE,
    source_memory_id VARCHAR(64),
    created_at DATETIME,
    updated_at DATETIME,
    INDEX idx_device_id (device_id),
    INDEX idx_memory_category (memory_category)
);
```

---

## 七、MQTT Topic规划

| Topic | 方向 | 说明 | QoS |
|-------|------|------|-----|
| `/miniclaw/{device_id}/up/voice` | up | 语音识别结果 | 1 |
| `/miniclaw/{device_id}/up/sensor` | up | 传感器数据 | 1 |
| `/miniclaw/{device_id}/up/status` | up | 设备状态 | 1 |
| `/miniclaw/{device_id}/up/ack` | up | 命令确认 | 1 |
| `/miniclaw/{device_id}/down/action` | down | 动作指令 | 1 |
| `/miniclaw/{device_id}/down/speech` | down | 语音播报 | 1 |
| `/miniclaw/{device_id}/down/config` | down | 配置下发 | 1 |

---

## 八、验收标准

### 8.1 功能验收

| 用例 | 验收条件 | 测试方法 |
|------|----------|----------|
| 宠物状态展示 | 获取设备状态，返回完整pet_status | 调用API验证返回 |
| 消息发送 | 发送消息，3秒内返回AI回复 | 计时验证 |
| 快捷指令 | 点击按钮，设备10秒内响应 | 实机测试 |
| WebSocket推送 | 连接后实时接收状态更新 | WebSocket连接测试 |
| 历史对话 | 分页查询会话列表和消息 | 分页测试 |
| 动作序列生成 | 给定意图，返回有效动作序列 | 调用API验证 |
| 传感器事件 | 上报事件，100ms内响应 | 计时验证 |

### 8.2 性能验收

| 验收点 | 标准 |
|--------|------|
| API响应时间 | <= 200ms |
| MQTT消息延迟 | <= 100ms |
| WebSocket推送延迟 | <= 500ms |
| 并发支持 | 支持100个并发连接 |

### 8.3 兼容性验收

| 验收点 | 标准 |
|--------|------|
| MQTT协议 | 支持 MQTT 3.1.1 / 5.0 |
| 设备型号 | 支持 M5Stack Basic/Gray/Stack |
| 固件版本 | 支持 >= V1.2.0 |

---

## 九、依赖与风险

### 依赖

| 依赖 | 说明 |
|------|------|
| Sprint 8 数据权限 | Repository层需要支持 data_scope 过滤 |
| MiniClaw 设备 | 需要设备固件支持 MQTT 连接 |
| MQTT Broker | 需要 EMQX 服务正常运行 |
| Redis | 需要 Redis 用于设备影子和缓存 |

### 风险

| 风险 | 影响 | 缓解措施 |
|------|------|----------|
| 设备MQTT连接不稳定 | 实时性无法保证 | 增加本地缓存和重连机制 |
| AI回复延迟高 | 用户体验差 | 增加超时降级逻辑 |
| WebSocket并发大 | 服务压力 | 水平扩展和连接限制 |

---

## 十、任务分配

### 后端任务分配

| 任务 | 负责人 | 完成 |
|------|--------|------|
| 宠物状态 + 消息API | agenthd | 2026-03-23 |
| MQTT动作下发 + 状态处理 | agenthd | 2026-03-23 |
| 行为引擎API | agenthd | 2026-03-23 |
| 记忆库API | agenthd | 2026-03-23 |
| WebSocket实现 | agenthd | 2026-03-23 |
| 数据库表创建 | agenthd | 2026-03-23 |

### 前端任务分配

| 任务 | 负责人 | 完成 |
|------|--------|------|
| 宠物控制台主页面 | agentqd | 2026-03-22 |
| 对话区域组件 | agentqd | 2026-03-22 |
| 状态卡片 + 快捷指令 | agentqd | 2026-03-22 |
| WebSocket集成 | agentqd | 2026-03-22 |
| 历史会话 + 设置 | agentqd | 2026-03-22 |

### 测试任务

| 任务 | 负责人 |
|------|--------|
| API单元测试 | agentcs |
| WebSocket集成测试 | agentcs |
| MQTT消息测试 | agentcs |

---

## 十一、Sprint 9 完成清单

- [x] pet_status 表创建
- [x] conversations + messages 表创建
- [x] action_library + decision_rules 表创建
- [x] short_term_memory + long_term_memory 表创建
- [x] GET /api/v1/pets/{device_id}/status API
- [x] POST /api/v1/pets/{device_id}/messages API
- [x] POST /api/v1/pets/{device_id}/actions API
- [x] MQTT /down/action 下发实现
- [x] MQTT /up/status 处理实现
- [x] WebSocket 实时推送
- [x] PetConsoleView.vue 主页面
- [x] ChatArea.vue + ChatMessage.vue
- [x] PetStatusCard.vue + QuickActions.vue
- [x] ConversationList.vue + PetSettingsModal.vue
- [ ] API 单元测试通过
