<template>
  <div class="container">
    <a-card class="general-card" title="短信渠道配置">
      <template #extra>
        <a-space :size="12">
          <a-button type="primary" @click="handleCreate"><icon-plus />新建</a-button>
          <a-button @click="handleSearch"><icon-refresh />刷新</a-button>
        </a-space>
      </template>
      <a-row :gutter="16">
        <a-col :span="8">
          <a-form-item label="渠道名称">
            <a-input v-model="form.channel_name" placeholder="请输入" @pressEnter="handleSearch" />
          </a-form-item>
        </a-col>
        <a-col :flex="'86px'" style="display: flex; align-items: flex-end">
          <a-space direction="vertical" :size="8">
            <a-button type="primary" @click="handleSearch">查询</a-button>
            <a-button @click="handleReset">重置</a-button>
          </a-space>
        </a-col>
      </a-row>
      <a-divider style="margin: 0 0 16px 0" />
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" row-key="id" />
    </a-card>
    <a-modal v-model:visible="modalVisible" :title="modalTitle" :width="520">
      <a-form :model="form" label-col-flex="100px">
        <a-form-item label="渠道名称"><a-input v-model="form.channel_name" /></a-form-item>
        <a-form-item label="服务商">
          <a-select v-model="form.sms_provider">
            <a-option value="aliyun">阿里云</a-option>
            <a-option value="tencent">腾讯云</a-option>
            <a-option value="huawei">华为云</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="AccessKey"><a-input v-model="form.access_key" /></a-form-item>
        <a-form-item label="SecretKey"><a-input-password v-model="form.secret_key" /></a-form-item>
        <a-form-item label="签名"><a-input v-model="form.sign_name" /></a-form-item>
      </a-form>
      <template #footer>
        <a-button @click="modalVisible = false">取消</a-button>
        <a-button type="primary" @click="handleSubmit">确定</a-button>
      </template>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import Breadcrumb from '@/components/Breadcrumb.vue'

const loading = ref(false)
const data = ref([])
const modalVisible = ref(false)
const modalTitle = ref('新建')
const form = reactive({ channel_name: '', sms_provider: 'aliyun', access_key: '', secret_key: '', sign_name: '' })
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })
const columns = [
  { title: '渠道名称', dataIndex: 'channel_name', width: 200 },
  { title: '服务商', dataIndex: 'config', customRender: ({ record }) => record.config?.sms_provider || '-', width: 120 },
  { title: '签名', dataIndex: 'config', customRender: ({ record }) => record.config?.sign_name || '-', width: 180 }
]

const handleSearch = () => { loadData() }
const handleReset = () => { form.channel_name = ''; loadData() }
const handleCreate = () => { modalTitle.value = '新建'; modalVisible.value = true }
const handleSubmit = async () => {
  modalVisible.value = false
  try {
    await fetch('/api/v1/alerts/channels', {
      method: 'POST',
      headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token'), 'Content-Type': 'application/json' },
      body: JSON.stringify({ channel_type: 'sms', channel_name: form.channel_name, config: { sms_provider: form.sms_provider, access_key: form.access_key, secret_key: form.secret_key, sign_name: form.sign_name } })
    })
    Message.success('保存成功')
    loadData()
  } catch (e) { Message.error('保存失败') }
}

const loadData = async () => {
  loading.value = true
  try {
    const res = await fetch('/api/v1/alerts/channels?channel_type=sms', {
      headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token') }
    }).then(r => r.json())
    if (res.code === 0) { data.value = res.data?.list || []; pagination.total = res.data?.total || 0 }
  } catch (e) { Message.error('加载失败') } finally { loading.value = false }
}
onMounted(() => { loadData() })
</script>
