<template>
  <div class="container">
    <a-card>
      <template #title>
        <a-space><icon-clock /> 定时发送配置</a-space>
      </template>
      <template #extra>
        <a-button type="primary" @click="handleCreate">
          <template #icon><icon-plus /></template>
          新建定时发送
        </a-button>
      </template>

      <a-table :columns="columns" :data="tableData">
        <template #status="{ record }">
          <a-tag :color="getStatusColor(record.status)">{{ getStatusLabel(record.status) }}</a-tag>
        </template>
        <template #scheduledTime="{ record }">
          <span>{{ record.scheduledAt }} {{ record.scheduledTime }}</span>
        </template>
        <template #actions="{ record }">
          <a-link v-if="record.status === 'pending'" @click="handleEdit(record)">编辑</a-link>
          <a-link v-if="record.status === 'pending'" @click="handlePause(record)">暂停</a-link>
          <a-link @click="handleView(record)">详情</a-link>
          <a-link @click="handleDelete(record)">删除</a-link>
        </template>
      </a-table>
    </a-card>

    <a-drawer v-model:visible="drawerVisible" :title="drawerTitle" :width="500" @ok="handleSubmit">
      <a-form :model="formData" layout="vertical">
        <a-form-item label="消息模板" required>
          <a-select v-model="formData.templateId" placeholder="选择消息模板">
            <a-option value="t1">系统通知</a-option>
            <a-option value="t2">活动提醒</a-option>
            <a-option value="t3">到期提醒</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="接收者类型">
          <a-radio-group v-model="formData.recipientType">
            <a-radio value="user">指定用户</a-radio>
            <a-radio value="device">指定设备</a-radio>
            <a-radio value="group">用户组</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item label="发送时间" required>
          <a-space direction="vertical">
            <a-date-picker v-model="formData.scheduledAt" />
            <a-time-picker v-model="formData.scheduledTime" format="HH:mm" />
          </a-space>
        </a-form-item>
        <a-form-item label="是否循环">
          <a-switch v-model="formData.recurring" />
        </a-form-item>
        <a-form-item v-if="formData.recurring" label="循环周期">
          <a-select v-model="formData.recurringType">
            <a-option value="daily">每天</a-option>
            <a-option value="weekly">每周</a-option>
            <a-option value="monthly">每月</a-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-drawer>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'

const drawerVisible = ref(false)
const drawerTitle = ref('新建定时发送')
const formData = reactive({
  templateId: '', recipientType: 'user', scheduledAt: '', scheduledTime: '',
  recurring: false, recurringType: 'daily'
})

const columns = [
  { title: '模板名称', dataIndex: 'templateName' },
  { title: '接收者', dataIndex: 'recipient' },
  { title: '发送时间', slotName: 'scheduledTime' },
  { title: '状态', slotName: 'status' },
  { title: '创建时间', dataIndex: 'createdAt' },
  { title: '操作', slotName: 'actions', width: 200 }
]

const tableData = ref([
  { id: 1, templateName: '系统升级通知', recipient: '全部用户', scheduledAt: '2026-03-29', scheduledTime: '10:00', status: 'pending', createdAt: '2026-03-28' },
  { id: 2, templateName: '活动提醒', recipient: 'VIP用户', scheduledAt: '2026-03-30', scheduledTime: '09:00', status: 'pending', createdAt: '2026-03-28' }
])

const getStatusColor = (s) => ({ pending: 'blue', sent: 'green', failed: 'red', cancelled: 'gray' }[s] || 'gray')
const getStatusLabel = (s) => ({ pending: '待发送', sent: '已发送', failed: '发送失败', cancelled: '已取消' }[s] || s)

const handleCreate = () => { drawerVisible.value = true; drawerTitle.value = '新建定时发送' }
const handleEdit = (r) => { drawerVisible.value = true; drawerTitle.value = '编辑定时发送' }
const handlePause = (r) => { }
const handleView = (r) => { }
const handleDelete = (r) => { }
const handleSubmit = () => { drawerVisible.value = false }
</script>

<style scoped>
.container { padding: 16px; }
</style>
