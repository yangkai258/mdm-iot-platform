<template>
  <div class="container">
    <a-card>
      <template #title><icon-cloud-download /> OTA管理</template>
      <template #extra><a-button type="primary" @click="handleUpgrade"><template #icon><icon-upload /></template>创建升级任务</a-button></template>
      <a-table :columns="columns" :data="tasks">
        <template #status="{ record }"><a-tag :color="getStatusColor(record.status)">{{ record.status }}</a-tag></template>
      </a-table>
    </a-card>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const columns = [{ title: '任务名称', dataIndex: 'name' }, { title: '固件版本', dataIndex: 'version' }, { title: '进度', dataIndex: 'progress' }, { title: '状态', slotName: 'status' }]
const tasks = ref([{ id: 1, name: '批量升级', version: 'v2.1.0', progress: '60%', status: '进行中' }])
const getStatusColor = (s) => ({ '进行中': 'blue', '已完成': 'green', '失败': 'red' }[s] || 'gray')
const handleUpgrade = () => {}
</script>
<style scoped>.container { padding: 16px; }</style>
