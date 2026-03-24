# OTA Worker 管理 PRD

## 1. 功能概述
OTA Worker 管理模块对固件升级的后台任务执行器进行配置和监控，包括并发数、速率限制、重试策略等参数的配置。

## 2. 页面布局与交互

### 页面路径
`/ota/worker` → `OtaWorkerConfigView.vue`（后端待实现）

### Worker 配置
| 字段 | 类型 | 说明 |
|------|------|------|
| 最大并发数 | Input | 同时升级设备数 |
| 升级速率 | Input | 每分钟升级设备数 |
| 重试次数 | Input | 失败重试次数 |
| 启用状态 | Switch | - |

### Worker 状态
- 运行中/已暂停
- 当前任务数
- CPU/内存使用率

## 3. API 契约

### 获取 Worker 配置
- 路径：`GET /api/v1/ota/worker/config`
- 响应：
```json
{
  "code": 0,
  "data": {
    "max_concurrency": 50,
    "rate_per_minute": 100,
    "retry_count": 3,
    "enabled": true
  }
}
```

### 更新 Worker 配置
- 路径：`PUT /api/v1/ota/worker/config`
- 请求体：同上

### Worker 状态
- 路径：`GET /api/v1/ota/worker/status`
- 响应：
```json
{
  "code": 0,
  "data": {
    "status": "running",
    "active_tasks": 25,
    "pending_tasks": 100,
    "cpu_usage": 15.5,
    "memory_usage": 30.2
  }
}
```

### 暂停/恢复 Worker
- 路径：`POST /api/v1/ota/worker/pause`
- 路径：`POST /api/v1/ota/worker/resume`

## 4. 验收标准
- [ ] Worker 配置保存成功
- [ ] Worker 状态实时显示
- [ ] 暂停/恢复操作正常
- [ ] 配置变更立即生效
