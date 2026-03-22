/**
 * 会员营销 API
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
  if (data.code !== 0 && data.code !== 200) {
    throw new Error(data.message || '请求失败')
  }
  return data
}

// ─── 会员标签 ────────────────────────────────────────────────

/** 标签列表 */
export async function getTagList(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`${API_BASE}/tag/list?${qs}`)
}

/** 创建标签 */
export async function createTag(data) {
  return request(`${API_BASE}/tag/create`, { method: 'POST', body: JSON.stringify(data) })
}

/** 更新标签 */
export async function updateTag(id, data) {
  return request(`${API_BASE}/tag/update/${id}`, { method: 'PUT', body: JSON.stringify(data) })
}

/** 删除标签 */
export async function deleteTag(id) {
  return request(`${API_BASE}/tag/delete/${id}`, { method: 'DELETE' })
}

/** 高频购买标签列表 */
export async function getHighFreqTagList(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`${API_BASE}/tag/high-freq/list?${qs}`)
}

/** 保存高频购买标签设置 */
export async function saveHighFreqTag(data) {
  return request(`${API_BASE}/tag/high-freq/save`, { method: 'POST', body: JSON.stringify(data) })
}

/** 低频购买标签列表 */
export async function getLowFreqTagList(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`${API_BASE}/tag/low-freq/list?${qs}`)
}

/** 保存低频购买标签设置 */
export async function saveLowFreqTag(data) {
  return request(`${API_BASE}/tag/low-freq/save`, { method: 'POST', body: JSON.stringify(data) })
}

/** 兴趣分类标签列表 */
export async function getInterestTagList(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`${API_BASE}/tag/interest/list?${qs}`)
}

/** 保存兴趣分类标签 */
export async function saveInterestTag(data) {
  return request(`${API_BASE}/tag/interest/save`, { method: 'POST', body: JSON.stringify(data) })
}

/** 标签自动清除设置 */
export async function getTagAutoCleanConfig() {
  return request(`${API_BASE}/tag/auto-clean/config`)
}

/** 保存标签自动清除设置 */
export async function saveTagAutoCleanConfig(data) {
  return request(`${API_BASE}/tag/auto-clean/config`, { method: 'PUT', body: JSON.stringify(data) })
}

/** 标签报表 */
export async function getTagReport(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`${API_BASE}/tag/report?${qs}`)
}

// ─── 促销活动 ────────────────────────────────────────────────

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

/** 红包记录 */
export async function getRedpacketRecords(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`${API_BASE}/redpacket/records?${qs}`)
}

/** 买赠促销列表 */
export async function getBuyGiftPromoList(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`${API_BASE}/promo/buy-gift/list?${qs}`)
}

/** 创建买赠促销 */
export async function createBuyGiftPromo(data) {
  return request(`${API_BASE}/promo/buy-gift/create`, { method: 'POST', body: JSON.stringify(data) })
}

/** 更新买赠促销 */
export async function updateBuyGiftPromo(id, data) {
  return request(`${API_BASE}/promo/buy-gift/update/${id}`, { method: 'PUT', body: JSON.stringify(data) })
}

/** 删除买赠促销 */
export async function deleteBuyGiftPromo(id) {
  return request(`${API_BASE}/promo/buy-gift/delete/${id}`, { method: 'DELETE' })
}

/** 直减促销列表 */
export async function getDirectReducePromoList(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`${API_BASE}/promo/direct-reduce/list?${qs}`)
}

/** 创建直减促销 */
export async function createDirectReducePromo(data) {
  return request(`${API_BASE}/promo/direct-reduce/create`, { method: 'POST', body: JSON.stringify(data) })
}

/** 更新直减促销 */
export async function updateDirectReducePromo(id, data) {
  return request(`${API_BASE}/promo/direct-reduce/update/${id}`, { method: 'PUT', body: JSON.stringify(data) })
}

/** 删除直减促销 */
export async function deleteDirectReducePromo(id) {
  return request(`${API_BASE}/promo/direct-reduce/delete/${id}`, { method: 'DELETE' })
}

/** 满额减促销列表 */
export async function getAmountReducePromoList(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`${API_BASE}/promo/amount-reduce/list?${qs}`)
}

/** 创建满额减促销 */
export async function createAmountReducePromo(data) {
  return request(`${API_BASE}/promo/amount-reduce/create`, { method: 'POST', body: JSON.stringify(data) })
}

/** 更新满额减促销 */
export async function updateAmountReducePromo(id, data) {
  return request(`${API_BASE}/promo/amount-reduce/update/${id}`, { method: 'PUT', body: JSON.stringify(data) })
}

/** 删除满额减促销 */
export async function deleteAmountReducePromo(id) {
  return request(`${API_BASE}/promo/amount-reduce/delete/${id}`, { method: 'DELETE' })
}

/** 满额折促销列表 */
export async function getAmountDiscountPromoList(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`${API_BASE}/promo/amount-discount/list?${qs}`)
}

/** 创建满额折促销 */
export async function createAmountDiscountPromo(data) {
  return request(`${API_BASE}/promo/amount-discount/create`, { method: 'POST', body: JSON.stringify(data) })
}

/** 更新满额折促销 */
export async function updateAmountDiscountPromo(id, data) {
  return request(`${API_BASE}/promo/amount-discount/update/${id}`, { method: 'PUT', body: JSON.stringify(data) })
}

/** 删除满额折促销 */
export async function deleteAmountDiscountPromo(id) {
  return request(`${API_BASE}/promo/amount-discount/delete/${id}`, { method: 'DELETE' })
}

/** 最高等级促销列表 */
export async function getVipExclusivePromoList(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`${API_BASE}/promo/vip-exclusive/list?${qs}`)
}

