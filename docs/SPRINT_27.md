# Sprint 27 规划

**时间**：2026-10-25
**状态**：待开始
**Sprint 周期**：2 周（2026-10-25 ～ 2026-11-07）

---

## 一、Sprint 目标

**目标：** 第三方集成 - 智能家居与宠物医疗

完成智能家居全平台对接（米家/天猫精灵/HomeKit）、宠物医疗对接（预约/病历同步）、宠物保险对接，实现MDM平台与外部生态的深度集成，构建完整的宠物生活服务闭环。

---

## 二、详细任务列表

### 后端 P0/P1/P2 任务表

| # | 任务 | 说明 | 交付物 | 优先级 |
|---|------|------|--------|--------|
| P0-1 | **米家对接** | 实现米家设备接入和控制 | miejia_adapter.go | P0 |
| P0-2 | **天猫精灵对接** | 实现天猫精灵设备接入和控制 | tmall_adapter.go | P0 |
| P0-3 | **HomeKit对接** | 实现HomeKit设备接入和控制 | homekit_adapter.go | P0 |
| P0-4 | **宠物医院预约 API** | 完成 `/api/v1/medical/appointments/*` 预约管理 | medical_appointment_controller.go | P0 |
| P1-1 | **病历同步 API** | 完成 `/api/v1/medical/records/*` 病历同步 | medical_record_controller.go | P1 |
| P1-2 | **宠物保险对接 API** | 完成 `/api/v1/insurance/*` 保险对接 | insurance_controller.go | P1 |
| P1-3 | **智能家居场景联动** | 实现宠物状态触发家居自动化 | smarthome_scene_linkage.go | P1 |
| P2-1 | **第三方健康设备同步** | 支持智能项圈等设备数据接入 | health_device_sync.go | P2 |
| P2-2 | **地图服务对接** | 高德/Google地图集成 | map_service_adapter.go | P2 |

### 前端 P0/P1/P2 任务表

| # | 任务 | 说明 | 交付物 | 优先级 |
|---|------|------|--------|--------|
| PF0-1 | **智能家居控制面板** | 完成 SmartHomeControlView.vue 家居设备控制 | SmartHomeControlView.vue | P0 |
| PF0-2 | **宠物医院预约页面** | 完成 MedicalAppointmentView.vue 预约和管理 | MedicalAppointmentView.vue | P0 |
| PF0-3 | **病历同步页面** | 完成 MedicalRecordsView.vue 病历查看和同步 | MedicalRecordsView.vue | P1 |
| PF1-1 | **宠物保险页面** | 完成 InsuranceView.vue 保险产品浏览和理赔 | InsuranceView.vue | P1 |
| PF1-2 | **场景联动配置页面** | 完成 SceneLinkageView.vue 家居自动化配置 | SceneLinkageView.vue | P1 |
| PF2-1 | **健康设备管理页面** | 完成 HealthDeviceView.vue 第三方设备绑定 | HealthDeviceView.vue | P2 |

---

## 三、技术方案

### API 路由设计

| 接口 | 方法 | 说明 |
|------|------|------|
| `GET /api/v1/smarthome/devices` | GET | 家居设备列表 |
| `POST /api/v1/smarthome/devices/:id/control` | POST | 控制设备 |
| `POST /api/v1/smarthome/linkages` | POST | 创建场景联动 |
| `GET /api/v1/smarthome/linkages` | GET | 联动列表 |
| `GET /api/v1/medical/clinics` | GET | 医院列表 |
| `POST /api/v1/medical/appointments` | POST | 创建预约 |
| `GET /api/v1/medical/appointments` | GET | 预约列表 |
| `PUT /api/v1/medical/appointments/:id` | PUT | 更新预约 |
| `DELETE /api/v1/medical/appointments/:id` | DELETE | 取消预约 |
| `GET /api/v1/medical/records` | GET | 病历列表 |
| `POST /api/v1/medical/records/sync` | POST | 同步病历 |
| `GET /api/v1/insurance/products` | GET | 保险产品列表 |
| `POST /api/v1/insurance/claims` | POST | 提交理赔 |
| `GET /api/v1/insurance/claims/:id` | GET | 理赔状态 |
| `POST /api/v1/health-devices/sync` | POST | 同步健康设备数据 |
| `GET /api/v1/maps/nearby` | GET | 附近宠物服务 |

### 数据库设计

