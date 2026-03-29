<template>
  <div class="page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item><a href="#/dashboard">首页</a></a-breadcrumb-item>
      <a-breadcrumb-item>租户管理</a-breadcrumb-item>
      <a-breadcrumb-item>租户管理</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 搜索框（左上角） -->
    <div class="search-bar">
      <a-input-search
        v-model="searchKey"
        placeholder="搜索..."
        style="width: 280px"
        @search="loadTenants"
      />
    </div>

    <!-- 操作按钮组（靠左） -->
    <div class="action-bar">
      <a-space :size="12">
        <a-button @click="showFilter = !showFilter">「筛选」</a-button>
        <a-button @click="loadTenants">「刷新」</a-button>
      </a-space>
    </div>

    <!-- 筛选面板 -->
    <a-card v-if="showFilter" :bordered="false" class="general-card" style="margin-bottom: 12px">
      <a-space wrap>
        <a-select v-model="filterStatus" placeholder="状态" style="width: 120px" allow-clear>
          <a-option value="active">正常</a-option>
          <a-option value="suspended">已禁用</a-option>
          <a-option value="expired">已到期</a-option>
        </a-select>
        <a-select v-model="filterPlan" placeholder="套餐" style="width: 140px" allow-clear>
          <a-option v-for="p in plans" :key="p.id" :value="p.plan_code">{{ p.plan_name }}</a-option>
        </a-select>
        <a-select v-model="filterExpire" placeholder="到期" style="width: 140px" allow-clear>
          <a-option value="expired">已到期</a-option>
          <a-option value="expiring_soon">7天内到期</a-option>
        </a-select>
        <a-button type="primary" @click="loadTenants">查询</a-button>
        <a-button @click="resetFilter">重置</a-button>
      </a-space>
    </a-card>

    <!-- 租户列表 -->
    <a-card :bordered="false" class="general-card">
      <a-table
        :columns="columns"
        :data="filteredData"
        :loading="loading"
        :pagination="{ pageSize: 10 }"
        row-key="id"
        scroll="x"
      >
        <template #tenantCode="{ record }">
          <a-link @click="openEditModal(record)">{{ record.tenant_code }}</a-link>
        </template>
        <template #planInfo="{ record }">
          <a-tag :color="getPlanColor(record.plan_code)">{{ record.plan_name }}</a-tag>
        </template>
        <template #userQuota="{ record }">
          <span :class="record.user_count > record.user_quota ? 'text-danger' : ''">
            {{ record.user_count }}/{{ record.user_quota === -1 ? '不限' : record.user_quota }}
          </span>
        </template>
        <template #deviceQuota="{ record }">
          <span :class="record.device_count > record.device_quota ? 'text-danger' : ''">
            {{ record.device_count }}/{{ record.device_quota === -1 ? '不限' : record.device_quota }}
          </span>
        </template>
        <template #statusBadge="{ record }">
          <a-tag :color="getStatusColor(record.status)">
            {{ getStatusText(record.status) }}
          </a-tag>
        </template>
        <template #expireTime="{ record }">
          <span :class="isExpiringSoon(record.expire_time) ? 'text-warning' : isExpired(record.expire_time) ? 'text-danger' : ''">
            {{ formatDate(record.expire_time) }}
          </span>
        </template>
        <template #createdAt="{ record }">
          {{ formatDate(record.created_at) }}
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="openEditModal(record)">「编辑」</a-button>
            <a-button
              v-if="record.status === 'active'"
              type="text"
              size="small"
              status="warning"
              @click="openDisableConfirm(record)"
            >「禁用」</a-button>
            <a-button
              v-if="record.status === 'suspended'"
              type="text"
              size="small"
              status="success"
              @click="handleActivate(record)"
            >「启用」</a-button>
            <a-button type="text" size="small" @click="openPlanModal(record)">「套餐变更」</a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 编辑租户 - 全屏模态 (风格D) -->
    <a-modal
      v-model:visible="editModalVisible"
      :title="`编辑租户 — ${currentRecord?.company_name || ''}`"
      :fullscreen="true"
      :mask-closable="false"
      @before-ok="submitEdit"
      @cancel="editModalVisible = false"
    >
      <a-form
        v-if="currentRecord"
        ref="editFormRef"
        :model="editForm"
        :rules="editFormRules"
        layout="vertical"
        style="max-width: 800px"
      >
        <a-divider>基本信息</a-divider>
        <a-row :gutter="24">
          <a-col :span="12">
            <a-form-item label="公司名称" field="company_name">
              <a-input v-model="editForm.company_name" placeholder="请输入公司名称" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="行业" field="industry">
              <a-select v-model="editForm.industry" placeholder="请选择行业">
                <a-option value="互联网">互联网</a-option>
                <a-option value="金融科技">金融科技</a-option>
                <a-option value="制造业">制造业</a-option>
                <a-option value="零售">零售</a-option>
                <a-option value="医疗">医疗</a-option>
                <a-option value="教育">教育</a-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="联系人" field="contact_name">
              <a-input v-model="editForm.contact_name" placeholder="请输入联系人姓名" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="联系电话" field="contact_phone">
              <a-input v-model="editForm.contact_phone" placeholder="请输入联系电话" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="联系邮箱" field="contact_email">
              <a-input v-model="editForm.contact_email" placeholder="请输入联系邮箱" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="企业规模" field="company_size">
              <a-select v-model="editForm.company_size" placeholder="请选择规模">
                <a-option value="1-50">1-50人</a-option>
                <a-option value="51-200">51-200人</a-option>
                <a-option value="201-500">201-500人</a-option>
                <a-option value="500+">500人以上</a-option>
              </a-select>
            </a-form-item>
          </a-col>
        </a-row>

        <a-divider>系统配置</a-divider>
        <a-row :gutter="24">
          <a-col :span="12">
            <a-form-item label="Logo URL" field="tenant_logo">
              <a-input v-model="editForm.tenant_logo" placeholder="https://..." />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="自定义域名" field="custom_domain">
              <a-input v-model="editForm.custom_domain" placeholder="yourcompany.mdm.com" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="到期时间" field="expire_time">
              <a-date-picker
                v-model="editForm.expire_time"
                show-time
                format="YYYY-MM-DD HH:mm"
                style="width: 100%"
              />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="当前套餐" field="plan_id">
              <a-select v-model="editForm.plan_id" placeholder="选择套餐">
                <a-option v-for="p in plans" :key="p.id" :value="p.id">{{ p.plan_name }}</a-option>
              </a-select>
            </a-form-item>
          </a-col>
        </a-row>
      </a-form>
    </a-modal>

    <!-- 套餐变更 - 全屏模态 (风格D) -->
    <a-modal
      v-model:visible="planModalVisible"
      title="套餐变更"
      :fullscreen="true"
      :mask-closable="false"
      @before-ok="submitPlanChange"
      @cancel="planModalVisible = false"
    >
      <a-form :model="planForm" layout="vertical" style="max-width: 600px">
        <a-descriptions :column="2" bordered size="small">
          <a-descriptions-item label="租户">{{ currentRecord?.company_name }}</a-descriptions-item>
          <a-descriptions-item label="当前套餐">
            <a-tag :color="getPlanColor(currentRecord?.plan_code)">{{ currentRecord?.plan_name }}</a-tag>
          </a-descriptions-item>
        </a-descriptions>

        <a-divider>选择新套餐</a-divider>
        <a-radio-group v-model="planForm.new_plan_id" direction="vertical">
          <a-card
            v-for="p in plans"
            :key="p.id"
            :class="{ 'plan-card-selected': planForm.new_plan_id === p.id }"
            class="plan-card"
            hoverable
            @click="planForm.new_plan_id = p.id"
          >
            <a-radio :value="p.id">
              <strong>{{ p.plan_name }}</strong>
              <span style="margin-left: 16px; color: #1650ff; font-size: 18px; font-weight: 700">
                ¥{{ p.price_yearly }}/年
              </span>
              <div style="margin-top: 8px; padding-left: 24px; color: #4e5969; font-size: 13px">
                设备上限: {{ p.device_quota === -1 ? '不限' : p.device_quota + '台' }} |
                用户上限: {{ p.user_quota === -1 ? '不限' : p.user_quota + '人' }}
              </div>
            </a-radio>
          </a-card>
        </a-radio-group>

        <a-divider>生效方式</a-divider>
        <a-form-item label="生效类型">
          <a-radio-group v-model="planForm.effective_type">
            <a-radio value="immediate">立即生效</a-radio>
            <a-radio value="end_of_cycle">周期末生效</a-radio>
          </a-radio-group>
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 禁用确认 -->
    <a-modal
      v-model:visible="disableModalVisible"
      title="禁用租户"
      @before-ok="submitDisable"
      @cancel="disableModalVisible = false"
    >
      <a-result status="warning" title="确定要禁用该租户吗？">
        <template #subtitle>
          禁用后，该租户下的所有用户将无法登录，但数据会保留。
        </template>
      </a-result>
      <a-form :model="disableForm" layout="vertical" style="max-width: 400px; margin-top: 16px">
        <a-form-item label="禁用原因（可选）">
          <a-textarea v-model="disableForm.reason" placeholder="请输入禁用原因" :rows="3" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import axios from 'axios'

