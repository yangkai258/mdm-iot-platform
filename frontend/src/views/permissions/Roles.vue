<template>
  <div class="page-container">
    <!-- 面包屑 + 搜索栏 -->
    <div class="page-header">
      <div class="header-left">
        <a-breadcrumb>
          <a-breadcrumb-item><a href="#/dashboard">首页</a></a-breadcrumb-item>
          <a-breadcrumb-item>权限管理</a-breadcrumb-item>
          <a-breadcrumb-item>角色管理</a-breadcrumb-item>
        </a-breadcrumb>
      </div>
      <div class="header-right">
        <a-input-search
          v-model="searchForm.keyword"
          placeholder="角色名称/编码"
          style="width: 240px"
          @search="doSearch"
        />
      </div>
    </div>

    <!-- 操作按钮栏（靠左） -->
    <div class="action-bar">
      <a-space :size="12">
        <a-button type="primary" @click="openCreateModal">「新建」</a-button>
        <a-button @click="toggleSearch">「筛选」</a-button>
        <a-button @click="loadRoles">「刷新」</a-button>
      </a-space>
    </div>

    <!-- 筛选面板 -->
    <a-card v-if="showSearch" :bordered="false" style="margin-bottom: 12px">
      <a-form :model="searchForm" layout="inline">
        <a-form-item field="keyword" label="角色名称/编码">
          <a-input v-model="searchForm.keyword" placeholder="请输入关键词" allow-clear style="width: 200px" />
        </a-form-item>
        <a-form-item field="status" label="状态">
          <a-select v-model="searchForm.status" placeholder="选择状态" allow-clear style="width: 120px">
            <a-option :value="1">启用</a-option>
            <a-option :value="0">禁用</a-option>
          </a-select>
        </a-form-item>
        <a-form-item>
          <a-button type="primary" @click="doSearch">查询</a-button>
          <a-button style="margin-left: 8px" @click="resetSearch">重置</a-button>
        </a-form-item>
      </a-form>
    </a-card>

      <!-- 角色列表 -->
      <a-table
        :columns="columns"
        :data="roles"
        :loading="loading"
        :pagination="pagination"
        @change="handleTableChange"
        row-key="id"
        style="margin-top: 16px"
      >
        <template #status="{ record }">
          <a-tag :color="record.status === 1 ? 'green' : 'gray'">
            {{ record.status === 1 ? '启用' : '禁用' }}
          </a-tag>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="openEditModal(record)">编辑</a-button>
            <a-button type="text" size="small" @click="openPermModal(record)">权限</a-button>
            <a-button type="text" size="small" status="danger" @click="handleDelete(record)">删除</a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 创建/编辑角色全屏模态 -->
    <a-modal
      v-model:visible="formVisible"
      :title="isEdit ? '编辑角色' : '创建角色'"
      :width="600"
      :mask-closable="false"
      @ok="submitForm"
      @cancel="formVisible = false"
    >
      <a-form :model="formData" layout="vertical">
        <a-form-item field="name" label="角色名称" required>
          <a-input v-model="formData.name" placeholder="请输入角色名称" />
        </a-form-item>
        <a-form-item field="code" label="角色编码" :required="!isEdit">
          <a-input v-model="formData.code" placeholder="请输入角色编码，如 admin" :disabled="isEdit" />
        </a-form-item>
        <a-form-item field="description" label="描述">
          <a-textarea v-model="formData.description" placeholder="请输入角色描述" :rows="3" />
        </a-form-item>
        <a-form-item field="sort" label="排序">
          <a-input-number v-model="formData.sort" :min="0" :max="9999" />
        </a-form-item>
        <a-form-item field="status" label="状态">
          <a-switch v-model="formData.status" :checked-value="1" :unchecked-value="0" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 权限分配弹窗 -->
    <a-modal
      v-model:visible="permVisible"
      title="分配权限"
      :width="700"
      :mask-closable="false"
      @ok="submitPerms"
      @cancel="permVisible = false"
    >
      <a-alert type="info" style="margin-bottom: 16px">
        当前角色: <b>{{ currentRole?.name }}</b>
      </a-alert>
      <a-tabs default-active-tab="tree">
        <a-tab-pane key="tree" title="权限树">
          <a-tree
            v-model:selected-keys="selectedPerms"
            v-model:checked-keys="selectedPerms"
            :data="permTreeData"
            :selectable="false"
            checkable
            :default-expand-all="true"
            field-names="{ key: 'id', title: 'name', children: 'children' }"
          />
        </a-tab-pane>
        <a-tab-pane key="list" title="权限列表">
          <a-checkbox-group v-model="selectedPerms">
            <a-space direction="vertical" fill>
              <div v-for="group in permGroups" :key="group.id">
                <div class="perm-group-title">{{ group.name }}</div>
                <a-space wrap>
                  <a-checkbox v-for="p in group.permissions" :key="p" :value="p">{{ p }}</a-checkbox>
                </a-space>
              </div>
            </a-space>
          </a-checkbox-group>
        </a-tab-pane>
      </a-tabs>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message, Modal } from '@arco-design/web-vue'

