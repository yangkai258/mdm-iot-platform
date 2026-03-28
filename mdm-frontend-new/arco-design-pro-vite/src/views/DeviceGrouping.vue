<template>
  <div class="container">
    <a-card>
      <template #title><icon-folder /> 设备分组</template>
      <template #extra><a-button type="primary" @click="handleCreate"><template #icon><icon-plus /></template>新建分组</a-button></template>
      <a-table :columns="columns" :data="groups">
        <template #actions="{ record }"><a-link @click="handleEdit(record)">编辑</a-link></template>
      </a-table>
    </a-card>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const columns = [{ title: '分组名称', dataIndex: 'name' }, { title: '设备数量', dataIndex: 'deviceCount' }, { title: '操作', slotName: 'actions' }]
const groups = ref([{ id: 1, name: '测试设备', deviceCount: 10 }])
const handleCreate = () => {}
const handleEdit = (r) => {}
</script>
<style scoped>.container { padding: 16px; }</style>
