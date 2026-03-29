<template>
  <div class="pro-page-container">
    <Breadcrumb :items="['menu.ai', 'menu.ai.emotion']" />
    <!-- 统计卡片 -->
    <a-row :gutter="16" class="stat-row">
      <a-col :span="6">
        <a-statistic title="今日情绪记录" :value="stats.today_count" :value-from="0" :animation="true">
          <template #prefix><icon-heart style="margin-right: 4px" /></template>
        </a-statistic>
      </a-col>
      <a-col :span="6">
        <a-statistic title="平均情绪强度" :value="stats.avg_intensity" suffix="%" :precision="1" :value-from="0" :animation="true" />
      </a-col>
      <a-col :span="6">
        <a-statistic title="负面情绪预警" :value="stats.negative_count" :value-from="0" :animation="true">
          <template #prefix><icon-warning style="margin-right: 4px" /></template>
        </a-statistic>
      </a-col>
      <a-col :span="6">
        <a-statistic title="已配置响应规则" :value="stats.rule_count" :value-from="0" :animation="true">
          <template #prefix><icon-settings style="margin-right: 4px" /></template>
        </a-statistic>
      </a-col>
    </a-row>

    <!-- 搜索表单 -->
    <div class="pro-search-bar">
      <a-form :model="searchForm" layout="inline">
        <a-form-item label="宠物">
          <a-select v-model="searchForm.pet_id" placeholder="选择宠物" allow-clear style="width: 160px">
            <a-option v-for="p in pets" :key="p.id" :value="p.id">{{ p.name }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="情绪类型">
          <a-select v-model="searchForm.emotion_type" placeholder="选择情绪" allow-clear style="width: 160px">
            <a-option value="happy">开心</a-option>
            <a-option value="sad">难过</a-option>
            <a-option value="angry">愤怒</a-option>
            <a-option value="fear">恐惧</a-option>
            <a-option value="surprise">惊讶</a-option>
            <a-option value="neutral">平静</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="强度范围">
          <a-input-number v-model="searchForm.minIntensity" :min="0" :max="100" placeholder="最低" style="width: 80px" />
          <span style="padding: 0 8px">-</span>
          <a-input-number v-model="searchForm.maxIntensity" :min="0" :max="100" placeholder="最高" style="width: 80px" />
        </a-form-item>
        <a-form-item label="响应动作">
          <a-select v-model="searchForm.response_action" placeholder="选择动作" allow-clear style="width: 140px">
            <a-option value="comfort">安慰</a-option>
            <a-option value="reward">奖励</a-option>
            <a-option value="ignore">忽略</a-option>
            <a-option value="alert">告警</a-option>
          </a-select>
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
        <a-button type="primary" @click="showConfigModal = true">配置响应规则</a-button>
        <a-button type="primary" @click="showTrendModal = true">情绪趋势</a-button>
        <a-button @click="loadData">刷新</a-button>
      </a-space>
    </div>

    <!-- 情绪趋势图 -->
    <div class="chart-panel">
      <div ref="trendChartRef" style="height: 260px"></div>
    </div>

    <!-- 情绪日志表格 -->
    <div class="pro-content-area">
      <a-table
        :columns="columns"
        :data="dataList"
        :loading="loading"
        :pagination="pagination"
        :scroll="{ x: 1100 }"
        @change="handleTableChange"
        row-key="id"
      >
        <template #emotion_type="{ record }">
          <a-tag :color="getEmotionColor(record.emotion_type)">
            {{ getEmotionText(record.emotion_type) }}
          </a-tag>
        </template>
        <template #intensity="{ record }">
          <a-progress :percent="record.intensity" size="small" :show-text="true"
            :color="record.intensity > 70 ? '#f53f3f' : record.intensity > 40 ? '#ff7d00' : '#0fc6c2'" />
        </template>
        <template #response_action="{ record }">
          <a-tag :color="getActionColor(record.response_action)">
            {{ getActionText(record.response_action) }}
          </a-tag>
        </template>
        <template #actions="{ record }">
          <a-button type="text" size="small" @click="viewDetail(record)">详情</a-button>
        </template>
      </a-table>
    </div>

    <!-- 详情弹窗 -->
    <a-modal v-model:visible="detailVisible" title="情绪记录详情" :width="640" footer="null">
      <a-descriptions :column="2" bordered size="large">
        <a-descriptions-item label="记录ID">{{ currentRecord?.id }}</a-descriptions-item>
        <a-descriptions-item label="情绪类型">{{ getEmotionText(currentRecord?.emotion_type) }}</a-descriptions-item>
        <a-descriptions-item label="宠物">{{ currentRecord?.pet_name }}</a-descriptions-item>
        <a-descriptions-item label="设备">{{ currentRecord?.device_id }}</a-descriptions-item>
        <a-descriptions-item label="情绪强度">{{ currentRecord?.intensity }}%</a-descriptions-item>
        <a-descriptions-item label="响应动作">{{ getActionText(currentRecord?.response_action) }}</a-descriptions-item>
        <a-descriptions-item label="触发时间">{{ currentRecord?.created_at }}</a-descriptions-item>
        <a-descriptions-item label="上下文">{{ currentRecord?.context }}</a-descriptions-item>
        <a-descriptions-item label="原始数据" :span="2">{{ currentRecord?.raw_data }}</a-descriptions-item>
      </a-descriptions>
    </a-modal>

    <!-- 响应配置弹窗 -->
    <a-modal v-model:visible="showConfigModal" title="情绪响应配置" :width="700">
      <a-form :model="configForm" layout="vertical">
        <a-form-item label="宠物">
          <a-select v-model="configForm.pet_id" placeholder="选择宠物" style="width: 100%">
            <a-option v-for="p in pets" :key="p.id" :value="p.id">{{ p.name }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="响应规则">
          <a-table :columns="configColumns" :data="configForm.rules" :pagination="false" size="small">
            <template #emotion_type="{ record }">
              <a-tag :color="getEmotionColor(record.emotion)">{{ getEmotionText(record.emotion) }}</a-tag>
            </template>
            <template #action="{ record }">
              <a-select v-model="record.action" style="width: 120px">
                <a-option value="comfort">安慰</a-option>
                <a-option value="reward">奖励</a-option>
                <a-option value="ignore">忽略</a-option>
                <a-option value="alert">告警</a-option>
              </a-select>
            </template>
            <template #enabled="{ record }">
              <a-switch v-model="record.enabled" />
            </template>
          </a-table>
        </a-form-item>
      </a-form>
      <template #footer>
        <a-button @click="showConfigModal = false">取消</a-button>
        <a-button type="primary" @click="saveConfig">保存</a-button>
      </template>
    </a-modal>

    <!-- 情绪趋势弹窗 -->
    <a-modal v-model:visible="showTrendModal" title="情绪趋势分析" :width="900">
      <div ref="fullTrendChartRef" style="height: 400px"></div>
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
const showConfigModal = ref(false)
const showTrendModal = ref(false)
const currentRecord = ref(null)
const trendChartRef = ref(null)
const fullTrendChartRef = ref(null)
let trendChart = null
let fullTrendChart = null

const searchForm = reactive({
  pet_id: '',
  emotion_type: '',
  minIntensity: null,
  maxIntensity: null,
  response_action: ''
})

const pets = ref([
  { id: '1', name: '小白' },
  { id: '2', name: '阿福' }
])

const stats = reactive({
  today_count: 0,
  avg_intensity: 0,
  negative_count: 0,
  rule_count: 0
})

const configForm = reactive({
  pet_id: '1',
  rules: [
    { emotion: 'happy', action: 'reward', enabled: true, threshold: 60 },
    { emotion: 'sad', action: 'comfort', enabled: true, threshold: 40 },
    { emotion: 'angry', action: 'comfort', enabled: true, threshold: 50 },
    { emotion: 'fear', action: 'comfort', enabled: true, threshold: 30 },
    { emotion: 'surprise', action: 'ignore', enabled: false, threshold: 70 },
    { emotion: 'neutral', action: 'ignore', enabled: false, threshold: 50 }
  ]
})

const configColumns = [
  { title: '情绪', slotName: 'emotion_type', width: 120 },
  { title: '阈值', dataIndex: 'threshold', width: 80 },
  { title: '动作', slotName: 'action', width: 140 },
  { title: '启用', slotName: 'enabled', width: 80 }
]

const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0
})

