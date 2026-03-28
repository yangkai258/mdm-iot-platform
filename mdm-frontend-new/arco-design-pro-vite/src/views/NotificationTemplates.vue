<template>
  <div class="container">
    <a-card>
      <template #title><icon-email /> 通知模板</template>
      <template #extra>
        <a-button type="primary" @click="handleCreate"><template #icon><icon-plus /></template>新建模板</a-button>
      </template>
      <a-table :columns="columns" :data="templates">
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
const columns = [{ title: '名称', dataIndex: 'name' }, { title: '类型', dataIndex: 'type' }, { title: '操作', slotName: 'actions' }]
const templates = ref([{ id: 1, name: '系统通知', type: 'system' }])
const handleCreate = () => {}
const handleEdit = (r) => {}
const handleDelete = (r) => {}
</script>
<style scoped>.container { padding: 16px; }</style>
