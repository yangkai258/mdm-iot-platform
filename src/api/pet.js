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

// ========== 宠物对话 ==========
export async function getConversations(petId, params = {}) {
  const res = await axios.get(`${BASE_URL}/pets/${petId}/conversations`, { params, headers: headers() })
  return res.data
}

export async function sendCommand(petId, data) {
  const res = await axios.post(`${BASE_URL}/pets/${petId}/conversations`, data, { headers: headers() })
  return res.data
}

export async function deleteConversation(petId, conversationId) {
  const res = await axios.delete(`${BASE_URL}/pets/${petId}/conversations/${conversationId}`, { headers: headers() })
  return res.data
}
