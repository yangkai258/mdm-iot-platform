<template>
  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>平台演进</a-breadcrumb-item>
      <a-breadcrumb-item>RTOS 优化</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 搜索栏 -->
    <div class="pro-search-bar">
      <a-space>
        <a-select v-model="selectedDevice" placeholder="选择设备" allow-search style="width: 200px" @change="onDeviceChange">
          <a-option v-for="d in devices" :key="d.id" :value="d.id">{{ d.name }}</a-option>
        </a-select>
        <a-button @click="loadAll"><template #icon><icon-refresh /></template>刷新</a-button>
      </a-space>
    </div>

    <!-- 操作按钮 -->
    <div class="pro-action-bar">
      <a-space>
        <a-button type="primary" @click="showConfigModal = true">
          <template #icon><icon-settings /></template>
          配置参数
        </a-button>
        <a-button @click="loadSuggestions">
          <template #icon><icon-light /></template>
          获取优化建议
        </a-button>
      </a-space>
    </div>

    <!-- 性能仪表板 -->
    <a-row :gutter="16">
      <a-col :span="6">
        <a-statistic title="CPU 使用率" :value="liveMetrics.cpu_percent" :loading="liveLoading" animation>
          <template #suffix>%</template>
          <template #extra>
            <a-progress :percent="liveMetrics.cpu_percent" :color="cpuColor(liveMetrics.cpu_percent)" :show-text="false" size="small" />
          </template>
        </a-statistic>
      </a-col>
      <a-col :span="6">
        <a-statistic title="内存使用率" :value="liveMetrics.mem_percent" :loading="liveLoading" animation>
          <template #suffix>%</template>
          <template #extra>
            <a-progress :percent="liveMetrics.mem_percent" color="arcoblue" :show-text="false" size="small" />
          </template>
        </a-statistic>
      </a-col>
      <a-col :span="6">
        <a-statistic title="任务数" :value="liveMetrics.task_count" :loading="liveLoading" animation>
          <template #extra><span class="text-secondary">运行中</span></template>
        </a-statistic>
      </a-col>
      <a-col :span="6">
        <a-statistic title="上下文切换" :value="liveMetrics.context_switches" :loading="liveLoading" animation>
          <template #extra><span class="text-secondary">次/s</span></template>
        </a-statistic>
      </a-col>
    </a-row>

    <!-- 实时指标图 -->
    <a-row :gutter="16" style="margin-top: 16px">
      <a-col :span="12">
        <a-card title="CPU 使用趋势" size="small">
          <a-table :data="cpuHistory" :loading="historyLoading" :pagination="false" row-key="ts" size="small">
            <template #columns>
              <a-table-column title="时间" data-index="ts" :width="160">
                <template #cell="{ record }">{{ formatTime(record.ts) }}</template>
              </a-table-column>
              <a-table-column title="CPU %" data-index="cpu">
                <template #cell="{ record }">
                  <a-progress :percent="record.cpu" :color="cpuColor(record.cpu)" :show-text="false" size="small" />
                </template>
              </a-table-column>
              <a-table-column title="User %" data-index="user" :width="80" />
              <a-table-column title="System %" data-index="system" :width="90" />
            </template>
          </a-table>
        </a-card>
      </a-col>
      <a-col :span="12">
        <a-card title="内存使用趋势" size="small">
          <a-table :data="memHistory" :loading="historyLoading" :pagination="false" row-key="ts" size="small">
            <template #columns>
              <a-table-column title="时间" data-index="ts" :width="160">
                <template #cell="{ record }">{{ formatTime(record.ts) }}</template>
              </a-table-column>
              <a-table-column title="内存 %" data-index="mem">
                <template #cell="{ record }">
                  <a-progress :percent="record.mem" color="arcoblue" :show-text="false" size="small" />
                </template>
              </a-table-column>
              <a-table-column title="已用" data-index="used_kb" :width="90">
                <template #cell="{ record }">{{ record.used_kb }} KB</template>
              </a-table-column>
              <a-table-column title="总量" data-index="total_kb" :width="90">
                <template #cell="{ record }">{{ record.total_kb }} KB</template>
              </a-table-column>
            </template>
          </a-table>
        </a-card>
      </a-col>
    </a-row>

    <!-- 任务列表 + 优化建议 -->
    <a-row :gutter="16" style="margin-top: 16px">
      <a-col :span="14">
        <a-card title="任务列表" size="small" :header-extra-text="`${tasks.length} 个任务`">
          <a-table :data="tasks" :loading="tasksLoading" :pagination="{ total: taskTotal, current: taskPage, pageSize: 10, showTotal: true }" @page-change="onTaskPageChange" row-key="id" size="small">
            <template #columns>
              <a-table-column title="任务名" data-index="name" :width="140" />
              <a-table-column title="优先级" data-index="priority" :width="80">
                <template #cell="{ record }">
                  <a-tag :color="priorityColor(record.priority)" size="small">{{ record.priority }}</a-tag>
                </template>
              </a-table-column>
              <a-table-column title="状态" data-index="state" :width="80">
                <template #cell="{ record }">
                  <a-badge :status="taskStateBadge(record.state)" :text="taskStateText(record.state)" />
                </template>
              </a-table-column>
              <a-table-column title="CPU %" data-index="cpu_percent" :width="100">
                <template #cell="{ record }">
                  <span :class="cpuClass(record.cpu_percent)">{{ record.cpu_percent }}%</span>
                </template>
              </a-table-column>
              <a-table-column title="栈使用" data-index="stack_usage" :width="100">
                <template #cell="{ record }">
                  <a-progress :percent="record.stack_usage" :color="stackColor(record.stack_usage)" :show-text="false" size="small" />
                </template>
              </a-table-column>
              <a-table-column title="执行时间" data-index="exec_time_ms" :width="90">
                <template #cell="{ record }">{{ record.exec_time_ms }} ms</template>
              </a-table-column>
            </template>
          </a-table>
        </a-card>
      </a-col>
      <a-col :span="10">
        <a-card title="优化建议" size="small" :header-extra-text="`${suggestions.length} 条`">
          <template #extra>
            <a-button size="small" @click="loadSuggestions"><icon-refresh /></a-button>
          </template>
          <a-list v-if="suggestions.length" :data-source="suggestions" size="small">
            <template #item="{ item }">
              <a-list-item>
                <a-list-item-meta :title="item.title" :description="item.description">
                  <template #avatar>
                    <a-badge :color="suggestionColor(item.priority)" />
                  </template>
                </a-list-item-meta>
                <template #actions>
                  <a-button v-if="!item.applied" size="small" type="primary" @click="applySuggestion(item)">应用</a-button>
                  <a-tag v-else color="green" size="small">已应用</a-tag>
                </template>
              </a-list-item>
            </template>
          </a-list>
          <a-empty v-else description="暂无优化建议" />
        </a-card>
      </a-col>
    </a-row>

    <!-- 调度策略 -->
    <a-divider>调度策略</a-divider>
    <div class="pro-action-bar">
      <a-space>
        <a-button type="primary" @click="showPolicyModal = true">
          <template #icon><icon-plus /></template>
          新增策略
        </a-button>
      </a-space>
    </div>
    <a-table :data="policies" :loading="policiesLoading" :pagination="{ total: policyTotal, current: policyPage, pageSize: 10, showTotal: true }" @page-change="onPolicyPageChange" row-key="id" style="margin-top: 12px">
      <template #columns>
        <a-table-column title="策略名称" data-index="name" :width="160" />
        <a-table-column title="调度算法" data-index="algorithm" :width="120">
          <template #cell="{ record }">{{ algorithmText(record.algorithm) }}</template>
        </a-table-column>
        <a-table-column title="时间片 (ms)" data-index="time_slice" :width="110" />
        <a-table-column title="优先级继承" data-index="priority_inherit" :width="100">
          <template #cell="{ record }">
            <a-tag :color="record.priority_inherit ? 'green' : 'gray'">
              {{ record.priority_inherit ? '启用' : '禁用' }}
            </a-tag>
          </template>
        </a-table-column>
        <a-table-column title="状态" data-index="active" :width="80">
          <template #cell="{ record }">
            <a-switch v-model="record.active" @change="togglePolicy(record)" />
          </template>
        </a-table-column>
        <a-table-column title="创建时间" data-index="created_at" :width="160">
          <template #cell="{ record }">{{ formatTime(record.created_at) }}</template>
        </a-table-column>
        <a-table-column title="操作" :width="140" fixed="right">
          <template #cell="{ record }">
            <a-space>
              <a-button size="small" @click="editPolicy(record)">编辑</a-button>
              <a-button size="small" status="danger" @click="deletePolicy(record)">删除</a-button>
            </a-space>
          </template>
        </a-table-column>
      </template>
    </a-table>

    <!-- 配置弹窗 -->
    <a-modal v-model:visible="showConfigModal" title="RTOS 配置" :width="520" @before-ok="handleConfig" @cancel="showConfigModal = false">
      <a-form :model="rtosConfigForm" layout="vertical">
        <a-form-item label="调度器类型">
          <a-select v-model="rtosConfigForm.scheduler_type">
            <a-option value="round_robin">时间片轮转</a-option>
            <a-option value="priority">优先级调度</a-option>
            <a-option value="edf">最早截止时间优先</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="最大优先级">
          <a-input-number v-model="rtosConfigForm.max_priority" :min="1" :max="255" style="width: 100%" />
        </a-form-item>
        <a-form-item label="时间片大小 (ms)">
          <a-input-number v-model="rtosConfigForm.time_slice_ms" :min="1" style="width: 100%" />
        </a-form-item>
        <a-form-item label="栈大小 (KB)">
          <a-input-number v-model="rtosConfigForm.default_stack_kb" :min="4" :step="4" style="width: 100%" />
        </a-form-item>
        <a-form-item label="内存池大小 (KB)">
          <a-input-number v-model="rtosConfigForm.memory_pool_kb" :min="4" :step="4" style="width: 100%" />
        </a-form-item>
        <a-form-item label="优先级继承">
          <a-switch v-model="rtosConfigForm.priority_inherit" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 策略弹窗 -->
    <a-modal v-model:visible="showPolicyModal" :title="editingPolicy ? '编辑策略' : '新增策略'" :width="480" @before-ok="handlePolicy" @cancel="showPolicyModal = false">
      <a-form :model="policyForm" layout="vertical">
        <a-form-item label="策略名称" required>
          <a-input v-model="policyForm.name" placeholder="请输入策略名称" />
        </a-form-item>
        <a-form-item label="调度算法" required>
          <a-select v-model="policyForm.algorithm" placeholder="请选择调度算法">
            <a-option value="round_robin">时间片轮转</a-option>
            <a-option value="priority_preempt">优先级抢占</a-option>
            <a-option value="edf">最早截止时间优先</a-option>
            <a-option value="multilevel">多级队列</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="时间片 (ms)">
          <a-input-number v-model="policyForm.time_slice" :min="1" style="width: 100%" />
        </a-form-item>
        <a-form-item label="优先级继承">
          <a-switch v-model="policyForm.priority_inherit" />
        </a-form-item>
        <a-form-item label="描述">
          <a-textarea v-model="policyForm.description" placeholder="策略描述" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const API_BASE = '/api/v1'

