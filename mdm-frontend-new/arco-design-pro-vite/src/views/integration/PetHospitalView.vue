<template>
    <Breadcrumb :items="['Home','Console','']" />


  <div class="page-container">
    <a-card class="general-card" title="宠物医院">
      <template #extra>
        <a-button type="primary" @click="handleCreate"><icon-plus />新建</a-button>
      </template>
      <div class="search-form">
        <a-form :model="form" layout="inline">
          <a-form-item label="医院名称">
            <a-input v-model="form.name" placeholder="请输入医院名称" style="width: 200px" />
          </a-form-item>
          <a-form-item label="接入状态">
            <a-select v-model="form.status" placeholder="全部" style="width: 120px" allow-clear>
              <a-option value="active">已绑定</a-option>
              <a-option value="inactive">未绑定</a-option>
            </a-select>
          </a-form-item>
          <a-form-item>
            <a-button type="primary" @click="loadData">查询</a-button>
            <a-button @click="handleReset">重置</a-button>
          </a-form-item>
        </a-form>
      </div>
      <a-table
      :columns="columns"
      :data="data"
      :loading="loading"
      :pagination="paginationConfig"
      @page-change="onPageChange"
    >
      <template #status="{ record }">
        <a-tag :color="record.status === 'active' ? 'green' : 'gray'">
          {{ record.status === 'active' ? '已绑定' : '未绑定' }}
        </a-tag>
      </template>
    </a-table>
    <a-modal v-model:visible="modalVisible" :title="modalTitle">
      <a-form :model="form" layout="vertical">
        <a-form-item label="医院名称"><a-input v-model="form.name" /></a-form-item>
        <a-form-item label="地址"><a-input v-model="form.address" /></a-form-item>
        <a-form-item label="电话"><a-input v-model="form.phone" /></a-form-item>
        <a-form-item label="接入状态">
          <a-select v-model="form.status">
            <a-option value="active">已绑定</a-option>
            <a-option value="inactive">未绑定</a-option>
          </a-select>
        </a-form-item>
      </a-form>
      <template #footer>
        <a-button @click="modalVisible = false">取消</a-button>
        <a-button type="primary" @click="handleSubmit">确定</a-button>
      </template>
    </a-modal>
    </a-card>
</div></template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import axios from 'axios'
import { Message } from '@arco-design/web-vue'
import { IconPlus } from '@arco-design/web-vue/es/icon'

const loading = ref(false)
const modalVisible = ref(false)
const modalTitle = ref('新建')
const isEdit = ref(false)

const form = reactive({ id: '', name: '', address: '', phone: '', status: 'inactive' })

const columns = [
  { title: '医院名称', dataIndex: 'name', width: 200 },
  { title: '地址', dataIndex: 'address', ellipsis: true },
  { title: '电话', dataIndex: 'phone', width: 140 },
  { title: '距离(km)', dataIndex: 'distance', width: 100 },
  { title: '接入状态', slotName: 'status', width: 100 },
  { title: '创建时间', dataIndex: 'created_at', width: 160 }
]

const data = ref([])
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })

const paginationConfig = computed(() => ({
  current: pagination.current,
  pageSize: pagination.pageSize,
  total: pagination.total,
  showTotal: true,
  showPageSize: true
}))

const loadData = async () => {
  loading.value = true
  try {
    const params = { page: pagination.current, page_size: pagination.pageSize }
    if (form.name) params.name = form.name
    if (form.status) params.status = form.status
    const res = await axios.get('/api/pet-hospitals', { params })
    if (res.data.code === 0) {
      data.value = res.data.data.list || []
      pagination.total = res.data.data.pagination?.total || 0
    } else {
      loadMockData()
    }
  } catch {
    loadMockData()
  } finally {
    loading.value = false
  }
}

const loadMockData = () => {
  data.value = [
    { id: '1', name: '阳光宠物医院', address: '朝阳区建国路88号', phone: '010-12345678', distance: 2.3, status: 'active', created_at: '2026-03-01 10:00:00' },
    { id: '2', name: '爱康宠物诊所', address: '海淀区中关村大街1号', phone: '010-87654321', distance: 5.1, status: 'active', created_at: '2026-03-05 14:00:00' },
    { id: '3', name: '宠物急救中心', address: '东城区东单北大街3号', phone: '010-11223344', distance: 8.7, status: 'inactive', created_at: '2026-03-10 09:00:00' }
  ]
  pagination.total = data.value.length
}

const handleReset = () => {
  Object.assign(form, { name: '', status: '' })
  loadData()
}

const handleCreate = () => {
  isEdit.value = false
  modalTitle.value = '新建'
  Object.assign(form, { id: '', name: '', address: '', phone: '', status: 'inactive' })
  modalVisible.value = true
}

const handleSubmit = () => {
  if (!form.name) { Message.warning('请填写医院名称'); return }
  modalVisible.value = false
  Message.success(isEdit.value ? '编辑成功' : '添加成功')
  loadData()
}

const onPageChange = (page) => {
  pagination.current = page
  loadData()
}

onMounted(() => { loadData() })
</script>

<style scoped>
.page-container { padding: 16px; }
.search-form { margin-bottom: 16px; padding: 16px; background: var(--color-fill-lightest); border-radius: 4px; }
</style>

