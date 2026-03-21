/**
 * 会员管理 API
 * Base: /api/v1/member
 */
const API_BASE = '/api/v1/member'

function getToken() {
  return localStorage.getItem('token') || ''
}

async function request(url, options = {}) {
  const res = await fetch(url, {
    ...options,
    headers: {
      'Authorization': `Bearer ${getToken()}`,
      'Content-Type': 'application/json',
      ...(options.headers || {})
    }
  })
  const data = await res.json()
  if (data.code !== 0) {
    throw new Error(data.message || '请求失败')
  }
  return data
}

// ─── 会员 CRUD ───────────────────────────────────────────────

/** 会员列表 */
export async function getMemberList(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`${API_BASE}/list?${qs}`)
}

/** 会员详情 */
export async function getMemberDetail(id) {
  return request(`${API_BASE}/detail/${id}`)
}

/** 创建会员 */
export async function createMember(data) {
  return request(`${API_BASE}/create`, { method: 'POST', body: JSON.stringify(data) })
}

/** 更新会员 */
export async function updateMember(id, data) {
  return request(`${API_BASE}/update/${id}`, { method: 'PUT', body: JSON.stringify(data) })
}

/** 删除会员 */
export async function deleteMember(id) {
  return request(`${API_BASE}/delete/${id}`, { method: 'DELETE' })
}

/** 会员状态变更 */
export async function updateMemberStatus(id, data) {
  return request(`${API_BASE}/status/${id}`, { method: 'PUT', body: JSON.stringify(data) })
}

/** 会员等级调整 */
export async function updateMemberLevel(id, data) {
  return request(`${API_BASE}/level/${id}`, { method: 'PUT', body: JSON.stringify(data) })
}

// ─── 会员等级 ────────────────────────────────────────────────

/** 等级列表 */
export async function getLevelList() {
  return request(`${API_BASE}/level/list`)
}

/** 创建等级 */
export async function createLevel(data) {
  return request(`${API_BASE}/level/create`, { method: 'POST', body: JSON.stringify(data) })
}

/** 更新等级 */
export async function updateLevel(id, data) {
  return request(`${API_BASE}/level/update/${id}`, { method: 'PUT', body: JSON.stringify(data) })
}

/** 删除等级 */
export async function deleteLevel(id) {
  return request(`${API_BASE}/level/delete/${id}`, { method: 'DELETE' })
}

/** 升级规则 */
export async function getUpgradeRules() {
  return request(`${API_BASE}/level/upgrade-rules`)
}

export async function updateUpgradeRules(data) {
  return request(`${API_BASE}/level/upgrade-rules`, { method: 'PUT', body: JSON.stringify(data) })
}

// ─── 优惠券 ─────────────────────────────────────────────────

/** 优惠券列表 */
export async function getCouponList(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`${API_BASE}/coupon/list?${qs}`)
}

/** 创建优惠券 */
export async function createCoupon(data) {
  return request(`${API_BASE}/coupon/create`, { method: 'POST', body: JSON.stringify(data) })
}

/** 更新优惠券 */
export async function updateCoupon(id, data) {
  return request(`${API_BASE}/coupon/update/${id}`, { method: 'PUT', body: JSON.stringify(data) })
}

/** 删除优惠券 */
export async function deleteCoupon(id) {
  return request(`${API_BASE}/coupon/delete/${id}`, { method: 'DELETE' })
}

/** 发放优惠券 */
export async function grantCoupon(data) {
  return request(`${API_BASE}/coupon/grant`, { method: 'POST', body: JSON.stringify(data) })
}

/** 会员优惠券列表 */
export async function getMemberCouponList(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`${API_BASE}/coupon/member-list?${qs}`)
}

// ─── 积分 ────────────────────────────────────────────────────

/** 积分规则查询 */
export async function getPointsRules() {
  return request(`${API_BASE}/points/rules`)
}

/** 积分规则设置 */
export async function updatePointsRules(data) {
  return request(`${API_BASE}/points/rules`, { method: 'PUT', body: JSON.stringify(data) })
}

/** 积分调整 */
export async function adjustPoints(data) {
  return request(`${API_BASE}/points/adjust`, { method: 'POST', body: JSON.stringify(data) })
}

/** 积分流水 */
export async function getPointsFlow(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`${API_BASE}/points/flow?${qs}`)
}

/** 会员积分余额 */
export async function getPointsBalance(memberId) {
  return request(`${API_BASE}/points/balance/${memberId}`)
}
