Agent HD - 后端开发任务

## Sprint 3.3: 告警系统增强（已完成 ✅）

| 验收标准 | 状态 | 实现文件 |
|----------|------|----------|
| 越狱检测告警触发 | ✅ | `mqtt/handler.go` - `JailbreakAlertHandler` + `StatusPayload` 增加 `is_jailbroken/root_status` 字段 |
| 地理围栏违规告警 | ✅ | `mqtt/handler.go` - `CheckGeofence` + `GeofenceRule` 模型 + Haversine 距离算法 |
| 告警处理流程完整 | ✅ | `controllers/alert_controller.go` - confirm/resolve/ignore/batchConfirm/batchResolve |

**新增文件:**
- `models/alert_models.go` - 新增 `GeofenceRule`, `GeofenceAlert`, `AlertNotification` 模型；`DeviceAlert` 增加 `ConfirmedAt/ConfirmedBy/ResolvedAt/ResolvedBy/IgnoredAt/IgnoredBy/ExtraData` 字段
- `services/notification_service.go` - 告警通知服务，支持 email/webhook/inapp 三种通知方式

**修改文件:**
- `models/device.go` - `DeviceShadow` 增加 `IsJailbroken`, `RootStatus`, `Latitude`, `Longitude` 字段
- `utils/redis.go` - `DeviceShadow` 结构体对应增加上述字段
- `mqtt/handler.go` - `InitMQTT` 增加 `GeofenceCallback`；`SetupSubscriber` 订阅 `jailbreak_alert` 主题；`StatusMessageHandler` 处理越狱和位置数据；新增 `JailbreakAlertHandler`、`CheckGeofence`、`haversineDistance`
- `controllers/alert_controller.go` - 完整重写，新增所有告警处理 API；`CheckAlerts` 支持 `jailbreak_detected` 类型
- `main.go` - 注册 `GeofenceRule`、`GeofenceAlert`、`AlertNotification` 模型迁移；`InitMQTT` 传入 `geofenceCallback`；注册新路由

**API路由（新增）:**
- `POST /api/v1/alerts/:id/confirm` - 确认告警
- `POST /api/v1/alerts/:id/resolve` - 解决告警
- `POST /api/v1/alerts/:id/ignore` - 忽略告警
- `POST /api/v1/alerts/batch/confirm` - 批量确认
- `POST /api/v1/alerts/batch/resolve` - 批量解决
- `GET /api/v1/alerts/:id/notifications` - 获取告警通知记录
- `PUT /api/v1/alerts/rules/:id` - 更新告警规则
- `DELETE /api/v1/alerts/rules/:id` - 删除告警规则
- `GET /api/v1/geofence/rules` - 地理围栏规则列表
- `POST /api/v1/geofence/rules` - 创建地理围栏规则
- `PUT /api/v1/geofence/rules/:id` - 更新地理围栏规则
- `DELETE /api/v1/geofence/rules/:id` - 删除地理围栏规则
- `GET /api/v1/geofence/alerts` - 地理围栏告警记录

**MQTT主题（新增）:**
- `/mdm/device/+/up/jailbreak_alert` - 订阅设备越狱/ROOT告警上报

---

## Sprint 2.2: 通知管理（已完成 ✅）

| 验收标准 | 状态 | 实现文件 |
|----------|------|----------|
| 通知通过MQTT下发到设备 | ✅ | `POST /devices/:device_id/notifications` 发布到 `/device/{id}/down/notification` |
| 通知模板支持变量替换 | ✅ | `ReplaceTemplateVariables()` 支持 `{{variable}}` 格式替换 |
| 公告CRUD完整 | ✅ | GET/POST/PUT/DELETE /api/v1/announcements + publish |

**新增文件:**
- `models/notification.go` - Notification, NotificationTemplate, Announcement
- `controllers/notification_controller.go` - NotificationController 实现所有接口

**修改文件:**
- `controllers/device_controller.go` - 注册通知管理路由
- `main.go` - 添加 Notification/NotificationTemplate/Announcement 自动迁移

**API路由:**
- `GET /api/v1/notifications` - 通知列表
- `GET /api/v1/notifications/:id` - 通知详情
- `DELETE /api/v1/notifications/:id` - 删除通知
- `POST /api/v1/devices/:device_id/notifications` - 发送通知（MQTT下发）
- `GET /api/v1/notification-templates` - 模板列表
- `POST /api/v1/notification-templates` - 创建模板
- `PUT /api/v1/notification-templates/:id` - 更新模板
- `DELETE /api/v1/notification-templates/:id` - 删除模板
- `GET /api/v1/announcements` - 公告列表
- `POST /api/v1/announcements` - 创建公告
- `PUT /api/v1/announcements/:id` - 更新公告
- `DELETE /api/v1/announcements/:id` - 删除公告
- `POST /api/v1/announcements/:id/publish` - 发布公告

---

## Sprint 2.1: 应用管理基础（已完成 ✅）

| 验收标准 | 状态 | 实现文件 |
|----------|------|----------|
| 应用CRUD完整 | ✅ | `models/app.go` + `controllers/app_controller.go` |
| 版本管理完整 | ✅ | `ListVersions`, `CreateVersion`, `DeleteVersion` |
| 分发任务支持设备/用户/组 | ✅ | `CreateDistribution`, `GetDistribution`, `CancelDistribution` |
| 安装统计接口 | ✅ | `GetStats` |

**新增文件:**
- `models/app.go` - App, AppVersion, AppDistribution, AppInstallRecord, AppLicense
- `controllers/app_controller.go` - AppController 实现所有接口

**修改文件:**
- `controllers/device_controller.go` - 注册应用管理路由
- `main.go` - 添加 App 相关模型自动迁移

