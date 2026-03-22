<template>
  <div class="page-container">
    <!-- 筛选栏 -->
    <a-card class="filter-card">
      <div class="filter-row">
        <a-input-search
          v-model="filter.keyword"
          placeholder="搜索API路径/名称"
          style="width: 260px"
          @search="handleSearch"
          @press-enter="handleSearch"
        />
        <a-select v-model="filter.method" placeholder="请求方法" style="width: 120px" allow-clear>
          <a-option value="GET">GET</a-option>
          <a-option value="POST">POST</a-option>
          <a-option value="PUT">PUT</a-option>
          <a-option value="DELETE">DELETE</a-option>
        </a-select>
        <a-select v-model="filter.status" placeholder="状态" style="width: 120px" allow-clear>
          <a-option :value="1">启用</a-option>
          <a-option :value="0">禁用</a-option>
        </a-select>
        <a-button type="primary" @click="handleSearch">
          <template #icon><icon-search /></template>
          查询
        </a-button>
        <a-button @click="handleReset">重置</a-button>
      </div>
    </a-card>

    <!-- API 权限列表 -->
    <a-card class="table-card">
      <template #title>
        <div class="table-title">
          <span>API 权限列表</span>
          <a-button type="primary" @click="openDrawer(null)">
            <template #icon><icon-plus /></template>
            新增API权限
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
        <template #method="{ record }">
          <a-tag :color="methodColor(record.method)">{{ record.method }}</a-tag>
        </template>
        <template #status="{ record }">
          <a-switch v-model="record._status" :checked-value="1" :unchecked-value="0" @change="handleToggleStatus(record)" />
        </template>
        <template #actions="{ record }">
          <a-button type="text" size="small" @click="openDrawer(record)">编辑</a-button>
          <a-divider direction="vertical" />
          <a-popconfirm content="确定删除该API权限？" @ok="handleDelete(record.id)">
            <a-button type="text" size="small" status="danger">删除</a-button>
          </a-popconfirm>
        </template>
      </a-table>
    </a-card>

    <!-- 创建/编辑 Drawer -->
    <a-drawer
      v-model:visible="drawerVisible"
      :title="isEdit ? '编辑API权限' : '新增API权限'"
      width="560px"
      @confirm="handleSubmit"
      @cancel="drawerVisible = false"
    >
      <a-form ref="formRef" :model="form" :rules="formRules" layout="vertical" label-align="left">
        <a-form-item label="API名称" field="api_name">
          <a-input v-model="form.api_name" placeholder="如：获取设备列表" />
        </a-form-item>
        <a-form-item label="请求方法" field="method">
          <a-select v-model="form.method" placeholder="请选择请求方法">
            <a-option value="GET">GET</a-option>
            <a-option value="POST">POST</a-option>
            <a-option value="PUT">PUT</a-option>
            <a-option value="DELETE">DELETE</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="API路径" field="api_path">
          <a-input v-model="form.api_path" placeholder="如：/api/v1/devices" />
        </a-form-item>
        <a-form-item label="所属模块" field="module">
          <a-select v-model="form.module" placeholder="请选择所属模块">
            <a-option value="device">设备管理</a-option>
            <a-option value="ota">OTA升级</a-option>
            <a-option value="alert">告警管理</a-option>
            <a-option value="system">系统设置</a-option>
            <a-option value="org">组织管理</a-option>
            <a-option value="permissions">权限管理</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="描述" field="description">
          <a-textarea v-model="form.description" placeholder="请输入描述" />
        </a-form-item>
        <a-form-item label="状态" field="status">
          <a-radio-group v-model="form.status">
            <a-radio :value="1">启用</a-radio>
            <a-radio :value="0">禁用</a-radio>
          </a-radio-group>
        </a-form-item>
      </a-form>
    </a-drawer>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import * as permApi from '@/api/permissions.js'

const loading = ref(false)
const drawerVisible = ref(false)
const isEdit = ref(false)
const editId = ref(null)
const formRef = ref(null)

