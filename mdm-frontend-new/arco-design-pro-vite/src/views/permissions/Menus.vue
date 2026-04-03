<template>
  <div class="page-container">
    <Breadcrumb :items="['Home', 'Console', '']" />
    <a-card class="general-card" title="菜单管理">
      <template #extra>
        <a-space :size="12">
          <a-button type="primary" @click="openCreateModal(null)"><icon-plus />新建</a-button>
          <a-button @click="loadMenus"><icon-refresh />刷新</a-button>
          <a-button type="primary" status="warning" :disabled="!hasDragged" @click="saveSort">
            <icon-save />保存排序
          </a-button>
        </a-space>
      </template>
      <div class="toolbar">
        <a-input-search
          v-model="searchForm.keyword"
          placeholder="菜单名称/编码"
          style="width: 280px"
          @search="doSearch"
          allow-clear
        />
        <a-button type="primary" status="warning" @click="batchSaveSort" :disabled="sortChangedCount === 0">
          <icon-save /> 保存排序 ({{ sortChangedCount }})
        </a-button>
      </div>

      <!-- 拖拽排序表格 -->
      <a-table
        :columns="columns"
        :data="menuTree"
        :loading="loading"
        :expanded-keys="expandedKeys"
        :pagination="false"
        :row-class="() => 'sortable-row'"
        row-key="id"
        @expand="onExpand"
        @dragrow="onDragRow"
      >
        <template #drag-handle="{ record }">
          <span class="drag-handle" title="拖拽排序">⋮⋮</span>
        </template>
        <template #menu_name="{ record }">
          <span class="menu-name-cell">
            <component :is="record.icon" v-if="record.icon && isValidIcon(record.icon)" class="menu-icon" />
            <span>{{ record.name }}</span>
          </span>
        </template>
        <template #menu_type="{ record }">
          <a-tag :color="getTypeColor(record.type)">{{ getTypeLabel(record.type) }}</a-tag>
        </template>
        <template #sort="{ record }">
          <a-input-number
            v-model="sortMap[record.id]"
            :min="0"
            :max="9999"
            size="small"
            style="width: 70px"
            @change="onSortChange(record.id)"
          />
        </template>
        <template #status="{ record }">
          <a-tag :color="record.status === 1 ? 'green' : 'gray'">
            {{ record.status === 1 ? '启用' : '禁用' }}
          </a-tag>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="openCreateModal(record)">子级</a-button>
            <a-button type="text" size="small" @click="openEditModal(record)">编辑</a-button>
            <a-button type="text" size="small" @click="moveUp(record)" :disabled="!canMoveUp(record)">上移</a-button>
            <a-button type="text" size="small" @click="moveDown(record)" :disabled="!canMoveDown(record)">下移</a-button>
            <a-button type="text" size="small" status="danger" @click="handleDelete(record)">删除</a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 创建/编辑弹窗 -->
    <a-modal
      v-model:visible="formVisible"
      :title="isEdit ? '编辑菜单' : '新建菜单'"
      :width="560"
      :mask-closable="false"
      @ok="submitForm"
      @cancel="formVisible = false"
    >
      <a-form :model="formData" layout="vertical">
        <a-form-item field="parent_id" label="上级菜单">
          <a-tree-select
            v-model="formData.parent_id"
            :data="menuTreeForSelect"
            placeholder="不选则为顶级"
            allow-clear
            :field-names="{ key: 'id', title: 'name', children: 'children' }"
            style="width: 100%"
          />
        </a-form-item>
        <a-form-item field="name" label="菜单名称" required>
          <a-input v-model="formData.name" placeholder="请输入菜单名称" />
        </a-form-item>
        <a-form-item field="code" label="菜单编码">
          <a-input v-model="formData.code" placeholder="如 system:menu" />
        </a-form-item>
        <a-form-item field="type" label="菜单类型" required>
          <a-radio-group v-model="formData.type">
            <a-radio :value="0">目录</a-radio>
            <a-radio :value="1">菜单</a-radio>
            <a-radio :value="2">按钮</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item field="icon" label="图标">
          <a-select v-model="formData.icon" placeholder="选择图标（可选）" allow-clear>
            <a-option value="IconHome" />
            <a-option value="IconUser" />
            <a-option value="IconUserGroup" />
            <a-option value="IconSettings" />
            <a-option value="IconMenu" />
            <a-option value="IconDashboard" />
            <a-option value="IconList" />
            <a-option value="IconFile" />
            <a-option value="IconFolder" />
            <a-option value="IconSearch" />
            <a-option value="IconPlus" />
            <a-option value="IconEdit" />
            <a-option value="IconDelete" />
            <a-option value="IconCheck" />
            <a-option value="IconClose" />
            <a-option value="IconRefresh" />
            <a-option value="IconExport" />
            <a-option value="IconImport" />
            <a-option value="IconArrowUp" />
            <a-option value="IconArrowDown" />
            <a-option value="IconSafe" />
            <a-option value="IconLock" />
            <a-option value="IconKey" />
            <a-option value="IconTool" />
            <a-option value="IconStar" />
            <a-option value="IconHeart" />
            <a-option value="IconMessage" />
            <a-option value="IconBell" />
            <a-option value="IconCalendar" />
            <a-option value="IconClockCircle" />
            <a-option value="IconWifi" />
            <a-option value="IconDevice" />
            <a-option value="IconRobot" />
          </a-select>
        </a-form-item>
        <a-form-item field="path" label="路由路径">
          <a-select v-model="formData.path" placeholder="选择路由（可选）" allow-clear allow-search>
            <a-option value="/dashboard" />
            <a-option value="/device" />
            <a-option value="/devices" />
            <a-option value="/member" />
            <a-option value="/members" />
            <a-option value="/order" />
            <a-option value="/orders" />
            <a-option value="/store" />
            <a-option value="/stores" />
            <a-option value="/ai" />
            <a-option value="/pet" />
            <a-option value="/ota" />
            <a-option value="/alerts" />
            <a-option value="/notifications" />
            <a-option value="/permission-manage" />
            <a-option value="/permissions" />
            <a-option value="/system" />
            <a-option value="/tenant" />
            <a-option value="/tenants" />
            <a-option value="/org" />
            <a-option value="/billing" />
            <a-option value="/subscription" />
            <a-option value="/analytics" />
            <a-option value="/knowledge" />
            <a-option value="/app" />
            <a-option value="/apps" />
            <a-option value="/market" />
            <a-option value="/marketing" />
            <a-option value="/content" />
            <a-option value="/webhooks" />
          </a-select>
        </a-form-item>
        <a-form-item field="component" label="组件路径">
          <a-select v-model="formData.component" placeholder="选择组件" allow-clear allow-search>
            <a-option value="/dashboard/Dashboard.vue" />
            <a-option value="/device/DeviceList.vue" />
            <a-option value="/device/DeviceDetail.vue" />
            <a-option value="/member/MemberList.vue" />
            <a-option value="/member/MemberDetail.vue" />
            <a-option value="/ota/FirmwareList.vue" />
            <a-option value="/ota/VersionList.vue" />
            <a-option value="/alerts/AlertsView.vue" />
            <a-option value="/notifications/NotificationList.vue" />
            <a-option value="/permissions/ApiPermissions.vue" />
            <a-option value="/permissions/Menus.vue" />
            <a-option value="/permissions/Roles.vue" />
            <a-option value="/permissions/PermissionGroups.vue" />
            <a-option value="/permissions/DataPermissionConfig.vue" />
            <a-option value="/system/SystemSettings.vue" />
            <a-option value="/system/Logs.vue" />
            <a-option value="/system/Monitor.vue" />
            <a-option value="/tenant/TenantList.vue" />
            <a-option value="/org/OrgList.vue" />
            <a-option value="/ai/AiConsole.vue" />
            <a-option value="/ai/ModelList.vue" />
            <a-option value="/pet/PetConsole.vue" />
            <a-option value="/pet/PetProfile.vue" />
            <a-option value="/store/StoreList.vue" />
            <a-option value="/billing/BillingList.vue" />
            <a-option value="/analytics/AnalyticsDashboard.vue" />
            <a-option value="/knowledge/KnowledgeBase.vue" />
            <a-option value="/content/ContentList.vue" />
            <a-option value="/webhooks/WebhookList.vue" />
            <a-option value="/list/BasicList.vue" />
            <a-option value="/list/CardList.vue" />
            <a-option value="/form/BasicForm.vue" />
            <a-option value="/profile/ProfileBasic.vue" />
          </a-select>
        </a-form-item>
        <a-form-item field="permission" label="权限标识">
          <a-select v-model="formData.permission" placeholder="选择权限（可选）" allow-clear allow-search>
            <a-option value="system:view" />
            <a-option value="system:create" />
            <a-option value="system:update" />
            <a-option value="system:delete" />
            <a-option value="menu:view" />
            <a-option value="menu:create" />
            <a-option value="menu:update" />
            <a-option value="menu:delete" />
            <a-option value="device:view" />
            <a-option value="device:create" />
            <a-option value="device:update" />
            <a-option value="device:delete" />
            <a-option value="device:control" />
            <a-option value="member:view" />
            <a-option value="member:create" />
            <a-option value="member:update" />
            <a-option value="member:delete" />
            <a-option value="ota:view" />
            <a-option value="ota:create" />
            <a-option value="ota:update" />
            <a-option value="ota:delete" />
            <a-option value="alert:view" />
            <a-option value="alert:create" />
            <a-option value="alert:update" />
            <a-option value="alert:delete" />
            <a-option value="role:view" />
            <a-option value="role:create" />
            <a-option value="role:update" />
            <a-option value="role:delete" />
            <a-option value="tenant:view" />
            <a-option value="tenant:create" />
            <a-option value="tenant:update" />
            <a-option value="tenant:delete" />
          </a-select>
        </a-form-item>
        <a-form-item field="sort" label="排序值">
          <a-input-number v-model="formData.sort" :min="0" :max="9999" style="width: 200px" />
        </a-form-item>
        <a-form-item field="status" label="状态">
          <a-switch v-model="formData.status" :checked-value="1" :unchecked-value="0" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { Message, Modal } from '@arco-design/web-vue'
