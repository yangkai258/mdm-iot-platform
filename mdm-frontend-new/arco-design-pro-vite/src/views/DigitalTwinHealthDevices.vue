<template>
  <div class="container">
    <a-card>
      <template #title>
        <a-space><icon-link /> 第三方健康设备绑定</a-space>
      </template>
      <template #extra>
        <a-button type="primary" @click="handleBind">
          <template #icon><icon-plus /></template>
          绑定新设备
        </a-button>
      </template>

      <a-row :gutter="16" style="margin-bottom: 16px">
        <a-col :span="6">
          <a-card size="small">
            <a-statistic title="已绑定设备" :value="stats.binded" />
          </a-card>
        </a-col>
        <a-col :span="6">
          <a-card size="small">
            <a-statistic title="同步中" :value="stats.syncing" :value-style="{ color: '#409EFF' }" />
          </a-card>
        </a-col>
        <a-col :span="6">
          <a-card size="small">
            <a-statistic title="最后同步" :value="stats.lastSync" />
          </a-card>
        </a-col>
      </a-row>

      <a-card title="已绑定设备">
        <a-table :columns="columns" :data="devices">
          <template #type="{ record }">
            <a-tag>{{ getDeviceTypeLabel(record.type) }}</a-tag>
          </template>
          <template #syncStatus="{ record }">
            <a-badge :status="getSyncStatusBadge(record.syncStatus)" :text="getSyncStatusText(record.syncStatus)" />
          </template>
          <template #actions="{ record }">
            <a-link @click="handleSync(record)">同步</a-link>
            <a-link @click="handleSettings(record)">设置</a-link>
            <a-link @click="handleUnbind(record)">解绑</a-link>
          </template>
        </a-table>
      </a-card>

      <a-card title="同步历史" style="margin-top: 16px">
        <a-table :columns="syncColumns" :data="syncHistory" size="small" :pagination="pagination" />
      </a-card>
    </a-card>

    <a-drawer v-model:visible="bindVisible" title="绑定设备" :width="400" @ok="handleBindConfirm">
      <a-form :model="bindForm" layout="vertical">
        <a-form-item label="设备类型" required>
          <a-select v-model="bindForm.type" placeholder="选择设备类型">
            <a-option value="smart_collar">智能项圈</a-option>
            <a-option value="weighing_scale">智能体重秤</a-option>
            <a-option value="blood_monitor">血压计</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="设备序列号" required>
          <a-input v-model="bindForm.deviceId" placeholder="请输入设备序列号" />
        </a-form-item>
        <a-form-item label="同步频率">
          <a-select v-model="bindForm.syncInterval">
            <a-option value="5min">每5分钟</a-option>
            <a-option value="15min">每15分钟</a-option>
            <a-option value="30min">每30分钟</a-option>
            <a-option value="1hour">每小时</a-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-drawer>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'

const stats = reactive({ binded: 3, syncing: 1, lastSync: '5分钟前' })
const bindVisible = ref(false)
const bindForm = reactive({ type: '', deviceId: '', syncInterval: '15min' })
const pagination = reactive({ current: 1, pageSize: 10, total: 20 })

const columns = [
  { title: '设备名称', dataIndex: 'name' },
  { title: '类型', slotName: 'type', width: 120 },
  { title: '序列号', dataIndex: 'deviceId' },
  { title: '同步状态', slotName: 'syncStatus', width: 120 },
  { title: '最后同步', dataIndex: 'lastSync' },
  { title: '操作', slotName: 'actions', width: 180 }
]

const devices = ref([
  { id: 1, name: '智能项圈-Pro', type: 'smart_collar', deviceId: 'SC-001-ABC', syncStatus: 'synced', lastSync: '2026-03-28 10:00' },
  { id: 2, name: '智能体重秤', type: 'weighing_scale', deviceId: 'WS-002-DEF', syncStatus: 'syncing', lastSync: '2026-03-28 09:55' }
])

const syncColumns = [
  { title: '时间', dataIndex: 'time' },
  { title: '设备', dataIndex: 'deviceName' },
  { title: '数据类型', dataIndex: 'dataType' },
  { title: '状态', dataIndex: 'status' }
]
const syncHistory = ref([
  { time: '2026-03-28 10:00', deviceName: '智能项圈-Pro', dataType: '心率/体温', status: '成功' }
])

const getDeviceTypeLabel = (t) => ({ smart_collar: '智能项圈', weighing_scale: '体重秤', blood_monitor: '血压计' }[t] || t)
const getSyncStatusBadge = (s) => ({ synced: 'success', syncing: 'processing', failed: 'error' }[s] || 'default')
const getSyncStatusText = (s) => ({ synced: '已同步', syncing: '同步中', failed: '失败' }[s] || s)

const handleBind = () => { bindVisible.value = true }
const handleSync = (r) => { }
const handleSettings = (r) => { }
const handleUnbind = (r) => { }
const handleBindConfirm = () => { bindVisible.value = false }
</script>

<style scoped>
.container { padding: 16px; }
</style>
