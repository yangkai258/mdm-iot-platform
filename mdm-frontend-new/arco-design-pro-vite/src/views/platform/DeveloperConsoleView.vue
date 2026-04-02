<template>
    <Breadcrumb :items="['Home','Console','']" />


  <div class="container">
    <a-card class="general-card" title="开发者控制台">
      <a-row :gutter="16" style="margin-bottom: 16px">
        <a-col :span="6">
          <a-statistic title="API调用次数" :value="stats.api_calls" />
        </a-col>
        <a-col :span="6">
          <a-statistic title="错误率" :value="stats.error_rate" suffix="%" color="red" />
        </a-col>
        <a-col :span="6">
          <a-statistic title="配额使用" :value="stats.quota_usage" suffix="%" />
        </a-col>
        <a-col :span="6">
          <a-statistic title="响应时间" :value="stats.avg_time" suffix="ms" />
        </a-col>
      </a-row>
      <a-tabs>
        <a-tab-pane key="logs" title="调用日志">
          <a-table :columns="logColumns" :data="logs" :pagination="pagination" row-key="id">
            <template #status="{ record }"><a-badge :color="record.status < 400 ? 'green' : 'red'" :text="record.status" /></template>
          </a-table>
        </a-tab-pane>
        <a-tab-pane key="quota" title="配额管理">
          <a-form :model="quotaForm" layout="vertical" style="max-width: 400px">
            <a-form-item label="日配额"><a-input-number v-model="quotaForm.daily" :min="0" style="width:100%" /></a-form-item>
            <a-form-item label="月配额"><a-input-number v-model="quotaForm.monthly" :min="0" style="width:100%" /></a-form-item>
            <a-button type="primary">保存</a-button>
          </a-form>
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>
      </a-table>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import Breadcrumb from '@/components/Breadcrumb.vue'

const stats = reactive({ api_calls: 0, error_rate: 0, quota_usage: 0, avg_time: 0 })
const quotaForm = reactive({ daily: 10000, monthly: 300000 })
const logs = ref([])
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })
const logColumns = [
  { title: '时间', dataIndex: 'time', width: 170 },
  { title: '接口', dataIndex: 'path', ellipsis: true },
  { title: '方法', dataIndex: 'method', width: 80 },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '耗时', dataIndex: 'duration', width: 100 }
]

onMounted(() => {
  logs.value = [{ id: 1, time: '2026-03-29 10:00:00', path: '/api/v1/devices', method: 'GET', status: 200, duration: '45ms' }]
  pagination.total = 1
})
</script>
