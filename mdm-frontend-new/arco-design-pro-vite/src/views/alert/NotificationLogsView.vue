<template>
  <Breadcrumb :items="['Home','Alert','NotificationLogs','']" />
  <div class="page-container">
    <a-card class="general-card" title="通知日志">
      <template #extra><a-button @click="loadData"><icon-refresh />刷新</a-button></template>
      <div class="search-form">
        <a-form :model="form" layout="inline">
          <a-form-item label="通知类型"><a-select v-model="form.channel" placeholder="选择类型" allow-clear style="width: 120px"><a-option value="email">邮件</a-option><a-option value="sms">短信</a-option><a-option value="webhook">Webhook</a-option></a-select></a-form-item>
          <a-form-item><a-button type="primary" @click="loadData">查询</a-button><a-button @click="handleReset">重置</a-button></a-form-item>
        </a-form>
      </div>
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" />
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import { IconRefresh } from '@arco-design/web-vue/es/icon'

const loading = ref(false)
const data = ref<any[]>([])
const form = ref<any>({ channel: '' })

const columns = [
  { title: 'ID', dataIndex: 'id', width: 70 },
  { title: '告警ID', dataIndex: 'alert_id', width: 90 },
  { title: '通知类型', dataIndex: 'channel', width: 100 },
  { title: '接收人', dataIndex: 'recipient', width: 160 },
  { title: '内容', dataIndex: 'message', ellipsis: true },
  { title: '状态', dataIndex: 'status', width: 90 },
  { title: '发送时间', dataIndex: 'sent_at', width: 170 }
]

const pagination = ref({ current: 1, pageSize: 20, total: 0, showTotal: true })

async function loadData() { loading.value = true; data.value = []; loading.value = false }
function handleReset() { form.value = { channel: '' }; loadData() }
onMounted(() => { loadData() })
</script>

<style scoped>
.page-container { padding: 16px; }
.search-form { margin-bottom: 16px; padding: 16px; background: var(--color-fill-lightest); border-radius: 4px; }
</style>
