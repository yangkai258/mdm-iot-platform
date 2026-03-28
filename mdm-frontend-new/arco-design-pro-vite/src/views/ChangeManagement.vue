<template>
  <div class="container">
    <a-card>
      <template #title><icon-edit /> 变更管理</template>
      <a-table :columns="columns" :data="changes">
        <template #status="{ record }"><a-tag :color="getStatusColor(record.status)">{{ record.status }}</a-tag></template>
      </a-table>
    </a-card>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const columns = [{ title: '变更ID', dataIndex: 'id' }, { title: '变更内容', dataIndex: 'content' }, { title: '状态', slotName: 'status' }]
const changes = ref([{ id: 1, content: '系统升级', status: '进行中' }])
const getStatusColor = (s) => ({ '进行中': 'blue', '已完成': 'green', '已撤销': 'gray' }[s] || 'gray')
</script>
<style scoped>.container { padding: 16px; }</style>