const loading = ref(false)
const showSearch = ref(false)
const formVisible = ref(false)
const permVisible = ref(false)
const isEdit = ref(false)
const roles = ref([])
const permTreeData = ref([])
const selectedPerms = ref([])
const currentRole = ref(null)

const searchForm = reactive({
  keyword: '',
  status: ''
})

const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0
})

const formData = reactive({
  id: null,
  name: '',
  code: '',
  description: '',
  sort: 0,
  status: 1
})

const columns = [
  { title: 'ID', dataIndex: 'id', width: 80 },
  { title: '角色名称', dataIndex: 'name', width: 150 },
  { title: '角色编码', dataIndex: 'code', width: 150 },
  { title: '描述', dataIndex: 'description', ellipsis: true },
  { title: '排序', dataIndex: 'sort', width: 80 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '创建时间', dataIndex: 'created_at', width: 180 },
  { title: '操作', slotName: 'actions', width: 200, fixed: 'right' }
]

const permGroups = ref([
  { id: 1, name: '租户管理', permissions: ['tenant:view', 'tenant:manage'] },
  { id: 2, name: '用户管理', permissions: ['user:view', 'user:manage'] },
  { id: 3, name: '设备管理', permissions: ['device:view', 'device:manage', 'device:control'] },
  { id: 4, name: 'OTA管理', permissions: ['ota:view', 'ota:deploy'] },
  { id: 5, name: '告警管理', permissions: ['alert:view', 'alert:manage'] },
  { id: 6, name: '会员管理', permissions: ['member:view', 'member:manage'] },
  { id: 7, name: '策略管理', permissions: ['policy:view', 'policy:manage'] },
  { id: 8, name: '通知管理', permissions: ['notification:view', 'notification:manage'] },
  { id: 9, name: '应用管理', permissions: ['app:view', 'app:manage'] },
  { id: 10, name: '系统管理', permissions: ['system:view', 'system:manage', 'role:manage'] },
  { id: 11, name: '数据操作', permissions: ['data:export', 'data:import'] },
])

const apiBase = '/api/v1'

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

const loadRoles = async () => {
  loading.value = true
  try {
    const params = new URLSearchParams({
      page: pagination.current,
      page_size: pagination.pageSize
    })
    if (searchForm.keyword) params.set('keyword', searchForm.keyword)
    if (searchForm.status !== '') params.set('status', searchForm.status)

    const data = await fetchApi(`/roles?${params}`)
    if (data.code === 0) {
      roles.value = data.data?.list || []
      pagination.total = data.data?.total || 0
    }
  } catch (e) {
    Message.error('加载角色列表失败')
  } finally {
    loading.value = false
  }
}

const loadPermTree = async () => {
  try {
    const data = await fetchApi('/permissions')
    if (data.code === 0) {
      permTreeData.value = data.data || []
    }
  } catch (e) {
    console.error('加载权限树失败', e)
  }
}

const toggleSearch = () => {
  showSearch.value = !showSearch.value
}

const doSearch = () => {
  pagination.current = 1
  loadRoles()
}

const resetSearch = () => {
  searchForm.keyword = ''
  searchForm.status = ''
  pagination.current = 1
  loadRoles()
}