function getToken() {
  return localStorage.getItem('token') || ''
}

async function apiRequest(url, options = {}) {
  const res = await fetch(url, {
    ...options,
    headers: {
      'Authorization': `Bearer ${getToken()}`,
      'Content-Type': 'application/json',
      ...(options.headers || {})
    }
  })
  const data = await res.json()
  if (data.code !== 0 && data.code !== 200) throw new Error(data.message || '请求失败')
  return data
}

const devices = ref([])
const selectedDevice = ref('')
const liveMetrics = ref({ cpu_percent: 0, mem_percent: 0, task_count: 0, context_switches: 0 })
const cpuHistory = ref([])
const memHistory = ref([])
const tasks = ref([])
const taskTotal = ref(0)
const taskPage = ref(1)
const suggestions = ref([])
const policies = ref([])
const policyTotal = ref(0)
const policyPage = ref(1)
const liveLoading = ref(false)
const historyLoading = ref(false)
const tasksLoading = ref(false)
const policiesLoading = ref(false)

// 配置
const showConfigModal = ref(false)
const rtosConfigForm = reactive({
  scheduler_type: 'round_robin', max_priority: 32, time_slice_ms: 10,
  default_stack_kb: 16, memory_pool_kb: 64, priority_inherit: true
})

// 策略
const showPolicyModal = ref(false)
const editingPolicy = ref(null)
const policyForm = reactive({ name: '', algorithm: 'round_robin', time_slice: 10, priority_inherit: true, description: '' })

