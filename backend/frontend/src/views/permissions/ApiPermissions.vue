<template>
  <div class="page-container">
    <!-- 面包屑 + 搜索栏 -->
    <div class="page-header">
      <div class="header-left">
        <a-breadcrumb>
          <a-breadcrumb-item><a href="#/dashboard">首页</a></a-breadcrumb-item>
          <a-breadcrumb-item>系统管理</a-breadcrumb-item>
          <a-breadcrumb-item>API权限管理</a-breadcrumb-item>
        </a-breadcrumb>
      </div>
      <div class="header-right">
        <a-input-search
          v-model="searchForm.keyword"
          placeholder="API路径/名称"
          style="width: 280px"
          @search="doSearch"
        />
      </div>
    </div>

    <!-- 操作按钮栏（靠左） -->
    <div class="action-bar">
      <a-space :size="12">
        <a-button type="primary" @click="openCreateModal">「新建」</a-button>
        <a-button @click="handleBatchImport">「批量导入」</a-button>
        <a-button @click="handleBatchExport">「批量导出」</a-button>
        <a-button @click="loadData">「刷新」</a-button>
      </a-space>
    </div>

    <!-- 数据表格 -->
    <a-card :bordered="false" style="border-radius: 4px">
      <a-table
        :columns="columns"
        :data="tableData"
        :loading="loading"
        :pagination="pagination"
        @change="handleTableChange"
        row-key="id"
      >
        <template #method="{ record }">
          <a-tag :color="methodColor(record.method)">
            {{ record.method }}
          </a-tag>
        </template>
        <template #status="{ record }">
          <a-tag :color="record.status === 1 ? 'green' : 'gray'">
            {{ record.status === 1 ? '启用' : '禁用' }}
          </a-tag>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="openEditModal(record)">编辑</a-button>
            <a-button type="text" size="small" status="danger" @click="handleDelete(record)">删除</a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 新建/编辑全屏模态 -->
    <a-modal
      v-model:visible="formVisible"
      :title="isEdit ? '编辑API权限' : '新建API权限'"
      :width="600"
      :mask-closable="false"
      @ok="handleSubmit"
      @cancel="formVisible = false"
    >
      <a-form :model="form" layout="vertical">
        <a-form-item field="apiPath" label="API路径" required>
          <a-input v-model="form.apiPath" placeholder="如 /api/v1/users" allow-clear />
        </a-form-item>
        <a-form-item field="apiName" label="API名称" required>
          <a-input v-model="form.apiName" placeholder="如 用户列表" allow-clear />
        </a-form-item>
        <a-form-item field="method" label="请求方法" required>
          <a-select v-model="form.method" placeholder="选择请求方法">
            <a-option value="GET">GET</a-option>
            <a-option value="POST">POST</a-option>
            <a-option value="PUT">PUT</a-option>
            <a-option value="DELETE">DELETE</a-option>
          </a-select>
        </a-form-item>
        <a-form-item field="path" label="路径">
          <a-input v-model="form.path" placeholder="路径描述" allow-clear />
        </a-form-item>
        <a-form-item field="permissionCode" label="权限码" required>
          <a-input v-model="form.permissionCode" placeholder="如 system:user:list" allow-clear />
        </a-form-item>
        <a-form-item field="menuId" label="关联菜单">
          <a-input v-model="form.menuId" placeholder="关联菜单ID" allow-clear />
        </a-form-item>
        <a-form-item field="status" label="状态">
          <a-switch v-model="form.status" checked-value="1" unchecked-value="0" />
          <span style="margin-left: 8px">{{ form.status == '1' ? '启用' : '禁用' }}</span>
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 批量导入弹窗 -->
    <a-modal
      v-model:visible="importVisible"
      title="批量导入API权限"
      :width="500"
      :mask-closable="false"
      @ok="handleImportConfirm"
      @cancel="importVisible = false"
    >
      <a-upload
        ref="uploadRef"
        :auto-upload="false"
        :limit="1"
        accept=".xlsx,.xls,.csv"
        @change="handleFileChange"
      >
        <template #upload-button>
          <div>
            <icon-upload />
            <div style="margin-top: 8px">点击选择文件 或 拖拽文件到此处</div>
            <div style="color: #86909c; font-size: 12px; margin-top: 4px">支持 .xlsx/.xls/.csv 格式</div>
          </div>
        </template>
      </a-upload>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

// 搜索表单
const searchForm = reactive({
  keyword: ''
})

