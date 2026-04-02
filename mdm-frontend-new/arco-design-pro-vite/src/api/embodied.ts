import axios from 'axios'

const BASE_URL = '/api'

function getToken() {
  return localStorage.getItem('token') || ''
}

function headers() {
  return { Authorization: `Bearer ${getToken()}` }
}

// ========== 环境感知 API ==========
export interface PerceptionResult {
  id: number
  device_id: string
  visual: {
    objects: Array<{ label: string; confidence: number; bbox?: number[] }>
    scene: string
    human_pose?: { keypoints: Array<{ x: number; y: number; score: number }> }
    gesture?: string
  }
  depth: {
    obstacles: Array<{ distance: number; direction: string }>
    confidence: number
  }
  touch: {
    touched: boolean
    position?: string
    force?: number
  }
  confidence: number
  created_at: string
}

export interface PerceptionEvent {
  id: number
  device_id: string
  event_type: string
  description: string
  data: Record<string, any>
  created_at: string
}

export async function getPerception(deviceId: string, params = {}) {
  const res = await axios.get(`/api/embodied/${deviceId}/perception`, { params, headers: headers() })
  return res.data
}

export async function reportVisualPerception(deviceId: string, data: any) {
  const res = await axios.post(`/api/embodied/${deviceId}/perception/visual`, data, { headers: headers() })
  return res.data
}

export async function reportDepthPerception(deviceId: string, data: any) {
  const res = await axios.post(`/api/embodied/${deviceId}/perception/depth`, data, { headers: headers() })
  return res.data
}

export async function reportTouchPerception(deviceId: string, data: any) {
  const res = await axios.post(`/api/embodied/${deviceId}/perception/touch`, data, { headers: headers() })
  return res.data
}

