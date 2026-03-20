# 模块 PRD：通知与消息 (Notification)

**版本：** V1.3  
**模块负责人：** agentcp  
**编制日期：** 2026-03-20  

---

## 1. 概述

通知与消息模块为 MDM 中台提供多渠道推送能力，支持向设备发送文本通知、企业公告，并提供命令反馈查看能力。

**业务目标：**
- 推送通知（向设备发送文本消息）
- 企业公告（全员/定向通知）
- 命令反馈查看（设备指令执行结果）
- 通知模板管理
- 通知历史记录

---

## 2. 功能列表

| 功能 | 描述 | 优先级 | 触发方式 | 前端入口/按钮 |
|------|------|--------|----------|--------------|
| 推送通知 | 向设备发送文本通知 | P0 | 人工 | 「发送通知」按钮 |
| 通知列表 | 查看所有发送的推送通知 | P0 | 自动 | 无按钮 |
| 通知模板 | 通知模板 CRUD | P1 | 人工 | 「新建模板」/「编辑」/「删除」按钮 |
| 公告管理 | 企业公告 CRUD（全员/定向） | P1 | 人工 | 「新建公告」/「编辑」/「删除」按钮 |
| 反馈查看 | 查看设备对命令的响应结果 | P1 | 自动 | 无按钮 |
| 通知统计 | 发送成功/失败/阅读统计 | P2 | 自动 | 无按钮 |

---

## 3. 数据模型

### 3.1 通知表 (notifications)

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| id | uint | 主键 | PK, auto |
| notification_code | string | 通知编码, unique | not null |
| title | string | 通知标题 | not null |
| content | string | 通知内容 | not null |
| notification_type | string | 类型 push/announcement/command_response | not null |
| target_type | string | 目标类型 device/user/all | not null |
| target_ids | jsonb | 目标 ID 列表 | nullable |
| status | string | 状态 pending/sent/failed/read | default 'pending' |
| sent_at | datetime | 发送时间 | nullable |
| read_count | int | 已读数 | default 0 |
| failed_count | int | 失败数 | default 0 |
| template_id | uint | 关联模板 ID, FK | nullable |
| created_by | string | 创建人 | not null |
| created_at | datetime | 创建时间 | auto |
| updated_at | datetime | 更新时间 | auto |

### 3.2 通知模板表 (notification_templates)

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| id | uint | 主键 | PK, auto |
| template_code | string | 模板编码, unique | not null |
| template_name | string | 模板名称 | not null |
| notification_type | string | 通知类型 | not null |
| title_template | string | 标题模板 | not null |
| content_template | string | 内容模板，支持变量占位符 {name} | not null |
| variables | jsonb | 变量定义 | nullable |
| status | int | 状态 1=启用 0=禁用 | default 1 |
| created_at | datetime | 创建时间 | auto |
| updated_at | datetime | 更新时间 | auto |

**variables 示例：**
```json
{
  "device_name": "设备名称",
  "owner_name": "主人名称",
  "current_time": "当前时间"
}
```

### 3.3 公告表 (announcements)

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| id | uint | 主键 | PK, auto |
| title | string | 公告标题 | not null |
| content | string | 公告内容，支持富文本 | not null |
| priority | string | 优先级 normal/important/urgent | default 'normal' |
| target_type | string | 目标类型 all/device/user/org_unit | not null |
| target_ids | jsonb | 目标 ID 列表 | nullable |
| effective_start | datetime | 生效开始时间 | not null |
| effective_end | datetime | 生效结束时间 | nullable |
| status | string | 状态 draft/published/expired | default 'draft' |
| published_at | datetime | 发布时间 | nullable |
| created_by | string | 创建人 | not null |
| created_at | datetime | 创建时间 | auto |
| updated_at | datetime | 更新时间 | auto |

### 3.4 设备通知记录表 (device_notifications)

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| id | uint | 主键 | PK, auto |
| notification_id | uint | 通知 ID, FK | not null |
| device_id | string | 设备 ID | not null, index |
| status | string | 状态 pending/delivered/read/failed | not null |
| delivered_at | datetime | 送达时间 | nullable |
| read_at | datetime | 已读时间 | nullable |
| error_message | string | 失败原因 | nullable |
| created_at | datetime | 创建时间 | auto |

---

## 4. 接口定义

### 4.1 推送通知

#### 4.1.1 发送推送通知
```
POST /api/v1/notifications/push
```
**请求体：**
```json
{
  "title": "固件升级通知",
  "content": "有新版本固件可用，请及时更新",
  "target_type": "device",
  "target_ids": ["550e8400-e29b-41d4-a716-446655440000"],
  "created_by": "admin"
}
```

