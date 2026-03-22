# 模块 PRD：内容与文档管理 (Content Management)

**版本：** V1.3  
**模块负责人：** agentcp  
**编制日期：** 2026-03-20  

---

## 1. 概述

内容与文档管理模块为 MDM 中台提供企业内部文件分发能力，支持上传和管理文档（PDF/PPT/Word等），并通过分发任务推送到设备或用户。

**业务目标：**
- 文件库统一管理（上传/分类/标签）
- 内容分发任务（推送到设备/用户/组）
- 安全容器保护（加密存储、禁止分享）
- 内容访问控制（按用户/设备/组织单元权限）

---

## 2. 功能列表

| 功能 | 描述 | 优先级 | 触发方式 | 前端入口/按钮 |
|------|------|--------|----------|--------------|
| 文件库 | 文件上传（PDF/PPT/Word/Excel/图片）、分类管理 | P0 | 人工 | 「上传文件」按钮 |
| 文件列表 | 分页查询文件，支持类型/分类/标签筛选 | P0 | 自动 | 无按钮 |
| 内容分发 | 创建分发任务，推送到设备/用户/组 | P0 | 人工 | 「新建分发」按钮 |
| 分类管理 | 文件分类 CRUD（制度/培训/公告等） | P1 | 人工 | 「新建分类」/「编辑」/「删除」按钮 |
| 标签管理 | 文件标签管理 | P1 | 人工 | 「新建标签」/「编辑」/「删除」按钮 |
| 访问控制 | 文件访问权限配置 | P1 | 人工 | 「权限配置」按钮 |
| 内容统计 | 下载次数/阅读次数统计 | P2 | 自动 | 无按钮 |

---

## 3. 数据模型

### 3.1 内容主表 (contents)

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| id | uint | 主键 | PK, auto |
| content_code | string | 内容编码, unique | not null |
| title | string | 文件标题 | not null |
| file_name | string | 原始文件名 | not null |
| file_type | string | 文件类型 pdf/ppt/pptx/doc/docx/xls/xlsx/zip/jpg/png | not null |
| file_size | int64 | 文件大小(字节) | default 0 |
| file_url | string | CDN 存储 URL | not null |
| file_md5 | string | MD5 校验码 | nullable |
| category_id | uint | 分类 ID, FK | nullable |
| description | string | 文件描述 | nullable |
| is_encrypted | bool | 是否加密存储 | default false |
| allow_share | bool | 是否允许分享 | default false |
| access_level | string | 访问级别 public/restricted/confidential | default 'public' |
| status | string | 状态 draft/published/archived | default 'draft' |
| download_count | int | 下载次数 | default 0 |
| view_count | int | 阅读次数 | default 0 |
| version | int | 版本号 | default 1 |
| created_by | string | 创建人 | not null |
| created_at | datetime | 创建时间 | auto |
| updated_at | datetime | 更新时间 | auto |

### 3.2 内容分类表 (content_categories)

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| id | uint | 主键 | PK, auto |
| name | string | 分类名称 | not null |
| parent_id | uint | 上级分类 | nullable |
| sort_order | int | 排序 | default 0 |
| created_at | datetime | 创建时间 | auto |

### 3.3 内容标签关联表 (content_tags)

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| id | uint | 主键 | PK, auto |
| tag_name | string | 标签名称 | not null |
| color | string | 标签颜色 HEX | nullable |
| created_at | datetime | 创建时间 | auto |

### 3.4 内容标签关系表 (content_tag_relations)

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| id | uint | 主键 | PK, auto |
| content_id | uint | 内容 ID, FK | not null |
| tag_id | uint | 标签 ID, FK | not null |

### 3.5 内容分发任务表 (content_distributions)

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| id | uint | 主键 | PK, auto |
| content_id | uint | 内容 ID, FK | not null |
| target_type | string | 目标类型 device/user/group/org_unit | not null |
| target_ids | jsonb | 目标 ID 列表 | not null |
| distribution_type | string | 分发类型 push（推送）/pull（用户可pull）| default 'push' |
| status | string | 状态 pending/sent/read/expired/cancelled | not null |
| sent_at | datetime | 发送时间 | nullable |
| expired_at | datetime | 过期时间 | nullable |
| created_by | string | 创建人 | not null |
| created_at | datetime | 创建时间 | auto |

