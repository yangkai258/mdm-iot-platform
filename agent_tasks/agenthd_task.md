Agent HD - 后端开发任务

## P0 修复任务（已完成 ✅）

| # | 问题 | 修复文件 | 修复内容 |
|---|------|----------|----------|
| 1 | JWT密钥硬编码 | `middleware/jwt.go` | 改为从环境变量 `JWT_SECRET` 读取，缺失时 panic |
| 2 | CORS全开放 | `main.go` | 白名单仅允许 `localhost:3000` 和 `127.0.0.1:3000`，支持 `CORS_ALLOWED_ORIGINS` 环境变量扩展 |
| 3 | MQTT注入为nil | `main.go` + `mqtt/handler.go` + `controllers/command_controller.go` | 添加 `GlobalMQTTClient` 全局变量，`CommandController` 改用 `mqtt.GlobalMQTTClient` |
| 4 | OTA无后台Worker | `main.go` | 添加 `startOTAWorker` goroutine，每5分钟检查 `pending` 部署并通过 MQTT 下发 OTA 指令 |
| 5 | CheckAlerts从未调用 | `mqtt/handler.go` + `main.go` | 通过 `AlertCallback` 接口将 `controllers.CheckAlerts` 注入到 MQTT 心跳处理中 |
| 6 | deleteRecord不删数据 | `controllers/device_crud.go` | GORM 软删除实现正确（`DeletedAt` 字段已配置），`db.Delete(&device)` 会设置 `DeletedAt` 时间戳，查询自动过滤 |

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

