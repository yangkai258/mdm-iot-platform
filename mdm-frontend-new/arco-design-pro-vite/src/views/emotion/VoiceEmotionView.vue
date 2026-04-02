<template>
    <Breadcrumb :items="['Home','Console','']" />


  <div class="container">
    <a-card class="general-card" title="语音情绪">
      <template #extra>
        <a-space :size="12">
          <a-button type="primary" @click="handleAnalyze"><icon-plus />分析录音</a-button>
          <a-button @click="handleSearch"><icon-refresh />刷新</a-button>
        </a-space>
      </template>
      <a-row :gutter="16">
        <a-col :span="8">
          <a-form-item label="情绪类型">
            <a-select v-model="form.emotion" placeholder="请选择" allow-clear style="width: 100%">
              <a-option value="calm">平静</a-option>
              <a-option value="happy">开心</a-option>
              <a-option value="anxious">焦虑</a-option>
              <a-option value="angry">愤怒</a-option>
            </a-select>
          </a-form-item>
        </a-col>
        <a-col :flex="'86px'" style="display: flex; align-items: flex-end">
          <a-space direction="vertical" :size="8">
            <a-button type="primary" @click="handleSearch">查询</a-button>
            <a-button @click="handleReset">重置</a-button>
          </a-space>
        </a-col>
      </a-row>
      <a-divider style="margin: 0 0 16px 0" />
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id" />
    </a-card>
      </a-table>
    <a-modal v-model:visible="modalVisible" title="分析结果" :width="520">
      <a-form :model="form" label-col-flex="100px">
        <a-form-item label="情绪类型"><a-input v-model="form.emotion" readonly /></a-form-item>
        <a-form-item label="强度"><a-input v-model="form.intensity" readonly /></a-form-item>
        <a-form-item label="置信度"><a-input v-model="form.confidence" readonly /></a-form-item>
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
import Breadcrumb from '@/components/Breadcrumb.vue'

const loading = ref(false)
const data = ref<any[]>([])
const modalVisible = ref(false)
const form = reactive({ emotion: '', intensity: '', confidence: '' })
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })
const columns = [
  { title: '宠物ID', dataIndex: 'pet_id', width: 100 },
  { title: '情绪类型', dataIndex: 'emotion_type', width: 120 },
  { title: '强度', dataIndex: 'intensity', width: 120 },
  { title: '置信度', dataIndex: 'confidence', width: 120 },
  { title: '转录', dataIndex: 'transcript', ellipsis: true },
  { title: '时间', dataIndex: 'created_at', width: 180 }
]

const loadData = async () => {
  loading.value = true
  try {
    const params = new URLSearchParams()
    params.append('page', String(pagination.current))
    params.append('page_size', String(pagination.pageSize))
    if (form.emotion) params.append('emotion', form.emotion)
    const res = await fetch(`/api/voice-emotion/records?${params}`, {
      headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token') }
    }).then(r => r.json())
    data.value = res.data?.list || []
    pagination.total = res.data?.total || 0
  } catch { Message.error('加载失败') } finally { loading.value = false }
}

const handleSearch = () => { pagination.current = 1; loadData() }
const handleReset = () => { form.emotion = ''; pagination.current = 1; loadData() }
const onPageChange = (page: number) => { pagination.current = page; loadData() }
const handleAnalyze = () => { Message.info('录音分析功能开发中') }

onMounted(() => loadData())
</script>
