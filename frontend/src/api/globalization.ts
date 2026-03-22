import axios from 'axios'

const BASE_URL = '/api/v1'

function getToken() {
  return localStorage.getItem('token') || ''
}

function headers() {
  return { Authorization: `Bearer ${getToken()}` }
}

// ============================================================
// 区域管理
// ============================================================
export async function getRegions(params?: any) {
  const res = await axios.get(`${BASE_URL}/regions`, { params, headers: headers() })
  return res.data
}

export async function getRegion(id: string) {
  const res = await axios.get(`${BASE_URL}/regions/${id}`, { headers: headers() })
  return res.data
}

export async function createRegion(data: any) {
  const res = await axios.post(`${BASE_URL}/regions`, data, { headers: headers() })
  return res.data
}

export async function updateRegion(id: string, data: any) {
  const res = await axios.put(`${BASE_URL}/regions/${id}`, data, { headers: headers() })
  return res.data
}

export async function deleteRegion(id: string) {
  const res = await axios.delete(`${BASE_URL}/regions/${id}`, { headers: headers() })
  return res.data
}

export async function getRegionNodes(regionId: string, params?: any) {
  const res = await axios.get(`${BASE_URL}/regions/${regionId}/nodes`, { params, headers: headers() })
  return res.data
}

export async function createRegionNode(regionId: string, data: any) {
  const res = await axios.post(`${BASE_URL}/regions/${regionId}/nodes`, data, { headers: headers() })
  return res.data
}

export async function deleteRegionNode(regionId: string, nodeId: string) {
  const res = await axios.delete(`${BASE_URL}/regions/${regionId}/nodes/${nodeId}`, { headers: headers() })
  return res.data
}

export async function getRegionHealth(regionId: string) {
  const res = await axios.get(`${BASE_URL}/regions/${regionId}/health`, { headers: headers() })
  return res.data
}

// ============================================================
// 时区配置
// ============================================================
export async function getTimezone() {
  const res = await axios.get(`${BASE_URL}/timezone`, { headers: headers() })
  return res.data
}

export async function updateTimezone(data: any) {
  const res = await axios.put(`${BASE_URL}/timezone`, data, { headers: headers() })
  return res.data
}

export async function getTimezoneList() {
  const res = await axios.get(`${BASE_URL}/timezone/list`, { headers: headers() })
  return res.data
}

// ============================================================
// 数据驻留
// ============================================================
export async function getDataResidencyRules(params?: any) {
  const res = await axios.get(`${BASE_URL}/data-residency/rules`, { params, headers: headers() })
  return res.data
}

export async function getDataResidencyRule(id: string) {
  const res = await axios.get(`${BASE_URL}/data-residency/rules/${id}`, { headers: headers() })
  return res.data
}

export async function createDataResidencyRule(data: any) {
  const res = await axios.post(`${BASE_URL}/data-residency/rules`, data, { headers: headers() })
  return res.data
}

export async function updateDataResidencyRule(id: string, data: any) {
  const res = await axios.put(`${BASE_URL}/data-residency/rules/${id}`, data, { headers: headers() })
  return res.data
}

export async function deleteDataResidencyRule(id: string) {
  const res = await axios.delete(`${BASE_URL}/data-residency/rules/${id}`, { headers: headers() })
  return res.data
}

// ============================================================
// 区域 AI 节点
// ============================================================
export async function getAINodes(params?: any) {
  const res = await axios.get(`${BASE_URL}/ai-nodes`, { params, headers: headers() })
  return res.data
}

export async function getAINode(id: string) {
  const res = await axios.get(`${BASE_URL}/ai-nodes/${id}`, { headers: headers() })
  return res.data
}

export async function registerAINode(data: any) {
  const res = await axios.post(`${BASE_URL}/ai-nodes`, data, { headers: headers() })
  return res.data
}

export async function deregisterAINode(id: string) {
  const res = await axios.delete(`${BASE_URL}/ai-nodes/${id}`, { headers: headers() })
  return res.data
}

// ============================================================
// 跨区域同步
// ============================================================
export async function getSyncRecords(params?: any) {
  const res = await axios.get(`${BASE_URL}/region-sync/records`, { params, headers: headers() })
  return res.data
}

export async function triggerSync(data: any) {
  const res = await axios.post(`${BASE_URL}/region-sync/trigger`, data, { headers: headers() })
  return res.data
}

export async function getSyncStatus() {
  const res = await axios.get(`${BASE_URL}/region-sync/status`, { headers: headers() })
  return res.data
}
