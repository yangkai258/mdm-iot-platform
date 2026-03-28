<template>
  <div class="container">
    <a-card>
      <template #title><icon-robot /> OpenClaw版本管理</template>
      <a-table :columns="columns" :data="versions">
        <template #status="{ record }"><a-tag :color="record.current ? 'green' : 'gray'">{{ record.current ? '当前' : '历史' }}</a-tag></template>
        <template #actions="{ record }"><a-button v-if="!record.current" size="small">回滚</a-button></template>
      </a-table>
    </a-card>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const columns = [{ title: '版本', dataIndex: 'version' }, { title: '发布时间', dataIndex: 'releaseTime' }, { title: '状态', slotName: 'status' }, { title: '操作', slotName: 'actions' }]
const versions = ref([{ id: 1, version: 'v2.1.0', releaseTime: '2026-03-28', current: true }])
</script>
<style scoped>.container { padding: 16px; }</style>
