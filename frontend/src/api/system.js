import axios from 'axios'

const BASE_URL = '/api/v1'

function getToken() {
  return localStorage.getItem('token') || ''
}

function headers() {
  return { Authorization: `Bearer ${getToken()}` }
}

// ========== 系统参数 ==========
export async function getSystemSettings() {
  const res = await axios.get(`${BASE_URL}/system/settings`, { headers: headers() })
  return res.data
}

export async function updateSystemSettings(data) {
  const res = await axios.put(`${BASE_URL}/system/settings`, data, { headers: headers() })
  return res.data
}

export async function getSystemParamsByCategory(category) {
  const res = await axios.get(`${BASE_URL}/system/params/${category}`, { headers: headers() })
  return res.data
}

export async function updateSystemParam(id, data) {
  const res = await axios.put(`${BASE_URL}/system/params/${id}`, data, { headers: headers() })
  return res.data
}

// ========== 邮件模板 ==========
export async function getEmailTemplates(params = {}) {
  const res = await axios.get(`${BASE_URL}/email-templates`, { params, headers: headers() })
  return res.data
}

export async function getEmailTemplate(id) {
  const res = await axios.get(`${BASE_URL}/email-templates/${id}`, { headers: headers() })
  return res.data
}

export async function createEmailTemplate(data) {
  const res = await axios.post(`${BASE_URL}/email-templates`, data, { headers: headers() })
  return res.data
}

export async function updateEmailTemplate(id, data) {
  const res = await axios.put(`${BASE_URL}/email-templates/${id}`, data, { headers: headers() })
  return res.data
}

export async function deleteEmailTemplate(id) {
  const res = await axios.delete(`${BASE_URL}/email-templates/${id}`, { headers: headers() })
  return res.data
}

export async function testEmailTemplate(id, params) {
  const res = await axios.post(`${BASE_URL}/email-templates/${id}/test`, params, { headers: headers() })
  return res.data
}

export async function getEmailTemplateVariables() {
  const res = await axios.get(`${BASE_URL}/email-templates/variables`, { headers: headers() })
  return res.data
}
