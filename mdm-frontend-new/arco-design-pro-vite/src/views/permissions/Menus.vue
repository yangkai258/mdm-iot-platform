<template>
    <Breadcrumb :items="['Home','Console','']" />


  <div class="page-container">
    <a-card class="general-card" title="菜单管理">
      <template #extra>
        <a-space :size="12">
          <a-button type="primary" @click="openCreateModal(null)"><icon-plus />新建</a-button>
          <a-button @click="loadMenus"><icon-refresh />刷新</a-button>
        </a-space>
      </template>
      <div class="search-bar">
        <a-input-search
          v-model="searchForm.keyword"
          placeholder="菜单名称"
          style="width: 280px"
          @search="doSearch"
          allow-clear
        />
      </div>
      <!-- 卡片容器 -->
      <!-- 树形表格 -->
      <a-table
        :columns="columns"
        :data="menuTree"
        :loading="loading"
        :expanded-keys="expandedKeys"
        :pagination="false"
        row-key="id"
        @expand="onExpand"
      >
        <template #menu_name="{ record }">
          <span class="menu-name-cell">
            <component :is="record.icon" v-if="record.icon && isValidIcon(record.icon)" class="menu-icon" />
            <span>{{ record.name }}</span>
          </span>
        </template>
        <template #menu_type="{ record }">
          <a-tag :color="getTypeColor(record.type)">
            {{ getTypeLabel(record.type) }}
          </a-tag>
        </template>
        <template #status="{ record }">
          <a-tag :color="record.status === 1 ? 'green' : 'gray'">
            {{ record.status === 1 ? '启用' : '禁用' }}
          </a-tag>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="openCreateModal(record)">新建子级</a-button>
            <a-button type="text" size="small" @click="openEditModal(record)">编辑</a-button>
            <a-button type="text" size="small" @click="toggleStatus(record)">
              {{ record.status === 1 ? '禁用' : '启用' }}
            </a-button>
            <a-button type="text" size="small" status="danger" @click="handleDelete(record)">删除</a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 创建/编辑菜单弹窗 -->
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
            placeholder="请选择上级菜单（不选则为顶级）"
            allow-clear
            :field-names="{ key: 'id', title: 'name', children: 'children' }"
            style="width: 100%"
          />
        </a-form-item>
        <a-form-item field="name" label="菜单名称" required>
          <a-input v-model="formData.name" placeholder="请输入菜单名称" />
        </a-form-item>
        <a-form-item field="code" label="菜单编码">
          <a-input v-model="formData.code" placeholder="请输入菜单编码，如 system:menu" />
        </a-form-item>
        <a-form-item field="type" label="菜单类型" required>
          <a-radio-group v-model="formData.type">
            <a-radio :value="0">目录</a-radio>
            <a-radio :value="1">菜单</a-radio>
            <a-radio :value="2">按钮</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item field="icon" label="图标">
          <a-input v-model="formData.icon" placeholder="请输入图标名称，如 IconHome" />
        </a-form-item>
        <a-form-item field="path" label="路由路径">
          <a-input v-model="formData.path" placeholder="请输入路由路径，如 /system/menus" />
        </a-form-item>
        <a-form-item field="component" label="组件路径">
          <a-input v-model="formData.component" placeholder="请输入组件路径，如 /views/system/Menus.vue" />
        </a-form-item>
        <a-form-item field="permission" label="权限标识">
          <a-input v-model="formData.permission" placeholder="请输入权限标识，如 system:menu:view" />
        </a-form-item>
        <a-form-item field="sort" label="排序">
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
import { IconPlus, IconRefresh } from '@arco-design/web-vue/es/icon'

const loading = ref(false)
const formVisible = ref(false)
const isEdit = ref(false)
const menuList = ref([])
const expandedKeys = ref([])
const currentEditRecord = ref(null)

const searchForm = reactive({
  keyword: ''
})

const formData = reactive({
  id: null,
  parent_id: null,
  name: '',
  code: '',
  type: 1,
  icon: '',
  path: '',
  component: '',
  permission: '',
  sort: 0,
  status: 1
})

// 扁平菜单转树形
const buildTree = (list, parentId = 0) => {
  return list.filter(item => item.parent_id === parentId).map(item => ({
    ...item,
    children: buildTree(list, item.id)
  })).sort((a, b) => a.sort - b.sort)
}

const menuTree = computed(() => buildTree(menuList.value))

// 用于下拉选择的树（包含自己作为父级的情况，需要排除自己和子级）
const getDisabledIds = (record) => {
  const ids = [record.id]
  if (record.children) {
    record.children.forEach(child => {
      ids.push(...getDisabledIds(child))
    })
  }
  return ids
}

const menuTreeForSelect = computed(() => {
  if (!isEdit.value) return menuTree.value
  // 编辑时排除自己和子级
  const disabledIds = currentEditRecord.value ? getDisabledIds(currentEditRecord.value) : []
  const filterTree = (list) => list.filter(item => !disabledIds.includes(item.id)).map(item => ({
    ...item,
    children: item.children ? filterTree(item.children) : []
  }))
  return filterTree(menuTree.value)
})

