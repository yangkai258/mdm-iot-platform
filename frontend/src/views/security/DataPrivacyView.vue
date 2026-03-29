<template>
  <div class="page-container">
    <Breadcrumb :items="['menu.security', 'menu.security.dataPrivacy']" />
    <!-- 脱敏配置 -->
    <a-card class="section-card">
      <template #title>
        <div class="card-title">
          <icon-eye-invisible />
          <span>数据脱敏配置</span>
        </div>
      </template>
      <a-form :model="desensitizationConfig" layout="vertical">
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="全局脱敏开关">
              <a-switch v-model="desensitizationConfig.enabled" />
              <span class="config-hint ml-3">{{ desensitizationConfig.enabled ? '已启用' : '已禁用' }}</span>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="默认脱敏规则">
              <a-select v-model="desensitizationConfig.defaultRule">
                <a-option value="full">完全脱敏</a-option>
                <a-option value="partial">部分脱敏</a-option>
                <a-option value="hash">哈希替换</a-option>
              </a-select>
            </a-form-item>
          </a-col>
        </a-row>
        <a-divider>字段脱敏规则</a-divider>
        <a-table
          :columns="fieldColumns"
          :data="desensitizationConfig.fields"
          :pagination="false"
          row-key="field"
          size="small"
        >
          <template #field="{ record }">
            <span>{{ fieldLabel(record.field) }}</span>
          </template>
          <template #rule="{ record, rowIndex }">
            <a-select v-model="record.rule" style="width: 140px" size="small">
              <a-option value="full">完全脱敏</a-option>
              <a-option value="partial">部分脱敏</a-option>
              <a-option value="hash">哈希替换</a-option>
              <a-option value="none">不脱敏</a-option>
            </a-select>
          </template>
          <template #enabled="{ record, rowIndex }">
            <a-switch v-model="record.enabled" size="small" />
          </template>
        </a-table>
        <a-form-item style="margin-top: 16px">
          <a-button type="primary" :loading="saving.desensitization" @click="saveDesensitization">
            保存配置
          </a-button>
        </a-form-item>
      </a-form>
    </a-card>

    <!-- GDPR 合规 -->
    <a-card class="section-card">
      <template #title>
        <div class="card-title">
          <icon-safe />
          <span>GDPR 合规</span>
        </div>
      </template>
      <a-form :model="gdprConfig" layout="vertical">
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="数据主体权利">
              <a-switch v-model="gdprConfig.rightToAccess" />
              <span class="config-hint">访问权</span>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label=" ">
              <a-switch v-model="gdprConfig.rightToRectification" />
              <span class="config-hint">更正权</span>
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label=" ">
              <a-switch v-model="gdprConfig.rightToErasure" />
              <span class="config-hint">被遗忘权</span>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label=" ">
              <a-switch v-model="gdprConfig.rightToPortability" />
              <span class="config-hint">数据可携权</span>
            </a-form-item>
          </a-col>
        </a-row>
        <a-divider>数据保留策略</a-divider>
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="用户数据保留期">
              <a-select v-model="gdprConfig.retentionPeriod">
                <a-option value="30">30 天</a-option>
                <a-option value="90">90 天</a-option>
                <a-option value="180">180 天</a-option>
                <a-option value="365">1 年</a-option>
                <a-option value="730">2 年</a-option>
                <a-option value="forever">永久保留</a-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="日志数据保留期">
              <a-select v-model="gdprConfig.logRetentionPeriod">
                <a-option value="30">30 天</a-option>
                <a-option value="90">90 天</a-option>
                <a-option value="180">180 天</a-option>
                <a-option value="365">1 年</a-option>
                <a-option value="730">2 年</a-option>
              </a-select>
            </a-form-item>
          </a-col>
        </a-row>
        <a-divider>同意管理</a-divider>
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="Cookie 同意">
              <a-switch v-model="gdprConfig.cookieConsent" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="隐私政策版本">
              <a-input v-model="gdprConfig.privacyPolicyVersion" readonly />
            </a-form-item>
          </a-col>
        </a-row>
        <a-form-item>
          <a-button type="primary" :loading="saving.gdpr" @click="saveGdpr">
            保存配置
          </a-button>
        </a-form-item>
      </a-form>
    </a-card>

    <!-- 数据导出 -->
    <a-card class="section-card">
      <template #title>
        <div class="card-title">
          <icon-download />
          <span>数据导出</span>
        </div>
      </template>
      <a-form :model="exportConfig" layout="vertical">
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="导出格式">
              <a-select v-model="exportConfig.format">
                <a-option value="json">JSON</a-option>
                <a-option value="csv">CSV</a-option>
                <a-option value="xlsx">Excel (XLSX)</a-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="数据范围">
              <a-select v-model="exportConfig.scope">
                <a-option value="all">全部数据</a-option>
                <a-option value="user">用户数据</a-option>
                <a-option value="device">设备数据</a-option>
                <a-option value="activity">活动日志</a-option>
              </a-select>
            </a-form-item>
          </a-col>
        </a-row>
        <a-form-item label="日期范围">
          <a-range-picker v-model="exportConfig.dateRange" style="width: 100%" />
        </a-form-item>
        <a-form-item>
          <a-space>
            <a-button type="primary" :loading="exporting" @click="handleExport">
              发起导出请求
            </a-button>
            <a-button @click="viewExportHistory">
              查看导出历史
            </a-button>
          </a-space>
        </a-form-item>
      </a-form>

      <a-divider>导出历史</a-divider>
      <a-table
        :columns="exportColumns"
        :data="exportHistory"
        :loading="loading.history"
        :pagination="pagination"
        row-key="id"
        @change="handleHistoryTableChange"
      >
        <template #created_at="{ record }">
          {{ formatDate(record.created_at) }}
        </template>
        <template #status="{ record }">
          <a-tag :color="exportStatusColor(record.status)">{{ exportStatusLabel(record.status) }}</a-tag>
        </template>
        <template #actions="{ record }">
          <a-button
            type="text"
            size="small"
            :disabled="record.status !== 'completed'"
            @click="downloadExport(record)"
          >
            下载
          </a-button>
          <a-divider direction="vertical" />
          <a-popconfirm content="确定删除该导出记录？" @ok="deleteExport(record)">
            <a-button type="text" size="small" status="danger">删除</a-button>
          </a-popconfirm>
        </template>
      </a-table>
    </a-card>

    <!-- GDPR 请求弹窗 -->
    <a-modal
      v-model:visible="gdprRequestVisible"
      title="提交 GDPR 请求"
      width="480px"
      @ok="handleGdprRequest"
      :ok-loading="submittingRequest"
    >
      <a-form :model="gdprRequestForm" layout="vertical">
        <a-form-item label="请求类型" required>
          <a-select v-model="gdprRequestForm.type" placeholder="请选择请求类型">
            <a-option value="access">数据访问请求</a-option>
            <a-option value="rectification">数据更正请求</a-option>
            <a-option value="erasure">数据删除请求</a-option>
            <a-option value="portability">数据可携请求</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="用户标识" required>
          <a-input v-model="gdprRequestForm.userId" placeholder="请输入用户 ID 或邮箱" />
        </a-form-item>
        <a-form-item label="附加说明">
          <a-textarea v-model="gdprRequestForm.description" :rows="3" placeholder="请输入其他说明" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import {
  getDesensitizationConfig,
  updateDesensitizationConfig,
  getGdprConfig,
  updateGdprConfig,
  exportData,
  getExportHistory,
  deleteExportRecord
} from '@/api/compliance'
import dayjs from 'dayjs'

