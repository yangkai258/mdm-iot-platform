# Agent QD - 前端开发任务

## 你是高级前端开发工程师 (agentqd)
你的任务是使用 Vue 3 (Composition API) + TypeScript + Arco Design 构建 MDM 管控台页面。

## UI 框架
**采用 Arco Design Vue (@arco-design/web-vue) 作为 UI 框架。**

---

## ✅ Sprint 2.1 - 应用管理页面 (已完成)

### 已完成文件：
- `src/views/apps/AppList.vue` - 应用列表，含统计卡片、筛选、创建应用抽屉、详情抽屉、审核操作
- `src/views/apps/AppVersions.vue` - 应用版本管理，含版本列表、上传新版本抽屉、版本详情抽屉
- `src/views/apps/AppDistributions.vue` - 应用分发任务，含统计卡片、分发任务列表、创建分发抽屉、任务详情进度
- `src/router.js` - 添加 `/apps/*` 路由
- `src/App.vue` - 添加「应用管理」侧边栏菜单

### API 对接：
- GET/POST `/api/v1/apps`
- GET/PUT `/api/v1/apps/:id`
- POST `/api/v1/apps/:id/approve`
- POST `/api/v1/apps/:id/reject`
- GET/POST `/api/v1/apps/:id/versions`
- DELETE `/api/v1/apps/:id/versions/:version_id`
- GET/POST `/api/v1/apps/distributions`
- GET `/api/v1/apps/distributions/:id`
- POST `/api/v1/apps/distributions/:id/cancel`

### Git Commit: `feat(apps): implement app management pages`

---

## ✅ Sprint 2.2 - 通知管理页面 (已完成)

### 已完成文件：
- `src/views/notifications/NotificationList.vue` - 推送通知列表，含统计卡片、发送通知抽屉、通知详情抽屉
- `src/views/notifications/NotificationTemplates.vue` - 通知模板管理，含模板CRUD、使用模板发送、变量说明
- `src/views/notifications/Announcements.vue` - 公告管理，含公告CRUD、发布/撤回、有效期设置
- `src/router/index.js` - 添加 /notifications/* 路由
- `src/App.vue` - 添加「通知管理」侧边栏菜单

### API 对接：
- GET/POST `/api/v1/notifications` / `/api/v1/notifications/push`
- GET/DELETE `/api/v1/notifications/:id`
- GET/POST `/api/v1/notification-templates`
- GET/PUT/DELETE `/api/v1/notification-templates/:id`
- GET/POST `/api/v1/announcements`
- POST `/api/v1/announcements/:id/publish`
- POST `/api/v1/announcements/:id/withdraw`
- POST `/api/v1/notifications/push/from-template`

### Git Commit: `feat(notifications): implement notification management pages`
