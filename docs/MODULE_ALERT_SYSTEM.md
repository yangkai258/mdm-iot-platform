# 模块 PRD：告警系统 (Alert System)

**版本：** V1.4
**模块负责人：** agentcp
**编制日期：** 2026-03-20

---

## 1. 概述

告警系统负责监控 M5Stack 设备运行状态，在设备异常（低电量、离线、温度过高等）时自动触发告警，并通过多种渠道通知运营人员。

**业务目标：**
- 灵活的告警规则配置
- 实时检测设备异常
- 多渠道通知（邮件/Webhook/SMS）
- 告警确认/解决流程

---

## 2. 功能列表

| 功能 | 描述 | 优先级 | 触发方式 | 前端入口/按钮 |
|------|------|--------|----------|--------------|
| 告警规则CRUD | 创建/更新/启用/禁用告警规则 | P0 | 人工 | 「新建规则」/「编辑」/「禁用」/「删除」按钮 |
| 告警触发 | 设备数据上报时自动检查规则 | P0 | 自动 | 无按钮 |
| 告警列表 | 查看所有告警记录 | P0 | 自动 | 无按钮 |
| 告警确认 | 运营人员确认告警 | P1 | 人工 | 「确认」按钮 / 「批量确认」按钮 |
| 告警解决 | 标记告警已处理 | P1 | 人工 | 「解决」按钮 / 「批量解决」按钮 |
| 告警通知 | 邮件/Webhook通知 | P1 | 自动 | 无按钮 |
| 告警统计 | Dashboard统计看板 | P2 | 自动 | 无按钮 |

---

## 3. 数据模型

### 3.1 告警规则表 (device_alert_rules)

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| name | string | 规则名称 |
| device_id | string | 设备ID，空表示所有设备 |
| alert_type | string | 告警类型 |
| condition | string | 条件 <, >, =, <=, >= |
| threshold | float64 | 阈值 |
| severity | int | 严重程度 1=低 2=中 3=高 4=严重 |
| enabled | bool | 是否启用 |
| notify_ways | string | 通知方式 email,webhook,sms |
| remark | string | 备注 |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |

**alert_type 枚举：** battery_low=低电量, offline=设备离线, temperature_high=温度过高, signal_weak=信号弱, jailbreak=越狱/Root检测, geofence_violation=地理围栏违规, compliance_violation=合规策略违规, storage_low=存储空间不足

### 3.2 告警记录表 (device_alerts)

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| rule_id | uint | 关联规则ID |
| device_id | string | 设备ID, index |
| alert_type | string | 告警类型 |
| severity | int | 严重程度 |
| message | string | 告警消息 |
| trigger_val | float64 | 触发值 |
| threshold | float64 | 阈值 |
| status | int | 1=未处理 2=已确认 3=已解决 |
| created_at | datetime | 创建时间 |

---

## 4. 接口定义

### 4.1 告警规则

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/alerts/rules | 规则列表 |
| POST | /api/v1/alerts/rules | 创建规则（name, device_id, alert_type, condition, threshold, severity, enabled, notify_ways, remark） |
| PUT | /api/v1/alerts/rules/:id | 更新规则 |
| PUT | /api/v1/alerts/rules/:id/enable | 启用规则 |
| PUT | /api/v1/alerts/rules/:id/disable | 禁用规则（body: {enabled: true/false}） |
| DELETE | /api/v1/alerts/rules/:id | 删除规则 |

### 4.2 告警记录

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/alerts | 告警列表（Query: device_id, alert_type, status, severity, start_time, end_time, page, page_size） |
| PUT | /api/v1/alerts/:id/confirm | 确认告警（status: 1→2） |
| PUT | /api/v1/alerts/:id/resolve | 解决告警（body: remark，status: 2→3） |
| PUT | /api/v1/alerts/batch/confirm | 批量确认（alert_ids数组） |
| PUT | /api/v1/alerts/batch/resolve | 批量解决（alert_ids数组） |

### 4.3 告警统计

