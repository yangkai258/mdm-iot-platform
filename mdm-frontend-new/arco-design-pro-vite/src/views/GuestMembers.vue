<template>
  <div class="container">
    <a-card>
      <template #title><icon-user-guest /> 临时会员/访客管理</template>
      <template #extra>
        <a-space>
          <a-button @click="handleExport"><template #icon><icon-download /></template>导出</a-button>
          <a-button type="primary" @click="handleConvert"><template #icon><icon-swap /></template>转为正式会员</a-button>
        </a-space>
      </template>
      <a-table :columns="columns" :data="guests" :pagination="pagination" :row-selection="{ type: 'checkbox' }">
        <template #status="{ record }">
          <a-tag :color="getStatusColor(record.status)">{{ getStatusText(record.status) }}</a-tag>
        </template>
        <template #type="{ record }">
          <a-tag :color="record.type === 'visitor' ? 'blue' : 'green'">{{ record.type === 'visitor' ? '访客' : '临时' }}</a-tag>
        </template>
        <template #actions="{ record }">
          <a-link @click="handleView(record)">查看</a-link>
          <a-divider direction="vertical" />
          <a-link @click="handleConvertOne(record)">转为正式</a-link>
          <a-divider direction="vertical" />
          <a-link status="error" @click="handleDelete(record)">删除</a-link>
        </template>
      </a-table>
    </a-card>
    <a-modal v-model:visible="detailVisible" title="访客详情" :footer="null">
      <a-descriptions :column="2" bordered>
        <a-descriptions-item label="访客ID">{{ currentGuest?.id }}</a-descriptions-item>
        <a-descriptions-item label="访客昵称">{{ currentGuest?.nickname }}</a-descriptions-item>
        <a-descriptions-item label="来访时间">{{ currentGuest?.visitTime }}</a-descriptions-item>
        <a-descriptions-item label="来访次数">{{ currentGuest?.visitCount }}</a-descriptions-item>
        <a-descriptions-item label="状态" :span="2">
          <a-tag :color="getStatusColor(currentGuest?.status)">{{ getStatusText(currentGuest?.status) }}</a-tag>
        </a-descriptions-item>
      </a-descriptions>
      <a-divider>行为记录</a-divider>
      <a-timeline>
        <a-timeline-item v-for="(log, idx) in behaviorLogs" :key="idx">
          <p>{{ log.action }}</p>
          <span class="time">{{ log.time }}</span>
        </a-timeline-item>
      </a-timeline>
    </a-modal>
  </div>
</template>
<script setup>
import { ref, reactive } from 'vue'
const columns = [
  { title: '访客ID', dataIndex: 'id', width: 100 },
  { title: '访客昵称', dataIndex: 'nickname' },
  { title: '类型', slotName: 'type', width: 80 },
  { title: '来访时间', dataIndex: 'visitTime' },
  { title: '来访次数', dataIndex: 'visitCount', width: 100 },
  { title: '最后活动', dataIndex: 'lastActivity' },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '操作', slotName: 'actions', width: 180 }
]
const pagination = { pageSize: 10 }
const guests = ref([
  { id: 'G001', nickname: '路过看看', type: 'visitor', visitTime: '2026-03-28 10:00', visitCount: 3, lastActivity: '2026-03-28 14:30', status: 'active' },
  { id: 'G002', nickname: '临时用户', type: 'temp', visitTime: '2026-03-27 09:00', visitCount: 1, lastActivity: '2026-03-27 09:30', status: 'inactive' },
  { id: 'G003', nickname: '游客123', type: 'visitor', visitTime: '2026-03-26 15:00', visitCount: 5, lastActivity: '2026-03-28 08:00', status: 'active' }
])
const detailVisible = ref(false)
const currentGuest = ref(null)
const behaviorLogs = ref([
  { action: '浏览了首页', time: '2026-03-28 14:30' },
  { action: '浏览了宠物列表', time: '2026-03-28 14:00' },
  { action: '首次访问', time: '2026-03-28 10:00' }
])
const getStatusColor = (s) => ({ active: 'green', inactive: 'gray', converted: 'blue' }[s] || 'gray')
const getStatusText = (s) => ({ active: '活跃', inactive: '沉默', converted: '已转化' }[s] || s)
const handleView = (r) => { currentGuest.value = r; detailVisible.value = true }
const handleConvert = () => {}
const handleConvertOne = (r) => {}
const handleDelete = (r) => {}
const handleExport = () => {}
</script>
<style scoped>
.container { padding: 16px; }
.time { font-size: 12px; color: #909399; }
</style>
