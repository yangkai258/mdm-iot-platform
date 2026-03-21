# MiniClaw通信协议

**版本：** V1.1  
**模块负责人：** agentcp  
**编制日期：** 2026-03-22  

---

## 1. 概述

MiniClaw通信协议定义了MDM中台与MiniClaw设备之间的MQTT通信规范，包括Topic设计、消息格式、命令类型、心跳管理和OTA升级协议。协议是设备与云端通信的标准语言，确保所有指令和数据的正确传递。

**业务目标：**
- 统一MQTT通信Topic结构
- 定义设备与云端的消息格式
- 规范动作指令下发和状态上报
- 实现设备心跳和在线管理

---

## 2. 功能列表

| 功能 | 描述 | 优先级 | 触发方式 | 前端入口/按钮 |
|------|------|--------|----------|--------------|
| MQTT连接管理 | 设备连接/断开/认证管理 | P0 | 自动 | 无按钮 |
| 心跳上报 | 设备定期心跳维持在线状态 | P0 | 自动 | 无按钮 |
| 语音识别上报 | 设备语音识别结果上报 | P0 | 自动 | 无按钮 |
| 传感器数据上报 | 设备传感器数据实时上报 | P0 | 自动 | 无按钮 |
| 设备状态上报 | 设备状态变化上报 | P0 | 自动 | 无按钮 |
| 动作指令下发 | 云端动作指令下发到设备 | P0 | 自动 | 快捷指令 |
| 语音播报下发 | 云端语音内容下发到设备TTS播放 | P0 | 人工 | 控制台对话 |
| 配置更新下发 | 云端配置参数下发到设备 | P1 | 人工 | 管理后台 |
| 命令确认机制 | 设备对云端命令的ACK确认 | P0 | 自动 | 无按钮 |

---

## 3. MQTT Topic设计

### 3.1 Topic结构规范

```
/miniclaw/{device_id}/{direction}/{type}
```

**方向定义：**
- `up` - 设备上行（设备 -> 云端）
- `down` - 云端下行（云端 -> 设备）

### 3.2 Topic详细定义

| Topic | 方向 | 描述 | QoS | 保留 |
|-------|------|------|-----|------|
| `/miniclaw/{device_id}/up/voice` | up | 语音识别结果上报 | 1 | 否 |
| `/miniclaw/{device_id}/up/sensor` | up | 传感器数据上报 | 1 | 否 |
| `/miniclaw/{device_id}/up/status` | up | 设备状态上报 | 1 | 是 |
| `/miniclaw/{device_id}/up/ack` | up | 命令执行结果确认 | 1 | 否 |
| `/miniclaw/{device_id}/up/log` | up | 设备日志上报 | 1 | 否 |
| `/miniclaw/{device_id}/down/action` | down | 云端动作指令下发 | 1 | 否 |
| `/miniclaw/{device_id}/down/speech` | down | 云端语音内容下发 | 1 | 否 |
| `/miniclaw/{device_id}/down/config` | down | 云端配置参数下发 | 1 | 否 |
| `/miniclaw/{device_id}/down/ping` | down | 心跳ping请求 | 1 | 否 |

### 3.3 通用Topic（无需device_id）

| Topic | 方向 | 描述 | QoS |
|-------|------|------|-----|
| `/miniclaw/broadcast/system` | down | 系统广播消息 | 1 |
| `/miniclaw/broadcast/ota` | down | OTA全量推送 | 1 |

---

## 4. 消息格式定义

### 4.1 通用消息头

所有消息统一使用以下JSON格式：

```json
{
  "msg_id": "msg-uuid-001",
  "device_id": "pet-001",
  "timestamp": 1710902400000,
  "type": "voice/sensor/status/action/speech/config",
  "version": "1.0",
  "payload": {}
}
```

### 4.2 协议帧格式

**完整协议帧结构：**

```
[帧头(2B)] [长度(2B)] [版本(1B)] [类型(1B)] [序列号(4B)] [Payload(N)] [校验(2B)]
```

| 字段 | 长度 | 说明 |
|------|------|------|
| 帧头 | 2字节 | 固定值 0xAA 0x55 |
| 长度 | 2字节 | 帧总长度（不含帧头和校验） |
| 版本 | 1字节 | 协议版本号，当前1.0 |
| 类型 | 1字节 | 消息类型（见下表） |
| 序列号 | 4字节 | 消息序号，用于追踪和ACK |
| Payload | N字节 | 实际数据内容 |
| 校验 | 2字节 | CRC16校验 |

**消息类型定义：**