### 3.6 内容访问权限表 (content_permissions)

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| id | uint | 主键 | PK, auto |
| content_id | uint | 内容 ID, FK | not null |
| permission_type | string | 权限类型 user/group/org_unit | not null |
| permission_id | string | 权限目标 ID | not null |
| permission_level | string | 权限级别 view/download/print | default 'view' |
| expires_at | datetime | 权限过期时间 | nullable |
| created_at | datetime | 创建时间 | auto |

---

## 4. 接口定义

### 4.1 文件库管理

#### 4.1.1 文件列表
```
GET /api/v1/contents
```
**Query:** category_id, file_type, status, keyword, page, page_size

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "content_code": "CONTENT001",
        "title": "员工手册2026",
        "file_type": "pdf",
        "file_size": 5242880,
        "category_name": "制度文件",
        "status": "published",
        "download_count": 150,
        "view_count": 500,
        "created_at": "2026-03-20T10:00:00Z"
      }
    ],
    "pagination": { "page": 1, "page_size": 20, "total": 100 }
  }
}
```

#### 4.1.2 文件详情
```
GET /api/v1/contents/:id
```

#### 4.1.3 上传文件
```
POST /api/v1/contents
```
**请求体：**
```json
{
  "content_code": "CONTENT001",
  "title": "员工手册2026",
  "file_name": "员工手册2026.pdf",
  "file_type": "pdf",
  "file_size": 5242880,
  "file_url": "https://cdn.example.com/contents/employee-handbook-2026.pdf",
  "file_md5": "d41d8cd98f00b204e9800998ecf8427e",
  "category_id": 1,
  "description": "2026年新版员工手册",
  "is_encrypted": false,
  "allow_share": false,
  "access_level": "public",
  "created_by": "admin"
}
```

#### 4.1.4 更新文件
```
PUT /api/v1/contents/:id
```

#### 4.1.5 发布/归档文件
```
POST /api/v1/contents/:id/publish
POST /api/v1/contents/:id/archive
```

#### 4.1.6 删除文件
```
DELETE /api/v1/contents/:id
```

---

### 4.2 分类管理

#### 4.2.1 分类列表
```
GET /api/v1/content/categories
```

#### 4.2.2 创建分类
```
POST /api/v1/content/categories
```
**请求体：**
```json
{
  "name": "制度文件",
  "parent_id": null,
  "sort_order": 1
}
```

---

### 4.3 内容分发

#### 4.3.1 创建分发任务
```
POST /api/v1/contents/distributions
```
**请求体：**
```json
{
  "content_id": 1,
  "target_type": "org_unit",
  "target_ids": ["dept-hr", "dept-tech"],
  "distribution_type": "push",
  "expired_at": "2026-12-31T23:59:59Z",
  "created_by": "admin"
}
```

#### 4.3.2 分发任务列表
```
GET /api/v1/contents/distributions
```
**Query:** content_id, status, page, page_size

#### 4.3.3 查看设备/用户内容列表
```
GET /api/v1/contents/my
```
返回当前用户/设备有权限访问的内容列表

---

### 4.4 访问权限

#### 4.4.1 设置内容权限
```
POST /api/v1/contents/:id/permissions
```
**请求体：**
```json
{
  "permissions": [
    {
      "permission_type": "org_unit",
      "permission_id": "dept-hr",
      "permission_level": "download",
      "expires_at": "2026-12-31T23:59:59Z"
    },
    {
      "permission_type": "user",
      "permission_id": "user-001",
      "permission_level": "print",
      "expires_at": null
    }
  ]
}
```

---

### 4.5 内容统计

#### 4.5.1 内容统计
```
GET /api/v1/contents/:id/stats
```
**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "content_id": 1,
    "title": "员工手册2026",
    "view_count": 500,
    "download_count": 150,
    "recent_views": [
      { "user_id": "user-001", "viewed_at": "2026-03-20T10:00:00Z" }
    ]
  }
}
```

---

## 5. 流程图

### 5.1 内容分发流程

```
管理员上传文件并发布
    │
    ▼
管理员创建分发任务
    │
    ▼
POST /api/v1/contents/distributions
    │
    ▼
创建 content_distributions 记录
    │
    ▼
目标设备/用户收到推送通知
    │
    ▼
用户点击查看/下载
    │
    ▼
检查权限（access_level + content_permissions）
    │
    ├─→ 有权限 ──► 允许查看/下载 ──► 更新 view_count / download_count
    │
    └─→ 无权限 ──► 返回 403 Forbidden
```

---

## 6. 模块联动

