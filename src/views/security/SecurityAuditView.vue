<template>
  <div class="page-container">
    <!-- 安全概览 -->
    <a-row :gutter="16">
      <a-col :span="6">
        <a-card class="stat-card">
          <div class="stat-icon security">
            <icon-shield />
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ overview.securityScore }}</div>
            <div class="stat-label">安全评分</div>
          </div>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <div class="stat-icon warning">
            <icon-warning />
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ overview.pendingAlerts }}</div>
            <div class="stat-label">待处理告警</div>
          </div>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <div class="stat-icon danger">
            <icon-exclamation />
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ overview.failedLogins24h }}</div>
            <div class="stat-label">24h 失败登录</div>
          </div>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <div class="stat-icon success">
            <icon-check-circle />
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ overview.complianceRate }}</div>
            <div class="stat-label">合规率</div>
          </div>
        </a-card>
      </a-col>
    </a-row>

    <!-- 标签页切换 -->
    <a-card class="main-card">
      <a-tabs v-model:active-tab="activeTab">
        <!-- 审计日志 -->
        <a-tab-pane key="logs" title="安全日志">
          <div class="toolbar-row">
            <div class="toolbar-left">
              <a-select
                v-model="filter.category"
                placeholder="安全类别"
                style="width: 130px"
                allow-clear
                @change="loadSecurityLogs"
              >
                <a-option value="authentication">身份认证</a-option>
                <a-option value="authorization">权限变更</a-option>
                <a-option value="data_access">数据访问</a-option>
                <a-option value="config_change">配置变更</a-option>
                <a-option value="threat">威胁检测</a-option>
                <a-option value="compliance">合规审计</a-option>
              </a-select>
              <a-select
                v-model="filter.severity"
                placeholder="严重级别"
                style="width: 110px"
                allow-clear
                @change="loadSecurityLogs"
              >
                <a-option value="critical">严重</a-option>
                <a-option value="high">高</a-option>
                <a-option value="medium">中</a-option>
                <a-option value="low">低</a-option>
                <a-option value="info">信息</a-option>
              </a-select>
              <a-select
                v-model="filter.event_type"
                placeholder="事件类型"
                style="width: 140px"
                allow-clear
                @change="loadSecurityLogs"
              >
                <a-option value="login_success">登录成功</a-option>
                <a-option value="login_failure">登录失败</a-option>
                <a-option value="logout">登出</a-option>
                <a-option value="password_change">密码修改</a-option>
                <a-option value="permission_change">权限变更</a-option>
                <a-option value="data_export">数据导出</a-option>
                <a-option value="api_access">API 访问</a-option>
                <a-option value="suspicious_activity">可疑活动</a-option>
              </a-select>
              <a-range-picker
                v-model="filter.dateRange"
                style="width: 260px"
                @change="loadSecurityLogs"
              />
              <a-input-search
                v-model="filter.keyword"
                placeholder="搜索事件/操作人..."
                style="width: 200px"
                @search="loadSecurityLogs"
                @press-enter="loadSecurityLogs"
              />
            </div>
            <div class="toolbar-right">
              <a-button @click="exportSecurityLogs">
                <template #icon><icon-download /></template>
                导出
              </a-button>
              <a-button @click="loadSecurityLogs">
                <template #icon><icon-refresh /></template>
                刷新
              </a-button>
            </div>
          </div>

          <a-table
            :columns="logColumns"
            :data="securityLogs"
            :loading="loading"
            :pagination="pagination"
            row-key="id"
            @change="handleTableChange"
            :scroll="{ x: 1300 }"
          >
            <template #timestamp="{ record }">
              <span class="mono-small">{{ formatTime(record.timestamp) }}</span>
            </template>
            <template #severity="{ record }">
              <a-tag :color="severityColor(record.severity)">{{ severityLabel(record.severity) }}</a-tag>
            </template>
            <template #category="{ record }">
              <span>{{ categoryLabel(record.category) }}</span>
            </template>
            <template #event_type="{ record }">
              <a-tag>{{ eventTypeLabel(record.event_type) }}</a-tag>
            </template>
            <template #user="{ record }">
              <span>{{ record.user_name || record.user_id || '-' }}</span>
            </template>
            <template #source_ip="{ record }">
              <span class="mono-small">{{ record.source_ip || '-' }}</span>
            </template>
            <template #result="{ record }">
              <a-tag :color="resultColor(record.result)">{{ resultLabel(record.result) }}</a-tag>
            </template>
            <template #actions="{ record }">
              <a-button type="text" size="small" @click="openLogDetail(record)">详情</a-button>
            </template>
          </a-table>
        </a-tab-pane>

        <!-- 安全报告 -->
        <a-tab-pane key="reports" title="安全报告">
          <div class="toolbar-row">
            <div class="toolbar-left">
              <a-select
                v-model="reportFilter.type"
                placeholder="报告类型"
                style="width: 140px"
                allow-clear
                @change="loadReports"
              >
                <a-option value="daily">日报</a-option>
                <a-option value="weekly">周报</a-option>
                <a-option value="monthly">月报</a-option>
                <a-option value="quarterly">季度报告</a-option>
                <a-option value="incident">事件报告</a-option>
              </a-select>
              <a-range-picker
                v-model="reportFilter.dateRange"
                style="width: 260px"
                @change="loadReports"
              />
            </div>
            <div class="toolbar-right">
              <a-button type="primary" @click="generateReport">
                <template #icon><icon-file-plus /></template>
                生成报告
              </a-button>
              <a-button @click="loadReports">
                <template #icon><icon-refresh /></template>
                刷新
              </a-button>
            </div>
          </div>

          <a-table
            :columns="reportColumns"
            :data="reports"
            :loading="loadingReports"
            :pagination="reportPagination"
            row-key="id"
            @change="handleReportTableChange"
          >
            <template #type="{ record }">
              <a-tag :color="reportTypeColor(record.type)">{{ reportTypeLabel(record.type) }}</a-tag>
            </template>
            <template #period="{ record }">
              <span>{{ record.period_start ? formatDate(record.period_start) : '-' }} ~ {{ record.period_end ? formatDate(record.period_end) : '-' }}</span>
            </template>
            <template #status="{ record }">
              <a-badge :status="reportStatusBadge(record.status)" />
              <span>{{ reportStatusLabel(record.status) }}</span>
            </template>
            <template #created_at="{ record }">
              {{ formatTime(record.created_at) }}
            </template>
            <template #actions="{ record }">
              <a-button type="text" size="small" @click="downloadReport(record)">
                <icon-download /> 下载
              </a-button>
              <a-button
                v-if="record.status === 'generating'"
                type="text"
                size="small"
                @click="checkReportStatus(record)"
              >
                刷新状态
              </a-button>
            </template>
          </a-table>
        </a-tab-pane>

        <!-- 威胁概览 -->
        <a-tab-pane key="threats" title="威胁概览">
          <a-row :gutter="16">
            <a-col :span="12">
              <a-card title="最近威胁">
                <a-table
                  :columns="threatColumns"
                  :data="recentThreats"
                  :loading="loadingThreats"
                  :pagination="false"
                  row-key="id"
                  size="small"
                >
                  <template #level="{ record }">
                    <a-tag :color="threatLevelColor(record.level)">{{ record.level }}</a-tag>
                  </template>
                  <template #status="{ record }">
                    <a-badge :status="threatStatusBadge(record.status)" />
                    {{ record.status }}
                  </template>
                  <template #detected_at="{ record }">
                    {{ formatTime(record.detected_at) }}
                  </template>
                </a-table>
              </a-card>
            </a-col>
            <a-col :span="12">
              <a-card title="威胁类型分布">
                <div class="threat-distribution">
                  <div
                    v-for="item in threatDistribution"
                    :key="item.type"
                    class="threat-item"
                  >
                    <div class="threat-label">
                      <span>{{ item.type }}</span>
                      <span>{{ item.count }}</span>
                    </div>
                    <a-progress
                      :percent="item.percent"
                      :color="threatLevelColor(item.level)"
                      :show-text="false"
                    />
                  </div>
                </div>
              </a-card>
            </a-col>
          </a-row>
        </a-tab-pane>
      </a-tabs>
    </a-card>

    <!-- 日志详情弹窗 -->
    <a-modal
      v-model:visible="detailVisible"
      title="安全日志详情"
      width="680px"
      :footer="null"
    >
      <div class="log-detail" v-if="currentLog">
        <a-descriptions :column="2" bordered size="small">
          <a-descriptions-item label="日志 ID" :span="2">
            <span class="mono-small">{{ currentLog.id }}</span>
          </a-descriptions-item>
          <a-descriptions-item label="时间">
            {{ formatTime(currentLog.timestamp) }}
          </a-descriptions-item>
          <a-descriptions-item label="严重级别">
            <a-tag :color="severityColor(currentLog.severity)">{{ severityLabel(currentLog.severity) }}</a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="类别">
            {{ categoryLabel(currentLog.category) }}
          </a-descriptions-item>
          <a-descriptions-item label="事件类型">
            <a-tag>{{ eventTypeLabel(currentLog.event_type) }}</a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="操作人">
            {{ currentLog.user_name || currentLog.user_id || '-' }}
          </a-descriptions-item>
          <a-descriptions-item label="来源 IP">
            <span class="mono-small">{{ currentLog.source_ip || '-' }}</span>
          </a-descriptions-item>
          <a-descriptions-item label="结果">
            <a-tag :color="resultColor(currentLog.result)">{{ resultLabel(currentLog.result) }}</a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="资源" :span="2">
            {{ currentLog.resource_type || '-' }} / {{ currentLog.resource_id || '-' }}
          </a-descriptions-item>
          <a-descriptions-item label="描述" :span="2">
            {{ currentLog.description || '-' }}
          </a-descriptions-item>
        </a-descriptions>

        <a-divider>事件详情</a-divider>
        <div class="code-block" v-if="currentLog.details">
          <pre>{{ formatJson(currentLog.details) }}</pre>
        </div>
        <div class="code-block" v-else>
          <span class="no-data">无详情</span>
        </div>
      </div>
    </a-modal>

    <!-- 生成报告弹窗 -->
    <a-modal
      v-model:visible="generateVisible"
      title="生成安全报告"
      width="480px"
      @ok="handleGenerateReport"
      :ok-loading="generating"
    >
      <a-form :model="generateForm" layout="vertical">
        <a-form-item label="报告类型" required>
          <a-select v-model="generateForm.type" placeholder="请选择报告类型">
            <a-option value="daily">日报</a-option>
            <a-option value="weekly">周报</a-option>
            <a-option value="monthly">月报</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="报告周期" required>
          <a-range-picker v-model="generateForm.period" style="width: 100%" />
        </a-form-item>
        <a-form-item label="包含内容">
          <a-checkbox-group v-model="generateForm.includes">
            <a-checkbox value="authentication">身份认证事件</a-checkbox>
            <a-checkbox value="authorization">权限变更记录</a-checkbox>
            <a-checkbox value="data_access">数据访问日志</a-checkbox>
            <a-checkbox value="threat">威胁检测报告</a-checkbox>
            <a-checkbox value="compliance">合规状态摘要</a-checkbox>
          </a-checkbox-group>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import {
  getSecurityLogs,
  getSecurityLogDetail,
  exportSecurityLogsApi,
  getSecurityOverview,
  getSecurityReports,
  generateSecurityReport,
  downloadSecurityReport,
  getRecentThreats,
  getThreatDistribution
} from '@/api/security'
import dayjs from 'dayjs'

