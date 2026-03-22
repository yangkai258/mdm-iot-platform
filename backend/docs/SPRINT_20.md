# Sprint 20 规划

**时间**：2026-08-23
**状态**：待开始
**Sprint 周期**：2 周（2026-08-23 ～ 2026-09-05）

---

## 一、Sprint 目标

**目标：** 家庭和多用户场景

在 Sprint 19（健康医疗）的基础上，实现家庭和多用户场景功能，包括多用户交互、家庭成员管理、儿童模式、老人陪伴模式、家庭相册，构建完整的家庭宠物陪伴体验。

---

## 二、详细任务列表

### 后端 P0/P1/P2 任务表

| # | 任务 | 说明 | 交付物 | 优先级 |
|---|------|------|--------|--------|
| P0-1 | **家庭成员管理 API** | 完成 `/api/v1/family/members/*` 成员 CRUD | family_member_controller.go | P0 |
| P0-2 | **家庭数据库** | 创建 households + household_members 表 | models/household.go | P0 |
| P0-3 | **多用户交互 API** | 完成 `/api/v1/interactions/*` 多用户交互记录 | interaction_controller.go | P0 |
| P0-4 | **家庭相册 API** | 完成 `/api/v1/family/albums/*` 家庭照片管理 | family_album_controller.go | P0 |
| P0-5 | **家庭配置 API** | 完成 `/api/v1/family/:family_id/settings` 家庭设置 | family_settings_controller.go | P0 |
| P1-1 | **儿童模式 API** | 完成 `/api/v1/family/:family_id/child-mode/*` 儿童模式 | child_mode_controller.go | P1 |
| P1-2 | **老人陪伴模式 API** | 完成 `/api/v1/family/:family_id/elder-mode/*` 老人模式 | elder_mode_controller.go | P1 |
| P1-3 | **家庭通知路由** | 实现家庭成员通知分发逻辑 | family/notification_router.go | P1 |
| P2-1 | **家庭动态订阅** | 家庭成员实时状态订阅 | family/realtime_subscription.go | P2 |
| P2-2 | **家庭分享邀请** | 完成家庭邀请链接生成和验证 | family/invite_service.go | P2 |

### 前端 P0/P1/P2 任务表

| # | 任务 | 说明 | 交付物 | 优先级 |
|---|------|------|--------|--------|
| PF0-1 | **家庭成员管理页面** | 完成 FamilyMembersView.vue 成员列表/邀请/移除 | FamilyMembersView.vue | P0 |
| PF0-2 | **儿童模式前端** | 完成 ChildModeView.vue 儿童模式配置和使用 | ChildModeView.vue | P0 |
| PF0-3 | **老人陪伴模式前端** | 完成 ElderModeView.vue 老人模式配置 | ElderModeView.vue | P0 |
| PF0-4 | **家庭相册前端** | 完成 FamilyAlbumView.vue 家庭照片管理 | FamilyAlbumView.vue | P0 |
| PF0-5 | **家庭设置页面** | 完成 FamilySettingsView.vue 家庭配置 | FamilySettingsView.vue | P0 |
| PF1-1 | **家庭动态页面** | 完成 FamilyActivityView.vue 家庭动态时间轴 | FamilyActivityView.vue | P1 |
| PF1-2 | **家庭邀请页面** | 完成 FamilyInviteView.vue 邀请成员 | FamilyInviteView.vue | P1 |
| PF2-1 | **家庭统计页面** | 完成 FamilyStatsView.vue 家庭宠物数据统计 | FamilyStatsView.vue | P2 |

---

## 三、技术方案

### API 路由设计

