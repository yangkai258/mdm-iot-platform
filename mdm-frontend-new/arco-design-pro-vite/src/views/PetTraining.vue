<template>
  <div class="container">
    <a-card>
      <template #title><icon-book-open /> 宠物训练</template>
      <template #extra><a-button type="primary" @click="handleCreate"><template #icon><icon-plus /></template>新建训练</a-button></template>
      <a-table :columns="columns" :data="trainings">
        <template #progress="{ record }"><a-progress :percent="record.progress" /></template>
        <template #actions="{ record }"><a-link @click="handleStart(record)">开始</a-link></template>
      </a-table>
    </a-card>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const columns = [{ title: '训练名称', dataIndex: 'name' }, { title: '宠物', dataIndex: 'pet' }, { title: '进度', slotName: 'progress' }, { title: '操作', slotName: 'actions' }]
const trainings = ref([{ id: 1, name: '坐下训练', pet: '小白', progress: 60 }])
const handleCreate = () => {}
const handleStart = (r) => {}
</script>
<style scoped>.container { padding: 16px; }</style>
