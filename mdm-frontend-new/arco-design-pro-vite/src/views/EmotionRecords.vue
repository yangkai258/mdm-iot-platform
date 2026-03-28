<template>
  <div class="container">
    <a-card>
      <template #title><icon-heart /> 情绪记录</template>
      <a-table :columns="columns" :data="records">
        <template #emotion="{ record }"><a-tag :color="getEmotionColor(record.emotion)">{{ record.emotion }}</a-tag></template>
      </a-table>
    </a-card>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const columns = [{ title: '时间', dataIndex: 'time' }, { title: '设备', dataIndex: 'deviceId' }, { title: '情绪', slotName: 'emotion' }, { title: '置信度', dataIndex: 'confidence' }]
const records = ref([{ time: '2026-03-28 10:00', deviceId: 'D001', emotion: '开心', confidence: '92%' }])
const getEmotionColor = (e) => ({ 开心: 'green', 平静: 'blue', 难过: 'orange', 生气: 'red' }[e] || 'gray')
</script>
<style scoped>.container { padding: 16px; }</style>
