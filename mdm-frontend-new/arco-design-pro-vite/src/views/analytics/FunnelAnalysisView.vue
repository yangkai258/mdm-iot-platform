<template>
    <Breadcrumb :items="['Home','Console','']" />


  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>数据分析</a-breadcrumb-item>
      <a-breadcrumb-item>漏斗分析</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 操作栏 -->
    <div class="pro-action-bar">
      <a-space>
        <a-button type="primary" @click="openCreateModal">创建漏斗</a-button>
      </a-space>
    </div>

    <!-- 筛选区 -->
    <div class="pro-filter-bar">
      <a-card class="filter-card">
        <a-space wrap>
          <a-input-search v-model="searchKeyword" placeholder="搜索漏斗名称" style="width: 240px" search-button @search="loadFunnels" />
          <a-select v-model="filterStatus" placeholder="状态" allow-clear style="width: 120px" @change="loadFunnels">
            <a-option value="active">启用</a-option>
            <a-option value="inactive">停用</a-option>
          </a-select>
        </a-space>
      </a-card>
    </div>

    <!-- 漏斗列表 -->
    <div class="pro-content-area">
      <a-table :columns="columns" :data="funnels" :loading="loading" row-key="id" :pagination="{ pageSize: 10 }" @page-change="onPageChange">
        <template #name="{ record }">
          <a-link @click="openFunnelDetail(record)">{{ record.name }}</a-link>
        </template>
      </a-table>
        <template #status="{ record }">
          <a-tag :color="record.status === 'active' ? 'green' : 'gray'">
            {{ record.status === 'active' ? '启用' : '停用' }}
          </a-tag>
        </template>
        <template #conversion_rate="{ record }">
          <span :class="record.conversion_rate < 30 ? 'text-danger' : ''">{{ record.conversion_rate || 0 }}%</span>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="openFunnelDetail(record)">查看</a-button>
            <a-button type="text" size="small" @click="openEditModal(record)">编辑</a-button>
            <a-button type="text" size="small" status="danger" @click="handleDelete(record)">删除</a-button>
          </a-space>
        </template>
      </a-table>
    </div>

    <!-- 漏斗详情面板 -->
    <a-drawer v-model:visible="detailVisible" :title="currentFunnel?.name || '漏斗详情'" :width="800" @before-ok="handleSaveFunnel">
      <div v-if="currentFunnel">
        <!-- 漏斗可视化 -->
        <a-card title="漏斗转化" class="funnel-visualization">
          <div class="funnel-container">
            <div v-for="(step, index) in funnelSteps" :key="index" class="funnel-step">
              <div class="funnel-bar-wrapper">
                <div
                  class="funnel-bar"
                  :style="{ width: step.percent + '%', backgroundColor: getStepColor(index) }"
                >
                  <span class="funnel-bar-label">{{ step.name }}</span>
                </div>
              </div>
              <div class="funnel-meta">
                <span class="funnel-count">{{ step.value }}</span>
                <span class="funnel-rate" v-if="index > 0">↓ {{ step.conversion_rate }}%</span>
              </div>
            </div>
          </div>
        </a-card>

        <!-- 漏斗步骤详情 -->
        <a-card title="步骤详情" class="steps-detail">
          <a-table :columns="stepColumns" :data="funnelSteps" :pagination="false" row-key="name" size="small">
            <template #conversion="{ record, rowIndex }">
              <span v-if="rowIndex === 0">—</span>
              <span v-else>{{ record.conversion_rate }}%</span>
            </template>
      </a-table>
        </a-card>

        <!-- 时间筛选 -->
        <div class="detail-filter">
          <a-space>
            <a-select v-model="detailTimeRange" style="width: 120px" @change="loadFunnelData">
              <a-option value="today">今日</a-option>
              <a-option value="week">近7天</a-option>
              <a-option value="month">近30天</a-option>
            </a-select>
          </a-space>
        </div>
      </div>
    </a-drawer>

    <!-- 创建/编辑漏斗弹窗 -->
    <a-modal v-model:visible="formVisible" :title="isEditing ? '编辑漏斗' : '创建漏斗'" :width="560" @before-ok="handleSaveFunnel" @cancel="formVisible = false">
      <a-form :model="funnelForm" layout="vertical" ref="formRef">
        <a-form-item label="漏斗名称" field="name" required>
          <a-input v-model="funnelForm.name" placeholder="请输入漏斗名称" />
        </a-form-item>
        <a-form-item label="描述" field="description">
          <a-textarea v-model="funnelForm.description" placeholder="请输入描述" :max-length="200" />
        </a-form-item>
        <a-form-item label="状态" field="status">
          <a-switch v-model="funnelForm.status" checked-value="active" unchecked-value="inactive" />
        </a-form-item>
        <a-form-item label="步骤定义" field="steps">
          <div v-for="(step, idx) in funnelForm.steps" :key="idx" class="step-item">
            <a-input v-model="step.name" placeholder="步骤名称" style="flex: 1" />
            <a-input-number v-model="step.value" placeholder="数值" style="width: 120px; margin-left: 8px" />
            <a-button type="text" status="danger" @click="removeStep(idx)" :disabled="funnelForm.steps.length <= 2">删除</a-button>
          </div>
          <a-button type="dashed" @click="addStep" style="margin-top: 8px; width: 100%">+ 添加步骤</a-button>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import * as analytics from '@/api/analytics'

