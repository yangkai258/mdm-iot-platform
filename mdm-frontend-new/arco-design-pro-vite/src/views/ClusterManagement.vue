<template>
  <div class="container">
    <a-card>
      <template #title><icon-node-tree /> 集群管理</template>
      <a-row :gutter="16">
        <a-col :span="6"><a-card><a-statistic title="节点数" :value="stats.nodes" /></a-card></a-col>
        <a-col :span="6"><a-card><a-statistic title="运行中" :value="stats.running" /></a-card></a-col>
      </a-row>
      <a-table :columns="columns" :data="clusters" style="margin-top: 16px">
        <template #status="{ record }"><a-tag :color="record.status === 'healthy' ? 'green' : 'red'">{{ record.status }}</a-tag></template>
      </a-table>
    </a-card>
  </div>
</template>
<script setup>
import { ref, reactive } from 'vue'
const stats = reactive({ nodes: 10, running: 9 })
const columns = [{ title: '节点名称', dataIndex: 'name' }, { title: 'IP', dataIndex: 'ip' }, { title: '状态', slotName: 'status' }]
const clusters = ref([{ id: 1, name: 'node-1', ip: '192.168.1.1', status: 'healthy' }])
</script>
<style scoped>.container { padding: 16px; }</style>