const columns = [
  { title: '菜单名称', slotName: 'menu_name', minWidth: 200 },
  { title: '菜单编码', dataIndex: 'code', width: 180 },
  { title: '图标', dataIndex: 'icon', width: 120 },
  { title: '路由路径', dataIndex: 'path', width: 200, ellipsis: true },
  { title: '菜单类型', slotName: 'menu_type', width: 100 },
  { title: '排序', dataIndex: 'sort', width: 80 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '操作', slotName: 'actions', width: 260, fixed: 'right' }
]

const apiBase = '/api'

const getToken = () => localStorage.getItem('token')

const fetchApi = async (url, options = {}) => {
  const res = await fetch(`${apiBase}${url}`, {
    ...options,
    headers: {
      'Authorization': `Bearer ${getToken()}`,
      'Content-Type': 'application/json',
      ...options.headers
    }
  })
  return res.json()
}

const loadMenus = async () => {
  loading.value = true
  try {
    const params = new URLSearchParams()
    if (searchForm.keyword) params.set('keyword', searchForm.keyword)

    const url = params.toString() ? `/menus?${params}` : '/menus'
    const data = await fetchApi(url)
    if (data.code === 0) {
      menuList.value = data.data || []
      // 默认展开前两层
      expandedKeys.value = menuList.value.filter(m => m.children && m.children.length > 0).map(m => m.id)
    } else {
      Message.error(data.message || '加载菜单列表失败')
    }
  } catch (e) {
    Message.error('加载菜单列表失败')
  } finally {
    loading.value = false
  }
}

const doSearch = () => {
  loadMenus()
}

const onExpand = (keys) => {
  expandedKeys.value = keys
}

const getTypeLabel = (type) => {
  const labels = { 0: '目录', 1: '菜单', 2: '按钮' }
  return labels[type] || '菜单'
}

const getTypeColor = (type) => {
  const colors = { 0: 'arcoblue', 1: 'green', 2: 'orange' }
  return colors[type] || 'gray'
}

const isValidIcon = (icon) => {
  // 简单判断是否看起来像图标名称
  return icon && icon.startsWith('Icon')
}

const openCreateModal = (parentRecord) => {
  isEdit.value = false
  currentEditRecord.value = null
  Object.assign(formData, {
    id: null,
    parent_id: parentRecord ? parentRecord.id : null,
    name: '',
    code: '',
    type: 1,
    icon: '',
    path: '',
    component: '',
    permission: '',
    sort: 0,
    status: 1
  })
  formVisible.value = true
}

const openEditModal = (record) => {
  isEdit.value = true
  currentEditRecord.value = record
  Object.assign(formData, {
    id: record.id,
    parent_id: record.parent_id || null,
    name: record.name,
    code: record.code || '',
    type: record.type ?? 1,
    icon: record.icon || '',
    path: record.path || '',
    component: record.component || '',
    permission: record.permission || '',
    sort: record.sort ?? 0,
    status: record.status ?? 1
  })
  formVisible.value = true
}

const submitForm = async () => {
  if (!formData.name) {
    Message.warning('请输入菜单名称')
    return
  }

  try {
    const url = isEdit.value ? `/menus/${formData.id}` : '/menus'
    const method = isEdit.value ? 'PUT' : 'POST'
    const payload = { ...formData }
    if (payload.parent_id === null) payload.parent_id = 0

    const data = await fetchApi(url, {
      method,
      body: JSON.stringify(payload)
    })
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
  try {
    const data = await fetchApi(`/menus/${record.id}`, {
      method: 'PUT',
      body: JSON.stringify({ ...record, status: newStatus })
    })
    if (data.code === 0) {
      Message.success(newStatus === 1 ? '启用成功' : '禁用成功')
      loadMenus()
    } else {
      Message.error(data.message || '操作失败')
    }
  } catch (e) {
    Message.error('操作失败')
  }
}

const handleDelete = (record) => {
  // 检查是否有子菜单
  const hasChildren = menuList.value.some(m => m.parent_id === record.id)
  if (hasChildren) {
    Message.warning('该菜单存在子菜单，请先删除子菜单')
    return
  }

  Modal.warning({
    title: '确认删除',
    content: `确定要删除菜单「${record.name}」吗？删除后不可恢复。`,
    okText: '删除',
    onOk: async () => {
      try {
        const data = await fetchApi(`/menus/${record.id}`, { method: 'DELETE' })
        if (data.code === 0) {
          Message.success('删除成功')
          loadMenus()
        } else {
          Message.error(data.message || '删除失败')
        }
      } catch (e) {
        Message.error('删除失败')
      }
    }
  })
}

onMounted(() => {
  loadMenus()
})
</script>

<style scoped>
.page-container { padding: 16px; }
.search-bar { margin-bottom: 16px; }
.menu-name-cell { display: flex; align-items: center; gap: 8px; }
.menu-icon { font-size: 16px; color: var(--color-text-2); }
</style>