function formatTime(t) {
  if (!t) return '-'
  return new Date(t).toLocaleString('zh-CN')
}

function cpuColor(p) {
  if (p < 50) return 'green'
  if (p < 80) return 'orange'
  return 'red'
}
function cpuClass(p) {
  if (p < 50) return 'text-success'
  if (p < 80) return 'text-warning'
  return 'text-danger'
}
function stackColor(s) {
  if (s < 60) return 'green'
  if (s < 80) return 'orange'
  return 'red'
}
function priorityColor(p) {
  if (p <= 5) return 'red'
  if (p <= 15) return 'orange'
  return 'green'
}
function taskStateBadge(s) {
  return { running: 'success', ready: 'processing', blocked: 'warning', suspended: 'default' }[s] || 'default'
}
function taskStateText(s) {
  return { running: '运行', ready: '就绪', blocked: '阻塞', suspended: '挂起' }[s] || s
}
function suggestionColor(p) {
  return { high: 'red', medium: 'orange', low: 'green' }[p] || 'gray'
}
function algorithmText(a) {
  return { round_robin: '时间片轮转', priority_preempt: '优先级抢占', edf: 'EDF', multilevel: '多级队列' }[a] || a
}

async function loadDevices() {
  try {
    const res = await apiRequest(`${API_BASE}/devices?page_size=100`)
    devices.value = res.data?.list || []
    if (devices.value.length && !selectedDevice.value) {
      selectedDevice.value = devices.value[0].id
      loadAll()
    }
  } catch (e) { /* ignore */ }
}

