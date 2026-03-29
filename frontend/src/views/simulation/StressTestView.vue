<template>
  <div class="pro-page-container">
    <Breadcrumb :items="['menu.simulation', 'menu.simulation.stressTest']" />

    <!-- 操作栏 -->
    <div class="pro-action-bar">
      <a-space>
        <a-button type="primary" @click="openCreateDrawer">创建测试</a-button>
        <a-button @click="loadTests">刷新</a-button>
      </a-space>
    </div>

    <!-- 测试列表 -->
    <div class="pro-content-area">
      <a-table :columns="columns" :data="tests" :loading="loading" :pagination="{ pageSize: 10 }" row-key="id" @page-change="onPageChange">
        <template #test_type="{ record }">
          <a-tag :color="getTestTypeColor(record.test_type)">{{ getTestTypeName(record.test_type) }}</a-tag>
        </template>
        <template #status="{ record }">
          <a-tag :color="getStatusColor(record.status)">{{ getStatusName(record.status) }}</a-tag>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button v-if="record.status === 'draft' || record.status === 'completed' || record.status === 'stopped'" type="text" size="small" status="success" @click="handleStart(record)">开始</a-button>
            <a-button v-if="record.status === 'running'" type="text" size="small" status="danger" @click="handleStop(record)">停止</a-button>
            <a-button type="text" size="small" @click="openReportDrawer(record)">报告</a-button>
            <a-button type="text" size="small" @click="openEditDrawer(record)">编辑</a-button>
            <a-button type="text" size="small" status="danger" @click="handleDelete(record)">删除</a-button>
          </a-space>
        </template>
      </a-table>

      <div class="pro-pagination" v-if="total > 0">
        <a-pagination :total="total" :current="page" :page-size="pageSize" show-total @page-size-change="onPageSizeChange" @change="onPageChange" />
      </div>
    </div>

    <!-- 创建/编辑测试抽屉 -->
    <a-drawer v-model:visible="drawerVisible" :title="isEditing ? '编辑压力测试' : '创建压力测试'" :width="600" @before-ok="handleSave">
      <a-form :model="testForm" layout="vertical" ref="formRef">
        <a-form-item label="测试名称" required>
          <a-input v-model="testForm.test_name" placeholder="请输入测试名称" />
        </a-form-item>
        <a-form-item label="测试类型" required>
          <a-select v-model="testForm.test_type" placeholder="请选择测试类型">
            <a-option value="concurrent">并发测试</a-option>
            <a-option value="performance">性能测试</a-option>
            <a-option value="stability">稳定性测试</a-option>
          </a-select>
        </a-form-item>
        <a-divider>目标配置</a-divider>
        <a-form-item label="请求路径" required>
          <a-input v-model="testForm.config.target.endpoint" placeholder="/api/v1/devices/status" />
        </a-form-item>
        <a-form-item label="请求方法">
          <a-select v-model="testForm.config.target.method">
            <a-option value="GET">GET</a-option>
            <a-option value="POST">POST</a-option>
            <a-option value="PUT">PUT</a-option>
            <a-option value="DELETE">DELETE</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="请求体模板">
          <a-textarea v-model="testForm.config.target.body_template_str" placeholder='{"device_id": "{{device_id}}", "status": "online"}' :auto-size="{ minRows: 2, maxRows: 4 }" />
        </a-form-item>
        <a-divider>负载模式</a-divider>
        <a-form-item label="负载类型">
          <a-select v-model="testForm.config.load_pattern.type">
            <a-option value="constant">恒定并发</a-option>
            <a-option value="ramp">逐步加压</a-option>
            <a-option value="spike">峰值测试</a-option>
            <a-option value="soak">长期稳压</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="初始 VU 数">
          <a-input-number v-model="testForm.config.load_pattern.initial_vus" :min="1" :max="1000" />
        </a-form-item>
        <a-form-item label="最大 VU 数">
          <a-input-number v-model="testForm.config.load_pattern.max_vus" :min="1" :max="10000" />
        </a-form-item>
        <a-form-item label="预热时长">
          <a-input v-model="testForm.config.load_pattern.ramp_up_duration" placeholder="2m" />
        </a-form-item>
        <a-form-item label="持续时长">
          <a-input v-model="testForm.config.load_pattern.hold_duration" placeholder="5m" />
        </a-form-item>
        <a-divider>阈值配置</a-divider>
        <a-form-item label="P50 响应时间 (ms)">
          <a-input-number v-model="testForm.config.thresholds.http_req_duration.p50" :min="0" />
        </a-form-item>
        <a-form-item label="P95 响应时间 (ms)">
          <a-input-number v-model="testForm.config.thresholds.http_req_duration.p95" :min="0" />
        </a-form-item>
        <a-form-item label="P99 响应时间 (ms)">
          <a-input-number v-model="testForm.config.thresholds.http_req_duration.p99" :min="0" />
        </a-form-item>
        <a-form-item label="请求失败率">
          <a-input-number v-model="testForm.config.thresholds.http_req_failed.rate" :min="0" :max="1" :step="0.01" />
        </a-form-item>
      </a-form>
    </a-drawer>

    <!-- 报告抽屉 -->
    <a-drawer v-model:visible="reportDrawerVisible" title="测试报告" :width="720">
      <div v-if="currentReport" class="report-content">
        <a-descriptions :column="2" bordered size="small" title="基本信息">
          <a-descriptions-item label="测试名称">{{ currentReport.test_name }}</a-descriptions-item>
          <a-descriptions-item label="测试类型">{{ getTestTypeName(currentReport.test_type) }}</a-descriptions-item>
          <a-descriptions-item label="状态">
            <a-tag :color="getStatusColor(currentReport.status)">{{ getStatusName(currentReport.status) }}</a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="持续时间">{{ currentReport.duration_seconds }}s</a-descriptions-item>
        </a-descriptions>

        <a-divider>汇总统计</a-divider>
        <a-row :gutter="16">
          <a-col :span="6">
            <a-statistic title="总请求数" :value="currentReport.summary?.total_requests || 0" />
          </a-col>
          <a-col :span="6">
            <a-statistic title="失败请求" :value="currentReport.summary?.failed_requests || 0" />
          </a-col>
          <a-col :span="6">
            <a-statistic title="成功率" :value="currentReport.summary?.success_rate || 0" suffix="%" :precision="2" />
          </a-col>
          <a-col :span="6">
            <a-statistic title="QPS" :value="currentReport.summary?.requests_per_second || 0" :precision="1" />
          </a-col>
        </a-row>

        <a-divider>响应时间</a-divider>
        <a-row :gutter="16">
          <a-col :span="8">
            <a-statistic title="平均" :value="currentReport.summary?.avg_response_time_ms || 0" suffix="ms" />
          </a-col>
          <a-col :span="8">
            <a-statistic title="P95" :value="currentReport.summary?.p95_response_time_ms || 0" suffix="ms" />
          </a-col>
          <a-col :span="8">
            <a-statistic title="P99" :value="currentReport.summary?.p99_response_time_ms || 0" suffix="ms" />
          </a-col>
        </a-row>

        <a-divider>阈值结果</a-divider>
        <a-result :title="currentReport.thresholds_passed ? '所有阈值通过' : '部分阈值未通过'" :status="currentReport.thresholds_passed ? 'success' : 'error'">
          <template #subtitle>
            {{ currentReport.thresholds_passed ? '性能指标均在设定阈值范围内' : '部分性能指标超出阈值范围' }}
          </template>
        </a-result>
      </div>
      <a-empty v-else description="暂无报告数据" />
    </a-drawer>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import { getStressTests, getStressTest, createStressTest, updateStressTest, deleteStressTest, startStressTest, stopStressTest, getStressTestReport } from '@/api/simulation'

