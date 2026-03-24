# 智能家居 PRD

## 1. 功能概述
智能家居集成模块支持宠物设备与主流智能家居平台（米家/天猫精灵/HomeKit）的互联互通，实现场景联动控制。

## 2. 页面布局与交互

### 页面路径
`/integration/smart-home` → `SmartHomeView.vue`

### 设备管理
- 第三方平台选择
- 设备授权
- 设备列表

### 场景配置
- 场景列表
- 触发条件配置
- 联动动作配置

## 3. API 契约

### 第三方集成列表
- 路径：`GET /api/v1/integration/platforms`
- 响应：
```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "platform": "xiaomi",
        "name": "米家",
        "status": "connected",
        "devices_count": 5
      }
    ]
  }
}
```

### 授权设备
- 路径：`POST /api/v1/integration/:platform/authorize`
- 请求体：`{ "auth_code": "xxx" }`

### 设备列表
- 路径：`GET /api/v1/integration/:platform/devices`
- 响应：第三方平台设备列表

### 发送设备控制
- 路径：`POST /api/v1/integration/:platform/devices/:id/control`
- 请求体：`{ "command": "turn_on" }`

## 4. 验收标准
- [ ] 平台授权流程正常
- [ ] 设备列表正确显示
- [ ] 设备控制正常
- [ ] 场景联动正常
