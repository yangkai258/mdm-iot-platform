<template>
  <div class="container">
    <a-card>
      <template #title><icon-message /> 反馈管理</template>
      <a-table :columns="columns" :data="feedbacks">
        <template #status="{ record }"><a-tag :color="getStatusColor(record.status)">{{ record.status }}</a-tag></template>
        <template #actions="{ record }"><a-link @click="handleReply(record)">回复</a-link></template>
      </a-table>
    </a-card>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const columns = [{ title: '用户', dataIndex: 'user' }, { title: '内容', dataIndex: 'content' }, { title: '状态', slotName: 'status' }, { title: '操作', slotName: 'actions' }]
const feedbacks = ref([{ id: 1, user: '张三', content: '建议增加xxx功能', status: '待处理' }])
const getStatusColor = (s) => ({ '待处理': 'orange', '已回复': 'green' }[s] || 'gray')
const handleReply = (r) => {}
</script>
<style scoped>.container { padding: 16px; }</style>