const tests = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const loading = ref(false)

const columns = [
  { title: '序号', width: 60, render: ({ rowIndex }) => (page.value - 1) * pageSize.value + rowIndex + 1 },
  { title: '测试名称', dataIndex: 'test_name', ellipsis: true },
  { title: '测试类型', dataIndex: 'test_type', slotName: 'test_type', width: 120 },
  { title: '状态', dataIndex: 'status', slotName: 'status', width: 100 },
  { title: '创建时间', dataIndex: 'created_at', width: 180 },
  { title: '操作', slotName: 'actions', width: 220, fixed: 'right' }
]

const drawerVisible = ref(false)
const isEditing = ref(false)
const editingId = ref(null)
const testForm = reactive({
  test_name: '',
  test_type: 'concurrent',
  config: {
    target: { endpoint: '/api/v1/devices/status', method: 'POST', body_template: {} },
    load_pattern: { type: 'ramp', initial_vus: 10, max_vus: 100, ramp_up_duration: '2m', hold_duration: '5m' },
    thresholds: { http_req_duration: { p50: 100, p95: 200, p99: 500 }, http_req_failed: { rate: 0.01 } }
  }
})

const reportDrawerVisible = ref(false)
const currentReport = ref(null)

async function loadTests() {
  loading.value = true
  try {
    const res = await getStressTests({ page: page.value, page_size: pageSize.value })
    tests.value = res.data?.items || res.data || []
    total.value = res.data?.total || 0
  } catch {
    Message.error('加载测试列表失败')
  } finally {
    loading.value = false
  }
}

