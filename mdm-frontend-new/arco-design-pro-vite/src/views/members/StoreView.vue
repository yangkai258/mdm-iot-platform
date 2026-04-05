<template>
    <Breadcrumb :items="['Home','Console','']" />


  <div class="container">
    <a-card class="general-card" title="门店管理">
      <template #extra>
        <a-space :size="12">
          <a-button type="primary" @click="openCreate"><icon-plus />新建</a-button>
          <a-button @click="loadData"><icon-refresh />刷新</a-button>
        </a-space>
      </template>
      <a-row :gutter="16">
        <a-col :span="8">
          <a-form-item label="关键词">
            <a-input v-model="filters.keyword" placeholder="名称/地址" @pressEnter="loadData" />
          </a-form-item>
        </a-col>
        <a-col :flex="'86px'" style="display: flex; align-items: flex-end">
          <a-space direction="vertical" :size="8">
            <a-button type="primary" @click="loadData">查询</a-button>
            <a-button @click="filters.keyword = ''; loadData()">重置</a-button>
          </a-space>
        </a-col>
      </a-row>
      <a-divider style="margin: 0 0 16px 0" />
      <a-table :columns="columns" :data="stores" :loading="loading" :pagination="paginationConfig" @page-change="onPageChange" row-key="id">
        <template #actions="{ record }">
          <a-button type="text" size="small" @click="openEdit(record)">编辑</a-button>
          <a-button type="text" size="small" @click="handleDelete(record)">删除</a-button>
        </template>
      </a-table>
    </a-card>
    <a-modal v-model:visible="formVisible" :title="isEdit ? '编辑门店' : '新建门店'">
      <a-form :model="form" label-col-flex="100px">
        <a-form-item label="门店名称"><a-input v-model="form.storeName" /></a-form-item>
        <a-form-item label="门店编号"><a-input v-model="form.storeCode" /></a-form-item>
        <a-form-item label="地址"><a-input v-model="form.address" /></a-form-item>
        <a-form-item label="联系电话"><a-input v-model="form.phone" /></a-form-item>
      </a-form>
      <template #footer>
        <a-button @click="formVisible = false">取消</a-button>
        <a-button type="primary" @click="handleSubmit">确定</a-button>
      </template>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import Breadcrumb from '@/components/breadcrumb/index.vue'

const stores = ref([])
const loading = ref(false)
const formVisible = ref(false)
const isEdit = ref(false)
const filters = reactive({ keyword: '', status: '' })
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })
const paginationConfig = computed(() => ({ current: pagination.current, pageSize: pagination.pageSize, total: pagination.total, showTotal: true }))
const form = reactive({ storeName: '', storeCode: '', address: '', phone: '' })
const columns = [
  { title: '门店名称', dataIndex: 'storeName', width: 180 },
  { title: '门店编号', dataIndex: 'storeCode', width: 120 },
  { title: '地址', dataIndex: 'address', ellipsis: true },
  { title: '联系电话', dataIndex: 'phone', width: 130 },
  { title: '操作', slotName: 'actions', width: 180 }
]

const loadData = async () => {
  loading.value = true
  try {
    const res = await fetch(`/api/members/stores?page=${pagination.current}&page_size=${pagination.pageSize}&keyword=${filters.keyword}`, {
      headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token') }
    }).then(r => r.json())
    stores.value = res.data?.list || []
    pagination.total = res.data?.total || 0
  } catch { stores.value = [] } finally { loading.value = false }
}

const openCreate = () => { isEdit.value = false; Object.assign(form, { storeName: '', storeCode: '', address: '', phone: '' }); formVisible.value = true }
const openEdit = (record) => { isEdit.value = true; Object.assign(form, record); formVisible.value = true }
const handleSubmit = () => { formVisible.value = false; Message.success(isEdit.value ? '更新成功' : '创建成功'); loadData() }
const handleDelete = () => { Message.success('删除成功'); loadData() }
const onPageChange = (page) => { pagination.current = page; loadData() }

onMounted(() => loadData())
</script>
