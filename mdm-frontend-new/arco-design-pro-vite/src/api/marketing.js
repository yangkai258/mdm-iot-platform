/**
 * 会员营销 API
 * Base: /api/v1/member
 */
const API_BASE = '/api/member'

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
  return request(`/api/tag/list?${qs}`)
}

/** 创建标签 */
export async function createTag(data) {
  return request(`/api/tag/create`, { method: 'POST', body: JSON.stringify(data) })
}

/** 更新标签 */
export async function updateTag(id, data) {
  return request(`/api/tag/update/${id}`, { method: 'PUT', body: JSON.stringify(data) })
}

/** 删除标签 */
export async function deleteTag(id) {
  return request(`/api/tag/delete/${id}`, { method: 'DELETE' })
}

/** 高频购买标签列表 */
export async function getHighFreqTagList(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`/api/tag/high-freq/list?${qs}`)
}

/** 保存高频购买标签设置 */
export async function saveHighFreqTag(data) {
  return request(`/api/tag/high-freq/save`, { method: 'POST', body: JSON.stringify(data) })
}

/** 低频购买标签列表 */
export async function getLowFreqTagList(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`/api/tag/low-freq/list?${qs}`)
}

/** 保存低频购买标签设置 */
export async function saveLowFreqTag(data) {
  return request(`/api/tag/low-freq/save`, { method: 'POST', body: JSON.stringify(data) })
}

/** 兴趣分类标签列表 */
export async function getInterestTagList(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`/api/tag/interest/list?${qs}`)
}

/** 保存兴趣分类标签 */
export async function saveInterestTag(data) {
  return request(`/api/tag/interest/save`, { method: 'POST', body: JSON.stringify(data) })
}

/** 标签自动清除设置 */
export async function getTagAutoCleanConfig() {
  return request(`/api/tag/auto-clean/config`)
}

/** 保存标签自动清除设置 */
export async function saveTagAutoCleanConfig(data) {
  return request(`/api/tag/auto-clean/config`, { method: 'PUT', body: JSON.stringify(data) })
}

/** 标签报表 */
export async function getTagReport(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`/api/tag/report?${qs}`)
}

// ─── 促销活动 ────────────────────────────────────────────────

/** 红包列表 */
export async function getRedpacketList(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`/api/redpacket/list?${qs}`)
}

/** 创建红包 */
export async function createRedpacket(data) {
  return request(`/api/redpacket/create`, { method: 'POST', body: JSON.stringify(data) })
}

/** 更新红包 */
export async function updateRedpacket(id, data) {
  return request(`/api/redpacket/update/${id}`, { method: 'PUT', body: JSON.stringify(data) })
}

/** 删除红包 */
export async function deleteRedpacket(id) {
  return request(`/api/redpacket/delete/${id}`, { method: 'DELETE' })
}

/** 发放红包 */
export async function grantRedpacket(data) {
  return request(`/api/redpacket/grant`, { method: 'POST', body: JSON.stringify(data) })
}

/** 红包记录 */
export async function getRedpacketRecords(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`/api/redpacket/records?${qs}`)
}

/** 买赠促销列表 */
export async function getBuyGiftPromoList(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`/api/promo/buy-gift/list?${qs}`)
}

/** 创建买赠促销 */
export async function createBuyGiftPromo(data) {
  return request(`/api/promo/buy-gift/create`, { method: 'POST', body: JSON.stringify(data) })
}

/** 更新买赠促销 */
export async function updateBuyGiftPromo(id, data) {
  return request(`/api/promo/buy-gift/update/${id}`, { method: 'PUT', body: JSON.stringify(data) })
}

/** 删除买赠促销 */
export async function deleteBuyGiftPromo(id) {
  return request(`/api/promo/buy-gift/delete/${id}`, { method: 'DELETE' })
}

/** 直减促销列表 */
export async function getDirectReducePromoList(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`/api/promo/direct-reduce/list?${qs}`)
}

/** 创建直减促销 */
export async function createDirectReducePromo(data) {
  return request(`/api/promo/direct-reduce/create`, { method: 'POST', body: JSON.stringify(data) })
}

/** 更新直减促销 */
export async function updateDirectReducePromo(id, data) {
  return request(`/api/promo/direct-reduce/update/${id}`, { method: 'PUT', body: JSON.stringify(data) })
}

/** 删除直减促销 */
export async function deleteDirectReducePromo(id) {
  return request(`/api/promo/direct-reduce/delete/${id}`, { method: 'DELETE' })
}

/** 满额减促销列表 */
export async function getAmountReducePromoList(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`/api/promo/amount-reduce/list?${qs}`)
}

/** 创建满额减促销 */
export async function createAmountReducePromo(data) {
  return request(`/api/promo/amount-reduce/create`, { method: 'POST', body: JSON.stringify(data) })
}

/** 更新满额减促销 */
export async function updateAmountReducePromo(id, data) {
  return request(`/api/promo/amount-reduce/update/${id}`, { method: 'PUT', body: JSON.stringify(data) })
}

/** 删除满额减促销 */
export async function deleteAmountReducePromo(id) {
  return request(`/api/promo/amount-reduce/delete/${id}`, { method: 'DELETE' })
}

/** 满额折促销列表 */
export async function getAmountDiscountPromoList(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`/api/promo/amount-discount/list?${qs}`)
}

/** 创建满额折促销 */
export async function createAmountDiscountPromo(data) {
  return request(`/api/promo/amount-discount/create`, { method: 'POST', body: JSON.stringify(data) })
}

/** 更新满额折促销 */
export async function updateAmountDiscountPromo(id, data) {
  return request(`/api/promo/amount-discount/update/${id}`, { method: 'PUT', body: JSON.stringify(data) })
}

