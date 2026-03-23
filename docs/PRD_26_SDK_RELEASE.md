# PRD：SDK发布

**版本：** V1.0
**所属Phase：** Phase 4（Sprint 25-26）
**优先级：** P2
**负责角色：** agentcp（产品）、agenthd（后端）、agentqd（前端）

---

## 一、概述

### 1.1 模块定位

SDK发布模块为第三方开发者提供MDM平台的官方SDK，支持iOS/Android/小程序/Web多端接入，帮助开发者快速集成设备管理、宠物控制、数据订阅等核心能力。

### 1.2 核心价值

- **降低接入成本**：SDK封装所有API调用，5行代码完成设备绑定
- **多端支持**：覆盖主流移动和Web平台
- **版本管理**：清晰的SDK版本和更新日志
- **示例代码**：完整的Demo项目和最佳实践

---

## 二、SDK规划

### 2.1 SDK列表

| SDK | 平台 | 语言 | 优先级 | 说明 |
|-----|------|------|--------|------|
| iOS SDK | iOS | Swift | P1 | iOS设备管理SDK |
| Android SDK | Android | Kotlin | P1 | Android设备管理SDK |
| MiniApp SDK | 微信/支付宝小程序 | JavaScript | P1 | 小程序SDK |
| Web SDK | Web | TypeScript | P2 | 前端Web SDK |

### 2.2 核心能力封装

| 能力 | 说明 | iOS | Android | MiniApp | Web |
|------|------|-----|---------|---------|-----|
| 设备绑定 | 扫码绑定设备 | ✅ | ✅ | ✅ | ✅ |
| 设备状态 | 实时获取设备在线状态 | ✅ | ✅ | ✅ | ✅ |
| 设备控制 | 下发控制指令 | ✅ | ✅ | ✅ | ✅ |
| OTA升级 | 查询/触发固件升级 | ✅ | ✅ | ✅ | ✅ |
| 宠物管理 | CRUD宠物档案 | ✅ | ✅ | ✅ | ✅ |
| 订阅管理 | 查询/购买订阅 | ✅ | ✅ | ✅ | ✅ |
| 消息订阅 | WebSocket实时消息 | ✅ | ✅ | ✅ | ✅ |
| 情绪数据 | 获取宠物情绪记录 | ✅ | ✅ | ✅ | ✅ |

---

## 三、SDK文档结构

### 3.1 文档目录

```
docs/
├── README.md                 # 快速开始
├── GETTING_STARTED.md        # 5分钟快速入门
├── AUTH.md                   # 认证授权
├── API_REFERENCE.md          # API参考
├── DEVICE.md                 # 设备管理
├── PET.md                    # 宠物管理
├── SUBSCRIPTION.md           # 订阅管理
├── CHANGELOG.md              # 更新日志
└── FAQ.md                    # 常见问题
```

### 3.2 示例代码结构

```
examples/
├── ios/                      # iOS示例项目
│   ├── MDMQuickStart/        # 快速开始Demo
│   ├── DeviceControl/        # 设备控制Demo
│   └── PetManagement/         # 宠物管理Demo
├── android/                  # Android示例项目
│   ├── MDMQuickStart/
│   └── DeviceControl/
├── miniapp/                  # 小程序示例
│   └── MDMQuickStart/
└── web/                      # Web示例
    └── MDMQuickStart/
```

---

## 四、SDK版本管理

### 4.1 版本号规范

遵循语义化版本 `主版本.次版本.修订号`：
- **主版本**：API不兼容的重大变更
- **次版本**：新增功能（向后兼容）
- **修订号**：Bug修复（向后兼容）

### 4.2 兼容矩阵

| SDK版本 | 最低API版本 | 说明 |
|---------|-------------|------|
| 1.x | API v1 | 初始版本 |
| 2.0 | API v1 | 新增WebSocket支持 |
| 2.1 | API v1 | Bug修复，性能优化 |

### 4.3 废弃策略

- 废弃功能提前2个版本通知
- 废弃后继续维护1年
- 废弃API保留指向新API的迁移指南

---

## 五、SDK质量标准

| 标准 | 要求 |
|------|------|
| 单元测试覆盖率 | >80% |
| API文档覆盖率 | 100% |
| 示例代码可运行 | 100% |
| iOS最低版本 | iOS 13.0 |
| Android最低版本 | API 21 (Android 5.0) |
| 包大小 | iOS<5MB, Android<3MB |

---

## 六、发布流程

| 阶段 | 说明 |
|------|------|
| 1. 代码冻结 | 确定发布日期前1周 |
| 2. 测试验证 | QA全量测试 |
| 3. 文档就绪 | 文档同步上线 |
| 4. SDK发布 | GitHub Release + CDN |
| 5. 开发者通知 | 邮件/社区公告 |

---

## 七、验收标准

| 验收点 | 标准 |
|--------|------|
| iOS SDK | 官方Demo可在Xcode运行，设备绑定<10秒 |
| Android SDK | 官方Demo可在Android Studio运行 |
| 小程序SDK | 官方Demo可在微信开发者工具运行 |
| Web SDK | CDN引入即可使用，无需构建 |
| API文档 | 每个API有完整参数说明和示例 |
| 示例代码 | 每个核心能力有独立示例 |
