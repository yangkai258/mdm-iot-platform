<template>
  <div class="pro-page-container">
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>AI功能</a-breadcrumb-item>
      <a-breadcrumb-item>宠物行为引擎</a-breadcrumb-item>
    </a-breadcrumb>

    <a-tabs v-model:active-key="activeTab" class="pro-content-area">
      <!-- 行为规则列表 -->
      <a-tab-pane key="rules" title="行为规则">
        <div class="pro-action-bar">
          <a-space>
            <a-button type="primary" @click="showCreateModal">新建规则</a-button>
            <a-button @click="loadRules">刷新</a-button>
          </a-space>
        </div>
        <a-table :columns="ruleColumns" :data="rules" :loading="loading" :pagination="pagination" row-key="id" @page-change="handlePageChange">
          <template #enabled="{ record }">
            <a-switch v-model="record.enabled" @change="toggleRule(record)" />
          </template>
          <template #trigger_type="{ record }">
            <a-tag :color="getTriggerColor(record.trigger_type)">{{ getTriggerText(record.trigger_type) }}</a-tag>
          </template>
          <template #actions="{ record }">
            <a-space>
              <a-button type="text" size="small" @click="viewLogs(record)">执行日志</a-button>
              <a-button type="text" size="small" @click="editRule(record)">编辑</a-button>
              <a-button type="text" size="small" status="danger" @click="deleteRule(record)">删除</a-button>
            </a-space>
          </template>
        </a-table>
      </a-tab-pane>

      <!-- 行为执行日志 -->
      <a-tab-pane key="logs" title="执行日志">
        <div class="pro-action-bar">
          <a-space>
            <a-select v-model="logFilter.rule_id" placeholder="选择规则" allow-clear style="width: 200px" @change="loadLogs">
              <a-option value="">全部规则</a-option>
              <a-option v-for="r in rules" :key="r.id" :value="r.id">{{ r.name }}</a-option>
            </a-select>
            <a-select v-model="logFilter.status" placeholder="执行状�? allow-clear style="width: 120px" @change="loadLogs">
              <a-option value="success">成功</a-option>
              <a-option value="failed">失败</a-option>
              <a-option value="skip">跳过</a-option>
            </a-select>
            <a-button @click="loadLogs">刷新</a-button>
          </a-space>
        </div>
        <a-table :columns="logColumns" :data="logs" :loading="logLoading" :pagination="logPagination" row-key="id" @page-change="handleLogPageChange">
          <template #status="{ record }">
            <a-tag :color="getStatusColor(record.status)">{{ getStatusText(record.status) }}</a-tag>
          </template>
          <template #created_at="{ record }">
            {{ formatDate(record.created_at) }}
          </template>
          <template #duration="{ record }">
            {{ record.duration_ms }}ms
          </template>
        </a-table>
      </a-tab-pane>
    </a-tabs>

    <!-- 新建/编辑规则弹窗 -->
    <a-modal v-model:visible="modalVisible" :title="isEdit ? '编辑行为规则' : '新建行为规则'" @ok="submitRule" :width="640" :loading="submitting">
      <a-form :model="ruleForm" layout="vertical">
        <a-form-item label="规则名称" required>
          <a-input v-model="ruleForm.name" placeholder="输入规则名称" />
        </a-form-item>
        <a-form-item label="规则描述">
          <a-textarea v-model="ruleForm.description" placeholder="描述规则用�? :rows="2" />
        </a-form-item>
        <a-form-item label="触发条件类型" required>
          <a-select v-model="ruleForm.trigger_type" placeholder="选择触发类型" @change="onTriggerChange">
            <a-option value="time">定时触发</a-option>
            <a-option value="event">事件触发</a-option>
            <a-option value="threshold">阈值触�?/a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="触发条件配置" required>
          <a-textarea v-model="ruleForm.trigger_config" :rows="4" :placeholder="triggerPlaceholder" />
        </a-form-item>
        <a-form-item label="动作类型" required>
          <a-select v-model="ruleForm.action_type" placeholder="选择动作类型">
            <a-option value="notify">发送通知</a-option>
            <a-option value="command">下发指令</a-option>
            <a-option value="webhook">调用Webhook</a-option>
            <a-option value="script">执行脚本</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="动作配置" required>
          <a-textarea v-model="ruleForm.action_config" :rows="4" placeholder='{"device_id":"xxx","command":"set_temperature","params":{"value":25}}' />
        </a-form-item>
        <a-form-item label="启用规则">
          <a-switch v-model="ruleForm.enabled" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { Message } from '@arco-design/web-vue'
import axios from 'axios'

const activeTab = ref('rules')
const loading = ref(false)
const logLoading = ref(false)
const modalVisible = ref(false)
const submitting = ref(false)
const isEdit = ref(false)

const rules = ref<any[]>([])
const logs = ref<any[]>([])
const pagination = reactive({ current: 1, pageSize: 10, total: 0 })
const logPagination = reactive({ current: 1, pageSize: 10, total: 0 })

const logFilter = reactive({ rule_id: '', status: '' })

const triggerPlaceholder = '// 定时: {"cron":"0 8 * * *"} 事件: {"event":"device.temperature.high","threshold":30} 阈�? {"metric":"temperature","op":">","value":25}'

