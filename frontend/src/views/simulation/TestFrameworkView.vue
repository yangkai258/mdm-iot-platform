<template>
  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>仿真测试</a-breadcrumb-item>
      <a-breadcrumb-item>自动化测试框架</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 操作栏 -->
    <div class="pro-action-bar">
      <a-space>
        <a-button type="primary" @click="openCreateModal">
          <template #icon><icon-plus /></template>
          新建测试用例
        </a-button>
        <a-button @click="runSelected" :disabled="!selectedKeys.length">
          <template #icon><icon-play-circle /></template>
          执行选中
        </a-button>
        <a-button @click="loadCases">
          <template #icon><icon-refresh /></template>
          刷新
        </a-button>
      </a-space>
      <a-space>
        <a-tag :color="runStats.running ? 'green' : 'gray'">
          {{ runStats.running ? '测试运行中' : '空闲' }}
        </a-tag>
      </a-space>
    </div>

    <!-- 统计卡片 -->
    <a-row :gutter="16" class="stats-row">
      <a-col :xs="12" :sm="6">
        <a-card class="stat-card">
          <a-statistic title="用例总数" :value="runStats.total" />
        </a-card>
      </a-col>
      <a-col :xs="12" :sm="6">
        <a-card class="stat-card">
          <a-statistic title="通过" :value="runStats.passed" color="green" />
        </a-card>
      </a-col>
      <a-col :xs="12" :sm="6">
        <a-card class="stat-card">
          <a-statistic title="失败" :value="runStats.failed" color="red" />
        </a-card>
      </a-col>
      <a-col :xs="12" :sm="6">
        <a-card class="stat-card">
          <a-statistic title="通过率" :suffix="'%'" :value="runStats.total ? Math.round(runStats.passed / runStats.total * 100) : 0" />
        </a-card>
      </a-col>
    </a-row>

    <!-- 筛选区 -->
    <div class="pro-filter-bar">
      <a-card>
        <a-space wrap>
          <a-input v-model="filterName" placeholder="用例名称" style="width: 160px" allow-clear @change="loadCases" />
          <a-select v-model="filterModule" placeholder="模块" allow-clear style="width: 140px" @change="loadCases">
            <a-option value="device">设备模块</a-option>
            <a-option value="pet">宠物模块</a-option>
            <a-option value="ota">OTA模块</a-option>
            <a-option value="ai">AI模块</a-option>
            <a-option value="api">API模块</a-option>
          </a-select>
          <a-select v-model="filterStatus" placeholder="状态" allow-clear style="width: 120px" @change="loadCases">
            <a-option value="passed">通过</a-option>
            <a-option value="failed">失败</a-option>
            <a-option value="pending">待执行</a-option>
          </a-select>
          <a-select v-model="filterPriority" placeholder="优先级" allow-clear style="width: 120px" @change="loadCases">
            <a-option value="P0">P0</a-option>
            <a-option value="P1">P1</a-option>
            <a-option value="P2">P2</a-option>
            <a-option value="P3">P3</a-option>
          </a-select>
        </a-space>
      </a-card>
    </div>

    <!-- 测试用例列表 -->
    <div class="pro-content-area">
      <a-table
        :columns="columns"
        :data="testCases"
        :loading="loading"
        :pagination="{ pageSize: 10 }"
        row-key="id"
        v-model:selected-keys="selectedKeys"
        :row-selection="{ type: 'checkbox', showCheckedAll: true }"
        @page-change="onPageChange"
      >
        <template #name="{ record }">
          <a-link @click="openDetail(record)">{{ record.name }}</a-link>
        </template>
        <template #priority="{ record }">
          <a-tag :color="getPriorityColor(record.priority)">{{ record.priority }}</a-tag>
        </template>
        <template #status="{ record }">
          <a-tag :color="getStatusColor(record.status)">
            {{ getStatusName(record.status) }}
          </a-tag>
        </template>
        <template #last_run="{ record }">
          <span :class="record.last_run === 'failed' ? 'text-danger' : ''">
            {{ record.last_run_at || '从未' }}
          </span>
        </template>
        <template #actions="{ record }">
          <a-button type="text" size="small" @click="runCase(record)">
            <template #icon><icon-play-circle /></template>
            执行
          </a-button>
          <a-button type="text" size="small" @click="openDetail(record)">
            <template #icon><icon-eye /></template>
            详情
          </a-button>
          <a-button type="text" size="small" @click="editCase(record)">
            <template #icon><icon-edit /></template>
            编辑
          </a-button>
          <a-button type="text" size="small" @click="deleteCase(record)" status="danger">
            <template #icon><icon-delete /></template>
          </a-button>
        </template>
      </a-table>
    </div>

    <!-- 执行历史 -->
    <a-divider>执行历史</a-divider>
    <div class="pro-content-area">
      <a-table
        :columns="historyColumns"
        :data="runHistory"
        :pagination="{ pageSize: 5 }"
        row-key="id"
        size="small"
      >
        <template #status="{ record }">
          <a-tag :color="record.status === 'passed' ? 'green' : 'red'">
            {{ record.status === 'passed' ? '通过' : '失败' }}
          </a-tag>
        </template>
        <template #duration="{ record }">
          {{ record.duration_ms }}ms
        </template>
        <template #actions="{ record }">
          <a-button type="text" size="small" @click="viewReport(record)">
            查看报告
          </a-button>
        </template>
      </a-table>
    </div>

    <!-- 新建/编辑用例弹窗 -->
    <a-modal v-model:visible="caseModalVisible" :title="editingCase ? '编辑测试用例' : '新建测试用例'" :width="720" @before-ok="handleSaveCase">
      <a-form :model="caseForm" layout="vertical" ref="formRef">
        <a-form-item label="用例名称" required field="name">
          <a-input v-model="caseForm.name" placeholder="请输入用例名称" />
        </a-form-item>
        <a-form-item label="所属模块" required field="module">
          <a-select v-model="caseForm.module" placeholder="请选择模块">
            <a-option value="device">设备模块</a-option>
            <a-option value="pet">宠物模块</a-option>
            <a-option value="ota">OTA模块</a-option>
            <a-option value="ai">AI模块</a-option>
            <a-option value="api">API模块</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="优先级" field="priority">
          <a-select v-model="caseForm.priority">
            <a-option value="P0">P0 - 核心</a-option>
            <a-option value="P1">P1 - 高</a-option>
            <a-option value="P2">P2 - 中</a-option>
            <a-option value="P3">P3 - 低</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="前置条件" field="precondition">
          <a-textarea v-model="caseForm.precondition" placeholder="请输入前置条件" :rows="2" />
        </a-form-item>
        <a-form-item label="测试步骤" required field="steps">
          <a-textarea v-model="caseForm.steps" placeholder="请输入测试步骤" :rows="4" />
        </a-form-item>
        <a-form-item label="预期结果" required field="expected">
          <a-textarea v-model="caseForm.expected" placeholder="请输入预期结果" :rows="2" />
        </a-form-item>
        <a-form-item label="超时时间(ms)" field="timeout">
          <a-input-number v-model="caseForm.timeout" :min="1000" :max="300000" :step="1000" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 用例详情弹窗 -->
    <a-modal v-model:visible="detailVisible" title="测试用例详情" :width="640" footer="null">
      <a-descriptions :column="1" size="small" bordered>
        <a-descriptions-item label="用例名称">{{ detailCase?.name }}</a-descriptions-item>
        <a-descriptions-item label="模块">{{ detailCase?.module }}</a-descriptions-item>
        <a-descriptions-item label="优先级">
          <a-tag :color="getPriorityColor(detailCase?.priority)">{{ detailCase?.priority }}</a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="状态">
          <a-tag :color="getStatusColor(detailCase?.status)">{{ getStatusName(detailCase?.status) }}</a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="前置条件">{{ detailCase?.precondition || '无' }}</a-descriptions-item>
        <a-descriptions-item label="测试步骤">{{ detailCase?.steps }}</a-descriptions-item>
        <a-descriptions-item label="预期结果">{{ detailCase?.expected }}</a-descriptions-item>
        <a-descriptions-item label="创建时间">{{ detailCase?.created_at }}</a-descriptions-item>
        <a-descriptions-item label="最后执行">{{ detailCase?.last_run_at || '从未' }}</a-descriptions-item>
      </a-descriptions>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const testCases = ref<any[]>([])
