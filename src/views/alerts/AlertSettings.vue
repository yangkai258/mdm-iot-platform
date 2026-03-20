<template>
  <a-layout class="alert-settings">
    <a-layout-sider v-model:collapsed="collapsed" :trigger="null" collapsible>
      <div class="logo">
        <span v-if="!collapsed">MDM 控制台</span>
      </div>
      <a-menu v-model:selectedKeys="selectedKeys" theme="dark" mode="inline" @click="handleMenuClick">
        <a-menu-item key="dashboard">
          <span>设备大盘</span>
        </a-menu-item>
        <a-menu-item key="alerts">
          <span>告警中心</span>
        </a-menu-item>
      </a-menu>
    </a-layout-sider>

    <a-layout>
      <a-layout-header class="header">
        <div class="header-left">
          <a-button type="text" @click="collapsed = !collapsed">
            <span v-if="collapsed">☰</span>
            <span v-else>✕</span>
          </a-button>
        </div>
        <div class="header-title">
          <a-breadcrumb>
            <a-breadcrumb-item>告警中心</a-breadcrumb-item>
            <a-breadcrumb-item>告警设置</a-breadcrumb-item>
          </a-breadcrumb>
        </div>
        <div class="header-right"></div>
      </a-layout-header>

      <a-layout-content class="content">
        <a-row :gutter="16">
          <!-- 通知渠道设置 -->
          <a-col :span="12">
            <a-card title="通知渠道设置" class="settings-card">
              <a-form :model="notifySettings" layout="vertical">
                <a-form-item label="邮件通知">
                  <a-switch v-model="notifySettings.email.enabled" />
                </a-form-item>
                <a-form-item v-if="notifySettings.email.enabled" label="收件邮箱">
                  <a-input v-model="notifySettings.email.recipients" placeholder="多个邮箱用逗号分隔" />
                </a-form-item>
                <a-form-item v-if="notifySettings.email.enabled" label="SMTP 服务器">
                  <a-input v-model="notifySettings.email.smtp_host" placeholder="smtp.example.com" />
                </a-form-item>
                <a-form-item v-if="notifySettings.email.enabled" label="SMTP 端口">
                  <a-input-number v-model="notifySettings.email.smtp_port" :min="1" :max="65535" placeholder="587" style="width: 100%" />
                </a-form-item>

                <a-divider />

                <a-form-item label="短信通知">
                  <a-switch v-model="notifySettings.sms.enabled" />
                </a-form-item>
                <a-form-item v-if="notifySettings.sms.enabled" label="接收手机号">
                  <a-input v-model="notifySettings.sms.recipients" placeholder="多个手机号用逗号分隔" />
                </a-form-item>
                <a-form-item v-if="notifySettings.sms.enabled" label="短信服务商">
                  <a-select v-model="notifySettings.sms.provider">
                    <a-option value="aliyun">阿里云</a-option>
                    <a-option value="tencent">腾讯云</a-option>
                    <a-option value="huawei">华为云</a-option>
                  </a-select>
                </a-form-item>

                <a-divider />

                <a-form-item label="Webhook 通知">
                  <a-switch v-model="notifySettings.webhook.enabled" />
                </a-form-item>
                <a-form-item v-if="notifySettings.webhook.enabled" label="Webhook URL">
                  <a-input v-model="notifySettings.webhook.url" placeholder="https://your-webhook-endpoint.com/alerts" />
                </a-form-item>
                <a-form-item v-if="notifySettings.webhook.enabled" label="Webhook 密钥">
                  <a-input-password v-model="notifySettings.webhook.secret" placeholder="用于签名验证" />
                </a-form-item>

                <a-divider />

                <a-form-item label="应用内通知">
                  <a-switch v-model="notifySettings.inApp.enabled" />
                </a-form-item>
              </a-form>
            </a-card>
          </a-col>

          <!-- 告警策略设置 -->
          <a-col :span="12">
            <a-card title="告警策略设置" class="settings-card">
              <a-form :model="policySettings" layout="vertical">
                <a-form-item label="告警升级">
                  <a-switch v-model="policySettings.escalation.enabled" />
                </a-form-item>
                <a-form-item v-if="policySettings.escalation.enabled" label="升级时间 (分钟)">
                  <a-input-number v-model="policySettings.escalation.timeout" :min="5" :max="1440" style="width: 100%" />
                  <template #extra>告警未处理超过此时间后自动升级</template>
                </a-form-item>
                <a-form-item v-if="policySettings.escalation.enabled" label="升级通知">
                  <a-checkbox-group v-model="policySettings.escalation.notify_levels">
                    <a-checkbox value="manager">主管</a-checkbox>
                    <a-checkbox value="admin">管理员</a-checkbox>
                    <a-checkbox value="all">所有人</a-checkbox>
                  </a-checkbox-group>
                </a-form-item>

                <a-divider />

                <a-form-item label="告警抑制 (相同告警)">
                  <a-switch v-model="policySettings.suppression.enabled" />
                </a-form-item>
                <a-form-item v-if="policySettings.suppression.enabled" label="抑制时间 (分钟)">
                  <a-input-number v-model="policySettings.suppression.minutes" :min="1" :max="1440" style="width: 100%" />
                  <template #extra>相同告警在此时间内不会重复通知</template>
                </a-form-item>

                <a-divider />

                <a-form-item label="告警静默期">
                  <a-switch v-model="policySettings.silentPeriod.enabled" />
                </a-form-item>
                <a-form-item v-if="policySettings.silentPeriod.enabled" label="静默开始时间">
                  <a-time-picker v-model="policySettings.silentPeriod.start" format="HH:mm" style="width: 100%" />
                </a-form-item>
                <a-form-item v-if="policySettings.silentPeriod.enabled" label="静默结束时间">
                  <a-time-picker v-model="policySettings.silentPeriod.end" format="HH:mm" style="width: 100%" />
                </a-form-item>
              </a-form>
            </a-card>
          </a-col>
        </a-row>

        <a-row :gutter="16" style="margin-top: 16px;">
          <!-- 告警阈值设置 -->
          <a-col :span="12">
            <a-card title="告警阈值设置" class="settings-card">
              <a-form :model="thresholdSettings" layout="vertical">
                <a-form-item label="设备离线告警">
                  <a-switch v-model="thresholdSettings.deviceOffline.enabled" />
                </a-form-item>
                <a-form-item v-if="thresholdSettings.deviceOffline.enabled" label="离线时间阈值 (秒)">
                  <a-input-number v-model="thresholdSettings.deviceOffline.threshold" :min="30" :max="3600" style="width: 100%" />
                </a-form-item>

                <a-divider />

                <a-form-item label="CPU 使用率告警">
                  <a-switch v-model="thresholdSettings.cpuUsage.enabled" />
                </a-form-item>
                <a-form-item v-if="thresholdSettings.cpuUsage.enabled" label="CPU 使用率阈值 (%)">
                  <a-input-number v-model="thresholdSettings.cpuUsage.threshold" :min="50" :max="100" style="width: 100%" />
                </a-form-item>

                <a-divider />

                <a-form-item label="内存使用率告警">
                  <a-switch v-model="thresholdSettings.memoryUsage.enabled" />
                </a-form-item>
                <a-form-item v-if="thresholdSettings.memoryUsage.enabled" label="内存使用率阈值 (%)">
                  <a-input-number v-model="thresholdSettings.memoryUsage.threshold" :min="50" :max="100" style="width: 100%" />
                </a-form-item>

                <a-divider />

                <a-form-item label="磁盘使用率告警">
                  <a-switch v-model="thresholdSettings.diskUsage.enabled" />
                </a-form-item>
                <a-form-item v-if="thresholdSettings.diskUsage.enabled" label="磁盘使用率阈值 (%)">
                  <a-input-number v-model="thresholdSettings.diskUsage.threshold" :min="50" :max="100" style="width: 100%" />
                </a-form-item>
              </a-form>
            </a-card>
          </a-col>

          <!-- 其他设置 -->
          <a-col :span="12">
            <a-card title="其他设置" class="settings-card">
              <a-form :model="otherSettings" layout="vertical">
                <a-form-item label="告警历史保留天数">
                  <a-input-number v-model="otherSettings.retentionDays" :min="7" :max="365" style="width: 100%" />
                </a-form-item>
                <a-form-item label="每日告警报告">
                  <a-switch v-model="otherSettings.dailyReport" />
                </a-form-item>
                <a-form-item v-if="otherSettings.dailyReport" label="报告发送时间">
                  <a-time-picker v-model="otherSettings.reportTime" format="HH:mm" style="width: 100%" />
                </a-form-item>
                <a-form-item label="允许移动端查看告警">
                  <a-switch v-model="otherSettings.mobileAccess" />
                </a-form-item>
                <a-form-item label="告警声音提醒">
                  <a-switch v-model="otherSettings.soundAlert" />
                </a-form-item>
              </a-form>
            </a-card>
          </a-col>
        </a-row>

        <!-- 保存按钮 -->
        <a-card class="save-card">
          <a-space>
            <a-button type="primary" @click="saveSettings">保存设置</a-button>
            <a-button @click="resetSettings">重置</a-button>
          </a-space>
        </a-card>
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'
import { Message } from '@arco-design/web-vue'

