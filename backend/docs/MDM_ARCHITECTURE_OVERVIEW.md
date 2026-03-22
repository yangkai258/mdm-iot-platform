# MDM 架构总览

**版本：** V1.3  
**编制角色：** 产品经理 (agentcp)  
**编制日期：** 2026-03-20  

---

## 一、系统架构图

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              Frontend (Vue 3)                               │
│                    Vue 3 + TypeScript + Arco Design Vue                     │
│                         Nginx 静态部署 / CDN 加速                            │
└──────────────────────────────────────┬──────────────────────────────────────┘
                                       │ HTTPS + REST API
                                       │ JWT Bearer Token
┌──────────────────────────────────────▼──────────────────────────────────────┐
│                              Backend (Go + Gin)                              │
│  ┌──────────────┬──────────────┬──────────────┬──────────────┐              │
│  │  Controllers │  Middleware  │   Models     │    MQTT      │              │
│  │              │  JWTAuth   │  GORM/Postgres│   Handler   │              │
│  │              │  OpLog     │              │              │              │
│  └──────────────┴──────────────┴──────┬───────┴──────────────┘              │
│                                         │                                     │
│  ┌──────────────────────────────────────▼──────────────────────────────┐      │
│  │                      Service Layer                                   │      │
│  │  device_service / member_service / ota_service / alert_service     │      │
│  │  policy_service / app_service / content_service / notify_service    │      │
│  └──────────────────────────────────────┬──────────────────────────────┘      │
└──────────────────────────────────────────┼───────────────────────────────────┘
                                           │
              ┌─────────────────────────────┼─────────────────────────────┐
              │                             │                             │
    ┌─────────▼─────────┐       ┌─────────▼─────────┐       ┌─────────▼─────────┐
    │    PostgreSQL     │       │       Redis        │       │    EMQX 5.x      │
    │  (主从流复制)      │       │   (设备影子缓存)    │       │   (MQTT Broker)   │
    │                    │       │                    │       │                   │
    │  • devices         │       │  shadow:{device_id}│       │  /device/{id}/up │
    │  • members         │       │  TTL: 90s          │       │  /device/{id}/down│
    │  • ota_packages   │       │                    │       │                   │
    │  • alerts          │       │  device_ids (Set) │       │                   │
    │  • policies        │       │                    │       │                   │
    │  • apps            │       │                    │       │                   │
    │  • contents        │       │                    │       │                   │
    │  • org_*          │       │                    │       │                   │
    └───────────────────┘       └─────────┬───────────┘       └─────────┬─────────┘
                                          │                             │
                                          │     ┌───────────────────────┘
                                          │     │ MQTT over TCP/TLS
                                    ┌──────▼──────▼──────┐
                                    │   M5Stack Devices  │
                                    │  MimiClaw Firmware │
                                    │                    │
                                    │  • 心跳上报 30s    │
                                    │  • OTA 升级回调    │
                                    │  • 指令响应        │
                                    │  • 越狱检测上报    │
                                    │  • 位置信息上报    │
                                    └────────────────────┘
