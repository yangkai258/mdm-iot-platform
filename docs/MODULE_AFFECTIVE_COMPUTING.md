# 模块 PRD：情感计算（MODULE_AFFECTIVE_COMPUTING）

**版本：** V1.0
**所属Phase：** Phase 3（Sprint 17-18）
**优先级：** P1
**负责角色：** agentcp（产品）、agenthd（后端）、agentqd（前端）

---

## 一、概述

### 1.1 模块定位

情感计算模块是宠物情感陪伴的核心能力，通过多模态感知（语音/文字/表情/行为）识别用户和宠物的情绪状态，生成情感响应策略，构建家庭情绪地图，为用户提供有温度的情感陪伴体验。

### 1.2 核心价值

- **情感共鸣**：宠物能感知用户情绪并做出恰当回应
- **心理健康**：长期情绪追踪，提供心理健康建议
- **家庭和谐**：感知家庭整体情绪状态，促进家庭和谐

### 1.3 范围边界

**包含：**
- 用户语音情绪识别
- 用户文字情绪识别
- 宠物表情/行为情绪识别
- 情绪响应策略
- 情绪日志与趋势
- 家庭情绪地图

**不包含：**
- 具身智能动作执行（MODULE_EMBODIED_AI）
- 数字孪生体征数据（MODULE_DIGITAL_TWIN）
- 前端表情动效渲染（UI层）

---

## 二、功能详情

### 2.1 情绪识别

#### 2.1.1 用户语音情绪识别

| 功能 | 说明 |
|------|------|
| 实时语音情绪 | 用户说话时实时分析情绪（愉快/悲伤/愤怒/惊讶/平静） |
| 情绪强度 | 分析情绪的强烈程度（0-100%） |
| 情绪转折点 | 识别对话中情绪变化时刻 |
| 多人对话 | 区分不同说话人的情绪 |
| 情绪原因推断 | 结合上下文推断情绪产生原因 |

#### 2.1.2 用户文字情绪识别

| 功能 | 说明 |
|------|------|
| 文本情绪分析 | 用户输入文字的情绪分析 |
| 表情符号解读 | 识别emoji表达的情绪 |
| 语义情绪 | 理解反语/讽刺等复杂语义 |
| 情绪关键词 | 提取情绪关键词 |
| 紧急情绪检测 | 识别用户表达自杀/自伤等紧急情绪 |

#### 2.1.3 宠物表情情绪识别

| 功能 | 说明 |
|------|------|
| 面部表情 | 识别宠物面部表情（开心/悲伤/害怕/好奇/疲惫） |
| 行为情绪 | 通过行为模式判断情绪（蹭/叫/躲/摇尾巴） |
| 声音情绪 | 通过叫声判断情绪（高兴/不满/害怕/求助） |
| 情绪组合 | 综合面部+行为+声音判断复合情绪 |
| 情绪置信度 | 给出情绪判断的置信度 |

#### 2.1.4 情绪评估指标

| 指标 | 说明 | 准确率要求 |
|------|------|------------|
| 语音情绪准确率 | 语音情绪分类准确率 | >85% |
| 文字情绪准确率 | 文字情绪分类准确率 | >88% |
| 宠物表情准确率 | 宠物表情分类准确率 | >80% |
| 情绪强度误差 | 预测强度与实际偏差 | <15% |

### 2.2 情绪响应策略

#### 2.2.1 情绪响应类型

| 用户情绪 | 宠物响应策略 | 说明 |
|----------|--------------|------|
| 愉快 | 共鸣回应 | 宠物也表现出开心的动作 |
| 悲伤 | 安慰陪伴 | 宠物主动靠近、舔舐、发出安抚声 |
| 愤怒 | 冷静等待 | 宠物保持安静，不激化情绪 |
| 惊讶 | 好奇回应 | 宠物表现出好奇/警觉 |
| 平静 | 正常陪伴 | 宠物保持正常活动 |
| 焦虑 | 安抚互动 | 宠物通过舔舐/依偎缓解焦虑 |
| 孤独 | 主动求互动 | 宠物主动找用户玩 |
| 疲劳 | 安静陪伴 | 宠物保持安静不打扰 |

