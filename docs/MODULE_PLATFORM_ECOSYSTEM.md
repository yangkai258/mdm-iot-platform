# 模块 PRD：开放平台生态（MODULE_PLATFORM_ECOSYSTEM）

**版本：** V1.0
**所属Phase：** Phase 4（Sprint 25-26）
**优先级：** P2
**负责角色：** agentcp（产品）、agenthd（后端）、agentqd（前端）

---

## 一、概述

### 1.1 模块定位

开放平台生态模块是MDM控制中台的生态扩展核心，通过开放API、插件市场、内容市场（表情包/动作/声音）、第三方集成，构建完整的开发者生态和内容生态，扩展产品边界，提升用户粘性和商业价值。

### 1.2 核心价值

- **生态扩展**：通过开放平台引入第三方开发者
- **内容丰富**：通过市场引入丰富的表情/动作/声音内容
- **商业变现**：通过市场分成模式创造新的收入来源

### 1.3 范围边界

**包含：**
- 开发者API管理
- App/插件市场
- 表情包市场
- 动作资源库
- 声音定制
- SDK发布
- 第三方集成

**不包含：**
- 订阅管理（MODULE_SUBSCRIPTION）
- 支付网关（支付平台）
- 设备固件（固件层）

---

## 二、功能详情

### 2.1 开发者API

#### 2.1.1 API分类

| API类型 | 说明 | 认证方式 |
|---------|------|----------|
| 设备管理API | 设备的增删改查 | API Key |
| 设备控制API | 下发指令/获取状态 | API Key |
| 订阅API | 订阅管理/用量查询 | OAuth 2.0 |
| Webhook API | 接收平台事件 | 签名验证 |
| 数据API | 设备数据/用户数据 | OAuth 2.0 |

#### 2.1.2 开发者门户

| 功能 | 说明 |
|------|------|
| 注册登录 | 开发者账号注册/登录 |
| 开发者认证 | 个人/企业认证 |
| 应用管理 | 创建/管理应用 |
| API Key管理 | 生成/管理API Key |
| 使用统计 | API调用统计 |
| 文档中心 | 完整的API文档 |
| SDK下载 | 各语言SDK下载 |

#### 2.1.3 API管理

| 功能 | 说明 |
|------|------|
| 限流配置 | 按API类型配置限流 |
| 配额管理 | 每日/每月调用配额 |
| 配额预警 | 配额使用超过80%预警 |
| 配额购买 | 购买额外API配额 |
| 审计日志 | API调用详细日志 |

### 2.2 App/插件市场

#### 2.2.1 插件类型

| 类型 | 说明 | 示例 |
|------|------|------|
| 功能插件 | 扩展设备功能 | 夜灯模式/温度传感 |
| 娱乐插件 | 娱乐内容 | 游戏/音乐 |
| 工具插件 | 效率工具 | 定时任务/自动化 |
| 主题插件 | UI主题 | 皮肤/图标包 |

#### 2.2.2 插件管理

| 功能 | 说明 |
|------|------|
| 插件开发 | 提供SDK和文档 |
| 插件上传 | 上传插件包 |
| 插件审核 | 管理员审核插件 |
| 插件发布 | 审核通过后发布 |
| 插件更新 | 开发者推送更新 |
| 插件下架 | 管理员/开发者下架 |

#### 2.2.3 插件市场

| 功能 | 说明 |
|------|------|
| 插件浏览 | 按分类浏览插件 |
| 插件搜索 | 关键词搜索 |
| 插件详情 | 查看插件详情 |
| 插件安装 | 一键安装到设备 |
| 插件评分 | 用户评分和评论 |
| 插件推荐 | 平台推荐插件 |

### 2.3 表情包市场

#### 2.3.1 表情包类型

| 类型 | 说明 | 价格 |
|------|------|------|
| 内置表情 | 预置表情包 | 免费 |
| 官方表情 | 官方出品 | ¥1-5/套 |
| 创作者表情 | 用户创作 | ¥0-10/套 |
| 节日表情 | 节日限定 | 免费/限时 |

#### 2.3.2 表情包管理

| 功能 | 说明 |
|------|------|
| 表情制作 | 提供表情制作工具 |
| 表情上传 | 上传表情包 |
| 表情审核 | 审核表情包内容 |
| 表情发布 | 上架表情市场 |
| 表情更新 | 更新表情内容 |

#### 2.3.3 表情推荐

| 功能 | 说明 |
|------|------|
| 场景推荐 | 基于场景推荐表情 |
| 情感推荐 | 基于情绪推荐表情 |
| 热门推荐 | 热门表情推荐 |
| 个性化推荐 | 基于用户喜好推荐 |

### 2.4 动作资源库

#### 2.4.1 动作类型

| 类型 | 说明 | 价格 |
|------|------|------|
| 内置动作 | 预置动作库 | 免费 |
| 官方动作 | 官方出品 | ¥2-10/个 |
| 创作者动作 | 用户创作 | ¥0-20/个 |
| 动作包 | 动作集合 | 折扣销售 |

