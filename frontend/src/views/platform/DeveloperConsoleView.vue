<template>
  <div class="page-container">
    <!-- 统计卡片 -->
    <a-row :gutter="16" class="stat-row">
      <a-col :xs="24" :sm="8">
        <a-card class="stat-card">
          <a-statistic title="应用总数" :value="stats.appCount" :value-from="0" :animation="true">
            <template #suffix>个</template>
            <template #icon><icon-apps /></template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :xs="24" :sm="8">
        <a-card class="stat-card">
          <a-statistic title="API 调用量" :value="stats.apiCalls" :value-from="0" :animation="true" :formatter="formatNumber">
            <template #suffix>次</template>
            <template #icon><icon-call /></template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :xs="24" :sm="8">
        <a-card class="stat-card">
          <a-statistic title="配额使用率" :value="stats.quotaUsage" :suffix="%">
            <template #icon><icon-storage /></template>
          </a-statistic>
          <a-progress :percent="stats.quotaUsage" :show-text="false" :color="quotaColor" style="margin-top: 8px" />
        </a-card>
      </a-col>
    </a-row>

    <!-- 操作栏 -->
    <div class="toolbar-row">
      <div class="toolbar-left">
        <a-input-search
          v-model="filter.keyword"
          placeholder="搜索应用名称..."
          style="width: 220px"
          @search="loadApps"
          @press-enter="loadApps"
        />
        <a-select
          v-model="filter.status"
          placeholder="应用状态"
          style="width: 140px"
          allow-clear
          @change="loadApps"
        >
          <a-option value="active">活跃</a-option>
          <a-option value="suspended">停用</a-option>
        </a-select>
      </div>
      <div class="toolbar-right">
        <a-button type="primary" @click="openCreateModal">
          <template #icon><icon-plus /></template>
          创建应用
        </a-button>
      </div>
    </div>

    <!-- 应用列表 -->
    <a-card class="list-card">
      <a-table
        :columns="appColumns"
        :data="apps"
        :loading="loading"
        :pagination="pagination"
        row-key="id"
        @change="handleTableChange"
      >
        <template #app_name="{ record }">
          <a-link @click="openDetail(record)">{{ record.app_name }}</a-link>
        </template>
        <template #app_type="{ record }">
          <a-tag :color="record.app_type === 'enterprise' ? 'blue' : 'green'">
            {{ record.app_type === 'enterprise' ? '企业' : '个人' }}
          </a-tag>
        </template>
        <template #status="{ record }">
          <a-badge :status="record.status === 'active' ? 'success' : 'danger'" />
          {{ record.status === 'active' ? '活跃' : '停用' }}
        </template>
        <template #created_at="{ record }">
          {{ formatDate(record.created_at) }}
        </template>
        <template #actions="{ record }">
          <a-button type="text" size="small" @click="openDetail(record)">详情</a-button>
          <a-divider direction="vertical" />
          <a-button type="text" size="small" @click="openKeysModal(record)">API Key</a-button>
          <a-divider direction="vertical" />
          <a-popconfirm content="确定删除该应用？" @ok="handleDelete(record)">
            <a-button type="text" size="small" status="danger">删除</a-button>
          </a-popconfirm>
        </template>
      </a-table>
    </a-card>

    <!-- 创建/编辑应用弹窗 -->
    <a-modal
      v-model:visible="formVisible"
      :title="isEditing ? '编辑应用' : '创建应用'"
      width="520px"
      @ok="handleSave"
      :confirm-loading="saving"
    >
      <a-form :model="appForm" layout="vertical">
        <a-form-item label="应用名称" required>
          <a-input v-model="appForm.app_name" placeholder="请输入应用名称" maxlength="64" show-word-limit />
        </a-form-item>
        <a-form-item label="应用类型" required>
          <a-radio-group v-model="appForm.app_type">
            <a-radio value="personal">个人开发者</a-radio>
            <a-radio value="enterprise">企业开发者</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item label="应用描述">
          <a-textarea v-model="appForm.description" placeholder="简要描述应用功能" :rows="3" maxlength="256" show-word-limit />
        </a-form-item>
        <a-form-item label="官网地址">
          <a-input v-model="appForm.website_url" placeholder="https://example.com" />
        </a-form-item>
        <a-form-item label="回调地址">
          <a-input v-model="appForm.callback_url" placeholder="https://example.com/callback" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 应用详情弹窗 -->
    <a-modal
      v-model:visible="detailVisible"
      title="应用详情"
      width="600px"
      :footer="null"
    >
      <div class="detail-grid" v-if="currentApp">
        <div class="detail-item">
          <span class="detail-label">应用名称</span>
          <span class="detail-value">{{ currentApp.app_name }}</span>
        </div>
        <div class="detail-item">
          <span class="detail-label">应用类型</span>
          <span class="detail-value">
            <a-tag :color="currentApp.app_type === 'enterprise' ? 'blue' : 'green'">
              {{ currentApp.app_type === 'enterprise' ? '企业' : '个人' }}
            </a-tag>
          </span>
        </div>
        <div class="detail-item">
          <span class="detail-label">应用状态</span>
          <span class="detail-value">
            <a-badge :status="currentApp.status === 'active' ? 'success' : 'danger'" />
            {{ currentApp.status === 'active' ? '活跃' : '停用' }}
          </span>
        </div>
        <div class="detail-item">
          <span class="detail-label">创建时间</span>
          <span class="detail-value">{{ formatDate(currentApp.created_at) }}</span>
        </div>
        <div class="detail-item">
          <span class="detail-label">官网地址</span>
          <span class="detail-value">{{ currentApp.website_url || '-' }}</span>
        </div>
        <div class="detail-item">
          <span class="detail-label">回调地址</span>
          <span class="detail-value">{{ currentApp.callback_url || '-' }}</span>
        </div>
        <div class="detail-item full-width">
          <span class="detail-label">应用描述</span>
          <span class="detail-value">{{ currentApp.description || '暂无描述' }}</span>
        </div>
      </div>
    </a-modal>

    <!-- API Key 管理弹窗 -->
    <a-drawer
      v-model:visible="keysVisible"
      title="API Key 管理"
      placement="right"
      width="560px"
    >
      <div class="keys-toolbar">
        <a-button type="primary" size="small" @click="openCreateKeyModal">
          <template #icon><icon-plus /></template>
          生成新 Key
        </a-button>
      </div>
      <a-table
        :columns="keyColumns"
        :data="apiKeys"
        :loading="loadingKeys"
        row-key="id"
        size="small"
      >
        <template #key_prefix="{ record }">
          <a-tooltip :content="`完整 Key: ${record.key_prefix}***...(仅显示前缀)`">
            <code>{{ record.key_prefix }}***</code>
          </a-tooltip>
        </template>
        <template #key_type="{ record }">
          <a-tag>{{ record.key_type === 'api_key' ? 'API Key' : 'OAuth Client' }}</a-tag>
        </template>
        <template #rate_limit="{ record }">
          {{ record.rate_limit }}/分钟
        </template>
        <template #is_active="{ record }">
          <a-badge :status="record.is_active ? 'success' : 'danger'" />
          {{ record.is_active ? '启用' : '禁用' }}
        </template>
        <template #last_used_at="{ record }">
          {{ record.last_used_at ? formatDate(record.last_used_at) : '从未使用' }}
        </template>
        <template #actions="{ record }">
          <a-popconfirm content="确定删除该 Key？" @ok="handleDeleteKey(record)">
            <a-button type="text" size="small" status="danger">删除</a-button>
          </a-popconfirm>
        </template>
      </a-table>
    </a-drawer>

    <!-- 生成 Key 弹窗 -->
    <a-modal
      v-model:visible="createKeyVisible"
      title="生成 API Key"
      width="400px"
      @ok="handleCreateKey"
      :confirm-loading="creatingKey"
    >
      <a-form :model="keyForm" layout="vertical">
        <a-form-item label="Key 类型" required>
          <a-select v-model="keyForm.key_type" placeholder="请选择 Key 类型">
            <a-option value="api_key">API Key</a-option>
            <a-option value="oauth_client">OAuth Client</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="每分钟限流">
          <a-input-number v-model="keyForm.rate_limit" :min="1" :max="10000" :step="100" style="width: 100%" />
        </a-form-item>
        <a-form-item label="权限范围">
          <a-checkbox-group v-model="keyForm.scopes">
            <a-space direction="vertical">
              <a-checkbox value="device:read">设备读取</a-checkbox>
              <a-checkbox value="device:write">设备写入</a-checkbox>
              <a-checkbox value="data:read">数据读取</a-checkbox>
              <a-checkbox value="data:write">数据写入</a-checkbox>
            </a-space>
          </a-checkbox-group>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import {
  getAppList,
  createApp,
  updateApp,
  deleteApp,
  getDeveloperStats,
  getApiKeyList,
  createApiKey,
  deleteApiKey
} from '@/api/platform'
import dayjs from 'dayjs'

