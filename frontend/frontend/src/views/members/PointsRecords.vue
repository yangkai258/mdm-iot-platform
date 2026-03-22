<template>
  <div class="pro-page-container">
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>会员管理</a-breadcrumb-item>
      <a-breadcrumb-item>积分流水</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 统计卡片 -->
    <a-row :gutter="16" class="stats-row">
      <a-col :span="6">
        <a-card class="stat-card"><a-statistic title="总流水笔数" :value="stats.total" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card"><a-statistic title="收入积分" :value="stats.income" :value-style="{ color: '#52c41a' }" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card"><a-statistic title="支出积分" :value="stats.expense" :value-style="{ color: '#ff4d4f' }" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card"><a-statistic title="本月笔数" :value="stats.monthCount" /></a-card>
      </a-col>
    </a-row>

    <!-- 搜索栏 -->
    <div class="pro-search-bar">
      <a-space wrap>
        <a-input-search v-model="filters.keyword" placeholder="搜索会员名称/手机" style="width: 220px" search-button @search="loadRecords" />
        <a-select v-model="filters.type" placeholder="变动类型" allow-clear style="width: 120px" @change="loadRecords">
          <a-option value="add">收入</a-option>
          <a-option value="deduct">支出</a-option>
        </a-select>
        <a-input v-model="filters.startDate" type="date" placeholder="开始日期" style="width: 150px" @change="loadRecords" />
        <a-input v-model="filters.endDate" type="date" placeholder="结束日期" style="width: 150px" @change="loadRecords" />
      </a-space>
    </div>

    <!-- 操作栏 -->
    <div class="pro-action-bar">
      <a-space>
        <a-button @click="loadRecords">刷新</a-button>
        <a-button @click="exportData">导出</a-button>
      </a-space>
    </div>

    <!-- 数据表格 -->
    <div class="pro-content-area">
      <a-table
        :columns="columns"
        :data="records"
        :loading="loading"
        :pagination="pagination"
        @page-change="onPageChange"
        row-key="id"
      >
        <template #type="{ record }">
          <a-tag :color="record.type === 'add' ? 'green' : 'red'">{{ record.type === 'add' ? '收入' : '支出' }}</a-tag>
        </template>
        <template #points="{ record }">
          <span :style="{ color: record.type === 'add' ? '#52c41a' : '#ff4d4f', fontWeight: 600 }">
            {{ record.type === 'add' ? '+' : '-' }}{{ record.points }}
          </span>
        </template>
        <template #member="{ record }">
          {{ record.member_name }} ({{ record.phone }})
        </template>
        <template #actions="{ record }">
          <a-button type="text" size="small" @click="showDetail(record)">详情</a-button>
        </template>
      </a-table>
    </div>

    <!-- 详情抽屉 -->
    <a-drawer v-model:visible="detailVisible" title="流水详情" :width="480">
      <template v-if="current">
        <a-descriptions :column="1" bordered size="small">
          <a-descriptions-item label="会员">{{ current.member_name }} ({{ current.phone }})</a-descriptions-item>
          <a-descriptions-item label="变动类型">
            <a-tag :color="current.type === 'add' ? 'green' : 'red'">{{ current.type === 'add' ? '收入' : '支出' }}</a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="积分变动">
            <span :style="{ color: current.type === 'add' ? '#52c41a' : '#ff4d4f', fontWeight: 600 }">
              {{ current.type === 'add' ? '+' : '-' }}{{ current.points }}
            </span>
          </a-descriptions-item>
          <a-descriptions-item label="变动前积分">{{ current.balance_before || '-' }}</a-descriptions-item>
          <a-descriptions-item label="变动后积分">{{ current.balance_after || '-' }}</a-descriptions-item>
          <a-descriptions-item label="来源/用途">{{ current.source || current.usage || '-' }}</a-descriptions-item>
          <a-descriptions-item label="备注">{{ current.remark || '-' }}</a-descriptions-item>
          <a-descriptions-item label="时间">{{ formatTime(current.created_at) }}</a-descriptions-item>
        </a-descriptions>
      </template>
    </a-drawer>
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
.pro-page-container { padding: 20px 24px; min-height: calc(100vh - 64px); background: #f5f7fa; }
.pro-breadcrumb { margin-bottom: 16px; }
.stats-row { margin-bottom: 16px; }
.stat-card { border-radius: 8px; text-align: center; }
.pro-search-bar { margin-bottom: 12px; }
.pro-action-bar { margin-bottom: 16px; }
.pro-content-area { background: #fff; border-radius: 8px; padding: 20px; box-shadow: 0 1px 3px rgba(0,0,0,0.04); }
</style>
