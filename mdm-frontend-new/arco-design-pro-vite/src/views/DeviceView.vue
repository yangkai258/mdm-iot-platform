<template>
  <div class="container">
    <a-card>
      <template #title><icon-computer /> 设备管理</template>
      <template #extra><a-button type="primary" @click="handleAdd"><template #icon><icon-plus /></template>添加设备</a-button></template>
      <a-table :columns="columns" :data="devices">
        <template #status="{ record }"><a-tag :color="record.online ? 'green' : 'gray'">{{ record.online ? '在线' : '离线' }}</a-tag></template>
        <template #actions="{ record }"><a-link @click="handleDetail(record)">详情</a-link></template>
      </a-table>
    </a-card>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const columns = [{ title: '设备ID', dataIndex: 'deviceId' }, { title: '名称', dataIndex: 'name' }, { title: '状态', slotName: 'status' }, { title: '操作', slotName: 'actions' }]
const devices = ref([{ deviceId: 'D001', name: 'M5Stack-001', online: true }])
const handleAdd = () => {}
const handleDetail = (r) => {}
</script>
<style scoped>.container { padding: 16px; }</style>
