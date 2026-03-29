<template>
  <div class="page-container">
    <a-card class="general-card" title="决策日志">
      <template #extra>
        <a-button @click="handleRefresh"><icon-refresh />刷新</a-button>
      </template>
      <div class="search-form">
        <a-form :model="form" layout="inline">
          <a-form-item label="决策类型">
            <a-select v-model="form.decision_type" placeholder="选择类型" allow-clear style="width: 140px">
              <a-option value="navigation">导航决策</a-option>
              <a-option value="interaction">交互决策</a-option>
              <a-option value="safety">安全决策</a-option>
              <a-option value="task">任务决策</a-option>
            </a-select>
          </a-form-item>
          <a-form-item>
            <a-button type="primary" @click="handleSearch">查询</a-button>
            <a-button @click="handleReset">重置</a-button>
          </a-form-item>
        </a-form>
      </div>
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { getDecisionLogs, getDecisionContext } from '@/api/embodied'
import { Message } from '@arco-design/web-vue'
import { IconRefresh } from '@arco-design/web-vue/es/icon'

const route = useRoute()
const deviceId = ref(route.params.device_id as string || '')

const loading = ref(false)
const data = ref<any[]>([])
const form = ref<any>({ decision_type: '' })

const columns = [
  { title: 'ID', dataIndex: 'id', width: 70 },
  { title: '决策类型', dataIndex: 'decision_type', width: 120 },
  { title: '执行动作', dataIndex: 'chosen_action', width: 140 },
  { title: '置信度', dataIndex: 'confidence', width: 100 },
  { title: '延迟', dataIndex: 'latency_ms', width: 90 },
  { title: '决策时间', dataIndex: 'decided_at', width: 170 }
]

const pagination = ref({ current: 1, pageSize: 20, total: 0, showTotal: true, showPageSize: true })

async function loadData() {
  try {
    loading.value = true
    const params: any = { page: pagination.value.current, page_size: pagination.value.pageSize }
    if (form.value.decision_type) params.decision_type = form.value.decision_type
    const res = await getDecisionLogs(deviceId.value, params)
    data.value = res.data?.logs || res.data || []
    pagination.value.total = data.value.length
  } catch (err: any) {
    Message.error('加载失败: ' + err.message)
  } finally {
    loading.value = false
  }
}

function handleSearch() {
  pagination.value.current = 1
  loadData()
}

function handleReset() {
  form.value = { decision_type: '' }
  handleSearch()
}

function handleRefresh() {
  loadData()
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.page-container { padding: 16px; }
.search-form { margin-bottom: 16px; padding: 16px; background: var(--color-fill-lightest); border-radius: 4px; }
</style>
