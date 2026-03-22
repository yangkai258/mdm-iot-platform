<template>
  <div class="page-container">
    <!-- 顶部筛选栏 -->
    <a-card class="filter-card">
      <div class="filter-row">
        <a-select
          v-model="selectedRoleId"
          placeholder="角色选择"
          style="width: 180px"
          allow-search
          @change="handleRoleChange"
        >
          <a-option v-for="r in roles" :key="r.id" :value="r.id">{{ r.name }}</a-option>
        </a-select>
        <a-select
          v-model="selectedUserId"
          placeholder="用户选择"
          style="width: 180px"
          allow-search
          allow-clear
          @change="handleUserChange"
        >
          <a-option v-for="u in users" :key="u.id" :value="u.id">{{ u.username || u.name }}</a-option>
        </a-select>
        <div class="filter-spacer" />
        <a-button @click="handleReset">重置</a-button>
        <a-button type="primary" @click="handleSave">保存配置</a-button>
      </div>
    </a-card>

    <!-- 内容区 -->
    <a-card class="content-card">
      <!-- 模块权限 -->
      <div class="section-title">模块权限</div>
      <a-table
        :columns="moduleColumns"
        :data="modulePermissions"
        :pagination="false"
        row-key="module"
        size="small"
        class="module-table"
      >
        <template #module="{ record }">
          <a-checkbox
            v-model="record.enabled"
            @change="handleModuleToggle(record)"
          >{{ record.label }}</a-checkbox>
        </template>
        <template #permissions="{ record }">
          <a-checkbox
            v-for="p in record.permissions"
            :key="p.key"
            :value="p.key"
            :disabled="!record.enabled"
            v-model="p.checked"
            class="perm-check"
          >{{ p.label }}</a-checkbox>
        </template>
      </a-table>

      <a-divider />

      <!-- 数据权限 -->
      <div class="section-title">数据权限</div>
      <div class="data-section">
        <div class="data-row">
          <span class="data-label">权限范围：</span>
          <a-radio-group v-model="dataScope" @change="handleDataScopeChange">
            <a-radio value="all">全部数据</a-radio>
            <a-radio value="department">本部门</a-radio>
            <a-radio value="self">本人</a-radio>
          </a-radio-group>
        </div>
        <div class="data-row">
          <span class="data-label">字段权限：</span>
          <a-checkbox
            v-for="field in fieldPermissions"
            :key="field.key"
            v-model="field.checked"
            class="perm-check"
          >{{ field.label }}</a-checkbox>
        </div>
      </div>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import { getRoles, getUsers, getDataPermissionRoles, updateDataPermissionRoles } from '@/api/security'

const selectedRoleId = ref(null)
const selectedUserId = ref(null)
const dataScope = ref('all')

const roles = ref([])
const users = ref([])

// 模块权限配置
const modulePermissions = reactive([
  {
    module: 'device',
    label: '设备管理',
    enabled: false,
    permissions: [
      { key: 'device_view', label: '查看', checked: false },
      { key: 'device_add', label: '新增', checked: false },
      { key: 'device_edit', label: '编辑', checked: false },
      { key: 'device_delete', label: '删除', checked: false },
      { key: 'device_export', label: '导出', checked: false }
    ]
  },
  {
    module: 'member',
    label: '会员管理',
    enabled: false,
    permissions: [
      { key: 'member_view', label: '查看', checked: false },
      { key: 'member_add', label: '新增', checked: false },
      { key: 'member_edit', label: '编辑', checked: false },
      { key: 'member_delete', label: '删除', checked: false },
      { key: 'member_export', label: '导出', checked: false }
    ]
  },
  {
    module: 'alert',
    label: '告警系统',
    enabled: false,
    permissions: [
      { key: 'alert_view', label: '查看', checked: false },
      { key: 'alert_confirm', label: '确认', checked: false },
      { key: 'alert_handle', label: '处理', checked: false },
      { key: 'alert_export', label: '导出', checked: false }
    ]
  },
  {
    module: 'system',
    label: '系统配置',
    enabled: false,
    permissions: [
      { key: 'system_view', label: '查看', checked: false },
      { key: 'system_edit', label: '编辑', checked: false }
    ]
  }
])