#### 2.2.2 响应动作库

| 动作类别 | 具体动作 | 触发条件 |
|----------|----------|----------|
| 肢体动作 | 蹭、舔、依偎、靠近、躲开 | 情绪类型 |
| 声音动作 | 叫声、哼唧、呜咽、欢叫 | 情绪类型+强度 |
| 表情动作 | 眼神变化、耳朵变化、尾巴变化 | 情绪类型 |
| 行为动作 | 叼玩具、送东西、蹭腿 | 情绪类型+用户偏好 |

#### 2.2.3 响应配置

| 功能 | 说明 |
|------|------|
| 响应灵敏度 | 用户可调整宠物响应灵敏度 |
| 响应延迟 | 用户可设置宠物响应延迟时间 |
| 响应频率 | 限制单位时间内响应次数 |
| 禁用某些响应 | 用户可禁用特定类型的响应 |
| 自定义响应 | 用户可自定义特定情绪的响应动作 |

### 2.3 情绪日志

#### 2.3.1 日志记录

| 功能 | 说明 |
|------|------|
| 自动记录 | 用户和宠物情绪自动记录到日志 |
| 情绪标签 | 每个情绪记录可打标签 |
| 情绪备注 | 用户可为情绪记录添加备注 |
| 事件关联 | 情绪记录关联触发事件 |

#### 2.3.2 日志查询

| 功能 | 说明 |
|------|------|
| 时间筛选 | 按日/周/月/自定义时间查询 |
| 类型筛选 | 按情绪类型/来源筛选 |
| 关键词搜索 | 通过备注关键词搜索 |
| 情绪对比 | 对比不同时间段的情绪变化 |

#### 2.3.3 情绪报告

| 功能 | 说明 |
|------|------|
| 日报 | 每日情绪总结 |
| 周报 | 每周情绪趋势分析 |
| 月报 | 每月情绪健康报告 |
| 异常报告 | 情绪异常时的专门报告 |

### 2.4 家庭情绪地图

#### 2.4.1 家庭情绪感知

| 功能 | 说明 |
|------|------|
| 家庭成员情绪 | 感知每个家庭成员的情绪状态 |
| 家庭整体情绪 | 计算家庭整体情绪指数 |
| 情绪传染分析 | 分析情绪在家庭成员间的传染 |
| 关键影响者 | 识别对家庭情绪影响最大的成员 |

#### 2.4.2 家庭情绪图谱

| 功能 | 说明 |
|------|------|
| 情绪关系图 | 可视化家庭成员情绪关系 |
| 情绪热力图 | 展示家庭情绪分布 |
| 情绪日历 | 家庭情绪日历视图 |
| 重要时刻 | 标记家庭情绪的重要时刻 |

#### 2.4.3 家庭情绪建议

| 功能 | 说明 |
|------|------|
| 氛围调节建议 | 提供改善家庭氛围的建议 |
| 成员互动建议 | 促进家庭成员互动的建议 |
| 冲突预警 | 预警可能出现的情绪冲突 |
| 庆祝提醒 | 提醒家庭情绪高涨适合庆祝 |

---

## 三、API接口定义

### 3.1 情绪识别

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/v1/emotion/recognize/voice | 语音情绪识别 |
| POST | /api/v1/emotion/recognize/text | 文字情绪识别 |
| POST | /api/v1/emotion/recognize/pet | 宠物情绪识别 |
| POST | /api/v1/emotion/recognize/batch | 批量情绪识别 |

### 3.2 情绪响应

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/emotion/:pet_id/response-strategies | 获取响应策略 |
| POST | /api/v1/emotion/:pet_id/response | 触发情绪响应 |
| GET | /api/v1/emotion/action-library | 动作库列表 |
| POST | /api/v1/emotion/action-library | 添加自定义动作 |
| PUT | /api/v1/emotion/:pet_id/response-config | 更新响应配置 |

