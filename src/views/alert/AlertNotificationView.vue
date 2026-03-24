<template>
  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>告警管理</a-breadcrumb-item>
      <a-breadcrumb-item>告警通知</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 页面标题 & 操作区 -->
    <div class="page-header">
      <div class="page-title">
        <icon-broadcast class="title-icon" />
        <span>告警通知</span>
      </div>
      <div class="page-actions">
        <a-button type="primary" @click="handleRefresh">
          <template #icon><icon-refresh /></template>
          「刷新」
        </a-button>
      </div>
    </div>

    <!-- 标签页 -->
    <div class="tab-container">
      <a-tabs v-model:activeKey="activeTab" @change="onTabChange">
        <a-tab-pane key="config" title="告警通知配置" />
        <a-tab-pane key="history" title="告警历史" />
        <a-tab-pane key="logs" title="通知日志" />
        <a-tab-pane key="stats" title="统计报表" />
      </a-tabs>
    </div>

    <!-- 内容区域 -->
    <div class="content-area">
      <!-- 阶段1: 告警通知配置 -->
      <div v-show="activeTab === 'config'" class="config-content">
        <!-- 渠道卡片区域 -->
        <div class="channels-grid">
          <!-- 邮件渠道 -->
          <div class="channel-card" :class="{ 'channel-disabled': !emailChannel }">
            <div class="channel-icon">
              <icon-email :size="32" />
            </div>
            <div class="channel-info">
              <div class="channel-title">
                <span class="channel-name">📧 邮件渠道</span>
                <a-tag :color="emailChannel?.enabled ? 'green' : 'gray'" size="small">
                  {{ emailChannel?.enabled ? '已启用' : '停用' }}
                </a-tag>
              </div>
              <div class="channel-status">
                {{ emailChannel ? '已配置' : '未配置' }} |
                <span v-if="emailChannel?.config?.smtp_host">
                  {{ emailChannel.config.smtp_host }}:{{ emailChannel.config.smtp_port || 587 }}
                </span>
                <span v-else class="text-muted">点击配置</span>
              </div>
            </div>
            <div class="channel-actions">
              <a-button size="small" @click="openEmailConfig">
                「{{ emailChannel ? '编辑' : '配置' }}」
              </a-button>
              <a-button v-if="emailChannel" size="small" @click="testChannel('email')">
                「测试」
              </a-button>
            </div>
          </div>

          <!-- 短信渠道 -->
          <div class="channel-card" :class="{ 'channel-disabled': !smsChannel }">
            <div class="channel-icon">
              <icon-message :size="32" />
            </div>
            <div class="channel-info">
              <div class="channel-title">
                <span class="channel-name">📱 短信渠道</span>
                <a-tag :color="smsChannel?.enabled ? 'green' : 'gray'" size="small">
                  {{ smsChannel?.enabled ? '已启用' : '停用' }}
                </a-tag>
              </div>
              <div class="channel-status">
                {{ smsChannel ? '已配置' : '未配置' }} |
                <span v-if="smsChannel?.config?.sms_provider">
                  {{ smsProviderLabel(smsChannel.config.sms_provider) }}
                </span>
                <span v-else class="text-muted">点击配置</span>
              </div>
            </div>
            <div class="channel-actions">
              <a-button size="small" @click="openSmsConfig">
                「{{ smsChannel ? '编辑' : '配置' }}」
              </a-button>
              <a-button v-if="smsChannel" size="small" @click="testChannel('sms')">
                「测试」
              </a-button>
            </div>
          </div>

          <!-- Webhook渠道 -->
          <div class="channel-card" :class="{ 'channel-disabled': !webhookChannel }">
            <div class="channel-icon">
              <icon-link :size="32" />
            </div>
            <div class="channel-info">
              <div class="channel-title">
                <span class="channel-name">🔗 Webhook</span>
                <a-tag :color="webhookChannel?.enabled ? 'green' : 'gray'" size="small">
                  {{ webhookChannel?.enabled ? '已启用' : '停用' }}
                </a-tag>
              </div>
              <div class="channel-status">
                {{ webhookChannel ? '已配置' : '未配置' }} |
                <span v-if="webhookChannel?.config?.webhook_url" class="url-text">
                  {{ truncateUrl(webhookChannel.config.webhook_url) }}
                </span>
                <span v-else class="text-muted">点击配置</span>
              </div>
            </div>
            <div class="channel-actions">
              <a-button size="small" @click="openWebhookConfig">
                「{{ webhookChannel ? '编辑' : '配置' }}」
              </a-button>
              <a-button v-if="webhookChannel" size="small" @click="testChannel('webhook')">
                「测试」
              </a-button>
            </div>
          </div>

          <!-- 通知时段 -->
          <div class="channel-card">
            <div class="channel-icon">
              <icon-clock :size="32" />
            </div>
            <div class="channel-info">
              <div class="channel-title">
                <span class="channel-name">⏰ 通知时段</span>
              </div>
              <div class="channel-status">
                <span v-if="periodText">{{ periodText }}</span>
                <span v-else class="text-muted">全天 00:00-23:59</span>
              </div>
            </div>
            <div class="channel-actions">
              <a-button size="small" @click="openPeriodConfig">
                「配置」
              </a-button>
            </div>
          </div>
        </div>
      </div>

      <!-- 告警历史 -->
      <div v-show="activeTab === 'history'" class="history-content">
        <AlertHistoryView />
      </div>

      <!-- 通知日志 -->
      <div v-show="activeTab === 'logs'" class="logs-content">
        <NotificationLogsView />
      </div>

      <!-- 统计报表 -->
      <div v-show="activeTab === 'stats'" class="stats-content">
        <NotificationStatsView />
      </div>
    </div>

    <!-- 邮件配置弹窗 -->
    <a-modal
      v-model:visible="emailModalVisible"
      title="邮件渠道配置"
      :width="640"
      :footer="null"
      @cancel="emailModalVisible = false"
    >
      <EmailChannelConfig
        :channel="emailChannel"
        @success="onChannelSaved"
        @cancel="emailModalVisible = false"
      />
    </a-modal>

    <!-- 短信配置弹窗 -->
    <a-modal
      v-model:visible="smsModalVisible"
      title="短信渠道配置"
      :width="640"
      :footer="null"
      @cancel="smsModalVisible = false"
    >
      <SMSChannelConfig
        :channel="smsChannel"
        @success="onChannelSaved"
        @cancel="smsModalVisible = false"
      />
    </a-modal>

    <!-- Webhook配置弹窗 -->
    <a-modal
      v-model:visible="webhookModalVisible"
      title="Webhook 配置"
      :width="640"
      :footer="null"
      @cancel="webhookModalVisible = false"
    >
      <WebhookChannelConfig
        :channel="webhookChannel"
        @success="onChannelSaved"
        @cancel="webhookModalVisible = false"
      />
    </a-modal>

    <!-- 通知时段配置弹窗 -->
    <a-modal
      v-model:visible="periodModalVisible"
      title="通知时段配置"
      :width="520"
      @ok="savePeriods"
      @cancel="periodModalVisible = false"
    >
      <div class="period-form">
        <div class="period-item" v-for="(period, index) in editingPeriods" :key="index">
          <a-input v-model="period.name" placeholder="时段名称" style="width: 120px" />
          <a-time-picker v-model="period.start_time" format="HH:mm" placeholder="开始时间" style="width: 120px" />
          <span>至</span>
          <a-time-picker v-model="period.end_time" format="HH:mm" placeholder="结束时间" style="width: 120px" />
          <a-switch v-model="period.enabled" />
          <a-button type="text" status="danger" @click="removePeriod(index)">
            <icon-delete />
          </a-button>
        </div>
        <a-button type="dashed" @click="addPeriod">
          <template #icon><icon-plus /></template>
          「添加时段」
        </a-button>
      </div>
    </a-modal>

    <!-- 测试结果弹窗 -->
    <a-modal
      v-model:visible="testResultVisible"
      title="测试结果"
      :footer="null"
      :width="400"
    >
      <a-result
        v-if="testResult"
        :status="testResult.success ? 'success' : 'error'"
        :title="testResult.success ? '测试成功' : '测试失败'"
        :sub-title="testResult.message || testResult.error || ''"
      />
    </a-modal>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import {
  IconLiveBroadcast,
  IconRefresh,
  IconEmail,
  IconMessage,
  IconLink,
  IconClockCircle,
  IconDelete,
  IconPlus
} from '@arco-design/web-vue/es/icon'
import { useNotificationChannels, useNotificationPeriods } from '@/composables/useNotification'
import EmailChannelConfig from './EmailChannelConfig.vue'
import SMSChannelConfig from './SMSChannelConfig.vue'
import WebhookChannelConfig from './WebhookChannelConfig.vue'
import AlertHistoryView from './AlertHistoryView.vue'
import NotificationLogsView from './NotificationLogsView.vue'

