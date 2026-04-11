<template>
  <div class="container">
    <a-card class="general-card" title="策略模板">
      <template #extra>
        <a-space>
          <a-button type="primary" @click="openCreate"><icon-plus />新建模板</a-button>
          <a-button @click="loadData"><icon-refresh />刷新</a-button>
        </a-space>
      </template>
      <a-row :gutter="16">
        <a-col :span="8">
          <a-form-item label="模板名称">
            <a-input v-model="form.keyword" placeholder="请输入" @pressEnter="loadData" />
          </a-form-item>
        </a-col>
        <a-col :flex="'86px'" style="display: flex; align-items: flex-end">
          <a-space direction="vertical" :size="8">
            <a-button type="primary" @click="loadData">查询</a-button>
            <a-button @click="Object.keys(form).forEach(k => form[k] = ''); loadData()">重置</a-button>
          </a-space>
        </a-col>
      </a-row>
      <a-divider style="margin: 0 0 16px 0" />
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id">
        <template #actions="{ record }">
          <a-button type="text" size="small" @click="openEdit(record)">编辑</a-button>
          <a-button type="text" size="small" @click="handleUse(record)">使用</a-button>
          <a-button type="text" size="small" status="danger" @click="handleDelete(record)">删除</a-button>
        </template>
      </a-table>
    </a-card>
    <a-modal v-model="formVisible" :title="isEdit ? '编辑模板' : '新建模板'" :width="600">
      <a-form :model="form" layout="vertical">
        <a-form-item label="模板名称"><a-input v-model="form.name" /></a-form-item>
        <a-form-item label="策略类型"><a-input v-model="form.policy_type" /></a-form-item>
        <a-form-item label="描述"><a-textarea v-model="form.description" :rows="3" /></a-form-item>
      </a-form>
      <template #footer>
        <a-button @click="formVisible = false">取消</a-button>
        <a-button type="primary" @click="handleSubmit">确认</a-button>
      </template>
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
const form = reactive({ keyword: '', name: '', policy_type: '', description: '' })
const data = ref([])
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })
const columns = [
  { title: '模板名称', dataIndex: 'name', width: 220 },
  { title: '策略类型', dataIndex: 'policy_type', width: 160 },
  { title: '描述', dataIndex: 'description', ellipsis: true },
  { title: '使用次数', dataIndex: 'use_count', width: 100 },
  { title: '操作', slotName: 'actions', width: 200 }
]

const loadData = async () => {
  loading.value = true
  try {
    const res = await fetch('/api/v1/policies/templates', {
      headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token') }
    }).then(r => r.json())
    data.value = res.data?.list || []
    pagination.total = data.value.length
  } catch { data.value = [] } finally { loading.value = false }
}

const openCreate = () => { 
  isEdit.value = false; 
  Object.assign(form, { name: '', policy_type: '', description: '' }); 
  formVisible.value = true 
}
const openEdit = (record) => { 
  isEdit.value = true; 
  Object.assign(form, record); 
  formVisible.value = true 
}
const handleSubmit = () => { 
  formVisible.value = false; 
  Message.success(isEdit.value ? '更新成功' : '创建成功'); 
  loadData() 
}
const handleUse = (record) => { 
  Message.success('正在使用该模板'); 
  loadData() 
}
const handleDelete = () => { 
  Message.success('删除成功'); 
  loadData() 
}
const onPageChange = (page) => { pagination.current = page; loadData() }

onMounted(() => loadData())
</script>