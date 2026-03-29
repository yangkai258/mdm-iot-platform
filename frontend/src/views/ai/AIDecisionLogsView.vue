<template>
  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>AI 功能</a-breadcrumb-item>
      <a-breadcrumb-item>决策日志</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 搜索表单 -->
    <div class="pro-search-bar">
      <a-form :model="searchForm" layout="inline">
        <a-form-item label="场景">
          <a-select v-model="searchForm.scene" placeholder="选择场景" allow-clear style="width: 160px">
            <a-option value="behavior">行为决策</a-option>
            <a-option value="emotion">情感响应</a-option>
            <a-option value="interaction">交互决策</a-option>
            <a-option value="health">健康监测</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="设备">
          <a-input v-model="searchForm.device_id" placeholder="设备ID" allow-clear style="width: 160px" />
        </a-form-item>
        <a-form-item label="置信度">
          <a-input-number v-model="searchForm.min_confidence" :min="0" :max="100" placeholder="最低" style="width: 80px" />
          <span style="padding: 0 8px">-</span>
          <a-input-number v-model="searchForm.max_confidence" :min="0" :max="100" placeholder="最高" style="width: 80px" />
        </a-form-item>
        <a-form-item label="决策结果">
          <a-select v-model="searchForm.decision" placeholder="选择结果" allow-clear style="width: 120px">
            <a-option value="approved">通过</a-option>
            <a-option value="rejected">拒绝</a-option>
            <a-option value="uncertain">不确定</a-option>
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
        <a-button type="primary" @click="exportLogs">导出日志</a-button>
        <a-button @click="loadData">刷新</a-button>
      </a-space>
    </div>

    <!-- 统计概览 -->
    <a-row :gutter="16" class="stat-row">
      <a-col :span="6">
        <a-statistic title="总决策数" :value="stats.total" :value-from="0" :animation="true" />
      </a-col>
      <a-col :span="6">
        <a-statistic title="平均置信度" :value="stats.avg_confidence" suffix="%" :precision="1" :value-from="0" :animation="true" />
      </a-col>
      <a-col :span="6">
        <a-statistic title="拒绝决策" :value="stats.rejected" :value-from="0" :animation="true">
          <template #prefix><icon-close style="color: #f53f3f; margin-right: 4px" /></template>
        </a-statistic>
      </a-col>
      <a-col :span="6">
        <a-statistic title="不确定决策" :value="stats.uncertain" :value-from="0" :animation="true">
          <template #prefix><icon-question style="color: #ff7d00; margin-right: 4px" /></template>
        </a-statistic>
      </a-col>
    </a-row>

    <!-- 决策日志表格 -->
    <div class="pro-content-area">
      <a-table
        :columns="columns"
        :data="dataList"
        :loading="loading"
        :pagination="pagination"
        :scroll="{ x: 1400 }"
        @change="handleTableChange"
        row-key="id"
      >
        <template #scene="{ record }">
          <a-tag :color="getSceneColor(record.scene)">{{ getSceneText(record.scene) }}</a-tag>
        </template>
        <template #decision="{ record }">
          <a-tag :color="getDecisionColor(record.decision)">{{ getDecisionText(record.decision) }}</a-tag>
        </template>
        <template #confidence="{ record }">
          <a-progress :percent="record.confidence * 100" size="small" :show-text="true"
            :color="record.confidence > 0.8 ? '#00b42a' : record.confidence > 0.5 ? '#ff7d00' : '#f53f3f'" />
        </template>
        <template #latency_ms="{ record }">
          {{ record.latency_ms }}ms
        </template>
        <template #actions="{ record }">
          <a-button type="text" size="small" @click="viewDetail(record)">详情</a-button>
        </template>
      </a-table>
    </div>

    <!-- 详情弹窗 -->
    <a-modal v-model:visible="detailVisible" title="决策详情" :width="720" footer="null">
      <a-descriptions :column="2" bordered size="large">
        <a-descriptions-item label="决策ID">{{ currentRecord?.id }}</a-descriptions-item>
        <a-descriptions-item label="场景">{{ getSceneText(currentRecord?.scene) }}</a-descriptions-item>
        <a-descriptions-item label="设备">{{ currentRecord?.device_id }}</a-descriptions-item>
        <a-descriptions-item label="决策结果">
          <a-tag :color="getDecisionColor(currentRecord?.decision)">{{ getDecisionText(currentRecord?.decision) }}</a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="置信度">{{ ((currentRecord?.confidence || 0) * 100).toFixed(1) }}%</a-descriptions-item>
        <a-descriptions-item label="延迟">{{ currentRecord?.latency_ms }}ms</a-descriptions-item>
        <a-descriptions-item label="触发时间">{{ currentRecord?.created_at }}</a-descriptions-item>
        <a-descriptions-item label="模型版本">{{ currentRecord?.model_version }}</a-descriptions-item>
        <a-descriptions-item label="输入数据" :span="2">
          <pre style="margin: 0; white-space: pre-wrap; font-size: 12px">{{ currentRecord?.input_data }}</pre>
        </a-descriptions-item>
        <a-descriptions-item label="决策输出" :span="2">
          <pre style="margin: 0; white-space: pre-wrap; font-size: 12px">{{ currentRecord?.output_data }}</pre>
        </a-descriptions-item>
      </a-descriptions>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const loading = ref(false)
