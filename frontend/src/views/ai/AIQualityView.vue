<template>
  <div class="pro-page-container">
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>AI 功能</a-breadcrumb-item>
      <a-breadcrumb-item>AI 质量报告</a-breadcrumb-item>
    </a-breadcrumb>

    <a-row :gutter="16" class="stat-row">
      <a-col :span="6">
        <a-card hoverable>
          <a-statistic title="模型准确率" :value="metrics.accuracy || 0" suffix="%" :precision="2" :value-from="0" :animation="true" />
          <div class="metric-change"><span class="trend-up">+1.2%</span> 较上周</div>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card hoverable>
          <a-statistic title="召回率" :value="metrics.recall || 0" suffix="%" :precision="2" :value-from="0" :animation="true" />
          <div class="metric-change"><span class="trend-up">+0.8%</span> 较上周</div>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card hoverable>
          <a-statistic title="F1 分数" :value="metrics.f1 || 0" suffix="" :precision="3" :value-from="0" :animation="true" />
          <div class="metric-change"><span class="trend-up">+0.015</span> 较上周</div>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card hoverable>
          <a-statistic title="平均延迟" :value="metrics.avg_latency || 0" suffix="ms" :precision="0" :value-from="0" :animation="true" />
          <div class="metric-change"><span class="trend-down">-5ms</span> 较上周</div>
        </a-card>
      </a-col>
    </a-row>

    <a-row :gutter="16" style="margin-top: 16px">
      <a-col :span="16">
        <a-card title="质量趋势">
          <div ref="trendChartRef" style="height: 300px"></div>
        </a-card>
      </a-col>
      <a-col :span="8">
        <a-card title="模型版本分布">
          <a-table :columns="versionColumns" :data="modelVersions" :pagination="false" size="small" row-key="id">
            <template #accuracy="{ record }"><span style="color: #00b42a">{{ record.accuracy }}%</span></template>
          </a-table>
        </a-card>
      </a-col>
    </a-row>

    <a-card title="质量详情" style="margin-top: 16px">
      <a-table :columns="detailColumns" :data="detailList" :loading="loading" :pagination="pagination" row-key="id" @page-change="onPageChange">
        <template #model_name="{ record }"><a-tag color="arcoblue">{{ record.model_name }}</a-tag></template>
        <template #status="{ record }"><a-tag :color="record.status === 'online' ? 'green' : 'gray'">{{ record.status === 'online' ? '在线' : '离线' }}</a-tag></template>
      </a-table>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import * as echarts from 'echarts'

const loading = ref(false)
const metrics = ref({ accuracy: 94.5, recall: 91.2, f1: 92.8, avg_latency: 45 })
const detailList = ref([])
const modelVersions = ref([
  { id: 1, version: 'v2.1.0', accuracy: 94.5, status: 'online' },
  { id: 2, version: 'v2.0.5', accuracy: 92.1, status: 'online' },
  { id: 3, version: 'v1.9.0', accuracy: 89.3, status: 'offline' }
])
const trendChartRef = ref(null)
let chartInstance = null

const pagination = reactive({ current: 1, pageSize: 20, total: 0 })
const detailColumns = [
  { title: '模型', slotName: 'model_name', width: 140 },
  { title: '版本', dataIndex: 'version', width: 100 },
  { title: '准确率', dataIndex: 'accuracy', width: 100 },
  { title: '召回率', dataIndex: 'recall', width: 100 },
  { title: 'F1分数', dataIndex: 'f1', width: 100 },
  { title: '延迟(ms)', dataIndex: 'latency', width: 100 },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '更新时间', dataIndex: 'updated_at', width: 170 }
]
const versionColumns = [
  { title: '版本', dataIndex: 'version', width: 80 },
  { title: '准确率', slotName: 'accuracy', width: 80 },
  { title: '状态', slotName: 'status', width: 60 }
]

const loadData = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const params = new URLSearchParams({ page: pagination.current, page_size: pagination.pageSize })
    const res = await fetch(`/api/v1/ai/quality/report?${params}`, { headers: { 'Authorization': `Bearer ${token}` } }).then(r => r.json())
    if (res.code === 0) { detailList.value = res.data?.list || []; pagination.total = res.data?.total || 0 }
    else detailList.value = []
  } catch (e) { Message.error('加载失败') } finally { loading.value = false }
}

const initChart = () => {
  if (!trendChartRef.value) return
  chartInstance = echarts.init(trendChartRef.value)
  chartInstance.setOption({
    tooltip: { trigger: 'axis' },
    legend: { data: ['准确率', '召回率', 'F1分数'] },
    xAxis: { type: 'category', data: ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'] },
    yAxis: { type: 'value', min: 80, max: 100 },
    series: [
      { name: '准确率', type: 'line', data: [93, 94, 94.2, 94.5, 94.3, 94.5, 94.5], smooth: true },
      { name: '召回率', type: 'line', data: [90, 90.5, 91, 91.2, 91, 91.2, 91.2], smooth: true },
      { name: 'F1分数', type: 'line', data: [91.5, 92.2, 92.5, 92.8, 92.6, 92.8, 92.8], smooth: true }
    ]
  })
}

const onPageChange = (page) => { pagination.current = page; loadData() }

onMounted(() => { loadData(); initChart() })
onUnmounted(() => { if (chartInstance) chartInstance.dispose() })
</script>

<style scoped>
.pro-page-container { padding: 16px; }
.stat-row { margin-bottom: 16px; }
.metric-change { font-size: 12px; color: #666; margin-top: 4px; }
.trend-up { color: #00b42a; font-weight: 600; }
.trend-down { color: #f53f3f; font-weight: 600; }
</style>
