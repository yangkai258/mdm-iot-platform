/**
 * 设备监控 API
 * 设备监控、传感器、日志、动作库、批量操作
 */

const BASE_URL = '/api/v1'

export interface DeviceMetric {
  device_id: string
  metric_type: string
  metric_name: string
  metric_value: number
  unit: string
  timestamp: string
}

export interface SensorEvent {
  id: number
  event_id: string
  device_id: string
  sensor_type: string
  sensor_value: number
  unit: string
  is_abnormal: boolean
  event_type: string
  created_at: string
}

export interface BatchTask {
  task_id: string
  total: number
  success: number
  failed: number
  status: string
  results: Array<{device_id: string; status: string; error?: string}>
}

const monitorApi = {
  // 设备监控
  getDeviceMetrics: (deviceId: string, params?: {start_time?: string; end_time?: string; metric_type?: string}) => {
    const token = localStorage.getItem('token')
    const qs = params ? '?' + new URLSearchParams(params as any).toString() : ''
    return fetch(`${BASE_URL}/monitor/devices/${deviceId}/metrics${qs}`, {
      headers: { 'Authorization': `Bearer ${token}` }
    }).then(r => r.json())
  },
  getDeviceRealtime: (deviceId: string) => {
    const token = localStorage.getItem('token')
    return fetch(`${BASE_URL}/monitor/devices/${deviceId}/realtime`, {
      headers: { 'Authorization': `Bearer ${token}` }
    }).then(r => r.json())
  },
  getDeviceHistory: (deviceId: string, params?: {start_time?: string; end_time?: string}) => {
    const token = localStorage.getItem('token')
    const qs = params ? '?' + new URLSearchParams(params as any).toString() : ''
    return fetch(`${BASE_URL}/monitor/devices/${deviceId}/history${qs}`, {
      headers: { 'Authorization': `Bearer ${token}` }
    }).then(r => r.json())
  },

  // 传感器
  getSensorEvents: (deviceId: string, params?: {page?: number; page_size?: number}) => {
    const token = localStorage.getItem('token')
    const qs = params ? '?' + new URLSearchParams(params as any).toString() : ''
    return fetch(`${BASE_URL}/sensors/${deviceId}/events${qs}`, {
      headers: { 'Authorization': `Bearer ${token}` }
    }).then(r => r.json())
  },
  reportSensorEvent: (deviceId: string, data: Partial<SensorEvent>) => {
    const token = localStorage.getItem('token')
    return fetch(`${BASE_URL}/sensors/${deviceId}/events`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(data)
    }).then(r => r.json())
  },
  getSensorLatest: (deviceId: string) => {
    const token = localStorage.getItem('token')
    return fetch(`${BASE_URL}/sensors/${deviceId}/latest`, {
      headers: { 'Authorization': `Bearer ${token}` }
    }).then(r => r.json())
  },

  // 设备日志
  getDeviceLogs: (params?: {device_id?: string; level?: string; keyword?: string; page?: number; page_size?: number}) => {
    const token = localStorage.getItem('token')
    const qs = params ? '?' + new URLSearchParams(params as any).toString() : ''
    return fetch(`${BASE_URL}/device-logs${qs}`, {
      headers: { 'Authorization': `Bearer ${token}` }
    }).then(r => r.json())
  },

  // 动作库
  getActionLibrary: (params?: {category?: string; page?: number; page_size?: number}) => {
    const token = localStorage.getItem('token')
    const qs = params ? '?' + new URLSearchParams(params as any).toString() : ''
    return fetch(`${BASE_URL}/action-library${qs}`, {
      headers: { 'Authorization': `Bearer ${token}` }
    }).then(r => r.json())
  },
  getActionDetail: (id: number) => {
    const token = localStorage.getItem('token')
    return fetch(`${BASE_URL}/action-library/${id}`, {
      headers: { 'Authorization': `Bearer ${token}` }
    }).then(r => r.json())
  },
  createAction: (data: any) => {
    const token = localStorage.getItem('token')
    return fetch(`${BASE_URL}/action-library`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(data)
    }).then(r => r.json())
  },
  updateAction: (id: number, data: any) => {
    const token = localStorage.getItem('token')
    return fetch(`${BASE_URL}/action-library/${id}`, {
      method: 'PUT',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(data)
    }).then(r => r.json())
  },
  deleteAction: (id: number) => {
    const token = localStorage.getItem('token')
    return fetch(`${BASE_URL}/action-library/${id}`, {
      method: 'DELETE',
      headers: { 'Authorization': `Bearer ${token}` }
    }).then(r => r.json())
  },

  // 批量操作
  batchDeviceAction: (data: {device_ids: string[]; action: string; params?: object}) => {
    const token = localStorage.getItem('token')
    return fetch(`${BASE_URL}/batch/devices/actions`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(data)
    }).then(r => r.json())
  },
  getBatchTaskStatus: (taskId: string) => {
    const token = localStorage.getItem('token')
    return fetch(`${BASE_URL}/batch/tasks/${taskId}`, {
      headers: { 'Authorization': `Bearer ${token}` }
    }).then(r => r.json())
  },
  getBatchTaskHistory: (params?: {page?: number; page_size?: number}) => {
    const token = localStorage.getItem('token')
    const qs = params ? '?' + new URLSearchParams(params as any).toString() : ''
    return fetch(`${BASE_URL}/batch/tasks${qs}`, {
      headers: { 'Authorization': `Bearer ${token}` }
    }).then(r => r.json())
  },
}

export default monitorApi
