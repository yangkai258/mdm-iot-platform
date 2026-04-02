<template>
    <Breadcrumb :items="['Home','Console','']" />


  <div class="container">
    <a-card class="general-card" title="情绪趋势">
      <template #extra>
        <a-space :size="12">
          <a-button type="primary" @click="handleExport"><icon-download />导出趋势</a-button>
          <a-button @click="handleSearch"><icon-refresh />刷新</a-button>
        </a-space>
      </template>
      <a-row :gutter="16">
        <a-col :span="8">
          <a-form-item label="时间周期">
            <a-select v-model="form.period" placeholder="请选择" style="width: 100%">
              <a-option value="week">近一周</a-option>
              <a-option value="month">近一月</a-option>
            </a-select>
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
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id" />
    </a-table>
  </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import Breadcrumb from '@/components/Breadcrumb.vue'

const loading = ref(false)
const data = ref<any[]>([])
const form = reactive({ period: 'week' })
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })
const columns = [
  { title: '日期', dataIndex: 'date', width: 140 },
  { title: '平均强度', dataIndex: 'avg_intensity', width: 120 },
  { title: '主要情绪', dataIndex: 'dominant_emotion', width: 120 },
  { title: '情绪趋势', dataIndex: 'trend', width: 120 },
  { title: '记录数', dataIndex: 'total_records', width: 100 }
]

const loadData = async () => {
  loading.value = true
  try {
    const res = await fetch(`/api/v1/emotions/records/stats?period=${form.period}`, {
      headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token') }
    }).then(r => r.json())
    data.value = res.weekly_data || []
    pagination.total = data.value.length
  } catch { Message.error('加载失败') } finally { loading.value = false }
}

const handleSearch = () => { pagination.current = 1; loadData() }
const handleReset = () => { form.period = 'week'; pagination.current = 1; loadData() }
const onPageChange = (page: number) => { pagination.current = page; loadData() }
const handleExport = () => { Message.info('导出功能开发中') }

onMounted(() => loadData())
</script>

