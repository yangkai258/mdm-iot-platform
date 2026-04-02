<template>
    <Breadcrumb :items="['Home','Console','']" />


  <div class="page-container">
    <!-- 同步状态卡片 -->
    <a-card class="status-card">
      <div class="status-grid">
        <div class="status-item">
          <div class="status-icon">
            <icon-clock-circle />
          </div>
          <div class="status-info">
            <div class="status-label">上次同步时间</div>
            <div class="status-value">{{ lastSyncTime || '从未同步' }}</div>
          </div>
        </div>
        <div class="status-item">
          <div class="status-icon">
            <icon-user />
          </div>
          <div class="status-info">
            <div class="status-label">LDAP 用户总数</div>
            <div class="status-value">{{ syncStats.total_users || 0 }}</div>
          </div>
        </div>
        <div class="status-item">
          <div class="status-icon success">
            <icon-check-circle />
          </div>
          <div class="status-info">
            <div class="status-label">同步成功</div>
            <div class="status-value">{{ syncStats.synced_users || 0 }}</div>
          </div>
        </div>
        <div class="status-item">
          <div class="status-icon danger">
            <icon-close-circle />
          </div>
          <div class="status-info">
            <div class="status-label">同步失败</div>
            <div class="status-value">{{ syncStats.failed_users || 0 }}</div>
          </div>
        </div>
        <div class="status-actions">
          <a-button
            type="primary"
            :loading="syncing"
            @click="handleSync"
          >
            <template #icon><icon-sync /></template>
            立即同步
          </a-button>
        </div>
      </div>
    </a-card>

    <!-- 标签页 -->
    <a-card class="content-card">
      <a-tabs default-active-key="users">
        <!-- LDAP 用户列表 -->
        <a-tab-pane key="users" title="用户列表">
          <a-table
            :columns="userColumns"
            :data="ldapUsers"
            :loading="usersLoading"
            :pagination="userPagination"
            row-key="ldap_dn"
            @change="handleUsersPageChange"
          >
            <template #status="{ record }">
              <a-tag :color="syncStatusColor(record.sync_status)">
                {{ syncStatusLabel(record.sync_status) }}
              </a-tag>
            </template>
            <template #last_synced="{ record }">
              {{ formatDate(record.last_synced_at) }}
            </template>
          </a-table>
        </a-tab-pane>

        <!-- 映射规则配置 -->
        <a-tab-pane key="mapping" title="映射规则">
          <div class="mapping-header">
            <div class="mapping-title">LDAP 分组 → 系统角色映射</div>
            <a-button type="primary" size="small" @click="openAddMapping">
              <template #icon><icon-plus /></template>
              添加映射
            </a-button>
          </div>
          <a-table
            :columns="mappingColumns"
            :data="mappings"
            :loading="mappingLoading"
            :pagination="false"
            row-key="id"
          >
            <template #actions="{ record }">
              <a-button type="text" size="small" @click="handleEditMapping(record)">编辑</a-button>
              <a-divider direction="vertical" />
              <a-popconfirm content="确定删除该映射？" @ok="handleDeleteMapping(record)">
                <a-button type="text" size="small" status="danger">删除</a-button>
              </a-popconfirm>
            </template>
          </a-table>
        </a-tab-pane>

        <!-- 同步历史 -->
        <a-tab-pane key="history" title="同步历史">
          <a-table
            :columns="historyColumns"
            :data="syncHistory"
            :loading="historyLoading"
            :pagination="historyPagination"
            row-key="id"
            @change="handleHistoryPageChange"
          >
            <template #status="{ record }">
              <a-tag :color="record.status === 'success' ? 'green' : 'red'">
                {{ record.status === 'success' ? '成功' : '失败' }}
              </a-tag>
            </template>
            <template #time="{ record }">
              {{ formatDate(record.created_at) }}
            </template>
          </a-table>
        </a-tab-pane>
      </a-tabs>
    </a-card>

    <!-- 添加/编辑映射弹窗 -->
    <a-modal
      v-model:visible="mappingModalVisible"
      :title="editingMapping ? '编辑映射' : '添加映射'"
      width="440px"
      @ok="handleSaveMapping"
      :ok-loading="mappingSaving"
    >
      <a-form :model="mappingForm" layout="vertical">
        <a-form-item label="LDAP 分组 DN" required>
          <a-input
            v-model="mappingForm.ldap_group_dn"
            placeholder="cn=admins,ou=groups,dc=example,dc=com"
          />
        </a-form-item>
        <a-form-item label="LDAP 分组名称">
          <a-input
            v-model="mappingForm.ldap_group_name"
            placeholder="系统管理员"
          />
        </a-form-item>
        <a-form-item label="映射到系统角色" required>
          <a-select
            v-model="mappingForm.role_id"
            placeholder="请选择角色"
            allow-search
          >
            <a-option v-for="r in roles" :key="r.id" :value="r.id">{{ r.name }}</a-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import {
  getLdapUsers,
  syncLdapUsers,
  getLdapGroups,
  setLdapGroupMapping,
  getRoles
} from '@/api/security'
import dayjs from 'dayjs'

const syncing = ref(false)
const usersLoading = ref(false)
const mappingLoading = ref(false)
const historyLoading = ref(false)
const mappingSaving = ref(false)

const lastSyncTime = ref('')
const syncStats = reactive({
  total_users: 0,
  synced_users: 0,
  failed_users: 0
})
const ldapUsers = ref([])
const syncHistory = ref([])
const mappings = ref([])
const roles = ref([])

const userPagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0
})

const historyPagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0
})

const mappingModalVisible = ref(false)
const editingMapping = ref(null)
const mappingForm = reactive({
  ldap_group_dn: '',
  ldap_group_name: '',
  role_id: null
})

