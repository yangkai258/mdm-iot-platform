<template>
  <div class="page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>AI 模型</a-breadcrumb-item>
      <a-breadcrumb-item>公平性测试</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 概览统计 -->
    <a-row :gutter="16" style="margin-top: 16px">
      <a-col :xs="24" :sm="12" :md="6">
        <a-card class="stat-card">
          <a-statistic title="测试任务总数" :value="stats.total_tasks">
            <template #prefix>🧪</template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :xs="24" :sm="12" :md="6">
        <a-card class="stat-card">
          <a-statistic title="通过率" :value="stats.pass_rate" suffix="%" :precision="1" :value-style="{ color: '#0fc6c2' }">
            <template #prefix>✅</template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :xs="24" :sm="12" :md="6">
        <a-card class="stat-card">
          <a-statistic title="模型数量" :value="stats.model_count">
            <template #prefix>🤖</template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :xs="24" :sm="12" :md="6">
        <a-card class="stat-card">
          <a-statistic title="最新测试时间" :value="stats.last_test" :value-style="{ fontSize: '14px' }">
            <template #prefix>🕐</template>
          </a-statistic>
        </a-card>
      </a-col>
    </a-row>

    <!-- 筛选区 -->
    <a-card class="section-card" style="margin-top: 16px">
      <a-space wrap>
        <a-select v-model="filterForm.model_id" placeholder="选择模型" style="width: 180px" allow-clear @change="loadTests">
          <a-option v-for="m in models" :key="m.model_id" :value="m.model_id">{{ m.model_name }}</a-option>
        </a-select>
        <a-select v-model="filterForm.dimension" placeholder="测试维度" style="width: 160px" allow-clear @change="loadTests">
          <a-option value="gender">性别</a-option>
          <a-option value="age">年龄</a-option>
          <a-option value="region">地区</a-option>
          <a-option value="language">语言</a-option>
        </a-select>
        <a-select v-model="filterForm.result" placeholder="测试结果" style="width: 140px" allow-clear @change="loadTests">
          <a-option value="pass">通过</a-option>
          <a-option value="warning">警告</a-option>
          <a-option value="fail">失败</a-option>
        </a-select>
        <a-button @click="loadTests">
          <template #icon><icon-refresh :spin="loading" /></template>
          刷新
        </a-button>
        <a-button type="primary" @click="showCreateTest">
          <template #icon><icon-plus /></template>
          新建测试
        </a-button>
      </a-space>
    </a-card>

    <!-- 公平性指标卡片 -->
    <a-row :gutter="16" style="margin-top: 16px">
      <a-col :xs="24" :md="8" v-for="dim in dimensionCards" :key="dim.key">
        <a-card class="section-card">
          <div class="dim-header">
            <span class="dim-icon">{{ dim.icon }}</span>
            <span class="dim-title">{{ dim.title }}</span>
            <a-tag :color="dim.color" size="small">{{ dim.score }}分</a-tag>
          </div>
          <a-progress :percent="dim.score" :stroke-color="dim.color" style="margin: 12px 0" />
          <div class="dim-desc">{{ dim.desc }}</div>
          <div class="dim-groups">
            <div v-for="group in dim.groups" :key="group.name" class="dim-group-item">
              <span>{{ group.name }}</span>
              <a-progress :percent="group.score" :stroke-color="group.score >= 80 ? '#0fc6c2' : group.score >= 60 ? '#ff7d00' : '#f53f3f'" size="small" style="width: 120px" />
              <span class="group-score">{{ group.score }}%</span>
            </div>
          </div>
        </a-card>
      </a-col>
    </a-row>

    <!-- 测试记录列表 -->
    <a-card class="section-card" style="margin-top: 16px">
      <template #title>测试记录</template>
      <a-spin :loading="loading" tip="加载中...">
        <a-table :data="testRecords" :pagination="{ pageSize: 10, total: testRecords.length, showTotal: true }" :columns="columns" row-key="test_id" stripe>
          <template #result="{ record }">
            <a-tag :color="getResultColor(record.result)">{{ getResultText(record.result) }}</a-tag>
          </template>
          <template #metrics="{ record }">
            <a-space size="small">
              <a-tooltip v-for="m in record.metrics" :key="m.key" :content="`${m.label}: ${m.value}`">
                <a-tag size="small">{{ m.label }} {{ m.value }}</a-tag>
              </a-tooltip>
            </a-space>
          </template>
          <template #operations="{ record }">
            <a-space>
              <a-button type="text" size="small" @click="viewDetail(record)">详情</a-button>
              <a-button type="text" size="small" @click="rerunTest(record)">重新测试</a-button>
            </a-space>
          </template>
        </a-table>
      </a-spin>
    </a-card>

    <!-- 新建测试弹窗 -->
    <a-modal v-model:visible="createVisible" title="新建公平性测试" :width="560" @before-ok="handleCreateTest" @cancel="createVisible = false">
      <a-form :model="testForm" layout="vertical">
        <a-form-item label="模型" field="model_id">
          <a-select v-model="testForm.model_id" placeholder="请选择模型">
            <a-option v-for="m in models" :key="m.model_id" :value="m.model_id">{{ m.model_name }} ({{ m.model_id }})</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="测试维度" field="dimensions">
          <a-checkbox-group v-model="testForm.dimensions">
            <a-space direction="vertical">
              <a-checkbox value="gender">性别公平性</a-checkbox>
              <a-checkbox value="age">年龄公平性</a-checkbox>
              <a-checkbox value="region">地区公平性</a-checkbox>
              <a-checkbox value="language">语言公平性</a-checkbox>
            </a-space>
          </a-checkbox-group>
        </a-form-item>
        <a-form-item label="样本数量" field="sample_size">
          <a-input-number v-model="testForm.sample_size" :min="100" :max="100000" :step="100" style="width: 100%" />
        </a-form-item>
        <a-form-item label="阈值设置" field="threshold">
          <a-input-number v-model="testForm.threshold" :min="0" :max="1" :precision="2" :step="0.05" style="width: 100%" />
          <div style="font-size: 12px; color: #86909c">低于此值将判定为不公平（0-1）</div>
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 测试详情弹窗 -->
    <a-drawer v-model:visible="detailVisible" :title="`测试详情: ${currentTest?.test_id}`" :width="680" unmountOnHide>
      <a-descriptions :column="2" bordered size="large" style="margin-bottom: 16px">
        <a-descriptions-item label="测试ID">{{ currentTest?.test_id }}</a-descriptions-item>
        <a-descriptions-item label="测试结果">
          <a-tag :color="getResultColor(currentTest?.result)">{{ getResultText(currentTest?.result) }}</a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="模型">{{ currentTest?.model_name }}</a-descriptions-item>
        <a-descriptions-item label="测试时间">{{ currentTest?.test_time }}</a-descriptions-item>
        <a-descriptions-item label="样本数量">{{ currentTest?.sample_size }}</a-descriptions-item>
        <a-descriptions-item label="阈值">{{ currentTest?.threshold }}</a-descriptions-item>
      </a-descriptions>

      <a-card title="分组测试结果" style="margin-bottom: 16px">
        <a-table :data="currentTest?.groupResults || []" :pagination="false" :columns="groupColumns" size="small">
          <template #score="{ record }">
            <a-progress :percent="record.score" size="small" :stroke-color="record.score >= 80 ? '#0fc6c2' : record.score >= 60 ? '#ff7d00' : '#f53f3f'" />
          </template>
        </a-table>
      </a-card>

      <a-card title="详细指标">
        <a-descriptions :column="1">
          <a-descriptions-item v-for="m in currentTest?.metrics || []" :key="m.key" :label="m.label">
            {{ m.value }} <a-tag size="small" :color="m.score >= 80 ? 'green' : m.score >= 60 ? 'orange' : 'red'">{{ m.score }}分</a-tag>
          </a-descriptions-item>
        </a-descriptions>
      </a-card>
    </a-drawer>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { Message } from '@arco-design/web-vue'

