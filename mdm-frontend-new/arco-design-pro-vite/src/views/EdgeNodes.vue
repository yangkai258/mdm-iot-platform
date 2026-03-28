<template>
  <div class="container">
    <a-card>
      <template #title><icon-cloud /> 边缘节点</template>
      <template #extra><a-button type="primary" @click="handleAdd"><template #icon><icon-plus /></template>添加节点</a-button></template>
      <a-table :columns="columns" :data="nodes">
        <template #status="{ record }"><a-tag :color="record.online ? 'green' : 'gray'">{{ record.online ? '在线' : '离线' }}</a-tag></template>
      </a-table>
    </a-card>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const columns = [{ title: '节点ID', dataIndex: 'id' }, { title: '名称', dataIndex: 'name' }, { title: '状态', slotName: 'status' }]
const nodes = ref([{ id: 'E001', name: '边缘节点1', online: true }])
const handleAdd = () => {}
</script>
<style scoped>.container { padding: 16px; }</style>