```sql
-- 智能家居设备表
CREATE TABLE smarthome_devices (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES users(id),
    platform        VARCHAR(50) NOT NULL,
    device_id       VARCHAR(100) NOT NULL,          -- 第三方设备ID
    device_name     VARCHAR(200),
    device_type     VARCHAR(50),                    -- 'light'/'switch'/'sensor'/'camera'
    capabilities    JSONB,                           -- 设备能力
    status          JSONB,                           -- 当前状态
    last_sync_at    TIMESTAMP,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW(),
    INDEX idx_user_devices (user_id)
);

-- 智能家居场景联动表
CREATE TABLE smarthome_linkages (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES users(id),
    linkage_name    VARCHAR(100) NOT NULL,
    trigger_type    VARCHAR(50) NOT NULL,            -- 'pet_status'/'time'/'location'
    trigger_config  JSONB NOT NULL,                   -- 触发条件配置
    actions         JSONB NOT NULL,                   -- 执行动作
    is_active       BOOLEAN DEFAULT TRUE,
    last_triggered_at TIMESTAMP,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

-- 宠物病历表
CREATE TABLE pet_medical_records (
    id              BIGSERIAL PRIMARY KEY,
    pet_id          BIGINT NOT NULL REFERENCES pets(id),
    clinic_id       VARCHAR(100) NOT NULL,
    clinic_name     VARCHAR(200),
    record_type     VARCHAR(50) NOT NULL,            -- 'checkup'/'vaccination'/'surgery'/'prescription'
    diagnosis       TEXT,
    treatment       TEXT,
    prescription    JSONB,                           -- 处方
    attachments     VARCHAR(500)[],                  -- 检查报告附件
    vet_name        VARCHAR(100),
    visit_date      TIMESTAMP NOT NULL,
    next_visit_date TIMESTAMP,
    synced_at       TIMESTAMP,
    created_at      TIMESTAMP DEFAULT NOW()
);

-- 保险产品表
CREATE TABLE insurance_products (
    id              BIGSERIAL PRIMARY KEY,
    product_name    VARCHAR(200) NOT NULL,
    insurer         VARCHAR(100) NOT NULL,
    coverage_type   VARCHAR(50) NOT NULL,            -- 'accident'/'health'/'comprehensive'
    coverage_amount DECIMAL(10,2),
    monthly_premium DECIMAL(10,2),
    pet_type        VARCHAR(50)[],
    age_range       VARCHAR(50),
    description     TEXT,
    policy_url      VARCHAR(500),
    is_active       BOOLEAN DEFAULT TRUE,
    created_at      TIMESTAMP DEFAULT NOW()
);

-- 保险理赔表
CREATE TABLE insurance_claims (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES users(id),
    pet_id          BIGINT NOT NULL REFERENCES pets(id),
    product_id      BIGINT REFERENCES insurance_products(id),
    claim_type      VARCHAR(50) NOT NULL,
    claim_amount    DECIMAL(10,2),
    status          VARCHAR(20) DEFAULT 'pending',   -- 'pending'/'approved'/'rejected'/'paid'
    claim_date      TIMESTAMP NOT NULL,
    description     TEXT,
    attachments     VARCHAR(500)[],
    reviewed_at     TIMESTAMP,
    paid_at         TIMESTAMP,
    created_at      TIMESTAMP DEFAULT NOW()
);
```

---

## 四、验收标准

### 4.1 功能验收

| 用例 | 验收条件 | 测试方法 |
|------|----------|----------|
| 米家对接 | 设备发现/控制正常，响应<2s | 对接测试 |
| 天猫精灵对接 | 设备同步正常，控制可用 | 对接测试 |
| HomeKit对接 | 设备同步正常，控制可用 | 对接测试 |
| 宠物医院预约 | 预约创建/取消/提醒正常 | 预约流程测试 |
| 病历同步 | 同步成功率>90% | 同步测试 |
| 宠物保险 | 保险产品浏览和理赔提交正常 | 全流程测试 |

### 4. 性能验收

| 验收点 | 标准 |
|--------|------|
| 家居设备控制响应 | < 2s |
| 预约提醒 | 提前>24h通知 |
| 病历同步 | < 30s |

---

## 五、依赖与风险

### 依赖

| 依赖 | 说明 |
|------|------|
| PRD_28_PET_SUPPLIES_ECOMMERCE | 电商能力 |
| 宠物档案 | 宠物基础信息 |
| 第三方API | 米家/天猫精灵/HomeKit API |

### 风险

| 风险 | 影响 | 缓解措施 |
|------|------|----------|
| 第三方API不稳定 | 集成不可用 | 本地缓存+降级策略 |
| 病历数据隐私 | 隐私合规风险 | 数据加密+授权机制 |
| 保险理赔纠纷 | 用户投诉 | 清晰的理赔条款 |
