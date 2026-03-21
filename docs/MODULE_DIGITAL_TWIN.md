# 模块 PRD：数字孪生（MODULE_DIGITAL_TWIN）

**版本：** V1.0
**所属Phase：** Phase 3（Sprint 19-20）
**优先级：** P1
**负责角色：** agentcp（产品）、agenthd（后端）、agentqd（前端）

---

## 一、概述

### 1.1 模块定位

数字孪生模块是宠物生命体的虚拟数字化映射，通过实时同步宠物设备的状态数据，构建宠物的虚拟生命体征、行为模式、历史记忆，支持跨设备状态同步和离线历史回放，为宠物健康管理提供数据基础。

### 1.2 核心价值

- **健康守护**：实时生命体征监测，异常提前预警
- **记忆永续**：宠物一生记录永不丢失
- **无缝体验**：切换设备不丢失任何状态

### 1.3 范围边界

**包含：**
- 实时生命体征数字孪生
- 行为预测
- 健康预警
- 历史回放与精彩瞬间
- 跨设备状态同步
- 离线支持

**不包含：**
- 情感计算（MODULE_AFFECTIVE_COMPUTING）
- 具身智能（MODULE_EMBODIED_AI）
- 前端3D渲染引擎（单独模块）

---

## 二、功能详情

### 2.1 实时生命体征数字孪生

#### 2.1.1 生命体征类型

| 指标 | 数据来源 | 更新频率 | 精度要求 |
|------|----------|----------|----------|
| 心跳 | 设备传感器 | 实时 | ±5 BPM |
| 呼吸频率 | 设备传感器 | 实时 | ±2 次/分钟 |
| 体温估算 | 行为/环境综合 | 1分钟 | ±0.5°C |
| 活动量 | 加速计 | 实时 | 95% |
| 睡眠状态 | 活动量推断 | 5分钟 | 90% |
| 情绪指数 | 多维度综合 | 1分钟 | 85% |

#### 2.1.2 数字心跳

| 功能 | 说明 |
|------|------|
| 实时心跳曲线 | 展示最近5分钟/1小时/24小时心电图曲线 |
| 心率变异性（HRV） | 分析心率间隔变化，评估压力水平 |
| 异常心率标记 | 心动过速/过缓自动标记 |
| 对比基准 | 与同品种/同年龄段基准对比 |

#### 2.1.3 数字呼吸

| 功能 | 说明 |
|------|------|
| 呼吸频率曲线 | 展示呼吸频率变化 |
| 呼吸深度分析 | 评估呼吸质量 |
| 呼吸暂停检测 | 检测睡眠呼吸暂停 |
| 呼吸模式识别 | 区分安静/活动/紧张等模式 |

#### 2.1.4 体温估算

| 功能 | 说明 |
|------|------|
| 体温曲线 | 基于行为和环境估算体温变化 |
| 发烧预警 | 体温超过阈值告警 |
| 体温分布图 | 全天候体温热力图 |

### 2.2 行为预测

#### 2.2.1 预测类型

| 预测类型 | 预测内容 | 准确率要求 | 应用场景 |
|----------|----------|------------|----------|
| 短期动作 | 接下来5分钟内可能的动作 | >75% | 实时互动 |
| 活动周期 | 即将进入活跃/休息状态 | >80% | 作息管理 |
| 意图识别 | 饿了/渴了/想玩/想出门 | >70% | 需求预判 |
| 行为异常 | 与日常模式显著偏离 | >85% | 健康预警 |

#### 2.2.2 预测展示

| 功能 | 说明 |
|------|------|
| 动作预测卡片 | 展示预测的动作 + 置信度 |
| 意图气泡 | App内展示宠物当前可能的意图 |
| 预测更新 | 实时更新预测结果 |

### 2.3 健康预警

#### 2.3.1 预警类型

| 预警类型 | 触发条件 | 预警级别 | 通知方式 |
|----------|----------|----------|----------|
| 发烧预警 | 体温>39.5°C持续5分钟 | 紧急 | 推送+电话 |
| 心率异常 | HR超出正常范围 | 警告 | 推送+短信 |
| 呼吸异常 | 呼吸暂停>15秒 | 紧急 | 推送+电话 |
| 行为异常 | 与基线偏差>30% | 提示 | 推送 |
| 饮水异常 | 饮水量突降>50% | 提示 | 推送 |
| 运动量异常 | 日运动量<正常30% | 提示 | 推送 |

