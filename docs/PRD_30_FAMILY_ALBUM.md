# PRD：家庭相册（AI筛选版）

**版本：** V1.0
**所属Phase：** Phase 4（Sprint 29-30）
**优先级：** P2
**负责角色：** agentcp（产品）、agenthd（后端）、agentqd（前端）

---

## 一、概述

### 1.1 模块定位

家庭相册模块为宠物主人提供宠物视角照片和视频的智能管理与分享功能，通过AI自动筛选高光时刻（精彩瞬间、可爱表情、互动场景），帮助用户轻松整理和回顾与宠物的美好回忆。

### 1.2 核心价值

- **AI自动筛选**：自动识别精彩瞬间、可爱表情、互动场景
- **宠物视角**：管理设备拍摄的宠物视角照片
- **家庭共享**：家庭成员共同维护相册
- **一键导出**：重要回忆导出到本地

### 1.3 范围边界

**包含：**
- 照片/视频上传和管理
- AI自动筛选（精彩瞬间、表情、互动）
- 家庭共享和权限管理
- 相册分类（按时间/事件/宠物）

**不包含：**
- 第三方云存储接入（已有云同步模块）
- 照片编辑和滤镜（可扩展）

---

## 二、功能详情

### 2.1 照片管理

| 功能 | 说明 |
|------|------|
| 照片上传 | 从设备端或Web端上传照片/视频 |
| 照片列表 | 按时间线展示所有照片 |
| 照片详情 | 查看大图、播放视频 |
| 照片删除 | 删除不需要的照片 |
| 照片下载 | 下载照片到本地 |

### 2.2 AI筛选

| 功能 | 说明 |
|------|------|
| 精彩瞬间筛选 | AI识别宠物高光时刻（跳跃、奔跑、打哈欠等） |
| 可爱表情筛选 | 识别宠物可爱表情（眯眼、歪头、吐舌头等） |
| 互动场景筛选 | 识别主人与宠物互动场景 |
| 自定义筛选 | 用户自定义标签进行筛选 |

### 2.3 相册分类

| 功能 | 说明 |
|------|------|
| 时间线视图 | 按日/月/年组织照片 |
| 事件相册 | 按「生日」「出游」「日常」等事件分类 |
| 宠物相册 | 多宠物家庭按宠物分别管理 |
| AI专辑 | AI自动生成的精选专辑 |

### 2.4 家庭共享

| 功能 | 说明 |
|------|------|
| 成员邀请 | 邀请家庭成员加入相册 |
| 权限管理 | 设置成员为「可浏览」或「可上传」 |
| 共享通知 | 新照片上传后通知家庭成员 |

---

## 三、API接口定义

### 3.1 照片管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/album/photos | 照片列表 |
| POST | /api/v1/album/photos | 上传照片 |
| GET | /api/v1/album/photos/:id | 照片详情 |
| DELETE | /api/v1/album/photos/:id | 删除照片 |
| GET | /api/v1/album/photos/:id/download | 下载照片 |
| GET | /api/v1/album/photos/export | 批量导出照片 |

### 3.2 AI筛选

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/album/photos/ai-filter | AI筛选照片列表 |
| GET | /api/v1/album/ai-albums | AI生成的精选专辑列表 |
| GET | /api/v1/album/ai-albums/:id | AI专辑详情 |
| POST | /api/v1/album/ai-albums/:id/regenerate | 重新生成AI专辑 |
| POST | /api/v1/album/photos/:id/tags | 手动打标签 |

### 3.3 相册分类

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/album/categories | 相册分类列表 |
| POST | /api/v1/album/categories | 创建分类 |
| GET | /api/v1/album/categories/:id/photos | 分类下的照片 |
| PUT | /api/v1/album/categories/:id | 更新分类 |
| DELETE | /api/v1/album/categories/:id | 删除分类 |

### 3.4 家庭共享

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/album/family/members | 家庭成员列表 |
| POST | /api/v1/album/family/invite | 邀请成员 |
| DELETE | /api/v1/album/family/members/:id | 移除成员 |
| PUT | /api/v1/album/family/members/:id/permission | 更新成员权限 |

---

## 四、数据库设计

### 4.1 照片表 (album_photos)