const userColumns = [
  { title: 'LDAP DN', dataIndex: 'ldap_dn', ellipsis: true },
  { title: '用户名', dataIndex: 'username' },
  { title: '邮箱', dataIndex: 'email' },
  { title: '显示名', dataIndex: 'display_name' },
  { title: '同步状态', slotName: 'status', width: 100 },
  { title: '最后同步', slotName: 'last_synced', width: 160 }
]

const mappingColumns = [
  { title: 'LDAP 分组 DN', dataIndex: 'ldap_group_dn', ellipsis: true },
  { title: '分组名称', dataIndex: 'ldap_group_name' },
  { title: '映射角色', dataIndex: 'role_name' },
  { title: '操作', slotName: 'actions', width: 140 }
]

const historyColumns = [
  { title: '同步时间', slotName: 'time', width: 160 },
  { title: '触发方式', dataIndex: 'trigger_type', width: 100 },
  { title: '新增用户', dataIndex: 'added_count', width: 90 },
  { title: '更新用户', dataIndex: 'updated_count', width: 90 },
  { title: '失败数量', dataIndex: 'failed_count', width: 90 },
  { title: '状态', slotName: 'status', width: 80 }
]

onMounted(async () => {
  await Promise.all([
    loadLdapUsers(),
    loadMappings(),
    loadRoles()
  ])
})

async function loadLdapUsers() {
  usersLoading.value = true
  try {
    const res = await getLdapUsers({ page: userPagination.current, page_size: userPagination.pageSize })
    const data = res.data || res
    ldapUsers.value = data.list || data.records || data || []
    userPagination.total = data.total || ldapUsers.value.length
    // 更新统计
    if (data.stats) {
      syncStats.total_users = data.stats.total || 0
      syncStats.synced_users = data.stats.synced || 0
      syncStats.failed_users = data.stats.failed || 0
    }
    if (data.last_sync_at) {
      lastSyncTime.value = formatDate(data.last_sync_at)
    }
  } catch (e) {
    console.error('加载LDAP用户失败', e)
  } finally {
    usersLoading.value = false
  }
}

async function loadMappings() {
  mappingLoading.value = true
  try {
    const res = await getLdapGroups()
    mappings.value = res.data || res || []
  } catch (e) {
    console.error('加载映射规则失败', e)
  } finally {
    mappingLoading.value = false
  }
}

async function loadRoles() {
  try {
    const res = await getRoles()
    roles.value = res.data || res || []
  } catch (e) {
    console.error('加载角色列表失败', e)
  }
}

async function handleSync() {
  syncing.value = true
  try {
    await syncLdapUsers()
    Message.success('同步任务已触发')
    await loadLdapUsers()
  } catch (e) {
    Message.error('同步失败')
  } finally {
    syncing.value = false
  }
}

function handleUsersPageChange(pag) {
  userPagination.current = pag.current
  userPagination.pageSize = pag.pageSize
  loadLdapUsers()
}

function handleHistoryPageChange(pag) {
  historyPagination.current = pag.current
  historyPagination.pageSize = pag.pageSize
  // loadSyncHistory()
}

function openAddMapping() {
  editingMapping.value = null
  mappingForm.ldap_group_dn = ''
  mappingForm.ldap_group_name = ''
  mappingForm.role_id = null
  mappingModalVisible.value = true
}

function handleEditMapping(record) {
  editingMapping.value = record
  mappingForm.ldap_group_dn = record.ldap_group_dn
  mappingForm.ldap_group_name = record.ldap_group_name
  mappingForm.role_id = record.role_id
  mappingModalVisible.value = true
}

async function handleSaveMapping() {
  mappingSaving.value = true
  try {
    await setLdapGroupMapping({
      id: editingMapping.value?.id,
      ldap_group_dn: mappingForm.ldap_group_dn,
      ldap_group_name: mappingForm.ldap_group_name,
      role_id: mappingForm.role_id
    })
    Message.success('保存成功')
    mappingModalVisible.value = false
    await loadMappings()
  } catch (e) {
    Message.error('保存失败')
  } finally {
    mappingSaving.value = false
  }
}

async function handleDeleteMapping(record) {
  try {
    await setLdapGroupMapping({ ...record, _delete: true })
    Message.success('删除成功')
    await loadMappings()
  } catch (e) {
    Message.error('删除失败')
  }
}

function syncStatusColor(status) {
  const map = { synced: 'green', pending: 'yellow', failed: 'red' }
  return map[status] || 'default'
}

function syncStatusLabel(status) {
  const map = { synced: '已同步', pending: '待同步', failed: '失败' }
  return map[status] || status
}

function formatDate(date) {
  if (!date) return '-'
  return dayjs(date).format('YYYY-MM-DD HH:mm:ss')
}
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

.status-card {
  flex-shrink: 0;
}

.status-grid {
  display: flex;
  align-items: center;
  gap: 24px;
  flex-wrap: wrap;
}

.status-item {
  display: flex;
  align-items: center;
  gap: 10px;
}

.status-icon {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: var(--color-fill-lightest);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--color-text-3);
}

.status-icon.success {
  color: #0fbf61;
  background: #e6fcf0;
}

.status-icon.danger {
  color: #f53f3f;
  background: #fff1f0;
}

.status-info {
  display: flex;
  flex-direction: column;
}

.status-label {
  font-size: 12px;
  color: var(--color-text-3);
}

.status-value {
  font-size: 15px;
  font-weight: 600;
  color: var(--color-text-1);
}

.status-actions {
  margin-left: auto;
}

.content-card {
  flex: 1;
  overflow: auto;
}

.mapping-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}

.mapping-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--color-text-1);
}
</style>

