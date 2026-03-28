<template>
  <div class="container">
    <a-card>
      <template #title><icon-history /> 订阅历史</template>
      <a-table :columns="columns" :data="history">
        <template #status="{ record }"><a-tag :color="record.status === 'active' ? 'green' : 'gray'">{{ record.status === 'active' ? '生效中' : '已过期' }}</a-tag></template>
      </a-table>
    </a-card>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const columns = [{ title: '订阅类型', dataIndex: 'type' }, { title: '开始时间', dataIndex: 'startTime' }, { title: '结束时间', dataIndex: 'endTime' }, { title: '状态', slotName: 'status' }]
const history = ref([{ id: 1, type: '专业版', startTime: '2026-03-01', endTime: '2026-04-01', status: 'active' }])
</script>
<style scoped>.container { padding: 16px; }</style>
