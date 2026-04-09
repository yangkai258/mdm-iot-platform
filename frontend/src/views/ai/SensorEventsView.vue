<template>
  <div class="page-container">
    <a-card class="general-card" title="传感器事�?>
      <template #extra>
        <a-button @click="loadData"><icon-refresh />刷新</a-button>
      </template>
      <div class="search-form">
        <a-form :model="form" layout="inline">
          <a-form-item label="设备ID"><a-input v-model="form.device_id" placeholder="请输�? /></a-form-item>
          <a-form-item label="事件类型">
            <a-select v-model="form.event_type" placeholder="选择类型" allow-clear style="width: 140px">
              <a-option value="temperature">温度</a-option>
              <a-option value="humidity">湿度</a-option>
              <a-option value="battery">电池</a-option>
              <a-option value="motion">运动</a-option>
              <a-option value="heartbeat">心跳</a-option>
            </a-select>
          </a-form-item>
          <a-form-item label="时间范围"><a-range-picker v-model="form.time_range" style="width: 240px" /></a-form-item>
          <a-form-item><a-button type="primary" @click="handleSearch">查询</a-button><a-button @click="handleReset">重置</a-button></a-form-item>
        </a-form>
      </div>
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id" />
    </a-card>
    <a-modal v-model:visible="detailVisible" title="事件详情" :width="560" :footer="null">
      <a-descriptions :column="1" bordered>
        <a-descriptions-item label="事件ID">{{ currentEvent?.id }}</a-descriptions-item>
        <a-descriptions-item label="设备ID">{{ currentEvent?.device_id }}</a-descriptions-item>
        <a-descriptions-item label="事件类型">{{ currentEvent?.event_type }}</a-descriptions-item>
        <a-descriptions-item label="传感器数�?>{{ currentEvent?.sensor_data }}</a-descriptions-item>
        <a-descriptions-item label="触发时间">{{ currentEvent?.created_at }}</a-descriptions-item>
      </a-descriptions>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const loading = ref(false)
const data = ref([])
const detailVisible = ref(false)
const currentEvent = ref(null)
const form = reactive({ device_id: '', event_type: '', time_range: [] })
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })

const columns = [
  { title: '事件ID', dataIndex: 'id', width: 80 },
  { title: '设备ID', dataIndex: 'device_id', width: 120 },
  { title: '设备名称', dataIndex: 'device_name', width: 140 },
  { title: '事件类型', dataIndex: 'event_type', width: 100 },
  { title: '传感器数�?, dataIndex: 'sensor_data', ellipsis: true },
  { title: '触发时间', dataIndex: 'created_at', width: 170 }
]

const loadData = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const params = new URLSearchParams({ page: pagination.current, page_size: pagination.pageSize })
    if (form.device_id) params.append('device_id', form.device_id)
    if (form.event_type) params.append('event_type', form.event_type)
    const res = await fetch(`/api/v1/ai/sensor-events?${params}`, { headers: { 'Authorization': `Bearer ${token}` } }).then(r => r.json())
    if (res.code === 0) { data.value = res.data?.list || []; pagination.total = res.data?.total || 0 }
    else { data.value = [] }
  } catch (e) { Message.error('加载失败') } finally { loading.value = false }
}

const handleSearch = () => { pagination.current = 1; loadData() }
const handleReset = () => { Object.assign(form, { device_id: '', event_type: '', time_range: [] }); loadData() }
const onPageChange = (page) => { pagination.current = page; loadData() }

onMounted(() => loadData())
</script>

<style scoped>
.page-container { padding: 16px; }
.search-form { margin-bottom: 16px; padding: 16px; background: var(--color-fill-lightest); border-radius: 4px; }
</style>