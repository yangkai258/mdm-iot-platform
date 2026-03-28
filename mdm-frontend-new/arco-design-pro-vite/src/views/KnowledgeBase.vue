<template>
  <div class="container">
    <a-card>
      <template #title><icon-book /> 知识库</template>
      <template #extra><a-button type="primary" @click="handleAdd"><template #icon><icon-plus /></template>添加知识</a-button></template>
      <a-table :columns="columns" :data="knowledge">
        <template #actions="{ record }"><a-link @click="handleView(record)">查看</a-link></template>
      </a-table>
    </a-card>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const columns = [{ title: '标题', dataIndex: 'title' }, { title: '分类', dataIndex: 'category' }, { title: '更新时间', dataIndex: 'updatedAt' }, { title: '操作', slotName: 'actions' }]
const knowledge = ref([{ id: 1, title: '宠物饲养指南', category: '科普', updatedAt: '2026-03-28' }])
const handleAdd = () => {}
const handleView = (r) => {}
</script>
<style scoped>.container { padding: 16px; }</style>
