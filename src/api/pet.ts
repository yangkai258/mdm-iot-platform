import axios from 'axios'

const BASE_URL = '/api/v1'

function getToken() {
  return localStorage.getItem('token') || ''
}

function headers() {
  return { Authorization: `Bearer ${getToken()}` }
}

// ========== 宠物档案 CRUD (Sprint 15) ==========
export async function getPets(params = {}) {
  const res = await axios.get(`${BASE_URL}/pets`, { params, headers: headers() })
  return res.data
}

export async function getPet(id) {
  const res = await axios.get(`${BASE_URL}/pets/${id}`, { headers: headers() })
  return res.data
}

export async function createPet(data) {
  const res = await axios.post(`${BASE_URL}/pets`, data, { headers: headers() })
  return res.data
}

export async function updatePet(id, data) {
  const res = await axios.put(`${BASE_URL}/pets/${id}`, data, { headers: headers() })
  return res.data
}

export async function deletePet(id) {
  const res = await axios.delete(`${BASE_URL}/pets/${id}`, { headers: headers() })
  return res.data
}

// ========== 宠物设备绑定 (Sprint 15) ==========
export async function getPetDevices(petId) {
  const res = await axios.get(`${BASE_URL}/pets/${petId}/devices`, { headers: headers() })
  return res.data
}

export async function bindDevice(petId, deviceId) {
  const res = await axios.post(`${BASE_URL}/pets/${petId}/devices`, { device_id: deviceId }, { headers: headers() })
  return res.data
}

export async function unbindDevice(petId, deviceId) {
  const res = await axios.delete(`${BASE_URL}/pets/${petId}/devices/${deviceId}`, { headers: headers() })
  return res.data
}

// ========== 寻回网络 (Sprint 15) ==========
export async function getLostFoundReports(params = {}) {
  const res = await axios.get(`${BASE_URL}/lost-found/reports`, { params, headers: headers() })
  return res.data
}

export async function createLostFoundReport(data) {
  const res = await axios.post(`${BASE_URL}/lost-found/reports`, data, { headers: headers() })
  return res.data
}

export async function getLostFoundReport(id) {
  const res = await axios.get(`${BASE_URL}/lost-found/reports/${id}`, { headers: headers() })
  return res.data
}

export async function updateLostFoundReport(id, data) {
  const res = await axios.put(`${BASE_URL}/lost-found/reports/${id}`, data, { headers: headers() })
  return res.data
}

export async function deleteLostFoundReport(id) {
  const res = await axios.delete(`${BASE_URL}/lost-found/reports/${id}`, { headers: headers() })
  return res.data
}

export async function reportSighting(reportId, data) {
  const res = await axios.post(`${BASE_URL}/lost-found/reports/${reportId}/sighting`, data, { headers: headers() })
  return res.data
}

export async function getNearbyLostFound(params = {}) {
  const res = await axios.get(`${BASE_URL}/lost-found/nearby`, { params, headers: headers() })
  return res.data
}

export async function resolveLostFoundReport(id, data = {}) {
  const res = await axios.post(`${BASE_URL}/lost-found/reports/${id}/resolve`, data, { headers: headers() })
  return res.data
}

// ========== 家庭宠物管理 (Sprint 15) ==========
export async function getHouseholdPets() {
  const res = await axios.get(`${BASE_URL}/household/pets`, { headers: headers() })
  return res.data
}

export async function getHouseholdMembers() {
  const res = await axios.get(`${BASE_URL}/household/members`, { headers: headers() })
  return res.data
}

export async function inviteHouseholdMember(data) {
  const res = await axios.post(`${BASE_URL}/household/members`, data, { headers: headers() })
  return res.data
}

export async function removeHouseholdMember(id) {
  const res = await axios.delete(`${BASE_URL}/household/members/${id}`, { headers: headers() })
  return res.data
}

// ========== 宠物设备绑定 (旧版兼容) ==========
export async function bindDeviceLegacy(petId, deviceId) {
  const res = await axios.post(`${BASE_URL}/pets/${petId}/bind`, { device_id: deviceId }, { headers: headers() })
  return res.data
}

export async function unbindDeviceLegacy(petId) {
  const res = await axios.post(`${BASE_URL}/pets/${petId}/unbind`, {}, { headers: headers() })
  return res.data
}

// ========== 宠物状态 (Sprint 9) ==========
export async function getPetStatus(deviceId, params = {}) {
  const res = await axios.get(`${BASE_URL}/pets/${deviceId}/status`, { params, headers: headers() })
  return res.data
}

// ========== 宠物对话 (Sprint 9) ==========
export async function sendMessage(deviceId, data) {
  const res = await axios.post(`${BASE_URL}/pets/${deviceId}/messages`, data, { headers: headers() })
  return res.data
}

export async function getConversations(deviceId, params = {}) {
  const res = await axios.get(`${BASE_URL}/conversations`, { params: { ...params, device_id: deviceId }, headers: headers() })
  return res.data
}

export async function createConversation(deviceId, data = {}) {
  const res = await axios.post(`${BASE_URL}/conversations`, { device_id: deviceId, ...data }, { headers: headers() })
  return res.data
}

export async function getConversationMessages(conversationId, params = {}) {
  const res = await axios.get(`${BASE_URL}/conversations/${conversationId}/messages`, { params, headers: headers() })
  return res.data
}

export async function deleteConversation(deviceId, conversationId) {
  const res = await axios.delete(`${BASE_URL}/conversations/${conversationId}`, {
    params: { device_id: deviceId },
    headers: headers()
  })
  return res.data
}

// ========== 宠物动作 (Sprint 9) ==========
export async function executeAction(deviceId, data) {
  const res = await axios.post(`${BASE_URL}/pets/${deviceId}/actions`, data, { headers: headers() })
  return res.data
}

// ========== 宠物设置 (Sprint 9) ==========
export async function updatePetSettings(deviceId, data) {
  const res = await axios.put(`${BASE_URL}/pets/${deviceId}/settings`, data, { headers: headers() })
  return res.data
}
