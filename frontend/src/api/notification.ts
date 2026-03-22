import axios from 'axios'

const BASE_URL = '/api/v1'

function getToken() {
  return localStorage.getItem('token') || ''
}

function headers() {
  return { Authorization: `Bearer ${getToken()}` }
}

// ========== 通知渠道 ==========
export interface NotificationChannel {
  id?: number
  channel_type: 'email' | 'sms' | 'webhook'
  channel_name: string
  config: ChannelConfig
  enabled: boolean
  is_default: boolean
  priority: number
  health_status?: 'healthy' | 'unhealthy' | 'unknown'
  last_checked_at?: string
  created_at?: string
  updated_at?: string
}

export interface ChannelConfig {
  // Email
  smtp_host?: string
  smtp_port?: number
  smtp_user?: string
  smtp_password?: string
  smtp_from?: string
  smtp_to?: string
  smtp_use_tls?: boolean
  smtp_use_ssl?: boolean
  // SMS
  sms_provider?: 'aliyun' | 'tencent' | 'huawei'
  access_key?: string
  secret_key?: string
  sign_name?: string
  test_phone?: string
  // Webhook
  webhook_url?: string
  webhook_secret?: string
  webhook_headers?: Record<string, string>
  retry_count?: number
  timeout_seconds?: number
  // 通知时段
  notification_periods?: NotificationPeriod[]
}

export interface NotificationPeriod {
  id?: number
  name: string
  start_time: string  // HH:mm
  end_time: string    // HH:mm
  enabled: boolean
  days_of_week?: number[]  // 0=周日, 1=周一...
}

export interface TestResult {
  success: boolean
  message?: string
  error?: string
}

export interface NotificationLog {
  id: number
  channel_id: number
  channel_type: string
  channel_name?: string
  alert_id?: number
  recipient: string
  subject?: string
  body?: string
  status: 'pending' | 'sent' | 'failed' | 'retrying'
  error_code?: string
  error_message?: string
  attempt_count: number
  sent_at?: string
  created_at: string
}

export interface AlertHistory {
  id: number
  original_id: number
  rule_id?: number
  device_id: string
  device_name?: string
  alert_type: string
  severity: number  // 1=提示, 2=警告, 3=严重, 4=紧急
  message: string
  trigger_value?: string
  threshold?: string
  status: number  // 0=未处理, 1=已确认, 2=已解决, 3=已忽略
  notified_channels?: string[]
  confirmed_at?: string
  confirmed_by?: number
  resolved_at?: string
  resolved_by?: number
  resolve_remark?: string
  created_at: string
  resolved_at_h?: string
}

export interface NotificationStats {
  total_sent: number
  total_failed: number
  success_rate: number
  by_channel: {
    channel_type: string
    sent: number
    failed: number
    success_rate: number
  }[]
  by_day: {
    date: string
    sent: number
    failed: number
  }[]
}

// ========== 通知渠道 API ==========
export async function getNotificationChannels(): Promise<{ code: number; data: NotificationChannel[] }> {
  const res = await axios.get(`${BASE_URL}/notification/channels`, { headers: headers() })
  return res.data
}

export async function getNotificationChannel(id: number): Promise<{ code: number; data: NotificationChannel }> {
  const res = await axios.get(`${BASE_URL}/notification/channels/${id}`, { headers: headers() })
  return res.data
}

export async function createNotificationChannel(data: Partial<NotificationChannel>): Promise<{ code: number; data: NotificationChannel }> {
  const res = await axios.post(`${BASE_URL}/notification/channels`, data, { headers: headers() })
  return res.data
}

export async function updateNotificationChannel(id: number, data: Partial<NotificationChannel>): Promise<{ code: number; data: NotificationChannel }> {
  const res = await axios.put(`${BASE_URL}/notification/channels/${id}`, data, { headers: headers() })
  return res.data
}

export async function deleteNotificationChannel(id: number): Promise<{ code: number }> {
  const res = await axios.delete(`${BASE_URL}/notification/channels/${id}`, { headers: headers() })
  return res.data
}

export async function testNotificationChannel(id: number, testData?: { recipient?: string }): Promise<{ code: number; data: TestResult }> {
  const res = await axios.post(`${BASE_URL}/notification/channels/${id}/test`, testData || {}, { headers: headers() })
  return res.data
}

export async function toggleNotificationChannel(id: number, enabled: boolean): Promise<{ code: number }> {
  const res = await axios.put(`${BASE_URL}/notification/channels/${id}`, { enabled }, { headers: headers() })
  return res.data
}

// ========== 通知时段 ==========
export async function getNotificationPeriods(): Promise<{ code: number; data: NotificationPeriod[] }> {
  const res = await axios.get(`${BASE_URL}/notification/periods`, { headers: headers() })
  return res.data
}

export async function updateNotificationPeriods(periods: NotificationPeriod[]): Promise<{ code: number }> {
  const res = await axios.put(`${BASE_URL}/notification/periods`, { periods }, { headers: headers() })
  return res.data
}

// ========== 通知日志 ==========
export interface NotificationLogParams {
  channel_type?: string
  status?: string
  start_time?: string
  end_time?: string
  page?: number
  page_size?: number
}

export async function getNotificationLogs(params: NotificationLogParams = {}): Promise<{ code: number; data: { list: NotificationLog[]; total: number } }> {
  const res = await axios.get(`${BASE_URL}/notification/logs`, { params, headers: headers() })
  return res.data
}

export async function getNotificationLogDetail(id: number): Promise<{ code: number; data: NotificationLog }> {
  const res = await axios.get(`${BASE_URL}/notification/logs/${id}`, { headers: headers() })
  return res.data
}

// ========== 通知统计 ==========
export async function getNotificationStats(startDate?: string, endDate?: string): Promise<{ code: number; data: NotificationStats }> {
  const res = await axios.get(`${BASE_URL}/notification/stats`, {
    params: { start_date: startDate, end_date: endDate },
    headers: headers()
  })
  return res.data
}

// ========== 告警历史 ==========
export interface AlertHistoryParams {
  device_id?: string
  alert_type?: string
  severity?: number
  status?: number
  start_time?: string
  end_time?: string
  page?: number
  page_size?: number
}

export async function getAlertHistory(params: AlertHistoryParams = {}): Promise<{ code: number; data: { list: AlertHistory[]; total: number } }> {
  const res = await axios.get(`${BASE_URL}/alerts/history`, { params, headers: headers() })
  return res.data
}

export async function confirmAlertHistory(id: number): Promise<{ code: number }> {
  const res = await axios.post(`${BASE_URL}/alerts/history/${id}/confirm`, {}, { headers: headers() })
  return res.data
}

export async function resolveAlertHistory(id: number, remark?: string): Promise<{ code: number }> {
  const res = await axios.post(`${BASE_URL}/alerts/history/${id}/resolve`, { remark }, { headers: headers() })
  return res.data
}

export async function batchExportAlertHistory(params: AlertHistoryParams = {}): Promise<Blob> {
  const res = await axios.get(`${BASE_URL}/alerts/history/export`, {
    params,
    headers: headers(),
    responseType: 'blob'
  })
  return res.data
}
