<template>
  <div class="page-container">
    <a-card class="general-card" title="Webhook日志">
      <div class="search-form">
        <a-form :model="form" layout="inline">
          <a-form-item label="Webhook">
            <a-select v-model="form.webhook_id" placeholder="请选择" allow-clear style="width: 200px" @change="loadData">
              <a-option v-for="wh in webhooks" :key="wh.id" :value="wh.id">{{ wh.url }}</a-option>
            </a-select>
          </a-form-item>
          <a-form-item label="状态">
            <a-select v-model="form.status" placeholder="请选择" allow-clear style="width: 120px" @change="loadData">
              <a-option value="success">成功</a-option>
              <a-option value="failed">失败</a-option>
            </a-select>
          </a-form-item>
          <a-form-item>
            <a-button type="primary" @click="loadData">查询</a-button>
            <a-button @click="handleReset">重置</a-button>
          </a-form-item>
        </a-form>
      </div>
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id">
      <template #status="{ record }">
        <a-tag :color="record.status === 'success' ? 'green' : 'red'">
          {{ record.status === 'success' ? '成功' : '失败' }}
        </a-tag>
      </template>
      </a-table>
      <template #response_time="{ record }">
        {{ record.response_time }}ms
      </template>
      <template #actions="{ record }">
        <a-button type="text" size="small" @click="handleViewDetail(record)">详情</a-button>
      </template>
    

    <!-- 详情弹窗 -->
    <a-modal v-model:visible="detailVisible" title="调用详情" :width="700">
      <a-form :model="currentRecord" layout="vertical" label-col-flex="100px">
        <a-form-item label="Webhook">
          {{ currentRecord?.webhook_url }}
        </a-form-item>
        <a-form-item label="事件类型">
          <a-tag>{{ currentRecord?.event_type }}</a-tag>
        </a-form-item>
        <a-form-item label="状态">
          <a-tag :color="currentRecord?.status === 'success' ? 'green' : 'red'">
            {{ currentRecord?.status === 'success' ? '成功' : '失败' }}
          </a-tag>
          ({{ currentRecord?.response_time }}ms)
        </a-form-item>
        <a-form-item label="触发时间">{{ currentRecord?.triggered_at }}</a-form-item>
        <a-form-item label="请求体">
          <a-textarea :model-value="JSON.stringify(JSON.parse(currentRecord?.request_body || '{}'), null, 2)" :rows="6" readonly />
        </a-form-item>
        <a-form-item label="响应">
          <a-textarea :model-value="currentRecord?.response_body || '-'" :rows="4" readonly />
        </a-form-item>
      </a-form>
    </a-modal>
    </a-card>`n</div></template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const props = defineProps({ webhookId: String })
const loading = ref(false)
const data = ref([])
const webhooks = ref([])
const detailVisible = ref(false)
const currentRecord = ref(null)
const form = reactive({ webhook_id: undefined, status: undefined })
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })

const columns = [
  { title: 'ID', dataIndex: 'id', width: 80 },
  { title: 'Webhook', dataIndex: 'webhook_url', ellipsis: true, width: 200 },
  { title: '事件类型', dataIndex: 'event_type', width: 140 },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '响应时间', slotName: 'response_time', width: 100 },
  { title: '触发时间', dataIndex: 'triggered_at', width: 160 },
  { title: '操作', slotName: 'actions', width: 80, fixed: 'right' },
]

const loadData = async () => {
  loading.value = true
  try {
    const res = await fetch(`/api/v1/webhooks/logs?page=${pagination.current}&page_size=${pagination.pageSize}&webhook_id=${form.webhook_id || ''}&status=${form.status || ''}`)
    const json = await res.json()
    data.value = json.data?.list || json.data || []
    pagination.total = json.data?.total || 0
  } catch {
    data.value = []
  } finally {
    loading.value = false
  }
}

const loadWebhooks = async () => {
  try {
    const res = await fetch('/api/v1/webhooks?page_size=100')
    const json = await res.json()
    webhooks.value = json.data?.list || json.data || []
    if (props.webhookId) {
      form.webhook_id = props.webhookId
    }
  } catch {}
}

const handleViewDetail = (record) => {
  currentRecord.value = record
  detailVisible.value = true
}

const handleReset = () => { form.webhook_id = undefined; form.status = undefined; loadData() }
const onPageChange = (page) => { pagination.current = page; loadData() }
onMounted(() => { loadWebhooks().then(() => loadData()) })
</script>

<style scoped>
.page-container { padding: 16px; }
.search-form { margin-bottom: 16px; padding: 16px; background: var(--color-fill-lightest); border-radius: 4px; }
</style>

