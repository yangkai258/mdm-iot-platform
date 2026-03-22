# MDM 控制中台 — OTA 模块产品需求文档

**版本：** V1.0  
**编制角色：** 产品经理 (agentcp)  
**编制日期：** 2026-03-20  
**文档状态：** 初稿  

---

## 一、模块概述

### 1.1 模块定位

OTA（Over-The-Air）模块是 MDM 控制中台的核心功能之一，负责管理设备固件包的生命周期，以及将固件安全、可控地推送到目标设备。模块支持灰度发布策略，提供完整的部署任务管理和设备升级进度追踪能力。

### 1.2 核心业务流程

```
固件包管理                          部署任务管理                      设备升级
    │                                    │                           │
    ├── 上传/创建固件包                    │                           │
    ├── 固件包列表/筛选                    │                           │
    ├── 编辑/删除固件包                    │                           │
    └── 固件包详情                        │                           │
                                        │                           │
                            创建部署任务 ────────────────┐            │
                            ├── 目标设备筛选              │            │
                            ├── 灰度策略配置              │            │
                            └── 部署状态流转              │            │
                                                    ▼                  │
                                        ┌─────────────────┐           │
                                        │  后台 OTA Worker │           │
                                        │  消费部署队列     │           │
                                        │  下发升级指令     │           │
                                        └────────┬────────┘           │
                                                 │ MQTT /device/{id}/down/cmd
                                                 ▼                  │
                                        ┌─────────────────┐           │
                                        │  M5Stack 设备端   │◄─────────┘
                                        │  下载固件 → 刷写  │
                                        │  上报升级进度    │
                                        └────────┬────────┘
                                                 │ MQTT /device/{id}/up/ota_progress
                                                 ▼
                                        设备升级进度入库 ← OTAProgress 表
```

### 1.3 关键设计约束

| 约束项 | 说明 |
|--------|------|
| 文件存储 | 固件包支持本地上传（`/uploads/firmware/`）或外部 URL 引用 |
| 灰度策略 | 全量 / 百分比 / 白名单三选一，互斥 |
| 部署状态机 | `pending → running → paused/completed/failed` |
| 设备筛选 | 按 `hardware_model` 精确匹配，不允许跨型号升级 |
| 版本校验 | 不允许降级，除非显式开启 `allow_downgrade` 标志 |
| 自动暂停 | 成功率 < 80% 时自动暂停部署，需人工确认后继续 |
| 并发限制 | 同一时刻同一设备只能有一个活跃部署任务 |

---

## 二、OTA 固件包管理

### 2.1 固件包列表查询

**接口描述：** 分页查询固件包列表，支持按硬件型号、版本号、状态等条件筛选。

**请求方法：** `GET`  
**请求路径：** `/api/v1/ota/packages`  
**是否需要认证：** 是

**Arco Design Vue 组件对应：** `a-table` + `a-select`（筛选硬件型号）+ `a-input-search`（版本号搜索）

**请求参数：**

| 参数名 | 类型 | 必填 | 说明 | Arco 组件映射 |
|--------|------|------|------|---------------|
| page | int | 否 | 页码，默认 1 | `a-pagination :current` |
| page_size | int | 否 | 每页条数，默认 20，最大 100 | `a-pagination :page-size` |
| hardware_model | string | 否 | 硬件型号（精确匹配） | `a-select` 筛选 |
| version | string | 否 | 固件版本号（模糊匹配） | `a-input-search` |
| is_active | bool | 否 | 是否启用：`true` / `false` | `a-select` |
| upload_source | string | 否 | 上传来源：`local` / `remote` | `a-select` |
| search | string | 否 | 综合搜索（匹配固件名称或版本号） | `a-input-search` |

**响应参数：**

| 字段 | 类型 | 说明 |
|------|------|------|
| code | int | 状态码，0=成功 |
| message | string | 提示信息 |
| data.list | array | 固件包列表 |
| data.pagination.page | int | 当前页码 |
| data.pagination.page_size | int | 每页条数 |
| data.pagination.total | int | 总记录数 |
| data.pagination.total_pages | int | 总页数 |

**响应 data.list 字段：**

| 字段 | 类型 | 说明 |
|------|------|------|
| id | int | 固件包 ID |
| name | string | 固件包名称 |
| version | string | 固件版本号，如 `v1.2.3` |
| hardware_model | string | 目标硬件型号，如 `M5Stack-Core2` |
| file_size | int | 文件大小（字节） |
| file_url | string | 固件包下载地址 |
| file_md5 | string | MD5 校验值（32位十六进制） |
| upload_source | string | 上传来源：`local` / `remote` |
| is_active | bool | 是否启用 |
| is_mandatory | bool | 是否为强制升级版本 |
| allow_downgrade | bool | 是否允许降级 |
| release_notes | string | 更新说明 |
| created_by | string | 创建人用户名 |
| created_at | string | 创建时间，ISO 8601 |
| updated_at | string | 更新时间，ISO 8601 |
| deployed_count | int | 已部署设备数 |
| success_count | int | 成功升级设备数 |
| success_rate | float | 升级成功率（百分比） |

