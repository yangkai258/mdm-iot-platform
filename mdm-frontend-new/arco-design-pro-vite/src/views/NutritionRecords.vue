<template>
  <div class="container">
    <a-card>
      <template #title><icon-food /> 营养记录</template>
      <a-table :columns="columns" :data="records">
        <template #actions="{ record }"><a-link @click="handleAdd">添加</a-link></template>
      </a-table>
    </a-card>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const columns = [{ title: '日期', dataIndex: 'date' }, { title: '宠物', dataIndex: 'pet' }, { title: '食物', dataIndex: 'food' }, { title: '热量', dataIndex: 'calories' }, { title: '操作', slotName: 'actions' }]
const records = ref([{ id: 1, date: '2026-03-28', pet: '小白', food: '狗粮100g', calories: '300kcal' }])
const handleAdd = () => {}
</script>
<style scoped>.container { padding: 16px; }</style>
