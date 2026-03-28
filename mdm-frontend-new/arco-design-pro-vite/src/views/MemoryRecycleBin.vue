<template>
  <div class="container">
    <a-card>
      <template #title>
        <a-space><icon-delete /> 记忆回收站</a-space>
      </template>

      <a-space style="margin-bottom: 16px">
        <a-button type="primary" @click="handleBatchRestore">
          <template #icon><icon-arrow-up /></template>
          批量恢复
        </a-button>
        <a-button type="primary" status="danger" @click="handleBatchDelete">
          <template #icon><icon-delete /></template>
          永久删除
        </a-button>
      </a-space>

      <a-table :columns="columns" :data="memories" :row-selection="{ type: 'checkbox' }">
        <template #content="{ record }">
          <div class="memory-preview">
            <a-tag :color="record.type === 'interaction' ? 'blue' : 'green'">{{ record.type }}</a-tag>
            <span>{{ record.preview }}</span>
          </div>
        </template>
        <template #expireDays="{ record }">
          <a-tag :color="record.expireDays <= 7 ? 'orange' : 'gray'">{{ record.expireDays }}天后过期</a-tag>
        </template>
        <template #actions="{ record }">
          <a-link @click="handleRestore(record)">恢复</a-link>
          <a-link @click="handlePermanentDelete(record)">永久删除</a-link>
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const columns = [
  { title: '记忆ID', dataIndex: 'id', width: 100 },
  { title: '内容预览', slotName: 'content' },
  { title: '删除时间', dataIndex: 'deletedAt', width: 180 },
  { title: '过期时间', slotName: 'expireDays', width: 120 },
  { title: '操作', slotName: 'actions', width: 180 }
]

const memories = ref([
  { id: 'M001', type: 'interaction', preview: '与主人互动，学习新动作', deletedAt: '2026-03-25 10:00', expireDays: 25 },
  { id: 'M002', type: 'behavior', preview: '建立对厨房区域的认知', deletedAt: '2026-03-26 15:00', expireDays: 26 }
])

const handleBatchRestore = () => { }
const handleBatchDelete = () => { }
const handleRestore = (r) => { }
const handlePermanentDelete = (r) => { }
</script>

<style scoped>
.container { padding: 16px; }
.memory-preview { display: flex; align-items: center; gap: 8px; }
</style>
