<template>
  <div class="container">
    <a-card>
      <template #title><icon-file-text /> 内容管理</template>
      <template #extra><a-button type="primary" @click="handleCreate"><template #icon><icon-plus /></template>创建内容</a-button></template>
      <a-table :columns="columns" :data="contents">
        <template #status="{ record }"><a-tag :color="record.published ? 'green' : 'gray'">{{ record.published ? '已发布' : '草稿' }}</a-tag></template>
      </a-table>
    </a-card>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const columns = [{ title: '标题', dataIndex: 'title' }, { title: '类型', dataIndex: 'type' }, { title: '状态', slotName: 'status' }]
const contents = ref([{ id: 1, title: '宠物饲养指南', type: 'article', published: true }])
const handleCreate = () => {}
</script>
<style scoped>.container { padding: 16px; }</style>
