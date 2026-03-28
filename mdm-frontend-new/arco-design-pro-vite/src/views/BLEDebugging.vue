<template>
  <div class="container">
    <a-card>
      <template #title><icon-bluetooth /> BLE调试</template>
      <a-button type="primary" @click="handleScan">扫描设备</a-button>
      <a-divider />
      <a-list>
        <a-list-item v-for="device in devices" :key="device.mac">
          <a-list-item-meta :title="device.name" :description="device.mac" />
          <template #actions><a-button size="small">连接</a-button></template>
        </a-list-item>
      </a-list>
    </a-card>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const devices = ref([{ name: 'BLE-Device-1', mac: 'AA:BB:CC:DD:EE:FF' }])
const handleScan = () => {}
</script>
<style scoped>.container { padding: 16px; }</style>
