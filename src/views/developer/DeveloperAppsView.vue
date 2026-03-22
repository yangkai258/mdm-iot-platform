<template>
  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>开发者平台</a-breadcrumb-item>
      <a-breadcrumb-item>应用管理</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 标签页 -->
    <a-tabs v-model:active-key="activeTab" class="pro-tabs">
      <!-- 应用列表 -->
      <a-tab-pane key="list" title="应用列表">
        <div class="pro-search-bar">
          <a-space>
            <a-input-search v-model="searchKeyword" placeholder="搜索应用名称或ID" style="width: 260px" @search="loadApps" search-button />
            <a-select v-model="filterStatus" placeholder="状态筛选" style="width: 140px" allow-clear>
              <a-option :value="1">启用</a-option>
              <a-option :value="0">禁用</a-option>
            </a-select>
          </a-space>
        </div>
        <div class="pro-action-bar">
          <a-space>
            <a-button type="primary" @click="showCreateModal">创建应用</a-button>
            <a-button @click="loadApps">刷新</a-button>
          </a-space>
        </div>
        <div class="pro-content-area">
          <a-table :columns="appColumns" :data="filteredApps" :loading="loading" :pagination="pagination" @change="handleTableChange" row-key="id">
            <template #logo="{ record }">
              <a-avatar :size="36" :style="{ backgroundColor: record.color || '#165dff' }">
                {{ record.name?.charAt(0) || 'A' }}
              </a-avatar>
            </template>
            <template #status="{ record }">
              <a-tag :color="record.status === 1 ? 'green' : 'gray'">{{ record.status === 1 ? '启用' : '禁用' }}</a-tag>
            </template>
            <template #createdAt="{ record }">
              {{ formatTime(record.created_at) }}
            </template>
            <template #actions="{ record }">
              <a-space>
                <a-button type="text" size="small" @click="viewApp(record)">详情</a-button>
                <a-button type="text" size="small" @click="editApp(record)">编辑</a-button>
                <a-button type="text" size="small" status="danger" @click="deleteApp(record)">删除</a-button>
              </a-space>
            </template>
          </a-table>
        </div>
      </a-tab-pane>

      <!-- API Key 管理 -->
      <a-tab-pane key="apikeys" title="API Key 管理">
        <div class="pro-search-bar">
          <a-space>
            <a-input-search v-model="searchKeyword" placeholder="搜索应用名称或Key ID" style="width: 280px" @search="loadApiKeys" search-button />
            <a-select v-model="filterKeyStatus" placeholder="Key 状态" style="width: 140px" allow-clear>
              <a-option :value="1">活跃</a-option>
              <a-option :value="0">已禁用</a-option>
            </a-select>
          </a-space>
        </div>
        <div class="pro-action-bar">
          <a-space>
            <a-button type="primary" @click="showCreateKeyModal">创建 API Key</a-button>
            <a-button @click="loadApiKeys">刷新</a-button>
          </a-space>
        </div>
        <div class="pro-content-area">
          <a-table :columns="keyColumns" :data="filteredKeys" :loading="keyLoading" :pagination="keyPagination" @change="handleKeyTableChange" row-key="id">
            <template #keyId="{ record }">
              <a-tag color="arcoblue">{{ record.key_id }}</a-tag>
            </template>
            <template #appName="{ record }">
              <a-tag>{{ record.app_name }}</a-tag>
            </template>
            <template #scopes="{ record }">
              <a-space wrap>
                <a-tag v-for="s in record.scopes" :key="s" size="small">{{ s }}</a-tag>
              </a-space>
            </template>
            <template #expiresAt="{ record }">
              <span :style="{ color: isExpiringSoon(record.expires_at) ? '#ff4d4f' : '' }">
                {{ record.expires_at ? formatTime(record.expires_at) : '永不过期' }}
              </span>
            </template>
            <template #lastUsed="{ record }">
              {{ record.last_used_at ? formatTime(record.last_used_at) : '从未使用' }}
            </template>
            <template #status="{ record }">
              <a-tag :color="record.status === 1 ? 'green' : 'gray'">{{ record.status === 1 ? '活跃' : '已禁用' }}</a-tag>
            </template>
            <template #actions="{ record }">
              <a-space>
                <a-button type="text" size="small" @click="copyKey(record)">复制</a-button>
                <a-button type="text" size="small" @click="regenerateKey(record)">重置</a-button>
                <a-button type="text" size="small" status="danger" @click="revokeKey(record)">撤销</a-button>
              </a-space>
            </template>
          </a-table>
        </div>
      </a-tab-pane>
    </a-tabs>

    <!-- 创建/编辑应用弹窗 -->
    <a-modal
      v-model:visible="appModalVisible"
      :title="editingApp ? '编辑应用' : '创建应用'"
      @ok="handleSaveApp"
      :confirm-loading="appSaving"
      :width="560"
    >
      <a-form :model="appForm" layout="vertical">
        <a-form-item label="应用名称" required>
          <a-input v-model="appForm.name" placeholder="请输入应用名称" :max-length="64" show-word-limit />
        </a-form-item>
        <a-form-item label="应用描述">
          <a-textarea v-model="appForm.description" placeholder="请输入应用描述（可选）" :max-length="256" show-word-limit :rows="3" />
        </a-form-item>
        <a-form-item label="应用图标颜色">
          <a-color-picker v-model="appForm.color" />
        </a-form-item>
        <a-form-item label="应用类型">
          <a-select v-model="appForm.app_type" placeholder="选择应用类型">
            <a-option value="web">Web 应用</a-option>
            <a-option value="mobile">移动应用</a-option>
            <a-option value="desktop">桌面应用</a-option>
            <a-option value="iot">IoT 设备</a-option>
            <a-option value="service">服务端服务</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="回调地址">
          <a-input v-model="appForm.redirect_uri" placeholder="https://example.com/callback" />
        </a-form-item>
        <a-form-item label="Homepage URL">
          <a-input v-model="appForm.homepage_url" placeholder="https://example.com" />
        </a-form-item>
        <a-form-item label="状态">
          <a-switch v-model="appForm.status" :checked-value="1" :unchecked-value="0" />
          <span style="margin-left: 8px">{{ appForm.status === 1 ? '启用' : '禁用' }}</span>
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 应用详情抽屉 -->
    <a-drawer v-model:visible="detailDrawerVisible" :width="560" :title="`应用详情: ${detailApp?.name}`" @close="detailDrawerVisible = false">
      <template v-if="detailApp">
        <a-descriptions :column="1" bordered size="small">
          <a-descriptions-item label="应用ID">{{ detailApp.id }}</a-descriptions-item>
          <a-descriptions-item label="应用名称">{{ detailApp.name }}</a-descriptions-item>
          <a-descriptions-item label="应用类型">{{ detailApp.app_type || '-' }}</a-descriptions-item>
          <a-descriptions-item label="描述">{{ detailApp.description || '-' }}</a-descriptions-item>
          <a-descriptions-item label="Client ID">{{ detailApp.client_id }}</a-descriptions-item>
          <a-descriptions-item label="Client Secret">
            <a-space>
              <span>••••••••</span>
              <a-button type="text" size="small" @click="showSecret(detailApp)">显示</a-button>
              <a-button type="text" size="small" @click="copyText(detailApp.client_secret)">复制</a-button>
            </a-space>
          </a-descriptions-item>
          <a-descriptions-item label="回调地址">{{ detailApp.redirect_uri || '-' }}</a-descriptions-item>
          <a-descriptions-item label="Homepage">{{ detailApp.homepage_url || '-' }}</a-descriptions-item>
          <a-descriptions-item label="状态">
            <a-tag :color="detailApp.status === 1 ? 'green' : 'gray'">{{ detailApp.status === 1 ? '启用' : '禁用' }}</a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="创建时间">{{ formatTime(detailApp.created_at) }}</a-descriptions-item>
          <a-descriptions-item label="更新时间">{{ formatTime(detailApp.updated_at) }}</a-descriptions-item>
        </a-descriptions>

        <a-divider>关联 API Keys</a-divider>
        <a-table :columns="keyColumns" :data="detailApp.api_keys || []" :pagination="false" row-key="id" size="small">
          <template #keyId="{ record }"><a-tag color="arcoblue">{{ record.key_id }}</a-tag></template>
          <template #appName="{ record }"><a-tag>{{ record.app_name }}</a-tag></template>
          <template #scopes="{ record }">
            <a-space wrap>
              <a-tag v-for="s in record.scopes" :key="s" size="small">{{ s }}</a-tag>
            </a-space>
          </template>
          <template #status="{ record }">
            <a-tag :color="record.status === 1 ? 'green' : 'gray'">{{ record.status === 1 ? '活跃' : '已禁用' }}</a-tag>
          </template>
          <template #expiresAt="{ record }">
            <span :style="{ color: isExpiringSoon(record.expires_at) ? '#ff4d4f' : '' }">
              {{ record.expires_at ? formatTime(record.expires_at) : '永不过期' }}
            </span>
          </template>
        </a-table>

        <a-divider>应用操作</a-divider>
        <a-space direction="vertical" style="width: 100%">
          <a-button type="primary" long @click="editApp(detailApp)">编辑应用信息</a-button>
          <a-button type="dashed" long @click="showCreateKeyForApp(detailApp)">为该应用创建 API Key</a-button>
        </a-space>
      </template>
    </a-drawer>

    <!-- 创建 API Key 弹窗 -->
    <a-modal v-model:visible="keyModalVisible" title="创建 API Key" @ok="handleSaveKey" :confirm-loading="keySaving" :width="520">
      <a-form :model="keyForm" layout="vertical">
        <a-form-item label="关联应用" required>
          <a-select v-model="keyForm.app_id" placeholder="选择应用" :disabled="!!keyForm.app_id" @change="onAppSelect">
            <a-option v-for="app in apps" :key="app.id" :value="app.id">{{ app.name }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="Key 名称">
          <a-input v-model="keyForm.name" placeholder="例如: 生产环境 Key" :max-length="64" show-word-limit />
        </a-form-item>
        <a-form-item label="权限范围">
          <a-checkbox-group v-model="keyForm.scopes">
            <a-space direction="vertical" style="width: 100%">
              <a-checkbox value="device:read">设备读取 (device:read)</a-checkbox>
              <a-checkbox value="device:write">设备写入 (device:write)</a-checkbox>
              <a-checkbox value="data:read">数据读取 (data:read)</a-checkbox>
              <a-checkbox value="data:write">数据写入 (data:write)</a-checkbox>
              <a-checkbox value="analytics:read">分析读取 (analytics:read)</a-checkbox>
              <a-checkbox value="admin">管理权限 (admin)</a-checkbox>
            </a-space>
          </a-checkbox-group>
        </a-form-item>
        <a-form-item label="过期时间">
          <a-select v-model="keyForm.expires_in" placeholder="选择过期时间">
            <a-option :value="0">永不过期</a-option>
            <a-option :value="30">30 天</a-option>
            <a-option :value="90">90 天</a-option>
            <a-option :value="180">180 天</a-option>
            <a-option :value="365">1 年</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="描述">
          <a-textarea v-model="keyForm.description" placeholder="该 Key 的用途说明" :max-length="256" show-word-limit :rows="2" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- API Key 展示弹窗 -->
    <a-modal v-model:visible="showKeyModalVisible" title="API Key 已创建" :width="520">
      <a-result status="success" title="请立即复制并妥善保管" subtitle="此 Key 仅显示一次，后续无法找回">
        <template #content>
          <div style="background: #f2f3f5; border-radius: 8px; padding: 16px; margin-bottom: 16px">
            <div style="margin-bottom: 8px; color: #86909c; font-size: 12px">Key ID</div>
            <a-tag color="arcoblue" style="font-size: 14px; margin-bottom: 12px">{{ newKeyData.key_id }}</a-tag>
            <div style="margin-bottom: 8px; color: #86909c; font-size: 12px">API Key (Secret)</div>
            <a-input-password v-model="newKeyData.api_key" readonly style="font-family: monospace; font-size: 14px" />
          </div>
          <div style="color: #ff4d4f; font-size: 12px">⚠️ 请立即复制并保存，关闭此窗口后将无法再次查看</div>
        </template>
      </a-result>
      <template #footer>
        <a-button type="primary" @click="copyNewKey">复制 Key</a-button>
        <a-button @click="showKeyModalVisible = false">关闭</a-button>
      </template>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { Message, Modal } from '@arco-design/web-vue'
import axios from 'axios'

const API_BASE = '/api/v1/developer'
const activeTab = ref('list')
const loading = ref(false)
const keyLoading = ref(false)
const appSaving = ref(false)
const keySaving = ref(false)
const searchKeyword = ref('')
const filterStatus = ref(null)
const filterKeyStatus = ref(null)

// 应用相关
const apps = ref([])
const appModalVisible = ref(false)
const detailDrawerVisible = ref(false)
const editingApp = ref(null)
const detailApp = ref(null)
const appForm = reactive({
  name: '',
  description: '',
  color: '#165dff',
  app_type: 'web',
  redirect_uri: '',
  homepage_url: '',
  status: 1
})

const pagination = reactive({ current: 1, pageSize: 20, total: 0 })

const appColumns = [
  { title: 'Logo', slotName: 'logo', width: 70 },
  { title: '应用名称', dataIndex: 'name', width: 180 },
  { title: '应用类型', dataIndex: 'app_type', width: 120 },
  { title: 'Client ID', dataIndex: 'client_id', width: 200 },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '创建时间', slotName: 'createdAt', width: 180 },
  { title: '操作', slotName: 'actions', width: 200, fixed: 'right' }
]

