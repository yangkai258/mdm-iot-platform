# App 市场 PRD

## 1. 功能概述
App 市场为宠物设备提供第三方应用分发平台，开发者可上传应用，用户可浏览、安装和管理应用。

## 2. 页面布局与交互

### 页面路径
`/market/apps` → `AppStoreView.vue`

### 应用商店
- 分类导航（游戏/教育/工具/社交）
- 应用卡片列表（图标+名称+评分）
- 搜索框
- 「我的应用」Tab

### 应用详情
- 应用介绍
- 版本列表
- 用户评分
- 「安装」「更新」「卸载」按钮

## 3. API 契约

### 应用列表
- 路径：`GET /api/v1/market/apps`
- 实际路径（market_controller.go）：通过 `RegisterRoutes` 注册
- 参数：`category`, `keyword`, `page`, `page_size`, `sort_by`
- 响应：
```json
{
  "code": 0,
  "data": {
    "items": [
      {
        "id": 1,
        "name": "宠物训练师",
        "app_key": "app-xxx",
        "icon_url": "https://...",
        "category": "education",
        "version": "1.2.0",
        "rating": 4.5,
        "downloads": 10000,
        "size_mb": 25.5,
        "price": 0,
        "status": "published"
      }
    ],
    "total": 100,
    "page": 1,
    "page_size": 20
  }
}
```

### 应用详情
- 路径：`GET /api/v1/market/apps/:id`
- 响应：包含完整应用信息和版本列表

### 创建应用
- 路径：`POST /api/v1/market/apps`
- 请求体：`{ "name": "...", "category": "...", "description": "...", "icon": "..." }`

### 发布应用
- 路径：`POST /api/v1/market/apps/:id/publish`

### 删除应用
- 路径：`DELETE /api/v1/market/apps/:id`

## 4. 验收标准
- [ ] 应用列表分类加载正常
- [ ] 搜索功能正常
- [ ] 应用详情显示正确
- [ ] 发布/下架正常
- [ ] 应用版本管理正常