#### 4.1.2 批量发送推送
```
POST /api/v1/notifications/push/batch
```
**请求体：**
```json
{
  "title": "系统维护通知",
  "content": "系统将于今晚22:00-23:00进行维护",
  "target_type": "all",
  "created_by": "admin"
}
```

#### 4.1.3 通知列表
```
GET /api/v1/notifications
```
**Query:** notification_type, status, page, page_size

#### 4.1.4 通知详情
```
GET /api/v1/notifications/:id
```

---

### 4.2 通知模板

#### 4.2.1 模板列表
```
GET /api/v1/notifications/templates
```

#### 4.2.2 创建模板
```
POST /api/v1/notifications/templates
```
**请求体：**
```json
{
  "template_code": "TPL_FIRMWARE_UPDATE",
  "template_name": "固件升级通知",
  "notification_type": "push",
  "title_template": "【固件更新】{device_name}",
  "content_template": "您的设备 {device_name} 有新版本固件 {version} 可用",
  "variables": {
    "device_name": "设备名称",
    "version": "版本号"
  }
}
```

#### 4.2.3 使用模板发送
```
POST /api/v1/notifications/push/from-template
```
**请求体：**
```json
{
  "template_id": 1,
  "variables": {
    "device_name": "Mimi一号",
    "version": "v1.3.0"
  },
  "target_type": "device",
  "target_ids": ["550e8400-e29b-41d4-a716-446655440000"],
  "created_by": "admin"
}
```

---

### 4.3 公告管理

#### 4.3.1 公告列表
```
GET /api/v1/announcements
```
**Query:** status, priority, page, page_size

#### 4.3.2 创建公告
```
POST /api/v1/announcements
```
**请求体：**
```json
{
  "title": "2026年度公司年会通知",
  "content": "<p>公司将于2026年12月31日举办年度年会...</p>",
  "priority": "important",
  "target_type": "all",
  "effective_start": "2026-03-20T00:00:00Z",
  "effective_end": "2026-03-25T23:59:59Z"
}
```

#### 4.3.3 发布公告
```
POST /api/v1/announcements/:id/publish
```

#### 4.3.4 撤销公告
```
POST /api/v1/announcements/:id/withdraw
```

---

### 4.4 设备通知

#### 4.4.1 设备通知列表
```
GET /api/v1/devices/:device_id/notifications
```

#### 4.4.2 标记已读
```
POST /api/v1/notifications/:id/read
```

---

### 4.5 通知统计

#### 4.5.1 通知统计
```
GET /api/v1/notifications/:id/stats
```
**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "notification_id": 1,
    "total_targets": 100,
    "sent_count": 98,
    "delivered_count": 95,
    "read_count": 50,
    "failed_count": 2,
    "delivery_rate": 95.0,
    "read_rate": 50.0
  }
}
```

---

## 5. 流程图

### 5.1 推送通知发送流程

```
管理员发送推送通知
    │
    ▼
POST /api/v1/notifications/push
    │
    ├─→ 创建 notifications 记录 (status=pending)
    ├─→ 解析 target_ids 批量创建 device_notifications
    │
    ▼
通知发送 Worker (后台处理)
    │
    ▼
遍历 device_notifications (status=pending)
    │
    ├─→ 通过 MQTT 发送到设备 /device/{id}/down/notification
    │
    ├─→ 成功 ──► status=delivered, delivered_at=now()
    │
    └─→ 失败 ──► status=failed, error_message=...
    │
    ▼
更新 notifications.sent_count / failed_count
```

### 5.2 设备接收通知流程

```
设备收到 MQTT 通知消息
    │
    ▼
设备端展示通知（弹窗/通知栏）
    │
    ▼
用户点击查看
    │
    ▼
设备上报已读状态
    │
    ▼
