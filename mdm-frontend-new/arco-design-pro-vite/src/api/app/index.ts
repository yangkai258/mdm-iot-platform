/**
 * App 端 API - Sprint 22
 * API 前缀: /api/v1/app
 */
import axios from 'axios'

const BASE_URL = '/api/app'

function getToken() {
  return localStorage.getItem('token') || ''
}

function headers() {
  return { Authorization: `Bearer ${getToken()}` }
}

// ========== 设备列表 ==========
export async function getDeviceList(params = {}) {
  const res = await axios.get(`/api/devices`, { params, headers: headers() })
  return res.data
}

// ========== 设备详情 ==========
export async function getDevice(id) {
  const res = await axios.get(`/api/devices/${id}`, { headers: headers() })
  return res.data
}

// ========== 设备状态 ==========
export async function getDeviceStatus(id) {
  const res = await axios.get(`/api/devices/${id}/status`, { headers: headers() })
  return res.data
}

// ========== 设备控制指令 ==========
export async function sendDeviceCommand(id, data) {
  const res = await axios.post(`/api/devices/${id}/commands`, data, { headers: headers() })
  return res.data
}

// ========== 设备快速操作 ==========
export async function quickAction(id, action) {
  const res = await axios.post(`/api/devices/${id}/quick-action`, { action }, { headers: headers() })
  return res.data
}

// ========== 设备心率/状态概览 ==========
export async function getDeviceOverview() {
  const res = await axios.get(`/api/overview`, { headers: headers() })
  return res.data
}
