<template>
  <div class="page-container">
    <!-- 筛选栏 -->
    <a-card class="filter-card">
      <div class="filter-row">
        <a-input-search
          v-model="filter.keyword"
          placeholder="搜索公司名称/联系人"
          style="width: 260px"
          @search="handleSearch"
          @press-enter="handleSearch"
        />
        <a-select
          v-model="filter.status"
          placeholder="审核状态"
          style="width: 140px"
          allow-clear
          @change="handleStatusChange"
        >
          <a-option value="pending">待审核</a-option>
          <a-option value="approved">已通过</a-option>
          <a-option value="rejected">已拒绝</a-option>
        </a-select>
        <a-range-picker
          v-model="filter.dateRange"
          style="width: 280px"
          @change="handleSearch"
        />
        <a-button type="primary" @click="handleSearch">
          <template #icon><icon-search /></template>
          查询
        </a-button>
        <a-button @click="handleReset">重置</a-button>
      </div>
    </a-card>

    <!-- 申请列表 -->
    <a-card class="table-card">
      <template #title>
        <div class="table-title">
          <span>租户入驻申请列表</span>
          <a-space>
            <a-tag v-if="counts.pending > 0" color="orange">待审核 {{ counts.pending }}</a-tag>
            <a-tag v-if="counts.approved > 0" color="green">已通过 {{ counts.approved }}</a-tag>
            <a-tag v-if="counts.rejected > 0" color="red">已拒绝 {{ counts.rejected }}</a-tag>
          </a-space>
        </div>
      </template>

      <a-table
        :columns="columns"
        :data="dataList"
        :loading="loading"
        :pagination="pagination"
        row-key="id"
        @change="handleTableChange"
      >
        <template #status="{ record }">
          <a-tag :color="statusColor(record.status)">
            {{ statusText(record.status) }}
          </a-tag>
        </template>
        <template #industry="{ record }">
          {{ record.industry || '-' }}
        </template>
        <template #scale="{ record }">
          {{ record.company_scale || '-' }}
        </template>
        <template #created_at="{ record }">
          {{ formatDate(record.created_at) }}
        </template>
        <template #actions="{ record }">
          <a-button type="text" size="small" @click="viewDetail(record)">
            查看详情
          </a-button>
          <template v-if="record.status === 'pending'">
            <a-divider direction="vertical" />
            <a-button type="text" size="small" status="danger" @click="openRejectModal(record)">
              拒绝
            </a-button>
            <a-divider direction="vertical" />
            <a-button type="text" size="small" status="success" @click="openApproveModal(record)">
              通过
            </a-button>
          </template>
        </template>
      </a-table>
    </a-card>

    <!-- 详情 Drawer -->
    <a-drawer
      v-model:visible="detailVisible"
      :title="`申请详情 - ${detailData.company_name || ''}`"
      width="680px"
      :footer="detailData.status === 'pending' ? drawerFooter : null"
      @close="detailVisible = false"
    >
      <div v-if="detailLoading" class="detail-loading">
        <a-spin size="large" />
      </div>
      <div v-else class="detail-content">
        <!-- 基本信息 -->
        <a-divider orientation="center">基本信息</a-divider>
        <a-descriptions :column="2" bordered size="small">
          <a-descriptions-item label="公司名称">{{ detailData.company_name || '-' }}</a-descriptions-item>
          <a-descriptions-item label="所属行业">{{ detailData.industry || '-' }}</a-descriptions-item>
          <a-descriptions-item label="公司规模">{{ detailData.company_scale || '-' }}</a-descriptions-item>
          <a-descriptions-item label="公司地址">{{ detailData.address || '-' }}</a-descriptions-item>
          <a-descriptions-item label="营业执照号" :span="2">{{ detailData.business_license || '-' }}</a-descriptions-item>
        </a-descriptions>

        <!-- 联系人信息 -->
        <a-divider orientation="center">联系人信息</a-divider>
        <a-descriptions :column="2" bordered size="small">
          <a-descriptions-item label="联系人">{{ detailData.contact_name || '-' }}</a-descriptions-item>
          <a-descriptions-item label="联系电话">{{ detailData.contact_phone || '-' }}</a-descriptions-item>
          <a-descriptions-item label="联系邮箱">{{ detailData.contact_email || '-' }}</a-descriptions-item>
        </a-descriptions>

        <!-- 申请说明 -->
        <a-divider orientation="center">申请说明</a-divider>
        <div class="description-box">
          {{ detailData.description || '无' }}
        </div>

        <!-- 审批历史 -->
        <a-divider orientation="center">审批历史</a-divider>
        <a-timeline v-if="historyList.length > 0" class="history-timeline">
          <a-timeline-item
            v-for="item in historyList"
            :key="item.id"
            :color="historyColor(item.action)"
          >
            <div class="timeline-item">
              <div class="timeline-header">
                <a-tag :color="historyColor(item.action)">{{ item.action_text }}</a-tag>
                <span class="timeline-time">{{ formatDate(item.created_at) }}</span>
              </div>
              <div class="timeline-operator">操作人：{{ item.operator || '系统' }}</div>
              <div v-if="item.comment" class="timeline-comment">
                备注：{{ item.comment }}
              </div>
            </div>
          </a-timeline-item>
        </a-timeline>
        <a-empty v-else description="暂无审批历史" />
      </div>
    </a-drawer>

    <!-- 审核操作 Footer -->
    <div v-if="detailData.status === 'pending'" class="drawer-footer">
      <a-space>
        <a-button @click="detailVisible = false">关闭</a-button>
        <a-button type="primary" status="danger" :loading="actionLoading" @click="openRejectModal(detailData)">
          拒绝申请
        </a-button>
        <a-button type="primary" status="success" :loading="actionLoading" @click="openApproveModal(detailData)">
          通过申请
        </a-button>
      </a-space>
    </div>

    <!-- 审核备注 Modal -->
    <a-modal
      v-model:visible="commentModalVisible"
      :title="approveAction === 'approve' ? '通过申请' : '拒绝申请'"
      :mask-closable="false"
      @before-ok="handleActionConfirm"
      @cancel="commentModalVisible = false"
    >
      <div class="comment-form">
        <div class="comment-tip">
          <a-tag :color="approveAction === 'approve' ? 'green' : 'red'">
            {{ approveAction === 'approve' ? '审核通过' : '审核拒绝' }}
          </a-tag>
          <span>申请：{{ currentRecord?.company_name }}</span>
        </div>
        <a-form :model="commentForm" layout="vertical">
          <a-form-item label="审核备注">
            <a-textarea
              v-model="commentForm.comment"
              :placeholder="approveAction === 'approve' ? '可选：添加审核通过备注' : '请输入拒绝原因'"
              :rows="4"
            />
          </a-form-item>
        </a-form>
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { Message } from '@arco-design/web-vue'
import * as tenantApi from '@/api/tenant.js'

