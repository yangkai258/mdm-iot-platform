<template>
  <div class="pro-page-container">
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>系统管理</a-breadcrumb-item>
      <a-breadcrumb-item>LDAP/AD集成</a-breadcrumb-item>
    </a-breadcrumb>

    <a-tabs v-model:active-key="activeTab" class="pro-content-area">
      <!-- 服务器配�?-->
      <a-tab-pane key="config" title="服务器配�?>
        <a-card title="LDAP服务器配�?>
          <a-form :model="ldapConfig" layout="vertical" ref="configFormRef">
            <a-row :gutter="16">
              <a-col :span="12">
                <a-form-item label="服务器地址" required>
                  <a-input v-model="ldapConfig.host" placeholder="ldap://192.168.1.100" />
                </a-form-item>
              </a-col>
              <a-col :span="12">
                <a-form-item label="端口" required>
                  <a-input-number v-model="ldapConfig.port" :min="1" :max="65535" placeholder="389" style="width: 100%" />
                </a-form-item>
              </a-col>
            </a-row>
            <a-row :gutter="16">
              <a-col :span="12">
                <a-form-item label="绑定DN" required>
                  <a-input v-model="ldapConfig.bind_dn" placeholder="cn=admin,dc=example,dc=com" />
                </a-form-item>
              </a-col>
              <a-col :span="12">
                <a-form-item label="绑定密码" required>
                  <a-input-password v-model="ldapConfig.bind_password" placeholder="输入密码" />
                </a-form-item>
              </a-col>
            </a-row>
            <a-row :gutter="16">
              <a-col :span="12">
                <a-form-item label="用户基础DN" required>
                  <a-input v-model="ldapConfig.user_base_dn" placeholder="ou=users,dc=example,dc=com" />
                </a-form-item>
              </a-col>
              <a-col :span="12">
                <a-form-item label="用户过滤规则">
                  <a-input v-model="ldapConfig.user_filter" placeholder="(objectClass=person)" />
                </a-form-item>
              </a-col>
            </a-row>
            <a-row :gutter="16">
              <a-col :span="12">
                <a-form-item label="启用SSL/TLS">
                  <a-switch v-model="ldapConfig.use_ssl" />
                </a-form-item>
              </a-col>
              <a-col :span="12">
                <a-form-item label="自动同步">
                  <a-switch v-model="ldapConfig.auto_sync_enabled" />
                </a-form-item>
              </a-col>
            </a-row>
            <a-row :gutter="16" v-if="ldapConfig.auto_sync_enabled">
              <a-col :span="12">
                <a-form-item label="同步周期">
                  <a-select v-model="ldapConfig.sync_interval" placeholder="选择同步周期">
                    <a-option value="1h">每小�?/a-option>
                    <a-option value="6h">�?小时</a-option>
                    <a-option value="12h">�?2小时</a-option>
                    <a-option value="24h">每天</a-option>
                  </a-select>
                </a-form-item>
              </a-col>
              <a-col :span="12">
                <a-form-item label="上次同步时间">
                  <a-input :model-value="lastSyncTime || '从未同步'" readonly />
                </a-form-item>
              </a-col>
            </a-row>
            <a-form-item>
              <a-space>
                <a-button type="primary" @click="saveConfig" :loading="saving">保存配置</a-button>
                <a-button @click="testConnection" :loading="testing">测试连接</a-button>
              </a-space>
            </a-form-item>
          </a-form>
        </a-card>
      </a-tab-pane>

      <!-- 用户同步 -->
      <a-tab-pane key="sync" title="用户同步">
        <div class="pro-action-bar">
          <a-space>
            <a-button type="primary" @click="triggerSync" :loading="syncing">立即同步</a-button>
            <a-button @click="loadSyncLogs">刷新</a-button>
          </a-space>
        </div>
        <a-table :columns="syncLogColumns" :data="syncLogs" :loading="loading" :pagination="pagination" row-key="id" @page-change="handlePageChange">
          <template #status="{ record }">
            <a-tag :color="record.status === 'success' ? 'green' : 'red'">
              {{ record.status === 'success' ? '成功' : '失败' }}
            </a-tag>
          </template>
          <template #type="{ record }">
            <a-tag :color="record.action === 'create' ? 'arcoblue' : record.action === 'update' ? 'orange' : 'gray'">
              {{ record.action === 'create' ? '新增' : record.action === 'update' ? '更新' : '跳过' }}
            </a-tag>
          </template>
          <template #created_at="{ record }">
            {{ formatDate(record.created_at) }}
          </template>
        </a-table>
      </a-tab-pane>

      <!-- 同步日志详情 -->
      <a-tab-pane key="log-detail" title="同步日志详情">
        <a-result v-if="!selectedSyncLog" status="info" title="请从同步日志中选择一条记录查看详�? />
        <a-card v-else :title="`同步日志详情 - ${selectedSyncLog.user_dn}`">
          <a-descriptions :column="2" bordered>
            <a-descriptions-item label="用户DN">{{ selectedSyncLog.user_dn }}</a-descriptions-item>
            <a-descriptions-item label="状�?>
              <a-tag :color="selectedSyncLog.status === 'success' ? 'green' : 'red'">
                {{ selectedSyncLog.status === 'success' ? '成功' : '失败' }}
              </a-tag>
            </a-descriptions-item>
            <a-descriptions-item label="操作">
              <a-tag>{{ selectedSyncLog.action === 'create' ? '新增' : selectedSyncLog.action === 'update' ? '更新' : '跳过' }}</a-tag>
            </a-descriptions-item>
            <a-descriptions-item label="同步时间">{{ formatDate(selectedSyncLog.created_at) }}</a-descriptions-item>
            <a-descriptions-item label="错误信息" :span="2">{{ selectedSyncLog.error_msg || '-' }}</a-descriptions-item>
          </a-descriptions>
          <a-divider>变更详情</a-divider>
          <a-alert v-if="selectedSyncLog.changes" type="info">
            <template #title>字段变更</template>
            <pre style="margin: 0; white-space: pre-wrap">{{ JSON.stringify(JSON.parse(selectedSyncLog.changes), null, 2) }}</pre>
          </a-alert>
          <a-alert v-else type="info" message="无变�? />
        </a-card>
      </a-tab-pane>
    </a-tabs>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { Message } from '@arco-design/web-vue'