const loading = ref(false)
const createVisible = ref(false)
const detailVisible = ref(false)
const currentTest = ref(null)

const stats = reactive({ total_tasks: 0, pass_rate: 0, model_count: 0, last_test: '-' })

const filterForm = reactive({ model_id: '', dimension: '', result: '' })

const testForm = reactive({ model_id: '', dimensions: [], sample_size: 1000, threshold: 0.8 })

const models = ref([
  { model_id: 'MODEL-001', model_name: '推荐模型-v3' },
  { model_id: 'MODEL-002', model_name: '风控模型-v2' },
  { model_id: 'MODEL-003', model_name: '情感分析模型-v1' }
])

const dimensionCards = ref([
  { key: 'gender', icon: '⚧', title: '性别公平性', score: 92, color: '#0fc6c2', desc: '男性/女性群体的预测结果差异', groups: [{ name: '男性', score: 94 }, { name: '女性', score: 90 }] },
  { key: 'age', icon: '👶', title: '年龄公平性', score: 78, color: '#ff7d00', desc: '不同年龄段群体的预测结果差异', groups: [{ name: '18-30岁', score: 85 }, { name: '31-50岁', score: 80 }, { name: '50岁以上', score: 68 }] },
  { key: 'region', icon: '🌍', title: '地区公平性', score: 65, color: '#f53f3f', desc: '不同地区群体的预测结果差异', groups: [{ name: '一线城市', score: 82 }, { name: '二线城市', score: 70 }, { name: '其他地区', score: 55 }] }
])

const testRecords = ref([])

const columns = [
  { title: '测试ID', dataIndex: 'test_id', width: 160 },
  { title: '模型', dataIndex: 'model_name', width: 140 },
  { title: '维度', dataIndex: 'dimensions', width: 120 },
  { title: '结果', dataIndex: 'result', width: 80, slotName: 'result' },
  { title: '关键指标', dataIndex: 'metrics', width: 300, slotName: 'metrics' },
  { title: '测试时间', dataIndex: 'test_time', width: 160 },
  { title: '操作', slotName: 'operations', width: 120 }
]

