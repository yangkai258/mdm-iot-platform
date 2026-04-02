import axios from 'axios'

const BASE_URL = '/api'

function getToken() {
  return localStorage.getItem('token') || ''
}

function headers() {
  return { Authorization: `Bearer ${getToken()}` }
}

// ============================================================
// LDAP 配置
// ============================================================
export async function getLdapConfig() {
  const res = await axios.get(`/api/ldap/config`, { headers: headers() })
  return res.data
}

export async function updateLdapConfig(data: any) {
  const res = await axios.put(`/api/ldap/config`, data, { headers: headers() })
  return res.data
}

export async function testLdapConnection(data: any) {
  const res = await axios.post(`/api/ldap/test`, data, { headers: headers() })
  return res.data
}

export async function getLdapUsers(params: any) {
  const res = await axios.get(`/api/ldap/users`, { params, headers: headers() })
  return res.data
}

export async function getLdapGroups() {
  const res = await axios.get(`/api/ldap/groups`, { headers: headers() })
  return res.data
}

export async function syncLdapUsers() {
  const res = await axios.post(`/api/ldap/sync`, {}, { headers: headers() })
  return res.data
}

export async function setLdapGroupMapping(data: any) {
  const res = await axios.post(`/api/ldap/group-role-mapping`, data, { headers: headers() })
  return res.data
}

export async function getLdapGroupMappings() {
  const res = await axios.get(`/api/ldap/group-role-mappings`, { headers: headers() })
  return res.data
}

// ============================================================
// 证书管理
// ============================================================
export async function getCertificates(params: any) {
  const res = await axios.get(`/api/certificates`, { params, headers: headers() })
  return res.data
}

export async function getCertificate(id: string) {
  const res = await axios.get(`/api/certificates/${id}`, { headers: headers() })
  return res.data
}

export async function uploadCertificate(data: any) {
  const res = await axios.post(`/api/certificates`, data, {
    headers: { ...headers(), 'Content-Type': 'multipart/form-data' }
  })
  return res.data
}

export async function deleteCertificate(id: string) {
  const res = await axios.delete(`/api/certificates/${id}`, { headers: headers() })
  return res.data
}

export async function revokeCertificate(id: string, reason?: string) {
  const res = await axios.post(`/api/certificates/${id}/revoke`, { reason }, { headers: headers() })
  return res.data
}

export async function renewCertificate(id: string) {
  const res = await axios.post(`/api/certificates/${id}/renew`, {}, { headers: headers() })
  return res.data
}

export async function issueCertificate(data: any) {
  const res = await axios.post(`/api/certificates/issue`, data, { headers: headers() })
  return res.data
}

export async function downloadCertificate(id: string) {
  const res = await axios.get(`/api/certificates/${id}/download`, {
    headers: headers(),
    responseType: 'blob'
  })
  return res.data
}

// ============================================================
// 设备安全操作
// ============================================================
export async function lockDevice(deviceId: string, reason?: string) {
  const res = await axios.post(`/api/devices/${deviceId}/security/lock`, { reason }, { headers: headers() })
  return res.data
}

export async function unlockDevice(deviceId: string, unlockCode?: string) {
  const res = await axios.post(`/api/devices/${deviceId}/security/unlock`, { unlock_code: unlockCode }, { headers: headers() })
  return res.data
}

export async function wipeDevice(deviceId: string, data: any) {
  const res = await axios.post(`/api/devices/${deviceId}/security/wipe`, data, { headers: headers() })
  return res.data
}

export async function getDeviceSecurityStatus(deviceId: string) {
  const res = await axios.get(`/api/devices/${deviceId}/security/status`, { headers: headers() })
  return res.data
}

export async function getWipeHistory(deviceId: string, params?: any) {
  const res = await axios.get(`/api/devices/${deviceId}/wipe-history`, { params, headers: headers() })
  return res.data
}

// ============================================================
// 数据权限
// ============================================================
export async function getDataPermissionRoles(roleId: string) {
  const res = await axios.get(`/api/data-permissions/roles/${roleId}`, { headers: headers() })
  return res.data
}

export async function updateDataPermissionRoles(roleId: string, data: any) {
  const res = await axios.put(`/api/data-permissions/roles/${roleId}`, data, { headers: headers() })
  return res.data
}

export async function getDataPermissionUsers(userId: string) {
  const res = await axios.get(`/api/data-permissions/users/${userId}`, { headers: headers() })
  return res.data
}

export async function updateDataPermissionUsers(userId: string, data: any) {
  const res = await axios.put(`/api/data-permissions/users/${userId}`, data, { headers: headers() })
  return res.data
}

export async function getDataPermissions(params?: any) {
  const res = await axios.get(`/api/data-permissions`, { params, headers: headers() })
  return res.data
}

export async function createDataPermission(data: any) {
  const res = await axios.post(`/api/data-permissions`, data, { headers: headers() })
  return res.data
}

export async function updateDataPermission(id: string, data: any) {
  const res = await axios.put(`/api/data-permissions/${id}`, data, { headers: headers() })
  return res.data
}

export async function deleteDataPermission(id: string) {
  const res = await axios.delete(`/api/data-permissions/${id}`, { headers: headers() })
  return res.data
}

// ============================================================
// 辅助：获取设备列表
// ============================================================
export async function getDevices(params?: any) {
  const res = await axios.get(`/api/devices`, { params, headers: headers() })
  return res.data
}

// ============================================================
// 辅助：获取角色列表
// ============================================================
export async function getRoles(params?: any) {
  const res = await axios.get(`/api/roles`, { params, headers: headers() })
  return res.data
}

// ============================================================
// 辅助：获取用户列表
// ============================================================
export async function getUsers(params?: any) {
  const res = await axios.get(`/api/users`, { params, headers: headers() })
  return res.data
}

// ============================================================
// 加密配置
// ============================================================
export async function getEncryptionConfig() {
  const res = await axios.get(`/api/security/encryption`, { headers: headers() })
  return res.data
}

export async function updateEncryptionConfig(data: any) {
  const res = await axios.put(`/api/security/encryption`, data, { headers: headers() })
  return res.data
}

// ============================================================
// 密钥管理
// ============================================================
export async function getKeys(params?: any) {
  const res = await axios.get(`/api/security/keys`, { params, headers: headers() })
  return res.data
}

export async function createKey(data: any) {
  const res = await axios.post(`/api/security/keys`, data, { headers: headers() })
  return res.data
}

export async function deleteKey(id: string) {
  const res = await axios.delete(`/api/security/keys/${id}`, { headers: headers() })
  return res.data
}

export async function toggleKeyStatus(id: string, enabled: boolean) {
  const res = await axios.put(`/api/security/keys/${id}/toggle`, { enabled }, { headers: headers() })
  return res.data
}

// ============================================================
// 密钥轮换
// ============================================================
export async function getRotationConfig() {
  const res = await axios.get(`/api/security/keys/rotation`, { headers: headers() })
  return res.data
}

export async function updateRotationConfig(data: any) {
  const res = await axios.put(`/api/security/keys/rotation`, data, { headers: headers() })
  return res.data
}

export async function rotateKeys() {
  const res = await axios.post(`/api/security/keys/rotate`, {}, { headers: headers() })
  return res.data
}
