<template>
  <div class="container">
    <a-card>
      <template #title><icon-team /> 角色管理</template>
      <template #extra><a-button type="primary" @click="handleCreate"><template #icon><icon-plus /></template>添加角色</a-button></template>
      <a-table :columns="columns" :data="roles">
        <template #actions="{ record }"><a-link @click="handleEdit(record)">编辑</a-link></template>
      </a-table>
    </a-card>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const columns = [{ title: '角色名称', dataIndex: 'name' }, { title: '描述', dataIndex: 'desc' }, { title: '操作', slotName: 'actions' }]
const roles = ref([{ id: 1, name: '管理员', desc: '系统管理员' }])
const handleCreate = () => {}
const handleEdit = (r) => {}
</script>
<style scoped>.container { padding: 16px; }</style>
