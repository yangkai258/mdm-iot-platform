<template>
  <div class="container">
    <a-card>
      <template #title><icon-global /> 区域配置</template>
      <template #extra><a-button type="primary" @click="handleCreate"><template #icon><icon-plus /></template>添加区域</a-button></template>
      <a-table :columns="columns" :data="regions">
        <template #actions="{ record }"><a-link @click="handleEdit(record)">编辑</a-link></template>
      </a-table>
    </a-card>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const columns = [{ title: '区域名称', dataIndex: 'name' }, { title: '代码', dataIndex: 'code' }, { title: '操作', slotName: 'actions' }]
const regions = ref([{ id: 1, name: '华北', code: 'CN-NORTH' }])
const handleCreate = () => {}
const handleEdit = (r) => {}
</script>
<style scoped>.container { padding: 16px; }</style>
