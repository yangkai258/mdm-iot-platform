<template>
  <div class="page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>DaaS 租赁管理</a-breadcrumb-item>
      <a-breadcrumb-item>租赁合同</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 统计卡片 -->
    <a-row :gutter="16" style="margin-top: 16px">
      <a-col :xs="24" :sm="12" :md="6">
        <a-card class="stat-card">
          <a-statistic title="合同总数" :value="stats.total">
            <template #prefix>📄</template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :xs="24" :sm="12" :md="6">
        <a-card class="stat-card">
          <a-statistic title="激活中" :value="stats.active" :value-style="{ color: '#0fc6c2' }">
            <template #prefix>✅</template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :xs="24" :sm="12" :md="6">
        <a-card class="stat-card">
          <a-statistic title="即将到期" :value="stats.expiring_soon" :value-style="{ color: '#ff7d00' }">
            <template #prefix>⏰</template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :xs="24" :sm="12" :md="6">
        <a-card class="stat-card">
          <a-statistic title="已过期" :value="stats.expired" :value-style="{ color: '#f53f3f' }">
            <template #prefix>❌</template>
          </a-statistic>
        </a-card>
      </a-col>
    </a-row>

    <!-- 搜索筛选区 -->
    <a-card class="section-card" style="margin-top: 16px">
      <a-space wrap>
        <a-input-search
          v-model="filterForm.keyword"
          placeholder="搜索合同号/设备ID/租户名"
          style="width: 260px"
          search-button
          @search="loadContracts"
          @change="e => !e.target.value && loadContracts()"
        />
        <a-select
          v-model="filterForm.status"
          placeholder="合同状态"
          style="width: 140px"
          allow-clear
          @change="loadContracts"
        >
          <a-option value="active">激活中</a-option>
          <a-option value="expiring">即将到期</a-option>
          <a-option value="expired">已过期</a-option>
          <a-option value="terminated">已终止</a-option>
        </a-select>
        <a-select
          v-model="filterForm.plan_type"
          placeholder="套餐类型"
          style="width: 140px"
          allow-clear
          @change="loadContracts"
        >
          <a-option value="basic">基础版</a-option>
          <a-option value="pro">专业版</a-option>
          <a-option value="enterprise">企业版</a-option>
        </a-select>
        <a-button @click="loadContracts">
          <template #icon><icon-refresh :spin="loading" /></template>
          刷新
        </a-button>
      </a-space>
      <div style="margin-top: 12px; text-align: right">
        <a-button type="primary" @click="showCreateModal">
          <template #icon><icon-plus /></template>
          新建合同
        </a-button>
      </div>
    </a-card>

    <!-- 合同列表 -->
    <a-card class="section-card" style="margin-top: 16px">
      <a-spin :loading="loading" tip="加载中...">
        <a-table
          :data="contracts"
          :pagination="{ pageSize: 10, total: contracts.length, showTotal: true, showPageSize: true }"
          :columns="columns"
          row-key="contract_id"
          stripe
        >
          <template #status="{ record }">
            <a-tag :color="getStatusColor(record.status)">
              {{ getStatusText(record.status) }}
            </a-tag>
          </template>
          <template #plan_type="{ record }">
            <a-tag :color="getPlanColor(record.plan_type)">
              {{ getPlanText(record.plan_type) }}
            </a-tag>
          </template>
          <template #monthly_fee="{ record }">
            ¥{{ record.monthly_fee?.toFixed(2) }}
          </template>
          <template #start_date="{ record }">
            {{ formatDate(record.start_date) }}
          </template>
          <template #end_date="{ record }">
            {{ formatDate(record.end_date) }}
          </template>
          <template #operations="{ record }">
            <a-space>
              <a-button type="text" size="small" @click="viewDetail(record)">详情</a-button>
              <a-button type="text" size="small" status="warning" @click="showRenewModal(record)">续费</a-button>
              <a-button type="text" size="small" status="danger" @click="handleTerminate(record)">终止</a-button>
            </a-space>
          </template>
        </a-table>
      </a-spin>
    </a-card>

    <!-- 新建合同弹窗 -->
    <a-modal v-model:visible="createVisible" title="新建租赁合同" :width="560" @before-ok="handleCreate" @cancel="createVisible = false">
      <a-form :model="createForm" layout="vertical">
        <a-form-item label="租户" field="tenant_id">
          <a-input v-model="createForm.tenant_id" placeholder="请输入租户ID" />
        </a-form-item>
        <a-form-item label="设备ID" field="device_id">
          <a-input v-model="createForm.device_id" placeholder="请输入设备ID" />
        </a-form-item>
        <a-form-item label="套餐类型" field="plan_type">
          <a-select v-model="createForm.plan_type" placeholder="请选择套餐">
            <a-option value="basic">基础版</a-option>
            <a-option value="pro">专业版</a-option>
            <a-option value="enterprise">企业版</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="租期（月）" field="duration_months">
          <a-input-number v-model="createForm.duration_months" :min="1" :max="60" placeholder="请输入租期" style="width: 100%" />
        </a-form-item>
        <a-form-item label="备注" field="remark">
          <a-textarea v-model="createForm.remark" placeholder="可选备注" :rows="3" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 合同详情弹窗 -->
    <a-modal v-model:visible="detailVisible" title="合同详情" :width="600" footer="null">
      <a-descriptions :column="2" bordered size="large">
        <a-descriptions-item label="合同ID">{{ currentContract?.contract_id }}</a-descriptions-item>
        <a-descriptions-item label="合同状态">
          <a-tag :color="getStatusColor(currentContract?.status)">{{ getStatusText(currentContract?.status) }}</a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="租户ID">{{ currentContract?.tenant_id }}</a-descriptions-item>
        <a-descriptions-item label="设备ID">{{ currentContract?.device_id }}</a-descriptions-item>
        <a-descriptions-item label="套餐类型">
          <a-tag :color="getPlanColor(currentContract?.plan_type)">{{ getPlanText(currentContract?.plan_type) }}</a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="月费用">¥{{ currentContract?.monthly_fee?.toFixed(2) }}</a-descriptions-item>
        <a-descriptions-item label="开始日期">{{ formatDate(currentContract?.start_date) }}</a-descriptions-item>
        <a-descriptions-item label="结束日期">{{ formatDate(currentContract?.end_date) }}</a-descriptions-item>
        <a-descriptions-item label="创建时间" :span="2">{{ formatDate(currentContract?.created_at) }}</a-descriptions-item>
        <a-descriptions-item label="备注" :span="2">{{ currentContract?.remark || '无' }}</a-descriptions-item>
      </a-descriptions>
    </a-modal>

    <!-- 续费弹窗 -->
    <a-modal v-model:visible="renewVisible" title="合同续费" @before-ok="handleRenew" @cancel="renewVisible = false">
      <a-form :model="renewForm" layout="vertical">
        <a-form-item label="合同ID">
          <a-input :model-value="currentContract?.contract_id" disabled />
        </a-form-item>
        <a-form-item label="续费时长（月）" field="months">
          <a-input-number v-model="renewForm.months" :min="1" :max="60" style="width: 100%" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { Message } from '@arco-design/web-vue'

