<template>
  <div class="member-page">
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>会员管理</a-breadcrumb-item>
      <a-breadcrumb-item>微信公众号设置</a-breadcrumb-item>
    </a-breadcrumb>

    <a-card class="action-card">
      <a-form :model="form" layout="vertical" style="max-width: 640px;">
        <a-divider>基本信息</a-divider>

        <a-form-item label="AppID" required>
          <a-input v-model="form.appId" placeholder="请输入微信公众号AppID" />
          <div style="color:#999;font-size:12px;margin-top:4px;">登录微信公众平台，在「开发-基本配置」中获取AppID</div>
        </a-form-item>

        <a-form-item label="AppSecret" required>
          <a-input-password v-model="form.appSecret" placeholder="请输入微信公众号AppSecret" />
          <div style="color:#999;font-size:12px;margin-top:4px;">请妥善保管，切勿泄露给他人</div>
        </a-form-item>

        <a-form-item label="服务器地址（URL）">
          <a-input v-model="form.serverUrl" readonly placeholder="请先保存以获取服务器地址">
            <template #suffix>
              <a-button type="text" size="small" @click="copyUrl">复制</a-button>
            </template>
          </a-input>
          <div style="color:#999;font-size:12px;margin-top:4px;">将此处地址填写到微信公众平台的服务器配置中</div>
        </a-form-item>

        <a-divider>接口配置</a-divider>

        <a-form-item label="Token" required>
          <a-input v-model="form.token" placeholder="请输入Token" />
          <div style="color:#999;font-size:12px;margin-top:4px;">必须与微信公众平台配置中的Token保持一致</div>
        </a-form-item>

        <a-form-item label="EncodingAESKey">
          <a-input v-model="form.encodingAesKey" placeholder="请输入EncodingAESKey（可选）" />
          <div style="color:#999;font-size:12px;margin-top:4px;">消息加解密密钥，非必填，建议填写以提高安全性</div>
        </a-form-item>

        <a-form-item label="消息加解密方式">
          <a-radio-group v-model="form.encryptType">
            <a-radio value="plaintext">明文模式</a-radio>
            <a-radio value="compatible">兼容模式</a-radio>
            <a-radio value="safe">安全模式（推荐）</a-radio>
          </a-radio-group>
        </a-form-item>

        <a-divider>功能开关</a-divider>

        <a-form-item label="自动回复">
          <a-switch v-model="form.autoReply" checked-value="1" unchecked-value="0" />
          <span style="margin-left:8px;color:#999;">开启后，当用户发送消息时自动回复</span>
        </a-form-item>

        <a-form-item label="模板消息">
          <a-switch v-model="form.templateMsg" checked-value="1" unchecked-value="0" />
          <span style="margin-left:8px;color:#999;">开启后可向用户发送模板消息（如订单通知、会员变动等）</span>
        </a-form-item>

        <a-form-item label="自定义菜单">
          <a-switch v-model="form.customMenu" checked-value="1" unchecked-value="0" />
          <span style="margin-left:8px;color:#999;">开启后在微信公众号底部显示自定义菜单</span>
        </a-form-item>

        <a-divider>连接状态</a-divider>

        <a-form-item label="当前状态">
          <a-tag :color="form.connected ? 'green' : 'red'">{{ form.connected ? '已连接' : '未连接' }}</a-tag>
          <a-button type="text" size="small" style="margin-left:8px;" @click="testConnection">测试连接</a-button>
        </a-form-item>

        <a-form-item>
          <a-space>
            <a-button type="primary" @click="handleSave">保存配置</a-button>
            <a-button @click="handleReset">重置</a-button>
          </a-space>
        </a-form-item>
      </a-form>
    </a-card>
  </div>
</template>

<script setup>
import { reactive } from 'vue'
import { Message } from '@arco-design/web-vue'

const form = reactive({
  appId: 'wx1234567890abcdef',
  appSecret: '************************',
  serverUrl: 'https://api.example.com/wechat/callback',
  token: 'my_wechat_token_2026',
  encodingAesKey: '',
  encryptType: 'safe',
  autoReply: '1',
  templateMsg: '1',
  customMenu: '1',
  connected: true
})

const copyUrl = () => {
  navigator.clipboard.writeText(form.serverUrl).then(() => {
    Message.success('已复制到剪贴板')
  })
}

const testConnection = () => {
  Message.loading({ content: '正在测试连接...', duration: 1500 })
  setTimeout(() => {
    Message.success('连接正常')
  }, 1600)
}

const handleSave = () => {
  if (!form.appId || !form.appSecret || !form.token) {
    Message.error('请填写必填项')
    return
  }
  Message.success('保存成功')
}

const handleReset = () => {
  Object.assign(form, {
    appId: '',
    appSecret: '',
    serverUrl: 'https://api.example.com/wechat/callback',
    token: '',
    encodingAesKey: '',
    encryptType: 'safe',
    autoReply: '1',
    templateMsg: '1',
    customMenu: '1',
    connected: false
  })
  Message.info('已重置')
}
</script>

<style scoped>
.member-page { padding: 20px; }
.breadcrumb { margin-bottom: 16px; }
.action-card { margin-bottom: 16px; }
</style>
