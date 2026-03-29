<template>
  <div class="gift-records-page">
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>会员管理</a-breadcrumb-item>
      <a-breadcrumb-item><router-link to="/member/gifts">会员礼包</router-link></a-breadcrumb-item>
      <a-breadcrumb-item>礼包发放明细</a-breadcrumb-item>
    </a-breadcrumb>

    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="8"><a-card hoverable><a-statistic title="发放总数" :value="stats.totalCount || 0" /></a-card></a-col>
      <a-col :span="8"><a-card hoverable><a-statistic title="已领取" :value="stats.claimedCount || 0" /></a-card></a-col>
      <a-col :span="8"><a-card hoverable><a-statistic title="未领取" :value="stats.unclaimedCount || 0" /></a-card></a-col>
    </a-row>

    <a-card class="action-card">
      <template #extra>
        <a-space>
          <a-select v-model="filterGiftId" placeholder="筛选礼包" style="width: 160px;" allow-clear>
            <a-option v-for="g in giftList" :key="g.id" :value="g.id">{{ g.name }}</a-option>
          </a-select>
          <a-select v-model="filterStatus" placeholder="状态" style="width: 120px;" allow-clear>
            <a-option value="pending">未领取</a-option>
            <a-option value="claimed">已领取</a-option>
            <a-option value="expired">已过期</a-option>
          </a-select>
          <a-button @click="loadData">筛选</a-button>
          <a-button @click="exportRecords">导出</a-button>
        </a-space>
      </template>
    </a-card>

    <a-card>
      <a-table :columns="columns" :data="recordList" :loading="loading" row-key="id" :pagination="{ pageSize: 10 }">
        <template #giftName="{ record }"><a-tag color="arcoblue">{{ record.giftName }}</a-tag></template>
        <template #memberName="{ record }">
          <a-space>
            <a-avatar :size="20">{{ record.memberName?.charAt(0) || '?' }}</a-avatar>
            {{ record.memberName || record.memberId }}
          </a-space>
        </template>
      </a-table>
        <template #grantTime="{ record }">{{ record.grantTime?.slice(0, 19) || '-' }}</template>
        <template #claimTime="{ record }">{{ record.claimTime?.slice(0, 19) || '-' }}</template>
        <template #status="{ record }">
          <a-tag :color="getStatusColor(record.status)">{{ getStatusName(record.status) }}</a-tag>
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import * as api from '@/api/marketing'

const recordList = ref([])
const giftList = ref([])
const stats = ref({})
const loading = ref(false)
const filterGiftId = ref(null)
const filterStatus = ref(null)

const columns = [
  { title: '礼包名称', slotName: 'giftName', width: 160 },
  { title: '会员', slotName: 'memberName', width: 150 },
  { title: '发放时间', slotName: 'grantTime', width: 180 },
  { title: '领取时间', slotName: 'claimTime', width: 180 },
  { title: '发放方式', dataIndex: 'grantType', width: 120 },
  { title: '状态', slotName: 'status', width: 100 }
]

const getStatusColor = (s) => ({ claimed: 'green', pending: 'orange', expired: 'gray' }[s] || 'gray')
const getStatusName = (s) => ({ claimed: '已领取', pending: '未领取', expired: '已过期' }[s] || s)

const loadData = async () => {
  loading.value = true
  try {
    const params = {}
    if (filterGiftId.value) params.giftId = filterGiftId.value
    if (filterStatus.value) params.status = filterStatus.value
    const res = await api.getGiftRecords(params)
    recordList.value = res.data?.list || []
    stats.value = res.data?.stats || {}
  } catch (err) { Message.error('加载数据失败: ' + err.message) } finally { loading.value = false }
}

const loadGifts = async () => {
  try {
    const res = await api.getMemberGiftList()
    giftList.value = res.data?.list || []
  } catch (err) { /* ignore */ }
}

const exportRecords = () => { Message.info('正在导出...') }

onMounted(async () => {
  await loadGifts()
  await loadData()
})
</script>

<style scoped>
.gift-records-page { padding: 20px 24px; min-height: calc(100vh - 64px); background: #f5f7fa; }
.breadcrumb { margin-bottom: 16px; }
.action-card { margin-bottom: 16px; }
</style>
