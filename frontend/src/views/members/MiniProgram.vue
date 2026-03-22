<template>
  <div class="pro-page-container">
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>会员管理</a-breadcrumb-item>
      <a-breadcrumb-item>小程序基础数据</a-breadcrumb-item>
    </a-breadcrumb>

    <div class="pro-content-area" style="margin-top: 16px;">
      <a-form :model="form" layout="vertical" style="max-width: 600px;">
        <a-divider>小程序基本信息</a-divider>

        <a-form-item label="小程序名称">
          <a-input v-model="form.miniprogram_name" placeholder="请输入小程序名称" />
        </a-form-item>

        <a-form-item label="AppID">
          <a-input v-model="form.app_id" placeholder="微信小程序AppID" />
        </a-form-item>

        <a-form-item label="AppSecret">
          <a-input v-model="form.app_secret" placeholder="微信小程序AppSecret" type="password" />
        </a-form-item>

        <a-form-item label="小程序logo">
          <a-space>
            <a-input v-model="form.logo_url" placeholder="请输入logo URL" style="width: 300px;" />
            <a-button v-if="form.logo_url" @click="previewLogo">预览</a-button>
          </a-space>
        </a-form-item>

        <a-divider>会员卡配置</a-divider>

        <a-form-item label="会员卡背景图">
          <a-input v-model="form.card_bg" placeholder="会员卡背景图URL" />
        </a-form-item>

        <a-form-item label="会员卡激活链接">
          <a-input v-model="form.card_activate_url" placeholder="https://..." />
        </a-form-item>

        <a-form-item label="显示积分">
          <a-switch v-model="form.show_points" checked-value="1" unchecked-value="0" />
        </a-form-item>

        <a-form-item label="显示等级">
          <a-switch v-model="form.show_level" checked-value="1" unchecked-value="0" />
        </a-form-item>

        <a-form-item label="显示优惠券">
          <a-switch v-model="form.show_coupons" checked-value="1" unchecked-value="0" />
        </a-form-item>

        <a-divider>消息推送配置</a-divider>

        <a-form-item label="关注回复">
          <a-textarea v-model="form.subscribe_reply" :rows="2" placeholder="用户关注时的自动回复" />
        </a-form-item>

        <a-form-item label="生日祝福">
          <a-textarea v-model="form.birthday_wish" :rows="2" placeholder="会员生日祝福语" />
        </a-form-item>

        <a-form-item>
          <a-space>
            <a-button type="primary" @click="handleSave">保存</a-button>
            <a-button @click="handleReset">重置</a-button>
          </a-space>
        </a-form-item>
      </a-form>
    </div>

    <a-modal v-model:visible="previewVisible" title="Logo预览" :width="400">
      <div style="text-align: center;">
        <img v-if="form.logo_url" :src="form.logo_url" style="max-width: 200px; max-height: 200px;" />
        <div v-else style="color: #999;">暂无logo</div>
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const API_BASE = '/api/v1'
const previewVisible = ref(false)

const form = reactive({
  miniprogram_name: '',
  app_id: '',
  app_secret: '',
  logo_url: '',
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
    card_bg: '', card_activate_url: '', show_points: '1', show_level: '1',
    show_coupons: '1', subscribe_reply: '', birthday_wish: ''
  })
  Message.info('已重置')
}

const previewLogo = () => { previewVisible.value = true }

onMounted(() => loadData())
</script>

<style scoped>
.pro-page-container { padding: 20px 24px; min-height: calc(100vh - 64px); background: #f5f7fa; }
.pro-breadcrumb { margin-bottom: 16px; }
.pro-content-area { background: #fff; border-radius: 8px; padding: 24px; box-shadow: 0 1px 3px rgba(0,0,0,0.04); }
</style>
