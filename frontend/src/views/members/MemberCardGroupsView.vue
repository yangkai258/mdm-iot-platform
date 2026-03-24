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
      <template #cardTypeCount="{ record }">
        <span>{{ record.cardTypeCount || 0 }}</span>
      </template>
      <template #actions="{ record }">
        <a-button type="text" size="small" @click="handleEdit(record)">编辑</a-button>
        <a-button type="text" size="small" @click="handleDelete(record)">删除</a-button>
      </template>
    </a-table>
    <a-modal v-model:visible="modalVisible" :title="modalTitle" @before-ok="handleSubmit" @cancel="modalVisible = false">
      <a-form :model="form" label-col-flex="100px">
        <a-form-item label="分组名称"><a-input v-model="form.name" placeholder="请输入分组名称" /></a-form-item>
        <a-form-item label="描述"><a-textarea v-model="form.description" :rows="3" placeholder="描述信息" /></a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import * as api from '@/api/member'

const loading = ref(false)
const data = ref([])
const modalVisible = ref(false)
const modalTitle = ref('新建分组')
const isEdit = ref(false)
const currentId = ref(null)

const form = reactive({ name: '', description: '' })

const pagination = reactive({ current: 1, pageSize: 20, total: 0 })

const columns = [
  { title: 'ID', dataIndex: 'id', width: 70 },
  { title: '分组名称', dataIndex: 'name', width: 200 },
  { title: '描述', dataIndex: 'description', ellipsis: true },
  { title: '包含卡类型数量', slotName: 'cardTypeCount', width: 150 },
  { title: '操作', slotName: 'actions', width: 120, fixed: 'right' }
]

const loadData = async () => {
  loading.value = true
  try {
    const params = { page: pagination.current, pageSize: pagination.pageSize }
    if (form.name) params.keyword = form.name
    const res = await api.getCardGroupList(params)
    const d = res.data || {}
    data.value = d.list || []
    pagination.total = d.total || 0
  } catch (err) {
    Message.error('加载失败: ' + err.message)
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.current = 1
  loadData()
}

const handleReset = () => {
  Object.assign(form, { name: '', description: '' })
  pagination.current = 1
  loadData()
}

const handleCreate = () => {
  isEdit.value = false
  currentId.value = null
  modalTitle.value = '新建分组'
  Object.assign(form, { name: '', description: '' })
  modalVisible.value = true
}

const handleEdit = (record) => {
  isEdit.value = true
  currentId.value = record.id
  modalTitle.value = '编辑分组'
  Object.assign(form, { name: record.name, description: record.description || '' })
  modalVisible.value = true
}

const handleDelete = async (record) => {
  try {
    await api.deleteCardGroup(record.id)
    Message.success('删除成功')
    loadData()
  } catch (err) {
    Message.error(err.message || '删除失败')
  }
}

const handleSubmit = async (done) => {
  try {
    if (isEdit.value) {
      await api.updateCardGroup(currentId.value, { ...form })
      Message.success('更新成功')
    } else {
      await api.createCardGroup({ ...form })
      Message.success('创建成功')
    }
    modalVisible.value = false
    loadData()
    done(true)
  } catch (err) {
    Message.error(err.message || '操作失败')
    done(false)
  }
}

const onPageChange = (page) => {
  pagination.current = page
  loadData()
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.page-container { background: #fff; border-radius: 4px; padding: 20px; }
.search-form { margin-bottom: 16px; padding: 16px; background: #f7f8fa; border-radius: 4px; }
.toolbar { margin-bottom: 16px; }
</style>
