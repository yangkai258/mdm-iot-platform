/**
 * 订阅和计费模块 API
 * 订阅计划、用户订阅、用量查询、Webhook、账单发票
 */

const BASE_URL = '/api/v1'

const getToken = () => localStorage.getItem('token') || ''

const request = (path: string, options: RequestInit = {}) => {
  return fetch(`${BASE_URL}${path}`, {
    ...options,
    headers: {
      'Authorization': `Bearer ${getToken()}`,
      'Content-Type': 'application/json',
      ...(options.headers || {})
    }
  }).then(r => r.json())
}

// ============ 订阅计划 ============

export interface SubscriptionPlan {
  id: number
  plan_code: string
  plan_name: string
  plan_type: 'free' | 'paid'
  price_monthly: number
  price_yearly: number
  features: string[]
  quotas: {
    devices: number
    api_calls: number
    storage_gb: number
  }
  is_active: boolean
  is_recommended: boolean
  sort_order: number
}

export interface UserSubscription {
  id: number
  user_id: number
  plan_id: number
  plan_name: string
  status: 'active' | 'cancelled' | 'expired' | 'trial'
  started_at: string
  expires_at: string
  next_billing_at: string
  auto_renew: boolean
  cancelled_at?: string
  cancel_reason?: string
}

export const billingApi = {
  // 订阅计划
  getSubscriptionPlans: (params?: { plan_type?: string }) => {
    const qs = params ? '?' + new URLSearchParams(params as any).toString() : ''
    return request(`/subscription/plans${qs}`)
  },
  getSubscriptionPlanById: (id: number) => request(`/subscription/plans/${id}`),

  // 用户订阅
  getCurrentSubscription: () => request('/subscriptions'),
  createSubscription: (data: { plan_id: number; billing_cycle: 'monthly' | 'yearly' }) =>
    request('/subscriptions', { method: 'POST', body: JSON.stringify(data) }),
  getSubscriptionById: (id: number) => request(`/subscriptions/${id}`),
  cancelSubscription: (id: number, data?: { reason?: string }) =>
    request(`/subscriptions/${id}/cancel`, { method: 'POST', body: JSON.stringify(data) }),
  renewSubscription: (id: number) =>
    request(`/subscriptions/${id}/renew`, { method: 'POST' }),
  upgradeSubscription: (id: number, data: { new_plan_id: number }) =>
    request(`/subscriptions/${id}/upgrade`, { method: 'POST', body: JSON.stringify(data) }),
  downgradeSubscription: (id: number, data: { new_plan_id: number }) =>
    request(`/subscriptions/${id}/downgrade`, { method: 'POST', body: JSON.stringify(data) }),

  // 用量查询
  getCurrentUsage: () => request('/usage/current'),
  getUsageHistory: (params?: { page?: number; page_size?: number; start_time?: string; end_time?: string }) => {
    const qs = params ? '?' + new URLSearchParams(params as any).toString() : ''
    return request(`/usage/history${qs}`)
  },
  getUsageQuotas: () => request('/usage/quotas'),
  getUsageStats: (params?: { period?: string }) => {
    const qs = params ? '?' + new URLSearchParams(params as any).toString() : ''
    return request(`/usage/stats${qs}`)
  },

  // Webhook
  getWebhooks: (params?: { page?: number; page_size?: number }) => {
    const qs = params ? '?' + new URLSearchParams(params as any).toString() : ''
    return request(`/webhooks${qs}`)
  },
  createWebhook: (data: { webhook_name: string; endpoint_url: string; events: string[]; secret_key?: string }) =>
    request('/webhooks', { method: 'POST', body: JSON.stringify(data) }),
  getWebhookById: (id: number) => request(`/webhooks/${id}`),
  updateWebhook: (id: number, data: Partial<{ webhook_name: string; endpoint_url: string; events: string[]; secret_key: string; is_active: boolean }>) =>
    request(`/webhooks/${id}`, { method: 'PUT', body: JSON.stringify(data) }),
  deleteWebhook: (id: number) => request(`/webhooks/${id}`, { method: 'DELETE' }),
  testWebhook: (id: number) => request(`/webhooks/${id}/test`, { method: 'POST' }),
  getWebhookLogs: (id: number, params?: { page?: number; page_size?: number }) => {
    const qs = params ? '?' + new URLSearchParams(params as any).toString() : ''
    return request(`/webhooks/${id}/logs${qs}`)
  },

  // 账单
  getBillingRecords: (params?: { page?: number; page_size?: number; status?: string }) => {
    const qs = params ? '?' + new URLSearchParams(params as any).toString() : ''
    return request(`/billing/records${qs}`)
  },
  getBillingRecordById: (id: number) => request(`/billing/records/${id}`),
  getBillingSummary: () => request('/billing/summary'),

  // 发票
  getInvoices: (params?: { page?: number; page_size?: number; status?: string }) => {
    const qs = params ? '?' + new URLSearchParams(params as any).toString() : ''
    return request(`/billing/invoices${qs}`)
  },
  createInvoice: (data: { title: string; tax_number?: string; amount: number; invoice_type?: string }) =>
    request('/billing/invoices', { method: 'POST', body: JSON.stringify(data) }),
  getInvoiceById: (id: number) => request(`/billing/invoices/${id}`),
}

export default billingApi
