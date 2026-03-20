# 模块 PRD：OTA升级 (OTA Updates)

**版本：** V1.4
**模块负责人：** agentcp
**编制日期：** 2026-03-20

---

## 1. 概述

OTA（Over-The-Air）固件升级是 MDM 中台的核心能力，支持 M5Stack 设备远程静默升级固件，支持灰度发布、进度追踪和异常自动回滚。

**业务目标：**
- 固件包完整生命周期管理
- 支持全量/百分比/白名单灰度策略
- 升级进度实时追踪
- 成功率低于阈值自动暂停

---

## 2. 功能列表

| 功能 | 描述 | 优先级 | 触发方式 | 前端入口/按钮 |
|------|------|--------|----------|--------------|
| 固件包上传 | 创建固件包记录，含版本/型号/MD5 | P0 | 人工 | 「上传固件」按钮 |
| 固件包列表 | 按硬件型号筛选固件包 | P0 | 自动 | 无按钮 |
| 部署任务创建 | 创建OTA部署任务 | P0 | 人工 | 「新建部署任务」按钮 |
| OTA Worker | 后台自动下发OTA指令 | P0 | 自动 | 无按钮 |
| 设备OTA检测 | 设备检查是否有新版本 | P1 | 自动 | 无按钮 |
| 升级进度追踪 | 设备上报进度，实时更新 | P1 | 自动 | 无按钮 |
| 灰度策略 | 支持百分比/白名单/全量发布 | P1 | 人工 | 「编辑」按钮（策略配置） |
| 任务暂停/取消 | 人工暂停或取消部署任务 | P1 | 人工 | 「暂停」按钮 / 「取消」按钮 |

---

## 3. 数据模型

### 3.1 固件包表 (ota_packages)

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| id | uint | 主键 | PK |
| name | string | 固件包名称 | not null |
| version | string | 版本号 | unique, not null |
| hardware_model | string | 目标硬件型号 | not null, index |
| file_size | int64 | 文件大小(字节) | default 0 |
| file_url | string | CDN下载链接 | not null |
| file_md5 | string | MD5校验码 | - |
| upload_source | string | 上传来源 local/remote | default 'local' |
| is_active | bool | 是否启用 | default true |
| is_mandatory | bool | 是否强制升级 | default false |
| allow_downgrade | bool | 是否允许降级 | default false |
| release_notes | string | 更新日志 | text |
| created_by | string | 创建人 | not null |
| created_at | datetime | 创建时间 | auto |
| updated_at | datetime | 更新时间 | auto |

**release_status:** 0=测试中, 1=灰度发布, 2=全量发布

### 3.2 部署任务表 (ota_deployments)

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| id | uint | 主键 | PK |
| name | string | 任务名称 | not null |
| package_id | uint | 固件包ID | FK → ota_packages |
| hardware_model | string | 目标硬件型号 | not null |
| strategy_type | string | 策略类型 full/percentage/whitelist | not null |
| strategy_config | jsonb | 策略配置 | JSON |
| target_device_count | int | 目标设备数 | default 0 |
| pending_count | int | 待升级数 | default 0 |
| running_count | int | 升级中数 | default 0 |
| success_count | int | 成功数 | default 0 |
| failed_count | int | 失败数 | default 0 |
| status | string | 任务状态 | pending/running/paused/completed/failed/cancelled |
| pause_reason | string | 暂停原因 | - |
| auto_paused | bool | 是否自动暂停 | default false |
| pause_on_failure_threshold | float64 | 失败率阈值 | default 20.0 |
| scheduled_at | datetime | 计划开始时间 | nullable |
| cancelled_by | string | 取消人 | - |
| cancelled_at | datetime | 取消时间 | - |
| completed_at | datetime | 完成时间 | - |
| created_by | string | 创建人 | not null |
| created_at | datetime | 创建时间 | auto |
| updated_at | datetime | 更新时间 | auto |

**strategy_config示例：**
```json
// 百分比策略: { "percentage": 30 }
// 白名单策略: { "device_ids": ["uuid1", "uuid2"] }
```

### 3.3 升级进度表 (ota_progress)

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| deployment_id | uint | 部署任务ID, index |
| device_id | string | 设备ID, index |
| package_id | uint | 固件包ID |
| from_version | string | 原版本 |
| to_version | string | 目标版本 |
| ota_status | string | pending/downloading/verifying/flashing/success/failed |
| ota_message | string | 状态消息 |
| progress_percent | int | 进度0-100 |
| retry_count | int | 重试次数 |
| started_at | datetime | 开始时间 |
| completed_at | datetime | 完成时间 |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |

---

## 4. 接口定义

### 4.1 创建固件包

```
POST /api/v1/ota/packages
```

**请求体：** name, version, hardware_model, file_size, file_url, file_md5, is_mandatory, allow_downgrade, release_notes, created_by

