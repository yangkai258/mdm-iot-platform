<template>
  <div class="page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>会员管理</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 搜索筛选区 -->
    <div class="search-form">
      <a-form :model="searchForm" layout="inline">
        <a-form-item label="关键词">
          <a-input-search
            v-model="searchForm.keyword"
            placeholder="搜索会员姓名/编号/手机"
            style="width: 280px"
            @search="handleSearch"
            search-button
          />
        </a-form-item>
        <a-form-item label="会员等级">
          <a-select v-model="searchForm.level" placeholder="选择等级" allow-clear style="width: 160px">
            <a-option v-for="lv in levels" :key="lv.id" :value="lv.level_code">{{ lv.level_name }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item>
          <a-button type="primary" @click="handleSearch">搜索</a-button>
          <a-button @click="handleReset">重置</a-button>
        </a-form-item>
      </a-form>
    </div>

    <!-- 统计卡片 -->
    <a-row :gutter="16" class="stats-row">
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="会员总数" :value="stats.total">
            <template #prefix><icon-user /></template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="活跃会员" :value="stats.active" :value-style="{ color: '#52c41a' }">
            <template #prefix><icon-check-circle /></template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="今日新增" :value="stats.todayNew">
            <template #prefix><icon-user-add /></template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="总积分" :value="stats.totalPoints" :value-style="{ color: '#faad14' }">
            <template #prefix><icon-star /></template>
          </a-statistic>
        </a-card>
      </a-col>
    </a-row>

    <!-- 操作栏 -->
    <div class="toolbar">
      <a-space>
        <a-button type="primary" @click="showAddModal">新增会员</a-button>
        <a-button @click="handleSearch">刷新</a-button>
      </a-space>
    </div>

    <!-- 表格 -->
    <a-table
      :columns="columns"
      :data="members"
      :loading="loading"
      :pagination="pagination"
      @page-change="onPageChange"
      row-key="id"
    >
      <template #avatar="{ record }">
        <a-avatar :style="{ backgroundColor: getLevelColor(record.member_level) }">
          {{ record.member_name?.charAt(0) || '?' }}
        </a-avatar>
      </template>
      <template #member_level="{ record }">
        <a-tag :color="getLevelColor(record.member_level)">{{ record.level_name || record.member_level }}</a-tag>
      </template>
      <template #status="{ record }">
        <a-tag :color="record.status === 1 ? 'green' : 'red'">{{ record.status === 1 ? '正常' : '禁用' }}</a-tag>
      </template>
      <template #actions="{ record }">
        <a-space>
          <a-button type="text" size="small" @click="showEdit(record)">编辑</a-button>
          <a-button type="text" size="small" status="danger" @click="handleDelete(record)">删除</a-button>
        </a-space>
      </template>
    </a-table>

    <!-- 新增/编辑弹窗 -->
    <a-modal v-model:visible="modalVisible" :title="isEdit ? '编辑会员' : '新增会员'" @ok="handleSubmit" :width="520">
      <a-form :model="form" layout="vertical">
        <a-form-item label="会员姓名" required>
          <a-input v-model="form.member_name" placeholder="请输入会员姓名" />
        </a-form-item>
        <a-form-item label="会员编号">
          <a-input v-model="form.member_code" placeholder="请输入会员编号" />
        </a-form-item>
        <a-form-item label="手机号">
          <a-input v-model="form.phone" placeholder="请输入手机号" />
        </a-form-item>
        <a-form-item label="会员等级">
          <a-select v-model="form.member_level" placeholder="请选择会员等级">
            <a-option v-for="lv in levels" :key="lv.id" :value="lv.level_code">{{ lv.level_name }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="状态">
          <a-switch v-model="form.status" checked-value="1" unchecked-value="0" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const API_BASE = '/api/v1'
const members = ref([])
const levels = ref([])
const loading = ref(false)
const modalVisible = ref(false)
const isEdit = ref(false)
const currentId = ref(null)

const searchForm = reactive({
  keyword: '',
  level: ''
})

const form = reactive({
  member_name: '', member_code: '', phone: '', member_level: '', status: '1'
})

const stats = ref({ total: 0, active: 0, todayNew: 0, totalPoints: 0 })

const pagination = reactive({ current: 1, pageSize: 10, total: 0 })

const columns = [
  { title: '头像', slotName: 'avatar', width: 70 },
  { title: '会员姓名', dataIndex: 'member_name' },
  { title: '会员编号', dataIndex: 'member_code' },
  { title: '手机号', dataIndex: 'phone' },
  { title: '等级', slotName: 'member_level' },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '积分', dataIndex: 'points', width: 80 },
  { title: '注册时间', dataIndex: 'created_at', width: 170 },
  { title: '操作', slotName: 'actions', width: 140 }
]

const getLevelColor = (lvl) => {
  const map = { '1': 'gold', '2': 'purple', '3': 'blue', '4': 'red' }
  return map[lvl] || 'gray'
}

const loadMembers = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const params = {
      page: pagination.current,
      page_size: pagination.pageSize
    }
    if (searchForm.keyword) params.keyword = searchForm.keyword
    if (searchForm.level) params.level = searchForm.level

    const res = await fetch(`${API_BASE}/members?page=${params.page}&page_size=${params.page_size}&keyword=${params.keyword || ''}&level=${params.level || ''}`, {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    const data = await res.json()
    if (data.code === 0) {
      members.value = data.data?.list || data.data || []
      pagination.total = data.data?.total || 0
    }
  } catch (e) {
    Message.error('加载会员列表失败')
  } finally {
    loading.value = false
  }
}

const loadMemberStats = async () => {
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`${API_BASE}/members/stats`, {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    const data = await res.json()
    if (data.code === 0) {
      const d = data.data || {}
      stats.value = {
        total: d.total ?? d.total_members ?? 0,
        active: d.active ?? d.active_members ?? 0,
        todayNew: d.todayNew ?? d.today_new ?? d.new_today ?? 0,
        totalPoints: d.totalPoints ?? d.total_points ?? 0
      }
    }
  } catch (e) {
    // stats 加载失败不影响列表
  }
}

const loadLevels = async () => {
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`${API_BASE}/member/levels`, {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    const data = await res.json()
    if (data.code === 0) levels.value = data.data || []
  } catch (e) {}
}

const handleSearch = () => {
  pagination.current = 1
  loadMembers()
}

const handleReset = () => {
  searchForm.keyword = ''
  searchForm.level = ''
  pagination.current = 1
  loadMembers()
}

const showAddModal = () => {
  isEdit.value = false
  Object.assign(form, { member_name: '', member_code: '', phone: '', member_level: '', status: '1' })
  modalVisible.value = true
}

const showEdit = (record) => {
  isEdit.value = true
  currentId.value = record.id
  Object.assign(form, record)
  modalVisible.value = true
}

const handleSubmit = async () => {
  try {
    const token = localStorage.getItem('token')
    const method = isEdit.value ? 'PUT' : 'POST'
    const url = isEdit.value ? `${API_BASE}/members/${currentId.value}` : `${API_BASE}/members`
    const res = await fetch(url, {
      method, headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' },
      body: JSON.stringify(form)
    })
    const data = await res.json()
    if (data.code === 0) {
      Message.success(isEdit.value ? '更新成功' : '创建成功')
      modalVisible.value = false
      loadMembers()
    } else Message.error(data.message || '操作失败')
  } catch (e) { Message.error('操作失败') }
}

const handleDelete = async (record) => {
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`${API_BASE}/members/${record.id}`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    })
    const data = await res.json()
    if (data && data.code === 0) {
      Message.success({ content: '删除成功', id: 'del-' + record.id })
      const idx = members.value.findIndex(m => m.id === record.id)
      if (idx !== -1) members.value.splice(idx, 1)
      pagination.total = Math.max(0, pagination.total - 1)
    } else {
      Message.error({ content: data?.message || '删除失败', id: 'del-' + record.id })
    }
  } catch (e) {
    Message.error({ content: '网络错误，删除失败', id: 'del-' + record.id })
  }
}

const onPageChange = (page) => { pagination.current = page; loadMembers() }

onMounted(() => { loadMembers(); loadLevels(); loadMemberStats() })
</script>

<style scoped>
.page-container {
  background: #fff;
  border-radius: 4px;
  padding: 20px;
}

.breadcrumb {
  margin-bottom: 16px;
}

.search-form {
  margin-bottom: 16px;
  padding: 16px;
  background: #f7f8fa;
  border-radius: 4px;
}

.stats-row {
  margin-bottom: 16px;
}

.stat-card {
  text-align: center;
}

.toolbar {
  margin-bottom: 16px;
}
</style>
