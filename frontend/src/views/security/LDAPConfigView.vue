<template>
  <div class="page-container">
    <!-- Tab 切换 -->
    <a-tabs v-model:active-tab="activeTab" class="ldap-tabs">
      <!-- Tab1: LDAP 配置 -->
      <a-tab-pane key="config" title="LDAP 配置">
        <a-card class="content-card">
          <template #title>
            <div class="card-header">
              <span>LDAP/AD 配置</span>
              <a-tag v-if="ldapConfig.is_enabled" color="green">已启用</a-tag>
              <a-tag v-else color="gray">未启用</a-tag>
            </div>
          </template>

          <a-form :model="form" layout="vertical" class="ldap-form">
            <div class="form-section">
              <div class="section-label">服务器配置</div>
              <a-form-item label="配置名称" required>
                <a-input v-model="form.config_name" placeholder="例如：企业LDAP" style="width: 300px" />
              </a-form-item>
              <a-form-item label="LDAP服务器地址" required>
                <a-input v-model="form.host" placeholder="ldap.company.com 或 ldaps://ldap.company.com" style="width: 360px" />
              </a-form-item>
              <a-form-item label="端口" required>
                <a-input-number v-model="form.port" :min="1" :max="65535" style="width: 120px" />
                <span class="form-tip">389（LDAP）或 636（LDAPS）</span>
              </a-form-item>
              <a-form-item label="使用 SSL">
                <a-switch v-model="form.use_ssl" />
              </a-form-item>
              <a-form-item label="使用 TLS">
                <a-switch v-model="form.use_tls" />
              </a-form-item>
            </div>

            <a-divider />

            <div class="form-section">
              <div class="section-label">认证配置</div>
              <a-form-item label="基准 DN" required>
                <a-input v-model="form.base_dn" placeholder="dc=company,dc=com" style="width: 360px" />
              </a-form-item>
              <a-form-item label="管理员 DN" required>
                <a-input v-model="form.bind_dn" placeholder="cn=admin,dc=company,dc=com" style="width: 360px" />
              </a-form-item>
              <a-form-item label="管理员密码">
                <a-input-password v-model="form.bind_password" placeholder="请输入密码（加密存储）" allow-clear style="width: 240px" />
              </a-form-item>
            </div>

            <a-divider />

            <div class="form-section">
              <div class="section-label">过滤规则</div>
              <a-form-item label="用户过滤器">
                <a-input v-model="form.user_filter" placeholder="(objectClass=user)" style="width: 400px" />
                <div class="form-tip">用于搜索用户，常用：(objectClass=user)</div>
              </a-form-item>
              <a-form-item label="分组过滤器">
                <a-input v-model="form.group_filter" placeholder="(objectClass=group)" style="width: 400px" />
              </a-form-item>
            </div>

            <a-divider />

            <div class="form-section">
              <div class="section-label">同步配置</div>
              <a-form-item label="同步间隔（秒）">
                <a-input-number v-model="form.sync_interval" :min="60" :step="60" style="width: 120px" />
                <span class="form-tip">默认 3600 秒（1小时），建议不低于 300 秒</span>
              </a-form-item>
              <a-form-item label="启用状态">
                <a-switch v-model="form.is_enabled" />
              </a-form-item>
              <a-form-item label="描述">
                <a-textarea v-model="form.description" placeholder="可选描述信息" :rows="2" style="width: 400px" />
              </a-form-item>
            </div>

            <div class="form-actions">
              <a-button @click="handleTest" :loading="testing" :disabled="!form.host || !form.port">
                测试连接
              </a-button>
              <a-button type="primary" @click="handleSave" :loading="saving">
                保存配置
              </a-button>
            </div>
          </a-form>
        </a-card>
      </a-tab-pane>

      <!-- Tab2: 用户同步 -->
      <a-tab-pane key="users" title="用户同步">
        <a-card class="content-card">
          <template #title>
            <span>LDAP 用户列表</span>
          </template>
          <template #extra>
            <a-space>
              <a-button type="primary" @click="handleSync" :loading="syncing">同步用户</a-button>
              <a-button @click="loadUsers">刷新</a-button>
            </a-space>
          </template>

          <!-- 同步结果统计 -->
          <a-alert v-if="syncResult" class="sync-alert" :type="syncResult.errors?.length > 0 ? 'warning' : 'success'" closable>
            <template #title>
              同步完成：共 {{ syncResult.total_users }} 用户，新增 {{ syncResult.added }}，更新 {{ syncResult.updated }}，跳过 {{ syncResult.skipped }}
            </template>
          </a-alert>

          <!-- 搜索 -->
          <div class="search-bar">
            <a-input-search v-model="userKeyword" placeholder="搜索用户名/邮箱" style="width: 240px" @search="loadUsers" search-button />
          </div>

          <a-table :columns="userColumns" :data="ldapUsers" :loading="usersLoading" :pagination="userPagination" row-key="dn" @page-change="handleUserPageChange" size="small">
            <template #actions="{ record }">
              <a-button type="text" size="small" @click="showUserDetail(record)">查看</a-button>
              <a-button type="text" size="small" status="danger" @click="removeUser(record)">移除</a-button>
            </template>
          </a-table>
        </a-card>
      </a-tab-pane>

      <!-- Tab3: 分组-角色映射 -->
      <a-tab-pane key="mapping" title="分组-角色映射">
        <a-card class="content-card">
          <template #title>
            <span>LDAP 分组 - 角色映射</span>
          </template>
          <template #extra>
            <a-button type="primary" @click="showMappingModal">添加映射</a-button>
          </template>

          <a-table :columns="mappingColumns" :data="groupMappings" :loading="mappingLoading" :pagination="false" row-key="id" size="small">
            <template #actions="{ record }">
              <a-button type="text" size="small" status="danger" @click="removeMapping(record)">删除</a-button>
            </template>
          </a-table>
        </a-card>
      </a-tab-pane>
    </a-tabs>

    <!-- 测试连接结果弹窗 -->
    <a-modal v-model:visible="testResultVisible" title="连接测试结果" width="440px" :footer="null">
      <a-result
        :status="testResult?.success ? 'success' : 'error'"
        :title="testResult?.success ? '连接成功' : '连接失败'"
      >
        <template #content>
          <div>{{ testResult?.message }}</div>
          <div v-if="testResult?.server_info" class="server-info">服务器信息：{{ testResult.server_info }}</div>
        </template>
      </a-result>
    </a-modal>

    <!-- 用户详情弹窗 -->
    <a-modal v-model:visible="userDetailVisible" title="用户详情" width="480px" :footer="null">
      <a-descriptions :column="2" bordered v-if="currentUser">
        <a-descriptions-item label="DN" :span="2">{{ currentUser.dn }}</a-descriptions-item>
        <a-descriptions-item label="用户名">{{ currentUser.username }}</a-descriptions-item>
        <a-descriptions-item label="显示名称">{{ currentUser.display_name }}</a-descriptions-item>
        <a-descriptions-item label="邮箱">{{ currentUser.email }}</a-descriptions-item>
        <a-descriptions-item label="分组" :span="2">
          <a-tag v-for="g in currentUser.groups" :key="g" size="small">{{ g }}</a-tag>
        </a-descriptions-item>
      </a-descriptions>
    </a-modal>

    <!-- 添加映射弹窗 -->
    <a-modal v-model:visible="mappingModalVisible" title="添加分组-角色映射" :width="480" :loading="mappingSubmitting" @before-ok="handleAddMapping" @cancel="mappingModalVisible = false">
      <a-form :model="mappingForm" layout="vertical">
        <a-form-item label="LDAP分组" required>
          <a-select v-model="mappingForm.ldap_group_dn" placeholder="选择LDAP分组" show-search>
            <a-option v-for="g in ldapGroups" :key="g.dn" :value="g.dn">{{ g.name }} ({{ g.dn }})</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="映射名称">
          <a-input v-model="mappingForm.ldap_group_name" placeholder="分组名称" />
        </a-form-item>
        <a-form-item label="关联角色" required>
          <a-select v-model="mappingForm.role_id" placeholder="选择角色">
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
import { getLdapConfig, updateLdapConfig, testLdapConnection, getLdapUsers, syncLdapUsers, getLdapGroups, getLdapGroupMappings, setLdapGroupMapping } from '@/api/security'
import { getRoles } from '@/api/security'

