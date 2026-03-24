# 固件兼容性 PRD

## 1. 功能概述
固件兼容性模块管理设备型号与固件版本之间的兼容关系，确保 OTA 升级只推送与设备硬件匹配的固件，防止变砖。

## 2. 页面布局与交互

### 页面路径
`/ota/compatibility` → `OtaCompatibilityView.vue`（后端待实现）

### 兼容性规则列表
| 列 | 说明 |
|----|------|
| 设备型号 | hardware_model |
| 固件版本 | firmware_version |
| 兼容状态 | compatible |
| 最低版本要求 | min_version |
| 操作 | 编辑/删除 |

### 添加规则弹窗
| 字段 | 类型 | 说明 |
|------|------|------|
| 设备型号 | Select | M5Stack/ESP32等 |
| 固件版本 | Input | 如 v1.2.3 |
| 兼容状态 | Switch | 是否兼容 |
| 最低可升级版本 | Input | - |

## 3. API 契约

### 兼容性规则列表
- 路径：`GET /api/v1/ota/compatibility`
- 参数：`hardware_model`, `page`, `page_size`
- 响应：
```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "id": 1,
        "hardware_model": "M5Stack",
        "firmware_version": "v1.2.3",
        "compatible": true,
        "min_upgrade_version": "v1.0.0",
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 20
  }
}
```

### 添加兼容性规则
- 路径：`POST /api/v1/ota/compatibility`
- 请求体：
```json
{
  "hardware_model": "M5Stack",
  "firmware_version": "v1.2.3",
  "compatible": true,
  "min_upgrade_version": "v1.0.0"
}
```

### 批量导入
- 路径：`POST /api/v1/ota/compatibility/import`
- 请求：CSV文件

### 检查兼容性
- 路径：`GET /api/v1/ota/compatibility/check`
- 参数：`hardware_model`, `current_version`, `target_version`
- 响应：
```json
{
  "code": 0,
  "data": {
    "compatible": true,
    "warnings": ["建议从v1.1.0升级"],
    "blocking_reasons": []
  }
}
```

## 4. 数据模型

### OTACompatibility（固件兼容性表）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| hardware_model | string | 设备型号 |
| firmware_version | string | 固件版本 |
| compatible | bool | 是否兼容 |
| min_upgrade_version | string | 最低可升级版本 |
| notes | string | 备注 |

## 5. 验收标准
- [ ] 兼容性规则列表加载正常
- [ ] 添加/编辑规则成功
- [ ] 批量导入CSV成功
- [ ] 兼容性检查正确返回结果
- [ ] 不兼容时阻止OTA升级
