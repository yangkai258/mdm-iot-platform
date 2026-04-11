<template>
  <div class="alert-list-container">
    <a-card class="general-card">
      <template #title><span class="card-title">告警查询</span></template>
      <a-row :gutter="16">
        <a-col :flex="1">
          <a-form :model="searchForm" layout="vertical" size="small">
            <a-row :gutter="16">
              <a-col :span="8">
                <a-form-item label="关键词">
                  <a-input v-model="searchForm.keyword" placeholder="输入告警关键词" allow-clear />
                </a-form-item>
              </a-col>
              <a-col :span="8">
                <a-form-item label="设备ID">
                  <a-input v-model="searchForm.device_id" placeholder="设备ID" allow-clear />
                </a-form-item>
              </a-col>
              <a-col :span="8">
                <a-form-item label="状态">
                  <a-select v-model="searchForm.status" placeholder="全部" allow-clear>
                    <a-option :value="1">未处理</a-option>
                    <a-option :value="0">已处理</a-option>
                  </a-select>
                </a-form-item>
              </a-col>
            </a-row>
          </a-form>
        </a-col>
        <a-divider style="height: 84px" direction="vertical" />
        <a-col :flex="'86px'" style="text-align: right">
          <a-space direction="vertical" :size="18">
            <a-button type="primary" @click="handleSearch">
              <template #icon><icon-search /></template>
              查询
            </a-button>
            <a-button @click="handleReset">
              <template #icon><icon-refresh /></template>
              重置
            </a-button>
          </a-space>
        </a-col>
      </a-row>
    </a-card>

    <a-card class="general-card" style="margin-top: 16px">
      <template #title><span class="card-title">告警记录</span></template>
      <template #extra>
        <a-button type="primary" @click="handleCreate">
          <template #icon><icon-plus /></template>
          新建
        </a-button>
      </template>
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="paginationConfig" @page-change="onPageChange" row-key="id">
        <template #severity="{ record }">
          <a-tag :color="getSeverityColor(record.severity)">{{ getSeverityText(record.severity) }}</a-tag>
        </template>
        <template #status="{ record }">
          <a-tag :color="record.status === 1 ? 'blue' : 'green'">{{ record.status === 1 ? '未处理' : '已处理' }}</a-tag>
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:visible="modalVisible" :title="modalTitle" @ok="handleSubmit" :width="520">
      <a-form :model="form" layout="vertical">
        <a-form-item label="告警名称">
          <a-input v-model="form.name" placeholder="输入告警名称" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const loading = ref(false)
const modalVisible = ref(false)
const isEdit = ref(false)
const modalTitle = computed(() => isEdit.value ? '编辑告警' : '新建告警')

const searchForm = reactive({ keyword: '', device_id: '', status: undefined })
const form = reactive({ id: null, name: '' })
const data = ref([])
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })
const paginationConfig = computed(() => ({ current: pagination.current, pageSize: pagination.pageSize, total: pagination.total, showTotal: true }))

const columns = [
  { title: 'ID', dataIndex: 'id', width: 80 },
  { title: '设备ID', dataIndex: 'device_id', width: 160 },
  { title: '告警消息', dataIndex: 'message', ellipsis: true },
  { title: '级别', slotName: 'severity', width: 100 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '时间', dataIndex: 'created_at', width: 170 }
]

const getSeverityColor = (s) => ({ warning: 'orange', error: 'red', info: 'blue' }[s] || 'gray')
const getSeverityText = (s) => ({ warning: '警告', error: '错误', info: '提示' }[s] || s)

const getMockData = () => [
  { id: 1, device_id: 'DEV001', message: '设备电量低于20%', severity: 'warning', status: 1, created_at: '2026-03-29 10:00:00' },
  { id: 2, device_id: 'DEV002', message: '设备离线超过5分钟', severity: 'error', status: 0, created_at: '2026-03-29 09:30:00' },
  { id: 3, device_id: 'DEV003', message: '温度异常', severity: 'warning', status: 1, created_at: '2026-03-29 08:00:00' }
]

const loadData = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const params = new URLSearchParams({ page: pagination.current, page_size: pagination.pageSize })
    if (searchForm.keyword) params.append('keyword', searchForm.keyword)
    if (searchForm.device_id) params.append('device_id', searchForm.device_id)
    if (searchForm.status !== undefined) params.append('status', searchForm.status)
    const res = await fetch(`/api/v1/alerts?${params}`, {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    const resData = await res.json()
    if (resData.code === 0) {
      data.value = resData.data?.list || []
      pagination.total = resData.data?.total || 0
    }
  } catch (e) {
    data.value = getMockData()
    pagination.total = data.value.length
  } finally {
    loading.value = false
  }
}

const handleSearch = () => { pagination.current = 1; loadData() }
const handleReset = () => { searchForm.keyword = ''; searchForm.device_id = ''; searchForm.status = undefined; pagination.current = 1; loadData() }
const onPageChange = (page) => { pagination.current = page; loadData() }
const handleCreate = () => { isEdit.value = false; Object.assign(form, { id: null, name: '' }); modalVisible.value = true }
const handleSubmit = async () => {
  try {
    const token = localStorage.getItem('token')
    const method = isEdit.value ? 'PUT' : 'POST'
    const url = isEdit.value ? `/api/v1/alerts/${form.id}` : '/api/v1/alerts'
    await fetch(url, { method, headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' }, body: JSON.stringify(form) })
    Message.success('操作成功')
    modalVisible.value = false
    loadData()
  } catch (e) { Message.error('操作失败') }
}

onMounted(() => { loadData() })
</script>

<style scoped>
.alert-list-container { padding: 20px 24px; min-height: calc(100vh - 64px); background: #f5f7fa; }
.general-card { border-radius: 8px; }
.card-title { font-weight: 600; font-size: 15px; }
</style>
