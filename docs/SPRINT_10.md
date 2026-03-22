# Sprint 10 规划

**时间**：2026-04-05
**状态**：✅ 前端已完成 (2026-04-10)
**Sprint 周期**：2 周（2026-04-05 ～ 2026-04-18）

---

## 一、Sprint 目标

**目标：** 完善设备监控面板、传感器事件处理、动作库管理

在 Sprint 9（OpenClaw AI层核心功能）的基础上，完善设备监控面板、传感器事件处理、动作库管理功能，提供完整的设备运维监控能力，支持设备批量操作和远程调试。

---

## 二、详细任务列表

### 后端 P0/P1/P2 任务表

| # | 任务 | 说明 | 交付物 | 优先级 |
|---|------|------|--------|--------|
| P0-1 | **传感器事件处理 API** | 完成 `POST /api/v1/sensors/events` 接收设备传感器数据 | sensor_controller.go | P0 |
| P0-2 | **传感器事件存储** | 创建 sensor_events 表存储传感器原始数据 | models/sensor_event.go | P0 |
| P0-3 | **动作库管理 API (CRUD)** | 完成 `/api/v1/action-library/*` 完整 CRUD | action_library_controller.go | P0 |
| P0-4 | **批量设备操作 API** | 完成 `POST /api/v1/devices/batch-actions` 批量下发指令 | device_controller.go | P0 |
| P0-5 | **设备监控指标 API** | 完成 `GET /api/v1/monitoring/metrics` 设备运行指标 | monitoring_controller.go | P0 |
| P1-1 | **告警规则引擎完善** | 完善 `CheckAlerts` 函数，支持自定义规则表达式 | alert_engine.go | P1 |
| P1-2 | **设备日志 API** | 完成 `GET /api/v1/devices/{device_id}/logs` 设备日志查询 | device_log_controller.go | P1 |
| P1-3 | **远程调试命令 API** | 完成 `POST /api/v1/devices/{device_id}/debug/command` | debug_controller.go | P1 |
| P1-4 | **传感器数据聚合 API** | 完成 `GET /api/v1/sensors/{device_id}/aggregations` 聚合统计 | sensor_controller.go | P1 |
| P2-1 | **告警规则模板** | 完成 `GET/POST /api/v1/alerts/templates` 规则模板管理 | alert_template_controller.go | P2 |
| P2-2 | **操作审计日志 API** | 完成 `GET /api/v1/audit/operations` 操作审计 | audit_controller.go | P2 |

### 前端 P0/P1/P2 任务表

| # | 任务 | 说明 | 交付物 | 优先级 |
|---|------|------|--------|--------|
| PF0-1 | **设备监控面板主页面** | 完成 DeviceMonitorView.vue 主布局和监控卡片 | DeviceMonitorView.vue | P0 |
| PF0-2 | **设备状态列表组件** | 完成 DeviceStatusList.vue 设备在线/离线状态 | DeviceStatusList.vue | P0 |
| PF0-3 | **设备日志查看页面** | 完成 DeviceLogsView.vue 日志查询和展示 | DeviceLogsView.vue | P0 |
| PF0-4 | **动作库管理页面** | 完成 ActionLibraryView.vue 动作库列表/搜索/编辑 | ActionLibraryView.vue | P0 |
| PF0-5 | **批量操作组件** | 完成 BatchActionModal.vue 批量操作确认弹窗 | BatchActionModal.vue | P0 |
| PF1-1 | **远程调试控制台** | 完成 RemoteDebugView.vue 远程命令下发和日志回显 | RemoteDebugView.vue | P1 |
| PF1-2 | **传感器数据图表** | 完成 SensorChart.vue 传感器数据趋势图 | SensorChart.vue | P1 |
| PF1-3 | **设备监控详情页** | 完成 DeviceMonitorDetail.vue 单设备详情监控 | DeviceMonitorDetail.vue | P1 |
| PF2-1 | **告警规则配置弹窗** | 完成 AlertRuleModal.vue 告警规则创建/编辑 | AlertRuleModal.vue | P2 |
| PF2-2 | **操作审计页面** | 完成 AuditLogView.vue 审计日志查看 | AuditLogView.vue | P2 |

