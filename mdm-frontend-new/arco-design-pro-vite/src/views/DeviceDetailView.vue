<template>
  <div class="container">
    <a-card>
      <template #title><icon-computer /> 设备详情</template>
      <a-descriptions :column="2" bordered>
        <a-descriptions-item label="设备ID">{{ device.id }}</a-descriptions-item>
        <a-descriptions-item label="设备名称">{{ device.name }}</a-descriptions-item>
        <a-descriptions-item label="型号">{{ device.model }}</a-descriptions-item>
        <a-descriptions-item label="固件版本">{{ device.firmware }}</a-descriptions-item>
        <a-descriptions-item label="状态"><a-tag :color="device.online ? 'green' : 'gray'">{{ device.online ? '在线' : '离线' }}</a-tag></a-descriptions-item>
        <a-descriptions-item label="最后活跃">{{ device.lastActive }}</a-descriptions-item>
      </a-descriptions>
    </a-card>
  </div>
</template>
<script setup>
import { ref, reactive } from 'vue'
const device = reactive({ id: 'D001', name: 'M5Stack-001', model: 'M5Stack-Grey', firmware: 'v2.1.0', online: true, lastActive: '2026-03-28 10:00' })
</script>
<style scoped>.container { padding: 16px; }</style>
