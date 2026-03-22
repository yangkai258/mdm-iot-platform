# Sprint 15 规划

**时间**：2026-06-14
**状态**：✅ 已完成
**Sprint 周期**：2 周（2026-06-14 ～ 2026-06-27）

---

## 一、Sprint 目标

**目标：** 宠物相关扩展功能

在 Sprint 14（AI 系统工程）的基础上，实现宠物相关扩展功能，包括宠物登记、寻回网络、多宠物管理，为宠物主人提供更完整的宠物生命周期管理服务。

---

## 二、详细任务列表

### 后端 P0/P1/P2 任务表

| # | 任务 | 说明 | 交付物 | 优先级 |
|---|------|------|--------|--------|
| P0-1 | **宠物登记 API** | 完成 `/api/v1/pets/*` 宠物档案 CRUD | pet_controller.go | P0 |
| P0-2 | **宠物档案数据库** | 创建 pets 表存储宠物完整档案 | models/pet.go | P0 |
| P0-3 | **寻回网络 API** | 完成 `/api/v1/lost-found/*` 寻回网络接口 | lost_found_controller.go | P0 |
| P0-4 | **多宠物管理 API** | 完成 `/api/v1/household/pets/*` 家庭宠物管理 | household_pet_controller.go | P0 |
| P0-5 | **宠物关系绑定** | 完成 `/api/v1/pets/{pet_id}/devices/*` 宠物-设备绑定 | pet_device_controller.go | P0 |
| P1-1 | **寻回网络扩散算法** | 实现失宠信息自动扩散到附近用户 | lost_found/spread_service.go | P1 |
| P1-2 | **宠物档案导入导出** | 完成宠物数据批量导入/导出 | pet_import_export_service.go | P1 |
| P1-3 | **宠物生日/疫苗提醒** | 完成宠物健康提醒 API | pet_health_reminder_controller.go | P1 |
| P2-1 | **宠物社区互动** | 完成 `/api/v1/pet-community/*` 社区接口 | community_controller.go | P2 |
| P2-2 | **宠物相册 API** | 完成 `/api/v1/pets/{pet_id}/albums/*` 相册管理 | pet_album_controller.go | P2 |

### 前端 P0/P1/P2 任务表

| # | 任务 | 说明 | 交付物 | 优先级 |
|---|------|------|--------|--------|
| PF0-1 | **宠物登记页面** | 完成 PetRegistrationView.vue 宠物信息录入 | PetRegistrationView.vue | P0 |
| PF0-2 | **宠物列表页面** | 完成 PetListView.vue 宠物列表/搜索 | PetListView.vue | P0 |
| PF0-3 | **宠物详情页面** | 完成 PetDetailView.vue 宠物档案详情 | PetDetailView.vue | P0 |
| PF0-4 | **寻回网络页面** | 完成 LostFoundView.vue 失宠发布/浏览 | LostFoundView.vue | P0 |
| PF0-5 | **多宠物管理页面** | 完成 MultiPetManageView.vue 家庭多宠物管理 | MultiPetManageView.vue | P0 |
| PF1-1 | **宠物绑定设备** | 完成 PetDeviceBindingView.vue 宠物-设备绑定 | PetDeviceBindingView.vue | P1 |
| PF1-2 | **寻回信息详情** | 完成 LostFoundDetailView.vue 寻宠详情页 | LostFoundDetailView.vue | P1 |
| PF2-1 | **宠物相册页面** | 完成 PetAlbumView.vue 宠物相册管理 | PetAlbumView.vue | P2 |

---

## 三、技术方案

### API 路由设计