const groupColumns = [
  { title: '分组', dataIndex: 'group_name', width: 120 },
  { title: '样本数', dataIndex: 'sample_count', width: 100 },
  { title: '正例率', dataIndex: 'positive_rate', width: 100 },
  { title: '公平性得分', dataIndex: 'score', slotName: 'score', width: 200 }
]

const getResultColor = (r) => ({ pass: 'green', warning: 'orange', fail: 'red' }[r] || 'gray')
const getResultText = (r) => ({ pass: '通过', warning: '警告', fail: '失败' }[r] || r)

const loadTests = async () => {
  loading.value = true
  await new Promise(r => setTimeout(r, 600))
  testRecords.value = [
    { test_id: 'FT-20260301', model_id: 'MODEL-001', model_name: '推荐模型-v3', dimensions: '性别/年龄', result: 'pass', sample_size: 5000, threshold: 0.8, test_time: '2026-03-22 10:30', metrics: [{ key: 'disparity', label: '差异度', value: '0.05', score: 95 }, { key: 'equalized_odds', label: '均等化', value: '0.08', score: 92 }], groupResults: [{ group_name: '男性', sample_count: 2500, positive_rate: 0.72, score: 94 }, { group_name: '女性', sample_count: 2500, positive_rate: 0.70, score: 90 }] },
    { test_id: 'FT-20260302', model_id: 'MODEL-002', model_name: '风控模型-v2', dimensions: '地区', result: 'warning', sample_size: 8000, threshold: 0.8, test_time: '2026-03-22 14:20', metrics: [{ key: 'disparity', label: '差异度', value: '0.18', score: 72 }, { key: 'demographic_parity', label: '人口公平', value: '0.22', score: 65 }], groupResults: [{ group_name: '一线城市', sample_count: 3000, positive_rate: 0.65, score: 82 }, { group_name: '二线城市', sample_count: 2800, positive_rate: 0.58, score: 70 }, { group_name: '其他地区', sample_count: 2200, positive_rate: 0.45, score: 55 }] },
    { test_id: 'FT-20260303', model_id: 'MODEL-003', model_name: '情感分析模型-v1', dimensions: '年龄', result: 'fail', sample_size: 3000, threshold: 0.8, test_time: '2026-03-21 16:00', metrics: [{ key: 'disparity', label: '差异度', value: '0.35', score: 45 }, { key: 'equalized_odds', label: '均等化', value: '0.40', score: 38 }], groupResults: [{ group_name: '18-30岁', sample_count: 1000, positive_rate: 0.80, score: 85 }, { group_name: '31-50岁', sample_count: 1200, positive_rate: 0.72, score: 80 }, { group_name: '50岁以上', sample_count: 800, positive_rate: 0.50, score: 48 }] },
    { test_id: 'FT-20260304', model_id: 'MODEL-001', model_name: '推荐模型-v3', dimensions: '语言', result: 'pass', sample_size: 4000, threshold: 0.8, test_time: '2026-03-21 09:15', metrics: [{ key: 'disparity', label: '差异度', value: '0.06', score: 93 }, { key: 'calibration', label: '校准度', value: '0.04', score: 96 }], groupResults: [] }
  ]
  stats.total_tasks = testRecords.value.length
  stats.pass_rate = Math.round(testRecords.value.filter(r => r.result === 'pass').length / testRecords.value.length * 100)
  stats.model_count = models.value.length
  stats.last_test = testRecords.value[0]?.test_time || '-'
  loading.value = false
}

const showCreateTest = () => {
  Object.assign(testForm, { model_id: '', dimensions: [], sample_size: 1000, threshold: 0.8 })
  createVisible.value = true
}

const handleCreateTest = async (done) => {
  if (!testForm.model_id || testForm.dimensions.length === 0) {
    Message.error('请填写必填字段')
    done(false)
    return
  }
  await new Promise(r => setTimeout(r, 800))
  Message.success('测试任务已创建，稍后可在测试记录中查看结果')
  createVisible.value = false
  done(true)
}

const viewDetail = (record) => {
  currentTest.value = record
  detailVisible.value = true
}

const rerunTest = async (record) => {
  await new Promise(r => setTimeout(r, 500))
  Message.success('重新测试已触发')
}

loadTests()
</script>

<style scoped>
.page-container { padding: 20px; }
.breadcrumb { margin-bottom: 12px; }
.stat-card { text-align: center; }
.section-card { margin-bottom: 16px; }
.dim-header { display: flex; align-items: center; gap: 8px; margin-bottom: 8px; }
.dim-icon { font-size: 20px; }
.dim-title { font-weight: 600; flex: 1; }
.dim-desc { font-size: 12px; color: #86909c; margin-bottom: 12px; }
.dim-groups { display: flex; flex-direction: column; gap: 8px; }
.dim-group-item { display: flex; align-items: center; gap: 8px; font-size: 13px; }
.dim-group-item span:first-child { width: 80px; color: #4e5969; }
.group-score { width: 40px; text-align: right; color: #86909c; }
</style>
