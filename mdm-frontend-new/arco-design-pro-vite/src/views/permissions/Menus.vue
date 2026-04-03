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
            <a-option value="IconHome">IconHome</a-option>
            <a-option value="IconUser">IconUser</a-option>
            <a-option value="IconUserGroup">IconUserGroup</a-option>
            <a-option value="IconSettings">IconSettings</a-option>
            <a-option value="IconMenu">IconMenu</a-option>
            <a-option value="IconDashboard">IconDashboard</a-option>
            <a-option value="IconList">IconList</a-option>
            <a-option value="IconFile">IconFile</a-option>
            <a-option value="IconFolder">IconFolder</a-option>
            <a-option value="IconSearch">IconSearch</a-option>
            <a-option value="IconPlus">IconPlus</a-option>
            <a-option value="IconEdit">IconEdit</a-option>
            <a-option value="IconDelete">IconDelete</a-option>
            <a-option value="IconCheck">IconCheck</a-option>
            <a-option value="IconClose">IconClose</a-option>
            <a-option value="IconRefresh">IconRefresh</a-option>
            <a-option value="IconExport">IconExport</a-option>
            <a-option value="IconImport">IconImport</a-option>
            <a-option value="IconDownload">IconDownload</a-option>
            <a-option value="IconUpload">IconUpload</a-option>
            <a-option value="IconArrowUp">IconArrowUp</a-option>
            <a-option value="IconArrowDown">IconArrowDown</a-option>
            <a-option value="IconLeft">IconLeft</a-option>
            <a-option value="IconRight">IconRight</a-option>
            <a-option value="IconSafe">IconSafe</a-option>
            <a-option value="IconLock">IconLock</a-option>
            <a-option value="IconUnlock">IconUnlock</a-option>
            <a-option value="IconKey">IconKey</a-option>
            <a-option value="IconBug">IconBug</a-option>
            <a-option value="IconTool">IconTool</a-option>
            <a-option value="IconConnection">IconConnection</a-option>
            <a-option value="IconStar">IconStar</a-option>
            <a-option value="IconStarFill">IconStarFill</a-option>
            <a-option value="IconHeart">IconHeart</a-option>
            <a-option value="IconHeartFill">IconHeartFill</a-option>
            <a-option value="IconMessage">IconMessage</a-option>
            <a-option value="IconMessageFill">IconMessageFill</a-option>
            <a-option value="IconBell">IconBell</a-option>
            <a-option value="IconBellFill">IconBellFill</a-option>
            <a-option value="IconCalendar">IconCalendar</a-option>
            <a-option value="IconHistory">IconHistory</a-option>
            <a-option value="IconClockCircle">IconClockCircle</a-option>
            <a-option value="IconLocation">IconLocation</a-option>
            <a-option value="IconPhone">IconPhone</a-option>
            <a-option value="IconEmail">IconEmail</a-option>
            <a-option value="IconPaper">IconPaper</a-option>
            <a-option value="IconCamera">IconCamera</a-option>
            <a-option value="IconSave">IconSave</a-option>
            <a-option value="IconMind">IconMind</a-option>
            <a-option value="IconWifi">IconWifi</a-option>
            <a-option value="IconHotspot">IconHotspot</a-option>
            <a-option value="IconDevice">IconDevice</a-option>
            <a-option value="IconRobot">IconRobot</a-option>
            <a-option value="IconPet">IconPet</a-option>
          </a-select>
        </a-form-item>
        <a-form-item field="path" label="路由路径">
          <a-select v-model="formData.path" placeholder="选择或输入路由" allow-clear allow-search>
            <a-option value="/dashboard">/dashboard</a-option>
            <a-option value="/device">/device</a-option>
            <a-option value="/devices">/devices</a-option>
            <a-option value="/member">/member</a-option>
            <a-option value="/members">/members</a-option>
            <a-option value="/order">/order</a-option>
            <a-option value="/orders">/orders</a-option>
            <a-option value="/store">/store</a-option>
            <a-option value="/stores">/stores</a-option>
            <a-option value="/ai">/ai</a-option>
            <a-option value="/pet">/pet</a-option>
            <a-option value="/ota">/ota</a-option>
            <a-option value="/alert">/alert</a-option>
            <a-option value="/alerts">/alerts</a-option>
            <a-option value="/notification">/notification</a-option>
            <a-option value="/notifications">/notifications</a-option>
            <a-option value="/permission-manage">/permission-manage</a-option>
            <a-option value="/permissions">/permissions</a-option>
            <a-option value="/system">/system</a-option>
            <a-option value="/tenant">/tenant</a-option>
            <a-option value="/tenants">/tenants</a-option>
            <a-option value="/org">/org</a-option>
            <a-option value="/billing">/billing</a-option>
            <a-option value="/subscription">/subscription</a-option>
            <a-option value="/analytics">/analytics</a-option>
            <a-option value="/knowledge">/knowledge</a-option>
            <a-option value="/app">/app</a-option>
            <a-option value="/apps">/apps</a-option>
            <a-option value="/market">/market</a-option>
            <a-option value="/marketing">/marketing</a-option>
            <a-option value="/content">/content</a-option>
            <a-option value="/webhook">/webhook</a-option>
            <a-option value="/webhooks">/webhooks</a-option>
          </a-select>
        </a-form-item>
        <a-form-item field="component" label="组件路径">
          <a-select v-model="formData.component" placeholder="选择组件" allow-clear allow-search>
            <a-option value="/dashboard/Dashboard.vue">/dashboard/Dashboard.vue</a-option>
            <a-option value="/device/DeviceList.vue">/device/DeviceList.vue</a-option>
            <a-option value="/device/DeviceDetail.vue">/device/DeviceDetail.vue</a-option>
            <a-option value="/member/MemberList.vue">/member/MemberList.vue</a-option>
            <a-option value="/member/MemberDetail.vue">/member/MemberDetail.vue</a-option>
            <a-option value="/ota/FirmwareList.vue">/ota/FirmwareList.vue</a-option>
            <a-option value="/ota/VersionList.vue">/ota/VersionList.vue</a-option>
            <a-option value="/alerts/AlertsView.vue">/alerts/AlertsView.vue</a-option>
            <a-option value="/notifications/NotificationList.vue">/notifications/NotificationList.vue</a-option>
            <a-option value="/permissions/ApiPermissions.vue">/permissions/ApiPermissions.vue</a-option>
            <a-option value="/permissions/Menus.vue">/permissions/Menus.vue</a-option>
            <a-option value="/permissions/Roles.vue">/permissions/Roles.vue</a-option>
            <a-option value="/permissions/PermissionGroups.vue">/permissions/PermissionGroups.vue</a-option>
            <a-option value="/system/SystemSettings.vue">/system/SystemSettings.vue</a-option>
            <a-option value="/system/Logs.vue">/system/Logs.vue</a-option>
            <a-option value="/tenant/TenantList.vue">/tenant/TenantList.vue</a-option>
            <a-option value="/org/OrgList.vue">/org/OrgList.vue</a-option>
            <a-option value="/ai/AiConsole.vue">/ai/AiConsole.vue</a-option>
            <a-option value="/ai/ModelList.vue">/ai/ModelList.vue</a-option>
            <a-option value="/pet/PetConsole.vue">/pet/PetConsole.vue</a-option>
            <a-option value="/pet/PetProfile.vue">/pet/PetProfile.vue</a-option>
            <a-option value="/store/StoreList.vue">/store/StoreList.vue</a-option>
            <a-option value="/billing/BillingList.vue">/billing/BillingList.vue</a-option>
            <a-option value="/analytics/AnalyticsDashboard.vue">/analytics/AnalyticsDashboard.vue</a-option>
            <a-option value="/knowledge/KnowledgeBase.vue">/knowledge/KnowledgeBase.vue</a-option>
            <a-option value="/content/ContentList.vue">/content/ContentList.vue</a-option>
            <a-option value="/webhooks/WebhookList.vue">/webhooks/WebhookList.vue</a-option>
            <a-option value="/list/BasicList.vue">/list/BasicList.vue</a-option>
            <a-option value="/list/CardList.vue">/list/CardList.vue</a-option>
            <a-option value="/form/BasicForm.vue">/form/BasicForm.vue</a-option>
            <a-option value="/profile/ProfileBasic.vue">/profile/ProfileBasic.vue</a-option>
          </a-select>
        </a-form-item>
        <a-form-item field="permission" label="权限标识">
          <a-select v-model="formData.permission" placeholder="选择或输入权限" allow-clear allow-search>
            <a-option value="system:view">system:view - 查看系统</a-option>
            <a-option value="system:create">system:create - 创建系统</a-option>
            <a-option value="system:update">system:update - 更新系统</a-option>
            <a-option value="system:delete">system:delete - 删除系统</a-option>
            <a-option value="system:export">system:export - 导出系统</a-option>
            <a-option value="system:import">system:import - 导入系统</a-option>
            <a-option value="menu:view">menu:view - 查看菜单</a-option>
            <a-option value="menu:create">menu:create - 创建菜单</a-option>
            <a-option value="menu:update">menu:update - 更新菜单</a-option>
            <a-option value="menu:delete">menu:delete - 删除菜单</a-option>
            <a-option value="device:view">device:view - 查看设备</a-option>
            <a-option value="device:create">device:create - 创建设备</a-option>
            <a-option value="device:update">device:update - 更新设备</a-option>
            <a-option value="device:delete">device:delete - 删除设备</a-option>
            <a-option value="device:control">device:control - 控制设备</a-option>
            <a-option value="member:view">member:view - 查看会员</a-option>
            <a-option value="member:create">member:create - 创建会员</a-option>
            <a-option value="member:update">member:update - 更新会员</a-option>
            <a-option value="member:delete">member:delete - 删除会员</a-option>
            <a-option value="ota:view">ota:view - 查看OTA</a-option>
            <a-option value="ota:create">ota:create - 创建OTA</a-option>
            <a-option value="ota:update">ota:update - 更新OTA</a-option>
            <a-option value="ota:delete">ota:delete - 删除OTA</a-option>
            <a-option value="alert:view">alert:view - 查看告警</a-option>
            <a-option value="alert:create">alert:create - 创建告警</a-option>
            <a-option value="alert:update">alert:update - 更新告警</a-option>
            <a-option value="alert:delete">alert:delete - 删除告警</a-option>
            <a-option value="role:view">role:view - 查看角色</a-option>
            <a-option value="role:create">role:create - 创建角色</a-option>
            <a-option value="role:update">role:update - 更新角色</a-option>
            <a-option value="role:delete">role:delete - 删除角色</a-option>
            <a-option value="tenant:view">tenant:view - 查看租户</a-option>
            <a-option value="tenant:create">tenant:create - 创建租户</a-option>
            <a-option value="tenant:update">tenant:update - 更新租户</a-option>
            <a-option value="tenant:delete">tenant:delete - 删除租户</a-option>
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