**响应：** 固件包完整信息含id/is_active/created_at

### 4.2 固件包列表

```
GET /api/v1/ota/packages
```

**Query:** hardware_model, page, page_size

**响应：** 固件包分页列表

### 4.3 创建部署任务

```
POST /api/v1/ota/deployments
```

**请求体：**
```json
{
  "name": "v1.3.0 灰度发布30%",
  "package_id": 1,
  "hardware_model": "CoreS3_Rover",
  "strategy_type": "percentage",
  "strategy_config": { "percentage": 30 },
  "pause_on_failure_threshold": 20.0,
  "scheduled_at": "2026-03-20T12:00:00Z",
  "created_by": "admin"
}
```

**响应：** deployment含status=pending, target_device_count, pending_count

### 4.4 部署任务列表

```
GET /api/v1/ota/deployments
```

**响应：** 部署任务分页列表，含success_rate=(success_count/(success_count+failed_count))

### 4.5 暂停/恢复/取消部署任务

```
POST /api/v1/ota/deployments/:id/pause
POST /api/v1/ota/deployments/:id/resume
POST /api/v1/ota/deployments/:id/cancel
```

**响应：** 更新后的status和pause_reason

### 4.6 设备检查OTA

```
GET /api/v1/ota/devices/:device_id/check
```

**响应(有更新)：** has_update=true, current_version, latest_version, package(含file_url/file_md5/is_mandatory)

**响应(无更新)：** has_update=false, message="当前已是最新版本"

### 4.7 设备上报OTA进度

```
POST /api/v1/ota/devices/:device_id/report
```

**请求体：** { "deployment_id": 1, "ota_status": "downloading", "progress_percent": 45, "ota_message": "正在下载..." }

**响应：** 更新后的ota_status/progress_percent

### 4.8 升级进度列表

```
GET /api/v1/ota/deployments/:id/progress
```

**响应：** summary(含pending/running/success/failed/success_rate) + 分页list

---

## 5. 流程图

### 5.1 OTA Worker工作流程

```
OTA Worker (后台goroutine, 每30s轮询)
    │
    ▼
SELECT * FROM ota_deployments WHERE status IN ('pending', 'running')
    │
    ▼
遍历每个部署任务
    │
    ├─► status='pending'
    │       ├─► 根据strategy_type计算目标设备列表
    │       │       full: 所有匹配型号设备
    │       │       percentage: 随机N%设备
    │       │       whitelist: 指定device_ids
    │       ├─► 批量创建ota_progress记录(pending)
    │       ├─► 更新pending_count/target_device_count
    │       └─► status='running'
    │
    ├─► status='running'
    │       ├─► SELECT * FROM ota_progress WHERE status='pending' LIMIT 10
    │       ├─► 对每条记录:
    │       │       MQTT下发 /device/{id}/down/cmd (OTA指令)
    │       │       更新ota_progress.ota_status='downloading'
    │       │       running_count++, pending_count--
    │       └─► 计算success_rate
    │               若success_rate < pause_on_failure_threshold
    │                       status='paused', auto_paused=true, pause_reason='成功率低于阈值'
    │
    └─► status IN ('completed', 'cancelled') ── 跳过

检查任务完成: pending_count=0 AND running_count=0 → status='completed'
```

### 5.2 设备OTA升级流程

```
设备上电/收到OTA指令
    │
    ▼
GET /api/v1/ota/devices/:device_id/check
    │
    ├─► 有更新 ──下载固件──► 刷写 ──验证──► POST /report (success)
    │               │
    │               │◄──── 失败 ───────────────┘
    │               │    POST /report (failed)
    │               ▼
    │         重试(最多3次)
    │               │
    │       ├─► 成功 ──► POST /report (success)
    │       └─► 失败 ──► POST /report (failed)
    │
    └─► 无更新 ── 流程结束
```

---

## 6. 模块联动

| 联动模块 | 联动方式 | 说明 |
|----------|----------|------|
| 设备影子 | OTA指令通过MQTT下发 | 设备收到指令后更新状态 |
| 设备管理 | 设备注册时CheckOTA | 新设备检测是否有待升级版本 |
| 数据分析 | OTAProgress表记录升级数据 | 统计升级成功率 |
| 告警系统 | 升级失败创建告警 | 通知运营人员 |
| 应用管理 | OTA版本可关联应用版本 | App依赖固件版本，固件升级后触发App更新检查 |
| 策略管理 | OTA可作为策略动作 | 合规设备可触发自动OTA升级 |
| 通知管理 | OTA结果通过通知推送 | notification_type=command_response |

---

## 7. 验收标准

### P0 验收标准

