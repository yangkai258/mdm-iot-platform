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

// ─── 会员卡类型 ──────────────────────────────────────────────

/** 会员卡类型列表 */
export async function getCardTypeList(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`${API_BASE}/card-types?${qs}`)
}

/** 创建会员卡类型 */
export async function createCardType(data) {
  return request(`${API_BASE}/card-types`, { method: 'POST', body: JSON.stringify(data) })
}

/** 更新会员卡类型 */
export async function updateCardType(id, data) {
  return request(`${API_BASE}/card-types/${id}`, { method: 'PUT', body: JSON.stringify(data) })
}

/** 删除会员卡类型 */
export async function deleteCardType(id) {
  return request(`${API_BASE}/card-types/${id}`, { method: 'DELETE' })
}

// ─── 会员卡分组 ──────────────────────────────────────────────

/** 会员卡分组列表 */
export async function getCardGroupList(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`${API_BASE}/card-groups?${qs}`)
}

/** 创建会员卡分组 */
export async function createCardGroup(data) {
  return request(`${API_BASE}/card-groups`, { method: 'POST', body: JSON.stringify(data) })
}

/** 更新会员卡分组 */
export async function updateCardGroup(id, data) {
  return request(`${API_BASE}/card-groups/${id}`, { method: 'PUT', body: JSON.stringify(data) })
}

/** 删除会员卡分组 */
export async function deleteCardGroup(id) {
  return request(`${API_BASE}/card-groups/${id}`, { method: 'DELETE' })
}

// ─── 会员参数设置 ────────────────────────────────────────────

/** 获取会员参数设置 */
export async function getMemberSettings() {
  return request(`${API_BASE}/settings`)
}

/** 更新会员参数设置 */
export async function updateMemberSettings(data) {
  return request(`${API_BASE}/settings`, { method: 'PUT', body: JSON.stringify(data) })
}

// ─── 会员订单 ────────────────────────────────────────────────

/** 会员订单列表 */
export async function getMemberOrderList(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`${API_BASE}/orders?${qs}`)
}

/** 订单详情 */
export async function getMemberOrderDetail(id) {
  return request(`${API_BASE}/orders/${id}`)
}

// ─── 从业类型 ────────────────────────────────────────────────

/** 从业类型列表 */
export async function getOccupationTypeList(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`${API_BASE}/occupation-types?${qs}`)
}

/** 创建从业类型 */
export async function createOccupationType(data) {
  return request(`${API_BASE}/occupation-types`, { method: 'POST', body: JSON.stringify(data) })
}

/** 更新从业类型 */
export async function updateOccupationType(id, data) {
  return request(`${API_BASE}/occupation-types/${id}`, { method: 'PUT', body: JSON.stringify(data) })
}

/** 删除从业类型 */
export async function deleteOccupationType(id) {
  return request(`${API_BASE}/occupation-types/${id}`, { method: 'DELETE' })
}

// ─── 优惠券库存 ──────────────────────────────────────────────

/** 优惠券库存列表 */
export async function getCouponInventoryList(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`${API_BASE}/coupon/inventory?${qs}`)
}

// ─── 优惠消息流水 ────────────────────────────────────────────

/** 优惠消息流水列表 */
export async function getCouponMessageList(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`${API_BASE}/coupon/messages?${qs}`)
}

// ─── 优惠券发放记录 ───────────────────────────────────────────

/** 优惠券发放记录列表 */
export async function getCouponGrantList(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`${API_BASE}/coupon/grants?${qs}`)
}

// ─── 优惠券库存充值 ───────────────────────────────────────────

/** 优惠券库存充值 */
export async function rechargeCouponInventory(data) {
  return request(`${API_BASE}/coupon/inventory/recharge`, { method: 'POST', body: JSON.stringify(data) })
}

// ─── 红包 ────────────────────────────────────────────────────

/** 红包列表 */
export async function getRedpacketList(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`${API_BASE}/redpacket/list?${qs}`)
}

/** 创建红包 */
export async function createRedpacket(data) {
  return request(`${API_BASE}/redpacket/create`, { method: 'POST', body: JSON.stringify(data) })
}

/** 更新红包 */
export async function updateRedpacket(id, data) {
  return request(`${API_BASE}/redpacket/update/${id}`, { method: 'PUT', body: JSON.stringify(data) })
}

/** 删除红包 */
export async function deleteRedpacket(id) {
  return request(`${API_BASE}/redpacket/delete/${id}`, { method: 'DELETE' })
}

/** 发放红包 */
export async function grantRedpacket(data) {
  return request(`${API_BASE}/redpacket/grant`, { method: 'POST', body: JSON.stringify(data) })
}

// ─── 促销活动 ─────────────────────────────────────────────────

/** 买赠促销列表 */
export async function getBuyGiftList(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`${API_BASE}/promotion/buy-gift?${qs}`)
}

/** 直减促销列表 */
export async function getDirectReduceList(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`${API_BASE}/promotion/direct-reduce?${qs}`)
}

/** 满额减促销列表 */
export async function getAmountReduceList(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`${API_BASE}/promotion/amount-reduce?${qs}`)
}

/** 满额折促销列表 */
export async function getAmountDiscountList(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`${API_BASE}/promotion/amount-discount?${qs}`)
}

/** 最高等级促销列表 */
export async function getVipExclusiveList(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`${API_BASE}/promotion/vip-exclusive?${qs}`)
}

/** 促销活动类型列表 */
export async function getPromotionTypeList(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`${API_BASE}/promotion/types?${qs}`)
}