const loading = ref(false)
const createVisible = ref(false)
const detailVisible = ref(false)
const renewVisible = ref(false)
const currentContract = ref(null)

const stats = reactive({
  total: 0,
  active: 0,
  expiring_soon: 0,
  expired: 0
})

const filterForm = reactive({
  keyword: '',
  status: '',
  plan_type: ''
})

const createForm = reactive({
  tenant_id: '',
  device_id: '',
  plan_type: '',
  duration_months: 1,
  remark: ''
})

const renewForm = reactive({
  months: 1
})

const contracts = ref([])

const columns = [
  { title: '合同ID', dataIndex: 'contract_id', width: 120 },
  { title: '租户', dataIndex: 'tenant_name', width: 120 },
  { title: '设备ID', dataIndex: 'device_id', width: 120 },
  { title: '套餐', dataIndex: 'plan_type', width: 100, slotName: 'plan_type' },
  { title: '月费用', dataIndex: 'monthly_fee', width: 100, slotName: 'monthly_fee' },
  { title: '开始日期', dataIndex: 'start_date', width: 120, slotName: 'start_date' },
  { title: '结束日期', dataIndex: 'end_date', width: 120, slotName: 'end_date' },
  { title: '状态', dataIndex: 'status', width: 100, slotName: 'status' },
  { title: '操作', slotName: 'operations', width: 180 }
]