#### 2.3.2 预警处理

| 功能 | 说明 |
|------|------|
| 预警确认 | 用户确认收到预警 |
| 预警忽略 | 用户可选择忽略并填写原因 |
| 预警建议 | 附带健康建议或就医指导 |
| 就医记录 | 预警关联就医记录 |

### 2.4 历史回放

#### 2.4.1 记录类型

| 记录类型 | 记录内容 | 保存期限 |
|----------|----------|----------|
| 生命体征 | 心跳/呼吸/活动量原始数据 | 永久 |
| 行为事件 | 吃饭/睡觉/玩耍等事件 | 永久 |
| 位置轨迹 | 室内定位轨迹（可选开启） | 1年 |
| 精彩瞬间 | AI筛选的高光时刻 | 永久 |
| 健康档案 | 体检/疫苗/就医记录 | 永久 |
| 成长记录 | 体重/身高/外观变化 | 永久 |

#### 2.4.2 历史回放功能

| 功能 | 说明 |
|------|------|
| 时间轴浏览 | 按日/周/月/年浏览事件 |
| 事件筛选 | 筛选特定类型事件 |
| 事件详情 | 点击查看事件详情 |
| 数据导出 | 导出为PDF/JSON格式 |
| 回忆生成 | AI自动生成"今天/这周/这月回忆" |

### 2.5 精彩瞬间AI筛选

#### 2.5.1 筛选条件

| 条件类型 | 说明 | 权重 |
|----------|------|------|
| 情绪峰值 | 宠物情绪达到高峰的瞬间 | 高 |
| 互动高峰 | 人宠互动最活跃的时刻 | 高 |
| 稀有动作 | 宠物做出稀有动作 | 中 |
| 家庭成员同框 | 宠物与主人同时出现 | 中 |
| 场景丰富 | 背景环境有特点 | 低 |

#### 2.5.2 精彩瞬间管理

| 功能 | 说明 |
|------|------|
| 自动生成 | 每日/每周自动生成精彩瞬间集 |
| 手动标记 | 用户可手动标记某个瞬间 |
| 相册导出 | 导出到手机相册 |
| 分享功能 | 分享到社交平台 |
| 精彩瞬间统计 | 每月精彩瞬间数量统计 |

### 2.6 跨设备状态同步

#### 2.6.1 同步内容

| 同步类型 | 说明 | 同步频率 |
|----------|------|----------|
| 实时状态 | 当前行为/情绪/生命体征 | 实时 |
| 历史数据 | 最近7天行为记录 | 每日同步 |
| 配置状态 | 宠物profile/性格/偏好 | 变更时同步 |
| 情感记忆 | 最近互动记忆 | 实时 |

#### 2.6.2 同步机制

| 功能 | 说明 |
|------|------|
| 云端同步 | 设备状态实时同步到云端 |
| 多设备订阅 | 同一宠物多个设备共享状态 |
| 冲突解决 | 多设备同时修改时以最新时间戳为准 |
| 同步状态 | 展示各设备同步状态 |

### 2.7 离线支持

#### 2.7.1 本地缓存

| 功能 | 说明 |
|------|------|
| 基础数据缓存 | 宠物profile/性格/偏好本地缓存 |
| 最近记录 | 最近100条行为记录本地缓存 |
| 模型缓存 | AI模型本地缓存（支持离线推理） |
| 同步队列 | 离线期间的变更记录到队列 |

#### 2.7.2 断网续传

| 功能 | 说明 |
|------|------|
| 离线记录 | 离线期间的行为记录本地存储 |
| 自动重连 | 网络恢复后自动上传离线数据 |
| 数据完整性 | 上传后校验数据完整性 |
| 状态恢复 | 上传完成后恢复完整状态 |

---

## 三、API接口定义

### 3.1 生命体征

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/digital-twin/:pet_id/vitals | 获取实时生命体征 |
| GET | /api/v1/digital-twin/:pet_id/vitals/history | 生命体征历史 |
| GET | /api/v1/digital-twin/:pet_id/vitals/heartbeat | 心跳曲线 |
| GET | /api/v1/digital-twin/:pet_id/vitals/respiration | 呼吸曲线 |
| GET | /api/v1/digital-twin/:pet_id/vitals/temperature | 体温曲线 |
| POST | /api/v1/digital-twin/:pet_id/vitals/report | 设备上报体征数据 |