| 接口 | 方法 | 说明 |
|------|------|------|
| `GET /api/v1/family` | GET | 获取家庭信息 |
| `POST /api/v1/family` | POST | 创建家庭 |
| `GET /api/v1/family/:family_id` | GET | 家庭详情 |
| `PUT /api/v1/family/:family_id` | PUT | 更新家庭信息 |
| `GET /api/v1/family/members` | GET | 家庭成员列表 |
| `POST /api/v1/family/members/invite` | POST | 邀请成员 |
| `PUT /api/v1/family/members/:member_id` | PUT | 更新成员信息/角色 |
| `DELETE /api/v1/family/members/:member_id` | DELETE | 移除成员 |
| `POST /api/v1/family/members/:member_id/accept` | POST | 接受邀请 |
| `POST /api/v1/family/members/:member_id/leave` | POST | 离开家庭 |
| `GET /api/v1/family/:family_id/settings` | GET | 家庭设置 |
| `PUT /api/v1/family/:family_id/settings` | PUT | 更新家庭设置 |
| `GET /api/v1/interactions` | GET | 交互记录列表 |
| `POST /api/v1/interactions` | POST | 记录交互 |
| `GET /api/v1/interactions/summary` | GET | 交互汇总 |
| `GET /api/v1/family/albums` | GET | 家庭相册列表 |
| `POST /api/v1/family/albums` | POST | 上传照片 |
| `DELETE /api/v1/family/albums/:id` | DELETE | 删除照片 |
| `PUT /api/v1/family/albums/:id/tags` | PUT | 更新照片标签 |
| `GET /api/v1/family/:family_id/child-mode` | GET | 儿童模式配置 |
| `PUT /api/v1/family/:family_id/child-mode` | PUT | 更新儿童模式配置 |
| `GET /api/v1/family/:family_id/elder-mode` | GET | 老人模式配置 |
| `PUT /api/v1/family/:family_id/elder-mode` | PUT | 更新老人模式配置 |
| `GET /api/v1/family/:family_id/activities` | GET | 家庭动态 |

### 数据库设计

```sql
-- 家庭表
CREATE TABLE households (
    id              BIGSERIAL PRIMARY KEY,
    household_name  VARCHAR(100) NOT NULL,
    household_uuid  VARCHAR(64) NOT NULL UNIQUE,
    owner_id        BIGINT NOT NULL REFERENCES users(id),
    household_type  VARCHAR(20) DEFAULT 'standard', -- 'standard'/'elder_care'/'child_friendly'
    avatar_url      VARCHAR(500),
    settings        JSONB,                          -- 家庭配置
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

-- 家庭成员表
CREATE TABLE household_members (
    id              BIGSERIAL PRIMARY KEY,
    household_id    BIGINT NOT NULL REFERENCES households(id),
    user_id         BIGINT NOT NULL REFERENCES users(id),
    role            VARCHAR(20) NOT NULL,           -- 'owner'/'admin'/'member'/'viewer'
    display_name    VARCHAR(50),
    avatar_url      VARCHAR(500),
    member_type     VARCHAR(20),                     -- 'adult'/'child'/'elder'
    notification_pref JSONB,                         -- 通知偏好
    is_active       BOOLEAN DEFAULT TRUE,
    joined_at       TIMESTAMP DEFAULT NOW(),
    UNIQUE(household_id, user_id)
);

-- 家庭邀请表
CREATE TABLE household_invites (
    id              BIGSERIAL PRIMARY KEY,
    household_id    BIGINT NOT NULL REFERENCES households(id),
    invite_code     VARCHAR(20) NOT NULL UNIQUE,
    invited_email   VARCHAR(255),
    invited_role    VARCHAR(20) DEFAULT 'member',
    invited_by      BIGINT NOT NULL REFERENCES users(id),
    expires_at      TIMESTAMP,
    status          VARCHAR(20) DEFAULT 'pending',  -- 'pending'/'accepted'/'declined'/'expired'
    created_at      TIMESTAMP DEFAULT NOW()
);

-- 家庭模式配置表
CREATE TABLE family_mode_configs (
    id              BIGSERIAL PRIMARY KEY,
    household_id    BIGINT NOT NULL REFERENCES households(id),
    mode_type       VARCHAR(30) NOT NULL,         -- 'child_mode'/'elder_mode'
    is_enabled      BOOLEAN DEFAULT FALSE,
    config          JSONB,                          -- 模式具体配置
    enabled_members BIGINT[],                       -- 启用的成员ID列表
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW(),
    UNIQUE(household_id, mode_type)
);

-- 交互记录表
CREATE TABLE interaction_records (
    id              BIGSERIAL PRIMARY KEY,
    household_id    BIGINT NOT NULL REFERENCES households(id),
    pet_id          BIGINT REFERENCES pets(id),
    user_id         BIGINT NOT NULL REFERENCES users(id),
    interaction_type VARCHAR(50) NOT NULL,         -- 'play'/'feed'/'walk'/'groom'/'medical'/'other'
    participant_ids BIGINT[],
    duration_minutes INT,
    notes           TEXT,
    mood_before     INT,
    mood_after      INT,
    media_urls      VARCHAR(500)[],
    occurred_at     TIMESTAMP NOT NULL,
    created_at      TIMESTAMP DEFAULT NOW(),
    INDEX idx_household_time (household_id, occurred_at DESC),
    INDEX idx_pet_time (pet_id, occurred_at DESC)
);

-- 家庭相册表
CREATE TABLE family_albums (
    id              BIGSERIAL PRIMARY KEY,
    household_id    BIGINT NOT NULL REFERENCES households(id),
    uploader_id     BIGINT NOT NULL REFERENCES users(id),
    photo_url       VARCHAR(500) NOT NULL,
    thumbnail_url   VARCHAR(500),
    caption         VARCHAR(255),
    tags            VARCHAR(50)[],
    location_name   VARCHAR(100),
    taken_at        TIMESTAMP,
    pet_ids         BIGINT[],
    visibility      VARCHAR(20) DEFAULT 'family',  -- 'family'/'public'
    created_at      TIMESTAMP DEFAULT NOW(),
    INDEX idx_household_time (household_id, taken_at DESC)
);

-- 家庭动态表
CREATE TABLE family_activities (
    id              BIGSERIAL PRIMARY KEY,
    household_id    BIGINT NOT NULL REFERENCES households(id),
    activity_type   VARCHAR(50) NOT NULL,
    actor_id        BIGINT REFERENCES users(id),
    target_type     VARCHAR(50),
    target_id       BIGINT,
    summary         VARCHAR(255),
    detail          JSONB,
    occurred_at     TIMESTAMP NOT NULL,
    created_at      TIMESTAMP DEFAULT NOW(),
    INDEX idx_household_time (household_id, occurred_at DESC)
);
```

