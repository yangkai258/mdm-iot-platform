<template>
  <div class="member-page">
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>会员管理</a-breadcrumb-item>
      <a-breadcrumb-item>会员短信通道</a-breadcrumb-item>
    </a-breadcrumb>

    <a-row :gutter="16" class="stats-row">
      <a-col :span="8">
        <a-card class="stat-card">
          <a-statistic title="通道总数" :value="stats.total || 0" />
        </a-card>
      </a-col>
      <a-col :span="8">
        <a-card class="stat-card">
          <a-statistic title="本月发送" :value="stats.monthSent || 0" suffix="条" />
        </a-card>
      </a-col>
      <a-col :span="8">
        <a-card class="stat-card">
          <a-statistic title="成功率" :value="stats.successRate || 0" suffix="%" />
        </a-card>
      </a-col>
    </a-row>

    <a-card class="action-card">
      <a-space wrap>
        <a-input-search v-model="filters.keyword" placeholder="搜索通道名称" style="width: 220px" search-button @search="loadData" />
        <a-select v-model="filters.type" placeholder="通道类型" allow-clear style="width: 140px" @change="loadData">
          <a-option value="gateway">网关通道</a-option>
          <a-option value="direct">直连通道</a-option>
          <a-option value="international">国际通道</a-option>
        </a-select>
        <a-button type="primary" @click="showCreate">新建</a-button>
        <a-button @click="loadData">刷新</a-button>
      </a-space>
    </a-card>

    <a-card class="table-card">
      <a-table :columns="columns" :data="dataList" :loading="loading" :pagination="paginationConfig" @page-change="onPageChange" row-key="id" :scroll="{ x: 1000 }">
        <template #status="{ record }">
          <a-tag :color="record.status === 1 ? 'green' : record.status === 0 ? 'orange' : 'red'">
            {{ ['', '正常', '欠费', '禁用'][record.status] }}
          </a-tag>
        </template>
        <template #actions="{ record }">
          <a-button type="text" size="small" @click="showEdit(record)">编辑</a-button>
          <a-button type="text" size="small" status="danger" @click="handleDelete(record)">删除</a-button>
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:visible="formVisible" :title="isEdit ? '编辑短信通道' : '新建短信通道'" :width="520px" @before-ok="handleSubmit" @cancel="formVisible = false">
      <a-form :model="form" layout="vertical">
        <a-form-item label="通道名称" required>
          <a-input v-model="form.name" placeholder="请输入通道名称" />
        </a-form-item>
        <a-form-item label="通道类型">
          <a-select v-model="form.type" placeholder="请选择通道类型" style="width: 100%;">
            <a-option value="gateway">网关通道</a-option>
            <a-option value="direct">直连通道</a-option>
            <a-option value="international">国际通道</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="服务商">
          <a-input v-model="form.provider" placeholder="请输入服务商名称" />
        </a-form-item>
        <a-form-item label="接口地址">
          <a-input v-model="form.apiUrl" placeholder="请输入短信接口地址" />
        </a-form-item>
        <a-form-item label="API Key">
          <a-input-password v-model="form.apiKey" placeholder="请输入API Key" />
        </a-form-item>
        <a-form-item label="单条成本（元）">
          <a-input-number v-model="form.costPerSms" :min="0" :precision="4" style="width: 100%;" />
        </a-form-item>
        <a-form-item label="日发送上限">
          <a-input-number v-model="form.dailyLimit" :min="0" style="width: 100%;" />
        </a-form-item>
        <a-form-item label="状态">
          <a-switch v-model="form.status" checked-value="1" unchecked-value="0" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { Message } from '@arco-design/web-vue'

const loading = ref(false)
const formVisible = ref(false)
const isEdit = ref(false)
const dataList = ref([])
const stats = ref({ total: 0, monthSent: 0, successRate: 0 })
const filters = reactive({ keyword: '', type: '' })
const paginationConfig = reactive({ current: 1, pageSize: 10, total: 0 })
const form = reactive({ id: null, name: '', type: 'gateway', provider: '', apiUrl: '', apiKey: '', costPerSms: 0.05, dailyLimit: 10000, status: '1' })

const columns = [
  { title: '通道名称', dataIndex: 'name', width: 180 },
  { title: '通道类型', dataIndex: 'typeName', width: 130 },
  { title: '服务商', dataIndex: 'provider', width: 150 },
  { title: '单条成本', dataIndex: 'costPerSms', width: 100 },
  { title: '今日已发', dataIndex: 'todaySent', width: 100 },
  { title: '日上限', dataIndex: 'dailyLimit', width: 100 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '操作', slotName: 'actions', width: 150, fixed: 'right' }
]

const typeNameMap = { gateway: '网关通道', direct: '直连通道', international: '国际通道' }

const mockData = () => [
  { id: 1, name: '阿里云短信', type: 'gateway', typeName: '网关通道', provider: '阿里云', costPerSms: 0.04, todaySent: 2345, dailyLimit: 10000, status: 1 },
  { id: 2, name: '腾讯云短信', type: 'gateway', typeName: '网关通道', provider: '腾讯云', costPerSms: 0.045, todaySent: 1234, dailyLimit: 8000, status: 1 },
  { id: 3, name: '国际短信通道', type: 'international', typeName: '国际通道', provider: 'Twilio', costPerSms: 0.20, todaySent: 56, dailyLimit: 1000, status: 1 }
]

const loadData = () => {
  loading.value = true
  setTimeout(() => {
    dataList.value = mockData()
    stats.value = { total: mockData().length, monthSent: 45678, successRate: 98.5 }
    paginationConfig.total = dataList.value.length
    loading.value = false
  }, 400)
}

const onPageChange = (page) => { paginationConfig.current = page; loadData() }
const showCreate = () => { isEdit.value = false; Object.assign(form, { id: null, name: '', type: 'gateway', provider: '', apiUrl: '', apiKey: '', costPerSms: 0.05, dailyLimit: 10000, status: '1' }); formVisible.value = true }
const showEdit = (record) => { isEdit.value = true; Object.assign(form, { ...record, status: String(record.status) }); formVisible.value = true }
const handleSubmit = (done) => {
  if (!form.name) { Message.error('请填写必填项'); done(false); return }
  setTimeout(() => { Message.success(isEdit.value ? '更新成功' : '创建成功'); formVisible.value = false; loadData(); done(true) }, 400)
}
const handleDelete = () => { Message.success('删除成功'); loadData() }

loadData()
</script>

<style scoped>
.member-page { padding: 20px; }
.breadcrumb { margin-bottom: 16px; }
.stats-row { margin-bottom: 16px; }
.stat-card { text-align: center; }
.action-card { margin-bottom: 16px; }
</style>
