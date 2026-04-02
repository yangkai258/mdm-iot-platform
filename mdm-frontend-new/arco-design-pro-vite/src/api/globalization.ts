import axios from 'axios'

const BASE_URL = '/api'

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
  const res = await axios.get(`/api/regions`, { params, headers: headers() })
  return res.data
}

export async function getRegion(id: string) {
  const res = await axios.get(`/api/regions/${id}`, { headers: headers() })
  return res.data
}

export async function createRegion(data: any) {
  const res = await axios.post(`/api/regions`, data, { headers: headers() })
  return res.data
}

export async function updateRegion(id: string, data: any) {
  const res = await axios.put(`/api/regions/${id}`, data, { headers: headers() })
  return res.data
}

export async function deleteRegion(id: string) {
  const res = await axios.delete(`/api/regions/${id}`, { headers: headers() })
  return res.data
}

export async function getRegionNodes(regionId: string, params?: any) {
  const res = await axios.get(`/api/regions/${regionId}/nodes`, { params, headers: headers() })
  return res.data
}

export async function createRegionNode(regionId: string, data: any) {
  const res = await axios.post(`/api/regions/${regionId}/nodes`, data, { headers: headers() })
  return res.data
}

export async function deleteRegionNode(regionId: string, nodeId: string) {
  const res = await axios.delete(`/api/regions/${regionId}/nodes/${nodeId}`, { headers: headers() })
  return res.data
}

export async function getRegionHealth(regionId: string) {
  const res = await axios.get(`/api/regions/${regionId}/health`, { headers: headers() })
  return res.data
}

// ============================================================
// 时区配置
// ============================================================
export async function getTimezone() {
  const res = await axios.get(`/api/timezone`, { headers: headers() })
  return res.data
}

export async function updateTimezone(data: any) {
  const res = await axios.put(`/api/timezone`, data, { headers: headers() })
  return res.data
}

export async function getTimezoneList() {
  const res = await axios.get(`/api/timezone/list`, { headers: headers() })
  return res.data
}

// ============================================================
// 数据驻留
// ============================================================
export async function getDataResidencyRules(params?: any) {
  const res = await axios.get(`/api/data-residency/rules`, { params, headers: headers() })
  return res.data
}

export async function getDataResidencyRule(id: string) {
  const res = await axios.get(`/api/data-residency/rules/${id}`, { headers: headers() })
  return res.data
}

export async function createDataResidencyRule(data: any) {
  const res = await axios.post(`/api/data-residency/rules`, data, { headers: headers() })
  return res.data
}

export async function updateDataResidencyRule(id: string, data: any) {
  const res = await axios.put(`/api/data-residency/rules/${id}`, data, { headers: headers() })
  return res.data
}

export async function deleteDataResidencyRule(id: string) {
  const res = await axios.delete(`/api/data-residency/rules/${id}`, { headers: headers() })
  return res.data
}

// ============================================================
// 区域 AI 节点
// ============================================================
export async function getAINodes(params?: any) {
  const res = await axios.get(`/api/ai-nodes`, { params, headers: headers() })
  return res.data
}

export async function getAINode(id: string) {
  const res = await axios.get(`/api/ai-nodes/${id}`, { headers: headers() })
  return res.data
}

export async function registerAINode(data: any) {
  const res = await axios.post(`/api/ai-nodes`, data, { headers: headers() })
  return res.data
}

export async function deregisterAINode(id: string) {
  const res = await axios.delete(`/api/ai-nodes/${id}`, { headers: headers() })
  return res.data
}

// ============================================================
// 跨区域同步
// ============================================================
export async function getSyncRecords(params?: any) {
  const res = await axios.get(`/api/region-sync/records`, { params, headers: headers() })
  return res.data
}

export async function triggerSync(data: any) {
  const res = await axios.post(`/api/region-sync/trigger`, data, { headers: headers() })
  return res.data
}

export async function getSyncStatus() {
  const res = await axios.get(`/api/region-sync/status`, { headers: headers() })
  return res.data
}