| 类型值 | 类型名称 | 方向 | 说明 |
|--------|----------|------|------|
| 0x01 | STATUS_REPORT | up | 设备状态上报 |
| 0x02 | VOICE_DATA | up | 语音数据上报 |
| 0x03 | SENSOR_DATA | up | 传感器数据上报 |
| 0x04 | ACK | up | 命令确认 |
| 0x05 | LOG_REPORT | up | 日志上报 |
| 0x11 | ACTION_CMD | down | 动作指令 |
| 0x12 | SPEECH_CMD | down | 语音播报指令 |
| 0x13 | CONFIG_CMD | down | 配置指令 |
| 0x14 | PING_REQ | down | 心跳请求 |
| 0x15 | PING_RSP | down | 心跳响应 |

### 4.3 语音识别上报 (up/voice)

**设备 -> 云端**

```json
{
  "msg_id": "msg-uuid-001",
  "device_id": "pet-001",
  "timestamp": 1710902400000,
  "type": "voice",
  "version": "1.0",
  "payload": {
    "text": "今天天气怎么样",
    "confidence": 0.95,
    "audio_duration_ms": 2500,
    "is_final": true
  }
}
```

### 4.4 传感器数据上报 (up/sensor)

**设备 -> 云端**

```json
{
  "msg_id": "msg-uuid-002",
  "device_id": "pet-001",
  "timestamp": 1710902400000,
  "type": "sensor",
  "version": "1.0",
  "payload": {
    "sensors": [
      {
        "type": "ultrasonic",
        "value": 15.5,
        "unit": "cm",
        "timestamp": 1710902400000
      },
      {
        "type": "accelerometer",
        "value": { "x": 0.1, "y": 0.0, "z": 9.8 },
        "unit": "m/s²",
        "timestamp": 1710902400000
      },
      {
        "type": "battery",
        "value": 85,
        "unit": "%",
        "timestamp": 1710902400000
      },
      {
        "type": "touch",
        "value": true,
        "timestamp": 1710902400000
      }
    ]
  }
}
```

### 4.5 设备状态上报 (up/status)

**设备 -> 云端**

```json
{
  "msg_id": "msg-uuid-003",
  "device_id": "pet-001",
  "timestamp": 1710902400000,
  "type": "status",
  "version": "1.0",
  "payload": {
    "online": true,
    "firmware_version": "1.3.0",
    "hardware_version": "v1.3",
    "battery_level": 85,
    "battery_status": "normal",
    "position": { "x": 120.5, "y": 80.3 },
    "current_expression": "happy",
    "current_action": "idle",
    "mood": 75,
    "energy": 80,
    "error_codes": [],
    "ota_progress": null
  }
}
```

### 4.6 命令执行确认 (up/ack)

**设备 -> 云端**

```json
{
  "msg_id": "msg-uuid-004",
  "device_id": "pet-001",
  "timestamp": 1710902400000,
  "type": "ack",
  "version": "1.0",
  "payload": {
    "original_msg_id": "cmd-uuid-001",
    "command": "action",
    "status": "success",
    "error_code": null,
    "error_message": null,
    "execution_time_ms": 350
  }
}
```

### 4.7 动作指令下发 (down/action)

**云端 -> 设备**

```json
{
  "msg_id": "cmd-uuid-001",
  "device_id": "pet-001",
  "timestamp": 1710902400000,
  "type": "action",
  "version": "1.0",
  "payload": {
    "action_id": "walk_forward",
    "action_name": "前进",
    "params": {
      "distance": 30,
      "speed": 3
    },
    "sequence_id": "seq-001",
    "priority": 5,
    "is_emergency": false
  }
}
```

### 4.8 语音播报下发 (down/speech)

**云端 -> 设备**

```json
{
  "msg_id": "cmd-uuid-002",
  "device_id": "pet-001",
  "timestamp": 1710902400000,
  "type": "speech",
  "version": "1.0",
  "payload": {
    "text": "今天天气晴朗，气温20到25度，适合出门散步哦",
    "voice_id": "xiaozhua",
    "volume": 80,
    "speed": 1.0,
    "emotion": "happy"
  }
}
```

### 4.9 配置更新下发 (down/config)

**云端 -> 设备**

```json
{
  "msg_id": "cmd-uuid-003",
  "device_id": "pet-001",
  "timestamp": 1710902400000,
  "type": "config",
  "version": "1.0",
  "payload": {
    "config_type": "ota",
    "config": {
      "ota_enabled": true,
      "ota_url": "https://cdn.example.com/firmware/v1.3.0.bin",
      "ota_version": "1.3.0",
      "ota_checksum": "d41d8cd98f00b204e9800998ecf8427e",
      "ota_file_size": 1048576
    }
  }
}
```

