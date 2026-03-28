<template>
  <div class="container">
    <a-card>
      <template #title><icon-bell /> 消息通知列表</template>
      <a-table :columns="columns" :data="notifications">
        <template #status="{ record }"><a-tag :color="record.read ? 'gray' : 'blue'">{{ record.read ? '已读' : '未读' }}</a-tag></template>
      </a-table>
    </a-card>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const columns = [{ title: '标题', dataIndex: 'title' }, { title: '时间', dataIndex: 'time' }, { title: '状态', slotName: 'status' }]
const notifications = ref([{ id: 1, title: '系统升级通知', time: '2026-03-28', read: false }])
</script>
<style scoped>.container { padding: 16px; }</style>
