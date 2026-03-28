<template>
  <div class="container">
    <a-card>
      <template #title><icon-user /> 用户管理</template>
      <template #extra><a-button type="primary" @click="handleCreate"><template #icon><icon-plus /></template>添加用户</a-button></template>
      <a-table :columns="columns" :data="users">
        <template #status="{ record }"><a-tag :color="record.status === 'active' ? 'green' : 'gray'">{{ record.status === 'active' ? '启用' : '禁用' }}</a-tag></template>
        <template #actions="{ record }"><a-link @click="handleEdit(record)">编辑</a-link></template>
      </a-table>
    </a-card>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const columns = [{ title: '用户名', dataIndex: 'username' }, { title: '邮箱', dataIndex: 'email' }, { title: '角色', dataIndex: 'role' }, { title: '状态', slotName: 'status' }, { title: '操作', slotName: 'actions' }]
const users = ref([{ id: 1, username: 'admin', email: 'admin@example.com', role: '管理员', status: 'active' }])
const handleCreate = () => {}
const handleEdit = (r) => {}
</script>
<style scoped>.container { padding: 16px; }</style>
