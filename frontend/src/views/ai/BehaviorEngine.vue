<template>
  <div class="pro-page-container">
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>棣栭〉</a-breadcrumb-item>
      <a-breadcrumb-item>AI鍔熻兘</a-breadcrumb-item>
      <a-breadcrumb-item>瀹犵墿琛屼负寮曟搸</a-breadcrumb-item>
    </a-breadcrumb>

    <a-tabs v-model:active-key="activeTab" class="pro-content-area">
      <!-- 琛屼负瑙勫垯鍒楄〃 -->
      <a-tab-pane key="rules" title="琛屼负瑙勫垯">
        <div class="pro-action-bar">
          <a-space>
            <a-button type="primary" @click="showCreateModal">鏂板缓瑙勫垯</a-button>
            <a-button @click="loadRules">鍒锋柊</a-button>
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
              <a-button type="text" size="small" @click="viewLogs(record)">鎵ц鏃ュ織</a-button>
              <a-button type="text" size="small" @click="editRule(record)">缂栬緫</a-button>
              <a-button type="text" size="small" status="danger" @click="deleteRule(record)">鍒犻櫎</a-button>
            </a-space>
          </template>
        </a-table>
      </a-tab-pane>

      <!-- 琛屼负鎵ц鏃ュ織 -->
      <a-tab-pane key="logs" title="鎵ц鏃ュ織">
        <div class="pro-action-bar">
          <a-space>
            <a-select v-model="logFilter.rule_id" placeholder="閫夋嫨瑙勫垯" allow-clear style="width: 200px" @change="loadLogs">
              <a-option value="">鍏ㄩ儴瑙勫垯</a-option>
              <a-option v-for="r in rules" :key="r.id" :value="r.id">{{ r.name }}</a-option>
            </a-select>
            <a-select v-model="logFilter.status" placeholder="鎵ц鐘讹拷" allow-clear style="width: 120px" @change="loadLogs">
              <a-option value="success">鎴愬姛</a-option>
              <a-option value="failed">澶辫触</a-option>
              <a-option value="skip">璺宠繃</a-option>
            </a-select>
            <a-button @click="loadLogs">鍒锋柊</a-button>
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

    <!-- 鏂板缓/缂栬緫瑙勫垯寮圭獥 -->
    <a-modal v-model:visible="modalVisible" :title="isEdit ? '缂栬緫琛屼负瑙勫垯' : '鏂板缓琛屼负瑙勫垯'" @ok="submitRule" :width="640" :loading="submitting">
      <a-form :model="ruleForm" layout="vertical">
        <a-form-item label="瑙勫垯鍚嶇О" required>
          <a-input v-model="ruleForm.name" placeholder="杈撳叆瑙勫垯鍚嶇О" />
        </a-form-item>
        <a-form-item label="瑙勫垯鎻忚堪">
          <a-textarea v-model="ruleForm.description" placeholder="鎻忚堪瑙勫垯鐢拷" :rows="2" />
        </a-form-item>
        <a-form-item label="瑙﹀彂鏉′欢绫诲瀷" required>
          <a-select v-model="ruleForm.trigger_type" placeholder="閫夋嫨瑙﹀彂绫诲瀷" @change="onTriggerChange">
            <a-option value="time">瀹氭椂瑙﹀彂</a-option>
            <a-option value="event">浜嬩欢瑙﹀彂</a-option>
            <a-option value="threshold">闃堝€艰Е锟?/a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="瑙﹀彂鏉′欢閰嶇疆" required>
          <a-textarea v-model="ruleForm.trigger_config" :rows="4" :placeholder="triggerPlaceholder" />
        </a-form-item>
        <a-form-item label="鍔ㄤ綔绫诲瀷" required>
          <a-select v-model="ruleForm.action_type" placeholder="閫夋嫨鍔ㄤ綔绫诲瀷">
            <a-option value="notify">鍙戦€侀€氱煡</a-option>
            <a-option value="command">涓嬪彂鎸囦护</a-option>
            <a-option value="webhook">璋冪敤Webhook</a-option>
            <a-option value="script">鎵ц鑴氭湰</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="鍔ㄤ綔閰嶇疆" required>
          <a-textarea v-model="ruleForm.action_config" :rows="4" placeholder='{"device_id":"xxx","command":"set_temperature","params":{"value":25}}' />
        </a-form-item>
        <a-form-item label="鍚敤瑙勫垯">
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

const triggerPlaceholder = '// 瀹氭椂: {"cron":"0 8 * * *"} 浜嬩欢: {"event":"device.temperature.high","threshold":30} 闃堬拷? {"metric":"temperature","op":">","value":25}'

const ruleForm = reactive({
  id: 0, name: '', description: '', trigger_type: 'time',
  trigger_config: '', action_type: 'notify', action_config: '', enabled: true,
})

