<template>
  <div class="container">
    <a-card>
      <template #title><icon-home /> 智能家居设备</template>
      <template #extra><a-button type="primary" @click="handleAdd"><template #icon><icon-plus /></template>添加设备</a-button></template>
      <a-table :columns="columns" :data="devices">
        <template #status="{ record }"><a-tag :color="record.online ? 'green' : 'gray'">{{ record.online ? '在线' : '离线' }}</a-tag></template>
      </a-table>
    </a-card>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const columns = [{ title: '设备名称', dataIndex: 'name' }, { title: '类型', dataIndex: 'type' }, { title: '状态', slotName: 'status' }]
const devices = ref([{ id: 1, name: '智能灯泡', type: 'light', online: true }])
const handleAdd = () => {}
</script>
<style scoped>.container { padding: 16px; }</style>
