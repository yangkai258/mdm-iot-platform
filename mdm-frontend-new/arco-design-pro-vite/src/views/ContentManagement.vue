<template>
  <div class="container">
    <a-card>
      <template #title><icon-app /> 内容管理</template>
      <template #extra><a-button type="primary" @click="handleAdd"><template #icon><icon-plus /></template>添加内容</a-button></template>
      <a-table :columns="columns" :data="contents">
        <template #type="{ record }"><a-tag>{{ record.type }}</a-tag></template>
        <template #actions="{ record }"><a-link @click="handleEdit(record)">编辑</a-link></template>
      </a-table>
    </a-card>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const columns = [{ title: '标题', dataIndex: 'title' }, { title: '类型', slotName: 'type' }, { title: '状态', dataIndex: 'status' }, { title: '操作', slotName: 'actions' }]
const contents = ref([{ id: 1, title: '宠物护理指南', type: 'article', status: '已发布' }])
const handleAdd = () => {}
const handleEdit = (r) => {}
</script>
<style scoped>.container { padding: 16px; }</style>
