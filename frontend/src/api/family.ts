import request from './request'

// Family Members
export const getFamilyMembers = (params?: any) =>
  request.get('/api/v1/family/members', { params })

export const inviteFamilyMember = (data: any) =>
  request.post('/api/v1/family/members/invite', data)

export const updateMemberRole = (id: number, data: any) =>
  request.put(`/api/v1/family/members/${id}/role`, data)

export const removeFamilyMember = (id: number) =>
  request.delete(`/api/v1/family/members/${id}`)

// Child Mode
export const getChildModes = (params?: any) =>
  request.get('/api/v1/family/child-mode', { params })

export const createChildMode = (data: any) =>
  request.post('/api/v1/family/child-mode', data)

export const updateChildMode = (id: number, data: any) =>
  request.put(`/api/v1/family/child-mode/${id}`, data)

export const deleteChildMode = (id: number) =>
  request.delete(`/api/v1/family/child-mode/${id}`)

// Elder Mode
export const getElderModes = (params?: any) =>
  request.get('/api/v1/family/elder-mode', { params })

export const createElderMode = (data: any) =>
  request.post('/api/v1/family/elder-mode', data)

export const updateElderMode = (id: number, data: any) =>
  request.put(`/api/v1/family/elder-mode/${id}`, data)

export const deleteElderMode = (id: number) =>
  request.delete(`/api/v1/family/elder-mode/${id}`)

// Family Album
export const getFamilyAlbum = (params?: any) =>
  request.get('/api/v1/family/album', { params })

export const uploadPhoto = (file: File) => {
  const formData = new FormData()
  formData.append('file', file)
  return request.post('/api/v1/family/album/upload', formData, {
    headers: { 'Content-Type': 'multipart/form-data' }
  })
}

export const deletePhoto = (id: number) =>
  request.delete(`/api/v1/family/album/${id}`)

// Family Settings
export const getFamilySettings = () =>
  request.get('/api/v1/family/settings')

export const updateBasicSettings = (data: any) =>
  request.put('/api/v1/family/settings/basic', data)

export const updateNotificationSettings = (data: any) =>
  request.put('/api/v1/family/settings/notification', data)

export const updatePrivacySettings = (data: any) =>
  request.put('/api/v1/family/settings/privacy', data)