const columns = [
  { title: '时间', dataIndex: 'created_at', width: 170 },
  { title: '宠物', dataIndex: 'pet_name', width: 100 },
  { title: '设备ID', dataIndex: 'device_id', width: 140 },
  { title: '情绪类型', dataIndex: 'emotion_type', width: 110, slotName: 'emotion_type' },
  { title: '强度', dataIndex: 'intensity', width: 160, slotName: 'intensity' },
  { title: '响应动作', dataIndex: 'response_action', width: 110, slotName: 'response_action' },
  { title: '操作', width: 80, slotName: 'actions', fixed: 'right' }
]

const getEmotionColor = (type) => ({
  happy: 'green', sad: 'blue', angry: 'red', fear: 'purple', surprise: 'orange', neutral: 'gray'
}[type] || 'gray')

const getEmotionText = (type) => ({
  happy: '开心', sad: '难过', angry: '愤怒', fear: '恐惧', surprise: '惊讶', neutral: '平静'
}[type] || type)

const getActionColor = (a) => ({
  comfort: 'blue', reward: 'green', ignore: 'gray', alert: 'red'
}[a] || 'gray')

const getActionText = (a) => ({
  comfort: '安慰', reward: '奖励', ignore: '忽略', alert: '告警'
}[a] || a)

