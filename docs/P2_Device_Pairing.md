# 设备配对管理 PRD

## 1. 功能概述
设备配对管理用于管理宠物设备与宠物之间的绑定关系，支持一设备一宠物、多设备一宠物等场景，提供配对码生成、配对状态跟踪和解绑功能。

## 2. 页面布局与交互

### 页面路径
`/devices/pairing` → `DevicePairing.vue`

### 页面布局
- 左侧：配对列表
- 右侧：配对详情/操作面板

### 搜索表单
| 字段 | 类型 | 说明 |
|------|------|------|
| 设备ID | Input | 模糊搜索 |
| 宠物ID | Input | 模糊搜索 |
| 配对状态 | Select | paired / unpaired / pending |

### 数据表格
| 列 | 说明 |
|----|------|
| 设备ID | device_id |
| 设备型号 | hardware_model |
| 宠物昵称 | pet_name |
| 宠物品种 | pet_breed |
| 配对状态 | status |
| 配对时间 | paired_at |
| 操作 | 查看详情/解绑 |

### 配对详情面板
- 设备信息：设备ID、型号、固件版本
- 宠物信息：宠物昵称、品种、年龄、照片
- 配对时间线：配对时间、最近互动时间

### 按钮
- 「生成配对码」主按钮
- 「扫码配对」（打开扫码器）
- 「解绑」操作按钮（需确认）

## 3. API 契约

### 配对列表
- 路径：`GET /api/v1/device-pairings`
- 请求参数：`device_id`, `pet_id`, `status`, `page`, `page_size`
- 响应：
```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "id": 1,
        "pairing_id": "uuid",
        "device_id": "device-001",
        "pet_id": "pet-001",
        "pet_name": "小白",
        "pet_breed": "柯基",
        "status": "paired",
        "paired_at": "2024-01-01T00:00:00Z",
        "last_active_at": "2024-01-02T00:00:00Z"
      }
    ],
    "total": 50,
    "page": 1,
    "page_size": 20
  }
}
```

### 生成配对码
- 路径：`POST /api/v1/device-pairings/pairing-code`
- 请求体：
```json
{
  "device_id": "device-001",
  "expires_in": 300
}
```
- 响应：
```json
{
  "code": 0,
  "data": {
    "pairing_code": "ABC123",
    "qr_code": "base64二维码图片",
    "expires_at": "2024-01-01T00:05:00Z"
  }
}
```

### 扫码配对（设备端）
- 路径：`POST /api/v1/devices/:device_id/pair`
- 请求体：
```json
{
  "pet_id": "pet-001",
  "pairing_code": "ABC123"
}
```
- 响应：
```json
{
  "code": 0,
  "data": {
    "pairing_id": "uuid",
    "status": "paired"
  }
}
```

### 解绑配对
- 路径：`DELETE /api/v1/device-pairings/:id`
- 响应：`{ "code": 0, "message": "解绑成功" }`

### 获取配对详情
- 路径：`GET /api/v1/device-pairings/:id`
- 响应：
```json
{
  "code": 0,
  "data": {
    "pairing_id": "uuid",
    "device_id": "device-001",
    "pet_id": "pet-001",
    "status": "paired",
    "paired_at": "2024-01-01T00:00:00Z",
    "last_active_at": "2024-01-02T00:00:00Z",
    "device": { "hardware_model": "M5Stack", "firmware_version": "v1.2.3" },
    "pet": { "name": "小白", "breed": "柯基", "age": 3 }
  }
}
```

## 4. 数据模型

### DevicePairing（设备配对表）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| pairing_id | string | UUID |
| device_id | string | 设备ID |
| pet_id | uint | 宠物ID |
| status | string | pending/paired/unpaired |
| pairing_code | string | 配对码（临时）|
| paired_at | datetime | 配对时间 |
| last_active_at | datetime | 最近活跃时间 |
| tenant_id | string | 租户ID |

## 5. 验收标准
- [ ] 配对列表分页加载正常
- [ ] 按设备/宠物/状态筛选有效
- [ ] 生成配对码成功（有效期可配置）
- [ ] 配对码二维码正确生成
- [ ] 解绑操作需二次确认
- [ ] 配对详情显示完整设备+宠物信息
- [ ] 配对时间线展示正确
