import axios from 'axios'

const BASE_URL = '/api/audit'

function getToken() {
  return localStorage.getItem('token') || ''
}

function headers() {
  return { Authorization: `Bearer ${getToken()}` }
}

// ============================================================
// 审计日志
// ============================================================
export async function getAuditLogs(params?: any) {
  const res = await axios.get(`/api/logs`, { params, headers: headers() })
  return res.data
}

export async function getAuditLogDetail(id: string) {
  const res = await axios.get(`/api/logs/${id}`, { headers: headers() })
  return res.data
}

export async function exportAuditLogs(params?: any) {
  const res = await axios.get(`/api/logs/export`, {
    params,
    headers: headers(),
    responseType: 'blob'
  })
  return res.data
}

export async function getAuditStatistics(params?: any) {
  const res = await axios.get(`/api/statistics`, { params, headers: headers() })
  return res.data
}