```

---

## 二、技术栈明细

| 层级 | 技术选型 | 版本 | 说明 |
|------|----------|------|------|
| 设备端 | M5Stack 硬件 + MimiClaw 固件 | - | 物联网终端设备 |
| 通信层 | MQTT (paho.mqtt.golang) | 5.x | 设备与云端双向通信 |
| MQTT Broker | EMQX | 5.3 | 企业级 MQTT Broker，支持集群 |
| 后端框架 | Go + Gin | 1.26/1.12 | 高性能 HTTP 框架 |
| ORM | GORM | - | Go ORM 库 |
| 数据库 | PostgreSQL | 15 | 主从流复制 |
| 缓存 | Redis | 7 | 设备影子状态缓存 |
| 前端 | Vue 3 + TypeScript + Vite | 3.x | 渐进式前端框架 |
| UI组件 | Arco Design Vue | - | 企业级 Vue UI 组件库 |
| 部署 | Docker + docker-compose + Nginx | - | 容器化部署 |

---

## 三、模块关系图（V1.3 完整版）

```
┌─────────────────────────────────────────────────────────────────────┐
│                        MDM 中台 (Go Backend)                        │
│                                                                     │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────────────────┐│
│  │  设备管理    │◄──►│  设备影子    │◄──►│  OTA 升级               ││
│  │ Device Mgmt │    │Device Shadow│    │  OTA Updates            ││
│  │             │    │             │    │                        ││
│  │ • 注册      │    │ • 心跳 30s  │    │ • 固件包管理           ││
│  │ • 绑定/解绑 │    │ • Redis TTL │    │ • 灰度发布              ││
│  │ • 生命周期  │    │ • 期望状态   │    │ • 进度追踪              ││
│  │ • 远程控制  │    │ • 越狱检测   │    │                        ││
│  │             │    │ • 地理围栏   │    │                        ││
│  └──────┬──────┘    └──────┬──────┘    └────────────┬────────────┘│
│         │                  │                        │             │
│         │          ┌───────▼───────┐                │             │
│         │          │  告警系统      │◄───────────────┘             │
│         │          │  Alert System │                              │
│         │          │               │                              │
│         │          │ • 规则配置     │                              │
│         │          │ • 越狱/围栏    │                              │
│         │          │ • 通知渠道     │                              │
│         └─────────►│ • 合规联动     │                              │
│                    └───────────────┘                              │
│                                                                     │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────────────────┐│
│  │  会员管理    │◄──►│  数据分析    │◄──►│  系统管理               ││
│  │ Member Mgmt │    │ Data Report │    │ System Mgmt             ││
│  │             │    │             │    │                         ││
│  │ • 会员信息  │    │ • 设备统计   │    │ • 用户/角色/权限        ││
│  │ • 会员卡    │    │ • OTA 报表   │    │ • 设备事件日志          ││
│  │ • 积分/券   │    │ • App 统计   │    │ • APNs 配置             ││
│  │ • 促销/订单 │    │ • 合规报表   │    │                         ││
│  └──────┬──────┘    └─────────────┘    └────────────┬────────────┘│
│         │                                             │             │
│         └─────────────────┬───────────────────────────┘             │
│                           │                                         │
│                    ┌──────▼──────┐                                  │
│                    │  组织架构   │                                  │
│                    │  Org Mgmt   │                                  │
│                    │            │                                  │
│                    │ • 公司/部门 │                                  │
│                    │ • 岗位/员工 │                                  │
│                    │ • 身份同步  │                                  │
│                    └──────┬──────┘                                  │
│                           │                                         │
│  ┌─────────────────────────────────────────────────────────────┐    │
│  │                    新增模块 (V1.3)                          │    │
│  │                                                             │    │
│  │  ┌──────────────┐   ┌──────────────┐   ┌──────────────┐   │    │
│  │  │  策略管理     │   │  应用管理     │   │  内容管理    │   │    │
│  │  │Policy Mgmt   │   │App Management│   │Content Mgmt  │   │    │
│  │  │              │   │              │   │              │   │    │
│  │  │• 配置文件库  │   │• 应用仓库    │   │• 文件库      │   │    │
│  │  │• 策略定义绑定│   │• 企业商店    │   │• 分发任务    │   │    │
│  │  │• 合规规则   │   │• 分发策略    │   │• 访问控制    │   │    │
│  │  │• 不合规动作 │   │• VPP许可证   │   │• 安全容器    │   │    │
│  │  └──────┬───────┘   └──────┬───────┘   └──────┬───────┘   │    │
│  │         │                  │                  │            │    │
│  │         └──────────────────┼──────────────────┘            │    │
│  │                            │                                │    │
│  │                    ┌───────▼───────┐                        │    │
│  │                    │  通知管理     │                        │    │
│  │                    │Notification  │                        │    │
│  │                    │              │                        │    │
│  │                    │• 推送通知    │                        │    │
│  │                    │• 企业公告    │                        │    │
│  │                    │• 通知模板    │                        │    │
│  │                    │• 命令反馈    │                        │    │
│  │                    └──────────────┘                        │    │
│  └─────────────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────────────┘
                              │
                              │ HTTP REST
                              ▼
                    ┌─────────────────────┐
                    │  Vue 3 管理控制台   │
                    │  (Arco Design Vue)  │
                    └─────────────────────┘
