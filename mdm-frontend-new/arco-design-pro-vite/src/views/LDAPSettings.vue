<template>
  <div class="container">
    <a-card>
      <template #title>
        <a-space><icon-user-group /> LDAP/AD集成</a-space>
      </template>

      <a-row :gutter="16">
        <a-col :span="12">
          <a-card title="连接配置">
            <a-form :model="formData" layout="vertical">
              <a-form-item label="服务器地址">
                <a-input v-model="formData.server" placeholder="ldap://your-server.com" />
              </a-form-item>
              <a-form-item label="端口">
                <a-input-number v-model="formData.port" :min="1" :max="65535" />
              </a-form-item>
              <a-form-item label="Base DN">
                <a-input v-model="formData.baseDn" placeholder="dc=example,dc=com" />
              </a-form-item>
              <a-form-item label="管理员DN">
                <a-input v-model="formData.adminDn" placeholder="cn=admin,dc=example,dc=com" />
              </a-form-item>
              <a-form-item label="密码">
                <a-input-password v-model="formData.password" />
              </a-form-item>
              <a-form-item>
                <a-button type="primary" @click="handleTestConnection">测试连接</a-button>
                <a-button @click="handleSave">保存</a-button>
              </a-form-item>
            </a-form>
          </a-card>
        </a-col>

        <a-col :span="12">
          <a-card title="同步状态">
            <a-descriptions :column="1" bordered>
              <a-descriptions-item label="连接状态">
                <a-badge status="success" text="已连接" />
              </a-descriptions-item>
              <a-descriptions-item label="最后同步">2026-03-28 10:00:00</a-descriptions-item>
              <a-descriptions-item label="同步用户数">256</a-descriptions-item>
              <a-descriptions-item label="同步部门数">12</a-descriptions-item>
            </a-descriptions>
            <a-space style="margin-top: 16px">
              <a-button @click="handleSyncNow">立即同步</a-button>
              <a-button @click="handleSyncSettings">同步设置</a-button>
            </a-space>
          </a-card>
        </a-col>
      </a-row>

      <a-card title="用户映射" style="margin-top: 16px">
        <a-form :model="mappingForm" layout="inline">
          <a-form-item label="LDAP属性">
            <a-input v-model="mappingForm.ldapAttr" style="width: 150px" />
          </a-form-item>
          <a-form-item label="映射到">
            <a-select v-model="mappingForm.targetField" style="width: 150px">
              <a-option value="username">用户名</a-option>
              <a-option value="email">邮箱</a-option>
              <a-option value="phone">手机</a-option>
            </a-select>
          </a-form-item>
          <a-form-item>
            <a-button type="primary">添加映射</a-button>
          </a-form-item>
        </a-form>
      </a-card>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'

const formData = reactive({ server: '', port: 389, baseDn: '', adminDn: '', password: '' })
const mappingForm = reactive({ ldapAttr: '', targetField: '' })

const handleTestConnection = () => { }
const handleSave = () => { }
const handleSyncNow = () => { }
const handleSyncSettings = () => { }
</script>

<style scoped>
.container { padding: 16px; }
</style>
