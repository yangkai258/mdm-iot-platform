import request from './request'

// Family Members
export const getFamilyMembers = (params?: any) =>
  request.get('/api/family/members', { params })

export const inviteFamilyMember = (data: any) =>
  request.post('/api/family/members/invite', data)

export const updateMemberRole = (id: number, data: any) =>
  request.put('/api/family/members/${id}/role`, data)

export const removeFamilyMember = (id: number) =>
  request.delete('/api/family/members/${id}`)

// Child Mode
export const getChildModes = (params?: any) =>
  request.get('/api/family/child-mode', { params })

export const createChildMode = (data: any) =>
  request.post('/api/family/child-mode', data)

export const updateChildMode = (id: number, data: any) =>
  request.put('/api/family/child-mode/${id}`, data)

export const deleteChildMode = (id: number) =>
  request.delete('/api/family/child-mode/${id}`)

// Elder Mode
export const getElderModes = (params?: any) =>
  request.get('/api/family/elder-mode', { params })

export const createElderMode = (data: any) =>
  request.post('/api/family/elder-mode', data)

export const updateElderMode = (id: number, data: any) =>
  request.put('/api/family/elder-mode/${id}`, data)

export const deleteElderMode = (id: number) =>
  request.delete('/api/family/elder-mode/${id}`)

// Family Album
export const getFamilyAlbum = (params?: any) =>
  request.get('/api/family/album', { params })

export const uploadPhoto = (file: File) => {
  const formData = new FormData()
  formData.append('file', file)
  return request.post('/api/family/album/upload', formData, {
    headers: { 'Content-Type': 'multipart/form-data' }
  })
}

export const deletePhoto = (id: number) =>
  request.delete('/api/family/album/${id}`)

// Family Settings
export const getFamilySettings = () =>
  request.get('/api/family/settings')

export const updateBasicSettings = (data: any) =>
  request.put('/api/family/settings/basic', data)

export const updateNotificationSettings = (data: any) =>
  request.put('/api/family/settings/notification', data)

export const updatePrivacySettings = (data: any) =>
  request.put('/api/family/settings/privacy', data)