// 动态导入统计报表（避免循环依赖）
const NotificationStatsView = defineAsyncComponent(() =>
  import('./NotificationStatsView.vue')
)

import { defineAsyncComponent as defineAsyncComponent } from 'vue'

const activeTab = ref('config')
const emailModalVisible = ref(false)
const smsModalVisible = ref(false)
const webhookModalVisible = ref(false)
const periodModalVisible = ref(false)
const testResultVisible = ref(false)
const testResult = ref(null)

const {
  channels,
  emailChannels,
  smsChannels,
  webhookChannels,
  loadChannels,
  testChannel: testNotificationChannel
} = useNotificationChannels()

const {
  periods,
  defaultPeriod,
  loadPeriods,
  savePeriods: savePeriodsApi
} = useNotificationPeriods()

// 取第一个配置的渠道
const emailChannel = computed(() => emailChannels.value[0] || null)
const smsChannel = computed(() => smsChannels.value[0] || null)
const webhookChannel = computed(() => webhookChannels.value[0] || null)

// 通知时段文本
const periodText = computed(() => {
  if (!periods.value?.length) return ''
  return periods.value.map(p => `${p.name} ${p.start_time}-${p.end_time}`).join(', ')
})

const editingPeriods = ref([])

function onTabChange(key) {
  // 切换标签页
}