**响应示例：**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "name": "M5Stack-Core2 固件 v1.2.3",
        "version": "v1.2.3",
        "hardware_model": "M5Stack-Core2",
        "file_size": 1048576,
        "file_url": "https://cdn.example.com/firmware/core2_v1.2.3.bin",
        "file_md5": "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4",
        "upload_source": "remote",
        "is_active": true,
        "is_mandatory": false,
        "allow_downgrade": false,
        "release_notes": "修复设备离线后无法重连的问题",
        "created_by": "admin",
        "created_at": "2026-03-20T08:00:00Z",
        "updated_at": "2026-03-20T08:00:00Z",
        "deployed_count": 50,
        "success_count": 48,
        "success_rate": 96.0
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 3,
      "total_pages": 1
    }
  }
}
```

---

### 2.2 固件包上传 / 创建

**接口描述：** 创建一个新的固件包记录，支持本地上传或填写外部 URL。

**请求方法：** `POST`  
**请求路径：** `/api/v1/ota/packages`  
**是否需要认证：** 是

**Arco Design Vue 组件对应：** `a-form` + `a-upload` + `a-input` + `a-select` + `a-textarea`

**请求参数（Content-Type: `multipart/form-data` 或 `application/json`）：**

| 参数名 | 类型 | 必填 | 说明 | Arco 组件映射 |
|--------|------|------|------|---------------|
| name | string | 是 | 固件包名称，长度 1-128 | `a-input` |
| version | string | 是 | 固件版本号，格式 `vX.Y.Z`，不能与同型号已有版本重复 | `a-input` |
| hardware_model | string | 是 | 目标硬件型号，精确匹配已有型号枚举 | `a-select` |
| file | file | 否 | 固件包文件（上传来回写 `file_url` 和 `file_md5`），与 `file_url` 二选一 | `a-upload` |
| file_url | string | 否 | 外部固件包 URL，与 `file` 二选一 | `a-input` |
| file_md5 | string | 否 | MD5 校验值（外部 URL 时必填；本地上传时自动计算） | `a-input` |
| is_active | bool | 否 | 是否启用，默认 `true` | `a-switch` |
| is_mandatory | bool | 否 | 是否强制升级，默认 `false` | `a-switch` |
| allow_downgrade | bool | 否 | 是否允许降级，默认 `false` | `a-switch` |
| release_notes | string | 否 | 更新说明，长度 0-1024 | `a-textarea` |

**请求示例（外部 URL 模式）：**

```json
{
  "name": "M5Stack-Core2 固件 v1.2.4",
  "version": "v1.2.4",
  "hardware_model": "M5Stack-Core2",
  "file_url": "https://cdn.example.com/firmware/core2_v1.2.4.bin",
  "file_md5": "b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5",
  "is_active": true,
  "is_mandatory": false,
  "allow_downgrade": false,
  "release_notes": "新增电量显示功能，优化内存占用"
}
```

**响应参数：**

| 字段 | 类型 | 说明 |
|------|------|------|
| code | int | 状态码，0=成功 |
| message | string | 提示信息 |
| data | object | 创建的固件包完整对象 |

**响应示例：**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 4,
    "name": "M5Stack-Core2 固件 v1.2.4",
    "version": "v1.2.4",
    "hardware_model": "M5Stack-Core2",
    "file_size": 1258291,
    "file_url": "https://cdn.example.com/firmware/core2_v1.2.4.bin",
    "file_md5": "b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5",
    "upload_source": "remote",
    "is_active": true,
    "is_mandatory": false,
    "allow_downgrade": false,
    "release_notes": "新增电量显示功能，优化内存占用",
    "created_by": "admin",
    "created_at": "2026-03-20T09:00:00Z",
    "updated_at": "2026-03-20T09:00:00Z",
    "deployed_count": 0,
    "success_count": 0,
    "success_rate": 0
  }
}
```

**错误响应（版本号重复）：**

```json
{
  "code": 4006,
  "message": "该硬件型号下已存在相同版本号",
  "error_code": "ERR_OTA_PACKAGE_DUPLICATE_VERSION"
}
```

---

### 2.3 固件包编辑

**接口描述：** 编辑指定固件包的元信息，不支持修改固件文件本身（需删除后重建）。

**请求方法：** `PUT`  
**请求路径：** `/api/v1/ota/packages/:id`  
**是否需要认证：** 是

**Arco Design Vue 组件对应：** `a-form` + `a-modal`（编辑弹窗）

**路径参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 固件包 ID |

**请求参数：**

| 参数名 | 类型 | 必填 | 说明 | Arco 组件映射 |
|--------|------|------|------|---------------|
| name | string | 否 | 固件包名称 | `a-input` |
| is_active | bool | 否 | 是否启用 | `a-switch` |
| is_mandatory | bool | 否 | 是否强制升级 | `a-switch` |
| allow_downgrade | bool | 否 | 是否允许降级 | `a-switch` |
| release_notes | string | 否 | 更新说明 | `a-textarea` |

> **注意：** `version`、`hardware_model`、`file_url`、`file_md5` 创建后不可修改，需删除重建。

