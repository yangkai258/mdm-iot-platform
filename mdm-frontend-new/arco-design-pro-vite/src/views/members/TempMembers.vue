<template>
    <Breadcrumb :items="['Home','Console','']" />


  <div class="container">
    <a-card class="general-card" title="临时会员">
      <template #extra>
        <a-button @click="loadData"><icon-refresh />刷新</a-button>
      </template>
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id" />
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
  { title: 'ID', dataIndex: 'id', width: 80 },
  { title: '昵称', dataIndex: 'nickname', width: 140 },
  { title: '手机号', dataIndex: 'phone', width: 130 },
  { title: '有效期至', dataIndex: 'expire_at', width: 170 }
]

const loadData = async () => {
  loading.value = true
  try {
    const res = await fetch('/api/v1/members/temp-members', { headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token') } }).then(r => r.json())
    data.value = res.data?.list || []
    pagination.total = data.value.length
  } catch { data.value = [] } finally { loading.value = false }
}
const onPageChange = (page) => { pagination.current = page; loadData() }
onMounted(() => loadData())
</script>
</a-card>
