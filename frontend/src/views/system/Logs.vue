<template>
  <div class="logs-container">
    <a-card>
      <template #title>
        <div class="card-title">
          <span>操作日志</span>
        </div>
      </template>
      
      <a-form :model="query" layout="inline">
        <a-form-item field="username" label="用户名">
          <a-input v-model="query.username" placeholder="请输入用户名" allow-clear style="width: 150px" />
        </a-form-item>
        <a-form-item field="module" label="模块">
          <a-select v-model="query.module" placeholder="选择模块" allow-clear style="width: 150px">
            <a-option value="devices">设备管理</a-option>
            <a-option value="ota">OTA管理</a-option>
            <a-option value="auth">认证</a-option>
          </a-select>
        </a-form-item>
        <a-form-item>
          <a-button type="primary" @click="loadLogs">查询</a-button>
        </a-form-item>
      </a-form>

      <a-table
        :columns="columns"
        :data="logs"
        :loading="loading"
        :pagination="pagination"
        @change="handleTableChange"
        row-key="id"
        style="margin-top: 16px"
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
    </a-card>
  </div>
</template>

<script setup>
import { reactive, ref, onMounted } from 'vue'

const loading = ref(false)
const logs = ref([])

const query = reactive({
  username: '',
  module: ''
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
    const res = await fetch('http://localhost:8080/api/v1/logs/operations', {
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

const handleTableChange = (pag) => {
  pagination.current = pag.current
  loadLogs()
}

onMounted(() => {
  loadLogs()
})
</script>

<style scoped>
.logs-container {
  padding: 16px;
}

.card-title {
  font-size: 16px;
  font-weight: bold;
}
</style>
