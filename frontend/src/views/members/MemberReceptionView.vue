<template>
  <div class="member-page">
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>会员管理</a-breadcrumb-item>
      <a-breadcrumb-item>会员接待</a-breadcrumb-item>
    </a-breadcrumb>

    <a-row :gutter="16" class="stats-row">
      <a-col :span="8">
        <a-card class="stat-card">
          <a-statistic title="今日接待" :value="stats.today || 0" />
        </a-card>
      </a-col>
      <a-col :span="8">
        <a-card class="stat-card">
          <a-statistic title="本周接待" :value="stats.week || 0" />
        </a-card>
      </a-col>
      <a-col :span="8">
        <a-card class="stat-card">
          <a-statistic title="本月接待" :value="stats.month || 0" />
        </a-card>
      </a-col>
    </a-row>

    <a-card class="action-card">
      <a-space wrap>
        <a-input-search v-model="filters.keyword" placeholder="搜索会员名称/接待人" style="width: 240px" search-button @search="loadData" />
        <a-range-picker v-model="filters.dateRange" style="width: 260px;" @change="loadData" />
        <a-button type="primary" @click="loadData">搜索</a-button>
        <a-button @click="loadData">刷新</a-button>
      </a-space>
    </a-card>

    <a-card class="table-card">
      <a-table :columns="columns" :data="dataList" :loading="loading" :pagination="paginationConfig" @page-change="onPageChange" row-key="id" :scroll="{ x: 1100 }">
        <template #type="{ record }">
          <a-tag :color="typeColor(record.type)">{{ record.type }}</a-tag>
        </template>
        <template #actions="{ record }">
          <a-button type="text" size="small" @click="showDetail(record)">查看详情</a-button>
        </template>
      </a-table>
    </a-card>

    <a-drawer v-model:visible="detailVisible" title="接待详情" :width="500" :footer="false">
      <a-descriptions :column="1" bordered>
        <a-descriptions-item label="会员名称">{{ detailRecord.memberName }}</a-descriptions-item>
        <a-descriptions-item label="接待时间">{{ detailRecord.receptionTime }}</a-descriptions-item>
        <a-descriptions-item label="接待人">{{ detailRecord.receptionist }}</a-descriptions-item>
        <a-descriptions-item label="接待类型">{{ detailRecord.type }}</a-descriptions-item>
        <a-descriptions-item label="接待内容">{{ detailRecord.content }}</a-descriptions-item>
        <a-descriptions-item label="备注">{{ detailRecord.remark }}</a-descriptions-item>
      </a-descriptions>
    </a-drawer>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { Message } from '@arco-design/web-vue'

const loading = ref(false)
const detailVisible = ref(false)
const detailRecord = ref({})
const dataList = ref([])
const stats = ref({ today: 5, week: 32, month: 128 })
const filters = reactive({ keyword: '', dateRange: [] })
const paginationConfig = reactive({ current: 1, pageSize: 10, total: 0 })

const columns = [
  { title: '接待记录', dataIndex: 'recordNo', width: 180 },
  { title: '会员名称', dataIndex: 'memberName', width: 150 },
  { title: '接待时间', dataIndex: 'receptionTime', width: 180 },
  { title: '接待人', dataIndex: 'receptionist', width: 120 },
  { title: '接待类型', slotName: 'type', width: 120 },
  { title: '操作', slotName: 'actions', width: 120, fixed: 'right' }
]

const typeColor = (t) => ({ '到店接待': 'blue', '电话回访': 'green', '线上咨询': 'purple', '活动接待': 'orange' }[t] || 'gray')

const mockData = () => [
  { id: 1, recordNo: 'RC20260324001', memberName: '张三', receptionTime: '2026-03-24 10:00:00', receptionist: '客服小李', type: '到店接待', content: 'VIP会员到店咨询，介绍新品', remark: '意向强烈' },
  { id: 2, recordNo: 'RC20260324002', memberName: '李四', receptionTime: '2026-03-24 11:30:00', receptionist: '客服小王', type: '电话回访', content: '生日会员电话回访', remark: '满意' },
  { id: 3, recordNo: 'RC20260324003', memberName: '王五', receptionTime: '2026-03-24 14:00:00', receptionist: '客服小张', type: '线上咨询', content: '线上咨询会员权益', remark: '-' }
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
const showDetail = (record) => { detailRecord.value = record; detailVisible.value = true }

loadData()
</script>

<style scoped>
.member-page { padding: 20px; }
.breadcrumb { margin-bottom: 16px; }
.stats-row { margin-bottom: 16px; }
.stat-card { text-align: center; }
.action-card { margin-bottom: 16px; }
</style>
