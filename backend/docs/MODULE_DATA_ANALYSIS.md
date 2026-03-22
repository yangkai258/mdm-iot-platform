# 模块 PRD：数据分析 (Data Analysis)

**版本：** V1.4
**模块负责人：** agentcp
**编制日期：** 2026-03-20

---

## 1. 概述

数据分析模块为 MDM 中台提供运营数据统计和报表能力，帮助运营人员了解设备运行状况、会员活跃度和 OTA 升级效果。

**业务目标：**
- 设备状态大盘统计
- 会员活跃度和消费分析
- OTA升级效果报表
- 告警趋势分析
- 自定义报表

---

## 2. 功能列表

| 功能 | 描述 | 优先级 | 触发方式 | 前端入口/按钮 |
|------|------|--------|----------|--------------|
| Dashboard大盘 | 设备/告警核心指标汇总 | P0 | 自动 | 无按钮 |
| 设备统计 | 设备在线/离线/生命周期统计 | P1 | 自动 | 无按钮 |
| OTA统计 | 升级成功率/进度统计 | P1 | 自动 | 无按钮 |
| 会员统计 | 会员增长/活跃/消费分析 | P2 | 自动 | 无按钮 |
| 告警统计 | 告警趋势/分布统计 | P2 | 自动 | 无按钮 |
| 自定义报表 | 支持自定义时间范围和维度 | P3 | 人工 | 「生成报表」按钮 |

---

## 3. 数据模型

数据分析模块不新增数据表，主要依赖现有模块的统计数据。

### 3.1 Dashboard统计数据结构

```json
{
  "devices": { "total": 1000, "online": 850, "offline": 150, "activating": 10, "maintaining": 5 },
  "alerts": { "total_today": 25, "pending": 8, "critical": 2 },
  "members": { "total": 5000, "active_today": 320, "active_week": 1200 },
  "ota": { "success_rate": 92.5, "running_tasks": 3 }
}
```

---

## 4. 接口定义

### 4.1 Dashboard大盘

```
GET /api/v1/dashboard/stats
```

**响应：** devices( total/online/offline/online_rate ), alerts( total_today/pending/critical/resolved_today ), members( total/active_today/new_today ), ota( total_tasks/running_tasks/avg_success_rate ), timestamp

### 4.2 设备统计

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/stats/devices/overview | 设备概览（Query: start_date, end_date） |
| GET | /api/v1/stats/devices/trend | 设备趋势（Query: start_date, end_date, granularity=day/week/month） |

**设备概览响应：** summary(total/online/offline/online_rate), by_lifecycle, by_hardware_model, by_firmware_version

**设备趋势响应：** granularity, list(date/total/online_avg/new_registered)

### 4.3 OTA统计

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/stats/ota/overview | OTA概览 |
| GET | /api/v1/stats/ota/tasks | OTA任务详情（Query: status, start_date, end_date, page, page_size） |
| GET | /api/v1/stats/ota/version-distribution | OTA版本分布 |

**OTA概览响应：** total_tasks/completed_tasks/running_tasks/paused_tasks/failed_tasks/total_upgraded/success_count/failed_count/avg_success_rate

### 4.4 会员统计

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/stats/members/overview | 会员概览 |
| GET | /api/v1/stats/members/level-distribution | 会员等级分布 |
| GET | /api/v1/stats/members/consumption | 会员消费统计（Query: start_date, end_date, granularity） |

**会员概览响应：** total/active_today/active_week/active_month/new_today/new_week/total_points/total_balance

### 4.5 告警统计

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/stats/alerts/overview | 告警概览 |
| GET | /api/v1/stats/alerts/trend | 告警趋势（Query: start_date, end_date, granularity） |
| GET | /api/v1/stats/alerts/distribution | 告警分布 |

**告警概览响应：** total_today/total_week/total_month/pending/critical/resolved_today/resolution_rate

---

## 5. 流程图

### 5.1 Dashboard数据聚合流程

```
前端请求GET /dashboard/stats
    │
    ▼
并行查询多个数据源:
    │
    ├─► 设备统计: SELECT COUNT(*) FROM devices + Redis SCAN shadow:* COUNT is_online
    ├─► 告警统计: SELECT COUNT(*) FROM device_alerts WHERE created_at>today
    ├─► 会员统计: SELECT COUNT(*) FROM members
    └─► OTA统计: SELECT AVG(success_rate) FROM ota_deployments
    │
    ▼
聚合结果返回前端
```

---

## 6. 模块联动

| 联动模块 | 联动方式 | 说明 |
|----------|----------|------|
| 设备管理 | 统计设备数量/在线率 | 查询devices表+Redis |
| 设备影子 | 实时在线设备数 | Redis SCAN计数 |
| OTA升级 | 升级成功率统计 | 查询ota_progress表 |
| 会员管理 | 会员活跃度/消费统计 | 查询members/orders表 |
| 告警系统 | 告警趋势/分布统计 | 查询device_alerts表 |
| 应用管理 | App安装统计/版本分布 | 查询app_installations表 |
| 策略管理 | 合规状态/违规趋势统计 | 查询device_compliance_status表 |
| 内容管理 | 内容下载/阅读统计 | 查询contents表download_count/view_count |

---

