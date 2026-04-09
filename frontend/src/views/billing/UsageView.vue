<template>
  <div class="page-container">
    <a-card class="general-card" title="用量计费">
      <template #extra>
        <a-button type="primary" @click="openAlertModal"><icon-bell />超量告警配置</a-button>
      </template>
      <div class="search-form">
        <a-form :model="form" layout="inline">
          <a-form-item label="时间范围"><a-range-picker v-model="form.time_range" style="width: 240px" /></a-form-item>
          <a-form-item><a-button type="primary" @click="loadData">查询</a-button><a-button @click="Object.keys(form).forEach(k => form[k] = ''); loadData()">重置</a-button></a-form-item>
        </a-form>
      </div>
      <!-- 用量概览卡片 -->
      <a-row :gutter="16" class="usage-cards">
        <a-col :span="6">
          <a-card class="stat-card">
            <div class="stat-label">API 调用量</div>
            <div class="stat-value">{{ summary.api_calls?.toLocaleString() || 0 }}</div>
            <div class="stat-unit">次</div>
          </a-card>
        </a-col>
        <a-col :span="6">
          <a-card class="stat-card">
            <div class="stat-label">存储用量</div>
            <div class="stat-value">{{ formatSize(summary.storage_bytes) }}</div>
            <div class="stat-unit">{{ summary.storage_limit ? `/ ${formatSize(summary.storage_limit)}` : '' }}</div>
          </a-card>
        </a-col>
        <a-col :span="6">
          <a-card class="stat-card">
            <div class="stat-label">流量用量</div>
            <div class="stat-value">{{ formatSize(summary.bandwidth_bytes) }}</div>
            <div class="stat-unit">{{ summary.bandwidth_limit ? `/ ${formatSize(summary.bandwidth_limit)}` : '' }}</div>
          </a-card>
        </a-col>
        <a-col :span="6">
          <a-card class="stat-card">
            <div class="stat-label">当前费用</div>
            <div class="stat-value" style="color: #f53f3f">¥{{ summary.total_cost?.toFixed(2) || '0.00' }}</div>
            <div class="stat-unit">本月</div>
          </a-card>
        </a-col>
      </a-row>
      <!-- 用量明细表格 -->
      <a-divider>用量明细</a-divider>
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id">
        <template #usage_type="{ record }">
          <a-tag :color="typeColor(record.usage_type)">{{ typeText(record.usage_type) }}</a-tag>
        </template>
        <template #amount="{ record }">
          <span v-if="record.usage_type === 'api'">{{ record.amount?.toLocaleString() }} 次</span>
          <span v-else>{{ formatSize(record.amount) }}</span>
        </template>
        <template #cost="{ record }">¥{{ record.cost?.toFixed(4) || '0' }}</template>
      </a-table>
    </a-card>
    <!-- 超量告警配置弹窗 -->
    <a-modal v-model:visible="alertVisible" title="超量告警配置" @before-ok="handleSaveAlert" :loading="submitting" :width="520">
      <a-form :model="alertForm" layout="vertical">
        <a-form-item label="API 调用告警阈值">
          <a-space>
            <a-input-number v-model="alertForm.api_calls_threshold" :min="0" style="width: 160px" />
            <span>次/月</span>
          </a-space>
        </a-form-item>
        <a-form-item label="存储告警阈值">
          <a-space>
            <a-input-number v-model="alertForm.storage_threshold_gb" :min="0" style="width: 160px" />
            <span>GB</span>
          </a-space>
        </a-form-item>
        <a-form-item label="流量告警阈值">
          <a-space>
            <a-input-number v-model="alertForm.bandwidth_threshold_gb" :min="0" style="width: 160px" />
            <span>GB</span>
          </a-space>
        </a-form-item>
        <a-form-item label="告警通知">
          <a-checkbox-group v-model="alertForm.notify_channels">
            <a-checkbox value="email">邮件</a-checkbox>
            <a-checkbox value="sms">短信</a-checkbox>
            <a-checkbox value="system">系统通知</a-checkbox>
          </a-checkbox-group>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const loading = ref(false)
const submitting = ref(false)
const data = ref([])
const summary = ref({})
const alertVisible = ref(false)
const form = reactive({ time_range: [] })
const alertForm = reactive({ api_calls_threshold: 100000, storage_threshold_gb: 100, bandwidth_threshold_gb: 500, notify_channels: ['email', 'system'] })
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })

const columns = [
  { title: 'ID', dataIndex: 'id', width: 70 },
  { title: '用量类型', slotName: 'usage_type', width: 100 },
  { title: '用量明细', slotName: 'amount', width: 160 },
  { title: '费用', slotName: 'cost', width: 100 },
  { title: '租户', dataIndex: 'tenant_name', width: 140 },
  { title: '时间', dataIndex: 'period', width: 170 },
  { title: '创建时间', dataIndex: 'created_at', width: 170 }
]

const typeColor = (t) => ({ api: 'arcoblue', storage: 'green', bandwidth: 'orange' }[t] || 'gray')
const typeText = (t) => ({ api: 'API调用', storage: '存储', bandwidth: '流量' }[t] || t)

const formatSize = (bytes) => {
  if (!bytes) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  let idx = 0
  let size = bytes
  while (size >= 1024 && idx < units.length - 1) { size /= 1024; idx++ }
  return `${size.toFixed(2)} ${units[idx]}`
}

const loadData = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const params = new URLSearchParams({ page: pagination.current, page_size: pagination.pageSize })
    const res = await fetch(`/api/v1/billing/usage?${params}`, { headers: { 'Authorization': `Bearer ${token}` } }).then(r => r.json())
    if (res.code === 0) {
      data.value = res.data?.list || []
      pagination.total = res.data?.total || 0
      summary.value = res.data?.summary || {}
    } else { data.value = [] }
  } catch (e) { Message.error('加载失败') } finally { loading.value = false }
}

const loadAlertConfig = async () => {
  try {
    const token = localStorage.getItem('token')
    const res = await fetch('/api/v1/billing/usage/alert-config', { headers: { 'Authorization': `Bearer ${token}` } }).then(r => r.json())
    if (res.code === 0) Object.assign(alertForm, res.data || {})
  } catch (e) { console.error('加载告警配置失败', e) }
}

const openAlertModal = () => { alertVisible.value = true }
const handleSaveAlert = async (done) => {
  submitting.value = true
  try {
    const token = localStorage.getItem('token')
    const res = await fetch('/api/v1/billing/usage/alert-config', { method: 'PUT', headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' }, body: JSON.stringify(alertForm) }).then(r => r.json())
    if (res.code === 0) { Message.success('保存成功'); alertVisible.value = false }
    else { Message.error(res.message || '保存失败') }
    done(true)
  } catch (e) { Message.error('保存失败'); done(false) } finally { submitting.value = false }
}

const onPageChange = (page) => { pagination.current = page; loadData() }
onMounted(() => { loadData(); loadAlertConfig() })
</script>

<style scoped>
.page-container { padding: 16px; }
.search-form { margin-bottom: 16px; padding: 16px; background: var(--color-fill-lightest); border-radius: 4px; }
.usage-cards { margin-bottom: 16px; }
.stat-card { text-align: center; }
.stat-label { font-size: 12px; color: #666; margin-bottom: 8px; }
.stat-value { font-size: 24px; font-weight: 600; color: #1a1a1a; }
.stat-unit { font-size: 12px; color: #999; margin-top: 4px; }
</style>
