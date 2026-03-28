<template>
  <div class="container">
    <a-card>
      <template #title>
        <a-space><icon-history /> 组织变更时间线</a-space>
      </template>

      <a-space style="margin-bottom: 16px">
        <a-select v-model="filterDept" placeholder="选择部门" allow-clear style="width: 200px">
          <a-option value="d1">研发部</a-option>
          <a-option value="d2">运维部</a-option>
        </a-select>
        <a-range-picker v-model="filterDate" />
        <a-button type="primary">查询</a-button>
        <a-button @click="handleExport">
          <template #icon><icon-download /></template>
          导出记录
        </a-button>
      </a-space>

      <a-timeline>
        <a-timeline-item v-for="change in changes" :key="change.id" :color="getChangeColor(change.type)">
          <template #label>{{ change.time }}</template>
          <a-card size="small">
            <template #title>
              <a-space>
                <a-tag :color="getChangeColor(change.type)">{{ getChangeTypeLabel(change.type) }}</a-tag>
                <span>{{ change.title }}</span>
              </a-space>
            </template>
            <a-descriptions :column="2" size="small">
              <a-descriptions-item label="执行人">{{ change.executor }}</a-descriptions-item>
              <a-descriptions-item label="影响人数">{{ change.affectedCount }} 人</a-descriptions-item>
            </a-descriptions>
            <a-space style="margin-top: 8px">
              <a-link @click="handleViewDetail(change)">详情</a-link>
              <a-link @click="handleViewAffected(change)">查看受影响人员</a-link>
            </a-space>
          </a-card>
        </a-timeline-item>
      </a-timeline>
    </a-card>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const filterDept = ref('')
const filterDate = ref([])

const changes = ref([
  { id: 1, time: '2026-03-28 10:00', type: 'merge', title: '前端组与UI组合并', executor: 'admin', affectedCount: 15 },
  { id: 2, time: '2026-03-27 15:00', type: 'split', title: '运维部分拆', executor: 'admin', affectedCount: 8 },
  { id: 3, time: '2026-03-26 09:00', type: 'transfer', title: '人员调岗', executor: 'hr', affectedCount: 5 }
])

const getChangeColor = (type) => ({ merge: 'green', split: 'orange', transfer: 'blue' }[type] || 'gray')
const getChangeTypeLabel = (type) => ({ merge: '合并', split: '拆分', transfer: '调岗' }[type] || type)

const handleExport = () => { }
const handleViewDetail = (c) => { }
const handleViewAffected = (c) => { }
</script>

<style scoped>
.container { padding: 16px; }
</style>