| 用例 | 验收条件 | 测试方法 |
|------|----------|----------|
| 固件包创建 | 正确创建固件包记录 | POST /api/v1/ota/packages |
| 部署任务创建 | 根据策略计算目标设备数 | 创建百分比30%任务，验证count |
| OTA Worker轮询 | 30s内自动选中pending设备下发指令 | 创建任务后观察MQTT消息 |
| 设备CheckOTA | 有更新返回固件信息，无更新返回提示 | 模拟不同版本 |
| 进度上报 | POST /report正确更新progress | 多次上报不同进度 |
| 自动暂停 | 成功率<阈值时任务自动暂停 | 制造30%失败率 |
| 白名单策略 | 只升级指定device_ids | 创建白名单任务 |

### P1 验收标准

| 用例 | 验收条件 | 测试方法 |
|------|----------|----------|
| 任务暂停 | 人工暂停后running设备继续，新设备不再下发 | 暂停任务 |
| 任务取消 | 取消后pending设备不再下发 | 取消任务 |
| 任务恢复 | 暂停任务恢复后继续下发 | 恢复并观察 |
| 强制升级 | is_mandatory=true设备必须升级 | 模拟非自愿升级场景 |

---

## 8. UI设计指引

### 8.1 页面结构
- **左侧菜单**：OTA升级 → 固件包管理 / 部署任务 / 升级进度
- **顶部区域**：统计卡片（总任务数/进行中/成功率/待升级设备数）
- **中间区域**：Tab页：固件包列表 / 部署任务列表 / 升级进度详情
- **底部区域**：分页组件

### 8.2 组件选用
| 组件 | 用途 |
|------|------|
| a-table | 固件包列表、部署任务列表、升级进度列表 |
| a-card | 顶部统计卡片，4列布局 |
| a-tabs | Tab切换：固件包/部署任务/进度详情 |
| a-drawer | 创建固件包、创建部署任务、查看任务详情 |
| a-upload | 固件包文件上传（支持.bin/.hex）|
| a-progress | 升级进度环形图/条形图 |
| a-select | 硬件型号筛选、任务状态筛选 |
| a-switch | 强制升级开关、启用/禁用开关 |
| a-tag | 任务状态标签（进行中=蓝，暂停=橙，完成=绿，失败=红）|

### 8.3 参考模板
```
┌──────────────────────────────────────────────────────────────┐
│  [统计卡片]  总任务:50  进行中:3  成功率:92%  待升级:150  │
├──────────────────────────────────────────────────────────────┤
│  [Tab: 固件包管理 | 部署任务 | 升级进度]                      │
├──────────────────────────────────────────────────────────────┤
│  【部署任务Tab】                                             │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  [筛选: 状态▼  型号▼]  [+新建部署任务]              │   │
│  ├──────────────────────────────────────────────────────┤   │
│  │ 任务名            │ 版本 │ 状态    │ 成功率 │ 操作  │   │
│  │ v1.3.0 灰度30%   │v1.3.0│🟡进行中│  94%   │详情暂停│   │
│  │ v1.2.0 全量发布   │v1.2.0│🟢完成  │  91%   │详情    │   │
│  └──────────────────────────────────────────────────────┘   │
│                                                              │
│  【升级进度Tab - 任务详情】                                  │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  目标设备: 150  成功: 80  失败: 5  进行中: 30       │   │
│  │  [████████████████░░░░░░░░]  56%                    │   │
│  └──────────────────────────────────────────────────────┘   │
└──────────────────────────────────────────────────────────────┘
```

### 8.4 交互流程
```
OTA固件包管理页
    ├── 点击「上传固件」──► a-upload弹窗 ──► 选择.bin文件 ──► 上传完成
    └── 点击「创建任务」──► a-drawer ──► 选择固件/策略/目标设备 ──► 创建成功

OTA部署任务列表页
    ├── 点击「暂停」──► a-modal二次确认 ──► 任务暂停
    ├── 点击「取消」──► a-modal二次确认 ──► 任务取消
    ├── 点击「详情」──► a-drawer ──► 查看进度详情
    └── 点击「恢复」──► 任务继续下发指令

升级进度Tab
    └── 查看进度列表，支持按状态筛选
```

### 8.5 关键状态显示
- **任务状态**：a-tag，蓝=进行中，橙=暂停，绿=完成，红=失败
- **升级进度**：a-progress环形图+条形图双视图
- **成功率**：低于阈值(默认20%)红色警示
- **设备升级状态**：a-tag，pending=灰，downloading=蓝，success=绿，failed=红

---

## 附录 B. 修订记录

| 版本 | 日期 | 修订人 | 修订内容 |
|------|------|--------|----------|
| V1.0 | 2026-03-20 | agentcp | 初稿，基于代码调研 |
| V1.2 | 2026-03-20 | agentcp | 修订功能列表，补充触发方式和前端入口按钮列 |
| V1.4 | 2026-03-20 | agentcp | 重建文档结构，统一使用8章节格式，合并重复的八、九章节