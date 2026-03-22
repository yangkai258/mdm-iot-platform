<template>
  <div class="pro-page-container">
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>会员管理</a-breadcrumb-item>
      <a-breadcrumb-item>积分规则</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 统计卡片 -->
    <a-row :gutter="16" class="stats-row">
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="积分规则总数" :value="stats.total" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="启用中" :value="stats.enabled" :value-style="{ color: '#52c41a' }" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="本月触发次数" :value="stats.triggerCount" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="本月发放积分" :value="stats.issuedPoints" :value-style="{ color: '#1890ff' }" />
        </a-card>
      </a-col>
    </a-row>

    <!-- 搜索栏 -->
    <div class="pro-search-bar">
      <a-space wrap>
        <a-input-search v-model="filters.keyword" placeholder="搜索规则名称" style="width: 240px" search-button @search="loadRules" />
        <a-select v-model="filters.type" placeholder="规则类型" allow-clear style="width: 140px" @change="loadRules">
          <a-option value="consume">消费积分</a-option>
          <a-option value="activity">活动积分</a-option>
          <a-option value="birthday">生日积分</a-option>
          <a-option value="level">等级积分</a-option>
        </a-select>
        <a-select v-model="filters.status" placeholder="状态" allow-clear style="width: 120px" @change="loadRules">
          <a-option :value="1">启用</a-option>
          <a-option :value="0">禁用</a-option>
        </a-select>
      </a-space>
    </div>

    <!-- 操作栏 -->
    <div class="pro-action-bar">
      <a-space>
        <a-button type="primary" @click="openCreate">新建规则</a-button>
        <a-button @click="loadRules">刷新</a-button>
      </a-space>
    </div>

    <!-- 数据表格 -->
    <div class="pro-content-area">
      <a-table
        :columns="columns"
        :data="rules"
        :loading="loading"
        :pagination="pagination"
        @page-change="onPageChange"
        row-key="id"
      >
        <template #type="{ record }">
          <a-tag :color="getTypeColor(record.rule_type)">{{ getTypeText(record.rule_type) }}</a-tag>
        </template>
        <template #ratio="{ record }">
          <span style="color: #ff6b00; font-weight: 600;">
            {{ record.ratio || record.points_per_yuan || 0 }} {{ record.rule_type === 'consume' ? '积分/元' : '积分' }}
          </span>
        </template>
        <template #status="{ record }">
          <a-tag :color="record.status === 1 ? 'green' : 'gray'">{{ record.status === 1 ? '启用' : '禁用' }}</a-tag>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="openDetail(record)">详情</a-button>
            <a-button type="text" size="small" @click="openEdit(record)">编辑</a-button>
            <a-button type="text" size="small" status="danger" @click="handleDelete(record)">删除</a-button>
          </a-space>
        </template>
      </a-table>
    </div>

    <!-- 新建/编辑弹窗 -->
    <a-modal
      v-model:visible="modalVisible"
      :title="isEdit ? '编辑积分规则' : '新建积分规则'"
      @ok="handleSubmit"
      :width="520"
      :mask-closable="false"
    >
      <a-form :model="form" layout="vertical">
        <a-form-item label="规则名称" required>
          <a-input v-model="form.rule_name" placeholder="请输入规则名称" />
        </a-form-item>
        <a-form-item label="规则类型" required>
          <a-select v-model="form.rule_type" placeholder="选择规则类型" @change="onTypeChange">
            <a-option value="consume">消费积分</a-option>
            <a-option value="activity">活动积分</a-option>
            <a-option value="birthday">生日积分</a-option>
            <a-option value="level">等级倍率</a-option>
          </a-select>
        </a-form-item>
        <a-form-item v-if="form.rule_type === 'consume'" label="每消费1元获得积分" required>
          <a-input-number v-model="form.points_per_yuan" :min="0" :precision="1" placeholder="请输入积分倍率" style="width: 100%" />
        </a-form-item>
        <a-form-item v-if="form.rule_type === 'consume'" label="适用会员等级" required>
          <a-select v-model="form.level_id" placeholder="选择适用等级" allow-clear>
            <a-option v-for="lv in levels" :key="lv.id" :value="lv.id">{{ lv.level_name }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item v-if="form.rule_type === 'activity'" label="活动名称" required>
          <a-input v-model="form.activity_name" placeholder="请输入活动名称" />
        </a-form-item>
        <a-form-item v-if="form.rule_type === 'activity'" label="奖励积分" required>
          <a-input-number v-model="form.points" :min="0" placeholder="请输入奖励积分" style="width: 100%" />
        </a-form-item>
        <a-form-item v-if="form.rule_type === 'birthday'" label="生日奖励积分" required>
          <a-input-number v-model="form.points" :min="0" placeholder="请输入生日奖励积分" style="width: 100%" />
        </a-form-item>
        <a-form-item v-if="form.rule_type === 'level'" label="等级积分倍率" required>
          <a-input-number v-model="form.ratio" :min="1" :precision="1" placeholder="请输入积分倍率" style="width: 100%" />
        </a-form-item>
        <a-form-item label="状态">
          <a-switch v-model="formStatus" checked-value="1" unchecked-value="0" />
        </a-form-item>
        <a-form-item label="备注">
          <a-textarea v-model="form.remark" :rows="2" placeholder="请输入备注" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 详情抽屉 -->
    <a-drawer v-model:visible="detailVisible" title="规则详情" :width="480">
      <template v-if="currentRule">
        <a-descriptions :column="1" bordered size="small">
          <a-descriptions-item label="规则名称">{{ currentRule.rule_name }}</a-descriptions-item>
          <a-descriptions-item label="规则类型">
            <a-tag :color="getTypeColor(currentRule.rule_type)">{{ getTypeText(currentRule.rule_type) }}</a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="积分倍率">{{ currentRule.points_per_yuan || currentRule.ratio || currentRule.points || '-' }}</a-descriptions-item>
          <a-descriptions-item label="状态">
            <a-tag :color="currentRule.status === 1 ? 'green' : 'gray'">{{ currentRule.status === 1 ? '启用' : '禁用' }}</a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="备注">{{ currentRule.remark || '-' }}</a-descriptions-item>
          <a-descriptions-item label="创建时间">{{ formatTime(currentRule.created_at) }}</a-descriptions-item>
        </a-descriptions>
      </template>
    </a-drawer>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

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
.pro-page-container { padding: 20px 24px; min-height: calc(100vh - 64px); background: #f5f7fa; }
.pro-breadcrumb { margin-bottom: 16px; }
.stats-row { margin-bottom: 16px; }
.stat-card { border-radius: 8px; text-align: center; }
.pro-search-bar { margin-bottom: 12px; }
.pro-action-bar { margin-bottom: 16px; }
.pro-content-area { background: #fff; border-radius: 8px; padding: 20px; box-shadow: 0 1px 3px rgba(0,0,0,0.04); }
</style>