#### 2.4.2 动作编辑器

| 功能 | 说明 |
|------|------|
| 动作录制 | 录制用户动作 |
| 动作编辑 | 调整动作参数 |
| 动作预览 | 实时预览效果 |
| 动作导出 | 导出为可执行文件 |

#### 2.4.3 动作市场

| 功能 | 说明 |
|------|------|
| 动作浏览 | 分类浏览动作 |
| 动作搜索 | 关键词搜索 |
| 动作详情 | 查看动作详情 |
| 动作购买 | 购买动作 |
| 动作下载 | 下载到本地/设备 |
| 动作评分 | 用户评分和评论 |

### 2.5 声音定制

#### 2.5.1 声音类型

| 类型 | 说明 | 价格 |
|------|------|------|
| 内置声音 | 预置声音 | 免费 |
| 官方声音 | 官方TTS声音 | ¥5-20/个 |
| 定制声音 | 克隆定制声音 | ¥50-100/个 |

#### 2.5.2 TTS个性化

| 功能 | 说明 |
|------|------|
| 声音预览 | 预览各声音效果 |
| 参数调整 | 调整语速/音调 |
| 声音购买 | 购买声音使用权 |
| 声音管理 | 管理已购声音 |

#### 2.5.3 语音克隆

| 功能 | 说明 |
|------|------|
| 录音采集 | 采集用户声音样本 |
| 声音训练 | 训练个性化声音 |
| 声音预览 | 预览克隆效果 |
| 声音发布 | 发布到声音市场 |

### 2.6 SDK发布

#### 2.6.1 SDK类型

| SDK | 说明 | 平台 |
|------|------|------|
| iOS SDK | iOS设备控制 | iOS 14+ |
| Android SDK | Android设备控制 | Android 8+ |
| 微信小程序SDK | 小程序集成 | 微信 |
| JavaScript SDK | Web端集成 | Browser |
| Python SDK | Python脚本集成 | Python 3.8+ |
| Go SDK | Go应用集成 | Go 1.18+ |

#### 2.6.2 SDK管理

| 功能 | 说明 |
|------|------|
| SDK下载 | 各平台SDK下载 |
| 快速开始 | 快速集成指南 |
| 示例代码 | 完整示例代码 |
| API文档 | API接口文档 |
| 更新日志 | SDK更新历史 |

### 2.7 第三方集成

#### 2.7.1 智能家居对接

| 功能 | 说明 | 优先级 |
|------|------|--------|
| 米家授权 | OAuth2.0授权米家账号 | P2 |
| 米家设备同步 | 同步米家设备列表和状态 | P2 |
| 设备控制 | 通过宠物行为触发米家设备 | P2 |
| 天猫精灵授权 | 授权天猫精灵控制宠物 | P2 |
| 语音指令 | 天猫精灵语音指令响应 | P2 |
| HomeKit配对 | MFi设备配对和控制 | P3 |
| Google Home | Google Assistant语音控制 | P3 |
| 智能联动 | 配置宠物行为触发智能家居规则 | P2 |

#### 2.7.2 宠物医疗对接

| 功能 | 说明 | 优先级 |
|------|------|--------|
| 医院查询 | 合作宠物医院查询和展示 | P2 |
| 预约挂号 | 在线预约宠物医生 | P2 |
| 预约管理 | 查看/修改/取消预约 | P2 |
| 病历同步 | 医院病历同步到宠物档案 | P3 |
| 疫苗管理 | 疫苗接种记录和提醒 | P2 |
| 健康报告 | 宠物健康数据报告 | P3 |
| 紧急联系 | 一键联系最近的宠物医院 | P3 |

#### 2.7.3 宠物保险对接

| 功能 | 说明 | 优先级 |
|------|------|--------|
| 保险方案展示 | 合作保险公司方案对比 | P3 |
| 在线投保 | 保险产品在线购买 | P3 |
| 保单管理 | 保单查询和状态管理 | P3 |
| 理赔申请 | 在线发起理赔并提交材料 | P3 |
| 理赔进度 | 理赔进度查询 | P3 |
| 宠物档案 | 宠物健康档案授权给保险公司 | P3 |

#### 2.7.4 地图服务对接

| 功能 | 说明 | 优先级 |
|------|------|--------|
| 室内地图 | 室内地图瓦片服务 | P2 |
| 室内导航 | 室内定位和导航 | P2 |
| 地址编码 | 经纬度与地址互转 | P3 |
| 路线规划 | 户外路线规划 | P3 |
| 周边搜索 | 宠物医院/公园等周边搜索 | P3 |
| 位置上报 | 宠物位置上报和历史轨迹 | P2 |

#### 2.7.5 宠物用品电商对接

