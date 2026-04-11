<template>
  <div class="container">
    <a-card class="general-card" title="策略分发记录">
      <template #extra>
        <a-button @click="loadData"><icon-refresh />刷新</a-button>
      </template>
      <a-row :gutter="16">
        <a-col :span="8">
          <a-form-item label="策略名称">
            <a-input v-model="form.keyword" placeholder="请输入" @pressEnter="loadData" />
          </a-form-item>
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
        <template #status="{ record }"><a-badge :color="record.status === 'success' ? 'green' : record.status === 'failed' ? 'red' : 'orange'" :text="record.status === 'success' ? '成功' : record.status === 'failed' ? '失败' : '进行中'" /></template>
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
  { title: '分发ID', dataIndex: 'id', width: 80 },
  { title: '策略名称', dataIndex: 'policy_name', width: 200 },
  { title: '目标设备', dataIndex: 'target_device', width: 160 },
  { title: '分发时间', dataIndex: 'distribute_time', width: 170 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '失败原因', dataIndex: 'error_msg', ellipsis: true }
]

const loadData = async () => {
  loading.value = true
  try {
    const res = await fetch('/api/v1/policies/distribution', {
      headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token') }
    }).then(r => r.json())
    data.value = res.data?.list || []
    pagination.total = data.value.length
  } catch { data.value = [] } finally { loading.value = false }
}

const onPageChange = (page) => { pagination.current = page; loadData() }

onMounted(() => loadData())
</script>