**响应示例：**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 4,
    "name": "M5Stack-Core2 固件 v1.2.4（正式版）",
    "version": "v1.2.4",
    "hardware_model": "M5Stack-Core2",
    "is_active": true,
    "is_mandatory": true,
    "allow_downgrade": false,
    "release_notes": "新增电量显示功能，优化内存占用。标记为强制升级",
    "updated_at": "2026-03-20T10:00:00Z"
  }
}
```

---

### 2.4 固件包删除

**接口描述：** 删除指定固件包。已关联未完成部署任务的固件包不允许删除。

**请求方法：** `DELETE`  
**请求路径：** `/api/v1/ota/packages/:id`  
**是否需要认证：** 是

**Arco Design Vue 组件对应：** `a-button`（删除）+ `a-modal`（确认弹窗）+ `a-message`

**路径参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 固件包 ID |

**响应示例（成功）：**

```json
{
  "code": 0,
  "message": "success"
}
```

**错误响应（有未完成部署）：**

```json
{
  "code": 4007,
  "message": "该固件包存在未完成的部署任务，无法删除",
  "error_code": "ERR_OTA_PACKAGE_HAS_ACTIVE_DEPLOYMENT"
}
```

---

### 2.5 固件包详情

**接口描述：** 查询指定固件包的完整信息。

**请求方法：** `GET`  
**请求路径：** `/api/v1/ota/packages/:id`  
**是否需要认证：** 是

**Arco Design Vue 组件对应：** `a-descriptions`（详情展示）

**路径参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 固件包 ID |

**响应参数：**

| 字段 | 类型 | 说明 |
|------|------|------|
| id | int | 固件包 ID |
| name | string | 固件包名称 |
| version | string | 固件版本号 |
| hardware_model | string | 目标硬件型号 |
| file_size | int | 文件大小（字节） |
| file_url | string | 固件包下载地址 |
| file_md5 | string | MD5 校验值 |
| upload_source | string | 上传来源：`local` / `remote` |
| is_active | bool | 是否启用 |
| is_mandatory | bool | 是否强制升级 |
| allow_downgrade | bool | 是否允许降级 |
| release_notes | string | 更新说明 |
| created_by | string | 创建人用户名 |
| created_at | string | 创建时间 |
| updated_at | string | 更新时间 |
| deployed_count | int | 已部署设备数 |
| running_count | int | 当前进行中设备数 |
| success_count | int | 成功升级设备数 |
| failed_count | int | 失败设备数 |
| success_rate | float | 升级成功率 |
| recent_deployments | array | 最近 5 条部署记录摘要 |

**响应示例：**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "name": "M5Stack-Core2 固件 v1.2.3",
    "version": "v1.2.3",
    "hardware_model": "M5Stack-Core2",
    "file_size": 1048576,
    "file_url": "https://cdn.example.com/firmware/core2_v1.2.3.bin",
    "file_md5": "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4",
    "upload_source": "remote",
    "is_active": true,
    "is_mandatory": false,
    "allow_downgrade": false,
    "release_notes": "修复设备离线后无法重连的问题",
    "created_by": "admin",
    "created_at": "2026-03-20T08:00:00Z",
    "updated_at": "2026-03-20T08:00:00Z",
    "deployed_count": 50,
    "running_count": 2,
    "success_count": 47,
    "failed_count": 1,
    "success_rate": 97.9,
    "recent_deployments": [
      {
        "deployment_id": 10,
        "name": "生产环境全量发布 v1.2.3",
        "status": "running",
        "created_at": "2026-03-20T09:00:00Z"
      }
    ]
  }
}
```

---

## 三、OTA 部署任务管理

### 3.1 创建部署任务

**接口描述：** 创建一个 OTA 部署任务，支持全量、百分比灰度和白名单三种发布策略。

**请求方法：** `POST`  
**请求路径：** `/api/v1/ota/deployments`  
**是否需要认证：** 是

**Arco Design Vue 组件对应：** `a-form` + `a-steps`（策略配置向导）+ `a-select` + `a-input-number` + `a-table`（白名单选择）

**请求参数：**

| 参数名 | 类型 | 必填 | 说明 | Arco 组件映射 |
|--------|------|------|------|---------------|
| name | string | 是 | 部署任务名称，长度 1-64 | `a-input` |
| package_id | int | 是 | 目标固件包 ID | `a-select`（联动硬件型号） |
| hardware_model | string | 是 | 目标硬件型号（必须与固件包匹配，自动校验） | `a-select`（只读，由 package_id 联动） |
| strategy_type | string | 是 | 灰度策略类型：`full` / `percentage` / `whitelist` | `a-radio-group` |
| strategy_config | object | 是 | 灰度策略配置（见下表） | — |
| scheduled_at | string | 否 | 定时发布时间，ISO 8601。为空则立即发布 | `a-date-picker` |
| notify_on_start | bool | 否 | 启动时是否发送通知，默认 `false` | `a-switch` |
| pause_on_failure_threshold | float | 否 | 失败率阈值，超过则自动暂停，默认 20.0（%） | `a-input-number` |

**strategy_config 灰度策略详细配置：**

| strategy_type | strategy_config 字段 | 类型 | 必填 | 说明 |
|---------------|---------------------|------|------|------|
| `full` | 无需配置 | — | — | 全量发布，忽略 `target_device_count` 和 `whitelist` 字段 |
| `percentage` | target_percentage | float | 是 | 灰度百分比，范围 1-100 |
| `whitelist` | whitelist | array[string] | 是 | 白名单设备 ID 数组，长度 1-1000 |

**strategy_config 请求示例：**

```json
// 全量发布
{
  "strategy_type": "full",
  "strategy_config": {}
}

// 百分比灰度
{
  "strategy_type": "percentage",
  "strategy_config": {
    "target_percentage": 30.0
  }
}

// 白名单
{
  "strategy_type": "whitelist",
  "strategy_config": {
    "whitelist": ["DEV001", "DEV002", "DEV003"]
  }
}
```

**完整请求示例（百分比灰度）：**

```json
{
  "name": "生产环境 30% 灰度发布 v1.2.4",
  "package_id": 4,
  "hardware_model": "M5Stack-Core2",
  "strategy_type": "percentage",
  "strategy_config": {
    "target_percentage": 30.0
  },
  "scheduled_at": "2026-03-21T00:00:00Z",
  "notify_on_start": false,
  "pause_on_failure_threshold": 20.0
}
```

**响应参数：**

