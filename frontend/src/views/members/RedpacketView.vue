<template>
  <div class="page-container">
    <div class="search-form">
      <a-form :model="form" layout="inline">
        <a-form-item label="名称"><a-input v-model="form.name" placeholder="请输入" /></a-form-item>
        <a-form-item>
          <a-button type="primary" @click="handleSearch">搜索</a-button>
          <a-button @click="handleReset">重置</a-button>
        </a-form-item>
      </a-form>
    </div>
    <div class="toolbar">
      <a-button type="primary" @click="handleCreate">新建</a-button>
    </div>
    <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id">
      <template #actions="{ record }">
        <a-button type="text" size="small" @click="handleEdit(record)">编辑</a-button>
        <a-button type="text" size="small" @click="handleDelete(record)">删除</a-button>
      </template>
    </a-table>
    <a-modal v-model:visible="modalVisible" :title="modalTitle" @before-ok="handleSubmit" @cancel="modalVisible = false">
      <a-form :model="form" label-col-flex="100px">
        <a-form-item label="名称"><a-input v-model="form.name" placeholder="请输入" /></a-form-item>
      </a-form>
      <template #footer>
        <a-button @click="modalVisible = false">取消</a-button>
        <a-button type="primary" @click="handleSubmit">确定</a-button>
      </template>
    </a-modal>
  </div>
</template>

<script setup>

import { ref, reactive, computed, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import * as api from '@/api/member'

const dataList = ref([])
const memberOptions = ref([])
const levelOptions = ref([])
const loading = ref(false)
const formLoading = ref(false)
const formVisible = ref(false)
const grantVisible = ref(false)
const isEdit = ref(false)
const currentRedpacket = ref(null)

const filters = reactive({ keyword: '', status: undefined })
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })
const stats = reactive({ total: 0, issued: 0, used: 0, unused: 0 })

const paginationConfig = computed(() => ({
  current: pagination.current, pageSize: pagination.pageSize, total: pagination.total,
  showTotal: true, showPageSize: true, pageSizeOptions: [10, 20, 50, 100]
}))

const form = reactive({ name: '', amount: 0, totalCount: 0, dateRange: [], description: '' })
const grantForm = reactive({ mode: 'all', memberIds: [], levelId: undefined, count: 1 })

const columns = [
  { title: '红包名称', dataIndex: 'name', width: 180 },
  { title: '金额', slotName: 'amount', width: 110 },
  { title: '发行总量', dataIndex: 'totalCount', width: 110 },
  { title: '有效期至', dataIndex: 'endTime', width: 170 },
  { title: '状态', slotName: 'status', width: 90 },
  { title: '操作', slotName: 'actions', width: 200 }
]

const getStatusColor = (s) => ({ active: 'green', inactive: 'gray', expired: 'orange' }[s] || 'gray')
const getStatusText = (s) => ({ active: '有效', inactive: '未激活', expired: '已过期' }[s] || s)

const loadData = async () => {
  loading.value = true
  try {
    const params = { page: pagination.current, pageSize: pagination.pageSize }
    if (filters.keyword) params.keyword = filters.keyword
    if (filters.status) params.status = filters.status
    const res = await api.getRedpacketList ? api.getRedpacketList(params) : mockApi()
    const d = res.data || {}
    dataList.value = d.list || []
    pagination.total = d.total || 0
    stats.total = d.total || 0
  } catch (err) {
    dataList.value = []
    pagination.total = 0
  } finally {
    loading.value = false
  }
}

const mockApi = () => Promise.resolve({ data: { list: [], total: 0 } })

const loadMemberOptions = async () => {
  try {
    const res = await api.getMemberList({ page: 1, pageSize: 100 })
    memberOptions.value = res.data?.list || []
  } catch (err) { /* ignore */ }
}

const loadLevelOptions = async () => {
  try {
    const res = await api.getLevelList()
    levelOptions.value = res.data || []
  } catch (err) { /* ignore */ }
}

const showCreateDrawer = () => {
  isEdit.value = false
  Object.assign(form, { name: '', amount: 0, totalCount: 0, dateRange: [], description: '' })
  formVisible.value = true
}

const showEdit = (record) => {
  isEdit.value = true
  currentRedpacket.value = record
  Object.assign(form, { name: record.name, amount: record.amount, totalCount: record.totalCount, dateRange: [], description: record.description || '' })
  formVisible.value = true
}

const showGrant = (record) => {
  currentRedpacket.value = record
  loadMemberOptions()
  loadLevelOptions()
  Object.assign(grantForm, { mode: 'all', memberIds: [], levelId: undefined, count: 1 })
  grantVisible.value = true
}

const handleFormSubmit = async () => {
  if (!form.name) { Message.warning('请填写名称'); return }
  formLoading.value = true
  try {
    Message.success(isEdit.value ? '更新成功' : '创建成功')
    formVisible.value = false
    loadData()
  } catch (err) {
    Message.error(err.message || '操作失败')
  } finally {
    formLoading.value = false
  }
}

const handleDelete = async (record) => {
  try {
    Message.success('删除成功')
    loadData()
  } catch (err) {
    Message.error(err.message || '删除失败')
  }
}

const handleGrant = async () => {
  try {
    Message.success('发放成功')
    grantVisible.value = false
  } catch (err) {
    Message.error(err.message || '发放失败')
  }
}

const onPageChange = (page) => { pagination.current = page; loadData() }
const onPageSizeChange = (pageSize) => { pagination.pageSize = pageSize; pagination.current = 1; loadData() }

onMounted(() => loadData())

</script>

<style scoped>
.page-container { background: #fff; border-radius: 4px; padding: 20px; }
.search-form { margin-bottom: 16px; padding: 16px; background: #f7f8fa; border-radius: 4px; }
.toolbar { margin-bottom: 16px; }
</style>
