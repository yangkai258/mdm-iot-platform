<template>
  <div class="page-container">
    <div class="search-form">
      <a-form :model="form" layout="inline">
        <a-form-item label="报告周期">
          <a-select v-model="form.period" placeholder="请选择" style="width: 140px">
            <a-option value="daily">日报</a-option>
            <a-option value="weekly">周报</a-option>
            <a-option value="monthly">月报</a-option>
          </a-select>
        </a-form-item>
        <a-form-item>
          <a-button type="primary" @click="handleSearch">搜索</a-button>
          <a-button @click="handleReset">重置</a-button>
        </a-form-item>
      </a-form>
    </div>
    <div class="toolbar">
      <a-button type="primary" @click="handleExport">导出报告</a-button>
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
  period: 'daily'
})

const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0
})

const columns = [
  { title: '日期/周期', dataIndex: 'period', width: 140 },
  { title: '主要情绪', dataIndex: 'primary_emotion', width: 120 },
  { title: '情绪占比', dataIndex: 'emotion_ratio', width: 120 },
  { title: '异常事件', dataIndex: 'anomaly_count', width: 100 },
  { title: '建议', dataIndex: 'suggestion', ellipsis: true }
]

async function loadData() {
  loading.value = true
  try {
    const res = await fetch(`/api/v1/emotion/reports?period=${form.period}`)
    const json = await res.json()
    data.value = json.data || []
    pagination.total = data.value.length
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
  form.period = 'daily'
  pagination.current = 1
  loadData()
}

function onPageChange(page: number) {
  pagination.current = page
  loadData()
}

function handleExport() {
  Message.info('导出功能开发中')
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