import { IconPlus, IconRefresh, IconSave } from '@arco-design/web-vue/es/icon'

const loading = ref(false)
const formVisible = ref(false)
const isEdit = ref(false)
const menuList = ref([])
const expandedKeys = ref([])
const currentEditRecord = ref(null)
const sortMap = ref({})     // id -> sort value
const originalSortMap = ref({})  // original sort for comparison
const sortChangedCount = ref(0)

const searchForm = reactive({ keyword: '' })

const formData = reactive({
  id: null, parent_id: null, name: '', code: '', type: 1,
  icon: '', path: '', component: '', permission: '', sort: 0, status: 1
})

const buildTree = (list, parentId = 0) => {
  return list.filter(item => item.parent_id === parentId)
    .sort((a, b) => (sortMap.value[a.id] ?? a.sort) - (sortMap.value[b.id] ?? b.sort))
    .map(item => ({ ...item, children: buildTree(list, item.id) }))
}

const menuTree = computed(() => buildTree(menuList.value))

const getDisabledIds = (record) => {
  const ids = [record.id]
  if (record.children) record.children.forEach(c => ids.push(...getDisabledIds(c)))
  return ids
}

const menuTreeForSelect = computed(() => {
  if (!isEdit.value) return menuTree.value
  const disabledIds = currentEditRecord.value ? getDisabledIds(currentEditRecord.value) : []
  const filter = (list) => list.filter(i => !disabledIds.includes(i.id))
    .map(i => ({ ...i, children: i.children ? filter(i.children) : [] }))
  return filter(menuTree.value)
})