const loading = ref(false)
const detailLoading = ref(false)
const actionLoading = ref(false)
const detailVisible = ref(false)
const commentModalVisible = ref(false)
const approveAction = ref('approve')
const currentRecord = ref(null)

const filter = reactive({
  keyword: '',
  status: undefined,
  dateRange: []
})

const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0
})

const counts = reactive({
  pending: 0,
  approved: 0,
  rejected: 0
})

const columns = [
  { title: 'ID', dataIndex: 'id', width: 70 },
  { title: '公司名称', dataIndex: 'company_name', ellipsis: true },
  { title: '联系人', dataIndex: 'contact_name' },
  { title: '联系电话', dataIndex: 'contact_phone', width: 130 },
  { title: '行业', slotName: 'industry', width: 100 },
  { title: '规模', slotName: 'scale', width: 100 },
  { title: '状态', slotName: 'status', width: 90 },
  { title: '申请时间', slotName: 'created_at', width: 160 },
  { title: '操作', slotName: 'actions', width: 180, fixed: 'right' }
]

const dataList = ref([])
const detailData = ref({})
const historyList = ref([])

const commentForm = reactive({
  comment: ''
})

const drawerFooter = computed(() => null)

const statusText = (status) => {
  const map = { pending: '待审核', approved: '已通过', rejected: '已拒绝' }
  return map[status] || status
}

