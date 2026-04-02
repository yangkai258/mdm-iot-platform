<template>
    <Breadcrumb :items="['Home','Console','']" />


  <div class="member-orders-page">
    <!-- 面包屑 -->
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>会员管理</a-breadcrumb-item>
      <a-breadcrumb-item>会员订单</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 搜索筛选 -->
    <a-card class="action-card">
      <a-space wrap>
        <a-input-search v-model="filters.keyword" placeholder="搜索订单号/会员名称" style="width: 240px" search-button @search="handleSearch" />
        <a-select v-model="filters.status" placeholder="订单状态" allow-clear style="width: 140px" @change="handleSearch">
          <a-option :value="1">待支付</a-option>
          <a-option :value="2">已支付</a-option>
          <a-option :value="3">已完成</a-option>
          <a-option :value="4">已取消</a-option>
          <a-option :value="5">已退款</a-option>
        </a-select>
        <a-range-picker v-model="filters.dateRange" style="width: 240px" @change="handleSearch" />
        <a-button @click="handleSearch">「搜索」</a-button>
        <a-button @click="resetFilters">「重置」</a-button>
      </a-space>
    </a-card>

    <!-- 操作栏 -->
    <a-card class="table-card">
      <template #extra>
        <a-space>
          <a-button @click="handleExport">「导出」</a-button>
          <a-button @click="loadData">🔄</a-button>
        </a-space>
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
        <template #orderAmount="{ record }">
          <span style="color: #ff4d4f; font-weight: 600;">¥{{ (record.orderAmount || 0).toFixed(2) }}</span>
        </template>
        <template #pointsUsed="{ record }">
          <span style="color: #ff6b00; font-weight: 600;">{{ record.pointsUsed || 0 }}</span>
        </template>
        <template #status="{ record }">
          <a-tag :color="getStatusColor(record.status)">{{ getStatusText(record.status) }}</a-tag>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="showDetail(record)">详情</a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 详情抽屉 -->
    <a-drawer v-model:visible="detailVisible" title="订单详情" :width="560">
      <template v-if="currentRecord">
        <a-descriptions :column="1" bordered size="small">
          <a-descriptions-item label="订单编号">{{ currentRecord.orderNo }}</a-descriptions-item>
          <a-descriptions-item label="会员名称">{{ currentRecord.memberName }}</a-descriptions-item>
          <a-descriptions-item label="会员手机">{{ currentRecord.memberMobile || '-' }}</a-descriptions-item>
          <a-descriptions-item label="订单金额">
            <span style="color: #ff4d4f; font-weight: 600;">¥{{ (currentRecord.orderAmount || 0).toFixed(2) }}</span>
          </a-descriptions-item>
          <a-descriptions-item label="实付金额">¥{{ (currentRecord.paidAmount || 0).toFixed(2) }}</a-descriptions-item>
          <a-descriptions-item label="使用积分">{{ currentRecord.pointsUsed || 0 }}分</a-descriptions-item>
          <a-descriptions-item label="积分抵扣">¥{{ (currentRecord.pointsDiscount || 0).toFixed(2) }}</a-descriptions-item>
          <a-descriptions-item label="优惠金额">¥{{ (currentRecord.couponDiscount || 0).toFixed(2) }}</a-descriptions-item>
          <a-descriptions-item label="订单状态">
            <a-tag :color="getStatusColor(currentRecord.status)">{{ getStatusText(currentRecord.status) }}</a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="所属门店">{{ currentRecord.storeName || '-' }}</a-descriptions-item>
          <a-descriptions-item label="支付方式">{{ currentRecord.payMethod || '-' }}</a-descriptions-item>
          <a-descriptions-item label="创建时间">{{ currentRecord.createdAt || '-' }}</a-descriptions-item>
          <a-descriptions-item label="支付时间">{{ currentRecord.paidAt || '-' }}</a-descriptions-item>
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
const currentRecord = ref(null)

const filters = reactive({ keyword: '', status: undefined, dateRange: [] })
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })

const paginationConfig = computed(() => ({
  current: pagination.current, pageSize: pagination.pageSize, total: pagination.total,
  showTotal: true, showPageSize: true, pageSizeOptions: [10, 20, 50, 100]
}))

const columns = [
  { title: '订单编号', dataIndex: 'orderNo', width: 180 },
  { title: '会员名称', dataIndex: 'memberName', width: 120 },
  { title: '会员手机', dataIndex: 'memberMobile', width: 130 },
  { title: '订单金额', slotName: 'orderAmount', width: 120 },
  { title: '实付金额', dataIndex: 'paidAmount', width: 110, render: (_, r) => `¥${(r.paidAmount || 0).toFixed(2)}` },
  { title: '使用积分', slotName: 'pointsUsed', width: 100 },
  { title: '状态', slotName: 'status', width: 90 },
  { title: '门店', dataIndex: 'storeName', width: 130 },
  { title: '创建时间', dataIndex: 'createdAt', width: 170 },
  { title: '操作', slotName: 'actions', width: 100, fixed: 'right' }
]

const statusColorMap = { 1: 'orange', 2: 'blue', 3: 'green', 4: 'gray', 5: 'red' }
const statusTextMap = { 1: '待支付', 2: '已支付', 3: '已完成', 4: '已取消', 5: '已退款' }

const getStatusColor = (s) => statusColorMap[s] || 'gray'
const getStatusText = (s) => statusTextMap[s] || '未知'

const loadData = async () => {
  loading.value = true
  try {
    const params = { page: pagination.current, pageSize: pagination.pageSize }
    if (filters.keyword) params.keyword = filters.keyword
    if (filters.status !== undefined) params.status = filters.status
    if (filters.dateRange && filters.dateRange.length === 2) {
      params.startDate = filters.dateRange[0]
      params.endDate = filters.dateRange[1]
    }
    const res = await api.getMemberOrderList(params)
    const d = res.data || {}
    dataList.value = d.list || []
    pagination.total = d.total || 0
  } catch (err) { Message.error('加载失败: ' + err.message) }
  finally { loading.value = false }
}

const handleSearch = () => { pagination.current = 1; loadData() }
const resetFilters = () => { filters.keyword = ''; filters.status = undefined; filters.dateRange = []; pagination.current = 1; loadData() }
const onPageChange = (page) => { pagination.current = page; loadData() }
const onPageSizeChange = (pageSize) => { pagination.pageSize = pageSize; pagination.current = 1; loadData() }

const showDetail = async (record) => {
  try {
    const res = await api.getMemberOrderDetail(record.id)
    currentRecord.value = res.data || record
  } catch { currentRecord.value = record }
  detailVisible.value = true
}

const handleExport = () => { Message.info('导出功能开发中') }

onMounted(() => loadData())
</script>

<style scoped>
.member-orders-page { padding: 20px 24px; min-height: calc(100vh - 64px); background: #f5f7fa; }
.breadcrumb { margin-bottom: 16px; }
.action-card { margin-bottom: 16px; }
.table-card { border-radius: 8px; }
</style>