// ========== 地图管理 API ==========
export interface EmbodiedMap {
  id: number
  device_id: string
  map_type: 'grid' | 'semantic' | 'topological'
  map_data: any
  resolution: number
  size: { width: number; height: number }
  version: number
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface MapAnnotation {
  id: number
  map_id: number
  annotation_type: string
  position: { x: number; y: number }
  label: string
  description?: string
}

export async function getMaps(deviceId: string, params = {}) {
  const res = await axios.get(`/api/embodied/${deviceId}/map`, { params, headers: headers() })
  return res.data
}

export async function updateMap(deviceId: string, data: any) {
  const res = await axios.post(`/api/embodied/${deviceId}/map/update`, data, { headers: headers() })
  return res.data
}

export async function getLocalization(deviceId: string) {
  const res = await axios.get(`/api/embodied/${deviceId}/localization`, { headers: headers() })
  return res.data
}

export async function calibrateLocalization(deviceId: string, data: any) {
  const res = await axios.post(`/api/embodied/${deviceId}/localization/calibrate`, data, { headers: headers() })
  return res.data
}

// ========== 导航控制 API ==========
export interface NavigationTarget {
  id: number
  device_id: string
  target_x: number
  target_y: number
  status: 'pending' | 'navigating' | 'arrived' | 'failed'
  path?: Array<{ x: number; y: number }>
  started_at?: string
  arrived_at?: string
}

export async function navigateTo(deviceId: string, data: { target_x: number; target_y: number }) {
  const res = await axios.post(`/api/embodied/${deviceId}/navigate`, data, { headers: headers() })
  return res.data
}

export async function stopNavigation(deviceId: string) {
  const res = await axios.post(`/api/embodied/${deviceId}/stop`, {}, { headers: headers() })
  return res.data
}

export async function startFollow(deviceId: string, data: { target_id: string }) {
  const res = await axios.post(`/api/embodied/${deviceId}/follow`, data, { headers: headers() })
  return res.data
}

export async function stopFollow(deviceId: string) {
  const res = await axios.post(`/api/embodied/${deviceId}/follow/stop`, {}, { headers: headers() })
  return res.data
}

export async function getExploreStatus(deviceId: string) {
  const res = await axios.get(`/api/embodied/${deviceId}/explore/status`, { headers: headers() })
  return res.data
}

export async function startExplore(deviceId: string) {
  const res = await axios.post(`/api/embodied/${deviceId}/explore/start`, {}, { headers: headers() })
  return res.data
}

export async function getNavigationHistory(deviceId: string, params = {}) {
  const res = await axios.get(`/api/embodied/${deviceId}/navigation/history`, { params, headers: headers() })
  return res.data
}

// ========== 动作库 API ==========
export interface ActionItem {
  id: number
  action_name: string
  action_type: 'built-in' | 'learned' | 'custom'
  description?: string
  duration_ms?: number
  difficulty?: 'easy' | 'medium' | 'hard'
  tags?: string[]
  motion_data?: any
  video_url?: string
  thumbnail_url?: string
  score?: number
  creator_id?: number
  is_public?: boolean
  downloads?: number
  created_at: string
}

export interface ActionExecution {
  id: number
  device_id: string
  action_id: number
  execution_type: 'triggered' | 'scheduled' | 'manual'
  start_time: string
  end_time?: string
  status: 'running' | 'completed' | 'interrupted' | 'failed'
  parameters?: any
  interruption_reason?: string
}

export async function getActionLibrary(params = {}) {
  const res = await axios.get(`/api/embodied/action-library`, { params, headers: headers() })
  return res.data
}

export async function getActionDetail(actionId: number) {
  const res = await axios.get(`/api/embodied/action-library/${actionId}`, { headers: headers() })
  return res.data
}

export async function recordAction(data: any) {
  const res = await axios.post(`/api/embodied/action-library/record`, data, { headers: headers() })
  return res.data
}

export async function learnAction(actionId: number, data: any) {
  const res = await axios.post(`/api/embodied/action-library/${actionId}/learn`, data, { headers: headers() })
  return res.data
}

export async function executeAction(deviceId: string, data: any) {
  const res = await axios.post(`/api/embodied/${deviceId}/action/execute`, data, { headers: headers() })
  return res.data
}

export async function stopAction(deviceId: string) {
  const res = await axios.post(`/api/embodied/${deviceId}/action/stop`, {}, { headers: headers() })
  return res.data
}

export async function shareAction(actionId: number) {
  const res = await axios.post(`/api/embodied/action-library/${actionId}/share`, {}, { headers: headers() })
  return res.data
}

export async function deleteAction(actionId: number) {
  const res = await axios.delete(`/api/embodied/action-library/${actionId}`, { headers: headers() })
  return res.data
}

export async function getActionExecutions(deviceId: string, params = {}) {
  const res = await axios.get(`/api/embodied/${deviceId}/action/executions`, { params, headers: headers() })
  return res.data
}

// ========== 决策引擎 API ==========
export interface DecisionContext {
  device_id: string
  user_state: {
    emotion?: string
    pose?: string
    position?: { x: number; y: number }
  }
  environment_state: {
    location?: string
    time?: string
    scene?: string
  }
  pet_state: {
    emotion?: string
    battery?: number
    mode?: string
  }
  task_state?: {
    progress?: number
    current_task?: string
  }
}

export interface DecisionLog {
  id: number
  device_id: string
  decision_type: string
  context: DecisionContext
  chosen_action: string
  action_params?: any
  confidence: number
  reasoning?: string
  execution_result?: string
  latency_ms?: number
  decided_at: string
}

export interface DecisionStrategy {
  id: number
  device_id: string
  strategy_type: 'safety_first' | 'task_oriented' | 'interaction_first' | 'exploration'
  is_active: boolean
  config?: any
  created_at: string
}

export async function getDecisionContext(deviceId: string) {
  const res = await axios.get(`/api/embodied/${deviceId}/decision/context`, { headers: headers() })
  return res.data
}

export async function setDecisionStrategy(deviceId: string, data: any) {
  const res = await axios.post(`/api/embodied/${deviceId}/decision/strategy`, data, { headers: headers() })
  return res.data
}

export async function getDecisionLogs(deviceId: string, params = {}) {
  const res = await axios.get(`/api/embodied/${deviceId}/decision/logs`, { params, headers: headers() })
  return res.data
}

export async function getDecisionStrategies(deviceId: string) {
  const res = await axios.get(`/api/embodied/${deviceId}/decision/strategies`, { headers: headers() })
  return res.data
}

// ========== 安全边界 API ==========
export interface SafetyZone {
  id: number
  device_id: string
  zone_type: 'forbidden' | 'caution' | 'safe'
  zone_shape: 'rectangle' | 'circle' | 'polygon'
  zone_data: any
  zone_name: string
  is_enabled: boolean
  created_by?: number
  created_at: string
}

export interface SafetyLog {
  id: number
  device_id: string
  event_type: 'collision' | 'emergency_stop' | 'zone_violation' | 'fall_prevention'
  severity: 'info' | 'warning' | 'critical'
  details?: any
  location?: { x: number; y: number }
  resolved: boolean
  resolved_at?: string
  created_at: string
}

export async function getSafetyZones(deviceId: string, params = {}) {
  const res = await axios.get(`/api/embodied/${deviceId}/safety/zones`, { params, headers: headers() })
  return res.data
}

export async function createSafetyZone(deviceId: string, data: any) {
  const res = await axios.post(`/api/embodied/${deviceId}/safety/zones`, data, { headers: headers() })
  return res.data
}

export async function updateSafetyZone(deviceId: string, zoneId: number, data: any) {
  const res = await axios.put(`/api/embodied/${deviceId}/safety/zones/${zoneId}`, data, { headers: headers() })
  return res.data
}

export async function deleteSafetyZone(deviceId: string, zoneId: number) {
  const res = await axios.delete(`/api/embodied/${deviceId}/safety/zones/${zoneId}`, { headers: headers() })
  return res.data
}

export async function emergencyStop(deviceId: string) {
  const res = await axios.post(`/api/embodied/${deviceId}/safety/emergency-stop`, {}, { headers: headers() })
  return res.data
}

export async function getSafetyLogs(deviceId: string, params = {}) {
  const res = await axios.get(`/api/embodied/${deviceId}/safety/logs`, { params, headers: headers() })
  return res.data
}
