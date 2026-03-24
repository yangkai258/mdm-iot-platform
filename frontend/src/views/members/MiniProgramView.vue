<template>
  <div class="page-container">
    <div class="toolbar">
      <a-button type="primary" @click="handleSave">保存设置</a-button>
      <a-button @click="handleReset">重置</a-button>
    </div>
    <a-form :model="form" layout="vertical" class="settings-form">
      <a-divider>基本信息</a-divider>
      <a-row :gutter="16">
        <a-col :span="12">
          <a-form-item label="小程序名称">
            <a-input v-model="form.miniprogram_name" placeholder="请输入小程序名称" />
          </a-form-item>
        </a-col>
        <a-col :span="12">
          <a-form-item label="AppID">
            <a-input v-model="form.app_id" placeholder="请输入AppID" />
          </a-form-item>
        </a-col>
        <a-col :span="24">
          <a-form-item label="AppSecret">
            <a-input v-model="form.app_secret" placeholder="请输入AppSecret" type="password" />
          </a-form-item>
        </a-col>
      </a-row>
      <a-divider>版本信息</a-divider>
      <a-row :gutter="16">
        <a-col :span="12">
          <a-form-item label="当前版本">
            <a-input v-model="form.version" placeholder="请输入版本号" />
          </a-form-item>
        </a-col>
        <a-col :span="12">
          <a-form-item label="最低版本">
            <a-input v-model="form.min_version" placeholder="请输入最低版本要求" />
          </a-form-item>
        </a-col>
        <a-col :span="24">
          <a-form-item label="更新说明">
            <a-textarea v-model="form.update_note" :rows="2" placeholder="请输入版本更新说明" />
          </a-form-item>
        </a-col>
      </a-row>
      <a-divider>卡面设置</a-divider>
      <a-row :gutter="16">
        <a-col :span="12">
          <a-form-item label="会员卡背景">
            <a-input v-model="form.card_bg" placeholder="请输入背景图片URL" />
          </a-form-item>
        </a-col>
        <a-col :span="12">
          <a-form-item label="激活链接">
            <a-input v-model="form.card_activate_url" placeholder="请输入激活链接" />
          </a-form-item>
        </a-col>
        <a-col :span="8">
          <a-form-item label="显示积分">
            <a-switch v-model="form.show_points" true-value="1" false-value="0" />
          </a-form-item>
        </a-col>
        <a-col :span="8">
          <a-form-item label="显示等级">
            <a-switch v-model="form.show_level" true-value="1" false-value="0" />
          </a-form-item>
        </a-col>
        <a-col :span="8">
          <a-form-item label="显示优惠券">
            <a-switch v-model="form.show_coupons" true-value="1" false-value="0" />
          </a-form-item>
        </a-col>
      </a-row>
      <a-divider>消息设置</a-divider>
      <a-row :gutter="16">
        <a-col :span="24">
          <a-form-item label="关注回复">
            <a-textarea v-model="form.subscribe_reply" :rows="2" placeholder="请输入关注回复内容" />
          </a-form-item>
        </a-col>
        <a-col :span="24">
          <a-form-item label="生日祝福">
            <a-textarea v-model="form.birthday_wish" :rows="2" placeholder="请输入生日祝福内容" />
          </a-form-item>
        </a-col>
      </a-row>
    </a-form>
  </div>
</template>

<script setup>

import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const API_BASE = '/api/v1'

const form = reactive({
  miniprogram_name: '',
  app_id: '',
  app_secret: '',
  logo_url: '',
  version: '',
  min_version: '',
  update_note: '',
  card_bg: '',
  card_activate_url: '',
  show_points: '1',
  show_level: '1',
  show_coupons: '1',
  subscribe_reply: '',
  birthday_wish: ''
})

const loadData = async () => {
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`${API_BASE}/member/miniprogram`, { headers: { 'Authorization': `Bearer ${token}` } })
    const data = await res.json()
    if (data.code === 0 && data.data) Object.assign(form, data.data)
  } catch (e) {}
}

const handleSave = async () => {
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`${API_BASE}/member/miniprogram`, {
      method: 'PUT',
      headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' },
      body: JSON.stringify(form)
    })
    const data = await res.json()
    if (data.code === 0) Message.success('保存成功')
    else Message.error(data.message || '保存失败')
  } catch (e) { Message.error('保存失败') }
}

const handleReset = () => {
  Object.assign(form, {
    miniprogram_name: '', app_id: '', app_secret: '', logo_url: '',
    version: '', min_version: '', update_note: '',
    card_bg: '', card_activate_url: '', show_points: '1', show_level: '1',
    show_coupons: '1', subscribe_reply: '', birthday_wish: ''
  })
  Message.info('已重置')
}

onMounted(() => loadData())

</script>

<style scoped>
.page-container { background: #fff; border-radius: 4px; padding: 20px; }
.toolbar { margin-bottom: 16px; }
.settings-form { max-width: 800px; }
</style>