| 功能 | 说明 | 优先级 |
|------|------|--------|
| 商品推荐 | 基于宠物特征推荐用品 | P2 |
| 购物车同步 | 购物车数据同步 | P3 |
| 一键购买 | 对接京东到家/饿了么 | P3 |
| 订单管理 | 订单查询和物流跟踪 | P3 |
| 食品推荐 | 基于体重/年龄/品种推荐 | P2 |
| 玩具推荐 | 基于宠物性格推荐 | P3 |

---

## 三、API接口定义

### 3.1 开发者API

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/v1/developer/apps | 创建应用 |
| GET | /api/v1/developer/apps | 应用列表 |
| GET | /api/v1/developer/apps/:id | 应用详情 |
| PUT | /api/v1/developer/apps/:id | 更新应用 |
| DELETE | /api/v1/developer/apps/:id | 删除应用 |
| POST | /api/v1/developer/apps/:id/keys | 生成API Key |
| DELETE | /api/v1/developer/apps/:id/keys/:key_id | 删除API Key |
| GET | /api/v1/developer/stats | 使用统计 |
| GET | /api/v1/developer/audit-logs | 审计日志 |

### 3.2 插件市场

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/marketplace/plugins | 插件列表 |
| GET | /api/v1/marketplace/plugins/:id | 插件详情 |
| POST | /api/v1/developer/plugins | 提交插件 |
| PUT | /api/v1/developer/plugins/:id | 更新插件 |
| DELETE | /api/v1/developer/plugins/:id | 删除插件 |
| POST | /api/v1/developer/plugins/:id/publish | 发布插件 |
| POST | /api/v1/plugins/:id/install | 安装插件 |
| POST | /api/v1/plugins/:id/uninstall | 卸载插件 |
| POST | /api/v1/plugins/:id/rate | 评分 |

### 3.3 表情包市场

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/marketplace/emotions | 表情包列表 |
| GET | /api/v1/marketplace/emotions/:id | 表情包详情 |
| POST | /api/v1/developer/emotions | 上传表情包 |
| PUT | /api/v1/developer/emotions/:id | 更新表情包 |
| POST | /api/v1/developer/emotions/:id/publish | 发布表情包 |
| POST | /api/v1/emotions/:id/purchase | 购买表情包 |
| POST | /api/v1/emotions/:id/rate | 评分 |

### 3.4 动作市场

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/marketplace/actions | 动作列表 |
| GET | /api/v1/marketplace/actions/:id | 动作详情 |
| POST | /api/v1/developer/actions | 上传动作 |
| PUT | /api/v1/developer/actions/:id | 更新动作 |
| POST | /api/v1/developer/actions/:id/publish | 发布动作 |
| POST | /api/v1/actions/:id/purchase | 购买动作 |
| POST | /api/v1/actions/:id/download | 下载动作 |
| POST | /api/v1/actions/:id/rate | 评分 |

### 3.5 声音市场

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/marketplace/voices | 声音列表 |
| GET | /api/v1/marketplace/voices/:id | 声音详情 |
| POST | /api/v1/developer/voices | 上传声音 |
| POST | /api/v1/voices/:id/purchase | 购买声音 |
| POST | /api/v1/voices/:id/preview | 预览声音 |
| POST | /api/v1/voices/clone/start | 开始语音克隆 |
| GET | /api/v1/voices/clone/:task_id | 克隆进度 |

### 3.6 第三方集成

#### 3.6.1 智能家居对接

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/integrations/smarthome/devices | 同步设备列表 |
| GET | /api/v1/integrations/smarthome/devices/:id | 设备详情 |
| POST | /api/v1/integrations/smarthome/devices/:id/control | 控制设备 |
| POST | /api/v1/integrations/smarthome/triggers | 智能联动配置 |
| GET | /api/v1/integrations/smarthome/triggers/:id | 联动规则详情 |
| PUT | /api/v1/integrations/smarthome/triggers/:id | 更新联动规则 |
| DELETE | /api/v1/integrations/smarthome/triggers/:id | 删除联动规则 |
| POST | /api/v1/integrations/mi-home/auth | 米家授权 |
| DELETE | /api/v1/integrations/mi-home/auth | 米家取消授权 |
| GET | /api/v1/integrations/mi-home/devices | 米家设备列表 |
| POST | /api/v1/integrations/tmall-genie/auth | 天猫精灵授权 |
| GET | /api/v1/integrations/tmall-genie/devices | 天猫设备列表 |
| POST | /api/v1/integrations/homekit/pair | HomeKit配对 |
| GET | /api/v1/integrations/homekit/devices | HomeKit设备列表 |
| POST | /api/v1/integrations/google-home/auth | Google Home授权 |
| GET | /api/v1/integrations/google-home/devices | Google Home设备 |

