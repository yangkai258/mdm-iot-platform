<template>
  <div class="container">
    <a-card class="general-card" title="积分规则">
      <template #extra>
        <a-space :size="12">
          <a-button type="primary" @click="openCreate"><icon-plus />新建</a-button>
          <a-button @click="loadRules"><icon-refresh />刷新</a-button>
        </a-space>
      </template>
      <a-row :gutter="16">
        <a-col :span="8">
          <a-form-item label="关键词">
            <a-input v-model="filters.keyword" placeholder="请输入" @pressEnter="loadRules" />
          </a-form-item>
        </a-col>
        <a-col :flex="'86px'" style="display: flex; align-items: flex-end">
          <a-space direction="vertical" :size="8">
            <a-button type="primary" @click="loadRules">查询</a-button>
            <a-button @click="Object.keys(filters).forEach(k => filters[k] = ''); loadRules()">重置</a-button>
          </a-space>
        </a-col>
      </a-row>
      <a-divider style="margin: 0 0 16px 0" />
      <a-table :columns="columns" :data="rules" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id">
        <template #actions="{ record }">
          <a-button type="text" size="small" @click="openEdit(record)">编辑</a-button>
          <a-button type="text" size="small" @click="handleDelete(record)">删除</a-button>
        </template>
      </a-table>
    </a-card>
    <a-modal v-model:visible="modalVisible" :title="isEdit ? '编辑规则' : '新建规则'">
      <a-form :model="form" label-col-flex="100px">
        <a-form-item label="规则名称"><a-input v-model="form.rule_name" /></a-form-item>
        <a-form-item label="规则类型">
          <a-select v-model="form.rule_type" style="width: 100%">
            <a-option value="consume">消费积分</a-option>
            <a-option value="activity">活动积分</a-option>
            <a-option value="birthday">生日积分</a-option>
          </a-select>
        </a-form-item>
      </a-form>
      <template #footer>
        <a-button @click="modalVisible = false">取消</a-button>
        <a-button type="primary" @click="handleSubmit">确定</a-button>
      </template>
    </a-modal>
  </div>
</template>

<script setup>

import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import Breadcrumb from '@/components/Breadcrumb.vue'

const API_BASE = '/api/v1'
const rules = ref([])
const levels = ref([])
const loading = ref(false)
const modalVisible = ref(false)
const detailVisible = ref(false)
const isEdit = ref(false)
const currentRule = ref(null)
const currentId = ref(null)

const filters = reactive({ keyword: '', type: '', status: '' })
const pagination = reactive({ current: 1, pageSize: 10, total: 0 })
const stats = reactive({ total: 0, enabled: 0, triggerCount: 0, issuedPoints: 0 })

const form = reactive({
  rule_name: '', rule_type: 'consume', points_per_yuan: 1,
  level_id: undefined, activity_name: '', points: 0, ratio: 1, remark: ''
})
const formStatus = ref('1')

const columns = [
  { title: '规则名称', dataIndex: 'rule_name' },
  { title: '规则类型', slotName: 'type', width: 120 },
  { title: '积分倍率', slotName: 'ratio', width: 130 },
  { title: '状态', slotName: 'status', width: 90 },
  { title: '创建时间', dataIndex: 'created_at', width: 170 },
  { title: '操作', slotName: 'actions', width: 180, fixed: 'right' }
]

const getTypeColor = (type) => ({ consume: 'blue', activity: 'purple', birthday: 'orange', level: 'green' }[type] || 'gray')
const getTypeText = (type) => ({ consume: '消费积分', activity: '活动积分', birthday: '生日积分', level: '等级倍率' }[type] || type)
const formatTime = (t) => t ? new Date(t).toLocaleString('zh-CN') : '-'

const loadRules = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const params = new URLSearchParams({ page: pagination.current, page_size: pagination.pageSize })
    if (filters.keyword) params.append('keyword', filters.keyword)
    if (filters.type) params.append('type', filters.type)
    if (filters.status !== '') params.append('status', filters.status)

    const res = await fetch(`${API_BASE}/member/points/rules?${params}`, {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    const data = await res.json()
    if (data.code === 0) {
      rules.value = data.data?.list || data.data || []
      pagination.total = data.data?.total || 0
      stats.total = data.data?.total || 0
      stats.enabled = rules.value.filter(r => r.status === 1).length
    }
  } catch (e) {
    Message.error('加载积分规则失败')
  } finally {
    loading.value = false
  }
}

const loadLevels = async () => {
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`${API_BASE}/member/levels`, { headers: { 'Authorization': `Bearer ${token}` } })
    const data = await res.json()
    if (data.code === 0) levels.value = data.data || []
  } catch (e) {}
}

const openCreate = () => {
  isEdit.value = false
  currentId.value = null
  Object.assign(form, { rule_name: '', rule_type: 'consume', points_per_yuan: 1, level_id: undefined, activity_name: '', points: 0, ratio: 1, remark: '' })
  formStatus.value = '1'
  modalVisible.value = true
}

const openEdit = (record) => {
  isEdit.value = true
  currentId.value = record.id
  Object.assign(form, record)
  formStatus.value = String(record.status || 1)
  modalVisible.value = true
}

const openDetail = (record) => {
  currentRule.value = record
  detailVisible.value = true
}

const handleSubmit = async () => {
  if (!form.rule_name) { Message.warning('请填写规则名称'); return }
  try {
    const token = localStorage.getItem('token')
    form.status = parseInt(formStatus.value)
    const method = isEdit.value ? 'PUT' : 'POST'
    const url = isEdit.value ? `${API_BASE}/member/points/rules/${currentId.value}` : `${API_BASE}/member/points/rules`
    const res = await fetch(url, {
      method,
      headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' },
      body: JSON.stringify(form)
    })
    const data = await res.json()
    if (data.code === 0) {
      Message.success(isEdit.value ? '更新成功' : '创建成功')
      modalVisible.value = false
      loadRules()
    } else {
      Message.error(data.message || '操作失败')
    }
  } catch (e) {
    Message.error('操作失败')
  }
}

const handleDelete = async (record) => {
  if (!confirm(`确定删除规则「${record.rule_name}」吗？`)) return
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`${API_BASE}/member/points/rules/${record.id}`, {
      method: 'DELETE', headers: { 'Authorization': `Bearer ${token}` }
    })
    const data = await res.json()
    if (data.code === 0) { Message.success('删除成功'); loadRules() }
    else Message.error(data.message || '删除失败')
  } catch (e) { Message.error('删除失败') }
}

const onPageChange = (page) => { pagination.current = page; loadRules() }
const onTypeChange = () => { /* reset fields */ }

onMounted(() => { loadRules(); loadLevels() })

</script>

<style scoped>
.page-container { background: #fff; border-radius: 4px; padding: 20px; }
.search-form { margin-bottom: 16px; padding: 16px; background: #f7f8fa; border-radius: 4px; }
.toolbar { margin-bottom: 16px; }
</style>
