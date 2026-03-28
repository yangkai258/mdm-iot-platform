<template>
  <div class="container">
    <a-card>
      <template #title>
        <a-space><icon-shield /> 合规策略</a-space>
      </template>
      <template #extra>
        <a-button type="primary" @click="handleCreate">
          <template #icon><icon-plus /></template>
          新建策略
        </a-button>
      </template>

      <a-table :columns="columns" :data="policies">
        <template #enabled="{ record }">
          <a-switch v-model="record.enabled" />
        </template>
        <template #actions="{ record }">
          <a-link @click="handleEdit(record)">编辑</a-link>
          <a-link @click="handleDelete(record)">删除</a-link>
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const columns = [
  { title: '策略名称', dataIndex: 'name' },
  { title: '类型', dataIndex: 'type' },
  { title: '启用', slotName: 'enabled', width: 100 },
  { title: '操作', slotName: 'actions', width: 150 }
]

const policies = ref([
  { id: 1, name: 'GDPR合规', type: '数据保护', enabled: true }
])

const handleCreate = () => { }
const handleEdit = (r) => { }
const handleDelete = (r) => { }
</script>

<style scoped>
.container { padding: 16px; }
</style>
