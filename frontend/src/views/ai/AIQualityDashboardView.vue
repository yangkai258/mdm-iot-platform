<template>
  <div class="page-container">
    <div class="search-form">
      <a-form :model="form" layout="inline">
        <a-form-item label="时间范围"><a-range-picker v-model="form.time_range" style="width: 280px" /></a-form-item>
        <a-form-item>
          <a-button type="primary" @click="handleSearch">搜索</a-button>
          <a-button @click="handleReset">重置</a-button>
        </a-form-item>
      </a-form>
    </div>
    <div class="toolbar">
      <a-button type="primary" @click="handleCreate">新建</a-button>
    </div>
    <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" />
    <a-modal v-model:visible="modalVisible" :title="modalTitle">
      <a-form :model="form" label-col-flex="100px">
        <a-form-item label="名称"><a-input v-model="form.name" /></a-form-item>
      </a-form>
      <template #footer>
        <a-button @click="modalVisible = false">取消</a-button>
        <a-button type="primary" @click="handleSubmit">确定</a-button>
      </template>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Message } from '@arco-design/web-vue'
import { useAIQuality } from '@/composables/useAIQuality'

const router = useRouter()

const loading = ref(false)
const data = ref([])
const modalVisible = ref(false)
const modalTitle = ref('新建')

const form = reactive({
  time_range: []
})

const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0
})

const columns = [
  { title: '推理次数', dataIndex: 'total_inferences', width: 120 },
  { title: '平均延迟', dataIndex: 'avg_latency_ms', width: 120 },
  { title: '错误率', dataIndex: 'error_rate', width: 100 },
  { title: '置信度', dataIndex: 'avg_confidence', width: 100 }
]

const handleSearch = () => {
  loadData()
}

const handleReset = () => {
  form.time_range = []
  loadData()
}

const handleCreate = () => {
  modalTitle.value = '新建'
  modalVisible.value = true
}

const handleSubmit = () => {
  modalVisible.value = false
  Message.success('保存成功')
}

const { metricsLoading, metrics, logLoading, logList, anomalyAlerts, loadMetrics, loadLogs } = useAIQuality()

const loadData = async () => {
  loading.value = true
  try {
    const params = {}
    if (form.time_range && form.time_range.length === 2) {
      params.start_time = form.time_range[0].toISOString()
      params.end_time = form.time_range[1].toISOString()
    }
    await Promise.all([
      loadMetrics(params),
      loadLogs({ ...params, page_size: 500 })
    ])
    data.value = [metrics.value]
    pagination.total = 1
  } catch (e) {
    Message.error('加载失败')
  } finally {
    loading.value = false
  }
}

const goToDetail = (id) => {
  router.push(`/ai/behavior-detail/${id}`)
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.page-container { background: #fff; border-radius: 4px; padding: 20px; }
.search-form { margin-bottom: 16px; padding: 16px; background: #f7f8fa; border-radius: 4px; }
.toolbar { margin-bottom: 16px; }
</style>
