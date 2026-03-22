<template>
  <div class="coupon-messages-page">
    <!-- 面包屑 -->
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>会员管理</a-breadcrumb-item>
      <a-breadcrumb-item>优惠消息流水</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 搜索筛选 -->
    <a-card class="action-card">
      <a-space wrap>
        <a-input-search v-model="filters.keyword" placeholder="搜索会员名称/优惠券名称" style="width: 260px" search-button @search="handleSearch" />
        <a-select v-model="filters.messageType" placeholder="消息类型" allow-clear style="width: 140px" @change="handleSearch">
          <a-option value="grant">发放</a-option>
          <a-option value="use">使用</a-option>
          <a-option value="expire">过期</a-option>
          <a-option value="remind">提醒</a-option>
        </a-select>
        <a-select v-model="filters.sendStatus" placeholder="发送状态" allow-clear style="width: 120px" @change="handleSearch">
          <a-option :value="1">成功</a-option>
          <a-option :value="0">失败</a-option>
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
        <template #messageType="{ record }">
          <a-tag :color="getMsgTypeColor(record.messageType)">{{ getMsgTypeText(record.messageType) }}</a-tag>
        </template>
        <template #sendStatus="{ record }">
          <a-tag :color="record.sendStatus === 1 ? 'green' : 'red'">{{ record.sendStatus === 1 ? '成功' : '失败' }}</a-tag>
        </template>
        <template #couponValue="{ record }">
          <span style="color: #ff6b00; font-weight: 600;">{{ record.couponValue || 0 }}</span>
        </template>
        <template #actions="{ record }">
          <a-button type="text" size="small" @click="showDetail(record)">详情</a-button>
        </template>
      </a-table>
    </a-card>

    <!-- 详情抽屉 -->
    <a-drawer v-model:visible="detailVisible" title="消息详情" :width="520">
      <template v-if="currentRecord">
        <a-descriptions :column="1" bordered size="small">
          <a-descriptions-item label="消息ID">{{ currentRecord.id }}</a-descriptions-item>
          <a-descriptions-item label="优惠券名称">{{ currentRecord.couponName }}</a-descriptions-item>
          <a-descriptions-item label="会员名称">{{ currentRecord.memberName }}</a-descriptions-item>
          <a-descriptions-item label="会员手机">{{ currentRecord.memberMobile || '-' }}</a-descriptions-item>
          <a-descriptions-item label="消息类型">
            <a-tag :color="getMsgTypeColor(currentRecord.messageType)">{{ getMsgTypeText(currentRecord.messageType) }}</a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="优惠券面值">{{ currentRecord.couponValue || 0 }}</a-descriptions-item>
          <a-descriptions-item label="发送状态">
            <a-tag :color="currentRecord.sendStatus === 1 ? 'green' : 'red'">{{ currentRecord.sendStatus === 1 ? '成功' : '失败' }}</a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="发送时间">{{ currentRecord.createdAt || '-' }}</a-descriptions-item>
          <a-descriptions-item label="失败原因">{{ currentRecord.failReason || '-' }}</a-descriptions-item>
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

const filters = reactive({ keyword: '', messageType: '', sendStatus: undefined, dateRange: [] })
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })

const paginationConfig = computed(() => ({
  current: pagination.current, pageSize: pagination.pageSize, total: pagination.total,
  showTotal: true, showPageSize: true, pageSizeOptions: [10, 20, 50, 100]
}))

const columns = [
  { title: '消息ID', dataIndex: 'id', width: 80 },
  { title: '优惠券名称', dataIndex: 'couponName', width: 180 },
  { title: '会员名称', dataIndex: 'memberName', width: 120 },
  { title: '会员手机', dataIndex: 'memberMobile', width: 130 },
  { title: '优惠券面值', slotName: 'couponValue', width: 110 },
  { title: '消息类型', slotName: 'messageType', width: 90 },
  { title: '发送状态', slotName: 'sendStatus', width: 90 },
  { title: '创建时间', dataIndex: 'createdAt', width: 170 },
  { title: '操作', slotName: 'actions', width: 100, fixed: 'right' }
]

const msgTypeColorMap = { grant: 'blue', use: 'green', expire: 'orange', remind: 'purple' }
const msgTypeTextMap = { grant: '发放', use: '使用', expire: '过期', remind: '提醒' }

const getMsgTypeColor = (t) => msgTypeColorMap[t] || 'gray'
const getMsgTypeText = (t) => msgTypeTextMap[t] || t

const loadData = async () => {
  loading.value = true
  try {
    const params = { page: pagination.current, pageSize: pagination.pageSize }
    if (filters.keyword) params.keyword = filters.keyword
    if (filters.messageType) params.messageType = filters.messageType
    if (filters.sendStatus !== undefined) params.sendStatus = filters.sendStatus
    if (filters.dateRange && filters.dateRange.length === 2) {
      params.startDate = filters.dateRange[0]
      params.endDate = filters.dateRange[1]
    }
    const res = await api.getCouponMessageList(params)
    const d = res.data || {}
    dataList.value = d.list || []
    pagination.total = d.total || 0
  } catch (err) { Message.error('加载失败: ' + err.message) }
  finally { loading.value = false }
}

const handleSearch = () => { pagination.current = 1; loadData() }
const resetFilters = () => { filters.keyword = ''; filters.messageType = ''; filters.sendStatus = undefined; filters.dateRange = []; pagination.current = 1; loadData() }
const onPageChange = (page) => { pagination.current = page; loadData() }
const onPageSizeChange = (pageSize) => { pagination.pageSize = pageSize; pagination.current = 1; loadData() }

const showDetail = (record) => { currentRecord.value = record; detailVisible.value = true }
const handleExport = () => { Message.info('导出功能开发中') }

onMounted(() => loadData())
</script>

<style scoped>
.coupon-messages-page { padding: 20px 24px; min-height: calc(100vh - 64px); background: #f5f7fa; }
.breadcrumb { margin-bottom: 16px; }
.action-card { margin-bottom: 16px; }
.table-card { border-radius: 8px; }
</style>
