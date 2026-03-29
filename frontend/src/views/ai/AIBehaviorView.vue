<template>
  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>AI 功能</a-breadcrumb-item>
      <a-breadcrumb-item>行为分析</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 统计卡片 -->
    <a-row :gutter="16" class="stat-row">
      <a-col :span="6">
        <a-statistic title="今日行为总数" :value="stats.today_total" :value-from="0" :animation="true">
          <template #prefix><icon-history style="margin-right: 4px" /></template>
        </a-statistic>
      </a-col>
      <a-col :span="6">
        <a-statistic title="平均响应延迟" :value="stats.avg_latency" suffix="ms" :value-from="0" :animation="true" />
      </a-col>
      <a-col :span="6">
        <a-statistic title="异常行为数" :value="stats.anomaly_count" :value-from="0" :animation="true">
          <template #prefix><icon-warning style="margin-right: 4px" /></template>
        </a-statistic>
      </a-col>
      <a-col :span="6">
        <a-statistic title="模型调用次数" :value="stats.model_calls" :value-from="0" :animation="true">
          <template #prefix><icon-robot style="margin-right: 4px" /></template>
        </a-statistic>
      </a-col>
    </a-row>

    <!-- 搜索表单 -->
    <div class="pro-search-bar">
      <a-form :model="searchForm" layout="inline">
        <a-form-item label="设备/用户">
          <a-input v-model="searchForm.keyword" placeholder="搜索设备/用户" allow-clear style="width: 200px" />
        </a-form-item>
        <a-form-item label="行为类型">
          <a-select v-model="searchForm.behavior_type" placeholder="选择行为类型" allow-clear style="width: 160px">
            <a-option value="move">移动</a-option>
            <a-option value="emote">表情</a-option>
            <a-option value="vocalize">发声</a-option>
            <a-option value="interact">互动</a-option>
            <a-option value="sleep">睡眠</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="宠物">
          <a-select v-model="searchForm.pet_id" placeholder="选择宠物" allow-clear style="width: 160px">
            <a-option v-for="p in pets" :key="p.id" :value="p.id">{{ p.name }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="结果">
          <a-select v-model="searchForm.result" placeholder="选择结果" allow-clear style="width: 120px">
            <a-option value="success">成功</a-option>
            <a-option value="failed">失败</a-option>
            <a-option value="anomaly">异常</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="时间范围">
          <a-range-picker v-model="searchForm.dateRange" style="width: 260px" />
        </a-form-item>
        <a-form-item>
          <a-space>
            <a-button type="primary" @click="handleSearch">搜索</a-button>
            <a-button @click="handleReset">重置</a-button>
          </a-space>
        </a-form-item>
      </a-form>
    </div>

    <!-- 操作栏 -->
    <div class="pro-action-bar">
      <a-space>
        <a-button type="primary" @click="handleExport">导出报告</a-button>
        <a-button @click="loadData">刷新</a-button>
      </a-space>
    </div>

    <!-- 行为统计图表 -->
    <div class="chart-panel">
      <a-tabs type="rounded" size="large">
        <a-tab-pane title="行为趋势">
          <div ref="trendChartRef" style="height: 260px"></div>
        </a-tab-pane>
        <a-tab-pane title="行为分布">
          <div ref="pieChartRef" style="height: 260px"></div>
        </a-tab-pane>
      </a-tabs>
    </div>

    <!-- 数据表格 -->
    <div class="pro-content-area general-card">
      <a-table
        :columns="columns"
        :data="dataList"
        :loading="loading"
        :pagination="pagination"
        :scroll="{ x: 1200 }"
        @change="handleTableChange"
        row-key="id"
      >
        <template #behavior_type="{ record }">
          <a-tag :color="getBehaviorColor(record.behavior_type)">
            {{ getBehaviorText(record.behavior_type) }}
          </a-tag>
        </template>
        <template #result="{ record }">
          <a-tag :color="getResultColor(record.result)">{{ getResultText(record.result) }}</a-tag>
        </template>
        <template #confidence="{ record }">
          <a-progress :percent="record.confidence * 100" size="small" :show-text="true" />
        </template>
        <template #latency_ms="{ record }">
          <span :style="{ color: record.latency_ms > 500 ? '#f53f3f' : 'inherit' }">{{ record.latency_ms }}ms</span>
        </template>
        <template #actions="{ record }">
          <a-button type="text" size="small" @click="viewDetail(record)">详情</a-button>
        </template>
      </a-table>
    </div>

    <!-- 详情弹窗 -->
    <a-modal v-model:visible="detailVisible" title="行为详情" :width="640" footer="null">
      <a-descriptions :column="2" bordered size="large">
        <a-descriptions-item label="行为ID">{{ currentRecord?.id }}</a-descriptions-item>
        <a-descriptions-item label="行为类型">{{ getBehaviorText(currentRecord?.behavior_type) }}</a-descriptions-item>
        <a-descriptions-item label="宠物">{{ currentRecord?.pet_name }}</a-descriptions-item>
        <a-descriptions-item label="设备">{{ currentRecord?.device_id }}</a-descriptions-item>
        <a-descriptions-item label="结果">{{ getResultText(currentRecord?.result) }}</a-descriptions-item>
        <a-descriptions-item label="置信度">{{ ((currentRecord?.confidence || 0) * 100).toFixed(1) }}%</a-descriptions-item>
        <a-descriptions-item label="延迟">{{ currentRecord?.latency_ms }}ms</a-descriptions-item>
        <a-descriptions-item label="时间">{{ currentRecord?.created_at }}</a-descriptions-item>
        <a-descriptions-item label="输入数据" :span="2">{{ currentRecord?.input_data }}</a-descriptions-item>
        <a-descriptions-item label="输出结果" :span="2">{{ currentRecord?.output_data }}</a-descriptions-item>
      </a-descriptions>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onUnmounted, nextTick } from 'vue'
