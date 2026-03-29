<template>
  <div class="container">
    <a-card class="general-card" title="标签报表">
      <template #extra>
        <a-button @click="loadData"><icon-refresh />刷新</a-button>
      </template>
      <a-row :gutter="16" style="margin-bottom: 16px">
        <a-col :span="6"><a-statistic title="标签总数" :value="stats.total" /></a-col>
        <a-col :span="6"><a-statistic title="高频标签" :value="stats.high_freq" color="green" /></a-col>
        <a-col :span="6"><a-statistic title="低频标签" :value="stats.low_freq" color="orange" /></a-col>
      </a-row>
      <a-divider style="margin: 0 0 16px 0" />
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id" />
    </a-card>
  </div>
</template>
      </a-table>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import Breadcrumb from '@/components/Breadcrumb.vue'

const loading = ref(false)
const stats = reactive({ total: 0, high_freq: 0, low_freq: 0 })
const data = ref([])
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })
const columns = [
  { title: '标签名称', dataIndex: 'name', width: 200 },
  { title: '使用次数', dataIndex: 'use_count', width: 120 },
  { title: '覆盖率', dataIndex: 'coverage', width: 120 },
  { title: '最后使用', dataIndex: 'last_used', width: 170 }
]

const loadData = async () => {
  loading.value = true
  try {
    const res = await fetch('/api/v1/members/tag-report', { headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token') } }).then(r => r.json())
    data.value = res.data?.list || []
    pagination.total = data.value.length
    if (res.data?.stats) Object.assign(stats, res.data.stats)
  } catch { data.value = [] } finally { loading.value = false }
}
const onPageChange = (page) => { pagination.current = page; loadData() }
onMounted(() => loadData())
</script>
</a-card>
