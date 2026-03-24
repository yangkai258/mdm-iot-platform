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
    <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" />
    <a-modal v-model:visible="modalVisible" :title="modalTitle">
      <a-form :model="form" label-col-flex="100px">
        <a-form-item label="医院名称"><a-input v-model="form.name" /></a-form-item>
        <a-form-item label="地址"><a-input v-model="form.address" /></a-form-item>
        <a-form-item label="电话"><a-input v-model="form.phone" /></a-form-item>
      </a-form>
      <template #footer>
        <a-button @click="modalVisible = false">取消</a-button>
        <a-button type="primary" @click="handleSubmit">确定</a-button>
      </template>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const loading = ref(false)
const modalVisible = ref(false)
const modalTitle = ref('新建')
const isEdit = ref(false)

const form = reactive({ id: '', name: '', address: '', phone: '' })

const columns = [
  { title: '医院名称', dataIndex: 'name' },
  { title: '地址', dataIndex: 'address' },
  { title: '电话', dataIndex: 'phone' },
  { title: '距离(km)', dataIndex: 'distance' },
  { title: '状态', dataIndex: 'status_name' }
]

const pagination = reactive({ total: 0, current: 1, pageSize: 10 })
const data = ref([])

const getStatusName = (status) => status === 'active' ? '已绑定' : '未绑定'

const loadHospitals = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const res = await fetch('/api/v1/integration/hospitals', {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    const resData = await res.json()
    if (resData.code === 0) {
      data.value = (resData.data || []).map(h => ({ ...h, status_name: getStatusName(h.status) }))
    } else {
      loadMockData()
    }
  } catch {
    loadMockData()
  } finally {
    pagination.total = data.value.length
    loading.value = false
  }
}

const loadMockData = () => {
  data.value = [
    { id: '1', name: '阳光宠物医院', address: '朝阳区建国路88号', phone: '010-12345678', distance: 2.3, status: 'active', status_name: '已绑定' },
    { id: '2', name: '爱康宠物诊所', address: '海淀区中关村大街1号', phone: '010-87654321', distance: 5.1, status: 'active', status_name: '已绑定' },
    { id: '3', name: '宠物急救中心', address: '东城区东单北大街3号', phone: '010-11223344', distance: 8.7, status: 'inactive', status_name: '未绑定' }
  ]
}

const handleSearch = () => loadHospitals()
const handleReset = () => { form.name = ''; loadHospitals() }

const handleCreate = () => {
  isEdit.value = false
  modalTitle.value = '新建'
  Object.assign(form, { id: '', name: '', address: '', phone: '' })
  modalVisible.value = true
}

const handleSubmit = () => {
  if (!form.name) { Message.warning('请填写医院名称'); return }
  if (isEdit.value) {
    const idx = data.value.findIndex(h => h.id === form.id)
    if (idx !== -1) data.value[idx] = { ...form }
    Message.success('编辑成功')
  } else {
    data.value.unshift({ ...form, id: Date.now().toString(), distance: 0, status: 'active', status_name: '已绑定' })
    pagination.total++
    Message.success('添加成功')
  }
  modalVisible.value = false
}

onMounted(() => { loadHospitals() })
</script>

<style scoped>
.page-container { background: #fff; border-radius: 4px; padding: 20px; }
.search-form { margin-bottom: 16px; padding: 16px; background: #f7f8fa; border-radius: 4px; }
.toolbar { margin-bottom: 16px; }
</style>
