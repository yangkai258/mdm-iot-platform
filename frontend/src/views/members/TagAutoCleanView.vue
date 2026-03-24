<template>
  <div class="member-page">
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>会员管理</a-breadcrumb-item>
      <a-breadcrumb-item>标签清除设置</a-breadcrumb-item>
    </a-breadcrumb>

    <a-card class="action-card">
      <a-form :model="form" layout="vertical" style="max-width: 600px;">
        <a-divider>自动清除规则</a-divider>

        <a-form-item label="自动清除周期">
          <a-select v-model="form.cleanCycle" style="width: 300px;">
            <a-option value="daily">每日清除</a-option>
            <a-option value="weekly">每周清除</a-option>
            <a-option value="monthly">每月清除</a-option>
            <a-option value="manual">手动清除</a-option>
          </a-select>
          <div style="color:#999;font-size:12px;margin-top:4px;">按设定周期自动清除过期标签</div>
        </a-form-item>

        <a-form-item label="不活跃时间阈值（天）">
          <a-input-number v-model="form.inactiveDays" :min="1" :max="365" style="width: 200px;" />
          <div style="color:#999;font-size:12px;margin-top:4px;">会员超过此天数未活跃，标签将被自动清除</div>
        </a-form-item>

        <a-form-item label="清除前通知">
          <a-switch v-model="form.notifyBefore" checked-value="1" unchecked-value="0" />
          <span style="margin-left:8px;color:#999;">清除前3天发送通知提醒会员</span>
        </a-form-item>

        <a-form-item label="允许手动恢复">
          <a-switch v-model="form.allowRestore" checked-value="1" unchecked-value="0" />
          <span style="margin-left:8px;color:#999;">标签清除后会员可手动申请恢复</span>
        </a-form-item>

        <a-divider>保护标签</a-divider>

        <a-form-item label="受保护的标签（不会被自动清除）">
          <a-select v-model="form.protectedTags" multiple placeholder="选择受保护的标签" style="width: 100%;">
            <a-option value="vip">VIP标签</a-option>
            <a-option value="birthday">生日标签</a-option>
            <a-option value="anniversary">周年标签</a-option>
          </a-select>
        </a-form-item>

        <a-divider>清除日志</a-divider>

        <a-form-item label="清除日志保留天数">
          <a-input-number v-model="form.logRetentionDays" :min="7" :max="365" style="width: 200px;" />
          <span style="margin-left:8px;color:#999;">天</span>
        </a-form-item>

        <a-form-item>
          <a-space>
            <a-button type="primary" @click="handleSave">保存设置</a-button>
            <a-button @click="handleReset">重置</a-button>
          </a-space>
        </a-form-item>
      </a-form>
    </a-card>

    <a-card class="table-card" style="margin-top: 16px;">
      <template #title><span style="font-weight:600;">清除记录</span></template>
      <a-table :columns="logColumns" :data="logList" :loading="logLoading" row-key="id" :scroll="{ x: 700 }">
        <template #status="{ record }">
          <a-tag :color="record.status === 1 ? 'green' : 'red'">{{ record.status === 1 ? '成功' : '失败' }}</a-tag>
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { Message } from '@arco-design/web-vue'

const form = reactive({
  cleanCycle: 'monthly',
  inactiveDays: 90,
  notifyBefore: '1',
  allowRestore: '1',
  protectedTags: ['vip', 'birthday'],
  logRetentionDays: 90
})

const logLoading = ref(false)
const logList = ref([
  { id: 1, time: '2026-03-01 00:00:00', tagName: '流失风险', cleanedCount: 234, status: 1 },
  { id: 2, time: '2026-02-01 00:00:00', tagName: '沉默会员', cleanedCount: 567, status: 1 },
  { id: 3, time: '2026-01-01 00:00:00', tagName: '边缘客户', cleanedCount: 123, status: 1 }
])

const logColumns = [
  { title: '清除时间', dataIndex: 'time', width: 200 },
  { title: '清除标签', dataIndex: 'tagName', width: 200 },
  { title: '清除数量', dataIndex: 'cleanedCount', width: 150 },
  { title: '状态', slotName: 'status', width: 150 }
]

const handleSave = () => {
  Message.success('保存成功')
}

const handleReset = () => {
  Object.assign(form, {
    cleanCycle: 'monthly',
    inactiveDays: 90,
    notifyBefore: '1',
    allowRestore: '1',
    protectedTags: ['vip', 'birthday'],
    logRetentionDays: 90
  })
  Message.info('已重置')
}
</script>

<style scoped>
.member-page { padding: 20px; }
.breadcrumb { margin-bottom: 16px; }
.action-card { margin-bottom: 16px; }
</style>
