<template>
  <div class="page-container">
    <div class="search-form">
      <a-form :model="form" layout="inline">
        <a-form-item label="情绪类型">
          <a-select v-model="form.emotion" placeholder="请选择" allow-clear style="width: 140px">
            <a-option value="calm">平静</a-option>
            <a-option value="happy">开心</a-option>
            <a-option value="anxious">焦虑</a-option>
            <a-option value="angry">愤怒</a-option>
          </a-select>
        </a-form-item>
        <a-form-item>
          <a-button type="primary" @click="handleSearch">搜索</a-button>
          <a-button @click="handleReset">重置</a-button>
        </a-form-item>
      </a-form>
    </div>
    <div class="toolbar">
      <a-button type="primary" @click="handleAnalyze">分析录音</a-button>
    </div>
    <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="onPageChange" />
    <a-modal v-model:visible="modalVisible" title="分析结果" :width="520">
      <a-form :model="form" label-col-flex="100px">
        <a-form-item label="情绪类型">
          <a-input v-model="form.emotion" readonly />
        </a-form-item>
        <a-form-item label="强度">
          <a-input v-model="form.intensity" readonly />
        </a-form-item>
        <a-form-item label="置信度">
          <a-input v-model="form.confidence" readonly />
        </a-form-item>
      </a-form>
      <template #footer>
        <a-button @click="modalVisible = false">关闭</a-button>
      </template>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const loading = ref(false)
const data = ref<any[]>([])
const modalVisible = ref(false)

const form = reactive({
  emotion: '',
  intensity: '',
  confidence: ''
})

const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0
})

const columns = [
  { title: '宠物ID', dataIndex: 'pet_id', width: 100 },
  { title: '情绪类型', dataIndex: 'emotion_type', width: 120 },
  { title: '强度', dataIndex: 'intensity', width: 120 },
  { title: '置信度', dataIndex: 'confidence', width: 120 },
  { title: '转录', dataIndex: 'transcript', ellipsis: true },
  { title: '时间', dataIndex: 'created_at', width: 180 }
]

async function loadData() {
  loading.value = true
  try {
    const params = new URLSearchParams()
    params.append('page', String(pagination.current))
    params.append('page_size', String(pagination.pageSize))
    if (form.emotion) params.append('emotion', form.emotion)

    const res = await fetch(`/api/v1/voice-emotion/records?${params}`)
    const json = await res.json()
    data.value = json.data?.list || []
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
  form.emotion = ''
  pagination.current = 1
  loadData()
}

function handleAnalyze() {
  Message.info('录音分析功能开发中')
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
