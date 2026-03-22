<template>
  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>告警管理</a-breadcrumb-item>
      <a-breadcrumb-item>告警设置</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 设置表单卡片 -->
    <a-card class="settings-card">
      <template #title>
        <div class="card-title-row">
          <icon-bell />
          <span>告警全局设置</span>
        </div>
      </template>

      <a-form :model="settings" layout="vertical" :loading="loading">
        <!-- 全局开关 -->
        <a-divider orientation="left">全局开关</a-divider>
        <a-form-item label="启用告警系统">
          <a-switch v-model="settings.alerts_enabled" />
          <span class="setting-hint">关闭后，所有告警规则将不再触发通知</span>
        </a-form-item>

        <!-- 通知渠道开关 -->
        <a-divider orientation="left">通知渠道</a-divider>
        <a-row :gutter="16">
          <a-col :span="6">
            <a-form-item label="邮件通知">
              <a-switch v-model="settings.email_enabled" />
            </a-form-item>
          </a-col>
          <a-col :span="6">
            <a-form-item label="短信通知">
              <a-switch v-model="settings.sms_enabled" />
            </a-form-item>
          </a-col>
          <a-col :span="6">
            <a-form-item label="Webhook通知">
              <a-switch v-model="settings.webhook_enabled" />
            </a-form-item>
          </a-col>
          <a-col :span="6">
            <a-form-item label="站内通知">
              <a-switch v-model="settings.inapp_enabled" />
            </a-form-item>
          </a-col>
        </a-row>

        <!-- 告警级别 -->
        <a-divider orientation="left">告警级别通知</a-divider>
        <a-form-item label="选择需要通知的告警级别">
          <a-space>
            <a-checkbox v-model="settings.notify_on_critical">
              <a-tag color="red">严重</a-tag> 严重告警
            </a-checkbox>
            <a-checkbox v-model="settings.notify_on_high">
              <a-tag color="orange">高</a-tag> 高优先级告警
            </a-checkbox>
            <a-checkbox v-model="settings.notify_on_medium">
              <a-tag color="arcoblue">中</a-tag> 中优先级告警
            </a-checkbox>
            <a-checkbox v-model="settings.notify_on_low">
              <a-tag color="green">低</a-tag> 低优先级告警
            </a-checkbox>
          </a-space>
        </a-form-item>

        <!-- 告警限流 -->
        <a-divider orientation="left">限流与自动处理</a-divider>
        <a-form-item label="每小时最大告警数">
          <a-space>
            <a-input-number v-model="settings.max_per_hour" :min="1" :max="10000" />
            <span class="setting-hint">超过后新告警将被抑制，防止告警风暴</span>
          </a-space>
        </a-form-item>

        <a-form-item label="自动解决超时（小时）">
          <a-space>
            <a-input-number v-model="settings.auto_resolve_hours" :min="0" :max="720" />
            <span class="setting-hint">0 表示不自动解决</span>
          </a-space>
        </a-form-item>

        <!-- 告警摘要 -->
        <a-divider orientation="left">告警摘要</a-divider>
        <a-form-item label="启用告警摘要">
          <a-switch v-model="settings.digest_enabled" />
          <span class="setting-hint">定时汇总发送告警，减少通知打扰</span>
        </a-form-item>
        <a-form-item v-if="settings.digest_enabled" label="摘要间隔（分钟）">
          <a-select v-model="settings.digest_interval" style="width: 200px">
            <a-option :value="15">每 15 分钟</a-option>
            <a-option :value="30">每 30 分钟</a-option>
            <a-option :value="60">每小时</a-option>
            <a-option :value="120">每 2 小时</a-option>
            <a-option :value="360">每 6 小时</a-option>
            <a-option :value="1440">每天</a-option>
          </a-select>
        </a-form-item>

        <!-- 静默时段 -->
        <a-divider orientation="left">静默时段</a-divider>
        <a-form-item label="启用静默时段">
          <a-switch v-model="settings.quiet_hours_enabled" />
          <span class="setting-hint">静默时段内生成的告警将被延迟通知</span>
        </a-form-item>
        <a-form-item v-if="settings.quiet_hours_enabled" label="静默时间段">
          <a-space>
            <a-time-picker v-model="quietStart" format="HH:mm" placeholder="开始时间" />
            <span>至</span>
            <a-time-picker v-model="quietEnd" format="HH:mm" placeholder="结束时间" />
          </a-space>
        </a-form-item>

        <!-- 快捷链接 -->
        <a-divider orientation="left">相关链接</a-divider>
        <a-space>
          <a-link @click="$router.push('/alert/rules')">
            <icon-list /> 告警规则
          </a-link>
          <a-link @click="$router.push('/alert/list')">
            <icon-burger /> 告警列表
          </a-link>
          <a-link @click="$router.push('/alert/channels')">
            <icon-email /> 通知渠道
          </a-link>
        </a-space>

        <!-- 保存按钮 -->
        <div class="form-actions">
          <a-button type="primary" :loading="saving" @click="handleSave">保存设置</a-button>
          <a-button @click="loadSettings">重置</a-button>
        </div>
      </a-form>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import * as api from '@/api/alerts'