const loading = reactive({ history: false })
const saving = reactive({ desensitization: false, gdpr: false })
const exporting = ref(false)
const submittingRequest = ref(false)

const desensitizationConfig = reactive({
  enabled: true,
  defaultRule: 'partial',
  fields: [
    { field: 'phone', label: '手机号', rule: 'partial', enabled: true },
    { field: 'email', label: '邮箱', rule: 'partial', enabled: true },
    { field: 'id_card', label: '身份证', rule: 'full', enabled: true },
    { field: 'name', label: '姓名', rule: 'partial', enabled: true },
    { field: 'address', label: '地址', rule: 'partial', enabled: true },
    { field: 'bank_card', label: '银行卡', rule: 'full', enabled: true }
  ]
})

const gdprConfig = reactive({
  rightToAccess: true,
  rightToRectification: true,
  rightToErasure: true,
  rightToPortability: true,
  retentionPeriod: '365',
  logRetentionPeriod: '180',
  cookieConsent: true,
  privacyPolicyVersion: 'v2.1'
})

const exportConfig = reactive({
  format: 'json',
  scope: 'all',
  dateRange: []
})

const exportHistory = ref([])
const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0
})

const gdprRequestVisible = ref(false)
const gdprRequestForm = reactive({
  type: 'access',
  userId: '',
  description: ''
})

