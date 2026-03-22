<template>
  <div class="page-container">
    <!-- 2FA 状态卡片 -->
    <a-row :gutter="16">
      <a-col :span="8">
        <a-card class="status-card">
          <div class="status-item">
            <div class="status-icon" :class="twoFactorEnabled ? 'enabled' : 'disabled'">
              <icon-shield v-if="twoFactorEnabled" />
              <icon-shield-disable v-else />
            </div>
            <div class="status-info">
              <span class="status-label">双因素认证</span>
              <span class="status-value">{{ twoFactorEnabled ? '已启用' : '未启用' }}</span>
            </div>
          </div>
        </a-card>
      </a-col>
      <a-col :span="8">
        <a-card class="status-card">
          <div class="status-item">
            <div class="status-icon warning">
              <icon-clock />
            </div>
            <div class="status-info">
              <span class="status-label">验证方式</span>
              <span class="status-value">{{ methodLabel(twoFactorMethod) }}</span>
            </div>
          </div>
        </a-card>
      </a-col>
      <a-col :span="8">
        <a-card class="status-card">
          <div class="status-item">
            <div class="status-icon info">
              <icon-calendar />
            </div>
            <div class="status-info">
              <span class="status-label">绑定时间</span>
              <span class="status-value">{{ bindTime || '未绑定' }}</span>
            </div>
          </div>
        </a-card>
      </a-col>
    </a-row>

    <!-- 未启用 2FA：显示启用向导 -->
    <a-card v-if="!twoFactorEnabled" class="setup-card">
      <template #title>
        <div class="card-title">
          <icon-shield />
          <span>启用双因素认证</span>
        </div>
      </template>

      <a-steps :current="setupStep" @change="handleStepChange">
        <a-step title="选择方式" />
        <a-step title="验证身份" />
        <a-step title="绑定成功" />
      </a-steps>

      <div class="step-content">
        <!-- Step 0: 选择验证方式 -->
        <div v-if="setupStep === 0" class="method-select">
          <a-radio-group v-model="selectedMethod" direction="vertical">
            <a-radio value="totp">
              <div class="method-option">
                <icon-qrcode class="method-icon" />
                <div>
                  <div class="method-name">身份验证器应用</div>
                  <div class="method-desc">使用 Google Authenticator、Microsoft Authenticator 等应用扫描二维码</div>
                </div>
              </div>
            </a-radio>
            <a-radio value="sms">
              <div class="method-option">
                <icon-message class="method-icon" />
                <div>
                  <div class="method-name">短信验证码</div>
                  <div class="method-desc">通过短信接收一次性验证码</div>
                </div>
              </div>
            </a-radio>
            <a-radio value="email">
              <div class="method-option">
                <icon-email class="method-icon" />
                <div>
                  <div class="method-name">邮件验证码</div>
                  <div class="method-desc">通过注册邮箱接收一次性验证码</div>
                </div>
              </div>
            </a-radio>
          </a-radio-group>
          <a-button type="primary" style="margin-top: 24px" :disabled="!selectedMethod" @click="setupStep = 1">
            下一步
          </a-button>
        </div>

        <!-- Step 1: TOTP - 显示二维码 -->
        <div v-if="setupStep === 1 && selectedMethod === 'totp'" class="verify-section">
          <a-alert type="info" style="margin-bottom: 16px">
            请使用身份验证器应用扫描下方二维码，然后在验证框中输入显示的 6 位数字。
          </a-alert>
          <a-row :gutter="32" align="center">
            <a-col :span="12" style="text-align: center">
              <div class="qr-wrapper">
                <img v-if="qrCodeUrl" :src="qrCodeUrl" alt="2FA QR Code" class="qr-image" />
                <div v-else class="qr-placeholder">
                  <a-spin />
                </div>
              </div>
              <div class="secret-key" v-if="totpSecret">
                密钥：<span class="mono">{{ totpSecret }}</span>
                <a-button type="text" size="small" @click="copySecret">
                  <icon-copy />
                </a-button>
              </div>
            </a-col>
            <a-col :span="12">
              <a-form :model="verifyForm" layout="vertical">
                <a-form-item label="请输入身份验证器中的 6 位数字" required>
                  <a-input
                    v-model="verifyForm.code"
                    placeholder="000000"
                    :maxlength="6"
                    style="font-size: 20px; letter-spacing: 8px; text-align: center"
                  />
                </a-form-item>
                <a-form-item>
                  <a-space>
                    <a-button @click="setupStep = 0">上一步</a-button>
                    <a-button type="primary" :loading="verifying" @click="verifyAndEnable">
                      验证并启用
                    </a-button>
                  </a-space>
                </a-form-item>
              </a-form>
            </a-col>
          </a-row>
        </div>

        <!-- Step 1: SMS - 发送验证码 -->
        <div v-if="setupStep === 1 && selectedMethod === 'sms'" class="verify-section">
          <a-alert type="info" style="margin-bottom: 16px">
            我们将向您的手机号码发送验证码。
          </a-alert>
          <a-form :model="smsForm" layout="vertical">
            <a-form-item label="手机号码">
              <a-input v-model="smsForm.phone" placeholder="请输入手机号码" />
            </a-form-item>
            <a-form-item label="验证码">
              <a-input-search
                v-model="smsForm.code"
                placeholder="请输入验证码"
                :maxlength="6"
                style="width: 200px"
                @search="sendSmsCode"
              >
                <template #search-button>
                  <a-button type="primary" :disabled="smsCooldown > 0" @click="sendSmsCode">
                    {{ smsCooldown > 0 ? `${smsCooldown}s` : '发送验证码' }}
                  </a-button>
                </template>
              </a-input-search>
            </a-form-item>
            <a-form-item>
              <a-space>
                <a-button @click="setupStep = 0">上一步</a-button>
                <a-button type="primary" :loading="verifying" @click="verifyAndEnableSms">
                  验证并启用
                </a-button>
              </a-space>
            </a-form-item>
          </a-form>
        </div>

        <!-- Step 1: Email - 发送验证码 -->
        <div v-if="setupStep === 1 && selectedMethod === 'email'" class="verify-section">
          <a-alert type="info" style="margin-bottom: 16px">
            我们将向您的注册邮箱发送验证码。
          </a-alert>
          <a-form :model="emailForm" layout="vertical">
            <a-form-item label="邮箱地址">
              <a-input v-model="emailForm.email" placeholder="请输入邮箱地址" />
            </a-form-item>
            <a-form-item label="验证码">
              <a-input-search
                v-model="emailForm.code"
                placeholder="请输入验证码"
                :maxlength="6"
                style="width: 200px"
                @search="sendEmailCode"
              >
                <template #search-button>
                  <a-button type="primary" :disabled="emailCooldown > 0" @click="sendEmailCode">
                    {{ emailCooldown > 0 ? `${emailCooldown}s` : '发送验证码' }}
                  </a-button>
                </template>
              </a-input-search>
            </a-form-item>
            <a-form-item>
              <a-space>
                <a-button @click="setupStep = 0">上一步</a-button>
                <a-button type="primary" :loading="verifying" @click="verifyAndEnableEmail">
                  验证并启用
                </a-button>
              </a-space>
            </a-form-item>
          </a-form>
        </div>

        <!-- Step 2: 成功 -->
        <div v-if="setupStep === 2" class="success-section">
          <a-result status="success" title="双因素认证已启用" subtitle="您的账户已获得额外的安全保护">
            <template #extra>
              <a-button type="primary" @click="setupStep = 0">完成</a-button>
            </template>
          </a-result>
        </div>
      </div>
    </a-card>

    <!-- 已启用 2FA：显示管理界面 -->
    <a-card v-else class="manage-card">
      <template #title>
        <div class="card-title">
          <icon-shield />
          <span>双因素认证管理</span>
        </div>
      </template>

      <!-- 恢复代码 -->
      <a-card class="sub-card">
        <template #title>
          <div class="sub-title">
            <icon-safe />
            <span>恢复代码</span>
          </div>
        </template>
        <a-alert type="warning" style="margin-bottom: 12px">
          请将恢复代码保存在安全的地方。如果您的身份验证器不可用，可以使用恢复代码访问您的账户。
        </a-alert>
        <div class="recovery-codes" v-if="recoveryCodes.length">
          <a-tag v-for="code in recoveryCodes" :key="code" class="recovery-code">{{ code }}</a-tag>
        </div>
        <div v-else class="no-codes">暂无恢复代码</div>
        <a-button type="outline" style="margin-top: 12px" :loading="regeneratingCodes" @click="regenerateCodes">
          <template #icon><icon-refresh /></template>
          重新生成恢复代码
        </a-button>
      </a-card>

      <!-- 信任设备管理 -->
      <a-card class="sub-card">
        <template #title>
          <div class="sub-title">
            <icon-computer />
            <span>信任设备</span>
          </div>
        </template>
        <a-table
          :columns="trustColumns"
          :data="trustDevices"
          :loading="loadingTrust"
          row-key="id"
          :pagination="false"
        >
          <template #device="{ record }">
            <div class="device-info">
              <icon-computer class="device-icon" />
              <div>
                <div>{{ record.device_name || '未知设备' }}</div>
                <div class="device-browser">{{ record.browser || '' }} / {{ record.os || '' }}</div>
              </div>
            </div>
          </template>
          <template #location="{ record }">
            <span>{{ record.location || '-' }}</span>
          </template>
          <template #lastActive="{ record }">
            <span>{{ formatTime(record.last_active) }}</span>
          </template>
          <template #actions="{ record }">
            <a-popconfirm content="确定移除该信任设备？" @ok="removeTrustDevice(record)">
              <a-button type="text" size="small" status="danger">移除</a-button>
            </a-popconfirm>
          </template>
        </a-table>
        <a-button type="text" style="margin-top: 8px" @click="loadTrustDevices">
          <icon-refresh />
          刷新
        </a-button>
      </a-card>

      <!-- 更换验证方式 -->
      <a-card class="sub-card">
        <template #title>
          <div class="sub-title">
            <icon-swap />
            <span>更换验证方式</span>
          </div>
        </template>
        <a-form :model="changeMethodForm" layout="vertical" style="max-width: 400px">
          <a-form-item label="新的验证方式">
            <a-select v-model="changeMethodForm.newMethod" placeholder="请选择新的验证方式">
              <a-option value="totp">身份验证器</a-option>
              <a-option value="sms">短信验证码</a-option>
              <a-option value="email">邮件验证码</a-option>
            </a-select>
          </a-form-item>
          <a-form-item label="当前密码" required>
            <a-input-password v-model="changeMethodForm.password" placeholder="请输入当前密码" />
          </a-form-item>
          <a-form-item>
            <a-button type="primary" :loading="changingMethod" @click="changeMethod">
              更换验证方式
            </a-button>
          </a-form-item>
        </a-form>
      </a-card>

      <!-- 禁用 2FA -->
      <a-card class="sub-card danger">
        <template #title>
          <div class="sub-title danger">
            <icon-warning />
            <span>禁用双因素认证</span>
          </div>
        </template>
        <a-alert type="error" style="margin-bottom: 12px">
          禁用双因素认证将降低您账户的安全性。请确保您了解相关风险。
        </a-alert>
        <a-form :model="disableForm" layout="vertical" style="max-width: 400px">
          <a-form-item label="当前密码" required>
            <a-input-password v-model="disableForm.password" placeholder="请输入当前密码" />
          </a-form-item>
          <a-form-item label="验证码" required>
            <a-input v-model="disableForm.code" placeholder="请输入 2FA 验证码" :maxlength="6" />
          </a-form-item>
          <a-form-item>
            <a-popconfirm content="确定要禁用双因素认证吗？" @ok="disable2FA">
              <a-button type="primary" status="danger" :loading="disabling" :disabled="!disableForm.password || !disableForm.code">
                禁用双因素认证
              </a-button>
            </a-popconfirm>
          </a-form-item>
        </a-form>
      </a-card>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import {
  getTwoFactorStatus,
  enableTwoFactor,
  disableTwoFactor,
  verifyTwoFactorCode,
  getRecoveryCodes,
  regenerateRecoveryCodes,
  getTrustDevices,
  removeTrustDevice as removeTrustDeviceApi,
  changeTwoFactorMethod,
  sendSmsVerifyCode,
  sendEmailVerifyCode,
  getTotpQrCode
} from '@/api/security'

