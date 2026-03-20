# 模块 PRD：应用管理 (App Management)

**版本：** V1.3  
**模块负责人：** agentcp  
**编制日期：** 2026-03-20  

---

## 1. 概述

应用管理模块为 MDM 中台提供企业应用分发能力，支持上传和管理企业内部应用（IPA/APK/AAB/MSI），并通过分发策略向设备推送安装、强制更新或卸载。

**业务目标：**
- 应用仓库统一管理（上传/审核/发布）
- 企业应用商店
- 应用分发策略（自动安装/强制安装/卸载黑名单）
- 应用托管配置（Managed App Config）
- 许可证管理（Apple VPP/Google 企业许可证）
- 应用使用统计（安装率/版本分布/使用时长）

---

## 2. 功能列表

| 功能 | 描述 | 优先级 | 触发方式 | 前端入口/按钮 |
|------|------|--------|----------|--------------|
| 应用仓库 | 应用上传（IPA/APK/AAB/MSI）、版本管理 | P0 | 人工 | 「上传应用」按钮 |
| 应用列表 | 分页查询应用，支持类型/状态筛选 | P0 | 自动 | 无按钮 |
| 应用分发 | 创建分发任务，推送到设备/用户/组 | P0 | 人工 | 「新建分发」按钮 |
| 强制安装 | 设备必须安装指定应用 | P1 | 自动 | 无按钮（策略触发） |
| 卸载黑名单 | 禁止用户卸载的应用列表 | P1 | 人工 | 「添加到黑名单」按钮 |
| 应用配置 | 托管参数配置（Managed App Config） | P1 | 人工 | 「应用配置」按钮 |
| VPP许可证 | Apple VPP / Google 许可证导入和管理 | P1 | 人工 | 「导入许可证」按钮 |
| 应用统计 | 安装率/版本分布/使用时长统计 | P2 | 自动 | 无按钮 |
| 应用审核 | 上传后需审核才能发布 | P2 | 人工 | 「审核通过」/「审核拒绝」按钮 |

---

## 3. 数据模型

### 3.1 应用主表 (apps)

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| id | uint | 主键 | PK, auto |
| app_code | string | 应用编码, unique | not null |
| app_name | string | 应用名称 | not null |
| app_type | string | 类型 ios/android/windows/macos | not null |
| bundle_id | string | Bundle ID / Package Name | not null |
| developer | string | 开发者名称 | nullable |
| description | string | 应用描述 | nullable |
| icon_url | string | 应用图标 URL | nullable |
| is_enterprise | bool | 是否企业应用 | default true |
| app_store_url | string | 应用商店链接 | nullable |
| status | string | 状态 pending/approved/rejected/archived | default 'pending' |
| blacklisted | bool | 是否在卸载黑名单 | default false |
| created_by | string | 创建人 | not null |
| created_at | datetime | 创建时间 | auto |
| updated_at | datetime | 更新时间 | auto |

### 3.2 应用版本表 (app_versions)

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| id | uint | 主键 | PK, auto |
| app_id | uint | 应用 ID, FK | not null |
| version | string | 版本号 | not null |
| version_code | int | 版本数字 | not null |
| file_url | string | 安装包 CDN URL | not null |
| file_size | int64 | 文件大小(字节) | default 0 |
| file_md5 | string | MD5 校验码 | nullable |
| min_os_version | string | 最低系统版本 | nullable |
| release_notes | string | 更新日志 | nullable |
| is_mandatory | bool | 是否强制更新 | default false |
| is_latest | bool | 是否最新版本 | default true |
| install_count | int | 安装次数 | default 0 |
| created_at | datetime | 创建时间 | auto |

### 3.3 应用分发任务表 (app_distributions)

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| id | uint | 主键 | PK, auto |
| app_id | uint | 应用 ID, FK | not null |
| app_version_id | uint | 分发的版本 ID, FK | not null |
| distribution_type | string | 分发类型 install/force_install/uninstall | not null |
| target_type | string | 目标类型 device/user/group/org_unit | not null |
| target_ids | jsonb | 目标 ID 列表 | not null |
| scheduled_at | datetime | 计划开始时间 | nullable |
| status | string | 状态 pending/running/completed/failed/cancelled | default 'pending' |
| total_count | int | 目标设备数 | default 0 |
| success_count | int | 成功数 | default 0 |
| failed_count | int | 失败数 | default 0 |
| created_by | string | 创建人 | not null |
| created_at | datetime | 创建时间 | auto |
| updated_at | datetime | 更新时间 | auto |