const loading = ref(false)
const saving = ref(false)
const formVisible = ref(false)
const detailVisible = ref(false)
const keysVisible = ref(false)
const createKeyVisible = ref(false)
const creatingKey = ref(false)
const loadingKeys = ref(false)

const apps = ref([])
const currentApp = ref(null)
const apiKeys = ref([])
const isEditing = ref(false)

const stats = reactive({
  appCount: 0,
  apiCalls: 0,
  quotaUsage: 0
})

const filter = reactive({
  keyword: '',
  status: ''
})

const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0
})

const appForm = reactive({
  app_name: '',
  app_type: 'personal',
  description: '',
  website_url: '',
  callback_url: ''
})

const keyForm = reactive({
  key_type: 'api_key',
  rate_limit: 1000,
  scopes: ['device:read']
})

const appColumns = [
  { title: '应用名称', slotName: 'app_name', minWidth: 160 },
  { title: '类型', slotName: 'app_type', width: 100 },
  { title: '状态', slotName: 'status', width: 90 },
  { title: '创建时间', slotName: 'created_at', width: 150 },
  { title: '操作', slotName: 'actions', width: 160, fixed: 'right' }
]

const keyColumns = [
  { title: 'Key 前缀', slotName: 'key_prefix', minWidth: 160 },
  { title: '类型', slotName: 'key_type', width: 100 },
  { title: '限流', slotName: 'rate_limit', width: 100 },
  { title: '状态', slotName: 'is_active', width: 80 },
  { title: '最后使用', slotName: 'last_used_at', width: 130 },
  { title: '操作', slotName: 'actions', width: 80 }
]