const loading = ref(false)
const searchKey = ref('')
const showFilter = ref(false)
const filterStatus = ref('')
const filterPlan = ref('')
const filterExpire = ref('')
const editModalVisible = ref(false)
const planModalVisible = ref(false)
const disableModalVisible = ref(false)
const currentRecord = ref<any>(null)
const editFormRef = ref()

const plans = ref<any[]>([])
const editForm = ref<any>({})
const planForm = ref({ new_plan_id: null as number | null, effective_type: 'immediate' })
const disableForm = ref({ reason: '' })

const editFormRules = {
  company_name: [{ required: true, message: '请输入公司名称' }],
  contact_name: [{ required: true, message: '请输入联系人' }],
  contact_phone: [{ required: true, message: '请输入联系电话' }],
}

const columns = [
  { title: '租户编号', slotName: 'tenantCode', width: 150 },
  { title: '公司名称', dataIndex: 'company_name' },
  { title: '联系人', dataIndex: 'contact_name', width: 120 },
  { title: '联系电话', dataIndex: 'contact_phone', width: 140 },
  { title: '套餐', slotName: 'planInfo', width: 120 },
  { title: '用户数/配额', slotName: 'userQuota', width: 120 },
  { title: '设备数/配额', slotName: 'deviceQuota', width: 120 },
  { title: '到期时间', slotName: 'expireTime', width: 140 },
  { title: '状态', slotName: 'statusBadge', width: 100 },
  { title: '创建时间', slotName: 'createdAt', width: 120 },
  { title: '操作', slotName: 'actions', width: 240, fixed: 'right' },
]

