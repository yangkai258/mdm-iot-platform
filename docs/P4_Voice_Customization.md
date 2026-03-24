# 声音定制 PRD

## 1. 功能概述
声音定制模块为宠物设备提供个性化的语音合成配置，包括音色选择、语速调节、音调控制等参数。

## 2. 页面布局与交互

### 页面路径
`/market/voices` → `VoiceConfigView.vue`

### 声音配置
- 内置音色选择（可爱/温柔/活泼/磁性）
- 语速滑块
- 音调滑块
- 预览播放按钮

### 我的声音
- 已配置声音列表
- 设为默认按钮

## 3. API 契约

### 声音配置列表
- 路径：`GET /api/v1/market/voices`
- 响应：
```json
{
  "code": 0,
  "items": [
    {
      "id": 1,
      "name": "可爱音色",
      "voice_type": "cute",
      "speed": 1.0,
      "pitch": 1.2,
      "is_default": true,
      "preview_url": "https://..."
    }
  ]
}
```

### 创建声音配置
- 路径：`POST /api/v1/market/voices`
- 请求体：`{ "name": "我的音色", "voice_type": "cute", "speed": 1.0, "pitch": 1.2 }`

### 应用声音配置
- 路径：`POST /api/v1/market/voices/:id/apply`
- 请求体：`{ "device_id": "xxx" }`

## 4. 验收标准
- [ ] 声音配置列表加载正常
- [ ] 预览播放正常
- [ ] 配置保存成功
- [ ] 设为默认成功
- [ ] 应用到设备成功
