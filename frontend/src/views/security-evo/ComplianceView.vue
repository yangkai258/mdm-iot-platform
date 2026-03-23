<template>
  <div class="page-container">
    <!-- 统计概览 -->
    <a-row :gutter="12" class="stat-row">
      <a-col :span="6">
        <a-card class="stat-card green" hoverable>
          <a-statistic title="合规率" :value="dashboard.compliance_rate" suffix="%" :value-from="0" :animation-duration="800" color="green">
            <template #icon><icon-check-circle /></template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card" hoverable>
          <a-statistic title="通过项" :value="dashboard.passed" :value-from="0" :animation-duration="600" color="green">
            <template #icon><icon-check /></template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card" hoverable>
          <a-statistic title="待修复" :value="dashboard.pending" :value-from="0" :animation-duration="600" :colored="dashboard.pending > 0 ? 'orange' : undefined">
            <template #icon><icon-clock /></template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card red" hoverable>
          <a-statistic title="不合规" :value="dashboard.failed" :value-from="0" :animation-duration="600" :colored="dashboard.failed > 0 ? 'red' : undefined">
            <template #icon><icon-close-circle /></template>
          </a-statistic>
        </a-card>
      </a-col>
    </a-row>

    <!-- Tab 切换 -->
    <a-card class="main-card">
      <a-tabs v-model:active-tab="activeTab">
        <!-- 合规检查项 -->
        <a-tab key="checks" title="合规检查">
          <div class="toolbar-row">
            <div class="toolbar-left">
              <a-select v-model="checkFilter.status" placeholder="状态" style="width: 110px" allow-clear @change="loadChecks">
                <a-option value="pass">通过</a-option>
                <a-option value="fail">失败</a-option>
                <a-option value="pending">待处理</a-option>
                <a-option value="warning">警告</a-option>
              </a-select>
              <a-select v-model="checkFilter.severity" placeholder="严重程度" style="width: 120px" allow-clear @change="loadChecks">
                <a-option value="critical">严重</a-option>
                <a-option value="high">高</a-option>
                <a-option value="medium">中</a-option>
                <a-option value="low">低</a-option>
              </a-select>
              <a-select v-model="checkFilter.regulation_id" placeholder="适用法规" style="width: 150px" allow-clear @change="loadChecks">
                <a-option v-for="r in regulations" :key="r.id" :value="r.id">{{ r.name }}</a-option>
              </a-select>
              <a-input-search v-model="checkFilter.keyword" placeholder="搜索检查项..." style="width: 180px" search-button @search="loadChecks" />
            </div>
            <div class="toolbar-right">
              <a-button @click="loadChecks"><template #icon><icon-refresh /></template>刷新</a-button>
            </div>
          </div>
          <a-table :columns="checkColumns" :data="checks" :loading="loading" :pagination="checkPagination" row-key="id" @change="handleCheckTableChange" :stripe="true" class="mt-3">
            <template #status="{ record }"><a-tag :color="statusColor(record.status)">{{ statusLabel(record.status) }}</a-tag></template>
            <template #severity="{ record }"><a-tag :color="severityColor(record.severity)">{{ severityLabel(record.severity) }}</a-tag></template>
            <template #regulation="{ record }"><span>{{ record.regulation_name || '-' }}</span></template>
            <template #last_check="{ record }"><span>{{ record.last_check_at ? formatDate(record.last_check_at) : '从未检查' }}</span></template>
            <template #actions="{ record }">
              <a-space>
                <a-button type="text" size="small" @click="runCheck(record)">执行检查</a-button>
                <a-button type="text" size="small" @click="openCheckDetail(record)">详情</a-button>
              </a-space>
            </template>
          </a-table>
        </a-tab>

        <!-- 合规报告 -->
        <a-tab key="reports" title="合规报告">
          <div class="toolbar-row">
            <div class="toolbar-left">
              <a-input-search v-model="reportKeyword" placeholder="搜索报告名称..." style="width: 200px" search-button @search="loadReports" />
            </div>
            <div class="toolbar-right">
              <a-button type="primary" @click="showReportModal = true"><template #icon><icon-plus /></template>生成报告</a-button>
              <a-button @click="loadReports"><template #icon><icon-refresh /></template>刷新</a-button>
            </div>
          </div>
          <a-table :columns="reportColumns" :data="reports" :loading="reportLoading" :pagination="reportPagination" row-key="id" @change="handleReportTableChange" :stripe="true" class="mt-3">
            <template #created_at="{ record }"><span>{{ formatDate(record.created_at) }}</span></template>
            <template #compliance_rate="{ record }">
              <a-progress :percent="record.compliance_rate" :color="record.compliance_rate >= 80 ? 'green' : record.compliance_rate >= 60 ? 'orange' : 'red'" />
            </template>
            <template #actions="{ record }">
              <a-space>
                <a-button type="text" size="small" @click="downloadReport(record)"><template #icon><icon-download /></template></a-button>
                <a-button type="text" size="small" @click="openReportDetail(record)">详情</a-button>
              </a-space>
            </template>
          </a-table>
        </a-tab>

        <!-- 法规清单 -->
        <a-tab key="regulations" title="法规清单">
          <div class="toolbar-row">
            <div class="toolbar-left">
              <a-input-search v-model="regKeyword" placeholder="搜索法规名称..." style="width: 200px" search-button @search="loadRegulations" />
            </div>
            <div class="toolbar-right">
              <a-button type="primary" @click="openRegulationModal()"><template #icon><icon-plus /></template>添加法规</a-button>
              <a-button @click="loadRegulations"><template #icon><icon-refresh /></template>刷新</a-button>
            </div>
          </div>
          <a-table :columns="regulationColumns" :data="regulations" :loading="regLoading" :pagination="regPagination" row-key="id" @change="handleRegTableChange" :stripe="true" class="mt-3">
            <template #type="{ record }"><a-tag>{{ typeLabel(record.type) }}</a-tag></template>
            <template #effective_date="{ record }"><span>{{ record.effective_date || '-' }}</span></template>
            <template #actions="{ record }">
              <a-space>
                <a-button type="text" size="small" @click="openRegulationModal(record)">编辑</a-button>
                <a-button type="text" size="small" status="danger" @click="handleDeleteRegulation(record)">删除</a-button>
              </a-space>
            </template>
          </a-table>
        </a-tab>
      </a-tabs>
    </a-card>

    <!-- 检查项详情弹窗 -->
    <a-modal v-model:visible="checkDetailVisible" title="检查项详情" width="640px" :footer="null">
      <div v-if="currentCheck" class="check-detail">
        <a-descriptions :column="2" bordered size="small">
          <a-descriptions-item label="检查项 ID" :span="2"><span class="mono">{{ currentCheck.id }}</span></a-descriptions-item>
          <a-descriptions-item label="名称" :span="2">{{ currentCheck.name }}</a-descriptions-item>
          <a-descriptions-item label="严重程度"><a-tag :color="severityColor(currentCheck.severity)">{{ severityLabel(currentCheck.severity) }}</a-tag></a-descriptions-item>
          <a-descriptions-item label="状态"><a-tag :color="statusColor(currentCheck.status)">{{ statusLabel(currentCheck.status) }}</a-tag></a-descriptions-item>
          <a-descriptions-item label="适用法规">{{ currentCheck.regulation_name || '-' }}</a-descriptions-item>
          <a-descriptions-item label="上次检查">{{ currentCheck.last_check_at ? formatDate(currentCheck.last_check_at) : '从未检查' }}</a-descriptions-item>
          <a-descriptions-item label="描述" :span="2">{{ currentCheck.description || '-' }}</a-descriptions-item>
        </a-descriptions>
        <div v-if="currentCheck.check_result" class="mt-3">
          <a-divider>检查结果</a-divider>
          <div class="code-block"><pre>{{ formatJson(currentCheck.check_result) }}</pre></div>
        </div>
        <div class="mt-3" style="text-align: right">
          <a-button type="primary" :loading="checkRunning" @click="runCheck(currentCheck)">重新执行检查</a-button>
        </div>
      </div>
    </a-modal>

    <!-- 报告详情弹窗 -->
    <a-modal v-model:visible="reportDetailVisible" title="合规报告详情" width="700px" :footer="null">
      <div v-if="currentReport">
        <a-descriptions :column="2" bordered size="small">
          <a-descriptions-item label="报告名称" :span="2">{{ currentReport.name }}</a-descriptions-item>
          <a-descriptions-item label="合规率"><a-progress :percent="currentReport.compliance_rate" /></a-descriptions-item>
          <a-descriptions-item label="生成时间">{{ formatDate(currentReport.created_at) }}</a-descriptions-item>
          <a-descriptions-item label="覆盖法规" :span="2">{{ currentReport.regulations_covered || '-' }}</a-descriptions-item>
          <a-descriptions-item label="通过项">{{ currentReport.passed_count || 0 }}</a-descriptions-item>
          <a-descriptions-item label="失败项">{{ currentReport.failed_count || 0 }}</a-descriptions-item>
        </a-descriptions>
        <div v-if="currentReport.details" class="mt-3">
          <a-divider>详细结果</a-divider>
          <div class="code-block"><pre>{{ formatJson(currentReport.details) }}</pre></div>
        </div>
      </div>
    </a-modal>

    <!-- 法规弹窗 -->
    <a-modal v-model:visible="regulationModalVisible" :title="editingRegulation.id ? '编辑法规' : '添加法规'" width="500px" @before-ok="handleSaveRegulation">
      <a-form :model="editingRegulation" layout="vertical">
        <a-form-item label="法规名称" required>
          <a-input v-model="editingRegulation.name" placeholder="如：GDPR、CCPA" />
        </a-form-item>
        <a-form-item label="法规类型" required>
          <a-select v-model="editingRegulation.type">
            <a-option value="gdpr">GDPR</a-option>
            <a-option value="ccpa">CCPA</a-option>
            <a-option value="sox">SOX</a-option>
            <a-option value="iso27001">ISO 27001</a-option>
            <a-option value="pci_dss">PCI DSS</a-option>
            <a-option value="hipaa">HIPAA</a-option>
            <a-option value="other">其他</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="生效日期">
          <a-date-picker v-model="editingRegulation.effective_date" style="width: 100%" />
        </a-form-item>
        <a-form-item label="描述">
          <a-textarea v-model="editingRegulation.description" :rows="3" placeholder="请输入法规描述" />
        </a-form-item>
        <a-form-item label="适用地区">
          <a-input v-model="editingRegulation.jurisdiction" placeholder="如：欧盟、美国加州" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import {
  getComplianceDashboard,
  getComplianceChecks,
  getComplianceCheckDetail,
  triggerComplianceCheck,
  getComplianceReports,
  getComplianceReportDetail,
  generateComplianceReport,
  downloadComplianceReport,
  getRegulations,
  saveRegulation,
  deleteRegulation
} from '@/api/security-evo'
import dayjs from 'dayjs'

