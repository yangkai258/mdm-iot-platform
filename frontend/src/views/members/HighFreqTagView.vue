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
    <a-table :columns="columns" :data="dataList" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id">
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

import { ref, reactive } from 'vue'
import { Message } from '@arco-design/web-vue'

const loading = ref(false)
const formVisible = ref(false)
const isEdit = ref(false)
const dataList = ref([])
const stats = ref({ total: 0, coverMembers: 0, active: 0 })
const filters = reactive({ keyword: '' })
const paginationConfig = reactive({ current: 1, pageSize: 10, total: 0 })
const form = reactive({ id: null, name: '', frequency: 3, description: '', status: '1' })

const columns = [
  { title: '标签名称', dataIndex: 'name', width: 200 },
  { title: '定义条件（消费频次/月）', dataIndex: 'frequency', width: 200 },
  { title: '包含会员数', dataIndex: 'memberCount', width: 150 },
  { title: '描述', dataIndex: 'description', ellipsis: true },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '操作', slotName: 'actions', width: 150, fixed: 'right' }
]

const mockData = () => {
  return [
    { id: 1, name: '月度活跃买家', frequency: 4, memberCount: 1234, description: '每月至少消费4次', status: 1 },
    { id: 2, name: '高频剁手党', frequency: 8, memberCount: 456, description: '每月消费8次以上', status: 1 },
    { id: 3, name: '周活跃会员', frequency: 12, memberCount: 89, description: '每周至少消费3次', status: 1 }
  ]
}

const loadData = () => {
  loading.value = true
  setTimeout(() => {
    dataList.value = mockData()
    stats.value = { total: mockData().length, coverMembers: mockData().reduce((s, d) => s + d.memberCount, 0), active: mockData().filter(d => d.status === 1).length }
    paginationConfig.total = dataList.value.length
    loading.value = false
  }, 400)
}

const onPageChange = (page) => {
  paginationConfig.current = page
  loadData()
}

const showCreate = () => {
  isEdit.value = false
  Object.assign(form, { id: null, name: '', frequency: 3, description: '', status: '1' })
  formVisible.value = true
}

const showEdit = (record) => {
  isEdit.value = true
  Object.assign(form, { ...record, status: String(record.status) })
  formVisible.value = true
}

const handleSubmit = (done) => {
  if (!form.name || !form.frequency) {
    Message.error('请填写必填项')
    done(false)
    return
  }
  setTimeout(() => {
    Message.success(isEdit.value ? '更新成功' : '创建成功')
    formVisible.value = false
    loadData()
    done(true)
  }, 400)
}

const handleDelete = (record) => {
  Message.success('删除成功')
  loadData()
}

loadData()

</script>

<style scoped>
.page-container { background: #fff; border-radius: 4px; padding: 20px; }
.search-form { margin-bottom: 16px; padding: 16px; background: #f7f8fa; border-radius: 4px; }
.toolbar { margin-bottom: 16px; }
</style>
