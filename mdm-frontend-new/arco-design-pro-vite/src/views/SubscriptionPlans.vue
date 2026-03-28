<template>
  <div class="container">
    <a-card>
      <template #title><icon-card /> 订阅套餐</template>
      <template #extra><a-button type="primary" @click="handleCreate"><template #icon><icon-plus /></template>新建套餐</a-button></template>
      <a-table :columns="columns" :data="plans">
        <template #actions="{ record }"><a-link @click="handleEdit(record)">编辑</a-link></template>
      </a-table>
    </a-card>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const columns = [{ title: '套餐名称', dataIndex: 'name' }, { title: '价格', dataIndex: 'price' }, { title: '操作', slotName: 'actions' }]
const plans = ref([{ id: 1, name: '基础版', price: '¥29/月' }])
const handleCreate = () => {}
const handleEdit = (r) => {}
</script>
<style scoped>.container { padding: 16px; }</style>
