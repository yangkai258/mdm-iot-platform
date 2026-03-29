<template>
  <div class="page-container">
    <Breadcrumb :items="['menu.integration', 'menu.integration.smartHome']" />
    <a-card class="general-card" title="智能家居">
      <template #extra>
        <a-button type="primary" @click="handleCreate"><icon-plus />新建</a-button>
      </template>
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
        <a-form-item label="设备名称"><a-input v-model="form.name" /></a-form-item>
        <a-form-item label="设备类型"><a-input v-model="form.type" /></a-form-item>
        <a-form-item label="房间"><a-input v-model="form.room" /></a-form-item>
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

const form = reactive({ id: '', name: '', type: '', room: '' })

const columns = [
  { title: '设备名称', dataIndex: 'name' },
  { title: '设备类型', dataIndex: 'type' },
  { title: '房间', dataIndex: 'room' },
  { title: '状态', dataIndex: 'status_name' }
]

const pagination = reactive({ total: 0, current: 1, pageSize: 10 })
const data = ref([])

const getTypeName = (type) => {
  const names = { light: '灯光', ac: '空调', humidifier: '加湿器', speaker: '音响' }
  return names[type] || type
}

const loadDevices = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const res = await fetch('/api/v1/integration/devices', {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    const resData = await res.json()
    if (resData.code === 0) {
      data.value = (resData.data || []).map(d => ({ ...d, type: getTypeName(d.type), status_name: d.status ? '开启' : '关闭' }))
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
    { id: '1', name: '客厅大灯', type: '灯光', room: '客厅', status: true, status_name: '开启' },
    { id: '2', name: '卧室空调', type: '空调', room: '卧室', status: false, status_name: '关闭' },
    { id: '3', name: '客厅加湿器', type: '加湿器', room: '客厅', status: true, status_name: '开启' },
    { id: '4', name: '背景音乐', type: '音响', room: '全屋', status: false, status_name: '关闭' }
  ]
}

const handleSearch = () => loadDevices()
const handleReset = () => { form.name = ''; loadDevices() }

const handleCreate = () => {
  isEdit.value = false
  modalTitle.value = '新建'
  Object.assign(form, { id: '', name: '', type: '', room: '' })
  modalVisible.value = true
}

const handleSubmit = () => {
  if (!form.name) { Message.warning('请填写设备名称'); return }
  if (isEdit.value) {
    const idx = data.value.findIndex(d => d.id === form.id)
    if (idx !== -1) data.value[idx] = { ...form, type: getTypeName(form.type) }
    Message.success('编辑成功')
  } else {
    data.value.unshift({ ...form, id: Date.now().toString(), status: false, status_name: '关闭' })
    pagination.total++
    Message.success('添加成功')
  }
  modalVisible.value = false
}

onMounted(() => loadDevices())
</script>

<style scoped>
.page-container { background: #fff; border-radius: 4px; padding: 20px; }
.search-form { margin-bottom: 16px; padding: 16px; background: #f7f8fa; border-radius: 4px; }
.toolbar { margin-bottom: 16px; }
</style>
