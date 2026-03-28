<template>
  <div class="container">
    <a-card>
      <template #title><icon-robot /> AI模型版本管理</template>
      <template #extra><a-button type="primary" @click="handleUpload"><template #icon><icon-upload /></template>上传模型</a-button></template>
      <a-table :columns="columns" :data="models">
        <template #status="{ record }"><a-tag :color="record.current ? 'green' : 'gray'">{{ record.current ? '当前' : '历史' }}</a-tag></template>
        <template #actions="{ record }"><a-button v-if="!record.current" size="small">回滚</a-button></template>
      </a-table>
    </a-card>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const columns = [{ title: '模型名称', dataIndex: 'name' }, { title: '版本', dataIndex: 'version' }, { title: '状态', slotName: 'status' }, { title: '操作', slotName: 'actions' }]
const models = ref([{ id: 1, name: 'EmotionNet', version: 'v2.1', current: true }])
const handleUpload = () => {}
</script>
<style scoped>.container { padding: 16px; }</style>
