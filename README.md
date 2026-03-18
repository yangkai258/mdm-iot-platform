# MDM 物联网设备管理中台

> Mobile Device Management Platform for M5Stack IoT Devices

## 项目概述

MDM (Mobile Device Management) 中台是整个"端云协同"架构的核心骨干（部署于 Mac mini 或云端服务器）。它负责统一管理分布在全球的 M5Stack 硬件终端，维护设备的生命周期（注册、在线、离线、升级、报废），并作为 OpenClaw 多智能体中枢与物理设备之间的"状态缓冲池"与"指令路由器"。

## 核心模块

### 1. 设备资产与生命周期管理 (Device Asset Management)
- 设备台账列表（UUID、MAC地址、SN码、硬件型号等）
- 用户绑定与解绑接口
- 设备状态追踪：待激活、正常服役、维修中、已挂失、已报废

### 2. 设备影子与状态监控 (Device Shadow & Monitoring)
- 实时状态机（电量、在线状态、当前模式）
- 期望状态控制（就寝时间、目标固件版本）
- 解决网络高延迟或断网问题

### 3. OTA 固件静默升级中心 (Over-The-Air Updates)
- 固件包管理（版本、MD5校验、强制升级）
- 发布任务与灰度管控（百分比灰度、白名单）
- 成功率监控与自动暂停

### 4. 宠物性格与管控策略 (Pet Configuration & Rules)
- 宠物系统设定（名字、性格、互动频率）
- 免打扰时间段配置
- 自定义规则集

## 技术栈

| 层级 | 技术 |
|------|------|
| 后端核心 | Golang + Gin 框架 |
| 前端后台 | Vue3 + TypeScript + Arco Design |
| 消息中间件 | EMQX (MQTT Broker) |
| 主数据库 | PostgreSQL 15+ |
| 状态缓存 | Redis 7.0 (设备影子) |
| 数据交换 | JSON |

## 快速开始

### 前置要求

- Docker & Docker Compose
- PostgreSQL 15+
- Redis 7+
- EMQX 5.0
- Go 1.21+ (开发)
- Node.js 18+ (开发)

### 1. 克隆项目

```bash
git clone https://github.com/yangkai258/mdm-iot-platform.git
cd mdm-iot-platform
```

### 2. 一键启动（推荐）

```bash
docker-compose up -d
```

服务启动后：
- **MDM Web 控制台**: http://localhost
- **EMQX Dashboard**: http://localhost:18083 (admin/public)

### 3. 手动启动

```bash
# 启动 PostgreSQL
docker run -d --name mdm_postgres -e POSTGRES_USER=mdm_user -e POSTGRES_PASSWORD=mdm_password -e POSTGRES_DB=mdm_db -p 5432:5432 postgres:15-alpine

# 启动 Redis
docker run -d --name mdm_redis -p 6379:6379 redis:7-alpine

# 启动 EMQX
docker run -d --name mdm_emqx -p 1883:1883 -p 8083:8083 -p 18083:18083 emqx/emqx:5.0-alpine

# 启动后端服务
cd backend
go run main.go

# 启动前端
cd frontend
npm install
npm run dev
```

## 项目结构

```
mdm-iot-platform/
├── backend/                    # Golang 后端
│   ├── models/               # 数据模型 (GORM)
│   ├── controllers/          # API 控制器
│   ├── mqtt/               # MQTT 消息处理
│   └── utils/              # 工具类
│
├── frontend/                 # Vue3 前端
│   └── src/views/          # 页面组件
│
├── agent_tasks/              # Agent 任务定义
│   ├── agentcp_task.md     # 产品经理
│   ├── agenthd_task.md     # 后端开发
│   ├── agentqd_task.md     # 前端开发
│   ├── agentcs_task.md     # 测试工程师
│   └── agentyw_task.md    # 运维工程师
│
├── test_scripts/             # 测试脚本
│   └── mqtt_stress_test.py  # MQTT 压测脚本
│
├── docker-compose.yml        # 一键部署配置
├── PRD_MDM_Control_Center.md  # 产品需求文档
├── Tech_Stack_and_DB_Design.md # 技术栈与数据库设计
└── MDM_API_Contract.md       # API 接口规范
```

## API 接口

| 接口 | 方法 | 说明 |
|------|------|------|
| `/api/v1/devices/register` | POST | 设备注册 |
| `/api/v1/devices/{sn_code}/bind` | POST | 用户绑定设备 |
| `/api/v1/devices/list` | GET | 设备列表（分页+筛选） |
| `/api/v1/devices/{device_id}` | GET | 设备详情 |
| `/api/v1/devices/{device_id}/status` | GET | 设备状态 |
| `/api/v1/ota/check` | GET | 检查固件更新 |
| `/api/v1/ota/deploy` | POST | 创建发布任务 |
| `/api/v1/pets/{device_id}/config` | GET/PUT | 宠物配置 |

### MQTT Topic 规范

```
/mdm/device/{device_id}/up/status        # 设备心跳上报
/mdm/device/{device_id}/up/property       # 设备属性上报
/mdm/device/{device_id}/down/cmd         # 云端指令下发
/mdm/device/{device_id}/down/config       # 配置下发
```

## Agent 团队

项目采用多 Agent 协作开发模式：

| Agent | 角色 | 职责 |
|-------|------|------|
| agentcp | 产品经理 | 定义 API 契约、PRD |
| agenthd | 后端开发 | Golang API 实现 |
| agentqd | 前端开发 | Vue3 组件开发 |
| agentcs | 测试工程师 | 压测、功能测试 |
| agentyw | 运维工程师 | Docker 部署、CI/CD |

## 开发指南

### 添加新设备类型

1. 在 `backend/models/device.go` 添加硬件型号枚举
2. 更新 `MDM_API_Contract.md` 中的字段定义
3. 前端添加对应的 UI 组件

### 自定义 OTA 策略

修改 `mdm-ota-packages` 表或通过管理界面配置灰度发布规则。

### 添加新的设备指标

1. 修改 MQTT 消息处理逻辑 (`backend/mqtt/handler.go`)
2. 更新 Redis 设备影子结构
3. 前端添加对应的展示组件

## 测试

### MQTT 压测

```bash
pip install paho-mqtt
python test_scripts/mqtt_stress_test.py --host localhost --port 1883 --devices 1000
```

### API 测试

```bash
# 使用 curl 测试设备注册
curl -X POST http://localhost:8080/api/v1/devices/register \
  -H "Content-Type: application/json" \
  -d '{
    "device_id": "550e8400-e29b-41d4-a716-446655440000",
    "mac_address": "00:1A:2B:3C:4D:5E",
    "sn_code": "SN123456789",
    "hardware_model": "M5_CoreS3",
    "firmware_version": "v1.0.0"
  }'
```

## 性能指标

- API 响应时间: < 200ms
- 支持设备数量: 10,000+
- MQTT 心跳处理: 500次/秒
- WebSocket 并发: 1000+

## 许可证

MIT License