const mockApis = [
  { id: 1, api_name: '获取设备列表', method: 'GET', api_path: '/api/v1/devices', module: 'device', description: '分页获取设备列表', status: 1 },
  { id: 2, api_name: '注册设备', method: 'POST', api_path: '/api/v1/devices', module: 'device', description: '注册新设备', status: 1 },
  { id: 3, api_name: '获取设备详情', method: 'GET', api_path: '/api/v1/devices/:id', module: 'device', description: '获取单个设备详情', status: 1 },
  { id: 4, api_name: '更新设备', method: 'PUT', api_path: '/api/v1/devices/:id', module: 'device', description: '更新设备信息', status: 1 },
  { id: 5, api_name: '删除设备', method: 'DELETE', api_path: '/api/v1/devices/:id', module: 'device', description: '删除设备', status: 0 },
  { id: 6, api_name: '获取固件包列表', method: 'GET', api_path: '/api/v1/ota/packages', module: 'ota', description: '获取OTA固件包列表', status: 1 },
  { id: 7, api_name: '创建部署任务', method: 'POST', api_path: '/api/v1/ota/deployments', module: 'ota', description: '创建OTA部署任务', status: 1 },
  { id: 8, api_name: '获取角色列表', method: 'GET', api_path: '/api/v1/permissions/roles', module: 'permissions', description: '获取角色列表', status: 1 },
  { id: 9, api_name: '创建角色', method: 'POST', api_path: '/api/v1/permissions/roles', module: 'permissions', description: '创建新角色', status: 1 }
]

const filter = reactive({
  keyword: '',
  method: undefined,
  status: undefined
})

const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0
})

const methodColor = (method) => {
  const map = { GET: 'green', POST: 'blue', PUT: 'orange', DELETE: 'red' }
  return map[method] || 'gray'
}

const columns = [
  { title: 'API名称', dataIndex: 'api_name', ellipsis: true },
  { title: '请求方法', slotName: 'method', width: 100 },
  { title: 'API路径', dataIndex: 'api_path', ellipsis: true },
  { title: '所属模块', dataIndex: 'module', width: 120, render: ({ record }) => {
    const map = { device: '设备管理', ota: 'OTA升级', alert: '告警管理', system: '系统设置', org: '组织管理', permissions: '权限管理' }
    return map[record.module] || record.module
  }},
  { title: '描述', dataIndex: 'description', ellipsis: true },
  { title: '启用', slotName: 'status', width: 80 },
  { title: '操作', slotName: 'actions', width: 140 }
]

const dataList = ref([])

const defaultForm = () => ({
  api_name: '',
  method: undefined,
  api_path: '',
  module: undefined,
  description: '',
  status: 1
})

const form = reactive(defaultForm())

const formRules = {
  api_name: [{ required: true, message: '请输入API名称' }],
  method: [{ required: true, message: '请选择请求方法' }],
  api_path: [{ required: true, message: '请输入API路径' }],
  module: [{ required: true, message: '请选择所属模块' }]
}

const loadData = async () => {
  loading.value = true
  try {
    const res = await permApi.getApiPermissions({
      page: pagination.current,
      page_size: pagination.pageSize,
      keyword: filter.keyword,
      method: filter.method,
      status: filter.status
    })
    if (res.code === 0) {
      dataList.value = res.data?.list || []
      pagination.total = res.data?.pagination?.total || 0
    }
  } catch {
    let list = [...mockApis]
    if (filter.keyword) {
      list = list.filter(a => a.api_name.includes(filter.keyword) || a.api_path.includes(filter.keyword))
    }
    if (filter.method) {
      list = list.filter(a => a.method === filter.method)
    }
    if (filter.status !== undefined) {
      list = list.filter(a => a.status === filter.status)
    }
    dataList.value = list.map(a => ({ ...a, _status: a.status }))
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
  filter.method = undefined
  filter.status = undefined
  handleSearch()
}

const handleTableChange = (pag) => {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  loadData()
}

const handleToggleStatus = (record) => {
  record.status = record._status
  permApi.updateApiPermission(record.id, { status: record.status }).catch(() => {
    record._status = record.status === 1 ? 0 : 1
    record.status = record._status
  })
}

const openDrawer = (record) => {
  if (record) {
    isEdit.value = true
    editId.value = record.id
    Object.assign(form, {
      api_name: record.api_name,
      method: record.method,
      api_path: record.api_path,
      module: record.module,
      description: record.description,
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
      res = await permApi.updateApiPermission(editId.value, { ...form })
    } else {
      res = await permApi.createApiPermission({ ...form })
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
    const res = await permApi.deleteApiPermission(id)
    if (res.code === 0) {
      Message.success('删除成功')
      loadData()
    }
  } catch {
    Message.success('删除成功（模拟）')
    loadData()
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.page-container { display: flex; flex-direction: column; gap: 16px; }
.filter-row { display: flex; gap: 12px; align-items: center; flex-wrap: wrap; }
.table-title { display: flex; justify-content: space-between; align-items: center; }
</style>