### 3.4 应用许可证表 (app_licenses)

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| id | uint | 主键 | PK, auto |
| app_id | uint | 应用 ID, FK | not null |
| license_type | string | 许可证类型 apple_vpp/google_enterprise | not null |
| license_count | int | 许可证总数 | default 0 |
| used_count | int | 已使用数 | default 0 |
| license_file | string | 许可证文件路径 | nullable |
| metadata | jsonb | 许可证元数据（来自 VPP） | nullable |
| expires_at | datetime | 过期时间 | nullable |
| created_at | datetime | 创建时间 | auto |

### 3.5 应用托管配置表 (app_configurations)

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| id | uint | 主键 | PK, auto |
| app_id | uint | 应用 ID, FK | not null |
| target_type | string | 目标类型 device/user/group | not null |
| target_id | string | 目标 ID | not null |
| config_data | jsonb | 托管配置数据 | not null |
| created_at | datetime | 创建时间 | auto |
| updated_at | datetime | 更新时间 | auto |

### 3.6 应用安装记录表 (app_installations)

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| id | uint | 主键 | PK, auto |
| app_id | uint | 应用 ID, FK | not null |
| app_version_id | uint | 版本 ID, FK | not null |
| device_id | string | 设备 ID | not null, index |
| user_id | string | 用户 ID | nullable |
| install_status | string | 安装状态 pending/downloading/installed/uninstall_failed | not null |
| installed_at | datetime | 安装时间 | nullable |
| uninstalled_at | datetime | 卸载时间 | nullable |
| created_at | datetime | 创建时间 | auto |
| updated_at | datetime | 更新时间 | auto |

---

## 4. 接口定义

### 4.1 应用管理

#### 4.1.1 应用列表
```
GET /api/v1/apps
```
**Query:** app_type, status, keyword, page, page_size

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "app_code": "APP001",
        "app_name": "企业IM",
        "app_type": "ios",
        "bundle_id": "com.company.im",
        "developer": "内部开发组",
        "status": "approved",
        "latest_version": "v2.1.0",
        "install_count": 500,
        "created_at": "2026-03-20T10:00:00Z"
      }
    ],
    "pagination": { "page": 1, "page_size": 20, "total": 50 }
  }
}
```

#### 4.1.2 应用详情
```
GET /api/v1/apps/:id
```
返回应用详情含所有版本列表

#### 4.1.3 创建应用
```
POST /api/v1/apps
```
**请求体：**
```json
{
  "app_code": "APP001",
  "app_name": "企业IM",
  "app_type": "ios",
  "bundle_id": "com.company.im",
  "developer": "内部开发组",
  "description": "企业内部即时通讯工具",
  "is_enterprise": true,
  "created_by": "admin"
}
```

#### 4.1.4 更新应用
```
PUT /api/v1/apps/:id
```

#### 4.1.5 审核应用
```
POST /api/v1/apps/:id/approve
POST /api/v1/apps/:id/reject
```
**请求体（reject）：**
```json
{
  "reason": "应用包签名无效"
}
```

---

### 4.2 应用版本管理

#### 4.2.1 上传新版本
```
POST /api/v1/apps/:id/versions
```
**请求体：**
```json
{
  "version": "v2.1.0",
  "version_code": 21,
  "file_url": "https://cdn.example.com/apps/enterprise-im-v2.1.0.ipa",
  "file_size": 52428800,
  "file_md5": "d41d8cd98f00b204e9800998ecf8427e",
  "min_os_version": "14.0",
  "release_notes": "修复消息推送延迟问题",
  "is_mandatory": true
}
```

#### 4.2.2 版本列表
```
GET /api/v1/apps/:id/versions
```

#### 4.2.3 删除版本
```
DELETE /api/v1/apps/:id/versions/:version_id
```
**限制：** 已有安装记录时不可删除

---

### 4.3 应用分发

#### 4.3.1 创建分发任务
```
POST /api/v1/apps/distributions
```
**请求体：**
```json
{
  "app_id": 1,
  "app_version_id": 5,
  "distribution_type": "install",
  "target_type": "group",
  "target_ids": ["group-devices-all"],
  "scheduled_at": "2026-03-20T14:00:00Z",
  "created_by": "admin"
}
```

**distribution_type 枚举：**
- `install` - 普通安装（用户可自行卸载）
- `force_install` - 强制安装（不可卸载）
- `uninstall` - 强制卸载

#### 4.3.2 分发任务列表
```
GET /api/v1/apps/distributions
```
**Query:** app_id, status, page, page_size

#### 4.3.3 分发任务详情
```
GET /api/v1/apps/distributions/:id
```
返回任务详情含分发进度

#### 4.3.4 取消分发任务
```
POST /api/v1/apps/distributions/:id/cancel
```

---

### 4.4 应用配置

#### 4.4.1 设置托管配置
```
POST /api/v1/apps/:id/configurations
```
**请求体：**
```json
{
  "target_type": "group",
  "target_id": "group-devices-all",
  "config_data": {
    "server_url": "https://api.company.com",
    "api_key": "xxx",
    "features": {
      "chat_enabled": true,
      "file_sharing": false
    }
  }
}
```

#### 4.4.2 查询托管配置
```
GET /api/v1/apps/:id/configurations
```

---

### 4.5 许可证管理

#### 4.5.1 许可证列表
```
GET /api/v1/apps/:id/licenses
```

#### 4.5.2 导入许可证
```
POST /api/v1/apps/:id/licenses
```
**请求体：**
```json
{
  "license_type": "apple_vpp",
  "license_file": "/path/to/vpp_license.csv",
  "expires_at": "2027-03-20T00:00:00Z"
}
```

---

### 4.6 应用统计

#### 4.6.1 安装统计
```
GET /api/v1/apps/:id/stats
```
**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total_installs": 500,
    "active_installs": 480,
    "uninstalls": 20,
    "install_rate": 96.0,
    "by_version": {
      "v2.1.0": 300,
      "v2.0.0": 150,
      "v1.9.0": 50
    },
    "by_status": {
      "installed": 480,
      "pending": 10,
      "failed": 10
    }
  }
}
```

