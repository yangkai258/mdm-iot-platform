# 用量计费 PRD

## 1. 功能概述
用量计费模块记录和管理用户的各项费用账单，包括订阅费用、API调用费用、存储费用等，提供费用明细查询、发票管理和在线支付功能。

## 2. 页面布局与交互

### 页面路径
`/billing` → `BillingList.vue`

### Tab1：账单记录
#### 搜索表单
| 字段 | 类型 | 说明 |
|------|------|------|
| 用户ID | Input | 精确搜索 |
| 计费类型 | Select | subscription/usage/API_quota |
| 账单状态 | Select | pending/paid/overdue |
| 时间范围 | DateRange | 计费周期 |

#### 数据表格
| 列 | 说明 |
|----|------|
| 账单号 | invoice_no |
| 用户ID | user_id |
| 计费类型 | billing_type |
| 金额 | amount |
| 货币 | currency |
| 周期 | period_start ~ period_end |
| 状态 | status（待支付/已支付/逾期）|
| 操作 | 查看详情/支付/开票 |

### Tab2：发票管理
#### 数据表格
| 列 | 说明 |
|----|------|
| 发票号 | invoice_no |
| 发票类型 | invoice_type（普通/VAT）|
| 抬头 | title |
| 金额 | total_amount（含税）|
| 状态 | status（待审核/已开票/已作废）|
| 开票日期 | issue_date |
| 操作 | 查看/寄送/作废 |

#### 新建发票弹窗
| 字段 | 类型 | 说明 |
|------|------|------|
| 账单 | Select | 选择关联账单 |
| 发票类型 | Select | 普通发票/VAT专用发票 |
| 发票抬头 | Input | 必填 |
| 税号 | Input | VAT必填 |
| 开户行 | Input | VAT必填 |
| 银行账号 | Input | VAT必填 |
| 地址电话 | Input | VAT必填 |
| 接收邮箱 | Input | 接收电子发票 |

### 费用汇总（顶部卡片）
- 本月应结：¥12,345.00
- 待支付：¥3,000.00
- 已支付（本年）：¥50,000.00
- 开票中：¥5,000.00

## 3. API 契约

### 账单记录列表
- 路径：`GET /api/v1/billing/records`
- 参数：`user_id`, `billing_type`, `status`, `page`, `page_size`
- 响应：
```json
{
  "code": 0,
  "data": [
    {
      "id": 1,
      "user_id": 100,
      "billing_type": "subscription",
      "reference_id": 1,
      "amount": 299.00,
      "currency": "CNY",
      "period_start": "2024-01-01T00:00:00Z",
      "period_end": "2024-02-01T00:00:00Z",
      "status": "paid",
      "paid_at": "2024-01-01T00:00:00Z",
      "invoice_id": 1
    }
  ]
}
```

### 获取账单详情
- 路径：`GET /api/v1/billing/records/:id`

### 支付账单
- 路径：`POST /api/v1/billing/pay`
- 请求体：`{ "record_id": 1, "payment_method": "alipay" }`
- 响应：`{ "code": 0, "data": { "payment_url": "https://..." } }`

### 发票列表
- 路径：`GET /api/v1/billing/invoices`
- 参数：`user_id`, `status`, `page`, `page_size`
- 响应：
```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "id": 1,
        "invoice_no": "FP2024010001",
        "user_id": 100,
        "amount": 299.00,
        "tax_amount": 29.00,
        "total_amount": 328.00,
        "status": "issued",
        "invoice_type": "normal",
        "title": "个人/公司名称",
        "tax_no": "",
        "issue_date": "2024-01-01T00:00:00Z",
        "pdf_url": "https://..."
      }
    ],
    "total": 10
  }
}
```

### 创建发票申请
- 路径：`POST /api/v1/billing/invoices`
- 请求体：`{ "user_id": 100, "subscription_id": 1, "invoice_type": "normal", "title": "张三", "email": "zhangsan@example.com" }`

### 审核发票
- 路径：`POST /api/v1/billing/invoices/:id/review`
- 请求体：`{ "approved": true, "comment": "审核通过" }`

### 开具发票
- 路径：`POST /api/v1/billing/invoices/:id/issue`

### 作废发票
- 路径：`POST /api/v1/billing/invoices/:id/void`
- 请求体：`{ "reason": "开票信息错误" }`

### 发票寄送
- 路径：`POST /api/v1/billing/invoices/:id/ship`
- 请求体：`{ "carrier": "顺丰", "tracking_no": "SF123456789", "recipient_name": "张三", "recipient_phone": "13800138000", "recipient_addr": "北京市朝阳区xxx" }`

### 费用汇总
- 路径：`GET /api/v1/billing/summary`
- 响应：
```json
{
  "code": 0,
  "data": {
    "monthly_total": 12345.00,
    "pending_amount": 3000.00,
    "paid_this_year": 50000.00,
    "invoice_pending": 5000.00
  }
}
```

## 4. 数据模型

### BillingRecord（计费记录）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| user_id | uint | 用户ID |
| billing_type | string | subscription/usage/API_quota |
| reference_id | uint | 关联订单/订阅ID |
| amount | float64 | 金额 |
| currency | string | 货币 |
| period_start | datetime | 周期开始 |
| period_end | datetime | 周期结束 |
| status | string | pending/paid/overdue |
| paid_at | datetime | 支付时间 |
| invoice_id | uint | 关联发票ID |

### Invoice（发票）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| invoice_no | string | 发票号 |
| user_id | uint | 用户ID |
| subscription_id | uint | 订阅ID |
| amount | float64 | 金额 |
| tax_amount | float64 | 税额 |
| total_amount | float64 | 含税总额 |
| status | string | pending/issued/void |
| invoice_type | string | normal/VAT |
| title | string | 发票抬头 |
| tax_no | string | 税号 |
| bank_name | string | 开户行 |
| bank_account | string | 银行账号 |
| email | string | 接收邮箱 |
| pdf_url | string | PDF链接 |
| issue_date | datetime | 开票日期 |

## 5. 验收标准
- [ ] 账单列表分页加载正常
- [ ] 账单按类型/状态筛选有效
- [ ] 创建发票申请成功
- [ ] 发票审核通过/拒绝正常
- [ ] 开具发票生成PDF
- [ ] 发票寄送记录完整
- [ ] 费用汇总数据正确
- [ ] 在线支付跳转正常
