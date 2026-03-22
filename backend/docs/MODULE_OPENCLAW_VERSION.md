# OpenClaw AI 版本管理

**版本：** V1.0  
**模块负责人：** agentcp  
**编制日期：** 2026-03-22  

---

## 1. 概述

OpenClaw AI 版本管理模块负责管理 OpenClaw AI 模型的版本发布、兼容性矩阵、以及与 MiniClaw 固件版本的匹配关系。该模块确保设备能够获取到兼容的 AI 能力。

**业务目标：**
- 统一管理 OpenClaw AI 模型版本
- 维护 AI 版本与固件版本的兼容性矩阵
- 实现 AI 版本的灰度发布和回滚

---

## 2. 功能列表

| 功能 | 描述 | 优先级 | 触发方式 |
|------|------|--------|----------|
| AI 版本列表 | 查看所有 AI 模型版本及状态 | P0 | 人工 |
| AI 版本详情 | 查看单个版本的详细信息 | P0 | 人工 |
| 版本兼容性矩阵 | AI 版本与固件版本的兼容关系 | P0 | 自动 |
| 版本发布 | 发布新 AI 版本 | P0 | 人工 |
| 版本禁用 | 禁用问题版本 | P0 | 人工 |
| 版本回滚 | 回滚到上一稳定版本 | P1 | 人工 |
| 设备版本查询 | 设备查询可用的 AI 版本 | P0 | 自动 |
| 版本升级推送 | 向设备推送 AI 版本更新 | P1 | 自动 |

---

## 3. 数据模型

### 3.1 AI 版本表 (openclaw_versions)

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| id | uint | 主键 | PK, auto_increment |
| version_id | string | 版本唯一标识 | unique, not null, semver格式 |
| version_name | string | 版本名称 | not null |
| ai_model_type | string | AI 模型类型 | openai/claude/gemini/local |
| model_name | string | 模型名称 | not null |
| model_version | string | 模型版本 | not null |
| capabilities | json | 支持的能力列表 | not null |
| min_firmware_version | string | 最低固件版本要求 | not null |
| max_firmware_version | string | 最高固件版本要求 | nullable |
| release_type | string | 发布类型 | stable/beta/dev |
| status | int | 状态 | 1=开发中 2=测试中 3=已发布 4=已禁用 |
| release_notes | text | 发布说明 | nullable |
| changelog | text | 变更日志 | nullable |
| published_at | datetime | 发布时间 | nullable |
| deprecated_at | datetime | 废弃时间 | nullable |
| created_by | uint | 创建人 | FK to sys_users |
| created_at | datetime | 创建时间 | auto |
| updated_at | datetime | 更新时间 | auto |

### 3.2 兼容性矩阵表 (version_compatibility)

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| id | uint | 主键 | PK, auto_increment |
| openclaw_version_id | uint | OpenClaw 版本 | FK to openclaw_versions, not null |
| miniclaw_firmware_version_id | uint | MiniClaw 固件版本 | FK to firmwares, not null |
| compatibility_level | int | 兼容级别 | 1=完全兼容 2=部分兼容 3=不兼容 |
| notes | string | 备注 | nullable |
| tested_at | datetime | 测试时间 | nullable |
| tested_by | uint | 测试人 | FK to sys_users |
| created_at | datetime | 创建时间 | auto |
| updated_at | datetime | 更新时间 | auto |

### 3.3 设备 AI 版本记录表 (device_openclaw_version)

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| id | uint | 主键 | PK, auto_increment |
| device_id | string | 设备ID | FK to devices, not null |
| current_version_id | uint | 当前 AI 版本 | FK to openclaw_versions, not null |
| target_version_id | uint | 目标 AI 版本 | FK to openclaw_versions, nullable |
| upgrade_status | int | 升级状态 | 1=无更新 2=下载中 3=安装中 4=成功 5=失败 |
| error_code | string | 错误码 | nullable |
| error_message | string | 错误信息 | nullable |
| started_at | datetime | 开始时间 | nullable |
| completed_at | datetime | 完成时间 | nullable |
| created_at | datetime | 创建时间 | auto |
| updated_at | datetime | 更新时间 | auto |

---

## 4. 接口定义

### 4.1 获取 AI 版本列表

```
GET /api/v1/openclaw/versions
```

**参数：**

| 参数 | 类型 | 位置 | 必填 | 说明 |
|------|------|------|------|------|
| status | int | query | 否 | 状态筛选 |
| release_type | string | query | 否 | 发布类型筛选 |

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {
        "version_id": "v2.1.0",
        "version_name": "OpenClaw 2.1.0",
        "ai_model_type": "openai",
        "model_name": "gpt-4o",
        "capabilities": ["chat", "vision", "action"],
        "release_type": "stable",
        "status": 3,
        "published_at": "2026-03-01T00:00:00Z"
      }
    ],
    "total": 5
  }
}
```

### 4.2 获取版本详情

```
GET /api/v1/openclaw/versions/:version_id
```

### 4.3 获取兼容性矩阵

```
GET /api/v1/openclaw/compatibility
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "openclaw_version": "v2.1.0",
      "firmware_versions": [
        { "version": "v1.3.0", "compatibility_level": 1 },
        { "version": "v1.2.0", "compatibility_level": 2 }
      ]
    }
  ]
}
```

### 4.4 发布新版本

```
POST /api/v1/openclaw/versions
```

### 4.5 禁用版本

```
POST /api/v1/openclaw/versions/:version_id/disable
```

### 4.6 设备查询可用版本

```
GET /api/v1/openclaw/versions/available
```

**参数：**

| 参数 | 类型 | 位置 | 必填 | 说明 |
|------|------|------|------|------|
| device_id | string | query | 是 | 设备ID |

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "current_version": "v2.0.0",
    "latest_stable": "v2.1.0",
    "available_versions": [
      { "version_id": "v2.1.0", "status": 3 }
    ]
  }
}
```

---

## 5. 前端页面

| 页面 | 路由 | 功能 |
|------|------|------|
| AI 版本列表 | /openclaw/versions | 版本列表、筛选、发布入口 |
| 版本详情 | /openclaw/versions/:id | 版本详情、兼容性、发布说明 |
| 兼容性配置 | /openclaw/compatibility | 兼容性矩阵配置 |
| 设备版本管理 | /openclaw/device-versions | 设备当前版本、升级管理 |