POST /api/v1/notifications/:id/read ──► 更新 device_notifications.read_at
```

---

## 6. 模块联动

| 联动模块 | 联动方式 | 说明 |
|----------|----------|------|
| 设备管理 | 通知发送到设备 | MQTT /device/{id}/down/notification |
| 会员管理 | 通知可关联用户 | target_type=user |
| OTA升级 | 升级结果通知 | notification_type=command_response |
| 应用管理 | 应用安装结果通知 | notification_type=command_response |
| 内容管理 | 新内容分发通知 | notification_type=push |
| 告警系统 | 告警通知 | notification_type=push |

---

## 7. 验收标准

### P0 验收标准

| 用例 | 验收条件 | 测试方法 |
|------|----------|----------|
| 推送通知发送 | 通知正确创建并发送到设备 | POST /notifications/push 后检查 MQTT |
| 批量发送 | 一次性发送到多个设备 | POST /notifications/push/batch 验证 |
| 通知列表 | 正确展示所有通知 | GET /notifications 验证 |

### P1 验收标准

| 用例 | 验收条件 | 测试方法 |
|------|----------|----------|
| 模板管理 | 模板创建后可用 | POST 模板后使用发送 |
| 公告发布 | 公告正确展示 | 发布后查看公告列表 |
| 已读标记 | 已读后状态更新 | 设备点击后验证 |

### P2 验收标准

| 用例 | 验收条件 | 测试方法 |
|------|----------|----------|
| 通知统计 | 发送/送达/已读数据正确 | 对比 device_notifications 数据 |

---

## 8. UI设计指引

### 页面结构
- **左侧菜单**：通知管理 → 推送通知 / 公告管理 / 通知模板 / 通知统计
- **顶部区域**：统计卡片（今日发送 / 送达率 / 已读率 / 待发送）
- **中间区域**：Tab 页签：推送通知 / 公告管理 / 通知模板
- **底部区域**：分页组件

### 组件选用
| 组件 | 用途 |
|------|------|
| a-table | 通知列表、公告列表、模板列表 |
| a-card | 顶部统计卡片 |
| a-tabs | Tab 切换：推送通知 / 公告管理 / 通知模板 |
| a-drawer | 发送通知、创建公告、创建模板 |
| a-modal | 删除确认 |
| a-select | 通知类型筛选、目标类型筛选、状态筛选 |
| a-tag | 通知状态标签（待发送=黄，已发送=绿，失败=红）|
| a-picker | 日期时间选择（定时发送）|
| a-radio-group | 优先级选择（普通/重要/紧急）|
| a-input | 通知标题、内容输入 |
| a-textarea | 富文本内容输入 |
| a-progress | 送达率/已读率进度条 |
| a-empty | 空数据状态 |

### 参考模板
```
┌──────────────────────────────────────────────────────────────┐
│  [统计卡片]  今日发送:50  送达率:98%  已读率:75%  待发送:5   │
├──────────────────────────────────────────────────────────────┤
│  [Tab: 推送通知 | 公告管理 | 通知模板]                        │
├──────────────────────────────────────────────────────────────┤
│  【推送通知 Tab】                                            │
│  ┌──────────────────────────────────────────────────────┐   │
│  │ [+发送通知]  [筛选: 类型▼  状态▼]                      │   │
│  ├──────────────────────────────────────────────────────┤   │
│  │ 通知标题       │类型│目标    │状态  │发送时间│操作       │   │
│  │ 固件升级通知   │推送│3台设备 │已发送│ 10:30 │详情       │   │
│  │ 系统维护通知   │推送│全部    │已发送│ 09:00 │详情       │   │
│  └──────────────────────────────────────────────────────┘   │
│                                                              │
│  【公告管理 Tab】                                            │
│  ┌──────────────────────────────────────────────────────┐   │
│  │ [+新建公告]                                            │   │
│  │ 公告标题         │优先级│状态  │有效期      │操作       │   │
│  │ 年会通知         │重要  │已发布│03/20-03/25 │详情撤回   │   │
│  └──────────────────────────────────────────────────────┘   │
└──────────────────────────────────────────────────────────────┘
```

### 交互流程
```
推送通知页
    │
    ├── 点击「发送通知」──► a-drawer ──► 填写标题/内容/选择目标 ──► 发送
    │
    ├── 点击「详情」──► a-drawer ──► 查看发送统计/目标列表
    │
    └── 支持定时发送（设置 effective_start）

公告管理页
    │
    ├── 点击「新建公告」──► a-drawer ──► 填写内容/选择目标/有效期 ──► 保存为草稿
    │
    ├── 点击「发布」──► 公告生效
    │
    └── 点击「撤回」──► 公告失效

通知模板页
    │
    ├── 点击「新建模板」──► a-drawer ──► 填写模板信息 ──► 保存
    │
    └── 点击「使用模板发送」──► a-drawer ──► 填写变量值+选择目标 ──► 发送
```

### 关键状态显示
- **通知状态**：a-tag，待发送=黄，已发送=绿，部分失败=橙，失败=红
- **公告状态**：a-tag，草稿=灰，已发布=绿，已过期=蓝，已撤回=红
- **优先级**：a-tag，普通=灰，重要=黄，紧急=红
- **送达率/已读率**：a-progress，成功率>95%绿色，>80%黄色，<80%红色

---

## 修订记录

| 版本 | 日期 | 修订人 | 修订内容 |
|------|------|--------|----------|
| V1.3 | 2026-03-20 | agentcp | 全新模块，基于新增功能需求重建 |
