<template>
  <div class="pro-page-container">
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>会员管理</a-breadcrumb-item>
      <a-breadcrumb-item>临时会员红包</a-breadcrumb-item>
    </a-breadcrumb>

    <a-row :gutter="16" class="stats-row">
      <a-col :span="8"><a-card class="stat-card"><a-statistic title="红包活动总数" :value="stats.total" /></a-card></a-col>
      <a-col :span="8"><a-card class="stat-card"><a-statistic title="进行中" :value="stats.active" :value-style="{ color: '#52c41a' }" /></a-card></a-col>
      <a-col :span="8"><a-card class="stat-card"><a-statistic title="已发放金额" :value="'￥' + stats.issuedAmount" /></a-card></a-col>
    </a-row>

    <div class="pro-search-bar">
      <a-space wrap>
        <a-input-search v-model="filters.keyword" placeholder="搜索红包名称" style="width: 240px" search-button @search="loadData" />
        <a-select v-model="filters.status" placeholder="状态" allow-clear style="width: 120px" @change="loadData">
          <a-option value="active">进行中</a-option>
          <a-option value="paused">已暂停</a-option>
          <a-option value="finished">已结束</a-option>
        </a-select>
      </a-space>
    </div>

    <div class="pro-action-bar">
      <a-space>
        <a-button type="primary" @click="openCreate">新建红包</a-button>
        <a-button @click="loadData">刷新</a-button>
      </a-space>
    </div>

    <div class="pro-content-area">
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id">
        <template #amount="{ record }"><span style="color: #ff6b00; font-weight: 600;">￥{{ record.amount }}</span></template>
        <template #remain="{ record }"><span>{{ record.remain_count }}/{{ record.total_count }}</span></template>
        <template #status="{ record }"><a-tag :color="getStatusColor(record.status)">{{ getStatusText(record.status) }}</a-tag></template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="openEdit(record)">编辑</a-button>
            <a-button type="text" size="small" @click="openGrantModal(record)" v-if="record.status === 'active'">发放</a-button>
            <a-button type="text" size="small" status="danger" @click="handleDelete(record)">删除</a-button>
          </a-space>
        </template>
      </a-table>
    </div>

    <a-modal v-model:visible="modalVisible" :title="isEdit ? '编辑红包' : '新建红包'" @ok="handleSubmit" :width="520" :mask-closable="false">
      <a-form :model="form" layout="vertical">
        <a-form-item label="红包名称" required><a-input v-model="form.redpacket_name" placeholder="请输入红包名称" /></a-form-item>
        <a-form-item label="红包金额" required>
          <a-input-number v-model="form.amount" :min="0.01" :precision="2" style="width: 100%;" placeholder="请输入红包金额" />
        </a-form-item>
        <a-form-item label="发放总数量" required>
          <a-input-number v-model="form.total_count" :min="1" style="width: 100%;" placeholder="请输入发放总数量" />
        </a-form-item>
        <a-form-item label="有效期">
          <a-range-picker v-model="form.dateRange" style="width: 100%;" />
        </a-form-item>
        <a-form-item label="使用门槛">
          <a-input-number v-model="form.min_amount" :min="0" :precision="2" style="width: 100%;" placeholder="最低消费金额，0为无门槛" />
        </a-form-item>
        <a-form-item label="备注"><a-textarea v-model="form.remark" :rows="2" placeholder="请输入备注" /></a-form-item>
      </a-form>
    </a-modal>

    <a-modal v-model:visible="grantModal" title="发放红包" @ok="handleGrant" :width="480" :mask-closable="false">
      <a-form :model="grantForm" layout="vertical">
        <a-form-item label="红包名称"><span style="font-weight: 600;">{{ grantRecord?.redpacket_name }}</span></a-form-item>
        <a-form-item label="红包金额"><span style="color: #ff6b00; font-weight: 600;">￥{{ grantRecord?.amount }}</span></a-form-item>
        <a-form-item label="剩余数量"><span>{{ grantRecord?.remain_count }}</span></a-form-item>
        <a-divider />
        <a-form-item label="发放方式">
          <a-radio-group v-model="grantForm.scope">
            <a-radio value="single">指定临时会员</a-radio>
            <a-radio value="batch">批量发放</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item v-if="grantForm.scope === 'single'" label="选择临时会员" required>
          <a-select v-model="grantForm.temp_member_id" placeholder="选择临时会员" filterable>
            <a-option v-for="m in tempMembers" :key="m.id" :value="m.id">{{ m.name }} ({{ m.phone }})</a-option>
          </a-select>
        </a-form-item>
        <a-form-item v-if="grantForm.scope === 'batch'" label="发放数量">
          <a-input-number v-model="grantForm.count" :min="1" :max="grantRecord?.remain_count || 1" style="width: 200px;" />
          <span style="margin-left: 8px; color: #999;">将发放给多个临时会员</span>
        </a-form-item>
      </a-form>
      <template #footer>
        <a-space>
          <a-button @click="grantModal = false">取消</a-button>
          <a-button type="primary" @click="handleGrant">确认发放</a-button>
        </a-space>
      </template>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const API_BASE = '/api/v1'
const data = ref([])
const tempMembers = ref([])
const loading = ref(false)
const modalVisible = ref(false)
const grantModal = ref(false)
const isEdit = ref(false)
const currentId = ref(null)
const grantRecord = ref(null)

