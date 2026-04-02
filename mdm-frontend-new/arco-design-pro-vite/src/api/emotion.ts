import axios from 'axios'

const BASE_URL = '/api'
const getHeaders = () => ({ Authorization: `Bearer ${localStorage.getItem('token') || ''}` })

async function request(method: string, url: string, data?: any, params?: any) {
  const res = await axios({ method, url, data, params, headers: getHeaders() })
  return res.data
}

// ============================================================
// 情绪日志
// ============================================================
export const getEmotionLogs = (params?: any) =>
  request('get', `/api/emotion/logs`, null, params)

// ============================================================
// 情绪识别配置
// ============================================================
export const getEmotionRecognizeConfig = () =>
  request('get', `/api/emotion/recognize-config`)

export const updateEmotionRecognizeConfig = (data: any) =>
  request('put', `/api/emotion/recognize-config`, data)

// ============================================================
// 情绪报告
// ============================================================
export const getEmotionReports = (params?: any) =>
  request('get', `/api/emotion/reports`, null, params)

// ============================================================
// 情绪响应配置
// ============================================================
export const getEmotionResponseConfigs = (params?: any) =>
  request('get', `/api/emotion/response-config`, null, params)

export const createEmotionResponseConfig = (data: any) =>
  request('post', `/api/emotion/response-config`, data)

export const updateEmotionResponseConfig = (id: number, data: any) =>
  request('put', `/api/emotion/response-config/${id}`, data)

export const deleteEmotionResponseConfig = (id: number) =>
  request('delete', `/api/emotion/response-config/${id}`)

// ============================================================
// 情绪趋势
// ============================================================
export const getEmotionTrendStats = (params?: any) =>
  request('get', `/api/emotions/records/stats`, null, params)

// ============================================================
// 家庭情绪地图
// ============================================================
export const getFamilyEmotionMap = (params?: any) =>
  request('get', `/api/emotions/family-map`, null, params)

// ============================================================
// 语音情绪
// ============================================================
export const getVoiceEmotionRecords = (params?: any) =>
  request('get', `/api/voice-emotion/records`, null, params)

export const analyzeVoiceEmotion = (data: any) =>
  request('post', `/api/voice-emotion/analyze`, data)
