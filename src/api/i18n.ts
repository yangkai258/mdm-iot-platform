import axios from 'axios'

const BASE_URL = '/api/v1/i18n'

function getToken() {
  return localStorage.getItem('token') || ''
}

function headers() {
  return { Authorization: `Bearer ${getToken()}` }
}

// ============================================================
// 翻译管理
// ============================================================

export interface TranslationItem {
  id?: number
  key: string
  locale: string
  value: string
  description?: string
  tags?: string[]
  updated_at?: string
  created_at?: string
}

export interface TranslationCreate {
  key: string
  locale: string
  value: string
  description?: string
  tags?: string[]
}

export interface TranslationUpdate {
  key?: string
  locale?: string
  value?: string
  description?: string
  tags?: string[]
}

export interface Language {
  code: string
  name: string
  native_name: string
  is_default: boolean
  is_active: boolean
  translation_count: number
}

// 获取翻译列表
export async function getTranslations(params?: {
  locale?: string
  key?: string
  search?: string
  page?: number
  page_size?: number
}) {
  const res = await axios.get(`${BASE_URL}/translations`, { params, headers: headers() })
  return res.data
}

// 获取单条翻译
export async function getTranslation(id: number) {
  const res = await axios.get(`${BASE_URL}/translations/${id}`, { headers: headers() })
  return res.data
}

// 创建翻译
export async function createTranslation(data: TranslationCreate) {
  const res = await axios.post(`${BASE_URL}/translations`, data, { headers: headers() })
  return res.data
}

// 更新翻译
export async function updateTranslation(id: number, data: TranslationUpdate) {
  const res = await axios.put(`${BASE_URL}/translations/${id}`, data, { headers: headers() })
  return res.data
}

// 删除翻译
export async function deleteTranslation(id: number) {
  const res = await axios.delete(`${BASE_URL}/translations/${id}`, { headers: headers() })
  return res.data
}

// 批量导入翻译
export async function importTranslations(data: { locale: string; translations: TranslationCreate[] }) {
  const res = await axios.post(`${BASE_URL}/translations/import`, data, { headers: headers() })
  return res.data
}

// 导出翻译
export async function exportTranslations(params?: { locale?: string }) {
  const res = await axios.get(`${BASE_URL}/translations/export`, { params, headers: headers() })
  return res.data
}

// ============================================================
// 语言管理
// ============================================================

// 获取语言列表
export async function getLanguages() {
  const res = await axios.get(`${BASE_URL}/languages`, { headers: headers() })
  return res.data
}

// 创建语言
export async function createLanguage(data: { code: string; name: string; native_name: string; is_default?: boolean }) {
  const res = await axios.post(`${BASE_URL}/languages`, data, { headers: headers() })
  return res.data
}

// 更新语言
export async function updateLanguage(code: string, data: Partial<Language>) {
  const res = await axios.put(`${BASE_URL}/languages/${code}`, data, { headers: headers() })
  return res.data
}

// 删除语言
export async function deleteLanguage(code: string) {
  const res = await axios.delete(`${BASE_URL}/languages/${code}`, { headers: headers() })
  return res.data
}

// 获取翻译键列表（所有唯一的 key）
export async function getTranslationKeys() {
  const res = await axios.get(`${BASE_URL}/keys`, { headers: headers() })
  return res.data
}

// ============================================================
// 区域设置（独立 API 前缀 /api/v1/regions）
// ============================================================

export interface RegionSetting {
  id?: number
  region_key: string
  region_name: string
  locale_default: string
  locale_supported: string[]
  timezone: string
  date_format: string
  time_format: string
  number_format: string
  currency_code: string
  is_active: boolean
  updated_at?: string
  created_at?: string
}

export async function getRegionSettings(params?: any) {
  const res = await axios.get('/api/v1/regions/settings', { params, headers: headers() })
  return res.data
}

export async function getRegionSetting(id: number) {
  const res = await axios.get(`/api/v1/regions/settings/${id}`, { headers: headers() })
  return res.data
}

export async function createRegionSetting(data: RegionSetting) {
  const res = await axios.post('/api/v1/regions/settings', data, { headers: headers() })
  return res.data
}

export async function updateRegionSetting(id: number, data: Partial<RegionSetting>) {
  const res = await axios.put(`/api/v1/regions/settings/${id}`, data, { headers: headers() })
  return res.data
}

export async function deleteRegionSetting(id: number) {
  const res = await axios.delete(`/api/v1/regions/settings/${id}`, { headers: headers() })
  return res.data
}
