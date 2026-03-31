<template>
  <div class="page-container">
    <div class="search-form">
      <a-form :model="form" layout="inline">
        <a-form-item label="小程序名称">
          <a-input v-model="form.keyword" placeholder="搜索小程序名称" style="width: 200px" />
        </a-form-item>
        <a-form-item>
          <a-button type="primary" @click="handleSearch">搜索</a-button>
          <a-button @click="handleReset">重置</a-button>
        </a-form-item>
      </a-form>
    </div>
    <div class="toolbar">
      <a-button type="primary" @click="handleSave">保存</a-button>
    </div>
    <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id">
      <template #status="{ record }"><a-tag :color="record.status === 1 ? 'green' : 'gray'">{{ record.status === 1 ? '启用' : '禁用' }}</a-tag></template>
      <template #actions="{ record }">
        <a-button type="text" size="small" @click="handleEdit(record)">编辑</a-button>
      </template>
    </a-table>
    <a-modal v-model:visible="modalVisible" :title="'编辑小程序配置'" @ok="handleSubmit" :width="520" :mask-closable="false">
      <a-form :model="form" layout="vertical">
        <a-form-item label="小程序名称"><a-input v-model="form.miniprogram_name" placeholder="请输入小程序名称" /></a-form-item>
        <a-form-item label="AppID"><a-input v-model="form.app_id" placeholder="微信小程序AppID" /></a-form-item>
        <a-form-item label="AppSecret"><a-input v-model="form.app_secret" placeholder="微信小程序AppSecret" type="password" /></a-form-item>
        <a-form-item label="小程序logo"><a-input v-model="form.logo_url" placeholder="请输入logo URL" /></a-form-item>
        <a-form-item label="会员卡背景图"><a-input v-model="form.card_bg" placeholder="会员卡背景图URL" /></a-form-item>
        <a-form-item label="会员卡激活链接"><a-input v-model="form.card_activate_url" placeholder="https://..." /></a-form-item>
        <a-form-item label="显示积分"><a-switch v-model="form.show_points" checked-value="1" unchecked-value="0" /></a-form-item>
        <a-form-item label="显示等级"><a-switch v-model="form.show_level" checked-value="1" unchecked-value="0" /></a-form-item>
        <a-form-item label="显示优惠券"><a-switch v-model="form.show_coupons" checked-value="1" unchecked-value="0" /></a-form-item>
        <a-form-item label="关注回复"><a-textarea v-model="form.subscribe_reply" :rows="2" placeholder="用户关注时的自动回复" /></a-form-item>
        <a-form-item label="生日祝福"><a-textarea v-model="form.birthday_wish" :rows="2" placeholder="会员生日祝福语" /></a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const API_BASE = '/api/v1'
const data = ref([])
const loading = ref(false)
const modalVisible = ref(false)

const form = reactive({
  keyword: '',
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

const pagination = reactive({ current: 1, pageSize: 10, total: 0 })

const columns = [
  { title: '配置项', dataIndex: 'name', width: 200 },
  { title: '值', dataIndex: 'value' },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '操作', slotName: 'actions', width: 100 }
]

const loadData = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`${API_BASE}/member/miniprogram`, { headers: { 'Authorization': `Bearer ${token}` } })
    const resp = await res.json()
    if (resp.code === 0) {
      if (resp.data) {
        Object.assign(form, resp.data)
        data.value = [{ id: 1, name: '小程序配置', value: resp.data.miniprogram_name || '-', status: 1 }]
      }
    }
  } catch (e) { Message.error('加载配置失败') } finally { loading.value = false }
}

const handleSearch = () => loadData()
const handleReset = () => { form.keyword = ''; loadData() }
const onPageChange = (page) => { pagination.current = page; loadData() }

const handleSave = async () => {
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`${API_BASE}/member/miniprogram`, {
      method: 'PUT',
      headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' },
      body: JSON.stringify(form)
    })
    const resp = await res.json()
    if (resp.code === 0) Message.success('保存成功')
    else Message.error(resp.message || '保存失败')
  } catch (e) { Message.error('保存失败') }
}

const handleEdit = (record) => { modalVisible.value = true }

const handleSubmit = async () => {
  modalVisible.value = false
  await handleSave()
}

onMounted(() => loadData())
</script>

<style scoped>
.page-container { background: #fff; border-radius: 4px; padding: 20px; }
.search-form { margin-bottom: 16px; padding: 16px; background: #f7f8fa; border-radius: 4px; }
.toolbar { margin-bottom: 16px; }
</style>
