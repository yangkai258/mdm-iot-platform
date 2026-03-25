# Sprint 21 规划

**时间**：2026-08-30
**状态**：待开始
**Sprint 周期**：2 周（2026-08-30 ～ 2026-09-12）

---

## 一、Sprint 目标

**目标：** 具身智能核心能力

在 Sprint 19-20（数字孪生）的基础上，实现具身智能核心能力，包括环境感知、空间认知、自主探索、动作模仿，构建宠物机器人在物理世界中的智能移动和交互能力。

---

## 二、详细任务列表

### 后端 P0/P1/P2 任务表

| # | 任务 | 说明 | 交付物 | 优先级 |
|---|------|------|--------|--------|
| P0-1 | **环境感知 API** | 完成 `/api/v1/embodied/:device_id/perception` 环境感知接口 | perception_controller.go | P0 |
| P0-2 | **空间认知 API** | 完成 `/api/v1/embodied/:device_id/map` 地图构建与定位接口 | spatial_controller.go | P0 |
| P0-3 | **自主导航 API** | 完成 `/api/v1/embodied/:device_id/navigate` 自主导航接口 | navigation_controller.go | P0 |
| P0-4 | **具身决策引擎 API** | 完成 `/api/v1/embodied/:device_id/decide` 感知-决策-执行闭环 | decision_controller.go | P0 |
| P1-1 | **动作模仿 API** | 完成 `/api/v1/embodied/:device_id/actions/learn` 动作学习接口 | action_controller.go | P1 |
| P1-2 | **安全边界 API** | 完成 `/api/v1/embodied/:device_id/safety-zones` 安全边界配置 | safety_controller.go | P1 |
| P2-1 | **感知融合服务** | 实现视觉/距离/触觉多传感器融合逻辑 | perception/fusion.go | P2 |
| P2-2 | **地图语义化服务** | 实现语义地图构建和更新 | spatial/semantic_map.go | P2 |

### 前端 P0/P1/P2 任务表

| # | 任务 | 说明 | 交付物 | 优先级 |
|---|------|------|--------|--------|
| PF0-1 | **环境感知面板** | 完成 EnvironmentalPerceptionView.vue 感知数据展示 | EnvironmentalPerceptionView.vue | P0 |
| PF0-2 | **空间地图页面** | 完成 SpatialMapView.vue 室内地图和定位展示 | SpatialMapView.vue | P0 |
| PF0-3 | **导航控制页面** | 完成 NavigationControlView.vue 自主导航目标设置 | NavigationControlView.vue | P0 |
| PF0-4 | **具身决策日志** | 完成 EmbodiedDecisionLogView.vue 决策过程可视化 | EmbodiedDecisionLogView.vue | P0 |
| PF1-1 | **动作模仿学习页面** | 完成 ActionLearningView.vue 动作学习记录和进度 | ActionLearningView.vue | P1 |
| PF1-2 | **安全边界配置页面** | 完成 SafetyZonesView.vue 安全区域配置可视化 | SafetyZonesView.vue | P1 |
| PF2-1 | **感知数据统计** | 完成 PerceptionStatsView.vue 感知准确率统计 | PerceptionStatsView.vue | P2 |

---

## 三、技术方案

### API 路由设计

| 接口 | 方法 | 说明 |
|------|------|------|
| `GET /api/v1/embodied/:device_id/perception` | GET | 获取当前环境感知数据 |
| `GET /api/v1/embodied/:device_id/perception/objects` | GET | 识别到的物体列表 |
| `GET /api/v1/embodied/:device_id/perception/depth` | GET | 深度/距离感知数据 |
| `GET /api/v1/embodied/:device_id/map` | GET | 获取当前地图 |
| `GET /api/v1/embodied/:device_id/map/position` | GET | 当前定位 |
| `POST /api/v1/embodied/:device_id/map/update` | POST | 更新地图 |
| `POST /api/v1/embodied/:device_id/navigate` | POST | 下发导航目标 |
| `GET /api/v1/embodied/:device_id/navigate/status` | GET | 导航状态 |
| `POST /api/v1/embodied/:device_id/navigate/cancel` | POST | 取消导航 |
| `GET /api/v1/embodied/:device_id/decide` | GET | 获取决策状态 |
| `GET /api/v1/embodied/:device_id/decide/history` | GET | 决策历史 |
| `POST /api/v1/embodied/:device_id/actions/learn` | POST | 学习新动作 |
| `GET /api/v1/embodied/:device_id/actions` | GET | 动作库列表 |
| `GET /api/v1/embodied/:device_id/safety-zones` | GET | 安全边界列表 |
| `POST /api/v1/embodied/:device_id/safety-zones` | POST | 创建安全区域 |
| `PUT /api/v1/embodied/:device_id/safety-zones/:id` | PUT | 更新安全区域 |
| `DELETE /api/v1/embodied/:device_id/safety-zones/:id` | DELETE | 删除安全区域 |

### 数据库设计

