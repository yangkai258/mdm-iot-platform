<template>
  <div class="container">
    <a-card>
      <template #title><icon-mind-mapping /> AI决策日志</template>
      <a-table :columns="columns" :data="logs">
        <template #decision="{ record }">{{ record.decision }}</template>
      </a-table>
    </a-card>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const columns = [{ title: '时间', dataIndex: 'time' }, { title: '设备', dataIndex: 'deviceId' }, { title: '决策', slotName: 'decision' }, { title: '置信度', dataIndex: 'confidence' }]
const logs = ref([{ time: '2026-03-28 10:00', deviceId: 'D001', decision: '执行动作A', confidence: '95%' }])
</script>
<style scoped>.container { padding: 16px; }</style>
