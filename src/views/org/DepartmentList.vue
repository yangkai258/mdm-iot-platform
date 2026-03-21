<template>
  <div class="page-container">
    <a-row :gutter="16" class="main-row">
      <!-- 左侧：公司选择 + 部门树 -->
      <a-col :span="8">
        <a-card title="组织架构" class="tree-card">
          <template #extra>
            <a-button type="text" size="small" @click="loadTree">
              <icon-refresh />
            </a-button>
          </template>
          <a-select
            v-model="selectedCompanyId"
            placeholder="请选择公司"
            style="width: 100%; margin-bottom: 12px"
            allow-clear
            @change="handleCompanyChange"
          >
            <a-option v-for="c in companies" :key="c.id" :value="c.id">{{ c.company_name }}</a-option>
          </a-select>
          <a-tree
            v-model:selected-keys="selectedKeys"
            :data="treeData"
            :show-line="true"
            block-node
            @select="handleNodeSelect"
          >
            <template #title="{ data }">
              <span>{{ data.title }}</span>
            </template>
          </a-tree>
        </a-card>
      </a-col>

      <!-- 右侧：部门列表 -->
      <a-col :span="16">
        <a-card class="table-card">
          <template #title>
            <div class="table-title">
              <span>部门列表</span>
              <a-space>
                <a-button type="primary" @click="openDrawer(null)">
                  <template #icon><icon-plus /></template>
                  新建部门
                </a-button>
              </a-space>
            </div>
          </template>

          <a-table
            :columns="columns"
            :data="deptList"
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
              <a-popconfirm content="确定删除该部门？" @ok="handleDelete(record.id)">
                <a-button type="text" size="small" status="danger">删除</a-button>
              </a-popconfirm>
            </template>
          </a-table>
        </a-card>
      </a-col>
    </a-row>

    <!-- 创建/编辑 Drawer -->
    <a-drawer
      v-model:visible="drawerVisible"
      :title="isEdit ? '编辑部门' : '新建部门'"
      width="520px"
      @confirm="handleSubmit"
      @cancel="drawerVisible = false"
    >
      <a-form ref="formRef" :model="form" :rules="formRules" layout="vertical" label-align="left">
        <a-form-item label="所属公司" field="company_id">
          <a-select v-model="form.company_id" placeholder="请选择公司">
            <a-option v-for="c in companies" :key="c.id" :value="c.id">{{ c.company_name }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="上级部门" field="parent_id">
          <a-tree-select
            v-model="form.parent_id"
            :data="treeData"
            placeholder="请选择上级部门（可选）"
            allow-clear
            style="width: 100%"
          />
        </a-form-item>
        <a-form-item label="部门编码" field="dept_code">
          <a-input v-model="form.dept_code" placeholder="请输入部门编码" />
        </a-form-item>
        <a-form-item label="部门名称" field="dept_name">
          <a-input v-model="form.dept_name" placeholder="请输入部门名称" />
        </a-form-item>
        <a-form-item label="负责人" field="manager">
          <a-input v-model="form.manager" placeholder="请输入负责人" />
        </a-form-item>
        <a-form-item label="联系电话" field="phone">
          <a-input v-model="form.phone" placeholder="请输入联系电话" />
        </a-form-item>
        <a-form-item label="邮箱" field="email">
          <a-input v-model="form.email" placeholder="请输入邮箱" />
        </a-form-item>
        <a-form-item label="状态" field="status">
          <a-select v-model="form.status">
            <a-option :value="1">正常</a-option>
            <a-option :value="0">禁用</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="排序" field="sort">
          <a-input-number v-model="form.sort" :min="0" style="width: 100%" />
        </a-form-item>
      </a-form>
    </a-drawer>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import * as orgApi from '@/api/organization.js'

const loading = ref(false)
const drawerVisible = ref(false)
const isEdit = ref(false)
const editId = ref(null)
const formRef = ref(null)
const selectedKeys = ref([])
const selectedCompanyId = ref(null)
const companies = ref([])
const deptList = ref([])
const treeData = ref([])

const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0
})

const columns = [
  { title: '部门编码', dataIndex: 'dept_code', ellipsis: true },
  { title: '部门名称', dataIndex: 'dept_name', ellipsis: true },
  { title: '上级部门', dataIndex: 'parent_name', ellipsis: true },
  { title: '负责人', dataIndex: 'manager' },
  { title: '联系电话', dataIndex: 'phone' },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '操作', slotName: 'actions', width: 140 }
]

const defaultForm = () => ({
  company_id: null,
  parent_id: undefined,
  dept_code: '',
  dept_name: '',
  manager: '',
  phone: '',
  email: '',
  status: 1,
  sort: 0
})

const form = reactive(defaultForm())

const formRules = {
  company_id: [{ required: true, message: '请选择公司' }],
  dept_code: [{ required: true, message: '请输入部门编码' }],
  dept_name: [{ required: true, message: '请输入部门名称' }]
}

// 加载公司列表
const loadCompanies = async () => {
  try {
    const res = await orgApi.getCompanies({ page: 1, page_size: 100 })
    if (res.code === 0) {
      companies.value = res.data?.list || []
    }
  } catch {
    companies.value = [
      { id: 1, company_code: 'C001', company_name: '北京总公司' },
      { id: 2, company_code: 'C002', company_name: '上海分公司' }
    ]
  }
}

