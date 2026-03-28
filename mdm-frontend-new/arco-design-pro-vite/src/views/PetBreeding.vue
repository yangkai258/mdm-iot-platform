<template>
  <div class="container">
    <a-card>
      <template #title><icon-team /> 宠物繁育</template>
      <a-button type="primary" @click="handleAdd">添加繁育记录</a-button>
      <a-table :columns="columns" :data="records" style="margin-top: 16px">
        <template #actions="{ record }"><a-link @click="handleView(record)">查看</a-link></template>
      </a-table>
    </a-card>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const columns = [{ title: '宠物名', dataIndex: 'petName' }, { title: '繁育日期', dataIndex: 'date' }, { title: '状态', dataIndex: 'status' }, { title: '操作', slotName: 'actions' }]
const records = ref([{ id: 1, petName: '小白', date: '2026-03-15', status: '待配对' }])
const handleAdd = () => {}
const handleView = (r) => {}
</script>
<style scoped>.container { padding: 16px; }</style>