// 状态
const funnels = ref([])
const loading = ref(false)
const searchKeyword = ref('')
const filterStatus = ref('')
const detailVisible = ref(false)
const formVisible = ref(false)
const isEditing = ref(false)
const currentFunnel = ref(null)
const funnelSteps = ref([])
const detailTimeRange = ref('week')

const funnelForm = reactive({
  name: '',
  description: '',
  status: 'active',
  steps: [{ name: '步骤1', value: 1000 }, { name: '步骤2', value: 500 }]
})

const columns = [
  { title: '漏斗名称', slotName: 'name' },
  { title: '描述', dataIndex: 'description', ellipsis: true },
  { title: '状态', slotName: 'status' },
  { title: '总用户数', dataIndex: 'total_users', width: 100 },
  { title: '最终转化率', slotName: 'conversion_rate', width: 110 },
  { title: '创建时间', dataIndex: 'created_at', width: 170 },
  { title: '操作', slotName: 'actions', width: 180 }
]

const stepColumns = [
  { title: '步骤', dataIndex: 'name' },
  { title: '用户数', dataIndex: 'value' },
  { title: '占总用户比例', dataIndex: 'percent', width: 130 },
  { title: '上一步转化率', slotName: 'conversion', width: 130 }
]

async function loadFunnels() {
  loading.value = true
  try {
    const res = await analytics.getFunnelList({
      keyword: searchKeyword.value,
      status: filterStatus.value
    })
    funnels.value = res.data?.list || []
  } catch (e) {
    console.error('loadFunnels error:', e)
  } finally {
    loading.value = false
  }
}

async function loadFunnelData() {
  if (!currentFunnel.value) return
  try {
    const res = await analytics.getFunnelData(currentFunnel.value.id, {
      time_range: detailTimeRange.value
    })
    const data = res.data?.steps || []
    const max = Math.max(...data.map(s => s.value), 1)
    funnelSteps.value = data.map((s, i) => ({
      ...s,
      percent: Math.round((s.value / max) * 100),
      conversion_rate: i === 0 ? 100 : Math.round((s.value / data[i - 1].value) * 100)
    }))
  } catch (e) {
    console.error('loadFunnelData error:', e)
  }
}

function openFunnelDetail(record) {
  currentFunnel.value = record
  detailVisible.value = true
  loadFunnelData()
}

function openCreateModal() {
  isEditing.value = false
  Object.assign(funnelForm, { name: '', description: '', status: 'active', steps: [{ name: '步骤1', value: 1000 }, { name: '步骤2', value: 500 }] })
  formVisible.value = true
}

function openEditModal(record) {
  isEditing.value = true
  currentFunnel.value = record
  Object.assign(funnelForm, {
    name: record.name,
    description: record.description,
    status: record.status,
    steps: record.steps || [{ name: '步骤1', value: 1000 }, { name: '步骤2', value: 500 }]
  })
  formVisible.value = true
}

async function handleSaveFunnel() {
  try {
    if (isEditing.value) {
      await analytics.updateFunnel(currentFunnel.value.id, funnelForm)
    } else {
      await analytics.createFunnel(funnelForm)
    }
    formVisible.value = false
    detailVisible.value = false
    loadFunnels()
  } catch (e) {
    console.error('handleSaveFunnel error:', e)
  }
}

async function handleDelete(record) {
  try {
    await analytics.deleteFunnel(record.id)
    loadFunnels()
  } catch (e) {
    console.error('handleDelete error:', e)
  }
}

function addStep() {
  funnelForm.steps.push({ name: `步骤${funnelForm.steps.length + 1}`, value: 0 })
}

function removeStep(idx) {
  funnelForm.steps.splice(idx, 1)
}

function onPageChange(page) {
  loadFunnels()
}

function getStepColor(index) {
  const colors = ['#1650ff', '#0bc6b0', '#7b61ff', '#ff7a00', '#f53f3f', '#0fbf60']
  return colors[index % colors.length]
}

onMounted(() => {
  loadFunnels()
})
</script>

<style scoped>
.funnel-visualization {
  margin-bottom: 16px;
}
.funnel-container {
  padding: 16px 0;
}
.funnel-step {
  margin-bottom: 12px;
}
.funnel-bar-wrapper {
  display: flex;
  align-items: center;
}
.funnel-bar {
  height: 40px;
  display: flex;
  align-items: center;
  padding: 0 16px;
  transition: width 0.6s ease;
  min-width: 60px;
}
.funnel-bar-label {
  color: #fff;
  font-size: 14px;
  font-weight: 500;
  white-space: nowrap;
}
.funnel-meta {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-top: 4px;
  padding-left: 8px;
}
.funnel-count {
  font-size: 14px;
  color: #333;
  font-weight: 500;
}
.funnel-rate {
  font-size: 12px;
  color: #f53f3f;
}
.steps-detail {
  margin-bottom: 16px;
}
.detail-filter {
  margin-bottom: 8px;
}
.step-item {
  display: flex;
  align-items: center;
  margin-bottom: 8px;
}
.text-danger {
  color: #f53f3f;
}
</style>
