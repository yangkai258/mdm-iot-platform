# 宠物记忆库

**版本：** V1.1  
**模块负责人：** agentcp  
**编制日期：** 2026-03-22  

---

## 1. 概述

宠物记忆库是OpenClaw AI层的持久化存储核心，负责保存宠物与主人交互的所有记忆数据。记忆库分为短期记忆（当前对话上下文）和长期记忆（重要事件和知识），支持对话上下文理解、宠物学习记录和记忆检索与召回，使宠物能够记住重要信息并持续学习。

**业务目标：**
- 实现对话上下文的短期记忆理解
- 持久化保存重要交互记忆
- 支持宠物长期学习与成长
- 提供记忆检索和召回能力

---

## 2. 功能列表

| 功能 | 描述 | 优先级 | 触发方式 | 前端入口/按钮 |
|------|------|--------|----------|--------------|
| 短期记忆存储 | 保存当前对话上下文（最近20轮） | P0 | 自动 | 无按钮 |
| 上下文加载 | 加载历史上下文用于理解当前对话 | P0 | 自动 | 无按钮 |
| 长期记忆写入 | 保存重要事件和知识到长期记忆 | P0 | 自动 | 无按钮 |
| 记忆检索 | 根据关键词/时间/类型检索记忆 | P1 | 人工 | 控制台「记忆查询」 |
| 记忆遗忘 | 自动清理过期的短期记忆（7天后） | P1 | 自动 | 无按钮 |
| 记忆强化 | 重要记忆自动强化权重 | P1 | 自动 | 无按钮 |
| 学习记录查看 | 查看宠物学习历史和能力成长 | P2 | 人工 | 控制台「学习记录」 |
| 记忆导出 | 导出记忆数据用于分析 | P2 | 人工 | 管理后台 |

---

## 3. 数据模型

### 3.1 短期记忆表 (short_term_memory)

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| id | uint | 主键 | PK, auto_increment |
| memory_id | string | 记忆唯一标识 | unique, not null |
| device_id | string | 设备ID | FK to devices, not null |
| user_id | uint | 用户ID | FK to sys_users, not null |
| session_id | string | 会话ID | not null |
| message_id | string | 关联消息ID | nullable |
| memory_type | string | 记忆类型 | intent/emotion/fact/action/event |
| content | json | 记忆内容 | not null |
| importance | float | 重要程度 | 0-1, default 0.5 |
| access_count | int | 访问次数 | default 0 |
| last_accessed_at | datetime | 最后访问时间 | auto |
| expires_at | datetime | 过期时间 | auto (7天后) |
| created_at | datetime | 创建时间 | auto |

### 3.2 长期记忆表 (long_term_memory)

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| id | uint | 主键 | PK, auto_increment |
| memory_id | string | 记忆唯一标识 | unique, not null |
| device_id | string | 设备ID | FK to devices, not null |
| user_id | uint | 用户ID | FK to sys_users, not null |
| memory_category | string | 记忆分类 | person/place/thing/event/preference/skill |
| content | json | 记忆内容 | not null |
| keywords | json | 关键词索引 | nullable |
| embedding | json | 向量化表示 | nullable |
| confidence | float | 置信度 | 0-1, default 0.8 |
| reinforcement_count | int | 强化次数 | default 1 |
| last_reinforced_at | datetime | 最后强化时间 | auto |
| decay_score | float | 衰减分数 | default 1.0 |
| is_locked | bool | 是否锁定（不自动删除） | default false |
| source_memory_id | string | 来源记忆ID | nullable |
| created_at | datetime | 创建时间 | auto |
| updated_at | datetime | 更新时间 | auto |

### 3.3 记忆索引表 (memory_index)

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| id | uint | 主键 | PK, auto_increment |
| index_id | string | 索引唯一标识 | unique, not null |
| memory_id | string | 关联记忆ID | FK to long_term_memory, not null |
| index_type | string | 索引类型 | keyword/entity/time/location |
| index_key | string | 索引键 | not null |
| index_value | json | 索引值 | not null |
| weight | float | 权重 | default 1.0 |
| created_at | datetime | 创建时间 | auto |

### 3.4 学习记录表 (learning_records)

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| id | uint | 主键 | PK, auto_increment |
| record_id | string | 记录唯一标识 | unique, not null |
| device_id | string | 设备ID | FK to devices, not null |
| user_id | uint | 用户ID | FK to sys_users, not null |
| skill_type | string | 技能类型 | language/trick/recognition/response |
| skill_name | string | 技能名称 | not null |
| learning_stage | int | 学习阶段 | 1=初学 2=练习 3=掌握 4=熟练 |
| proficiency | float | 熟练度 | 0-100, auto |
| practice_count | int | 练习次数 | default 0 |
| success_rate | float | 成功率 | 0-1, auto |
| milestone | json | 里程碑记录 | nullable |
| started_at | datetime | 开始学习时间 | auto |
| mastered_at | datetime | 熟练掌握时间 | nullable |
| created_at | datetime | 创建时间 | auto |
| updated_at | datetime | 更新时间 | auto |

