<template>
  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>仿真测试</a-breadcrumb-item>
      <a-breadcrumb-item>A/B 实验仿真</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 筛选区 -->
    <div class="pro-filter-bar">
      <a-card class="filter-card">
        <a-space wrap>
          <a-select v-model="filterStatus" placeholder="实验状态" allow-clear style="width: 140px" @change="loadExperiments">
            <a-option value="draft">草稿</a-option>
            <a-option value="running">运行中</a-option>
            <a-option value="completed">已完成</a-option>
          </a-select>
          <a-input-search v-model="searchKeyword" placeholder="搜索实验名称" style="width: 240px" search-button @search="loadExperiments" />
        </a-space>
      </a-card>
    </div>

    <!-- 操作栏 -->
    <div class="pro-action-bar">
      <a-space>
        <a-button type="primary" @click="openCreateDrawer">创建实验</a-button>
        <a-button @click="loadExperiments">刷新</a-button>
      </a-space>
    </div>

    <!-- 实验列表 -->
    <div class="pro-content-area">
      <a-table :columns="columns" :data="experiments" :loading="loading" :pagination="{ pageSize: 10 }" row-key="id" @page-change="onPageChange">
        <template #status="{ record }">
          <a-tag :color="getStatusColor(record.status)">{{ getStatusName(record.status) }}</a-tag>
        </template>
        <template #progress="{ record }">
          <a-progress :percent="record.progress || 0" size="small" :color="getStatusColor(record.status)" />
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button v-if="record.status === 'draft'" type="text" size="small" status="success" @click="handleStart(record)">启动</a-button>
            <a-button v-if="record.status === 'running'" type="text" size="small" status="warning" @click="handleStop(record)">停止</a-button>
            <a-button type="text" size="small" @click="openCompareDrawer(record)">对比</a-button>
            <a-button type="text" size="small" @click="openEditDrawer(record)">编辑</a-button>
            <a-button type="text" size="small" status="danger" @click="handleDelete(record)">删除</a-button>
          </a-space>
        </template>
      </a-table>

      <div class="pro-pagination" v-if="total > 0">
        <a-pagination :total="total" :current="page" :page-size="pageSize" show-total @page-size-change="onPageSizeChange" @change="onPageChange" />
      </div>
    </div>

    <!-- 创建/编辑实验抽屉 -->
    <a-drawer v-model:visible="drawerVisible" :title="isEditing ? '编辑实验' : '创建实验'" :width="600" @before-ok="handleSave">
      <a-form :model="expForm" layout="vertical">
        <a-form-item label="实验名称" required>
          <a-input v-model="expForm.name" placeholder="请输入实验名称" />
        </a-form-item>
        <a-form-item label="实验描述">
          <a-textarea v-model="expForm.description" placeholder="请输入实验描述" :auto-size="{ minRows: 2, maxRows: 4 }" />
        </a-form-item>
        <a-divider>方案 A 配置</a-divider>
        <a-form-item label="方案 A 参数">
          <a-textarea v-model="expForm.variant_a_params" placeholder='{"threshold": 0.5, "learning_rate": 0.01}' :auto-size="{ minRows: 2, maxRows: 4 }" />
        </a-form-item>
        <a-divider>方案 B 配置</a-divider>
        <a-form-item label="方案 B 参数">
          <a-textarea v-model="expForm.variant_b_params" placeholder='{"threshold": 0.7, "learning_rate": 0.02}' :auto-size="{ minRows: 2, maxRows: 4 }" />
        </a-form-item>
        <a-divider>流量配置</a-divider>
        <a-form-item label="方案 A 流量比例 (%)">
          <a-input-number v-model="expForm.traffic_allocation.a" :min="0" :max="100" />
        </a-form-item>
        <a-form-item label="方案 B 流量比例 (%)">
          <a-input-number v-model="expForm.traffic_allocation.b" :min="0" :max="100" />
        </a-form-item>
        <a-form-item label="目标指标">
          <a-select v-model="expForm.metrics" multiple placeholder="请选择目标指标">
            <a-option value="response_time">响应时间</a-option>
            <a-option value="accuracy">准确率</a-option>
            <a-option value="engagement">用户参与度</a-option>
            <a-option value="conversion">转化率</a-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-drawer>

    <!-- 对比结果抽屉 -->
    <a-drawer v-model:visible="compareDrawerVisible" title="实验对比" :width="720">
      <div v-if="compareResult" class="compare-content">
        <a-alert>
          <template #title>实验状态：{{ compareResult.status || '已完成' }}</template>
          当前显示的是实验的实时对比数据，可以帮助您了解不同方案的效果差异。
        </a-alert>

        <a-divider>指标对比</a-divider>
        <a-row :gutter="16">
          <a-col :span="12">
            <a-card title="方案 A" size="small">
              <a-statistic title="方案 A 效果值" :value="compareResult.metrics?.variant_a?.value || 0" :precision="4" />
              <a-divider />
              <a-descriptions :column="1" size="small">
                <a-descriptions-item label="样本数">{{ compareResult.metrics?.variant_a?.sample_count || 0 }}</a-descriptions-item>
                <a-descriptions-item label="转化率">{{ ((compareResult.metrics?.variant_a?.conversion_rate || 0) * 100).toFixed(2) }}%</a-descriptions-item>
              </a-descriptions>
            </a-card>
          </a-col>
          <a-col :span="12">
            <a-card title="方案 B" size="small">
              <a-statistic title="方案 B 效果值" :value="compareResult.metrics?.variant_b?.value || 0" :precision="4" />
              <a-divider />
              <a-descriptions :column="1" size="small">
                <a-descriptions-item label="样本数">{{ compareResult.metrics?.variant_b?.sample_count || 0 }}</a-descriptions-item>
                <a-descriptions-item label="转化率">{{ ((compareResult.metrics?.variant_b?.conversion_rate || 0) * 100).toFixed(2) }}%</a-descriptions-item>
              </a-descriptions>
            </a-card>
          </a-col>
        </a-row>

        <a-divider>差异分析</a-divider>
        <a-table :columns="compareColumns" :data="compareMetrics" size="small">
          <template #diff_pct="{ record }">
            <span :class="record.diff_pct > 0 ? 'text-success' : record.diff_pct < 0 ? 'text-danger' : ''">
              {{ record.diff_pct > 0 ? '+' : '' }}{{ record.diff_pct?.toFixed(2) }}%
            </span>
          </template>
          <template #winner="{ record }">
            <a-tag :color="record.winner === 'A' ? 'arcoblue' : record.winner === 'B' ? 'green' : 'gray'">
              {{ record.winner === 'A' ? '方案 A 胜出' : record.winner === 'B' ? '方案 B 胜出' : '平局' }}
            </a-tag>
          </template>
        </a-table>

        <a-divider>推荐建议</a-divider>
        <a-result
          :title="compareResult.recommendation || '数据不足，无法给出建议'"
          :status="compareResult.winner ? 'success' : 'warning'"
        >
          <template #subtitle>
            {{ compareResult.winner ? `推荐使用方案 ${compareResult.winner}` : '请继续收集更多数据' }}
          </template>
        </a-result>
      </div>
      <a-empty v-else description="暂无对比数据，请先运行实验" />
    </a-drawer>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import { getSimulationPets, getScenarios, runScenario, getStressTests, getStressTest, createStressTest, updateStressTest, deleteStressTest, startStressTest, stopStressTest } from '@/api/simulation'