const ruleColumns = [
  { title: '瑙勫垯鍚嶇О', dataIndex: 'name', ellipsis: true },
  { title: '鎻忚堪', dataIndex: 'description', ellipsis: true },
  { title: '瑙﹀彂绫诲瀷', dataIndex: 'trigger_type', slotName: 'trigger_type' },
  { title: '鍚敤', slotName: 'enabled', width: 80 },
  { title: '鎿嶄綔', slotName: 'actions', fixed: 'right', width: 200 },
]

const logColumns = [
  { title: '瑙勫垯鍚嶇О', dataIndex: 'rule_name', ellipsis: true },
  { title: '瑙﹀彂鏉′欢', dataIndex: 'trigger_summary', ellipsis: true },
  { title: '鎵ц鍔ㄤ綔', dataIndex: 'action_summary', ellipsis: true },
  { title: '鐘讹拷?, dataIndex: 'status', slotName: 'status' },
  { title: '鑰楁椂', slotName: 'duration' },
  { title: '鎵ц鏃堕棿', dataIndex: 'created_at', slotName: 'created_at' },
]

const loadRules = async () => {
  loading.value = true
  try {
    const res = await axios.get('/api/v1/ai/behavior-rules', { params: { page: pagination.current, page_size: pagination.pageSize } })
    rules.value = res.data.items || []
    pagination.total = res.data.total || 0
  } catch {
    rules.value = [
      { id: 1, name: '锟?鐐瑰敜锟?, description: '姣忓ぉ鏃╀笂8鐐瑰彂閫佸敜閱掗€氱煡', trigger_type: 'time', trigger_config: '{"cron":"0 8 * * *"}', action_type: 'notify', action_config: '{"title":"鏃╁畨"}', enabled: true },
      { id: 2, name: '娓╁害杩囬珮鍛婅', description: '娓╁害瓒呰繃30搴﹀彂閫佸憡锟?, trigger_type: 'threshold', trigger_config: '{"metric":"temperature","op":">","value":30}', action_type: 'command', action_config: '{"device_id":"dev-001","command":"set_mode","params":{"mode":"cool"}}', enabled: true },
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
      { id: 1, rule_name: '锟?鐐瑰敜锟?, rule_id: 1, trigger_summary: 'cron:0 8 * * *', action_summary: 'notify:鏃╁畨', status: 'success', duration_ms: 45, created_at: new Date().toISOString() },
      { id: 2, rule_name: '娓╁害杩囬珮鍛婅', rule_id: 2, trigger_summary: 'temperature>30', action_summary: 'command:dev-001', status: 'success', duration_ms: 120, created_at: new Date(Date.now() - 600000).toISOString() },
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
      Message.success('鏇存柊鎴愬姛')
    } else {
      await axios.post('/api/v1/ai/behavior-rules', ruleForm)
      Message.success('鍒涘缓鎴愬姛')
    }
    modalVisible.value = false
    loadRules()
  } catch (e) { Message.error('鎿嶄綔澶辫触') } finally { submitting.value = false }
}

const toggleRule = async (record: any) => {
  try {
    await axios.put(`/api/v1/ai/behavior-rules/${record.id}`, { enabled: record.enabled })
    Message.success(record.enabled ? '瑙勫垯宸插惎锟? : '瑙勫垯宸茬锟?)
  } catch { record.enabled = !record.enabled }
}

const deleteRule = async (record: any) => {
  try {
    await axios.delete(`/api/v1/ai/behavior-rules/${record.id}`)
    Message.success('鍒犻櫎鎴愬姛')
    loadRules()
  } catch { Message.error('鍒犻櫎澶辫触') }
}

const viewLogs = (record: any) => {
  activeTab.value = 'logs'
  logFilter.rule_id = String(record.id)
  loadLogs()
}

const handlePageChange = (page: number) => { pagination.current = page; loadRules() }
const handleLogPageChange = (page: number) => { logPagination.current = page; loadLogs() }

const getTriggerColor = (type: string) => ({ time: 'arcoblue', event: 'orange', threshold: 'purple' }[type] || 'gray')
const getTriggerText = (type: string) => ({ time: '瀹氭椂', event: '浜嬩欢', threshold: '闃堬拷? }[type] || type)
const getStatusColor = (s: string) => ({ success: 'green', failed: 'red', skip: 'gray' }[s] || 'gray')
const getStatusText = (s: string) => ({ success: '鎴愬姛', failed: '澶辫触', skip: '璺宠繃' }[s] || s)
const onTriggerChange = () => { ruleForm.trigger_config = '' }
const formatDate = (d: string) => d ? new Date(d).toLocaleString('zh-CN') : '-'

loadRules()
</script>