// API Key 相关
const apiKeys = ref([])
const keyModalVisible = ref(false)
const showKeyModalVisible = ref(false)
const keyForm = reactive({
  app_id: null,
  name: '',
  scopes: ['device:read'],
  expires_in: 365,
  description: ''
})
const newKeyData = reactive({ key_id: '', api_key: '' })
const keyPagination = reactive({ current: 1, pageSize: 20, total: 0 })

const keyColumns = [
  { title: 'Key ID', slotName: 'keyId', width: 200 },
  { title: '关联应用', slotName: 'appName', width: 140 },
  { title: '名称', dataIndex: 'name', width: 140 },
  { title: '权限范围', slotName: 'scopes', width: 240 },
  { title: '过期时间', slotName: 'expiresAt', width: 160 },
  { title: '最后使用', slotName: 'lastUsed', width: 160 },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '操作', slotName: 'actions', width: 200, fixed: 'right' }
]

const filteredApps = computed(() => {
  let result = apps.value
  if (filterStatus.value !== null) result = result.filter(a => a.status === filterStatus.value)
  if (searchKeyword.value) {
    const kw = searchKeyword.value.toLowerCase()
    result = result.filter(a => a.name?.toLowerCase().includes(kw) || a.client_id?.toLowerCase().includes(kw) || String(a.id)?.includes(kw))
  }
  return result
})