const filters = reactive({ keyword: '', status: '' })
const pagination = reactive({ current: 1, pageSize: 10, total: 0 })
const stats = reactive({ total: 0, active: 0, issuedAmount: 0 })
const form = reactive({ redpacket_name: '', amount: 0, total_count: 0, min_amount: 0, dateRange: [], remark: '' })
const grantForm = reactive({ scope: 'single', temp_member_id: undefined, count: 1 })

const columns = [
  { title: 'ID', dataIndex: 'id', width: 80 },
  { title: '红包名称', dataIndex: 'redpacket_name' },
  { title: '金额', slotName: 'amount', width: 100 },
  { title: '剩余数量', slotName: 'remain', width: 120 },
  { title: '有效期', dataIndex: 'expire_time', width: 170 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '操作', slotName: 'actions', width: 200 }
]

const getStatusColor = (s) => ({ active: 'green', paused: 'orange', finished: 'gray' }[s] || 'gray')
const getStatusText = (s) => ({ active: '进行中', paused: '已暂停', finished: '已结束' }[s] || s)

const loadData = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const params = new URLSearchParams({ page: pagination.current, page_size: pagination.pageSize })
    if (filters.keyword) params.append('keyword', filters.keyword)
    if (filters.status) params.append('status', filters.status)
    const res = await fetch(`${API_BASE}/member/temp-redpackets?${params}`, { headers: { 'Authorization': `Bearer ${token}` } })
    const resp = await res.json()
    if (resp.code === 0) {
      data.value = resp.data?.list || resp.data || []
      pagination.total = resp.data?.total || 0
      stats.total = pagination.total
      stats.active = data.value.filter(d => d.status === 'active').length
      stats.issuedAmount = data.value.reduce((sum, d) => sum + (d.total_count - d.remain_count) * (d.amount || 0), 0)
    }
  } catch (e) { Message.error('加载红包失败') } finally { loading.value = false }
}

const loadTempMembers = async () => {
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`${API_BASE}/member/temp-members?page_size=100&status=pending`, { headers: { 'Authorization': `Bearer ${token}` } })
    const d = await res.json()
    if (d.code === 0) tempMembers.value = d.data?.list || []
  } catch (e) {}
}

const openCreate = () => { isEdit.value = false; currentId.value = null; Object.assign(form, { redpacket_name: '', amount: 0, total_count: 0, min_amount: 0, dateRange: [], remark: '' }); modalVisible.value = true }
const openEdit = (r) => { isEdit.value = true; currentId.value = r.id; Object.assign(form, r); modalVisible.value = true }
const openGrantModal = (r) => { grantRecord.value = r; grantForm.scope = 'single'; grantForm.temp_member_id = undefined; grantForm.count = 1; grantModal.value = true }

const handleSubmit = async () => {
  if (!form.redpacket_name || !form.amount || !form.total_count) { Message.warning('请填写完整信息'); return }
  try {
    const token = localStorage.getItem('token')
    const method = isEdit.value ? 'PUT' : 'POST'
    const url = isEdit.value ? `${API_BASE}/member/temp-redpackets/${currentId.value}` : `${API_BASE}/member/temp-redpackets`
    const res = await fetch(url, { method, headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' }, body: JSON.stringify(form) })
    const d = await res.json()
    if (d.code === 0) { Message.success(isEdit.value ? '更新成功' : '创建成功'); modalVisible.value = false; loadData() }
    else Message.error(d.message || '操作失败')
  } catch (e) { Message.error('操作失败') }
}

const handleGrant = async () => {
  if (grantForm.scope === 'single' && !grantForm.temp_member_id) { Message.warning('请选择临时会员'); return }
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`${API_BASE}/member/temp-redpackets/${grantRecord.value.id}/grant`, {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' },
      body: JSON.stringify(grantForm)
    })
    const d = await res.json()
    if (d.code === 0) { Message.success('发放成功'); grantModal.value = false; loadData() }
    else Message.error(d.message || '发放失败')
  } catch (e) { Message.error('发放失败') }
}

const handleDelete = async (r) => {
  if (!confirm(`确定删除红包「${r.redpacket_name}」吗？`)) return
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`${API_BASE}/member/temp-redpackets/${r.id}`, { method: 'DELETE', headers: { 'Authorization': `Bearer ${token}` } })
    const d = await res.json()
    if (d.code === 0) { Message.success('删除成功'); loadData() } else Message.error(d.message || '删除失败')
  } catch (e) { Message.error('删除失败') }
}

const onPageChange = (page) => { pagination.current = page; loadData() }
onMounted(() => { loadData(); loadTempMembers() })
</script>

<style scoped>
.pro-page-container { padding: 20px 24px; min-height: calc(100vh - 64px); background: #f5f7fa; }
.pro-breadcrumb { margin-bottom: 16px; }
.stats-row { margin-bottom: 16px; }
.stat-card { border-radius: 8px; text-align: center; }
.pro-search-bar { margin-bottom: 12px; }
.pro-action-bar { margin-bottom: 16px; }
.pro-content-area { background: #fff; border-radius: 8px; padding: 20px; box-shadow: 0 1px 3px rgba(0,0,0,0.04); }
</style>