const router = useRouter()
const collapsed = ref(false)
const selectedKeys = ref(['alerts'])
const loading = ref(false)

const notifySettings = reactive({
  email: { enabled: true, recipients: '', smtp_host: '', smtp_port: 587 },
  sms: { enabled: false, recipients: '', provider: 'aliyun' },
  webhook: { enabled: false, url: '', secret: '' },
  inApp: { enabled: true }
})

const policySettings = reactive({
  escalation: { enabled: false, timeout: 30, notify_levels: ['admin'] },
  suppression: { enabled: true, minutes: 5 },
  silentPeriod: { enabled: false, start: '', end: '' }
})

const thresholdSettings = reactive({
  deviceOffline: { enabled: true, threshold: 90 },
  cpuUsage: { enabled: true, threshold: 80 },
  memoryUsage: { enabled: true, threshold: 85 },
  diskUsage: { enabled: true, threshold: 90 }
})

const otherSettings = reactive({
  retentionDays: 90,
  dailyReport: false,
  reportTime: '',
  mobileAccess: true,
  soundAlert: true
})

const handleMenuClick = ({ key }) => {
  const routes = { dashboard: '/dashboard', alerts: '/alerts/settings' }
  if (routes[key]) router.push(routes[key])
  selectedKeys.value = [key]
}

