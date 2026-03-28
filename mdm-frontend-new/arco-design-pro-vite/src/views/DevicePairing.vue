<template>
  <div class="container">
    <a-card>
      <template #title><icon-link /> 设备配对</template>
      <a-button type="primary" @click="handleStartPairing">开始配对</a-button>
      <a-divider />
      <a-card title="配对记录">
        <a-table :columns="columns" :data="records" size="small" />
      </a-card>
    </a-card>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const columns = [{ title: '设备ID', dataIndex: 'deviceId' }, { title: '配对时间', dataIndex: 'time' }, { title: '状态', dataIndex: 'status' }]
const records = ref([{ deviceId: 'D001', time: '2026-03-28 10:00', status: '成功' }])
const handleStartPairing = () => {}
</script>
<style scoped>.container { padding: 16px; }</style>
