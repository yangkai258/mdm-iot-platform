<template>
  <div class="member-page">
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>会员管理</a-breadcrumb-item>
      <a-breadcrumb-item>礼包发放明细</a-breadcrumb-item>
    </a-breadcrumb>

    <a-row :gutter="16" class="stats-row">
      <a-col :span="8">
        <a-card class="stat-card">
          <a-statistic title="发放总次数" :value="stats.total || 0" />
        </a-card>
      </a-col>
      <a-col :span="8">
        <a-card class="stat-card">
          <a-statistic title="已领取" :value="stats.claimed || 0" />
        </a-card>
      </a-col>
      <a-col :span="8">
        <a-card class="stat-card">
          <a-statistic title="未领取" :value="stats.unclaimed || 0" />
        </a-card>
      </a-col>
    </a-row>

    <a-card class="action-card">
      <a-space wrap>
        <a-input-search v-model="filters.keyword" placeholder="搜索会员名称/手机号" style="width: 240px" search-button @search="loadData" />
        <a-select v-model="filters.giftName" placeholder="礼包名称" allow-clear style="width: 160px" @change="loadData">
          <a-option v-for="g in giftOptions" :key="g" :value="g">{{ g }}</a-option>
        </a-select>
        <a-select v-model="filters.status" placeholder="状态" allow-clear style="width: 120px" @change="loadData">
          <a-option value="claimed">已领取</a-option>
          <a-option value="unclaimed">未领取</a-option>
          <a-option value="expired">已过期</a-option>
        </a-select>
        <a-range-picker v-model="filters.dateRange" style="width: 260px;" @change="loadData" />
        <a-button type="primary" @click="loadData">搜索</a-button>
        <a-button @click="handleExport">导出</a-button>
        <a-button @click="loadData">刷新</a-button>
      </a-space>
    </a-card>

    <a-card class="table-card">
      <a-table :columns="columns" :data="dataList" :loading="loading" :pagination="paginationConfig" @page-change="onPageChange" row-key="id" :scroll="{ x: 1100 }">
        <template #status="{ record }">
          <a-tag :color="statusColor(record.status)">{{ statusText(record.status) }}</a-tag>
        </template>
        <template #actions="{ record }">
          <a-button type="text" size="small" @click="showDetail(record)">查看详情</a-button>
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { Message } from '@arco-design/web-vue'

const loading = ref(false)
const dataList = ref([])
const giftOptions = ref(['新会员礼包', '生日专属礼包', '节日礼包'])
const stats = ref({ total: 0, claimed: 0, unclaimed: 0 })
const filters = reactive({ keyword: '', giftName: '', status: '', dateRange: [] })
const paginationConfig = reactive({ current: 1, pageSize: 10, total: 0 })

const columns = [
  { title: '会员名称', dataIndex: 'memberName', width: 150 },
  { title: '手机号', dataIndex: 'phone', width: 130 },
  { title: '礼包名称', dataIndex: 'giftName', width: 180 },
  { title: '礼包价值', dataIndex: 'giftValue', width: 100 },
  { title: '发放时间', dataIndex: 'grantTime', width: 180 },
  { title: '领取时间', dataIndex: 'claimTime', width: 180 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '操作', slotName: 'actions', width: 120, fixed: 'right' }
]

const mockData = () => [
  { id: 1, memberName: '张三', phone: '13800138001', giftName: '新会员礼包', giftValue: 50, grantTime: '2026-03-01 10:00:00', claimTime: '2026-03-01 11:30:00', status: 'claimed' },
  { id: 2, memberName: '李四', phone: '13800138002', giftName: '生日专属礼包', giftValue: 100, grantTime: '2026-03-05 09:00:00', claimTime: '-', status: 'unclaimed' },
  { id: 3, memberName: '王五', phone: '13800138003', giftName: '节日礼包', giftValue: 80, grantTime: '2026-02-14 08:00:00', claimTime: '-', status: 'expired' }
]

const statusColor = (s) => ({ claimed: 'green', unclaimed: 'orange', expired: 'gray' }[s] || 'gray')
const statusText = (s) => ({ claimed: '已领取', unclaimed: '未领取', expired: '已过期' }[s] || s)
const showDetail = (record) => { Message.info(`查看 ${record.memberName} 的 ${record.giftName} 明细`) }

const loadData = () => {
  loading.value = true
  setTimeout(() => {
    dataList.value = mockData()
    stats.value = { total: 3, claimed: 1, unclaimed: 1 }
    paginationConfig.total = dataList.value.length
    loading.value = false
  }, 400)
}

const onPageChange = (page) => { paginationConfig.current = page; loadData() }
const handleExport = () => { Message.success('导出成功') }

loadData()
</script>

<style scoped>
.member-page { padding: 20px; }
.breadcrumb { margin-bottom: 16px; }
.stats-row { margin-bottom: 16px; }
.stat-card { text-align: center; }
.action-card { margin-bottom: 16px; }
</style>
