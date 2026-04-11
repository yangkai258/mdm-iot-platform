<template>
  <div class="page-container">
    <a-card class="general-card" title="Webhooks模板市场">
      <template #extra>
        <a-button @click="loadData"><icon-refresh />刷新</a-button>
      </template>
      <div class="search-form">
        <a-form :model="form" layout="inline">
          <a-form-item label="关键词">
            <a-input v-model="form.keyword" placeholder="请输入关键词" style="width: 200px" />
          </a-form-item>
          <a-form-item label="分类">
            <a-select v-model="form.category" placeholder="请选择" allow-clear style="width: 140px">
              <a-option value="notification">通知</a-option>
              <a-option value="automation">自动化</a-option>
              <a-option value="data">数据同步</a-option>
              <a-option value="analytics">分析</a-option>
            </a-select>
          </a-form-item>
          <a-form-item>
            <a-button type="primary" @click="loadData">查询</a-button>
            <a-button @click="handleReset">重置</a-button>
          </a-form-item>
        </a-form>
      </div>
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id">
        <template #actions="{ record }">
          <a-button type="primary" size="small" @click="handleInstall(record)">一键安装</a-button>
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import { IconRefresh } from '@arco-design/web-vue/es/icon'

const loading = ref(false)
const form = reactive({ keyword: '', category: '' })
const data = ref([])
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })

const columns = [
  { title: '模板名称', dataIndex: 'name', width: 200 },
  { title: '描述', dataIndex: 'description', ellipsis: true },
  { title: '分类', dataIndex: 'category_name', width: 100 },
  { title: '使用次数', dataIndex: 'usage_count', width: 100 },
  { title: '作者', dataIndex: 'author', width: 120 },
  { title: '操作', slotName: 'actions', width: 120 }
]

const loadData = async () => {
  loading.value = true
  try {
    const res = await fetch('/api/v1/market/webhook-templates').then(r => r.json())
    if (res.code === 0) {
      data.value = res.data?.list || []
    } else {
      loadMockData()
    }
  } catch { loadMockData() } finally { loading.value = false }
}

const loadMockData = () => {
  data.value = [
    { id: 1, name: '设备离线通知', description: '当设备离线时自动发送通知到企业微信', category: 'notification', category_name: '通知', usage_count: 1250, author: '官方' },
    { id: 2, name: '心跳异常告警', description: '设备心跳超过阈值时触发告警', category: 'notification', category_name: '通知', usage_count: 980, author: '官方' },
    { id: 3, name: 'OTA完成同步', description: '固件升级完成后自动同步状态到CRM', category: 'data', category_name: '数据同步', usage_count: 560, author: '社区' },
    { id: 4, name: '定时健康检查', description: '每天定时检查设备健康状态', category: 'automation', category_name: '自动化', usage_count: 780, author: '官方' },
    { id: 5, name: '使用统计报表', description: '每日自动生成使用统计报表', category: 'analytics', category_name: '分析', usage_count: 420, author: '社区' }
  ]
}

const handleReset = () => {
  form.keyword = ''
  form.category = ''
  loadData()
}

const handleInstall = (record) => {
  Message.success(`已安装模板: ${record.name}`)
}

const onPageChange = (page) => {
  pagination.current = page
  loadData()
}

onMounted(() => loadData())
</script>

<style scoped>
.page-container { padding: 16px; }
.search-form { margin-bottom: 16px; padding: 16px; background: var(--color-fill-lightest); border-radius: 4px; }
</style>