#### 3.6.2 宠物医疗对接

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/integrations/vet/hospitals | 合作医院列表 |
| GET | /api/v1/integrations/vet/hospitals/:id | 医院详情 |
| POST | /api/v1/integrations/vet/appointments | 创建预约 |
| GET | /api/v1/integrations/vet/appointments/:id | 预约详情 |
| PUT | /api/v1/integrations/vet/appointments/:id | 更新预约 |
| DELETE | /api/v1/integrations/vet/appointments/:id | 取消预约 |
| GET | /api/v1/integrations/vet/appointments | 预约列表 |
| POST | /api/v1/integrations/vet/records/sync | 病历同步 |
| GET | /api/v1/integrations/vet/records/:pet_id | 宠物病历 |
| POST | /api/v1/integrations/vet/vaccines/reminder | 设置疫苗提醒 |
| GET | /api/v1/integrations/vet/vaccines/:pet_id | 疫苗记录 |
| POST | /api/v1/integrations/vet/health/check | 健康检查预约 |

#### 3.6.3 宠物保险对接

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/integrations/insurance/plans | 保险方案列表 |
| GET | /api/v1/integrations/insurance/plans/:id | 方案详情 |
| POST | /api/v1/integrations/insurance/policies | 投保申请 |
| GET | /api/v1/integrations/insurance/policies/:id | 保单详情 |
| GET | /api/v1/integrations/insurance/policies | 保单列表 |
| POST | /api/v1/integrations/insurance/claims | 发起理赔 |
| GET | /api/v1/integrations/insurance/claims/:id | 理赔详情 |
| GET | /api/v1/integrations/insurance/claims | 理赔列表 |
| PUT | /api/v1/integrations/insurance/claims/:id/upload | 上传理赔材料 |
| GET | /api/v1/integrations/insurance/claims/:id/status | 理赔状态 |

#### 3.6.4 地图服务对接

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/integrations/maps/tiles/:z/:x/:y | 地图瓦片 |
| POST | /api/v1/integrations/maps/geocode | 地址编码 |
| POST | /api/v1/integrations/maps/reverse-geocode | 逆地址编码 |
| POST | /api/v1/integrations/maps/route | 路线规划 |
| GET | /api/v1/integrations/maps/poi/search | 周边搜索 |
| GET | /api/v1/integrations/maps/poi/:id | POI详情 |
| POST | /api/v1/integrations/maps/indoor/map | 室内地图数据 |
| GET | /api/v1/integrations/maps/indoor/:floor_id | 室内地图查询 |
| POST | /api/v1/integrations/maps/location/report | 位置上报 |
| GET | /api/v1/integrations/maps/location/history | 位置历史 |

#### 3.6.5 宠物用品电商对接

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/integrations/ecommerce/products | 商品推荐 |
| GET | /api/v1/integrations/ecommerce/products/:id | 商品详情 |
| POST | /api/v1/integrations/ecommerce/cart | 加入购物车 |
| GET | /api/v1/integrations/ecommerce/cart | 购物车列表 |
| POST | /api/v1/integrations/ecommerce/orders | 创建订单 |
| GET | /api/v1/integrations/ecommerce/orders/:id | 订单详情 |
| GET | /api/v1/integrations/ecommerce/food/recommend | 食品推荐 |
| GET | /api/v1/integrations/ecommerce/toy/recommend | 玩具推荐 |

---

## 四、数据库设计

### 4.1 开发者应用表 (developer_apps)

```sql
CREATE TABLE developer_apps (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES users(id),
    app_name        VARCHAR(255) NOT NULL,
    app_type        VARCHAR(50),                  -- 'personal'/'enterprise'
    description     TEXT,
    app_icon        VARCHAR(500),
    website_url     VARCHAR(500),
    callback_url    VARCHAR(500),
    status          VARCHAR(20) DEFAULT 'active',  -- 'active'/'suspended'/'deleted'
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);
```

### 4.2 API Key表 (api_keys)

```sql
CREATE TABLE api_keys (
    id              BIGSERIAL PRIMARY KEY,
    app_id          BIGINT NOT NULL REFERENCES developer_apps(id) ON DELETE CASCADE,
    key_prefix      VARCHAR(20) NOT NULL,
    key_hash        VARCHAR(255) NOT NULL,
    key_type        VARCHAR(20) NOT NULL,         -- 'api_key'/'oauth_client'
    scopes          VARCHAR(100)[],                -- 权限范围
    rate_limit      INT DEFAULT 1000,              -- 每分钟限制
    is_active       BOOLEAN DEFAULT TRUE,
    last_used_at    TIMESTAMP,
    expires_at      TIMESTAMP,
    created_at      TIMESTAMP DEFAULT NOW()
);
```

### 4.3 插件表 (plugins)

