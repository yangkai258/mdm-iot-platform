<template>
  <div class="container">
    <a-card>
      <template #title><icon-certificate /> 证书管理</template>
      <template #extra><a-button type="primary" @click="handleCreate"><template #icon><icon-plus /></template>创建证书</a-button></template>
      <a-table :columns="columns" :data="certificates">
        <template #status="{ record }"><a-tag :color="record.valid ? 'green' : 'red'">{{ record.valid ? '有效' : '无效' }}</a-tag></template>
      </a-table>
    </a-card>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const columns = [{ title: '证书名称', dataIndex: 'name' }, { title: '类型', dataIndex: 'type' }, { title: '状态', slotName: 'status' }]
const certificates = ref([{ id: 1, name: '设备证书', type: 'device', valid: true }])
const handleCreate = () => {}
</script>
<style scoped>.container { padding: 16px; }</style>
