<template>
  <div class="emotion-log-view">
    <a-card title="情绪日志">
      <template #extra>
        <a-space>
          <a-select v-model="filterType" placeholder="情绪类型" style="width: 120px" allow-clear>
            <a-option value="happy">开心</a-option>
            <a-option value="sad">难过</a-option>
            <a-option value="angry">生气</a-option>
            <a-option value="fear">害怕</a-option>
            <a-option value="neutral">平静</a-option>
          </a-select>
          <a-button type="primary" @click="loadData">刷新</a-button>
        </a-space>
      </template>
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="onPageChange">
        <template #device="{ record }">
          <a-tag>{{ record.device_id }}</a-tag>
        </template>
        <template #emotion="{ record }">
          <a-tag :color="getEmotionColor(record.emotion)">{{ record.emotion_label }}</a-tag>
        </template>
        <template #intensity="{ record }">
          <a-progress :percent="record.intensity" :color="getEmotionColor(record.emotion)" size="small" />
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'

const loading = ref(false)
const filterType = ref('')
const data = ref([])

const columns = [
  { title: '时间', dataIndex: 'created_at', width: 180 },
  { title: '设备', dataIndex: 'device_id', slotName: 'device' },
  { title: '情绪类型', dataIndex: 'emotion_label', slotName: 'emotion' },
  { title: '强度', dataIndex: 'intensity', slotName: 'intensity' },
  { title: '触发原因', dataIndex: 'trigger' }
]

const pagination = {
  total: 0,
  current: 1,
  pageSize: 20
}

function getEmotionColor(emotion) {
  const colors = { happy: 'green', sad: 'blue', angry: 'red', fear: 'orange', neutral: 'gray' }
  return colors[emotion] || 'gray'
}

async function loadData() {
  loading.value = true
  try {
    const res = await fetch('/api/v1/emotion/logs?page=' + pagination.current + '&page_size=' + pagination.pageSize + (filterType.value ? '&emotion=' + filterType.value : ''), {
      headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token') }
    })
    const json = await res.json()
    data.value = json.data?.list || []
    pagination.total = json.data?.total || 0
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

function onPageChange(page) {
  pagination.current = page
  loadData()
}

onMounted(() => loadData())
</script>
