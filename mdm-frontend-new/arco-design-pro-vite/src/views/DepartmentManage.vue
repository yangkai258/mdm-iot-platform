<template>
  <div class="container">
    <a-card>
      <template #title><icon-business /> 部门管理</template>
      <template #extra><a-button type="primary" @click="handleCreate"><template #icon><icon-plus /></template>新建部门</a-button></template>
      <a-table :columns="columns" :data="departments">
        <template #actions="{ record }"><a-link @click="handleEdit(record)">编辑</a-link></template>
      </a-table>
    </a-card>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const columns = [{ title: '部门名称', dataIndex: 'name' }, { title: '负责人', dataIndex: 'manager' }, { title: '操作', slotName: 'actions' }]
const departments = ref([{ id: 1, name: '研发部', manager: '张三' }])
const handleCreate = () => {}
const handleEdit = (r) => {}
</script>
<style scoped>.container { padding: 16px; }</style>
