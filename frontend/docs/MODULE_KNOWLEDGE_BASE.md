# 知识库

**版本：** V1.0  
**模块负责人：** agentcp  
**编制日期：** 2026-03-20  

---

## 1. 概述

知识库是OpenClaw AI层的智能信息中枢，为宠物提供天气查询、新闻推送、常识问答和自定义知识管理能力。知识库通过集成外部API和本地知识库，使宠物能够回答主人关于生活信息的问题，成为真正有用的智能伙伴。

**业务目标：**
- 提供实时天气查询和穿衣/出行建议
- 支持每日新闻推送（可配置）
- 覆盖常见常识问答
- 支持用户自定义知识问答

---

## 2. 功能列表

| 功能 | 描述 | 优先级 | 触发方式 | 前端入口/按钮 |
|------|------|--------|----------|--------------|
| 天气查询 | 查询当前天气和预报，提供出行建议 | P0 | 自动+人工 | 宠物主动播报/用户询问 |
| 新闻推送 | 每日定时推送精选新闻 | P1 | 自动 | 「新闻设置」开关 |
| 常识问答 | 回答常见百科问题 | P0 | 自动 | 无按钮 |
| 自定义知识 | 用户添加自定义问答知识 | P1 | 人工 | 「知识管理」按钮 |
| 知识分类管理 | 管理知识条目分类结构 | P2 | 人工 | 管理后台 |
| 知识导入导出 | 批量导入导出知识数据 | P2 | 人工 | 管理后台 |
| 查询日志 | 记录所有知识查询用于分析 | P1 | 自动 | 无按钮 |

---

## 3. 数据模型

### 3.1 知识条目表 (knowledge_entries)

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| id | uint | 主键 | PK, auto_increment |
| entry_id | string | 知识唯一标识 | unique, not null |
| category_id | uint | 所属分类 | FK → knowledge_categories |
| question | text | 问题/触发语句 | not null, indexed |
| answer | text | 回答内容 | not null |
| answer_type | string | 回答类型 | text/weather/news/action |
| keywords | json | 关键词列表 | nullable |
| slot_entities | json | 槽位实体 | nullable |
| api_config | json | 关联API配置 | nullable |
| confidence_threshold | float | 置信度阈值 | default 0.6 |
| is_active | bool | 是否启用 | default true |
| is_system | bool | 是否系统知识 | default false |
| usage_count | int | 使用次数 | default 0 |
| success_count | int | 成功回答次数 | default 0 |
| user_id | uint | 创建用户（null=系统） | FK → sys_users, nullable |
| created_at | datetime | 创建时间 | auto |
| updated_at | datetime | 更新时间 | auto |

### 3.2 知识分类表 (knowledge_categories)

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| id | uint | 主键 | PK, auto_increment |
| category_id | string | 分类唯一标识 | unique, not null |
| category_name | string | 分类名称 | not null |
| parent_id | uint | 父分类ID | FK → knowledge_categories, nullable |
| description | string | 分类描述 | nullable |
| icon | string | 图标名称 | nullable |
| sort_order | int | 排序 | default 0 |
| is_active | bool | 是否启用 | default true |
| created_at | datetime | 创建时间 | auto |

### 3.3 查询日志表 (knowledge_queries)

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| id | uint | 主键 | PK, auto_increment |
| query_id | string | 查询唯一标识 | unique, not null |
| user_id | uint | 用户ID | FK → sys_users, nullable |
| device_id | string | 设备ID | FK → devices, nullable |
| query_text | text | 查询文本 | not null |
| intent | string | 识别的意图 | nullable |
| category_id | string | 匹配的知识分类 | nullable |
| entry_id | string | 匹配的知识条目 | nullable |
| api_source | string | 调用的外部API | nullable |
| response_text | text | 返回回答 | nullable |
| response_time_ms | int | 响应时间 | nullable |
| is_successful | bool | 是否成功 | default false |
| feedback | int | 用户反馈 | -1=差 0=无 1=好 |
| created_at | datetime | 查询时间 | auto |

---

## 4. 接口定义

### 4.1 天气查询

```
GET /api/v1/knowledge/weather
```

**参数：**

