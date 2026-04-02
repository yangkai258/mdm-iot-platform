import axios from 'axios'

const BASE_URL = '/api/tenant-approvals'

function getToken() {
  return localStorage.getItem('token') || ''
}

function headers() {
  return { Authorization: `Bearer ${getToken()}` }
}

// ========== 租户入驻申请 ==========

/**
 * 提交租户入驻申请
 * @param {Object} data - 申请数据
 * @param {string} data.company_name - 公司名称
 * @param {string} data.contact_name - 联系人姓名
 * @param {string} data.contact_phone - 联系电话
 * @param {string} data.contact_email - 联系邮箱
 * @param {string} data.industry - 所属行业
 * @param {string} data.company_scale - 公司规模
 * @param {string} data.business_license - 营业执照号
 * @param {string} data.address - 公司地址
 * @param {string} data.description - 申请说明
 */
export async function submitApplication(data) {
  const res = await axios.post(BASE_URL, data, { headers: headers() })
  return res.data
}

/**
 * 获取申请列表
 * @param {Object} params
 * @param {number} params.page
 * @param {number} params.page_size
 * @param {string} params.status - pending/approved/rejected
 * @param {string} params.keyword - 搜索关键词
 */
export async function getApprovalList(params = {}) {
  const res = await axios.get(BASE_URL, { params, headers: headers() })
  return res.data
}

/**
 * 获取申请详情
 * @param {number} id - 申请ID
 */
export async function getApprovalDetail(id) {
  const res = await axios.get(`/api/${id}`, { headers: headers() })
  return res.data
}

/**
 * 审核通过
 * @param {number} id - 申请ID
 * @param {string} comment - 审核备注
 */
export async function approveApplication(id, comment = '') {
  const res = await axios.post(`/api/${id}/approve`, { comment }, { headers: headers() })
  return res.data
}

/**
 * 审核拒绝
 * @param {number} id - 申请ID
 * @param {string} comment - 拒绝原因
 */
export async function rejectApplication(id, comment = '') {
  const res = await axios.post(`/api/${id}/reject`, { comment }, { headers: headers() })
  return res.data
}

/**
 * 获取审批历史
 * @param {number} id - 申请ID
 */
export async function getApprovalHistory(id) {
  const res = await axios.get(`/api/${id}/history`, { headers: headers() })
  return res.data
}
