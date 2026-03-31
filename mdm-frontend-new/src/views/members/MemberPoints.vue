<template>
  <div class="page-container">
    <div class="search-form">
      <a-form :model="form" layout="inline">
        <a-form-item label="会员名称">
          <a-input v-model="form.keyword" placeholder="搜索会员名称/手机号" style="width: 200px" />
        </a-form-item>
        <a-form-item label="会员等级">
          <a-select v-model="form.level" placeholder="选择等级" allow-clear style="width: 120px">
            <a-option value="gold">黄金</a-option>
            <a-option value="silver">白银</a-option>
            <a-option value="bronze">青铜</a-option>
          </a-select>
        </a-form-item>
        <a-form-item>
          <a-button type="primary" @click="handleSearch">搜索</a-button>
          <a-button @click="handleReset">重置</a-button>
        </a-form-item>
      </a-form>
    </div>
    <div class="toolbar">
      <a-button type="primary" @click="handleCreate">积分调整</a-button>
    </div>
    <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id">
      <template #level="{ record }"><a-tag :color="getLevelColor(record.level)">{{ getLevelText(record.level) }}</a-tag></template>
      <template #points="{ record }"><span style="color: #ff6b00; font-weight: 600;">{{ record.points || 0 }}</span></template>
      <template #actions="{ record }">
        <a-space>
          <a-button type="text" size="small" @click="handleEdit(record)">调整</a-button>
          <a-button type="text" size="small" @click="handleView(record)">详情</a-button>
        </a-space>
      </template>
    </a-table>
    <a-modal v-model:visible="modalVisible" :title="modalTitle" @ok="handleSubmit" :width="400" :mask-closable="false">
      <a-form :model="form" layout="vertical">
        <a-form-item label="会员" required>
          <a-select v-model="form.member_id" placeholder="选择会员" filterable>
            <a-option v-for="m in data" :key="m.id" :value="m.id">{{ m.member_name }} ({{ m.phone }})</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="调整类型" required>
          <a-radio-group v-model="form.type">
            <a-radio value="add">增加</a-radio>
            <a-radio value="deduct">扣除</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item label="积分数量" required>
          <a-input-number v-model="form.points" :min="1" :max="100000" placeholder="请输入积分数量" style="width: 100%" />
        </a-form-item>
        <a-form-item label="调整原因" required>
          <a-textarea v-model="form.reason" :rows="3" placeholder="请输入调整原因" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const API_BASE = '/api/v1'
const data = ref([])
const loading = ref(false)
const modalVisible = ref(false)
const isEdit = ref(false)

const form = reactive({ keyword: '', level: '', member_id: undefined, type: 'add', points: 0, reason: '' })
const pagination = reactive({ current: 1, pageSize: 10, total: 0 })

const modalTitle = computed(() => isEdit.value ? '编辑积分' : '积分调整')

const columns = [
  { title: '会员名称', dataIndex: 'member_name', width: 150 },
  { title: '手机号', dataIndex: 'phone', width: 130 },
  { title: '会员等级', slotName: 'level', width: 100 },
  { title: '当前积分', slotName: 'points', width: 120 },
  { title: '成长值', dataIndex: 'growth_value', width: 100 },
  { title: '注册时间', dataIndex: 'created_at', width: 160 },
  { title: '操作', slotName: 'actions', width: 120 }
]

const getLevelColor = (level) => ({ gold: '#FFD700', silver: '#C0C0C0', bronze: '#CD7F32' }[level] || 'gray')
const getLevelText = (level) => ({ gold: '黄金', silver: '白银', bronze: '青铜' }[level] || level)

const loadData = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const params = new URLSearchParams({ page: pagination.current, page_size: pagination.pageSize })
    if (form.keyword) params.append('keyword', form.keyword)
    if (form.level) params.append('level', form.level)
    const res = await fetch(`${API_BASE}/member/points?${params}`, { headers: { 'Authorization': `Bearer ${token}` } })
    const resp = await res.json()
    if (resp.code === 0) { data.value = resp.data?.list || resp.data || []; pagination.total = resp.data?.total || 0 }
  } catch (e) { Message.error('加载会员积分失败') } finally { loading.value = false }
}

const handleSearch = () => { pagination.current = 1; loadData() }
const handleReset = () => { form.keyword = ''; form.level = ''; handleSearch() }
const onPageChange = (page) => { pagination.current = page; loadData() }

const handleCreate = () => {
  isEdit.value = false
  Object.assign(form, { member_id: undefined, type: 'add', points: 0, reason: '' })
  modalVisible.value = true
}

const handleEdit = (record) => {
  isEdit.value = true
  Object.assign(form, { member_id: record.id, type: 'add', points: 0, reason: '' })
  modalVisible.value = true
}

const handleView = (record) => { Message.info(`当前积分: ${record.points || 0}`) }

const handleSubmit = async () => {
  if (!form.member_id || !form.points || !form.reason) { Message.warning('请填写完整信息'); return }
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`${API_BASE}/member/points/adjust`, {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' },
      body: JSON.stringify(form)
    })
    const resp = await res.json()
    if (resp.code === 0) { Message.success('积分调整成功'); modalVisible.value = false; loadData() }
    else Message.error(resp.message || '调整失败')
  } catch (e) { Message.error('调整失败') }
}

onMounted(() => loadData())
</script>

<style scoped>
.page-container { background: #fff; border-radius: 4px; padding: 20px; }
.search-form { margin-bottom: 16px; padding: 16px; background: #f7f8fa; border-radius: 4px; }
.toolbar { margin-bottom: 16px; }
</style>