/** 创建最高等级促销 */
export async function createVipExclusivePromo(data) {
  return request(`${API_BASE}/promo/vip-exclusive/create`, { method: 'POST', body: JSON.stringify(data) })
}

/** 更新最高等级促销 */
export async function updateVipExclusivePromo(id, data) {
  return request(`${API_BASE}/promo/vip-exclusive/update/${id}`, { method: 'PUT', body: JSON.stringify(data) })
}

/** 删除最高等级促销 */
export async function deleteVipExclusivePromo(id) {
  return request(`${API_BASE}/promo/vip-exclusive/delete/${id}`, { method: 'DELETE' })
}

/** 促销活动类型列表 */
export async function getPromotionTypes(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`${API_BASE}/promo/types?${qs}`)
}

// ─── 会员礼包 ────────────────────────────────────────────────

/** 会员礼包列表 */
export async function getMemberGiftList(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`${API_BASE}/gift/list?${qs}`)
}

/** 创建会员礼包 */
export async function createMemberGift(data) {
  return request(`${API_BASE}/gift/create`, { method: 'POST', body: JSON.stringify(data) })
}

/** 更新会员礼包 */
export async function updateMemberGift(id, data) {
  return request(`${API_BASE}/gift/update/${id}`, { method: 'PUT', body: JSON.stringify(data) })
}

/** 删除会员礼包 */
export async function deleteMemberGift(id) {
  return request(`${API_BASE}/gift/delete/${id}`, { method: 'DELETE' })
}

/** 发放礼包 */
export async function grantMemberGift(data) {
  return request(`${API_BASE}/gift/grant`, { method: 'POST', body: JSON.stringify(data) })
}

/** 礼包发放明细 */
export async function getGiftRecords(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`${API_BASE}/gift/records?${qs}`)
}

/** 礼包统计 */
export async function getGiftStats() {
  return request(`${API_BASE}/gift/stats`)
}

// ─── 会员服务 ────────────────────────────────────────────────

/** 会员接待列表 */
export async function getMemberReceptionList(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`${API_BASE}/reception/list?${qs}`)
}

/** 创建接待记录 */
export async function createMemberReception(data) {
  return request(`${API_BASE}/reception/create`, { method: 'POST', body: JSON.stringify(data) })
}

/** 更新接待记录 */
export async function updateMemberReception(id, data) {
  return request(`${API_BASE}/reception/update/${id}`, { method: 'PUT', body: JSON.stringify(data) })
}

/** 删除接待记录 */
export async function deleteMemberReception(id) {
  return request(`${API_BASE}/reception/delete/${id}`, { method: 'DELETE' })
}

/** 会员推文流水 */
export async function getMemberArticleList(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`${API_BASE}/article/list?${qs}`)
}

/** 创建推文 */
export async function createMemberArticle(data) {
  return request(`${API_BASE}/article/create`, { method: 'POST', body: JSON.stringify(data) })
}

/** 删除推文 */
export async function deleteMemberArticle(id) {
  return request(`${API_BASE}/article/delete/${id}`, { method: 'DELETE' })
}

/** 短信模板列表 */
export async function getSmsTemplateList(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`${API_BASE}/sms/template/list?${qs}`)
}

/** 创建短信模板 */
export async function createSmsTemplate(data) {
  return request(`${API_BASE}/sms/template/create`, { method: 'POST', body: JSON.stringify(data) })
}

/** 更新短信模板 */
export async function updateSmsTemplate(id, data) {
  return request(`${API_BASE}/sms/template/update/${id}`, { method: 'PUT', body: JSON.stringify(data) })
}

/** 删除短信模板 */
export async function deleteSmsTemplate(id) {
  return request(`${API_BASE}/sms/template/delete/${id}`, { method: 'DELETE' })
}

/** 微信公众号设置 */
export async function getWechatSettings() {
  return request(`${API_BASE}/wechat/settings`)
}

/** 保存微信公众号设置 */
export async function saveWechatSettings(data) {
  return request(`${API_BASE}/wechat/settings`, { method: 'PUT', body: JSON.stringify(data) })
}

/** 会员权益列表 */
export async function getMemberBenefitList(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`${API_BASE}/benefit/list?${qs}`)
}

/** 创建会员权益 */
export async function createMemberBenefit(data) {
  return request(`${API_BASE}/benefit/create`, { method: 'POST', body: JSON.stringify(data) })
}

/** 更新会员权益 */
export async function updateMemberBenefit(id, data) {
  return request(`${API_BASE}/benefit/update/${id}`, { method: 'PUT', body: JSON.stringify(data) })
}

/** 删除会员权益 */
export async function deleteMemberBenefit(id) {
  return request(`${API_BASE}/benefit/delete/${id}`, { method: 'DELETE' })
}

/** 短信通道列表 */
export async function getSmsChannelList(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`${API_BASE}/sms/channel/list?${qs}`)
}

/** 创建短信通道 */
export async function createSmsChannel(data) {
  return request(`${API_BASE}/sms/channel/create`, { method: 'POST', body: JSON.stringify(data) })
}

/** 更新短信通道 */
export async function updateSmsChannel(id, data) {
  return request(`${API_BASE}/sms/channel/update/${id}`, { method: 'PUT', body: JSON.stringify(data) })
}

/** 删除短信通道 */
export async function deleteSmsChannel(id) {
  return request(`${API_BASE}/sms/channel/delete/${id}`, { method: 'DELETE' })
}

/** 测试短信通道 */
export async function testSmsChannel(id) {
  return request(`${API_BASE}/sms/channel/test/${id}`, { method: 'POST' })
}
