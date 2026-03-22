<template>
  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>告警管理</a-breadcrumb-item>
      <a-breadcrumb-item>通知渠道</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 页面标题 -->
    <div class="page-header">
      <a-tabs v-model:activeKey="activeTab" @change="onTabChange">
        <a-tab-pane key="smtp" title="邮件 (SMTP)" />
        <a-tab-pane key="webhook" title="Webhook" />
        <a-tab-pane key="sms" title="短信 (SMS)" />
      </a-tabs>
    </div>

    <!-- 渠道卡片列表 -->
    <div class="channels-list">
      <a-row :gutter="16">
        <a-col :span="8" v-for="ch in filteredChannels" :key="ch.id">
          <a-card class="channel-card" :class="{ 'channel-disabled': !ch.enabled }">
            <template #title>
              <div class="channel-card-title">
                <a-space>
                  <icon-email v-if="ch.channel_type === 'smtp'" />
                  <icon-link v-else-if="ch.channel_type === 'webhook'" />
                  <icon-message v-else />
                  <span>{{ ch.name }}</span>
                </a-space>
                <a-tag :color="ch.enabled ? 'green' : 'gray'" size="small">
                  {{ ch.enabled ? '已启用' : '已停用' }}
                </a-tag>
              </div>
            </template>
            <template #extra>
              <a-space>
                <a-switch v-model="ch.enabled" size="small" @change="handleToggle(ch)" />
                <a-button type="text" size="small" @click="editChannel(ch)">编辑</a-button>
                <a-popconfirm content="确定删除该渠道？" @ok="handleDelete(ch.id)">
                  <a-button type="text" size="small" status="danger">删除</a-button>
                </a-popconfirm>
              </a-space>
            </template>

            <a-descriptions :column="1" size="small">
              <a-descriptions-item v-if="ch.channel_type === 'smtp'" label="SMTP 服务器">
                {{ ch.smtp_host || '-' }}
              </a-descriptions-item>
              <a-descriptions-item v-if="ch.channel_type === 'smtp'" label="端口">
                {{ ch.smtp_port || '-' }}
              </a-descriptions-item>
              <a-descriptions-item v-if="ch.channel_type === 'smtp'" label="发件人">
                {{ ch.smtp_from || '-' }}
              </a-descriptions-item>
              <a-descriptions-item v-if="ch.channel_type === 'smtp'" label="收件人">
                {{ ch.smtp_to || '-' }}
              </a-descriptions-item>
              <a-descriptions-item v-if="ch.channel_type === 'webhook'" label="URL">
                <a-tooltip :content="ch.webhook_url">{{ truncateUrl(ch.webhook_url) }}</a-tooltip>
              </a-descriptions-item>
              <a-descriptions-item v-if="ch.channel_type === 'webhook'" label="请求方式">
                {{ ch.webhook_method || 'POST' }}
              </a-descriptions-item>
              <a-descriptions-item v-if="ch.channel_type === 'sms'" label="Provider">
                {{ ch.sms_provider || '-' }}
              </a-descriptions-item>
              <a-descriptions-item v-if="ch.channel_type === 'sms'" label="发件号">
                {{ ch.sms_from || '-' }}
              </a-descriptions-item>
              <a-descriptions-item label="备注">
                {{ ch.remark || '-' }}
              </a-descriptions-item>
            </a-descriptions>

            <div class="card-actions">
              <a-button type="primary" size="small" @click="testChannel(ch)">测试</a-button>
              <a-button size="small" @click="editChannel(ch)">修改配置</a-button>
            </div>
          </a-card>
        </a-col>

        <!-- 添加渠道卡片 -->
        <a-col :span="8">
          <a-card class="add-channel-card" @click="showAddModal">
            <div class="add-channel-inner">
              <icon-plus-circle :size="40" />
              <span>添加{{ tabLabel }}渠道</span>
            </div>
          </a-card>
        </a-col>
      </a-row>
    </div>

    <!-- 添加/编辑渠道弹窗 -->
    <a-modal
      v-model:visible="modalVisible"
      :title="isEdit ? '编辑通知渠道' : `添加${tabLabel}渠道`"
      :width="560"
      @ok="handleSubmit"
      @cancel="modalVisible = false"
    >
      <a-form :model="form" layout="vertical" :rules="formRules">
        <a-form-item label="渠道名称" field="name">
          <a-input v-model="form.name" placeholder="请输入渠道名称，如：生产环境邮件" />
        </a-form-item>
        <a-form-item label="启用该渠道">
          <a-switch v-model="form.enabled" />
        </a-form-item>

        <!-- SMTP 配置 -->
        <template v-if="activeTab === 'smtp'">
          <a-divider>邮件服务器配置</a-divider>
          <a-form-item label="SMTP 服务器" field="smtp_host" required>
            <a-input v-model="form.smtp_host" placeholder="smtp.example.com" />
          </a-form-item>
          <a-form-item label="端口" field="smtp_port">
            <a-input-number v-model="form.smtp_port" :min="1" :max="65535" placeholder="587" style="width: 100%" />
          </a-form-item>
          <a-form-item label="用户名" field="smtp_user">
            <a-input v-model="form.smtp_user" placeholder="alert@example.com" />
          </a-form-item>
          <a-form-item label="密码/授权码" field="smtp_password">
            <a-input-password v-model="form.smtp_password" placeholder="请输入密码或授权码" />
          </a-form-item>
          <a-form-item label="使用 TLS">
            <a-switch v-model="form.smtp_use_tls" />
          </a-form-item>
          <a-form-item label="发件人邮箱" field="smtp_from">
            <a-input v-model="form.smtp_from" placeholder="alert@example.com" />
          </a-form-item>
          <a-form-item label="默认收件人（多个用逗号分隔）" field="smtp_to">
            <a-input v-model="form.smtp_to" placeholder="admin@example.com,ops@example.com" />
          </a-form-item>
        </template>

        <!-- Webhook 配置 -->
        <template v-else-if="activeTab === 'webhook'">
          <a-divider>Webhook 配置</a-divider>
          <a-form-item label="Webhook URL" field="webhook_url" required>
            <a-input v-model="form.webhook_url" placeholder="https://your-webhook-endpoint.com/alert" />
          </a-form-item>
          <a-form-item label="请求方式" field="webhook_method">
            <a-select v-model="form.webhook_method">
              <a-option value="POST">POST</a-option>
              <a-option value="PUT">PUT</a-option>
            </a-select>
          </a-form-item>
          <a-form-item label="Token（可选）" field="webhook_token">
            <a-input-password v-model="form.webhook_token" placeholder="Bearer Token（如有）" />
          </a-form-item>
        </template>

        <!-- SMS 配置 -->
        <template v-else>
          <a-divider>短信配置</a-divider>
          <a-form-item label="短信服务商" field="sms_provider">
            <a-select v-model="form.sms_provider" placeholder="选择短信服务商">
              <a-option value="aliyun">阿里云</a-option>
              <a-option value="tencent">腾讯云</a-option>
              <a-option value="twilio">Twilio</a-option>
              <a-option value="custom">自定义</a-option>
            </a-select>
          </a-form-item>
          <a-form-item label="Account SID / Access Key" field="sms_account">
            <a-input v-model="form.sms_account" placeholder="请输入 Account SID 或 Access Key" />
          </a-form-item>
          <a-form-item label="Secret / Auth Token" field="sms_secret">
            <a-input-password v-model="form.sms_secret" placeholder="请输入 Secret" />
          </a-form-item>
          <a-form-item label="发件号码" field="sms_from">
            <a-input v-model="form.sms_from" placeholder="+86xxxxxxxxxxx" />
          </a-form-item>
        </template>

        <a-divider />
        <a-form-item label="备注">
          <a-textarea v-model="form.remark" :rows="2" placeholder="可选备注信息" />
        </a-form-item>
        <a-form-item label="设为默认渠道">
          <a-switch v-model="form.is_default" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 测试结果反馈 -->
    <a-modal v-model:visible="testResultVisible" title="测试结果" :footer="null" :width="400">
      <a-result
        v-if="testResult"
        :status="testResult.test_status === 'success' ? 'success' : 'error'"
        :title="testResult.test_status === 'success' ? '测试成功' : '测试失败'"
        :sub-title="testResult.message || testResult.error || ''"
      />
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import * as api from '@/api/alerts'

