<template>
  <div class="coupon-inventory-page">
    <!-- 面包屑 -->
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>会员管理</a-breadcrumb-item>
      <a-breadcrumb-item>优惠券库存</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 搜索筛选 -->
    <a-card class="action-card">
      <a-space wrap>
        <a-input-search v-model="filters.keyword" placeholder="搜索优惠券名称/批次号" style="width: 260px" search-button @search="handleSearch" />
        <a-select v-model="filters.status" placeholder="库存状态" allow-clear style="width: 140px" @change="handleSearch">
          <a-option :value="1">正常</a-option>
          <a-option :value="2">库存不足</a-option>
          <a-option :value="0">已用完</a-option>
        </a-select>
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
        <template #totalCount="{ record }">
          <span>{{ record.totalCount || 0 }}</span>
        </template>
        <template #claimedCount="{ record }">
          <span style="color: #1890ff; font-weight: 600;">{{ record.claimedCount || 0 }}</span>
        </template>
        <template #remainCount="{ record }">
          <a-tag :color="getRemainColor(record)">{{ record.remainCount || 0 }}</a-tag>
        </template>
        <template #usageRate="{ record }">
          <span>{{ getUsageRate(record) }}%</span>
        </template>
        <template #status="{ record }">
          <a-tag :color="getInventoryStatusColor(record)">{{ getInventoryStatusText(record) }}</a-tag>
        </template>
        <template #actions="{ record }">
          <a-button type="text" size="small" @click="showDetail(record)">详情</a-button>
        </template>
      </a-table>
    </a-card>

    <!-- 详情抽屉 -->
    <a-drawer v-model:visible="detailVisible" title="优惠券库存详情" :width="560">
      <template v-if="currentRecord">
        <a-descriptions :column="1" bordered size="small">
          <a-descriptions-item label="优惠券名称">{{ currentRecord.couponName }}</a-descriptions-item>
          <a-descriptions-item label="批次号">{{ currentRecord.batchNo }}</a-descriptions-item>
          <a-descriptions-item label="总库存">{{ currentRecord.totalCount || 0 }}</a-descriptions-item>
          <a-descriptions-item label="已发放">{{ currentRecord.claimedCount || 0 }}</a-descriptions-item>
          <a-descriptions-item label="剩余库存">
            <a-tag :color="getRemainColor(currentRecord)">{{ currentRecord.remainCount || 0 }}</a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="使用率">{{ getUsageRate(currentRecord) }}%</a-descriptions-item>
          <a-descriptions-item label="库存状态">
            <a-tag :color="getInventoryStatusColor(currentRecord)">{{ getInventoryStatusText(currentRecord) }}</a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="有效期">{{ currentRecord.startDate || '-' }} 至 {{ currentRecord.endDate || '-' }}</a-descriptions-item>
          <a-descriptions-item label="创建时间">{{ currentRecord.createdAt || '-' }}</a-descriptions-item>
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

const filters = reactive({ keyword: '', status: undefined })
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })

const paginationConfig = computed(() => ({
  current: pagination.current, pageSize: pagination.pageSize, total: pagination.total,
  showTotal: true, showPageSize: true, pageSizeOptions: [10, 20, 50, 100]
}))

const columns = [
  { title: '优惠券名称', dataIndex: 'couponName', width: 180 },
  { title: '批次号', dataIndex: 'batchNo', width: 160 },
  { title: '总库存', slotName: 'totalCount', width: 100 },
  { title: '已发放', slotName: 'claimedCount', width: 100 },
  { title: '剩余库存', slotName: 'remainCount', width: 110 },
  { title: '使用率', slotName: 'usageRate', width: 90 },
  { title: '有效期', dataIndex: 'endDate', width: 170 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '操作', slotName: 'actions', width: 100, fixed: 'right' }
]

const getUsageRate = (record) => {
  const total = record.totalCount || 0
  if (!total) return 0
  const used = record.claimedCount || 0
  return ((used / total) * 100).toFixed(1)
}

const getRemainColor = (record) => {
  const remain = record.remainCount || 0
  const total = record.totalCount || 1
  const rate = remain / total
  if (rate <= 0) return 'red'
  if (rate <= 0.2) return 'orange'
  return 'green'
}

const getInventoryStatusColor = (record) => {
  const remain = record.remainCount || 0
  if (remain <= 0) return 'red'
  const total = record.totalCount || 1
  if (remain / total <= 0.2) return 'orange'
  return 'green'
}

const getInventoryStatusText = (record) => {
  const remain = record.remainCount || 0
  if (remain <= 0) return '已用完'
  const total = record.totalCount || 1
  if (remain / total <= 0.2) return '库存不足'
  return '正常'
}

const loadData = async () => {
  loading.value = true
  try {
    const params = { page: pagination.current, pageSize: pagination.pageSize }
    if (filters.keyword) params.keyword = filters.keyword
    if (filters.status !== undefined) params.status = filters.status
    const res = await api.getCouponInventoryList(params)
    const d = res.data || {}
    dataList.value = d.list || []
    pagination.total = d.total || 0
  } catch (err) { Message.error('加载失败: ' + err.message) }
  finally { loading.value = false }
}

const handleSearch = () => { pagination.current = 1; loadData() }
const resetFilters = () => { filters.keyword = ''; filters.status = undefined; pagination.current = 1; loadData() }
const onPageChange = (page) => { pagination.current = page; loadData() }
const onPageSizeChange = (pageSize) => { pagination.pageSize = pageSize; pagination.current = 1; loadData() }

const showDetail = (record) => { currentRecord.value = record; detailVisible.value = true }
const handleExport = () => { Message.info('导出功能开发中') }

onMounted(() => loadData())
</script>

<style scoped>
.coupon-inventory-page { padding: 20px 24px; min-height: calc(100vh - 64px); background: #f5f7fa; }
.breadcrumb { margin-bottom: 16px; }
.action-card { margin-bottom: 16px; }
.table-card { border-radius: 8px; }
</style>
