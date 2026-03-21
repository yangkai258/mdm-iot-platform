<template>
  <div class="page-container">
    <a-card class="filter-card">
      <div class="filter-row">
        <a-select
          v-model="filter.company_id"
          placeholder="所属公司"
          style="width: 160px"
          allow-clear
          @change="handleCompanyChange"
        >
          <a-option v-for="c in companies" :key="c.id" :value="c.id">{{ c.company_name }}</a-option>
        </a-select>
        <a-select
          v-model="filter.dept_id"
          placeholder="所属部门"
          style="width: 160px"
          allow-clear
          @change="handleSearch"
        >
          <a-option v-for="d in departments" :key="d.id" :value="d.id">{{ d.dept_name }}</a-option>
        </a-select>
        <a-select
          v-model="filter.emp_status"
          placeholder="员工状态"
          style="width: 140px"
          allow-clear
          @change="handleSearch"
        >
          <a-option value="在职">在职</a-option>
          <a-option value="离职">离职</a-option>
          <a-option value="退休">退休</a-option>
        </a-select>
        <a-input-search
          v-model="filter.keyword"
          placeholder="搜索工号/姓名/手机号"
          style="width: 240px"
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

    <a-card class="table-card">
      <template #title>
        <div class="table-title">
          <span>员工列表</span>
          <a-space>
            <a-button @click="handleExport">
              <template #icon><icon-download /></template>
              导出
            </a-button>
            <a-button type="primary" @click="openDrawer(null)">
              <template #icon><icon-plus /></template>
              新建员工
            </a-button>
          </a-space>
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
        <template #gender="{ record }">
          {{ record.gender === '男' ? '男' : record.gender === '女' ? '女' : '-' }}
        </template>
        <template #emp_status="{ record }">
          <a-tag :color="getEmpStatusColor(record.emp_status)">
            {{ record.emp_status || '在职' }}
          </a-tag>
        </template>
        <template #status="{ record }">
          <a-tag :color="record.status === 1 ? 'green' : 'gray'">
            {{ record.status === 1 ? '正常' : '禁用' }}
          </a-tag>
        </template>
        <template #actions="{ record }">
          <a-button type="text" size="small" @click="openDrawer(record)">编辑</a-button>
          <a-divider direction="vertical" />
          <a-popconfirm content="确定删除该员工？" @ok="handleDelete(record.id)">
            <a-button type="text" size="small" status="danger">删除</a-button>
          </a-popconfirm>
        </template>
      </a-table>
    </a-card>

    <!-- 创建/编辑 Drawer -->
    <a-drawer
      v-model:visible="drawerVisible"
      :title="isEdit ? '编辑员工' : '新建员工'"
      width="560px"
      @confirm="handleSubmit"
      @cancel="drawerVisible = false"
    >
      <a-form ref="formRef" :model="form" :rules="formRules" layout="vertical" label-align="left">
        <a-form-item label="工号" field="emp_code">
          <a-input v-model="form.emp_code" placeholder="请输入工号" />
        </a-form-item>
        <a-form-item label="姓名" field="emp_name">
          <a-input v-model="form.emp_name" placeholder="请输入姓名" />
        </a-form-item>
        <a-form-item label="性别" field="gender">
          <a-radio-group v-model="form.gender">
            <a-radio value="男">男</a-radio>
            <a-radio value="女">女</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item label="手机号" field="phone">
          <a-input v-model="form.phone" placeholder="请输入手机号" />
        </a-form-item>
        <a-form-item label="邮箱" field="email">
          <a-input v-model="form.email" placeholder="请输入邮箱" />
        </a-form-item>
        <a-form-item label="身份证号" field="id_card">
          <a-input v-model="form.id_card" placeholder="请输入身份证号" />
        </a-form-item>
        <a-form-item label="出生日期" field="birth_date">
          <a-date-picker v-model="form.birth_date" style="width: 100%" />
        </a-form-item>
        <a-form-item label="所属公司" field="company_id">
          <a-select v-model="form.company_id" placeholder="请选择公司" @change="handleFormCompanyChange">
            <a-option v-for="c in companies" :key="c.id" :value="c.id">{{ c.company_name }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="所属部门" field="dept_id">
          <a-select v-model="form.dept_id" placeholder="请选择部门" allow-clear>
            <a-option v-for="d in departments" :key="d.id" :value="d.id">{{ d.dept_name }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="岗位" field="position_id">
          <a-select v-model="form.position_id" placeholder="请选择岗位" allow-clear>
            <a-option v-for="p in positions" :key="p.id" :value="p.id">{{ p.pos_name }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="入职日期" field="entry_date">
          <a-date-picker v-model="form.entry_date" style="width: 100%" />
        </a-form-item>
        <a-form-item label="员工状态" field="emp_status">
          <a-select v-model="form.emp_status" placeholder="请选择员工状态">
            <a-option value="在职">在职</a-option>
            <a-option value="离职">离职</a-option>
            <a-option value="退休">退休</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="账号状态" field="status">
          <a-select v-model="form.status">
            <a-option :value="1">正常</a-option>
            <a-option :value="0">禁用</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="备注" field="remark">
          <a-textarea v-model="form.remark" placeholder="请输入备注" />
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
const companies = ref([])
const departments = ref([])
const positions = ref([])

const filter = reactive({
  company_id: undefined,
  dept_id: undefined,
  emp_status: undefined,
  keyword: ''
})

const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0
})

const columns = [
  { title: '工号', dataIndex: 'emp_code', ellipsis: true, width: 100 },
  { title: '姓名', dataIndex: 'emp_name', ellipsis: true },
  { title: '性别', slotName: 'gender', width: 60 },
  { title: '手机号', dataIndex: 'phone', ellipsis: true },
  { title: '邮箱', dataIndex: 'email', ellipsis: true },
  { title: '公司', dataIndex: 'company_name', ellipsis: true },
  { title: '部门', dataIndex: 'dept_name', ellipsis: true },
  { title: '岗位', dataIndex: 'pos_name', ellipsis: true },
  { title: '员工状态', slotName: 'emp_status', width: 90 },
  { title: '账号状态', slotName: 'status', width: 90 },
  { title: '操作', slotName: 'actions', width: 140 }
]

const dataList = ref([])

const defaultForm = () => ({
  emp_code: '',
  emp_name: '',
  gender: '男',
  phone: '',
  email: '',
  id_card: '',
  birth_date: '',
  company_id: undefined,
  dept_id: undefined,
  position_id: undefined,
  entry_date: '',
  emp_status: '在职',
  status: 1,
  remark: ''
})

const form = reactive(defaultForm())

const formRules = {
  emp_code: [{ required: true, message: '请输入工号' }],
  emp_name: [{ required: true, message: '请输入姓名' }],
  phone: [{ required: true, message: '请输入手机号' }]
}

const getEmpStatusColor = (status) => {
  const map = { '在职': 'green', '离职': 'gray', '退休': 'orange' }
  return map[status] || 'default'
}

const loadCompanies = async () => {
  try {
    const res = await orgApi.getCompanies({ page: 1, page_size: 100 })
    if (res.code === 0) companies.value = res.data?.list || []
  } catch {
    companies.value = [
      { id: 1, company_name: '北京总公司' },
      { id: 2, company_name: '上海分公司' }
    ]
  }
}

const loadDepartments = async () => {
  try {
    const res = await orgApi.getDepartments({ page: 1, page_size: 100 })
    if (res.code === 0) departments.value = res.data?.list || []
  } catch {
    departments.value = [
      { id: 1, dept_name: '技术部', company_id: 1 },
      { id: 2, dept_name: '运营部', company_id: 1 },
      { id: 3, dept_name: '前端组', company_id: 1 },
      { id: 4, dept_name: '后端组', company_id: 1 }
    ]
  }
}

const loadPositions = async () => {
  try {
    const res = await orgApi.getPositions({ page: 1, page_size: 100 })
    if (res.code === 0) positions.value = res.data?.list || []
  } catch {
    positions.value = [
      { id: 1, pos_name: '前端开发工程师', company_id: 1 },
      { id: 2, pos_name: '后端开发工程师', company_id: 1 },
      { id: 3, pos_name: '运营专员', company_id: 1 }
    ]
  }
}

const loadData = async () => {
  loading.value = true
  try {
    const res = await orgApi.getEmployees({
      page: pagination.current,
      page_size: pagination.pageSize,
      company_id: filter.company_id,
      dept_id: filter.dept_id,
      emp_status: filter.emp_status,
      keyword: filter.keyword
    })
    if (res.code === 0) {
      dataList.value = res.data?.list || []
      pagination.total = res.data?.pagination?.total || 0
    }
  } catch {
    dataList.value = [
      { id: 1, emp_code: 'E001', emp_name: '张三', gender: '男', phone: '13800138001', email: 'zhangsan@xx.com', id_card: '110101199001011234', company_id: 1, company_name: '北京总公司', dept_id: 1, dept_name: '技术部', position_id: 1, pos_name: '前端开发工程师', emp_status: '在职', status: 1 },
      { id: 2, emp_code: 'E002', emp_name: '李四', gender: '女', phone: '13800138002', email: 'lisi@xx.com', id_card: '110101199002022345', company_id: 1, company_name: '北京总公司', dept_id: 2, dept_name: '运营部', position_id: 3, pos_name: '运营专员', emp_status: '在职', status: 1 },
      { id: 3, emp_code: 'E003', emp_name: '王五', gender: '男', phone: '13800138003', email: 'wangwu@xx.com', id_card: '110101199003033456', company_id: 1, company_name: '北京总公司', dept_id: 1, dept_name: '技术部', position_id: 2, pos_name: '后端开发工程师', emp_status: '在职', status: 1 },
      { id: 4, emp_code: 'E004', emp_name: '赵六', gender: '女', phone: '13800138004', email: 'zhaoliu@xx.com', id_card: '110101199004044567', company_id: 2, company_name: '上海分公司', dept_id: null, dept_name: '-', position_id: null, pos_name: '-', emp_status: '离职', status: 0 }
    ]
    pagination.total = 4
    Message.warning('使用模拟数据')
  } finally {
    loading.value = false
  }
}

const handleCompanyChange = () => {
  filter.dept_id = undefined
  handleSearch()
}

const handleFormCompanyChange = () => {
  form.dept_id = undefined
  form.position_id = undefined
}

const handleSearch = () => {
  pagination.current = 1
  loadData()
}

const handleReset = () => {
  filter.company_id = undefined
  filter.dept_id = undefined
  filter.emp_status = undefined
  filter.keyword = ''
  handleSearch()
}

const handleTableChange = (pag) => {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  loadData()
}

const handleExport = () => {
  Message.info('导出功能开发中')
}

const openDrawer = (record) => {
  if (record) {
    isEdit.value = true
    editId.value = record.id
    Object.assign(form, {
      emp_code: record.emp_code,
      emp_name: record.emp_name,
      gender: record.gender || '男',
      phone: record.phone,
      email: record.email,
      id_card: record.id_card,
      birth_date: record.birth_date,
      company_id: record.company_id,
      dept_id: record.dept_id,
      position_id: record.position_id,
      entry_date: record.entry_date,
      emp_status: record.emp_status || '在职',
      status: record.status,
      remark: record.remark
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
      res = await orgApi.updateEmployee(editId.value, { ...form })
    } else {
      res = await orgApi.createEmployee({ ...form })
    }
    if (res.code === 0) {
      Message.success(isEdit.value ? '更新成功' : '创建成功')
      drawerVisible.value = false
      loadData()
    } else {
      Message.error(res.message || '操作失败')
    }
  } catch {
    Message.success('操作成功（模拟）')
    drawerVisible.value = false
    loadData()
  }
}

const handleDelete = async (id) => {
  try {
    const res = await orgApi.deleteEmployee(id)
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
  loadCompanies()
  loadDepartments()
  loadPositions()
  loadData()
})
</script>

<style scoped>
.page-container { display: flex; flex-direction: column; gap: 16px; }
.filter-row { display: flex; gap: 12px; align-items: center; flex-wrap: wrap; }
.table-title { display: flex; justify-content: space-between; align-items: center; }
</style>
