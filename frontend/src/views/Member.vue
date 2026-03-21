<template>
  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>会员管理</a-breadcrumb-item>
    </a-breadcrumb>

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

    <!-- 搜索框 -->
    <div class="pro-search-bar">
      <a-space>
        <a-input-search
          v-model="keyword"
          placeholder="搜索会员姓名/编号/手机"
          style="width: 280px"
          @search="loadMembers"
          search-button
        />
        <a-select v-model="level" placeholder="会员等级" allow-clear style="width: 160px" @change="loadMembers">
          <a-option v-for="lv in levels" :key="lv.id" :value="lv.level_code">{{ lv.level_name }}</a-option>
        </a-select>
      </a-space>
    </div>

    <!-- 操作按钮组 -->
    <div class="pro-action-bar">
      <a-space>
        <a-button type="primary" @click="showAddModal">新增会员</a-button>
        <a-button @click="loadMembers">刷新</a-button>
      </a-space>
    </div>

    <!-- 数据表格 -->
    <div class="pro-content-area">
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
            <a-button type="text" size="small" @click="showDetail(record)">详情</a-button>
            <a-button type="text" size="small" @click="showEdit(record)">编辑</a-button>
            <a-button type="text" size="small" status="danger" @click="handleDelete(record)">删除</a-button>
          </a-space>
        </template>
      </a-table>
    </div>

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
const keyword = ref('')
const level = ref('')
const modalVisible = ref(false)
const isEdit = ref(false)
const currentId = ref(null)

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
  { title: '操作', slotName: 'actions', width: 180 }
]

const getLevelColor = (lvl) => {
  const map = { '1': 'gold', '2': 'purple', '3': 'blue', '4': 'red' }
  return map[lvl] || 'gray'
}

const loadMembers = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`${API_BASE}/members?page=${pagination.current}&page_size=${pagination.pageSize}&keyword=${keyword.value}&level=${level.value}`, {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    const data = await res.json()
    if (data.code === 0) {
      members.value = data.data.list || []
      pagination.total = data.data.total || 0
    }
  } catch (e) {
    Message.error('加载会员列表失败')
  } finally {
    loading.value = false
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

const showDetail = (record) => showEdit(record)

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
  if (!confirm('确定要删除该会员吗？')) return
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`${API_BASE}/members/${record.id}`, {
      method: 'DELETE', headers: { 'Authorization': `Bearer ${token}` }
    })
    const data = await res.json()
    if (data.code === 0) { Message.success('删除成功'); loadMembers() }
    else Message.error(data.message || '删除失败')
  } catch (e) { Message.error('删除失败') }
}

const onPageChange = (page) => { pagination.current = page; loadMembers() }

onMounted(() => { loadMembers(); loadLevels() })
</script>

<style scoped>
.pro-page-container {
  padding: 20px 24px;
  min-height: calc(100vh - 64px);
  background: #f5f7fa;
}

.pro-breadcrumb { margin-bottom: 16px; }
.stats-row { margin-bottom: 16px; }
.stat-card { border-radius: 8px; text-align: center; }
.pro-search-bar { margin-bottom: 12px; }
.pro-action-bar { margin-bottom: 16px; }
.pro-content-area {
  background: #fff;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
}
</style>
