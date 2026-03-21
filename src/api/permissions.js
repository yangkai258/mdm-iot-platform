import axios from 'axios'

const BASE_URL = '/api/v1/permissions'

function getToken() {
  return localStorage.getItem('token') || ''
}

function headers() {
  return { Authorization: `Bearer ${getToken()}` }
}

// ========== 角色管理 ==========
export async function getRoles(params = {}) {
  const res = await axios.get(`${BASE_URL}/roles`, { params, headers: headers() })
  return res.data
}

export async function getRole(id) {
  const res = await axios.get(`${BASE_URL}/roles/${id}`, { headers: headers() })
  return res.data
}

export async function createRole(data) {
  const res = await axios.post(`${BASE_URL}/roles`, data, { headers: headers() })
  return res.data
}

export async function updateRole(id, data) {
  const res = await axios.put(`${BASE_URL}/roles/${id}`, data, { headers: headers() })
  return res.data
}

export async function deleteRole(id) {
  const res = await axios.delete(`${BASE_URL}/roles/${id}`, { headers: headers() })
  return res.data
}

export async function assignPermissions(roleId, permissionIds) {
  const res = await axios.post(`${BASE_URL}/roles/${roleId}/permissions`, { permission_ids: permissionIds }, { headers: headers() })
  return res.data
}

export async function getRolePermissions(roleId) {
  const res = await axios.get(`${BASE_URL}/roles/${roleId}/permissions`, { headers: headers() })
  return res.data
}

// ========== 权限组管理 ==========
export async function getPermissionGroups(params = {}) {
  const res = await axios.get(`${BASE_URL}/groups`, { params, headers: headers() })
  return res.data
}

export async function getPermissionGroup(id) {
  const res = await axios.get(`${BASE_URL}/groups/${id}`, { headers: headers() })
  return res.data
}

export async function createPermissionGroup(data) {
  const res = await axios.post(`${BASE_URL}/groups`, data, { headers: headers() })
  return res.data
}

export async function updatePermissionGroup(id, data) {
  const res = await axios.put(`${BASE_URL}/groups/${id}`, data, { headers: headers() })
  return res.data
}

export async function deletePermissionGroup(id) {
  const res = await axios.delete(`${BASE_URL}/groups/${id}`, { headers: headers() })
  return res.data
}

// ========== 菜单管理 ==========
export async function getMenus(params = {}) {
  const res = await axios.get(`${BASE_URL}/menus`, { params, headers: headers() })
  return res.data
}

export async function getMenuTree() {
  const res = await axios.get(`${BASE_URL}/menus/tree`, { headers: headers() })
  return res.data
}

export async function createMenu(data) {
  const res = await axios.post(`${BASE_URL}/menus`, data, { headers: headers() })
  return res.data
}

export async function updateMenu(id, data) {
  const res = await axios.put(`${BASE_URL}/menus/${id}`, data, { headers: headers() })
  return res.data
}

export async function deleteMenu(id) {
  const res = await axios.delete(`${BASE_URL}/menus/${id}`, { headers: headers() })
  return res.data
}

// ========== API 权限 ==========
export async function getApiPermissions(params = {}) {
  const res = await axios.get(`${BASE_URL}/api-permissions`, { params, headers: headers() })
  return res.data
}

export async function createApiPermission(data) {
  const res = await axios.post(`${BASE_URL}/api-permissions`, data, { headers: headers() })
  return res.data
}

export async function updateApiPermission(id, data) {
  const res = await axios.put(`${BASE_URL}/api-permissions/${id}`, data, { headers: headers() })
  return res.data
}

export async function deleteApiPermission(id) {
  const res = await axios.delete(`${BASE_URL}/api-permissions/${id}`, { headers: headers() })
  return res.data
}

// ========== 数据权限规则 ==========
export async function getDataPermissionRules(params = {}) {
  const res = await axios.get(`${BASE_URL}/data-rules`, { params, headers: headers() })
  return res.data
}

export async function createDataPermissionRule(data) {
  const res = await axios.post(`${BASE_URL}/data-rules`, data, { headers: headers() })
  return res.data
}

export async function updateDataPermissionRule(id, data) {
  const res = await axios.put(`${BASE_URL}/data-rules/${id}`, data, { headers: headers() })
  return res.data
}

export async function deleteDataPermissionRule(id) {
  const res = await axios.delete(`${BASE_URL}/data-rules/${id}`, { headers: headers() })
  return res.data
}

// ========== 获取所有可用权限（树形） ==========
export async function getAllPermissions() {
  const res = await axios.get(`${BASE_URL}/all`, { headers: headers() })
  return res.data
}
