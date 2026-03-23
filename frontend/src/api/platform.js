/**
 * 开放平台 API
 * 开发者控制台 / Webhook 市场 / API 文档
 * Base: /api/v1
 */

const API_BASE = '/api/v1'

function getToken() {
  return localStorage.getItem('token') || ''
}

async function request(url, options = {}) {
  const res = await fetch(url, {
    ...options,
    headers: {
      'Authorization': `Bearer ${getToken()}`,
      'Content-Type': 'application/json',
      ...(options.headers || {})
    }
  })
  const data = await res.json()
  if (data.code !== 0 && data.code !== 200) {
    throw new Error(data.message || '请求失败')
  }
  return data
}

// ════════════════════════════════════════════════════════════
// 开发者应用管理
// ════════════════════════════════════════════════════════════

/** 创建应用 */
export async function createApp(data) {
  return request(`${API_BASE}/developer/apps`, { method: 'POST', body: JSON.stringify(data) })
}

/** 应用列表 */
export async function getAppList(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`${API_BASE}/developer/apps?${qs}`)
}

/** 应用详情 */
export async function getAppDetail(id) {
  return request(`${API_BASE}/developer/apps/${id}`)
}

/** 更新应用 */
export async function updateApp(id, data) {
  return request(`${API_BASE}/developer/apps/${id}`, { method: 'PUT', body: JSON.stringify(data) })
}

/** 删除应用 */
export async function deleteApp(id) {
  return request(`${API_BASE}/developer/apps/${id}`, { method: 'DELETE' })
}

/** 生成 API Key */
export async function createApiKey(appId, data = {}) {
  return request(`${API_BASE}/developer/apps/${appId}/keys`, { method: 'POST', body: JSON.stringify(data) })
}

/** 删除 API Key */
export async function deleteApiKey(appId, keyId) {
  return request(`${API_BASE}/developer/apps/${appId}/keys/${keyId}`, { method: 'DELETE' })
}

/** API Key 列表 */
export async function getApiKeyList(appId) {
  return request(`${API_BASE}/developer/apps/${appId}/keys`)
}

/** 使用统计 */
export async function getDeveloperStats(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`${API_BASE}/developer/stats?${qs}`)
}

/** 审计日志 */
export async function getAuditLogs(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`${API_BASE}/developer/audit-logs?${qs}`)
}

// ════════════════════════════════════════════════════════════
// Webhook 管理
// ════════════════════════════════════════════════════════════

/** Webhook 模板列表 */
export async function getWebhookTemplates(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`${API_BASE}/marketplace/webhooks/templates?${qs}`)
}

/** Webhook 模板详情 */
export async function getWebhookTemplateDetail(id) {
  return request(`${API_BASE}/marketplace/webhooks/templates/${id}`)
}

/** 订阅 Webhook */
export async function subscribeWebhook(data) {
  return request(`${API_BASE}/marketplace/webhooks/subscriptions`, { method: 'POST', body: JSON.stringify(data) })
}

/** 订阅列表 */
export async function getWebhookSubscriptions(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`${API_BASE}/marketplace/webhooks/subscriptions?${qs}`)
}

/** 取消订阅 */
export async function unsubscribeWebhook(id) {
  return request(`${API_BASE}/marketplace/webhooks/subscriptions/${id}`, { method: 'DELETE' })
}

/** 订阅详情 */
export async function getSubscriptionDetail(id) {
  return request(`${API_BASE}/marketplace/webhooks/subscriptions/${id}`)
}

/** 更新订阅 */
export async function updateSubscription(id, data) {
  return request(`${API_BASE}/marketplace/webhooks/subscriptions/${id}`, { method: 'PUT', body: JSON.stringify(data) })
}

/** 投递日志列表 */
export async function getDeliveryLogs(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`${API_BASE}/marketplace/webhooks/delivery-logs?${qs}`)
}

/** 投递日志详情 */
export async function getDeliveryLogDetail(id) {
  return request(`${API_BASE}/marketplace/webhooks/delivery-logs/${id}`)
}

/** 重试投递 */
export async function retryDelivery(logId) {
  return request(`${API_BASE}/marketplace/webhooks/delivery-logs/${logId}/retry`, { method: 'POST' })
}

// ════════════════════════════════════════════════════════════
// API 文档
// ════════════════════════════════════════════════════════════

/** 获取 API 文档列表 */
export async function getApiDocs(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`${API_BASE}/developer/docs?${qs}`)
}

/** 获取 API 文档详情 */
export async function getApiDocDetail(id) {
  return request(`${API_BASE}/developer/docs/${id}`)
}

/** API 测试请求 */
export async function testApiRequest(data) {
  return request(`${API_BASE}/developer/docs/test`, { method: 'POST', body: JSON.stringify(data) })
}

/** 获取 SDK 列表 */
export async function getSdkList(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`${API_BASE}/developer/sdk?${qs}`)
}