| 字段 | 类型 | 说明 |
|------|------|------|
| code | int | 状态码，0=成功 |
| message | string | 提示信息 |
| data | object | 部署任务对象 |

**data 字段详情：**

| 字段 | 类型 | 说明 |
|------|------|------|
| deployment_id | int | 部署任务 ID |
| name | string | 部署任务名称 |
| package_id | int | 固件包 ID |
| hardware_model | string | 目标硬件型号 |
| strategy_type | string | 灰度策略类型 |
| strategy_config | object | 灰度策略配置 |
| target_device_count | int | 本次部署目标设备数（实际命中数） |
| status | string | 初始状态：`pending` |
| scheduled_at | string | 定时发布时间（可为 null） |
| pause_on_failure_threshold | float | 自动暂停失败率阈值 |
| created_by | string | 创建人 |
| created_at | string | 创建时间 |
| updated_at | string | 更新时间 |

**响应示例：**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "deployment_id": 15,
    "name": "生产环境 30% 灰度发布 v1.2.4",
    "package_id": 4,
    "hardware_model": "M5Stack-Core2",
    "strategy_type": "percentage",
    "strategy_config": {
      "target_percentage": 30.0
    },
    "target_device_count": 30,
    "status": "pending",
    "scheduled_at": "2026-03-21T00:00:00Z",
    "pause_on_failure_threshold": 20.0,
    "created_by": "admin",
    "created_at": "2026-03-20T10:00:00Z",
    "updated_at": "2026-03-20T10:00:00Z"
  }
}
```

**错误响应（固件包与硬件型号不匹配）：**

```json
{
  "code": 4008,
  "message": "固件包与目标硬件型号不匹配",
  "error_code": "ERR_OTA_STRATEGY_MODEL_MISMATCH"
}
```

**错误响应（白名单设备型号不匹配）：**

```json
{
  "code": 4009,
  "message": "白名单中存在不属于目标硬件型号的设备",
  "error_code": "ERR_OTA_WHITELIST_MODEL_MISMATCH"
}
```

---

### 3.2 部署任务列表

**接口描述：** 分页查询部署任务列表，支持按状态、硬件型号、时间范围等条件筛选。

**请求方法：** `GET`  
**请求路径：** `/api/v1/ota/deployments`  
**是否需要认证：** 是

**Arco Design Vue 组件对应：** `a-table`（主列表）+ `a-select`（状态/型号筛选）+ `a-range-picker`（时间范围）+ `a-tabs`（快捷状态筛选）

**请求参数：**

| 参数名 | 类型 | 必填 | 说明 | Arco 组件映射 |
|--------|------|------|------|---------------|
| page | int | 否 | 页码，默认 1 | `a-pagination :current` |
| page_size | int | 否 | 每页条数，默认 20，最大 100 | `a-pagination :page-size` |
| status | string | 否 | 部署状态筛选：`pending` / `running` / `paused` / `completed` / `failed` | `a-select` / `a-tabs` 快捷筛选 |
| hardware_model | string | 否 | 硬件型号（精确匹配） | `a-select` |
| package_id | int | 否 | 固件包 ID | `a-select` |
| strategy_type | string | 否 | 灰度策略类型：`full` / `percentage` / `whitelist` | `a-select` |
| created_by | string | 否 | 创建人 | `a-input` |
| created_at_start | string | 否 | 创建时间起，ISO 8601 | `a-range-picker` |
| created_at_end | string | 否 | 创建时间止，ISO 8601 | `a-range-picker` |
| search | string | 否 | 搜索关键字（匹配任务名称） | `a-input-search` |

**status 状态说明：**

| 状态值 | 说明 | 触发条件 |
|--------|------|----------|
| `pending` | 待发布 | 任务创建后尚未开始，或定时任务未到时 |
| `running` | 进行中 | 后台 Worker 已开始下发升级指令 |
| `paused` | 已暂停 | 管理员手动暂停，或自动暂停（失败率超阈值） |
| `completed` | 已完成 | 所有目标设备升级完成（成功+失败均算完成） |
| `failed` | 全部失败 | 所有目标设备升级均失败（极少出现） |

**响应 data.list 字段：**

| 字段 | 类型 | 说明 |
|------|------|------|
| deployment_id | int | 部署任务 ID |
| name | string | 部署任务名称 |
| package_id | int | 固件包 ID |
| package_version | string | 固件版本号（冗余展示） |
| hardware_model | string | 目标硬件型号 |
| strategy_type | string | 灰度策略类型 |
| target_device_count | int | 目标设备总数 |
| pending_count | int | 待升级设备数 |
| running_count | int | 升级中设备数 |
| success_count | int | 成功设备数 |
| failed_count | int | 失败设备数 |
| success_rate | float | 当前成功率（百分比） |
| status | string | 当前状态 |
| scheduled_at | string | 定时发布时间（可为 null） |
| paused_at | string | 暂停时间（可为 null） |
| completed_at | string | 完成时间（可为 null） |
| created_by | string | 创建人 |
| created_at | string | 创建时间 |
| updated_at | string | 更新时间 |

**响应示例：**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "deployment_id": 15,
        "name": "生产环境 30% 灰度发布 v1.2.4",
        "package_id": 4,
        "package_version": "v1.2.4",
        "hardware_model": "M5Stack-Core2",
        "strategy_type": "percentage",
        "target_device_count": 30,
        "pending_count": 5,
        "running_count": 10,
        "success_count": 14,
        "failed_count": 1,
        "success_rate": 93.3,
        "status": "running",
        "scheduled_at": null,
        "paused_at": null,
        "completed_at": null,
        "created_by": "admin",
        "created_at": "2026-03-20T10:00:00Z",
        "updated_at": "2026-03-20T11:00:00Z"
      },
      {
        "deployment_id": 14,
        "name": "测试环境全量 v1.2.3",
        "package_id": 1,
        "package_version": "v1.2.3",
        "hardware_model": "M5Stack-Core2",
        "strategy_type": "full",
        "target_device_count": 50,
        "pending_count": 0,
        "running_count": 0,
        "success_count": 48,
        "failed_count": 2,
        "success_rate": 96.0,
        "status": "completed",
        "scheduled_at": null,
        "paused_at": null,
        "completed_at": "2026-03-20T09:30:00Z",
        "created_by": "admin",
        "created_at": "2026-03-20T08:30:00Z",
        "updated_at": "2026-03-20T09:30:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 5,
      "total_pages": 1
    }
  }
}
```

