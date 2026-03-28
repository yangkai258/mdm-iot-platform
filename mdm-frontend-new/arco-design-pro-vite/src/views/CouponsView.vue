<template>
  <div class="container">
    <a-card>
      <template #title><icon-ticket /> 优惠券管理</template>
      <template #extra><a-button type="primary" @click="handleCreate"><template #icon><icon-plus /></template>创建优惠券</a-button></template>
      <a-table :columns="columns" :data="coupons">
        <template #status="{ record }"><a-tag :color="record.status === 'active' ? 'green' : 'gray'">{{ record.status === 'active' ? '生效中' : '已过期' }}</a-tag></template>
        <template #actions="{ record }"><a-link @click="handleEdit(record)">编辑</a-link></template>
      </a-table>
    </a-card>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const columns = [{ title: '名称', dataIndex: 'name' }, { title: '面值', dataIndex: 'value' }, { title: '状态', slotName: 'status' }, { title: '操作', slotName: 'actions' }]
const coupons = ref([{ id: 1, name: '新人券', value: '¥10', status: 'active' }])
const handleCreate = () => {}
const handleEdit = (r) => {}
</script>
<style scoped>.container { padding: 16px; }</style>
