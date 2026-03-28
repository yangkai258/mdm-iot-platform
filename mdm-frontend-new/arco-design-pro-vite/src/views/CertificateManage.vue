<template>
  <div class="container">
    <a-card>
      <template #title>
        <a-space><icon-certificate /> 证书管理</a-space>
      </template>
      <template #extra>
        <a-button type="primary" @click="handleCreate">
          <template #icon><icon-plus /></template>
          创建证书
        </a-button>
      </template>

      <a-row :gutter="16" style="margin-bottom: 16px">
        <a-col :span="6">
          <a-card>
            <a-statistic title="有效证书" :value="stats.valid" />
          </a-card>
        </a-col>
        <a-col :span="6">
          <a-card>
            <a-statistic title="即将过期" :value="stats.expiring" :value-style="{ color: '#E6A23C' }" />
          </a-card>
        </a-col>
        <a-col :span="6">
          <a-card>
            <a-statistic title="已过期" :value="stats.expired" :value-style="{ color: '#F56C6C' }" />
          </a-card>
        </col>
      </a-row>

      <a-table :columns="columns" :data="certificates">
        <template #status="{ record }">
          <a-tag :color="getStatusColor(record.status)">{{ getStatusLabel(record.status) }}</a-tag>
        </template>
        <template #expiry="{ record }">
          <span :class="{ expired: record.status === 'expired' }">{{ record.expiresAt }}</span>
        </template>
        <template #actions="{ record }">
          <a-link @click="handleView(record)">详情</a-link>
          <a-link @click="handleRenew(record)">续期</a-link>
          <a-link @click="handleRevoke(record)">吊销</a-link>
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'

const stats = reactive({ valid: 45, expiring: 5, expired: 2 })

const columns = [
  { title: '证书ID', dataIndex: 'id' },
  { title: '证书名称', dataIndex: 'name' },
  { title: '设备', dataIndex: 'deviceId' },
  { title: '状态', slotName: 'status', width: 120 },
  { title: '过期时间', slotName: 'expiry', width: 180 },
  { title: '操作', slotName: 'actions', width: 180 }
]

const certificates = ref([
  { id: 'C001', name: '设备证书', deviceId: 'D001', status: 'valid', expiresAt: '2027-03-28' }
])

const getStatusColor = (s) => ({ valid: 'green', expiring: 'orange', expired: 'red' }[s] || 'gray')
const getStatusLabel = (s) => ({ valid: '有效', expiring: '即将过期', expired: '已过期' }[s] || s)

const handleCreate = () => { }
const handleView = (r) => { }
const handleRenew = (r) => { }
const handleRevoke = (r) => { }
</script>

<style scoped>
.container { padding: 16px; }
.expired { color: #F56C6C; }
</style>