const handleSearch = () => {
  pagination.current = 1
  loadData()
}

const handleReset = () => {
  searchForm.pet_id = ''
  searchForm.emotion_type = ''
  searchForm.minIntensity = null
  searchForm.maxIntensity = null
  searchForm.response_action = ''
  pagination.current = 1
  loadData()
}

const handleTableChange = (pag) => {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  loadData()
}

const viewDetail = (record) => {
  currentRecord.value = record
  detailVisible.value = true
}

const saveConfig = () => {
  showConfigModal.value = false
  Message.success('配置保存成功')
}

const initCharts = () => {
  trendChart = echarts.init(trendChartRef.value)
  trendChart.setOption({
    tooltip: { trigger: 'axis' },
    legend: { data: ['开心', '难过', '愤怒', '平静'] },
    xAxis: { type: 'category', data: ['06:00', '09:00', '12:00', '15:00', '18:00', '21:00', '24:00'] },
    yAxis: { type: 'value', max: 100 },
    series: [
      { name: '开心', type: 'line', smooth: true, data: [30, 60, 45, 70, 55, 80, 65] },
      { name: '难过', type: 'line', smooth: true, data: [10, 5, 15, 8, 12, 5, 8] },
      { name: '愤怒', type: 'line', smooth: true, data: [5, 3, 8, 5, 10, 4, 5] },
      { name: '平静', type: 'line', smooth: true, data: [55, 32, 32, 17, 23, 11, 22] }
    ]
  })

  fullTrendChart = echarts.init(fullTrendChartRef.value)
  fullTrendChart.setOption({
    tooltip: { trigger: 'axis' },
    legend: { data: ['开心', '难过', '愤怒', '平静', '恐惧'] },
    xAxis: { type: 'category', data: ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'] },
    yAxis: { type: 'value', max: 100 },
    series: [
      { name: '开心', type: 'line', smooth: true, data: [62, 65, 70, 68, 72, 78, 75] },
      { name: '难过', type: 'line', smooth: true, data: [12, 15, 10, 18, 14, 8, 10] },
      { name: '愤怒', type: 'line', smooth: true, data: [5, 8, 6, 12, 9, 5, 7] },
      { name: '平静', type: 'line', smooth: true, data: [21, 12, 14, 2, 5, 9, 8] },
      { name: '恐惧', type: 'line', smooth: true, data: [0, 0, 0, 0, 0, 0, 0] }
    ]
  })
}

const loadData = async () => {
  loading.value = true
  try {
    const types = ['happy', 'sad', 'angry', 'fear', 'surprise', 'neutral']
    const actions = ['comfort', 'reward', 'ignore', 'alert']
    dataList.value = Array.from({ length: 20 }, (_, i) => {
      const intensity = Math.floor(20 + Math.random() * 80)
      return {
        id: `emo_${Date.now()}_${i}`,
        created_at: new Date(Date.now() - i * 1800000).toLocaleString('zh-CN'),
        pet_name: i % 2 === 0 ? '小白' : '阿福',
        device_id: `device_${String(i + 1).padStart(4, '0')}`,
        emotion_type: types[i % 6],
        intensity,
        response_action: actions[i % 4],
        context: `上下文_${i}`,
        raw_data: `{ "audio_features": [], "visual_cues": [] }`
      }
    })
    pagination.total = 234
    stats.today_count = 234
    stats.avg_intensity = 52.3
    stats.negative_count = 18
    stats.rule_count = 6
  } catch (e) {
    Message.error('加载失败')
  } finally {
    loading.value = false
  }
}

const handleResize = () => {
  trendChart?.resize()
  fullTrendChart?.resize()
}

onMounted(() => {
  loadData()
  nextTick(() => initCharts())
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  trendChart?.dispose()
  fullTrendChart?.dispose()
})
</script>

<style scoped lang="less">
.stat-row { margin-bottom: 16px; }
.chart-panel { background: #fff; border-radius: 4px; padding: 16px; margin-bottom: 16px; }
</style>