const columns = [
  { title: '', slotName: 'drag-handle', width: 40, align: 'center' },
  { title: '菜单名称', slotName: 'menu_name', minWidth: 200 },
  { title: '编码', dataIndex: 'code', width: 160 },
  { title: '类型', slotName: 'menu_type', width: 90 },
  { title: '排序', slotName: 'sort', width: 100 },
  { title: '状态', slotName: 'status', width: 90 },
  { title: '操作', slotName: 'actions', width: 280, fixed: 'right' }
]

const apiBase = '/api'
const getToken = () => localStorage.getItem('token')

const fetchApi = async (url, options = {}) => {
  const res = await fetch(`${apiBase}${url}`, {
    ...options,
    headers: { 'Authorization': `Bearer ${getToken()}`, 'Content-Type': 'application/json', ...options.headers }
  })
  return res.json()
}

const loadMenus = async () => {
  loading.value = true
  sortChangedCount.value = 0
  sortMap.value = {}
  originalSortMap.value = {}
  try {
    const url = searchForm.keyword ? `/menus?keyword=${encodeURIComponent(searchForm.keyword)}` : '/menus'
    const data = await fetchApi(url)
    if (data.code === 0) {
      // 转换后端字段名 -> 前端字段名
      menuList.value = (data.data || []).map(m => ({
        id: m.id,
        parent_id: m.parent_id,
        name: m.menu_name,
        code: m.menu_code,
        icon: m.icon,
        path: m.route_path,
        component: m.component,
        permission: m.permission,
        type: m.menu_type,
        sort: m.sort,
        status: m.status,
        children: m.children,
      }))
      // Init sort maps
      menuList.value.forEach(m => {
        sortMap.value[m.id] = m.sort ?? 0
        originalSortMap.value[m.id] = m.sort ?? 0
      })
      expandedKeys.value = menuList.value.filter(m => m.children && m.children.length > 0).map(m => m.id)
    } else {
      Message.error(data.message || '加载失败')
    }
  } catch (e) {
    Message.error('加载菜单失败')
  } finally {
    loading.value = false
  }
}

