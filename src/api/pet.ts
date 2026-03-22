import axios from 'axios'

const BASE_URL = '/api/v1'

function getToken() {
  return localStorage.getItem('token') || ''
}

function headers() {
  return { Authorization: `Bearer ${getToken()}` }
}

// ========== 宠物配置 CRUD ==========
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

// ========== 宠物设备绑定 ==========
export async function bindDevice(petId, deviceId) {
  const res = await axios.post(`${BASE_URL}/pets/${petId}/bind`, { device_id: deviceId }, { headers: headers() })
  return res.data
}

export async function unbindDevice(petId) {
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
