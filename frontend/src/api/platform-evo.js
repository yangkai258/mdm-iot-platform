/**
 * 平台演进 API（Sprint 30-31）
 * 端侧 AI / 模型分片 / BLE Mesh / RTOS 优化
 * Base: /api/v1
 */

const API_BASE = '/api/v1'

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
// 端侧 AI 推理
// ════════════════════════════════════════════════════════════

/** 模型列表 */
export async function getEdgeModels(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`${API_BASE}/edge-ai/models?${qs}`)
}

/** 模型详情 */
export async function getEdgeModelDetail(id) {
  return request(`${API_BASE}/edge-ai/models/${id}`)
}

/** 部署模型 */
export async function deployEdgeModel(data) {
  return request(`${API_BASE}/edge-ai/models/${data.model_id}/deploy`, { method: 'POST', body: JSON.stringify(data) })
}

/** 卸载模型 */
export async function undeployEdgeModel(modelId, deviceId) {
  return request(`${API_BASE}/edge-ai/models/${modelId}/undeploy`, { method: 'POST', body: JSON.stringify({ device_id: deviceId }) })
}

/** 推理性能指标 */
export async function getInferenceMetrics(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`${API_BASE}/edge-ai/metrics?${qs}`)
}

/** 推理日志 */
export async function getInferenceLogs(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`${API_BASE}/edge-ai/logs?${qs}`)
}

/** 设备推理状态 */
export async function getDeviceInferenceStatus(deviceId) {
  return request(`${API_BASE}/edge-ai/devices/${deviceId}/inference-status`)
}

// ════════════════════════════════════════════════════════════
// 模型分片
// ════════════════════════════════════════════════════════════

/** 分片列表 */
export async function getModelShards(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`${API_BASE}/model-shards?${qs}`)
}

/** 分片详情 */
export async function getShardDetail(id) {
  return request(`${API_BASE}/model-shards/${id}`)
}

/** 创建分片 */
export async function createShard(data) {
  return request(`${API_BASE}/model-shards`, { method: 'POST', body: JSON.stringify(data) })
}

/** 更新分片 */
export async function updateShard(id, data) {
  return request(`${API_BASE}/model-shards/${id}`, { method: 'PUT', body: JSON.stringify(data) })
}

/** 删除分片 */
export async function deleteShard(id) {
  return request(`${API_BASE}/model-shards/${id}`, { method: 'DELETE' })
}

/** 加载分片到设备 */
export async function loadShardToDevice(shardId, deviceId) {
  return request(`${API_BASE}/model-shards/${shardId}/load`, { method: 'POST', body: JSON.stringify({ device_id: deviceId }) })
}

/** 从设备卸载分片 */
export async function unloadShardFromDevice(shardId, deviceId) {
  return request(`${API_BASE}/model-shards/${shardId}/unload`, { method: 'POST', body: JSON.stringify({ device_id: deviceId }) })
}

/** 分片版本列表 */
export async function getShardVersions(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`${API_BASE}/model-shards/versions?${qs}`)
}

/** 创建分片版本 */
export async function createShardVersion(data) {
  return request(`${API_BASE}/model-shards/versions`, { method: 'POST', body: JSON.stringify(data) })
}

/** 发布分片版本 */
export async function publishShardVersion(id) {
  return request(`${API_BASE}/model-shards/versions/${id}/publish`, { method: 'POST' })
}

/** 设备分片映射 */
export async function getDeviceShardMappings(deviceId) {
  return request(`${API_BASE}/model-shards/devices/${deviceId}/mappings`)
}

// ════════════════════════════════════════════════════════════
// BLE Mesh 网络
// ════════════════════════════════════════════════════════════

/** 网络拓扑 */
export async function getMeshTopology() {
  return request(`${API_BASE}/ble-mesh/topology`)
}

/** 节点列表 */
export async function getMeshNodes(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`${API_BASE}/ble-mesh/nodes?${qs}`)
}

/** 节点详情 */
export async function getMeshNodeDetail(id) {
  return request(`${API_BASE}/ble-mesh/nodes/${id}`)
}

/** 节点配置更新 */
export async function updateMeshNode(id, data) {
  return request(`${API_BASE}/ble-mesh/nodes/${id}`, { method: 'PUT', body: JSON.stringify(data) })
}

/** 网络管理操作 */
export async function meshNetworkAction(action, data = {}) {
  return request(`${API_BASE}/ble-mesh/network/${action}`, { method: 'POST', body: JSON.stringify(data) })
}

/** Mesh 配置 */
export async function getMeshConfig() {
  return request(`${API_BASE}/ble-mesh/config`)
}

/** 更新 Mesh 配置 */
export async function updateMeshConfig(data) {
  return request(`${API_BASE}/ble-mesh/config`, { method: 'PUT', body: JSON.stringify(data) })
}

/** 节点配网 */
export async function provisionMeshNode(nodeId, data) {
  return request(`${API_BASE}/ble-mesh/nodes/${nodeId}/provision`, { method: 'POST', body: JSON.stringify(data) })
}

// ════════════════════════════════════════════════════════════
// RTOS 优化
// ════════════════════════════════════════════════════════════

/** RTOS 配置 */
export async function getRtosConfig(deviceId) {
  return request(`${API_BASE}/rtos/config/${deviceId}`)
}

/** 更新 RTOS 配置 */
export async function updateRtosConfig(deviceId, data) {
  return request(`${API_BASE}/rtos/config/${deviceId}`, { method: 'PUT', body: JSON.stringify(data) })
}

/** 性能仪表板 */
export async function getRtosMetrics(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`${API_BASE}/rtos/metrics?${qs}`)
}

/** 实时性能数据 */
export async function getRtosLiveMetrics(deviceId) {
  return request(`${API_BASE}/rtos/devices/${deviceId}/live-metrics`)
}

/** 优化建议 */
export async function getOptimizationSuggestions(deviceId) {
  return request(`${API_BASE}/rtos/devices/${deviceId}/suggestions`)
}

/** 应用优化建议 */
export async function applyOptimization(deviceId, suggestionId) {
  return request(`${API_BASE}/rtos/devices/${deviceId}/apply-suggestion`, {
    method: 'POST',
    body: JSON.stringify({ suggestion_id: suggestionId })
  })
}

/** 调度策略列表 */
export async function getSchedulerPolicies(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`${API_BASE}/rtos/scheduler-policies?${qs}`)
}

/** 创建调度策略 */
export async function createSchedulerPolicy(data) {
  return request(`${API_BASE}/rtos/scheduler-policies`, { method: 'POST', body: JSON.stringify(data) })
}

/** 更新调度策略 */
export async function updateSchedulerPolicy(id, data) {
  return request(`${API_BASE}/rtos/scheduler-policies/${id}`, { method: 'PUT', body: JSON.stringify(data) })
}

/** 内存使用分析 */
export async function getMemoryAnalysis(deviceId) {
  return request(`${API_BASE}/rtos/devices/${deviceId}/memory-analysis`)
}

/** CPU 分析 */
export async function getCpuAnalysis(deviceId) {
  return request(`${API_BASE}/rtos/devices/${deviceId}/cpu-analysis`)
}
