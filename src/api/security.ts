import axios from 'axios'

const BASE_URL = '/api/v1'

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
  const res = await axios.get(`${BASE_URL}/ldap/config`, { headers: headers() })
  return res.data
}

export async function updateLdapConfig(data: any) {
  const res = await axios.put(`${BASE_URL}/ldap/config`, data, { headers: headers() })
  return res.data
}

export async function testLdapConnection(data: any) {
  const res = await axios.post(`${BASE_URL}/ldap/test`, data, { headers: headers() })
  return res.data
}

export async function getLdapUsers(params: any) {
  const res = await axios.get(`${BASE_URL}/ldap/users`, { params, headers: headers() })
  return res.data
}

export async function getLdapGroups() {
  const res = await axios.get(`${BASE_URL}/ldap/groups`, { headers: headers() })
  return res.data
}

export async function syncLdapUsers() {
  const res = await axios.post(`${BASE_URL}/ldap/sync`, {}, { headers: headers() })
  return res.data
}

export async function setLdapGroupMapping(data: any) {
  const res = await axios.post(`${BASE_URL}/ldap/group-mapping`, data, { headers: headers() })
  return res.data
}

// ============================================================
// 证书管理
// ============================================================
export async function getCertificates(params: any) {
  const res = await axios.get(`${BASE_URL}/certificates`, { params, headers: headers() })
  return res.data
}

export async function getCertificate(id: string) {
  const res = await axios.get(`${BASE_URL}/certificates/${id}`, { headers: headers() })
  return res.data
}

export async function uploadCertificate(data: any) {
  const res = await axios.post(`${BASE_URL}/certificates`, data, {
    headers: { ...headers(), 'Content-Type': 'multipart/form-data' }
  })
  return res.data
}

export async function deleteCertificate(id: string) {
  const res = await axios.delete(`${BASE_URL}/certificates/${id}`, { headers: headers() })
  return res.data
}

export async function revokeCertificate(id: string, reason?: string) {
  const res = await axios.post(`${BASE_URL}/certificates/${id}/revoke`, { reason }, { headers: headers() })
  return res.data
}

export async function renewCertificate(id: string) {
  const res = await axios.post(`${BASE_URL}/certificates/${id}/renew`, {}, { headers: headers() })
  return res.data
}

export async function issueCertificate(data: any) {
  const res = await axios.post(`${BASE_URL}/certificates/issue`, data, { headers: headers() })
  return res.data
}

export async function downloadCertificate(id: string) {
  const res = await axios.get(`${BASE_URL}/certificates/${id}/download`, {
    headers: headers(),
    responseType: 'blob'
  })
  return res.data
}

// ============================================================
// 设备安全操作
// ============================================================
export async function lockDevice(deviceId: string) {
  const res = await axios.post(`${BASE_URL}/devices/${deviceId}/lock`, {}, { headers: headers() })
  return res.data
}

export async function unlockDevice(deviceId: string) {
  const res = await axios.post(`${BASE_URL}/devices/${deviceId}/unlock`, {}, { headers: headers() })
  return res.data
}

export async function wipeDevice(deviceId: string) {
  const res = await axios.post(`${BASE_URL}/devices/${deviceId}/wipe`, {}, { headers: headers() })
  return res.data
}

export async function confirmWipeDevice(deviceId: string, data: any) {
  const res = await axios.post(`${BASE_URL}/devices/${deviceId}/wipe/confirm`, data, { headers: headers() })
  return res.data
}

export async function getWipeHistory(deviceId: string, params?: any) {
  const res = await axios.get(`${BASE_URL}/devices/${deviceId}/wipe-history`, { params, headers: headers() })
  return res.data
}

// ============================================================
// 数据权限
// ============================================================
export async function getDataPermissionRoles(roleId: string) {
  const res = await axios.get(`${BASE_URL}/data-permissions/roles/${roleId}`, { headers: headers() })
  return res.data
}

export async function updateDataPermissionRoles(roleId: string, data: any) {
  const res = await axios.put(`${BASE_URL}/data-permissions/roles/${roleId}`, data, { headers: headers() })
  return res.data
}

export async function getDataPermissionUsers(userId: string) {
  const res = await axios.get(`${BASE_URL}/data-permissions/users/${userId}`, { headers: headers() })
  return res.data
}

export async function updateDataPermissionUsers(userId: string, data: any) {
  const res = await axios.put(`${BASE_URL}/data-permissions/users/${userId}`, data, { headers: headers() })
  return res.data
}

export async function getDataPermissions(params?: any) {
  const res = await axios.get(`${BASE_URL}/data-permissions`, { params, headers: headers() })
  return res.data
}

export async function createDataPermission(data: any) {
  const res = await axios.post(`${BASE_URL}/data-permissions`, data, { headers: headers() })
  return res.data
}

export async function updateDataPermission(id: string, data: any) {
  const res = await axios.put(`${BASE_URL}/data-permissions/${id}`, data, { headers: headers() })
  return res.data
}

export async function deleteDataPermission(id: string) {
  const res = await axios.delete(`${BASE_URL}/data-permissions/${id}`, { headers: headers() })
  return res.data
}

