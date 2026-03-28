<template>
  <div class="container">
    <a-card>
      <template #title><icon-star /> 积分规则配置</template>
      <template #extra><a-button type="primary" @click="handleCreate"><template #icon><icon-plus /></template>添加规则</a-button></template>
      <a-table :columns="columns" :data="rules">
        <template #actions="{ record }"><a-link @click="handleEdit(record)">编辑</a-link></template>
      </a-table>
    </a-card>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const columns = [{ title: '规则名称', dataIndex: 'name' }, { title: '积分', dataIndex: 'points' }, { title: '操作', slotName: 'actions' }]
const rules = ref([{ id: 1, name: '每日登录', points: '+10' }])
const handleCreate = () => {}
const handleEdit = (r) => {}
</script>
<style scoped>.container { padding: 16px; }</style>
