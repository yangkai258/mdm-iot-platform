<template>
  <div class="container">
    <a-card>
      <template #title><icon-gift /> 行为奖励</template>
      <template #extra>
        <a-button type="primary" @click="handleCreate"><template #icon><icon-plus /></template>新建奖励</a-button>
      </template>
      <a-table :columns="columns" :data="rewards">
        <template #actions="{ record }">
          <a-link @click="handleEdit(record)">编辑</a-link>
          <a-link @click="handleDelete(record)">删除</a-link>
        </template>
      </a-table>
    </a-card>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const columns = [{ title: '行为', dataIndex: 'behavior' }, { title: '奖励', dataIndex: 'reward' }, { title: '操作', slotName: 'actions' }]
const rewards = ref([{ id: 1, behavior: '完成训练', reward: '积分x10' }])
const handleCreate = () => {}
const handleEdit = (r) => {}
const handleDelete = (r) => {}
</script>
<style scoped>.container { padding: 16px; }</style>
