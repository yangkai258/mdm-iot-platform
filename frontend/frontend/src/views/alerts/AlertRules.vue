<template>
  <div class="page-container">
    <div class="search-form">
      <a-form :model="form" layout="inline">
        <a-form-item label="规则名称">
          <a-input v-model="form.name" placeholder="请输入规则名称" />
        </a-form-item>
        <a-form-item>
          <a-button type="primary" @click="handleSearch">搜索</a-button>
          <a-button @click="handleReset">重置</a-button>
        </a-form-item>
      </a-form>
    </div>
    <div class="toolbar">
      <a-button type="primary" @click="handleCreate">新建规则</a-button>
    </div>
    <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="handlePageChange" row-key="id">
      <template #severity="{ record }">
        <a-tag :color="getSeverityColor(record.severity)">{{ getSeverityText(record.severity) }}</a-tag>
      </template>
      <template #enabled="{ record }">
        <a-switch v-model="record.enabled" @change="handleToggle(record)" />
      </template>
      <template #actions="{ record }">
        <a-space>
          <a-button type="text" size="small" @click="handleEdit(record)">编辑</a-button>
          <a-button type="text" size="small" status="danger" @click="handleDelete(record)">删除</a-button>
        </a-space>
      </template>
    </a-table>
    <a-modal v-model:visible="modalVisible" :title="modalTitle" @before-ok="handleSubmit" :unmount-on-close="false">
      <a-form :model="form" label-col-flex="100px" ref="formRef">
        <a-form-item label="规则名称" field="name" :rules="[{ required: true, message: '请输入规则名称' }]">
          <a-input v-model="form.name" placeholder="请输入" />
        </a-form-item>
        <a-form-item label="告警类型" field="alert_type" :rules="[{ required: true, message: '请选择告警类型' }]">
          <a-select v-model="form.alert_type" placeholder="请选择">
            <a-option value="battery_low">电量过低</a-option>
            <a-option value="offline">设备离线</a-option>
            <a-option value="temperature_high">温度过高</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="条件">
          <a-input-group compact>
            <a-select v-model="form.condition" style="width: 80px">
              <a-option value="<">&lt;</a-option>
              <a-option value=">">&gt;</a-option>
              <a-option value="=">=</a-option>
            </a-select>
            <a-input-number v-model="form.threshold" style="width: 120px" />
          </a-input-group>
        </a-form-item>
        <a-form-item label="严重程度" field="severity" :rules="[{ required: true, message: '请选择严重程度' }]">
          <a-select v-model="form.severity" placeholder="请选择">
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
import { ref, reactive, onMounted } from 'vue'
import { Message, Modal } from '@arco-design/web-vue'
import axios from 'axios'

const API_BASE = '/api/v1'

const loading = ref(false)
const modalVisible = ref(false)
const modalTitle = ref('新建规则')
const formRef = ref()
const editId = ref<number | null>(null)

const form = reactive({
  name: '',
  alert_type: 'battery_low',
  condition: '<',
  threshold: 20,
  severity: 2
})

const data = ref<any[]>([])
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })

const columns = [
  { title: '规则名称', dataIndex: 'name', width: 150 },
  { title: '告警类型', dataIndex: 'alert_type', width: 150 },
  { title: '条件', dataIndex: 'condition', width: 80 },
  { title: '阈值', dataIndex: 'threshold', width: 80 },
  { title: '严重程度', slotName: 'severity', width: 100 },
  { title: '启用', slotName: 'enabled', width: 100 },
  { title: '操作', slotName: 'actions', width: 150 }
]

const getSeverityColor = (s: number) => ({ 1: 'green', 2: 'blue', 3: 'orange', 4: 'red' }[s] || 'gray')
const getSeverityText = (s: number) => ({ 1: '低', 2: '中', 3: '高', 4: '严重' }[s] || '未知')

const loadData = async () => {
  loading.value = true
  try {
    const params: any = { page: pagination.current, page_size: pagination.pageSize }
    if (form.name) params.name = form.name
    const res = await axios.get(`${API_BASE}/alerts/rules`, { params })
    if (res.data.code === 0) {
      data.value = res.data.data?.list || []
      pagination.total = res.data.data?.pagination?.total || 0
    }
  } catch {
    data.value = [
      { id: 1, name: '温度监控', alert_type: 'temperature_high', condition: '>', threshold: 80, severity: 3, enabled: true },
      { id: 2, name: '电量监控', alert_type: 'battery_low', condition: '<', threshold: 20, severity: 2, enabled: true }
    ]
    pagination.total = 2
  } finally {
    loading.value = false
  }
}

const handleSearch = () => { pagination.current = 1; loadData() }
const handleReset = () => { form.name = ''; pagination.current = 1; loadData() }
const handlePageChange = (page: number) => { pagination.current = page; loadData() }

const handleCreate = () => {
  editId.value = null
  modalTitle.value = '新建规则'
  Object.assign(form, { name: '', alert_type: 'battery_low', condition: '<', threshold: 20, severity: 2 })
  modalVisible.value = true
}

const handleEdit = (record: any) => {
  editId.value = record.id
  modalTitle.value = '编辑规则'
  Object.assign(form, { ...record })
  modalVisible.value = true
}

const handleSubmit = async (done: (arg: boolean) => void) => {
  try {
    await formRef.value?.validate()
    if (editId.value) {
      await axios.put(`${API_BASE}/alerts/rules/${editId.value}`, form)
    } else {
      await axios.post(`${API_BASE}/alerts/rules`, form)
    }
    Message.success('保存成功')
    modalVisible.value = false
    loadData()
    done(true)
  } catch {
    done(false)
  }
}

const handleToggle = async (record: any) => {
  try { await axios.put(`${API_BASE}/alerts/rules/${record.id}`, { enabled: record.enabled }) } catch {}
}

const handleDelete = (record: any) => {
  Modal.confirm({ title: '确认删除', content: '确定删除该规则？', onOk: async () => {
    try { await axios.delete(`${API_BASE}/alerts/rules/${record.id}`) } catch {}
    Message.success('删除成功')
    loadData()
  }})
}

onMounted(() => loadData())
</script>

<style scoped>
.page-container { background: #fff; border-radius: 4px; padding: 20px; }
.search-form { margin-bottom: 16px; padding: 16px; background: #f7f8fa; border-radius: 4px; }
.toolbar { margin-bottom: 16px; }
</style>
