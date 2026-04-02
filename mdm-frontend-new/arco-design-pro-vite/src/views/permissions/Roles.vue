<template>
    <Breadcrumb :items="['Home','Console','']" />


  <div class="page-container">
    <!-- 筛选栏 -->
    <a-card class="general-card">
      <div class="filter-row">
        <a-input-search
          v-model="filter.keyword"
          placeholder="搜索角色名称/编码"
          style="width: 260px"
          @search="handleSearch"
          @press-enter="handleSearch"
        />
        <a-select
          v-model="filter.status"
          placeholder="角色状态"
          style="width: 140px"
          allow-clear
          @change="loadData"
        >
          <a-option :value="1">正常</a-option>
          <a-option :value="0">禁用</a-option>
        </a-select>
        <a-button type="primary" @click="handleSearch">查询</a-button>
        <a-button @click="handleReset">重置</a-button>
      </div>
    </a-card>

    <!-- 角色列表 -->
    <a-card class="general-card">
      <template #title>
        <div class="table-title">
          <span>角色列表</span>
          <a-button type="primary" @click="openDrawer(null)">
            <template #icon><icon-plus /></template>
            新建角色
          </a-button>
        </div>
      </template>

      <a-table
        :columns="columns"
        :data="dataList"
        :loading="loading"
        :pagination="pagination"
        row-key="id"
        @change="handleTableChange"
      >
        <template #status="{ record }">
          <a-tag :color="record.status === 1 ? 'green' : 'gray'">
            {{ record.status === 1 ? '正常' : '禁用' }}
          </a-tag>
        </template>
        <template #is_system="{ record }">
          <a-tag :color="record.is_system ? 'blue' : 'default'">
            {{ record.is_system ? '系统' : '自定义' }}
          </a-tag>
        </template>
        <template #actions="{ record }">
          <a-button type="text" size="small" @click="openPermissionModal(record)">分配权限</a-button>
          <a-divider direction="vertical" />
          <a-button type="text" size="small" @click="openDrawer(record)">编辑</a-button>
          <a-divider direction="vertical" />
          <a-popconfirm content="确定删除该角色？" @ok="handleDelete(record.id)">
            <a-button type="text" size="small" status="danger" :disabled="record.is_system">删除</a-button>
          </a-popconfirm>
        </template>
      </a-table>
    </a-card>

    <!-- 创建/编辑 Drawer -->
    <a-drawer
      v-model:visible="drawerVisible"
      :title="isEdit ? '编辑角色' : '新建角色'"
      width="520px"
      @confirm="handleSubmit"
      @cancel="drawerVisible = false"
    >
      <a-form ref="formRef" :model="form" :rules="formRules" layout="vertical" label-align="left">
        <a-form-item label="角色编码" field="role_code">
          <a-input v-model="form.role_code" placeholder="请输入角色编码，如 admin" />
        </a-form-item>
        <a-form-item label="角色名称" field="role_name">
          <a-input v-model="form.role_name" placeholder="请输入角色名称" />
        </a-form-item>
        <a-form-item label="角色描述" field="description">
          <a-textarea v-model="form.description" placeholder="请输入角色描述" />
        </a-form-item>
        <a-form-item label="排序" field="sort">
          <a-input-number v-model="form.sort" :min="0" style="width: 100%" />
        </a-form-item>
        <a-form-item label="状态" field="status">
          <a-radio-group v-model="form.status">
            <a-radio :value="1">正常</a-radio>
            <a-radio :value="0">禁用</a-radio>
          </a-radio-group>
        </a-form-item>
      </a-form>
    </a-drawer>

    <!-- 分配权限 Modal -->
    <a-modal
      v-model:visible="permModalVisible"
      title="分配权限"
      width="600px"
      @before-ok="handlePermSubmit"
      @cancel="permModalVisible = false"
    >
      <div class="perm-tree-wrapper">
        <a-spin v-if="permLoading" />
        <a-tree
          v-else
          v-model:checked-keys="selectedPermKeys"
          :data="permTreeData"
          :checkable="true"
          :show-line="true"
          node-key="id"
          @check="handlePermCheck"
        />
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import * as permApi from '@/api/permissions.js'

const loading = ref(false)
const permLoading = ref(false)
const drawerVisible = ref(false)
const permModalVisible = ref(false)
const isEdit = ref(false)
const editId = ref(null)
const formRef = ref(null)
const permTreeData = ref([])
const selectedPermKeys = ref([])
const currentRoleId = ref(null)

// 模拟数据
const mockRoles = [
  { id: 1, role_code: 'super_admin', role_name: '超级管理员', description: '拥有系统所有权限', status: 1, sort: 1, is_system: true },
  { id: 2, role_code: 'admin', role_name: '管理员', description: '拥有管理权限', status: 1, sort: 2, is_system: true },
  { id: 3, role_code: 'operator', role_name: '运维人员', description: '负责设备运维管理', status: 1, sort: 3, is_system: false },
  { id: 4, role_code: 'viewer', role_name: '访客', description: '只读权限', status: 0, sort: 4, is_system: false }
]

