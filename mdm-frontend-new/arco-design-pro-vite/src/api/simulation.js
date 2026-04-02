import axios from 'axios'

const BASE_URL = '/api'

function getToken() {
  return localStorage.getItem('token') || ''
}

function headers() {
  return { Authorization: `Bearer ${getToken()}` }
}

// ========== 虚拟宠物 API ==========
export async function getSimulationPets(params = {}) {
  const res = await axios.get(`/api/simulation/pets`, { params, headers: headers() })
  return res.data
}

export async function getSimulationPet(id) {
  const res = await axios.get(`/api/simulation/pets/${id}`, { headers: headers() })
  return res.data
}

export async function createSimulationPet(data) {
  const res = await axios.post(`/api/simulation/pets`, data, { headers: headers() })
  return res.data
}

export async function updateSimulationPet(id, data) {
  const res = await axios.put(`/api/simulation/pets/${id}`, data, { headers: headers() })
  return res.data
}

export async function deleteSimulationPet(id) {
  const res = await axios.delete(`/api/simulation/pets/${id}`, { headers: headers() })
  return res.data
}

export async function interactWithPet(id, data) {
  const res = await axios.post(`/api/simulation/pets/${id}/interact`, data, { headers: headers() })
  return res.data
}

// ========== 测试用例 API ==========
export async function getTestCases(params = {}) {
  const res = await axios.get(`/api/simulation/testcases`, { params, headers: headers() })
  return res.data
}

export async function getTestCase(id) {
  const res = await axios.get(`/api/simulation/testcases/${id}`, { headers: headers() })
  return res.data
}

export async function createTestCase(data) {
  const res = await axios.post(`/api/simulation/testcases`, data, { headers: headers() })
  return res.data
}

export async function updateTestCase(id, data) {
  const res = await axios.put(`/api/simulation/testcases/${id}`, data, { headers: headers() })
  return res.data
}

export async function deleteTestCase(id) {
  const res = await axios.delete(`/api/simulation/testcases/${id}`, { headers: headers() })
  return res.data
}

export async function executeTestCase(id, data) {
  const res = await axios.post(`/api/simulation/testcases/${id}/execute`, data, { headers: headers() })
  return res.data
}

export async function batchExecuteTestCases(data) {
  const res = await axios.post(`/api/simulation/testcases/batch-execute`, data, { headers: headers() })
  return res.data
}

export async function getExecution(id) {
  const res = await axios.get(`/api/simulation/executions/${id}`, { headers: headers() })
  return res.data
}

// ========== 测试报告 API ==========
export async function getReports(params = {}) {
  const res = await axios.get(`/api/simulation/reports`, { params, headers: headers() })
  return res.data
}

export async function getReport(id) {
  const res = await axios.get(`/api/simulation/reports/${id}`, { headers: headers() })
  return res.data
}

// ========== 回放系统 API ==========
export async function getPlaybacks(params = {}) {
  const res = await axios.get(`/api/simulation/playbacks`, { params, headers: headers() })
  return res.data
}

export async function getPlayback(id) {
  const res = await axios.get(`/api/simulation/playbacks/${id}`, { headers: headers() })
  return res.data
}

export async function createPlayback(data) {
  const res = await axios.post(`/api/simulation/playbacks`, data, { headers: headers() })
  return res.data
}

export async function playPlayback(id, data) {
  const res = await axios.post(`/api/simulation/playbacks/${id}/play`, data, { headers: headers() })
  return res.data
}

export async function stopPlayback(id) {
  const res = await axios.post(`/api/simulation/playbacks/${id}/stop`, {}, { headers: headers() })
  return res.data
}

export async function comparePlaybacks(id, data) {
  const res = await axios.post(`/api/simulation/playbacks/${id}/compare`, data, { headers: headers() })
  return res.data
}

export async function deletePlayback(id) {
  const res = await axios.delete(`/api/simulation/playbacks/${id}`, { headers: headers() })
  return res.data
}

// ========== 仿真场景 API ==========
export async function getScenarios(params = {}) {
  const res = await axios.get(`/api/simulation/scenarios`, { params, headers: headers() })
  return res.data
}

export async function getScenario(id) {
  const res = await axios.get(`/api/simulation/scenarios/${id}`, { headers: headers() })
  return res.data
}

