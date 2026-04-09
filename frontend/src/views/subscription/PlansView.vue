<template>
  <div class="page-container">
    <a-card class="general-card" title="订阅套餐">
      <template #extra>
        <a-button type="primary" @click="openCreateModal"><icon-plus />新建套餐</a-button>
      </template>
      <div class="search-form">
        <a-form :model="form" layout="inline">
          <a-form-item label="套餐名称"><a-input v-model="form.plan_name" placeholder="请输入" /></a-form-item>
          <a-form-item label="状态">
            <a-select v-model="form.status" placeholder="选择状态" allow-clear style="width: 120px">
              <a-option value="active">生效中</a-option>
              <a-option value="inactive">停用</a-option>
            </a-select>
          </a-form-item>
          <a-form-item><a-button type="primary" @click="loadData">查询</a-button><a-button @click="Object.keys(form).forEach(k => form[k] = ''); loadData()">重置</a-button></a-form-item>
        </a-form>
      </div>
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id">
        <template #price="{ record }">
          <span style="color: #f53f3f; font-weight: 600">¥{{ record.price }}</span>
          <span style="color: #999; font-size: 12px">/{{ record.billing_cycle === 'monthly' ? '月' : '年' }}</span>
        </template>
        <template #features="{ record }">
          <a-tooltip :content="record.features?.join('\n') || '无'" placement="top">
            <span class="features-text">{{ record.features?.length || 0 }} 项功能</span>
          </a-tooltip>
        </template>
        <template #status="{ record }">
          <a-tag :color="record.status === 'active' ? 'green' : 'gray'">{{ record.status === 'active' ? '生效中' : '停用' }}</a-tag>
        </template>
        <template #actions="{ record }">
          <a-button type="text" size="small" @click="editPlan(record)">编辑</a-button>
          <a-button type="text" size="small" @click="viewSubscribers(record)">订阅者</a-button>
          <a-button type="text" size="small" status="danger" @click="deletePlan(record)">删除</a-button>
        </template>
      </a-table>
    </a-card>
    <!-- 新建/编辑套餐弹窗 -->
    <a-modal v-model:visible="formVisible" :title="isEdit ? '编辑套餐' : '新建套餐'" @before-ok="handleSubmit" :loading="submitting" :width="600">
      <a-form :model="planForm" layout="vertical">
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="套餐名称" required><a-input v-model="planForm.plan_name" placeholder="请输入套餐名称" /></a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="套餐代码"><a-input v-model="planForm.plan_code" placeholder="如: basic" /></a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="价格(元)" required><a-input-number v-model="planForm.price" :min="0" :precision="2" style="width: 100%" /></a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="计费周期" required>
              <a-select v-model="planForm.billing_cycle" placeholder="选择周期">
                <a-option value="monthly">月付</a-option>
                <a-option value="yearly">年付</a-option>
                <a-option value="one-time">一次性</a-option>
              </a-select>
            </a-form-item>
          </a-col>
        </a-row>
        <a-form-item label="套餐描述"><a-textarea v-model="planForm.description" :rows="2" placeholder="套餐描述" /></a-form-item>
        <a-form-item label="套餐功能">
          <a-select v-model="planForm.features" multiple placeholder="选择包含的功能" allow-create style="width: 100%">
            <a-option value="device_management">设备管理</a-option>
            <a-option value="ai_features">AI功能</a-option>
            <a-option value="ota_upgrade">OTA升级</a-option>
            <a-option value="analytics">数据分析</a-option>
            <a-option value="api_access">API访问</a-option>
            <a-option value="priority_support">优先支持</a-option>
            <a-option value="custom_branding">自定义品牌</a-option>
            <a-option value="multi_user">多用户</a-option>
          </a-select>
        </a-form-item>
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="设备上限"><a-input-number v-model="planForm.max_devices" :min="1" placeholder="无限制" style="width: 100%" /></a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="存储上限(GB)"><a-input-number v-model="planForm.max_storage_gb" :min="0" placeholder="无限制" style="width: 100%" /></a-form-item>
          </a-col>
        </a-row>
        <a-form-item label="状态"><a-switch v-model="planForm.is_active" /></a-form-item>
      </a-form>
    </a-modal>
    <!-- 订阅者列表 -->
    <a-modal v-model:visible="subscriberVisible" title="订阅者列表" :width="700" :footer="null">
      <a-table :columns="subColumns" :data="subscribers" :loading="subLoading" :pagination="subPagination" @page-change="onSubPageChange" row-key="id">
        <template #status="{ record }">
          <a-tag :color="record.subscription_status === 'active' ? 'green' : 'gray'">{{ record.subscription_status === 'active' ? '活跃' : '已过期' }}</a-tag>
        </template>
        <template #expires_at="{ record }">{{ record.expires_at || '永久' }}</template>
      </a-table>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const loading = ref(false)