const activeTab = ref('checks')
const loading = ref(false)
const reportLoading = ref(false)
const regLoading = ref(false)
const checkRunning = ref(false)

const dashboard = ref({ compliance_rate: 0, passed: 0, pending: 0, failed: 0 })
const checks = ref([])
const reports = ref([])
const regulations = ref([])

const checkDetailVisible = ref(false)
const reportDetailVisible = ref(false)
const regulationModalVisible = ref(false)
const showReportModal = ref(false)
const currentCheck = ref(null)
const currentReport = ref(null)
const editingRegulation = ref({})

const checkFilter = reactive({ status: '', severity: '', regulation_id: '', keyword: '' })
const reportKeyword = ref('')
const regKeyword = ref('')

const checkPagination = reactive({ current: 1, pageSize: 20, total: 0, showTotal: true })
const reportPagination = reactive({ current: 1, pageSize: 20, total: 0, showTotal: true })
const regPagination = reactive({ current: 1, pageSize: 20, total: 0, showTotal: true })

const checkColumns = [
  { title: '状态', slotName: 'status', width: 90 },
  { title: '检查项', dataIndex: 'name', width: 200 },
  { title: '严重程度', slotName: 'severity', width: 90 },
  { title: '适用法规', slotName: 'regulation', width: 150 },
  { title: '上次检查', slotName: 'last_check', width: 160 },
  { title: '描述', dataIndex: 'description', ellipsis: true },
  { title: '操作', slotName: 'actions', width: 140, fixed: 'right' }
]

