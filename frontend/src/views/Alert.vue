<template>
  <div class="alert-container">
    <a-card>
      <template #title>
        <div class="card-title">
          <span>告警规则管理</span>
          <a-button type="primary" @click="showAddModal">添加规则</a-button>
        </div>
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

    <a-card title="告警记录" style="margin-top: 16px;">
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
    
    const res = await fetch('http://localhost:8080/api/v1/alerts/rules', {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    const data = await res.json()
    if (data.code === 0) {
      rules.value = data.data.list || []
    }

    const alertRes = await fetch('http://localhost:8080/api/v1/alerts', {
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

const showAddModal = () => {
  modalVisible.value = true
}

const addRule = async () => {
  try {
    const token = localStorage.getItem('token')
    const res = await fetch('http://localhost:8080/api/v1/alerts/rules', {
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
.alert-container {
  padding: 16px;
}
.card-title {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
