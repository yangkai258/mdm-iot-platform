# Sprint 26 规划

**时间**：2026-10-11
**状态**：待开始
**Sprint 周期**：2 周（2026-10-11 ～ 2026-10-24）

---

## 一、Sprint 目标

**目标：** 开放平台完善 - 市场与第三方集成基础

完善开放平台生态，包括App/插件市场、表情包市场、动作资源库、声音定制，并开始第三方集成（智能家居、宠物医疗）的能力建设。

---

## 二、详细任务列表

### 后端 P0/P1/P2 任务表

| # | 任务 | 说明 | 交付物 | 优先级 |
|---|------|------|--------|--------|
| P0-1 | **App/插件市场 API** | 完成 `/api/v1/market/apps/*` 应用市场 | market_app_controller.go | P0 |
| P0-2 | **表情包市场 API** | 完成 `/api/v1/market/emojis/*` 表情包市场 | emoji_market_controller.go | P0 |
| P0-3 | **动作资源库 API** | 完成 `/api/v1/market/actions/*` 动作资源库 | action_market_controller.go | P0 |
| P0-4 | **声音定制 API** | 完成 `/api/v1/market/voices/*` 声音定制 | voice_market_controller.go | P0 |
| P1-1 | **智能家居对接基础** | 实现米家/天猫精灵/HomeKit对接架构 | smarthome_adapter.go | P1 |
| P1-2 | **宠物医疗对接基础** | 实现宠物医疗API对接架构 | pet_medical_adapter.go | P1 |
| P1-3 | **市场审核工作流** | 实现应用/表情/动作审核流程 | market_review_workflow.go | P1 |
| P2-1 | **开发者等级体系** | 实现初级/中级/高级开发者分级 | developer_tier.go | P2 |
| P2-2 | **市场数据分析** | 实现市场数据统计分析 | market_analytics.go | P2 |

### 前端 P0/P1/P2 任务表

| # | 任务 | 说明 | 交付物 | 优先级 |
|---|------|------|--------|--------|
| PF0-1 | **App市场页面** | 完成 AppMarketView.vue 应用浏览和详情 | AppMarketView.vue | P0 |
| PF0-2 | **表情包市场页面** | 完成 EmojiMarketView.vue 表情包浏览和购买 | EmojiMarketView.vue | P0 |
| PF0-3 | **动作资源库页面** | 完成 ActionLibraryView.vue 动作浏览和管理 | ActionLibraryView.vue | P0 |
| PF0-4 | **声音定制页面** | 完成 VoiceCustomizeView.vue 声音定制配置 | VoiceCustomizeView.vue | P0 |
| PF1-1 | **智能家居配置页面** | 完成 SmartHomeConfigView.vue 设备对接配置 | SmartHomeConfigView.vue | P1 |
| PF1-2 | **宠物医疗对接页面** | 完成 PetMedicalView.vue 宠物医疗对接配置 | PetMedicalView.vue | P1 |
| PF2-1 | **开发者等级页面** | 完成 DeveloperTierView.vue 等级权益展示 | DeveloperTierView.vue | P2 |

---

## 三、技术方案

### API 路由设计

| 接口 | 方法 | 说明 |
|------|------|------|
| `GET /api/v1/market/apps` | GET | 应用列表 |
| `GET /api/v1/market/apps/:id` | GET | 应用详情 |
| `POST /api/v1/market/apps/:id/install` | POST | 安装应用 |
| `GET /api/v1/market/emojis` | GET | 表情包列表 |
| `GET /api/v1/market/emojis/:id` | GET | 表情包详情 |
| `POST /api/v1/market/emojis/:id/download` | POST | 下载表情包 |
| `GET /api/v1/market/actions` | GET | 动作列表 |
| `GET /api/v1/market/actions/:id` | GET | 动作详情 |
| `POST /api/v1/market/actions/:id/download` | POST | 下载动作 |
| `GET /api/v1/market/voices` | GET | 声音列表 |
| `POST /api/v1/market/voices/tts` | POST | TTS合成 |
| `POST /api/v1/market/voices/clone` | POST | 声音克隆 |
| `GET /api/v1/market/smarthome/devices` | GET | 智能家居设备列表 |
| `POST /api/v1/market/smarthome/connect` | POST | 连接智能家居 |
| `GET /api/v1/market/medical/clinics` | GET | 宠物医院列表 |
| `POST /api/v1/market/medical/sync` | POST | 同步病历 |

### 数据库设计