const reportColumns = [
  { title: '报告名称', dataIndex: 'name', width: 200 },
  { title: '生成时间', slotName: 'created_at', width: 170 },
  { title: '合规率', slotName: 'compliance_rate', width: 180 },
  { title: '通过/失败', dataIndex: 'summary', width: 120, ellipsis: true },
  { title: '操作', slotName: 'actions', width: 120, fixed: 'right' }
]

const regulationColumns = [
  { title: '法规名称', dataIndex: 'name', width: 180 },
  { title: '类型', slotName: 'type', width: 120 },
  { title: '适用地区', dataIndex: 'jurisdiction', width: 150 },
  { title: '生效日期', slotName: 'effective_date', width: 140 },
  { title: '描述', dataIndex: 'description', ellipsis: true },
  { title: '操作', slotName: 'actions', width: 120, fixed: 'right' }
]

function statusColor(s) { const m = { pass: 'green', fail: 'red', pending: 'gray', warning: 'orange' }; return m[s] || 'default' }
function statusLabel(s) { const m = { pass: '通过', fail: '失败', pending: '待处理', warning: '警告' }; return m[s] || s }
function severityColor(s) { const m = { critical: 'red', high: 'orange', medium: 'arcoblue', low: 'green' }; return m[s] || 'default' }
function severityLabel(s) { const m = { critical: '严重', high: '高', medium: '中', low: '低' }; return m[s] || s }
function typeLabel(t) { const m = { gdpr: 'GDPR', ccpa: 'CCPA', sox: 'SOX', iso27001: 'ISO 27001', pci_dss: 'PCI DSS', hipaa: 'HIPAA', other: '其他' }; return m[t] || t }
function formatDate(d) { return d ? dayjs(d).format('YYYY-MM-DD HH:mm') : '-' }
function formatJson(obj) { try { return JSON.stringify(obj, null, 2) } catch { return String(obj) } }

