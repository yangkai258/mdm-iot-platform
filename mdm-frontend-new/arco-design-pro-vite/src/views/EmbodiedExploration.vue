<template>
  <div class="container">
    <a-card>
      <template #title>
        <a-space><icon-compass /> 探索任务</a-space>
      </template>
      <template #extra>
        <a-button type="primary" @click="handleCreateTask">
          <template #icon><icon-plus /></template>
          新建任务
        </a-button>
      </template>

      <a-row :gutter="16">
        <a-col :span="8">
          <a-card title="任务列表" size="small">
            <a-list>
              <a-list-item v-for="task in tasks" :key="task.id" :class="{ selected: selectedTask?.id === task.id }" @click="handleSelectTask(task)">
                <a-list-item-meta :title="task.name" :description="task.description" />
                <template #actions>
                  <a-badge :status="getTaskStatusBadge(task.status)" />
                </template>
              </a-list-item>
            </a-list>
          </a-card>
        </a-col>

        <a-col :span="10">
          <a-card title="任务详情">
            <a-descriptions :column="2" bordered v-if="selectedTask">
              <a-descriptions-item label="任务名称">{{ selectedTask.name }}</a-descriptions-item>
              <a-descriptions-item label="状态">
                <a-badge :status="getTaskStatusBadge(selectedTask.status)" :text="getTaskStatusText(selectedTask.status)" />
              </a-descriptions-item>
              <a-descriptions-item label="优先级">
                <a-tag :color="getPriorityColor(selectedTask.priority)">{{ selectedTask.priority }}</a-tag>
              </a-descriptions-item>
              <a-descriptions-item label="预计时长">{{ selectedTask.duration }}</a-descriptions-item>
              <a-descriptions-item label="起点" :span="2>{{ selectedTask.startPoint }}</a-descriptions-item>
              <a-descriptions-item label="终点" :span="2>{{ selectedTask.endPoint }}</a-descriptions-item>
            </a-descriptions>

            <a-divider>执行进度</a-divider>
            <a-progress :percent="selectedTask?.progress || 0" :show-text="true" />

            <a-space style="margin-top: 16px">
              <a-button v-if="selectedTask?.status === 'pending'" type="primary" @click="handleStart">开始</a-button>
              <a-button v-if="selectedTask?.status === 'running'" @click="handlePause">暂停</a-button>
              <a-button v-if="selectedTask?.status === 'paused'" type="primary" @click="handleResume">继续</a-button>
              <a-button v-if="selectedTask?.status === 'running' || selectedTask?.status === 'paused'" @click="handleCancel">取消</a-button>
            </a-space>
          </a-card>

          <a-card title="任务历史" style="margin-top: 16px">
            <a-table :columns="historyColumns" :data="taskHistory" size="small" :pagination="false" />
          </a-card>
        </a-col>

        <a-col :span="6">
          <a-card title="地图预览" size="small">
            <div class="mini-map">
              <svg width="200" height="150" viewBox="0 0 200 150">
                <rect x="0" y="0" width="200" height="150" fill="#e8e8e8" />
                <circle cx="100" cy="75" r="10" fill="#409EFF" />
              </svg>
            </div>
          </a-card>
        </a-col>
      </a-row>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'

const tasks = ref([
  { id: 1, name: '探索客厅', description: '探索客厅区域', status: 'pending', priority: 'medium', duration: '10分钟', startPoint: '卧室', endPoint: '客厅', progress: 0 },
  { id: 2, name: '巡逻厨房', description: '巡逻厨房区域', status: 'running', priority: 'low', duration: '15分钟', startPoint: '客厅', endPoint: '厨房', progress: 45 }
])
const selectedTask = ref(tasks.value[0])

const historyColumns = [
  { title: '任务', dataIndex: 'name' },
  { title: '状态', dataIndex: 'status' },
  { title: '完成时间', dataIndex: 'completedAt' }
]
const taskHistory = ref([
  { name: '探索书房', status: '完成', completedAt: '2026-03-28 09:00' }
])

const getTaskStatusBadge = (s) => ({ pending: 'default', running: 'processing', paused: 'warning', completed: 'success' }[s] || 'default')
const getTaskStatusText = (s) => ({ pending: '待执行', running: '执行中', paused: '已暂停', completed: '已完成' }[s] || s)
const getPriorityColor = (p) => ({ high: 'red', medium: 'orange', low: 'blue' }[p] || 'gray')

const handleSelectTask = (t) => { selectedTask.value = t }
const handleCreateTask = () => { }
const handleStart = () => { }
const handlePause = () => { }
const handleResume = () => { }
const handleCancel = () => { }
</script>

<style scoped>
.container { padding: 16px; }
.selected { background: #e6f7ff; cursor: pointer; }
.mini-map { background: #f0f0f0; border-radius: 4px; display: flex; justify-content: center; padding: 8px; }
</style>
