<template>
  <Breadcrumb :items="['Home','Alert','NotificationStats','']" />
  <div class="page-container">
    <a-card class="general-card" title="通知统计">
      <a-row :gutter="16" style="margin-bottom: 16px">
        <a-col :span="6"><a-card><a-statistic title="今日发送" :value="stats.today" /></a-card></a-col>
        <a-col :span="6"><a-card><a-statistic title="本周发送" :value="stats.week" /></a-card></a-col>
        <a-col :span="6"><a-card><a-statistic title="本月发送" :value="stats.month" /></a-card></a-col>
        <a-col :span="6"><a-card><a-statistic title="失败数" :value="stats.failed" :value-style="{ color: '#F56C6C' }" /></a-card></a-col>
      </a-row>
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" />
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const loading = ref(false)
const data = ref<any[]>([])
const stats = reactive({ today: 0, week: 0, month: 0, failed: 0 })

const columns = [
  { title: '日期', dataIndex: 'date', width: 120 },
  { title: '邮件', dataIndex: 'email', width: 80 },
  { title: '短信', dataIndex: 'sms', width: 80 },
  { title: 'Webhook', dataIndex: 'webhook', width: 80 },
  { title: '成功数', dataIndex: 'success', width: 90 },
  { title: '失败数', dataIndex: 'failed', width: 90 }
]

const pagination = ref({ current: 1, pageSize: 20, total: 0, showTotal: true })

async function loadData() { loading.value = true; data.value = []; loading.value = false }
onMounted(() => { loadData() })
</script>

<style scoped>
.page-container { padding: 16px; }
</style>