| 接口 | 方法 | 说明 |
|------|------|------|
| `GET /api/v1/pets` | GET | 宠物列表（支持用户筛选） |
| `POST /api/v1/pets` | POST | 创建宠物档案 |
| `GET /api/v1/pets/:pet_id` | GET | 宠物详情 |
| `PUT /api/v1/pets/:pet_id` | PUT | 更新宠物档案 |
| `DELETE /api/v1/pets/:pet_id` | DELETE | 删除宠物档案 |
| `GET /api/v1/pets/:pet_id/devices` | GET | 宠物绑定的设备列表 |
| `POST /api/v1/pets/:pet_id/devices` | POST | 绑定设备到宠物 |
| `DELETE /api/v1/pets/:pet_id/devices/:device_id` | DELETE | 解除宠物设备绑定 |
| `GET /api/v1/lost-found` | GET | 寻回信息列表 |
| `POST /api/v1/lost-found` | POST | 发布失宠信息 |
| `GET /api/v1/lost-found/:id` | GET | 寻回信息详情 |
| `PUT /api/v1/lost-found/:id` | PUT | 更新寻回信息 |
| `DELETE /api/v1/lost-found/:id` | DELETE | 删除/撤销寻回信息 |
| `POST /api/v1/lost-found/:id/sighting` | POST | 上报目击信息 |
| `POST /api/v1/lost-found/:id/resolve` | POST | 标记已找到 |
| `GET /api/v1/household/pets` | GET | 家庭宠物列表 |
| `POST /api/v1/household/pets/invite` | POST | 邀请家庭成员管理宠物 |
| `GET /api/v1/pets/:pet_id/albums` | GET | 宠物相册 |
| `POST /api/v1/pets/:pet_id/albums` | POST | 上传照片 |
| `DELETE /api/v1/pets/:pet_id/albums/:id` | DELETE | 删除照片 |

### 数据库设计

```sql
-- 宠物档案表
CREATE TABLE pets (
    id              BIGSERIAL PRIMARY KEY,
    pet_uuid        VARCHAR(64) NOT NULL UNIQUE,
    pet_name        VARCHAR(50) NOT NULL,
    pet_type        VARCHAR(30) NOT NULL,         -- 'dog'/'cat'/'bird'/'rabbit'/'other'
    breed           VARCHAR(50),
    gender          VARCHAR(10),
    birth_date      DATE,
    weight          DECIMAL(5,2),
    height          DECIMAL(5,2),
    color           VARCHAR(50),
    microchip_no    VARCHAR(100),
    description     TEXT,
    avatar_url      VARCHAR(500),
    owner_id        BIGINT NOT NULL REFERENCES users(id),
    household_id    BIGINT REFERENCES households(id),
    status          VARCHAR(20) DEFAULT 'active',  -- 'active'/'lost'/'found'/'deceased'
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

-- 宠物-设备绑定表
CREATE TABLE pet_device_bindings (
    id              BIGSERIAL PRIMARY KEY,
    pet_id          BIGINT NOT NULL REFERENCES pets(id),
    device_id       VARCHAR(64) NOT NULL REFERENCES devices(id),
    binding_type    VARCHAR(20) DEFAULT 'primary', -- 'primary'/'secondary'
    is_active       BOOLEAN DEFAULT TRUE,
    bound_at        TIMESTAMP DEFAULT NOW(),
    unbound_at      TIMESTAMP,
    UNIQUE(pet_id, device_id)
);

-- 寻回信息表
CREATE TABLE lost_found_reports (
    id              BIGSERIAL PRIMARY KEY,
    pet_id          BIGINT REFERENCES pets(id),
    reporter_id      BIGINT NOT NULL REFERENCES users(id),
    report_type      VARCHAR(10) NOT NULL,         -- 'lost'/'found'
    title           VARCHAR(200) NOT NULL,
    description     TEXT,
    lost_location    VARCHAR(255),
    lost_lat        DECIMAL(10,8),
    lost_lng        DECIMAL(11,8),
    last_seen_time  TIMESTAMP,
    contact_method  VARCHAR(100),
    reward          DECIMAL(10,2),
    photos          VARCHAR(500)[],
    status          VARCHAR(20) DEFAULT 'active',  -- 'active'/'resolved'/'expired'
    resolved_at     TIMESTAMP,
    resolved_note   TEXT,
    spread_radius_km DECIMAL(5,2) DEFAULT 10,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW(),
    INDEX idx_location (lost_lat, lost_lng),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at DESC)
);

-- 目击信息表
CREATE TABLE sighting_reports (
    id              BIGSERIAL PRIMARY KEY,
    report_id       BIGINT NOT NULL REFERENCES lost_found_reports(id),
    sighting_location VARCHAR(255),
    sighting_lat    DECIMAL(10,8),
    sighting_lng    DECIMAL(11,8),
    sighting_time   TIMESTAMP NOT NULL,
    reporter_id     BIGINT NOT NULL REFERENCES users(id),
    description     TEXT,
    photo_url       VARCHAR(500),
    is_credible     BOOLEAN DEFAULT TRUE,
    created_at      TIMESTAMP DEFAULT NOW()
);

-- 家庭宠物邀请表
CREATE TABLE household_pet_invites (
    id              BIGSERIAL PRIMARY KEY,
    pet_id          BIGINT NOT NULL REFERENCES pets(id),
    household_id    BIGINT NOT NULL REFERENCES households(id),
    invite_code     VARCHAR(20) NOT NULL UNIQUE,
    invited_email   VARCHAR(255),
    role            VARCHAR(20) DEFAULT 'member',  -- 'owner'/'member'/'viewer'
    status          VARCHAR(20) DEFAULT 'pending', -- 'pending'/'accepted'/'declined'/'expired'
    invited_by      BIGINT NOT NULL,
    expires_at      TIMESTAMP,
    created_at      TIMESTAMP DEFAULT NOW()
);
```

