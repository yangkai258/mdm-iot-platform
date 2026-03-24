# 环境感知 PRD

## 1. 功能概述
环境感知模块使设备具备感知周围环境的能力，包括视觉感知（摄像头数据）、深度感知（TOF/结构光）、触觉感知等，为具身智能提供环境输入。

## 2. 页面布局与交互

### 页面路径
`/embodied/perception` → `PerceptionView.vue`

### 感知数据展示
- 设备选择器
- 感知类型切换（视觉/深度/触觉）
- 实时画面/点云/触觉数据
- 感知结果标注

## 3. API 契约

### 获取当前感知
- 路径：`GET /api/v1/embodied/:device_id/perception`
- 实际路径（embodied_controller.go）：`GET /api/v1/embodied/:device_id/perception`
- 参数：`type` (visual/depth/touch)
- 响应：
```json
{
  "code": 0,
  "data": {
    "device_id": "device-001",
    "type": "visual",
    "data": "base64编码的图像数据",
    "objects_detected": [
      { "class": "person", "bbox": [x, y, w, h], "confidence": 0.95 }
    ],
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

### 上报视觉感知
- 路径：`POST /api/v1/embodied/:device_id/perception/visual`
- 请求体：`{ "data": "base64...", "objects_detected": [...] }`

### 上报深度感知
- 路径：`POST /api/v1/embodied/:device_id/perception/depth`
- 请求体：`{ "data": "点云数据", "distance_map": [...] }`

### 上报触觉感知
- 路径：`POST /api/v1/embodied/:device_id/perception/touch`
- 请求体：`{ "pressure": [...], "temperature": [...] }`

## 4. 验收标准
- [ ] 感知类型切换正常
- [ ] 实时数据正确展示
- [ ] 数据上报接口正常
- [ ] 感知结果正确解析
