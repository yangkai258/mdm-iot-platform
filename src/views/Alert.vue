<template>
  <div class="page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>告警管理</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 搜索筛选区 -->
    <div class="search-form">
      <a-form :model="searchForm" layout="inline">
        <a-form-item label="设备ID">
          <a-input v-model="searchForm.device_id" placeholder="请输入设备ID" allow-clear />
        </a-form-item>
        <a-form-item label="告警类型">
          <a-select v-model="searchForm.alert_type" placeholder="选择类型" allow-clear style="width: 160px">
            <a-option value="battery_low">电量过低</a-option>
            <a-option value="offline">设备离线</a-option>
            <a-option value="temperature_high">温度过高</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="状态">
          <a-select v-model="searchForm.status" placeholder="选择状态" allow-clear style="width: 120px">
            <a-option :value="1">未处理</a-option>
            <a-option :value="2">已解决</a-option>
          </a-select>
        </a-form-item>
        <a-form-item>
          <a-space>
            <a-button type="primary" @click="handleSearch">搜索</a-button>
            <a-button @click="handleReset">重置</a-button>
          </a-space>
        </a-form-item>
      </a-form>
    </div>

    <!-- 操作栏 -->
    <div class="toolbar">
      <a-button type="primary" @click="showAddModal">添加规则</a-button>
    </div>

    <!-- 表格 -->
    <a-card title="告警记录" class="table-card">
      <a-table :columns="alertColumns" :data="alerts" :loading="loading" row-key="id">
        <template #status="{ record }">
          <a-tag :color="record.status === 1 ? 'orange' : 'green'">
            {{ record.status === 1 ? '未处理' : '已解决' }}
          </a-tag>
        </template>
        <template #severity="{ record }">
          <a-tag :color="getSeverityColor(record.severity)">
            {{ record.severity }}级
          </a-tag>
        </template>
      </a-table>
    </a-card>

    <!-- 规则管理卡片 -->
    <a-card title="告警规则" class="table-card">
      <template #actions>
        <a-button type="primary" size="small" @click="showAddModal">添加规则</a-button>
      </template>
      <a-table :columns="columns" :data="rules" :loading="loading" row-key="id">
        <template #enabled="{ record }">
          <a-switch v-model="record.enabled" @change="toggleRule(record)" />
        </template>
        <template #severity="{ record }">
          <a-tag :color="getSeverityColor(record.severity)">
            {{ getSeverityText(record.severity) }}
          </a-tag>
        </template>
        <template #condition="{ record }">
          {{ record.condition }} {{ record.threshold }}
        </template>
      </a-table>
    </a-card>

    <!-- 添加规则弹窗 -->
    <a-modal v-model:visible="modalVisible" title="添加告警规则" @ok="addRule">
      <a-form :model="form" layout="vertical">
        <a-form-item label="规则名称" required>
          <a-input v-model="form.name" placeholder="请输入规则名称" />
        </a-form-item>
        <a-form-item label="设备ID（留空表示所有设备）">
          <a-input v-model="form.device_id" placeholder="设备ID" />
        </a-form-item>
        <a-form-item label="告警类型" required>
          <a-select v-model="form.alert_type" placeholder="选择告警类型">
            <a-option value="battery_low">电量过低</a-option>
            <a-option value="offline">设备离线</a-option>
            <a-option value="temperature_high">温度过高</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="条件" required>
          <a-input-group compact>
            <a-select v-model="form.condition" style="width: 80px">
              <a-option value="<">&lt;</a-option>
              <a-option value=">">&gt;</a-option>
              <a-option value="=">=</a-option>
            </a-select>
            <a-input-number v-model="form.threshold" style="width: 120px" />
          </a-input-group>
        </a-form-item>
        <a-form-item label="严重程度" required>
          <a-select v-model="form.severity">
            <a-option :value="1">低</a-option>
            <a-option :value="2">中</a-option>
            <a-option :value="3">高</a-option>
            <a-option :value="4">严重</a-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const loading = ref(false)
const rules = ref([])
const alerts = ref([])
const modalVisible = ref(false)

const searchForm = reactive({
  device_id: '',
  alert_type: '',
  status: ''
})

const form = reactive({
  name: '',
  device_id: '',
  alert_type: 'battery_low',
  condition: '<',
  threshold: 20,
  severity: 2
})

const columns = [
  { title: '规则名称', dataIndex: 'name' },
  { title: '设备ID', dataIndex: 'device_id' },
  { title: '告警类型', dataIndex: 'alert_type' },
  { title: '触发条件', slotName: 'condition' },
  { title: '严重程度', slotName: 'severity' },
  { title: '启用', slotName: 'enabled' }
]

const alertColumns = [
  { title: '设备ID', dataIndex: 'device_id' },
  { title: '告警类型', dataIndex: 'alert_type' },
  { title: '消息', dataIndex: 'message', ellipsis: true },
  { title: '触发值', dataIndex: 'trigger_val' },
  { title: '严重程度', slotName: 'severity' },
  { title: '状态', slotName: 'status' },
  { title: '时间', dataIndex: 'created_at' }
]

const getSeverityColor = (s) => {
  const c = { 1: 'green', 2: 'blue', 3: 'orange', 4: 'red' }
  return c[s] || 'gray'
}

const getSeverityText = (s) => {
  const t = { 1: '低', 2: '中', 3: '高', 4: '严重' }
  return t[s] || '未知'
}

const loadData = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    
    const res = await fetch('/api/v1/alerts/rules', {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    const data = await res.json()
    if (data.code === 0) {
      rules.value = data.data.list || []
    }

    const params = new URLSearchParams()
    if (searchForm.device_id) params.append('device_id', searchForm.device_id)
    if (searchForm.alert_type) params.append('alert_type', searchForm.alert_type)
    if (searchForm.status !== '') params.append('status', searchForm.status)

    const alertRes = await fetch(`/api/v1/alerts?${params.toString()}`, {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    const alertData = await alertRes.json()
    if (alertData.code === 0) {
      alerts.value = alertData.data.list || []
    }
  } catch (e) {
    Message.error('加载失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  loadData()
}

const handleReset = () => {
  searchForm.device_id = ''
  searchForm.alert_type = ''
  searchForm.status = ''
  loadData()
}

const showAddModal = () => {
  modalVisible.value = true
}

const addRule = async () => {
  try {
    const token = localStorage.getItem('token')
    const res = await fetch('/api/v1/alerts/rules', {
      method: 'POST',
      headers: { 
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(form)
    })
    const data = await res.json()
    if (data.code === 0) {
      Message.success('添加成功')
      modalVisible.value = false
      loadData()
    }
  } catch (e) {
    Message.error('添加失败')
  }
}

const toggleRule = async (record) => {
  // Already updated via v-model
}

onMounted(() => {
  loadData()
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

.table-card {
  margin-bottom: 16px;
}
</style>
