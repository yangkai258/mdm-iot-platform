/**
 * AI 模块 API
 * AI 行为监控、模型版本管理、沙箱测试、质量指标
 */

const BASE_URL = '/api/ai'

// ============ AI 行为监控 ============

/**
 * 上报 AI 行为事件
 */
export async function postAiMonitorEvent(data) {
  const token = localStorage.getItem('token')
  const res = await fetch(`/api/monitor/events`, {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(data)
  })
  return res.json()
}

/**
 * 获取 AI 行为事件列表
 * @param {Object} params - { device_id, user_id, model_version, behavior_type, status, start_time, end_time, page, page_size }
 */
export async function getAiMonitorEvents(params = {}) {
  const token = localStorage.getItem('token')
  const qs = new URLSearchParams(params).toString()
  const res = await fetch(`/api/monitor/events?${qs}`, {
    headers: { 'Authorization': `Bearer ${token}` }
  })
  return res.json()
}

/**
 * 获取 AI 行为事件详情
 * @param {string|number} id
 */
export async function getAiMonitorEventById(id) {
  const token = localStorage.getItem('token')
  const res = await fetch(`/api/monitor/events/${id}`, {
    headers: { 'Authorization': `Bearer ${token}` }
  })
  return res.json()
}

/**
 * 获取 AI 监控统计数据
 */
export async function getAiMonitorStats(params = {}) {
  const token = localStorage.getItem('token')
  const qs = new URLSearchParams(params).toString()
  const res = await fetch(`/api/monitor/stats?${qs}`, {
    headers: { 'Authorization': `Bearer ${token}` }
  })
  return res.json()
}

// ============ 模型版本管理 ============

/**
 * 获取模型版本列表
 * @param {Object} params - { model_id, status, page, page_size }
 */
export async function getAiModels(params = {}) {
  const token = localStorage.getItem('token')
  const qs = new URLSearchParams(params).toString()
  const res = await fetch(`/api/models?${qs}`, {
    headers: { 'Authorization': `Bearer ${token}` }
  })
  return res.json()
}

/**
 * 注册新模型
 * @param {Object} data - { model_id, version_id, version_name, model_type, model_name, ... }
 */
export async function postAiModel(data) {
  const token = localStorage.getItem('token')
  const res = await fetch(`/api/models`, {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(data)
  })
  return res.json()
}

/**
 * 获取模型详情
 * @param {string|number} id
 */
export async function getAiModelById(id) {
  const token = localStorage.getItem('token')
  const res = await fetch(`/api/models/${id}`, {
    headers: { 'Authorization': `Bearer ${token}` }
  })
  return res.json()
}

/**
 * 获取模型版本历史
 * @param {string|number} id
 */
export async function getAiModelVersions(id) {
  const token = localStorage.getItem('token')
  const res = await fetch(`/api/models/${id}/versions`, {
    headers: { 'Authorization': `Bearer ${token}` }
  })
  return res.json()
}

/**
 * 发布新版本
 * @param {string|number} id
 * @param {Object} data - { version_id, version_name, ... }
 */
export async function postAiModelVersion(id, data) {
  const token = localStorage.getItem('token')
  const res = await fetch(`/api/models/${id}/versions`, {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(data)
  })
  return res.json()
}

/**
 * 回滚模型版本
 * @param {string|number} id
 * @param {Object} data - { target_version, reason }
 */
export async function postAiModelRollback(id, data) {
  const token = localStorage.getItem('token')
  const res = await fetch(`/api/models/${id}/rollback`, {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(data)
  })
  return res.json()
}

/**
 * 废弃模型版本
 * @param {string|number} id
 * @param {Object} data - { reason }
 */
export async function postAiModelDeprecate(id, data) {
  const token = localStorage.getItem('token')
  const res = await fetch(`/api/models/${id}/deprecate`, {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(data)
  })
  return res.json()
}

// ============ 沙箱测试 ============

/**
 * 创建沙箱测试任务
 * @param {Object} data - { test_name, model_version_id, test_type, test_cases, description }
 */
export async function postAiSandboxTest(data) {
  const token = localStorage.getItem('token')
  const res = await fetch(`/api/sandbox/test`, {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(data)
  })
  return res.json()
}

/**
 * 获取沙箱测试结果
 * @param {string} taskId
 */
export async function getAiSandboxTestResult(taskId) {
  const token = localStorage.getItem('token')
  const res = await fetch(`/api/sandbox/test/${taskId}`, {
    headers: { 'Authorization': `Bearer ${token}` }
  })
  return res.json()
}

/**
 * 获取沙箱测试历史列表
 * @param {Object} params - { model_id, test_type, status, page, page_size }
 */
export async function getAiSandboxTests(params = {}) {
  const token = localStorage.getItem('token')
  const qs = new URLSearchParams(params).toString()
  const res = await fetch(`/api/sandbox/tests?${qs}`, {
    headers: { 'Authorization': `Bearer ${token}` }
  })
  return res.json()
}

// ============ 质量指标 ============

/**
 * 获取 AI 质量指标
 * @param {Object} params - { start_time, end_time, time_range }
 */
export async function getAiQualityMetrics(params = {}) {
  const token = localStorage.getItem('token')
  const qs = new URLSearchParams(params).toString()
  const res = await fetch(`/api/quality/metrics?${qs}`, {
    headers: { 'Authorization': `Bearer ${token}` }
  })
  return res.json()
}

/**
 * 获取 AI 质量趋势数据
 * @param {Object} params - { start_time, end_time, time_range }
 */
export async function getAiQualityMetricsTrend(params = {}) {
  const token = localStorage.getItem('token')
  const qs = new URLSearchParams(params).toString()
  const res = await fetch(`/api/quality/metrics/trend?${qs}`, {
    headers: { 'Authorization': `Bearer ${token}` }
  })
  return res.json()
}
