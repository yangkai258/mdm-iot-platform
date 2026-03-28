<template>
  <div class="container">
    <a-card>
      <template #title><icon-database /> 数据集管理</template>
      <template #extra><a-button type="primary" @click="handleCreate"><template #icon><icon-plus /></template>创建数据集</a-button></template>
      <a-table :columns="columns" :data="datasets">
        <template #actions="{ record }"><a-link @click="handleEdit(record)">编辑</a-link></template>
      </a-table>
    </a-card>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const columns = [{ title: '名称', dataIndex: 'name' }, { title: '类型', dataIndex: 'type' }, { title: '大小', dataIndex: 'size' }, { title: '操作', slotName: 'actions' }]
const datasets = ref([{ id: 1, name: '宠物行为数据集', type: '训练数据', size: '256MB' }])
const handleCreate = () => {}
const handleEdit = (r) => {}
</script>
<style scoped>.container { padding: 16px; }</style>
