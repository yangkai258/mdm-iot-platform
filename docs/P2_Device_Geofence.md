# 设备地理围栏 PRD

## 1. 功能概述
地理围栏（Geofence）功能允许管理员为设备设置虚拟边界区域，当设备进入或离开预设区域时，系统自动触发告警或执行特定动作。

## 2. 页面布局与交互

### 页面路径
`/devices/geofence` → `DeviceGeofence.vue`

### 页面布局
- 顶部：设备选择器 + 地图视图
- 左侧：围栏列表（已配置的围栏）
- 右侧：围栏配置面板

### 围栏配置表单
| 字段 | 类型 | 说明 |
|------|------|------|
| 围栏名称 | Input | 必填 |
| 设备 | Select | 关联的设备 |
| 围栏类型 | Select | 圆形 / 矩形 / 多边形 |
| 中心点坐标 | Lat/Lng Picker | 地图点击选择 |
| 半径/范围 | Input | 圆形为半径(米)，矩形为对角点 |
| 触发动作 | Multi-Select | 进入告警 / 离开告警 / 限制功能 |
| 启用状态 | Switch | 是否启用 |

### 数据表格（围栏列表）
| 列 | 说明 |
|----|------|
| 围栏名称 | name |
| 关联设备 | device_id |
| 围栏类型 | type（圆形/矩形/多边形）|
| 范围 | radius/size |
| 触发动作 | trigger_actions |
| 状态 | is_enabled（开关）|
| 操作 | 编辑/删除 |

### 按钮
- 「新建围栏」主按钮
- 地图上「绘制围栏」工具
- 「编辑」「删除」操作按钮

## 3. API 契约

### 围栏列表
- 路径：`GET /api/v1/geofences`
- 请求参数：
  - `device_id` (string)
  - `page` (int)
  - `page_size` (int)
- 响应：
```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "id": 1,
        "geofence_id": "uuid",
        "name": "家庭区域",
        "device_id": "device-001",
        "fence_type": "circle",
        "center": { "lat": 39.9, "lng": 116.4 },
        "radius": 100,
        "trigger_actions": ["enter_alert", "exit_alert"],
        "is_enabled": true,
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 10,
    "page": 1,
    "page_size": 20
  }
}
```

### 创建围栏
- 路径：`POST /api/v1/geofences`
- 请求体：
```json
{
  "name": "家庭区域",
  "device_id": "device-001",
  "fence_type": "circle",
  "center": { "lat": 39.9, "lng": 116.4 },
  "radius": 100,
  "trigger_actions": ["enter_alert", "exit_alert"],
  "is_enabled": true
}
```

### 更新围栏
- 路径：`PUT /api/v1/geofences/:id`
- 请求体：同上

### 删除围栏
- 路径：`DELETE /api/v1/geofences/:id`

### 获取设备当前围栏状态
- 路径：`GET /api/v1/devices/:device_id/geofence/status`
- 响应：
```json
{
  "code": 0,
  "data": {
    "device_id": "device-001",
    "inside_fences": ["家庭区域"],
    "last_location": { "lat": 39.9, "lng": 116.4 },
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

### 围栏事件记录
- 路径：`GET /api/v1/geofences/events`
- 参数：`device_id`, `fence_id`, `event_type`, `start_time`, `end_time`
- 响应：
```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "event_id": "uuid",
        "fence_name": "家庭区域",
        "event_type": "enter",
        "device_id": "device-001",
        "location": { "lat": 39.9, "lng": 116.4 },
        "occurred_at": "2024-01-01T08:00:00Z"
      }
    ]
  }
}
```

## 4. 数据模型

### Geofence（围栏表）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| geofence_id | string | UUID |
| name | string | 围栏名称 |
| device_id | string | 关联设备 |
| fence_type | string | circle/rectangle/polygon |
| center_lat | float | 中心纬度 |
| center_lng | float | 中心经度 |
| radius | float | 半径（米）|
| polygon_coords | JSON | 多边形坐标数组 |
| trigger_actions | JSON | 触发动作 |
| is_enabled | bool | 是否启用 |
| tenant_id | string | 租户ID |

### GeofenceEvent（围栏事件）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| geofence_id | uint | 围栏ID |
| device_id | string | 设备ID |
| event_type | string | enter/exit |
| location_lat | float | 触发位置纬度 |
| location_lng | float | 触发位置经度 |
| occurred_at | datetime | 发生时间 |

## 5. 验收标准
- [ ] 地图上绘制圆形/矩形/多边形围栏
- [ ] 围栏列表分页加载正常
- [ ] 创建/编辑/删除围栏成功
- [ ] 围栏启用/禁用开关生效
- [ ] 设备进入/离开围栏触发事件记录
- [ ] 围栏事件列表可查询
- [ ] 与设备监控面板联动显示位置
