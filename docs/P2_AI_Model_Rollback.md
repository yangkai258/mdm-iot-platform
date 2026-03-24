# 模型热回滚 PRD

## 1. 功能概述
模型热回滚功能允许在发现模型异常时，快速将正在服务的模型回滚到之前稳定版本，最小化故障影响时间。

## 2. 页面布局与交互

### 页面路径
`/ai/rollback` → `ModelRollbackView.vue`

### 回滚操作面板
- 设备/模型选择
- 当前版本显示
- 可回滚版本列表
- 回滚原因（必填）
- 「确认回滚」按钮（红色警告样式）

### 回滚历史
| 列 | 说明 |
|----|------|
| 时间 | created_at |
| 模型 | model_name |
| 从版本 | from_version |
| 到版本 | to_version |
| 操作人 | operator |
| 原因 | reason |

## 3. API 契约

### 回滚操作
- 路径：`POST /api/v1/ai/models/:model_id/rollback`
- 请求体：
```json
{
  "target_version": "1.2.3",
  "reason": "线上发现异常预测"
}
```
- 响应：
```json
{
  "code": 0,
  "data": {
    "task_id": "rollback-task-uuid",
    "status": "started"
  }
}
```

### 回滚状态查询
- 路径：`GET /api/v1/ai/rollback/:task_id`
- 响应：包含 task_id, status, progress, started_at, completed_at

### 回滚历史
- 路径：`GET /api/v1/ai/rollback/history`
- 参数：`model_id`, `page`, `page_size`

## 4. 验收标准
- [ ] 回滚版本列表正确
- [ ] 回滚原因必填验证
- [ ] 回滚任务状态可追踪
- [ ] 回滚历史记录完整
- [ ] 回滚进度实时更新
