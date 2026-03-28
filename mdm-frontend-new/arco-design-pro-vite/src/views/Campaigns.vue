<template>
  <div class="container">
    <a-card>
      <template #title><icon-promotion /> 营销活动</template>
      <template #extra><a-button type="primary" @click="handleCreate"><template #icon><icon-plus /></template>新建活动</a-button></template>
      <a-table :columns="columns" :data="campaigns">
        <template #status="{ record }"><a-tag :color="record.status === 'running' ? 'green' : 'gray'">{{ record.status === 'running' ? '进行中' : '已结束' }}</a-tag></template>
        <template #actions="{ record }"><a-link @click="handleEdit(record)">编辑</a-link></template>
      </a-table>
    </a-card>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const columns = [{ title: '活动名称', dataIndex: 'name' }, { title: '时间', dataIndex: 'time' }, { title: '状态', slotName: 'status' }, { title: '操作', slotName: 'actions' }]
const campaigns = ref([{ id: 1, name: '春节促销', time: '2026-02-01', status: 'ended' }])
const handleCreate = () => {}
const handleEdit = (r) => {}
</script>
<style scoped>.container { padding: 16px; }</style>