const twoFactorEnabled = ref(false)
const twoFactorMethod = ref('')
const bindTime = ref('')
const setupStep = ref(0)
const selectedMethod = ref('')
const verifying = ref(false)
const disabling = ref(false)
const changingMethod = ref(false)
const regeneratingCodes = ref(false)
const loadingTrust = ref(false)

const totpSecret = ref('')
const qrCodeUrl = ref('')

const recoveryCodes = ref([])
const trustDevices = ref([])

const verifyForm = reactive({ code: '' })
const smsForm = reactive({ phone: '', code: '' })
const emailForm = reactive({ email: '', code: '' })
const disableForm = reactive({ password: '', code: '' })
const changeMethodForm = reactive({ newMethod: '', password: '' })

let smsTimer = null
let emailTimer = null
const smsCooldown = ref(0)
const emailCooldown = ref(0)

const trustColumns = [
  { title: '设备', slotName: 'device', minWidth: 200 },
  { title: '位置', slotName: 'location', minWidth: 120 },
  { title: '最后活动', slotName: 'lastActive', minWidth: 150 },
  { title: '操作', slotName: 'actions', width: 100 }
]

function methodLabel(method) {
  const map = { totp: '身份验证器', sms: '短信', email: '邮件' }
  return map[method] || method || '未设置'
}