const onSortChange = (id) => {
  const orig = originalSortMap.value[id] ?? 0
  const curr = sortMap.value[id] ?? 0
  if (orig !== curr) {
    sortChangedCount.value++
  } else {
    sortChangedCount.value = Math.max(0, sortChangedCount.value - 1)
  }
}

const batchSaveSort = async () => {
  const items = Object.entries(sortMap.value).map(([id, sort]) => ({
    id: Number(id), sort
  }))
  try {
    const data = await fetchApi('/menus/batch-sort', {
      method: 'PUT',
      body: JSON.stringify({ items })
    })
    if (data.code === 0) {
      Message.success('排序已保存')
      // Update originals
      Object.assign(originalSortMap.value, sortMap.value)
      sortChangedCount.value = 0
      loadMenus()
    } else {
      Message.error(data.message || '保存失败')
    }
  } catch (e) {
    Message.error('保存排序失败')
  }
}

// Move up/down among siblings
const canMoveUp = (record) => {
  const siblings = menuList.value.filter(m => m.parent_id === record.parent_id)
  siblings.sort((a, b) => (sortMap.value[a.id] ?? a.sort) - (sortMap.value[b.id] ?? b.sort))
  const idx = siblings.findIndex(m => m.id === record.id)
  return idx > 0
}

const canMoveDown = (record) => {
  const siblings = menuList.value.filter(m => m.parent_id === record.parent_id)
  siblings.sort((a, b) => (sortMap.value[a.id] ?? a.sort) - (sortMap.value[b.id] ?? b.sort))
  const idx = siblings.findIndex(m => m.id === record.id)
  return idx < siblings.length - 1
}

const moveUp = (record) => {
  const siblings = menuList.value.filter(m => m.parent_id === record.parent_id)
  siblings.sort((a, b) => (sortMap.value[a.id] ?? a.sort) - (sortMap.value[b.id] ?? b.sort))
  const idx = siblings.findIndex(m => m.id === record.id)
  if (idx > 0) {
    const currSort = sortMap.value[record.id] ?? 0
    const prevSort = sortMap.value[siblings[idx - 1].id] ?? 0
    sortMap.value[record.id] = prevSort
    sortMap.value[siblings[idx - 1].id] = currSort
    sortChangedCount.value = Object.keys(sortMap.value).filter(k => sortMap.value[k] !== originalSortMap.value[k]).length
  }
}

const moveDown = (record) => {
  const siblings = menuList.value.filter(m => m.parent_id === record.parent_id)
  siblings.sort((a, b) => (sortMap.value[a.id] ?? a.sort) - (sortMap.value[b.id] ?? b.sort))
  const idx = siblings.findIndex(m => m.id === record.id)
  if (idx < siblings.length - 1) {
    const currSort = sortMap.value[record.id] ?? 0
    const nextSort = sortMap.value[siblings[idx + 1].id] ?? 0
    sortMap.value[record.id] = nextSort
    sortMap.value[siblings[idx + 1].id] = currSort
    sortChangedCount.value = Object.keys(sortMap.value).filter(k => sortMap.value[k] !== originalSortMap.value[k]).length
  }
}

