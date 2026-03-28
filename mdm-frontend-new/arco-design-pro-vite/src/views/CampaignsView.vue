<template>
  <div class="container">
    <a-card>
      <template #title><icon-promotion /> 活动管理</template>
      <template #extra><a-button type="primary" @click="handleCreate"><template #icon><icon-plus /></template>创建活动</a-button></template>
      <a-table :columns="columns" :data="campaigns">
        <template #status="{ record }"><a-tag :color="record.active ? 'green' : 'gray'">{{ record.active ? '进行中' : '已结束' }}</a-tag></template>
      </a-table>
    </a-card>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const columns = [{ title: '活动名称', dataIndex: 'name' }, { title: '时间', dataIndex: 'time' }, { title: '状态', slotName: 'status' }]
const campaigns = ref([{ id: 1, name: '春节促销', time: '2026-02-01', active: false }])
const handleCreate = () => {}
</script>
<style scoped>.container { padding: 16px; }</style>
