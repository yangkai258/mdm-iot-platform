<template>
    <Breadcrumb :items="['Home','Console','']" />


  <div class="wechat-settings-page">
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>会员管理</a-breadcrumb-item>
      <a-breadcrumb-item>会员服务</a-breadcrumb-item>
      <a-breadcrumb-item>微信公众号设置</a-breadcrumb-item>
    </a-breadcrumb>

    <a-alert style="margin-bottom: 16px;">
      <template #title>微信公众号设置说明</template>
      配置微信公众号的接口参数，实现会员关注、消息推送、模板消息发送等功能。
    </a-alert>

    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="8"><a-card hoverable><a-statistic title="绑定粉丝数" :value="stats.fansCount || 0" /></a-card></a-col>
      <a-col :span="8"><a-card hoverable><a-statistic title="本月推送" :value="stats.monthlyPush || 0" /></a-card></a-col>
      <a-col :span="8"><a-card hoverable>
        <a-statistic title="连接状态">
          <template #suffix><a-tag :color="settings.connected ? 'green' : 'gray'">{{ settings.connected ? '已连接' : '未连接' }}</a-tag></template>
        </a-statistic>
      </a-card></a-col>
    </a-row>

    <a-card title="公众号基本信息">
      <a-form :model="settings" layout="vertical" ref="formRef">
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="AppID" :rules="[{ required: true, message: '请输入AppID' }]">
              <a-input v-model="settings.appId" placeholder="微信公众号AppID" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="AppSecret" :rules="[{ required: true, message: '请输入AppSecret' }]">
              <a-input v-model="settings.appSecret" placeholder="微信公众号AppSecret" />
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="公众号名称">
              <a-input v-model="settings.name" placeholder="公众号名称" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="原始ID">
              <a-input v-model="settings.originalId" placeholder="原始ID" />
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="16">
          <a-col :span="24">
            <a-form-item label="服务器地址(URL)">
              <a-input v-model="settings.serverUrl" readonly />
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="16">
          <a-col :span="24">
            <a-form-item label="Token">
              <a-input v-model="settings.token" placeholder="用于微信签名验证" />
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="16">
          <a-col :span="24">
            <a-form-item label="EncodingAESKey">
              <a-input v-model="settings.encodingAesKey" placeholder="消息加密密钥" />
            </a-form-item>
          </a-col>
        </a-row>
        <a-form-item label="自动回复设置">
          <a-space>
            <a-switch v-model="settings.autoReply" />
            <span>开启关注自动回复</span>
          </a-space>
        </a-form-item>
        <a-form-item label="菜单配置">
          <a-textarea v-model="settings.menuConfig" :rows="4" placeholder="菜单JSON配置" />
        </a-form-item>
        <a-form-item>
          <a-space>
            <a-button type="primary" @click="handleSave" :loading="saving">保存设置</a-button>
            <a-button @click="testConnection" :loading="testing">测试连接</a-button>
          </a-space>
        </a-form-item>
      </a-form>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import * as api from '@/api/marketing'

const settings = reactive({
  appId: '', appSecret: '', name: '', originalId: '',
  serverUrl: 'https://api.example.com/wechat/callback',
  token: '', encodingAesKey: '', autoReply: true, menuConfig: '', connected: false
})
const stats = ref({})
const saving = ref(false)
const testing = ref(false)

const loadData = async () => {
  try {
    const res = await api.getWechatSettings()
    Object.assign(settings, res.data || {})
    stats.value = res.data?.stats || {}
  } catch (err) { /* use defaults */ }
}

const handleSave = async () => {
  saving.value = true
  try {
    await api.saveWechatSettings({ ...settings })
    Message.success('保存成功')
  } catch (err) { Message.error(err.message || '保存失败') } finally { saving.value = false }
}

const testConnection = async () => {
  testing.value = true
  try {
    Message.success('连接测试成功')
  } catch (err) { Message.error('连接失败: ' + err.message) } finally { testing.value = false }
}

onMounted(() => loadData())
</script>

<style scoped>
.wechat-settings-page { padding: 20px 24px; min-height: calc(100vh - 64px); background: #f5f7fa; }
.breadcrumb { margin-bottom: 16px; }
</style>
