<template>
  <div class="container">
    <a-card>
      <template #title><icon-office-building /> 租户管理</template>
      <template #extra><a-button type="primary" @click="handleCreate"><template #icon><icon-plus /></template>添加租户</a-button></template>
      <a-table :columns="columns" :data="tenants">
        <template #status="{ record }"><a-tag :color="record.status === 'active' ? 'green' : 'gray'">{{ record.status === 'active' ? '正常' : '停用' }}</a-tag></template>
        <template #actions="{ record }"><a-link @click="handleEdit(record)">编辑</a-link></template>
      </a-table>
    </a-card>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const columns = [{ title: '租户名称', dataIndex: 'name' }, { title: '联系人', dataIndex: 'contact' }, { title: '状态', slotName: 'status' }, { title: '操作', slotName: 'actions' }]
const tenants = ref([{ id: 1, name: '租户A', contact: '张三', status: 'active' }])
const handleCreate = () => {}
const handleEdit = (r) => {}
</script>
<style scoped>.container { padding: 16px; }</style>
