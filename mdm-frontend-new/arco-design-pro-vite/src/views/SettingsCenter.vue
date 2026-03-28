<template>
  <div class="settings-center-container">
    <a-card>
      <template #title>
        <span>系统设置</span>
      </template>
      
      <a-tabs>
        <a-tab-pane key="general" title="基本设置">
          <a-form layout="vertical" style="max-width: 600px;">
            <a-form-item label="系统名称">
              <a-input v-model="settings.systemName" />
            </a-form-item>
            <a-form-item label="系统Logo">
              <a-upload action="#" :limit="1" />
            </a-form-item>
            <a-form-item label="系统描述">
              <a-textarea v-model="settings.description" :rows="3" />
            </a-form-item>
            <a-form-item label="时区">
              <a-select v-model="settings.timezone">
                <a-option value="Asia/Shanghai">Asia/Shanghai (UTC+8)</a-option>
                <a-option value="America/New_York">America/New_York (UTC-5)</a-option>
              </a-select>
            </a-form-item>
            <a-form-item>
              <a-button type="primary" @click="handleSaveGeneral">保存</a-button>
            </a-form-item>
          </a-form>
        </a-tab-pane>
        
        <a-tab-pane key="notification" title="通知设置">
          <a-form layout="vertical" style="max-width: 600px;">
            <a-form-item label="邮件通知">
              <a-switch v-model="settings.emailNotification" />
            </a-form-item>
            <a-form-item label="邮件服务器">
              <a-input v-model="settings.smtpHost" placeholder="smtp.example.com" />
            </a-form-item>
            <a-form-item label="端口">
              <a-input-number v-model="settings.smtpPort" :min="1" :max="65535" />
            </a-form-item>
            <a-divider />
            <a-form-item label="短信通知">
              <a-switch v-model="settings.smsNotification" />
            </a-form-item>
            <a-form-item label="短信服务商">
              <a-select v-model="settings.smsProvider">
                <a-option value="aliyun">阿里云</a-option>
                <a-option value="tencent">腾讯云</a-option>
              </a-select>
            </a-form-item>
            <a-form-item>
              <a-button type="primary" @click="handleSaveNotification">保存</a-button>
            </a-form-item>
          </a-form>
        </a-tab-pane>
        
        <a-tab-pane key="security" title="安全设置">
          <a-form layout="vertical" style="max-width: 600px;">
            <a-form-item label="登录尝试限制">
              <a-input-number v-model="settings.maxLoginAttempts" :min="1" />
            </a-form-item>
            <a-form-item label="密码最小长度">
              <a-input-number v-model="settings.minPasswordLength" :min="6" />
            </a-form-item>
            <a-form-item label="会话超时时间">
              <a-input-number v-model="settings.sessionTimeout" :min="5" /> 分钟
            </a-form-item>
            <a-form-item label="双因素认证">
              <a-switch v-model="settings.twoFactorAuth" />
            </a-form-item>
            <a-form-item>
              <a-button type="primary" @click="handleSaveSecurity">保存</a-button>
            </a-form-item>
          </a-form>
        </a-tab-pane>
        
        <a-tab-pane key="backup" title="备份设置">
          <a-form layout="vertical" style="max-width: 600px;">
            <a-form-item label="自动备份">
              <a-switch v-model="settings.autoBackup" />
            </a-form-item>
            <a-form-item label="备份间隔">
              <a-select v-model="settings.backupInterval">
                <a-option value="daily">每天</a-option>
                <a-option value="weekly">每周</a-option>
                <a-option value="monthly">每月</a-option>
              </a-select>
            </a-form-item>
            <a-form-item label="保留备份数量">
              <a-input-number v-model="settings.backupRetention" :min="1" />
            </a-form-item>
            <a-form-item label="备份存储位置">
              <a-input v-model="settings.backupPath" />
            </a-form-item>
            <a-form-item>
              <a-space>
                <a-button type="primary" @click="handleBackupNow">立即备份</a-button>
                <a-button @click="handleRestore">从备份恢复</a-button>
              </a-space>
            </a-form-item>
          </a-form>
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const settings = reactive({
  systemName: 'MDM智能宠物平台',
  description: 'AI电子宠物MDM控制中台',
  timezone: 'Asia/Shanghai',
  emailNotification: true,
  smtpHost: 'smtp.example.com',
  smtpPort: 587,
  smsNotification: false,
  smsProvider: 'aliyun',
  maxLoginAttempts: 5,
  minPasswordLength: 8,
  sessionTimeout: 30,
  twoFactorAuth: false,
  autoBackup: true,
  backupInterval: 'daily',
  backupRetention: 7,
  backupPath: '/backup',
});

const handleSaveGeneral = () => {};
const handleSaveNotification = () => {};
const handleSaveSecurity = () => {};
const handleBackupNow = () => {};
const handleRestore = () => {};
</script>

<style scoped>
.settings-center-container { padding: 20px; }
</style>
