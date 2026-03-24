# 家庭相册 PRD

## 1. 功能概述
家庭相册模块为家庭成员提供共享的宠物照片和视频存储空间，支持按时间/事件分类、家庭成员协作管理。

## 2. 页面布局与交互

### 页面路径
`/family/album` → `FamilyAlbumView.vue`

### 相册视图
- 时间轴视图/网格视图切换
- 照片墙
- 上传入口

### 照片详情
- 大图预览
- 拍摄时间/地点
- 所属成员
- 标签管理
- 下载/删除按钮

### 家庭成员管理
- 成员列表
- 权限配置（上传/编辑/删除）

## 3. API 契约

### 相册列表
- 路径：`GET /api/v1/family/albums`
- 参数：`family_id`, `page`, `page_size`
- 响应：
```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "id": 1,
        "media_url": "https://...",
        "thumbnail_url": "https://...",
        "media_type": "photo",
        "taken_at": "2024-01-01T08:00:00Z",
        "location": "家里",
        "owner_id": 100,
        "tags": ["可爱", "吃饭"]
      }
    ],
    "total": 200,
    "page": 1,
    "page_size": 20
  }
}
```

### 上传照片
- 路径：`POST /api/v1/family/albums`
- 请求体：`{ "family_id": 1, "media_url": "...", "media_type": "photo", "taken_at": "...", "location": "...", "tags": [...] }`

### 删除照片
- 路径：`DELETE /api/v1/family/albums/:id`

### 更新标签
- 路径：`PUT /api/v1/family/albums/:id`
- 请求体：`{ "tags": ["可爱", "玩耍"] }`

### 家庭成员管理
- 路径：`GET /api/v1/family/members`
- 路径：`POST /api/v1/family/members`
- 路径：`DELETE /api/v1/family/members/:id`

## 4. 验收标准
- [ ] 照片列表加载正常
- [ ] 上传照片成功
- [ ] 大图预览正常
- [ ] 标签管理正常
- [ ] 家庭成员权限正常
- [ ] 删除照片正常