```sql
CREATE TABLE album_photos (
    id                  BIGSERIAL PRIMARY KEY,
    user_id             BIGINT NOT NULL,
    device_id           VARCHAR(100),
    pet_id              BIGINT,
    file_url            VARCHAR(500) NOT NULL,
    thumbnail_url       VARCHAR(500),
    media_type          VARCHAR(20) NOT NULL,         -- 'photo'/'video'
    ai_tags             VARCHAR(100)[],               -- AI识别标签
    ai_highlight_score  DECIMAL(5,2),                 -- 精彩瞬间评分
    ai_emotion_score    DECIMAL(5,2),                 -- 可爱表情评分
    ai_interaction_score DECIMAL(5,2),                -- 互动场景评分
    taken_at            TIMESTAMP,                    -- 拍摄时间
    file_size           BIGINT,
    width               INT,
    height              INT,
    is_family_shared    BOOLEAN DEFAULT FALSE,
    created_at          TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_album_photos_user ON album_photos(user_id, created_at DESC);
CREATE INDEX idx_album_photos_pet ON album_photos(pet_id);
CREATE INDEX idx_album_photos_ai_tags ON album_photos USING GIN(ai_tags);
CREATE INDEX idx_album_photos_highlight ON album_photos(ai_highlight_score DESC);
```

### 4.2 AI精选专辑表 (album_ai_albums)

```sql
CREATE TABLE album_ai_albums (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL,
    album_name      VARCHAR(255) NOT NULL,
    album_type      VARCHAR(50) NOT NULL,              -- 'highlight'/'emotion'/'interaction'/'custom'
    cover_photo_id  BIGINT REFERENCES album_photos(id),
    photo_count     INT DEFAULT 0,
    description     TEXT,
    generated_at    TIMESTAMP DEFAULT NOW(),
    created_at      TIMESTAMP DEFAULT NOW()
);
```

### 4.3 相册分类表 (album_categories)

```sql
CREATE TABLE album_categories (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL,
    category_name   VARCHAR(255) NOT NULL,
    category_type   VARCHAR(50) NOT NULL,              -- 'event'/'pet'/'custom'
    cover_url       VARCHAR(500),
    pet_id          BIGINT,
    photo_count     INT DEFAULT 0,
    created_at      TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_album_categories_user ON album_categories(user_id);
```

### 4.4 家庭成员表 (album_family_members)

```sql
CREATE TABLE album_family_members (
    id              BIGSERIAL PRIMARY KEY,
    album_id        BIGINT NOT NULL,
    user_id         BIGINT NOT NULL,
    role            VARCHAR(20) DEFAULT 'viewer',      -- 'viewer'/'uploader'/'admin'
    invited_at      TIMESTAMP DEFAULT NOW(),
    joined_at       TIMESTAMP,
    status          VARCHAR(20) DEFAULT 'pending'
);
```

---

## 五、前端页面

### 5.1 照片管理

| 页面 | 路由 | 说明 |
|------|------|------|
| 相册首页 | /album | 时间线照片流 |
| 照片详情 | /album/photos/:id | 照片大图/视频播放 |
| 上传照片 | /album/upload | 批量上传 |

### 5.2 AI筛选

| 页面 | 路由 | 说明 |
|------|------|------|
| 精彩瞬间 | /album/highlights | AI筛选的精彩瞬间 |
| 可爱表情 | /album/emotions | AI筛选的可爱表情 |
| AI专辑 | /album/ai-albums | AI生成的精选专辑 |
| 自定义筛选 | /album/filter | 按标签筛选照片 |

### 5.3 相册分类

| 页面 | 路由 | 说明 |
|------|------|------|
| 分类管理 | /album/categories | 查看/创建/编辑分类 |
| 分类详情 | /album/categories/:id | 分类下的照片 |

### 5.4 家庭共享

| 页面 | 路由 | 说明 |
|------|------|------|
| 家庭成员 | /album/family | 成员管理 |
| 分享设置 | /album/family/settings | 共享权限配置 |

---

## 六、验收标准

### 6.1 照片管理

| 验收点 | 标准 |
|--------|------|
| 照片上传 | 单次最多上传50张照片，支持JPG/PNG/MP4 |
| 照片加载 | 照片列表首屏加载<1秒 |
| 照片下载 | 支持批量下载，ZIP压缩包 |

### 6.2 AI筛选

| 验收点 | 标准 |
|--------|------|
| 精彩瞬间识别 | 识别准确率>75% |
| 可爱表情识别 | 识别准确率>80% |
| 互动场景识别 | 识别准确率>70% |
| AI专辑生成 | 新照片入库后24小时内自动更新专辑 |

### 6.3 相册分类

| 验收点 | 标准 |
|--------|------|
| 时间线 | 按日/月切换正常 |
| 事件分类 | 创建/编辑/删除正常 |

### 6.4 家庭共享

| 验收点 | 标准 |
|--------|------|
| 成员邀请 | 邀请链接有效期24小时 |
| 权限控制 | 上传权限可精准控制 |