/** 删除满额折促销 */
export async function deleteAmountDiscountPromo(id) {
  return request(`/api/promo/amount-discount/delete/${id}`, { method: 'DELETE' })
}

/** 最高等级促销列表 */
export async function getVipExclusivePromoList(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`/api/promo/vip-exclusive/list?${qs}`)
}

/** 创建最高等级促销 */
export async function createVipExclusivePromo(data) {
  return request(`/api/promo/vip-exclusive/create`, { method: 'POST', body: JSON.stringify(data) })
}

/** 更新最高等级促销 */
export async function updateVipExclusivePromo(id, data) {
  return request(`/api/promo/vip-exclusive/update/${id}`, { method: 'PUT', body: JSON.stringify(data) })
}

/** 删除最高等级促销 */
export async function deleteVipExclusivePromo(id) {
  return request(`/api/promo/vip-exclusive/delete/${id}`, { method: 'DELETE' })
}

/** 促销活动类型列表 */
export async function getPromotionTypes(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`/api/promo/types?${qs}`)
}

// ─── 会员礼包 ────────────────────────────────────────────────

/** 会员礼包列表 */
export async function getMemberGiftList(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`/api/gift/list?${qs}`)
}

/** 创建会员礼包 */
export async function createMemberGift(data) {
  return request(`/api/gift/create`, { method: 'POST', body: JSON.stringify(data) })
}

/** 更新会员礼包 */
export async function updateMemberGift(id, data) {
  return request(`/api/gift/update/${id}`, { method: 'PUT', body: JSON.stringify(data) })
}

/** 删除会员礼包 */
export async function deleteMemberGift(id) {
  return request(`/api/gift/delete/${id}`, { method: 'DELETE' })
}

/** 发放礼包 */
export async function grantMemberGift(data) {
  return request(`/api/gift/grant`, { method: 'POST', body: JSON.stringify(data) })
}

/** 礼包发放明细 */
export async function getGiftRecords(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`/api/gift/records?${qs}`)
}

/** 礼包统计 */
export async function getGiftStats() {
  return request(`/api/gift/stats`)
}

// ─── 会员服务 ────────────────────────────────────────────────

/** 会员接待列表 */
export async function getMemberReceptionList(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`/api/reception/list?${qs}`)
}

/** 创建接待记录 */
export async function createMemberReception(data) {
  return request(`/api/reception/create`, { method: 'POST', body: JSON.stringify(data) })
}

/** 更新接待记录 */
export async function updateMemberReception(id, data) {
  return request(`/api/reception/update/${id}`, { method: 'PUT', body: JSON.stringify(data) })
}

/** 删除接待记录 */
export async function deleteMemberReception(id) {
  return request(`/api/reception/delete/${id}`, { method: 'DELETE' })
}

/** 会员推文流水 */
export async function getMemberArticleList(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`/api/article/list?${qs}`)
}

/** 创建推文 */
export async function createMemberArticle(data) {
  return request(`/api/article/create`, { method: 'POST', body: JSON.stringify(data) })
}

/** 删除推文 */
export async function deleteMemberArticle(id) {
  return request(`/api/article/delete/${id}`, { method: 'DELETE' })
}

/** 短信模板列表 */
export async function getSmsTemplateList(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`/api/sms/template/list?${qs}`)
}

/** 创建短信模板 */
export async function createSmsTemplate(data) {
  return request(`/api/sms/template/create`, { method: 'POST', body: JSON.stringify(data) })
}

/** 更新短信模板 */
export async function updateSmsTemplate(id, data) {
  return request(`/api/sms/template/update/${id}`, { method: 'PUT', body: JSON.stringify(data) })
}

/** 删除短信模板 */
export async function deleteSmsTemplate(id) {
  return request(`/api/sms/template/delete/${id}`, { method: 'DELETE' })
}

/** 微信公众号设置 */
export async function getWechatSettings() {
  return request(`/api/wechat/settings`)
}

/** 保存微信公众号设置 */
export async function saveWechatSettings(data) {
  return request(`/api/wechat/settings`, { method: 'PUT', body: JSON.stringify(data) })
}

/** 会员权益列表 */
export async function getMemberBenefitList(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`/api/benefit/list?${qs}`)
}

/** 创建会员权益 */
export async function createMemberBenefit(data) {
  return request(`/api/benefit/create`, { method: 'POST', body: JSON.stringify(data) })
}

/** 更新会员权益 */
export async function updateMemberBenefit(id, data) {
  return request(`/api/benefit/update/${id}`, { method: 'PUT', body: JSON.stringify(data) })
}

/** 删除会员权益 */
export async function deleteMemberBenefit(id) {
  return request(`/api/benefit/delete/${id}`, { method: 'DELETE' })
}

/** 短信通道列表 */
export async function getSmsChannelList(params = {}) {
  const qs = new URLSearchParams(params).toString()
  return request(`/api/sms/channel/list?${qs}`)
}

/** 创建短信通道 */
export async function createSmsChannel(data) {
  return request(`/api/sms/channel/create`, { method: 'POST', body: JSON.stringify(data) })
}

/** 更新短信通道 */
export async function updateSmsChannel(id, data) {
  return request(`/api/sms/channel/update/${id}`, { method: 'PUT', body: JSON.stringify(data) })
}

/** 删除短信通道 */
export async function deleteSmsChannel(id) {
  return request(`/api/sms/channel/delete/${id}`, { method: 'DELETE' })
}

/** 测试短信通道 */
export async function testSmsChannel(id) {
  return request(`/api/sms/channel/test/${id}`, { method: 'POST' })
}
