<template>
  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>告警管理</a-breadcrumb-item>
      <a-breadcrumb-item>告警规则</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 搜索框 -->
    <div class="pro-search-bar">
      <a-input-search v-model="searchKeyword" placeholder="搜索规则名称" style="width: 280px" search-button @search="loadRules" />
    </div>

    <!-- 操作按钮组 -->
    <div class="pro-action-bar">
      <a-space>
        <a-button type="primary" @click="showAddModal">新建规则</a-button>
        <a-button @click="loadRules">刷新</a-button>
      </a-space>
    </div>

    <!-- 数据表格 -->
    <div class="pro-content-area">
      <a-table :columns="columns" :data="rules" :loading="loading" row-key="id">
        <template #enabled="{ record }">
          <a-switch v-model="record.enabled" @change="toggleRule(record)" />
        </template>
        <template #severity="{ record }">
          <a-tag :color="getSeverityColor(record.severity)">{{ getSeverityText(record.severity) }}</a-tag>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="editRule(record)">编辑</a-button>
            <a-button type="text" size="small" status="danger" @click="deleteRule(record)">删除</a-button>
          </a-space>
        </template>
      </a-table>
    </div>

    <!-- 添加/编辑规则弹窗 -->
    <a-modal v-model:visible="modalVisible" :title="isEdit ? '编辑规则' : '新建规则'" @ok="handleSubmit" :width="520">
      <a-form :model="form" layout="vertical">
        <a-form-item label="规则名称" required>
          <a-input v-model="form.name" placeholder="请输入规则名称" />
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
const modalVisible = ref(false)
const isEdit = ref(false)
const searchKeyword = ref('')

const form = reactive({
  name: '', alert_type: 'battery_low', condition: '<', threshold: 20, severity: 2
})

const columns = [
  { title: '规则名称', dataIndex: 'name' },
  { title: '告警类型', dataIndex: 'alert_type' },
  { title: '条件', dataIndex: 'condition' },
  { title: '阈值', dataIndex: 'threshold' },
  { title: '严重程度', slotName: 'severity' },
  { title: '启用', slotName: 'enabled', width: 100 },
  { title: '操作', slotName: 'actions', width: 150 }
]

const getSeverityColor = (s) => ({ 1: 'green', 2: 'blue', 3: 'orange', 4: 'red' }[s] || 'gray')
const getSeverityText = (s) => ({ 1: '低', 2: '中', 3: '高', 4: '严重' }[s] || '未知')

const loadRules = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const res = await fetch('/api/v1/alerts/rules', { headers: { 'Authorization': `Bearer ${token}` } })
    const data = await res.json()
    if (data.code === 0) rules.value = data.data.list || []
  } catch (e) { Message.error('加载失败') }
  finally { loading.value = false }
}

const showAddModal = () => {
  isEdit.value = false
  Object.assign(form, { name: '', alert_type: 'battery_low', condition: '<', threshold: 20, severity: 2 })
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
    if (data.code === 0) { Message.success('保存成功'); modalVisible.value = false; loadRules() }
  } catch (e) { Message.error('操作失败') }
}

const toggleRule = async (record) => {}
const deleteRule = (record) => { rules.value = rules.value.filter(r => r.id !== record.id); Message.success('删除成功') }

onMounted(() => { loadRules() })
</script>

<style scoped>
.pro-page-container { padding: 20px 24px; min-height: calc(100vh - 64px); background: #f5f7fa; }
.pro-breadcrumb { margin-bottom: 16px; }
.pro-search-bar { margin-bottom: 12px; }
.pro-action-bar { margin-bottom: 16px; }
.pro-content-area {
  background: #fff; border-radius: 8px; padding: 20px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
}
</style>