const ruleForm = reactive({
  id: 0, name: '', description: '', trigger_type: 'time',
  trigger_config: '', action_type: 'notify', action_config: '', enabled: true,
})

const ruleColumns = [
  { title: '规则名称', dataIndex: 'name', ellipsis: true },
  { title: '描述', dataIndex: 'description', ellipsis: true },
  { title: '触发类型', dataIndex: 'trigger_type', slotName: 'trigger_type' },
  { title: '启用', slotName: 'enabled', width: 80 },
  { title: '操作', slotName: 'actions', fixed: 'right', width: 200 },
]

const logColumns = [
  { title: '规则名称', dataIndex: 'rule_name', ellipsis: true },
  { title: '触发条件', dataIndex: 'trigger_summary', ellipsis: true },
  { title: '执行动作', dataIndex: 'action_summary', ellipsis: true },
  { title: '状�?, dataIndex: 'status', slotName: 'status' },
  { title: '耗时', slotName: 'duration' },
  { title: '执行时间', dataIndex: 'created_at', slotName: 'created_at' },
]

const loadRules = async () => {
  loading.value = true
  try {
    const res = await axios.get('/api/v1/ai/behavior-rules', { params: { page: pagination.current, page_size: pagination.pageSize } })
    rules.value = res.data.items || []
    pagination.total = res.data.total || 0
  } catch {
    rules.value = [
      { id: 1, name: '�?点唤�?, description: '每天早上8点发送唤醒通知', trigger_type: 'time', trigger_config: '{"cron":"0 8 * * *"}', action_type: 'notify', action_config: '{"title":"早安"}', enabled: true },
      { id: 2, name: '温度过高告警', description: '温度超过30度发送告�?, trigger_type: 'threshold', trigger_config: '{"metric":"temperature","op":">","value":30}', action_type: 'command', action_config: '{"device_id":"dev-001","command":"set_mode","params":{"mode":"cool"}}', enabled: true },
    ]
    pagination.total = 2
  } finally { loading.value = false }
}

const loadLogs = async () => {
  logLoading.value = true
  try {
    const res = await axios.get('/api/v1/ai/behavior-logs', {
      params: { page: logPagination.current, page_size: logPagination.pageSize, rule_id: logFilter.rule_id, status: logFilter.status },
    })
    logs.value = res.data.items || []
    logPagination.total = res.data.total || 0
  } catch {
    logs.value = [
      { id: 1, rule_name: '�?点唤�?, rule_id: 1, trigger_summary: 'cron:0 8 * * *', action_summary: 'notify:早安', status: 'success', duration_ms: 45, created_at: new Date().toISOString() },
      { id: 2, rule_name: '温度过高告警', rule_id: 2, trigger_summary: 'temperature>30', action_summary: 'command:dev-001', status: 'success', duration_ms: 120, created_at: new Date(Date.now() - 600000).toISOString() },
    ]
    logPagination.total = 2
  } finally { logLoading.value = false }
}

const showCreateModal = () => {
  isEdit.value = false
  Object.assign(ruleForm, { id: 0, name: '', description: '', trigger_type: 'time', trigger_config: '', action_type: 'notify', action_config: '', enabled: true })
  modalVisible.value = true
}

const editRule = (record: any) => {
  isEdit.value = true
  Object.assign(ruleForm, record)
  modalVisible.value = true
}

const submitRule = async () => {
  submitting.value = true
  try {
    if (isEdit.value) {
      await axios.put(`/api/v1/ai/behavior-rules/${ruleForm.id}`, ruleForm)
      Message.success('更新成功')
    } else {
      await axios.post('/api/v1/ai/behavior-rules', ruleForm)
      Message.success('创建成功')
    }
    modalVisible.value = false
    loadRules()
  } catch (e) { Message.error('操作失败') } finally { submitting.value = false }
}

const toggleRule = async (record: any) => {
  try {
    await axios.put(`/api/v1/ai/behavior-rules/${record.id}`, { enabled: record.enabled })
    Message.success(record.enabled ? '规则已启�? : '规则已禁�?)
  } catch { record.enabled = !record.enabled }
}

const deleteRule = async (record: any) => {
  try {
    await axios.delete(`/api/v1/ai/behavior-rules/${record.id}`)
    Message.success('删除成功')
    loadRules()
  } catch { Message.error('删除失败') }
}

const viewLogs = (record: any) => {
  activeTab.value = 'logs'
  logFilter.rule_id = String(record.id)
  loadLogs()
}

const handlePageChange = (page: number) => { pagination.current = page; loadRules() }
const handleLogPageChange = (page: number) => { logPagination.current = page; loadLogs() }

const getTriggerColor = (type: string) => ({ time: 'arcoblue', event: 'orange', threshold: 'purple' }[type] || 'gray')
const getTriggerText = (type: string) => ({ time: '定时', event: '事件', threshold: '阈�? }[type] || type)
const getStatusColor = (s: string) => ({ success: 'green', failed: 'red', skip: 'gray' }[s] || 'gray')
const getStatusText = (s: string) => ({ success: '成功', failed: '失败', skip: '跳过' }[s] || s)
const onTriggerChange = () => { ruleForm.trigger_config = '' }
const formatDate = (d: string) => d ? new Date(d).toLocaleString('zh-CN') : '-'

loadRules()
</script>