| 联动模块 | 联动方式 | 说明 |
|----------|----------|------|
| 设备管理 | 内容分发到设备 | target_type=device |
| 会员管理 | 内容分发到用户 | target_type=user，用户查看记录关联会员 |
| 组织架构 | 内容分发到组织单元 | target_type=org_unit |
| 策略管理 | 策略可包含内容访问权限 | 机密文件受策略保护 |
| 通知系统 | 分发时发送推送通知 | 通知用户有新内容 |

---

## 7. 验收标准

### P0 验收标准

| 用例 | 验收条件 | 测试方法 |
|------|----------|----------|
| 文件上传 | 文件正确创建，状态 draft | POST /contents 后验证 |
| 文件发布 | 发布后状态变为 published | POST /contents/:id/publish 验证 |
| 内容分发 | 分发后目标收到通知 | 创建分发后检查通知 |
| 权限检查 | 无权限用户无法访问 | 用无权限账号访问文件 |

### P1 验收标准

| 用例 | 验收条件 | 测试方法 |
|------|----------|----------|
| 分类管理 | 分类树正确展示 | 创建多级分类验证 |
| 标签管理 | 文件可打多个标签 | 添加标签后查询验证 |
| 访问统计 | 下载/阅读次数正确累加 | 下载后检查 count |

### P2 验收标准

| 用例 | 验收条件 | 测试方法 |
|------|----------|----------|
| 过期权限 | 过期后无法访问 | 设置过期时间后验证 |

---

## 8. UI设计指引

### 页面结构
- **左侧菜单**：内容管理 → 文件库 / 内容分发 / 分类管理 / 内容统计
- **顶部区域**：统计卡片（总文件数 / 本月新增 / 下载次数 / 阅读次数）
- **中间区域**：Tab 页签：文件库 / 内容分发 / 分类管理
- **底部区域**：分页组件

### 组件选用
| 组件 | 用途 |
|------|------|
| a-table | 文件列表、分发任务列表 |
| a-card | 顶部统计卡片 |
| a-tabs | Tab 切换：文件库 / 内容分发 / 分类管理 |
| a-drawer | 上传文件、创建分发任务、权限配置 |
| a-upload | 文件上传组件 |
| a-tree | 分类树展示 |
| a-tag | 文件类型标签、状态标签 |
| a-icon | 文件类型图标（PDF/PPT/Word等）|
| a-progress | 分发进度展示 |
| a-empty | 空数据状态 |

### 参考模板
```
┌──────────────────────────────────────────────────────────────┐
│  [统计卡片]  总文件:200  本月新增:15  下载:1,500  阅读:5,000   │
├──────────────────────────────────────────────────────────────┤
│  [Tab: 文件库 | 内容分发 | 分类管理]                           │
├──────────────────────────────────────────────────────────────┤
│  【文件库 Tab】                                              │
│  ┌──────────────────────────────────────────────────────┐   │
│  │ [分类▼] [类型▼] [状态▼] [+上传文件]                    │   │
│  ├──────────────────────────────────────────────────────┤   │
│  │ 文件名         │类型│大小  │分类   │下载│状态│操作      │   │
│  │ 员工手册2026  │PDF │5MB   │制度   │150 │已发布│详情分发│   │
│  │ 产品培训PPT   │PPTX│20MB  │培训   │ 80 │已发布│详情分发│   │
│  └──────────────────────────────────────────────────────┘   │
└──────────────────────────────────────────────────────────────┘
```

### 交互流程
```
文件库页
    │
    ├── 点击「上传文件」──► a-drawer ──► 选择文件+填写信息 ──► 上传
    │
    ├── 点击「分发」──► a-drawer ──► 选择目标 ──► 创建分发任务
    │
    └── 点击「权限配置」──► a-drawer ──► 配置访问权限

内容分发页
    │
    ├── 查看分发任务列表
    │
    └── 查看分发状态统计
```

### 关键状态显示
- **文件状态**：a-tag，草稿=灰，已发布=绿，已归档=蓝
- **文件类型**：a-tag，PDF=红，Word=蓝，Excel=绿，PPT=橙，图片=紫
- **分发状态**：a-tag，进行中=蓝，已完成=绿，已取消=灰

---

## 修订记录

| 版本 | 日期 | 修订人 | 修订内容 |
|------|------|--------|----------|
| V1.3 | 2026-03-20 | agentcp | 全新模块，基于新增功能需求重建 |