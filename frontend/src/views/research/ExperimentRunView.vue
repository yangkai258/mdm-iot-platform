<template>
  <div class="container">
    <a-card class="general-card" title="实验运行">
      <template #extra>
        <a-button type="primary" @click="handleRun"><icon-play />运行</a-button>
      </template>
      <a-row :gutter="16" style="margin-bottom: 16px">
        <a-col :span="6"><a-statistic title="运行中" :value="stats.running" color="blue" /></a-col>
        <a-col :span="6"><a-statistic title="已完成" :value="stats.completed" color="green" /></a-col>
        <a-col :span="6"><a-statistic title="失败" :value="stats.failed" color="red" /></a-col>
      </a-row>
      <a-divider style="margin: 0 0 16px 0" />
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id" />
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import Breadcrumb from '@/components/Breadcrumb.vue'

const loading = ref(false)
const stats = reactive({ running: 0, completed: 0, failed: 0 })
const data = ref([])
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })
const columns = [
  { title: '实验ID', dataIndex: 'id', width: 80 },
  { title: '实验名称', dataIndex: 'name', width: 200 },
  { title: '状态', dataIndex: 'status', width: 100 },
  { title: '开始时间', dataIndex: 'started_at', width: 170 },
  { title: '耗时', dataIndex: 'duration', width: 100 }
]

const loadData = async () => {
  loading.value = true
  try {
    const res = await fetch('/api/v1/research/experiment-runs', { headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token') } }).then(r => r.json())
    data.value = res.data?.list || []
    pagination.total = data.value.length
  } catch { data.value = [] } finally { loading.value = false }
}
const handleRun = () => Message.success('实验已启动')
const onPageChange = (page) => { pagination.current = page; loadData() }
onMounted(() => loadData())
</script>
