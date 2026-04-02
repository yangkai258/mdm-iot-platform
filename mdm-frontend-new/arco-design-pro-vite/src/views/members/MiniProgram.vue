<template>
    <Breadcrumb :items="['Home','Console','']" />


  <div class="container">
    <a-card class="general-card" title="小程序管理">
      <template #extra>
        <a-button @click="loadData"><icon-refresh />刷新</a-button>
      </template>
      <a-row :gutter="16" style="margin-bottom: 16px">
        <a-col :span="6"><a-statistic title="小程序总数" :value="stats.total" /></a-col>
        <a-col :span="6"><a-statistic title="在线" :value="stats.online" color="green" /></a-col>
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
const stats = reactive({ total: 0, online: 0 })
const data = ref([])
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })
const columns = [
  { title: '小程序名称', dataIndex: 'name', width: 180 },
  { title: 'AppID', dataIndex: 'app_id', width: 200 },
  { title: '版本', dataIndex: 'version', width: 100 },
  { title: '状态', dataIndex: 'status', width: 90 }
]

const loadData = async () => {
  loading.value = true
  try {
    const res = await fetch('/api/members/mini-programs', { headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token') } }).then(r => r.json())
    data.value = res.data?.list || []
    pagination.total = data.value.length
  } catch { data.value = [] } finally { loading.value = false }
}
const onPageChange = (page) => { pagination.current = page; loadData() }
onMounted(() => loadData())
</script>
</a-card>
