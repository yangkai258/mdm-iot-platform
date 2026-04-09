<template>
  <div class="pro-page-container">

    <!-- 搜索表单 -->
    <div class="pro-search-bar">
      <a-space wrap>
        <a-input-search v-model="form.user_id" placeholder="用户ID" style="width: 160px" @search="loadData" search-button />
        <a-select v-model="form.status" placeholder="订阅状态" allow-clear style="width: 150px" @change="loadData">
          <a-option value="active">生效中</a-option>
          <a-option value="expired">已过期</a-option>
          <a-option value="cancelled">已取消</a-option>
          <a-option value="suspended">已暂停</a-option>
        </a-select>
        <a-input-search v-model="form.plan_name" placeholder="计划名称" style="width: 160px" @search="loadData" search-button />
        <a-button @click="handleReset">重置</a-button>
      </a-space>
    </div>

    <!-- 操作按钮 -->
    <div class="pro-action-bar">
      <a-button type="primary" @click="showCreateModal">新建订阅</a-button>
      <a-button @click="loadData">刷新</a-button>
    </div>

    <!-- 订阅列表 -->
    <div class="pro-content-area">
      <a-table
        :columns="columns"
        :data="subscriptions"
        :loading="loading"
        :pagination="pagination"
        row-key="id"
        @page-change="handlePageChange"
      >
        <template #plan_type="{ record }">
          <a-tag :color="record.plan_type === 'yearly' ? 'blue' : 'green'">
            {{ record.plan_type === 'yearly' ? '年付' : '月付' }}
          </a-tag>
        </template>
        <template #status="{ record }">
          <a-tag :color="getStatusColor(record.status)">{{ getStatusText(record.status) }}</a-tag>
        </template>
        <template #price="{ record }">
          <span class="price">¥{{ record.price }}</span>
        </template>
        <template #validity="{ record }">
          <span :class="{ 'text-expired': record.status === 'expired' }">
            {{ formatDate(record.start_date) }} ~ {{ formatDate(record.end_date) }}
          </span>
        </template>
        <template #auto_renew="{ record }">
          <a-switch :model-value="record.auto_renew" disabled size="small" />
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="showDetailModal(record)">详情</a-button>
            <a-button type="text" size="small" @click="handleRenew(record)" v-if="record.status === 'expired' || record.status === 'suspended'">续费</a-button>
            <a-button type="text" size="small" @click="handleCancelRenewal(record)" v-if="record.auto_renew && record.status === 'active'">取消自动续费</a-button>
            <a-button type="text" size="small" @click="handleResume(record)" v-if="record.status === 'suspended'">恢复</a-button>
            <a-button type="text" size="small" status="danger" @click="handleCancel(record)" v-if="record.status === 'active' || record.status === 'suspended'">取消订阅</a-button>
          </a-space>
        </template>
      </a-table>
    </div>

    <!-- 新建订阅弹窗 -->
    <a-modal v-model:visible="createModalVisible" title="新建订阅" :width="520" :loading="submitting" @before-ok="handleCreate" @cancel="createModalVisible = false">
      <a-form :model="createForm" layout="vertical">
        <a-form-item label="用户ID" required>
          <a-input-number v-model="createForm.user_id" :min="1" placeholder="请输入用户ID" style="width: 100%" />
        </a-form-item>
        <a-form-item label="计划名称" required>
          <a-input v-model="createForm.plan_name" placeholder="例如：年度高级版" />
        </a-form-item>
        <a-form-item label="计划类型" required>
          <a-select v-model="createForm.plan_type" placeholder="选择计划类型">
            <a-option value="monthly">月付</a-option>
            <a-option value="yearly">年付</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="价格（元）" required>
          <a-input-number v-model="createForm.price" :min="0" :precision="2" placeholder="例如：299.00" style="width: 100%" />
        </a-form-item>
        <a-form-item label="时长（天）">
          <a-input-number v-model="createForm.duration" :min="1" placeholder="月付30天，年付365天" style="width: 100%" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 续费弹窗 -->
    <a-modal v-model:visible="renewModalVisible" title="续费订阅" :width="480" :loading="submitting" @before-ok="handleRenewSubmit" @cancel="renewModalVisible = false">
      <a-form :model="renewForm" layout="vertical">
        <a-form-item label="订阅ID">
          <a-input :value="String(selectedSub?.id)" disabled />
        </a-form-item>
        <a-form-item label="当前到期日期">
          <a-input :value="formatDate(selectedSub?.end_date)" disabled />
        </a-form-item>
        <a-form-item label="续费时长（天）" required>
          <a-input-number v-model="renewForm.duration" :min="1" style="width: 100%" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 详情弹窗 -->
    <a-modal v-model:visible="detailModalVisible" title="订阅详情" :width="620" :footer="null">
      <a-descriptions :column="2" bordered v-if="currentSub">
        <a-descriptions-item label="订阅ID">{{ currentSub.id }}</a-descriptions-item>
        <a-descriptions-item label="用户ID">{{ currentSub.user_id }}</a-descriptions-item>
        <a-descriptions-item label="计划名称">{{ currentSub.plan_name }}</a-descriptions-item>
        <a-descriptions-item label="计划类型">
          <a-tag :color="currentSub.plan_type === 'yearly' ? 'blue' : 'green'">{{ currentSub.plan_type === 'yearly' ? '年付' : '月付' }}</a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="价格">¥{{ currentSub.price }}</a-descriptions-item>
        <a-descriptions-item label="时长">{{ currentSub.duration }}天</a-descriptions-item>
        <a-descriptions-item label="开始日期">{{ formatDate(currentSub.start_date) }}</a-descriptions-item>
        <a-descriptions-item label="结束日期">{{ formatDate(currentSub.end_date) }}</a-descriptions-item>
        <a-descriptions-item label="状态">
          <a-tag :color="getStatusColor(currentSub.status)">{{ getStatusText(currentSub.status) }}</a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="自动续费">
          <a-switch :model-value="currentSub.auto_renew" disabled />
        </a-descriptions-item>
        <a-descriptions-item label="续费次数">{{ currentSub.renew_count || 0 }}</a-descriptions-item>
        <a-descriptions-item label="最后续费时间">{{ formatDate(currentSub.last_renew_at) }}</a-descriptions-item>
      </a-descriptions>

      <!-- 续费历史 -->
      <a-divider>续费历史</a-divider>
      <a-table :columns="renewalLogColumns" :data="renewalLogs" size="small" :pagination="{ pageSize: 5 }" row-key="id">
        <template #status="{ record }">
          <a-tag :color="record.status === 'success' ? 'green' : 'red'">{{ record.status === 'success' ? '成功' : '失败' }}</a-tag>
        </template>
      </a-table>
      </a-table>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message, Modal } from '@arco-design/web-vue'

