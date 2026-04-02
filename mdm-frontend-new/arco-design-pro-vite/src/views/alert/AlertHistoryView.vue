<template>
    <Breadcrumb :items="['Home','Console','']" />


  <div class="container">
    <a-card class="general-card" title="告警历史">
      <template #extra>
        <a-space :size="12">
          <a-button type="primary" @click="handleCreate"><icon-plus />新建</a-button>
          <a-button @click="handleSearch"><icon-refresh />刷新</a-button>
        </a-space>
      </template>
      <a-row :gutter="16">
        <a-col :span="6">
          <a-form-item label="设备">
            <a-input v-model="form.device_id" placeholder="设备ID" @pressEnter="handleSearch" />
          </a-form-item>
        </a-col>
        <a-col :span="6">
          <a-form-item label="告警类型">
            <a-select v-model="form.alert_type" placeholder="选择类型" allow-clear style="width: 100%">
              <a-option value="temperature">温度告警</a-option>
              <a-option value="battery">电量告警</a-option>
              <a-option value="offline">离线告警</a-option>
            </a-select>
          </a-form-item>
        </a-col>
        <a-col :span="6">
          <a-form-item label="级别">
            <a-select v-model="form.severity" placeholder="选择级别" allow-clear style="width: 100%">
              <a-option value="info">提示</a-option>
              <a-option value="warning">警告</a-option>
              <a-option value="error">错误</a-option>
            </a-select>
          </a-form-item>
        </a-col>
        <a-col :span="6">
          <a-form-item label="状态">
            <a-select v-model="form.status" placeholder="选择状态" allow-clear style="width: 100%">
              <a-option value="pending">待处理</a-option>
              <a-option value="resolved">已处理</a-option>
            </a-select>
          </a-form-item>
        </a-col>
      </a-row>
      <a-divider style="margin: 0 0 16px 0" />
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" row-key="id" />
    </a-card>
      </a-table>
    <a-modal v-model:visible="modalVisible" :title="modalTitle">
      <a-form :model="form" label-col-flex="100px">
        <a-form-item label="名称"><a-input v-model="form.name" /></a-form-item>
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
const form = reactive({ name: '', device_id: '', alert_type: '', severity: undefined, status: undefined })
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })
const columns = [
  { title: '时间', dataIndex: 'created_at', width: 170 },
  { title: '设备', dataIndex: 'device_id', width: 140 },
  { title: '告警类型', dataIndex: 'alert_type', width: 120 },
  { title: '级别', dataIndex: 'severity', width: 80 },
  { title: '告警消息', dataIndex: 'message', ellipsis: true },
  { title: '触发值', dataIndex: 'trigger_value', width: 100 },
  { title: '状态', dataIndex: 'status', width: 90 }
]

const handleSearch = () => {
  pagination.current = 1
  loadData()
}
const handleReset = () => {
  Object.keys(form).forEach(k => { form[k] = k === 'severity' || k === 'status' ? undefined : '' })
  pagination.current = 1
  loadData()
}
const handleCreate = () => { modalTitle.value = '新建'; modalVisible.value = true }
const handleSubmit = () => { modalVisible.value = false; Message.success('保存成功') }

const loadData = async () => {
  loading.value = true
  try {
    const params = { page: pagination.current, page_size: pagination.pageSize }
    Object.keys(form).forEach(k => { if (form[k]) params[k] = form[k] })
    const res = await fetch('/api/alerts/history?' + new URLSearchParams(params), {
      headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token') }
    }).then(r => r.json())
    if (res.code === 0) { data.value = res.data?.list || []; pagination.total = res.data?.total || 0 }
  } catch (e) { Message.error('加载失败') } finally { loading.value = false }
}
onMounted(() => { loadData() })
</script>