---

### 3.3 部署任务详情 / 进度

**接口描述：** 查询指定部署任务的完整信息，包括整体进度和设备维度升级明细。

**请求方法：** `GET`  
**请求路径：** `/api/v1/ota/deployments/:deployment_id`  
**是否需要认证：** 是

**Arco Design Vue 组件对应：** `a-descriptions`（任务概览）+ `a-progress`（整体进度）+ `a-table`（设备明细）+ `a-steps`

**路径参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| deployment_id | int | 是 | 部署任务 ID |

**响应参数：**

| 字段 | 类型 | 说明 |
|------|------|------|
| deployment_id | int | 部署任务 ID |
| name | string | 部署任务名称 |
| package_id | int | 固件包 ID |
| package_info | object | 固件包摘要（id/name/version/file_size） |
| hardware_model | string | 目标硬件型号 |
| strategy_type | string | 灰度策略类型 |
| strategy_config | object | 灰度策略配置 |
| target_device_count | int | 目标设备总数 |
| status | string | 当前状态 |
| scheduled_at | string | 定时发布时间 |
| pause_reason | string | 暂停原因（status=paused 时有值） |
| auto_paused | bool | 是否为自动暂停 |
| pause_on_failure_threshold | float | 自动暂停失败率阈值 |
| pending_count | int | 待升级设备数 |
| running_count | int | 升级中设备数 |
| success_count | int | 成功设备数 |
| failed_count | int | 失败设备数 |
| success_rate | float | 当前成功率 |
| created_by | string | 创建人 |
| created_at | string | 创建时间 |
| updated_at | string | 更新时间 |
| completed_at | string | 完成时间 |
| device_progress | array | 设备维度升级明细（分页） |

**device_progress 分页参数：**

| 参数名 | 类型 | 必填 | 说明 | Arco 组件映射 |
|--------|------|------|------|---------------|
| page | int | 否 | 页码，默认 1 | `a-pagination` |
| page_size | int | 否 | 每页条数，默认 50，最大 200 | `a-pagination` |
| status | string | 否 | 设备升级状态筛选 | `a-select` |
| search | string | 否 | 搜索设备 ID | `a-input-search` |

**device_progress 字段：**

| 字段 | 类型 | 说明 |
|------|------|------|
| device_id | string | 设备 ID |
| sn_code | string | 设备序列号 |
| ota_status | string | 升级状态：`pending` / `downloading` / `verifying` / `flashing` / `success` / `failed` |
| ota_message | string | 升级说明或错误信息 |
| progress_percent | int | 下载进度 0-100 |
| started_at | string | 开始升级时间 |
| completed_at | string | 升级完成时间 |
| retry_count | int | 重试次数 |

**ota_status 状态说明：**

| 状态值 | 说明 |
|--------|------|
| `pending` | 等待升级指令 |
| `downloading` | 正在下载固件 |
| `verifying` | 正在校验 MD5 |
| `flashing` | 正在刷写固件 |
| `success` | 升级成功 |
| `failed` | 升级失败 |

**响应示例：**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "deployment_id": 15,
    "name": "生产环境 30% 灰度发布 v1.2.4",
    "package_id": 4,
    "package_info": {
      "id": 4,
      "name": "M5Stack-Core2 固件 v1.2.4",
      "version": "v1.2.4",
      "file_size": 1258291
    },
    "hardware_model": "M5Stack-Core2",
    "strategy_type": "percentage",
    "strategy_config": {
      "target_percentage": 30.0
    },
    "target_device_count": 30,
    "status": "running",
    "scheduled_at": null,
    "pause_reason": null,
    "auto_paused": false,
    "pause_on_failure_threshold": 20.0,
    "pending_count": 5,
    "running_count": 10,
    "success_count": 14,
    "failed_count": 1,
    "success_rate": 93.3,
    "created_by": "admin",
    "created_at": "2026-03-20T10:00:00Z",
    "updated_at": "2026-03-20T11:00:00Z",
    "completed_at": null,
    "device_progress": {
      "list": [
        {
          "device_id": "DEV001",
          "sn_code": "SN2026032000001",
          "ota_status": "success",
          "ota_message": "升级成功",
          "progress_percent": 100,
          "started_at": "2026-03-20T10:05:00Z",
          "completed_at": "2026-03-20T10:07:30Z",
          "retry_count": 0
        },
        {
          "device_id": "DEV002",
          "sn_code": "SN2026032000002",
          "ota_status": "downloading",
          "ota_message": "正在下载固件",
          "progress_percent": 45,
          "started_at": "2026-03-20T10:08:00Z",
          "completed_at": null,
          "retry_count": 0
        }
      ],
      "pagination": {
        "page": 1,
        "page_size": 50,
        "total": 30,
        "total_pages": 1
      }
    }
  }
}
```

---

### 3.4 暂停部署

**接口描述：** 暂停一个进行中的部署任务。暂停后，后台 Worker 停止下发新的升级指令，但已在下载或刷写中的设备不受影响。

**请求方法：** `POST`  
**请求路径：** `/api/v1/ota/deployments/:deployment_id/pause`  
**是否需要认证：** 是

**Arco Design Vue 组件对应：** `a-button`（暂停按钮）+ `a-modal`（确认弹窗，填写暂停原因）+ `a-message`

**路径参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| deployment_id | int | 是 | 部署任务 ID |

**请求参数：**

| 参数名 | 类型 | 必填 | 说明 | Arco 组件映射 |
|--------|------|------|------|---------------|
| reason | string | 是 | 暂停原因，长度 1-256 | `a-input`（多行） |

**响应示例：**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "deployment_id": 15,
    "status": "paused",
    "pause_reason": "线上发现严重bug，需要回滚",
    "auto_paused": false,
    "paused_at": "2026-03-20T12:00:00Z"
  }
}
```

