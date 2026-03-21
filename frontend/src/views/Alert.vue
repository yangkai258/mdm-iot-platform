<template>
  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>告警管理</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 搜索框 -->
    <div class="pro-search-bar">
      <a-space>
        <a-input-search
          v-model="searchKeyword"
          placeholder="搜索告警内容"
          style="width: 280px"
          @search="loadAlerts"
          search-button
        />
        <a-select v-model="filterSeverity" placeholder="严重级别" allow-clear style="width: 120px" @change="loadAlerts">
          <a-option value="critical">严重</a-option>
          <a-option value="high">高</a-option>
          <a-option value="medium">中</a-option>
          <a-option value="low">低</a-option>
        </a-select>
      </a-space>
    </div>

    <!-- 操作按钮组 -->
    <div class="pro-action-bar">
      <a-space>
        <a-button type="primary" @click="showAddModal">新建规则</a-button>
        <a-button @click="loadAlerts">刷新</a-button>
      </a-space>
    </div>

    <!-- 告警规则表格 -->
    <div class="pro-content-area">
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
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="editRule(record)">编辑</a-button>
            <a-button type="text" size="small" status="danger" @click="deleteRule(record)">删除</a-button>
          </a-space>
        </template>
      </a-table>
    </div>

    <!-- 告警记录 -->
    <a-card class="alert-record-card">
      <template #title>
        <span>告警记录</span>
      </template>
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

    <!-- 添加/编辑规则弹窗 -->
    <a-modal v-model:visible="modalVisible" :title="isEdit ? '编辑规则' : '添加规则'" @ok="handleSubmit" :width="520">
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
const isEdit = ref(false)
const searchKeyword = ref('')
const filterSeverity = ref('')

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
  { title: '启用', slotName: 'enabled' },
  { title: '操作', slotName: 'actions', width: 150 }
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

const loadAlerts = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const res = await fetch('/api/v1/alerts/rules', {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    const data = await res.json()
    if (data.code === 0) rules.value = data.data.list || []

    const alertRes = await fetch('/api/v1/alerts', {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    const alertData = await alertRes.json()
    if (alertData.code === 0) alerts.value = alertData.data.list || []
  } catch (e) {
    Message.error('加载失败')
  } finally {
    loading.value = false
  }
}

const showAddModal = () => {
  isEdit.value = false
  Object.assign(form, { name: '', device_id: '', alert_type: 'battery_low', condition: '<', threshold: 20, severity: 2 })
  modalVisible.value = true
}

const editRule = (record) => {
  isEdit.value = true
  Object.assign(form, record)
  modalVisible.value = true
}

const handleSubmit = async () => {
  try {
    const token = localStorage.getItem('token')
    const res = await fetch('/api/v1/alerts/rules', {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' },
      body: JSON.stringify(form)
    })
    const data = await res.json()
    if (data.code === 0) {
      Message.success('保存成功')
      modalVisible.value = false
      loadAlerts()
    }
  } catch (e) {
    Message.error('操作失败')
  }
}

const toggleRule = async (record) => {}
const deleteRule = (record) => {
  rules.value = rules.value.filter(r => r.id !== record.id)
  Message.success('删除成功')
}

onMounted(() => {
  loadAlerts()
})
</script>

<style scoped>
.pro-page-container {
  padding: 20px 24px;
  min-height: calc(100vh - 64px);
  background: #f5f7fa;
}

.pro-breadcrumb {
  margin-bottom: 16px;
}

.pro-search-bar {
  margin-bottom: 12px;
}

.pro-action-bar {
  margin-bottom: 16px;
}

.pro-content-area {
  background: #fff;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
  margin-bottom: 16px;
}

.alert-record-card {
  border-radius: 8px;
}
</style>
