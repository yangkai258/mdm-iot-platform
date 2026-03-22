<template>
  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>家庭管理</a-breadcrumb-item>
      <a-breadcrumb-item>家庭设置</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 页面标题 -->
    <div class="pro-page-header">
      <h2 class="pro-page-title">家庭设置</h2>
      <p class="pro-page-desc">管理家庭基本信息、通知偏好和隐私设置</p>
    </div>

    <!-- 内容区 -->
    <div class="pro-content-area">
      <a-tabs v-model:active-key="activeTab" type="line" :animated="false">
        <!-- 基本信息 -->
        <a-tab-pane key="basic" title="基本信息">
          <a-card class="settings-card">
            <a-form :model="basicForm" layout="vertical" @submit="saveBasic">
              <a-form-item label="家庭名称" required>
                <a-input v-model="basicForm.name" placeholder="请输入家庭名称" style="max-width: 400px" />
              </a-form-item>
              <a-form-item label="家庭描述">
                <a-textarea v-model="basicForm.description" placeholder="请输入家庭描述" :rows="3" style="max-width: 500px" />
              </a-form-item>
              <a-form-item label="家庭地址">
                <a-input v-model="basicForm.address" placeholder="请输入家庭地址" style="max-width: 500px" />
              </a-form-item>
              <a-form-item label="家庭头像">
                <div class="avatar-upload">
                  <a-avatar :size="80" :style="{ backgroundColor: '#1650ff' }">
                    {{ basicForm.name?.charAt(0) || '?' }}
                  </a-avatar>
                  <a-button style="margin-left: 16px">更换头像</a-button>
                </div>
              </a-form-item>
              <a-form-item>
                <a-button type="primary" html-type="submit" :loading="saving.basic">保存修改</a-button>
              </a-form-item>
            </a-form>
          </a-card>
        </a-tab-pane>

        <!-- 通知设置 -->
        <a-tab-pane key="notification" title="通知设置">
          <a-card class="settings-card">
            <a-form :model="notifForm" layout="vertical">
              <a-form-item label="设备状态变更通知">
                <a-switch v-model="notifForm.device_status_notify" />
              </a-form-item>
              <a-form-item label="成员加入通知">
                <a-switch v-model="notifForm.member_join_notify" />
              </a-form-item>
              <a-form-item label="儿童使用报告">
                <a-switch v-model="notifForm.child_usage_notify" />
                <div v-if="notifForm.child_usage_notify" style="margin-top:8px">
                  <a-select v-model="notifForm.child_report_frequency" style="width:200px">
                    <a-option value="daily">每日报告</a-option>
                    <a-option value="weekly">每周报告</a-option>
                    <a-option value="monthly">每月报告</a-option>
                  </a-select>
                </div>
              </a-form-item>
              <a-form-item label="老人陪伴模式提醒">
                <a-switch v-model="notifForm.elder_reminder_notify" />
              </a-form-item>
              <a-form-item label="相册更新通知">
                <a-switch v-model="notifForm.album_update_notify" />
              </a-form-item>
              <a-form-item label="固件升级提醒">
                <a-switch v-model="notifForm.ota_notify" />
              </a-form-item>
              <a-form-item label="通知方式">
                <a-checkbox-group v-model="notifForm.notify_channels">
                  <a-checkbox value="app">APP推送</a-checkbox>
                  <a-checkbox value="sms">短信</a-checkbox>
                  <a-checkbox value="email">邮件</a-checkbox>
                </a-checkbox-group>
              </a-form-item>
              <a-form-item>
                <a-button type="primary" @click="saveNotification" :loading="saving.notif">保存设置</a-button>
              </a-form-item>
            </a-form>
          </a-card>
        </a-tab-pane>

        <!-- 隐私设置 -->
        <a-tab-pane key="privacy" title="隐私设置">
          <a-card class="settings-card">
            <a-form :model="privacyForm" layout="vertical">
              <a-divider>位置信息</a-divider>
              <a-form-item label="共享设备位置">
                <a-switch v-model="privacyForm.share_location" />
              </a-form-item>
              <a-form-item label="位置更新频率" v-if="privacyForm.share_location">
                <a-select v-model="privacyForm.location_update_interval" style="width:200px">
                  <a-option value="realtime">实时</a-option>
                  <a-option value="5min">每5分钟</a-option>
                  <a-option value="15min">每15分钟</a-option>
                  <a-option value="1hour">每小时</a-option>
                </a-select>
              </a-form-item>

              <a-divider>数据共享</a-divider>
              <a-form-item label="允许数据统计">
                <a-switch v-model="privacyForm.allow_analytics" />
              </a-form-item>
              <a-form-item label="发送使用诊断数据">
                <a-switch v-model="privacyForm.send_diagnostics" />
              </a-form-item>

              <a-divider>相册隐私</a-divider>
              <a-form-item label="新照片需审核">
                <a-switch v-model="privacyForm.photo_review_required" />
                <div style="margin-top:4px;color:var(--color-text-3, #86909c);font-size:12px">
                  启用后，新上传的照片需要户主审核才能对所有成员可见
                </div>
              </a-form-item>
              <a-form-item label="照片自动标注人脸">
                <a-switch v-model="privacyForm.auto_face_tag" />
              </a-form-item>

              <a-divider>账户安全</a-divider>
              <a-form-item label="允许其他成员邀请">
                <a-switch v-model="privacyForm.allow_member_invite" />
              </a-form-item>
              <a-form-item label="双因素认证">
                <a-switch v-model="privacyForm.two_factor_enabled" />
              </a-form-item>
              <a-form-item label="登录设备管理">
                <a-button @click="showDeviceList">管理已登录设备</a-button>
              </a-form-item>

              <a-form-item>
                <a-button type="primary" @click="savePrivacy" :loading="saving.privacy">保存设置</a-button>
              </a-form-item>
            </a-form>
          </a-card>
        </a-tab-pane>
      </a-tabs>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const activeTab = ref('basic')