// 字段权限配置
const fieldPermissions = reactive([
  { key: 'phone', label: '手机号', checked: true },
  { key: 'address', label: '地址', checked: true },
  { key: 'email', label: '邮箱', checked: true },
  { key: 'name', label: '姓名', checked: true }
])

const moduleColumns = [
  { title: '模块', slotName: 'module', width: 180 },
  { title: '权限', slotName: 'permissions' }
]

onMounted(async () => {
  try {
    const [rolesRes, usersRes] = await Promise.all([
      getRoles(),
      getUsers()
    ])
    roles.value = rolesRes.data || rolesRes || []
    users.value = usersRes.data || usersRes || []
  } catch (e) {
    console.error('加载数据失败', e)
  }
})

function handleModuleToggle(record) {
  if (!record.enabled) {
    record.permissions.forEach(p => { p.checked = false })
  }
}

async function handleRoleChange(roleId) {
  if (!roleId) return
  try {
    const res = await getDataPermissionRoles(roleId)
    const data = res.data || res
    if (data) {
      dataScope.value = data.data_scope || 'all'
      // 还原模块权限
      modulePermissions.forEach(mod => {
        const modPerms = data.module_permissions?.[mod.module] || {}
        mod.enabled = modPerms.enabled || false
        mod.permissions.forEach(p => {
          p.checked = modPerms[p.key] || false
        })
      })
      // 还原字段权限
      const allowedFields = data.allowed_fields || []
      fieldPermissions.forEach(f => {
        f.checked = allowedFields.includes(f.key)
      })
    }
  } catch (e) {
    console.error('加载角色权限失败', e)
  }
}

function handleUserChange(userId) {
  // 如果选择了用户，后续可以加载用户特定权限
}

function handleDataScopeChange() {
  // 范围变化
}

function handleReset() {
  selectedRoleId.value = null
  selectedUserId.value = null
  dataScope.value = 'all'
  modulePermissions.forEach(mod => {
    mod.enabled = false
    mod.permissions.forEach(p => { p.checked = false })
  })
  fieldPermissions.forEach(f => { f.checked = true })
  Message.success('已重置')
}

async function handleSave() {
  if (!selectedRoleId.value) {
    Message.warning('请先选择角色')
    return
  }
  try {
    const payload = {
      role_id: selectedRoleId.value,
      data_scope: dataScope.value,
      allowed_fields: fieldPermissions.filter(f => f.checked).map(f => f.key),
      module_permissions: {}
    }
    modulePermissions.forEach(mod => {
      payload.module_permissions[mod.module] = {
        enabled: mod.enabled,
        ...Object.fromEntries(mod.permissions.map(p => [p.key, p.checked]))
      }
    })
    await updateDataPermissionRoles(selectedRoleId.value, payload)
    Message.success('保存成功')
  } catch (e) {
    Message.error('保存失败')
  }
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

.filter-card {
  flex-shrink: 0;
}

.filter-row {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.filter-spacer {
  flex: 1;
}

.content-card {
  flex: 1;
  overflow: auto;
}

.section-title {
  font-size: 15px;
  font-weight: 600;
  margin-bottom: 12px;
  color: var(--color-text-1);
}

.module-table {
  margin-bottom: 8px;
}

.perm-check {
  margin-right: 16px;
}

.data-section {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.data-row {
  display: flex;
  align-items: flex-start;
  flex-wrap: wrap;
  gap: 8px;
}

.data-label {
  color: var(--color-text-2);
  font-size: 14px;
  min-width: 80px;
  line-height: 32px;
}
</style>