function formatTime(timestamp) {
  if (!timestamp) return '-'
  return new Date(timestamp).toLocaleString('zh-CN')
}

function handleStepChange(step) {
  setupStep.value = step
}

async function loadStatus() {
  try {
    const data = await getTwoFactorStatus()
    const res = data.data || data
    twoFactorEnabled.value = res.enabled
    twoFactorMethod.value = res.method || ''
    bindTime.value = res.bind_time ? formatTime(res.bind_time) : ''
  } catch (e) {
    console.error('加载 2FA 状态失败', e)
  }
}

async function loadTrustDevices() {
  loadingTrust.value = true
  try {
    const data = await getTrustDevices()
    trustDevices.value = (data.data || data).list || (data.data || data) || []
  } catch (e) {
    console.error('加载信任设备失败', e)
  } finally {
    loadingTrust.value = false
  }
}

async function loadRecoveryCodes() {
  try {
    const data = await getRecoveryCodes()
    recoveryCodes.value = (data.data || data).codes || []
  } catch (e) {
    console.error('加载恢复代码失败', e)
  }
}

async function verifyAndEnable() {
  if (!verifyForm.code || verifyForm.code.length < 6) {
    Message.warning('请输入 6 位验证码')
    return
  }
  verifying.value = true
  try {
    await verifyTwoFactorCode({ code: verifyForm.code })
    await enableTwoFactor({ method: 'totp', code: verifyForm.code })
    Message.success('双因素认证已启用')
    setupStep.value = 2
    twoFactorEnabled.value = true
    twoFactorMethod.value = 'totp'
    await loadRecoveryCodes()
    await loadTrustDevices()
    bindTime.value = formatTime(Date.now())
  } catch (e) {
    Message.error('验证失败: ' + (e.message || '请检查验证码是否正确'))
  } finally {
    verifying.value = false
  }
}

