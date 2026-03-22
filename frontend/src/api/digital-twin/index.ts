/**
 * 宠物数字孪生 API
 * Sprint 18 - 宠物数字孪生前端
 */

import axios from 'axios'

const BASE_URL = '/api/v1/digital-twin'

function getToken() {
  return localStorage.getItem('token') || ''
}

function headers() {
  return { Authorization: `Bearer ${getToken()}` }
}

// ========== 生命体征 ==========
export interface VitalData {
  device_id: string
  heart_rate: number        // 心率 bpm
  respiratory_rate: number   // 呼吸频率 次/分
  body_temp: number          // 体温 ℃
  activity_level: number     // 活动量 0-100
  is_abnormal: boolean
  abnormal_items: string[]
  timestamp: string
}

export interface VitalTrend {
  time: string
  heart_rate: number
  respiratory_rate: number
  body_temp: number
  activity_level: number
}

// ========== 历史事件 ==========
export interface HistoryEvent {
  id: number
  device_id: string
  event_type: string         // eating, drinking, sleeping, playing,异常
  event_name: string
  description: string
  occurred_at: string
  duration?: number          // 持续时间(秒)
  severity?: string          // normal, warning, critical
}

// ========== 行为预测 ==========
export interface BehaviorPrediction {
  device_id: string
  current_state: {
    state: string            // awake, sleeping, eating, playing
    confidence: number
    started_at: string
  }
  short_term_predictions: Array<{
    behavior: string
    probability: number
    expected_time: string
    duration_estimate: number
  }>
  intent_recognition: {
    intent: string           // hunger, thirst, play, rest, seek_attention
    confidence: number
    suggested_action: string
  }
  updated_at: string
}

// ========== API 方法 ==========

const digitalTwinApi = {
  // 获取当前生命体征
  getCurrentVitals: (deviceId: string) => {
    return axios.get(`${BASE_URL}/vitals/current/${deviceId}`, { headers: headers() }).then(r => r.data)
  },

  // 获取生命体征趋势（最近24小时）
  getVitalsTrend: (deviceId: string, params?: { hours?: number }) => {
    return axios.get(`${BASE_URL}/vitals/trend/${deviceId}`, { params, headers: headers() }).then(r => r.data)
  },

  // 获取历史事件
  getHistoryEvents: (deviceId: string, params?: {
    start_date?: string
    end_date?: string
    event_type?: string
    page?: number
    page_size?: number
  }) => {
    return axios.get(`${BASE_URL}/history/events/${deviceId}`, { params, headers: headers() }).then(r => r.data)
  },

  // 获取行为预测
  getBehaviorPrediction: (deviceId: string) => {
    return axios.get(`${BASE_URL}/behavior/prediction/${deviceId}`, { headers: headers() }).then(r => r.data)
  },

  // 获取宠物列表（用于选择）
  getPets: (params?: { page?: number; page_size?: number }) => {
    return axios.get(`${BASE_URL}/pets`, { params, headers: headers() }).then(r => r.data)
  }
}

export default digitalTwinApi