### 3.2 行为预测

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/digital-twin/:pet_id/predictions | 获取行为预测 |
| GET | /api/v1/digital-twin/:pet_id/predictions/short-term | 短期动作预测 |
| GET | /api/v1/digital-twin/:pet_id/predictions/intent | 意图识别结果 |

### 3.3 健康预警

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/digital-twin/:pet_id/alerts | 健康预警列表 |
| POST | /api/v1/digital-twin/:pet_id/alerts/:id/ack | 确认预警 |
| POST | /api/v1/digital-twin/:pet_id/alerts/:id/ignore | 忽略预警 |
| GET | /api/v1/digital-twin/:pet_id/health-report | 健康报告 |

### 3.4 历史回放

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/digital-twin/:pet_id/timeline | 时间轴事件 |
| GET | /api/v1/digital-twin/:pet_id/events/:event_id | 事件详情 |
| GET | /api/v1/digital-twin/:pet_id/memories | 回忆集 |
| GET | /api/v1/digital-twin/:pet_id/memories/:id/highlights | 精彩瞬间 |
| POST | /api/v1/digital-twin/:pet_id/export | 导出数据 |

### 3.5 跨设备同步

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/digital-twin/:pet_id/sync/status | 同步状态 |
| POST | /api/v1/digital-twin/:pet_id/sync/pull | 拉取最新状态 |
| POST | /api/v1/digital-twin/:pet_id/sync/push | 推送本地变更 |
| POST | /api/v1/digital-twin/:pet_id/sync/batch | 批量同步离线数据 |

---

## 四、数据库设计

### 4.1 宠物数字孪生表 (digital_twin_pets)

```sql
CREATE TABLE digital_twin_pets (
    id              BIGSERIAL PRIMARY KEY,
    pet_id          BIGINT NOT NULL REFERENCES pets(id),
    current_state   JSONB,                          -- 当前状态快照
    last_sync_at    TIMESTAMP,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW(),
    UNIQUE(pet_id)
);
```

### 4.2 生命体征记录表 (vital_records)

```sql
CREATE TABLE vital_records (
    id              BIGSERIAL PRIMARY KEY,
    pet_id          BIGINT NOT NULL REFERENCES pets(id),
    vital_type      VARCHAR(30) NOT NULL,         -- 'heartbeat'/'respiration'/'temperature'/'activity'
    value           JSONB NOT NULL,               -- 根据类型不同存储不同数据
    recorded_at      TIMESTAMP NOT NULL,
    device_id       VARCHAR(100),
    created_at      TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_vital_records_pet_time ON vital_records(pet_id, vital_type, recorded_at DESC);
CREATE INDEX idx_vital_records_device ON vital_records(device_id, recorded_at DESC);
```

### 4.3 行为事件表 (behavior_events)

```sql
CREATE TABLE behavior_events (
    id              BIGSERIAL PRIMARY KEY,
    pet_id          BIGINT NOT NULL REFERENCES pets(id),
    event_type      VARCHAR(50) NOT NULL,         -- 'eating'/'sleeping'/'playing'/'walking'/'other'
    event_data      JSONB,                          -- 事件详情
    start_time      TIMESTAMP NOT NULL,
    end_time       TIMESTAMP,
    duration_seconds INT,
    confidence      DECIMAL(5,4),
    device_id       VARCHAR(100),
    created_at      TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_behavior_events_pet_time ON behavior_events(pet_id, start_time DESC);
```

### 4.4 健康预警表 (health_alerts)

```sql
CREATE TABLE health_alerts (
    id              BIGSERIAL PRIMARY KEY,
    pet_id          BIGINT NOT NULL REFERENCES pets(id),
    alert_type      VARCHAR(50) NOT NULL,         -- 'fever'/'heart_rate'/'breathing'/'behavior'
    severity        VARCHAR(20) NOT NULL,          -- 'info'/'warning'/'critical'
    trigger_value   VARCHAR(100),
    threshold       VARCHAR(100),
    status          VARCHAR(20) DEFAULT 'pending', -- 'pending'/'acknowledged'/'ignored'/'resolved'
    suggestion      TEXT,
    acknowledged_at TIMESTAMP,
    acknowledged_by BIGINT REFERENCES users(id),
    ignore_reason   TEXT,
    resolved_at     TIMESTAMP,
    created_at      TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_health_alerts_pet ON health_alerts(pet_id, created_at DESC);
CREATE INDEX idx_health_alerts_status ON health_alerts(status, created_at DESC);
```

