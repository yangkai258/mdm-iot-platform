<template>
  <div class="container">
    <a-card>
      <template #title>
        <a-space><icon-grid /> 固件兼容矩阵</a-space>
      </template>
      <a-table :columns="columns" :data="data" :pagination="false" size="small">
        <template #device-type="{ record }">
          <a-tag>{{ record.deviceType }}</a-tag>
        </template>
        <template #compatible="{ record }">
          <a-tag :color="record.compatible ? 'green' : 'red'">
            {{ record.compatible ? '兼容' : '不兼容' }}
          </a-tag>
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const columns = [
  { title: '设备型号', dataIndex: 'deviceType', slotName: 'device-type' },
  { title: '固件版本', dataIndex: 'firmwareVersion' },
  { title: '最低要求版本', dataIndex: 'minRequired' },
  { title: '状态', slotName: 'compatible' },
]
const data = ref([
  { deviceType: 'M5Stack-Core2', firmwareVersion: 'v2.1.0', minRequired: 'v2.0.0', compatible: true },
  { deviceType: 'M5Stack-Atom', firmwareVersion: 'v1.8.5', minRequired: 'v1.8.0', compatible: true },
  { deviceType: 'Generic-ESP32', firmwareVersion: 'v2.0.0', minRequired: 'v2.1.0', compatible: false },
])
</script>

<style scoped>
.container { padding: 16px; }
</style>
