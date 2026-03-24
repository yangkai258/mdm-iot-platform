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

import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const API_BASE = '/api/v1'
const records = ref([])
const loading = ref(false)
const detailVisible = ref(false)
const current = ref(null)

const filters = reactive({ keyword: '', type: '', startDate: '', endDate: '' })
const pagination = reactive({ current: 1, pageSize: 10, total: 0 })
const stats = reactive({ total: 0, income: 0, expense: 0, monthCount: 0 })

const columns = [
  { title: 'ID', dataIndex: 'id', width: 80 },
  { title: '会员', slotName: 'member' },
  { title: '变动类型', slotName: 'type', width: 100 },
  { title: '积分变动', slotName: 'points', width: 110 },
  { title: '变动前', dataIndex: 'balance_before', width: 100 },
  { title: '变动后', dataIndex: 'balance_after', width: 100 },
  { title: '来源/用途', dataIndex: 'source', ellipsis: true },
  { title: '时间', dataIndex: 'created_at', width: 170 },
  { title: '操作', slotName: 'actions', width: 100 }
]

const formatTime = (t) => t ? new Date(t).toLocaleString('zh-CN') : '-'

const loadRecords = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const params = new URLSearchParams({ page: pagination.current, page_size: pagination.pageSize })
    if (filters.keyword) params.append('keyword', filters.keyword)
    if (filters.type) params.append('type', filters.type)
    if (filters.startDate) params.append('start_date', filters.startDate)
    if (filters.endDate) params.append('end_date', filters.endDate)
    const res = await fetch(`${API_BASE}/member/points/records?${params}`, {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    const resp = await res.json()
    if (resp.code === 0) {
      records.value = resp.data?.list || resp.data || []
      pagination.total = resp.data?.total || 0
    }
  } catch (e) {
    Message.error('加载积分流水失败')
  } finally {
    loading.value = false
  }
}

const showDetail = (record) => { current.value = record; detailVisible.value = true }
const onPageChange = (page) => { pagination.current = page; loadRecords() }
const exportData = () => Message.info('导出功能开发中')

onMounted(() => loadRecords())

</script>

<style scoped>
.page-container { background: #fff; border-radius: 4px; padding: 20px; }
.search-form { margin-bottom: 16px; padding: 16px; background: #f7f8fa; border-radius: 4px; }
.toolbar { margin-bottom: 16px; }
</style>
