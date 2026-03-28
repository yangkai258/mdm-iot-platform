<template>
  <div class="container">
    <a-card>
      <template #title>
        <a-space><icon-location /> 地理围栏</a-space>
      </template>
      <template #extra>
        <a-button type="primary" @click="handleCreate">
          <template #icon><icon-plus /></template>
          新建围栏
        </a-button>
      </template>

      <a-row :gutter="16" style="margin-bottom: 16px">
        <a-col :span="6">
          <a-card><a-statistic title="围栏数量" :value="stats.total" /></a-card>
        </a-col>
        <a-col :span="6">
          <a-card><a-statistic title="触发次数" :value="stats.triggered" /></a-card>
        </a-col>
      </a-row>

      <a-table :columns="columns" :data="fences">
        <template #type="{ record }">
          <a-tag :color="record.type === 'enter' ? 'green' : 'orange'">{{ record.type === 'enter' ? '进入' : '离开' }}</a-tag>
        </template>
        <template #enabled="{ record }">
          <a-switch v-model="record.enabled" />
        </template>
        <template #actions="{ record }">
          <a-link @click="handleEdit(record)">编辑</a-link>
          <a-link @click="handleDelete(record)">删除</a-link>
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'

const stats = reactive({ total: 5, triggered: 12 })
const columns = [
  { title: '名称', dataIndex: 'name' },
  { title: '类型', slotName: 'type' },
  { title: '范围', dataIndex: 'range' },
  { title: '启用', slotName: 'enabled' },
  { title: '操作', slotName: 'actions' }
]
const fences = ref([
  { id: 1, name: '家', type: 'enter', range: '500米', enabled: true }
])

const handleCreate = () => { }
const handleEdit = (r) => { }
const handleDelete = (r) => { }
</script>

<style scoped>
.container { padding: 16px; }
</style>