---

## 5. 流程图

### 5.1 应用分发流程

```
管理员创建分发任务
    │
    ▼
POST /api/v1/apps/distributions
    │
    ├─→ 验证 app_id / app_version_id 存在
    ├─→ 解析 target_ids 计算目标设备列表
    ├─→ 创建 app_distributions 记录 (status=pending)
    └─→ 创建 app_installations 记录 (install_status=pending)
    │
    ▼
App Distribution Worker (后台轮询)
    │
    ▼
SELECT * FROM app_distributions WHERE status='pending' OR status='running'
    │
    ├─→ 遍历每个分发任务
    │       │
    │       ├─→ status='pending'
    │       │       ├─→ 更新 status='running'
    │       │       └─→ 触发 App Install Notification
    │       │
    │       └─→ status='running'
    │               ├─→ 处理 pending 安装请求
    │               └─→ 更新 success_count / failed_count
    │
    ▼
设备收到推送通知
    │
    ▼
设备端自动下载并安装应用
    │
    ▼
POST /api/v1/apps/installations/:id/report
    │
    ▼
更新 app_installations.install_status
    │
    ▼
更新分发任务 success_count / failed_count
    │若 success_count + failed_count = total_count
    │       └─→ status = 'completed'
```

### 5.2 强制更新流程

```
管理员上传新版本并标记 is_mandatory=true
    │
    ▼
创建强制分发任务
    │
    ▼
设备收到推送
    │
    ▼
设备端检查版本，若低于目标版本
    │
    ├─→ 强制下载安装（用户无法取消）
    │
    ▼
安装完成后上报
```

---

## 6. 模块联动

| 联动模块 | 联动方式 | 说明 |
|----------|----------|------|
| 设备管理 | 应用分发到设备 | target_type=device |
| 会员管理 | 应用分发到用户 | target_type=user，应用使用数据关联会员 |
| 组织架构 | 应用分发到组织单元 | target_type=org_unit |
| 策略管理 | 策略可包含强制安装/卸载黑名单 | 应用分发作为策略动作 |
| OTA升级 | 应用版本可关联固件版本 | App-OS 版本兼容性 |
| 数据分析 | 应用安装统计 | 安装率/版本分布 |
| 告警系统 | 应用安装失败创建告警 | 通知运营人员 |

---

## 7. 验收标准

### P0 验收标准

| 用例 | 验收条件 | 测试方法 |
|------|----------|----------|
| 应用上传 | 应用正确创建，状态为 pending | POST /apps 后验证 |
| 版本上传 | 正确创建版本，旧版本 is_latest=false | POST /apps/:id/versions 验证 |
| 应用分发 | 分发任务创建后目标设备收到安装推送 | 创建任务后检查设备通知 |
| 应用审核 | 审核通过后 status=approved | POST /apps/:id/approve 验证 |

### P1 验收标准

