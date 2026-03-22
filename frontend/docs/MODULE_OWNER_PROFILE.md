# 主人画像库

**版本：** V1.0  
**模块负责人：** agentcp  
**编制日期：** 2026-03-20  

---

## 1. 概述

主人画像库是OpenClaw AI层的个性化核心模块，负责学习和管理每位用户与宠物交互的偏好数据。通过分析交互历史和行为模式，画像库为对话引擎提供千人千面的个性化回复策略，实现宠物与主人之间更亲密的互动体验。

**业务目标：**
- 学习主人的称呼偏好（喜欢被叫什么）
- 分析主人的活跃时间和话题偏好
- 为宠物生成个性化回复风格
- 支持宠物性格养成系统

---

## 2. 功能列表

| 功能 | 描述 | 优先级 | 触发方式 | 前端入口/按钮 |
|------|------|--------|----------|--------------|
| 主人信息管理 | 维护主人基本信息（昵称/头像/生日） | P0 | 人工 | 控制台「个人设置」 |
| 称呼偏好学习 | 分析用户互动，学习喜欢被称呼的方式 | P0 | 自动 | 无按钮 |
| 活跃时间分析 | 统计用户活跃时间段，优化互动时机 | P0 | 自动 | 无按钮 |
| 话题偏好分析 | 分析用户常聊话题，个性化推荐 | P1 | 自动 | 无按钮 |
| 个性化回复生成 | 根据画像生成符合主人风格的回复 | P0 | 自动 | 无按钮 |
| 宠物性格配置 | 配置宠物初始性格参数 | P1 | 人工 | 控制台「宠物设置」 |
| 宠物性格养成 | 随交互时长，宠物性格逐渐变化 | P1 | 自动 | 无按钮 |
| 画像数据导出 | 导出主人画像统计数据 | P2 | 人工 | 管理后台 |

---

## 3. 数据模型

### 3.1 主人画像表 (owner_profiles)

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| id | uint | 主键 | PK, auto_increment |
| user_id | uint | 用户ID | FK → sys_users, unique, not null |
| nickname | string | 昵称 | not null |
| avatar_url | string | 头像URL | nullable |
| birthday | date | 生日 | nullable |
| gender | string | 性别 | nullable |
| preferred_name | string | 希望被称呼的名字 | nullable |
| name_style | int | 称呼风格 | 1=正式 2=昵称 3=宝宝式 |
| activity_pattern | json | 活跃时段分布 | auto |
| interaction_count | int | 累计交互次数 | default 0 |
| last_interaction_at | datetime | 最后交互时间 | auto |
| personality_config | json | 个性化配置 | default {} |
| pet_personality | json | 宠物性格参数 | auto |
| tags | json | 用户标签 | default [] |
| created_at | datetime | 创建时间 | auto |
| updated_at | datetime | 更新时间 | auto |

### 3.2 偏好规则表 (preference_rules)

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| id | uint | 主键 | PK, auto_increment |
| rule_id | string | 规则唯一标识 | unique, not null |
| user_id | uint | 关联用户 | FK → sys_users, nullable（null表示通用） |
| preference_type | string | 偏好类型 | call_name/topic/time/style/response_length |
| preference_key | string | 偏好键 | not null |
| preference_value | json | 偏好值 | not null |
| confidence | float | 置信度 | 0-1, default 0.5 |
| sample_count | int | 样本数量 | default 1 |
| is_active | bool | 是否启用 | default true |
| expires_at | datetime | 过期时间 | nullable |
| created_at | datetime | 创建时间 | auto |
| updated_at | datetime | 更新时间 | auto |

### 3.3 交互历史表 (interaction_history)

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| id | uint | 主键 | PK, auto_increment |
| interaction_id | string | 交互唯一标识 | unique, not null |
| user_id | uint | 用户ID | FK → sys_users, not null |
| device_id | string | 设备ID | FK → devices, nullable |
| interaction_type | string | 交互类型 | chat/command/sensor/auto |
| user_message | text | 用户消息 | nullable |
| pet_response | text | 宠物回复 | nullable |
| intent | string | 识别意图 | nullable |
| topics | json | 涉及话题 | nullable |
| sentiment_score | float | 情感分值 | -1到1 |
| response_style | string | 回复风格 | nullable |
| duration_ms | int | 交互耗时 | nullable |
| time_bucket | string | 时间段标识 | morning/afternoon/evening/night |
| created_at | datetime | 交互时间 | auto |

### 3.4 宠物性格表 (pet personalities - 内嵌于owner_profiles.pet_personality)