const filteredKeys = computed(() => {
  let result = apiKeys.value
  if (filterKeyStatus.value !== null) result = result.filter(k => k.status === filterKeyStatus.value)
  if (searchKeyword.value) {
    const kw = searchKeyword.value.toLowerCase()
    result = result.filter(k => k.key_id?.toLowerCase().includes(kw) || k.app_name?.toLowerCase().includes(kw))
  }
  return result
})

const loadApps = async () => {
  loading.value = true
  try {
    const res = await axios.get(`${API_BASE}/apps`)
    if (res.data.code === 0) {
      apps.value = res.data.data.list || []
      pagination.total = res.data.data.pagination?.total || apps.value.length
    }
  } catch {
    apps.value = getMockApps()
    pagination.total = apps.value.length
    Message.warning('使用模拟数据')
  } finally { loading.value = false }
}

const loadApiKeys = async () => {
  keyLoading.value = true
  try {
    const res = await axios.get(`${API_BASE}/api-keys`)
    if (res.data.code === 0) {
      apiKeys.value = res.data.data.list || []
      keyPagination.total = res.data.data.pagination?.total || apiKeys.value.length
    }
  } catch {
    apiKeys.value = getMockApiKeys()
    keyPagination.total = apiKeys.value.length
    Message.warning('使用模拟数据')
  } finally { keyLoading.value = false }
}