### 4.5 精彩瞬间表 (highlight_moments)

```sql
CREATE TABLE highlight_moments (
    id              BIGSERIAL PRIMARY KEY,
    pet_id          BIGINT NOT NULL REFERENCES pets(id),
    moment_type     VARCHAR(30) NOT NULL,         -- 'emotion_peak'/'interaction'/'rare_action'/'family_moment'
    media_url       VARCHAR(500),
    thumbnail_url   VARCHAR(500),
    description     TEXT,
    emotion_score   DECIMAL(5,4),
    captured_at     TIMESTAMP NOT NULL,
    is_auto_generated BOOLEAN DEFAULT TRUE,
    is_exported     BOOLEAN DEFAULT FALSE,
    created_at      TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_highlight_moments_pet ON highlight_moments(pet_id, captured_at DESC);
```

### 4.6 同步记录表 (sync_records)

```sql
CREATE TABLE sync_records (
    id              BIGSERIAL PRIMARY KEY,
    pet_id          BIGINT NOT NULL REFERENCES pets(id),
    device_id       VARCHAR(100) NOT NULL,
    sync_type       VARCHAR(20) NOT NULL,         -- 'push'/'pull'/'batch'
    sync_data       JSONB,
    sync_status     VARCHAR(20) NOT NULL,          -- 'success'/'failed'/'partial'
    records_synced  INT DEFAULT 0,
    error_message   TEXT,
    synced_at       TIMESTAMP,
    created_at      TIMESTAMP DEFAULT NOW()
);
```

---

## 五、前端页面清单

### 5.1 数字孪生主页

| 页面 | 路由 | 说明 |
|------|------|------|
| 数字孪生首页 | /digital-twin/:pet_id | 宠物3D虚拟形象 + 核心体征 |
| 实时体征面板 | /digital-twin/:pet_id/vitals | 心跳/呼吸/体温实时曲线 |

### 5.2 健康预警

| 页面 | 路由 | 说明 |
|------|------|------|
| 健康预警列表 | /digital-twin/:pet_id/alerts | 预警列表 |
| 健康报告 | /digital-twin/:pet_id/health-report | 健康周报/月报 |

### 5.3 历史回放

| 页面 | 路由 | 说明 |
|------|------|------|
| 时间轴 | /digital-twin/:pet_id/timeline | 历史事件时间轴 |
| 回忆集 | /digital-twin/:pet_id/memories | AI生成的回忆集 |
| 精彩瞬间 | /digital-twin/:pet_id/highlights | 精彩瞬间相册 |
| 数据导出 | /digital-twin/:pet_id/export | 数据导出配置 |

### 5.4 行为预测

| 页面 | 路由 | 说明 |
|------|------|------|
| 行为预测 | /digital-twin/:pet_id/predictions | 当前预测 + 历史预测 |

### 5.5 跨设备同步

| 页面 | 路由 | 说明 |
|------|------|------|
| 同步状态 | /digital-twin/:pet_id/sync | 各设备同步状态 |

---

## 六、验收标准

### 6.1 实时生命体征

| 验收点 | 标准 |
|--------|------|
| 体征延迟 | 设备上报到展示延迟<2秒 |
| 心跳曲线 | 最近24小时曲线正确展示 |
| 体征准确性 | 与设备传感器数据偏差在允许范围内 |

### 6.2 行为预测

| 验收点 | 标准 |
|--------|------|
| 短期动作准确率 | >75% |
| 意图识别准确率 | >70% |
| 异常检测召回率 | >85% |

### 6.3 健康预警

| 验收点 | 标准 |
|--------|------|
| 预警触发 | 符合条件时5分钟内发出预警 |
| 预警通知 | 紧急预警电话接通 |
| 预警遗漏 | 紧急预警不遗漏 |

### 6.4 历史回放

| 验收点 | 标准 |
|--------|------|
| 数据完整性 | 历史事件记录不丢失 |
| 查询性能 | 1000条记录内查询<1秒 |
| 数据导出 | 导出PDF内容完整 |

### 6.5 跨设备同步

| 验收点 | 标准 |
|--------|------|
| 同步延迟 | 状态变更到其他设备<5秒 |
| 离线记录 | 离线100条记录内可完整保存 |
| 续传成功率 | 断网恢复后自动续传成功率>99% |
