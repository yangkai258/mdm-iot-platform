<template>
  <Breadcrumb :items="['Home','Advanced','CertificateManage','']" />
  <div class="page-container">
    <a-card class="general-card" title="证书管理">
      <template #extra>
        <a-button type="primary" @click="handleCreate"><icon-plus />新建证书</a-button>
      </template>
      <a-row :gutter="16" style="margin-bottom: 16px">
        <a-col :span="6"><a-card><a-statistic title="证书总数" :value="stats.total" /></a-card></a-col>
        <a-col :span="6"><a-card><a-statistic title="有效证书" :value="stats.active" /></a-card></a-col>
        <a-col :span="6"><a-card><a-statistic title="即将到期" :value="stats.expiring" :value-style="{ color: '#E6A23C' }" /></a-card></a-col>
        <a-col :span="6"><a-card><a-statistic title="已吊销" :value="stats.revoked" :value-style="{ color: '#F56C6C' }" /></a-card></a-col>
      </a-row>
      <div class="search-form">
        <a-form :model="form" layout="inline">
          <a-form-item label="搜索"><a-input v-model="form.keyword" placeholder="证书名称/序列号/主题" /></a-form-item>
          <a-form-item label="类型">
            <a-select v-model="form.cert_type" placeholder="选择类型" allow-clear style="width: 140px">
              <a-option value="device">设备证书</a-option>
              <a-option value="client">客户端证书</a-option>
              <a-option value="server">服务器证书</a-option>
              <a-option value="ca">CA证书</a-option>
            </a-select>
          </a-form-item>
          <a-form-item label="状态">
            <a-select v-model="form.status" placeholder="选择状态" allow-clear style="width: 130px">
              <a-option value="active">有效</a-option>
              <a-option value="expired">已过期</a-option>
              <a-option value="revoked">已吊销</a-option>
              <a-option value="pending">待激活</a-option>
            </a-select>
          </a-form-item>
          <a-form-item><a-button type="primary" @click="loadCertificates">查询</a-button></a-form-item>
        </a-form>
      </div>
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" />
      <a-modal v-model:visible="modalVisible" :title="modalTitle" @before-ok="handleSubmit">
        <a-form :model="form" layout="vertical">
          <a-form-item label="证书名称" required><a-input v-model="form.cert_name" placeholder="请输入证书名称" /></a-form-item>
          <a-form-item label="证书类型">
            <a-select v-model="form.cert_type" placeholder="选择类型">
              <a-option value="device">设备证书</a-option>
              <a-option value="client">客户端证书</a-option>
              <a-option value="server">服务器证书</a-option>
              <a-option value="ca">CA证书</a-option>
            </a-select>
          </a-form-item>
        </a-form>
      </a-modal>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import { IconPlus } from '@arco-design/web-vue/es/icon'

const loading = ref(false)
const data = ref<any[]>([])
const modalVisible = ref(false)
const modalTitle = ref('新建证书')
const stats = reactive({ total: 0, active: 0, expiring: 0, revoked: 0 })
const form = ref<any>({ keyword: '', cert_type: '', status: '' })

const columns = [
  { title: '证书名称', dataIndex: 'cert_name', width: 180 },
  { title: '类型', dataIndex: 'cert_type', width: 110 },
  { title: '状态', dataIndex: 'status', width: 90 },
  { title: '到期时间', dataIndex: 'not_after', width: 170 }
]

const pagination = ref({ current: 1, pageSize: 20, total: 0, showTotal: true })

function getStatusTag(s: string) {
  const map: any = { active: 'green', expired: 'red', revoked: 'orange', pending: 'blue' }
  return map[s] || 'gray'
}

async function loadCertificates() {
  try {
    loading.value = true
    data.value = [{ cert_name: '设备证书', cert_type: '设备证书', status: 'active', not_after: '2027-03-28' }]
    pagination.value.total = data.value.length
  } catch (err: any) {
    Message.error('加载失败: ' + err.message)
  } finally {
    loading.value = false
  }
}

function handleCreate() {
  modalTitle.value = '新建证书'
  form.value = { cert_name: '', cert_type: 'device' }
  modalVisible.value = true
}

async function handleSubmit(done: (val: boolean) => void) {
  Message.success('创建成功')
  modalVisible.value = false
  loadCertificates()
  done(true)
}

onMounted(() => { loadCertificates() })
</script>

<style scoped>
.page-container { padding: 16px; }
.search-form { margin-bottom: 16px; padding: 16px; background: var(--color-fill-lightest); border-radius: 4px; }
</style>
