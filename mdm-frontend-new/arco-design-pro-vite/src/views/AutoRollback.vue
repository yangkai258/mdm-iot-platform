<template>
  <div class="container">
    <a-card>
      <template #title><icon-rollback /> 自动回滚规则</template>
      <template #extra><a-button type="primary" @click="handleCreate"><template #icon><icon-plus /></template>新建规则</a-button></template>
      <a-table :columns="columns" :data="rules">
        <template #actions="{ record }"><a-link @click="handleEdit(record)">编辑</a-link></template>
      </a-table>
    </a-card>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const columns = [{ title: '规则名称', dataIndex: 'name' }, { title: '触发条件', dataIndex: 'condition' }, { title: '操作', slotName: 'actions' }]
const rules = ref([{ id: 1, name: '错误率>5%', condition: '错误率>5%持续5分钟' }])
const handleCreate = () => {}
const handleEdit = (r) => {}
</script>
<style scoped>.container { padding: 16px; }</style>
