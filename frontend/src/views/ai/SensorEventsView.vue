<template>
  <div class="page-container">
    <a-card class="general-card" title="浼犳劅鍣ㄤ簨锟?>
      <template #extra>
        <a-button @click="loadData"><icon-refresh />鍒锋柊</a-button>
      </template>
      <div class="search-form">
        <a-form :model="form" layout="inline">
          <a-form-item label="璁惧ID"><a-input v-model="form.device_id" placeholder="璇疯緭锟" /></a-form-item>
          <a-form-item label="浜嬩欢绫诲瀷">
            <a-select v-model="form.event_type" placeholder="閫夋嫨绫诲瀷" allow-clear style="width: 140px">
              <a-option value="temperature">娓╁害</a-option>
              <a-option value="humidity">婀垮害</a-option>
              <a-option value="battery">鐢垫睜</a-option>
              <a-option value="motion">杩愬姩</a-option>
              <a-option value="heartbeat">蹇冭烦</a-option>
            </a-select>
          </a-form-item>
          <a-form-item label="鏃堕棿鑼冨洿"><a-range-picker v-model="form.time_range" style="width: 240px" /></a-form-item>
          <a-form-item><a-button type="primary" @click="handleSearch">鏌ヨ</a-button><a-button @click="handleReset">閲嶇疆</a-button></a-form-item>
        </a-form>
      </div>
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id" />
    </a-card>
    <a-modal v-model:visible="detailVisible" title="浜嬩欢璇︽儏" :width="560" :footer="null">
      <a-descriptions :column="1" bordered>
        <a-descriptions-item label="浜嬩欢ID">{{ currentEvent?.id }}</a-descriptions-item>
        <a-descriptions-item label="璁惧ID">{{ currentEvent?.device_id }}</a-descriptions-item>
        <a-descriptions-item label="浜嬩欢绫诲瀷">{{ currentEvent?.event_type }}</a-descriptions-item>
        <a-descriptions-item label="浼犳劅鍣ㄦ暟锟?>{{ currentEvent?.sensor_data }}</a-descriptions-item>
        <a-descriptions-item label="瑙﹀彂鏃堕棿">{{ currentEvent?.created_at }}</a-descriptions-item>
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
  { title: '浜嬩欢ID', dataIndex: 'id', width: 80 },
  { title: '璁惧ID', dataIndex: 'device_id', width: 120 },
  { title: '璁惧鍚嶇О', dataIndex: 'device_name', width: 140 },
  { title: '浜嬩欢绫诲瀷', dataIndex: 'event_type', width: 100 },
  { title: '浼犳劅鍣ㄦ暟锟?, dataIndex: 'sensor_data', ellipsis: true },
  { title: '瑙﹀彂鏃堕棿', dataIndex: 'created_at', width: 170 }
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
  } catch (e) { Message.error('鍔犺浇澶辫触') } finally { loading.value = false }
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