const runHistory = ref<any[]>([])
const loading = ref(false)
const selectedKeys = ref<string[]>([])
const caseModalVisible = ref(false)
const detailVisible = ref(false)
const editingCase = ref<any>(null)
const detailCase = ref<any>(null)
const filterName = ref('')
const filterModule = ref('')
const filterStatus = ref('')
const filterPriority = ref('')

const runStats = ref({ total: 0, passed: 0, failed: 0, running: false })

const caseForm = ref({
  name: '',
  module: 'device',
  priority: 'P2',
  precondition: '',
  steps: '',
  expected: '',
  timeout: 30000,
})

const columns = [
  { title: '用例名称', dataIndex: 'name', slotName: 'name' },
  { title: '模块', dataIndex: 'module' },
  { title: '优先级', dataIndex: 'priority', slotName: 'priority' },
  { title: '状态', dataIndex: 'status', slotName: 'status' },
  { title: '最后执行', dataIndex: 'last_run', slotName: 'last_run' },
  { title: '操作', slotName: 'actions', width: 200 },
]

const historyColumns = [
  { title: '执行ID', dataIndex: 'id' },
  { title: '用例', dataIndex: 'name' },
  { title: '状态', dataIndex: 'status', slotName: 'status' },
  { title: '耗时', dataIndex: 'duration_ms', slotName: 'duration' },
  { title: '执行时间', dataIndex: 'run_at' },
  { title: '操作', slotName: 'actions' },
]