const getMockApps = () => [
  { id: 1, name: 'PetCare App', app_type: 'mobile', client_id: 'cli_a1b2c3d4e5f6', color: '#165dff', description: '宠物健康监测应用', status: 1, created_at: '2026-03-01 10:00:00', updated_at: '2026-03-20 14:00:00', api_keys: [] },
  { id: 2, name: 'HomeDashboard', app_type: 'web', client_id: 'cli_b2c3d4e5f6g7', color: '#00d4ff', description: '家庭设备管理后台', status: 1, created_at: '2026-03-05 11:00:00', updated_at: '2026-03-18 09:00:00', api_keys: [] },
  { id: 3, name: 'EdgeGateway', app_type: 'iot', client_id: 'cli_c3d4e5f6g7h8', color: '#52c41a', description: '边缘网关服务', status: 1, created_at: '2026-03-10 08:00:00', updated_at: '2026-03-22 10:00:00', api_keys: [] },
  { id: 4, name: 'DataAnalytics', app_type: 'service', client_id: 'cli_d4e5f6g7h8i9', color: '#722ed1', description: '数据分析服务', status: 0, created_at: '2026-03-12 15:00:00', updated_at: '2026-03-19 16:00:00', api_keys: [] }
]

const getMockApiKeys = () => [
  { id: 1, key_id: 'ak_prod_xxxxxxxxxxxx', app_id: 1, app_name: 'PetCare App', name: '生产环境 Key', scopes: ['device:read', 'device:write', 'data:read'], expires_at: '2027-03-22 00:00:00', last_used_at: '2026-03-22 10:30:00', status: 1 },
  { id: 2, key_id: 'ak_test_yyyyyyyyyyyy', app_id: 1, app_name: 'PetCare App', name: '测试环境 Key', scopes: ['device:read', 'data:read'], expires_at: '2026-06-22 00:00:00', last_used_at: '2026-03-21 18:00:00', status: 1 },
  { id: 3, key_id: 'ak_home_zzzzzzzzzzzz', app_id: 2, app_name: 'HomeDashboard', name: 'Home Key', scopes: ['device:read', 'device:write', 'data:read', 'data:write'], expires_at: null, last_used_at: '2026-03-22 11:00:00', status: 1 },
  { id: 4, key_id: 'ak_edge_wwwwwwwwwwww', app_id: 3, app_name: 'EdgeGateway', name: 'Edge Key', scopes: ['device:read', 'analytics:read'], expires_at: '2026-04-22 00:00:00', last_used_at: null, status: 1 },
  { id: 5, key_id: 'ak_old_uuuuuuuuuuuu', app_id: 2, app_name: 'HomeDashboard', name: '旧版 Key (已过期)', scopes: ['device:read'], expires_at: '2026-01-01 00:00:00', last_used_at: '2026-01-01 00:00:00', status: 0 }
]

