<template>
  <div class="container">
    <a-card class="general-card" title="通知日志">
      <template #extra>
        <a-button @click="handleSearch"><icon-refresh />刷新</a-button>
      </template>
      <a-row :gutter="16">
        <a-col :span="6">
          <a-form-item label="渠道类型">
            <a-select v-model="form.channel_type" placeholder="选择渠道" allow-clear>
              <a-option value="email">邮件</a-option>
              <a-option value="sms">短信</a-option>
              <a-option value="webhook">Webhook</a-option>
            </a-select>
          </a-form-item>
        </a-col>
        <a-col :span="6">
          <a-form-item label="发送状态">
            <a-select v-model="form.status" placeholder="选择状态" allow-clear>
              <a-option value="success">成功</a-option>
              <a-option value="failed">失败</a-option>
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
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" row-key="id" />
    </a-table>
  </a-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import Breadcrumb from '@/components/Breadcrumb.vue'

const loading = ref(false)
const data = ref([])
const modalVisible = ref(false)
const modalTitle = ref('新建')
const form = reactive({ channel_type: '', status: '' })
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })
const columns = [
  { title: '时间', dataIndex: 'created_at', width: 170 },
  { title: '渠道', dataIndex: 'channel_type', width: 90 },
  { title: '渠道名称', dataIndex: 'channel_name', width: 140 },
  { title: '接收人', dataIndex: 'recipient', width: 160 },
  { title: '主题', dataIndex: 'subject', ellipsis: true },
  { title: '状态', dataIndex: 'status', width: 90 },
  { title: '重试次数', dataIndex: 'attempt_count', width: 90 }
]

const handleSearch = () => { pagination.current = 1; loadData() }
const handleReset = () => { form.channel_type = ''; form.status = ''; pagination.current = 1; loadData() }

const loadData = async () => {
  loading.value = true
  try {
    const params = { page: pagination.current, page_size: pagination.pageSize }
    if (form.channel_type) params.channel_type = form.channel_type
    if (form.status) params.status = form.status
    const res = await fetch('/api/v1/alerts/notification-logs?' + new URLSearchParams(params), {
      headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token') }
    }).then(r => r.json())
    if (res.code === 0) { data.value = res.data?.list || []; pagination.total = res.data?.total || 0 }
  } catch (e) { Message.error('加载失败') } finally { loading.value = false }
}
onMounted(() => { loadData() })
</script>

