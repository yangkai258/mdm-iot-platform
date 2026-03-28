<template>
  <div class="container">
    <a-card>
      <template #title><icon-app /> 动作库</template>
      <template #extra><a-button type="primary" @click="handleCreate"><template #icon><icon-plus /></template>添加动作</a-button></template>
      <a-table :columns="columns" :data="actions">
        <template #actions="{ record }"><a-link @click="handleEdit(record)">编辑</a-link></template>
      </a-table>
    </a-card>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const columns = [{ title: '动作名称', dataIndex: 'name' }, { title: '类型', dataIndex: 'type' }, { title: '操作', slotName: 'actions' }]
const actions = ref([{ id: 1, name: '挥手', type: 'greeting' }])
const handleCreate = () => {}
const handleEdit = (r) => {}
</script>
<style scoped>.container { padding: 16px; }</style>