**错误响应（任务状态不允许暂停）：**

```json
{
  "code": 4010,
  "message": "只有 running 状态的部署任务可以暂停",
  "error_code": "ERR_OTA_DEPLOYMENT_NOT_PAUSEABLE"
}
```

---

### 3.5 恢复部署

**接口描述：** 恢复一个已暂停的部署任务，继续下发升级指令。

**请求方法：** `POST`  
**请求路径：** `/api/v1/ota/deployments/:deployment_id/resume`  
**是否需要认证：** 是

**Arco Design Vue 组件对应：** `a-button`（恢复按钮）+ `a-modal`（确认弹窗）

**路径参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| deployment_id | int | 是 | 部署任务 ID |

**请求参数：** 无

**响应示例：**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "deployment_id": 15,
    "status": "running",
    "paused_at": null,
    "pause_reason": null
  }
}
```

**错误响应（任务状态不允许恢复）：**

```json
{
  "code": 4011,
  "message": "只有 paused 状态的部署任务可以恢复",
  "error_code": "ERR_OTA_DEPLOYMENT_NOT_RESUMABLE"
}
```

---

### 3.6 取消部署

**接口描述：** 取消一个部署任务（pending / paused 状态可取消；running 状态需先暂停再取消）。取消后视为终止，不会再次调度。

**请求方法：** `POST`  
**请求路径：** `/api/v1/ota/deployments/:deployment_id/cancel`  
**是否需要认证：** 是

**Arco Design Vue 组件对应：** `a-button`（取消按钮）+ `a-modal`（确认弹窗，危险操作提示）

**路径参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| deployment_id | int | 是 | 部署任务 ID |

**请求参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| reason | string | 是 | 取消原因，长度 1-256 |

**响应示例：**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "deployment_id": 15,
    "status": "cancelled",
    "cancel_reason": "固件包存在严重问题，需重新制作后发布",
    "cancelled_by": "admin",
    "cancelled_at": "2026-03-20T12:30:00Z"
  }
}
```

**错误响应（任务已完成无法取消）：**

```json
{
  "code": 4012,
  "message": "已完成或已取消的部署任务无法操作",
  "error_code": "ERR_OTA_DEPLOYMENT_NOT_CANCELLABLE"
}
```

---

## 四、设备 OTA 状态

### 4.1 设备升级进度查询

**接口描述：** 查询单个设备在所有部署任务中的升级历史和当前状态。

**请求方法：** `GET`  
**请求路径：** `/api/v1/ota/devices/:device_id/progress`  
**是否需要认证：** 是

**Arco Design Vue 组件对应：** `a-timeline`（升级历史）+ `a-descriptions`（当前状态）+ `a-table`（历史记录分页）

**路径参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| device_id | string | 是 | 设备 ID |

**请求参数：**

| 参数名 | 类型 | 必填 | 说明 | Arco 组件映射 |
|--------|------|------|------|---------------|
| page | int | 否 | 页码，默认 1 | `a-pagination` |
| page_size | int | 否 | 每页条数，默认 20，最大 100 | `a-pagination` |

**响应参数：**

| 字段 | 类型 | 说明 |
|------|------|------|
| device_id | string | 设备 ID |
| sn_code | string | 设备序列号 |
| hardware_model | string | 硬件型号 |
| current_firmware_version | string | 当前固件版本 |
| latest_package_id | int | 最新可用固件包 ID |
| latest_package_version | string | 最新可用固件版本 |
| has_available_update | bool | 是否有可用更新 |
| is_mandatory_update | bool | 是否有强制更新 |
| current_ota | object/null | 当前进行中的升级任务（如有） |
| history | array | 历史升级记录（分页） |

**current_ota 字段：**

| 字段 | 类型 | 说明 |
|------|------|------|
| deployment_id | int | 部署任务 ID |
| package_version | string | 目标版本 |
| ota_status | string | 当前升级状态 |
| ota_message | string | 升级说明 |
| progress_percent | int | 下载进度 |
| started_at | string | 开始时间 |
| retry_count | int | 重试次数 |

**history 单条记录字段：**

| 字段 | 类型 | 说明 |
|------|------|------|
| deployment_id | int | 部署任务 ID |
| deployment_name | string | 部署任务名称 |
| from_version | string | 原版本 |
| to_version | string | 目标版本 |
| ota_status | string | 最终状态 |
| ota_message | string | 说明或错误信息 |
| started_at | string | 开始时间 |
| completed_at | string | 完成时间 |
| duration_seconds | int | 耗时（秒） |

