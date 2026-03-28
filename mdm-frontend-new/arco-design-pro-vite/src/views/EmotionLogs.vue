<template>
  <div class="container">
    <a-card>
      <template #title>
        <a-space><icon-history /> 情绪日志</a-space>
      </template>
      <template #extra>
        <a-button @click="handleExport">
          <template #icon><icon-download /></template>
          导出日志
        </a-button>
      </template>

      <a-row :gutter="16" style="margin-bottom: 16px">
        <a-col :span="4">
          <a-card size="small">
            <a-statistic title="总记录数" :value="stats.total" />
          </a-card>
        </a-col>
        <a-col :span="4">
          <a-card size="small">
            <a-statistic title="今日记录" :value="stats.today" />
          </a-card>
        </a-col>
        <a-col :span="4">
          <a-card size="small">
            <a-statistic title="平均情绪值" :value="stats.avgScore" suffix="分" />
          </a-card>
        </a-col>
      </a-row>

      <a-row :gutter="16">
        <a-col :span="8">
          <a-card title="情绪时间轴" size="small">
            <a-space direction="vertical" fill>
              <a-date-picker v-model="selectedDate" style="width: 100%" />
              <a-timeline>
                <a-timeline-item v-for="log in timelineLogs" :key="log.id" :color="getEmotionColor(log.type)">
                  <a-space>
                    <a-tag :color="getEmotionColor(log.type)">{{ log.type }}</a-tag>
                    <span>{{ log.time }}</span>
                  </a-space>
                  <p class="emotion-desc">{{ log.description }}</p>
                </a-timeline-item>
              </a-timeline>
            </a-space>
          </a-card>
        </a-col>

        <a-col :span="16">
          <a-card title="情绪详情">
            <a-form :model="filterForm" layout="inline">
              <a-form-item label="情绪类型">
                <a-select v-model="filterForm.type" placeholder="全部" allow-clear style="width: 120px">
                  <a-option value="happy">开心</a-option>
                  <a-option value="sad">难过</a-option>
                  <a-option value="angry">生气</a-option>
                  <a-option value="neutral">平静</a-option>
                </a-select>
              </a-form-item>
              <a-form-item label="置信度">
                <a-input-number v-model="filterForm.minConfidence" :min="0" :max="100" style="width: 100px" />
              </a-form-item>
              <a-form-item>
                <a-button type="primary">查询</a-button>
              </a-form-item>
            </a-form>

            <a-table :columns="columns" :data="tableData" style="margin-top: 16px">
              <template #emotion="{ record }">
                <a-tag :color="getEmotionColor(record.type)">{{ record.type }}</a-tag>
              </template>
              <template #confidence="{ record }">
                <a-progress :percent="record.confidence" :show-text="true" size="small" />
              </template>
            </a-table>
          </a-card>
        </a-col>
      </a-row>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'

const stats = reactive({ total: 1256, today: 45, avgScore: 72 })
const selectedDate = ref('')
const filterForm = reactive({ type: '', minConfidence: 0 })

const timelineLogs = ref([
  { id: 1, type: 'happy', time: '10:30:00', description: '与主人互动' },
  { id: 2, type: 'neutral', time: '10:15:00', description: '待机状态' },
  { id: 3, type: 'sad', time: '09:45:00', description: '主人离开' }
])

const columns = [
  { title: '时间', dataIndex: 'time', width: 180 },
  { title: '情绪', slotName: 'emotion', width: 120 },
  { title: '置信度', slotName: 'confidence' },
  { title: '描述', dataIndex: 'description' },
  { title: '触发来源', dataIndex: 'source' }
]
const tableData = ref([
  { time: '2026-03-28 10:30:00', type: 'happy', confidence: 85, description: '与主人互动', source: '语音' },
  { time: '2026-03-28 10:15:00', type: 'neutral', confidence: 72, description: '待机状态', source: '表情' },
  { time: '2026-03-28 09:45:00', type: 'sad', confidence: 65, description: '主人离开', source: '语音' }
])

const getEmotionColor = (type) => ({
  happy: '#67C23A', sad: '#409EFF', angry: '#F56C6C', neutral: '#909399'
}[type] || 'gray')

const handleExport = () => { }
</script>

<style scoped>
.container { padding: 16px; }
.emotion-desc { font-size: 12px; color: #909399; margin: 4px 0 0 0; }
</style>