```json
{
  "playfulness": 70,
  "friendliness": 85,
  "curiosity": 60,
  "loyalty": 90,
  "bravery": 50,
  "energy_level": 75,
  "mood_stability": 65,
  "learning_speed": 80,
  "bond_strength": 60
}
```

---

## 4. 接口定义

### 4.1 获取主人画像

```
GET /api/v1/owner-profile/{user_id}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "user_id": 10001,
    "nickname": "小明",
    "preferred_name": "小明",
    "name_style": 2,
    "activity_pattern": {
      "morning": 20,
      "afternoon": 30,
      "evening": 45,
      "night": 5
    },
    "interaction_count": 1250,
    "pet_personality": {
      "playfulness": 75,
      "friendliness": 88,
      "loyalty": 92
    },
    "tags": ["宠物爱好者", "上班族", "早晨型"]
  }
}
```

### 4.2 更新主人偏好

```
PUT /api/v1/owner-profile/{user_id}/preferences
```

**请求示例：**
```json
{
  "preferred_name": "老板",
  "name_style": 1,
  "personality_config": {
    "response_length": "short",
    "tone": "formal"
  }
}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": { "user_id": 10001, "updated_at": "2026-03-20T10:40:00Z" }
}
```

### 4.3 记录交互

```
POST /api/v1/owner-profile/interactions
```

**请求示例：**
```json
{
  "user_id": 10001,
  "device_id": "pet-001",
  "interaction_type": "chat",
  "user_message": "今天心情不好",
  "pet_response": "怎么了呀？要不要听我唱歌？",
  "intent": " venting",
  "topics": ["情绪", "音乐"],
  "sentiment_score": -0.3,
  "time_bucket": "evening"
}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": { "interaction_id": "int-001", "profile_updated": true }
}
```

### 4.4 获取个性化回复参数

```
GET /api/v1/owner-profile/{user_id}/reply-config
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "user_id": 10001,
    "preferred_name": "小明",
    "name_style": 2,
    "response_style": {
      "length": "medium",
      "tone": "warm",
      "emoji_frequency": "high",
      "use_pet_name": true
    },
    "active_hours": ["morning", "evening"],
    "preferred_topics": ["宠物", "音乐", "电影"]
  }
}
```

### 4.5 宠物性格养成更新

```
POST /api/v1/owner-profile/{user_id}/personality-evolve
```

**请求示例：**
```json
{
  "event_type": "daily_interaction",
  "positive_interaction": true,
  "bonding_moment": true
}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "pet_personality": {
      "playfulness": 76,
      "friendliness": 89,
      "loyalty": 93,
      "bond_strength": 62
    },
    "evolved_traits": ["bond_strength"],
    "message": "宠物与你的感情更深了！"
  }
}
```

---

## 5. 流程图

### 5.1 偏好学习流程

```
┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│   用户消息   │     │  对话引擎    │     │ 主人画像库   │
│             │────>│              │────>│              │
└──────────────┘     └──────────────┘     └──────────────┘
                                                  │
                                                  ▼
                                     ┌────────────────────────┐
                                     │  1. 称呼偏好提取        │
                                     │  - 检测名字/昵称出现    │
                                     │  - 更新preferred_name   │
                                     └───────────┬────────────┘
                                                 │
                                                 ▼
                                     ┌────────────────────────┐
                                     │  2. 时间模式分析        │
                                     │  - 记录交互时间点       │
                                     │  - 更新activity_pattern │
                                     └───────────┬────────────┘
                                                 │
                                                 ▼
                                     ┌────────────────────────┐
                                     │  3. 话题偏好提取        │
                                     │  - NLP话题分类         │
                                     │  - 更新preferred_topics│
                                     └───────────┬────────────┘
                                                 │
                                                 ▼
                                     ┌────────────────────────┐
                                     │  4. 情感分析            │
                                     │  - 计算sentiment_score │
                                     │  - 影响宠物情绪参数    │
                                     └───────────┬────────────┘
                                                 │
                                                 ▼
                                     ┌────────────────────────┐
                                     │  5. 性格养成            │
                                     │  - bond_strength++     │
                                     │  - loyalty适应变化      │
                                     └────────────────────────┘
```

### 5.2 个性化回复生成流程