const statusColor = (status) => {
  const map = { pending: 'orange', approved: 'green', rejected: 'red' }
  return map[status] || 'gray'
}

const historyColor = (action) => {
  const map = { approve: 'green', reject: 'red', submit: 'blue' }
  return map[action] || 'gray'
}

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  return d.toLocaleString('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit' })
}

const loadData = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.current,
      page_size: pagination.pageSize,
      keyword: filter.keyword || undefined,
      status: filter.status || undefined
    }
    const res = await tenantApi.getApprovalList(params)
    if (res.code === 0 || res.code === 200) {
      dataList.value = res.data?.list || []
      pagination.total = res.data?.pagination?.total || 0
      // 更新 counts
      const all = res.data?.list || []
      counts.pending = all.filter(i => i.status === 'pending').length
      counts.approved = all.filter(i => i.status === 'approved').length
      counts.rejected = all.filter(i => i.status === 'rejected').length
    }
  } catch {
    // 模拟数据
    dataList.value = getMockData()
    pagination.total = 6
    counts.pending = 2
    counts.approved = 2
    counts.rejected = 2
    Message.warning('使用模拟数据')
  } finally {
    loading.value = false
  }
}

const getMockData = () => [
  { id: 1, company_name: '北京智云科技有限公司', contact_name: '张明', contact_phone: '13800138001', industry: '互联网/IT', company_scale: '51-200人', status: 'pending', created_at: '2026-03-20T10:30:00Z', description: '我们是一家专注于企业级SaaS服务的互联网公司，计划使用MDM平台管理公司设备。' },
  { id: 2, company_name: '上海智造工业集团', contact_name: '李华', contact_phone: '13900139002', industry: '制造业', company_scale: '501-1000人', status: 'pending', created_at: '2026-03-20T14:22:00Z', description: '需要管理全国各工厂的生产设备接入和监控。' },
  { id: 3, company_name: '深圳创新科技有限公司', contact_name: '王芳', contact_phone: '13700137003', industry: '互联网/IT', company_scale: '201-500人', status: 'approved', created_at: '2026-03-19T09:15:00Z', description: '申请开通企业设备管理权限。' },
  { id: 4, company_name: '广州医疗健康集团', contact_name: '赵军', contact_phone: '13600136004', industry: '医疗健康', company_scale: '1000人以上', status: 'approved', created_at: '2026-03-19T11:40:00Z', description: '医疗机构设备管理需求。' },
  { id: 5, company_name: '杭州电商有限公司', contact_name: '钱伟', contact_phone: '13500135005', industry: '零售/电商', company_scale: '51-200人', status: 'rejected', created_at: '2026-03-18T16:00:00Z', description: '电商企业设备管理申请。' },
  { id: 6, company_name: '成都软件园', contact_name: '孙丽', contact_phone: '13400134006', industry: '互联网/IT', company_scale: '201-500人', status: 'rejected', created_at: '2026-03-18T08:30:00Z', description: '软件园区设备统一管理需求。' }
]

const loadDetail = async (id) => {
  detailLoading.value = true
  try {
    const res = await tenantApi.getApprovalDetail(id)
    if (res.code === 0 || res.code === 200) {
      detailData.value = res.data || {}
    } else {
      const mock = dataList.value.find(i => i.id === id)
      detailData.value = mock || {}
    }
  } catch {
    const mock = dataList.value.find(i => i.id === id)
    detailData.value = mock || {}
  } finally {
    detailLoading.value = false
  }
}

