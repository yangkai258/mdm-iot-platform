<template>
  <div class="page-container">
    <div class="search-form">
      <a-form :model="form" layout="inline">
        <a-form-item label="家庭ID">
          <a-input v-model="form.family_id" placeholder="请输入" style="width: 160px" />
        </a-form-item>
        <a-form-item>
          <a-button type="primary" @click="handleSearch">搜索</a-button>
          <a-button @click="handleReset">重置</a-button>
        </a-form-item>
      </a-form>
    </div>
    <div class="toolbar">
      <a-button type="primary" @click="handleRefresh">刷新</a-button>
    </div>
    <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="onPageChange" />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

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

    const res = await fetch(`/api/v1/emotions/family-map?${params}`)
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
.page-container {
  background: #fff;
  border-radius: 4px;
  padding: 20px;
}
.search-form {
  margin-bottom: 16px;
  padding: 16px;
  background: #f7f8fa;
  border-radius: 4px;
}
.toolbar {
  margin-bottom: 16px;
}
</style>