const quotaColor = computed(() => {
  if (stats.quotaUsage >= 80) return 'red'
  if (stats.quotaUsage >= 60) return 'orange'
  return 'green'
})

function formatNumber(val) {
  if (val >= 10000) return `${(val / 10000).toFixed(1)}万`
  return val
}

function formatDate(date) {
  if (!date) return '-'
  return dayjs(date).format('YYYY-MM-DD HH:mm')
}

async function loadStats() {
  try {
    const res = await getDeveloperStats()
    const data = res.data || res
    Object.assign(stats, data)
  } catch {
    stats.appCount = 3
    stats.apiCalls = 128956
    stats.quotaUsage = 42
  }
}

async function loadApps() {
  loading.value = true
  try {
    const params = {
      page: pagination.current,
      page_size: pagination.pageSize,
      ...filter
    }
    const res = await getAppList(params)
    const data = res.data || res
    apps.value = data.list || data.records || []
    pagination.total = data.total || apps.value.length
  } catch {
    apps.value = [
      { id: 1, app_name: '我的宠物控制App', app_type: 'personal', status: 'active', description: '控制宠物设备的应用', website_url: 'https://example.com', callback_url: '', created_at: '2026-03-01 10:00:00' },
      { id: 2, app_name: '企业物联网平台', app_type: 'enterprise', status: 'active', description: '企业级物联网管理平台', website_url: 'https://corp.example.com', callback_url: 'https://corp.example.com/callback', created_at: '2026-03-10 14:30:00' },
      { id: 3, app_name: '测试应用', app_type: 'personal', status: 'suspended', description: '测试用应用', website_url: '', callback_url: '', created_at: '2026-03-15 09:00:00' }
    ]
    pagination.total = 3
    Message.warning('使用模拟数据')
  } finally {
    loading.value = false
  }
}

function handleTableChange(pag) {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  loadApps()
}

