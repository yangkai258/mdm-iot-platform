<template>
    <Breadcrumb :items="['Home','Console','']" />


  <div class="container">
    <a-card class="general-card" title="区域同步状态">
      <template #extra>
        <a-button @click="loadData"><icon-refresh />刷新</a-button>
      </template>
      <a-row :gutter="16" style="margin-bottom: 16px">
        <a-col :span="6"><a-statistic title="同步中" :value="stats.syncing" color="blue" /></a-col>
        <a-col :span="6"><a-statistic title="已同步" :value="stats.synced" color="green" /></a-col>
        <a-col :span="6"><a-statistic title="失败" :value="stats.failed" color="red" /></a-col>
      </a-row>
      <a-divider style="margin: 0 0 16px 0" />
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id" />
    </a-table>
  </a-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import Breadcrumb from '@/components/Breadcrumb.vue'

const loading = ref(false)
const stats = reactive({ syncing: 0, synced: 0, failed: 0 })
const data = ref([])
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })
const columns = [
  { title: '区域', dataIndex: 'region', width: 120 },
  { title: '同步类型', dataIndex: 'sync_type', width: 120 },
  { title: '状态', dataIndex: 'status', width: 90 },
  { title: '开始时间', dataIndex: 'started_at', width: 170 },
  { title: '完成时间', dataIndex: 'completed_at', width: 170 }
]

const loadData = async () => {
  loading.value = true
  try {
    const res = await fetch('/api/v1/globalization/region-sync-status', { headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token') } }).then(r => r.json())
    data.value = res.data?.list || []
    pagination.total = data.value.length
  } catch { data.value = [] } finally { loading.value = false }
}
const onPageChange = (page) => { pagination.current = page; loadData() }
onMounted(() => loadData())
</script>