**常规配置下发：**
```json
{
  "msg_id": "cmd-uuid-004",
  "device_id": "pet-001",
  "timestamp": 1710902400000,
  "type": "config",
  "version": "1.0",
  "payload": {
    "config_type": "general",
    "config": {
      "heartbeat_interval": 30,
      "sensor_report_interval": 1000,
      "expression": "happy",
      "volume": 80
    }
  }
}
```

---

## 5. 心跳与在线管理

### 5.1 心跳机制

```
┌──────────────┐                              ┌──────────────┐
│  MiniClaw    │                              │   MDM中台    │
│    设备      │                              │              │
└──────┬───────┘                              └──────┬───────┘
       │                                             │
       │  1. MQTT CONNECT (clientId=pet-001)        │
       │────────────────────────────────────────────>
       │                                             │
       │  2. CONNACK (ack=0)                        │
       │<────────────────────────────────────────────│
       │                                             │
       │  3. PUBLISH /up/status (online:true)       │
       │     QoS=1, RETAIN=true                     │
       │────────────────────────────────────────────>
       │                                             │
       │  4. PUBCOMP (msg_id)                       │
       │<────────────────────────────────────────────│
       │                                             │
       │  [每30秒重复步骤3-4]                        │
       │                                             │
       │  5. MQTT DISCONNECT                        │
       │────────────────────────────────────────────>
       │                                             │
       │  6. 设备离线，MDM检测心跳超时(90秒)         │
       │     更新设备影子 online=false               │
```

### 5.2 心跳参数

| 参数 | 值 | 说明 |
|------|-----|------|
| 心跳间隔 | 30秒 | 设备上报状态间隔 |
| 心跳超时 | 90秒 | 3倍间隔内无消息判定离线 |
| 在线阈值 | 连续2次心跳 | 判定设备在线 |
| 离线阈值 | 连续3次无心跳 | 判定设备离线 |

### 5.3 设备影子

设备在线状态存储在Redis设备影子中：

```
Key: device_shadow:{device_id}
TTL: 120秒
Value: {
  "online": true,
  "last_heartbeat": 1710902400000,
  "firmware_version": "1.3.0",
  "battery_level": 85,
  "position": {"x": 120.5, "y": 80.3}
}
```

---

## 6. 错误码定义

### 6.1 设备端错误码

| 错误码 | 名称 | 说明 |
|--------|------|------|
| 0 | SUCCESS | 执行成功 |
| 1001 | ERR_INVALID_CMD | 无效命令 |
| 1002 | ERR_PARAM | 参数错误 |
| 1003 | ERR_EXEC_TIMEOUT | 执行超时 |
| 2001 | ERR_MOTOR | 电机故障 |
| 2002 | ERR_SENSOR | 传感器故障 |
| 2003 | ERR_BATTERY_LOW | 电量过低 |
| 3001 | ERR_OTA_DOWNLOAD | OTA下载失败 |
| 3002 | ERR_OTA_FLASH | OTA刷写失败 |
| 3003 | ERR_OTA_VERIFY | OTA校验失败 |

---

## 7. 接口定义

### 7.1 设备连接认证

```
POST /api/v1/devices/{device_id}/auth
```

**请求参数：**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| device_secret | string | 是 | 设备密钥 |
| firmware_version | string | 是 | 固件版本 |
| hardware_version | string | 是 | 硬件版本 |

**请求示例：**
```json
{
  "device_secret": "xxxx",
  "firmware_version": "1.3.0",
  "hardware_version": "v1.3"
}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "mqtt_broker": "mqtt.mdm.example.com",
    "mqtt_port": 8883,
    "mqtt_username": "pet-001",
    "mqtt_password": "token-xxxxx",
    "mqtt_topic_prefix": "/miniclaw/pet-001"
  }
}
```

### 7.2 手动触发设备上报

```
POST /api/v1/devices/{device_id}/report
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "msg_id": "cmd-uuid-100",
    "type": "ping",
    "sent_at": "2026-03-20T10:55:00Z"
  }
}
```

### 7.3 下发动作指令

```
POST /api/v1/devices/{device_id}/commands/action
```

**请求示例：**
```json
{
  "action_id": "walk_forward",
  "params": { "distance": 30, "speed": 3 }
}
```

### 7.4 下发语音播报

```
POST /api/v1/devices/{device_id}/commands/speech
```

**请求示例：**
```json
{
  "text": "今天天气晴朗",
  "volume": 80
}
```

### 7.5 下发配置

```
POST /api/v1/devices/{device_id}/commands/config
```

**请求示例：**
```json
{
  "config": {
    "heartbeat_interval": 30,
    "volume": 80
  }
}
```

---

## 8. 流程图

### 8.1 完整通信流程

