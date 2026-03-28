<template>
  <div class="container">
    <a-card>
      <template #title><icon-bell /> 告警列表</template>
      <a-table :columns="columns" :data="alerts">
        <template #level="{ record }"><a-tag :color="getLevelColor(record.level)">{{ record.level }}</a-tag></template>
        <template #actions="{ record }"><a-link @click="handleProcess(record)">处理</a-link></template>
      </a-table>
    </a-card>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const columns = [{ title: '时间', dataIndex: 'time' }, { title: '告警内容', dataIndex: 'content' }, { title: '级别', slotName: 'level' }, { title: '操作', slotName: 'actions' }]
const alerts = ref([{ time: '2026-03-28 10:00', content: '设备离线', level: 'warning' }])
const getLevelColor = (l) => ({ warning: 'orange', major: 'red', info: 'blue' }[l] || 'gray')
const handleProcess = (r) => {}
</script>
<style scoped>.container { padding: 16px; }</style>