// Mock data
const mockData = ref([
  {
    id: 1,
    tenant_code: 'TEN-000001',
    company_name: '深圳市智联科技有限公司',
    contact_name: '张经理',
    contact_phone: '13800138000',
    contact_email: 'zhang@zhilian.com',
    plan_name: '企业版',
    plan_code: 'enterprise',
    plan_id: 4,
    user_count: 23,
    user_quota: 50,
    device_count: 128,
    device_quota: 500,
    status: 'active',
    expire_time: '2026-06-20T00:00:00Z',
    created_at: '2026-01-15T08:00:00Z',
    industry: '互联网',
    company_size: '100-500人',
    tenant_logo: '',
    custom_domain: 'zhilian.mdm.com',
  },
  {
    id: 2,
    tenant_code: 'TEN-000002',
    company_name: '广州云智数据服务公司',
    contact_name: '李总监',
    contact_phone: '13900139000',
    contact_email: 'li@yunzhi.com',
    plan_name: '专业版',
    plan_code: 'pro',
    plan_id: 3,
    user_count: 12,
    user_quota: 20,
    device_count: 87,
    device_quota: 200,
    status: 'active',
    expire_time: '2026-04-10T00:00:00Z',
    created_at: '2026-02-01T10:00:00Z',
    industry: '金融科技',
    company_size: '50-100人',
    tenant_logo: '',
    custom_domain: 'yunzhi.mdm.com',
  },
  {
    id: 3,
    tenant_code: 'TEN-000003',
    company_name: '杭州创新科技集团',
    contact_name: '王总',
    contact_phone: '13700137000',
    contact_email: 'wang@cxtech.com',
    plan_name: '企业版',
    plan_code: 'enterprise',
    plan_id: 4,
    user_count: 5,
    user_quota: 50,
    device_count: 12,
    device_quota: 500,
    status: 'active',
    expire_time: '2026-03-25T00:00:00Z',
    created_at: '2026-03-18T14:00:00Z',
    industry: '制造业',
    company_size: '201-500人',
    tenant_logo: '',
    custom_domain: '',
  },
  {
    id: 4,
    tenant_code: 'TEN-000004',
    company_name: '北京华信信息技术有限公司',
    contact_name: '赵主管',
    contact_phone: '13600136000',
    contact_email: 'zhao@huaxin.com',
    plan_name: '基础版',
    plan_code: 'basic',
    plan_id: 2,
    user_count: 9,
    user_quota: 10,
    device_count: 48,
    device_quota: 50,
    status: 'suspended',
    expire_time: '2026-03-25T00:00:00Z',
    created_at: '2025-12-01T09:00:00Z',
    industry: '互联网',
    company_size: '50-100人',
    tenant_logo: '',
    custom_domain: 'huaxin.mdm.com',
  },
  {
    id: 5,
    tenant_code: 'TEN-000005',
    company_name: '上海鼎新数据工作室',
    contact_name: '陈老板',
    contact_phone: '13500135000',
    contact_email: 'chen@dxdata.com',
    plan_name: '免费版',
    plan_code: 'free',
    plan_id: 1,
    user_count: 3,
    user_quota: 5,
    device_count: 15,
    device_quota: 10,
    status: 'expired',
    expire_time: '2026-03-01T00:00:00Z',
    created_at: '2025-09-01T11:00:00Z',
    industry: '零售',
    company_size: '1-50人',
    tenant_logo: '',
    custom_domain: '',
  },
])