```

---

## 四、模块联动关系图（V1.3 新增模块联动）

```
                        ┌─────────────┐
                        │  设备注册    │
                        └──────┬──────┘
                               │
              ┌────────────────┼────────────────┐
              │                │                │
              ▼                ▼                ▼
       ┌──────────┐     ┌──────────┐     ┌──────────┐
       │策略检查  │     │应用分发  │     │内容分发  │
       │Policy    │     │App Dist  │     │Content   │
       └────┬─────┘     └────┬─────┘     └────┬─────┘
            │                │                │
            ▼                ▼                ▼
       ┌──────────┐     ┌──────────┐     ┌──────────┐
       │合规策略  │     │企业商店  │     │文件库    │
       │Compliance│     │App Store │     │File Repo │
       └────┬─────┘     └──────────┘     └──────────┘
            │
            ▼
       ┌──────────┐     ┌──────────┐     ┌──────────┐
       │告警系统  │◄────│安全检测  │     │会员管理  │
       │Alert     │     │Security  │     │Member    │
       │(越狱/围栏│     │(越狱/位置│     │◄──应用/内容│
       │合规联动) │     │检测)    │     │分发入口  │
       └──────────┘     └──────────┘     └──────────┘
            │
            ▼
       ┌──────────┐
       │通知管理  │
       │Notifica- │
       │tion      │
       └──────────┘
            │
            ▼
       ┌──────────┐     ┌──────────┐
       │数据分析  │◄────│OTA升级   │
       │Data      │     │OTA       │
       │Analytics │     │(关联App版本)│
       │(含App统计 │     └──────────┘
       │合规报表) │
       └──────────┘
```

---

## 五、数据流架构

### 5.1 设备数据流

```
M5Stack 设备上电
    │
    ├─→ POST /api/v1/devices/register (首次注册，无需Token)
    │
    ▼
DeviceController.Register()
    │
    ├─→ 写 PostgreSQL (devices 表, lifecycle_status=1)
    ├─→ Redis SADD device_ids Set
    ├─→ 创建 Redis shadow:{device_id} (TTL 90s)
    │
    ├─→ 触发合规检查 CheckCompliance()
    │       │
    │       ├─→ 查询设备绑定的策略
    │       ├─→ 执行合规规则检查（加密/越狱/密码/地理围栏）
    │       ├─→ 合规 ──► 正常服役
    │       └─→ 不合规 ──► 触发 remediation_action (隔离/擦除/告警)
    │
    ▼
用户扫码 → POST /api/v1/devices/bind/:sn_code
    │
    ▼
DeviceController.Bind()
    │
    ├─→ 写 PostgreSQL (bind_user_id, lifecycle_status=2)
    └─→ 触发 OpenClaw 初始化宠物记忆库
            │
            ▼
设备 MQTT 连接 → /device/{id}/up/status (心跳 30s)
            │
            ▼
MQTTHandler.StatusMessageHandler()
            │
            ├─→ 解析 JSON payload
            ├─→ 写 Redis shadow:{device_id} (TTL 90s)
            │     device_id, is_online, battery_level, current_mode, last_heartbeat
            │     jailbreak_detected, location (经纬度)
            ├─→ 调用 CheckAlerts() 检查告警规则（电量低/离线/越狱/围栏）
            ├─→ 调用 CheckCompliance() 检查合规策略
            ├─→ 若 desired_config 存在，通过 /down/desired 下发
            └─→ 若有 App/内容分发指令，通过 /down/desired 下发
