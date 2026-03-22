<template>
  <div class="page-container">
    <a-card class="filter-card">
      <div class="filter-row">
        <a-input-search
          v-model="filter.keyword"
          placeholder="搜索公司名称/编码"
          style="width: 280px"
          @search="handleSearch"
          @press-enter="handleSearch"
        />
        <a-select
          v-model="filter.status"
          placeholder="公司状态"
          style="width: 140px"
          allow-clear
          @change="loadData"
        >
          <a-option :value="1">正常</a-option>
          <a-option :value="0">禁用</a-option>
        </a-select>
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
          <span>公司列表</span>
          <a-button type="primary" @click="openDrawer(null)">
            <template #icon><icon-plus /></template>
            新建公司
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
          <a-popconfirm content="确定删除该公司？" @ok="handleDelete(record.id)">
            <a-button type="text" size="small" status="danger">删除</a-button>
          </a-popconfirm>
        </template>
      </a-table>
    </a-card>

    <!-- 创建/编辑 Drawer -->
    <a-drawer
      v-model:visible="drawerVisible"
      :title="isEdit ? '编辑公司' : '新建公司'"
      width="520px"
      @confirm="handleSubmit"
      @cancel="drawerVisible = false"
    >
      <a-form ref="formRef" :model="form" :rules="formRules" layout="vertical" label-align="left">
        <a-form-item label="公司编码" field="company_code">
          <a-input v-model="form.company_code" placeholder="请输入公司编码" />
        </a-form-item>
        <a-form-item label="公司名称" field="company_name">
          <a-input v-model="form.company_name" placeholder="请输入公司名称" />
        </a-form-item>
        <a-form-item label="简称" field="short_name">
          <a-input v-model="form.short_name" placeholder="请输入简称" />
        </a-form-item>
        <a-form-item label="省" field="province">
          <a-input v-model="form.province" placeholder="请输入省" />
        </a-form-item>
        <a-form-item label="市" field="city">
          <a-input v-model="form.city" placeholder="请输入市" />
        </a-form-item>
        <a-form-item label="区" field="district">
          <a-input v-model="form.district" placeholder="请输入区" />
        </a-form-item>
        <a-form-item label="详细地址" field="address">
          <a-input v-model="form.address" placeholder="请输入详细地址" />
        </a-form-item>
        <a-form-item label="法人代表" field="legal_person">
          <a-input v-model="form.legal_person" placeholder="请输入法人代表" />
        </a-form-item>
        <a-form-item label="联系人" field="contact">
          <a-input v-model="form.contact" placeholder="请输入联系人" />
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
  { title: '公司编码', dataIndex: 'company_code', ellipsis: true },
  { title: '公司名称', dataIndex: 'company_name', ellipsis: true },
  { title: '简称', dataIndex: 'short_name', ellipsis: true },
  { title: '省市区', ellipsis: true, render: ({ record }) => `${record.province || ''}${record.city || ''}${record.district || ''}` },
  { title: '联系人', dataIndex: 'contact' },
  { title: '联系电话', dataIndex: 'phone' },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '操作', slotName: 'actions', width: 140 }
]

const dataList = ref([])

const defaultForm = () => ({
  company_code: '',
  company_name: '',
  short_name: '',
  province: '',
  city: '',
  district: '',
  address: '',
  legal_person: '',
  contact: '',
  phone: '',
  email: '',
  status: 1,
  sort: 0,
  remark: ''
})

const form = reactive(defaultForm())

const formRules = {
  company_code: [{ required: true, message: '请输入公司编码' }],
  company_name: [{ required: true, message: '请输入公司名称' }]
}

const loadData = async () => {
  loading.value = true
  try {
    const res = await orgApi.getCompanies({
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
    dataList.value = [
      { id: 1, company_code: 'C001', company_name: '北京总公司', short_name: '北京', province: '北京市', city: '北京市', district: '朝阳区', address: 'xxx路1号', legal_person: '张三', contact: '李四', phone: '13800138000', email: 'bjj@xx.com', status: 1, sort: 1 },
      { id: 2, company_code: 'C002', company_name: '上海分公司', short_name: '上海', province: '上海市', city: '上海市', district: '浦东新区', address: 'yyy路2号', legal_person: '王五', contact: '赵六', phone: '13900139000', email: 'sh@xx.com', status: 1, sort: 2 },
      { id: 3, company_code: 'C003', company_name: '深圳分公司', short_name: '深圳', province: '广东省', city: '深圳市', district: '南山区', address: 'zzz路3号', legal_person: '孙七', contact: '周八', phone: '13700137000', email: 'sz@xx.com', status: 0, sort: 3 }
    ]
    pagination.total = 3
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
      res = await orgApi.updateCompany(editId.value, { ...form })
    } else {
      res = await orgApi.createCompany({ ...form })
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
    const res = await orgApi.deleteCompany(id)
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
.filter-card { }
.filter-row { display: flex; gap: 12px; align-items: center; flex-wrap: wrap; }
.table-card { }
.table-title { display: flex; justify-content: space-between; align-items: center; }
</style>