```sql
-- 应用市场表
CREATE TABLE market_apps (
    id              BIGSERIAL PRIMARY KEY,
    developer_id    BIGINT REFERENCES developers(id),
    app_name        VARCHAR(100) NOT NULL,
    app_type        VARCHAR(30) NOT NULL,
    category        VARCHAR(50),
    description     TEXT,
    icon_url        VARCHAR(500),
    screenshots     VARCHAR(500)[],
    version         VARCHAR(20),
    price           DECIMAL(10,2) DEFAULT 0,
    rating          DECIMAL(3,2),
    install_count   INT DEFAULT 0,
    review_status   VARCHAR(20) DEFAULT 'pending', -- 'pending'/'approved'/'rejected'
    published_at    TIMESTAMP,
    created_at      TIMESTAMP DEFAULT NOW()
);

-- 表情包市场表
CREATE TABLE market_emojis (
    id              BIGSERIAL PRIMARY KEY,
    developer_id    BIGINT REFERENCES developers(id),
    emoji_name      VARCHAR(100) NOT NULL,
    emoji_type      VARCHAR(30) NOT NULL,            -- 'static'/'animated'/'3d'
    preview_urls    VARCHAR(500)[],
    file_url        VARCHAR(500) NOT NULL,
    file_size       BIGINT,
    price           DECIMAL(10,2) DEFAULT 0,
    download_count  INT DEFAULT 0,
    rating          DECIMAL(3,2),
    tags            VARCHAR(50)[],
    review_status   VARCHAR(20) DEFAULT 'pending',
    created_at      TIMESTAMP DEFAULT NOW()
);

-- 动作资源库表
CREATE TABLE market_actions (
    id              BIGSERIAL PRIMARY KEY,
    developer_id    BIGINT REFERENCES developers(id),
    action_name     VARCHAR(100) NOT NULL,
    action_category VARCHAR(50),
    description     TEXT,
    preview_video   VARCHAR(500),
    action_file_url VARCHAR(500) NOT NULL,
    compatibility   VARCHAR(50)[],                   -- 兼容设备型号
    difficulty      VARCHAR(20),                    -- 'easy'/'medium'/'hard'
    price           DECIMAL(10,2) DEFAULT 0,
    download_count  INT DEFAULT 0,
    rating          DECIMAL(3,2),
    review_status   VARCHAR(20) DEFAULT 'pending',
    created_at      TIMESTAMP DEFAULT NOW()
);

-- 声音定制表
CREATE TABLE market_voices (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES users(id),
    voice_name      VARCHAR(100) NOT NULL,
    voice_type      VARCHAR(30) NOT NULL,            -- 'tts'/'clone'
    voice_config    JSONB,                           -- TTS配置或克隆配置
    audio_sample_url VARCHAR(500),
    is_active       BOOLEAN DEFAULT TRUE,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

-- 智能家居对接配置表
CREATE TABLE smarthome_configs (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES users(id),
    platform        VARCHAR(50) NOT NULL,            -- 'xiaomi'/'天猫精灵'/'homekit'
    auth_config     JSONB,                           -- 授权配置
    linked_devices  JSONB,                           -- 关联的设备
    is_active       BOOLEAN DEFAULT TRUE,
    last_sync_at    TIMESTAMP,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

-- 宠物医疗对接表
CREATE TABLE pet_medical_configs (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES users(id),
    clinic_id       VARCHAR(100),
    clinic_name     VARCHAR(200),
    auth_token      VARCHAR(255),
    linked_pet_ids  BIGINT[],
    is_active       BOOLEAN DEFAULT TRUE,
    created_at      TIMESTAMP DEFAULT NOW()
);
```

---

## 四、验收标准

### 4.1 功能验收

| 用例 | 验收条件 | 测试方法 |
|------|----------|----------|
| App市场 | 应用上架/审核/下架流程完整 | 上架测试 |
| 表情包市场 | 表情包上传/审核/下载正常 | 全流程测试 |
| 动作资源库 | 动作上传/审核/下载正常 | 全流程测试 |
| 声音定制 | TTS合成正常，声音克隆效果可接受 | 主观评测 |
| 智能家居对接 | 米家/天猫精灵/HomeKit至少1个对接完成 | 对接测试 |

### 4.2 性能验收

| 验收点 | 标准 |
|--------|------|
| 市场列表加载 | < 1s |
| 应用安装 | < 3s |
| 声音TTS合成 | < 2s |

---

## 五、依赖与风险

### 依赖

| 依赖 | 说明 |
|------|------|
| Sprint 25 开放平台核心 | 开发者API基础 |
| MODULE_PLATFORM_ECOSYSTEM | 开放平台生态完整 PRD |
| 第三方API | 米家/天猫精灵/HomeKit API |

### 风险

| 风险 | 影响 | 缓解措施 |
|------|------|----------|
| 第三方API不稳定 | 集成体验差 | 本地缓存+降级策略 |
| 内容审核效率低 | 市场质量差 | AI辅助审核 |
| 版权侵权 | 法律风险 | 内容审核+举报机制 |
