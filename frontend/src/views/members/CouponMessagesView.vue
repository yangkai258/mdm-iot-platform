<template>
  <div class="coupon-messages-page">
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>会员管理</a-breadcrumb-item>
      <a-breadcrumb-item>优惠消息流水</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 搜索筛选 -->
    <a-card class="action-card">
      <a-space wrap>
        <a-select v-model="filters.type" placeholder="消息类型" allow-clear style="width: 140px" @change="loadData">
          <a-option value="grant">发放通知</a-option>
          <a-option value="remind">到期提醒</a-option>
          <a-option value="use">核销通知</a-option>
          <a-option value="expire">过期通知</a-option>
        </a-select>
        <a-input-search v-model="filters.memberName" placeholder="搜索会员名称" style="width: 200px" search-button @search="loadData" />
        <a-select v-model="filters.status" placeholder="状态" allow-clear style="width: 120px" @change="loadData">
          <a-option value="success">已发送</a-option>
          <a-option value="failed">发送失败</a-option>
          <a-option value="pending">待发送</a-option>
        </a-select>
        <a-range-picker v-model="filters.dateRange" style="width: 260px" @change="loadData" />
        <a-button @click="handleExport">导出</a-button>
        <a-button @click="loadData">刷新</a-button>
      </a-space>
    </a-card>

    <!-- 消息流水列表 -->
    <a-card class="table-card">
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
        <template #type="{ record }">
          <a-tag :color="getTypeColor(record.type)">{{ getTypeText(record.type) }}</a-tag>
        </template>
        <template #status="{ record }">
          <a-tag :color="getStatusColor(record.status)">{{ getStatusText(record.status) }}</a-tag>
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import * as api from '@/api/member'

const dataList = ref([])
const loading = ref(false)

const filters = reactive({ type: undefined, memberName: '', status: undefined, dateRange: [] })
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })

const paginationConfig = computed(() => ({
  current: pagination.current, pageSize: pagination.pageSize, total: pagination.total,
  showTotal: true, showPageSize: true, pageSizeOptions: [10, 20, 50, 100]
}))

const columns = [
  { title: '消息类型', slotName: 'type', width: 120 },
  { title: '会员名称', dataIndex: 'memberName', width: 150 },
  { title: '优惠券名称', dataIndex: 'couponName', width: 160 },
  { title: '消息内容', dataIndex: 'content', width: 280, ellipsis: true },
  { title: '发送时间', dataIndex: 'sendTime', width: 170 },
  { title: '状态', slotName: 'status', width: 100 }
]

const getTypeColor = (t) => ({ grant: 'blue', remind: 'orange', use: 'green', expire: 'purple' }[t] || 'gray')
const getTypeText = (t) => ({ grant: '发放通知', remind: '到期提醒', use: '核销通知', expire: '过期通知' }[t] || t)
const getStatusColor = (s) => ({ success: 'green', failed: 'red', pending: 'orange' }[s] || 'gray')
const getStatusText = (s) => ({ success: '已发送', failed: '发送失败', pending: '待发送' }[s] || s)

const loadData = async () => {
  loading.value = true
  try {
    const params = { page: pagination.current, pageSize: pagination.pageSize }
    if (filters.type) params.type = filters.type
    if (filters.memberName) params.memberName = filters.memberName
    if (filters.status) params.status = filters.status
    if (filters.dateRange && filters.dateRange.length === 2) {
      params.startDate = filters.dateRange[0]
      params.endDate = filters.dateRange[1]
    }
    const res = await api.getCouponMessageList(params)
    const d = res.data || {}
    dataList.value = d.list || []
    pagination.total = d.total || 0
  } catch (err) {
    Message.error('加载消息流水失败: ' + err.message)
  } finally {
    loading.value = false
  }
}

const handleExport = () => {
  Message.info('导出功能开发中')
}

const onPageChange = (page) => { pagination.current = page; loadData() }
const onPageSizeChange = (pageSize) => { pagination.pageSize = pageSize; pagination.current = 1; loadData() }

onMounted(() => loadData())
</script>

<style scoped>
.coupon-messages-page { padding: 20px 24px; min-height: calc(100vh - 64px); background: #f5f7fa; }
.breadcrumb { margin-bottom: 16px; }
.action-card { margin-bottom: 16px; }
.table-card { border-radius: 8px; }
</style>
