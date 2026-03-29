<template>
  <div class="container">
    <a-card class="general-card" title="AI 质量仪表盘">
      <template #extra>
        <a-button @click="handleSearch"><icon-refresh />刷新</a-button>
      </template>
      <a-row :gutter="16">
        <a-col :span="8">
          <a-form-item label="时间范围">
            <a-range-picker v-model="form.time_range" style="width: 100%" />
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
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" row-key="id" />
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import Breadcrumb from '@/components/Breadcrumb.vue'

const loading = ref(false)
const data = ref([])
const form = reactive({ time_range: [] })
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })
const columns = [
  { title: '推理次数', dataIndex: 'total_inferences', width: 120 },
  { title: '平均延迟(ms)', dataIndex: 'avg_latency_ms', width: 120 },
  { title: '错误率', dataIndex: 'error_rate', width: 100 },
  { title: '置信度', dataIndex: 'avg_confidence', width: 100 }
]

const handleSearch = () => { loadData() }
const handleReset = () => { form.time_range = []; loadData() }

const loadData = async () => {
  loading.value = true
  try {
    const params = {}
    if (form.time_range && form.time_range.length === 2) {
      params.start_time = form.time_range[0].toISOString()
      params.end_time = form.time_range[1].toISOString()
    }
    const res = await fetch('/api/v1/ai/quality/metrics?' + new URLSearchParams(params), {
      headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token') }
    }).then(r => r.json())
    if (res.code === 0) { data.value = res.data?.list || []; pagination.total = res.data?.total || 0 }
  } catch (e) { Message.error('加载失败') } finally { loading.value = false }
}
onMounted(() => { loadData() })
</script>