function buildCheckParams() { return { page: checkPagination.current, page_size: checkPagination.pageSize, status: checkFilter.status || undefined, severity: checkFilter.severity || undefined, regulation_id: checkFilter.regulation_id || undefined, keyword: checkFilter.keyword || undefined } }

onMounted(() => { loadDashboard(); loadChecks(); loadReports(); loadRegulations() })

async function loadDashboard() {
  try { const r = await getComplianceDashboard(); dashboard.value = r.data || r } catch (e) { console.error(e) }
}
async function loadChecks() {
  loading.value = true
  try { const r = await getComplianceChecks(buildCheckParams()); const d = r.data || r; checks.value = d.list || d.records || []; checkPagination.total = d.total || checks.value.length } catch (e) { Message.error('加载检查项失败') } finally { loading.value = false }
}
async function loadReports() {
  reportLoading.value = true
  try { const r = await getComplianceReports({ page: reportPagination.current, page_size: reportPagination.pageSize, keyword: reportKeyword.value || undefined }); const d = r.data || r; reports.value = d.list || d.records || []; reportPagination.total = d.total || reports.value.length } catch (e) { Message.error('加载报告失败') } finally { reportLoading.value = false }
}
async function loadRegulations() {
  regLoading.value = true
  try { const r = await getRegulations({ page: regPagination.current, page_size: regPagination.pageSize, keyword: regKeyword.value || undefined }); const d = r.data || r; regulations.value = d.list || d.records || []; regPagination.total = d.total || regulations.value.length } catch (e) { Message.error('加载法规失败') } finally { regLoading.value = false }
}

function handleCheckTableChange(p) { checkPagination.current = p.current; checkPagination.pageSize = p.pageSize; loadChecks() }
function handleReportTableChange(p) { reportPagination.current = p.current; reportPagination.pageSize = p.pageSize; loadReports() }
function handleRegTableChange(p) { regPagination.current = p.current; regPagination.pageSize = p.pageSize; loadRegulations() }

async function runCheck(record) {
  checkRunning.value = true
  try { await triggerComplianceCheck(record.id); Message.success('检查执行中'); openCheckDetail(record) } catch (e) { Message.error('执行失败') } finally { checkRunning.value = false }
}

function openCheckDetail(record) { currentCheck.value = record; checkDetailVisible.value = true }
function openReportDetail(record) { currentReport.value = record; reportDetailVisible.value = true }

async function downloadReport(record) {
  try {
    const blob = await downloadComplianceReport(record.id)
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a'); a.href = url; a.download = `${record.name}_${dayjs().format('YYYYMMDD')}.pdf`; a.click()
    URL.revokeObjectURL(url)
    Message.success('下载成功')
  } catch (e) { Message.error('下载失败') }
}

function openRegulationModal(record = null) {
  editingRegulation.value = record ? { ...record } : { type: 'gdpr' }
  regulationModalVisible.value = true
}

async function handleSaveRegulation(done) {
  try {
    await saveRegulation(editingRegulation.value)
    Message.success('保存成功')
    regulationModalVisible.value = false
    loadRegulations()
    done(true)
  } catch (e) { Message.error('保存失败'); done(false) }
}

async function handleDeleteRegulation(record) {
  try { await deleteRegulation(record.id); Message.success('删除成功'); loadRegulations() } catch (e) { Message.error('删除失败') }
}
</script>

<style scoped>
.page-container { padding: 16px; display: flex; flex-direction: column; gap: 12px; height: 100%; box-sizing: border-box; }
.stat-row .stat-card { text-align: center; }
.main-card { flex: 1; overflow: auto; }
.toolbar-row { display: flex; align-items: center; justify-content: space-between; gap: 10px; }
.toolbar-left { display: flex; align-items: center; gap: 10px; flex-wrap: wrap; }
.toolbar-right { display: flex; gap: 8px; flex-shrink: 0; }
.mt-3 { margin-top: 12px; }
.mono { font-family: 'Courier New', monospace; font-size: 12px; }
.check-detail { display: flex; flex-direction: column; gap: 12px; }
.code-block { background: var(--color-fill-1); border-radius: 4px; padding: 12px; max-height: 200px; overflow: auto; }
.code-block pre { margin: 0; font-family: 'Courier New', monospace; font-size: 12px; white-space: pre-wrap; word-break: break-all; }
</style>