| 参数 | 类型 | 位置 | 必填 | 说明 |
|------|------|------|------|------|
| location | string | query | 否 | 城市名，默认自动获取用户位置 |
| device_id | string | query | 否 | 设备ID（用于定位） |

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "location": "上海",
    "weather": "晴",
    "temperature": 22,
    "temperature_range": "18-25°C",
    "humidity": 65,
    "wind": "东南风 3级",
    "aqi": 58,
    "aqi_level": "良",
    "suggestion": "今天天气不错，适合出门散步，记得带上水哦！",
    "forecast": [
      { "date": "今天", "weather": "晴", "temp": "18-25°C" },
      { "date": "明天", "weather": "多云", "temp": "17-24°C" },
      { "date": "后天", "weather": "小雨", "temp": "16-22°C" }
    ],
    "updated_at": "2026-03-20T10:00:00Z"
  }
}
```

### 4.2 新闻查询

```
GET /api/v1/knowledge/news
```

**参数：**

| 参数 | 类型 | 位置 | 必填 | 说明 |
|------|------|------|------|------|
| category | string | query | 否 | 新闻分类 hot/tech/sports/ent，默认hot |
| limit | int | query | 否 | 返回条数，默认5 |

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "category": "hot",
    "items": [
      { "title": "AI技术新突破", "source": "科技日报", "time": "2小时前", "summary": "..." },
      { "title": "周末天气适宜出行", "source": "气象局", "time": "4小时前", "summary": "..." }
    ],
    "updated_at": "2026-03-20T08:00:00Z"
  }
}
```

### 4.3 知识问答查询

```
POST /api/v1/knowledge/query
```

**参数：**

| 参数 | 类型 | 位置 | 必填 | 说明 |
|------|------|------|------|------|
| query | string | body | 是 | 查询文本 |
| user_id | uint | body | 否 | 用户ID |
| device_id | string | body | 否 | 设备ID |

**请求示例：**
```json
{
  "query": "今天上海适合出门吗？",
  "user_id": 10001,
  "device_id": "pet-001"
}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "query_id": "q-001",
    "answer": "今天上海天气晴朗，气温18-25度，空气质量良好，非常适合出门散步哦！",
    "answer_type": "weather",
    "confidence": 0.95,
    "related_entries": ["weather_query_001"]
  }
}
```

### 4.4 添加自定义知识

```
POST /api/v1/knowledge/entries
```

**请求示例：**
```json
{
  "category_id": "custom",
  "question": "你叫什么名字",
  "answer": "我叫小爪，是你的智能宠物！",
  "keywords": ["名字", "叫", "小爪"],
  "is_active": true
}
```

### 4.5 获取知识分类

```
GET /api/v1/knowledge/categories
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      { "category_id": "weather", "category_name": "天气", "entry_count": 12 },
      { "category_id": "news", "category_name": "新闻", "entry_count": 8 },
      { "category_id": "general", "category_name": "常识", "entry_count": 156 },
      { "category_id": "custom", "category_name": "自定义", "entry_count": 23 }
    ]
  }
}
```

### 4.6 查询日志统计

```
GET /api/v1/knowledge/query-stats
```

**参数：**

| 参数 | 类型 | 位置 | 必填 | 说明 |
|------|------|------|------|------|
| start_date | date | query | 否 | 开始日期 |
| end_date | date | query | 否 | 结束日期 |

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total_queries": 12580,
    "successful_queries": 11890,
    "success_rate": 94.5,
    "top_categories": [
      { "category": "weather", "count": 5230 },
      { "category": "general", "count": 4120 },
      { "category": "news", "count": 3230 }
    ],
    "avg_response_time_ms": 120
  }
}
```

---

## 5. 流程图

### 5.1 知识查询处理流程

```
┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│   用户消息   │     │  对话引擎    │     │   知识库     │
│             │────>│              │────>│              │
└──────────────┘     └──────────────┘     └──────┬───────┘
                                                  │
                                                  ▼
                                     ┌────────────────────────┐
                                     │  意图分类匹配           │
                                     │  weather/news/general  │
                                     └───────────┬────────────┘
                                                 │
                          ┌──────────────────────┼──────────────────────┐
                          ▼                      ▼                      ▼
                 ┌───────────────┐        ┌───────────────┐        ┌───────────────┐
                 │  天气查询     │        │  新闻查询     │        │  常识问答     │
                 │  调用天气API  │        │  调用新闻API  │        │  匹配知识库   │
                 └───────┬───────┘        └───────┬───────┘        └───────┬───────┘
                         │                        │                        │
                         └─────────────────────────┼────────────────────────┘
                                                   ▼
                                      ┌────────────────────────┐
                                      │  答案组装               │
                                      │  - 格式化回复           │
                                      │  - 注入宠物风格         │
                                      │  - 添加建议             │
                                      └───────────┬────────────┘
                                                  │
                                                  ▼
                                     ┌────────────────────────┐
                                     │  记录查询日志           │
                                     │  - 保存query_text       │
                                     │  - 保存response         │
                                     │  - 统计usage_count      │
                                     └────────────────────────┘
