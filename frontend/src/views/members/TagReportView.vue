<template>
  <div class="member-page">
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>会员管理</a-breadcrumb-item>
      <a-breadcrumb-item>会员标签报表</a-breadcrumb-item>
    </a-breadcrumb>

    <a-card class="action-card">
      <a-space wrap>
        <a-select v-model="filters.tagType" placeholder="标签类型" allow-clear style="width: 160px" @change="loadData">
          <a-option value="highfreq">高频购买</a-option>
          <a-option value="lowfreq">低频购买</a-option>
          <a-option value="interest">兴趣分类</a-option>
        </a-select>
        <a-range-picker v-model="filters.dateRange" style="width: 260px;" @change="loadData" />
        <a-button type="primary" @click="loadData">搜索</a-button>
        <a-button @click="handleExport">导出</a-button>
        <a-button @click="loadData">刷新</a-button>
      </a-space>
    </a-card>

    <a-row :gutter="16" class="stats-row">
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="标签总数" :value="stats.total || 0" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="本月新增" :value="stats.monthNew || 0" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="本月清除" :value="stats.monthClean || 0" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="活跃标签" :value="stats.active || 0" />
        </a-card>
      </a-col>
    </a-row>

    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="12">
        <a-card>
          <template #title><span style="font-weight:600;">标签分布</span></template>
          <div ref="pieChartRef" style="height: 280px;"></div>
        </a-card>
      </a-col>
      <a-col :span="12">
        <a-card>
          <template #title><span style="font-weight:600;">标签趋势（近6月）</span></template>
          <div ref="lineChartRef" style="height: 280px;"></div>
        </a-card>
      </a-col>
    </a-row>

    <a-card class="table-card">
      <a-table :columns="columns" :data="dataList" :loading="loading" :pagination="paginationConfig" @page-change="onPageChange" row-key="id" :scroll="{ x: 1000 }">
        <template #trend="{ record }">
          <span :style="{ color: record.trend > 0 ? '#52c41a' : record.trend < 0 ? '#ff4d4f' : '#999' }">
            {{ record.trend > 0 ? '↑' : record.trend < 0 ? '↓' : '-' }}{{ Math.abs(record.trend) }}
          </span>
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const loading = ref(false)
const dataList = ref([])
const stats = ref({ total: 12, monthNew: 3, monthClean: 1, active: 11 })
const filters = reactive({ tagType: '', dateRange: [] })
const paginationConfig = reactive({ current: 1, pageSize: 10, total: 0 })
const pieChartRef = ref(null)
const lineChartRef = ref(null)

const columns = [
  { title: '标签名称', dataIndex: 'name', width: 180 },
  { title: '标签类型', dataIndex: 'type', width: 150 },
  { title: '当前会员数', dataIndex: 'memberCount', width: 130 },
  { title: '上月会员数', dataIndex: 'lastMonthCount', width: 130 },
  { title: '环比变化', slotName: 'trend', width: 120 },
  { title: '创建时间', dataIndex: 'createTime', width: 180 },
  { title: '状态', dataIndex: 'status', width: 100 }
]

const mockData = () => [
  { id: 1, name: '月度活跃买家', type: '高频购买', memberCount: 1234, lastMonthCount: 1100, trend: 134, createTime: '2026-01-01', status: '启用' },
  { id: 2, name: '沉睡会员', type: '低频购买', memberCount: 3456, lastMonthCount: 3200, trend: 256, createTime: '2026-01-01', status: '启用' },
  { id: 3, name: '美食爱好者', type: '兴趣分类', memberCount: 2345, lastMonthCount: 2100, trend: 245, createTime: '2026-02-01', status: '启用' }
]

const loadData = () => {
  loading.value = true
  setTimeout(() => {
    dataList.value = mockData()
    paginationConfig.total = dataList.value.length
    loading.value = false
    renderCharts()
  }, 400)
}

const renderCharts = () => {
  // simple visual placeholder - in production would usearco charts
}

const onPageChange = (page) => { paginationConfig.current = page; loadData() }

const handleExport = () => {
  Message.success('导出成功')
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.member-page { padding: 20px; }
.breadcrumb { margin-bottom: 16px; }
.stats-row { margin-bottom: 16px; }
.stat-card { text-align: center; }
.action-card { margin-bottom: 16px; }
</style>
