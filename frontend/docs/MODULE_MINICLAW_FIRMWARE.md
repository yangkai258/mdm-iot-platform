# MiniClaw固件管理

**版本：** V1.1  
**模块负责人：** agentcp  
**编制日期：** 2026-03-22  

---

## 1. 概述

MiniClaw固件管理模块负责管理M5Stack设备的固件版本仓库、固件与设备的关联关系、版本兼容性检查以及OTA（Over-The-Air）升级流程。该模块是MDM中台设备管理层的重要组成部分，确保设备固件的安全、可控更新。

**业务目标：**
- 建立统一的固件版本仓库
- 维护设备与固件的关联关系
- 实现远程OTA升级
- 确保固件兼容性和安全性

---

## 2. 功能列表

| 功能 | 描述 | 优先级 | 触发方式 | 前端入口/按钮 |
|------|------|--------|----------|--------------|
| 固件上传 | 上传新固件文件到仓库 | P0 | 人工 | 管理后台「固件上传」 |
| 固件版本列表 | 查看所有固件版本及详情 | P0 | 人工 | 管理后台「固件管理」 |
| 固件详情 | 查看单个固件的详细信息 | P0 | 人工 | 管理后台点击版本号 |
| 固件仓库API | 提供固件查询、版本管理等RESTful API | P0 | 自动 | 无按钮 |
| 固件绑定 | 将固件版本绑定到特定设备 | P1 | 自动 | 无按钮（系统自动） |
| 兼容性检查 | 检查设备型号与固件兼容性 | P0 | 自动 | 无按钮（OTA时） |
| OTA升级任务 | 创建并执行设备OTA升级 | P0 | 人工 | 管理后台「设备升级」 |
| 升级进度查看 | 查看设备升级进度和结果 | P0 | 人工 | 设备详情页「OTA历史」 |
| 固件删除 | 删除废弃固件版本 | P2 | 人工 | 管理后台「删除」按钮 |
| 固件统计 | 统计固件下载次数、升级成功率等 | P1 | 人工 | 管理后台「统计」 |

---

## 3. 数据模型

### 3.1 固件表 (firmwares)

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| id | uint | 主键 | PK, auto_increment |
| firmware_id | string | 固件唯一标识 | unique, not null |
| version | string | 固件版本号 | not null, semver格式 |
| device_model | string | 适用设备型号 | not null |
| file_path | string | 固件文件存储路径 | not null |
| file_size | int | 文件大小(字节) | not null |
| checksum | string | 文件MD5校验 | not null |
| description | string | 版本说明 | nullable |
| min_hardware_version | string | 最低硬件版本 | nullable |
| release_type | string | 发布类型 | stable/beta/dev |
| status | int | 状态 | 1=草稿 2=测试中 3=已发布 4=已停用 |
| download_count | int | 下载次数 | default 0 |
| failure_count | int | 升级失败次数 | default 0 |
| created_by | uint | 上传人 | FK to sys_users |
| published_at | datetime | 发布时间 | nullable |
| created_at | datetime | 创建时间 | auto |
| updated_at | datetime | 更新时间 | auto |

### 3.2 固件与设备关联表 (firmware_device_bindings)

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| id | uint | 主键 | PK, auto_increment |
| binding_id | string | 关联唯一标识 | unique, not null |
| firmware_id | string | 固件ID | FK to firmwares, not null |
| device_id | string | 设备ID | FK to devices, not null |
| target_version | string | 目标版本 | not null |
| current_version | string | 当前版本 | not null |
| upgrade_status | int | 升级状态 | 1=待升级 2=升级中 3=成功 4=失败 5=已回滚 |
| upgrade_task_id | string | 升级任务ID | nullable |
| started_at | datetime | 开始时间 | nullable |
| completed_at | datetime | 完成时间 | nullable |
| failure_reason | string | 失败原因 | nullable |
| created_at | datetime | 创建时间 | auto |
| updated_at | datetime | 更新时间 | auto |

