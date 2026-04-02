<template>
    <Breadcrumb :items="['Home','Console','']" />


  <div class="page-container">
    <a-card class="general-card" title="租户审批">
      <template #extra>
        <a-button @click="loadApplications"><icon-refresh />刷新</a-button>
      </template>
      <div class="search-bar">
        <a-input-search
          v-model="searchKey"
          placeholder="搜索租户名称/联系人"
          style="width: 280px"
          @search="loadApplications"
        />
      </div>
      <!-- 待审核列表 -->
      <a-table
        :columns="columns"
        :data="filteredData"
        :loading="loading"
        :pagination="{ pageSize: 10 }"
        row-key="id"
      >
        <template #applicationCode="{ record }">
          <a-link @click="viewDetail(record)">{{ record.application_code }}</a-link>
        </template>
        <template #status="{ record }">
          <a-tag :color="getStatusColor(record.status)">
            {{ getStatusText(record.status) }}
          </a-tag>
        </template>
        <template #applyTime="{ record }">
          {{ formatDate(record.apply_time) }}
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="viewDetail(record)">「查看」</a-button>
            <a-button
              v-if="record.status === 'pending'"
              type="text"
              size="small"
              status="success"
              @click="handleApprove(record)"
            >「通过」</a-button>
            <a-button
              v-if="record.status === 'pending'"
              type="text"
              size="small"
              status="danger"
              @click="openRejectModal(record)"
            >「拒绝」</a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 审核详情抽屉 -->
    <a-drawer
      v-model:visible="detailVisible"
      title="申请详情"
      :width="560"
      @close="detailVisible = false"
    >
      <a-descriptions :column="1" bordered size="large">
        <a-descriptions-item label="申请编号">{{ currentRecord?.application_code }}</a-descriptions-item>
        <a-descriptions-item label="公司名称">{{ currentRecord?.company_name }}</a-descriptions-item>
        <a-descriptions-item label="联系人">{{ currentRecord?.contact_name }}</a-descriptions-item>
        <a-descriptions-item label="联系电话">{{ currentRecord?.contact_phone }}</a-descriptions-item>
        <a-descriptions-item label="联系邮箱">{{ currentRecord?.contact_email }}</a-descriptions-item>
        <a-descriptions-item label="所属行业">{{ currentRecord?.industry || '-' }}</a-descriptions-item>
        <a-descriptions-item label="企业规模">{{ currentRecord?.company_size || '-' }}</a-descriptions-item>
        <a-descriptions-item label="申请套餐">{{ currentRecord?.plan_name }}</a-descriptions-item>
        <a-descriptions-item label="使用场景">{{ currentRecord?.use_case || '-' }}</a-descriptions-item>
        <a-descriptions-item label="申请时间">{{ formatDate(currentRecord?.apply_time) }}</a-descriptions-item>
      </a-descriptions>

      <!-- 审核配置（仅pending状态显示） -->
      <a-divider v-if="currentRecord?.status === 'pending'">审核配置</a-divider>
      <a-form
        v-if="currentRecord?.status === 'pending'"
        :model="approvalForm"
        layout="vertical"
        style="max-width: 400px"
      >
        <a-form-item label="实际分配套餐">
          <a-select v-model="approvalForm.effective_plan_id" placeholder="选择套餐">
            <a-option v-for="p in plans" :key="p.id" :value="p.id">{{ p.plan_name }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="试用天数">
          <a-input-number v-model="approvalForm.trial_days" :min="0" :max="365" placeholder="不填则无试用" style="width: 100%" />
        </a-form-item>
        <a-form-item label="审核备注">
          <a-textarea v-model="approvalForm.admin_notes" placeholder="可选备注" :rows="3" />
        </a-form-item>
      </a-form>

      <a-divider>审核历史</a-divider>
      <a-timeline v-if="approvalHistory.length">
        <a-timeline-item
          v-for="item in approvalHistory"
          :key="item.id"
          :color="item.action === 'approve' ? 'green' : item.action === 'reject' ? 'red' : 'gray'"
        >
          <p><strong>{{ item.action_text }}</strong> — {{ item.operator }}</p>
          <p>{{ item.comment || '无' }}</p>
          <p class="text-gray">{{ formatDate(item.created_at) }}</p>
        </a-timeline-item>
      </a-timeline>
      <a-empty v-else description="暂无审核记录" />

      <template v-if="currentRecord?.status === 'pending'" #footer>
        <a-space>
          <a-button status="danger" @click="openRejectModal(currentRecord)">「拒绝」</a-button>
          <a-button type="primary" status="success" @click="handleApprove(currentRecord)">「通过」</a-button>
        </a-space>
      </template>
    </a-drawer>

    <!-- 拒绝原因弹窗 -->
    <a-modal
      v-model:visible="rejectModalVisible"
      title="拒绝申请"
      @before-ok="submitReject"
      @cancel="rejectModalVisible = false"
    >
      <a-form :model="rejectForm" layout="vertical">
        <a-form-item label="拒绝原因" required>
          <a-textarea v-model="rejectForm.reject_reason" placeholder="请输入拒绝原因" :rows="4" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import { IconRefresh } from '@arco-design/web-vue/es/icon'
import axios from 'axios'

const loading = ref(false)
const searchKey = ref('')
const detailVisible = ref(false)
const rejectModalVisible = ref(false)
const currentRecord = ref<any>(null)
const approvalHistory = ref<any[]>([])
const plans = ref<any[]>([])

const approvalForm = ref({
  effective_plan_id: null as number | null,
  trial_days: 30,
  admin_notes: ''
})

const rejectForm = ref({ reject_reason: '' })

const columns = [
  { title: '申请编号', slotName: 'applicationCode', width: 180 },
  { title: '公司名称', dataIndex: 'company_name' },
  { title: '联系人', dataIndex: 'contact_name', width: 120 },
  { title: '联系电话', dataIndex: 'contact_phone', width: 140 },
  { title: '申请套餐', dataIndex: 'plan_name', width: 100 },
  { title: '申请时间', slotName: 'applyTime', width: 170 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '操作', slotName: 'actions', width: 180, fixed: 'right' },
]

// Mock data
const mockData = ref([
  {
    id: 1,
    application_code: 'APP-20260320-001',
    company_name: '深圳市智联科技有限公司',
    contact_name: '张经理',
    contact_phone: '13800138000',
    contact_email: 'zhang@zhilian.com',
    industry: '互联网',
    company_size: '100-500人',
    plan_name: '企业版',
    plan_id: 4,
    use_case: '用于公司物联网设备管理',
    status: 'pending',
    apply_time: '2026-03-20T10:30:00Z',
  },
  {
    id: 2,
    application_code: 'APP-20260319-001',
    company_name: '广州云智数据服务公司',
    contact_name: '李总监',
    contact_phone: '13900139000',
    contact_email: 'li@yunzhi.com',
    industry: '金融科技',
    company_size: '50-100人',
    plan_name: '专业版',
    plan_id: 3,
    use_case: '数据采集设备管理',
    status: 'pending',
    apply_time: '2026-03-19T15:20:00Z',
  },
])

const mockPlans = [
  { id: 1, plan_name: '免费版' },
  { id: 2, plan_name: '基础版' },
  { id: 3, plan_name: '专业版' },
  { id: 4, plan_name: '企业版' },
]

const filteredData = computed(() => {
  if (!searchKey.value) return mockData.value
  return mockData.value.filter(item =>
    item.company_name.includes(searchKey.value) ||
    item.contact_name.includes(searchKey.value)
  )
})

const getStatusColor = (status: string) => {
  const map: Record<string, string> = { pending: 'orange', approved: 'green', rejected: 'red' }
  return map[status] || 'gray'
}

const getStatusText = (status: string) => {
  const map: Record<string, string> = { pending: '待审核', approved: '已通过', rejected: '已拒绝' }
  return map[status] || status
}

const formatDate = (dateStr: string) => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleString('zh-CN')
}