const doSearch = () => loadMenus()
const onExpand = (keys) => { expandedKeys.value = keys }
const getTypeLabel = (t) => ({ 0: '目录', 1: '菜单', 2: '按钮' })[t] || '菜单'
const getTypeColor = (t) => ({ 0: 'arcoblue', 1: 'green', 2: 'orange' })[t] || 'gray'
const isValidIcon = (icon) => icon && icon.startsWith('Icon')

const openCreateModal = (parentRecord) => {
  isEdit.value = false; currentEditRecord.value = null
  Object.assign(formData, { id: null, parent_id: parentRecord ? parentRecord.id : null, name: '', code: '', type: 1, icon: '', path: '', component: '', permission: '', sort: 0, status: 1 })
  formVisible.value = true
}

const openEditModal = (record) => {
  isEdit.value = true; currentEditRecord.value = record
  Object.assign(formData, { id: record.id, parent_id: record.parent_id || null, name: record.name, code: record.code || '', type: record.type ?? 1, icon: record.icon || '', path: record.path || '', component: record.component || '', permission: record.permission || '', sort: record.sort ?? 0, status: record.status ?? 1 })
  formVisible.value = true
}

const submitForm = async () => {
  if (!formData.name) { Message.warning('请输入菜单名称'); return }
  try {
    const url = isEdit.value ? `/menus/${formData.id}` : '/menus'
    const method = isEdit.value ? 'PUT' : 'POST'
    // 转换前端字段名 -> 后端字段名
    const payload = {
      menu_name: formData.name,
      menu_code: formData.code,
      icon: formData.icon,
      route_path: formData.path,
      component: formData.component,
      permission: formData.permission,
      menu_type: formData.type,
      sort: formData.sort,
      status: formData.status,
      parent_id: formData.parent_id === null ? 0 : formData.parent_id,
    }
    const data = await fetchApi(url, { method, body: JSON.stringify(payload) })
    if (data.code === 0) {
      Message.success(isEdit.value ? '更新成功' : '创建成功')
      formVisible.value = false
      loadMenus()
    } else {
      Message.error(data.message || '操作失败')
    }
  } catch (e) {
    Message.error('操作失败')
  }
}

const toggleStatus = async (record) => {
  const newStatus = record.status === 1 ? 0 : 1
  const data = await fetchApi(`/menus/${record.id}`, { method: 'PUT', body: JSON.stringify({ menu_name: record.name, menu_code: record.code, menu_type: record.type, status: newStatus, sort: record.sort, icon: record.icon, route_path: record.path, component: record.component, permission: record.permission, parent_id: record.parent_id }) })
  if (data.code === 0) { Message.success(newStatus === 1 ? '启用' : '禁用'); loadMenus() }
  else Message.error(data.message || '操作失败')
}

const handleDelete = (record) => {
  if (menuList.value.some(m => m.parent_id === record.id)) { Message.warning('存在子菜单，请先删除'); return }
  Modal.warning({ title: '确认删除', content: `确定删除「${record.name}」？`, okText: '删除', onOk: async () => {
    const data = await fetchApi(`/menus/${record.id}`, { method: 'DELETE' })
    if (data.code === 0) { Message.success('删除成功'); loadMenus() }
    else Message.error(data.message || '删除失败')
  }})
}

onMounted(() => { loadMenus() })
</script>

<style scoped>
.page-container { padding: 16px; }
.toolbar { display: flex; align-items: center; gap: 12px; margin-bottom: 16px; }
.drag-handle { cursor: grab; color: #888; font-size: 18px; user-select: none; padding: 0 4px; }
.drag-handle:active { cursor: grabbing; }
.sortable-row:hover .drag-handle { color: #165dff; }
.menu-name-cell { display: flex; align-items: center; gap: 8px; }
.menu-icon { font-size: 16px; }
</style>