```
用户输入 -> 控制台/云端 -> MQTT Broker
                          |
           ┌──────────────┼──────────────┐
           ▼              ▼              ▼
      /down/speech   /down/action   /down/config
        语音下发       动作下发        配置下发
           │              │              │
           └──────────────┼──────────────┘
                          │
                          ▼
               ┌────────────────────────┐
               │    MiniClaw 设备        │
               │    - TTS 语音播放       │
               │    - 执行动作           │
               │    - 传感器采集         │
               └───────────┬────────────┘
                           │
          ┌────────────────┼────────────────┐
          ▼                ▼                ▼
    /up/voice         /up/sensor        /up/status
      语音识别          传感器数据         状态上报
          │                │                │
          └────────────────┴────────────────┘
                            │
                            ▼
                 ┌────────────────────────┐
                 │     OpenClaw AI层     │
                 │  - 对话引擎            │
                 │  - 行为引擎            │
                 │  - 记忆库              │
                 └────────────────────────┘
```

### 8.2 心跳在线管理流程

```
MiniClaw设备 -> MQTT Broker -> Redis设备影子
       │
       MQTT CONNECT
       │
       PUBLISH /up/status (online:true)
       │
       SET device_shadow (TTL:120s)
       │
       [每30秒重复]
       │
       PUBLISH /up/status
       │
       EXPIRE更新
       │
       [设备关机]
       PUBLISH /up/status (online:false)
       │
       DEL device_shadow
       │
       [90秒无心跳]
       GET device_shadow (nil)
       判定离线，告警
```

---

## 9. 模块联动

### 9.1 与固件管理(MINICLAW_FIRMWARE)联动

- **触发时机：** OTA升级时
- **联动内容：** 固件管理通过 `/down/config` 下发OTA配置，设备通过 `/up/status` 上报升级进度
- **数据流向：** 固件管理 -> MQTT Broker -> MiniClaw

### 9.2 与行为引擎(PET_BEHAVIOR_ENGINE)联动

- **触发时机：** 动作下发时
- **联动内容：** 行为引擎通过 `/down/action` 下发动作指令，接收设备 `/up/ack` 确认执行结果
- **数据流向：** 行为引擎 -> MQTT Broker -> MiniClaw

### 9.3 与对话引擎联动

- **触发时机：** 语音播报时
- **联动内容：** 对话引擎通过 `/down/speech` 下发TTS文本，设备通过 `/up/voice` 上报用户语音
- **数据流向：** 对话引擎 <-> MQTT Broker <-> MiniClaw

### 9.4 与记忆库(PET_MEMORY)联动

- **触发时机：** 传感器数据上报时
- **联动内容：** 记忆库记录传感器时序数据，用于上下文理解和长期学习
- **数据流向：** MiniClaw -> MQTT -> 记忆库

---

## 10. 验收标准

### 10.1 功能验收

| 功能 | 验收条件 |
|------|----------|
| MQTT连接 | 设备可在5秒内完成连接认证 |
| 心跳维持 | 设备在线状态准确率>=99.9% |
| 指令下发 | QoS=1保证指令可靠到达 |
| 命令确认 | 设备在3秒内返回ACK |
| 语音上报 | 语音识别结果延迟<=500ms |

### 10.2 性能验收

- MQTT消息吞吐量：10000 msg/s
- 单设备消息延迟：<=100ms
- 支持同时在线设备：100000台

### 10.3 安全性验收

- MQTT连接使用用户名+Token认证
- TLS加密传输（预留）
- 消息签名验证（预留）

---

## 11. 技术参数汇总

### 11.1 MQTT配置

| 参数 | 值 |
|------|-----|
| Broker | EMQX 5.0 |
| 协议 | MQTT 3.1.1 / MQTT 5.0 |
| 端口 | 8883 (TLS) / 1883 (Plain) |
| QoS | 0=最多一次, 1=至少一次, 2=恰好一次 |
| Keepalive | 60秒 |
| Clean Session | false (持久会话) |

### 11.2 Topic配额

| 类型 | 限制 |
|------|------|
| 单设备Topic数 | <=100 |
| 单用户订阅数 | <=500 |
| Topic最大长度 | 256字符 |
| 消息最大Payload | 64KB |

### 11.3 传感器类型定义

| 类型 | 说明 | 单位 |
|------|------|------|
| ultrasonic | 超声波测距 | cm |
| accelerometer | 加速度计 | m/s² |
| gyroscope | 陀螺仪 | deg/s |
| battery | 电池电量 | % |
| temperature | 温度 | ℃ |
| touch | 触摸传感器 | bool |
| fall | 跌落检测 | bool |
| collision | 碰撞检测 | bool |