async function loadAll() {
  if (!selectedDevice.value) return
  await Promise.all([loadLiveMetrics(), loadHistory(), loadTasks(), loadSuggestions(), loadPolicies()])
}

async function loadLiveMetrics() {
  liveLoading.value = true
  try {
    const res = await apiRequest(`${API_BASE}/rtos/devices/${selectedDevice.value}/live-metrics`)
    liveMetrics.value = res.data || {}
  } catch (e) { /* ignore */ }
  finally { liveLoading.value = false }
}

async function loadHistory() {
  historyLoading.value = true
  try {
    const res = await apiRequest(`${API_BASE}/rtos/devices/${selectedDevice.value}/history?metrics=cpu,mem&range=1h`)
    cpuHistory.value = res.data?.cpu || []
    memHistory.value = res.data?.mem || []
  } catch (e) { /* ignore */ }
  finally { historyLoading.value = false }
}

async function loadTasks() {
  tasksLoading.value = true
  try {
    const params = { page: taskPage.value, page_size: 10 }
    const qs = new URLSearchParams(params).toString()
    const res = await apiRequest(`${API_BASE}/rtos/devices/${selectedDevice.value}/tasks?${qs}`)
    tasks.value = res.data?.list || []
    taskTotal.value = res.data?.total || 0
  } catch (e) { /* ignore */ }
  finally { tasksLoading.value = false }
}

async function loadSuggestions() {
  try {
    const res = await apiRequest(`${API_BASE}/rtos/devices/${selectedDevice.value}/suggestions`)
    suggestions.value = res.data?.list || []
  } catch (e) { /* ignore */ }
}