### 3.3 OTA升级任务表 (ota_upgrade_tasks)

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| id | uint | 主键 | PK, auto_increment |
| task_id | string | 任务唯一标识 | unique, not null |
| task_name | string | 任务名称 | not null |
| firmware_id | string | 目标固件 | FK to firmwares, not null |
| target_type | string | 目标类型 | device/group/all |
| target_devices | json | 目标设备列表 | nullable |
| device_group_id | string | 设备组ID | nullable |
| strategy | string | 升级策略 | immediate/scheduled/rolling |
| scheduled_at | datetime | 计划升级时间 | nullable |
| rolling_interval | int | 滚动升级间隔(秒) | default 0 |
| status | int | 任务状态 | 1=待执行 2=执行中 3=已完成 4=已取消 |
| total_devices | int | 总设备数 | auto |
| success_count | int | 成功数 | default 0 |
| failure_count | int | 失败数 | default 0 |
| in_progress_count | int | 进行中数 | default 0 |
| created_by | uint | 创建人 | FK to sys_users |
| started_at | datetime | 开始时间 | nullable |
| completed_at | datetime | 完成时间 | nullable |
| created_at | datetime | 创建时间 | auto |
| updated_at | datetime | 更新时间 | auto |

### 3.4 固件版本历史表 (firmware_versions)

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| id | uint | 主键 | PK, auto_increment |
| version_id | string | 版本历史ID | unique, not null |
| firmware_id | string | 固件ID | FK to firmwares, not null |
| version | string | 版本号 | not null |
| change_type | string | 变更类型 | feature/bugfix/optimization/security |
| change_description | string | 变更描述 | not null |
| created_at | datetime | 创建时间 | auto |

---

## 4. 接口定义

### 4.1 上传固件

```
POST /api/v1/firmware/upload
```

**参数：** multipart/form-data

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| file | file | 是 | 固件文件(.bin) |
| version | string | 是 | 版本号 |
| device_model | string | 是 | 设备型号 |
| description | string | 否 | 版本说明 |
| release_type | string | 否 | stable/beta/dev |

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "firmware_id": "fw-001",
    "version": "1.3.0",
    "device_model": "M5Stack-Basic",
    "file_size": 1048576,
    "checksum": "d41d8cd98f00b204e9800998ecf8427e",
    "status": 1,
    "created_at": "2026-03-20T10:45:00Z"
  }
}
```

### 4.2 获取固件列表

```
GET /api/v1/firmware
```

**参数：**

| 参数 | 类型 | 位置 | 必填 | 说明 |
|------|------|------|------|------|
| device_model | string | query | 否 | 设备型号筛选 |
| release_type | string | query | 否 | 发布类型筛选 |
| status | int | query | 否 | 状态筛选 |
| page | int | query | 否 | 页码 |
| page_size | int | query | 否 | 每页条数 |

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {
        "firmware_id": "fw-001",
        "version": "1.3.0",
        "device_model": "M5Stack-Basic",
        "release_type": "stable",
        "status": 3,
        "file_size": 1048576,
        "download_count": 523,
        "failure_count": 2,
        "published_at": "2026-03-15T00:00:00Z"
      }
    ],
    "pagination": { "page": 1, "page_size": 20, "total": 8, "total_pages": 1 }
  }
}
```

### 4.3 获取固件详情

```
GET /api/v1/firmware/{firmware_id}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "firmware_id": "fw-001",
    "version": "1.3.0",
    "device_model": "M5Stack-Basic",
    "min_hardware_version": "v1.2",
    "file_size": 1048576,
    "checksum": "d41d8cd98f00b204e9800998ecf8427e",
    "description": "修复了传感器掉线问题，优化了电池续航",
    "release_type": "stable",
    "status": 3,
    "download_count": 523,
    "failure_count": 2,
    "failure_rate": 0.38,
    "device_stats": {
      "total": 523,
      "up_to_date": 480,
      "outdated": 43
    },
    "created_by": { "user_id": 10001, "username": "admin" },
    "published_at": "2026-03-15T00:00:00Z",
    "created_at": "2026-03-10T00:00:00Z"
  }
}
```

### 4.4 更新固件

```
PUT /api/v1/firmware/{firmware_id}
```

**请求参数：**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| description | string | 否 | 版本说明 |
| status | int | 否 | 状态 |

### 4.5 删除固件

```
DELETE /api/v1/firmware/{firmware_id}
```

### 4.6 发布固件

```
POST /api/v1/firmware/{firmware_id}/publish
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "firmware_id": "fw-001",
    "version": "1.3.0",
    "status": 3,
    "published_at": "2026-03-22T10:00:00Z"
  }
}
```

