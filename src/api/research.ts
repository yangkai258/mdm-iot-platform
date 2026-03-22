/**
 * 研究平台 API
 * 数据集管理 + AI 行为实验
 */

const BASE_URL = '/api/v1/research'

const getToken = () => localStorage.getItem('token') || ''

const request = async (path: string, options: RequestInit = {}) => {
  const token = getToken()
  const res = await fetch(`${BASE_URL}${path}`, {
    ...options,
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json',
      ...options.headers
    }
  })
  return res.json()
}

// ============ 数据集管理 ============

/**
 * 获取数据集列表
 */
export async function getDataSets(params: {
  page?: number
  page_size?: number
  keyword?: string
  data_type?: string
  status?: string
} = {}) {
  const qs = new URLSearchParams(params as any).toString()
  return request(`/datasets?${qs}`)
}

/**
 * 创建数据集
 */
export async function postDataSet(data: {
  name: string
  description?: string
  data_type: string
  tags?: string[]
}) {
  return request('/datasets', { method: 'POST', body: JSON.stringify(data) })
}

/**
 * 获取数据集详情
 */
export async function getDataSetById(id: string | number) {
  return request(`/datasets/${id}`)
}

/**
 * 更新数据集
 */
export async function putDataSet(id: string | number, data: {
  name?: string
  description?: string
  tags?: string[]
  status?: string
}) {
  return request(`/datasets/${id}`, { method: 'PUT', body: JSON.stringify(data) })
}

/**
 * 删除数据集
 */
export async function deleteDataSet(id: string | number) {
  return request(`/datasets/${id}`, { method: 'DELETE' })
}

/**
 * 导出数据集
 */
export async function exportDataSet(id: string | number, format: 'csv' | 'json' | 'parquet' = 'csv') {
  const token = getToken()
  const res = await fetch(`${BASE_URL}/datasets/${id}/export?format=${format}`, {
    method: 'GET',
    headers: { 'Authorization': `Bearer ${token}` }
  })
  const blob = await res.blob()
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `dataset-${id}-${Date.now()}.${format}`
  a.click()
  URL.revokeObjectURL(url)
  return { code: 0 }
}

// ============ 匿名数据 ============

/**
 * 获取匿名数据列表
 */
export async function getAnonymizedData(params: {
  page?: number
  page_size?: number
  dataset_id?: string | number
  keyword?: string
  start_time?: string
  end_time?: string
} = {}) {
  const qs = new URLSearchParams(params as any).toString()
  return request(`/data/anonymized?${qs}`)
}

/**
 * 删除匿名数据
 */
export async function deleteAnonymizedData(id: string | number) {
  return request(`/data/anonymized/${id}`, { method: 'DELETE' })
}

/**
 * 批量删除匿名数据
 */
export async function deleteAnonymizedDataBatch(ids: (string | number)[]) {
  return request(`/data/anonymized/batch`, { method: 'DELETE', body: JSON.stringify({ ids }) })
}

// ============ AI 行为实验 ============

/**
 * 获取实验列表
 */
export async function getExperiments(params: {
  page?: number
  page_size?: number
  keyword?: string
  status?: string
  model_version?: string
} = {}) {
  const qs = new URLSearchParams(params as any).toString()
  return request(`/experiments?${qs}`)
}

/**
 * 创建实验
 */
export async function postExperiment(data: {
  name: string
  description?: string
  model_version: string
  config: Record<string, any>
  dataset_ids?: (string | number)[]
}) {
  return request('/experiments', { method: 'POST', body: JSON.stringify(data) })
}

/**
 * 获取实验详情
 */
export async function getExperimentById(id: string | number) {
  return request(`/experiments/${id}`)
}

/**
 * 更新实验
 */
export async function putExperiment(id: string | number, data: {
  name?: string
  description?: string
  config?: Record<string, any>
  dataset_ids?: (string | number)[]
}) {
  return request(`/experiments/${id}`, { method: 'PUT', body: JSON.stringify(data) })
}

/**
 * 删除实验
 */
export async function deleteExperiment(id: string | number) {
  return request(`/experiments/${id}`, { method: 'DELETE' })
}

/**
 * 开始实验
 */
export async function startExperiment(id: string | number) {
  return request(`/experiments/${id}/start`, { method: 'POST' })
}

/**
 * 停止实验
 */
export async function stopExperiment(id: string | number) {
  return request(`/experiments/${id}/stop`, { method: 'POST' })
}

/**
 * 获取实验运行状态
 */
export async function getExperimentStatus(id: string | number) {
  return request(`/experiments/${id}/status`)
}

/**
 * 获取实验日志
 */
export async function getExperimentLogs(id: string | number, params: {
  page?: number
  page_size?: number
} = {}) {
  const qs = new URLSearchParams(params as any).toString()
  return request(`/experiments/${id}/logs?${qs}`)
}

/**
 * 获取实验结果
 */
export async function getExperimentResults(id: string | number) {
  return request(`/experiments/${id}/results`)
}