const getPriorityColor = (p: string) => ({ P0: 'red', P1: 'orange', P2: 'blue', P3: 'gray' }[p] || 'gray')
const getStatusColor = (s: string) => ({ passed: 'green', failed: 'red', pending: 'gray' }[s] || 'gray')
const getStatusName = (s: string) => ({ passed: '通过', failed: '失败', pending: '待执行' }[s] || s)

const loadCases = () => {
  loading.value = true
  setTimeout(() => {
    testCases.value = [
      { id: 'TC001', name: '设备注册流程', module: '设备模块', priority: 'P0', status: 'passed', last_run: 'passed', last_run_at: '2026-03-23 10:00', created_at: '2026-03-01', precondition: '设备未注册', steps: '1. 上电\n2. 发送注册请求\n3. 验证响应', expected: '返回设备ID和Token', timeout: 10000 },
      { id: 'TC002', name: 'OTA升级流程', module: 'OTA模块', priority: 'P1', status: 'passed', last_run: 'passed', last_run_at: '2026-03-23 10:05', created_at: '2026-03-01', precondition: '设备在线', steps: '1. 下发升级指令\n2. 等待下载\n3. 验证版本', expected: '版本号更新', timeout: 120000 },
      { id: 'TC003', name: 'AI情感识别', module: 'AI模块', priority: 'P1', status: 'failed', last_run: 'failed', last_run_at: '2026-03-22 15:30', created_at: '2026-03-05', precondition: '宠物存在', steps: '1. 发送图片\n2. 调用AI接口\n3. 验证结果', expected: '返回情感标签', timeout: 5000 },
      { id: 'TC004', name: '心率异常告警', module: 'AI模块', priority: 'P0', status: 'pending', last_run: null, last_run_at: null, created_at: '2026-03-10', precondition: '宠物在线', steps: '1. 模拟心率数据\n2. 验证告警触发', expected: '触发告警通知', timeout: 30000 },
      { id: 'TC005', name: 'API限流测试', module: 'API模块', priority: 'P2', status: 'passed', last_run: 'passed', last_run_at: '2026-03-21 09:00', created_at: '2026-03-08', precondition: '无', steps: '1. 发送100次请求\n2. 统计响应', expected: '限流返回429', timeout: 60000 },
    ]
    runHistory.value = [
      { id: 'RUN001', name: '设备注册流程', status: 'passed', duration_ms: 1250, run_at: '2026-03-23 10:00' },
      { id: 'RUN002', name: 'OTA升级流程', status: 'passed', duration_ms: 45200, run_at: '2026-03-23 10:05' },
      { id: 'RUN003', name: 'AI情感识别', status: 'failed', duration_ms: 5010, run_at: '2026-03-22 15:30' },
      { id: 'RUN004', name: 'API限流测试', status: 'passed', duration_ms: 32000, run_at: '2026-03-21 09:00' },
    ]
    runStats.value = { total: 5, passed: 3, failed: 1, running: false }
    loading.value = false
  }, 500)
}

