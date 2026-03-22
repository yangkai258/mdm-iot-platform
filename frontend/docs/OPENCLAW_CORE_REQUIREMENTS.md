# OpenClaw 核心功能需求清单

**版本：** V1.1  
**编制日期：** 2026-03-22  
**状态：** 补充到 76 个功能点中

---

## 一、OpenClaw 核心功能（14个）

| # | 功能 | 优先级 | 说明 |
|---|------|--------|------|
| 1 | **动作序列生成 API** | P0 | POST /api/v1/behavior/generate |
| 2 | **决策规则管理** | P1 | CRUD + 优先级冲突处理 |
| 3 | **传感器事件处理** | P1 | 跌落/碰撞/触摸事件记录 |
| 4 | **动作库管理** | P1 | 动作的增删改查 |
| 5 | **宠物状态同步** | P0 | MQTT 上报状态，期望状态下发 |
| 6 | **OpenClaw AI 版本管理** | P0 | AI模型版本记录、兼容性矩阵 |
| 7 | **设备配对/注册流程** | P0 | 设备首次开机配对、AI授权绑定 |
| 8 | **设备监控面板** | P1 | 实时传感器数据、在线设备地图 |
| 9 | **设备日志系统** | P1 | 日志上报、存储、查询 |
| 10 | **远程调试接口** | P1 | 远程诊断、问题排查 |
| 11 | **批量操作** | P1 | 批量升级、批量指令下发 |
| 12 | **固件兼容性矩阵** | P0 | AI版本 vs 固件版本兼容性 |
| 13 | **告警规则引擎** | P1 | 基于传感器阈值触发告警 |
| 14 | **设备模拟器** | P3 | 开发测试用模拟器 |

---

## 二、OpenClaw 版本管理详细功能

### 2.1 AI 版本管理

| 功能 | 说明 |
|------|------|
| AI 模型版本记录 | OpenClaw AI 版本列表 |
| 版本与固件兼容性矩阵 | AI 版本 vs 固件版本兼容性 |
| 版本发布历史 | 每个版本的发布时间/说明 |
| 设备升级路径 | 设备从 v1 → v2 的升级路径 |

### 2.2 API 接口

```
GET  /api/v1/openclaw/versions           # 版本列表
GET  /api/v1/openclaw/versions/:id      # 版本详情
GET  /api/v1/openclaw/compatibility      # 兼容性矩阵
GET  /api/v1/openclaw/versions/available?device_id=xxx  # 设备可用版本
POST /api/v1/openclaw/versions          # 发布新版本
POST /api/v1/openclaw/versions/:id/disable  # 禁用版本
```

---

## 三、设备配对流程

### 3.1 配对流程

```
设备首次开机
    ↓
生成配对码（6位数字）
    ↓
用户在 App 输入配对码
    ↓
设备与用户账号绑定
    ↓
AI 授权（下载 OpenClaw 模型）
    ↓
配对完成
```

### 3.2 API 接口

```
POST /api/v1/devices/pairing/code     # 生成配对码
POST /api/v1/devices/pairing/verify   # 设备配对验证
POST /api/v1/devices/:id/unbind      # 解绑设备
GET  /api/v1/devices/pairing/history # 配对历史
```

---

## 四、设备影子（期望状态）

### 4.1 MQTT Topic

```
/miniclaw/{device_id}/down/desired_state  # 期望状态下发
/miniclaw/{device_id}/up/reported_state   # 实际状态上报
```

### 4.2 API 接口

```
PUT  /api/v1/devices/:device_id/desired-state   # 设置期望状态
GET  /api/v1/devices/:device_id/desired-state   # 获取期望状态
GET  /api/v1/devices/:device_id/reported-state  # 获取实际状态
GET  /api/v1/devices/:device_id/state-diff      # 状态差异
```

---

## 五、动作库管理

### 5.1 动作属性

| 属性 | 说明 |
|------|------|
| 动作 ID | 唯一标识 |
| 动作名称 | 中文名称 |
| 动作类型 | 表情/运动/声音/灯光 |
| 优先级 | 0-100 |
| 兼容硬件型号 | 适用设备列表 |
| 持续时间 | 秒 |
| 资源文件 | 动画/音频路径 |

### 5.2 API 接口

```
GET    /api/v1/actions              # 动作列表
GET    /api/v1/actions/:id         # 动作详情
POST   /api/v1/actions              # 创建动作
PUT    /api/v1/actions/:id          # 更新动作
DELETE /api/v1/actions/:id          # 删除动作
POST   /api/v1/actions/generate     # 生成动作序列
```

---

## 六、传感器事件

### 6.1 事件类型

| 类型 | 说明 |
|------|------|
| fall_detected | 跌落检测 |
| collision_detected | 碰撞检测 |
| touch_detected | 触摸检测 |
| low_battery | 低电量 |
| high_temperature | 高温告警 |
| offline | 设备离线 |

### 6.2 API 接口

```
GET  /api/v1/devices/:device_id/events      # 事件列表
GET  /api/v1/devices/:device_id/events/:id  # 事件详情
POST /api/v1/devices/:device_id/events     # 上报事件（设备端）
```

---

## 七、告警规则引擎

### 7.1 规则配置

| 字段 | 说明 |
|------|------|
| 规则名称 | 告警名称 |
| 条件 | 传感器阈值（如 battery < 20） |
| 级别 | info/warning/critical |
| 触发动作 | 通知/记录/自动化 |

### 7.2 API 接口

```
GET    /api/v1/alert/rules           # 规则列表
GET    /api/v1/alert/rules/:id      # 规则详情
POST   /api/v1/alert/rules           # 创建规则
PUT    /api/v1/alert/rules/:id       # 更新规则
DELETE /api/v1/alert/rules/:id       # 删除规则
GET    /api/v1/alert/rules/trigger   # 触发检查
```

---

## 八、设备模拟器

### 8.1 功能

| 功能 | 说明 |
|------|------|
| 虚拟设备 | 无需硬件即可测试 |
| 传感器模拟 | 模拟各种传感器数据 |
| 行为模拟 | 模拟宠物行为 |
| 故障模拟 | 模拟各种故障场景 |

### 8.2 用途

- 开发阶段无需真实设备
- CI/CD 自动化测试
- 演示环境
- 调试问题
