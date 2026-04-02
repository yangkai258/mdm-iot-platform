<template>
    <Breadcrumb :items="['Home','Console','']" />


  <div class="page-container">
    <a-card class="general-card" title="情绪日志">
      <template #extra>
        <a-button @click="handleExport"><icon-download />导出</a-button>
      </template>
      <div class="search-form">
        <a-form :model="form" layout="inline">
          <a-form-item label="情绪类型">
            <a-select v-model="form.emotion" placeholder="请选择" allow-clear style="width: 140px">
              <a-option value="happy">开心</a-option>
              <a-option value="sad">难过</a-option>
              <a-option value="angry">生气</a-option>
              <a-option value="fear">害怕</a-option>
              <a-option value="neutral">平静</a-option>
            </a-select>
          </a-form-item>
          <a-form-item>
            <a-button type="primary" @click="handleSearch">查询</a-button>
            <a-button @click="handleReset">重置</a-button>
          </a-form-item>
        </a-form>
      </div>
      <a-table
      :columns="columns"
      :data="data"
      :loading="loading"
      :pagination="pagination"
      @page-change="onPageChange"
    />
    </a-table>
  </a-card>`n</div></template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import { IconDownload } from '@arco-design/web-vue/es/icon'
import { getEmotionLogs } from '@/api/emotion'

const loading = ref(false)
const data = ref<any[]>([])

const form = reactive({
  emotion: ''
})

const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0
})

const columns = [
  { title: '时间', dataIndex: 'created_at', width: 180 },
  { title: '设备', dataIndex: 'device_id', width: 120 },
  { title: '情绪类型', dataIndex: 'emotion_label', width: 120 },
  { title: '强度', dataIndex: 'intensity', width: 150 },
  { title: '触发原因', dataIndex: 'trigger', ellipsis: true }
]

async function loadData() {
  try {
    loading.value = true
    const params: any = { page: pagination.current, page_size: pagination.pageSize }
    if (form.emotion) params.emotion = form.emotion
    const res = await getEmotionLogs(params)
    data.value = res.data?.list || res.data || []
    pagination.total = res.data?.total || 0
  } catch (err: any) {
    Message.error('加载失败: ' + err.message)
  } finally {
    loading.value = false
  }
}

function handleSearch() {
  pagination.current = 1
  loadData()
}

function handleReset() {
  form.emotion = ''
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
.page-container { padding: 16px; }
.search-form { margin-bottom: 16px; padding: 16px; background: var(--color-fill-lightest); border-radius: 4px; }
</style>