import { Message } from '@arco-design/web-vue'
import * as echarts from 'echarts'

const loading = ref(false)
const dataList = ref([])
const detailVisible = ref(false)
const currentRecord = ref(null)
const trendChartRef = ref(null)
const pieChartRef = ref(null)
let trendChart = null
let pieChart = null

const searchForm = reactive({
  keyword: '',
  behavior_type: '',
  pet_id: '',
  result: '',
  dateRange: []
})

const pets = ref([
  { id: '1', name: '小白' },
  { id: '2', name: '阿福' }
])

const stats = reactive({
  today_total: 0,
  avg_latency: 0,
  anomaly_count: 0,
  model_calls: 0
})

const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0
})

const columns = [
  { title: '时间', dataIndex: 'created_at', width: 170 },
  { title: '宠物', dataIndex: 'pet_name', width: 100 },
  { title: '设备ID', dataIndex: 'device_id', width: 140 },
  { title: '行为类型', dataIndex: 'behavior_type', width: 110, slotName: 'behavior_type' },
  { title: '结果', dataIndex: 'result', width: 90, slotName: 'result' },
  { title: '置信度', dataIndex: 'confidence', width: 140, slotName: 'confidence' },
  { title: '延迟', dataIndex: 'latency_ms', width: 90, slotName: 'latency_ms' },
  { title: '操作', width: 80, slotName: 'actions', fixed: 'right' }
]

const getBehaviorColor = (type) => {
  const colors = { move: 'arcoblue', emote: 'purple', vocalize: 'orange', interact: 'green', sleep: 'gray' }
  return colors[type] || 'gray'
}

const getBehaviorText = (type) => {
  const texts = { move: '移动', emote: '表情', vocalize: '发声', interact: '互动', sleep: '睡眠' }
  return texts[type] || type
}

const getResultColor = (r) => {
  return r === 'success' ? 'green' : r === 'anomaly' ? 'red' : 'orange'
}

const getResultText = (r) => {
  return { success: '成功', failed: '失败', anomaly: '异常' }[r] || r
}

const handleSearch = () => {
  pagination.current = 1
  loadData()
}

const handleReset = () => {
  searchForm.keyword = ''
  searchForm.behavior_type = ''
  searchForm.pet_id = ''
  searchForm.result = ''
  searchForm.dateRange = []
  pagination.current = 1
  loadData()
}

const handleTableChange = (pag) => {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  loadData()
}

const handleExport = () => {
  Message.success('导出功能开发中')
}

const viewDetail = (record) => {
  currentRecord.value = record
  detailVisible.value = true
}

const initCharts = () => {
  // 趋势图
  trendChart = echarts.init(trendChartRef.value)
  trendChart.setOption({
    tooltip: { trigger: 'axis' },
    legend: { data: ['移动', '表情', '发声', '互动'] },
    xAxis: { type: 'category', data: ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'] },
    yAxis: { type: 'value' },
    series: [
      { name: '移动', type: 'line', smooth: true, data: [120, 132, 101, 134, 90, 230, 210] },
      { name: '表情', type: 'line', smooth: true, data: [80, 92, 71, 94, 60, 150, 130] },
      { name: '发声', type: 'line', smooth: true, data: [50, 72, 41, 64, 30, 90, 80] },
      { name: '互动', type: 'line', smooth: true, data: [30, 42, 21, 44, 20, 60, 50] }
    ]
  })

  // 饼图
  pieChart = echarts.init(pieChartRef.value)
  pieChart.setOption({
    tooltip: { trigger: 'item' },
    legend: { bottom: 0 },
    series: [{
      type: 'pie',
      radius: ['40%', '70%'],
      data: [
        { value: 1048, name: '移动' },
        { value: 735, name: '表情' },
        { value: 580, name: '发声' },
        { value: 484, name: '互动' },
        { value: 300, name: '睡眠' }
      ]
    }]
  })
}

const loadData = async () => {
  loading.value = true
  try {
    // Mock data
    dataList.value = Array.from({ length: 20 }, (_, i) => ({
      id: `beh_${Date.now()}_${i}`,
      created_at: new Date(Date.now() - i * 3600000).toLocaleString('zh-CN'),
      pet_name: i % 2 === 0 ? '小白' : '阿福',
      device_id: `device_${String(i + 1).padStart(4, '0')}`,
      behavior_type: ['move', 'emote', 'vocalize', 'interact', 'sleep'][i % 5],
      result: i % 10 === 0 ? 'anomaly' : i % 5 === 0 ? 'failed' : 'success',
      confidence: 0.7 + Math.random() * 0.29,
      latency_ms: Math.floor(100 + Math.random() * 800),
      input_data: `{ "context": "pet_interaction", "intensity": ${(i % 5 + 1) * 2} }`,
      output_data: `{ "action": "execute", "target": "position_${i % 8}" }`
    }))
    pagination.total = 158
    stats.today_total = 158
    stats.avg_latency = 245
    stats.anomaly_count = 12
    stats.model_calls = 4821
  } catch (e) {
    Message.error('加载失败')
  } finally {
    loading.value = false
  }
}

const handleResize = () => {
  trendChart?.resize()
  pieChart?.resize()
}

onMounted(() => {
  loadData()
  nextTick(() => initCharts())
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  trendChart?.dispose()
  pieChart?.dispose()
})
</script>

<style scoped lang="less">
.stat-row {
  margin-bottom: 16px;
}
.chart-panel {
  background: #fff;
  border-radius: 4px;
  padding: 16px;
  margin-bottom: 16px;
}

.general-card {
  border-radius: 8px;
}
</style>