```

### 5.2 OTA 升级数据流

```
运营人员在控制台创建部署任务
    │
    ├─→ POST /api/v1/ota/deployments
    │
    ▼
OTAController.CreateDeployment()
    │
    ├─→ 写 PostgreSQL (ota_deployments 表, status=pending)
    ├─→ 根据策略计算目标设备列表
    └─→ 更新 ota_progress 表 (pending 状态)
            │
            ▼
OTA Worker (后台 goroutine, 30s 轮询)
    │
    ├─→ SELECT * FROM ota_deployments WHERE status='pending' OR status='running'
    ├─→ 对 pending 设备，下发 MQTT /device/{id}/down/cmd (OTA指令)
    ├─→ 更新 OTAProgress.ota_status = 'downloading'
    │
    ▼
设备端固件接收指令 → 下载固件 → 刷写 → 上报进度
    │
    ├─→ MQTT /device/{id}/up/ota_progress
    │
    ▼
MQTTHandler.OTAProgressHandler()
    │
    ├─→ 更新 OTAProgress 表 (progress_percent, ota_status)
    ├─→ 计算部署任务成功率
    └─→ 若成功率 < pause_on_failure_threshold，暂停部署
```

### 5.3 策略合规检查数据流

```
设备注册/心跳触发合规检查
    │
    ▼
CheckCompliance(db, device_id, shadow_data)
    │
    ├─→ 查询设备绑定的所有策略
    │       SELECT p.* FROM policies p
    │       JOIN policy_bindings pb ON p.id = pb.policy_id
    │       WHERE pb.binding_type='device' AND pb.binding_id=:device_id
    │       AND p.is_active = true
    │       ORDER BY p.priority DESC
    │
    ├─→ 取最高优先级策略的 compliance_rules
    │
    ▼
遍历合规规则执行检查:
    │
    ├─→ require_encryption ──► 检查设备存储加密状态
    │
    ├─→ jailbreak_detection ──► 检查 shadow_data.jailbreak_detected
    │
    ├─→ min_password_length ──► 检查密码策略
    │
    ├─→ geofencing ──► 检查 shadow_data.location 是否在允许区域内
    │
    ▼
所有规则检查完成
    │
    ├─→ 无违规 ──► is_compliant=true ──► 更新 device_compliance_status
    │
    └─→ 有违规 ──► is_compliant=false
            ├─► 记录 violation 到 device_compliance_status.last_check_result
            ├─► 创建 device_alerts (compliance_violation)
            └─► 执行 remediation_action:
                    ├─► notify ──► 告警通知
                    ├─► quarantine ──► 限制设备功能
                    └─► wipe ──► 远程擦除
