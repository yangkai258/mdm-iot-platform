<template>
    <Breadcrumb :items="['Home','Console','']" />


  <div class="container">
    <a-card class="general-card" title="离线数据同步">
      <template #extra>
        <a-button type="primary" @click="handleSync"><icon-refresh />立即同步</a-button>
      </template>
      <a-row :gutter="16">
        <a-col :span="6">
          <a-statistic title="待同步" :value="stats.pending" color="orange" />
        </a-col>
        <a-col :span="6">
          <a-statistic title="同步成功" :value="stats.success" color="green" />
        </a-col>
        <a-col :span="6">
          <a-statistic title="同步失败" :value="stats.failed" color="red" />
        </a-col>
        <a-col :span="6">
          <a-statistic title="最后同步" :value="stats.last_sync" />
        </a-col>
      </a-row>
      <a-divider style="margin: 16px 0" />
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id" />
    </a-table>
  </a-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import Breadcrumb from '@/components/Breadcrumb.vue'

const loading = ref(false)
const data = ref([])
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })
const stats = reactive({ pending: 0, success: 0, failed: 0, last_sync: '-' })
const columns = [
  { title: '记录ID', dataIndex: 'id', width: 80 },
  { title: '数据类型', dataIndex: 'data_type', width: 120 },
  { title: '设备ID', dataIndex: 'device_id', width: 120 },
  { title: '同步状态', dataIndex: 'status', width: 100 },
  { title: '创建时间', dataIndex: 'created_at', width: 170 }
]

const loadData = async () => {
  loading.value = true
  try {
    const res = await fetch('/api/v1/offline/sync-records', { headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token') } }).then(r => r.json())
    data.value = res.data?.list || []
    pagination.total = data.value.length
    if (res.data?.stats) Object.assign(stats, res.data.stats)
  } catch { data.value = [] } finally { loading.value = false }
}
const handleSync = () => { Message.success('同步任务已启动'); loadData() }
const onPageChange = (page) => { pagination.current = page; loadData() }
onMounted(() => loadData())
</script>