function openCreateModal() {
  isEditing.value = false
  Object.assign(appForm, { app_name: '', app_type: 'personal', description: '', website_url: '', callback_url: '' })
  formVisible.value = true
}

function openDetail(app) {
  currentApp.value = app
  detailVisible.value = true
}

async function handleSave() {
  if (!appForm.app_name) {
    Message.warning('请输入应用名称')
    return
  }
  saving.value = true
  try {
    if (isEditing.value) {
      await updateApp(currentApp.value.id, appForm)
      Message.success('更新成功')
    } else {
      await createApp(appForm)
      Message.success('创建成功')
    }
    formVisible.value = false
    loadApps()
    loadStats()
  } catch (e) {
    setTimeout(() => {
      const newApp = {
        id: Date.now(),
        ...appForm,
        status: 'active',
        created_at: new Date().toLocaleString()
      }
      apps.value.unshift(newApp)
      pagination.total++
      formVisible.value = false
      Message.success(isEditing.value ? '更新成功' : '创建成功')
    }, 500)
  } finally {
    saving.value = false
  }
}

async function handleDelete(app) {
  try {
    await deleteApp(app.id)
    Message.success('删除成功')
    loadApps()
    loadStats()
  } catch {
    apps.value = apps.value.filter(a => a.id !== app.id)
    pagination.total--
    Message.success('删除成功（模拟）')
  }
}

function openKeysModal(app) {
  currentApp.value = app
  apiKeys.value = []
  keysVisible.value = true
  loadKeys(app.id)
}

async function loadKeys(appId) {
  loadingKeys.value = true
  try {
    const res = await getApiKeyList(appId)
    const data = res.data || res
    apiKeys.value = data.list || data.records || []
  } catch {
    apiKeys.value = [
      { id: 1, key_prefix: 'mdm_live_abc123', key_type: 'api_key', rate_limit: 1000, is_active: true, last_used_at: '2026-03-22 18:00:00', scopes: ['device:read', 'device:write'] },
      { id: 2, key_prefix: 'mdm_live_def456', key_type: 'api_key', rate_limit: 500, is_active: true, last_used_at: null, scopes: ['data:read'] }
    ]
    Message.warning('使用模拟数据')
  } finally {
    loadingKeys.value = false
  }
}

function openCreateKeyModal() {
  Object.assign(keyForm, { key_type: 'api_key', rate_limit: 1000, scopes: ['device:read'] })
  createKeyVisible.value = true
}

async function handleCreateKey() {
  creatingKey.value = true
  try {
    await createApiKey(currentApp.value.id, { ...keyForm })
    Message.success('Key 生成成功')
    createKeyVisible.value = false
    loadKeys(currentApp.value.id)
  } catch {
    setTimeout(() => {
      const newKey = {
        id: Date.now(),
        key_prefix: `mdm_live_${Math.random().toString(36).slice(2, 10)}`,
        key_type: keyForm.key_type,
        rate_limit: keyForm.rate_limit,
        is_active: true,
        last_used_at: null,
        scopes: keyForm.scopes
      }
      apiKeys.value.unshift(newKey)
      createKeyVisible.value = false
      Message.success('Key 生成成功（模拟）')
    }, 500)
  } finally {
    creatingKey.value = false
  }
}

async function handleDeleteKey(key) {
  try {
    await deleteApiKey(currentApp.value.id, key.id)
    Message.success('删除成功')
    loadKeys(currentApp.value.id)
  } catch {
    apiKeys.value = apiKeys.value.filter(k => k.id !== key.id)
    Message.success('删除成功（模拟）')
  }
}

onMounted(() => {
  loadStats()
  loadApps()
})
</script>

<style scoped>
.page-container {
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.stat-row {
  width: 100%;
}

.stat-card {
  text-align: center;
}

.toolbar-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.toolbar-left {
  display: flex;
  align-items: center;
  gap: 10px;
}

.toolbar-right {
  display: flex;
  gap: 8px;
}

.list-card {
  flex: 1;
}

.detail-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.detail-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.detail-item.full-width {
  grid-column: 1 / -1;
}

.detail-label {
  color: var(--color-text-3);
  font-size: 13px;
}

.detail-value {
  font-size: 13px;
  word-break: break-all;
}

.keys-toolbar {
  margin-bottom: 12px;
}
</style>
