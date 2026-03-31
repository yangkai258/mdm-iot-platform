// Permissions API
const API_BASE = '/api/v1'

export const getRoles = (params) => {
  const token = localStorage.getItem('token')
  return fetch(`${API_BASE}/roles?${new URLSearchParams(params)}`, {
    headers: { 'Authorization': `Bearer ${token}` }
  }).then(r => r.json())
}

export const createRole = (data) => {
  const token = localStorage.getItem('token')
  return fetch(`${API_BASE}/roles`, {
    method: 'POST',
    headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' },
    body: JSON.stringify(data)
  }).then(r => r.json())
}

export const updateRole = (id, data) => {
  const token = localStorage.getItem('token')
  return fetch(`${API_BASE}/roles/${id}`, {
    method: 'PUT',
    headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' },
    body: JSON.stringify(data)
  }).then(r => r.json())
}

export const deleteRole = (id) => {
  const token = localStorage.getItem('token')
  return fetch(`${API_BASE}/roles/${id}`, {
    method: 'DELETE',
    headers: { 'Authorization': `Bearer ${token}` }
  }).then(r => r.json())
}

export const getAllPermissions = () => {
  const token = localStorage.getItem('token')
  return fetch(`${API_BASE}/permissions`, {
    headers: { 'Authorization': `Bearer ${token}` }
  }).then(r => r.json())
}

export const getRolePermissions = (id) => {
  const token = localStorage.getItem('token')
  return fetch(`${API_BASE}/roles/${id}/permissions`, {
    headers: { 'Authorization': `Bearer ${token}` }
  }).then(r => r.json())
}

export const assignPermissions = (id, data) => {
  const token = localStorage.getItem('token')
  return fetch(`${API_BASE}/roles/${id}/permissions`, {
    method: 'POST',
    headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' },
    body: JSON.stringify(data)
  }).then(r => r.json())
}