### 4.7 获取固件下载链接

```
GET /api/v1/firmware/{firmware_id}/download
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "download_url": "https://cdn.example.com/firmware/fw-001.bin",
    "expires_at": "2026-03-22T11:00:00Z",
    "file_size": 1048576,
    "checksum": "d41d8cd98f00b204e9800998ecf8427e"
  }
}
```

### 4.8 固件仓库查询API

```
GET /api/v1/firmware/repository/query
```

**参数：**

| 参数 | 类型 | 位置 | 必填 | 说明 |
|------|------|------|------|------|
| device_model | string | query | 是 | 设备型号 |
| current_version | string | query | 否 | 当前版本号 |
| release_type | string | query | 否 | 发布类型，默认stable |

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "device_model": "M5Stack-Basic",
    "current_version": "1.2.5",
    "latest_version": "1.3.0",
    "has_update": true,
    "firmware_info": {
      "firmware_id": "fw-001",
      "version": "1.3.0",
      "release_type": "stable",
      "description": "修复了传感器掉线问题",
      "file_size": 1048576,
      "published_at": "2026-03-15T00:00:00Z"
    },
    "breaking_changes": false
  }
}
```

### 4.9 创建OTA升级任务

```
POST /api/v1/ota/tasks
```

**请求参数：**

| 参数 | 类型 | 位置 | 必填 | 说明 |
|------|------|------|------|------|
| task_name | string | body | 是 | 任务名称 |
| firmware_id | string | body | 是 | 目标固件ID |
| target_type | string | body | 是 | 目标类型 |
| target_devices | array | body | 否 | 目标设备列表 |
| device_group_id | string | body | 否 | 设备组ID |
| strategy | string | body | 是 | 升级策略 |
| scheduled_at | datetime | body | 否 | 计划升级时间 |
| rolling_interval | int | body | 否 | 滚动升级间隔 |

**请求示例：**
```json
{
  "task_name": "V1.3.0统一升级",
  "firmware_id": "fw-001",
  "target_type": "device",
  "target_devices": ["pet-001", "pet-002", "pet-003"],
  "strategy": "rolling",
  "rolling_interval": 60
}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "task_id": "ota-001",
    "task_name": "V1.3.0统一升级",
    "status": 1,
    "total_devices": 3,
    "created_at": "2026-03-20T10:50:00Z"
  }
}
```

### 4.10 查询升级任务状态

```
GET /api/v1/ota/tasks/{task_id}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "task_id": "ota-001",
    "task_name": "V1.3.0统一升级",
    "firmware_id": "fw-001",
    "status": 2,
    "total_devices": 3,
    "success_count": 1,
    "failure_count": 0,
    "in_progress_count": 2,
    "devices": [
      { "device_id": "pet-001", "status": 3, "current_version": "1.3.0" },
      { "device_id": "pet-002", "status": 2, "progress": 45 },
      { "device_id": "pet-003", "status": 2, "progress": 20 }
    ]
  }
}
```

### 4.11 暂停/恢复/取消升级任务

```
POST /api/v1/ota/tasks/{task_id}/pause
POST /api/v1/ota/tasks/{task_id}/resume
POST /api/v1/ota/tasks/{task_id}/cancel
```

### 4.12 获取设备固件状态

```
GET /api/v1/devices/{device_id}/firmware
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "device_id": "pet-001",
    "current_version": "1.2.5",
    "hardware_version": "v1.3",
    "latest_stable_version": "1.3.0",
    "upgrade_available": true,
    "binding": {
      "firmware_id": "fw-001",
      "target_version": "1.3.0",
      "upgrade_status": 1
    }
  }
}
```

### 4.13 设备检查固件更新

```
GET /api/v1/devices/{device_id}/firmware/check
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "device_id": "pet-001",
    "has_update": true,
    "current_version": "1.2.5",
    "available_version": "1.3.0",
    "download_url": "/api/v1/firmware/fw-001/download",
    "file_size": 1048576,
    "checksum": "d41d8cd98f00b204e9800998ecf8427e"
  }
}
```

### 4.14 固件兼容性检查

```
POST /api/v1/firmware/{firmware_id}/compatibility-check
```

**请求示例：**
```json
{
  "device_model": "M5Stack-Basic",
  "hardware_version": "v1.3",
  "current_firmware_version": "1.2.0"
}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "compatible": true,
    "warnings": ["建议在电量>50%时升级"],
    "required": true,
    "min_version_met": true
  }
}
```

### 4.15 错误码定义

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 40001 | 固件文件无效 |
| 40002 | 固件版本格式错误 |
| 40003 | 设备型号不匹配 |
| 40004 | 固件不存在 |
| 40005 | OTA任务不存在 |
| 40006 | 设备不在线无法升级 |
| 40007 | 升级任务已达最大重试次数 |

---

## 5. OTA完整流程

### 5.1 OTA升级完整流程

```
管理员创建升级任务
       │
       ▼
