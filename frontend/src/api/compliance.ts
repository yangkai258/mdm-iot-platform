import axios from 'axios'

const BASE_URL = '/api/v1/compliance'

function getToken() {
  return localStorage.getItem('token') || ''
}

function headers() {
  return { Authorization: `Bearer ${getToken()}` }
}

// ============================================================
// 数据脱敏配置
// ============================================================
export async function getDesensitizationConfig() {
  const res = await axios.get(`${BASE_URL}/desensitization`, { headers: headers() })
  return res.data
}

export async function updateDesensitizationConfig(data: any) {
  const res = await axios.put(`${BASE_URL}/desensitization`, data, { headers: headers() })
  return res.data
}

// ============================================================
// GDPR 合规
// ============================================================
export async function getGdprConfig() {
  const res = await axios.get(`${BASE_URL}/gdpr`, { headers: headers() })
  return res.data
}

export async function updateGdprConfig(data: any) {
  const res = await axios.put(`${BASE_URL}/gdpr`, data, { headers: headers() })
  return res.data
}

// ============================================================
// 数据导出
// ============================================================
export async function exportData(params: any) {
  const res = await axios.post(`${BASE_URL}/export`, params, { headers: headers() })
  return res.data
}

export async function getExportHistory(params?: any) {
  const res = await axios.get(`${BASE_URL}/export/history`, { params, headers: headers() })
  return res.data
}

export async function deleteExportRecord(id: string) {
  const res = await axios.delete(`${BASE_URL}/export/${id}`, { headers: headers() })
  return res.data
}

export async function getExportFile(id: string) {
  const res = await axios.get(`${BASE_URL}/export/${id}/file`, {
    headers: headers(),
    responseType: 'blob'
  })
  return res.data
}

// ============================================================
// GDPR 请求
// ============================================================
export async function submitGdprRequest(data: any) {
  const res = await axios.post(`${BASE_URL}/gdpr/requests`, data, { headers: headers() })
  return res.data
}

export async function getGdprRequests(params?: any) {
  const res = await axios.get(`${BASE_URL}/gdpr/requests`, { params, headers: headers() })
  return res.data
}

export async function processGdprRequest(id: string, data: any) {
  const res = await axios.post(`${BASE_URL}/gdpr/requests/${id}/process`, data, { headers: headers() })
  return res.data
}