**响应示例：**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "device_id": "DEV001",
    "sn_code": "SN2026032000001",
    "hardware_model": "M5Stack-Core2",
    "current_firmware_version": "v1.2.3",
    "latest_package_id": 4,
    "latest_package_version": "v1.2.4",
    "has_available_update": true,
    "is_mandatory_update": false,
    "current_ota": null,
    "history": {
      "list": [
        {
          "deployment_id": 14,
          "deployment_name": "测试环境全量 v1.2.3",
          "from_version": "v1.2.2",
          "to_version": "v1.2.3",
          "ota_status": "success",
          "ota_message": "升级成功",
          "started_at": "2026-03-20T08:35:00Z",
          "completed_at": "2026-03-20T08:37:30Z",
          "duration_seconds": 150
        }
      ],
      "pagination": {
        "page": 1,
        "page_size": 20,
        "total": 2,
        "total_pages": 1
      }
    }
  }
}
```

---

### 4.2 单设备强制升级

**接口描述：** 立即向单个设备下发强制升级指令，跳过部署任务排队。常用于紧急修复场景。

**请求方法：** `POST`  
**请求路径：** `/api/v1/ota/devices/:device_id/force-upgrade`  
**是否需要认证：** 是

**Arco Design Vue 组件对应：** `a-button`（强制升级按钮）+ `a-modal`（选择目标固件包）+ `a-select`

**路径参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| device_id | string | 是 | 设备 ID |

**请求参数：**

| 参数名 | 类型 | 必填 | 说明 | Arco 组件映射 |
|--------|------|------|------|---------------|
| package_id | int | 是 | 目标固件包 ID（系统校验型号匹配） | `a-select` |
| reason | string | 是 | 强制升级原因（记入操作日志） | `a-input` |

**响应参数：**

| 字段 | 类型 | 说明 |
|------|------|------|
| code | int | 状态码，0=成功 |
| message | string | 提示信息 |
| data | object | 下发结果 |

**data 字段：**

| 字段 | 类型 | 说明 |
|------|------|------|
| device_id | string | 设备 ID |
| package_id | int | 固件包 ID |
| package_version | string | 固件版本 |
| command_sent_at | string | 指令下发时间 |
| ota_status | string | 初始状态（pending） |

**响应示例：**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "device_id": "DEV001",
    "package_id": 4,
    "package_version": "v1.2.4",
    "command_sent_at": "2026-03-20T13:00:00Z",
    "ota_status": "pending"
  }
}
```

**错误响应（设备不在线）：**

```json
{
  "code": 4013,
  "message": "设备不在线，无法下发强制升级指令",
  "error_code": "ERR_OTA_DEVICE_OFFLINE"
}
```

**错误响应（设备已有进行中的升级任务）：**

```json
{
  "code": 4014,
  "message": "设备存在进行中的升级任务，请先等待或取消",
  "error_code": "ERR_OTA_DEVICE_BUSY"
}
```

**错误响应（固件包与设备型号不匹配）：**

```json
{
  "code": 4008,
  "message": "固件包与目标硬件型号不匹配",
  "error_code": "ERR_OTA_STRATEGY_MODEL_MISMATCH"
}
```

---

### 4.3 设备 OTA 检查（设备端直连）

**接口描述：** 设备主动查询是否有可用更新。此接口为设备端直连，无需 JWT 认证。

**请求方法：** `GET`  
**请求路径：** `/api/v1/ota/devices/:device_id/check`  
**是否需要认证：** 否（设备端直连）

**响应参数：**

| 字段 | 类型 | 说明 |
|------|------|------|
| has_update | bool | 是否有可用更新 |
| is_mandatory | bool | 是否强制升级 |
| package_id | int | 可用固件包 ID（无更新时为 null） |
| version | string | 可用固件版本号 |
| file_url | string | 固件下载 URL |
| file_md5 | string | MD5 校验值 |
| file_size | int | 文件大小（字节） |
| release_notes | string | 更新说明 |

**响应示例（有更新）：**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "has_update": true,
    "is_mandatory": false,
    "package_id": 4,
    "version": "v1.2.4",
    "file_url": "https://cdn.example.com/firmware/core2_v1.2.4.bin",
    "file_md5": "b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5",
    "file_size": 1258291,
    "release_notes": "新增电量显示功能，优化内存占用"
  }
}
```

**响应示例（无更新）：**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "has_update": false,
    "is_mandatory": false,
    "package_id": null
  }
}
```

---

### 4.4 设备升级进度上报（设备端直连）

**接口描述：** 设备端在升级过程中主动上报进度，上报 downloading / verifying / flashing / success / failed 状态。

**请求方法：** `POST`  
**请求路径：** `/api/v1/ota/devices/:device_id/report`  
**是否需要认证：** 否（设备端直连）

**请求参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| deployment_id | int | 否 | 关联的部署任务 ID（强制升级时可为空） |
| package_id | int | 是 | 固件包 ID |
| ota_status | string | 是 | 当前状态：`downloading` / `verifying` / `flashing` / `success` / `failed` |
| ota_message | string | 否 | 说明或错误信息（失败时必填） |
| progress_percent | int | 否 | 下载进度 0-100（downloading 时传） |

**请求示例：**

```json
{
  "deployment_id": 15,
  "package_id": 4,
  "ota_status": "downloading",
  "progress_percent": 45
}
```

**响应示例：**

```json
{
  "code": 0,
  "message": "success"
}
```

---

## 五、数据库模型