// ============================================================
// 辅助：获取设备列表
// ============================================================
export async function getDevices(params?: any) {
  const res = await axios.get(`${BASE_URL}/devices`, { params, headers: headers() })
  return res.data
}

// ============================================================
// 辅助：获取角色列表
// ============================================================
export async function getRoles(params?: any) {
  const res = await axios.get(`${BASE_URL}/roles`, { params, headers: headers() })
  return res.data
}

// ============================================================
// 辅助：获取用户列表
// ============================================================
export async function getUsers(params?: any) {
  const res = await axios.get(`${BASE_URL}/users`, { params, headers: headers() })
  return res.data
}

// ============================================================
// 双因素认证 (2FA)
// ============================================================
export async function getTwoFactorStatus() {
  const res = await axios.get(`${BASE_URL}/two-factor/status`, { headers: headers() })
  return res.data
}

export async function enableTwoFactor(data: any) {
  const res = await axios.post(`${BASE_URL}/two-factor/enable`, data, { headers: headers() })
  return res.data
}

export async function disableTwoFactor(data: any) {
  const res = await axios.post(`${BASE_URL}/two-factor/disable`, data, { headers: headers() })
  return res.data
}

export async function verifyTwoFactorCode(data: any) {
  const res = await axios.post(`${BASE_URL}/two-factor/verify`, data, { headers: headers() })
  return res.data
}

export async function getTotpQrCode() {
  const res = await axios.get(`${BASE_URL}/two-factor/totp/qr`, { headers: headers() })
  return res.data
}

export async function getRecoveryCodes() {
  const res = await axios.get(`${BASE_URL}/two-factor/recovery-codes`, { headers: headers() })
  return res.data
}

export async function regenerateRecoveryCodes() {
  const res = await axios.post(`${BASE_URL}/two-factor/recovery-codes/regenerate`, {}, { headers: headers() })
  return res.data
}

export async function getTrustDevices() {
  const res = await axios.get(`${BASE_URL}/two-factor/trust-devices`, { headers: headers() })
  return res.data
}

export async function removeTrustDevice(id: string) {
  const res = await axios.delete(`${BASE_URL}/two-factor/trust-devices/${id}`, { headers: headers() })
  return res.data
}

export async function changeTwoFactorMethod(data: any) {
  const res = await axios.post(`${BASE_URL}/two-factor/change-method`, data, { headers: headers() })
  return res.data
}

export async function sendSmsVerifyCode(data: any) {
  const res = await axios.post(`${BASE_URL}/two-factor/sms/send`, data, { headers: headers() })
  return res.data
}

export async function sendEmailVerifyCode(data: any) {
  const res = await axios.post(`${BASE_URL}/two-factor/email/send`, data, { headers: headers() })
  return res.data
}

// ============================================================
// 会话管理
// ============================================================
export async function getSessions(params?: any) {
  const res = await axios.get(`${BASE_URL}/sessions`, { params, headers: headers() })
  return res.data
}

export async function terminateSession(id: string) {
  const res = await axios.delete(`${BASE_URL}/sessions/${id}`, { headers: headers() })
  return res.data
}

export async function terminateOtherSessions() {
  const res = await axios.delete(`${BASE_URL}/sessions/others`, { headers: headers() })
  return res.data
}

export async function terminateAllSessions() {
  const res = await axios.delete(`${BASE_URL}/sessions/all`, { headers: headers() })
  return res.data
}

export async function getSessionStats() {
  const res = await axios.get(`${BASE_URL}/sessions/stats`, { headers: headers() })
  return res.data
}

// ============================================================
// 安全审计
// ============================================================
export async function getSecurityLogs(params?: any) {
  const res = await axios.get(`${BASE_URL}/security/logs`, { params, headers: headers() })
  return res.data
}

export async function getSecurityLogDetail(id: string) {
  const res = await axios.get(`${BASE_URL}/security/logs/${id}`, { headers: headers() })
  return res.data
}

export async function exportSecurityLogsApi(params?: any) {
  const res = await axios.get(`${BASE_URL}/security/logs/export`, {
    params,
    headers: headers(),
    responseType: 'blob'
  })
  return res.data
}

export async function getSecurityOverview() {
  const res = await axios.get(`${BASE_URL}/security/overview`, { headers: headers() })
  return res.data
}

export async function getSecurityReports(params?: any) {
  const res = await axios.get(`${BASE_URL}/security/reports`, { params, headers: headers() })
  return res.data
}

export async function generateSecurityReport(data: any) {
  const res = await axios.post(`${BASE_URL}/security/reports`, data, { headers: headers() })
  return res.data
}

export async function downloadSecurityReport(id: string) {
  const res = await axios.get(`${BASE_URL}/security/reports/${id}/download`, {
    headers: headers(),
    responseType: 'blob'
  })
  return res.data
}

export async function getRecentThreats() {
  const res = await axios.get(`${BASE_URL}/security/threats`, { headers: headers() })
  return res.data
}

export async function getThreatDistribution() {
  const res = await axios.get(`${BASE_URL}/security/threats/distribution`, { headers: headers() })
  return res.data
}
