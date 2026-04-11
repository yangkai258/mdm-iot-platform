<template>
  <div class="page-container">
    <a-card class="general-card" title="SDK发布">
      <template #extra>
        <a-button type="primary"><icon-plus />发布SDK</a-button>
      </template>
      <a-tabs v-model:activeTab="activeTab" @change="loadData">
        <a-tab-pane key="all" tab="全部版本" />
        <a-tab-pane key="latest" tab="最新版" />
        <a-tab-pane key="deprecated" tab="已废弃" />
      </a-tabs>
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id">
        <template #status="{ record }">
          <a-tag :color="getStatusColor(record.status)">{{ getStatusText(record.status) }}</a-tag>
        </template>
        <template #actions="{ record }">
          <a-button type="primary" size="small" @click="handleDownload(record)">下载</a-button>
          <a-button size="small" @click="handleViewDocs(record)">文档</a-button>
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import { IconPlus } from '@arco-design/web-vue/es/icon'

const loading = ref(false)
const activeTab = ref('all')
const data = ref([])
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })

const columns = [
  { title: 'SDK名称', dataIndex: 'name', width: 160 },
  { title: '版本', dataIndex: 'version', width: 100 },
  { title: '语言', dataIndex: 'language', width: 80 },
  { title: '发布说明', dataIndex: 'changelog', ellipsis: true },
  { title: '下载量', dataIndex: 'downloads', width: 100 },
  { title: '发布日期', dataIndex: 'release_date', width: 120 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '操作', slotName: 'actions', width: 160 }
]

const getStatusColor = (status) => {
  const colors = { latest: 'green', stable: 'blue', deprecated: 'orange' }
  return colors[status] || 'gray'
}

const getStatusText = (status) => {
  const texts = { latest: '最新', stable: '稳定', deprecated: '已废弃' }
  return texts[status] || status
}

const loadData = async () => {
  loading.value = true
  try {
    const res = await fetch('/api/v1/market/sdk').then(r => r.json())
    if (res.code === 0) {
      data.value = res.data?.list || []
    } else {
      loadMockData()
    }
  } catch { loadMockData() } finally { loading.value = false }
}

const loadMockData = () => {
  data.value = [
    { id: 1, name: 'Go SDK', version: 'v2.3.0', language: 'Go', changelog: '新增设备注册API,优化连接池性能', downloads: 12500, release_date: '2026-03-20', status: 'latest' },
    { id: 2, name: 'Python SDK', version: 'v1.8.5', language: 'Python', changelog: '修复MQTT重连bug,增加示例代码', downloads: 9800, release_date: '2026-03-18', status: 'latest' },
    { id: 3, name: 'JavaScript SDK', version: 'v3.1.2', language: 'JavaScript', changelog: '支持TypeScript,优化打包体积', downloads: 15600, release_date: '2026-03-15', status: 'stable' },
    { id: 4, name: 'Java SDK', version: 'v2.5.0', language: 'Java', changelog: 'Spring Boot 3.0支持,增加WebSocket', downloads: 8200, release_date: '2026-03-10', status: 'stable' },
    { id: 5, name: 'Go SDK', version: 'v1.5.0', language: 'Go', changelog: '旧版本已废弃', downloads: 5600, release_date: '2025-12-01', status: 'deprecated' }
  ]
}

const handleDownload = (record) => {
  Message.success(`下载SDK: ${record.name} ${record.version}`)
}

const handleViewDocs = (record) => {
  Message.info(`查看文档: ${record.name}`)
}

const onPageChange = (page) => {
  pagination.current = page
  loadData()
}

onMounted(() => loadData())
</script>

<style scoped>
.page-container { padding: 16px; }
</style>