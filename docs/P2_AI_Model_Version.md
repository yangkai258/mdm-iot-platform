# 模型版本管理 PRD

## 1. 功能概述
模型版本管理模块提供 AI 模型版本的上传、发布、回滚和生命周期管理，支持灰度发布和 A/B 测试。

## 2. 页面布局与交互

### 页面路径
`/ai/versions` → `ModelVersionView.vue`

### 模型列表
| 列 | 说明 |
|----|------|
| 模型名称 | model_name |
| 当前版本 | current_version |
| 状态 | status |
| 更新时间 | updated_at |
| 操作 | 版本管理 |

### 版本管理弹窗
- 版本列表（历史版本）
- 「发布新版本」按钮
- 「设为默认」按钮
- 「回滚」按钮

## 3. API 契约

### 模型列表
- 路径：`GET /api/v1/ai/models`
- 响应：包含 model_name, current_version, status

### 上传模型版本
- 路径：`POST /api/v1/ai/models/:id/versions`
- 请求体：
```json
{
  "version": "1.2.4",
  "file_url": "https://...",
  "checksum": "sha256:...",
  "release_notes": "修复xxx问题",
  "config": { "threshold": 0.8 }
}
```

### 切换默认版本
- 路径：`POST /api/v1/ai/models/:id/switch`
- 请求体：`{ "version": "1.2.4" }`

### 回滚模型
- 路径：`POST /api/v1/ai/models/:id/rollback`
- 请求体：`{ "target_version": "1.2.3" }`

## 4. 验收标准
- [ ] 版本列表加载正常
- [ ] 上传新版本成功
- [ ] 切换默认版本成功
- [ ] 回滚成功
- [ ] 历史版本可查看
