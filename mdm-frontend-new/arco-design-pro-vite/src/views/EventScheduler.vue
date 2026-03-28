<template>
  <div class="container">
    <a-card>
      <template #title><icon-clock /> 事件调度器</template>
      <template #extra>
        <a-button type="primary" @click="handleCreate"><template #icon><icon-plus /></template>新建事件</a-button>
      </template>
      <a-table :columns="columns" :data="events">
        <template #enabled="{ record }"><a-switch v-model="record.enabled" /></template>
        <template #actions="{ record }">
          <a-link @click="handleEdit(record)">编辑</a-link>
          <a-link @click="handleDelete(record)">删除</a-link>
        </template>
      </a-table>
    </a-card>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const columns = [{ title: '事件名', dataIndex: 'name' }, { title: '触发时间', dataIndex: 'time' }, { title: '启用', slotName: 'enabled' }, { title: '操作', slotName: 'actions' }]
const events = ref([{ id: 1, name: '每日报告', time: '09:00', enabled: true }])
const handleCreate = () => {}
const handleEdit = (r) => {}
const handleDelete = (r) => {}
</script>
<style scoped>.container { padding: 16px; }</style>
