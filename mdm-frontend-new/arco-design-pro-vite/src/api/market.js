/**
 * 内容生态市场 API
 * 表情包市场 / 动作资源库 / 声音定制
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
// 表情包市场
// ════════════════════════════════════════════════════════════

/** 表情包分类 */
export async function getEmoticonCategories() {
  return request(`/api/emoticons/categories`)
}

/** 表情包列表 */
export async function getEmoticonList(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`/api/emoticons?${qs}`)
}

/** 表情包详情 */
export async function getEmoticonDetail(id) {
  return request(`/api/emoticons/${id}`)
}

/** 上传表情包 */
export async function uploadEmoticon(data) {
  const formData = new FormData()
  Object.entries(data).forEach(([k, v]) => formData.append(k, v))
  const res = await fetch(`/api/emoticons`, {
    method: 'POST',
    headers: { 'Authorization': `Bearer ${getToken()}` },
    body: formData
  })
  const json = await res.json()
  if (json.code !== 0 && json.code !== 200) throw new Error(json.message || '上传失败')
  return json
}

/** 删除表情包 */
export async function deleteEmoticon(id) {
  return request(`/api/emoticons/${id}`, { method: 'DELETE' })
}

/** 下载表情包 */
export async function downloadEmoticon(id) {
  return request(`/api/emoticons/${id}/download`)
}

// ════════════════════════════════════════════════════════════
// 动作资源库
// ════════════════════════════════════════════════════════════

/** 动作分类 */
export async function getActionCategories() {
  return request(`/api/actions/categories`)
}

/** 动作列表 */
export async function getActionList(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`/api/actions?${qs}`)
}

/** 动作详情 */
export async function getActionDetail(id) {
  return request(`/api/actions/${id}`)
}

/** 创建自定义动作 */
export async function createAction(data) {
  return request(`/api/actions`, { method: 'POST', body: JSON.stringify(data) })
}

/** 更新动作 */
export async function updateAction(id, data) {
  return request(`/api/actions/${id}`, { method: 'PUT', body: JSON.stringify(data) })
}

/** 删除动作 */
export async function deleteAction(id) {
  return request(`/api/actions/${id}`, { method: 'DELETE' })
}

/** 发布动作到市场 */
export async function publishAction(id) {
  return request(`/api/actions/${id}/publish`, { method: 'POST' })
}

/** 预览动作动画 */
export async function previewAction(id) {
  return request(`/api/actions/${id}/preview`)
}

// ════════════════════════════════════════════════════════════
// 声音定制
// ════════════════════════════════════════════════════════════

/** 声音类型列表 */
export async function getVoiceList(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`/api/voices?${qs}`)
}

/** 声音详情 */
export async function getVoiceDetail(id) {
  return request(`/api/voices/${id}`)
}

/** 创建声音配置 */
export async function createVoice(data) {
  return request(`/api/voices`, { method: 'POST', body: JSON.stringify(data) })
}

/** 更新声音配置 */
export async function updateVoice(id, data) {
  return request(`/api/voices/${id}`, { method: 'PUT', body: JSON.stringify(data) })
}

/** 删除声音配置 */
export async function deleteVoice(id) {
  return request(`/api/voices/${id}`, { method: 'DELETE' })
}

/** 预览声音 */
export async function previewVoice(id) {
  return request(`/api/voices/${id}/preview`)
}

// ════════════════════════════════════════════════════════════
// 内容审核
// ════════════════════════════════════════════════════════════

/** 审核内容列表 */
export async function getReviewList(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`/api/market/reviews?${qs}`)
}

/** 待审核内容数量 */
export async function getReviewPendingCount() {
  return request(`/api/market/reviews/pending-count`)
}

/** 审核内容详情 */
export async function getReviewDetail(id) {
  return request(`/api/market/reviews/${id}`)
}

/** 审核内容（通过/拒绝） */
export async function reviewContent(id, data) {
  return request(`/api/market/reviews/${id}/review`, {
    method: 'POST',
    body: JSON.stringify(data)
  })
}

/** 批量审核 */
export async function batchReviewContent(data) {
  return request(`/api/market/reviews/batch-review`, {
    method: 'POST',
    body: JSON.stringify(data)
  })
}

/** 审核历史 */
export async function getReviewHistory(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`/api/market/reviews/history?${qs}`)
}
