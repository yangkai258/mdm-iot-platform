<template>
  <div class="container">
    <a-row :gutter="16" style="margin-bottom: 16px">
      <a-col :span="4">
        <a-card>
          <a-statistic title="心率" :value="vitals.heartRate" suffix="bpm" />
        </a-card>
      </a-col>
      <a-col :span="4">
        <a-card>
          <a-statistic title="呼吸" :value="vitals.breathing" suffix="次/分" />
        </a-card>
      </a-col>
      <a-col :span="4">
        <a-card>
          <a-statistic title="体温" :value="vitals.temperature" suffix="°C" />
        </a-card>
      </a-col>
      <a-col :span="4">
        <a-card>
          <a-statistic title="活动量" :value="vitals.activity" suffix="步" />
        </a-card>
      </a-col>
      <a-col :span="4">
        <a-card>
          <a-statistic title="睡眠质量" :value="vitals.sleepScore" suffix="分" />
        </a-card>
      </a-col>
      <a-col :span="4">
        <a-space direction="vertical">
          <a-button @click="handleFullscreen">
            <template #icon><icon-fullscreen /></template>
            全屏
          </a-button>
          <a-button @click="handleExport">
            <template #icon><icon-download /></template>
            导出数据
          </a-button>
        </a-space>
      </a-col>
    </a-row>

    <a-row :gutter="16">
      <a-col :span="16">
        <a-card title="实时生命体征">
          <a-tabs>
            <a-tab-pane key="heart" title="心率">
              <a-chart :option="heartChart" style="height: 200px" />
            </a-tab-pane>
            <a-tab-pane key="breathing" title="呼吸">
              <a-chart :option="breathingChart" style="height: 200px" />
            </a-tab-pane>
            <a-tab-pane key="temp" title="体温">
              <a-chart :option="tempChart" style="height: 200px" />
            </a-tab-pane>
          </a-tabs>
        </a-card>
      </a-col>
      <a-col :span="8">
        <a-card title="健康雷达">
          <a-chart :option="radarChart" style="height: 250px" />
        </a-card>
      </a-col>
    </a-row>

    <a-card title="历史数据" style="margin-top: 16px">
      <template #extra>
        <a-space>
          <a-radio-group v-model="timeRange" type="button">
            <a-radio value="24h">24小时</a-radio>
            <a-radio value="7d">7天</a-radio>
            <a-radio value="30d">30天</a-radio>
          </a-radio-group>
        </a-space>
      </template>
      <a-table :columns="historyColumns" :data="historyData" :pagination="pagination" />
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'

const vitals = reactive({ heartRate: 85, breathing: 18, temperature: 38.5, activity: 2500, sleepScore: 88 })
const timeRange = ref('24h')
const pagination = reactive({ current: 1, pageSize: 10, total: 50 })

const historyColumns = [
  { title: '时间', dataIndex: 'time' },
  { title: '心率', dataIndex: 'heartRate' },
  { title: '呼吸', dataIndex: 'breathing' },
  { title: '体温', dataIndex: 'temperature' },
  { title: '活动量', dataIndex: 'activity' }
]
const historyData = ref([])

const realtimeData = reactive({ labels: [], heartData: [], breathingData: [] })

const heartChart = reactive({
  tooltip: { trigger: 'axis' },
  xAxis: { type: 'category', data: realtimeData.labels },
  yAxis: { type: 'value', name: 'bpm' },
  series: [{ type: 'line', smooth: true, data: realtimeData.heartData, areaStyle: {} }]
})

const breathingChart = reactive({
  tooltip: { trigger: 'axis' },
  xAxis: { type: 'category', data: realtimeData.labels },
  yAxis: { type: 'value', name: '次/分' },
  series: [{ type: 'line', smooth: true, data: realtimeData.breathingData, areaStyle: {} }]
})

const tempChart = reactive({
  tooltip: { trigger: 'axis' },
  xAxis: { type: 'category', data: realtimeData.labels },
  yAxis: { type: 'value', name: '°C' },
  series: [{ type: 'line', smooth: true, data: [38.2, 38.3, 38.4, 38.5, 38.4, 38.5], areaStyle: {} }]
})

const radarChart = reactive({
  radar: {
    indicator: [
      { name: '心率', max: 150 }, { name: '呼吸', max: 30 },
      { name: '体温', max: 42 }, { name: '活动', max: 10000 },
      { name: '睡眠', max: 100 }
    ]
  },
  series: [{ type: 'radar', data: [{ value: [85, 18, 38.5, 2500, 88], name: '健康评分' }] }]
})

const handleFullscreen = () => { }
const handleExport = () => { }
</script>

<style scoped>
.container { padding: 16px; }
</style>