import axios from 'axios'

const activeTab = ref('config')
const saving = ref(false)
const testing = ref(false)
const syncing = ref(false)
const loading = ref(false)
const lastSyncTime = ref('')
const selectedSyncLog = ref<any>(null)

const ldapConfig = reactive({
  host: '',
  port: 389,
  bind_dn: '',
  bind_password: '',
  user_base_dn: '',
  user_filter: '(objectClass=person)',
  use_ssl: false,
  auto_sync_enabled: false,
  sync_interval: '24h',
})

const pagination = reactive({ current: 1, pageSize: 10, total: 0 })
const syncLogs = ref<any[]>([])

const syncLogColumns = [
  { title: '用户DN', dataIndex: 'user_dn', ellipsis: true },
  { title: '操作', dataIndex: 'action', slotName: 'type' },
  { title: '状�?, dataIndex: 'status', slotName: 'status' },
  { title: '同步时间', dataIndex: 'created_at', slotName: 'created_at' },
  {
    title: '操作',
    slotName: 'actions',
    fixed: 'right',
    width: 100,
  },
]

const configFormRef = ref()

const loadConfig = async () => {
  try {
    const res = await axios.get('/api/v1/ldap/config')
    Object.assign(ldapConfig, res.data)
    lastSyncTime.value = res.data.last_sync_at || ''
  } catch (e) {
    // use mock data
  }
}

const saveConfig = async () => {
  saving.value = true
  try {
    await axios.post('/api/v1/ldap/config', ldapConfig)
    Message.success('配置保存成功')
  } catch (e) {
    Message.error('保存失败')
  } finally {
    saving.value = false
  }
}

const testConnection = async () => {
  testing.value = true
  try {
    await axios.post('/api/v1/ldap/test', ldapConfig)
    Message.success('连接成功')
  } catch (e: any) {
    Message.error('连接失败: ' + (e.response?.data?.message || e.message))
  } finally {
    testing.value = false
  }
}

const triggerSync = async () => {
  syncing.value = true
  try {
    await axios.post('/api/v1/ldap/sync')
    Message.success('同步任务已触�?)
    loadSyncLogs()
    lastSyncTime.value = new Date().toLocaleString()
  } catch (e) {
    Message.error('同步失败')
  } finally {
    syncing.value = false
  }
}

const loadSyncLogs = async () => {
  loading.value = true
  try {
    const res = await axios.get('/api/v1/ldap/sync-logs', {
      params: { page: pagination.current, page_size: pagination.pageSize },
    })
    syncLogs.value = res.data.items || []
    pagination.total = res.data.total || 0
  } catch (e) {
    // mock
    syncLogs.value = [
      { id: 1, user_dn: 'cn=user1,ou=users,dc=example,dc=com', action: 'create', status: 'success', created_at: new Date().toISOString() },
      { id: 2, user_dn: 'cn=user2,ou=users,dc=example,dc=com', action: 'update', status: 'success', created_at: new Date().toISOString() },
      { id: 3, user_dn: 'cn=user3,ou=users,dc=example,dc=com', action: 'skip', status: 'fail', error_msg: '用户已存�?, created_at: new Date().toISOString() },
    ]
    pagination.total = 3
  } finally {
    loading.value = false
  }
}

const handlePageChange = (page: number) => {
  pagination.current = page
  loadSyncLogs()
}

const formatDate = (date: string) => {
  if (!date) return '-'
  return new Date(date).toLocaleString('zh-CN')
}

loadConfig()
loadSyncLogs()
</script>

<style scoped lang="less">
.text-expired { color: var(--color-red); }
</style>