const showCreateModal = () => {
  editingApp.value = null
  Object.assign(appForm, { name: '', description: '', color: '#165dff', app_type: 'web', redirect_uri: '', homepage_url: '', status: 1 })
  appModalVisible.value = true
}

const editApp = (record) => {
  editingApp.value = record
  Object.assign(appForm, {
    name: record.name,
    description: record.description || '',
    color: record.color || '#165dff',
    app_type: record.app_type || 'web',
    redirect_uri: record.redirect_uri || '',
    homepage_url: record.homepage_url || '',
    status: record.status
  })
  appModalVisible.value = true
}

const viewApp = async (record) => {
  detailApp.value = { ...record }
  try {
    const res = await axios.get(`${API_BASE}/apps/${record.id}/api-keys`)
    if (res.data.code === 0) detailApp.value.api_keys = res.data.data.list || []
  } catch { detailApp.value.api_keys = [] }
  detailDrawerVisible.value = true
}

const handleSaveApp = () => {
  if (!appForm.name?.trim()) { Message.warning('请填写应用名称'); return }
  appSaving.value = true
  const payload = { ...appForm }
  const apiCall = editingApp.value
    ? axios.put(`${API_BASE}/apps/${editingApp.value.id}`, payload)
    : axios.post(`${API_BASE}/apps`, payload)

  apiCall.then(res => {
    if (res.data.code === 0) {
      Message.success(editingApp.value ? '应用已更新' : '应用已创建')
      appSaving.value = false
      appModalVisible.value = false
      loadApps()
    } else {
      Message.error(res.data.message || '保存失败')
      appSaving.value = false
    }
  }).catch(() => {
    setTimeout(() => {
      if (editingApp.value) {
        const idx = apps.value.findIndex(a => a.id === editingApp.value.id)
        if (idx !== -1) Object.assign(apps.value[idx], { ...editingApp.value, ...appForm })
        Message.success('应用已更新 (模拟)')
      } else {
        apps.value.unshift({ id: Date.now(), ...appForm, client_id: 'cli_' + Math.random().toString(36).slice(2, 14), created_at: new Date().toLocaleString(), updated_at: new Date().toLocaleString() })
        Message.success('应用已创建 (模拟)')
      }
      appSaving.value = false
      appModalVisible.value = false
    }, 600)
  })
}

const deleteApp = (record) => {
  Modal.confirm({ title: '确认删除', content: `确定要删除应用「${record.name}」吗？此操作不可恢复。`, onOk: () => {
    return axios.delete(`${API_BASE}/apps/${record.id}`).then(res => {
      if (res.data.code === 0) { apps.value = apps.value.filter(a => a.id !== record.id); Message.success('应用已删除') }
    }).catch(() => {
      apps.value = apps.value.filter(a => a.id !== record.id)
      Message.success('应用已删除 (模拟)')
    })
  }})
}