const loadApplications = async () => {
  loading.value = true
  try {
    // API: GET /api/admin/tenants/applications
    // const res = await axios.get('/api/admin/tenants/applications', { params: { keyword: searchKey.value } })
    // if (res.data.code === 'SUCCESS') { ... }
    // 使用mock
    plans.value = mockPlans
  } catch {
    Message.error('加载申请列表失败')
  } finally {
    loading.value = false
  }
}

const viewDetail = async (record: any) => {
  currentRecord.value = record
  approvalForm.value.effective_plan_id = record.plan_id
  approvalForm.value.trial_days = 30
  approvalForm.value.admin_notes = ''

  // Mock history
  approvalHistory.value = [
    {
      id: 1,
      action: 'apply',
      action_text: '提交申请',
      operator: record.contact_name,
      comment: '',
      created_at: record.apply_time,
    },
  ]

  // API: GET /api/admin/tenants/applications/{id}/history
  detailVisible.value = true
}

const handleApprove = async (record: any) => {
  try {
    // API: POST /api/admin/tenants/applications/{id}/approve
    // await axios.post(`/api/admin/tenants/applications/${record.id}/approve`, approvalForm.value)
    record.status = 'approved'
    Message.success('已通过审核')
    detailVisible.value = false
    loadApplications()
  } catch {
    Message.error('审核操作失败')
  }
}

const openRejectModal = (record: any) => {
  currentRecord.value = record
  rejectForm.value.reject_reason = ''
  rejectModalVisible.value = true
}

const submitReject = async (done: (val: boolean) => void) => {
  if (!rejectForm.value.reject_reason.trim()) {
    Message.warning('请输入拒绝原因')
    done(false)
    return
  }
  try {
    // API: POST /api/admin/tenants/applications/{id}/reject
    // await axios.post(`/api/admin/tenants/applications/${currentRecord.value.id}/reject`, {
    //   reject_reason: rejectForm.value.reject_reason
    // })
    currentRecord.value.status = 'rejected'
    Message.success('已拒绝')
    rejectModalVisible.value = false
    detailVisible.value = false
    loadApplications()
    done(true)
  } catch {
    Message.error('操作失败')
    done(false)
  }
}

onMounted(() => {
  loadApplications()
})
</script>

<style scoped>
.page-container { padding: 16px; }
.search-bar { margin-bottom: 16px; }
.text-gray { color: #999; font-size: 12px; }
</style>