```

---

## 六、模块联动矩阵（V1.3 完整版）

| 触发模块 | 被触发模块 | 联动方式 | 说明 |
|----------|------------|----------|------|
| 设备管理 | 设备影子 | MQTT 消息 | 设备注册/绑定时初始化影子 |
| 设备管理 | 设备影子 | 解绑时清除 | 解绑时清除 Redis shadow |
| 设备影子 | 告警系统 | 函数调用 | CheckAlerts() 在心跳处理中调用 |
| 设备影子 | 策略管理 | desired_config | 策略通过 desired_config 下发到设备 |
| 设备影子 | OTA 升级 | MQTT 消息 | OTA 指令下发改变设备状态 |
| 设备影子 | 应用管理 | desired_config | App 分发指令下发 |
| 设备影子 | 内容管理 | desired_config | 内容分发指令下发 |
| 设备影子 | 通知管理 | MQTT 消息 | 通知通过 /device/{id}/down/notification 下发 |
| OTA 升级 | 设备影子 | MQTT 消息 | OTA 指令改变设备固件版本 |
| OTA 升级 | 数据分析 | 数据库写入 | OTAProgress 表记录升级数据 |
| OTA 升级 | 告警系统 | 函数调用 | 升级失败时创建告警 |
| OTA 升级 | 应用管理 | 版本关联 | OTA 可关联应用版本依赖 |
| 会员管理 | 数据分析 | 数据库查询 | 会员消费/积分数据统计 |
| 会员管理 | 应用管理 | target_type=user | 应用分发到用户 |
| 会员管理 | 内容管理 | target_type=user | 内容分发到用户 |
| 会员管理 | 通知管理 | target_type=user | 通知推送给用户 |
| 会员管理 | 设备管理 | bind_user_id | 设备绑定会员 |
| 策略管理 | 设备影子 | desired_config | 策略配置下发 |
| 策略管理 | 告警系统 | 函数调用 | 不合规时触发告警 |
| 策略管理 | 会员管理 | binding_type=user | 用户级策略绑定 |
| 策略管理 | 组织架构 | binding_type=org_unit | 组织单元策略继承 |
| 应用管理 | 设备管理 | target_type=device | App 分发到设备 |
| 应用管理 | 数据分析 | 数据库查询 | App 安装统计 |
| 应用管理 | 告警系统 | 函数调用 | App 安装失败时创建告警 |
| 应用管理 | 通知管理 | 通知推送 | App 安装结果通知 |
| 内容管理 | 设备管理 | target_type=device | 内容分发到设备 |
| 内容管理 | 会员管理 | target_type=user | 内容分发到用户 |
| 内容管理 | 数据分析 | 数据库查询 | 内容下载/阅读统计 |
| 通知管理 | 设备管理 | MQTT 消息 | 通知推送到设备 |
| 通知管理 | 会员管理 | HTTP API | 通知推送给用户 |
| 告警系统 | 通知管理 | 函数调用 | 告警通知复用通知模块 |
| 告警系统 | 数据分析 | 数据库查询 | 告警趋势统计 |
| 系统管理 | 设备管理 | 用户权限 | 运营人员管理设备权限 |
| 系统管理 | 设备影子 | 设备事件日志 | 设备事件上报记录 |
| 系统管理 | 通知管理 | APNs 配置 | iOS 推送配置 |
| 组织架构 | 策略管理 | 策略继承 | 组织单元绑定策略 |
| 组织架构 | 应用管理 | target_type=org_unit | 组织单元应用分发 |
| 组织架构 | 内容管理 | target_type=org_unit | 组织单元内容分发 |
| 组织架构 | 系统管理 | sys_user_ext | 用户关联员工/部门 |
| 数据分析 | 设备管理 | 数据库查询 | 设备统计 |
| 数据分析 | OTA 升级 | 数据库查询 | OTA 统计 |
| 数据分析 | 应用管理 | 数据库查询 | App 统计 |
| 数据分析 | 策略管理 | 数据库查询 | 合规报表 |
| 数据分析 | 内容管理 | 数据库查询 | 内容统计 |

---

## 七、接口规范

### 7.1 统一响应格式

```json
// 成功
{
  "code": 0,
  "message": "success",
  "data": { ... }
}