```
GET /api/v1/alerts/stats
```

**响应：** total_alerts, pending_alerts, confirmed_alerts, resolved_alerts, today_alerts, critical_alerts, by_type, by_severity

---

## 5. 流程图

### 5.1 CheckAlerts 触发时机

`CheckAlerts` 函数在以下时机被调用：

```
┌──────────────────────────────────────────────────────────────┐
│                    CheckAlerts 触发时机                        │
├──────────────────────────────────────────────────────────────┤
│                                                              │
│  触发点 1：设备 MQTT 心跳上报                                 │
│  ───────────────────────────────────────────                 │
│  Topic: /device/{device_id}/up/status                        │
│  Handler: StatusMessageHandler()                             │
│  调用链: StatusMessageHandler() → CheckAlerts(db, deviceID, payload)
│  触发频率: 设备心跳间隔（默认30秒）                           │
│  检查内容: battery_low, temperature_high, signal_weak 等     │
│                                                              │
│  触发点 2：设备离线检测（定时任务）                           │
│  ───────────────────────────────────────────                 │
│  Cron: 每 60 秒执行一次                                       │
│  检查逻辑: devices 表中 last_heartbeat < now-90s 的设备     │
│  针对规则: alert_type='offline' 的告警规则                   │
│                                                              │
│  触发点 3：定时合规检测（策略引擎）                           │
│  ───────────────────────────────────────────                 │
│  Cron: 每 5 分钟执行一次                                       │
│  检查逻辑: 合规策略违规 (compliance_violation, geofence_violation)
│  调用方: PolicyEngine 或 ComplianceChecker                   │
│                                                              │
│  触发点 4：OTA 升级失败回调                                  │
│  ───────────────────────────────────────────                 │
│  调用方: OTAProgressHandler 或 OTA Worker                    │
│  检查内容: alert_type='ota_failure'（自动创建）             │
│                                                              │
│  触发点 5：手动触发（API）                                    │
│  ───────────────────────────────────────────                 │
│  POST /api/v1/alerts/devices/:device_id/check                │
│  用途: 手动对指定设备执行一次告警规则检查                     │
│                                                              │
└──────────────────────────────────────────────────────────────┘
```

**防抖机制：**
- 同一 device_id + alert_type 组合，在 `cooldown_minutes`（默认 30 分钟）内不重复创建告警
- 已解决（status=3）的告警不计入防抖

### 5.2 通知渠道配置

通知渠道在系统级别配置，支持多渠道叠加。

#### 5.2.1 配置存储

通知配置存储于 `system_config` 表（或独立 `notification_channels` 表）：

| 字段 | 类型 | 说明 |
|------|------|------|
| channel_type | string | email / webhook / sms |
| config | jsonb | 渠道配置 |
| enabled | bool | 是否启用 |
| updated_at | datetime | 更新时间 |

**config 示例：**

```json
// email 渠道
{
  "smtp_host": "smtp.example.com",
  "smtp_port": 465,
  "smtp_user": "alerts@example.com",
  "smtp_password": "encrypted_xxx",
  "from_address": "alerts@example.com",
  "to_addresses": ["admin@example.com", "ops@example.com"],
  "use_tls": true
}

// webhook 渠道
{
  "url": "https://hooks.example.com/alert",
  "method": "POST",
  "headers": {
    "Authorization": "Bearer xxx",
    "Content-Type": "application/json"
  },
  "template": "json",
  "timeout_seconds": 10
}

// sms 渠道（以阿里云为例）
{
  "provider": "aliyun",
  "access_key_id": "xxx",
  "access_key_secret": "encrypted_xxx",
  "sign_name": "MDM告警",
  "template_code": "SMS_xxx",
  "phone_numbers": ["13800138000"]
}
```

#### 5.2.2 通知发送流程

