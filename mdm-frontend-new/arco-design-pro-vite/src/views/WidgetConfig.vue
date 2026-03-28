<template>
  <div class="container">
    <a-card>
      <template #title><icon-app /> 小部件配置</template>
      <template #extra>
        <a-button type="primary" @click="handleCreate"><template #icon><icon-plus /></template>添加小部件</a-button>
      </template>
      <a-table :columns="columns" :data="widgets">
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
const widgets = ref([{ id: 1, name: '设备状态', type: 'chart' }])
const handleCreate = () => {}
const handleEdit = (r) => {}
const handleDelete = (r) => {}
</script>
<style scoped>.container { padding: 16px; }</style>