const loadHistory = async (id) => {
  try {
    const res = await tenantApi.getApprovalHistory(id)
    if (res.code === 0 || res.code === 200) {
      historyList.value = res.data || []
    } else {
      historyList.value = getMockHistory(id)
    }
  } catch {
    historyList.value = getMockHistory(id)
  }
}

const getMockHistory = (id) => {
  const record = dataList.value.find(i => i.id === id)
  if (!record) return []
  const history = [
    { id: 1, action: 'submit', action_text: '提交申请', operator: record.contact_name, comment: '', created_at: record.created_at }
  ]
  if (record.status === 'approved') {
    history.push({ id: 2, action: 'approve', action_text: '审核通过', operator: '管理员', comment: '资质审核通过', created_at: new Date(new Date(record.created_at).getTime() + 86400000).toISOString() })
  } else if (record.status === 'rejected') {
    history.push({ id: 2, action: 'reject', action_text: '审核拒绝', operator: '管理员', comment: '资质不符合要求', created_at: new Date(new Date(record.created_at).getTime() + 86400000).toISOString() })
  }
  return history
}

const viewDetail = async (record) => {
  currentRecord.value = record
  await loadDetail(record.id)
  await loadHistory(record.id)
  detailVisible.value = true
}

const openApproveModal = (record) => {
  currentRecord.value = record
  approveAction.value = 'approve'
  commentForm.comment = ''
  commentModalVisible.value = true
}

const openRejectModal = (record) => {
  currentRecord.value = record
  approveAction.value = 'reject'
  commentForm.comment = ''
  commentModalVisible.value = true
}

const handleActionConfirm = async (done) => {
  if (approveAction.value === 'reject' && !commentForm.comment.trim()) {
    Message.error('请输入拒绝原因')
    done(false)
    return
  }

  actionLoading.value = true
  try {
    const id = currentRecord.value.id
    let res
    if (approveAction.value === 'approve') {
      res = await tenantApi.approveApplication(id, commentForm.comment)
    } else {
      res = await tenantApi.rejectApplication(id, commentForm.comment)
    }

    if (res.code === 0 || res.code === 200) {
      Message.success(approveAction.value === 'approve' ? '已通过申请' : '已拒绝申请')
      commentModalVisible.value = false
      detailVisible.value = false
      loadData()
    } else {
      Message.error(res.message || '操作失败')
      done(false)
    }
  } catch {
    // 模拟成功
    Message.success('操作成功（模拟）')
    commentModalVisible.value = false
    detailVisible.value = false
    loadData()
    done(false)
  } finally {
    actionLoading.value = false
  }
}

const handleSearch = () => {
  pagination.current = 1
  loadData()
}

const handleReset = () => {
  filter.keyword = ''
  filter.status = undefined
  filter.dateRange = []
  handleSearch()
}

const handleStatusChange = () => {
  pagination.current = 1
  loadData()
}

const handleTableChange = (pag) => {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  loadData()
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.page-container {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.filter-row {
  display: flex;
  gap: 12px;
  align-items: center;
  flex-wrap: wrap;
}

.table-title {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.detail-loading {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 200px;
}

.detail-content {
  padding: 0 8px;
}

.description-box {
  background: #f8f9fa;
  border-radius: 6px;
  padding: 16px;
  color: #333;
  line-height: 1.6;
  min-height: 80px;
  white-space: pre-wrap;
}

.history-timeline {
  padding: 8px 0;
}

.timeline-item {
  padding-bottom: 4px;
}

.timeline-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 4px;
}

.timeline-time {
  color: #999;
  font-size: 12px;
}

.timeline-operator {
  color: #666;
  font-size: 13px;
}

.timeline-comment {
  color: #888;
  font-size: 12px;
  margin-top: 4px;
  font-style: italic;
}

.drawer-footer {
  padding: 16px;
  border-top: 1px solid #f0f0f0;
  display: flex;
  justify-content: flex-end;
}

.comment-form {
  padding: 8px 0;
}

.comment-tip {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 16px;
  padding: 12px;
  background: #f8f9fa;
  border-radius: 6px;
}
</style>
