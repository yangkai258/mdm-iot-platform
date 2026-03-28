<template>
  <div class="container">
    <a-card>
      <template #title><icon-file /> 固件兼容矩阵</template>
      <a-table :columns="columns" :data="matrix" />
    </a-card>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const columns = [{ title: '型号', dataIndex: 'model' }, { title: 'v2.0.0', dataIndex: 'v200' }, { title: 'v2.1.0', dataIndex: 'v210' }]
const matrix = ref([{ model: 'M5Stack-1', v200: '✓', v210: '✓' }])
</script>
<style scoped>.container { padding: 16px; }</style>