```sql
CREATE TABLE plugins (
    id              BIGSERIAL PRIMARY KEY,
    developer_id    BIGINT NOT NULL REFERENCES users(id),
    plugin_name     VARCHAR(255) NOT NULL,
    plugin_type     VARCHAR(50) NOT NULL,         -- 'feature'/'entertainment'/'tool'/'theme'
    description     TEXT,
    icon_url        VARCHAR(500),
    screenshots     VARCHAR(500)[],
    version         VARCHAR(20),
    package_url     VARCHAR(500),                  -- 插件包地址
    price           DECIMAL(10,2) DEFAULT 0,
    is_free         BOOLEAN DEFAULT TRUE,
    status          VARCHAR(20) DEFAULT 'draft',   -- 'draft'/'pending'/'approved'/'rejected'/'published'/'removed'
    reviewed_at     TIMESTAMP,
    reviewed_by    BIGINT,
    published_at    TIMESTAMP,
    downloads       INT DEFAULT 0,
    rating_avg      DECIMAL(3,2) DEFAULT 0,
    rating_count    INT DEFAULT 0,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);
```

### 4.4 表情包表 (emotion_packs)

```sql
CREATE TABLE emotion_packs (
    id              BIGSERIAL PRIMARY KEY,
    developer_id    BIGINT NOT NULL REFERENCES users(id),
    pack_name       VARCHAR(255) NOT NULL,
    pack_type       VARCHAR(20) NOT NULL,         -- 'built-in'/'official'/'creator'/'seasonal'
    description     TEXT,
    thumbnail_url   VARCHAR(500),
    preview_url     VARCHAR(500),
    emotions        JSONB,                          -- 表情列表
    price           DECIMAL(10,2) DEFAULT 0,
    is_free         BOOLEAN DEFAULT TRUE,
    status          VARCHAR(20) DEFAULT 'draft',
    downloads       INT DEFAULT 0,
    rating_avg      DECIMAL(3,2) DEFAULT 0,
    rating_count    INT DEFAULT 0,
    tags            VARCHAR(100)[],
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);
```

### 4.5 动作表 (actions)

```sql
CREATE TABLE actions (
    id              BIGSERIAL PRIMARY KEY,
    developer_id    BIGINT NOT NULL REFERENCES users(id),
    action_name     VARCHAR(255) NOT NULL,
    description     TEXT,
    action_type     VARCHAR(20) NOT NULL,         -- 'built-in'/'official'/'creator'
    difficulty      VARCHAR(20),                   -- 'easy'/'medium'/'hard'
    thumbnail_url   VARCHAR(500),
    video_url       VARCHAR(500),
    motion_data     JSONB,                          -- 动作运动数据
    price           DECIMAL(10,2) DEFAULT 0,
    is_free         BOOLEAN DEFAULT TRUE,
    status          VARCHAR(20) DEFAULT 'draft',
    downloads       INT DEFAULT 0,
    rating_avg      DECIMAL(3,2) DEFAULT 0,
    rating_count    INT DEFAULT 0,
    tags            VARCHAR(100)[],
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);
```

### 4.6 声音表 (voices)

```sql
CREATE TABLE voices (
    id              BIGSERIAL PRIMARY KEY,
    developer_id    BIGINT NOT NULL REFERENCES users(id),
    voice_name      VARCHAR(255) NOT NULL,
    voice_type      VARCHAR(20) NOT NULL,         -- 'built-in'/'official'/'custom'
    description     TEXT,
    preview_url     VARCHAR(500),
    audio_samples   JSONB,                          -- 示例音频
    voice_params    JSONB,                          -- 声音参数
    price           DECIMAL(10,2) DEFAULT 0,
    is_free         BOOLEAN DEFAULT TRUE,
    status          VARCHAR(20) DEFAULT 'draft',
    downloads       INT DEFAULT 0,
    rating_avg      DECIMAL(3,2) DEFAULT 0,
    rating_count    INT DEFAULT 0,
    tags            VARCHAR(100)[],
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);
```

### 4.7 用户购买表 (user_purchases)

```sql
CREATE TABLE user_purchases (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES users(id),
    item_type       VARCHAR(20) NOT NULL,         -- 'plugin'/'emotion'/'action'/'voice'
    item_id        BIGINT NOT NULL,
    price           DECIMAL(10,2) NOT NULL,
    payment_method  VARCHAR(20),
    payment_id      VARCHAR(100),
    status          VARCHAR(20) DEFAULT 'completed',
    purchased_at    TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_user_purchases_user ON user_purchases(user_id, purchased_at DESC);
```

### 4.8 评分表 (ratings)

```sql
CREATE TABLE ratings (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES users(id),
    item_type       VARCHAR(20) NOT NULL,
    item_id        BIGINT NOT NULL,
    rating          INT NOT NULL,                  -- 1-5
    review          TEXT,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_id, item_type, item_id)
);
```

### 4.9 第三方集成表 (integrations)

