/**
 * 高级安全功能 API（Sprint 32）
 * 审计日志 / 合规管理 / 数据导出 / GDPR 请求
 * Base: /api/v1
 */

const API_BASE = '/api'

function getToken() {
  return localStorage.getItem('token') || ''
}

async function request(url, options = {}) {
  const res = await fetch(url, {
    ...options,
    headers: {
      'Authorization': `Bearer ${getToken()}`,
      'Content-Type': 'application/json',
      ...(options.headers || {})
    }
  })
  const data = await res.json()
  if (data.code !== 0 && data.code !== 200) {
    throw new Error(data.message || '请求失败')
  }
  return data
}

// ════════════════════════════════════════════════════════════
// 审计日志
// ════════════════════════════════════════════════════════════

/** 审计日志列表 */
export async function getAuditLogs(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`/api/security-evo/audit/logs?${qs}`)
}

/** 审计日志详情 */
export async function getAuditLogDetail(id) {
  return request(`/api/security-evo/audit/logs/${id}`)
}

/** 导出审计日志 */
export async function exportAuditLogs(params = {}) {
  const qs = new URLSearchParams(params).toString()
  const res = await fetch(`/api/security-evo/audit/logs/export?${qs}`, {
    headers: { 'Authorization': `Bearer ${getToken()}` },
    responseType: 'blob'
  })
  return res.blob()
}

/** 审计统计 */
export async function getAuditStatistics(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`/api/security-evo/audit/statistics?${qs}`)
}

/** 生成审计报告 */
export async function generateAuditReport(data) {
  return request(`/api/security-evo/audit/reports`, { method: 'POST', body: JSON.stringify(data) })
}

/** 审计报告列表 */
export async function getAuditReports(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`/api/security-evo/audit/reports?${qs}`)
}

/** 下载审计报告 */
export async function downloadAuditReport(reportId) {
  const res = await fetch(`/api/security-evo/audit/reports/${reportId}/download`, {
    headers: { 'Authorization': `Bearer ${getToken()}` },
    responseType: 'blob'
  })
  return res.blob()
}

// ════════════════════════════════════════════════════════════
// 合规管理
// ════════════════════════════════════════════════════════════

/** 合规状态仪表板 */
export async function getComplianceDashboard() {
  return request(`/api/security-evo/compliance/dashboard`)
}

/** 合规检查项列表 */
export async function getComplianceChecks(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`/api/security-evo/compliance/checks?${qs}`)
}

/** 合规检查项详情 */
export async function getComplianceCheckDetail(id) {
  return request(`/api/security-evo/compliance/checks/${id}`)
}

/** 触发合规检查 */
export async function triggerComplianceCheck(id) {
  return request(`/api/security-evo/compliance/checks/${id}/run`, { method: 'POST' })
}

/** 合规报告列表 */
export async function getComplianceReports(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`/api/security-evo/compliance/reports?${qs}`)
}

/** 合规报告详情 */
export async function getComplianceReportDetail(id) {
  return request(`/api/security-evo/compliance/reports/${id}`)
}

/** 生成合规报告 */
export async function generateComplianceReport(data) {
  return request(`/api/security-evo/compliance/reports`, { method: 'POST', body: JSON.stringify(data) })
}

/** 下载合规报告 */
export async function downloadComplianceReport(reportId) {
  const res = await fetch(`/api/security-evo/compliance/reports/${reportId}/download`, {
    headers: { 'Authorization': `Bearer ${getToken()}` },
    responseType: 'blob'
  })
  return res.blob()
}

/** 法规清单 */
export async function getRegulations(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`/api/security-evo/compliance/regulations?${qs}`)
}

/** 法规详情 */
export async function getRegulationDetail(id) {
  return request(`/api/security-evo/compliance/regulations/${id}`)
}

/** 创建/更新法规 */
export async function saveRegulation(data) {
  return request(`/api/security-evo/compliance/regulations`, { method: data.id ? 'PUT' : 'POST', body: JSON.stringify(data) })
}

/** 删除法规 */
export async function deleteRegulation(id) {
  return request(`/api/security-evo/compliance/regulations/${id}`, { method: 'DELETE' })
}

/** 合规历史趋势 */
export async function getComplianceTrend(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`/api/security-evo/compliance/trend?${qs}`)
}

// ════════════════════════════════════════════════════════════
// 数据导出
// ════════════════════════════════════════════════════════════

/** 导出任务列表 */
export async function getExportTasks(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`/api/security-evo/data-export/tasks?${qs}`)
}

/** 导出任务详情 */
export async function getExportTaskDetail(id) {
  return request(`/api/security-evo/data-export/tasks/${id}`)
}

/** 创建导出任务 */
export async function createExportTask(data) {
  return request(`/api/security-evo/data-export/tasks`, { method: 'POST', body: JSON.stringify(data) })
}

/** 取消导出任务 */
export async function cancelExportTask(id) {
  return request(`/api/security-evo/data-export/tasks/${id}/cancel`, { method: 'POST' })
}

/** 删除导出任务 */
export async function deleteExportTask(id) {
  return request(`/api/security-evo/data-export/tasks/${id}`, { method: 'DELETE' })
}

/** 下载导出文件 */
export async function downloadExportFile(taskId) {
  const res = await fetch(`/api/security-evo/data-export/tasks/${taskId}/download`, {
    headers: { 'Authorization': `Bearer ${getToken()}` },
    responseType: 'blob'
  })
  return res.blob()
}

/** 导出模板列表 */
export async function getExportTemplates(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`/api/security-evo/data-export/templates?${qs}`)
}

/** 导出格式列表 */
export async function getExportFormats() {
  return request(`/api/security-evo/data-export/formats`)
}

// ════════════════════════════════════════════════════════════
// GDPR 请求
// ════════════════════════════════════════════════════════════

/** GDPR 请求列表 */
export async function getGdprRequests(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`/api/security-evo/gdpr/requests?${qs}`)
}

/** GDPR 请求详情 */
export async function getGdprRequestDetail(id) {
  return request(`/api/security-evo/gdpr/requests/${id}`)
}

/** 创建 GDPR 请求 */
export async function createGdprRequest(data) {
  return request(`/api/security-evo/gdpr/requests`, { method: 'POST', body: JSON.stringify(data) })
}

/** 更新 GDPR 请求状态 */
export async function updateGdprRequest(id, data) {
  return request(`/api/security-evo/gdpr/requests/${id}`, { method: 'PUT', body: JSON.stringify(data) })
}

/** 审批 GDPR 请求 */
export async function approveGdprRequest(id, data) {
  return request(`/api/security-evo/gdpr/requests/${id}/approve`, { method: 'POST', body: JSON.stringify(data) })
}

/** 拒绝 GDPR 请求 */
export async function rejectGdprRequest(id, data) {
  return request(`/api/security-evo/gdpr/requests/${id}/reject`, { method: 'POST', body: JSON.stringify(data) })
}

/** 完成 GDPR 请求 */
export async function completeGdprRequest(id, data) {
  return request(`/api/security-evo/gdpr/requests/${id}/complete`, { method: 'POST', body: JSON.stringify(data) })
}

/** GDPR 请求统计 */
export async function getGdprStatistics() {
  return request(`/api/security-evo/gdpr/statistics`)
}

/** GDPR 请求类型列表 */
export async function getGdprRequestTypes() {
  return request(`/api/security-evo/gdpr/request-types`)
}