const mockPermTree = [
  {
    id: 'p1', title: '设备管理',
    children: [
      { id: 'p1-1', title: '设备列表' },
      { id: 'p1-2', title: '设备详情' },
      { id: 'p1-3', title: '设备注册' }
    ]
  },
  {
    id: 'p2', title: 'OTA管理',
    children: [
      { id: 'p2-1', title: '固件包管理' },
      { id: 'p2-2', title: '部署任务' }
    ]
  },
  {
    id: 'p3', title: '告警管理',
    children: [
      { id: 'p3-1', title: '告警规则' },
      { id: 'p3-2', title: '告警列表' }
    ]
  },
  {
    id: 'p4', title: '系统设置',
    children: [
      { id: 'p4-1', title: '用户管理' },
      { id: 'p4-2', title: '角色管理' },
      { id: 'p4-3', title: '菜单管理' }
    ]
  }
]

const filter = reactive({
  keyword: '',
  status: undefined
})

const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0
})

const columns = [
  { title: '角色编码', dataIndex: 'role_code', ellipsis: true, width: 140 },
  { title: '角色名称', dataIndex: 'role_name', ellipsis: true },
  { title: '描述', dataIndex: 'description', ellipsis: true },
  { title: '类型', slotName: 'is_system', width: 90 },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '排序', dataIndex: 'sort', width: 70 },
  { title: '操作', slotName: 'actions', width: 200 }
]

const dataList = ref([])

const defaultForm = () => ({
  role_code: '',
  role_name: '',
  description: '',
  sort: 0,
  status: 1
})

const form = reactive(defaultForm())

const formRules = {
  role_code: [{ required: true, message: '请输入角色编码' }],
  role_name: [{ required: true, message: '请输入角色名称' }]
}

const loadData = async () => {
  loading.value = true
  try {
    const res = await permApi.getRoles({
      page: pagination.current,
      page_size: pagination.pageSize,
      keyword: filter.keyword,
      status: filter.status
    })
    if (res.code === 0) {
      dataList.value = res.data?.list || []
      pagination.total = res.data?.pagination?.total || 0
    }
  } catch {
    // 使用模拟数据
    let list = [...mockRoles]
    if (filter.keyword) {
      list = list.filter(r => r.role_name.includes(filter.keyword) || r.role_code.includes(filter.keyword))
    }
    if (filter.status !== undefined) {
      list = list.filter(r => r.status === filter.status)
    }
    dataList.value = list
    pagination.total = list.length
    Message.warning('使用模拟数据')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.current = 1
  loadData()
}

const handleReset = () => {
  filter.keyword = ''
  filter.status = undefined
  handleSearch()
}

const handleTableChange = (pag) => {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  loadData()
}

const openDrawer = (record) => {
  if (record) {
    isEdit.value = true
    editId.value = record.id
    Object.assign(form, {
      role_code: record.role_code,
      role_name: record.role_name,
      description: record.description,
      sort: record.sort,
      status: record.status
    })
  } else {
    isEdit.value = false
    editId.value = null
    Object.assign(form, defaultForm())
  }
  drawerVisible.value = true
}

const handleSubmit = async () => {
  try {
    await formRef.value?.validate()
  } catch {
    return
  }
  try {
    let res
    if (isEdit.value) {
      res = await permApi.updateRole(editId.value, { ...form })
    } else {
      res = await permApi.createRole({ ...form })
    }
    if (res.code === 0) {
      Message.success(isEdit.value ? '更新成功' : '创建成功')
      drawerVisible.value = false
      loadData()
    }
  } catch {
    Message.success('操作成功（模拟）')
    drawerVisible.value = false
    loadData()
  }
}

const handleDelete = async (id) => {
  try {
    const res = await permApi.deleteRole(id)
    if (res.code === 0) {
      Message.success('删除成功')
      loadData()
    }
  } catch {
    Message.success('删除成功（模拟）')
    loadData()
  }
}

const loadPermTree = async () => {
  permLoading.value = true
  try {
    const res = await permApi.getAllPermissions()
    if (res.code === 0) {
      permTreeData.value = res.data || []
    }
  } catch {
    permTreeData.value = mockPermTree
  } finally {
    permLoading.value = false
  }
}

const loadRolePerms = async (roleId) => {
  try {
    const res = await permApi.getRolePermissions(roleId)
    if (res.code === 0) {
      selectedPermKeys.value = res.data || []
    }
  } catch {
    // 模拟已有权限
    selectedPermKeys.value = roleId === 1 ? ['p1', 'p1-1', 'p2', 'p2-1', 'p3', 'p3-1'] : []
  }
}

const openPermissionModal = async (record) => {
  currentRoleId.value = record.id
  permModalVisible.value = true
  selectedPermKeys.value = []
  await loadPermTree()
  await loadRolePerms(record.id)
}

const handlePermCheck = (checked) => {
  selectedPermKeys.value = checked
}

const handlePermSubmit = async (done) => {
  try {
    await permApi.assignPermissions(currentRoleId.value, selectedPermKeys.value)
    Message.success('权限分配成功')
  } catch {
    Message.success('权限分配成功（模拟）')
  }
  done(true)
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.page-container { padding: 16px; display: flex; flex-direction: column; gap: 16px; }
.filter-row { display: flex; gap: 12px; align-items: center; flex-wrap: wrap; }
.table-title { display: flex; justify-content: space-between; align-items: center; }
.perm-tree-wrapper { max-height: 500px; overflow-y: auto; padding: 8px 0; }
</style>