function openCreateDrawer() {
  isEditing.value = false
  editingId.value = null
  Object.assign(testForm, {
    test_name: '',
    test_type: 'concurrent',
    config: {
      target: { endpoint: '/api/v1/devices/status', method: 'POST', body_template: {} },
      load_pattern: { type: 'ramp', initial_vus: 10, max_vus: 100, ramp_up_duration: '2m', hold_duration: '5m' },
      thresholds: { http_req_duration: { p50: 100, p95: 200, p99: 500 }, http_req_failed: { rate: 0.01 } }
    }
  })
  drawerVisible.value = true
}

function openEditDrawer(test) {
  isEditing.value = true
  editingId.value = test.id
  try {
    const bodyTemplate = test.config?.target?.body_template
    const bodyTemplateStr = typeof bodyTemplate === 'string' ? bodyTemplate : JSON.stringify(bodyTemplate || {}, null, 2)
    Object.assign(testForm, {
      test_name: test.test_name,
      test_type: test.test_type,
      config: {
        target: { ...test.config?.target, body_template: bodyTemplate || {} },
        load_pattern: { ...test.config?.load_pattern },
        thresholds: { ...test.config?.thresholds }
      }
    })
  } catch {
    testForm.test_name = test.test_name
    testForm.test_type = test.test_type
  }
  drawerVisible.value = true
}

async function handleSave() {
  try {
    const data = {
      test_name: testForm.test_name,
      test_type: testForm.test_type,
      config: testForm.config
    }
    if (isEditing.value) {
      await updateStressTest(editingId.value, data)
      Message.success('更新成功')
    } else {
      await createStressTest(data)
      Message.success('创建成功')
    }
    drawerVisible.value = false
    loadTests()
  } catch {
    Message.error('保存失败')
    return false
  }
}

async function handleStart(test) {
  try {
    await startStressTest(test.id)
    Message.success('测试开始')
    loadTests()
  } catch {
    Message.error('启动失败')
  }
}

async function handleStop(test) {
  try {
    await stopStressTest(test.id)
    Message.success('测试已停止')
    loadTests()
  } catch {
    Message.error('停止失败')
  }
}

async function handleDelete(test) {
  try {
    await deleteStressTest(test.id)
    Message.success('删除成功')
    loadTests()
  } catch {
    Message.error('删除失败')
  }
}

async function openReportDrawer(test) {
  try {
    const res = await getStressTestReport(test.id)
    currentReport.value = res.data
    reportDrawerVisible.value = true
  } catch {
    // 如果报告不存在，尝试获取基本信息
    try {
      const res = await getStressTest(test.id)
      currentReport.value = res.data
      reportDrawerVisible.value = true
    } catch {
      Message.error('加载报告失败')
    }
  }
}

function getTestTypeColor(type) {
  const colors = { concurrent: 'arcoblue', performance: 'green', stability: 'orange' }
  return colors[type] || 'gray'
}

function getTestTypeName(type) {
  const names = { concurrent: '并发测试', performance: '性能测试', stability: '稳定性测试' }
  return names[type] || type
}

function getStatusColor(status) {
  const colors = { draft: 'gray', running: 'green', completed: 'arcoblue', stopped: 'orange', failed: 'red' }
  return colors[status] || 'gray'
}

function getStatusName(status) {
  const names = { draft: '草稿', running: '运行中', completed: '已完成', stopped: '已停止', failed: '失败' }
  return names[status] || status
}

function onPageChange(p) {
  page.value = p
  loadTests()
}

function onPageSizeChange(s) {
  pageSize.value = s
  page.value = 1
  loadTests()
}

onMounted(() => {
  loadTests()
})
</script>

<style scoped>
.report-content {
  min-height: 400px;
}
</style>