| 用例 | 验收条件 | 测试方法 |
|------|----------|----------|
| 强制安装 | is_mandatory=true 时用户无法取消 | 分发强制任务后设备端验证 |
| 卸载黑名单 | blacklisted=true 时不可卸载 | 尝试卸载黑名单应用 |
| 托管配置 | 配置正确推送到设备 | 检查设备收到的 config_data |
| VPP许可证 | 许可证导入后 used_count 正确扣减 | 导入+分发后验证 |

### P2 验收标准

| 用例 | 验收条件 | 测试方法 |
|------|----------|----------|
| 安装统计 | 安装率/版本分布正确 | 对比 installations 数据 |
| 应用分发进度 | 实时更新 success/failed count | 创建分发后多次查询 |

---

## 8. UI设计指引

### 页面结构
- **左侧菜单**：应用管理 → 应用列表 / 应用分发 / 许可证管理 / 应用统计
- **顶部区域**：统计卡片（总应用数 / 已发布 / 安装次数 / 成功率）
- **中间区域**：Tab 页签：应用列表 / 应用分发 / 许可证管理
- **底部区域**：分页组件

### 组件选用
| 组件 | 用途 |
|------|------|
| a-table | 应用列表、版本列表、分发任务列表 |
| a-card | 顶部统计卡片，4列布局 |
| a-tabs | Tab 切换：应用列表 / 应用分发 / 许可证管理 |
| a-drawer | 创建/编辑应用、上传新版本、托管配置 |
| a-modal | 审核确认、删除确认 |
| a-upload | 应用包文件上传 |
| a-steps | 版本发布流程（上传→审核→发布）|
| a-progress | 分发进度展示 |
| a-select | 应用类型筛选、状态筛选 |
| a-tag | 应用状态标签（待审核=黄，已发布=绿，已拒绝=红）|
| a-statistic | 安装次数/成功率展示 |

### 参考模板
```
┌──────────────────────────────────────────────────────────────┐
│  [统计卡片]  总应用:50  已发布:45  安装次数:10,000  成功率:98% │
├──────────────────────────────────────────────────────────────┤
│  [Tab: 应用列表 | 应用分发 | 许可证管理]                       │
├──────────────────────────────────────────────────────────────┤
│  【应用列表 Tab】                                             │
│  ┌──────────────────────────────────────────────────────┐   │
│  │ [iOS▼ Android▼] [待审核▼] [+上传应用]                  │   │
│  ├──────────────────────────────────────────────────────┤   │
│  │ 应用名称   │类型│版本│状态  │安装数│更新时间│操作      │   │
│  │ 企业IM    │iOS │v2.1│🟢已发布│ 500  │ 10:30  │详情分发编辑│   │
│  │ 企业邮箱  │iOS │v1.5│🟡待审核│  --  │ 09:00  │详情审核编辑│   │
│  └──────────────────────────────────────────────────────┘   │
│                                                              │
│  【应用分发 Tab】                                            │
│  ┌──────────────────────────────────────────────────────┐   │
│  │ [+新建分发任务]                                        │   │
│  │ 任务名称   │应用  │类型│目标数│成功│失败│状态  │操作  │   │
│  │ Dev组安装  │企业IM│安装│ 50   │ 45 │  2  │进行中│详情取消│   │
│  └──────────────────────────────────────────────────────┘   │
└──────────────────────────────────────────────────────────────┘
```

### 交互流程
```
应用列表页
    │
    ├── 点击「上传应用」──► a-drawer ──► 填写应用信息+上传包文件 ──► 创建（待审核）
    │
    ├── 点击「审核通过」──► 状态变为已发布
    │
    ├── 点击「审核拒绝」──► a-modal ──► 填写拒绝原因 ──► 状态变为已拒绝
    │
    ├── 点击「上传新版本」──► a-drawer ──► 填写版本信息 ──► 上传
    │
    └── 点击「分发」──► a-drawer ──► 选择版本+分发类型+目标 ──► 创建分发任务

应用分发页
    │
    ├── 查看分发任务列表和进度
    │
    ├── 点击「详情」──► a-drawer ──► 查看分发进度详情
    │
    └── 点击「取消」──► 取消分发任务
```

### 关键状态显示
- **应用状态**：a-tag，待审核=黄，已发布=绿，已拒绝=红，已归档=灰
- **分发类型**：a-tag，普通安装=蓝，强制安装=橙，卸载=红
- **分发状态**：a-tag，进行中=蓝，已完成=绿，已取消=灰，失败=红
- **安装状态**：a-tag，pending=黄，downloading=蓝，installed=绿，failed=红

---

## 修订记录

| 版本 | 日期 | 修订人 | 修订内容 |
|------|------|--------|----------|
| V1.3 | 2026-03-20 | agentcp | 全新模块，基于新增功能需求重建 |