export async function createScenario(data) {
  const res = await axios.post(`/api/simulation/scenarios`, data, { headers: headers() })
  return res.data
}

export async function updateScenario(id, data) {
  const res = await axios.put(`/api/simulation/scenarios/${id}`, data, { headers: headers() })
  return res.data
}

export async function deleteScenario(id) {
  const res = await axios.delete(`/api/simulation/scenarios/${id}`, { headers: headers() })
  return res.data
}

export async function runScenario(id, data) {
  const res = await axios.post(`/api/simulation/scenarios/${id}/run`, data, { headers: headers() })
  return res.data
}

export async function importScenarios(data) {
  const res = await axios.post(`/api/simulation/scenarios/import`, data, { headers: headers() })
  return res.data
}

export async function exportScenario(id) {
  const res = await axios.post(`/api/simulation/scenarios/export/${id}`, {}, { headers: headers() })
  return res.data
}

// ========== 压力测试 API ==========
export async function getStressTests(params = {}) {
  const res = await axios.get(`/api/simulation/stress-tests`, { params, headers: headers() })
  return res.data
}

export async function getStressTest(id) {
  const res = await axios.get(`/api/simulation/stress-tests/${id}`, { headers: headers() })
  return res.data
}

export async function createStressTest(data) {
  const res = await axios.post(`/api/simulation/stress-tests`, data, { headers: headers() })
  return res.data
}

export async function updateStressTest(id, data) {
  const res = await axios.put(`/api/simulation/stress-tests/${id}`, data, { headers: headers() })
  return res.data
}

export async function deleteStressTest(id) {
  const res = await axios.delete(`/api/simulation/stress-tests/${id}`, { headers: headers() })
  return res.data
}

export async function startStressTest(id, data) {
  const res = await axios.post(`/api/simulation/stress-tests/${id}/start`, data, { headers: headers() })
  return res.data
}

export async function stopStressTest(id) {
  const res = await axios.post(`/api/simulation/stress-tests/${id}/stop`, {}, { headers: headers() })
  return res.data
}

export async function getStressTestStatus(id) {
  const res = await axios.get(`/api/simulation/stress-tests/${id}/status`, { headers: headers() })
  return res.data
}

export async function getStressTestReport(id) {
  const res = await axios.get(`/api/simulation/stress-tests/${id}/report`, { headers: headers() })
  return res.data
}

// ========== 仿真数据集 API ==========
export async function getDatasets(params = {}) {
  const res = await axios.get(`/api/simulation/datasets`, { params, headers: headers() })
  return res.data
}

export async function getDataset(id) {
  const res = await axios.get(`/api/simulation/datasets/${id}`, { headers: headers() })
  return res.data
}

export async function createDataset(data) {
  const res = await axios.post(`/api/simulation/datasets`, data, { headers: headers() })
  return res.data
}

export async function updateDataset(id, data) {
  const res = await axios.put(`/api/simulation/datasets/${id}`, data, { headers: headers() })
  return res.data
}

export async function deleteDataset(id) {
  const res = await axios.delete(`/api/simulation/datasets/${id}`, { headers: headers() })
  return res.data
}

export async function createDatasetVersion(datasetId, data) {
  const res = await axios.post(`/api/simulation/datasets/${datasetId}/versions`, data, { headers: headers() })
  return res.data
}

export async function compareDatasetVersions(datasetId, data) {
  const res = await axios.post(`/api/simulation/datasets/${datasetId}/versions/compare`, data, { headers: headers() })
  return res.data
}

// ========== CI/CD 集成 API ==========
export async function getIntegrations(params = {}) {
  const res = await axios.get(`/api/simulation/integrations`, { params, headers: headers() })
  return res.data
}

export async function createIntegration(data) {
  const res = await axios.post(`/api/simulation/integrations`, data, { headers: headers() })
  return res.data
}

export async function updateIntegration(id, data) {
  const res = await axios.put(`/api/simulation/integrations/${id}`, data, { headers: headers() })
  return res.data
}

export async function deleteIntegration(id) {
  const res = await axios.delete(`/api/simulation/integrations/${id}`, { headers: headers() })
  return res.data
}
