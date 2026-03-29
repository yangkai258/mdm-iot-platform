<template>
  <div class="container">
    <a-card class="general-card" title="公司管理">
      <template #extra>
        <a-space>
          <a-input-search v-model="filter.keyword" placeholder="搜索公司名称/编码" @search="handleSearch" @press-enter="handleSearch" style="width: 240px" />
          <a-select v-model="filter.status" placeholder="公司状态" allow-clear style="width: 120px" @change="loadData">
            <a-option :value="1">正常</a-option>
            <a-option :value="0">禁用</a-option>
          </a-select>
          <a-button type="primary" @click="handleSearch">查询</a-button>
          <a-button @click="handleReset">重置</a-button>
          <a-button type="primary" @click="openCreateModal"><icon-plus />新建公司</a-button>
        </a-space>
      </template>
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id">
        <template #status="{ record }">
          <a-badge :color="record.status === 1 ? 'green' : 'red'" :text="record.status === 1 ? '正常' : '禁用'" />
        </template>
      </a-table>
        <template #actions="{ record }">
          <a-button type="text" size="small" @click="openEditModal(record)">编辑</a-button>
          <a-button type="text" size="small" status="danger" @click="handleDelete(record)">删除</a-button>
        </template>
      </a-table>
    </a-card>
    <a-drawer v-model="drawerVisible" :title="drawerTitle" :width="560" @before-ok="handleSubmit">
      <a-form :model="form" :rules="formRules" layout="vertical">
        <a-form-item label="公司名称" field="name"><a-input v-model="form.name" placeholder="请输入公司名称" /></a-form-item>
        <a-form-item label="统一社会信用代码" field="license_no"><a-input v-model="form.license_no" placeholder="请输入统一社会信用代码" /></a-form-item>
        <a-form-item label="联系人" field="contact_name"><a-input v-model="form.contact_name" placeholder="请输入联系人姓名" /></a-form-item>
        <a-form-item label="联系电话" field="contact_phone"><a-input v-model="form.contact_phone" placeholder="请输入联系电话" /></a-form-item>
        <a-form-item label="公司地址" field="address"><a-textarea v-model="form.address" placeholder="请输入公司详细地址" /></a-form-item>
        <a-form-item label="状态" field="status"><a-switch v-model="form.status" /></a-form-item>
      </a-form>
    </a-drawer>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import Breadcrumb from '@/components/Breadcrumb.vue'

const loading = ref(false)
const data = ref([])
const drawerVisible = ref(false)
const drawerTitle = ref('新建公司')
const filter = reactive({ keyword: '', status: null })
const form = reactive({ name: '', license_no: '', contact_name: '', contact_phone: '', address: '', status: 1 })
const formRules = { name: { required: true, message: '请输入公司名称' } }
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })
const columns = [
  { title: 'ID', dataIndex: 'id', width: 80 },
  { title: '公司名称', dataIndex: 'name', width: 220 },
  { title: '统一社会信用代码', dataIndex: 'license_no', width: 180 },
  { title: '联系人', dataIndex: 'contact_name', width: 120 },
  { title: '联系电话', dataIndex: 'contact_phone', width: 130 },
  { title: '状态', slotName: 'status', width: 90 },
  { title: '操作', slotName: 'actions', width: 120 }
]

const loadData = async () => {
  loading.value = true
  try {
    const params = { page: pagination.current, page_size: pagination.pageSize }
    if (filter.keyword) params.keyword = filter.keyword
    if (filter.status !== null) params.status = filter.status
    const res = await fetch(`/api/v1/org/companies?${new URLSearchParams(params)}`, {
      headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token') }
    }).then(r => r.json())
    data.value = res.data?.list || []
    pagination.total = res.data?.total || 0
  } catch { data.value = [] } finally { loading.value = false }
}

const handleSearch = () => { pagination.current = 1; loadData() }
const handleReset = () => { filter.keyword = ''; filter.status = null; pagination.current = 1; loadData() }
const openCreateModal = () => { Object.assign(form, { name: '', license_no: '', contact_name: '', contact_phone: '', address: '', status: 1 }); drawerTitle.value = '新建公司'; drawerVisible.value = true }
const openEditModal = (record) => { Object.assign(form, record); drawerTitle.value = '编辑公司'; drawerVisible.value = true }
const handleSubmit = async () => { drawerVisible.value = false; Message.success(drawerTitle.value === '新建公司' ? '创建成功' : '更新成功'); loadData() }
const handleDelete = () => { Message.success('删除成功'); loadData() }
const onPageChange = (page) => { pagination.current = page; loadData() }

onMounted(() => loadData())
</script>