### 5.1 固件包表（ota_package）

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | serial | PK | 固件包 ID |
| name | varchar(128) | NOT NULL | 固件包名称 |
| version | varchar(32) | NOT NULL | 固件版本号 |
| hardware_model | varchar(64) | NOT NULL | 目标硬件型号 |
| file_size | bigint | DEFAULT 0 | 文件大小（字节） |
| file_url | varchar(512) | NOT NULL | 固件下载地址 |
| file_md5 | varchar(32) | | MD5 校验值 |
| upload_source | varchar(16) | NOT NULL, DEFAULT 'local' | 上传来源：local/remote |
| is_active | boolean | DEFAULT true | 是否启用 |
| is_mandatory | boolean | DEFAULT false | 是否强制升级 |
| allow_downgrade | boolean | DEFAULT false | 是否允许降级 |
| release_notes | text | | 更新说明 |
| created_by | varchar(64) | NOT NULL | 创建人 |
| created_at | timestamptz | DEFAULT NOW() | 创建时间 |
| updated_at | timestamptz | DEFAULT NOW() | 更新时间 |

**唯一索引：** `(hardware_model, version)`

### 5.2 部署任务表（ota_deployment）

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | serial | PK | 部署任务 ID |
| name | varchar(64) | NOT NULL | 部署任务名称 |
| package_id | int | FK → ota_package | 固件包 ID |
| hardware_model | varchar(64) | NOT NULL | 目标硬件型号 |
| strategy_type | varchar(16) | NOT NULL | 灰度策略：full/percentage/whitelist |
| strategy_config | jsonb | NOT NULL DEFAULT '{}' | 灰度策略配置 |
| target_device_count | int | DEFAULT 0 | 目标设备数 |
| pending_count | int | DEFAULT 0 | 待升级设备数 |
| running_count | int | DEFAULT 0 | 升级中设备数 |
| success_count | int | DEFAULT 0 | 成功设备数 |
| failed_count | int | DEFAULT 0 | 失败设备数 |
| status | varchar(16) | NOT NULL, DEFAULT 'pending' | 任务状态 |
| pause_reason | varchar(256) | | 暂停/取消原因 |
| auto_paused | boolean | DEFAULT false | 是否自动暂停 |
| pause_on_failure_threshold | float | DEFAULT 20.0 | 自动暂停失败率阈值（%） |
| scheduled_at | timestamptz | | 定时发布时间 |
| created_by | varchar(64) | NOT NULL | 创建人 |
| created_at | timestamptz | DEFAULT NOW() | 创建时间 |
| updated_at | timestamptz | DEFAULT NOW() | 更新时间 |
| completed_at | timestamptz | | 完成时间 |

### 5.3 设备升级进度表（ota_progress）

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | serial | PK | 进度记录 ID |
| deployment_id | int | FK → ota_deployment, INDEX | 部署任务 ID |
| device_id | varchar(64) | NOT NULL, INDEX | 设备 ID |
| package_id | int | FK → ota_package | 固件包 ID |
| from_version | varchar(32) | | 原固件版本 |
| to_version | varchar(32) | NOT NULL | 目标固件版本 |
| ota_status | varchar(16) | NOT NULL, DEFAULT 'pending' | 升级状态 |
| ota_message | varchar(256) | | 说明或错误信息 |
| progress_percent | int | DEFAULT 0 | 下载进度 0-100 |
| retry_count | int | DEFAULT 0 | 重试次数 |
| started_at | timestamptz | | 开始时间 |
| completed_at | timestamptz | | 完成时间 |
| created_at | timestamptz | DEFAULT NOW() | 创建时间 |
| updated_at | timestamptz | DEFAULT NOW() | 更新时间 |

**唯一索引：** `(deployment_id, device_id)`

---

## 六、错误码定义

| 错误码 | error_code | 说明 |
|--------|------------|------|
| 0 | — | 成功 |
| 400 | — | 参数错误（通用） |
| 401 | — | 未授权 / Token 过期 |
| 403 | — | 账号禁用 |
| 404 | — | 资源不存在 |
| 4006 | ERR_OTA_PACKAGE_DUPLICATE_VERSION | 同型号下版本号重复 |
| 4007 | ERR_OTA_PACKAGE_HAS_ACTIVE_DEPLOYMENT | 固件包有关联未完成部署任务 |
| 4008 | ERR_OTA_STRATEGY_MODEL_MISMATCH | 固件包与目标硬件型号不匹配 |
| 4009 | ERR_OTA_WHITELIST_MODEL_MISMATCH | 白名单中存在不属于目标型号的设备 |
| 4010 | ERR_OTA_DEPLOYMENT_NOT_PAUSEABLE | 只有 running 状态的任务可以暂停 |
| 4011 | ERR_OTA_DEPLOYMENT_NOT_RESUMABLE | 只有 paused 状态的任务可以恢复 |
| 4012 | ERR_OTA_DEPLOYMENT_NOT_CANCELLABLE | 已完成/已取消的任务无法操作 |
| 4013 | ERR_OTA_DEVICE_OFFLINE | 设备不在线 |
| 4014 | ERR_OTA_DEVICE_BUSY | 设备存在进行中的升级任务 |
| 4015 | ERR_OTA_NO_UPDATE_AVAILABLE | 设备没有可用更新 |
| 4016 | ERR_OTA_DOWNGRADE_NOT_ALLOWED | 不允许降级 |
| 5001 | ERR_INTERNAL | 服务器内部错误 |

---

## 七、修订记录

| 日期 | 版本 | 修改内容 |
|------|------|----------|
| 2026-03-20 | V1.0 | 初稿完成，定义 OTA 固件包管理、部署任务管理、设备 OTA 状态三大模块完整接口 |