function handleRefresh() {
  loadChannels()
  loadPeriods()
  Message.success('已刷新')
}

function openEmailConfig() {
  emailModalVisible.value = true
}

function openSmsConfig() {
  smsModalVisible.value = true
}

function openWebhookConfig() {
  webhookModalVisible.value = true
}

function openPeriodConfig() {
  editingPeriods.value = JSON.parse(JSON.stringify(periods.value || [defaultPeriod]))
  periodModalVisible.value = true
}

function addPeriod() {
  editingPeriods.value.push({
    name: `时段${editingPeriods.value.length + 1}`,
    start_time: '09:00',
    end_time: '18:00',
    enabled: true
  })
}

function removePeriod(index) {
  editingPeriods.value.splice(index, 1)
}

async function savePeriods() {
  await savePeriodsApi(editingPeriods.value)
  periodModalVisible.value = false
}

async function testChannel(type) {
  const ch = { email: emailChannel, sms: smsChannel, webhook: webhookChannel }[type]?.value
  if (!ch) {
    Message.warning('请先配置渠道')
    return
  }
  testResult.value = null
  testResultVisible.value = true
  const result = await testNotificationChannel(ch.id)
  testResult.value = result
}

function onChannelSaved() {
  emailModalVisible.value = false
  smsModalVisible.value = false
  webhookModalVisible.value = false
  loadChannels()
}

function truncateUrl(url) {
  if (!url) return ''
  return url.length > 30 ? url.substring(0, 30) + '...' : url
}

function smsProviderLabel(provider) {
  const map = { aliyun: '阿里云', tencent: '腾讯云', huawei: '华为云' }
  return map[provider] || provider
}

onMounted(() => {
  loadChannels()
  loadPeriods()
})
</script>

<style scoped>
.pro-page-container {
  padding: 20px 24px;
  min-height: calc(100vh - 64px);
  background: #f5f7fa;
}
.pro-breadcrumb {
  margin-bottom: 16px;
}
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: #fff;
  border-radius: 8px 8px 0 0;
  padding: 16px 20px;
}
.page-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 18px;
  font-weight: 600;
  color: #1d2129;
}
.title-icon {
  font-size: 24px;
  color: #1650d8;
}
.page-actions {
  display: flex;
  gap: 8px;
}
.tab-container {
  background: #fff;
  margin-bottom: 0;
}
.tab-container :deep(.arco-tabs-nav) {
  margin-bottom: 0;
  padding: 0 20px;
}
.content-area {
  background: #f5f7fa;
  padding: 20px;
  min-height: 400px;
}
.channels-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
}
.channel-card {
  background: #fff;
  border-radius: 8px;
  padding: 20px;
  display: flex;
  align-items: flex-start;
  gap: 16px;
  transition: all 0.2s;
}
.channel-card:hover {
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
}
.channel-disabled {
  opacity: 0.7;
}
.channel-icon {
  width: 48px;
  height: 48px;
  border-radius: 8px;
  background: #e6f4ff;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #1650d8;
  flex-shrink: 0;
}
.channel-info {
  flex: 1;
  min-width: 0;
}
.channel-title {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}
.channel-name {
  font-size: 16px;
  font-weight: 600;
  color: #1d2129;
}
.channel-status {
  font-size: 14px;
  color: #86909c;
}
.text-muted {
  color: #c9cdd4;
}
.url-text {
  max-width: 200px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  display: inline-block;
}
.channel-actions {
  display: flex;
  flex-direction: column;
  gap: 8px;
  flex-shrink: 0;
}
.period-form {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
.period-item {
  display: flex;
  align-items: center;
  gap: 8px;
}
.history-content,
.logs-content,
.stats-content {
  background: #fff;
  border-radius: 8px;
}
</style>