const dataList = ref([])
const detailVisible = ref(false)
const currentRecord = ref(null)

const searchForm = reactive({
  scene: '',
  device_id: '',
  min_confidence: null,
  max_confidence: null,
  decision: '',
  dateRange: []
})

const stats = reactive({
  total: 0,
  avg_confidence: 0,
  rejected: 0,
  uncertain: 0
})

const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0
})

const columns = [
  { title: '时间', dataIndex: 'created_at', width: 170 },
  { title: '场景', dataIndex: 'scene', width: 120, slotName: 'scene' },
  { title: '设备ID', dataIndex: 'device_id', width: 140 },
  { title: '决策结果', dataIndex: 'decision', width: 110, slotName: 'decision' },
  { title: '置信度', dataIndex: 'confidence', width: 160, slotName: 'confidence' },
  { title: '延迟', dataIndex: 'latency_ms', width: 90 },
  { title: '模型版本', dataIndex: 'model_version', width: 100 },
  { title: '操作', width: 80, slotName: 'actions', fixed: 'right' }
]

const getSceneColor = (s) => ({ behavior: 'arcoblue', emotion: 'purple', interaction: 'green', health: 'orange' }[s] || 'gray')
const getSceneText = (s) => ({ behavior: '行为决策', emotion: '情感响应', interaction: '交互决策', health: '健康监测' }[s] || s)
const getDecisionColor = (d) => ({ approved: 'green', rejected: 'red', uncertain: 'orange' }[d] || 'gray')
const getDecisionText = (d) => ({ approved: '通过', rejected: '拒绝', uncertain: '不确定' }[d] || d)

const handleSearch = () => {
  pagination.current = 1
  loadData()
}

const handleReset = () => {
  searchForm.scene = ''
  searchForm.device_id = ''
  searchForm.min_confidence = null
  searchForm.max_confidence = null
  searchForm.decision = ''
  searchForm.dateRange = []
  pagination.current = 1
  loadData()
}

const handleTableChange = (pag) => {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  loadData()
}

const exportLogs = () => {
  Message.success('日志导出中...')
}

const viewDetail = (record) => {
  currentRecord.value = record
  detailVisible.value = true
}

const loadData = async () => {
  loading.value = true
  try {
    const scenes = ['behavior', 'emotion', 'interaction', 'health']
    const decisions = ['approved', 'rejected', 'uncertain']
    dataList.value = Array.from({ length: 20 }, (_, i) => ({
      id: `dec_${Date.now()}_${i}`,
      created_at: new Date(Date.now() - i * 600000).toLocaleString('zh-CN'),
      scene: scenes[i % 4],
      device_id: `device_${String(i + 1).padStart(4, '0')}`,
      decision: decisions[i % 3],
      confidence: 0.3 + Math.random() * 0.69,
      latency_ms: Math.floor(50 + Math.random() * 500),
      model_version: `v${['2.1.0', '2.0.5', '2.0.0'][i % 3]}`,
      input_data: JSON.stringify({ context: `scene_${scenes[i % 4]}`, pet_id: i % 2 === 0 ? '小白' : '阿福', timestamp: Date.now() }, null, 2),
      output_data: JSON.stringify({ action: decisions[i % 3] === 'approved' ? 'execute' : 'reject', target: `action_${i}` }, null, 2)
    }))
    pagination.total = 3892
    stats.total = 3892
    stats.avg_confidence = 76.4
    stats.rejected = 128
    stats.uncertain = 245
  } catch (e) {
    Message.error('加载失败')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped lang="less">
.stat-row { margin-bottom: 16px; }
</style>
