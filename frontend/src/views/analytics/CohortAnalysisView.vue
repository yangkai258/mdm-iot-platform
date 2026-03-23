<template>
  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>数据分析</a-breadcrumb-item>
      <a-breadcrumb-item>群组分析</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 操作栏 -->
    <div class="pro-action-bar">
      <a-space>
        <a-button type="primary" @click="openCreateModal">创建群组</a-button>
        <a-button @click="loadComparison">群组对比</a-button>
      </a-space>
    </div>

    <!-- 筛选区 -->
    <div class="pro-filter-bar">
      <a-card class="filter-card">
        <a-space wrap>
          <a-input-search v-model="searchKeyword" placeholder="搜索群组名称" style="width: 240px" search-button @search="loadCohorts" />
          <a-select v-model="filterType" placeholder="群组类型" allow-clear style="width: 140px" @change="loadCohorts">
            <a-option value="registration">注册群组</a-option>
            <a-option value="behavior">行为群组</a-option>
            <a-option value="custom">自定义</a-option>
          </a-select>
          <a-range-picker v-model="dateRange" style="width: 260px" @change="loadCohorts" />
        </a-space>
      </a-card>
    </div>

    <!-- 群组列表 -->
    <div class="pro-content-area">
      <a-row :gutter="[16, 16]" v-if="cohorts.length">
        <a-col :xs="24" :sm="12" :md="8" v-for="cohort in cohorts" :key="cohort.id">
          <a-card class="cohort-card" @click="openCohortDetail(cohort)">
            <template #title>
              <div class="cohort-title">
                <span>{{ cohort.name }}</span>
                <a-tag :color="getCohortTypeColor(cohort.type)" size="small">{{ getCohortTypeText(cohort.type) }}</a-tag>
              </div>
            </template>
            <template #extra>
              <span class="cohort-size">{{ cohort.size || 0 }} 人</span>
            </template>
            <a-descriptions :column="1" size="small">
              <a-descriptions-item label="创建时间">{{ cohort.created_at }}</a-descriptions-item>
              <a-descriptions-item label="条件描述">{{ cohort.description || '无' }}</a-descriptions-item>
            </a-descriptions>
            <template #actions>
              <a-button type="text" size="small" @click.stop="openCohortDetail(cohort)">详情</a-button>
              <a-button type="text" size="small" @click.stop="openEditModal(cohort)">编辑</a-button>
              <a-button type="text" size="small" status="danger" @click.stop="handleDelete(cohort)">删除</a-button>
            </template>
          </a-card>
        </a-col>
      </a-row>
      <a-empty v-else description="暂无群组" />
    </div>

    <!-- 群组详情抽屉 -->
    <a-drawer v-model:visible="detailVisible" :title="currentCohort?.name || '群组详情'" :width="900">
      <div v-if="currentCohort">
        <!-- 群组基本信息 -->
        <a-card title="群组信息" class="detail-card">
          <a-descriptions :column="2" size="small">
            <a-descriptions-item label="群组名称">{{ currentCohort.name }}</a-descriptions-item>
            <a-descriptions-item label="群组类型">{{ getCohortTypeText(currentCohort.type) }}</a-descriptions-item>
            <a-descriptions-item label="成员数量">{{ currentCohort.size }}</a-descriptions-item>
            <a-descriptions-item label="创建时间">{{ currentCohort.created_at }}</a-descriptions-item>
          </a-descriptions>
        </a-card>

        <!-- 热力图 -->
        <a-card title="留存热力图" class="detail-card">
          <div class="heatmap-filter">
            <a-space>
              <a-select v-model="heatmapMetric" style="width: 120px" @change="loadHeatmap">
                <a-option value="retention">留存率</a-option>
                <a-option value="activity">活跃度</a-option>
                <a-option value="revenue">收入</a-option>
              </a-select>
            </a-space>
          </div>
          <div class="heatmap-container" ref="heatmapRef"></div>
        </a-card>

        <!-- 群组对比 -->
        <a-card title="群组对比" class="detail-card">
          <div class="comparison-chart" ref="comparisonChartRef"></div>
        </a-card>
      </div>
    </a-drawer>

    <!-- 创建/编辑群组弹窗 -->
    <a-modal v-model:visible="formVisible" :title="isEditing ? '编辑群组' : '创建群组'" :width="560" @before-ok="handleSaveCohort" @cancel="formVisible = false">
      <a-form :model="cohortForm" layout="vertical" ref="formRef">
        <a-form-item label="群组名称" field="name" required>
          <a-input v-model="cohortForm.name" placeholder="请输入群组名称" />
        </a-form-item>
        <a-form-item label="群组类型" field="type" required>
          <a-select v-model="cohortForm.type" placeholder="请选择群组类型">
            <a-option value="registration">注册群组</a-option>
            <a-option value="behavior">行为群组</a-option>
            <a-option value="custom">自定义</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="描述" field="description">
          <a-textarea v-model="cohortForm.description" placeholder="请输入群组描述" :max-length="200" />
        </a-form-item>
        <a-form-item label="筛选条件" field="conditions">
          <a-textarea v-model="cohortForm.conditions" placeholder="JSON 格式筛选条件" :max-length="500" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 群组对比弹窗 -->
    <a-modal v-model:visible="comparisonVisible" title="群组对比" :width="800">
      <div class="comparison-modal-chart" ref="comparisonModalChartRef"></div>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, nextTick } from 'vue'