### 3.5 对话历史表 (conversation_history)

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| id | uint | 主键 | PK, auto_increment |
| history_id | string | 历史记录ID | unique, not null |
| device_id | string | 设备ID | FK to devices, not null |
| user_id | uint | 用户ID | FK to sys_users, not null |
| conversation_id | string | 会话ID | FK to conversations |
| messages | json | 消息列表 | not null |
| summary | string | 对话摘要 | nullable |
| context_window | int | 上下文窗口大小 | default 20 |
| created_at | datetime | 创建时间 | auto |
| updated_at | datetime | 更新时间 | auto |

---

## 4. 接口定义

### 4.1 写入短期记忆

```
POST /api/v1/memory/short-term
```

**请求参数：**

| 参数 | 类型 | 位置 | 必填 | 说明 |
|------|------|------|------|------|
| device_id | string | body | 是 | 设备ID |
| user_id | uint | body | 是 | 用户ID |
| session_id | string | body | 是 | 会话ID |
| message_id | string | body | 否 | 关联消息ID |
| memory_type | string | body | 是 | 记忆类型 |
| content | json | body | 是 | 记忆内容 |
| importance | float | body | 否 | 重要程度 0-1 |

**请求示例：**
```json
{
  "device_id": "pet-001",
  "user_id": 10001,
  "session_id": "sess-001",
  "message_id": "msg-001",
  "memory_type": "intent",
  "content": {
    "user_message": "今天是我的生日",
    "recognized_intent": "celebration",
    "entities": { "occasion": "birthday" }
  },
  "importance": 0.8
}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "memory_id": "stm-001",
    "expires_at": "2026-03-29T10:00:00Z"
  }
}
```

### 4.2 获取对话上下文

```
GET /api/v1/memory/context/{device_id}
```

**请求参数：**

| 参数 | 类型 | 位置 | 必填 | 说明 |
|------|------|------|------|------|
| device_id | string | path | 是 | 设备ID |
| session_id | string | query | 否 | 会话ID，不传则获取最新会话 |
| limit | int | query | 否 | 返回条数，默认20 |

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "device_id": "pet-001",
    "session_id": "sess-001",
    "context": [
      { "memory_id": "stm-020", "memory_type": "intent", "content": { "user_message": "今天天气不错" }, "created_at": "2026-03-20T10:55:00Z" },
      { "memory_id": "stm-019", "memory_type": "action", "content": { "action": "walk_forward", "result": "success" }, "created_at": "2026-03-20T10:54:30Z" }
    ],
    "total_count": 15
  }
}
```

### 4.3 写入长期记忆

```
POST /api/v1/memory/long-term
```

**请求参数：**

| 参数 | 类型 | 位置 | 必填 | 说明 |
|------|------|------|------|------|
| device_id | string | body | 是 | 设备ID |
| user_id | uint | body | 是 | 用户ID |
| memory_category | string | body | 是 | 记忆分类 |
| content | json | body | 是 | 记忆内容 |
| keywords | json | body | 否 | 关键词索引 |
| confidence | float | body | 否 | 置信度 0-1 |
| is_locked | bool | body | 否 | 是否锁定 |

**请求示例：**
```json
{
  "device_id": "pet-001",
  "user_id": 10001,
  "memory_category": "preference",
  "content": {
    "type": "favorite_topic",
    "value": "宠物训练",
    "evidence": ["用户多次询问宠物训练方法"]
  },
  "keywords": ["宠物", "训练", "喜好"],
  "confidence": 0.85,
  "is_locked": true
}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "memory_id": "ltm-001",
    "memory_category": "preference",
    "confidence": 0.85,
    "created_at": "2026-03-22T10:00:00Z"
  }
}
```

### 4.4 检索长期记忆

```
GET /api/v1/memory/long-term/search
```

**请求参数：**

| 参数 | 类型 | 位置 | 必填 | 说明 |
|------|------|------|------|------|
| device_id | string | query | 是 | 设备ID |
| keyword | string | query | 否 | 关键词搜索 |
| category | string | query | 否 | 记忆分类 |
| start_date | string | query | 否 | 开始日期 |
| end_date | string | query | 否 | 结束日期 |
| page | int | query | 否 | 页码 |
| page_size | int | query | 否 | 每页条数 |

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {
        "memory_id": "ltm-001",
        "memory_category": "preference",
        "content": { "type": "favorite_topic", "value": "宠物训练" },
        "confidence": 0.85,
        "reinforcement_count": 3,
        "created_at": "2026-03-15T10:00:00Z"
      }
    ],
    "pagination": { "page": 1, "page_size": 20, "total": 45, "total_pages": 3 }
  }
}
```