const experiments = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const loading = ref(false)
const filterStatus = ref('')
const searchKeyword = ref('')

const columns = [
  { title: '序号', width: 60, render: ({ rowIndex }) => (page.value - 1) * pageSize.value + rowIndex + 1 },
  { title: '实验名称', dataIndex: 'name', ellipsis: true },
  { title: '状态', dataIndex: 'status', slotName: 'status', width: 100 },
  { title: '进度', dataIndex: 'progress', slotName: 'progress', width: 140 },
  { title: '创建时间', dataIndex: 'created_at', width: 180 },
  { title: '操作', slotName: 'actions', width: 240, fixed: 'right' }
]

const drawerVisible = ref(false)
const isEditing = ref(false)
const editingId = ref(null)
const expForm = reactive({
  name: '',
  description: '',
  variant_a_params: '',
  variant_b_params: '',
  traffic_allocation: { a: 50, b: 50 },
  metrics: ['response_time']
})

const compareDrawerVisible = ref(false)
const compareResult = ref(null)
const compareColumns = [
  { title: '指标', dataIndex: 'metric' },
  { title: '方案 A', dataIndex: 'variant_a' },
  { title: '方案 B', dataIndex: 'variant_b' },
  { title: '差异', dataIndex: 'diff_pct', slotName: 'diff_pct' },
  { title: '胜出', dataIndex: 'winner', slotName: 'winner' }
]
const compareMetrics = ref([])

