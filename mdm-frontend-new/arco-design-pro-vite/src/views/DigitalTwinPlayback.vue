<template>
  <div class="container">
    <a-card>
      <template #title>
        <a-space><icon-history /> 历史回放</a-space>
      </template>

      <a-row :gutter="16">
        <a-col :span="6">
          <a-card title="时间轴" size="small">
            <a-form :model="timeForm" layout="vertical">
              <a-form-item label="选择日期">
                <a-date-picker v-model="timeForm.date" style="width: 100%" />
              </a-form-item>
              <a-form-item label="时间范围">
                <a-space direction="vertical" fill>
                  <a-time-picker v-model="timeForm.startTime" format="HH:mm" />
                  <a-time-picker v-model="timeForm.endTime" format="HH:mm" />
                </a-space>
              </a-form-item>
              <a-form-item>
                <a-button type="primary" long @click="handleLoad">加载数据</a-button>
              </a-form-item>
            </a-form>
          </a-card>

          <a-card title="播放控制" size="small" style="margin-top: 16px">
            <a-space direction="vertical" fill>
              <a-slider v-model="playbackProgress" :show-input="true" />
              <a-space>
                <a-button @click="handlePlay">
                  <template #icon><icon-play /></template>
                </a-button>
                <a-button @click="handlePause">
                  <template #icon><icon-pause /></template>
                </a-button>
                <a-button @click="handleStop">
                  <template #icon><icon-stop /></template>
                </a-button>
              </a-space>
              <a-radio-group v-model="playbackSpeed" type="button">
                <a-radio value="0.5x">0.5x</a-radio>
                <a-radio value="1x">1x</a-radio>
                <a-radio value="2x">2x</a-radio>
              </a-radio-group>
            </a-space>
          </a-card>
        </a-col>

        <a-col :span="18">
          <a-card title="回放视图">
            <a-row :gutter="16">
              <a-col :span="12">
                <a-card size="small" title="时间点标注">
                  <a-timeline>
                    <a-timeline-item v-for="event in events" :key="event.time" :color="getEventColor(event.type)">
                      <a-space>
                        <a-tag :color="getEventColor(event.type)">{{ event.type }}</a-tag>
                        <span>{{ event.time }}</span>
                      </a-space>
                      <p>{{ event.description }}</p>
                    </a-timeline-item>
                  </a-timeline>
                </a-card>
              </a-col>
              <a-col :span="12">
                <a-card size="small" title="关键生命体征">
                  <a-chart :option="vitalChart" style="height: 300px" />
                </a-card>
              </a-col>
            </a-row>
          </a-card>

          <a-space style="margin-top: 16px">
            <a-button @click="handleExport">
              <template #icon><icon-download /></template>
              导出回放
            </a-button>
            <a-button @click="handleJumpPrev">
              <template #icon><icon-left /></template>
              上一事件
            </a-button>
            <a-button @click="handleJumpNext">
              下一事件
              <template #icon><icon-right /></template>
            </a-button>
          </a-space>
        </a-col>
      </a-row>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'

const timeForm = reactive({ date: '', startTime: '', endTime: '' })
const playbackProgress = ref(0)
const playbackSpeed = ref('1x')

const events = ref([
  { time: '10:00', type: '进食', description: '早餐时间', color: 'orange' },
  { time: '11:30', type: '活动', description: '户外活动', color: 'green' },
  { time: '14:00', type: '睡眠', description: '午休', color: 'blue' }
])

const vitalChart = reactive({
  tooltip: { trigger: 'axis' },
  legend: { data: ['心率', '呼吸'] },
  xAxis: { type: 'category', data: ['10:00', '11:00', '12:00', '13:00', '14:00'] },
  yAxis: [{ type: 'value', name: '心率' }, { type: 'value', name: '呼吸' }],
  series: [
    { name: '心率', type: 'line', data: [85, 88, 82, 78, 65] },
    { name: '呼吸', type: 'line', yAxisIndex: 1, data: [18, 20, 19, 16, 12] }
  ]
})

const getEventColor = (type) => ({
  '进食': 'orange', '活动': 'green', '睡眠': 'blue', '异常': 'red'
}[type] || 'gray')

const handleLoad = () => { }
const handlePlay = () => { }
const handlePause = () => { }
const handleStop = () => { }
const handleExport = () => { }
const handleJumpPrev = () => { }
const handleJumpNext = () => { }
</script>

<style scoped>
.container { padding: 16px; }
</style>
