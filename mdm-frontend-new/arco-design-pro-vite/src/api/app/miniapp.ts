/**
 * 微信小程序 H5 API - Sprint 22
 * API 前缀: /api/v1/miniapp
 */
import axios from 'axios'

const BASE_URL = '/api/miniapp'

function getToken() {
  return localStorage.getItem('token') || ''
}

function headers() {
  return { Authorization: `Bearer ${getToken()}` }
}

// ========== 首页设备列表 ==========
export async function getHomeDevices(params = {}) {
  const res = await axios.get(`/api/home/devices`, { params, headers: headers() })
  return res.data
}

// ========== 设备状态 ==========
export async function getDeviceStatus(id) {
  const res = await axios.get(`/api/devices/${id}/status`, { headers: headers() })
  return res.data
}

// ========== 设备控制 ==========
export async function controlDevice(id, data) {
  const res = await axios.post(`/api/devices/${id}/control`, data, { headers: headers() })
  return res.data
}

// ========== 语音指令 ==========
export async function sendVoiceCommand(id, data) {
  const res = await axios.post(`/api/devices/${id}/voice`, data, { headers: headers() })
  return res.data
}

// ========== 首页状态概览 ==========
export async function getHomeOverview() {
  const res = await axios.get(`/api/home/overview`, { headers: headers() })
  return res.data
}

// ========== 设备快速操作 ==========
export async function quickAction(id, action) {
  const res = await axios.post(`/api/devices/${id}/quick-action`, { action }, { headers: headers() })
  return res.data
}