async function applySuggestion(item) {
  try {
    await apiRequest(`${API_BASE}/rtos/devices/${selectedDevice.value}/apply-suggestion`, {
      method: 'POST',
      body: JSON.stringify({ suggestion_id: item.id })
    })
    Message.success('应用成功')
    item.applied = true
    loadAll()
  } catch (e) {
    Message.error('应用失败: ' + e.message)
  }
}

async function loadPolicies() {
  policiesLoading.value = true
  try {
    const params = { page: policyPage.value, page_size: 10 }
    const qs = new URLSearchParams(params).toString()
    const res = await apiRequest(`${API_BASE}/rtos/scheduler-policies?${qs}`)
    policies.value = res.data?.list || []
    policyTotal.value = res.data?.total || 0
  } catch (e) { /* ignore */ }
  finally { policiesLoading.value = false }
}

async function loadRtosConfig() {
  try {
    const res = await apiRequest(`${API_BASE}/rtos/config/${selectedDevice.value}`)
    Object.assign(rtosConfigForm, res.data || {})
  } catch (e) { /* ignore */ }
}

async function handleConfig(done) {
  try {
    await apiRequest(`${API_BASE}/rtos/config/${selectedDevice.value}`, {
      method: 'PUT',
      body: JSON.stringify(rtosConfigForm)
    })
    Message.success('配置保存成功')
    showConfigModal.value = false
    done(true)
  } catch (e) {
    Message.error('保存失败: ' + e.message)
    done(false)
  }
}

function editPolicy(record) {
  editingPolicy.value = record
  Object.assign(policyForm, {
    id: record.id, name: record.name, algorithm: record.algorithm,
    time_slice: record.time_slice, priority_inherit: record.priority_inherit, description: record.description || ''
  })
  showPolicyModal.value = true
}

async function handlePolicy(done) {
  if (!policyForm.name || !policyForm.algorithm) {
    Message.warning('请填写必填项')
    done(false)
    return
  }
  try {
    const url = editingPolicy.value
      ? `${API_BASE}/rtos/scheduler-policies/${editingPolicy.value.id}`
      : `${API_BASE}/rtos/scheduler-policies`
    const method = editingPolicy.value ? 'PUT' : 'POST'
    await apiRequest(url, { method, body: JSON.stringify(policyForm) })
    Message.success(editingPolicy.value ? '更新成功' : '创建成功')
    showPolicyModal.value = false
    editingPolicy.value = null
    Object.assign(policyForm, { name: '', algorithm: 'round_robin', time_slice: 10, priority_inherit: true, description: '' })
    loadPolicies()
    done(true)
  } catch (e) {
    Message.error('操作失败: ' + e.message)
    done(false)
  }
}

async function togglePolicy(record) {
  try {
    await apiRequest(`${API_BASE}/rtos/scheduler-policies/${record.id}`, {
      method: 'PUT',
      body: JSON.stringify({ active: record.active })
    })
    Message.success(record.active ? '策略已启用' : '策略已禁用')
  } catch (e) {
    record.active = !record.active
    Message.error('操作失败: ' + e.message)
  }
}

async function deletePolicy(record) {
  try {
    await apiRequest(`${API_BASE}/rtos/scheduler-policies/${record.id}`, { method: 'DELETE' })
    Message.success('删除成功')
    loadPolicies()
  } catch (e) {
    Message.error('删除失败: ' + e.message)
  }
}

function onDeviceChange() {
  loadAll()
  loadRtosConfig()
}

function onTaskPageChange(p) {
  taskPage.value = p
  loadTasks()
}
function onPolicyPageChange(p) {
  policyPage.value = p
  loadPolicies()
}

onMounted(() => {
  loadDevices()
})
</script>

<style scoped>
.text-success { color: var(--color-success); }
.text-warning { color: var(--color-warning); }
.text-danger { color: var(--color-danger); }
.text-secondary { color: var(--color-text-3); font-size: 12px; }
</style>