---

## 三、技术方案

### API 路由设计

| 接口 | 方法 | 说明 |
|------|------|------|
| `POST /api/v1/sensors/events` | POST | 接收设备传感器事件 |
| `GET /api/v1/sensors/{device_id}/data` | GET | 获取设备传感器数据 |
| `GET /api/v1/sensors/{device_id}/aggregations` | GET | 传感器数据聚合统计 |
| `GET /api/v1/action-library` | GET | 动作库列表 |
| `POST /api/v1/action-library` | POST | 创建动作 |
| `GET /api/v1/action-library/:id` | GET | 动作详情 |
| `PUT /api/v1/action-library/:id` | PUT | 更新动作 |
| `DELETE /api/v1/action-library/:id` | DELETE | 删除动作 |
| `POST /api/v1/devices/batch-actions` | POST | 批量设备操作 |
| `GET /api/v1/monitoring/metrics` | GET | 设备监控指标 |
| `GET /api/v1/devices/{device_id}/logs` | GET | 设备日志 |
| `POST /api/v1/devices/{device_id}/debug/command` | POST | 远程调试命令 |
| `GET /api/v1/alerts/templates` | GET | 告警规则模板列表 |
| `POST /api/v1/alerts/templates` | POST | 创建告警规则模板 |
| `GET /api/v1/audit/operations` | GET | 操作审计日志 |

### 数据库设计

```sql
-- 传感器事件表
CREATE TABLE sensor_events (
    id              BIGSERIAL PRIMARY KEY,
    device_id       VARCHAR(64) NOT NULL,
    sensor_type     VARCHAR(50) NOT NULL,
    value           JSONB NOT NULL,
    unit            VARCHAR(20),
    quality         VARCHAR(20),
    timestamp       TIMESTAMP NOT NULL,
    created_at      TIMESTAMP DEFAULT NOW(),
    INDEX idx_device_sensor_time (device_id, sensor_type, timestamp DESC)
);

-- 设备操作日志表
CREATE TABLE device_operation_logs (
    id              BIGSERIAL PRIMARY KEY,
    device_id       VARCHAR(64) NOT NULL,
    operator_id     BIGINT NOT NULL,
    operator_type   VARCHAR(20),
    operation_type  VARCHAR(50) NOT NULL,
    operation_data  JSONB,
    result          VARCHAR(20),
    error_message   TEXT,
    ip_address      VARCHAR(45),
    created_at      TIMESTAMP DEFAULT NOW(),
    INDEX idx_device_id (device_id),
    INDEX idx_operator_id (operator_id),
    INDEX idx_created_at (created_at DESC)
);

-- 告警规则模板表
CREATE TABLE alert_templates (
    id              BIGSERIAL PRIMARY KEY,
    template_name   VARCHAR(100) NOT NULL,
    alert_type      VARCHAR(50) NOT NULL,
    condition_expr  TEXT NOT NULL,
    threshold_value VARCHAR(100),
    severity        INT DEFAULT 2,
    notify_ways     VARCHAR(50)[],
    remark          TEXT,
    created_by      BIGINT,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

-- 操作审计表
CREATE TABLE audit_logs (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL,
    action          VARCHAR(50) NOT NULL,
    resource_type   VARCHAR(50),
    resource_id     VARCHAR(100),
    request_data    JSONB,
    response_data   JSONB,
    ip_address      VARCHAR(45),
    user_agent      VARCHAR(255),
    created_at      TIMESTAMP DEFAULT NOW(),
    INDEX idx_user_id (user_id),
    INDEX idx_action (action),
    INDEX idx_created_at (created_at DESC)
);
```

---

## 四、验收标准

### 4.1 功能验收

