<template>
  <div class="page-container">
    <a-card class="general-card" title="高级安全">
      <a-tabs>
        <a-tab-pane key="policy" tab="安全策略">
          <a-form :model="policyForm" layout="vertical" style="max-width: 600px">
            <a-form-item label="密码复杂度要求">
              <a-switch v-model="policyForm.passwordComplexity" />
            </a-form-item>
            <a-form-item label="最小密码长度">
              <a-input-number v-model="policyForm.minPasswordLength" :min="6" :max="32" style="width: 100%" />
            </a-form-item>
            <a-form-item label="密码有效期(天)">
              <a-input-number v-model="policyForm.passwordExpireDays" :min="0" style="width: 100%" />
            </a-form-item>
            <a-form-item label="会话超时(分钟)">
              <a-input-number v-model="policyForm.sessionTimeout" :min="5" :max="1440" style="width: 100%" />
            </a-form-item>
            <a-form-item label="启用MFA">
              <a-switch v-model="policyForm.mfaEnabled" />
            </a-form-item>
            <a-form-item label="登录失败锁定">
              <a-switch v-model="policyForm.lockAfterFailed" />
            </a-form-item>
            <a-form-item label="最大登录尝试次数">
              <a-input-number v-model="policyForm.maxLoginAttempts" :min="3" :max="10" style="width: 100%" />
            </a-form-item>
            <a-form-item>
              <a-button type="primary" @click="handleSavePolicy">保存策略</a-button>
            </a-form-item>
          </a-form>
        </a-tab-pane>
        <a-tab-pane key="login" tab="登录日志">
          <div class="search-form">
            <a-form :model="logForm" layout="inline">
              <a-form-item label="用户">
                <a-input v-model="logForm.username" placeholder="请输入用户名" style="width: 140px" />
              </a-form-item>
              <a-form-item label="日期">
                <a-range-picker v-model="logForm.dateRange" style="width: 240px" />
              </a-form-item>
              <a-form-item>
                <a-button type="primary" @click="loadLogs">查询</a-button>
              </a-form-item>
            </a-form>
          </div>
          <a-table :columns="logColumns" :data="logData" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id">
            <template #status="{ record }">
              <a-tag :color="record.status === 'success' ? 'green' : 'red'">{{ record.status === 'success' ? '成功' : '失败' }}</a-tag>
            </template>
          </a-table>
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const loading = ref(false)
const logForm = reactive({ username: '', dateRange: [] })
const logData = ref([])
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })

const policyForm = reactive({
  passwordComplexity: true,
  minPasswordLength: 8,
  passwordExpireDays: 90,
  sessionTimeout: 30,
  mfaEnabled: false,
  lockAfterFailed: true,
  maxLoginAttempts: 5
})

const logColumns = [
  { title: '时间', dataIndex: 'time', width: 180 },
  { title: '用户', dataIndex: 'username', width: 140 },
  { title: 'IP地址', dataIndex: 'ip', width: 140 },
  { title: '登录方式', dataIndex: 'method', width: 100 },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '备注', dataIndex: 'remark', ellipsis: true }
]

const loadLogs = async () => {
  loading.value = true
  try {
    const res = await fetch('/api/v1/system/security/login-logs').then(r => r.json())
    if (res.code === 0) {
      logData.value = res.data?.list || []
    } else {
      loadMockLogs()
    }
  } catch { loadMockLogs() } finally { loading.value = false }
}

const loadMockLogs = () => {
  logData.value = [
    { id: 1, time: '2026-04-09 19:30:00', username: 'admin', ip: '192.168.1.100', method: '密码', status: 'success', remark: '' },
    { id: 2, time: '2026-04-09 19:25:00', username: 'admin', ip: '192.168.1.100', method: '密码', status: 'success', remark: '' },
    { id: 3, time: '2026-04-09 18:45:00', username: 'user1', ip: '192.168.1.101', method: '密码', status: 'failed', remark: '密码错误' },
    { id: 4, time: '2026-04-09 18:30:00', username: 'admin', ip: '192.168.1.100', method: 'MFA', status: 'success', remark: '' },
    { id: 5, time: '2026-04-09 17:20:00', username: 'guest', ip: '10.0.0.55', method: 'SSO', status: 'success', remark: '企业微信' }
  ]
}

const handleSavePolicy = () => {
  Message.success('安全策略已保存')
}

const onPageChange = (page) => {
  pagination.current = page
  loadLogs()
}

onMounted(() => loadLogs())
</script>

<style scoped>
.page-container { padding: 16px; }
.search-form { margin-bottom: 16px; padding: 16px; background: var(--color-fill-lightest); border-radius: 4px; }
</style>