import * as analytics from '@/api/analytics'
import * as echarts from 'echarts'

// 状态
const cohorts = ref([])
const loading = ref(false)
const searchKeyword = ref('')
const filterType = ref('')
const dateRange = ref([])
const detailVisible = ref(false)
const formVisible = ref(false)
const isEditing = ref(false)
const comparisonVisible = ref(false)
const currentCohort = ref(null)
const heatmapMetric = ref('retention')

const cohortForm = reactive({
  name: '',
  type: 'registration',
  description: '',
  conditions: ''
})

// 图表 ref
const heatmapRef = ref(null)
const comparisonChartRef = ref(null)
const comparisonModalChartRef = ref(null)

async function loadCohorts() {
  loading.value = true
  try {
    const res = await analytics.getCohortList({
      keyword: searchKeyword.value,
      type: filterType.value
    })
    cohorts.value = res.data?.list || []
  } catch (e) {
    console.error('loadCohorts error:', e)
  } finally {
    loading.value = false
  }
}

async function loadHeatmap() {
  if (!currentCohort.value) return
  try {
    const res = await analytics.getCohortHeatmap(currentCohort.value.id, {
      metric: heatmapMetric.value
    })
    const data = res.data || []
    renderHeatmap(data)
  } catch (e) {
    console.error('loadHeatmap error:', e)
  }
}

async function loadComparison() {
  comparisonVisible.value = true
  await nextTick()
  try {
    const res = await analytics.getCohortComparison()
    const data = res.data || []
    renderComparisonModal(data)
  } catch (e) {
    console.error('loadComparison error:', e)
  }
}

function renderHeatmap(data) {
  if (!heatmapRef.value) return
  const chart = echarts.init(heatmapRef.value)
  const days = data.map(d => d.date)
  const cohorts = data[0] ? Object.keys(data[0]).filter(k => k !== 'date') : []
  const heatmapData = []
  days.forEach((day, i) => {
    cohorts.forEach((cohort, j) => {
      heatmapData.push([j, i, data[i][cohort] || 0])
    })
  })
  chart.setOption({
    tooltip: { position: 'top' },
    grid: { top: 40, bottom: 60, left: 100, right: 30 },
    xAxis: { type: 'category', data: days, axisLabel: { rotate: 30 } },
    yAxis: { type: 'category', data: cohorts },
    visualMap: { min: 0, max: 100, calculable: true, orient: 'horizontal', left: 'center', bottom: 10 },
    series: [{ type: 'heatmap', data: heatmapData, label: { show: true, formatter: '{c}%' }, emphasis: { itemStyle: { shadowBlur: 10 } } }]
  })
}

function renderComparisonModal(data) {
  if (!comparisonModalChartRef.value) return
  const chart = echarts.init(comparisonModalChartRef.value)
  const names = data.map(d => d.name)
  const values = data.map(d => d.value)
  chart.setOption({
    tooltip: { trigger: 'axis' },
    xAxis: { type: 'category', data: names },
    yAxis: { type: 'value' },
    series: [{ type: 'bar', data: values }]
  })
}

function openCohortDetail(record) {
  currentCohort.value = record
  detailVisible.value = true
  nextTick(() => loadHeatmap())
}

function openCreateModal() {
  isEditing.value = false
  Object.assign(cohortForm, { name: '', type: 'registration', description: '', conditions: '' })
  formVisible.value = true
}

function openEditModal(record) {
  isEditing.value = true
  currentCohort.value = record
  Object.assign(cohortForm, {
    name: record.name,
    type: record.type,
    description: record.description,
    conditions: record.conditions || ''
  })
  formVisible.value = true
}

async function handleSaveCohort() {
  try {
    if (isEditing.value) {
      await analytics.updateCohort(currentCohort.value.id, cohortForm)
    } else {
      await analytics.createCohort(cohortForm)
    }
    formVisible.value = false
    detailVisible.value = false
    loadCohorts()
  } catch (e) {
    console.error('handleSaveCohort error:', e)
  }
}

async function handleDelete(record) {
  try {
    await analytics.deleteCohort(record.id)
    loadCohorts()
  } catch (e) {
    console.error('handleDelete error:', e)
  }
}

function getCohortTypeColor(type) {
  return { registration: 'arcoblue', behavior: 'green', custom: 'purple' }[type] || 'gray'
}

function getCohortTypeText(type) {
  return { registration: '注册', behavior: '行为', custom: '自定义' }[type] || type
}

onMounted(() => {
  loadCohorts()
})
</script>

<style scoped>
.cohort-card {
  cursor: pointer;
  transition: box-shadow 0.2s;
}
.cohort-card:hover {
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
}
.cohort-title {
  display: flex;
  align-items: center;
  gap: 8px;
}
.cohort-size {
  color: #1650ff;
  font-weight: 500;
}
.detail-card {
  margin-bottom: 16px;
}
.heatmap-container {
  height: 300px;
  width: 100%;
}
.heatmap-filter {
  margin-bottom: 12px;
}
.comparison-chart {
  height: 300px;
  width: 100%;
}
.comparison-modal-chart {
  height: 400px;
  width: 100%;
}
</style>