```sql
-- 环境感知表
CREATE TABLE embodied_perceptions (
    id              BIGSERIAL PRIMARY KEY,
    device_id       VARCHAR(100) NOT NULL,
    perception_type VARCHAR(50) NOT NULL,           -- 'visual'/'depth'/'tactile'/'fusion'
    data            JSONB NOT NULL,                   -- 感知数据
    confidence      DECIMAL(5,4),                    -- 置信度 0-1
    detected_at     TIMESTAMP NOT NULL,
    created_at      TIMESTAMP DEFAULT NOW(),
    INDEX idx_device_perception (device_id, detected_at DESC)
);

-- 空间地图表
CREATE TABLE embodied_maps (
    id              BIGSERIAL PRIMARY KEY,
    device_id       VARCHAR(100) NOT NULL UNIQUE,
    map_type        VARCHAR(20) NOT NULL,            -- 'grid'/'semantic'
    map_data        JSONB NOT NULL,                   -- 地图数据
    resolution      DECIMAL(5,2),                    -- 分辨率 cm/格
    size_width      INT,
    size_height     INT,
    version         INT DEFAULT 1,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

-- 空间位置表
CREATE TABLE spatial_positions (
    id              BIGSERIAL PRIMARY KEY,
    device_id       VARCHAR(100) NOT NULL,
    position_x      DECIMAL(10,4),
    position_y      DECIMAL(10,4),
    position_z      DECIMAL(10,4),
    orientation     DECIMAL(6,2),                    -- 朝向角度 0-360
    map_id          BIGINT REFERENCES embodied_maps(id),
    located_at      TIMESTAMP NOT NULL,
    accuracy        DECIMAL(6,2),                    -- 定位精度 cm
    created_at      TIMESTAMP DEFAULT NOW(),
    INDEX idx_device_position (device_id, located_at DESC)
);

-- 动作库表
CREATE TABLE embodied_action_library (
    id              BIGSERIAL PRIMARY KEY,
    device_id       VARCHAR(100) NOT NULL,
    action_name     VARCHAR(100) NOT NULL,
    action_type     VARCHAR(50),                    -- 'movement'/'expression'/'interaction'
    action_data     JSONB NOT NULL,                   -- 动作定义
    learned_from    VARCHAR(50),                    -- 'human_demo'/'auto_learn'/'preset'
    confidence      DECIMAL(5,4),
    execution_count INT DEFAULT 0,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW(),
    INDEX idx_device_actions (device_id)
);

-- 安全边界表
CREATE TABLE safety_zones (
    id              BIGSERIAL PRIMARY KEY,
    device_id       VARCHAR(100) NOT NULL,
    zone_name       VARCHAR(100) NOT NULL,
    zone_type       VARCHAR(30) NOT NULL,            -- 'restricted'/'allowed'/'danger'
    boundary        JSONB NOT NULL,                   -- 边界坐标
    priority        INT DEFAULT 0,
    is_active       BOOLEAN DEFAULT TRUE,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW(),
    INDEX idx_device_safety (device_id)
);

-- 具身决策日志表
CREATE TABLE embodied_decision_logs (
    id              BIGSERIAL PRIMARY KEY,
    device_id       VARCHAR(100) NOT NULL,
    decision_id     VARCHAR(100) NOT NULL,
    perception_data JSONB,
    decision_result JSONB,
    executed_action JSONB,
    safety_triggered BOOLEAN DEFAULT FALSE,
    latency_ms      INT,
    decided_at      TIMESTAMP NOT NULL,
    created_at      TIMESTAMP DEFAULT NOW(),
    INDEX idx_device_decisions (device_id, decided_at DESC)
);

-- 安全日志表
CREATE TABLE embodied_safety_logs (
    id              BIGSERIAL PRIMARY KEY,
    device_id       VARCHAR(100) NOT NULL,
    trigger_type    VARCHAR(50) NOT NULL,            -- 'boundary'/'collision'/'unauthorized'
    zone_id         BIGINT,
    action_taken    VARCHAR(100),
    severity        VARCHAR(20),                    -- 'warning'/'critical'
    triggered_at    TIMESTAMP NOT NULL,
    created_at      TIMESTAMP DEFAULT NOW(),
    INDEX idx_device_safety (device_id, triggered_at DESC)
);
```

---

## 四、验收标准

### 4.1 功能验收

| 用例 | 验收条件 | 测试方法 |
|------|----------|----------|
| 环境感知 | 物体识别准确率>90%，障碍物检测准确率>95% | 标准化测试场景 |
| 空间认知 | 室内定位精度<10cm，地图构建覆盖95%可活动区域 | 实地测试 |
| 自主导航 | 导航成功率>85%，避障成功率>90% | 迷宫测试 |
| 动作模仿 | 动作识别准确率>80%，动作执行延迟<1s | 动作示范测试 |
| 安全边界 | 危险区域触发准确率100%，响应延迟<100ms | 安全边界测试 |

### 4.2 性能验收

| 验收点 | 标准 |
|--------|------|
| 感知数据处理延迟 | < 200ms |
| 决策延迟 | < 1s |
| 地图更新延迟 | < 500ms |
| 安全响应延迟 | < 100ms |

---

## 五、依赖与风险

### 依赖

| 依赖 | 说明 |
|------|------|
| Sprint 19-20 数字孪生 | 设备状态数据和传感器数据 |
| MODULE_EMBODIED_AI | 具身智能完整 PRD |
| 设备端固件支持 | 具身感知和控制指令支持 |

### 风险

| 风险 | 影响 | 缓解措施 |
|------|------|----------|
| 视觉感知准确性不足 | 物体识别失败 | 增加训练数据集 |
| 地图更新不及时 | 定位漂移 | 定期自动校准 |
| 安全边界漏检 | 可能碰撞受伤 | 多传感器冗余 |