### 3.3 情绪日志

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/emotion/logs | 情绪日志列表 |
| GET | /api/v1/emotion/logs/:id | 日志详情 |
| POST | /api/v1/emotion/logs | 创建日志（含备注） |
| PUT | /api/v1/emotion/logs/:id | 更新日志 |
| DELETE | /api/v1/emotion/logs/:id | 删除日志 |
| GET | /api/v1/emotion/:pet_id/reports/daily | 每日情绪报告 |
| GET | /api/v1/emotion/:pet_id/reports/weekly | 每周情绪报告 |
| GET | /api/v1/emotion/:pet_id/reports/monthly | 每月情绪报告 |

### 3.4 家庭情绪地图

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/emotion/family/:family_id/map | 家庭情绪地图 |
| GET | /api/v1/emotion/family/:family_id/trends | 家庭情绪趋势 |
| GET | /api/v1/emotion/family/:family_id/suggestions | 家庭情绪建议 |
| GET | /api/v1/emotion/family/:family_id/calendar | 家庭情绪日历 |

---

## 四、数据库设计

### 4.1 情绪记录表 (emotion_records)

```sql
CREATE TABLE emotion_records (
    id              BIGSERIAL PRIMARY KEY,
    subject_type    VARCHAR(20) NOT NULL,         -- 'user'/'pet'
    subject_id      BIGINT NOT NULL,               -- user_id 或 pet_id
    emotion_type    VARCHAR(30) NOT NULL,         -- 'happy'/'sad'/'angry'/'surprised'/'calm'/'anxious'/'lonely'/'tired'
    intensity       DECIMAL(5,2) NOT NULL,        -- 情绪强度 0-100
    source          VARCHAR(20) NOT NULL,         -- 'voice'/'text'/'face'/'behavior'
    confidence      DECIMAL(5,4),
    context         JSONB,                          -- 上下文信息
    trigger_event   TEXT,
    tags            VARCHAR(100)[],
    note            TEXT,
    recorded_at     TIMESTAMP NOT NULL,
    created_at      TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_emotion_records_subject ON emotion_records(subject_type, subject_id, recorded_at DESC);
CREATE INDEX idx_emotion_records_type ON emotion_records(emotion_type, recorded_at DESC);
```

### 4.2 宠物情绪动作表 (pet_emotion_actions)

```sql
CREATE TABLE pet_emotion_actions (
    id              BIGSERIAL PRIMARY KEY,
    pet_id          BIGINT NOT NULL REFERENCES pets(id),
    emotion_type    VARCHAR(30) NOT NULL,
    action_type     VARCHAR(30) NOT NULL,         -- 'body'/'sound'/'expression'/'behavior'
    action_name     VARCHAR(100) NOT NULL,
    action_data     JSONB,                          -- 动作具体参数
    priority        INT DEFAULT 0,
    min_intensity   DECIMAL(5,2) DEFAULT 0,
    max_intensity   DECIMAL(5,2) DEFAULT 100,
    is_custom       BOOLEAN DEFAULT FALSE,
    created_by      BIGINT REFERENCES users(id),
    created_at      TIMESTAMP DEFAULT NOW()
);
```

### 4.3 情绪响应配置表 (emotion_response_configs)

```sql
CREATE TABLE emotion_response_configs (
    id              BIGSERIAL PRIMARY KEY,
    pet_id          BIGINT NOT NULL REFERENCES pets(id),
    sensitivity     DECIMAL(3,2) DEFAULT 0.5,      -- 灵敏度 0-1
    response_delay_ms INT DEFAULT 0,
    response_freq_limit INT DEFAULT 10,           -- 每小时最大响应次数
    disabled_actions VARCHAR(50)[],              -- 禁用的动作类型
    enabled         BOOLEAN DEFAULT TRUE,
    updated_at      TIMESTAMP DEFAULT NOW()
);
```

### 4.4 家庭情绪表 (family_emotions)

