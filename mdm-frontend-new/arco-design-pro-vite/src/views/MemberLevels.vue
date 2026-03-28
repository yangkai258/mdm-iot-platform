<template>
  <div class="container">
    <a-card>
      <template #title><icon-level /> 会员等级管理</template>
      <template #extra><a-button type="primary" @click="handleCreate"><template #icon><icon-plus /></template>添加等级</a-button></template>
      <a-table :columns="columns" :data="levels">
        <template #actions="{ record }"><a-link @click="handleEdit(record)">编辑</a-link></template>
      </a-table>
    </a-card>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const columns = [{ title: '等级名称', dataIndex: 'name' }, { title: '所需成长值', dataIndex: 'threshold' }, { title: '折扣', dataIndex: 'discount' }, { title: '操作', slotName: 'actions' }]
const levels = ref([{ id: 1, name: '白银', threshold: 1000, discount: '95%' }])
const handleCreate = () => {}
const handleEdit = (r) => {}
</script>
<style scoped>.container { padding: 16px; }</style>
