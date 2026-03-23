<template>
  <div class="ai-fairness">
    <a-card title="AI 公平性测试">
      <template #extra>
        <a-button type="primary" @click="runTest">运行测试</a-button>
      </template>
      
      <a-tabs>
        <a-tab-pane key="bias" title="偏见检测">
          <a-spin :loading="loading">
            <a-table :columns="biasColumns" :data="biasResults" :pagination="false">
              <template #level="{ record }">
                <a-tag :color="getLevelColor(record.level)">{{ record.level }}</a-tag>
              </template>
              <template #actions="{ record }">
                <a-button size="small" @click="viewDetail(record)">详情</a-button>
              </template>
            </a-table>
          </a-spin>
        </a-tab-pane>
        
        <a-tab-pane key="metrics" title="公平性指标">
          <a-row :gutter="16" style="margin-top: 16px">
            <a-col :span="6">
              <a-statistic title="性别公平指数" :value="metrics.gender_fairness" suffix="%" />
            </a-col>
            <a-col :span="6">
              <a-statistic title="年龄公平指数" :value="metrics.age_fairness" suffix="%" />
            </a-col>
            <a-col :span="6">
              <a-statistic title="地区公平指数" :value="metrics.region_fairness" suffix="%" />
            </a-col>
            <a-col :span="6">
              <a-statistic title="综合评分" :value="metrics.overall_score" />
            </a-col>
          </a-row>
          
          <a-divider>详细指标</a-divider>
          <a-table :columns="metricColumns" :data="detailedMetrics" :pagination="false" />
        </a-tab-pane>
        
        <a-tab-pane key="reports" title="测试报告">
          <a-list :data="reports" :loading="loading">
            <template #item="{ item }">
              <a-list-item>
                <a-list-item-meta :title="item.name" :description="'生成时间: ' + formatTime(item.created_at)" />
                <template #actions>
                  <a-button size="small" type="text" @click="downloadReport(item)">下载</a-button>
                  <a-button size="small" type="text" @click="viewReport(item)">查看</a-button>
                </template>
              </a-list-item>
            </template>
          </a-list>
        </a-tab-pane>
      </a-tabs>
    </a-card>
    
    <a-modal v-model:visible="detailVisible" title="偏见详情" width="700px">
      <a-descriptions :data="detailData" :column="2" />
    </a-modal>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const API_BASE = '/api/v1'
const loading = ref(false)
const biasResults = ref([])
const metrics = ref({ gender_fairness: 0, age_fairness: 0, region_fairness: 0, overall_score: 0 })
const detailedMetrics = ref([])
const reports = ref([])
const detailVisible = ref(false)
const detailData = ref([])

const biasColumns = [
  { title: '检测项', dataIndex: 'check_item' },
  { title: '受影响群体', dataIndex: 'affected_group' },
  { title: '偏差程度', dataIndex: 'bias_score' },
  { title: '风险等级', slotName: 'level' },
  { title: '建议', dataIndex: 'suggestion' },
  { title: '操作', slotName: 'actions', width: 100 },
]

const metricColumns = [
  { title: '指标名称', dataIndex: 'name' },
  { title: '指标值', dataIndex: 'value' },
  { title: '基准值', dataIndex: 'baseline' },
  { title: '偏差', dataIndex: 'deviation' },
  { title: '状态', dataIndex: 'status' },
]

const loadData = async () => {
  loading.value = true
  try {
    const [biasRes, metricsRes, reportsRes] = await Promise.all([
      fetch(`${API_BASE}/ai-fairness/bias-detection`),
      fetch(`${API_BASE}/ai-fairness/metrics`),
      fetch(`${API_BASE}/ai-fairness/reports`)
    ])
    
    const biasData = await biasRes.json()
    const metricsData = await metricsRes.json()
    const reportsData = await reportsRes.json()
    
    biasResults.value = biasData.bias_results || []
    metrics.value = metricsData.metrics || metrics.value
    detailedMetrics.value = metricsData.detailed || []
    reports.value = reportsData.reports || []
  } catch (e) {
    Message.error('加载失败')
  } finally {
    loading.value = false
  }
}

const runTest = async () => {
  loading.value = true
  try {
    const res = await fetch(`${API_BASE}/ai-fairness/run`, { method: 'POST' })
    const data = await res.json()
    Message.success('测试已启动')
    setTimeout(loadData, 2000)
  } catch (e) {
    Message.error('启动失败')
  }
}

const viewDetail = (record) => {
  detailData.value = Object.entries(record).map(([key, value]) => ({ label: key, value }))
  detailVisible.value = true
}

const getLevelColor = (level) => ({ high: 'red', medium: 'orange', low: 'green' }[level] || 'default')

const formatTime = (t) => new Date(t).toLocaleString()

onMounted(loadData)
</script>

<style scoped>
.ai-fairness { padding: 16px; }
</style>