const activeTab = ref('logs')
const loading = ref(false)
const loadingReports = ref(false)
const loadingThreats = ref(false)
const securityLogs = ref([])
const reports = ref([])
const recentThreats = ref([])
const threatDistribution = ref([])
const detailVisible = ref(false)
const currentLog = ref(null)
const generateVisible = ref(false)
const generating = ref(false)

const overview = reactive({
  securityScore: 0,
  pendingAlerts: 0,
  failedLogins24h: 0,
  complianceRate: '0%'
})

const filter = reactive({
  category: '',
  severity: '',
  event_type: '',
  dateRange: [],
  keyword: ''
})

const reportFilter = reactive({
  type: '',
  dateRange: []
})

const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0,
  showTotal: true,
  showPageSize: true
})

const reportPagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0,
  showTotal: true,
  showPageSize: true
})

const generateForm = reactive({
  type: 'weekly',
  period: [],
  includes: ['authentication', 'authorization', 'threat']
})

const logColumns = [
  { title: '时间', slotName: 'timestamp', width: 170 },
  { title: '级别', slotName: 'severity', width: 80 },
  { title: '类别', slotName: 'category', width: 110 },
  { title: '事件类型', slotName: 'event_type', width: 110 },
  { title: '操作人', slotName: 'user', width: 120 },
  { title: '来源 IP', slotName: 'source_ip', width: 130 },
  { title: '结果', slotName: 'result', width: 80 },
  { title: '操作', slotName: 'actions', width: 80, fixed: 'right' }
]