async function sendSmsCode() {
  if (!smsForm.phone) {
    Message.warning('请输入手机号码')
    return
  }
  try {
    await sendSmsVerifyCode({ phone: smsForm.phone })
    Message.success('验证码已发送')
    smsCooldown.value = 60
    smsTimer = setInterval(() => {
      smsCooldown.value--
      if (smsCooldown.value <= 0) clearInterval(smsTimer)
    }, 1000)
  } catch (e) {
    Message.error('发送失败')
  }
}

async function verifyAndEnableSms() {
  if (!smsForm.phone || !smsForm.code || smsForm.code.length < 6) {
    Message.warning('请填写完整信息')
    return
  }
  verifying.value = true
  try {
    await enableTwoFactor({ method: 'sms', phone: smsForm.phone, code: smsForm.code })
    Message.success('双因素认证已启用')
    setupStep.value = 2
    twoFactorEnabled.value = true
    twoFactorMethod.value = 'sms'
    bindTime.value = formatTime(Date.now())
    await loadRecoveryCodes()
  } catch (e) {
    Message.error('验证失败')
  } finally {
    verifying.value = false
  }
}

async function sendEmailCode() {
  if (!emailForm.email) {
    Message.warning('请输入邮箱地址')
    return
  }
  try {
    await sendEmailVerifyCode({ email: emailForm.email })
    Message.success('验证码已发送')
    emailCooldown.value = 60
    emailTimer = setInterval(() => {
      emailCooldown.value--
      if (emailCooldown.value <= 0) clearInterval(emailTimer)
    }, 1000)
  } catch (e) {
    Message.error('发送失败')
  }
}

async function verifyAndEnableEmail() {
  if (!emailForm.email || !emailForm.code || emailForm.code.length < 6) {
    Message.warning('请填写完整信息')
    return
  }
  verifying.value = true
  try {
    await enableTwoFactor({ method: 'email', email: emailForm.email, code: emailForm.code })
    Message.success('双因素认证已启用')
    setupStep.value = 2
    twoFactorEnabled.value = true
    twoFactorMethod.value = 'email'
    bindTime.value = formatTime(Date.now())
    await loadRecoveryCodes()
  } catch (e) {
    Message.error('验证失败')
  } finally {
    verifying.value = false
  }
}

async function regenerateCodes() {
  regeneratingCodes.value = true
  try {
    const data = await regenerateRecoveryCodes()
    recoveryCodes.value = (data.data || data).codes || []
    Message.success('恢复代码已重新生成')
  } catch (e) {
    Message.error('生成失败')
  } finally {
    regeneratingCodes.value = false
  }
}

async function removeTrustDevice(record) {
  try {
    await removeTrustDeviceApi(record.id)
    Message.success('已移除')
    loadTrustDevices()
  } catch (e) {
    Message.error('移除失败')
  }
}

