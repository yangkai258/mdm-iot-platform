<template>
  <div class="container">
    <a-card>
      <template #title><icon-message /> 内容评论</template>
      <a-table :columns="columns" :data="comments">
        <template #status="{ record }">
          <a-tag :color="record.status === '显示' ? 'green' : 'gray'">{{ record.status }}</a-tag>
        </template>
        <template #actions="{ record }">
          <a-link @click="handleApprove(record)">通过</a-link>
          <a-link @click="handleDelete(record)">删除</a-link>
        </template>
      </a-table>
    </a-card>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const columns = [{ title: '内容', dataIndex: 'content' }, { title: '用户', dataIndex: 'user' }, { title: '状态', slotName: 'status' }, { title: '操作', slotName: 'actions' }]
const comments = ref([{ id: 1, content: '很好用', user: '张三', status: '显示' }])
const handleApprove = (r) => {}
const handleDelete = (r) => {}
</script>
<style scoped>.container { padding: 16px; }</style>
