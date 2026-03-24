# 设备证书管理 PRD

## 1. 功能概述
设备证书管理用于管理 MQTT/CoAP 等设备通信所需的 X.509 证书，支持证书的创建、上传、分配、续期和吊销操作，确保设备身份的安全认证。

## 2. 页面布局与交互

### 页面路径
`/devices/certificates` → `DeviceCertificates.vue`

### 搜索表单
| 字段 | 类型 | 说明 |
|------|------|------|
| 证书名称 | Input | 模糊搜索 |
| 证书类型 | Select | device / client / server / ca |
| 状态 | Select | active / expired / revoked / pending |
| 关键词 | Input | 搜索序列号/主题 |

### 数据表格
| 列 | 说明 |
|----|------|
| 证书名称 | cert_name |
| 证书类型 | cert_type（设备/客户端/服务器/CA） |
| 序列号 | serial_number |
| 主题 | subject (CN=xxx) |
| 颁发者 | issuer |
| 有效期 | not_before ~ not_after |
| 状态 | status 标签 |
| 操作 | 查看/下载/吊销/删除 |

### 新建/编辑弹窗
- 证书名称（必填）
- 证书类型（必填 Select）
- 证书文件（上传 PEM）
- 私钥文件（上传 PEM，仅创建时）
- 描述（可选）
- 到期提醒天数（默认30天）

### 按钮
- 「新建证书」主按钮（左上）
- 「下载」操作按钮
- 「吊销」操作按钮
- 「删除」操作按钮（需二次确认）
- 右侧「刷新」按钮

## 3. API 契约

### 证书列表
- 路径：`GET /api/v1/certificates`
- 请求参数：
  - `page` (int, default 1)
  - `page_size` (int, default 20)
  - `status` (string)
  - `cert_type` (string)
  - `keyword` (string)
- 响应：
```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "id": 1,
        "cert_id": "uuid",
        "cert_name": "设备证书-A区",
        "cert_type": "device",
        "serial_number": "00:1B:44:11:3A:B7",
        "subject": "CN=device-001",
        "issuer": "CN=Root CA",
        "thumbprint": "SHA1指纹",
        "not_before": "2024-01-01T00:00:00Z",
        "not_after": "2025-01-01T00:00:00Z",
        "status": "active",
        "description": ""
      }
    ],
    "total": 100,
    "page": 1,
    "page_size": 20
  },
  "message": "success"
}
```

### 获取单个证书
- 路径：`GET /api/v1/certificates/:id`
- 响应：同上 data 单条

### 创建证书
- 路径：`POST /api/v1/certificates`
- 请求体：
```json
{
  "cert_name": "设备证书-A区",
  "cert_type": "device",
  "cert_file_data": "base64编码的PEM证书",
  "key_file_data": "base64编码的PEM私钥",
  "description": "可选描述"
}
```
- 响应：`{ "code": 0, "data": {...}, "message": "创建成功" }`

### 更新证书
- 路径：`PUT /api/v1/certificates/:id`
- 请求体：
```json
{
  "cert_name": "新名称",
  "cert_type": "client",
  "description": "更新描述"
}
```

### 删除证书
- 路径：`DELETE /api/v1/certificates/:id`
- 响应：`{ "code": 0, "message": "删除成功" }`

### 吊销证书
- 路径：`POST /api/v1/certificates/:id/revoke`
- 响应：`{ "code": 0, "data": {...}, "message": "吊销成功" }`

### 获取即将到期证书
- 路径：`GET /api/v1/certificates/expiring`
- 参数：`days` (int, default 30)
- 响应：返回30天内到期的证书列表

### 验证证书
- 路径：`POST /api/v1/certificates/validate`
- 请求体：`{ "cert_id": "xxx" }` 或直接上传证书内容
- 响应：
```json
{
  "code": 0,
  "data": {
    "valid": true,
    "message": "证书有效",
    "not_before": "...",
    "not_after": "..."
  }
}
```

### 证书统计
- 路径：`GET /api/v1/certificates/stats`
- 响应：
```json
{
  "code": 0,
  "data": {
    "total": 100,
    "active": 80,
    "expired": 10,
    "revoked": 5,
    "expiring": 5
  }
}
```

## 4. 数据模型

### Certificate（证书表）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| cert_id | string | 证书唯一标识UUID |
| cert_name | string | 证书名称 |
| cert_type | string | device/client/server/ca |
| serial_number | string | 证书序列号 |
| subject | string | 主题CN |
| issuer | string | 颁发者 |
| thumbprint | string | SHA1指纹 |
| not_before | datetime | 生效时间 |
| not_after | datetime | 过期时间 |
| status | string | active/expired/revoked/pending |
| cert_file | string | 证书文件路径 |
| key_file | string | 私钥文件路径（不返回前端）|
| tenant_id | string | 租户ID |
| description | string | 描述 |

## 5. 验收标准
- [ ] 证书列表分页加载正常
- [ ] 支持按类型/状态/关键词搜索
- [ ] 新建证书成功（上传PEM文件）
- [ ] 查看证书详情正确（不暴露私钥）
- [ ] 吊销证书状态变更
- [ ] 删除证书（需二次确认）
- [ ] 即将到期证书提醒显示
- [ ] 证书验证功能正常