---

## 四、验收标准

### 4.1 功能验收

| 用例 | 验收条件 | 测试方法 |
|------|----------|----------|
| 家庭创建 | 用户成功创建家庭并成为 owner | 调用 API 验证 |
| 邀请成员 | 被邀请者收到邀请并成功加入 | 邀请流程测试 |
| 角色权限 | owner/admin/member/viewer 权限隔离正确 | 各角色操作测试 |
| 儿童模式 | 内容过滤和时间控制正确生效 | 配置验证 |
| 老人模式 | 界面简化和语音交互正确 | 配置验证 |
| 家庭相册 | 照片上传/删除/标签正确 | 上传测试 |
| 家庭动态 | 动态实时更新且正确分发 | 多人操作测试 |

### 4.2 性能验收

| 验收点 | 标准 |
|--------|------|
| 家庭成员查询 | <= 100ms（100 成员内） |
| 动态推送延迟 | <= 2s |
| 相册加载 | <= 1s（100 张照片） |

---

## 五、依赖与风险

### 依赖

| 依赖 | 说明 |
|------|------|
| Sprint 15 宠物生态 | 宠物档案和家庭宠物 |
| Sprint 17 情感计算 | 情绪记录和响应 |
| 文件存储服务 | 相册照片存储 |

### 风险

| 风险 | 影响 | 缓解措施 |
|------|------|----------|
| 儿童模式绕过 | 内容过滤失效 | 持续监控+用户反馈 |
| 隐私泄露 | 家庭照片意外公开 | 权限验证+数据审核 |
| 家庭成员纠纷 | 宠物管理权争议 | 明确的角色权限体系 |