const activeTab = ref('smtp')
const channels = ref([])
const modalVisible = ref(false)
const isEdit = ref(false)
const testResultVisible = ref(false)
const testResult = ref(null)
const editingId = ref(null)

const form = reactive({
  name: '',
  channel_type: 'smtp',
  enabled: true,
  // SMTP
  smtp_host: '',
  smtp_port: 587,
  smtp_user: '',
  smtp_password: '',
  smtp_from: '',
  smtp_to: '',
  smtp_use_tls: true,
  // Webhook
  webhook_url: '',
  webhook_method: 'POST',
  webhook_token: '',
  // SMS
  sms_provider: '',
  sms_account: '',
  sms_secret: '',
  sms_from: '',
  // Common
  remark: '',
  is_default: false
})

const formRules = {
  name: [{ required: true, message: '请输入渠道名称' }],
  smtp_host: [{ required: true, message: '请输入 SMTP 服务器' }],
  webhook_url: [{ required: true, message: '请输入 Webhook URL' }]
}

const tabLabel = computed(() => ({
  smtp: '邮件',
  webhook: 'Webhook',
  sms: '短信'
}[activeTab.value] || ''))

const filteredChannels = computed(() =>
  channels.value.filter(ch => ch.channel_type === activeTab.value)
)

function onTabChange() {
  // reset form channel type when tab changes
  form.channel_type = activeTab.value
}

