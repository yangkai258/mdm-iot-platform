<template>
    <Breadcrumb :items="['Home','Console','']" />


  <div class="page-container">
    <a-card class="general-card" title="家庭情绪地图">
      <template #extra>
        <a-button @click="handleRefresh"><icon-refresh />刷新</a-button>
      </template>
      <div class="search-form">
        <a-form :model="form" layout="inline">
          <a-form-item label="家庭ID">
            <a-input v-model="form.family_id" placeholder="请输入" style="width: 160px" />
          </a-form-item>
          <a-form-item>
            <a-button type="primary" @click="handleSearch">查询</a-button>
            <a-button @click="handleReset">重置</a-button>
          </a-form-item>
        </a-form>
      </div>
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="onPageChange" />
    </a-table>
  </a-card>`n</div></template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import { IconRefresh } from '@arco-design/web-vue/es/icon'

const loading = ref(false)
const data = ref<any[]>([])

const form = reactive({
  family_id: ''
})

const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0
})

const columns = [
  { title: '宠物ID', dataIndex: 'pet_id', width: 100 },
  { title: '当前情绪', dataIndex: 'current_mood', width: 120 },
  { title: '情绪强度', dataIndex: 'intensity', width: 120 },
  { title: '互动次数', dataIndex: 'interaction_count', width: 120 },
  { title: '更新时间', dataIndex: 'updated_at', width: 180 }
]

async function loadData() {
  loading.value = true
  try {
    const params = new URLSearchParams()
    if (form.family_id) params.append('family_id', form.family_id)
    params.append('page', String(pagination.current))
    params.append('page_size', String(pagination.pageSize))

    const res = await fetch(`/api/emotions/family-map?${params}`)
    const json = await res.json()
    data.value = json.data?.members || []
    pagination.total = json.data?.total || 0
  } catch {
    Message.error('加载失败')
  } finally {
    loading.value = false
  }
}

function handleSearch() {
  pagination.current = 1
  loadData()
}

function handleReset() {
  form.family_id = ''
  pagination.current = 1
  loadData()
}

function handleRefresh() {
  loadData()
}

function onPageChange(page: number) {
  pagination.current = page
  loadData()
}

onMounted(() => loadData())
</script>

<style scoped>
.page-container { padding: 16px; }
.search-form { margin-bottom: 16px; padding: 16px; background: var(--color-fill-lightest); border-radius: 4px; }
</style>