async function changeMethod() {
  if (!changeMethodForm.newMethod || !changeMethodForm.password) {
    Message.warning('请填写完整信息')
    return
  }
  changingMethod.value = true
  try {
    await changeTwoFactorMethod(changeMethodForm)
    Message.success('验证方式已更换')
    changeMethodForm.newMethod = ''
    changeMethodForm.password = ''
    await loadStatus()
  } catch (e) {
    Message.error('更换失败')
  } finally {
    changingMethod.value = false
  }
}

async function disable2FA() {
  if (!disableForm.password || !disableForm.code) {
    Message.warning('请填写完整信息')
    return
  }
  disabling.value = true
  try {
    await disableTwoFactor({ password: disableForm.password, code: disableForm.code })
    Message.success('双因素认证已禁用')
    disableForm.password = ''
    disableForm.code = ''
    twoFactorEnabled.value = false
    twoFactorMethod.value = ''
    bindTime.value = ''
    recoveryCodes.value = []
    trustDevices.value = []
  } catch (e) {
    Message.error('禁用失败: ' + (e.message || ''))
  } finally {
    disabling.value = false
  }
}

function copySecret() {
  navigator.clipboard.writeText(totpSecret.value)
  Message.success('密钥已复制')
}

// Load TOTP QR code when switching to TOTP method
import { watch } from 'vue'
watch([setupStep, selectedMethod], async ([step, method]) => {
  if (step === 1 && method === 'totp' && !qrCodeUrl.value) {
    try {
      const data = await getTotpQrCode()
      const res = data.data || data
      totpSecret.value = res.secret || ''
      qrCodeUrl.value = res.qr_code_url || res.url || ''
    } catch (e) {
      console.error('获取二维码失败', e)
    }
  }
})

onMounted(async () => {
  await loadStatus()
  if (twoFactorEnabled.value) {
    await loadRecoveryCodes()
    await loadTrustDevices()
  }
})

onUnmounted(() => {
  if (smsTimer) clearInterval(smsTimer)
  if (emailTimer) clearInterval(emailTimer)
})
</script>

<style scoped>
.page-container {
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.status-card {
  cursor: default;
}

.status-item {
  display: flex;
  align-items: center;
  gap: 16px;
}

.status-icon {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
}

.status-icon.enabled {
  background: rgba(22, 93, 255, 0.1);
  color: rgb(var(--primary-6));
}

.status-icon.disabled {
  background: rgba(var(--gray-3));
  color: var(--color-text-3);
}

.status-icon.warning {
  background: rgba(240, 152, 0, 0.1);
  color: rgb(240, 152, 0);
}

.status-icon.info {
  background: rgba(var(--arcoblue-1));
  color: rgb(var(--arcoblue-6));
}

.status-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.status-label {
  font-size: 13px;
  color: var(--color-text-3);
}

.status-value {
  font-size: 16px;
  font-weight: 600;
}

.card-title {
  display: flex;
  align-items: center;
  gap: 8px;
}

.setup-card,
.manage-card {
  flex-shrink: 0;
}

.step-content {
  margin-top: 32px;
  min-height: 300px;
}

.method-option {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 12px 0;
}

.method-icon {
  font-size: 32px;
  color: rgb(var(--primary-6));
}

.method-name {
  font-size: 15px;
  font-weight: 500;
}

.method-desc {
  font-size: 13px;
  color: var(--color-text-3);
  margin-top: 4px;
}

.qr-wrapper {
  display: inline-block;
  padding: 16px;
  background: white;
  border: 1px solid var(--color-border);
  border-radius: 8px;
}

.qr-image {
  width: 200px;
  height: 200px;
}

.qr-placeholder {
  width: 200px;
  height: 200px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.secret-key {
  margin-top: 12px;
  font-size: 13px;
  color: var(--color-text-3);
}

.mono {
  font-family: monospace;
  background: var(--color-fill-1);
  padding: 2px 6px;
  border-radius: 4px;
}

.success-section {
  padding: 40px 0;
}

.sub-card {
  margin-bottom: 16px;
}

.sub-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
}

.sub-title.danger {
  color: rgb(var(--danger-6));
}

.recovery-codes {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.recovery-code {
  font-family: monospace;
  font-size: 14px;
  letter-spacing: 2px;
}

.no-codes {
  color: var(--color-text-3);
  font-size: 13px;
}

.device-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.device-icon {
  font-size: 24px;
  color: var(--color-text-3);
}

.device-browser {
  font-size: 12px;
  color: var(--color-text-3);
}

.verify-section,
.method-select {
  max-width: 600px;
}
</style>