## 7. 验收标准

### P0 验收标准

| 用例 | 验收条件 | 测试方法 |
|------|----------|----------|
| Dashboard大盘 | 设备/告警/会员/OTA核心指标正确 | 对比各模块实际数据 |
| 实时性 | 数据延迟<5秒 | 请求时间戳验证 |

### P1 验收标准

| 用例 | 验收条件 | 测试方法 |
|------|----------|----------|
| 设备趋势 | 支持day/week/month维度 | 不同granularity请求 |
| 设备分布 | 按型号/版本/状态正确分组 | 对比实际设备数据 |
| OTA成功率 | 正确计算success/(success+failed) | 对比ota_progress数据 |
| 会员活跃度 | 正确区分日/周/月活跃 | 验证时间范围计算 |

### P2 验收标准

| 用例 | 验收条件 | 测试方法 |
|------|----------|----------|
| 告警趋势 | 正确按类型分组统计 | 对比device_alerts数据 |
| 自定义报表 | 支持自定义时间范围 | 不同时间范围请求 |

---

## 8. UI设计指引

### 8.1 页面结构
- **左侧菜单**：数据分析 → Dashboard / 设备统计 / OTA统计 / 会员分析 / 告警统计
- **顶部区域**：Dashboard全局统计概览（大数字卡片+趋势图）
- **中间区域**：Tab页签：Dashboard/设备/OTA/会员/告警
- **底部区域**：分页组件

### 8.2 组件选用
| 组件 | 用途 |
|------|------|
| a-card | 统计数字卡片（多个并排）|
| a-row / a-col | 统计卡片栅格布局（4列）|
| a-tabs | Tab切换不同统计模块 |
| a-table | 详细数据表格（趋势列表、分布列表）|
| a-line | 趋势折线图（设备在线趋势、告警趋势）|
| a-pie | 分布饼图（设备型号分布、告警类型分布）|
| a-bar | 柱状图（OTA版本分布、会员等级分布）|
| a-date-picker | 时间范围选择器 |
| a-range-picker | 趋势图时间范围选择 |
| a-select | 粒度选择（天/周/月）、维度选择 |
| a-statistic | Dashboard数字展示（支持动画计数）|

### 8.3 参考模板
```
┌──────────────────────────────────────────────────────────────┐
│  [Tab: Dashboard | 设备统计 | OTA统计 | 会员分析 | 告警统计]       │
├──────────────────────────────────────────────────────────────┤
│  【Dashboard Tab】                                           │
│  ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐              │
│  │ 总设备 │ │ 在线   │ │ 离线   │ │今日告警│              │
│  │ 1,000  │ │  850   │ │  150   │ │   25   │              │
│  └────────┘ └────────┘ └────────┘ └────────┘              │
│                                                              │
│  ┌─────────────────────────┐ ┌─────────────────────────┐    │
│  │  设备在线趋势 (折线图)   │ │  设备型号分布 (饼图)   │    │
│  │  📈                    │ │   CoreS3  ████████ 60% │    │
│  │  900 ─────────────     │ │   StickC  █████   30% │    │
│  │  800 ─────────    ─────│ │   AtomS3  ███    10% │    │
│  │  [近7天 / 近30天 / 自定义]│ │                      │    │
│  └─────────────────────────┘ └─────────────────────────┘    │
│                                                              │
│  ┌─────────────────────────┐ ┌─────────────────────────┐    │
│  │  会员活跃趋势            │ │  OTA升级成功率趋势      │    │
│  │  📈                    │ │   📈                    │    │
│  └─────────────────────────┘ └─────────────────────────┘    │
└──────────────────────────────────────────────────────────────┘
```

### 8.4 交互流程
```
Dashboard首页
    ├── 页面加载 ──► 并行请求多个统计接口 ──► 聚合展示
    ├── 点击数字卡片 ──► 跳转对应统计Tab页
    ├── 时间范围切换 ──► 重新请求数据 ──► 图表刷新
    └── 支持全屏放大查看图表

设备统计页
    ├── 选择时间范围（近7天/30天/自定义）
    ├── 查看趋势折线图（设备总数/在线平均/新增注册）
    └── 查看分布饼图/柱状图（按型号/版本/生命周期）

OTA统计页
    ├── 查看升级成功率趋势
    └── 查看版本分布柱状图

会员分析页
    ├── 查看会员增长趋势
    └── 查看会员等级/消费分布

告警统计页
    ├── 查看告警趋势折线图
    └── 查看告警类型/严重程度分布饼图
```

### 8.5 关键状态显示
- **统计数字**：a-statistic，支持数字动画（从0递增到目标值）
- **在线率**：a-progress显示百分比
- **图表加载态**：图表区域skeleton骨架屏
- **空数据态**：a-empty组件，提示「暂无数据」

---

## 附录 B. 修订记录

| 版本 | 日期 | 修订人 | 修订内容 |
|------|------|--------|----------|
| V1.0 | 2026-03-20 | agentcp | 初稿，基于代码调研 |
| V1.2 | 2026-03-20 | agentcp | 修订功能列表，补充触发方式和前端入口按钮列 |
| V1.4 | 2026-03-20 | agentcp | 重建文档结构，统一使用8章节格式，合并重复的八、九章节