```
┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│   用户消息   │     │  对话引擎    │     │ 主人画像库   │
│             │────>│              │────>│              │
└──────────────┘     └──────────────┘     └──────┬───────┘
                                                  │
                                                  ▼
                                     ┌────────────────────────┐
                                     │  查询用户画像           │
                                     │  - preferred_name     │
                                     │  - response_style      │
                                     │  - active_hours        │
                                     └───────────┬────────────┘
                                                 │
                    ┌────────────────────────────┼────────────────────────────┐
                    ▼                            ▼                            ▼
           ┌───────────────┐          ┌───────────────┐          ┌───────────────┐
           │称呼注入        │          │风格调整        │          │话题推荐        │
           │加入用户         │          │length/tone    │          │优先推荐        │
           │喜欢的称呼       │          │/emoji频率      │          │用户喜欢的话题  │
           └───────┬───────┘          └───────┬───────┘          └───────┬───────┘
                   │                          │                          │
                   └──────────────────────────┼──────────────────────────┘
                                              ▼
                                   ┌────────────────────────┐
                                   │  生成个性化回复         │
                                   │  "小明，今天累了吗？    │
                                   │   我给你唱首歌吧~"     │
                                   └────────────────────────┘
```

---

## 6. 模块联动

### 6.1 与对话引擎联动

- **触发时机：** 每次生成回复前
- **联动内容：**
  - 对话引擎调用 `GET /api/v1/owner-profile/{user_id}/reply-config`
  - 根据返回的个性化参数调整回复内容
- **数据流向：** 主人画像库 → 对话引擎

### 6.2 与行为引擎(PET_BEHAVIOR_ENGINE)联动

- **触发时机：** 生成动作序列时
- **联动内容：**
  - 行为引擎查询主人的动作风格偏好
  - 活泼型主人 → 动作幅度大、速度快
  - 安静型主人 → 动作轻柔、缓慢
- **数据流向：** 主人画像库 → 行为引擎

### 6.3 与记忆库(PET_MEMORY)联动

- **触发时机：** 对话结束后
- **联动内容：**
  - 记忆库记录交互历史
  - 画像库从历史中学习偏好
- **数据流向：** 记忆库 → 主人画像库

### 6.4 与控制台(OPENCLAW_CONSOLE)联动

- **触发时机：** 用户修改个人设置
- **联动内容：**
  - 控制台调用 `PUT /api/v1/owner-profile/{user_id}/preferences`
  - 实时更新画像数据
- **数据流向：** 控制台 → 主人画像库

---

## 7. 验收标准

### 7.1 功能验收

| 功能 | 验收条件 |
|------|----------|
| 称呼偏好学习 | 3次以上出现相同称呼后，识别准确率≥90% |
| 活跃时间分析 | 可正确识别用户的主要活跃时段 |
| 个性化回复 | 回复包含用户喜欢的称呼，风格与偏好一致 |
| 宠物性格养成 | 30天连续互动后，bond_strength提升≥10 |

### 7.2 性能验收

- 画像查询响应时间 ≤ 50ms
- 偏好更新延迟 ≤ 200ms
- 支持10万用户并发画像查询

### 7.3 数据质量验收

- 用户交互数据完整率 ≥ 99%
- 画像数据7天内完成初次学习

---

## 8. UI设计指引

### 8.1 个人设置界面

```
┌─────────────────────────────────────────────────────────────┐
│  个人设置                                                    │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌─────────┐                                                │
│  │  头像   │   昵称: [小明____________]                      │
│  │  [上传] │   生日: [1990-01-01____]                      │
│  └─────────┘   性别: ○男 ○女 ○保密                         │
│                                                             │
│  ─────────────────────────────────────────────────────────  │
│                                                             │
│  称呼偏好                                                   │
│  希望宠物怎么称呼你？                                        │
│  ○ 正式称呼（小明先生/小姐）                                 │
│  ● 昵称（小明）                                              │
│  ○ 亲昵称呼（宝贝/小可爱）                                   │
│                                                             │
│  ─────────────────────────────────────────────────────────  │
│                                                             │
│  宠物性格配置                                                │
│                                                             │
│  活泼度    [━━━━━━━━━━░░░░] 70%                            │
│  亲密度    [━━━━━━━━━━━━━░] 90%                            │
│  好奇心    [━━━━━━━━━░░░░░] 60%                            │
│                                                             │
│                                    [保存设置]                │
└─────────────────────────────────────────────────────────────┘
```

### 8.2 组件规范

- **称呼选择器：** 单选按钮组，展示不同称呼风格示例
- **性格滑块：** 0-100滑块，可实时预览宠物表现变化
- **统计卡片：** 显示交互次数/活跃时段/话题分布

### 8.3 配色方案

沿用控制台配色，增加以下扩展：

| 用途 | 颜色 | 色值 |
|------|------|------|
| 亲密指示 | 粉红渐变 | #FF6B9D → #FFB6C1 |
| 活跃时段 | 活力橙 | #FF9F43 |
| 话题标签 | 柔和紫 | #A29BFE |