const showCreateKeyModal = () => {
  keyForm.app_id = null
  keyForm.name = ''
  keyForm.scopes = ['device:read']
  keyForm.expires_in = 365
  keyForm.description = ''
  keyModalVisible.value = true
}

const showCreateKeyForApp = (app) => {
  detailDrawerVisible.value = false
  showCreateKeyModal()
  keyForm.app_id = app.id
}

const onAppSelect = (val) => { keyForm.app_id = val }

const handleSaveKey = () => {
  if (!keyForm.app_id) { Message.warning('请选择应用'); return }
  keySaving.value = true
  axios.post(`${API_BASE}/api-keys`, { ...keyForm }).then(res => {
    if (res.data.code === 0) {
      newKeyData.key_id = res.data.data.key_id
      newKeyData.api_key = res.data.data.api_key
      Message.success('API Key 已创建')
      keySaving.value = false
      keyModalVisible.value = false
      showKeyModalVisible.value = true
      loadApiKeys()
    } else {
      Message.error(res.data.message || '创建失败')
      keySaving.value = false
    }
  }).catch(() => {
    setTimeout(() => {
      newKeyData.key_id = 'ak_' + Math.random().toString(36).slice(2, 14)
      newKeyData.api_key = 'sk_' + Math.random().toString(36).slice(2) + Math.random().toString(36).slice(2)
      keySaving.value = false
      keyModalVisible.value = false
      showKeyModalVisible.value = true
      Message.warning('使用模拟数据')
      loadApiKeys()
    }, 600)
  })
}

const copyNewKey = () => {
  navigator.clipboard.writeText(`${newKeyData.key_id}:${newKeyData.api_key}`)
  Message.success('已复制到剪贴板')
}

const copyKey = (record) => {
  navigator.clipboard.writeText(record.key_id)
  Message.success('Key ID 已复制')
}

const copyText = (text) => {
  navigator.clipboard.writeText(text || '')
  Message.success('已复制')
}

const showSecret = (app) => {
  Message.info(`Client Secret: ${app.client_secret || '模拟 Secret'}`)
}

const regenerateKey = (record) => {
  Modal.confirm({ title: '确认重置', content: `确定要重置 Key 「${record.name}」吗？旧 Key 将立即失效。`, onOk: () => {
    axios.post(`${API_BASE}/api-keys/${record.id}/regenerate`).then(res => {
      if (res.data.code === 0) { Message.success('Key 已重置'); loadApiKeys() }
    }).catch(() => { Message.success('Key 已重置 (模拟)') })
  }})
}

const revokeKey = (record) => {
  Modal.confirm({ title: '确认撤销', content: `确定要撤销 Key 「${record.name}」吗？此操作不可恢复。`, onOk: () => {
    axios.delete(`${API_BASE}/api-keys/${record.id}`).then(res => {
      if (res.data.code === 0) { apiKeys.value = apiKeys.value.filter(k => k.id !== record.id); Message.success('Key 已撤销') }
    }).catch(() => { apiKeys.value = apiKeys.value.filter(k => k.id !== record.id); Message.success('Key 已撤销 (模拟)') })
  }})
}

const isExpiringSoon = (expiresAt) => {
  if (!expiresAt) return false
  const days = (new Date(expiresAt) - new Date()) / (1000 * 60 * 60 * 24)
  return days > 0 && days < 30
}

const formatTime = (t) => t || '-'
const handleTableChange = (pag) => { pagination.current = pag.current }
const handleKeyTableChange = (pag) => { keyPagination.current = pag.current }

onMounted(() => { loadApps(); loadApiKeys() })
</script>

<style scoped>
.pro-page-container { padding: 20px 24px; min-height: calc(100vh - 64px); background: #f5f7fa; }
.pro-breadcrumb { margin-bottom: 16px; }
.pro-tabs { background: #fff; border-radius: 8px; padding: 16px; box-shadow: 0 1px 3px rgba(0,0,0,0.04); }
.pro-search-bar { margin-bottom: 12px; }
.pro-action-bar { margin-bottom: 16px; }
.pro-content-area { background: #fff; border-radius: 8px; padding: 20px; box-shadow: 0 1px 3px rgba(0,0,0,0.04); }
</style>