const mockPlans = [
  { id: 1, plan_name: '免费版', plan_code: 'free', price_yearly: 0, device_quota: 10, user_quota: 5 },
  { id: 2, plan_name: '基础版', plan_code: 'basic', price_yearly: 299, device_quota: 50, user_quota: 10 },
  { id: 3, plan_name: '专业版', plan_code: 'pro', price_yearly: 799, device_quota: 200, user_quota: 20 },
  { id: 4, plan_name: '企业版', plan_code: 'enterprise', price_yearly: 1999, device_quota: 500, user_quota: 50 },
]

const filteredData = computed(() => {
  let data = mockData.value
  if (filterStatus.value) {
    data = data.filter(item => item.status === filterStatus.value)
  }
  if (filterPlan.value) {
    data = data.filter(item => item.plan_code === filterPlan.value)
  }
  if (filterExpire.value === 'expired') {
    data = data.filter(item => isExpired(item.expire_time))
  } else if (filterExpire.value === 'expiring_soon') {
    data = data.filter(item => isExpiringSoon(item.expire_time))
  }
  if (searchKey.value) {
    data = data.filter(item =>
      item.company_name.includes(searchKey.value) ||
      item.contact_name.includes(searchKey.value)
    )
  }
  return data
})

const getStatusColor = (status: string) => {
  const map: Record<string, string> = {
    pending: 'orange',
    active: 'green',
    suspended: 'yellow',
    expired: 'red',
  }
  return map[status] || 'gray'
}

const getStatusText = (status: string) => {
  const map: Record<string, string> = {
    pending: '待审核',
    active: '正常',
    suspended: '已禁用',
    expired: '已到期',
  }
  return map[status] || status
}

const getPlanColor = (planCode: string) => {
  const map: Record<string, string> = {
    free: 'default',
    basic: 'gray',
    pro: 'arcoblue',
    enterprise: 'purple',
  }
  return map[planCode] || 'gray'
}

const formatDate = (dateStr: string | null) => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleDateString('zh-CN')
}

const isExpiringSoon = (dateStr: string | null) => {
  if (!dateStr) return false
  const diff = new Date(dateStr).getTime() - Date.now()
  return diff > 0 && diff < 7 * 24 * 60 * 60 * 1000
}

