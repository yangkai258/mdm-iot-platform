<template>
  <div class="container">
    <a-card>
      <template #title><icon-bug /> 错误追踪</template>
      <a-table :columns="columns" :data="errors">
        <template #level="{ record }"><a-tag :color="getColor(record.level)">{{ record.level }}</a-tag></template>
      </a-table>
    </a-card>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const columns = [{ title: '时间', dataIndex: 'time' }, { title: '错误信息', dataIndex: 'message' }, { title: '级别', slotName: 'level' }]
const errors = ref([{ time: '2026-03-28', message: 'NullPointerException', level: 'error' }])
const getColor = (l) => ({ error: 'red', warning: 'orange', info: 'blue' }[l] || 'gray')
</script>
<style scoped>.container { padding: 16px; }</style>
