<template>
  <div class="container">
    <Breadcrumb :items="['menu.members', 'menu.members.storeLocations']" />
    <a-card class="general-card" title="门店位置">
      <template #extra>
        <a-button @click="loadData"><icon-refresh />刷新</a-button>
      </template>
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id" />
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import Breadcrumb from '@/components/Breadcrumb.vue'

const loading = ref(false)
const data = ref([])
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })
const columns = [
  { title: '门店名称', dataIndex: 'name', width: 180 },
  { title: '地址', dataIndex: 'address', ellipsis: true },
  { title: '状态', dataIndex: 'status', width: 90 }
]

const loadData = async () => {
  loading.value = true
  try {
    const res = await fetch('/api/v1/members/store-locations', { headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token') } }).then(r => r.json())
    data.value = res.data?.list || []
    pagination.total = data.value.length
  } catch { data.value = [] } finally { loading.value = false }
}
const onPageChange = (page) => { pagination.current = page; loadData() }
onMounted(() => loadData())
</script>
