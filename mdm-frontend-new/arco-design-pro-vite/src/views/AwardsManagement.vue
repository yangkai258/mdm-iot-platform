<template>
  <div class="container">
    <a-card>
      <template #title><icon-trophy /> 奖项管理</template>
      <template #extra>
        <a-button type="primary" @click="handleCreate"><template #icon><icon-plus /></template>添加奖项</a-button>
      </template>
      <a-table :columns="columns" :data="awards">
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
const columns = [{ title: '奖项名称', dataIndex: 'name' }, { title: '等级', dataIndex: 'level' }, { title: '操作', slotName: 'actions' }]
const awards = ref([{ id: 1, name: '最佳表现', level: '金牌' }])
const handleCreate = () => {}
const handleEdit = (r) => {}
const handleDelete = (r) => {}
</script>
<style scoped>.container { padding: 16px; }</style>
