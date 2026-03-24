<template>
  <div class="member-orders-page">
    <!-- 面包屑 -->
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>会员管理</a-breadcrumb-item>
      <a-breadcrumb-item>会员订单</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 搜索筛选区 -->
    <a-card class="action-card">
      <a-space wrap>
        <a-input-search
          v-model="filters.keyword"
          placeholder="搜索订单编号/会员名称"
          style="width: 220px"
          search-button
          @search="handleSearch"
        />
        <a-select v-model="filters.payStatus" placeholder="支付状态" allow-clear style="width: 130px" @change="handleSearch">
          <a-option :value="1">已支付</a-option>
          <a-option :value="0">未支付</a-option>
          <a-option :value="2">已退款</a-option>
        </a-select>
        <a-range-picker v-model="filters.dateRange" style="width: 260px" @change="handleSearch" />
        <a-button @click="handleSearch">筛选</a-button>
        <a-button @click="resetFilters">重置</a-button>
      </a-space>
    </a-card>

    <!-- 操作+表格 -->
    <a-card class="table-card">
      <template #title>
        <a-space>
          <span style="font-weight: 600; font-size: 15px;">会员订单</span>
          <a-badge :count="pagination.total" :max-count="99999" />
        </a-space>
      </template>
      <template #extra>
        <a-button @click="handleExport">导出</a-button>
      </template>

      <a-table
        :columns="columns"
        :data="dataList"
        :loading="loading"
        :pagination="paginationConfig"
        @page-change="onPageChange"
        @page-size-change="onPageSizeChange"
        row-key="id"
        :scroll="{ x: 1100 }"
      >
        <template #payStatus="{ record }">
          <a-tag :color="getPayStatusColor(record.payStatus)">{{ getPayStatusText(record.payStatus) }}</a-tag>
        </template>
        <template #amount="{ record }">
          <span style="color: #ff4d4f; font-weight: 600;">¥{{ (record.amount || 0).toFixed(2) }}</span>
        </template>
        <template #actions="{ record }">
          <a-button type="text" size="small" @click="showDetail(record)">查看详情</a-button>
        </template>
      </a-table>
    </a-card>

    <!-- 订单详情抽屉 -->
    <a-drawer v-model:visible="detailVisible" title="订单详情" :width="520">
      <template v-if="currentOrder">
        <a-descriptions :column="1" bordered size="small">
          <a-descriptions-item label="订单编号">{{ currentOrder.orderNo }}</a-descriptions-item>
          <a-descriptions-item label="会员名称">{{ currentOrder.memberName }}</a-descriptions-item>
          <a-descriptions-item label="会员手机">{{ currentOrder.memberMobile || '-' }}</a-descriptions-item>
          <a-descriptions-item label="商品信息">{{ currentOrder.goodsName || '-' }}</a-descriptions-item>
          <a-descriptions-item label="订单金额">
            <span style="color: #ff4d4f; font-weight: 600;">¥{{ (currentOrder.amount || 0).toFixed(2) }}</span>
          </a-descriptions-item>
          <a-descriptions-item label="支付状态">
            <a-tag :color="getPayStatusColor(currentOrder.payStatus)">{{ getPayStatusText(currentOrder.payStatus) }}</a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="支付方式">{{ currentOrder.payMethod || '-' }}</a-descriptions-item>
          <a-descriptions-item label="支付时间">{{ currentOrder.payTime || '-' }}</a-descriptions-item>
          <a-descriptions-item label="门店">{{ currentOrder.storeName || '-' }}</a-descriptions-item>
          <a-descriptions-item label="创建时间">{{ currentOrder.createdAt || '-' }}</a-descriptions-item>
        </a-descriptions>
      </template>
    </a-drawer>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import * as api from '@/api/member'

const dataList = ref([])
const loading = ref(false)
const detailVisible = ref(false)
const currentOrder = ref(null)

const filters = reactive({
  keyword: '',
  payStatus: undefined,
  dateRange: []
})
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })

const paginationConfig = computed(() => ({
  current: pagination.current,
  pageSize: pagination.pageSize,
  total: pagination.total,
  showTotal: true,
  showPageSize: true,
  pageSizeOptions: [10, 20, 50, 100]
}))

const columns = [
  { title: '订单编号', dataIndex: 'orderNo', width: 200 },
  { title: '会员名称', dataIndex: 'memberName', width: 120 },
  { title: '商品', dataIndex: 'goodsName', width: 160, ellipsis: true },
  { title: '金额', slotName: 'amount', width: 100 },
  { title: '支付状态', slotName: 'payStatus', width: 90 },
  { title: '支付方式', dataIndex: 'payMethod', width: 100 },
  { title: '创建时间', dataIndex: 'createdAt', width: 170 },
  { title: '操作', slotName: 'actions', width: 100, fixed: 'right' }
]

const getPayStatusColor = (s) => ({ 1: 'green', 0: 'orange', 2: 'gray' }[s] || 'gray')
const getPayStatusText = (s) => ({ 1: '已支付', 0: '未支付', 2: '已退款' }[s] || '未知')

const loadData = async () => {
  loading.value = true
  try {
    const params = { page: pagination.current, pageSize: pagination.pageSize }
    if (filters.keyword) params.keyword = filters.keyword
    if (filters.payStatus !== undefined) params.payStatus = filters.payStatus
    if (filters.dateRange && filters.dateRange.length === 2) {
      params.startDate = filters.dateRange[0]
      params.endDate = filters.dateRange[1]
    }

    const res = await api.getMemberOrderList(params)
    const d = res.data || {}
    dataList.value = d.list || []
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

const resetFilters = () => {
  filters.keyword = ''
  filters.payStatus = undefined
  filters.dateRange = []
  pagination.current = 1
  loadData()
}

const onPageChange = (page) => {
  pagination.current = page
  loadData()
}

const onPageSizeChange = (pageSize) => {
  pagination.pageSize = pageSize
  pagination.current = 1
  loadData()
}

const showDetail = async (record) => {
  try {
    const res = await api.getMemberOrderDetail(record.id)
    currentOrder.value = res.data || record
  } catch {
    currentOrder.value = record
  }
  detailVisible.value = true
}

const handleExport = () => {
  Message.info('导出功能开发中')
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.member-orders-page {
  padding: 20px 24px;
  min-height: calc(100vh - 64px);
  background: #f5f7fa;
}
.breadcrumb { margin-bottom: 16px; }
.action-card { margin-bottom: 16px; }
.table-card { border-radius: 8px; }
</style>