const reportColumns = [
  { title: '报告名称', dataIndex: 'name', minWidth: 200 },
  { title: '类型', slotName: 'type', width: 100 },
  { title: '周期', slotName: 'period', width: 220 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '生成时间', slotName: 'created_at', width: 170 },
  { title: '操作', slotName: 'actions', width: 120 }
]

const threatColumns = [
  { title: '威胁', dataIndex: 'threat_type', minWidth: 150 },
  { title: '级别', slotName: 'level', width: 80 },
  { title: '状态', slotName: 'status', width: 90 },
  { title: '检测时间', slotName: 'detected_at', width: 170 }
]

function severityColor(severity) {
  const map = {
    critical: 'red',
    high: 'orangered',
    medium: 'orange',
    low: 'blue',
    info: 'gray'
  }
  return map[severity] || 'default'
}

function severityLabel(severity) {
  const map = {
    critical: '严重',
    high: '高',
    medium: '中',
    low: '低',
    info: '信息'
  }
  return map[severity] || severity
}

function categoryLabel(category) {
  const map = {
    authentication: '身份认证',
    authorization: '权限变更',
    data_access: '数据访问',
    config_change: '配置变更',
    threat: '威胁检测',
    compliance: '合规审计'
  }
  return map[category] || category
}

function eventTypeLabel(type) {
  const map = {
    login_success: '登录成功',
    login_failure: '登录失败',
    logout: '登出',
    password_change: '密码修改',
    permission_change: '权限变更',
    data_export: '数据导出',
    api_access: 'API 访问',
    suspicious_activity: '可疑活动'
  }
  return map[type] || type
}

