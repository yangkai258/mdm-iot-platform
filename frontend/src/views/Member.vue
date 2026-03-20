<template>
  <div class="member-container">
    <a-row :gutter="16" class="toolbar">
      <a-col :span="8">
        <a-input-search v-model="keyword" placeholder="搜索会员姓名/编号/手机" search-button @search="loadMembers" />
      </a-col>
      <a-col :span="8">
        <a-select v-model="level" placeholder="会员等级" allow-clear style="width: 160px" @change="loadMembers">
          <a-option v-for="lv in levels" :key="lv.id" :value="lv.level_code">{{ lv.level_name }}</a-option>
        </a-select>
      </a-col>
      <a-col :span="8" style="text-align: right;">
        <a-button type="primary" @click="showAddModal">
          <template #icon><icon-plus /></template>
          新增会员
        </a-button>
      </a-col>
    </a-row>

    <a-row :gutter="16" style="margin-top: 16px;">
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

    <a-table :columns="columns" :data="members" :loading="loading" :pagination="pagination" @page-change="onPageChange" style="margin-top: 16px;" row-key="id">
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
        <a-button type="text" size="small" @click="showDetail(record)">详情</a-button>
        <a-button type="text" size="small" @click="showEdit(record)">编辑</a-button>
        <a-button type="text" size="small" status="danger" @click="handleDelete(record)">删除</a-button>
      </template>
    </a-table>

    <!-- 新增/编辑弹窗 -->
    <a-modal v-model:visible="modalVisible" :title="isEdit ? '编辑会员' : '新增会员'" @ok="handleSubmit" @cancel="modalVisible = false" :width="520">
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
  member_name: '',
  member_code: '',
  phone: '',
  member_level: '',
  status: '1'
})

const stats = ref({ total: 0, active: 0, todayNew: 0, totalPoints: 0 })

const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0
})

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
    if (data.code === 0) {
      levels.value = data.data || []
    }
  } catch (e) {
    // ignore
  }
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

const showDetail = (record) => {
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
      method,
      headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' },
      body: JSON.stringify(form)
    })
    const data = await res.json()
    if (data.code === 0) {
      Message.success(isEdit.value ? '更新成功' : '创建成功')
      modalVisible.value = false
      loadMembers()
    } else {
      Message.error(data.message || '操作失败')
    }
  } catch (e) {
    Message.error('操作失败')
  }
}

const handleDelete = async (record) => {
  if (!confirm('确定要删除该会员吗？')) return
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`${API_BASE}/members/${record.id}`, {
      method: 'DELETE',
      headers: { 'Authorization': `Bearer ${token}` }
    })
    const data = await res.json()
    if (data.code === 0) {
      Message.success('删除成功')
      loadMembers()
    } else {
      Message.error(data.message || '删除失败')
    }
  } catch (e) {
    Message.error('删除失败')
  }
}

const onPageChange = (page) => {
  pagination.current = page
  loadMembers()
}

onMounted(() => {
  loadMembers()
  loadLevels()
})
</script>

<style scoped>
.member-container {
  padding: 16px;
}
.toolbar {
  align-items: center;
}
.stat-card {
  border-radius: 8px;
  text-align: center;
}
</style>