import dayjs from 'dayjs'

const loading = ref(false)
const saving = ref(false)

const settings = reactive({
  alerts_enabled: true,
  email_enabled: false,
  sms_enabled: false,
  webhook_enabled: false,
  inapp_enabled: true,
  notify_on_critical: true,
  notify_on_high: true,
  notify_on_medium: true,
  notify_on_low: false,
  digest_enabled: false,
  digest_interval: 60,
  quiet_hours_enabled: false,
  quiet_hours_start: '22:00',
  quiet_hours_end: '08:00',
  max_per_hour: 100,
  auto_resolve_hours: 24
})

// 静默时段用 dayjs 对象方便绑定
const quietStart = ref(null)
const quietEnd = ref(null)

async function loadSettings() {
  loading.value = true
  try {
    const data = await api.getAlertSettings()
    const s = data.data || {}
    Object.assign(settings, {
      alerts_enabled: s.alerts_enabled ?? true,
      email_enabled: s.email_enabled ?? false,
      sms_enabled: s.sms_enabled ?? false,
      webhook_enabled: s.webhook_enabled ?? false,
      inapp_enabled: s.inapp_enabled ?? true,
      notify_on_critical: s.notify_on_critical ?? true,
      notify_on_high: s.notify_on_high ?? true,
      notify_on_medium: s.notify_on_medium ?? true,
      notify_on_low: s.notify_on_low ?? false,
      digest_enabled: s.digest_enabled ?? false,
      digest_interval: s.digest_interval ?? 60,
      quiet_hours_enabled: s.quiet_hours_enabled ?? false,
      quiet_hours_start: s.quiet_hours_start ?? '22:00',
      quiet_hours_end: s.quiet_hours_end ?? '08:00',
      max_per_hour: s.max_per_hour ?? 100,
      auto_resolve_hours: s.auto_resolve_hours ?? 24
    })
    quietStart.value = s.quiet_hours_start ? dayjs(`2024-01-01 ${s.quiet_hours_start}`) : null
    quietEnd.value = s.quiet_hours_end ? dayjs(`2024-01-01 ${s.quiet_hours_end}`) : null
  } catch (e) {
    Message.error('加载告警设置失败')
  } finally {
    loading.value = false
  }
}

async function handleSave() {
  saving.value = true
  try {
    const payload = {
      ...settings,
      quiet_hours_start: quietStart.value ? quietStart.value.format('HH:mm') : '22:00',
      quiet_hours_end: quietEnd.value ? quietEnd.value.format('HH:mm') : '08:00'
    }
    await api.updateAlertSettings(payload)
    Message.success('设置已保存')
  } catch (e) {
    Message.error('保存失败')
  } finally {
    saving.value = false
  }
}

onMounted(() => { loadSettings() })
</script>

<style scoped>
.pro-page-container { padding: 20px 24px; min-height: calc(100vh - 64px); background: #f5f7fa; }
.pro-breadcrumb { margin-bottom: 16px; }
.settings-card { border-radius: 8px; max-width: 800px; }
.card-title-row { display: flex; align-items: center; gap: 8px; font-weight: 600; font-size: 15px; }
.setting-hint { margin-left: 12px; font-size: 12px; color: #888; }
.form-actions { margin-top: 24px; display: flex; gap: 12px; }
</style>
