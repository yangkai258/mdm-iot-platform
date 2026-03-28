<template>
  <div class="container">
    <a-card>
      <template #title><icon-file-text /> API日志</template>
      <a-table :columns="columns" :data="logs">
        <template #method="{ record }"><a-tag :color="getMethodColor(record.method)">{{ record.method }}</a-tag></template>
      </a-table>
    </a-card>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const columns = [{ title: '时间', dataIndex: 'time' }, { title: 'Method', slotName: 'method' }, { title: '路径', dataIndex: 'path' }, { title: '状态码', dataIndex: 'status' }]
const logs = ref([{ time: '2026-03-28 10:00', method: 'GET', path: '/api/v1/devices', status: 200 }])
const getMethodColor = (m) => ({ GET: 'blue', POST: 'green', PUT: 'orange', DELETE: 'red' }[m] || 'gray')
</script>
<style scoped>.container { padding: 16px; }</style>