// 表格数据（模拟）
const tableData = ref([
  {
    id: 1,
    apiPath: '/api/v1/users',
    apiName: '用户列表',
    method: 'GET',
    path: '/users',
    permissionCode: 'system:user:list',
    menuId: '系统管理',
    status: 1
  },
  {
    id: 2,
    apiPath: '/api/v1/users',
    apiName: '创建用户',
    method: 'POST',
    path: '/users',
    permissionCode: 'system:user:create',
    menuId: '系统管理',
    status: 1
  },
  {
    id: 3,
    apiPath: '/api/v1/users/:id',
    apiName: '更新用户',
    method: 'PUT',
    path: '/users/:id',
    permissionCode: 'system:user:update',
    menuId: '系统管理',
    status: 1
  },
  {
    id: 4,
    apiPath: '/api/v1/users/:id',
    apiName: '删除用户',
    method: 'DELETE',
    path: '/users/:id',
    permissionCode: 'system:user:delete',
    menuId: '系统管理',
    status: 0
  },
  {
    id: 5,
    apiPath: '/api/v1/roles',
    apiName: '角色列表',
    method: 'GET',
    path: '/roles',
    permissionCode: 'system:role:list',
    menuId: '权限管理',
    status: 1
  }
])

const loading = ref(false)

// 表格列定义
const columns = [
  {
    title: 'API路径',
    dataIndex: 'apiPath',
    width: 200,
    ellipsis: true
  },
  {
    title: 'API名称',
    dataIndex: 'apiName',
    width: 120,
    ellipsis: true
  },
  {
    title: '方法',
    dataIndex: 'method',
    width: 90,
    slotName: 'method'
  },
  {
    title: '路径',
    dataIndex: 'path',
    width: 120,
    ellipsis: true
  },
  {
    title: '权限码',
    dataIndex: 'permissionCode',
    width: 180,
    ellipsis: true
  },
  {
    title: '关联菜单',
    dataIndex: 'menuId',
    width: 120,
    ellipsis: true
  },
  {
    title: '状态',
    dataIndex: 'status',
    width: 80,
    slotName: 'status'
  },
  {
    title: '操作',
    width: 140,
    slotName: 'actions',
    fixed: 'right'
  }
]

// 分页
const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 5,
  showTotal: true,
  showPageSize: true,
  showJumper: true
})

// 方法Tag颜色
const methodColor = (method) => {
  const colors = {
    GET: 'green',
    POST: 'blue',
    PUT: 'oranges',
    DELETE: 'red'
  }
  return colors[method] || 'gray'
}

// 新建/编辑
const formVisible = ref(false)
const isEdit = ref(false)
const form = reactive({
  id: null,
  apiPath: '',
  apiName: '',
  method: '',
  path: '',
  permissionCode: '',
  menuId: '',
  status: '1'
})

// 批量导入
const importVisible = ref(false)
const uploadFile = ref(null)

const doSearch = () => {
  pagination.current = 1
  loadData()
}

const loadData = () => {
  loading.value = true
  // 模拟加载
  setTimeout(() => {
    loading.value = false
  }, 500)
}

const handleTableChange = (pag) => {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  loadData()
}

const openCreateModal = () => {
  isEdit.value = false
  Object.assign(form, {
    id: null,
    apiPath: '',
    apiName: '',
    method: '',
    path: '',
    permissionCode: '',
    menuId: '',
    status: '1'
  })
  formVisible.value = true
}

const openEditModal = (record) => {
  isEdit.value = true
  Object.assign(form, {
    id: record.id,
    apiPath: record.apiPath,
    apiName: record.apiName,
    method: record.method,
    path: record.path,
    permissionCode: record.permissionCode,
    menuId: record.menuId,
    status: String(record.status)
  })
  formVisible.value = true
}

const handleSubmit = () => {
  if (!form.apiPath || !form.apiName || !form.method || !form.permissionCode) {
    Message.warning('请填写必填项')
    return
  }
  // 模拟提交
  Message.success(isEdit.value ? '编辑成功' : '创建成功')
  formVisible.value = false
  loadData()
}

const handleDelete = (record) => {
  Message.success(`删除API权限：${record.apiName}`)
  loadData()
}

const handleBatchImport = () => {
  importVisible.value = true
}

const handleBatchExport = () => {
  Message.info('正在导出...')
}

const handleFileChange = (fileList) => {
  if (fileList.length > 0) {
    uploadFile.value = fileList[0]
  }
}

const handleImportConfirm = () => {
  if (!uploadFile.value) {
    Message.warning('请选择导入文件')
    return
  }
  Message.success('导入成功')
  importVisible.value = false
  loadData()
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.page-container {
  background-color: #f2f3f5;
  min-height: 100vh;
  padding: 16px;
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}

.header-left {
  flex: 1;
}

.header-right {
  display: flex;
  align-items: center;
}

.action-bar {
  display: flex;
  align-items: center;
  margin-bottom: 12px;
}

:deep(.arco-card) {
  background: #ffffff;
  border-radius: 4px;
}
</style>
