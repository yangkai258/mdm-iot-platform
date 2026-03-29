<template>
  <div class="container">
    <a-card class="general-card" title="工作台">
      <a-row :gutter="16" style="margin-bottom: 16px">
        <a-col :span="6"><a-statistic title="待办事项" :value="stats.todo" color="blue" /></a-col>
        <a-col :span="6"><a-statistic title="进行中" :value="stats.in_progress" color="orange" /></a-col>
        <a-col :span="6"><a-statistic title="已完成" :value="stats.done" color="green" /></a-col>
        <a-col :span="6"><a-statistic title="总任务" :value="stats.total" /></a-col>
      </a-row>
      <a-divider style="margin: 0 0 16px 0" />
      <a-table :columns="columns" :data="tasks" :pagination="pagination" row-key="id">
        <template #priority="{ record }"><a-tag :color="record.priority === 'high' ? 'red' : record.priority === 'medium' ? 'orange' : 'blue'">{{ record.priority }}</a-tag></template>
      </a-table>
    </a-card>
  </div>
</template>
      </a-table>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import Breadcrumb from '@/components/Breadcrumb.vue'

const stats = reactive({ todo: 0, in_progress: 0, done: 0, total: 0 })
const tasks = ref([])
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })
const columns = [
  { title: '任务名称', dataIndex: 'title', width: 240 },
  { title: '优先级', slotName: 'priority', width: 100 },
  { title: '状态', dataIndex: 'status', width: 100 },
  { title: '截止时间', dataIndex: 'deadline', width: 170 }
]

onMounted(() => {
  tasks.value = []
  pagination.total = 0
})
</script>