```

### 5.2 宠物主动推送流程

```
┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│   定时触发   │     │  知识库      │     │ 主人画像库   │
│  (每日早晚)  │────>│              │────>│              │
└──────────────┘     └──────────────┘     └──────┬───────┘
                                                  │
                                                  ▼
                                     ┌────────────────────────┐
                                     │  判断推送内容           │
                                     │  - 天气主动播报         │
                                     │  - 新闻精选推送         │
                                     │  - 纪念日提醒           │
                                     └───────────┬────────────┘
                                                 │
                                                 ▼
                                      ┌────────────────────────┐
                                      │  生成推送消息           │
                                      │  "主人早上好！今天是    │
                                      │   晴天，适合散步哦~"    │
                                      └───────────┬────────────┘
                                                 │
                                                 ▼
                                      ┌────────────────────────┐
                                      │  通过MQTT下发           │
                                      │  /miniclaw/{device_id} │
                                      │    /down/speech        │
                                      └────────────────────────┘
```

---

## 6. 模块联动

### 6.1 与对话引擎联动

- **触发时机：** 用户询问天气/新闻/常识时
- **联动内容：**
  - 对话引擎调用 `POST /api/v1/knowledge/query`
  - 知识库返回结构化答案
  - 对话引擎将答案转化为自然语言
- **数据流向：** 对话引擎 → 知识库 → 对话引擎

### 6.2 与主人画像库(OWNER_PROFILE)联动

- **触发时机：** 生成主动推送时
- **联动内容：**
  - 知识库查询主人的活跃时段
  - 根据偏好选择推送内容类型
- **数据流向：** 主人画像库 → 知识库

### 6.3 与行为引擎(PET_BEHAVIOR_ENGINE)联动

- **触发时机：** 需要边说边动时
- **联动内容：**
  - 知识库返回的回答类型为action时
  - 触发行为引擎生成配套动作
- **数据流向：** 知识库 → 行为引擎

### 6.4 与MiniClaw通信协议(MINICLAW_PROTOCOL)联动

- **触发时机：** 主动推送和语音播报
- **联动内容：**
  - 知识库调用MQTT的 `/down/speech` 下发语音
  - 接收设备状态确认
- **数据流向：** 知识库 → MQTT Broker → MiniClaw

---

## 7. 验收标准

### 7.1 功能验收

| 功能 | 验收条件 |
|------|----------|
| 天气查询 | 支持全国主要城市，返回完整天气信息 |
| 新闻推送 | 每日早8点、晚6点可配置推送 |
| 常识问答 | 常见问题回答准确率≥85% |
| 自定义知识 | 用户可添加问题-答案对，实时生效 |

### 7.2 性能验收

- 天气API响应时间 ≤ 500ms
- 新闻API响应时间 ≤ 1s
- 知识库匹配响应时间 ≤ 100ms
- 支持1000次/分钟查询

### 7.3 数据质量验收

- 天气数据更新频率 ≥ 每小时1次
- 新闻数据更新频率 ≥ 每日4次
- 知识库条目准确率验证 ≥ 95%

---

## 8. UI设计指引

### 8.1 知识管理界面

```
┌─────────────────────────────────────────────────────────────┐
│  知识管理                                                    │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│ 分类: [全部 ▼]  搜索: [________________________] [搜索]     │
│                                                             │
│ ┌─────────────────────────────────────────────────────────┐ │
│ │ ☀ 天气知识                    共12条    [添加] [导入]    │ │
│ ├─────────────────────────────────────────────────────────┤ │
│ │ 问: 上海天气怎么样？                                      │ │
│ │ 答: 今天上海{weather}，气温{temp}...                    │ │
│ │ 关键词: 上海, 天气     使用: 523次  ✓启用  [编辑] [删除] │ │
│ ├─────────────────────────────────────────────────────────┤ │
│ │ 问: 明天要带伞吗？                                        │ │
│ │ 答: 明天{weather}，建议{umbrella_suggestion}            │ │
│ │ 关键词: 伞, 明天      使用: 312次  ✓启用  [编辑] [删除]   │ │
│ └─────────────────────────────────────────────────────────┘ │
│                                                             │
│ ┌─────────────────────────────────────────────────────────┐ │
│ │ 📚 自定义知识                    共23条                  │ │
│ │ ...                                                      │ │
│ └─────────────────────────────────────────────────────────┘ │
│                                                             │
│                         页码: < 1 2 3 4 5 >                  │
└─────────────────────────────────────────────────────────────┘
```

### 8.2 组件规范

- **知识卡片：** 问答预览，支持展开全文
- **分类树：** 左侧知识分类导航
- **API配置面板：** 配置外部API的key和参数

### 8.3 配色方案

| 用途 | 颜色 | 色值 |
|------|------|------|
| 天气 | 阳光黄 | #FFD93D |
| 新闻 | 资讯蓝 | #74B9FF |
| 常识 | 知识绿 | #6BCB77 |
| 自定义 | 创意紫 | #A29BFE |
