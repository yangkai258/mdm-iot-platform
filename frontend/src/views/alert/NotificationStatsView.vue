<template>
  <div class="container">
    <Breadcrumb :items="['menu.alerts', 'menu.alerts.stats']" />
    <a-card class="general-card" title="通知统计">
      <template #extra>
        <a-space :size="12">
          <a-button type="primary" @click="handleCreate"><icon-plus />新建</a-button>
          <a-button @click="handleSearch"><icon-refresh />刷新</a-button>
        </a-space>
      </template>
      <a-row :gutter="16">
        <a-col :span="8">
          <a-form-item label="时间范围">
            <a-range-picker v-model="form.dateRange" style="width: 100%" />
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
const form = reactive({ name: '', dateRange: [] })
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })
const columns = [
  { title: '渠道类型', dataIndex: 'channel_type', width: 120 },
  { title: '发送成功', dataIndex: 'sent', width: 100 },
  { title: '发送失败', dataIndex: 'failed', width: 100 },
  { title: '成功率', dataIndex: 'success_rate', width: 100 }
]

const handleSearch = () => { loadData() }
const handleReset = () => { form.dateRange = []; loadData() }
const handleCreate = () => { modalTitle.value = '新建'; modalVisible.value = true }
const handleSubmit = () => { modalVisible.value = false; Message.success('保存成功') }

const loadData = async () => {
  loading.value = true
  try {
    const params = {}
    if (form.dateRange && form.dateRange.length === 2) {
      params.start_time = form.dateRange[0].toISOString()
      params.end_time = form.dateRange[1].toISOString()
    }
    const res = await fetch('/api/v1/alerts/notification-stats?' + new URLSearchParams(params), {
      headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token') }
    }).then(r => r.json())
    if (res.code === 0) { data.value = res.data?.by_channel || []; pagination.total = data.value.length }
  } catch (e) { Message.error('加载失败') } finally { loading.value = false }
}
onMounted(() => { loadData() })
</script>