function resultColor(result) {
  return result === 'success' ? 'green' : 'red'
}

function resultLabel(result) {
  return result === 'success' ? '成功' : '失败'
}

function formatTime(timestamp) {
  if (!timestamp) return '-'
  return dayjs(timestamp).format('YYYY-MM-DD HH:mm:ss')
}

function formatDate(date) {
  if (!date) return '-'
  return dayjs(date).format('YYYY-MM-DD')
}

function formatJson(obj) {
  if (!obj) return ''
  try {
    return JSON.stringify(obj, null, 2)
  } catch {
    return String(obj)
  }
}

function reportTypeColor(type) {
  const map = {
    daily: 'blue',
    weekly: 'arcoblue',
    monthly: 'purple',
    quarterly: 'magenta',
    incident: 'red'
  }
  return map[type] || 'default'
}

function reportTypeLabel(type) {
  const map = {
    daily: '日报',
    weekly: '周报',
    monthly: '月报',
    quarterly: '季度',
    incident: '事件'
  }
  return map[type] || type
}

function reportStatusBadge(status) {
  const map = {
    ready: 'success',
    generating: 'processing',
    failed: 'danger'
  }
  return map[status] || 'default'
}

function reportStatusLabel(status) {
  const map = {
    ready: '就绪',
    generating: '生成中',
    failed: '失败'
  }
  return map[status] || status
}

function threatLevelColor(level) {
  const map = {
    critical: 'red',
    high: 'orangered',
    medium: 'orange',
    low: 'blue'
  }
  return map[level] || 'default'
}

function threatStatusBadge(status) {
  const map = {
    resolved: 'success',
    investigating: 'warning',
    open: 'danger'
  }
  return map[status] || 'default'
}

function handleTableChange(pag) {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  loadSecurityLogs()
}

function handleReportTableChange(pag) {
  reportPagination.current = pag.current
  reportPagination.pageSize = pag.pageSize
  loadReports()
}

async function loadOverview() {
  try {
    const data = await getSecurityOverview()
    const res = data.data || data
    overview.securityScore = res.security_score || 0
    overview.pendingAlerts = res.pending_alerts || 0
    overview.failedLogins24h = res.failed_logins_24h || 0
    overview.complianceRate = res.compliance_rate || '0%'
  } catch (e) {
    console.error('加载安全概览失败', e)
  }
}

async function loadSecurityLogs() {
  loading.value = true
  try {
    const params = {
      page: pagination.current,
      page_size: pagination.pageSize,
      category: filter.category || undefined,
      severity: filter.severity || undefined,
      event_type: filter.event_type || undefined,
      keyword: filter.keyword || undefined,
      start_date: filter.dateRange[0] ? dayjs(filter.dateRange[0]).format('YYYY-MM-DD') : undefined,
      end_date: filter.dateRange[1] ? dayjs(filter.dateRange[1]).format('YYYY-MM-DD') : undefined
    }
    const data = await getSecurityLogs(params)
    const res = data.data || data
    securityLogs.value = res.list || res.records || []
    pagination.total = res.total || securityLogs.value.length
  } catch (e) {
    console.error('加载安全日志失败', e)
    Message.error('加载安全日志失败')
  } finally {
    loading.value = false
  }
}

async function loadReports() {
  loadingReports.value = true
  try {
    const params = {
      page: reportPagination.current,
      page_size: reportPagination.pageSize,
      type: reportFilter.type || undefined,
      start_date: reportFilter.dateRange[0] ? dayjs(reportFilter.dateRange[0]).format('YYYY-MM-DD') : undefined,
      end_date: reportFilter.dateRange[1] ? dayjs(reportFilter.dateRange[1]).format('YYYY-MM-DD') : undefined
    }
    const data = await getSecurityReports(params)
    const res = data.data || data
    reports.value = res.list || res.records || []
    reportPagination.total = res.total || reports.value.length
  } catch (e) {
    console.error('加载报告列表失败', e)
  } finally {
    loadingReports.value = false
  }
}

