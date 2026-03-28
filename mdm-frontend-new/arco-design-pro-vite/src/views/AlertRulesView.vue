<template>
  <div class="container">
    <a-card>
      <template #title><icon-warning /> 告警规则</template>
      <template #extra><a-button type="primary" @click="handleCreate"><template #icon><icon-plus /></template>新建规则</a-button></template>
      <a-table :columns="columns" :data="rules">
        <template #enabled="{ record }"><a-switch v-model="record.enabled" /></template>
        <template #actions="{ record }"><a-link @click="handleEdit(record)">编辑</a-link></template>
      </a-table>
    </a-card>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const columns = [{ title: '规则名称', dataIndex: 'name' }, { title: '条件', dataIndex: 'condition' }, { title: '启用', slotName: 'enabled' }, { title: '操作', slotName: 'actions' }]
const rules = ref([{ id: 1, name: '温度告警', condition: '>38°C', enabled: true }])
const handleCreate = () => {}
const handleEdit = (r) => {}
</script>
<style scoped>.container { padding: 16px; }</style>