```
告警触发
    │
    ▼
查询告警规则的 notify_ways（email,webhook,sms）
    │
    ▼
对每个渠道并行发送：
    │
    ├─► email ── 构建邮件内容（告警类型/设备/触发值/时间）
    │           ── SMTP发送 ── 运营人员邮箱
    │           ── 记录 notification_logs（成功/失败）
    │
    ├─► webhook ── 构造POST Body（JSON格式）
    │            ── HTTP POST ── Webhook URL
    │            ── 记录 notification_logs（含响应状态码）
    │
    └─► sms ── 按模板替换变量（设备ID/告警类型/触发值）
             ── 短信网关API发送
             ── 记录 notification_logs
```

#### 5.2.3 通知模板

**邮件主题格式：**
```
【MDM告警】{severity_label} - {alert_type_label} - {device_id}
```

**邮件正文模板：**
```
设备管理系统告警通知

告警级别：{severity_label}
告警类型：{alert_type_label}
设备ID：{device_id}
触发时间：{created_at}
触发值：{trigger_val}
阈值：{threshold}
告警消息：{message}

请及时处理。
```

**Webhook POST Body：**
```json
{
  "alert_id": 1,
  "device_id": "550e8400-e29b-41d4-a716-446655440000",
  "alert_type": "battery_low",
  "severity": 2,
  "message": "设备电量低于15%",
  "trigger_val": 12.5,
  "threshold": 15.0,
  "created_at": "2026-03-20T12:05:00Z",
  "severity_label": "中"
}
```

#### 5.2.4 通知失败重试

- 重试次数：3 次
- 重试间隔：30 秒、60 秒、180 秒（指数退避）
- 重试仍失败：记录 `notification_logs.status=failed`，不阻塞告警创建

### 5.3 告警处理流程

```
告警产生 → [1:未处理] → 确认告警 PUT /alerts/:id/confirm → [2:已确认]
         → 运营人员处理问题 → 解决告警 PUT /alerts/:id/resolve → [3:已解决] → 告警关闭
```

---

## 6. 模块联动

| 联动模块 | 联动方式 | 说明 |
|----------|----------|------|
| 设备影子 | CheckAlerts在心跳处理中调用 | 设备状态数据触发告警 |
| OTA升级 | 升级失败创建告警 | 通知运营人员 |
| 系统管理 | 通知渠道配置 | 邮件/Webhook配置 |
| 数据分析 | 告警统计数据 | Dashboard展示 |

---

## 7. 验收标准

### P0 验收标准

| 用例 | 验收条件 | 测试方法 |
|------|----------|----------|
| 规则CRUD | 完整增删改查 | 调用各接口验证 |
| 规则启用/禁用 | 禁用后不再触发 | 禁用后模拟数据检查 |
| 告警触发 | 电量<15%创建告警 | 模拟低电量心跳 |
| 告警触发 | 离线>90s创建告警 | 停止心跳等待超时 |
| 告警列表 | 支持多条件筛选 | 按type/status/severity筛选 |
| 全局规则 | device_id=''规则对所有设备生效 | 创建全局规则验证 |

### P1 验收标准

| 用例 | 验收条件 | 测试方法 |
|------|----------|----------|
| 告警确认 | 状态1→2 | 调用确认接口 |
| 告警解决 | 状态2→3 | 调用解决接口 |
| 批量确认 | 一次确认多条 | 调用批量接口 |
| 邮件通知 | 触发告警后收到邮件 | 配置邮件后触发 |
| Webhook通知 | 触发告警后收到HTTP请求 | 配置Webhook后触发 |

### P2 验收标准

| 用例 | 验收条件 | 测试方法 |
|------|----------|----------|
| 告警统计 | Dashboard数据正确 | 统计数据对比DB |
| 严重程度分类 | 按severity分组统计 | 验证统计结果 |

---

## 8. UI设计指引

### 8.1 页面结构
- **左侧菜单**：告警管理 → 告警规则 / 告警记录 / 告警统计
- **顶部区域**：统计卡片（今日告警/待处理/严重告警/已解决）
- **中间区域**：Tab页签：告警规则 / 告警记录 / 告警统计
- **底部区域**：分页组件

