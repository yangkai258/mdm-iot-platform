<template>
  <div class="container">
    <Breadcrumb :items="['menu.members', 'menu.members.points', 'menu.members.pointsRecords']" />
    <a-card class="general-card" title="积分流水">
      <template #extra>
        <a-button @click="loadRecords"><icon-refresh />刷新</a-button>
      </template>
      <a-row :gutter="16">
        <a-col :span="8">
          <a-form-item label="关键词">
            <a-input v-model="filters.keyword" placeholder="请输入" @pressEnter="loadRecords" />
          </a-form-item>
        </a-col>
        <a-col :flex="'86px'" style="display: flex; align-items: flex-end">
          <a-space direction="vertical" :size="8">
            <a-button type="primary" @click="loadRecords">查询</a-button>
            <a-button @click="Object.keys(filters).forEach(k => filters[k] = ''); loadRecords()">重置</a-button>
          </a-space>
        </a-col>
      </a-row>
      <a-divider style="margin: 0 0 16px 0" />
      <a-table :columns="columns" :data="records" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id" />
    </a-card>
  </div>
</template>

<script setup>

import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import Breadcrumb from '@/components/Breadcrumb.vue'

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
