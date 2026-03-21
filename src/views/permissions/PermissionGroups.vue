<template>
  <div class="page-container">
    <!-- 筛选栏 -->
    <a-card class="filter-card">
      <div class="filter-row">
        <a-input-search
          v-model="filter.keyword"
          placeholder="搜索权限组名称"
          style="width: 260px"
          @search="handleSearch"
          @press-enter="handleSearch"
        />
        <a-button type="primary" @click="handleSearch">
          <template #icon><icon-search /></template>
          查询
        </a-button>
        <a-button @click="handleReset">重置</a-button>
      </div>
    </a-card>

    <!-- 权限组列表 -->
    <a-card class="table-card">
      <template #title>
        <div class="table-title">
          <span>权限组列表</span>
          <a-button type="primary" @click="openDrawer(null)">
            <template #icon><icon-plus /></template>
            新建权限组
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
        <template #actions="{ record }">
          <a-button type="text" size="small" @click="openDrawer(record)">编辑</a-button>
          <a-divider direction="vertical" />
          <a-popconfirm content="确定删除该权限组？" @ok="handleDelete(record.id)">
            <a-button type="text" size="small" status="danger">删除</a-button>
          </a-popconfirm>
        </template>
      </a-table>
    </a-card>

    <!-- 创建/编辑 Drawer -->
    <a-drawer
      v-model:visible="drawerVisible"
      :title="isEdit ? '编辑权限组' : '新建权限组'"
      width="520px"
      @confirm="handleSubmit"
      @cancel="drawerVisible = false"
    >
      <a-form ref="formRef" :model="form" :rules="formRules" layout="vertical" label-align="left">
        <a-form-item label="权限组编码" field="group_code">
          <a-input v-model="form.group_code" placeholder="请输入权限组编码" />
        </a-form-item>
        <a-form-item label="权限组名称" field="group_name">
          <a-input v-model="form.group_name" placeholder="请输入权限组名称" />
        </a-form-item>
        <a-form-item label="所属模块" field="module">
          <a-select v-model="form.module" placeholder="请选择所属模块">
            <a-option value="device">设备管理</a-option>
            <a-option value="ota">OTA升级</a-option>
            <a-option value="alert">告警管理</a-option>
            <a-option value="system">系统设置</a-option>
            <a-option value="org">组织管理</a-option>
            <a-option value="policy">策略管理</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="描述" field="description">
          <a-textarea v-model="form.description" placeholder="请输入描述" />
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

const mockGroups = [
  { id: 1, group_code: 'device_view', group_name: '设备查看', module: 'device', description: '设备查看权限', status: 1, sort: 1 },
  { id: 2, group_code: 'device_manage', group_name: '设备管理', module: 'device', description: '设备增删改查权限', status: 1, sort: 2 },
  { id: 3, group_code: 'ota_view', group_name: 'OTA查看', module: 'ota', description: 'OTA固件包查看权限', status: 1, sort: 3 },
  { id: 4, group_code: 'ota_deploy', group_name: 'OTA部署', module: 'ota', description: 'OTA升级部署权限', status: 1, sort: 4 },
  { id: 5, group_code: 'alert_view', group_name: '告警查看', module: 'alert', description: '告警查看权限', status: 0, sort: 5 }
]

const filter = reactive({
  keyword: ''
})

const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0
})

const columns = [
  { title: '权限组编码', dataIndex: 'group_code', ellipsis: true, width: 150 },
  { title: '权限组名称', dataIndex: 'group_name', ellipsis: true },
  { title: '所属模块', dataIndex: 'module', width: 120, render: ({ record }) => {
    const map = { device: '设备管理', ota: 'OTA升级', alert: '告警管理', system: '系统设置', org: '组织管理', policy: '策略管理' }
    return map[record.module] || record.module
  }},
  { title: '描述', dataIndex: 'description', ellipsis: true },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '排序', dataIndex: 'sort', width: 70 },
  { title: '操作', slotName: 'actions', width: 140 }
]

const dataList = ref([])

const defaultForm = () => ({
  group_code: '',
  group_name: '',
  module: undefined,
  description: '',
  sort: 0,
  status: 1
})

const form = reactive(defaultForm())

const formRules = {
  group_code: [{ required: true, message: '请输入权限组编码' }],
  group_name: [{ required: true, message: '请输入权限组名称' }],
  module: [{ required: true, message: '请选择所属模块' }]
}

const loadData = async () => {
  loading.value = true
  try {
    const res = await permApi.getPermissionGroups({
      page: pagination.current,
      page_size: pagination.pageSize,
      keyword: filter.keyword
    })
    if (res.code === 0) {
      dataList.value = res.data?.list || []
      pagination.total = res.data?.pagination?.total || 0
    }
  } catch {
    let list = [...mockGroups]
    if (filter.keyword) {
      list = list.filter(g => g.group_name.includes(filter.keyword) || g.group_code.includes(filter.keyword))
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
    Object.assign(form, record)
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
      res = await permApi.updatePermissionGroup(editId.value, { ...form })
    } else {
      res = await permApi.createPermissionGroup({ ...form })
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
    const res = await permApi.deletePermissionGroup(id)
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
