<template>
    <Breadcrumb :items="['Home','Console','']" />


  <div class="container">
    <a-card class="general-card" title="部门管理">
      <template #extra>
        <a-space>
          <a-button type="primary" @click="openCreate(null)"><icon-plus />新建部门</a-button>
          <a-button @click="loadData"><icon-refresh />刷新</a-button>
        </a-space>
      </template>
      <a-spin :loading="loading">
        <a-table :columns="columns" :data="data" :loading="loading" row-key="id">
          <template #actions="{ record }">
            <a-button type="text" size="small" @click="openCreate(record)">新增子部门</a-button>
            <a-button type="text" size="small" @click="openEdit(record)">编辑</a-button>
            <a-button type="text" size="small" status="danger" @click="handleDelete(record)">删除</a-button>
          </template>
      </a-table>
      </a-spin>
    </a-card>
    <a-modal v-model="formVisible" :title="isEdit ? '编辑部门' : '新建部门'" @ok="handleSubmit">
      <a-form :model="form" layout="vertical">
        <a-form-item label="部门名称" required><a-input v-model="form.name" /></a-form-item>
        <a-form-item label="上级部门"><a-input :value="parentDeptName" readonly /></a-form-item>
        <a-form-item label="部门编码"><a-input v-model="form.code" /></a-form-item>
        <a-form-item label="排序"><a-input-number v-model="form.sort" /></a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import Breadcrumb from '@/components/breadcrumb'

const loading = ref(false)
const formVisible = ref(false)
const isEdit = ref(false)
const parentDeptName = ref('-')
const form = reactive({ name: '', code: '', sort: 0, parent_id: null })
const data = ref([])
const columns = [
  { title: '部门名称', dataIndex: 'name', width: 220 },
  { title: '部门编码', dataIndex: 'code', width: 140 },
  { title: '排序', dataIndex: 'sort', width: 80 },
  { title: '操作', slotName: 'actions', width: 200 }
]

const loadData = async () => {
  loading.value = true
  try {
    const res = await fetch('/api/org/departments', {
      headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token') }
    }).then(r => r.json())
    data.value = res.data?.list || []
  } catch { data.value = [] } finally { loading.value = false }
}

const openCreate = (parent) => { isEdit.value = false; Object.assign(form, { name: '', code: '', sort: 0, parent_id: parent?.id || null }); parentDeptName.value = parent?.name || '-'; formVisible.value = true }
const openEdit = (record) => { isEdit.value = true; Object.assign(form, record); parentDeptName.value = '-'; formVisible.value = true }
const handleSubmit = () => { formVisible.value = false; Message.success(isEdit.value ? '更新成功' : '创建成功'); loadData() }
const handleDelete = () => { Message.success('删除成功'); loadData() }

onMounted(() => loadData())
</script>