const fieldColumns = [
  { title: '字段', slotName: 'field', width: 120 },
  { title: '脱敏规则', slotName: 'rule', width: 160 },
  { title: '启用', slotName: 'enabled', width: 80 }
]

const exportColumns = [
  { title: 'ID', dataIndex: 'id', width: 80 },
  { title: '请求时间', slotName: 'created_at', width: 160 },
  { title: '格式', dataIndex: 'format', width: 80 },
  { title: '范围', dataIndex: 'scope', width: 100 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '操作', slotName: 'actions', width: 140, fixed: 'right' }
]

function fieldLabel(field) {
  const map = {
    phone: '手机号',
    email: '邮箱',
    id_card: '身份证',
    name: '姓名',
    address: '地址',
    bank_card: '银行卡'
  }
  return map[field] || field
}

function exportStatusColor(status) {
  const map = { pending: 'yellow', processing: 'blue', completed: 'green', failed: 'red' }
  return map[status] || 'default'
}

function exportStatusLabel(status) {
  const map = { pending: '待处理', processing: '处理中', completed: '已完成', failed: '失败' }
  return map[status] || status
}

function formatDate(date) {
  if (!date) return '-'
  return dayjs(date).format('YYYY-MM-DD HH:mm')
}

onMounted(() => {
  loadDesensitizationConfig()
  loadGdprConfig()
  loadExportHistory()
})

async function loadDesensitizationConfig() {
  try {
    const res = await getDesensitizationConfig()
    const data = res.data || res
    if (data.fields && data.fields.length > 0) {
      Object.assign(desensitizationConfig, data)
    }
  } catch (e) {
    console.error('加载脱敏配置失败', e)
  }
}

async function loadGdprConfig() {
  try {
    const res = await getGdprConfig()
    const data = res.data || res
    Object.assign(gdprConfig, data)
  } catch (e) {
    console.error('加载GDPR配置失败', e)
  }
}

async function loadExportHistory() {
  loading.history = true
  try {
    const params = {
      page: pagination.current,
      page_size: pagination.pageSize
    }
    const res = await getExportHistory(params)
    const data = res.data || res
    exportHistory.value = data.list || data.records || []
    pagination.total = data.total || exportHistory.value.length
  } catch (e) {
    console.error('加载导出历史失败', e)
  } finally {
    loading.history = false
  }
}

function handleHistoryTableChange(pag) {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  loadExportHistory()
}

async function saveDesensitization() {
  saving.desensitization = true
  try {
    await updateDesensitizationConfig(desensitizationConfig)
    Message.success('保存成功')
  } catch (e) {
    Message.error('保存失败')
  } finally {
    saving.desensitization = false
  }
}

async function saveGdpr() {
  saving.gdpr = true
  try {
    await updateGdprConfig(gdprConfig)
    Message.success('保存成功')
  } catch (e) {
    Message.error('保存失败')
  } finally {
    saving.gdpr = false
  }
}

async function handleExport() {
  exporting.value = true
  try {
    const params = {
      format: exportConfig.format,
      scope: exportConfig.scope,
      start_date: exportConfig.dateRange[0] || null,
      end_date: exportConfig.dateRange[1] || null
    }
    await exportData(params)
    Message.success('导出请求已提交')
    loadExportHistory()
  } catch (e) {
    Message.error('导出请求失败')
  } finally {
    exporting.value = false
  }
}

function viewExportHistory() {
  loadExportHistory()
}

async function downloadExport(record) {
  try {
    const blob = await fetch(record.download_url).then(r => r.blob())
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `export_${record.id}.${record.format}`
    a.click()
    URL.revokeObjectURL(url)
    Message.success('下载成功')
  } catch (e) {
    Message.error('下载失败')
  }
}

async function deleteExport(record) {
  try {
    await deleteExportRecord(record.id)
    Message.success('删除成功')
    loadExportHistory()
  } catch (e) {
    Message.error('删除失败')
  }
}
</script>

<style scoped>
</style>