const getStatusColor = (status) => {
  const map = { active: 'green', expiring: 'orange', expired: 'red', terminated: 'gray' }
  return map[status] || 'gray'
}

const getStatusText = (status) => {
  const map = { active: '激活中', expiring: '即将到期', expired: '已过期', terminated: '已终止' }
  return map[status] || status
}

const getPlanColor = (type) => {
  const map = { basic: 'arcoblue', pro: 'purple', enterprise: 'gold' }
  return map[type] || 'gray'
}

const getPlanText = (type) => {
  const map = { basic: '基础版', pro: '专业版', enterprise: '企业版' }
  return map[type] || type
}

const formatDate = (date) => {
  if (!date) return '-'
  return new Date(date).toLocaleDateString('zh-CN')
}

const loadContracts = async () => {
  loading.value = true
  // Mock data
  await new Promise(r => setTimeout(r, 600))
  contracts.value = [
    { contract_id: 'DC-202603001', tenant_id: 'T-001', tenant_name: '科技公司A', device_id: 'DEV-1001', plan_type: 'pro', monthly_fee: 299.00, start_date: '2026-01-01', end_date: '2026-12-31', status: 'active', remark: '长期合作客户', created_at: '2025-12-20' },
    { contract_id: 'DC-202603002', tenant_id: 'T-002', tenant_name: '物联网实验室', device_id: 'DEV-1002', plan_type: 'enterprise', monthly_fee: 899.00, start_date: '2025-10-01', end_date: '2026-03-28', status: 'expiring', remark: '即将到期需续费', created_at: '2025-09-15' },
    { contract_id: 'DC-202603003', tenant_id: 'T-003', tenant_name: '智能家居厂商', device_id: 'DEV-1003', plan_type: 'basic', monthly_fee: 99.00, start_date: '2025-06-01', end_date: '2026-05-31', status: 'active', remark: '', created_at: '2025-05-20' },
    { contract_id: 'DC-202603004', tenant_id: 'T-004', tenant_name: '教育机构B', device_id: 'DEV-1004', plan_type: 'basic', monthly_fee: 99.00, start_date: '2025-01-01', end_date: '2025-12-31', status: 'expired', remark: '已过期', created_at: '2024-12-15' },
    { contract_id: 'DC-202603005', tenant_id: 'T-005', tenant_name: '云服务公司', device_id: 'DEV-1005', plan_type: 'enterprise', monthly_fee: 1299.00, start_date: '2026-02-01', end_date: '2027-01-31', status: 'active', remark: '企业大客户', created_at: '2026-01-25' }
  ]
  stats.total = contracts.value.length
  stats.active = contracts.value.filter(c => c.status === 'active').length
  stats.expiring_soon = contracts.value.filter(c => c.status === 'expiring').length
  stats.expired = contracts.value.filter(c => c.status === 'expired').length
  loading.value = false
}

const showCreateModal = () => {
  Object.assign(createForm, { tenant_id: '', device_id: '', plan_type: '', duration_months: 1, remark: '' })
  createVisible.value = true
}

const handleCreate = async (done) => {
  if (!createForm.tenant_id || !createForm.device_id || !createForm.plan_type) {
    Message.error('请填写必填字段')
    done(false)
    return
  }
  await new Promise(r => setTimeout(r, 500))
  Message.success('合同创建成功')
  createVisible.value = false
  loadContracts()
  done(true)
}

const viewDetail = (record) => {
  currentContract.value = record
  detailVisible.value = true
}

const showRenewModal = (record) => {
  currentContract.value = record
  renewForm.months = 1
  renewVisible.value = true
}

const handleRenew = async (done) => {
  await new Promise(r => setTimeout(r, 500))
  Message.success('续费成功')
  renewVisible.value = false
  loadContracts()
  done(true)
}

const handleTerminate = async (record) => {
  currentContract.value = record
  await new Promise(r => setTimeout(r, 300))
  Message.warning('合同已终止')
  loadContracts()
}

loadContracts()
</script>

<style scoped>
.page-container { padding: 20px; }
.breadcrumb { margin-bottom: 12px; }
.stat-card { text-align: center; }
.section-card { margin-bottom: 16px; }
</style>
