<template>
  <div class="container">
    <a-card>
      <template #title><icon-list /> 设备列表</template>
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
const handleDetail = (r) => {}
</script>
<style scoped>.container { padding: 16px; }</style>
