<template>
  <div class="container">
    <a-card class="general-card" title="租户申请">
      <template #extra>
        <a-button @click="loadData"><icon-refresh />刷新</a-button>
      </template>
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id">
        <template #status="{ record }"><a-badge :color="record.status === 'pending' ? 'orange' : record.status === 'approved' ? 'green' : 'red'" :text="record.status === 'pending' ? '待审核' : record.status === 'approved' ? '已通过' : '已拒绝'" /></template>
      </a-table>
    </a-card>
  </div>
</template>
      </a-table>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import Breadcrumb from '@/components/Breadcrumb.vue'

const loading = ref(false)
const data = ref([])
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })
const columns = [
  { title: '公司名称', dataIndex: 'company_name', width: 200 },
  { title: '联系人', dataIndex: 'contact', width: 120 },
  { title: '手机号', dataIndex: 'phone', width: 130 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '申请时间', dataIndex: 'created_at', width: 170 }
]

const loadData = async () => {
  loading.value = true
  try {
    const res = await fetch('/api/v1/tenants/applications', { headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token') } }).then(r => r.json())
    data.value = res.data?.list || []
    pagination.total = data.value.length
  } catch { data.value = [] } finally { loading.value = false }
}
const onPageChange = (page) => { pagination.current = page; loadData() }
onMounted(() => loadData())
</script>