### 4.5 强化记忆

```
POST /api/v1/memory/{memory_id}/reinforce
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": { 
    "memory_id": "ltm-001", 
    "reinforcement_count": 5, 
    "confidence": 0.92 
  }
}
```

### 4.6 删除记忆

```
DELETE /api/v1/memory/{memory_id}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": { "memory_id": "ltm-001", "deleted": true }
}
```

### 4.7 获取学习记录

```
GET /api/v1/memory/learning/{device_id}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "device_id": "pet-001",
    "skills": [
      { "skill_type": "recognition", "skill_name": "认识主人", "learning_stage": 4, "proficiency": 95 },
      { "skill_type": "language", "skill_name": "日常对话", "learning_stage": 3, "proficiency": 78 }
    ],
    "total_learning_time_days": 45
  }
}
```

### 4.8 更新学习记录

```
PUT /api/v1/memory/learning/{record_id}
```

**请求示例：**
```json
{
  "practice_count": 10,
  "success_rate": 0.85
}
```

### 4.9 保存对话历史

```
POST /api/v1/memory/conversation-history
```

**请求参数：**

| 参数 | 类型 | 位置 | 必填 | 说明 |
|------|------|------|------|------|
| device_id | string | body | 是 | 设备ID |
| user_id | uint | body | 是 | 用户ID |
| conversation_id | string | body | 是 | 会话ID |
| messages | json | body | 是 | 消息列表 |

**请求示例：**
```json
{
  "device_id": "pet-001",
  "user_id": 10001,
  "conversation_id": "conv-001",
  "messages": [
    { "sender_type": 1, "content": "今天天气怎么样？", "created_at": "2026-03-22T10:00:00Z" },
    { "sender_type": 2, "content": "今天天气晴朗！", "created_at": "2026-03-22T10:00:05Z" }
  ]
}
```

### 4.10 获取对话历史

```
GET /api/v1/memory/conversation-history
```

**请求参数：**

| 参数 | 类型 | 位置 | 必填 | 说明 |
|------|------|------|------|------|
| device_id | string | query | 是 | 设备ID |
| conversation_id | string | query | 否 | 会话ID |
| page | int | query | 否 | 页码 |

### 4.11 错误码定义

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 30001 | 记忆不存在 |
| 30002 | 记忆写入失败 |
| 30003 | 检索参数错误 |
| 30004 | 学习记录不存在 |

---

## 5. 流程图

### 5.1 短期记忆生命周期

```
用户消息 -> 对话引擎 -> 记忆库
       │
       ▼
1. 写入短期记忆
   - 记忆类型识别
   - 重要性评分
   - 设置过期时间(7天)
       │
       ▼
2. 更新上下文窗口
   - 保留最近20轮
   - 更新last_accessed
       │
       ▼
3. 判断是否转入长期
   - 重要性>0.8?
   - 重复出现3次以上?
   - 是主人偏好信息?
       │
       ├─ 是 -> 转入长期记忆
       │    - 生成embedding
       │    - 建立关键词索引
       └─ 否 -> 保持短期记忆
            - 等待自然遗忘
            - 或被强化
```

### 5.2 记忆遗忘与强化流程

```
定时任务(每日)
       │
       ▼
扫描短期记忆过期项
expires_at < now
       │
       ├─ 重要性 >= 0.7 -> 转入长期记忆
       └─ 重要性 < 0.7 -> 直接删除
```

---

## 6. 模块联动

### 6.1 与对话引擎联动

- **触发时机：** 每次对话时
- **联动内容：** 写入短期记忆保存对话内容，加载上下文理解当前对话
- **数据流向：** 对话引擎 -> 记忆库 -> 对话引擎

### 6.2 与主人画像库(OWNER_PROFILE)联动

- **触发时机：** 识别到主人偏好时
- **联动内容：** 记忆库中的偏好信息同步到画像库，画像库更新偏好规则
- **数据流向：** 记忆库 -> 主人画像库

### 6.3 与行为引擎(PET_BEHAVIOR_ENGINE)联动

- **触发时机：** 动作执行完成时
- **联动内容：** 记录动作执行结果到记忆库，用于学习动作效果
- **数据流向：** 行为引擎 -> 记忆库

### 6.4 与知识库(KNOWLEDGE_BASE)联动

- **触发时机：** 知识查询时
- **联动内容：** 记忆库提供上下文辅助理解问题，知识库结果存入长期记忆
- **数据流向：** 知识库 <-> 记忆库

---

## 7. 验收标准

### 7.1 功能验收

