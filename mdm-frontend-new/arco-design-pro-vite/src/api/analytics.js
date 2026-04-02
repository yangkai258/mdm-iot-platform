/**
 * 数据分析模块 API
 * Dashboard / 设备统计 / OTA统计 / 会员统计 / 告警统计
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
// Dashboard 大盘
// ════════════════════════════════════════════════════════════

/** Dashboard 核心指标统计 */
export async function getDashboardStats() {
  return request(`/api/dashboard/stats`)
}

/** 设备概览 */
export async function getDeviceStatsOverview(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`/api/stats/devices/overview?${qs}`)
}

/** 设备趋势 */
export async function getDeviceStatsTrend(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`/api/stats/devices/trend?${qs}`)
}

/** OTA 概览 */
export async function getOtaStatsOverview() {
  return request(`/api/stats/ota/overview`)
}

/** OTA 任务详情 */
export async function getOtaStatsTasks(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`/api/stats/ota/tasks?${qs}`)
}

/** OTA 版本分布 */
export async function getOtaVersionDistribution() {
  return request(`/api/stats/ota/version-distribution`)
}

/** 会员概览 */
export async function getMemberStatsOverview() {
  return request(`/api/stats/members/overview`)
}

/** 会员等级分布 */
export async function getMemberLevelDistribution() {
  return request(`/api/stats/members/level-distribution`)
}

/** 会员消费统计 */
export async function getMemberConsumption(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`/api/stats/members/consumption?${qs}`)
}

/** 告警概览 */
export async function getAlertStatsOverview() {
  return request(`/api/stats/alerts/overview`)
}

/** 告警趋势 */
export async function getAlertStatsTrend(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`/api/stats/alerts/trend?${qs}`)
}

/** 告警分布 */
export async function getAlertStatsDistribution() {
  return request(`/api/stats/alerts/distribution`)
}

// ════════════════════════════════════════════════════════════
// 漏斗分析
// ════════════════════════════════════════════════════════════

/** 漏斗列表 */
export async function getFunnelList(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`/api/analytics/funnels?${qs}`)
}

/** 漏斗详情 */
export async function getFunnelDetail(id) {
  return request(`/api/analytics/funnels/${id}`)
}

/** 创建漏斗 */
export async function createFunnel(data) {
  return request(`/api/analytics/funnels`, { method: 'POST', body: JSON.stringify(data) })
}

/** 更新漏斗 */
export async function updateFunnel(id, data) {
  return request(`/api/analytics/funnels/${id}`, { method: 'PUT', body: JSON.stringify(data) })
}

/** 删除漏斗 */
export async function deleteFunnel(id) {
  return request(`/api/analytics/funnels/${id}`, { method: 'DELETE' })
}

/** 漏斗数据（转化率） */
export async function getFunnelData(id, params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`/api/analytics/funnels/${id}/data?${qs}`)
}

// ════════════════════════════════════════════════════════════
// 群组分析
// ════════════════════════════════════════════════════════════

/** 群组列表 */
export async function getCohortList(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`/api/analytics/cohorts?${qs}`)
}

/** 群组详情 */
export async function getCohortDetail(id) {
  return request(`/api/analytics/cohorts/${id}`)
}

/** 创建群组 */
export async function createCohort(data) {
  return request(`/api/analytics/cohorts`, { method: 'POST', body: JSON.stringify(data) })
}

/** 更新群组 */
export async function updateCohort(id, data) {
  return request(`/api/analytics/cohorts/${id}`, { method: 'PUT', body: JSON.stringify(data) })
}

/** 删除群组 */
export async function deleteCohort(id) {
  return request(`/api/analytics/cohorts/${id}`, { method: 'DELETE' })
}

/** 群组对比数据 */
export async function getCohortComparison(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`/api/analytics/cohorts/comparison?${qs}`)
}

/** 群组热力图数据 */
export async function getCohortHeatmap(id, params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`/api/analytics/cohorts/${id}/heatmap?${qs}`)
}

// ════════════════════════════════════════════════════════════
// 留存分析
// ════════════════════════════════════════════════════════════

/** 留存概览 */
export async function getRetentionOverview(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`/api/analytics/retention/overview?${qs}`)
}

/** 留存曲线数据 */
export async function getRetentionCurve(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`/api/analytics/retention/curve?${qs}`)
}

/** 留存报表 */
export async function getRetentionReport(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`/api/analytics/retention/report?${qs}`)
}

/** 自定义留存分析 */
export async function getRetentionCustom(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`/api/analytics/retention/custom?${qs}`)
}

// ════════════════════════════════════════════════════════════
// 事件分析
// ════════════════════════════════════════════════════════════

/** 事件列表 */
export async function getEventList(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`/api/analytics/events?${qs}`)
}

/** 事件详情 */
export async function getEventDetail(id) {
  return request(`/api/analytics/events/${id}`)
}

/** 创建事件 */
export async function createEvent(data) {
  return request(`/api/analytics/events`, { method: 'POST', body: JSON.stringify(data) })
}

/** 更新事件 */
export async function updateEvent(id, data) {
  return request(`/api/analytics/events/${id}`, { method: 'PUT', body: JSON.stringify(data) })
}

/** 删除事件 */
export async function deleteEvent(id) {
  return request(`/api/analytics/events/${id}`, { method: 'DELETE' })
}

/** 事件趋势 */
export async function getEventTrend(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`/api/analytics/events/trend?${qs}`)
}

/** 事件分布 */
export async function getEventDistribution(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`/api/analytics/events/distribution?${qs}`)
}

/** 事件漏斗 */
export async function getEventFunnel(id, params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`/api/analytics/events/${id}/funnel?${qs}`)
}