// 失败
{
  "code": 4001,
  "message": "设备不存在",
  "error_code": "ERR_DEVICE_002"
}
```

### 7.2 认证方式

- **HTTP API：** JWT Bearer Token
- **MQTT：** 用户名/密码 + Topic ACL
- **JWT 有效期：** 24小时（无 refresh_token）
- **Token 格式：** `Authorization: Bearer <token>`

### 7.3 分页格式

```json
{
  "data": {
    "list": [...],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 100,
      "total_pages": 5
    }
  }
}
```

---

## 八、安全设计

| 安全项 | 当前状态 | 规划 |
|--------|----------|------|
| JWT 密钥 | ⚠️ 硬编码在源码 | 环境变量 JWT_SECRET |
| CORS | ⚠️ 全通配 * | 配置化 CORS_ALLOWED_ORIGINS |
| 数据库凭证 | ⚠️ docker-compose 硬编码 | K8s Secret / .env |
| 登录限流 | 🔜 未实现 | 每 IP 5次/分钟 |
| PostgreSQL SSL | ⚠️ sslmode=disable | sslmode=require |
| MQTT TLS | 🔜 未实现 | tls:// 协议 |
| APNs TLS | 🔜 未实现 | APNs 证书双向认证 |
| 配置文件加密 | 🔜 未实现 | 敏感配置（VPN密码等）加密存储 |

---

## 九、数据库表清单

### 9.1 设备模块
- `devices` - 设备主表
- `device_shadows` - 设备影子
- `pet_profiles` - 宠物配置
- `command_histories` - 指令历史
- `device_groups` - 设备分组
- `device_group_relations` - 分组关联
- `device_tags` - 设备标签
- `device_tag_relations` - 标签关联
- `device_operation_logs` - 设备操作日志
- `device_event_logs` - 设备事件日志
- `device_compliance_status` - 设备合规状态

### 9.2 OTA模块
- `ota_packages` - OTA 固件包
- `ota_deployments` - OTA 部署任务
- `ota_progress` - 设备升级进度

### 9.3 会员模块
- `members` - 会员信息
- `member_cards` - 会员卡
- `member_card_groups` - 会员卡分组
- `member_levels` - 会员等级
- `member_upgrade_rules` - 升级规则
- `member_tags` / `member_tag_records` - 标签
- `coupons` / `coupon_grants` - 优惠券
- `promotions` - 促销活动
- `stores` - 店铺
- `points_rules` - 积分规则
- `member_points_records` - 积分流水
- `member_orders` - 会员订单
- `member_upgrade_records` - 等级调整流水
- `temp_members` - 临时会员

### 9.4 告警模块
- `device_alert_rules` - 告警规则
- `device_alerts` - 告警记录

### 9.5 策略管理模块（新增）
- `policy_configs` - 配置文件
- `policies` - 策略主表
- `policy_bindings` - 策略绑定
- `compliance_rules` - 合规规则
- `policy_versions` - 策略版本历史

### 9.6 应用管理模块（新增）
- `apps` - 应用主表
- `app_versions` - 应用版本
- `app_distributions` - 应用分发任务
- `app_licenses` - 应用许可证
- `app_configurations` - 应用托管配置
- `app_installations` - 应用安装记录

### 9.7 内容管理模块（新增）
- `contents` - 内容主表
- `content_categories` - 内容分类
- `content_tags` - 内容标签
- `content_tag_relations` - 标签关系
- `content_distributions` - 内容分发任务
- `content_permissions` - 内容访问权限

### 9.8 通知管理模块（新增）
- `notifications` - 通知主表
- `notification_templates` - 通知模板
- `announcements` - 公告
- `device_notifications` - 设备通知记录

### 9.9 系统模块
- `sys_users` - 系统用户
- `sys_roles` - 角色
- `sys_menus` - 菜单权限
- `sys_dictionaries` - 字典表
- `sys_operation_logs` - 操作日志
- `sys_login_logs` - 登录日志
- `sys_user_exts` - 用户扩展
- `apns_config` - APNs 配置

### 9.10 组织模块
- `companies` - 公司
- `departments` - 部门
- `positions` - 岗位
- `employees` - 员工
- `standard_positions` - 基准岗位

---

## 十、修订记录

| 版本 | 日期 | 修订人 | 修订内容 |
|------|------|--------|----------|
| V1.0 | 2026-03-20 | agentcp | 初稿，基于 ARCHITECTURE_ANALYSIS.md 和代码调研 |
| V1.3 | 2026-03-20 | agentcp | 重建版，新增策略管理、应用管理、内容管理、通知管理模块联动关系，完整更新模块关系图和数据流 |
