<template>
    <Breadcrumb :items="['Home','Console','']" />


  <div class="container">
    <a-card class="general-card" title="岗位管理">
      <template #extra>
        <a-space>
          <a-button type="primary" @click="openCreate"><icon-plus />新建岗位</a-button>
          <a-button @click="loadData"><icon-refresh />刷新</a-button>
        </a-space>
      </template>
      <a-row :gutter="16">
        <a-col :span="8">
          <a-input-search v-model="filter.keyword" placeholder="搜索岗位名称/编码" @search="loadData" />
        </a-col>
        <a-col :flex="'86px'" style="display: flex; align-items: flex-end">
          <a-space direction="vertical" :size="8">
            <a-button type="primary" @click="loadData">查询</a-button>
            <a-button @click="Object.keys(filter).forEach(k => filter[k] = ''); loadData()">重置</a-button>
          </a-space>
        </a-col>
      </a-row>
      <a-divider style="margin: 16px 0 0 0" />
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id" style="margin-top: 16px">
        <template #status="{ record }"><a-badge :color="record.status === 1 ? 'green' : 'red'" :text="record.status === 1 ? '正常' : '禁用'" /></template>
        <template #actions="{ record }">
          <a-button type="text" size="small" @click="openEdit(record)">编辑</a-button>
          <a-button type="text" size="small" status="danger" @click="handleDelete(record)">删除</a-button>
        </template>
      </a-table>
      </a-table>
    </a-card>
    <a-modal v-model="formVisible" :title="isEdit ? '编辑岗位' : '新建岗位'" @ok="handleSubmit">
      <a-form :model="form" layout="vertical">
        <a-form-item label="岗位名称" required><a-input v-model="form.name" /></a-form-item>
        <a-form-item label="岗位编码"><a-input v-model="form.code" /></a-form-item>
        <a-form-item label="所属部门"><a-select v-model="form.department_id" :options="deptOptions" style="width: 100%" /></a-form-item>
        <a-form-item label="状态"><a-switch v-model="form.status" /></a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import Breadcrumb from '@/components/Breadcrumb.vue'

const loading = ref(false)
const formVisible = ref(false)
const isEdit = ref(false)
const filter = reactive({ keyword: '' })
const form = reactive({ name: '', code: '', department_id: null, status: true })
const data = ref([])
const deptOptions = ref([])
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })
const columns = [
  { title: '岗位名称', dataIndex: 'name', width: 200 },
  { title: '岗位编码', dataIndex: 'code', width: 140 },
  { title: '所属部门', dataIndex: 'department_name', width: 180 },
  { title: '状态', slotName: 'status', width: 90 },
  { title: '操作', slotName: 'actions', width: 120 }
]

const loadData = async () => {
  loading.value = true
  try {
    const params = { page: pagination.current, page_size: pagination.pageSize }
    if (filter.keyword) params.keyword = filter.keyword
    const res = await fetch(`/api/org/posts?${new URLSearchParams(params)}`, {
      headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token') }
    }).then(r => r.json())
    data.value = res.data?.list || []
    pagination.total = res.data?.total || 0
  } catch { data.value = [] } finally { loading.value = false }
}

const loadDepts = async () => {
  const res = await fetch('/api/org/departments', {
    headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token') }
  }).then(r => r.json())
  deptOptions.value = (res.data?.list || []).map(d => ({ label: d.name, value: d.id }))
}

const openCreate = () => { isEdit.value = false; Object.assign(form, { name: '', code: '', department_id: null, status: true }); formVisible.value = true }
const openEdit = (record) => { isEdit.value = true; Object.assign(form, record); formVisible.value = true }
const handleSubmit = () => { formVisible.value = false; Message.success(isEdit.value ? '更新成功' : '创建成功'); loadData() }
const handleDelete = () => { Message.success('删除成功'); loadData() }
const onPageChange = (page) => { pagination.current = page; loadData() }

onMounted(() => { loadData(); loadDepts() })
</script>
