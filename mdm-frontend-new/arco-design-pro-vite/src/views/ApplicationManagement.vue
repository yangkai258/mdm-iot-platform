<template>
  <div class="container">
    <a-card>
      <template #title><icon-app /> 应用管理</template>
      <template #extra><a-button type="primary" @click="handleCreate"><template #icon><icon-plus /></template>创建应用</a-button></template>
      <a-table :columns="columns" :data="apps">
        <template #status="{ record }"><a-tag :color="record.status === 'active' ? 'green' : 'gray'">{{ record.status }}</a-tag></template>
        <template #actions="{ record }"><a-link @click="handleEdit(record)">编辑</a-link></template>
      </a-table>
    </a-card>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const columns = [{ title: '名称', dataIndex: 'name' }, { title: 'AppKey', dataIndex: 'appKey' }, { title: '状态', slotName: 'status' }, { title: '操作', slotName: 'actions' }]
const apps = ref([{ id: 1, name: '我的应用', appKey: 'abc123', status: 'active' }])
const handleCreate = () => {}
const handleEdit = (r) => {}
</script>
<style scoped>.container { padding: 16px; }</style>
