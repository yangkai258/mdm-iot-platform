import axios from 'axios'

const BASE_URL = '/api/v1/org'

function getToken() {
  return localStorage.getItem('token') || ''
}

function headers() {
  return { Authorization: `Bearer ${getToken()}` }
}

// ========== 公司管理 ==========
export async function getCompanies(params = {}) {
  const res = await axios.get(`${BASE_URL}/companies`, { params, headers: headers() })
  return res.data
}

export async function getCompany(id) {
  const res = await axios.get(`${BASE_URL}/companies/${id}`, { headers: headers() })
  return res.data
}

export async function createCompany(data) {
  const res = await axios.post(`${BASE_URL}/companies`, data, { headers: headers() })
  return res.data
}

export async function updateCompany(id, data) {
  const res = await axios.put(`${BASE_URL}/companies/${id}`, data, { headers: headers() })
  return res.data
}

export async function deleteCompany(id) {
  const res = await axios.delete(`${BASE_URL}/companies/${id}`, { headers: headers() })
  return res.data
}

// ========== 部门管理 ==========
export async function getDepartments(params = {}) {
  const res = await axios.get(`${BASE_URL}/departments`, { params, headers: headers() })
  return res.data
}

export async function getDepartmentTree(companyId) {
  const res = await axios.get(`${BASE_URL}/departments/tree`, { params: { company_id: companyId }, headers: headers() })
  return res.data
}

export async function createDepartment(data) {
  const res = await axios.post(`${BASE_URL}/departments`, data, { headers: headers() })
  return res.data
}

export async function updateDepartment(id, data) {
  const res = await axios.put(`${BASE_URL}/departments/${id}`, data, { headers: headers() })
  return res.data
}

export async function deleteDepartment(id) {
  const res = await axios.delete(`${BASE_URL}/departments/${id}`, { headers: headers() })
  return res.data
}

// ========== 岗位管理 ==========
export async function getPositions(params = {}) {
  const res = await axios.get(`${BASE_URL}/positions`, { params, headers: headers() })
  return res.data
}

export async function createPosition(data) {
  const res = await axios.post(`${BASE_URL}/positions`, data, { headers: headers() })
  return res.data
}

export async function updatePosition(id, data) {
  const res = await axios.put(`${BASE_URL}/positions/${id}`, data, { headers: headers() })
  return res.data
}

export async function deletePosition(id) {
  const res = await axios.delete(`${BASE_URL}/positions/${id}`, { headers: headers() })
  return res.data
}

// ========== 员工管理 ==========
export async function getEmployees(params = {}) {
  const res = await axios.get(`${BASE_URL}/employees`, { params, headers: headers() })
  return res.data
}

export async function getEmployee(id) {
  const res = await axios.get(`${BASE_URL}/employees/${id}`, { headers: headers() })
  return res.data
}

export async function createEmployee(data) {
  const res = await axios.post(`${BASE_URL}/employees`, data, { headers: headers() })
  return res.data
}

export async function updateEmployee(id, data) {
  const res = await axios.put(`${BASE_URL}/employees/${id}`, data, { headers: headers() })
  return res.data
}

export async function deleteEmployee(id) {
  const res = await axios.delete(`${BASE_URL}/employees/${id}`, { headers: headers() })
  return res.data
}

// ========== 基准岗位 ==========
export async function getStandardPositions(params = {}) {
  const res = await axios.get(`${BASE_URL}/standard-positions`, { params, headers: headers() })
  return res.data
}

export async function createStandardPosition(data) {
  const res = await axios.post(`${BASE_URL}/standard-positions`, data, { headers: headers() })
  return res.data
}

export async function updateStandardPosition(id, data) {
  const res = await axios.put(`${BASE_URL}/standard-positions/${id}`, data, { headers: headers() })
  return res.data
}

export async function deleteStandardPosition(id) {
  const res = await axios.delete(`${BASE_URL}/standard-positions/${id}`, { headers: headers() })
  return res.data
}