const isExpired = (dateStr: string | null) => {
  if (!dateStr) return false
  return new Date(dateStr).getTime() < Date.now()
}

const loadTenants = async () => {
  loading.value = true
  try {
    // API: GET /api/v1/admin/tenants
    // const res = await axios.get('/api/v1/admin/tenants', { params: { ... } })
    plans.value = mockPlans
  } catch {
    Message.error('加载租户列表失败')
  } finally {
    loading.value = false
  }
}

const resetFilter = () => {
  filterStatus.value = ''
  filterPlan.value = ''
  filterExpire.value = ''
  searchKey.value = ''
  loadTenants()
}

const openEditModal = (record: any) => {
  currentRecord.value = record
  editForm.value = { ...record }
  editFormRef.value?.clearValidate()
  editModalVisible.value = true
}

const submitEdit = async (done: (val: boolean) => void) => {
  try {
    // API: PUT /api/v1/admin/tenants/{id}
    // await axios.put(`/api/v1/admin/tenants/${currentRecord.value.id}`, editForm.value)
    Object.assign(currentRecord.value, editForm.value)
    Message.success('租户信息已保存')
    editModalVisible.value = false
    loadTenants()
    done(true)
  } catch {
    Message.error('保存失败')
    done(false)
  }
}

const openPlanModal = (record: any) => {
  currentRecord.value = record
  planForm.value.new_plan_id = record.plan_id
  planForm.value.effective_type = 'immediate'
  planModalVisible.value = true
}

const submitPlanChange = async (done: (val: boolean) => void) => {
  if (!planForm.value.new_plan_id) {
    Message.warning('请选择新套餐')
    done(false)
    return
  }
  try {
    // API: POST /api/v1/admin/tenants/{id}/change-plan
    // await axios.post(`/api/v1/admin/tenants/${currentRecord.value.id}/change-plan`, planForm.value)
    currentRecord.value.plan_id = planForm.value.new_plan_id
    const newPlan = plans.value.find(p => p.id === planForm.value.new_plan_id)
    if (newPlan) {
      currentRecord.value.plan_name = newPlan.plan_name
      currentRecord.value.plan_code = newPlan.plan_code
      currentRecord.value.device_quota = newPlan.device_quota
      currentRecord.value.user_quota = newPlan.user_quota
    }
    Message.success('套餐变更已提交')
    planModalVisible.value = false
    done(true)
  } catch {
    Message.error('套餐变更失败')
    done(false)
  }
}

const openDisableConfirm = (record: any) => {
  currentRecord.value = record
  disableForm.value.reason = ''
  disableModalVisible.value = true
}

const submitDisable = async (done: (val: boolean) => void) => {
  try {
    // API: PUT /api/v1/admin/tenants/{id}/suspend
    // await axios.put(`/api/v1/admin/tenants/${currentRecord.value.id}/suspend`, { reason: disableForm.value.reason })
    currentRecord.value.status = 'suspended'
    Message.warning(`已禁用租户 ${currentRecord.value.company_name}`)
    disableModalVisible.value = false
    done(true)
  } catch {
    Message.error('操作失败')
    done(false)
  }
}

const handleActivate = async (record: any) => {
  try {
    // API: PUT /api/v1/admin/tenants/{id}/activate
    // await axios.put(`/api/v1/admin/tenants/${record.id}/activate`)
    record.status = 'active'
    Message.success(`已启用租户 ${record.company_name}`)
  } catch {
    Message.error('操作失败')
  }
}

onMounted(() => {
  loadTenants()
})
</script>

<style scoped>
.page-container {
  padding: 20px 24px;
  min-height: calc(100vh - 64px);
  background: #f5f7fa;
}

.breadcrumb {
  margin-bottom: 16px;
}
.search-bar {
  margin-bottom: 12px;
}
.action-bar {
  margin-bottom: 16px;
}
.general-card {
  border-radius: 8px;
}
.text-danger {
  color: #f53f3f;
  font-weight: 500;
}
.text-warning {
  color: #ff7d00;
  font-weight: 500;
}
.plan-card {
  margin-bottom: 12px;
  cursor: pointer;
}
.plan-card-selected {
  border-color: #1650ff;
  background: #f2f3f5;
}
</style>