const loading = ref(false)
const submitting = ref(false)
const subscriptions = ref([])
const currentSub = ref(null)
const createModalVisible = ref(false)
const detailModalVisible = ref(false)
const renewModalVisible = ref(false)
const selectedSub = ref(null)
const renewalLogs = ref([])

const form = reactive({ user_id: '', status: '', plan_name: '' })
const createForm = reactive({ user_id: null, plan_name: '', plan_type: 'yearly', price: 0, duration: 365 })
const renewForm = reactive({ duration: 365 })

const pagination = reactive({ current: 1, pageSize: 20, total: 0 })

const columns = [
  { title: '订阅ID', dataIndex: 'id', width: 80 },
  { title: '用户ID', dataIndex: 'user_id', width: 80 },
  { title: '计划名称', dataIndex: 'plan_name', width: 140, ellipsis: true },
  { title: '计划类型', slotName: 'plan_type', width: 90 },
  { title: '价格', slotName: 'price', width: 100 },
  { title: '有效期', slotName: 'validity', width: 220 },
  { title: '状态', slotName: 'status', width: 90 },
  { title: '自动续费', slotName: 'auto_renew', width: 90 },
  { title: '操作', slotName: 'actions', width: 280, fixed: 'right' }
]

const renewalLogColumns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '续费时间', dataIndex: 'renewed_at', width: 170 },
  { title: '金额', dataIndex: 'amount', width: 100 },
  { title: '状态', slotName: 'status' }
]

const getToken = () => localStorage.getItem('token')
const getStatusColor = (s) => ({ active: 'green', expired: 'red', cancelled: 'gray', suspended: 'orange' }[s] || 'gray')
const getStatusText = (s) => ({ active: '生效中', expired: '已过期', cancelled: '已取消', suspended: '已暂停' }[s] || s)
const formatDate = (d) => d ? new Date(d).toLocaleString('zh-CN') : '-'

const loadData = async () => {
  loading.value = true
  try {
    const params = { page: pagination.current, page_size: pagination.pageSize }
    if (form.user_id) params.user_id = form.user_id
    if (form.status) params.status = form.status
    if (form.plan_name) params.plan_name = form.plan_name

    const res = await fetch(`/api/v1/subscriptions?${new URLSearchParams(params)}`, {
      headers: { 'Authorization': `Bearer ${getToken()}` }
    })
    const json = await res.json()
    if (json.code === 0) {
      subscriptions.value = json.data?.list || json.data || []
      pagination.total = json.data?.total || 0
    }
  } catch (e) {
    Message.error('加载订阅列表失败')
  } finally {
    loading.value = false
  }
}