| 功能 | 验收条件 |
|------|----------|
| 短期记忆存储 | 每次对话自动保存，保留最近20轮 |
| 上下文加载 | 支持加载最近7天上下文 |
| 长期记忆 | 重要记忆永久保存，支持检索 |
| 记忆遗忘 | 7天后自动过期清理 |
| 记忆强化 | 重复出现3次以上自动强化 |
| 对话历史 | 自动保存会话消息，支持分页查询 |

### 7.2 性能验收

- 短期记忆写入延迟 <= 50ms
- 上下文加载延迟 <= 100ms
- 支持10000次/分钟记忆操作

### 7.3 数据质量验收

- 记忆数据完整率 >= 99.9%
- 上下文理解准确率 >= 90%

---

## 8. 前端页面设计

### 8.1 记忆查询界面

```
+---------------------------------------------------------------+
| 记忆查询                                                       |
+---------------------------------------------------------------+
| 设备: [pet-001 v]  分类: [全部 v]  关键词: [____________] [搜索]|
| 日期: [2026-03-01] - [2026-03-22]                            |
+---------------------------------------------------------------+
|                                                               |
| [v] 长期记忆  [v] 短期记忆  [ ] 学习记录                       |
|                                                               |
+---------------------------------------------------------------+
| 检索结果: 12条                                                 |
|                                                               |
| +-----------------------------------------------------------+ |
| | [锁定] 偏好 - 宠物训练                                     | |
| | 内容: 用户喜欢讨论宠物训练相关话题                          | |
| | 置信度: 92% | 强化次数: 5 | 创建: 2026-03-15              | |
| | [查看详情] [强化] [删除]                                   | |
| +-----------------------------------------------------------+ |
|                                                               |
| +-----------------------------------------------------------+ |
| | 事件 - 生日                                                | |
| | 内容: 今天是我的生日                                        | |
| | 置信度: 85% | 来源: 短期记忆转入 | 创建: 2026-03-22        | |
| | [查看详情] [强化] [删除]                                   | |
| +-----------------------------------------------------------+ |
+---------------------------------------------------------------+
```

### 8.2 学习记录界面

```
+---------------------------------------------------------------+
| 宠物学习记录                                  pet-001         |
+---------------------------------------------------------------+
|                                                               |
| 总学习时长: 45天  当前技能: 8个  熟练度: 78%                  |
|                                                               |
| 技能概览                                                       |
| +-----------------------------------------------------------+ |
| | 认识主人    [████████████████████] 95%  熟练              | |
| | 日常对话    [██████████████░░░░░░] 78%  掌握              | |
| | 情绪识别    [██████████░░░░░░░░░░░] 55%  练习中            | |
| | 定点动作    [██████░░░░░░░░░░░░░░░░] 35%  初学              | |
| +-----------------------------------------------------------+ |
|                                                               |
| 学习里程碑                                                     |
| - 2026-03-01: 学会识别主人声音                                |
| - 2026-03-10: 掌握"过来"指令                                   |
| - 2026-03-15: 能够回应简单问候                                 |
+---------------------------------------------------------------+
```

### 8.3 配色方案

| 用途 | 颜色 | 色值 |
|------|------|------|
| 短期记忆 | 浅蓝 | #74B9FF |
| 长期记忆 | 金色 | #FFD93D |
| 锁定记忆 | 紫色 | #A29BFE |
| 学习技能 | 绿色 | #6BCB77 |
| 遗忘记忆 | 灰色 | #B2BEC3 |


---

## 9. 页面布局规范

### 9.1 记忆查询页面

**布局结构：**
1. 面包屑 → 页面标题
2. 筛选区（浅灰卡片 #F2F3F5）：设备 / 记忆类型 / 关键词 / 日期范围
3. 操作栏（导出靠右）
4. 记忆检索结果列表（卡片形式：长期记忆/短期记忆/学习记录 Tab）

**按钮规范：**
- [导出] — 右对齐

**卡片内容：**
- 记忆类型标签（锁定图标 / 重要程度）
- 记忆内容摘要
- 置信度 / 强化次数 / 创建时间
- 操作：[查看详情] [强化] [删除]

**分页：** 右下角，10/20/50/100 条

### 9.2 学习记录页面

**布局结构：**
1. 面包屑 → 页面标题 + 设备信息
2. 统计卡片（总学习时长 / 当前技能数 / 熟练度）—— 白色
3. 技能概览（进度条列表）
4. 学习里程碑时间线

**分页：** 不适用（里程碑为时间线形式）

### 9.3 弹窗规范

| 类型 | 使用场景 |
|------|----------|
| Drawer 抽屉 | 记忆详情查看 |
| Dialog 对话框 | 确认删除记忆 |
| 全屏模态 | 暂无复杂表单场景 |
