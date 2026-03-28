<template>
  <div class="container">
    <a-card>
      <template #title><icon-card /> 订阅管理</template>
      <template #extra><a-button type="primary" @click="handleCreate"><template #icon><icon-plus /></template>新建订阅</a-button></template>
      <a-table :columns="columns" :data="subscriptions">
        <template #status="{ record }"><a-tag :color="record.status === 'active' ? 'green' : 'gray'">{{ record.status === 'active' ? '生效中' : '已过期' }}</a-tag></template>
      </a-table>
    </a-card>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const columns = [{ title: '用户', dataIndex: 'user' }, { title: '套餐', dataIndex: 'plan' }, { title: '状态', slotName: 'status' }]
const subscriptions = ref([{ id: 1, user: '张三', plan: '专业版', status: 'active' }])
const handleCreate = () => {}
</script>
<style scoped>.container { padding: 16px; }</style>
