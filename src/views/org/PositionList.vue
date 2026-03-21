<template>
  <div class="page-container">
    <a-card class="filter-card">
      <div class="filter-row">
        <a-select
          v-model="filter.company_id"
          placeholder="所属公司"
          style="width: 160px"
          allow-clear
          @change="handleSearch"
        >
          <a-option v-for="c in companies" :key="c.id" :value="c.id">{{ c.company_name }}</a-option>
        </a-select>
        <a-select
          v-model="filter.category"
          placeholder="岗位类别"
          style="width: 140px"
          allow-clear
          @change="handleSearch"
        >
          <a-option value="技术">技术</a-option>
          <a-option value="运营">运营</a-option>
          <a-option value="销售">销售</a-option>
          <a-option value="职能">职能</a-option>
        </a-select>
        <a-select
          v-model="filter.status"
          placeholder="状态"
          style="width: 120px"
          allow-clear
          @change="handleSearch"
        >
          <a-option :value="1">正常</a-option>
          <a-option :value="0">禁用</a-option>
        </a-select>
        <a-input-search
          v-model="filter.keyword"
          placeholder="搜索岗位名称/编码"
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
          <span>岗位列表</span>
          <a-button type="primary" @click="openDrawer(null)">
            <template #icon><icon-plus /></template>
            新建岗位
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
          <a-popconfirm content="确定删除该岗位？" @ok="handleDelete(record.id)">
            <a-button type="text" size="small" status="danger">删除</a-button>
          </a-popconfirm>
        </template>
      </a-table>
    </a-card>

    <!-- 创建/编辑 Drawer -->
    <a-drawer
      v-model:visible="drawerVisible"
      :title="isEdit ? '编辑岗位' : '新建岗位'"
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
        <a-form-item label="所属部门" field="dept_id">
          <a-select v-model="form.dept_id" placeholder="请选择部门" allow-clear>
            <a-option v-for="d in departments" :key="d.id" :value="d.id">{{ d.dept_name }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="岗位编码" field="pos_code">
          <a-input v-model="form.pos_code" placeholder="请输入岗位编码" />
        </a-form-item>
        <a-form-item label="岗位名称" field="pos_name">
          <a-input v-model="form.pos_name" placeholder="请输入岗位名称" />
        </a-form-item>
        <a-form-item label="岗位类别" field="category">
          <a-select v-model="form.category" placeholder="请选择岗位类别">
            <a-option value="技术">技术</a-option>
            <a-option value="运营">运营</a-option>
            <a-option value="销售">销售</a-option>
            <a-option value="职能">职能</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="级别" field="level">
          <a-input-number v-model="form.level" :min="1" :max="10" style="width: 100%" />
        </a-form-item>
        <a-form-item label="岗位描述" field="description">
          <a-textarea v-model="form.description" placeholder="请输入岗位描述" />
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
const companies = ref([])
const departments = ref([])

const filter = reactive({
  company_id: undefined,
  category: undefined,
  status: undefined,
  keyword: ''
})

const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0
})

const columns = [
  { title: '岗位编码', dataIndex: 'pos_code', ellipsis: true },
  { title: '岗位名称', dataIndex: 'pos_name', ellipsis: true },
  { title: '所属公司', dataIndex: 'company_name', ellipsis: true },
  { title: '所属部门', dataIndex: 'dept_name', ellipsis: true },
  { title: '类别', dataIndex: 'category', width: 80 },
  { title: '级别', dataIndex: 'level', width: 60 },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '操作', slotName: 'actions', width: 140 }
]

const dataList = ref([])

const defaultForm = () => ({
  company_id: undefined,
  dept_id: undefined,
  pos_code: '',
  pos_name: '',
  category: '',
  level: 1,
  description: '',
  status: 1,
  sort: 0
})

const form = reactive(defaultForm())

const formRules = {
  company_id: [{ required: true, message: '请选择公司' }],
  pos_code: [{ required: true, message: '请输入岗位编码' }],
  pos_name: [{ required: true, message: '请输入岗位名称' }]
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

const loadData = async () => {
  loading.value = true
  try {
    const res = await orgApi.getPositions({
      page: pagination.current,
      page_size: pagination.pageSize,
      company_id: filter.company_id,
      category: filter.category,
      status: filter.status,
      keyword: filter.keyword
    })
    if (res.code === 0) {
      dataList.value = res.data?.list || []
      pagination.total = res.data?.pagination?.total || 0
    }
  } catch {
    dataList.value = [
      { id: 1, pos_code: 'P001', pos_name: '前端开发工程师', company_id: 1, company_name: '北京总公司', dept_id: 3, dept_name: '前端组', category: '技术', level: 2, description: '负责前端开发', status: 1, sort: 1 },
      { id: 2, pos_code: 'P002', pos_name: '后端开发工程师', company_id: 1, company_name: '北京总公司', dept_id: 4, dept_name: '后端组', category: '技术', level: 2, description: '负责后端开发', status: 1, sort: 2 },
      { id: 3, pos_code: 'P003', pos_name: '运营专员', company_id: 1, company_name: '北京总公司', dept_id: 2, dept_name: '运营部', category: '运营', level: 1, description: '负责运营工作', status: 1, sort: 3 },
      { id: 4, pos_code: 'P004', pos_name: '销售经理', company_id: 2, company_name: '上海分公司', dept_id: null, dept_name: '-', category: '销售', level: 3, description: '负责销售管理', status: 0, sort: 4 }
    ]
    pagination.total = 4
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
  filter.company_id = undefined
  filter.category = undefined
  filter.status = undefined
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
      res = await orgApi.updatePosition(editId.value, { ...form })
    } else {
      res = await orgApi.createPosition({ ...form })
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
    const res = await orgApi.deletePosition(id)
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
  loadData()
})
</script>

<style scoped>
.page-container { display: flex; flex-direction: column; gap: 16px; }
.filter-row { display: flex; gap: 12px; align-items: center; flex-wrap: wrap; }
.table-title { display: flex; justify-content: space-between; align-items: center; }
</style>