const saving = reactive({ basic: false, notif: false, privacy: false })

const basicForm = reactive({
  name: '',
  description: '',
  address: '',
  avatar_url: ''
})

const notifForm = reactive({
  device_status_notify: true,
  member_join_notify: true,
  child_usage_notify: false,
  child_report_frequency: 'weekly',
  elder_reminder_notify: true,
  album_update_notify: true,
  ota_notify: true,
  notify_channels: ['app'] as string[]
})

const privacyForm = reactive({
  share_location: false,
  location_update_interval: '15min',
  allow_analytics: false,
  send_diagnostics: false,
  photo_review_required: false,
  auto_face_tag: true,
  allow_member_invite: true,
  two_factor_enabled: false
})

async function loadSettings() {
  try {
    const res = await fetch('/api/v1/family/settings', { credentials: 'include' })
    const data = await res.json()
    if (data.code === 0 || data.code === 200) {
      const s = data.data || {}
      basicForm.name = s.name || ''
      basicForm.description = s.description || ''
      basicForm.address = s.address || ''
      basicForm.avatar_url = s.avatar_url || ''

      if (s.notification) {
        Object.assign(notifForm, s.notification)
      }
      if (s.privacy) {
        Object.assign(privacyForm, s.privacy)
      }
    }
  } catch {
    Message.error('加载设置失败')
  }
}

async function saveBasic() {
  saving.basic = true
  try {
    const res = await fetch('/api/v1/family/settings/basic', {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify(basicForm)
    })
    const data = await res.json()
    if (data.code === 0 || data.code === 200) {
      Message.success('基本信息已保存')
    } else {
      Message.error(data.message || '保存失败')
    }
  } catch {
    Message.error('网络错误')
  } finally {
    saving.basic = false
  }
}

async function saveNotification() {
  saving.notif = true
  try {
    const res = await fetch('/api/v1/family/settings/notification', {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify(notifForm)
    })
    const data = await res.json()
    if (data.code === 0 || data.code === 200) {
      Message.success('通知设置已保存')
    } else {
      Message.error(data.message || '保存失败')
    }
  } catch {
    Message.error('网络错误')
  } finally {
    saving.notif = false
  }
}

async function savePrivacy() {
  saving.privacy = true
  try {
    const res = await fetch('/api/v1/family/settings/privacy', {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify(privacyForm)
    })
    const data = await res.json()
    if (data.code === 0 || data.code === 200) {
      Message.success('隐私设置已保存')
    } else {
      Message.error(data.message || '保存失败')
    }
  } catch {
    Message.error('网络错误')
  } finally {
    saving.privacy = false
  }
}

function showDeviceList() {
  Message.info('设备管理功能开发中')
}

onMounted(() => {
  loadSettings()
})
</script>

<style scoped>
.settings-card {
  max-width: 700px;
}

.avatar-upload {
  display: flex;
  align-items: center;
}
</style>
