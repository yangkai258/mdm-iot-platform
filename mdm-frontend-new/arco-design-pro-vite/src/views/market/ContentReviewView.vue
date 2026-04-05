<template>
    <Breadcrumb :items="['Home','Console','']" />


  <div class="container">
    <a-card class="general-card" title="内容审核">
      <template #extra>
        <a-button @click="loadData"><icon-refresh />刷新</a-button>
      </template>
      <a-row :gutter="16">
        <a-col :span="8">
          <a-form-item label="内容ID"><a-input v-model="form.keyword" placeholder="请输入" @pressEnter="loadData" /></a-form-item>
        </a-col>
        <a-col :flex="'86px'" style="display: flex; align-items: flex-end">
          <a-space direction="vertical" :size="8">
            <a-button type="primary" @click="loadData">查询</a-button>
            <a-button @click="Object.keys(form).forEach(k => form[k] = ''); loadData()">重置</a-button>
          </a-space>
        </a-col>
      </a-row>
      <a-divider style="margin: 0 0 16px 0" />
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id">
        <template #status="{ record }"><a-badge :color="record.status === 'approved' ? 'green' : record.status === 'rejected' ? 'red' : 'orange'" :text="record.status === 'approved' ? '通过' : record.status === 'rejected' ? '拒绝' : '待审核'" /></template>
        <template #actions="{ record }">
          <a-button v-if="record.status === 'pending'" type="text" size="small" @click="handleApprove(record)">通过</a-button>
          <a-button v-if="record.status === 'pending'" type="text" size="small" status="danger" @click="handleReject(record)">拒绝</a-button>
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import Breadcrumb from '@/components/Breadcrumb.vue'

const loading = ref(false)
const form = reactive({ keyword: '' })
const data = ref([])
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })
const columns = [
  { title: '内容ID', dataIndex: 'id', width: 80 },
  { title: '内容类型', dataIndex: 'content_type', width: 120 },
  { title: '内容摘要', dataIndex: 'summary', ellipsis: true },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '提交时间', dataIndex: 'created_at', width: 170 },
  { title: '操作', slotName: 'actions', width: 160 }
]

const loadData = async () => {
  loading.value = true
  try {
    const res = await fetch('/api/market/content-review', { headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token') } }).then(r => r.json())
    data.value = res.data?.list || []
    pagination.total = data.value.length
  } catch { data.value = [] } finally { loading.value = false }
}
const handleApprove = (record) => { record.status = 'approved'; Message.success('已通过') }
const handleReject = (record) => { record.status = 'rejected'; Message.success('已拒绝') }
const onPageChange = (page) => { pagination.current = page; loadData() }
onMounted(() => loadData())
</script>
