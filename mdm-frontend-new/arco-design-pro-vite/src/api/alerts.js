import axios from 'axios'

const BASE_URL = '/api'

function getToken() {
  return localStorage.getItem('token') || ''
}

function headers() {
  return { Authorization: `Bearer ${getToken()}` }
}

// ========== 告警规则 ==========
export async function getAlertRules() {
  const res = await axios.get(`/api/alert-rules`, { headers: headers() })
  return res.data
}

export async function createAlertRule(data) {
  const res = await axios.post(`/api/alert-rules`, data, { headers: headers() })
  return res.data
}

export async function updateAlertRule(id, data) {
  const res = await axios.put(`/api/alert-rules/${id}`, data, { headers: headers() })
  return res.data
}

export async function deleteAlertRule(id) {
  const res = await axios.delete(`/api/alert-rules/${id}`, { headers: headers() })
  return res.data
}

export async function batchDeleteAlertRules(ids) {
  const res = await axios.post(`/api/alert-rules/batch-delete`, { ids }, { headers: headers() })
  return res.data
}

// ========== 告警记录 ==========
export async function getAlerts(params = {}) {
  const res = await axios.get(`/api/alerts`, { params, headers: headers() })
  return res.data
}

export async function confirmAlert(id) {
  const res = await axios.post(`/api/alerts/${id}/confirm`, {}, { headers: headers() })
  return res.data
}

export async function resolveAlert(id) {
  const res = await axios.post(`/api/alerts/${id}/resolve`, {}, { headers: headers() })
  return res.data
}

export async function ignoreAlert(id) {
  const res = await axios.post(`/api/alerts/${id}/ignore`, {}, { headers: headers() })
  return res.data
}

export async function batchConfirmAlerts(ids) {
  const res = await axios.post(`/api/alerts/batch/confirm`, { alert_ids: ids }, { headers: headers() })
  return res.data
}

export async function batchResolveAlerts(ids) {
  const res = await axios.post(`/api/alerts/batch/resolve`, { alert_ids: ids }, { headers: headers() })
  return res.data
}

export async function getAlertNotifications(id) {
  const res = await axios.get(`/api/alerts/${id}/notifications`, { headers: headers() })
  return res.data
}

// ========== 地理围栏 ==========
export async function getGeofenceRules() {
  const res = await axios.get(`/api/geofence/rules`, { headers: headers() })
  return res.data
}

export async function createGeofenceRule(data) {
  const res = await axios.post(`/api/geofence/rules`, data, { headers: headers() })
  return res.data
}

export async function updateGeofenceRule(id, data) {
  const res = await axios.put(`/api/geofence/rules/${id}`, data, { headers: headers() })
  return res.data
}

export async function deleteGeofenceRule(id) {
  const res = await axios.delete(`/api/geofence/rules/${id}`, { headers: headers() })
  return res.data
}

export async function getGeofenceAlerts(params = {}) {
  const res = await axios.get(`/api/geofence/alerts`, { params, headers: headers() })
  return res.data
}

// ========== 告警设置 ==========
export async function getAlertSettings() {
  const res = await axios.get(`/api/alerts/settings`, { headers: headers() })
  return res.data
}

export async function updateAlertSettings(data) {
  const res = await axios.put(`/api/alerts/settings`, data, { headers: headers() })
  return res.data
}

// ========== 通知渠道 ==========
export async function getNotificationChannels(params = {}) {
  const res = await axios.get(`/api/notification-channels`, { params, headers: headers() })
  return res.data
}

export async function createNotificationChannel(data) {
  const res = await axios.post(`/api/notification-channels`, data, { headers: headers() })
  return res.data
}

export async function updateNotificationChannel(id, data) {
  const res = await axios.put(`/api/notification-channels/${id}`, data, { headers: headers() })
  return res.data
}

export async function deleteNotificationChannel(id) {
  const res = await axios.delete(`/api/notification-channels/${id}`, { headers: headers() })
  return res.data
}

export async function toggleNotificationChannel(id) {
  const res = await axios.post(`/api/notification-channels/${id}/toggle`, {}, { headers: headers() })
  return res.data
}

export async function testNotificationChannel(id) {
  const res = await axios.post(`/api/notification-channels/${id}/test`, {}, { headers: headers() })
  return res.data
}
