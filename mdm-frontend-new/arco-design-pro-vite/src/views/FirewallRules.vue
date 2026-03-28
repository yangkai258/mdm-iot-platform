<template>
  <div class="container">
    <a-card>
      <template #title><icon-fire /> 防火墙规则</template>
      <template #extra>
        <a-button type="primary" @click="handleCreate"><template #icon><icon-plus /></template>添加规则</a-button>
      </template>
      <a-table :columns="columns" :data="rules">
        <template #action="{ record }"><a-tag :color="record.action === 'allow' ? 'green' : 'red'">{{ record.action }}</a-tag></template>
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
const columns = [{ title: '名称', dataIndex: 'name' }, { title: '源IP', dataIndex: 'srcIp' }, { title: '目标', dataIndex: 'dst' }, { title: '动作', slotName: 'action' }, { title: '操作', slotName: 'actions' }]
const rules = ref([{ id: 1, name: '允许HTTP', srcIp: '0.0.0.0/0', dst: '80,443', action: 'allow' }])
const handleCreate = () => {}
const handleEdit = (r) => {}
const handleDelete = (r) => {}
</script>
<style scoped>.container { padding: 16px; }</style>
