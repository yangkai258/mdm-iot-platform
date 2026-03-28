<template>
  <div class="container">
    <a-card>
      <template #title><icon-settings /> 任务自动化</template>
      <template #extra>
        <a-button type="primary" @click="handleCreate"><template #icon><icon-plus /></template>新建任务</a-button>
      </template>
      <a-table :columns="columns" :data="tasks">
        <template #trigger="{ record }">
          <a-tag>{{ record.trigger }}</a-tag>
        </template>
        <template #actions="{ record }">
          <a-switch v-model="record.enabled" />
          <a-link @click="handleEdit(record)">编辑</a-link>
        </template>
      </a-table>
    </a-card>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const columns = [{ title: '任务名称', dataIndex: 'name' }, { title: '触发条件', slotName: 'trigger' }, { title: '操作', slotName: 'actions' }]
const tasks = ref([{ id: 1, name: '自动喂食', trigger: '定时', enabled: true }])
const handleCreate = () => {}
const handleEdit = (r) => {}
</script>
<style scoped>.container { padding: 16px; }</style>