const activeTab = ref('config')
const saving = ref(false)
const testing = ref(false)
const syncing = ref(false)
const testResultVisible = ref(false)
const testResult = ref(null)
const userDetailVisible = ref(false)
const currentUser = ref(null)
const mappingModalVisible = ref(false)
const mappingSubmitting = ref(false)

const ldapConfig = reactive({ is_enabled: false, last_sync_at: null })
const ldapUsers = ref([])
const ldapGroups = ref([])
const groupMappings = ref([])
const roles = ref([])
const syncResult = ref(null)
const usersLoading = ref(false)
const mappingLoading = ref(false)
const userKeyword = ref('')

const userPagination = reactive({ current: 1, pageSize: 20, total: 0 })

const form = reactive({
  config_name: '',
  host: '',
  port: 389,
  use_ssl: false,
  use_tls: false,
  base_dn: '',
  bind_dn: '',
  bind_password: '',
  user_filter: '(objectClass=user)',
  group_filter: '(objectClass=group)',
  sync_interval: 3600,
  is_enabled: false,
  description: ''
})

const mappingForm = reactive({ ldap_group_dn: '', ldap_group_name: '', role_id: null })

const userColumns = [
  { title: 'DN', dataIndex: 'dn', ellipsis: true },
  { title: '用户名', dataIndex: 'username', width: 120 },
  { title: '显示名称', dataIndex: 'display_name', width: 120 },
  { title: '邮箱', dataIndex: 'email', width: 180 },
  { title: '操作', slotName: 'actions', width: 120 }
]

const mappingColumns = [
  { title: 'LDAP分组DN', dataIndex: 'ldap_group_dn', ellipsis: true },
  { title: '分组名称', dataIndex: 'ldap_group_name', width: 140 },
  { title: '关联角色', dataIndex: 'role_name', width: 120 },
  { title: '操作', slotName: 'actions', width: 80 }
]

