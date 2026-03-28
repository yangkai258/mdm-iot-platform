<template>
  <div class="container">
    <a-card>
      <template #title><icon-heart /> 健康记录</template>
      <a-row :gutter="16">
        <a-col :span="6"><a-card><a-statistic title="体重" value="12.5kg" /></a-card></a-col>
        <a-col :span="6"><a-card><a-statistic title="心率" value="85bpm" /></a-card></a-col>
      </a-row>
      <a-table :columns="columns" :data="records" style="margin-top: 16px" />
    </a-card>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const columns = [{ title: '日期', dataIndex: 'date' }, { title: '体重', dataIndex: 'weight' }, { title: '心率', dataIndex: 'heartRate' }]
const records = ref([{ date: '2026-03-28', weight: '12.5kg', heartRate: '85bpm' }])
</script>
<style scoped>.container { padding: 16px; }</style>