```sql
CREATE TABLE integrations (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES users(id),
    integration_type VARCHAR(50) NOT NULL,         -- 'mi_home'/'tmall_genie'/'homekit'/'google_home'/'vet'/'insurance'/'ecommerce'
    status          VARCHAR(20) DEFAULT 'disconnected', -- 'connected'/'disconnected'/'error'
    config          JSONB,                          -- 集成配置
    connected_at    TIMESTAMP,
    last_sync_at    TIMESTAMP,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_id, integration_type)
);

### 4.10 智能家居设备表 (smart_home_devices)

```sql
CREATE TABLE smart_home_devices (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES users(id),
    integration_id  BIGINT REFERENCES integrations(id),
    platform        VARCHAR(50) NOT NULL,          -- 'mi_home'/'tmall_genie'/'homekit'/'google_home'
    platform_device_id VARCHAR(100) NOT NULL,
    device_name     VARCHAR(255),
    device_type     VARCHAR(50),                   -- 'light'/'switch'/'sensor'/'camera'
    status          JSONB,                          -- 设备状态
    is_online       BOOLEAN DEFAULT TRUE,
    last_control_at TIMESTAMP,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW(),
    UNIQUE(integration_id, platform_device_id)
);

CREATE INDEX idx_smart_home_devices_user ON smart_home_devices(user_id);
```

### 4.11 智能联动规则表 (smart_home_triggers)

```sql
CREATE TABLE smart_home_triggers (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES users(id),
    trigger_name    VARCHAR(255) NOT NULL,
    trigger_type    VARCHAR(50) NOT NULL,          -- 'pet_action'/'schedule'/'sensor'/'voice'
    condition       JSONB NOT NULL,                 -- 触发条件
    action          JSONB NOT NULL,                 -- 执行动作
    device_id       BIGINT REFERENCES smart_home_devices(id),
    is_enabled      BOOLEAN DEFAULT TRUE,
    run_count       INT DEFAULT 0,
    last_run_at     TIMESTAMP,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_smart_home_triggers_user ON smart_home_triggers(user_id);
```

### 4.12 宠物医疗预约表 (vet_appointments)

```sql
CREATE TABLE vet_appointments (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES users(id),
    pet_id          BIGINT REFERENCES pets(id),
    hospital_id     VARCHAR(100) NOT NULL,
    hospital_name   VARCHAR(255),
    doctor_name     VARCHAR(100),
    appointment_time TIMESTAMP NOT NULL,
    appointment_type VARCHAR(50),                  -- 'checkup'/'vaccine'/'emergency'/'followup'
    reason          TEXT,
    status          VARCHAR(20) DEFAULT 'pending',  -- 'pending'/'confirmed'/'completed'/'cancelled'
    notes           TEXT,
    external_id     VARCHAR(100),
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_vet_appointments_user ON vet_appointments(user_id, appointment_time DESC);
```

### 4.13 宠物病历表 (pet_medical_records)

```sql
CREATE TABLE pet_medical_records (
    id              BIGSERIAL PRIMARY KEY,
    pet_id          BIGINT NOT NULL REFERENCES pets(id),
    record_type     VARCHAR(50) NOT NULL,           -- 'diagnosis'/'vaccine'/'surgery'/'checkup'
    hospital_name   VARCHAR(255),
    doctor_name     VARCHAR(100),
    diagnosis       TEXT,
    treatment       TEXT,
    prescription    TEXT,
    record_date     TIMESTAMP NOT NULL,
    attachments     VARCHAR(500)[],
    external_id     VARCHAR(100),
    sync_status     VARCHAR(20) DEFAULT 'synced',
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_pet_medical_records_pet ON pet_medical_records(pet_id, record_date DESC);
```

### 4.14 宠物保险保单表 (insurance_policies)

```sql
CREATE TABLE insurance_policies (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES users(id),
    pet_id          BIGINT REFERENCES pets(id),
    plan_id         VARCHAR(100) NOT NULL,
    plan_name       VARCHAR(255),
    insurer         VARCHAR(100),
    policy_number   VARCHAR(100) UNIQUE,
    premium         DECIMAL(10,2),
    coverage_amount DECIMAL(10,2),
    start_date      DATE,
    end_date        DATE,
    status          VARCHAR(20) DEFAULT 'active',
    external_id     VARCHAR(100),
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_insurance_policies_user ON insurance_policies(user_id);
```

### 4.15 保险理赔表 (insurance_claims)

```sql
CREATE TABLE insurance_claims (
    id              BIGSERIAL PRIMARY KEY,
    policy_id       BIGINT NOT NULL REFERENCES insurance_policies(id),
    claim_number    VARCHAR(100) UNIQUE,
    claim_type      VARCHAR(50),
    amount          DECIMAL(10,2) NOT NULL,
    status          VARCHAR(20) DEFAULT 'pending',
    description     TEXT,
    documents       VARCHAR(500)[],
    claim_date      TIMESTAMP,
    processed_at    TIMESTAMP,
    approved_amount DECIMAL(10,2),
    external_id     VARCHAR(100),
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_insurance_claims_policy ON insurance_claims(policy_id);
```

### 4.16 地图服务配置表 (map_configs)

```sql
CREATE TABLE map_configs (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES users(id),
    provider        VARCHAR(50) NOT NULL,
    api_key         VARCHAR(500),
    is_active       BOOLEAN DEFAULT TRUE,
    quota_used      INT DEFAULT 0,
    quota_limit     INT DEFAULT 10000,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_map_configs_user ON map_configs(user_id);
```

---

## 五、前端页面清单

### 5.1 开发者门户

| 页面 | 路由 | 说明 |
|------|------|------|
| 开发者首页 | /developer | 开发者概览 |
| 应用管理 | /developer/apps | 应用列表 |
| 创建应用 | /developer/apps/create | 创建应用 |
| 应用详情 | /developer/apps/:id | 应用详情 |
| API Key管理 | /developer/apps/:id/keys | API Key管理 |
| 使用统计 | /developer/stats | API使用统计 |
| 审计日志 | /developer/audit-logs | 调用日志 |

### 5.2 插件市场

| 页面 | 路由 | 说明 |
|------|------|------|
| 插件市场 | /marketplace/plugins | 插件列表 |
| 插件详情 | /marketplace/plugins/:id | 插件详情 |
| 开发者上传 | /developer/plugins/create | 上传插件 |
| 我的插件 | /developer/plugins | 我的插件 |
| 插件审核 | /admin/plugins/review | 审核插件 |

### 5.3 表情包市场

| 页面 | 路由 | 说明 |
|------|------|------|
| 表情包市场 | /marketplace/emotions | 表情包列表 |
| 表情包详情 | /marketplace/emotions/:id | 表情包详情 |
| 表情制作 | /developer/emotions/create | 制作表情包 |
| 我的表情包 | /developer/emotions | 我的表情包 |

### 5.4 动作市场

| 页面 | 路由 | 说明 |
|------|------|------|
| 动作市场 | /marketplace/actions | 动作列表 |
| 动作详情 | /marketplace/actions/:id | 动作详情 |
| 动作编辑器 | /developer/actions/create | 编辑动作 |
| 我的动作 | /developer/actions | 我的动作 |

### 5.5 声音市场

| 页面 | 路由 | 说明 |
|------|------|------|
| 声音市场 | /marketplace/voices | 声音列表 |
| 声音详情 | /marketplace/voices/:id | 声音详情 |
| 语音克隆 | /voices/clone | 语音克隆 |
| 我的声音 | /developer/voices | 我的声音 |

### 5.6 第三方集成

| 页面 | 路由 | 说明 |
|------|------|------|
| 集成中心 | /integrations | 集成列表 |
| 智能家居 | /integrations/smarthome | 智能家居管理 |
| 设备列表 | /integrations/smarthome/devices | 设备列表 |
| 联动规则 | /integrations/smarthome/triggers | 联动规则管理 |
| 宠物医疗 | /integrations/vet | 宠物医疗 |
| 医院查询 | /integrations/vet/hospitals | 合作医院列表 |
| 预约管理 | /integrations/vet/appointments | 预约列表 |
| 病历同步 | /integrations/vet/records/:pet_id | 病历查看 |
| 宠物保险 | /integrations/insurance | 保险管理 |
| 保险方案 | /integrations/insurance/plans | 方案对比 |
| 保单管理 | /integrations/insurance/policies | 我的保单 |
| 理赔申请 | /integrations/insurance/claims | 理赔申请 |
| 地图服务 | /integrations/maps | 地图配置 |
| 地图配置 | /integrations/maps/config | API Key配置 |
| 宠物电商 | /integrations/ecommerce | 宠物用品 |

### 5.7 管理后台

| 页面 | 路由 | 说明 |
|------|------|------|
| 市场管理 | /admin/marketplace | 市场概览 |
| 插件审核 | /admin/plugins/review | 插件审核 |
| 表情审核 | /admin/emotions/review | 表情审核 |
| 动作审核 | /admin/actions/review | 动作审核 |
| 声音审核 | /admin/voices/review | 声音审核 |

---

## 六、验收标准

### 6.1 开发者API

| 验收点 | 标准 |
|--------|------|
| 应用创建 | 开发者成功创建应用并获取API Key |
| API调用 | API按限流正常工作 |
| 文档完整性 | 文档覆盖100%的API |
| SDK可用性 | 各语言SDK可正常运行示例代码 |

### 6.2 插件市场

| 验收点 | 标准 |
|--------|------|
| 插件上传 | 开发者可成功上传插件包 |
| 审核流程 | 审核流程完整，3个工作日内完成 |
| 插件安装 | 安装成功率>99% |
| 插件市场加载 | 页面加载<2秒 |

### 6.3 表情包市场

| 验收点 | 标准 |
|--------|------|
| 表情购买 | 购买流程完整，支付成功后立即可用 |
| 场景推荐 | 推荐准确率>70% |
| 评分展示 | 评分实时更新 |

### 6.4 动作市场

| 验收点 | 标准 |
|--------|------|
| 动作下载 | 下载成功率>99% |
| 动作评分 | 评分功能正常 |
| 开发者分成 | 分成结算准确 |

### 6.5 声音市场

| 验收点 | 标准 |
|--------|------|
| 声音预览 | 预览播放流畅 |
| 语音克隆 | 克隆任务成功完成 |
| 声音使用 | 购买后立即可用于TTS |

### 6.6 第三方集成

| 验收点 | 标准 |
|--------|------|
| 米家对接 | 米家设备可被宠物触发 |
| 天猫精灵 | 语音指令正常响应 |
| 预约功能 | 预约流程完整 |
| 地图服务 | 室内地图加载正常 |
| 保险理赔 | 理赔流程可完整提交 |


---

### 6.7 地图服务

| 验收点 | 标准 |
|--------|------|
| 地图加载 | 室内地图瓦片加载<500ms |
| 位置精度 | 室内定位精度<1m |
| 地址编码 | 地址编码准确率>95% |

### 6.8 宠物保险

| 验收点 | 标准 |
|--------|------|
| 方案展示 | 保险方案信息准确 |
| 投保流程 | 投保流程完整顺畅 |
| 理赔进度 | 理赔状态实时更新 |


---

## 七、页面布局规范

### 7.1 开发者首页（/developer）

**布局结构：**
1. 面包屑 → 页面标题
2. 开发者统计卡片（应用数/API调用量/配额使用率）—— 白色
3. 快捷入口卡片网格（应用管理/API Key/文档中心/SDK下载）

### 7.2 应用管理页面（/developer/apps）

**布局结构：**
1. 面包屑 → 页面标题
2. 操作栏（创建应用靠左）
3. 应用列表表格

**按钮规范：**
- [创建应用] — 左对齐
- [详情] [编辑] [删除] — 行内右对齐

**表格列：**
| 列名 | 宽度 | 说明 |
|------|------|------|
| 序号 | 60px | - |
| 应用名称 | 200px | - |
| 应用类型 | 100px | personal/enterprise |
| API Key | 200px | - |
| 状态 | 80px | active/suspended |
| 创建时间 | 150px | - |
| 操作 | 120px | 详情/编辑/删除 |

**分页：** 右下角，10/20/50/100 条

### 7.3 插件市场页面（/marketplace/plugins）

**布局结构：**
1. 面包屑 → 页面标题
2. 筛选区（浅灰卡片 #F2F3F5）：类型 / 价格 / 评分排序
3. 操作栏（我的插件靠右）
4. 插件卡片网格

**按钮规范：**
- [我的插件] — 右对齐（开发者视图）
- [安装] [详情] — 卡片内

### 7.4 表情包市场页面（/marketplace/emotions）

**布局结构：**
1. 面包屑 → 页面标题
2. 筛选区（浅灰卡片）：类型 / 价格 / 排序
3. 操作栏（我的表情包靠右）
4. 表情包卡片网格

**按钮规范：**
- [我的表情包] — 右对齐（创作者视图）
- [购买] [详情] — 卡片内

### 7.5 动作市场页面（/marketplace/actions）

**布局结构：**
1. 面包屑 → 页面标题
2. 筛选区（浅灰卡片）：类型 / 难度 / 价格
3. 操作栏（我的动作靠右）
4. 动作卡片网格

**按钮规范：**
- [我的动作] — 右对齐（创作者视图）
- [购买] [下载] [详情] — 卡片内

### 7.6 声音市场页面（/marketplace/voices）

**布局结构：**
1. 面包屑 → 页面标题
2. 筛选区（浅灰卡片）：类型 / 价格
3. 声音列表（支持在线预览播放）

**按钮规范：**
- [预览] [购买] — 行内

**分页：** 右下角，10/20/50/100 条

### 7.7 第三方集成页面（/integrations）

**布局结构：**
1. 面包屑 → 页面标题
2. 集成分类 Tab（智能家居/宠物医疗/宠物保险/宠物用品）
3. 集成卡片网格（Logo/名称/功能描述/连接状态）

**按钮规范：**
- [连接] [配置] [断开] — 卡片内右对齐

### 7.8 管理后台页面（/admin/marketplace）

**布局结构：**
1. 面包屑 → 页面标题
2. Tab 页签：插件审核 / 表情审核 / 动作审核 / 声音审核
3. 待审核内容列表表格

**按钮规范：**
- [通过] [拒绝] — 行内右对齐

**表格列：**
| 列名 | 宽度 | 说明 |
|------|------|------|
| 序号 | 60px | - |
| 名称 | 200px | - |
| 开发者 | 150px | - |
| 提交时间 | 150px | - |
| 状态 | 100px | pending/approved/rejected |
| 操作 | 120px | 通过/拒绝 |

**分页：** 右下角，10/20/50/100 条

### 7.9 弹窗规范

| 类型 | 使用场景 |
|------|----------|
| Drawer 抽屉 | 创建/编辑应用、创建Webhook、插件/表情/动作详情 |
| Dialog 对话框 | 确认删除、确认审核结果 |
| 全屏模态 | 暂无复杂表单场景 |