onMounted(async () => {
  await loadConfig()
  loadUsers()
  loadMappings()
})

const loadConfig = async () => {
  try {
    const res = await getLdapConfig()
    const data = res.data || res
    if (data && data.id) {
      Object.assign(ldapConfig, data)
      Object.assign(form, {
        config_name: data.config_name,
        host: data.host,
        port: data.port,
        use_ssl: data.use_ssl,
        use_tls: data.use_tls,
        base_dn: data.base_dn,
        bind_dn: data.bind_dn,
        bind_password: data.bind_password || '',
        user_filter: data.user_filter || '(objectClass=user)',
        group_filter: data.group_filter || '(objectClass=group)',
        sync_interval: data.sync_interval || 3600,
        is_enabled: data.is_enabled,
        description: data.description || ''
      })
    }
  } catch (e) {
    console.error('加载LDAP配置失败', e)
  }
}

const loadUsers = async () => {
  usersLoading.value = true
  try {
    const res = await getLdapUsers({ query: userKeyword.value, page: userPagination.current, page_size: userPagination.pageSize })
    const data = res.data || res
    ldapUsers.value = data.list || []
    userPagination.total = data.total || 0
  } catch (e) {
    console.error('加载用户失败', e)
  } finally {
    usersLoading.value = false
  }
}

const loadMappings = async () => {
  mappingLoading.value = true
  try {
    const res = await getLdapGroupMappings()
    groupMappings.value = (res.data || res)?.list || []
    // 加载分组
    const grpRes = await getLdapGroups()
    ldapGroups.value = (grpRes.data || grpRes)?.list || []
    // 加载角色
    const roleRes = await getRoles()
    roles.value = (roleRes.data || roleRes)?.list || roleRes || []
  } catch (e) {
    console.error('加载映射失败', e)
  } finally {
    mappingLoading.value = false
  }
}

const handleSave = async () => {
  if (!form.config_name || !form.host || !form.base_dn || !form.bind_dn) {
    Message.warning('请填写必填字段')
    return
  }
  saving.value = true
  try {
    await updateLdapConfig({ ...form })
    Object.assign(ldapConfig, { is_enabled: form.is_enabled })
    Message.success('保存成功')
  } catch (e) {
    Message.error('保存失败')
  } finally {
    saving.value = false
  }
}

const handleTest = async () => {
  testing.value = true
  try {
    const res = await testLdapConnection({
      host: form.host,
      port: form.port,
      use_ssl: form.use_ssl,
      use_tls: form.use_tls,
      base_dn: form.base_dn,
      bind_dn: form.bind_dn,
      bind_password: form.bind_password
    })
    testResult.value = { success: res.success !== false, message: res.message || '连接成功', server_info: res.server_info }
  } catch (e) {
    testResult.value = { success: false, message: e.message || '连接失败' }
  } finally {
    testing.value = false
    testResultVisible.value = true
  }
}

const handleSync = async () => {
  syncing.value = true
  syncResult.value = null
  try {
    const res = await syncLdapUsers()
    syncResult.value = res.data || res
    Message.success('同步完成')
    loadUsers()
  } catch (e) {
    Message.error('同步失败')
  } finally {
    syncing.value = false
  }
}

const handleUserPageChange = (page) => { userPagination.current = page; loadUsers() }

const showUserDetail = (user) => { currentUser.value = user; userDetailVisible.value = true }

const removeUser = async (user) => {
  // 调用移除接口（后端实现）
  Message.info('移除用户：' + user.username)
}

const showMappingModal = () => {
  Object.assign(mappingForm, { ldap_group_dn: '', ldap_group_name: '', role_id: null })
  mappingModalVisible.value = true
}

const handleAddMapping = async (done) => {
  if (!mappingForm.ldap_group_dn || !mappingForm.role_id) {
    Message.warning('请填写分组和角色')
    done(false)
    return
  }
  mappingSubmitting.value = true
  try {
    await setLdapGroupMapping(mappingForm)
    Message.success('映射添加成功')
    mappingModalVisible.value = false
    loadMappings()
    done(true)
  } catch (e) {
    Message.error('添加失败')
    done(false)
  } finally {
    mappingSubmitting.value = false
  }
}

const removeMapping = async (record) => {
  // 简化处理，实际应调用删除API
  Message.info('删除映射：' + record.ldap_group_name)
  loadMappings()
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
.ldap-tabs { background: #fff; border-radius: 8px; }
.content-card { border-radius: 8px; }
.card-header { display: flex; align-items: center; gap: 10px; }
.ldap-form { max-width: 600px; }
.form-section { display: flex; flex-direction: column; gap: 4px; }
.section-label { font-size: 13px; font-weight: 600; color: var(--color-text-1); margin-bottom: 8px; }
.form-tip { color: var(--color-text-3); font-size: 12px; margin-left: 8px; }
.form-actions { display: flex; gap: 10px; margin-top: 8px; }
.search-bar { margin: 12px 0; }
.sync-alert { margin-bottom: 12px; }
.server-info { margin-top: 8px; font-size: 12px; color: #666; }
</style>