const submitting = ref(false)
const subLoading = ref(false)
const data = ref([])
const subscribers = ref([])
const formVisible = ref(false)
const subscriberVisible = ref(false)
const isEdit = ref(false)
const selectedPlan = ref(null)
const form = reactive({ plan_name: '', status: '' })
const planForm = reactive({ id: null, plan_name: '', plan_code: '', price: 0, billing_cycle: 'monthly', description: '', features: [], max_devices: null, max_storage_gb: null, is_active: true })
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })
const subPagination = reactive({ current: 1, pageSize: 10, total: 0 })

const columns = [
  { title: 'ID', dataIndex: 'id', width: 70 },
  { title: '套餐名称', dataIndex: 'plan_name', width: 160 },
  { title: '套餐代码', dataIndex: 'plan_code', width: 100 },
  { title: '价格', slotName: 'price', width: 120 },
  { title: '计费周期', dataIndex: 'billing_cycle', width: 100 },
  { title: '功能数', slotName: 'features', width: 100 },
  { title: '设备上限', dataIndex: 'max_devices', width: 100 },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '创建时间', dataIndex: 'created_at', width: 170 },
  { title: '操作', slotName: 'actions', width: 160 }
]

const subColumns = [
  { title: '用户', dataIndex: 'user_name', width: 120 },
  { title: '邮箱', dataIndex: 'user_email', ellipsis: true },
  { title: '订阅状态', slotName: 'status', width: 90 },
  { title: '开始时间', dataIndex: 'started_at', width: 170 },
  { title: '到期时间', slotName: 'expires_at', width: 170 }
]

const loadData = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const params = new URLSearchParams({ page: pagination.current, page_size: pagination.pageSize })
    if (form.plan_name) params.append('plan_name', form.plan_name)
    if (form.status) params.append('status', form.status)
    const res = await fetch(`/api/v1/subscription/plans?${params}`, { headers: { 'Authorization': `Bearer ${token}` } }).then(r => r.json())
    if (res.code === 0) { data.value = res.data?.list || []; pagination.total = res.data?.total || 0 }
    else { data.value = [] }
  } catch (e) { Message.error('加载失败') } finally { loading.value = false }
}

const openCreateModal = () => { isEdit.value = false; Object.assign(planForm, { id: null, plan_name: '', plan_code: '', price: 0, billing_cycle: 'monthly', description: '', features: [], max_devices: null, max_storage_gb: null, is_active: true }); formVisible.value = true }
const editPlan = (record) => { isEdit.value = true; Object.assign(planForm, record); planForm.features = record.features || []; formVisible.value = true }

const handleSubmit = async (done) => {
  if (!planForm.plan_name) { Message.warning('请输入套餐名称'); done(false); return }
  submitting.value = true
  try {
    const token = localStorage.getItem('token')
    const url = isEdit.value ? `/api/v1/subscription/plans/${planForm.id}` : '/api/v1/subscription/plans'
    const res = await fetch(url, { method: isEdit.value ? 'PUT' : 'POST', headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' }, body: JSON.stringify(planForm) }).then(r => r.json())
    if (res.code === 0) { Message.success(isEdit.value ? '更新成功' : '创建成功'); formVisible.value = false; loadData() }
    else { Message.error(res.message || '操作失败') }
    done(true)
  } catch (e) { Message.error('操作失败'); done(false) } finally { submitting.value = false }
}

const deletePlan = async (record) => {
  try {
    const token = localStorage.getItem('token')
    await fetch(`/api/v1/subscription/plans/${record.id}`, { method: 'DELETE', headers: { 'Authorization': `Bearer ${token}` } })
    Message.success('删除成功'); loadData()
  } catch (e) { Message.error('删除失败') }
}

const viewSubscribers = async (record) => {
  selectedPlan.value = record
  subscriberVisible.value = true
  subLoading.value = true
  try {
    const token = localStorage.getItem('token')
    const params = new URLSearchParams({ plan_id: record.id, page: subPagination.current, page_size: subPagination.pageSize })
    const res = await fetch(`/api/v1/subscription/subscribers?${params}`, { headers: { 'Authorization': `Bearer ${token}` } }).then(r => r.json())
    if (res.code === 0) { subscribers.value = res.data?.list || []; subPagination.total = res.data?.total || 0 }
    else { subscribers.value = [] }
  } catch (e) { Message.error('加载失败') } finally { subLoading.value = false }
}

const onPageChange = (page) => { pagination.current = page; loadData() }
const onSubPageChange = (page) => { subPagination.current = page; viewSubscribers(selectedPlan.value) }

onMounted(() => loadData())
</script>

<style scoped>
.page-container { padding: 16px; }
.search-form { margin-bottom: 16px; padding: 16px; background: var(--color-fill-lightest); border-radius: 4px; }
.features-text { cursor: pointer; color: #165dff; }
</style>