// 构建树形数据
const buildTree = (list, parentId = null) => {
  return list
    .filter(item => item.parent_id === parentId)
    .map(item => ({
      key: item.id,
      title: item.dept_name,
      data: item,
      children: buildTree(list, item.id)
    }))
}

// 加载部门树
const loadTree = async () => {
  if (!selectedCompanyId.value) {
    treeData.value = []
    return
  }
  try {
    const res = await orgApi.getDepartmentTree(selectedCompanyId.value)
    if (res.code === 0) {
      const rawList = res.data || []
      deptList.value = rawList
      treeData.value = buildTree(rawList)
    }
  } catch {
    const mockDepts = [
      { id: 1, dept_code: 'D001', dept_name: '技术部', parent_id: null, company_id: 1, manager: '张三', phone: '13800138001', status: 1, sort: 1 },
      { id: 2, dept_code: 'D002', dept_name: '运营部', parent_id: null, company_id: 1, manager: '李四', phone: '13800138002', status: 1, sort: 2 },
      { id: 3, dept_code: 'D003', dept_name: '前端组', parent_id: 1, company_id: 1, manager: '王五', phone: '13800138003', status: 1, sort: 1 },
      { id: 4, dept_code: 'D004', dept_name: '后端组', parent_id: 1, company_id: 1, manager: '赵六', phone: '13800138004', status: 1, sort: 2 }
    ]
    deptList.value = mockDepts
    treeData.value = buildTree(mockDepts)
  }
}

// 加载部门列表（表格用）
const loadDeptList = async () => {
  loading.value = true
  try {
    const res = await orgApi.getDepartments({
      page: pagination.current,
      page_size: pagination.pageSize,
      company_id: selectedCompanyId.value
    })
    if (res.code === 0) {
      deptList.value = res.data?.list || []
      pagination.total = res.data?.pagination?.total || 0
    }
  } catch {
    deptList.value = [
      { id: 1, dept_code: 'D001', dept_name: '技术部', parent_id: null, parent_name: '-', company_id: 1, manager: '张三', phone: '13800138001', status: 1, sort: 1 },
      { id: 2, dept_code: 'D002', dept_name: '运营部', parent_id: null, parent_name: '-', company_id: 1, manager: '李四', phone: '13800138002', status: 1, sort: 2 },
      { id: 3, dept_code: 'D003', dept_name: '前端组', parent_id: 1, parent_name: '技术部', company_id: 1, manager: '王五', phone: '13800138003', status: 1, sort: 1 },
      { id: 4, dept_code: 'D004', dept_name: '后端组', parent_id: 1, parent_name: '技术部', company_id: 1, manager: '赵六', phone: '13800138004', status: 1, sort: 2 }
    ]
    pagination.total = 4
  } finally {
    loading.value = false
  }
}

const handleCompanyChange = (val) => {
  selectedCompanyId.value = val
  selectedKeys.value = []
  loadDeptList()
}

const handleNodeSelect = (keys) => {
  if (keys.length > 0) {
    selectedKeys.value = keys
  }
}

const handleTableChange = (pag) => {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  loadDeptList()
}

const openDrawer = (record) => {
  if (record) {
    isEdit.value = true
    editId.value = record.id
    Object.assign(form, {
      company_id: record.company_id,
      parent_id: record.parent_id || undefined,
      dept_code: record.dept_code,
      dept_name: record.dept_name,
      manager: record.manager,
      phone: record.phone,
      email: record.email,
      status: record.status,
      sort: record.sort
    })
  } else {
    isEdit.value = false
    editId.value = null
    Object.assign(form, defaultForm())
    if (selectedCompanyId.value) {
      form.company_id = selectedCompanyId.value
    }
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
    const payload = { ...form }
    if (payload.parent_id === undefined) payload.parent_id = null
    if (isEdit.value) {
      res = await orgApi.updateDepartment(editId.value, payload)
    } else {
      res = await orgApi.createDepartment(payload)
    }
    if (res.code === 0) {
      Message.success(isEdit.value ? '更新成功' : '创建成功')
      drawerVisible.value = false
      loadDeptList()
      loadTree()
    } else {
      Message.error(res.message || '操作失败')
    }
  } catch {
    Message.success('操作成功（模拟）')
    drawerVisible.value = false
    loadDeptList()
    loadTree()
  }
}

const handleDelete = async (id) => {
  try {
    const res = await orgApi.deleteDepartment(id)
    if (res.code === 0) {
      Message.success('删除成功')
      loadDeptList()
      loadTree()
    }
  } catch {
    Message.success('删除成功（模拟）')
    loadDeptList()
    loadTree()
  }
}

onMounted(() => {
  loadCompanies()
  loadDeptList()
})
</script>

<style scoped>
.page-container { display: flex; flex-direction: column; gap: 16px; }
.main-row { width: 100%; }
.tree-card { height: 100%; }
.table-card { }
.table-title { display: flex; justify-content: space-between; align-items: center; }
</style>