| 用例 | 验收条件 | 测试方法 |
|------|----------|----------|
| 传感器事件接收 | 上报事件后100ms内存储成功 | 调用API验证 |
| 动作库CRUD | 完整增删改查 | 调用各接口验证 |
| 批量设备操作 | 一次操作最多100台设备 | 调用批量接口 |
| 设备监控指标 | 返回设备在线数/离线数/告警数 | 调用API验证 |
| 设备日志查询 | 支持按时间/类型筛选 | 分页+筛选测试 |
| 远程调试 | 命令下发后设备响应 | 实机测试 |
| 告警规则引擎 | 规则表达式正确匹配 | 模拟数据测试 |

### 4.2 性能验收

| 验收点 | 标准 |
|--------|------|
| 传感器事件处理吞吐量 | >= 1000 events/s |
| 批量操作响应时间 | <= 5s（100台设备） |
| 监控指标查询延迟 | <= 500ms |

---

## 五、依赖与风险

### 依赖

| 依赖 | 说明 |
|------|------|
| Sprint 9 设备基础 | 设备注册/心跳机制 |
| Sprint 8 数据权限 | Repository层权限过滤 |
| EMQX MQTT Broker | MQTT消息接收 |

### 风险

| 风险 | 影响 | 缓解措施 |
|------|------|----------|
| 传感器数据量过大 | DB存储压力 | 分表+TTL过期 |
| 批量操作超时 | 用户体验差 | 增加超时配置+进度反馈 |

---

## 六、前端完成清单

### ✅ 阶段1: 设备监控面板
- [x] `views/monitor/DeviceDashboard.vue` - 设备监控面板主页面
  - 统计卡片（在线/离线/告警/平均电量）
  - SVG折线图（CPU使用率趋势、内存使用趋势）
  - 电量分布图（电池条形图）
  - 网络状态列表
  - 设备列表（在线状态、在线时长、最后活跃）
  - 筛选工具栏（设备选择、时间范围、刷新频率）
  - 导出报表功能

### ✅ 阶段2: 设备日志前端
- [x] `views/monitor/DeviceLogs.vue` - 设备日志查看页面
  - 日志级别筛选（info/warn/error）
  - 时间范围筛选
  - 关键词搜索
  - 日志列表（时间、设备、内容、级别）
  - 导出日志功能
  - 分页支持

### ✅ 阶段3: 远程调试前端
- [x] `views/monitor/RemoteDebug.vue` - 远程调试控制台
  - 设备选择
  - 终端输出区（黑色背景，等宽字体）
  - 命令输入区
  - 快捷命令按钮（重启设备、查看日志、抓取堆栈、性能分析）
  - 设备信息展示

### ✅ 阶段4: 动作库管理前端
- [x] `views/action/ActionLibrary.vue` - 动作库管理页面
  - 动作列表（分类展示、卡片布局）
  - 动作详情弹窗（参数配置预览）
  - 创建/编辑动作表单
  - 分类统计卡片

### ✅ 阶段5: API 层
- [x] `api/monitor.ts` - 设备监控 API
  - 设备监控 API (monitoring/metrics, device/realtime, device/history, device/alerts)
  - 设备日志 API (device-logs)
  - 动作库 CRUD API
  - 批量操作 API
  - 远程调试 API
- [x] `composables/useDeviceMonitor.ts` - 设备监控 Hook
- [x] `composables/useDeviceLogs.ts` - 设备日志 Hook

### ✅ 阶段6: 路由配置
- [x] `router/index.js` - 添加 Sprint 10 路由
  - `/monitor/dashboard` → DeviceDashboard
  - `/monitor/logs` → DeviceLogs
  - `/monitor/debug` → RemoteDebug
  - `/action/library` → ActionLibrary

### 📝 Git 提交
- Commit: `b4837e8` - `feat(frontend): Sprint 10 - All stages (DeviceDashboard/DeviceLogs/RemoteDebug/ActionLibrary)`
- Branch: `master`
- Pushed: ✅