---

## 四、验收标准

### 4.1 功能验收

| 用例 | 验收条件 | 测试方法 |
|------|----------|----------|
| 宠物档案 CRUD | 完整增删改查 | 调用各接口验证 |
| 宠物-设备绑定 | 一宠物可绑定多设备 | 绑定测试 |
| 寻回信息发布 | 失宠信息成功发布并可见 | 发布测试 |
| 目击上报 | 目击信息正确关联 | 上报测试 |
| 寻回扩散 | 附近用户能收到推送 | 距离测试 |
| 多宠物管理 | 家庭内宠物统一管理 | 邀请+接受测试 |

### 4.2 性能验收

| 验收点 | 标准 |
|--------|------|
| 寻回信息查询 | <= 200ms（方圆 50km） |
| 寻回扩散推送 | <= 30s（1000 用户） |

---

## 五、依赖与风险

### 依赖

| 依赖 | 说明 |
|------|------|
| Sprint 9 宠物基础 | 宠物模型和关系 |
| 地理位置服务 | 寻回网络扩散 |

### 风险

| 风险 | 影响 | 缓解措施 |
|------|------|----------|
| 寻回信息虚假 | 用户体验差 | 审核机制+举报功能 |
| 隐私泄露 | 位置信息暴露 | 数据脱敏+授权控制 |
| 宠物数据量大 | DB 存储压力 | 分表+定期清理 |

---

## 六、Sprint 15 完成清单

### 后端 (agenthd)
- [x] 宠物档案 CRUD API - `controllers/pet_controller.go` + `models/pet.go`
- [x] 设备绑定 API - `PetCtrl.RegisterPetRoutes()` (pets/:pet_id/devices)
- [x] 寻回网络 API - `controllers/lost_found_controller.go`
- [x] 家庭宠物管理 API - `controllers/household_pet_controller.go`
- [x] 宠物健康提醒 API - `controllers/pet_health_controller.go`
- [x] 数据库迁移 - `main.go` AutoMigrate 已注册 Sprint 15 所有模型
- [x] 路由注册 - `main.go` 已注册所有 Sprint 15 路由
