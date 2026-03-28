<template>
  <div class="container">
    <a-card>
      <template #title><icon-book-open /> AI训练任务</template>
      <template #extra><a-button type="primary" @click="handleCreate"><template #icon><icon-plus /></template>新建任务</a-button></template>
      <a-table :columns="columns" :data="tasks">
        <template #status="{ record }"><a-tag :color="getStatusColor(record.status)">{{ record.status }}</a-tag></template>
        <template #actions="{ record }"><a-button size="small">查看</a-button></template>
      </a-table>
    </a-card>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const columns = [{ title: '任务名称', dataIndex: 'name' }, { title: '模型', dataIndex: 'model' }, { title: '进度', dataIndex: 'progress' }, { title: '状态', slotName: 'status' }, { title: '操作', slotName: 'actions' }]
const tasks = ref([{ id: 1, name: '情感识别训练', model: 'EmotionNet', progress: '75%', status: '进行中' }])
const getStatusColor = (s) => ({ '进行中': 'blue', '已完成': 'green', '失败': 'red' }[s] || 'gray')
const handleCreate = () => {}
</script>
<style scoped>.container { padding: 16px; }</style>
