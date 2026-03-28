<template>
  <div class="settings-container">
    <a-card>
      <template #title>
        <span>系统设置</span>
      </template>
      
      <a-tabs default-active-key="general">
        <a-tab-pane key="general" title="基本设置">
          <a-form :model="form" layout="vertical">
            <a-form-item label="系统名称">
              <a-input v-model="form.systemName" />
            </a-form-item>
            <a-form-item label="系统Logo">
              <a-input v-model="form.logo" />
            </a-form-item>
            <a-form-item label="默认时区">
              <a-select v-model="form.timezone">
                <a-option value="Asia/Shanghai">Asia/Shanghai (UTC+8)</a-option>
                <a-option value="UTC">UTC</a-option>
              </a-select>
            </a-form-item>
            <a-form-item>
              <a-button type="primary">保存</a-button>
            </a-form-item>
          </a-form>
        </a-tab-pane>
        <a-tab-pane key="notification" title="通知设置">
          <a-form :model="form" layout="vertical">
            <a-form-item label="邮件通知">
              <a-switch v-model="form.emailEnabled" />
            </a-form-item>
            <a-form-item label="SMS通知">
              <a-switch v-model="form.smsEnabled" />
            </a-form-item>
            <a-form-item label="Webhook通知">
              <a-switch v-model="form.webhookEnabled" />
            </a-form-item>
            <a-form-item>
              <a-button type="primary">保存</a-button>
            </a-form-item>
          </a-form>
        </a-tab-pane>
        <a-tab-pane key="security" title="安全设置">
          <a-form layout="vertical">
            <a-form-item label="会话超时(分钟)">
              <a-input-number v-model="form.sessionTimeout" :min="5" :max="1440" />
            </a-form-item>
            <a-form-item label="密码强度要求">
              <a-checkbox-group v-model="form.passwordRequirements">
                <a-checkbox value="length">至少8位</a-checkbox>
                <a-checkbox value="uppercase">包含大写字母</a-checkbox>
                <a-checkbox value="number">包含数字</a-checkbox>
              </a-checkbox-group>
            </a-form-item>
            <a-form-item>
              <a-button type="primary">保存</a-button>
            </a-form-item>
          </a-form>
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { reactive } from 'vue';

const form = reactive({
  systemName: 'MDM 控制中台',
  logo: '/logo.png',
  timezone: 'Asia/Shanghai',
  emailEnabled: true,
  smsEnabled: true,
  webhookEnabled: false,
  sessionTimeout: 60,
  passwordRequirements: ['length', 'number'],
});
</script>

<style scoped>
.settings-container { padding: 20px; }
</style>
