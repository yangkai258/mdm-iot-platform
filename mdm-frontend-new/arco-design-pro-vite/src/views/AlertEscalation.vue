<template>
  <div class="container">
    <a-card>
      <template #title>
        <a-space><icon-up-circle /> 告警升级规则配置</a-space>
      </template>

      <a-tabs default-active-key="rules">
        <a-tab-pane key="rules" title="升级规则">
          <a-space style="margin-bottom: 16px">
            <a-button type="primary" @click="handleCreate">
              <template #icon><icon-plus /></template>
              新建规则
            </a-button>
            <a-button @click="handleBatchEnable">批量启用</a-button>
            <a-button @click="handleBatchDisable">批量禁用</a-button>
          </a-space>

          <a-table :columns="columns" :data="tableData" :row-selection="{ type: 'checkbox' }">
            <template #enabled="{ record }">
              <a-switch v-model="record.enabled" @change="handleToggle(record)" />
            </template>
            <template #condition="{ record }">
              {{ record.conditionType === 'time' ? `持续${record.timeThreshold}分钟` : `累计${record.countThreshold}次` }}
            </template>
            <template #actions="{ record }">
              <a-link @click="handleEdit(record)">编辑</a-link>
              <a-link @click="handleTest(record)">测试</a-link>
              <a-link @click="handleDelete(record)">删除</a-link>
            </template>
          </a-table>
        </a-tab-pane>

        <a-tab-pane key="logs" title="升级记录">
          <a-table :columns="logColumns" :data="logData" :pagination="pagination" />
        </a-tab-pane>
      </a-tabs>
    </a-card>

    <a-modal v-model:visible="modalVisible" :title="modalTitle" @ok="handleSubmit" width="600px">
      <a-form :model="formData" layout="vertical">
        <a-form-item label="规则名称" required>
          <a-input v-model="formData.ruleName" placeholder="请输入规则名称" />
        </a-form-item>
        <a-form-item label="告警类型">
          <a-select v-model="formData.alertType" placeholder="选择告警类型">
            <a-option value="temperature">温度异常</a-option>
            <a-option value="battery">电量不足</a-option>
            <a-option value="offline">设备离线</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="升级条件">
          <a-radio-group v-model="formData.conditionType">
            <a-radio value="time">时间触发</a-radio>
            <a-radio value="count">次数触发</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item v-if="formData.conditionType === 'time'" label="时间阈值(分钟)">
          <a-input-number v-model="formData.timeThreshold" :min="1" />
        </a-form-item>
        <a-form-item v-if="formData.conditionType === 'count'" label="次数阈值">
          <a-input-number v-model="formData.countThreshold" :min="1" />
        </a-form-item>
        <a-form-item label="升级目标级别">
          <a-select v-model="formData.escalationLevel">
            <a-option value="warning">警告</a-option>
            <a-option value="major">重要</a-option>
            <a-option value="critical">紧急</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="通知渠道">
          <a-checkbox-group v-model="formData.channels">
            <a-checkbox value="email">邮件</a-checkbox>
            <a-checkbox value="sms">短信</a-checkbox>
            <a-checkbox value="push">推送</a-checkbox>
            <a-checkbox value="webhook">Webhook</a-checkbox>
          </a-checkbox-group>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'

const modalVisible = ref(false)
const modalTitle = ref('新建规则')
const formData = reactive({
  ruleName: '', alertType: '', conditionType: 'time',
  timeThreshold: 30, countThreshold: 3, escalationLevel: 'major', channels: ['email']
})
const pagination = reactive({ current: 1, pageSize: 10, total: 30 })
const columns = [
  { title: '规则名称', dataIndex: 'ruleName' },
  { title: '告警类型', dataIndex: 'alertType' },
  { title: '升级条件', slotName: 'condition' },
  { title: '目标级别', dataIndex: 'escalationLevel' },
  { title: '启用状态', slotName: 'enabled' },
  { title: '操作', slotName: 'actions', width: 180 }
]
const tableData = ref([
  { id: 1, ruleName: '温度持续异常', alertType: 'temperature', conditionType: 'time', timeThreshold: 30, escalationLevel: 'major', enabled: true },
  { id: 2, ruleName: '频繁离线', alertType: 'offline', conditionType: 'count', countThreshold: 5, escalationLevel: 'critical', enabled: false }
])
const logColumns = [
  { title: '时间', dataIndex: 'time' },
  { title: '原告警', dataIndex: 'originalAlert' },
  { title: '升级后', dataIndex: 'escalatedTo' },
  { title: '通知渠道', dataIndex: 'channels' },
  { title: '状态', dataIndex: 'status' }
]
const logData = ref([
  { time: '2026-03-28 10:00', originalAlert: '温度异常-warning', escalatedTo: '温度异常-major', channels: '邮件,短信', status: '已发送' }
])

const handleCreate = () => { modalVisible.value = true; modalTitle.value = '新建规则' }
const handleEdit = (record) => { modalVisible.value = true; modalTitle.value = '编辑规则' }
const handleTest = (record) => { }
const handleDelete = (record) => { }
const handleToggle = (record) => { }
const handleBatchEnable = () => { }
const handleBatchDisable = () => { }
const handleSubmit = () => { modalVisible.value = false }
</script>

<style scoped>
.container { padding: 16px; }
</style>
