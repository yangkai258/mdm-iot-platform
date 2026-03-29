<template>
  <div class="alert-rules-container">

    <a-card class="general-card">
      <template #title><span class="card-title">规则查询</span></template>
      <a-row :gutter="16">
        <a-col :flex="1">
          <a-form :model="searchForm" layout="vertical" size="small">
            <a-row :gutter="16">
              <a-col :span="8">
                <a-form-item label="规则名称">
                  <a-input v-model="searchForm.keyword" placeholder="搜索规则名称" allow-clear />
                </a-form-item>
              </a-col>
              <a-col :span="8">
                <a-form-item label="严重程度">
                  <a-select v-model="searchForm.severity" placeholder="全部" allow-clear>
                    <a-option :value="1">低</a-option>
                    <a-option :value="2">中</a-option>
                    <a-option :value="3">高</a-option>
                    <a-option :value="4">严重</a-option>
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

    <a-card class="general-card" style="margin-top: 16px">
      <template #title><span class="card-title">告警规则</span></template>
      <template #extra>
        <a-button type="primary" @click="handleCreate">
          <template #icon><icon-plus /></template>
          新建规则
        </a-button>
      </template>
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="paginationConfig" @page-change="onPageChange" row-key="id">
        <template #severity="{ record }">
          <a-tag :color="getSeverityColor(record.severity)">{{ getSeverityText(record.severity) }}</a-tag>
        </template>
      </a-table>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="handleEdit(record)">编辑</a-button>
            <a-button type="text" size="small" status="danger" @click="handleDelete(record)">删除</a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:visible="modalVisible" :title="modalTitle" @ok="handleSubmit" :width="520">
      <a-form :model="form" layout="vertical">
        <a-form-item label="规则名称">
          <a-input v-model="form.name" placeholder="请输入规则名称" />
        </a-form-item>
        <a-form-item label="告警类型">
          <a-input v-model="form.alert_type" placeholder="如 battery_low" />
        </a-form-item>
        <a-form-item label="条件">
          <a-input-group compact>
            <a-select v-model="form.condition" style="width: 80px">
              <a-option value="<">&lt;</a-option>
              <a-option value=">">&gt;</a-option>
              <a-option value="=">=</a-option>
            </a-select>
            <a-input-number v-model="form.threshold" :min="0" style="width: 120px" />
          </a-input-group>
        </a-form-item>
        <a-form-item label="严重程度">
          <a-input-number v-model="form.severity" :min="1" :max="4" style="width: 120px" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const loading = ref(false)
const modalVisible = ref(false)
const isEdit = ref(false)
const modalTitle = computed(() => isEdit.value ? '编辑规则' : '新建规则')

const searchForm = reactive({ keyword: '', severity: undefined })
const form = reactive({ id: null, name: '', alert_type: '', condition: '<', threshold: 20, severity: 2 })
const data = ref([])
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })
const paginationConfig = computed(() => ({ current: pagination.current, pageSize: pagination.pageSize, total: pagination.total, showTotal: true }))

const columns = [
  { title: 'ID', dataIndex: 'id', width: 80 },
  { title: '规则名称', dataIndex: 'name', ellipsis: true },
  { title: '告警类型', dataIndex: 'alert_type', width: 130 },
  { title: '条件', dataIndex: 'condition', width: 80 },
  { title: '阈值', dataIndex: 'threshold', width: 80 },
  { title: '严重程度', slotName: 'severity', width: 100 },
  { title: '操作', slotName: 'actions', width: 150, fixed: 'right' }
]

const getSeverityColor = (s) => ({ 1: 'green', 2: 'blue', 3: 'orange', 4: 'red' }[s] || 'gray')
const getSeverityText = (s) => ({ 1: '低', 2: '中', 3: '高', 4: '严重' }[s] || '未知')

const getMockData = () => [
  { id: 1, name: '电量过低告警', alert_type: 'battery_low', condition: '<', threshold: 20, severity: 3 },
  { id: 2, name: '设备离线告警', alert_type: 'offline', condition: '=', threshold: 1, severity: 4 }
]

const loadData = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`${'/api/v1'}/alerts/rules?page=${pagination.current}&page_size=${pagination.pageSize}`, {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    const resData = await res.json()
    if (resData.code === 0) { data.value = resData.data?.list || []; pagination.total = resData.data?.total || 0 }
  } catch (e) { data.value = getMockData(); pagination.total = data.value.length }
  finally { loading.value = false }
}

const handleSearch = () => { pagination.current = 1; loadData() }
const handleReset = () => { searchForm.keyword = ''; searchForm.severity = undefined; pagination.current = 1; loadData() }
const onPageChange = (page) => { pagination.current = page; loadData() }
const handleCreate = () => { isEdit.value = false; Object.assign(form, { id: null, name: '', alert_type: '', condition: '<', threshold: 20, severity: 2 }); modalVisible.value = true }
const handleEdit = (record) => { isEdit.value = true; Object.assign(form, record); modalVisible.value = true }
const handleSubmit = async () => {
  try {
    const token = localStorage.getItem('token')
    const method = isEdit.value ? 'PUT' : 'POST'
    const url = isEdit.value ? `${'/api/v1'}/alerts/rules/${form.id}` : `${'/api/v1'}/alerts/rules`
    await fetch(url, { method, headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' }, body: JSON.stringify(form) })
    Message.success('保存成功'); modalVisible.value = false; loadData()
  } catch (e) { Message.error('操作失败') }
}
const handleDelete = (record) => { data.value = data.value.filter(d => d.id !== record.id); Message.success('删除成功') }

onMounted(() => { loadData() })
</script>

<style scoped>
.alert-rules-container { padding: 20px 24px; min-height: calc(100vh - 64px); background: #f5f7fa; }
.general-card { border-radius: 8px; }
.card-title { font-weight: 600; font-size: 15px; }
</style>
