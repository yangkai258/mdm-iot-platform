# Sprint 29 规划

**时间**：2026-11-15
**状态**：待开始
**Sprint 周期**：2 周（2026-11-15 ～ 2026-11-28）

---

## 一、Sprint 目标

**目标：** 高级用户功能 - 家庭相册AI筛选与模式完善

在 Sprint 20 家庭功能基础上，完善家庭相册AI筛选版（精彩瞬间/可爱表情/互动场景）、儿童模式、老人陪伴模式、寻回网络，构建完整的家庭宠物陪伴体验。

---

## 二、详细任务列表

### 后端 P0/P1/P2 任务表

| # | 任务 | 说明 | 交付物 | 优先级 |
|---|------|------|--------|--------|
| P0-1 | **家庭相册 API** | 完成 `/api/v1/album/*` 相册管理（基于 PRD_30_FAMILY_ALBUM） | album_controller.go | P0 |
| P0-2 | **AI筛选服务** | 实现精彩瞬间/可爱表情/互动场景识别 | ai_filter_service.go | P0 |
| P0-3 | **AI精选专辑生成** | 自动生成AI精选专辑 | ai_album_generator.go | P0 |
| P1-1 | **寻回网络 API** | 完成 `/api/v1/pet-finder/*` 寻回网络 | pet_finder_controller.go | P1 |
| P1-2 | **寻回协查通知** | 实现丢失宠物协查广播 | finder_notify.go | P1 |
| P2-1 | **儿童模式完善** | 完善内容过滤和使用时长控制 | child_mode_enhanced.go | P2 |
| P2-2 | **老人陪伴模式完善** | 简化UI+主动问候+紧急求助 | elder_mode_enhanced.go | P2 |

### 前端 P0/P1/P2 任务表

| # | 任务 | 说明 | 交付物 | 优先级 |
|---|------|------|--------|--------|
| PF0-1 | **相册首页** | 完成 AlbumHomeView.vue 时间线照片流 | AlbumHomeView.vue | P0 |
| PF0-2 | **AI筛选页面** | 完成 AIFilterView.vue 精彩瞬间/表情/互动 | AIFilterView.vue | P0 |
| PF0-3 | **AI专辑页面** | 完成 AIAlbumView.vue AI精选专辑 | AIAlbumView.vue | P0 |
| PF0-4 | **相册分类管理页面** | 完成 AlbumCategoriesView.vue 分类管理 | AlbumCategoriesView.vue | P0 |
| PF1-1 | **寻回网络页面** | 完成 PetFinderView.vue 丢失协查和定位 | PetFinderView.vue | P1 |
| PF1-2 | **儿童模式完善页面** | 完成 ChildModeEnhancedView.vue 内容过滤配置 | ChildModeEnhancedView.vue | P2 |
| PF2-1 | **老人陪伴模式页面** | 完成 ElderModeEnhancedView.vue 简化UI配置 | ElderModeEnhancedView.vue | P2 |

---

## 三、技术方案

### API 路由设计

| 接口 | 方法 | 说明 |
|------|------|------|
| `GET /api/v1/album/photos` | GET | 照片列表 |
| `POST /api/v1/album/photos` | POST | 上传照片 |
| `GET /api/v1/album/photos/ai-filter` | GET | AI筛选照片 |
| `GET /api/v1/album/ai-albums` | GET | AI精选专辑列表 |
| `POST /api/v1/album/ai-albums/:id/regenerate` | POST | 重新生成专辑 |
| `GET /api/v1/album/categories` | GET | 分类列表 |
| `POST /api/v1/album/family/invite` | POST | 邀请家庭成员 |
| `GET /api/v1/pet-finder/report` | POST | 报告丢失宠物 |
| `GET /api/v1/pet-finder/reports` | GET | 协查列表 |
| `GET /api/v1/pet-finder/reports/:id` | GET | 协查详情 |
| `POST /api/v1/pet-finder/reports/:id/sighting` | POST | 报告发现 |
| `GET /api/v1/child-mode/config` | GET | 儿童模式配置 |
| `PUT /api/v1/child-mode/config` | PUT | 更新儿童模式配置 |
| `GET /api/v1/elder-mode/config` | GET | 老人模式配置 |
| `PUT /api/v1/elder-mode/config` | PUT | 更新老人模式配置 |

### 数据库设计

```sql
-- 相册表（已在 PRD_30_FAMILY_ALBUM.md 定义）
-- 复用 album_photos / album_ai_albums / album_categories / album_family_members

-- 寻回网络表
CREATE TABLE pet_finder_reports (
    id              BIGSERIAL PRIMARY KEY,
    pet_id          BIGINT NOT NULL REFERENCES pets(id),
    reporter_id     BIGINT NOT NULL REFERENCES users(id),
    report_type     VARCHAR(20) NOT NULL,            -- 'lost'/'found'
    last_seen_location VARCHAR(500),
    last_seen_time  TIMESTAMP,
    description     TEXT,
    reward          DECIMAL(10,2),
    status          VARCHAR(20) DEFAULT 'active',  -- 'active'/'resolved'/'closed'
    contact_phone   VARCHAR(20),
    photo_urls      VARCHAR(500)[],
    resolved_at     TIMESTAMP,
    created_at      TIMESTAMP DEFAULT NOW(),
    INDEX idx_active_reports (status, created_at DESC),
    INDEX idx_pet_reports (pet_id)
);

-- 寻回网络发现记录表
CREATE TABLE pet_finder_sightings (
    id              BIGSERIAL PRIMARY KEY,
    report_id       BIGINT NOT NULL REFERENCES pet_finder_reports(id),
    reporter_id     BIGINT REFERENCES users(id),
    sighting_location VARCHAR(500) NOT NULL,
    sighting_time   TIMESTAMP NOT NULL,
    description     TEXT,
    confidence      DECIMAL(5,2),
    verified        BOOLEAN DEFAULT FALSE,
    created_at      TIMESTAMP DEFAULT NOW()
);
```

---

## 四、验收标准

### 4.1 功能验收

| 用例 | 验收条件 | 测试方法 |
|------|----------|----------|
| 家庭相册 | 照片上传/删除/AI筛选正常，响应<1s | 相册功能测试 |
| AI筛选 | 精彩瞬间识别>75%，可爱表情>80%，互动场景>70% | 识别评测 |
| AI专辑 | 新照片入库后24h内自动更新 | 自动化测试 |
| 寻回网络 | 协查发布<1min，定位精度<50m | GPS测试 |

### 4.2 性能验收

| 验收点 | 标准 |
|--------|------|
| 照片加载 | 首屏<1s |
| AI筛选 | < 5s/张照片 |
| 协查广播 | < 1min触达所有用户 |

---

## 五、依赖与风险

### 依赖

| 依赖 | 说明 |
|------|------|
| PRD_30_FAMILY_ALBUM | 家庭相册 AI筛选 PRD |
| Sprint 20 家庭功能 | 家庭基础能力 |
| AI识别服务 | 图像识别能力 |

### 风险

| 风险 | 影响 | 缓解措施 |
|------|------|----------|
| AI识别准确率不足 | 用户体验差 | 持续优化模型 |
| 寻回网络隐私 | 位置信息泄露 | 数据脱敏+授权机制 |
| 儿童模式绕过 | 内容过滤失效 | 持续监控+用户反馈 |