### 8.2 组件选用
| 组件 | 用途 |
|------|------|
| a-table | 告警规则列表、告警记录列表 |
| a-card | 顶部统计卡片，4列布局 |
| a-tabs | Tab切换：告警规则/告警记录/告警统计 |
| a-drawer | 创建/编辑告警规则、告警详情 |
| a-modal | 确认告警对话框、解决告警对话框 |
| a-form | 告警规则配置表单（类型/条件/阈值/通知方式）|
| a-select | 告警类型筛选、状态筛选、严重程度筛选 |
| a-tag | 严重程度标签（低=蓝，中=黄，高=橙，严重=红）|
| a-alert | 告警统计页面Banner告警提示 |
| a-timeline | 告警时间线（时间倒序排列）|

### 8.3 参考模板
```
┌──────────────────────────────────────────────────────────────┐
│  [统计卡片]  今日告警:25  待处理:8  严重:2  已解决:20        │
├──────────────────────────────────────────────────────────────┤
│  [Tab: 告警规则 | 告警记录 | 告警统计]                        │
├──────────────────────────────────────────────────────────────┤
│  【告警规则Tab】                                             │
│  ┌──────────────────────────────────────────────────────┐   │
│  │ [+新建规则]                                          │   │
│  │ 规则名称       │类型      │条件  │阈值│严重│状态│操作  │   │
│  │ 低电量告警     │battery_low│  <   │15% │中  │启用│编辑禁用│   │
│  │ 设备离线告警   │offline   │  =   │ -- │高  │启用│编辑禁用│   │
│  └──────────────────────────────────────────────────────┘   │
│                                                              │
│  【告警记录Tab】                                             │
│  ┌──────────────────────────────────────────────────────┐   │
│  │ [设备筛选▼] [类型▼] [状态▼] [严重▼]                   │   │
│  │ 告警类型  │设备ID      │触发值│阈值│严重│状态│时间    │   │
│  │ 🔴低电量  │550e8400-..│ 12%  │15% │中  │待处理│10:05 │   │
│  │ 🟠设备离线│660e8400-..│  --  │ -- │高  │已确认│09:30 │   │
│  └──────────────────────────────────────────────────────┘   │
│  [批量确认]  [批量解决]                                       │
└──────────────────────────────────────────────────────────────┘
```

### 8.4 交互流程
```
告警规则页
    ├── 点击「新建规则」──► a-drawer ──► 选择类型/条件/阈值/通知渠道 ──► 保存
    ├── 点击「编辑」──► a-drawer ──► 修改规则 ──► 保存
    └── 点击「禁用」──► 规则不触发但保留记录

告警记录页
    ├── 点击「确认」──► a-modal ──► 确认处理 ──► 状态更新
    ├── 点击「解决」──► a-modal ──► 填写备注 ──► 已解决
    ├── 勾选多条 ──► 点击「批量确认」──► 批量更新状态
    └── 点击「详情」──► a-drawer ──► 查看告警详情+设备信息

告警统计页
    └── 查看告警趋势图、分布饼图、设备TOP10
```

### 8.5 关键状态显示
- **严重程度**：a-tag，红=严重(4)，橙=高(3)，黄=中(2)，蓝=低(1)
- **告警状态**：a-tag，待处理=红，已确认=橙，已解决=绿
- **告警类型图标**：红色=电量低，橙色=离线，蓝色=温度高，紫色=信号弱

---

## 附录 B. 修订记录

| 版本 | 日期 | 修订人 | 修订内容 |
|------|------|--------|----------|
| V1.0 | 2026-03-20 | agentcp | 初稿，基于代码调研 |
| V1.2 | 2026-03-20 | agentcp | 修订功能列表，补充触发方式和前端入口按钮列 |
| V1.4 | 2026-03-20 | agentcp | 重建文档结构，统一使用8章节格式，合并重复的八、九章节 |
| V1.5 | 2026-03-21 | agentcp | 补充 CheckAlerts 触发时机（5.1）、通知渠道配置（5.2，含模板和重试） |