# 空间认知 PRD

## 1. 功能概述
空间认知模块为设备提供环境地图构建、定位和导航能力，使设备能够在未知环境中自主探索和定位。

## 2. 页面布局与交互

### 页面路径
`/embodied/spatial` → `MapView.vue` + `NavigationView.vue`

### 地图视图
- 2D/3D地图切换
- 设备位置实时显示
- 禁区标注（红色区域）
- 导航路径显示

### 地图管理
- 地图列表
- 「新建地图」按钮
- 地图版本历史

## 3. API 契约

### 获取地图
- 路径：`GET /api/v1/embodied/:device_id/map`
- 实际路径（embodied_controller.go）：`GET /api/v1/embodied/:device_id/map`
- 响应：
```json
{
  "code": 0,
  "data": {
    "device_id": "device-001",
    "map_type": "grid",
    "map_data": { "cells": [...] },
    "resolution": 0.05,
    "version": 3,
    "is_active": true
  }
}
```

### 更新地图
- 路径：`POST /api/v1/embodied/:device_id/map/update`
- 请求体：`{ "map_data": {...}, "version": 3 }`

### 获取定位
- 路径：`GET /api/v1/embodied/:device_id/localization`
- 响应：`{ "device_id": "xxx", "position": {x, y, z}, "orientation": 90, "confidence": 0.98 }`

### 校准定位
- 路径：`POST /api/v1/embodied/:device_id/localization/calibrate`
- 请求体：`{ "reference_points": [...] }`

## 4. 数据模型

### EmbodiedMap（环境地图）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| device_id | string | 设备ID |
| map_type | string | grid/semantic/topological |
| map_data | JSONB | 地图数据 |
| resolution | float | 分辨率 |
| size | JSONB | 地图尺寸 |
| version | int | 版本号 |
| is_active | bool | 是否当前地图 |

### SpatialPosition（空间位置）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| device_id | string | 设备ID |
| position_x/y/z | float | 位置坐标 |
| orientation | float | 朝向角度 |
| confidence | float | 置信度 |
| recorded_at | datetime | 记录时间 |

## 5. 验收标准
- [ ] 地图正确加载和显示
- [ ] 设备位置实时更新
- [ ] 地图版本管理正常
- [ ] 定位校准成功
- [ ] 导航路径正确显示
