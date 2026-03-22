<template>
  <div class="page-container">
    <a-card class="content-card">
      <template #title>
        <div class="card-header">
          <span>LDAP 配置</span>
          <a-tag v-if="ldapConfig.is_enabled" color="green">已启用</a-tag>
          <a-tag v-else color="gray">未启用</a-tag>
        </div>
      </template>

      <a-form :model="form" layout="vertical" class="ldap-form">
        <div class="form-section">
          <div class="section-label">服务器配置</div>
          <a-form-item label="服务器地址" required>
            <a-input
              v-model="form.server_url"
              placeholder="ldap://ldap.example.com"
              style="width: 360px"
            />
          </a-form-item>
          <a-form-item label="端口" required>
            <a-input-number
              v-model="form.port"
              :min="1"
              :max="65535"
              style="width: 120px"
            />
            <span class="form-tip">389（LDAP）或 636（LDAPS）</span>
          </a-form-item>
          <a-form-item label="使用 SSL">
            <a-switch v-model="form.use_ssl" />
          </a-form-item>
        </div>

        <a-divider />

        <div class="form-section">
          <div class="section-label">认证配置</div>
          <a-form-item label="Base DN" required>
            <a-input
              v-model="form.base_dn"
              placeholder="dc=example,dc=com"
              style="width: 360px"
            />
          </a-form-item>
          <a-form-item label="Bind DN">
            <a-input
              v-model="form.bind_dn"
              placeholder="cn=admin,dc=example,dc=com"
              style="width: 360px"
            />
          </a-form-item>
          <a-form-item label="Bind 密码">
            <a-input
              v-model="form.bind_password"
              type="password"
              placeholder="请输入密码"
              allow-clear
              style="width: 240px"
            />
          </a-form-item>
        </div>

        <a-divider />

        <div class="form-section">
          <div class="section-label">过滤规则</div>
          <a-form-item label="用户过滤器">
            <a-input
              v-model="form.user_filter"
              placeholder="(objectClass=person)"
              style="width: 400px"
            />
            <div class="form-tip">用于搜索用户，常用：(objectClass=person)、(objectClass=user)</div>
          </a-form-item>
          <a-form-item label="分组过滤器">
            <a-input
              v-model="form.group_filter"
              placeholder="(objectClass=groupOfNames)"
              style="width: 400px"
            />
          </a-form-item>
        </div>

        <a-divider />

        <div class="form-section">
          <div class="section-label">同步配置</div>
          <a-form-item label="同步间隔">
            <a-input-number
              v-model="form.sync_interval"
              :min="60"
              :step="60"
              style="width: 120px"
            />
            <span class="form-tip">秒，建议不低于 300 秒（5分钟）</span>
          </a-form-item>
        </div>

        <div class="form-actions">
          <a-button @click="handleTest" :loading="testing">
            测试连接
          </a-button>
          <a-button type="primary" @click="handleSave" :loading="saving">
            保存配置
          </a-button>
          <a-button
            v-if="!ldapConfig.is_enabled"
            type="primary"
            status="success"
            @click="handleEnable"
            :loading="enabling"
          >
            保存并启用
          </a-button>
          <a-button
            v-else
            status="warning"
            @click="handleDisable"
          >
            禁用
          </a-button>
        </div>
      </a-form>
    </a-card>

    <!-- 测试结果 -->
    <a-modal
      v-model:visible="testResultVisible"
      title="连接测试结果"
      width="440px"
      :footer="null"
    >
      <div class="test-result" v-if="testResult">
        <a-result
          :status="testResult.success ? 'success' : 'error'"
          :title="testResult.success ? '连接成功' : '连接失败'"
        >
          <template #content>
            <div class="result-detail">{{ testResult.message }}</div>
          </template>
        </a-result>
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import { getLdapConfig, updateLdapConfig, testLdapConnection } from '@/api/security'

const saving = ref(false)
const testing = ref(false)
const enabling = ref(false)
const testResultVisible = ref(false)
const testResult = ref(null)

const ldapConfig = reactive({
  is_enabled: false
})

const form = reactive({
  server_url: '',
  port: 389,
  use_ssl: false,
  base_dn: '',
  bind_dn: '',
  bind_password: '',
  user_filter: '(objectClass=person)',
  group_filter: '(objectClass=groupOfNames)',
  sync_interval: 3600
})

onMounted(async () => {
  try {
    const res = await getLdapConfig()
    const data = res.data || res
    if (data) {
      Object.assign(ldapConfig, data)
      Object.assign(form, data)
    }
  } catch (e) {
    console.error('加载LDAP配置失败', e)
  }
})

async function handleSave() {
  saving.value = true
  try {
    await updateLdapConfig({ ...form, is_enabled: ldapConfig.is_enabled })
    Message.success('保存成功')
  } catch (e) {
    Message.error('保存失败')
  } finally {
    saving.value = false
  }
}

async function handleEnable() {
  enabling.value = true
  try {
    await updateLdapConfig({ ...form, is_enabled: true })
    ldapConfig.is_enabled = true
    Message.success('LDAP已启用')
  } catch (e) {
    Message.error('启用失败')
  } finally {
    enabling.value = false
  }
}

async function handleDisable() {
  try {
    await updateLdapConfig({ ...form, is_enabled: false })
    ldapConfig.is_enabled = false
    Message.success('LDAP已禁用')
  } catch (e) {
    Message.error('禁用失败')
  }
}

async function handleTest() {
  testing.value = true
  try {
    const res = await testLdapConnection({
      server_url: form.server_url,
      port: form.port,
      use_ssl: form.use_ssl,
      base_dn: form.base_dn,
      bind_dn: form.bind_dn,
      bind_password: form.bind_password
    })
    testResult.value = { success: res.success !== false, message: res.message || '连接成功' }
  } catch (e) {
    testResult.value = { success: false, message: e.message || '连接失败' }
  } finally {
    testing.value = false
    testResultVisible.value = true
  }
}
</script>

<style scoped>
.page-container {
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
  height: 100%;
  box-sizing: border-box;
}

.content-card {
  flex: 1;
  overflow: auto;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 10px;
}

.ldap-form {
  max-width: 600px;
}

.form-section {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.section-label {
  font-size: 13px;
  font-weight: 600;
  color: var(--color-text-1);
  margin-bottom: 8px;
}

.form-tip {
  color: var(--color-text-3);
  font-size: 12px;
  margin-left: 8px;
}

.form-actions {
  display: flex;
  gap: 10px;
  margin-top: 8px;
}

.test-result {
  padding: 16px 0;
}

.result-detail {
  font-size: 13px;
  color: var(--color-text-2);
  margin-top: 8px;
}
</style>