async function loadThreats() {
  loadingThreats.value = true
  try {
    const [threatsData, distData] = await Promise.all([
      getRecentThreats(),
      getThreatDistribution()
    ])
    recentThreats.value = (threatsData.data || threatsData).list || []
    threatDistribution.value = (distData.data || distData).distribution || []
  } catch (e) {
    console.error('加载威胁数据失败', e)
  } finally {
    loadingThreats.value = false
  }
}

function openLogDetail(log) {
  currentLog.value = log
  detailVisible.value = true
}

async function exportSecurityLogs() {
  try {
    const params = {
      category: filter.category || undefined,
      severity: filter.severity || undefined,
      event_type: filter.event_type || undefined,
      keyword: filter.keyword || undefined,
      start_date: filter.dateRange[0] ? dayjs(filter.dateRange[0]).format('YYYY-MM-DD') : undefined,
      end_date: filter.dateRange[1] ? dayjs(filter.dateRange[1]).format('YYYY-MM-DD') : undefined
    }
    const blob = await exportSecurityLogsApi(params)
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `security_audit_${dayjs().format('YYYYMMDD_HHmmss')}.json`
    a.click()
    URL.revokeObjectURL(url)
    Message.success('导出成功')
  } catch (e) {
    Message.error('导出失败')
  }
}

function generateReport() {
  generateForm.type = 'weekly'
  generateForm.period = []
  generateForm.includes = ['authentication', 'authorization', 'threat']
  generateVisible.value = true
}

async function handleGenerateReport() {
  if (!generateForm.type || !generateForm.period || generateForm.period.length < 2) {
    Message.warning('请填写完整信息')
    return
  }
  generating.value = true
  try {
    await generateSecurityReport({
      type: generateForm.type,
      period_start: dayjs(generateForm.period[0]).format('YYYY-MM-DD'),
      period_end: dayjs(generateForm.period[1]).format('YYYY-MM-DD'),
      includes: generateForm.includes
    })
    Message.success('报告生成中...')
    generateVisible.value = false
    loadReports()
  } catch (e) {
    Message.error('生成失败')
  } finally {
    generating.value = false
  }
}

async function downloadReport(record) {
  try {
    const blob = await downloadSecurityReport(record.id)
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = record.name || `security_report_${record.id}`
    a.click()
    URL.revokeObjectURL(url)
    Message.success('下载成功')
  } catch (e) {
    Message.error('下载失败')
  }
}

async function checkReportStatus(record) {
  loadReports()
}

onMounted(() => {
  loadOverview()
  loadSecurityLogs()
  loadReports()
  loadThreats()
})
</script>

<style scoped>
.page-container {
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
  height: 100%;
  box-sizing: border-box;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 16px;
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
}

.stat-icon.security {
  background: rgba(22, 93, 255, 0.1);
  color: rgb(var(--primary-6));
}

.stat-icon.warning {
  background: rgba(240, 152, 0, 0.1);
  color: rgb(240, 152, 0);
}

.stat-icon.danger {
  background: rgba(255, 77, 79, 0.1);
  color: rgb(255, 77, 79);
}

.stat-icon.success {
  background: rgba(22, 196, 112, 0.1);
  color: rgb(22, 196, 112);
}

.stat-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.stat-value {
  font-size: 24px;
  font-weight: 700;
}

.stat-label {
  font-size: 13px;
  color: var(--color-text-3);
}

.main-card {
  flex: 1;
  overflow: auto;
}

.toolbar-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  margin-bottom: 16px;
  flex-wrap: wrap;
}

.toolbar-left {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.toolbar-right {
  display: flex;
  gap: 8px;
}

.mono-small {
  font-family: monospace;
  font-size: 12px;
}

.log-detail {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.code-block {
  background: var(--color-fill-1);
  border-radius: 4px;
  padding: 12px;
  max-height: 200px;
  overflow: auto;
}

.code-block pre {
  margin: 0;
  font-family: 'Courier New', monospace;
  font-size: 12px;
  white-space: pre-wrap;
  word-break: break-all;
}

.no-data {
  color: var(--color-text-3);
  font-size: 13px;
}

.threat-distribution {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.threat-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.threat-label {
  display: flex;
  justify-content: space-between;
  font-size: 13px;
}
</style>