固件管理模块
       │
       ▼
1. 兼容性检查
   - 设备型号匹配
   - 硬件版本满足
   - 当前版本可升级
       │
       ▼
2. 创建升级任务
   - 生成task_id
   - 初始化设备状态
   - 记录目标版本
       │
       ▼
┌──────┼──────┬──────┐
▼      ▼      ▼      ▼
立即升级  定时升级  滚动升级
       │
       ▼
3. MQTT通知设备
   /down/config
   {ota: true, url: ..., version: ...}
       │
       ▼
4. MiniClaw设备
   下载固件
   执行升级
   重启
       │
       ▼
5. 状态上报
   /up/status
   {ota_result: success}
       │
       ▼
6. 更新任务状态
   success_count++
   记录升级历史
```

### 5.2 设备固件检查更新流程

```
MiniClaw定期检查
       │
       ▼
MDM中台查询最新固件
       │
       ▼
比对当前版本
       │
       ▼
返回更新信息
{has_update: true, version: "1.3.0"}
       │
       ▼
MiniClaw确认
保存版本信息
等待OTA指令
```

---

## 6. 模块联动

### 6.1 与设备管理(DEVICE_MANAGEMENT)联动

- **触发时机：** 设备注册和OTA升级时
- **联动内容：** 设备注册时检查固件版本，记录到设备表；OTA升级时查询设备列表
- **数据流向：** 设备管理 -> 固件管理

### 6.2 与MiniClaw通信协议(MINICLAW_PROTOCOL)联动

- **触发时机：** OTA升级执行
- **联动内容：** 通过MQTT的 `/down/config` 下发固件信息；设备通过 `/up/status` 上报升级结果
- **数据流向：** 固件管理 -> MQTT Broker -> MiniClaw

### 6.3 与告警系统(ALERT_SYSTEM)联动

- **触发时机：** OTA升级失败时
- **联动内容：** 升级失败自动触发告警，通知管理员处理
- **数据流向：** 固件管理 -> 告警系统

---

## 7. 验收标准

### 7.1 功能验收

| 功能 | 验收条件 |
|------|----------|
| 固件上传 | 支持.bin格式，最大100MB，MD5校验 |
| 固件列表 | 支持按型号/状态/类型筛选，分页展示 |
| 固件仓库API | 提供完整RESTful接口供设备调用 |
| OTA升级 | 支持单设备/批量/定时/滚动升级 |
| 兼容性检查 | 设备型号和硬件版本双重校验 |
| 升级回滚 | 失败后可回滚到上一版本 |

### 7.2 性能验收

- 单固件上传最大100MB，超时30分钟
- 支持同时执行10个OTA任务
- 单任务最大设备数：10000台
- 升级状态实时更新，延迟<=5秒

### 7.3 安全性验收

- 固件文件MD5校验
- 固件签名验证（预留）
- OTA过程加密传输（HTTPS）

---

## 8. 前端页面设计

### 8.1 固件管理列表

```
┌────────────────────────────────────────────────────────────────────┐
│  固件管理                                           [上传固件]    │
├────────────────────────────────────────────────────────────────────┤
│  设备型号: [M5Stack-Basic v]  发布类型: [全部 v]  状态: [全部 v]│
├────────────────────────────────────────────────────────────────────┤
│                                                                    │
│  版本        │ 设备型号    │ 类型    │ 状态   │ 下载 │ 失败率 │ 操作│
│ ──────────────────────────────────────────────────────────────────│
│  v1.3.0     │ Basic       │ Stable  │ 已发布 │ 523  │ 0.4%   │ ...│
│  v1.4.0-beta│ Basic       │ Beta    │ 测试中 │  45  │ 0.0%   │ ...│
│  v1.2.5     │ Basic       │ Stable  │ 已发布 │ 1234 │ 1.2%   │ ...│
│  v1.3.0     │ Gray        │ Stable  │ 已发布 │  89  │ 0.0%   │ ...│
│                                                                    │
├────────────────────────────────────────────────────────────────────┤
│  说明: ● = 当前推荐版本                                            │
└────────────────────────────────────────────────────────────────────┘
```

### 8.2 OTA升级任务创建

```
┌────────────────────────────────────────────────────────────────────┐
│  创建OTA升级任务                                            [×]  │
├────────────────────────────────────────────────────────────────────┤
│  任务名称: [V1.3.0统一升级________________________________]       │
│                                                                    │
│  目标固件: [v1.3.0 (M5Stack-Basic) v]  ✓ 兼容性通过                │
│                                                                    │
│  升级范围: ○ 单设备  ● 批量设备  ○ 全部设备  ○ 设备组            │
│                                                                    │
│  选择设备: [pet-001] [pet-002] [pet-003] [+ 添加更多]              │
│                                                                    │
│  升级策略: ● 立即升级  ○ 定时升级  ○ 滚动升级                    │
│                                                                    │
│  滚动间隔: [60] 秒 (滚动升级时有效)                                │
│                                                                    │
│  ─────────────────────────────────────────────────────────────────│
│  预计升级: 3台设备                                                  │
│                                                                    │
│                                        [取消]  [创建任务]          │
└────────────────────────────────────────────────────────────────────┘
```

### 8.3 配色方案

| 用途 | 颜色 | 色值 |
|------|------|------|
| 稳定版 | 稳定绿 | #6BCB77 |
| 测试版 | 警告黄 | #FFD93D |
| 开发版 | 灰阶 | #B2BEC3 |
| 升级中 | 进行蓝 | #74B9FF |
| 升级失败 | 错误红 | #FF6B6B |
| 升级成功 | 成功绿 | #00B894 |


---

## 9. 页面布局规范

### 9.1 固件管理列表页面

**布局结构：**
1. 面包屑 → 页面标题
2. 筛选区（浅灰卡片 #F2F3F5）：设备型号 / 发布类型 / 状态
3. 操作栏（上传固件靠左）
4. 固件列表表格

**按钮规范：**
- [上传固件] — 左对齐
- [详情] [编辑] [发布] [删除] — 行内右对齐

**表格列：**
| 列名 | 宽度 | 说明 |
|------|------|------|
| 序号 | 60px | - |
| 版本 | 120px | - |
| 设备型号 | 120px | - |
| 文件大小 | 100px | - |
| 发布类型 | 100px | stable/beta/dev |
| 状态 | 100px | 草稿/测试中/已发布/已停用 |
| 下载次数 | 100px | - |
| 失败率 | 80px | - |
| 操作 | 120px | 详情/编辑/发布/删除 |

**分页：** 右下角，10/20/50/100 条

### 9.2 OTA升级任务页面

**布局结构：**
1. 面包屑 → 页面标题
2. 筛选区（浅灰卡片）：任务状态 / 设备型号
3. 操作栏（创建任务靠左）
4. OTA任务列表表格

**按钮规范：**
- [创建OTA任务] — 左对齐
- [详情] [暂停/恢复] [取消] — 行内右对齐

**表格列：**
| 列名 | 宽度 | 说明 |
|------|------|------|
| 序号 | 60px | - |
| 任务名称 | 200px | - |
| 目标固件 | 120px | - |
| 升级策略 | 100px | 立即/定时/滚动 |
| 状态 | 100px | 待执行/执行中/已完成/已取消 |
| 设备数 | 80px | - |
| 成功数 | 80px | - |
| 失败数 | 80px | - |
| 操作 | 120px | 详情/暂停/取消 |

**分页：** 右下角，10/20/50/100 条

### 9.3 弹窗规范

| 类型 | 使用场景 |
|------|----------|
| Drawer 抽屉 | 新建/编辑固件详情、创建OTA升级任务、任务详情 |
| Dialog 对话框 | 确认发布固件、确认删除固件、确认取消任务 |
| 全屏模态 | 暂无复杂表单场景 |
