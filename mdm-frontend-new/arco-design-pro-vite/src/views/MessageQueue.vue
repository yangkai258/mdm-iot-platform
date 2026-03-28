<template>
  <div class="container">
    <a-card>
      <template #title><icon-list /> 消息队列</template>
      <a-row :gutter="16" style="margin-bottom: 16px">
        <a-col :span="6"><a-card><a-statistic title="队列数量" :value="stats.queues" /></a-card></a-col>
        <a-col :span="6"><a-card><a-statistic title="消息总数" :value="stats.messages" /></a-card></a-col>
      </a-row>
      <a-table :columns="columns" :data="queues" />
    </a-card>
  </div>
</template>
<script setup>
import { ref, reactive } from 'vue'
const stats = reactive({ queues: 10, messages: 5000 })
const columns = [{ title: '队列名', dataIndex: 'name' }, { title: '消息数', dataIndex: 'count' }, { title: '状态', dataIndex: 'status' }]
const queues = ref([{ name: 'device.up', count: 120, status: 'active' }])
</script>
<style scoped>.container { padding: 16px; }</style>