**API路由:**
- `GET /api/v1/apps` - 应用列表
- `POST /api/v1/apps` - 创建应用
- `GET /api/v1/apps/:id` - 应用详情
- `PUT /api/v1/apps/:id` - 更新应用
- `DELETE /api/v1/apps/:id` - 删除应用
- `GET /api/v1/apps/:id/versions` - 版本列表
- `POST /api/v1/apps/:id/versions` - 添加版本
- `DELETE /api/v1/apps/:id/versions/:version_id` - 删除版本
- `POST /api/v1/app/distributions` - 创建分发任务
- `GET /api/v1/app/distributions/:id` - 分发详情
- `POST /api/v1/app/distributions/:id/cancel` - 取消分发
- `GET /api/v1/apps/:id/stats` - 安装统计

---

## Sprint 1.1: OTA Worker 实现（已完成 ✅）

| 验收标准 | 状态 | 实现文件 |
|----------|------|----------|
| OTA Worker 每5分钟轮询一次 | ✅ | `services/ota_worker.go` |
| 支持全量/百分比/白名单灰度策略 | ✅ | `SelectTargetDevices()` 方法 |
| 设备可以通过 MQTT 接收 OTA 指令 | ✅ | `PublishOTACommand()` 发布到 `/device/{device_id}/down/cmd` |
| 设备上报进度后更新数据库 | ✅ | `DeviceOTAReport()` 回调接口 + `ota_progress` 表 |
| 成功率<80%自动暂停 | ✅ | `CheckAndAutoPause()` 失败率阈值检查 |

**新增文件:**
- `services/ota_worker.go` - OTA后台Worker服务

**修改文件:**
- `controllers/ota_controller.go` - 新增 `SetOTAWorkerRef`, `DeviceOTAReport`, `PauseDeployment`, `ResumeDeployment`, `CancelDeployment`, `GetDeploymentProgress`
- `controllers/device_controller.go` - 新增 OTA 部署管理路由
- `main.go` - 集成 `services.NewOTAWorker`

**API路由:**
- `POST /api/v1/ota/devices/:device_id/report` - 设备回调上报进度
- `GET /api/v1/ota/deployments` - 部署列表
- `GET /api/v1/ota/deployments/:id` - 部署详情
- `POST /api/v1/ota/deployments/:id/pause` - 暂停部署
- `POST /api/v1/ota/deployments/:id/resume` - 恢复部署
- `POST /api/v1/ota/deployments/:id/cancel` - 取消部署
- `GET /api/v1/ota/deployments/:id/progress` - 部署进度详情

**提交:** `dafa4e3` Sprint 1.1: OTA Worker 实现

---

## P0 修复任务（已完成 ✅）

| # | 问题 | 修复文件 | 修复内容 |
|---|------|----------|----------|
| 1 | JWT密钥硬编码 | `middleware/jwt.go` | 改为从环境变量 `JWT_SECRET` 读取，缺失时 panic |
| 2 | CORS全开放 | `main.go` | 白名单仅允许 `localhost:3000` 和 `127.0.0.1:3000`，支持 `CORS_ALLOWED_ORIGINS` 环境变量扩展 |
| 3 | MQTT注入为nil | `main.go` + `mqtt/handler.go` + `controllers/command_controller.go` | 添加 `GlobalMQTTClient` 全局变量，`CommandController` 改用 `mqtt.GlobalMQTTClient` |
| 4 | OTA无后台Worker | `main.go` | 添加 `startOTAWorker` goroutine，每5分钟检查 `pending` 部署并通过 MQTT 下发 OTA 指令 |
| 5 | CheckAlerts从未调用 | `mqtt/handler.go` + `main.go` | 通过 `AlertCallback` 接口将 `controllers.CheckAlerts` 注入到 MQTT 心跳处理中 |
| 6 | deleteRecord不删数据 | `controllers/device_crud.go` | GORM 软删除实现正确（`DeletedAt` 字段已配置），`db.Delete(&device)` 会设置 `DeletedAt` 时间戳，查询自动过滤 |
| 7 | fmt.Sscanf解析百分比 | `services/ota_worker.go` | 替换 `fmt.Sscanf` 为 `strconv.Atoi`，避免解析失败时 percentage 保持初始值导致意外行为 |

---

## 历史任务（原始需求）

你是资深 Golang 后端开发工程师 (agenthd)

你的任务是使用 Gin 框架和 GORM 编写高并发的设备管理核心逻辑。

数据库已使用 PostgreSQL

核心表有 mdm\_devices (设备台账) 和基于 Redis 的设备影子。

任务 1：数据模型定义 (models/device.go)

使用 GORM tag 严格定义 Device 结构体。字段包括：

* ID
* DeviceID (UUID)
* MacAddress
* SnCode
* HardwareModel
* FirmwareVersion
* BindUserID
* LifecycleStatus
* CreatedAt
* UpdatedAt

任务 2：API 路由与控制器 (controllers/device\_controller.go)

实现设备注册/绑定的处理函数。要求：

* 接收 JSON
* 校验必填字段（如 MacAddress）
* 如果 Mac 存在则更新绑定，不存在则新建
* 最终存入 PostgreSQL

任务 3：MQTT 消息桥接框架 (mqtt/handler.go)

使用 paho.mqtt.golang 库，写一个订阅设备上报 Topic 的回调函数。当收到心跳 JSON 时，将其解析并更新到 Redis 中（Key 为 shadow:{device\_id}，设置 90 秒 TTL）。

约束

* 代码必须包含充分的错误处理 (Error Handling) 和日志打印
* 不要给出项目初始化步骤，只输出核心源码

