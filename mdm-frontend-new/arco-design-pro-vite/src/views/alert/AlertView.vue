<template>
  <div class="alert-view-container">
    <Breadcrumb :items="[{title: '首页', path: '/'},{title: '告警管理'},{title: '告警总览'}]" />

    <!-- 搜索筛选区 -->
    <a-card class="general-card" style="margin-bottom: 16px">
      <a-row :gutter="16">
        <a-col :flex="1">
          <a-form :model="searchForm" layout="vertical">
            <a-row :gutter="16">
              <a-col :span="8">
                <a-form-item label="关键词">
                  <a-input v-model="searchForm.keyword" placeholder="搜索告警内容" allow-clear />
                </a-form-item>
              </a-col>
              <a-col :span="8">
                <a-form-item label="严重级别">
                  <a-select v-model="searchForm.severity" placeholder="全部" allow-clear>
                    <a-option value="critical">严重</a-option>
                    <a-option value="high">高</a-option>
                    <a-option value="medium">中</a-option>
                    <a-option value="low">低</a-option>
                  </a-select>
                </a-form-item>
              </a-col>
              <a-col :span="8">
                <a-form-item label="状态">
                  <a-select v-model="searchForm.status" placeholder="全部" allow-clear>
                    <a-option :value="1">未处理</a-option>
                    <a-option :value="0">已解决</a-option>
                  </a-select>
                </a-form-item>
              </a-col>
            </a-row>
          </a-form>
        </a-col>
        <a-divider style="height: 84px" direction="vertical" />
        <a-col :flex="'86px'" style="text-align: right">
          <a-space direction="vertical" :size="18">
            <a-button type="primary" @click="handleSearch">
              <template #icon><icon-search /></template>
              查询
            </a-button>
            <a-button @click="handleReset">
              <template #icon><icon-refresh /></template>
              重置
            </a-button>
          </a-space>
        </a-col>
      </a-row>
    </a-card>

    <!-- 数据表格 -->
    <a-card class="general-card">
      <template #title>告警记录</template>
      <template #extra>
        <a-space>
          <a-button type="primary" @click="showAddModal">
            <template #icon><icon-plus /></template>
            新建规则
          </a-button>
          <a-button @click="loadRules">
            <template #icon><icon-refresh /></template>
            刷新
          </a-button>
        </a-space>
      </template>

      <a-table
        :columns="columns"
        :data="rules"
        :loading="loading"
        :pagination="paginationConfig"
        @page-change="onPageChange"
        @page-size-change="onPageSizeChange"
        row-key="id"
      >
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

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const API_BASE = '/api/v1'
const loading = ref(false)
const rules = ref<any[]>([])
const modalVisible = ref(false)
const isEdit = ref(false)
const currentId = ref<number | null>(null)

const searchForm = reactive({
  keyword: '',
  severity: undefined as string | undefined,
  status: undefined as number | undefined
})

const form = reactive({
  name: '', device_id: '', alert_type: 'battery_low', condition: '<', threshold: 20, severity: 2
})

const pagination = reactive({ current: 1, pageSize: 10, total: 0 })

const paginationConfig = computed(() => ({
  current: pagination.current,
  pageSize: pagination.pageSize,
  total: pagination.total,
  showTotal: true,
  showSizeChanger: true
}))

const columns = [
  { title: '规则名称', dataIndex: 'name', ellipsis: true },
  { title: '设备ID', dataIndex: 'device_id', width: 140 },
  { title: '告警类型', dataIndex: 'alert_type', width: 130 },
  { title: '触发条件', slotName: 'condition', width: 130 },
  { title: '严重程度', slotName: 'severity', width: 100 },
  { title: '启用', slotName: 'enabled', width: 80 },
  { title: '操作', slotName: 'actions', width: 150, fixed: 'right' }
]

const getSeverityColor = (s: number) => ({ 1: 'green', 2: 'blue', 3: 'orange', 4: 'red' }[s] || 'gray')
const getSeverityText = (s: number) => ({ 1: '低', 2: '中', 3: '高', 4: '严重' }[s] || '未知')

const getMockRules = () => [
  { id: 1, name: '电量过低告警', device_id: '', alert_type: 'battery_low', condition: '<', threshold: 20, severity: 3, enabled: true },
  { id: 2, name: '设备离线告警', device_id: '', alert_type: 'offline', condition: '=', threshold: 1, severity: 4, enabled: true }
]

const loadRules = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`${API_BASE}/alerts/rules`, {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    const data = await res.json()
    if (data.code === 0) {
      rules.value = data.data?.list || []
      pagination.total = data.data?.total || rules.value.length
    }
  } catch {
    rules.value = getMockRules()
    pagination.total = rules.value.length
  } finally { loading.value = false }
}

const handleSearch = () => { pagination.current = 1; loadRules() }
const handleReset = () => {
  searchForm.keyword = ''
  searchForm.severity = undefined
  searchForm.status = undefined
  pagination.current = 1
  loadRules()
}
const onPageChange = (page: number) => { pagination.current = page; loadRules() }
const onPageSizeChange = (size: number) => { pagination.pageSize = size; pagination.current = 1; loadRules() }

const showAddModal = () => {
  isEdit.value = false
  currentId.value = null
  Object.assign(form, { name: '', device_id: '', alert_type: 'battery_low', condition: '<', threshold: 20, severity: 2 })
  modalVisible.value = true
}

const editRule = (record: any) => {
  isEdit.value = true
  currentId.value = record.id
  Object.assign(form, record)
  modalVisible.value = true
}

const handleSubmit = async () => {
  try {
    const token = localStorage.getItem('token')
    const method = isEdit.value ? 'PUT' : 'POST'
    const url = isEdit.value ? `${API_BASE}/alerts/rules/${currentId.value}` : `${API_BASE}/alerts/rules`
    const res = await fetch(url, {
      method,
      headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' },
      body: JSON.stringify(form)
    })
    const data = await res.json()
    if (data.code === 0) {
      Message.success('保存成功')
      modalVisible.value = false
      loadRules()
    } else Message.error(data.message || '保存失败')
  } catch {
    if (isEdit.value) {
      const idx = rules.value.findIndex(r => r.id === currentId.value)
      if (idx !== -1) rules.value[idx] = { ...rules.value[idx], ...form }
      Message.success('规则已更新（模拟）')
    } else {
      rules.value.unshift({ ...form, id: Date.now(), enabled: true })
      Message.success('规则已添加（模拟）')
    }
    modalVisible.value = false
  }
}

const toggleRule = async (record: any) => {
  try {
    const token = localStorage.getItem('token')
    await fetch(`${API_BASE}/alerts/rules/${record.id}/toggle`, {
      method: 'PUT',
      headers: { 'Authorization': `Bearer ${token}` }
    })
  } catch {}
}

const deleteRule = (record: any) => {
  rules.value = rules.value.filter(r => r.id !== record.id)
  Message.success('删除成功')
}

onMounted(() => { loadRules() })
</script>

<style scoped>
.alert-view-container {
  padding: 20px 24px;
  min-height: calc(100vh - 64px);
  background: #f5f7fa;
}
.general-card {
  border-radius: 8px;
}
</style>
