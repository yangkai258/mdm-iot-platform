import axios from 'axios'

const BASE_URL = '/api/v1'
const token = localStorage.getItem('token') || ''
const headers = { Authorization: `Bearer ${token}` }

// ========== Child Mode ==========
export async function getChildModes(params = {}) {
  const res = await axios.get(`${BASE_URL}/advanced/child-mode`, { params, headers })
  return res.data
}

export async function createChildMode(data) {
  const res = await axios.post(`${BASE_URL}/advanced/child-mode`, data, { headers })
  return res.data
}

export async function updateChildMode(id, data) {
  const res = await axios.put(`${BASE_URL}/advanced/child-mode/${id}`, data, { headers })
  return res.data
}

export async function deleteChildMode(id) {
  const res = await axios.delete(`${BASE_URL}/advanced/child-mode/${id}`, { headers })
  return res.data
}

export async function getChildModeUsageReport(params = {}) {
  const res = await axios.get(`${BASE_URL}/advanced/child-mode/usage-report`, { params, headers })
  return res.data
}

// ========== Elder Mode ==========
export async function getElderModes(params = {}) {
  const res = await axios.get(`${BASE_URL}/advanced/elder-mode`, { params, headers })
  return res.data
}

export async function createElderMode(data) {
  const res = await axios.post(`${BASE_URL}/advanced/elder-mode`, data, { headers })
  return res.data
}

export async function updateElderMode(id, data) {
  const res = await axios.put(`${BASE_URL}/advanced/elder-mode/${id}`, data, { headers })
  return res.data
}

export async function deleteElderMode(id) {
  const res = await axios.delete(`${BASE_URL}/advanced/elder-mode/${id}`, { headers })
  return res.data
}

export async function getElderHealthData(params = {}) {
  const res = await axios.get(`${BASE_URL}/advanced/elder-mode/health-data`, { params, headers })
  return res.data
}

// ========== Family Album ==========
export async function getAlbumPhotos(params = {}) {
  const res = await axios.get(`${BASE_URL}/advanced/album`, { params, headers })
  return res.data
}

export async function uploadAlbumPhoto(file, albumId) {
  const formData = new FormData()
  formData.append('file', file)
  if (albumId) formData.append('album_id', String(albumId))
  const res = await axios.post(`${BASE_URL}/advanced/album/upload`, formData, {
    headers: { ...headers, 'Content-Type': 'multipart/form-data' }
  })
  return res.data
}

export async function deleteAlbumPhoto(id) {
  const res = await axios.delete(`${BASE_URL}/advanced/album/${id}`, { headers })
  return res.data
}

export async function sharePhoto(data) {
  const res = await axios.post(`${BASE_URL}/advanced/album/share`, data, { headers })
  return res.data
}

export async function getAlbums(params = {}) {
  const res = await axios.get(`${BASE_URL}/advanced/album/albums`, { params, headers })
  return res.data
}

export async function createAlbum(data) {
  const res = await axios.post(`${BASE_URL}/advanced/album/albums`, data, { headers })
  return res.data
}

export async function updateAlbum(id, data) {
  const res = await axios.put(`${BASE_URL}/advanced/album/albums/${id}`, data, { headers })
  return res.data
}

export async function deleteAlbum(id) {
  const res = await axios.delete(`${BASE_URL}/advanced/album/albums/${id}`, { headers })
  return res.data
}

// ========== Pet Finder ==========
export async function getPetFinderReports(params = {}) {
  const res = await axios.get(`${BASE_URL}/advanced/pet-finder/reports`, { params, headers })
  return res.data
}

export async function createPetFinderReport(data) {
  const res = await axios.post(`${BASE_URL}/advanced/pet-finder/reports`, data, { headers })
  return res.data
}

export async function updatePetFinderReport(id, data) {
  const res = await axios.put(`${BASE_URL}/advanced/pet-finder/reports/${id}`, data, { headers })
  return res.data
}

export async function deletePetFinderReport(id) {
  const res = await axios.delete(`${BASE_URL}/advanced/pet-finder/reports/${id}`, { headers })
  return res.data
}

export async function sharePetFinderReport(reportId) {
  const res = await axios.post(`${BASE_URL}/advanced/pet-finder/reports/share`, { report_id: reportId }, { headers })
  return res.data
}

export async function getPetFinderReportDetail(id) {
  const res = await axios.get(`${BASE_URL}/advanced/pet-finder/reports/${id}`, { headers })
  return res.data
}

// ========== Vaccination ==========
export async function getVaccinationRecords(params = {}) {
  const res = await axios.get(`${BASE_URL}/advanced/vaccination`, { params, headers })
  return res.data
}

export async function createVaccinationRecord(data) {
  const res = await axios.post(`${BASE_URL}/advanced/vaccination`, data, { headers })
  return res.data
}

export async function updateVaccinationRecord(id, data) {
  const res = await axios.put(`${BASE_URL}/advanced/vaccination/${id}`, data, { headers })
  return res.data
}

export async function deleteVaccinationRecord(id) {
  const res = await axios.delete(`${BASE_URL}/advanced/vaccination/${id}`, { headers })
  return res.data
}

export async function getVaccinationReminders(params = {}) {
  const res = await axios.get(`${BASE_URL}/advanced/vaccination/reminders`, { params, headers })
  return res.data
}

export async function setVaccinationReminder(data) {
  const res = await axios.post(`${BASE_URL}/advanced/vaccination/reminders`, data, { headers })
  return res.data
}

export async function updateVaccinationReminder(id, data) {
  const res = await axios.put(`${BASE_URL}/advanced/vaccination/reminders/${id}`, data, { headers })
  return res.data
}

export async function deleteVaccinationReminder(id) {
  const res = await axios.delete(`${BASE_URL}/advanced/vaccination/reminders/${id}`, { headers })
  return res.data
}

// ========== Diet Record ==========
export async function getDietRecords(params = {}) {
  const res = await axios.get(`${BASE_URL}/advanced/diet`, { params, headers })
  return res.data
}

export async function createDietRecord(data) {
  const res = await axios.post(`${BASE_URL}/advanced/diet`, data, { headers })
  return res.data
}

export async function updateDietRecord(id, data) {
  const res = await axios.put(`${BASE_URL}/advanced/diet/${id}`, data, { headers })
  return res.data
}

export async function deleteDietRecord(id) {
  const res = await axios.delete(`${BASE_URL}/advanced/diet/${id}`, { headers })
  return res.data
}

export async function getDietSummary(params = {}) {
  const res = await axios.get(`${BASE_URL}/advanced/diet/summary`, { params, headers })
  return res.data
}