const openCreateModal = () => {
  isEdit.value = false
  formData.id = null
  formData.name = ''
  formData.code = ''
  formData.description = ''
  formData.sort = 0
  formData.status = 1
  formVisible.value = true
}

const openEditModal = (record) => {
  isEdit.value = true
  Object.assign(formData, {
    id: record.id,
    name: record.name,
    code: record.code,
    description: record.description,
    sort: record.sort,
    status: record.status
  })
  formVisible.value = true
}

const submitForm = async () => {
  if (!formData.name) {
    Message.warning('请输入角色名称')
    return
  }
  if (!isEdit.value && !formData.code) {
    Message.warning('请输入角色编码')
    return
  }

  try {
    const url = isEdit.value ? `/roles/${formData.id}` : '/roles'
    const method = isEdit.value ? 'PUT' : 'POST'
    const data = await fetchApi(url, {
      method,
      body: JSON.stringify(formData)
    })
    if (data.code === 0) {
      Message.success(isEdit.value ? '更新成功' : '创建成功')
      formVisible.value = false
      loadRoles()
    } else {
      Message.error(data.message || '操作失败')
    }
  } catch (e) {
    Message.error('操作失败')
  }
}

const openPermModal = async (record) => {
  currentRole.value = record
  permVisible.value = true
  selectedPerms.value = []

  // 加载角色已有权限
  try {
    const data = await fetchApi(`/roles/${record.id}/permissions`)
    if (data.code === 0) {
      selectedPerms.value = data.data?.permission_ids || []
    }
  } catch (e) {
    console.error('加载角色权限失败', e)
  }

  // 确保权限树已加载
  if (permTreeData.value.length === 0) {
    await loadPermTree()
  }
}

const submitPerms = async () => {
  if (!currentRole.value) return

  try {
    // 将选中权限的code转为ID（这里简化处理，前端传code列表，后端处理）
    const permIds = selectedPerms.value.map(p => typeof p === 'number' ? p : p)
    const data = await fetchApi(`/roles/${currentRole.value.id}/permissions`, {
      method: 'POST',
      body: JSON.stringify({ permission_ids: permIds })
    })
    if (data.code === 0) {
      Message.success('权限分配成功')
      permVisible.value = false
    } else {
      Message.error(data.message || '分配失败')
    }
  } catch (e) {
    Message.error('权限分配失败')
  }
}

const handleDelete = (record) => {
  Modal.warning({
    title: '确认删除',
    content: `确定要删除角色「${record.name}」吗？删除后不可恢复。`,
    okText: '删除',
    onOk: async () => {
      try {
        const data = await fetchApi(`/roles/${record.id}`, { method: 'DELETE' })
        if (data.code === 0) {
          Message.success('删除成功')
          loadRoles()
        } else {
          Message.error(data.message || '删除失败')
        }
      } catch (e) {
        Message.error('删除失败')
      }
    }
  })
}

const handleTableChange = (pag) => {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  loadRoles()
}

onMounted(() => {
  loadRoles()
  loadPermTree()
})
</script>

<style scoped>
.page-container {
  padding: 20px;
}
.breadcrumb {
  margin-bottom: 12px;
}
.search-bar {
  margin-bottom: 12px;
}
.action-bar {
  margin-bottom: 16px;
}
.search-area {
  padding: 16px;
  background: #f7f8fa;
  border-radius: 4px;
  margin-bottom: 16px;
}
.perm-group-title {
  font-weight: bold;
  margin-bottom: 8px;
  color: var(--color-text-1);
}
</style>
;
}
.header-left {
  display: flex;
  align-items: center;
}
.header-right {
  display: flex;
  align-items: center;
}
.breadcrumb {
  margin-bottom: 16px;
}
.search-bar {
  margin-bottom: 16px;
}
.action-bar {
  margin-bottom: 16px;
}
.search-area {
  padding: 16px;
  background: #f7f8fa;
  border-radius: 4px;
  margin-bottom: 16px;
}
.perm-group-title {
  font-weight: bold;
  margin-bottom: 8px;
  color: var(--color-text-1);
}
</style>
