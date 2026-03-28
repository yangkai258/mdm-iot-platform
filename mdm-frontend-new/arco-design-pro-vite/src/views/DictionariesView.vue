<template>
  <div class="container">
    <a-card>
      <template #title><icon-book /> 数据字典</template>
      <template #extra><a-button type="primary" @click="handleCreate"><template #icon><icon-plus /></template>添加字典</a-button></template>
      <a-table :columns="columns" :data="dictionaries">
        <template #actions="{ record }"><a-link @click="handleEdit(record)">编辑</a-link></template>
      </a-table>
    </a-card>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const columns = [{ title: '字典名称', dataIndex: 'name' }, { title: '字典值', dataIndex: 'value' }, { title: '操作', slotName: 'actions' }]
const dictionaries = ref([{ id: 1, name: '状态', value: '启用/禁用' }])
const handleCreate = () => {}
const handleEdit = (r) => {}
</script>
<style scoped>.container { padding: 16px; }</style>
