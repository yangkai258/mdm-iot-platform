<template>
  <div class="container">
    <Breadcrumb :items="['menu.members', 'menu.members.giftRecords']" />
    <a-card class="general-card" title="礼包记录">
      <template #extra>
        <a-button @click="loadData"><icon-refresh />刷新</a-button>
      </template>
      <a-row :gutter="16">
        <a-col :span="8">
          <a-form-item label="关键词">
            <a-input v-model="filters.keyword" placeholder="请输入" @pressEnter="loadData" />
          </a-form-item>
        </a-col>
        <a-col :flex="'86px'" style="display: flex; align-items: flex-end">
          <a-space direction="vertical" :size="8">
            <a-button type="primary" @click="loadData">查询</a-button>
            <a-button @click="Object.keys(filters).forEach(k => filters[k] = ''); loadData()">重置</a-button>
          </a-space>
        </a-col>
      </a-row>
      <a-divider style="margin: 0 0 16px 0" />
      <a-table :columns="columns" :data="dataList" :loading="loading" :pagination="paginationConfig" @page-change="onPageChange" row-key="id" />
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { Message } from '@arco-design/web-vue'
import Breadcrumb from '@/components/Breadcrumb.vue'

const loading = ref(false)
const dataList = ref([])
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
  { title: '操作', slotName: 'actions', width: 120 }
]

const statusColor = (s) => ({ claimed: 'green', unclaimed: 'orange', expired: 'gray' }[s] || 'gray')
const statusText = (s) => ({ claimed: '已领取', unclaimed: '未领取', expired: '已过期' }[s] || s)

const mockData = () => [
  { id: 1, memberName: '张三', phone: '13800138001', giftName: '新会员礼包', giftValue: 50, grantTime: '2026-03-01 10:00:00', claimTime: '2026-03-01 11:30:00', status: 'claimed' },
  { id: 2, memberName: '李四', phone: '13800138002', giftName: '生日专属礼包', giftValue: 100, grantTime: '2026-03-05 09:00:00', claimTime: '-', status: 'unclaimed' },
  { id: 3, memberName: '王五', phone: '13800138003', giftName: '节日礼包', giftValue: 80, grantTime: '2026-02-14 08:00:00', claimTime: '-', status: 'expired' }
]

const loadData = () => {
  loading.value = true
  setTimeout(() => {
    dataList.value = mockData()
    paginationConfig.total = dataList.value.length
    loading.value = false
  }, 400)
}

const onPageChange = (page) => { paginationConfig.current = page; loadData() }

loadData()
</script>