```sql
CREATE TABLE family_emotions (
    id              BIGSERIAL PRIMARY KEY,
    family_id       BIGINT NOT NULL REFERENCES families(id),
    recorded_at     TIMESTAMP NOT NULL,
    overall_index   DECIMAL(5,2),                 -- 整体情绪指数 0-100
    dominant_emotion VARCHAR(30),
    member_states   JSONB,                          -- 各成员情绪状态
    influence_map   JSONB,                        -- 情绪影响关系
    created_at      TIMESTAMP DEFAULT NOW(),
    UNIQUE(family_id, recorded_at)
);

CREATE INDEX idx_family_emotions_family ON family_emotions(family_id, recorded_at DESC);
```

### 4.5 情绪报告表 (emotion_reports)

```sql
CREATE TABLE emotion_reports (
    id              BIGSERIAL PRIMARY KEY,
    subject_type    VARCHAR(20) NOT NULL,
    subject_id      BIGINT NOT NULL,
    report_type     VARCHAR(20) NOT NULL,         -- 'daily'/'weekly'/'monthly'/'alert'
    period_start    TIMESTAMP NOT NULL,
    period_end      TIMESTAMP NOT NULL,
    summary         JSONB,                          -- 报告摘要
    emotion_trends  JSONB,                          -- 情绪趋势数据
    anomaly_events  JSONB,                          -- 异常事件
    suggestions     TEXT,
    generated_at    TIMESTAMP DEFAULT NOW()
);
```

---

## 五、前端页面清单

### 5.1 情绪识别

| 页面 | 路由 | 说明 |
|------|------|------|
| 情绪识别面板 | /emotion/recognize | 多模态情绪识别入口 |
| 语音情绪 | /emotion/recognize/voice | 实时语音情绪识别 |
| 文字情绪 | /emotion/recognize/text | 文字输入情绪分析 |

### 5.2 情绪响应

| 页面 | 路由 | 说明 |
|------|------|------|
| 响应策略配置 | /emotion/:pet_id/response-config | 响应灵敏度/频率配置 |
| 动作库 | /emotion/action-library | 动作库管理 |
| 自定义动作 | /emotion/action-library/create | 添加自定义动作 |

### 5.3 情绪日志

| 页面 | 路由 | 说明 |
|------|------|------|
| 情绪日志 | /emotion/logs | 情绪历史记录 |
| 日志详情 | /emotion/logs/:id | 单条记录详情 |
| 情绪报告 | /emotion/reports | 日报/周报/月报 |

### 5.4 家庭情绪

| 页面 | 路由 | 说明 |
|------|------|------|
| 家庭情绪地图 | /emotion/family/:family_id | 家庭情绪可视化 |
| 情绪趋势 | /emotion/family/:family_id/trends | 家庭情绪趋势 |
| 情绪建议 | /emotion/family/:family_id/suggestions | 改善建议 |

---

## 六、验收标准

### 6.1 情绪识别

| 验收点 | 标准 |
|--------|------|
| 语音情绪准确率 | >85%（与人工标注对比） |
| 文字情绪准确率 | >88% |
| 宠物表情准确率 | >80% |
| 识别延迟 | 语音<500ms，文字<200ms |
| 紧急情绪检测 | 自杀/自伤内容100%预警 |

### 6.2 情绪响应

| 验收点 | 标准 |
|--------|------|
| 响应触发 | 情绪识别后5秒内触发响应 |
| 响应准确率 | 响应动作与情绪匹配度>80% |
| 响应频率 | 不超过配置的频率限制 |
| 禁用生效 | 禁用的动作类型不触发 |

### 6.3 情绪日志

| 验收点 | 标准 |
|--------|------|
| 自动记录 | 所有情绪变化自动记录 |
| 查询性能 | 1000条记录查询<1秒 |
| 报告生成 | 日报/周报/月报按时自动生成 |

### 6.4 家庭情绪地图

| 验收点 | 标准 |
|--------|------|
| 整体情绪指数 | 正确反映家庭情绪状态 |
| 情绪建议采纳率 | >30%（用户采纳建议比例） |
| 冲突预警准确率 | >70% |
