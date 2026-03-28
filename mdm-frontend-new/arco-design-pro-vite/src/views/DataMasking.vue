<template>
  <div class="container">
    <a-card>
      <template #title><icon-eye-invisible /> 数据脱敏规则</template>
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
const columns = [{ title: '规则名称', dataIndex: 'name' }, { title: '字段', dataIndex: 'field' }, { title: '启用', slotName: 'enabled' }, { title: '操作', slotName: 'actions' }]
const rules = ref([{ id: 1, name: '手机号脱敏', field: 'phone', enabled: true }])
const handleCreate = () => {}
const handleEdit = (r) => {}
</script>
<style scoped>.container { padding: 16px; }</style>