function truncateUrl(url) {
  if (!url) return '-'
  return url.length > 40 ? url.substring(0, 40) + '...' : url
}

async function loadChannels() {
  try {
    const data = await api.getNotificationChannels()
    channels.value = data.data?.list || []
  } catch (e) {
    Message.error('加载渠道列表失败')
  }
}

function showAddModal() {
  isEdit.value = false
  editingId.value = null
  Object.assign(form, {
    name: '',
    channel_type: activeTab.value,
    enabled: true,
    smtp_host: '', smtp_port: 587, smtp_user: '', smtp_password: '',
    smtp_from: '', smtp_to: '', smtp_use_tls: true,
    webhook_url: '', webhook_method: 'POST', webhook_token: '',
    sms_provider: '', sms_account: '', sms_secret: '', sms_from: '',
    remark: '', is_default: false
  })
  modalVisible.value = true
}

function editChannel(ch) {
  isEdit.value = true
  editingId.value = ch.id
  Object.assign(form, {
    name: ch.name,
    channel_type: ch.channel_type,
    enabled: ch.enabled,
    smtp_host: ch.smtp_host || '',
    smtp_port: ch.smtp_port || 587,
    smtp_user: ch.smtp_user || '',
    smtp_password: ch.smtp_password || '',
    smtp_from: ch.smtp_from || '',
    smtp_to: ch.smtp_to || '',
    smtp_use_tls: ch.smtp_use_tls !== false,
    webhook_url: ch.webhook_url || '',
    webhook_method: ch.webhook_method || 'POST',
    webhook_token: ch.webhook_token || '',
    sms_provider: ch.sms_provider || '',
    sms_account: ch.sms_account || '',
    sms_secret: ch.sms_secret || '',
    sms_from: ch.sms_from || '',
    remark: ch.remark || '',
    is_default: ch.is_default || false
  })
  activeTab.value = ch.channel_type
  modalVisible.value = true
}

async function handleSubmit() {
  if (!form.name) {
    Message.warning('请输入渠道名称')
    return
  }
  if (activeTab.value === 'smtp' && !form.smtp_host) {
    Message.warning('请输入 SMTP 服务器')
    return
  }
  if (activeTab.value === 'webhook' && !form.webhook_url) {
    Message.warning('请输入 Webhook URL')
    return
  }
  form.channel_type = activeTab.value

  try {
    if (isEdit.value) {
      await api.updateNotificationChannel(editingId.value, { ...form })
      Message.success('更新成功')
    } else {
      await api.createNotificationChannel({ ...form })
      Message.success('添加成功')
    }
    modalVisible.value = false
    loadChannels()
  } catch (e) {
    Message.error(isEdit.value ? '更新失败' : '添加失败')
  }
}

async function handleToggle(ch) {
  try {
    await api.toggleNotificationChannel(ch.id)
    Message.success(ch.enabled ? '已启用' : '已停用')
  } catch (e) {
    ch.enabled = !ch.enabled // revert
    Message.error('操作失败')
  }
}

async function handleDelete(id) {
  try {
    await api.deleteNotificationChannel(id)
    Message.success('删除成功')
    loadChannels()
  } catch (e) {
    Message.error('删除失败')
  }
}

async function testChannel(ch) {
  try {
    testResult.value = null
    testResultVisible.value = true
    const data = await api.testNotificationChannel(ch.id)
    testResult.value = data.data
  } catch (e) {
    testResult.value = { test_status: 'error', error: '测试请求失败' }
    testResultVisible.value = true
  }
}

onMounted(() => {
  loadChannels()
})
</script>

<style scoped>
.pro-page-container { padding: 20px 24px; min-height: calc(100vh - 64px); background: #f5f7fa; }
.pro-breadcrumb { margin-bottom: 16px; }
.page-header { background: #fff; border-radius: 8px 8px 0 0; padding: 0 20px; margin-bottom: 0; }
.page-header :deep(.arco-tabs-nav) { margin-bottom: 0; }
.channels-list { background: #f5f7fa; padding: 16px 0; }
.channel-card { border-radius: 8px; }
.channel-disabled { opacity: 0.7; }
.channel-card-title { display: flex; align-items: center; gap: 8px; }
.card-actions { margin-top: 12px; display: flex; gap: 8px; }
.add-channel-card {
  border-radius: 8px; border: 2px dashed #c0c0c0; cursor: pointer;
  min-height: 200px; display: flex; align-items: center; justify-content: center;
  transition: all 0.2s;
}
.add-channel-card:hover { border-color: #1650d8; }
.add-channel-inner {
  display: flex; flex-direction: column; align-items: center; gap: 8px;
  color: #888; font-size: 14px;
}
</style>