const loadRenewalLogs = async (subId) => {
  try {
    const res = await fetch(`/api/v1/subscriptions/${subId}/renewal-logs?page=1&page_size=5`, {
      headers: { 'Authorization': `Bearer ${getToken()}` }
    })
    const json = await res.json()
    if (json.code === 0) {
      renewalLogs.value = json.data?.list || []
    }
  } catch (e) { /* silent */ }
}

const handlePageChange = (page) => { pagination.current = page; loadData() }

const showCreateModal = () => {
  Object.assign(createForm, { user_id: null, plan_name: '', plan_type: 'yearly', price: 0, duration: 365 })
  createModalVisible.value = true
}

const handleCreate = async (done) => {
  if (!createForm.user_id || !createForm.plan_name || !createForm.price) {
    Message.warning('请填写必填字段')
    done(false)
    return
  }
  submitting.value = true
  try {
    const res = await fetch('/api/v1/subscriptions', {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${getToken()}`, 'Content-Type': 'application/json' },
      body: JSON.stringify(createForm)
    })
    const json = await res.json()
    if (json.code === 0) {
      Message.success('订阅创建成功')
      createModalVisible.value = false
      loadData()
      done(true)
    } else {
      Message.error(json.message || '创建失败')
      done(false)
    }
  } catch (e) {
    Message.error('创建失败')
    done(false)
  } finally {
    submitting.value = false
  }
}

const showDetailModal = async (record) => {
  currentSub.value = record
  await loadRenewalLogs(record.id)
  detailModalVisible.value = true
}

const handleRenew = (record) => {
  selectedSub.value = record
  renewForm.duration = record.plan_type === 'yearly' ? 365 : 30
  renewModalVisible.value = true
}

const handleRenewSubmit = async (done) => {
  submitting.value = true
  try {
    const res = await fetch(`/api/v1/subscriptions/${selectedSub.value.id}/renew`, {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${getToken()}`, 'Content-Type': 'application/json' },
      body: JSON.stringify({ duration: renewForm.duration })
    })
    const json = await res.json()
    if (json.code === 0) {
      Message.success('续费成功')
      renewModalVisible.value = false
      loadData()
      done(true)
    } else {
      Message.error(json.message || '续费失败')
      done(false)
    }
  } catch (e) {
    Message.error('续费失败')
    done(false)
  } finally {
    submitting.value = false
  }
}

const handleCancelRenewal = async (record) => {
  try {
    const res = await fetch(`/api/v1/subscriptions/${record.id}/cancel-renewal`, {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${getToken()}` }
    })
    const json = await res.json()
    if (json.code === 0) {
      Message.success('已取消自动续费')
      loadData()
    } else {
      Message.error(json.message || '操作失败')
    }
  } catch (e) {
    Message.error('操作失败')
  }
}

const handleResume = async (record) => {
  try {
    const res = await fetch(`/api/v1/subscriptions/${record.id}/resume`, {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${getToken()}` }
    })
    const json = await res.json()
    if (json.code === 0) {
      Message.success('订阅已恢复')
      loadData()
    } else {
      Message.error(json.message || '恢复失败')
    }
  } catch (e) {
    Message.error('恢复失败')
  }
}

const handleCancel = async (record) => {
  Modal.warning({
    title: '确认取消订阅',
    content: `确定要取消该订阅吗？`,
    okText: '确认取消',
    onOk: async () => {
      try {
        const res = await fetch(`/api/v1/subscriptions/${record.id}`, {
          method: 'DELETE',
          headers: { 'Authorization': `Bearer ${getToken()}` }
        })
        const json = await res.json()
        if (json.code === 0) {
          Message.success('订阅已取消')
          loadData()
        } else {
          Message.error(json.message || '取消失败')
        }
      } catch (e) {
        Message.error('取消失败')
      }
    }
  })
}

const handleReset = () => {
  form.user_id = ''
  form.status = ''
  form.plan_name = ''
  pagination.current = 1
  loadData()
}

onMounted(() => { loadData() })
</script>

<style scoped>
.pro-page-container { padding: 20px 24px; min-height: calc(100vh - 64px); background: #f5f7fa; }
.pro-breadcrumb { margin-bottom: 16px; }
.pro-search-bar { margin-bottom: 12px; }
.pro-action-bar { margin-bottom: 16px; }
.pro-content-area { background: #fff; border-radius: 8px; padding: 20px; box-shadow: 0 1px 3px rgba(0,0,0,0.04); }
.price { font-weight: 600; color: #f53f3f; }
.text-expired { color: #f53f3f; }
</style>