const openCreateModal = () => {
  editingCase.value = null
  caseForm.value = { name: '', module: 'device', priority: 'P2', precondition: '', steps: '', expected: '', timeout: 30000 }
  caseModalVisible.value = true
}

const editCase = (record: any) => {
  editingCase.value = record
  caseForm.value = { ...record }
  caseModalVisible.value = true
}

const handleSaveCase = () => {
  if (editingCase.value) {
    const idx = testCases.value.findIndex((c) => c.id === editingCase.value.id)
    if (idx >= 0) testCases.value[idx] = { ...editingCase.value, ...caseForm.value }
    Message.success('用例已更新')
  } else {
    const newCase = { ...caseForm.value, id: `TC${String(testCases.value.length + 1).padStart(3, '0')}`, status: 'pending', created_at: new Date().toISOString().slice(0, 10), last_run: null, last_run_at: null }
    testCases.value.push(newCase)
    Message.success('用例已创建')
  }
  caseModalVisible.value = false
}

const openDetail = (record: any) => {
  detailCase.value = record
  detailVisible.value = true
}

const runCase = (record: any) => {
  runStats.value.running = true
  Message.info(`开始执行: ${record.name}`)
  setTimeout(() => {
    const passed = Math.random() > 0.3
    record.status = passed ? 'passed' : 'failed'
    record.last_run = record.status
    record.last_run_at = new Date().toLocaleString()
    runStats.value = { ...runStats.value, running: false, passed: testCases.value.filter((c) => c.status === 'passed').length, failed: testCases.value.filter((c) => c.status === 'failed').length }
    Message.success(`用例执行完成: ${passed ? '通过' : '失败'}`)
  }, 2000)
}

const runSelected = () => {
  if (!selectedKeys.value.length) return
  Message.info(`开始执行 ${selectedKeys.value.length} 个用例`)
  runStats.value.running = true
  setTimeout(() => {
    runStats.value.running = false
    Message.success('批量执行完成')
    loadCases()
  }, 3000)
}

const deleteCase = (record: any) => {
  testCases.value = testCases.value.filter((c) => c.id !== record.id)
  Message.success('用例已删除')
}

const viewReport = (record: any) => {
  Message.info(`查看执行报告: ${record.id}`)
}

const onPageChange = () => {
  loadCases()
}

onMounted(() => {
  loadCases()
})
</script>

<style scoped lang="less">
.pro-page-container {
  padding: 16px;
}
.pro-breadcrumb {
  margin-bottom: 12px;
}
.pro-action-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}
.pro-filter-bar {
  margin-bottom: 12px;
}
.stats-row {
  margin-bottom: 16px;
}
.stat-card {
  text-align: center;
}
.text-danger {
  color: rgb(var(--danger-6));
}
.pro-content-area {
  margin-bottom: 16px;
}
</style>
