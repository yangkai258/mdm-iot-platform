<template>
  <div class="page-container">
    <a-card class="general-card" title="设备远程擦除">
      <template #extra>
        <a-button @click="loadData"><icon-refresh />刷新</a-button>
      </template>
      <div class="search-form">
        <a-form :model="form" layout="inline">
          <a-form-item label="设备ID"><a-input v-model="form.device_id" placeholder="请输�? /></a-form-item>
          <a-form-item label="状�?>
            <a-select v-model="form.status" placeholder="选择状�? allow-clear style="width: 140px">
              <a-option value="pending">待处�?/a-option>
              <a-option value="wiping">擦除�?/a-option>
              <a-option value="completed">已完�?/a-option>
              <a-option value="failed">失败</a-option>
            </a-select>
          </a-form-item>
          <a-form-item><a-button type="primary" @click="handleSearch">查询</a-button><a-button @click="handleReset">重置</a-button></a-form-item>
        </a-form>
      </div>
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id">
        <template #status="{ record }">
          <a-tag :color="getStatusColor(record.status)">{{ getStatusText(record.status) }}</a-tag>
        </template>
        <template #progress="{ record }">
          <a-progress :percent="record.progress || 0" :status="record.status === 'failed' ? 'exception' : 'normal'" size="small" />
        </template>
        <template #actions="{ record }">
          <a-button v-if="record.status === 'pending'" type="primary" size="small" @click="confirmWipe(record)">确认擦除</a-button>
          <a-button v-if="record.status === 'wiping'" type="text" size="small" @click="viewProgress(record)">查看进度</a-button>
        </template>
      </a-table>
    </a-card>
    <a-modal v-model:visible="confirmVisible" title="远程擦除确认" @before-ok="executeWipe" :loading="submitting">
      <a-result status="warning" title="即将执行远程擦除">
        <template #subtitle>
          <div>
            <p>设备ID: <strong>{{ selectedDevice?.device_id }}</strong></p>
            <p>设备名称: {{ selectedDevice?.device_name }}</p>
            <p style="color: #f53f3f">此操作将擦除设备上的所有数据，且无法恢复！</p>
          </div>
        </template>
      </a-result>
      <a-form layout="vertical" style="margin-top: 16px">
        <a-form-item label="擦除模式">
          <a-radio-group v-model="wipeMode">
            <a-radio value="fast">快速擦除（仅用户数据）</a-radio>
            <a-radio value="full">完全擦除（全部数据）</a-radio>
          </a-radio-group>
        </a-form-item>
      </a-form>
    </a-modal>
    <a-modal v-model:visible="progressVisible" title="擦除进度" :footer="null" :width="560">
      <a-descriptions :column="1" bordered>
        <a-descriptions-item label="设备ID">{{ selectedDevice?.device_id }}</a-descriptions-item>
        <a-descriptions-item label="当前状�?><a-tag :color="getStatusColor(selectedDevice?.status)">{{ getStatusText(selectedDevice?.status) }}</a-tag></a-descriptions-item>
        <a-descriptions-item label="擦除进度"><a-progress :percent="selectedDevice?.progress || 0" size="large" /></a-descriptions-item>
        <a-descriptions-item label="开始时�?>{{ selectedDevice?.started_at }}</a-descriptions-item>
        <a-descriptions-item label="预计完成">{{ selectedDevice?.estimated_time }}</a-descriptions-item>
      </a-descriptions>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const loading = ref(false)
const submitting = ref(false)
const data = ref([])
const confirmVisible = ref(false)
const progressVisible = ref(false)
const selectedDevice = ref(null)
const wipeMode = ref('fast')
const form = reactive({ device_id: '', status: '' })
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })

const columns = [
  { title: 'ID', dataIndex: 'id', width: 70 },
  { title: '设备ID', dataIndex: 'device_id', width: 120 },
  { title: '设备名称', dataIndex: 'device_name', width: 140 },
  { title: '擦除状�?, slotName: 'status', width: 100 },
  { title: '擦除进度', slotName: 'progress', width: 160 },
  { title: '发起�?, dataIndex: 'initiator', width: 100 },
  { title: '发起时间', dataIndex: 'created_at', width: 170 },
  { title: '操作', slotName: 'actions', width: 120 }
]

const getStatusColor = (status) => ({ pending: 'orange', wiping: 'arcoblue', completed: 'green', failed: 'red' }[status] || 'gray')
const getStatusText = (status) => ({ pending: '待处�?, wiping: '擦除�?, completed: '已完�?, failed: '失败' }[status] || status)

const loadData = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const params = new URLSearchParams({ page: pagination.current, page_size: pagination.pageSize })
    if (form.device_id) params.append('device_id', form.device_id)
    if (form.status) params.append('status', form.status)
    const res = await fetch(`/api/v1/device/remote-wipe?${params}`, { headers: { 'Authorization': `Bearer ${token}` } }).then(r => r.json())
    if (res.code === 0) { data.value = res.data?.list || []; pagination.total = res.data?.total || 0 }
    else { data.value = [] }
  } catch (e) { Message.error('加载失败') } finally { loading.value = false }
}

const handleSearch = () => { pagination.current = 1; loadData() }
const handleReset = () => { Object.assign(form, { device_id: '', status: '' }); loadData() }
const onPageChange = (page) => { pagination.current = page; loadData() }
const confirmWipe = (record) => { selectedDevice.value = record; confirmVisible.value = true }
const viewProgress = (record) => { selectedDevice.value = record; progressVisible.value = true }

const executeWipe = async (done) => {
  submitting.value = true
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`/api/v1/device/remote-wipe/${selectedDevice.value.id}/execute`, { method: 'POST', headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' }, body: JSON.stringify({ mode: wipeMode.value }) }).then(r => r.json())
    if (res.code === 0) { Message.success('擦除任务已下�?); confirmVisible.value = false; loadData() }
    else { Message.error(res.message || '操作失败') }
    done(true)
  } catch (e) { Message.error('操作失败'); done(false) } finally { submitting.value = false }
}

onMounted(() => loadData())
</script>

<style scoped>
.page-container { padding: 16px; }
.search-form { margin-bottom: 16px; padding: 16px; background: var(--color-fill-lightest); border-radius: 4px; }
</style>