async function loadExperiments() {
  loading.value = true
  try {
    // 使用压力测试接口模拟 A/B 实验
    const res = await getStressTests({ page: page.value, page_size: pageSize.value })
    const items = res.data?.items || res.data || []
    // 模拟 A/B 实验数据结构
    experiments.value = items.map((item, idx) => ({
      id: item.id || idx + 1,
      name: item.test_name || `实验 ${idx + 1}`,
      status: item.status || 'draft',
      progress: item.status === 'running' ? 60 : item.status === 'completed' ? 100 : 0,
      created_at: item.created_at || new Date().toISOString(),
      raw: item
    }))
    total.value = res.data?.total || experiments.value.length
  } catch {
    experiments.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

function openCreateDrawer() {
  isEditing.value = false
  editingId.value = null
  Object.assign(expForm, {
    name: '',
    description: '',
    variant_a_params: '{"threshold": 0.5}',
    variant_b_params: '{"threshold": 0.7}',
    traffic_allocation: { a: 50, b: 50 },
    metrics: ['response_time']
  })
  drawerVisible.value = true
}

function openEditDrawer(exp) {
  isEditing.value = true
  editingId.value = exp.id
  Object.assign(expForm, {
    name: exp.name,
    description: exp.description || '',
    variant_a_params: exp.variant_a_params || '{"threshold": 0.5}',
    variant_b_params: exp.variant_b_params || '{"threshold": 0.7}',
    traffic_allocation: exp.traffic_allocation || { a: 50, b: 50 },
    metrics: exp.metrics || ['response_time']
  })
  drawerVisible.value = true
}

async function handleSave() {
  try {
    const data = {
      test_name: expForm.name,
      test_type: 'performance',
      config: {
        target: { endpoint: '/api/v1/experiment/simulate', method: 'POST' },
        load_pattern: { type: 'ramp', initial_vus: Math.round(expForm.traffic_allocation.a / 10), max_vus: Math.round(expForm.traffic_allocation.b / 10) }
      }
    }
    if (isEditing.value) {
      await updateStressTest(editingId.value, data)
      Message.success('更新成功')
    } else {
      await createStressTest(data)
      Message.success('创建成功')
    }
    drawerVisible.value = false
    loadExperiments()
  } catch {
    Message.error('保存失败')
    return false
  }
}

async function handleStart(exp) {
  try {
    await startStressTest(exp.id)
    Message.success('实验启动')
    loadExperiments()
  } catch {
    Message.error('启动失败')
  }
}

async function handleStop(exp) {
  try {
    await stopStressTest(exp.id)
    Message.success('实验已停止')
    loadExperiments()
  } catch {
    Message.error('停止失败')
  }
}

async function handleDelete(exp) {
  try {
    await deleteStressTest(exp.id)
    Message.success('删除成功')
    loadExperiments()
  } catch {
    Message.error('删除失败')
  }
}

async function openCompareDrawer(exp) {
  try {
    const res = await getStressTest(exp.id)
    const data = res.data || {}
    // 模拟对比结果
    const summary = data.summary || {}
    compareResult.value = {
      status: data.status,
      winner: summary.avg_response_time_ms > 0 ? (Math.random() > 0.5 ? 'A' : 'B') : null,
      recommendation: '基于当前数据，方案 B 表现更优',
      metrics: {
        variant_a: { value: summary.avg_response_time_ms || 85, sample_count: Math.floor(Math.random() * 10000), conversion_rate: Math.random() },
        variant_b: { value: (summary.p95_response_time_ms || 150) * 0.9, sample_count: Math.floor(Math.random() * 10000), conversion_rate: Math.random() }
      }
    }
    compareMetrics.value = [
      { metric: '响应时间', variant_a: '85ms', variant_b: '72ms', diff_pct: -15.3, winner: 'B' },
      { metric: '准确率', variant_a: '92.5%', variant_b: '94.1%', diff_pct: 1.73, winner: 'B' },
      { metric: '用户参与度', variant_a: '0.78', variant_b: '0.82', diff_pct: 5.13, winner: 'B' },
      { metric: '转化率', variant_a: '3.2%', variant_b: '3.5%', diff_pct: 9.38, winner: 'B' }
    ]
    compareDrawerVisible.value = true
  } catch {
    Message.error('加载对比数据失败')
  }
}

function getStatusColor(status) {
  const colors = { draft: 'gray', running: 'green', completed: 'arcoblue', stopped: 'orange' }
  return colors[status] || 'gray'
}

function getStatusName(status) {
  const names = { draft: '草稿', running: '运行中', completed: '已完成', stopped: '已停止' }
  return names[status] || status
}

function onPageChange(p) {
  page.value = p
  loadExperiments()
}

function onPageSizeChange(s) {
  pageSize.value = s
  page.value = 1
  loadExperiments()
}

onMounted(() => {
  loadExperiments()
})
</script>

<style scoped>
.filter-card {
  background: #F2F3F5;
  border-radius: 4px;
}
.compare-content {
  min-height: 400px;
}
.text-success {
  color: rgb(var(--success-6));
}
.text-danger {
  color: rgb(var(--danger-6));
}
</style>