const loadSettings = async () => {
  loading.value = true
  try {
    const res = await axios.get('/api/v1/alerts/settings')
    const data = res.data
    if (data.code === 0 && data.data) {
      const settings = data.data
      if (settings.notify) Object.assign(notifySettings, settings.notify)
      if (settings.policy) Object.assign(policySettings, settings.policy)
      if (settings.threshold) Object.assign(thresholdSettings, settings.threshold)
      if (settings.other) Object.assign(otherSettings, settings.other)
    }
  } catch (err) {
    Message.error('加载告警设置失败')
  } finally {
    loading.value = false
  }
}

const saveSettings = async () => {
  try {
    const settings = {
      notify: notifySettings,
      policy: policySettings,
      threshold: thresholdSettings,
      other: otherSettings
    }
    const res = await axios.put('/api/v1/alerts/settings', settings)
    if (res.data.code === 0) {
      Message.success('设置保存成功')
    } else {
      Message.error(res.data.message || '保存失败')
    }
  } catch (err) {
    Message.error('保存失败')
  }
}

const resetSettings = () => {
  loadSettings()
  Message.info('已重置为保存的设置')
}

onMounted(() => {
  loadSettings()
})
</script>

<style scoped>
.alert-settings { min-height: 100vh; }
.header { background: #fff; padding: 0 16px; display: flex; align-items: center; gap: 16px; box-shadow: 0 1px 4px rgba(0,0,0,0.1); }
.header-left { display: flex; align-items: center; }
.header-title { font-size: 16px; font-weight: 500; }
.content { padding: 16px; background: #f0f2f5; }
.settings-card { margin-bottom: 16px; }
.save-card { margin-top: 16px; }
</style>
