<template>
  <div class="container">
    <a-card>
      <template #title><icon-safe /> 安全审计</template>
      <a-table :columns="columns" :data="logs">
        <template #level="{ record }"><a-tag :color="getColor(record.level)">{{ record.level }}</a-tag></template>
      </a-table>
    </a-card>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const columns = [{ title: '时间', dataIndex: 'time' }, { title: '用户', dataIndex: 'user' }, { title: '操作', dataIndex: 'action' }, { title: '级别', slotName: 'level' }]
const logs = ref([{ time: '2026-03-28', user: 'admin', action: '登录系统', level: 'info' }])
const getColor = (l) => ({ info: 'blue', warning: 'orange', error: 'red' }[l] || 'gray')
</script>
<style scoped>.container { padding: 16px; }</style>
