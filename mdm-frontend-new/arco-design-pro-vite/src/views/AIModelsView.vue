<template>
  <div class="container">
    <a-card>
      <template #title><icon-robot /> AI模型管理</template>
      <template #extra><a-button type="primary" @click="handleAdd"><template #icon><icon-plus /></template>添加模型</a-button></template>
      <a-table :columns="columns" :data="models">
        <template #status="{ record }"><a-tag :color="record.status === 'active' ? 'green' : 'gray'">{{ record.status === 'active' ? '启用' : '停用' }}</a-tag></template>
        <template #actions="{ record }"><a-link @click="handleEdit(record)">编辑</a-link></template>
      </a-table>
    </a-card>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const columns = [{ title: '模型名称', dataIndex: 'name' }, { title: '版本', dataIndex: 'version' }, { title: '状态', slotName: 'status' }, { title: '操作', slotName: 'actions' }]
const models = ref([{ id: 1, name: 'EmotionNet', version: 'v2.1', status: 'active' }])
const handleAdd = () => {}
const handleEdit = (r) => {}
</script>
<style scoped>.container { padding: 16px; }</style>
