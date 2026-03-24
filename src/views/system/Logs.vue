<template>
  <div class="page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>系统管理</a-breadcrumb-item>
      <a-breadcrumb-item>操作日志</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 搜索筛选区 -->
    <div class="search-form">
      <a-form :model="searchForm" layout="inline">
        <a-form-item label="用户名">
          <a-input v-model="searchForm.username" placeholder="请输入用户名" allow-clear style="width: 150px" />
        </a-form-item>
        <a-form-item label="模块">
          <a-select v-model="searchForm.module" placeholder="选择模块" allow-clear style="width: 150px">
            <a-option value="devices">设备管理</a-option>
            <a-option value="ota">OTA管理</a-option>
            <a-option value="auth">认证</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="状态">
          <a-select v-model="searchForm.status" placeholder="选择状态" allow-clear style="width: 120px">
            <a-option :value="1">成功</a-option>
            <a-option :value="2">失败</a-option>
          </a-select>
        </a-form-item>
        <a-form-item>
          <a-space>
            <a-button type="primary" @click="handleSearch">查询</a-button>
            <a-button @click="handleReset">重置</a-button>
          </a-space>
        </a-form-item>
      </a-form>
    </div>

    <!-- 操作栏 -->
    <div class="toolbar">
      <a-button @click="handleExport">导出日志</a-button>
    </div>

    <!-- 表格 -->
    <a-table
      :columns="columns"
      :data="logs"
      :loading="loading"
      :pagination="pagination"
      @change="handleTableChange"
      row-key="id"
    >
      <template #status="{ record }">
        <a-tag :color="record.status === 1 ? 'green' : 'red'">
          {{ record.status === 1 ? '成功' : '失败' }}
        </a-tag>
      </template>
      <template #duration="{ record }">
        {{ record.duration }}ms
      </template>
    </a-table>
  </div>
</template>

<script setup>
import { reactive, ref, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const loading = ref(false)
const logs = ref([])

const searchForm = reactive({
  username: '',
  module: '',
  status: ''
})

const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0
})

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '用户名', dataIndex: 'username', width: 100 },
  { title: '模块', dataIndex: 'module', width: 100 },
  { title: '操作', dataIndex: 'operation', ellipsis: true },
  { title: '方法', dataIndex: 'method', width: 60 },
  { title: 'IP', dataIndex: 'ip', width: 120 },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '耗时', slotName: 'duration', width: 80 },
  { title: '时间', dataIndex: 'created_at', width: 180 }
]

const loadLogs = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const res = await fetch('/api/v1/logs/operations', {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    const data = await res.json()
    if (data.code === 0) {
      logs.value = data.data?.list || []
      pagination.total = logs.value.length
    }
  } catch (e) {
    console.error('加载日志失败:', e)
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.current = 1
  loadLogs()
}

const handleReset = () => {
  searchForm.username = ''
  searchForm.module = ''
  searchForm.status = ''
  pagination.current = 1
  loadLogs()
}

const handleTableChange = (pag) => {
  pagination.current = pag.current
  loadLogs()
}

const handleExport = () => {
  Message.success('日志导出功能开发中')
}

onMounted(() => {
  loadLogs()
})
</script>

<style scoped>
.page-container {
  background: #fff;
  border-radius: 4px;
  padding: 20px;
}

.breadcrumb {
  margin-bottom: 16px;
}

.search-form {
  margin-bottom: 16px;
  padding: 16px;
  background: #f7f8fa;
  border-radius: 4px;
}

.toolbar {
  margin-bottom: 16px;
}
</style>
