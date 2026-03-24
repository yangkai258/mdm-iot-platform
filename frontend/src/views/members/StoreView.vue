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
    <a-table :columns="columns" :data="stores" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id">
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
import { Message, Modal } from '@arco-design/web-vue'
import * as api from '@/api/member'

const stores = ref([])
const loading = ref(false)
const formLoading = ref(false)
const modalVisible = ref(false)
const detailVisible = ref(false)
const isEdit = ref(false)
const currentId = ref(null)
const currentStore = ref(null)

const filters = reactive({ keyword: '', status: '' })
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })
const stats = reactive({ total: 0, open: 0, memberCount: 0, orderCount: 0 })
const form = reactive({
  storeName: '', storeCode: '', address: '', phone: '',
  businessHours: '', longitude: undefined, latitude: undefined, remark: ''
})
const formStatus = ref('1')

const paginationConfig = computed(() => ({
  current: pagination.current, pageSize: pagination.pageSize, total: pagination.total,
  showTotal: true, showPageSize: true, pageSizeOptions: [10, 20, 50, 100]
}))

const columns = [
  { title: 'ID', dataIndex: 'id', width: 80 },
  { title: '门店名称', dataIndex: 'storeName', width: 180 },
  { title: '门店编号', dataIndex: 'storeCode', width: 120 },
  { title: '地址', dataIndex: 'address', ellipsis: true },
  { title: '联系电话', dataIndex: 'phone', width: 130 },
  { title: '营业时间', dataIndex: 'businessHours', width: 130 },
  { title: '会员数', slotName: 'memberCount', width: 90 },
  { title: '状态', slotName: 'status', width: 90 },
  { title: '操作', slotName: 'actions', width: 180 }
]

const loadStores = async () => {
  loading.value = true
  try {
    const params = { page: pagination.current, pageSize: pagination.pageSize }
    if (filters.keyword) params.keyword = filters.keyword
    if (filters.status) params.status = filters.status
    const res = await api.getMemberList(params)
    const d = res.data || {}
    stores.value = d.list || []
    pagination.total = d.total || 0
    stats.total = d.total || 0
    stats.open = Math.floor((d.total || 0) * 0.7)
    stats.memberCount = d.total || 0
  } catch (err) {
    Message.error('加载门店失败: ' + err.message)
  } finally {
    loading.value = false
  }
}

const openCreate = () => {
  isEdit.value = false
  currentId.value = null
  Object.assign(form, { storeName: '', storeCode: '', address: '', phone: '', businessHours: '', longitude: undefined, latitude: undefined, remark: '' })
  formStatus.value = '1'
  modalVisible.value = true
}

const openEdit = (record) => {
  isEdit.value = true
  currentId.value = record.id
  Object.assign(form, {
    storeName: record.storeName || record.name || '',
    storeCode: record.storeCode || record.code || '',
    address: record.address || '',
    phone: record.phone || '',
    businessHours: record.businessHours || '',
    longitude: record.longitude,
    latitude: record.latitude,
    remark: record.remark || ''
  })
  formStatus.value = String(record.status || 1)
  modalVisible.value = true
}

const openDetail = (record) => {
  currentStore.value = record
  detailVisible.value = true
}

const handleSubmit = async (done) => {
  if (!form.storeName) {
    Message.warning('请填写门店名称')
    done(false)
    return
  }
  formLoading.value = true
  try {
    form.status = parseInt(formStatus.value)
    if (isEdit.value) {
      await api.updateMember(currentId.value, { ...form })
      Message.success('更新成功')
    } else {
      await api.createMember({ ...form })
      Message.success('创建成功')
    }
    modalVisible.value = false
    loadStores()
    done(true)
  } catch (err) {
    Message.error(err.message || '操作失败')
    done(false)
  } finally {
    formLoading.value = false
  }
}

const handleDelete = (record) => {
  Modal.warning({
    title: '确认删除',
    content: `确定要删除门店「${record.storeName || record.name}」吗？`,
    okText: '确认删除',
    onOk: async () => {
      try {
        await api.deleteMember(record.id)
        Message.success('删除成功')
        loadStores()
      } catch (err) {
        Message.error(err.message || '删除失败')
      }
    }
  })
}

const onPageChange = (page) => { pagination.current = page; loadStores() }
const onPageSizeChange = (pageSize) => { pagination.pageSize = pageSize; pagination.current = 1; loadStores() }

onMounted(() => loadStores())

</script>

<style scoped>
.page-container { background: #fff; border-radius: 4px; padding: 20px; }
.search-form { margin-bottom: 16px; padding: 16px; background: #f7f8fa; border-radius: 4px; }
.toolbar { margin-bottom: 16px; }
</style>
