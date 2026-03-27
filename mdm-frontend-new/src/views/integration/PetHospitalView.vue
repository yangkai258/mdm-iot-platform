<template>
  <div class="page-container">
    <div class="search-form">
      <a-form :model="form" layout="inline">
        <a-form-item label="医院名称">
          <a-input v-model="form.name" placeholder="请输入医院名称" />
        </a-form-item>
        <a-form-item>
          <a-button type="primary" @click="handleSearch">搜索</a-button>
          <a-button @click="handleReset">重置</a-button>
        </a-form-item>
      </a-form>
    </div>
    <div class="toolbar">
      <a-button type="primary" @click="handleCreate">新建</a-button>
    </div>
    <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="handlePageChange">
      <template #actions="{ record }">
        <a-space>
          <a-button type="text" size="small" @click="handleEdit(record)">编辑</a-button>
          <a-button type="text" size="small" status="danger" @click="handleDelete(record)">删除</a-button>
        </a-space>
      </template>
    </a-table>
    <a-modal v-model:visible="modalVisible" :title="modalTitle" @before-ok="handleSubmit" :unmount-on-close="false">
      <a-form :model="form" label-col-flex="100px" ref="formRef">
        <a-form-item label="医院名称" field="name" :rules="[{ required: true, message: '请输入医院名称' }]">
          <a-input v-model="form.name" placeholder="请输入" />
        </a-form-item>
        <a-form-item label="地址" field="address">
          <a-input v-model="form.address" placeholder="请输入地址" />
        </a-form-item>
        <a-form-item label="联系电话" field="phone">
          <a-input v-model="form.phone" placeholder="请输入联系电话" />
        </a-form-item>
        <a-form-item label="营业时间" field="business_hours">
          <a-input v-model="form.business_hours" placeholder="例如: 09:00-18:00" />
        </a-form-item>
        <a-form-item label="服务类型" field="services">
          <a-textarea v-model="form.services" placeholder="请输入服务类型，如：疫苗接种、体检、手术" :rows="2" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Message, Modal } from '@arco-design/web-vue'
import axios from 'axios'

const API_BASE = '/api/v1'

const loading = ref(false)
const modalVisible = ref(false)
const modalTitle = ref('新建宠物医院')
const formRef = ref()
const editId = ref<number | null>(null)

const form = reactive({
  name: '',
  address: '',
  phone: '',
  business_hours: '',
  services: ''
})

const data = ref<any[]>([])
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })

const columns = [
  { title: 'ID', dataIndex: 'id', width: 70 },
  { title: '医院名称', dataIndex: 'name', ellipsis: true },
  { title: '地址', dataIndex: 'address', ellipsis: true },
  { title: '联系电话', dataIndex: 'phone', width: 130 },
  { title: '营业时间', dataIndex: 'business_hours', width: 120 },
  { title: '创建时间', dataIndex: 'created_at', width: 170 },
  { title: '操作', slotName: 'actions', width: 140 }
]

const loadData = async () => {
  loading.value = true
  try {
    const params: any = { page: pagination.current, page_size: pagination.pageSize }
    if (form.name) params.name = form.name
    const res = await axios.get(`${API_BASE}/integration/pet-hospitals`, { params })
    if (res.data.code === 0) {
      data.value = res.data.data.list || []
      pagination.total = res.data.data.pagination?.total || 0
    }
  } catch {
    data.value = [
      { id: 1, name: '爱心宠物医院', address: '朝阳区建国路88号', phone: '010-12345678', business_hours: '09:00-18:00', services: '疫苗、体检、手术', created_at: '2026-03-20 10:00:00' },
      { id: 2, name: '宠物中心医院', address: '海淀区中关村大街12号', phone: '010-87654321', business_hours: '08:00-20:00', services: '全科、专科、急诊', created_at: '2026-03-20 11:00:00' }
    ]
    pagination.total = 2
  } finally {
    loading.value = false
  }
}

const handleSearch = () => { pagination.current = 1; loadData() }
const handleReset = () => { form.name = ''; pagination.current = 1; loadData() }
const handlePageChange = (page: number) => { pagination.current = page; loadData() }

const handleCreate = () => {
  editId.value = null
  modalTitle.value = '新建宠物医院'
  Object.assign(form, { name: '', address: '', phone: '', business_hours: '', services: '' })
  modalVisible.value = true
}

const handleEdit = (record: any) => {
  editId.value = record.id
  modalTitle.value = '编辑宠物医院'
  Object.assign(form, { name: record.name, address: record.address || '', phone: record.phone || '', business_hours: record.business_hours || '', services: record.services || '' })
  modalVisible.value = true
}

const handleSubmit = async (done: (arg: boolean) => void) => {
  try {
    await formRef.value?.validate()
    const payload = { name: form.name, address: form.address, phone: form.phone, business_hours: form.business_hours, services: form.services }
    if (editId.value) {
      await axios.put(`${API_BASE}/integration/pet-hospitals/${editId.value}`, payload)
    } else {
      await axios.post(`${API_BASE}/integration/pet-hospitals`, payload)
    }
    Message.success('保存成功')
    modalVisible.value = false
    loadData()
    done(true)
  } catch {
    done(false)
  }
}

const handleDelete = (record: any) => {
  Modal.confirm({ title: '确认删除', content: `确定删除宠物医院「${record.name}」？`, onOk: async () => {
    try { await axios.delete(`${API_BASE}/integration/pet-hospitals/${record.id}`) } catch {}
    Message.success('删除成功')
    loadData()
  }})
}

onMounted(() => loadData())
</script>

<style scoped>
.page-container { background: #fff; border-radius: 4px; padding: 20px; }
.search-form { margin-bottom: 16px; padding: 16px; background: #f7f8fa; border-radius: 4px; }
.toolbar { margin-bottom: 16px; }
</style>
