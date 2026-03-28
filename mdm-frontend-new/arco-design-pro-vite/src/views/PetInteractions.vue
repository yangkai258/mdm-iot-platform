<template>
  <div class="container">
    <a-card>
      <template #title><icon-message /> 宠物互动记录</template>
      <a-table :columns="columns" :data="interactions" />
    </a-card>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const columns = [{ title: '时间', dataIndex: 'time' }, { title: '宠物', dataIndex: 'pet' }, { title: '互动类型', dataIndex: 'type' }, { title: '内容', dataIndex: 'content' }]
const interactions = ref([{ time: '2026-03-28', pet: '小白', type: '语音', content: '今天天气很好' }])
</script>
<style scoped>.container { padding: 16px; }</style>
