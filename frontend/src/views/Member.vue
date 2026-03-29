<template>
  <div class="member-page-container">
    <Breadcrumb :items="[{ label: '首页', href: '/' }, { label: '会员管理' }]" />

    <a-card class="general-card">
      <template #title>
        <span class="card-title">会员统计</span>
      </template>
      <a-row :gutter="16">
        <a-col :span="6">
          <a-statistic title="会员总数" :value="stats.total" :value-style="{ color: '#1650d8' }">
            <template #prefix>
              <icon-user style="font-size: 20px; margin-right: 6px; color: #1650d8" />
            </template>
          </a-statistic>
        </a-col>
        <a-col :span="6">
          <a-statistic title="活跃会员" :value="stats.active" :value-style="{ color: '#52c41a' }">
            <template #prefix>
              <icon-check-circle style="font-size: 20px; margin-right: 6px; color: #52c41a" />
            </template>
          </a-statistic>
        </a-col>
        <a-col :span="6">
          <a-statistic title="今日新增" :value="stats.todayNew" :value-style="{ color: '#1650d8' }">
            <template #prefix>
              <icon-user-add style="font-size: 20px; margin-right: 6px; color: #1650d8" />
            </template>
          </a-statistic>
        </a-col>
        <a-col :span="6">
          <a-statistic title="总积分" :value="stats.totalPoints" :value-style="{ color: '#faad14' }">
            <template #prefix>
              <icon-star style="font-size: 20px; margin-right: 6px; color: #faad14" />
            </template>
          </a-statistic>
        </a-col>
      </a-row>
    </a-card>

    <a-card class="general-card" style="margin-top: 16px">
      <template #title>
        <span class="card-title">会员查询</span>
      </template>
      <a-row :gutter="16">
        <a-col :flex="1">
          <a-form :model="searchForm" layout="vertical" size="small">
            <a-row :gutter="16">
              <a-col :span="8">
                <a-form-item label="关键词">
                  <a-input v-model="searchForm.keyword" placeholder="会员姓名/编号/手机" allow-clear />
                </a-form-item>
              </a-col>
              <a-col :span="8">
                <a-form-item label="会员等级">
                  <a-select v-model="searchForm.level" placeholder="全部" allow-clear>
                    <a-option v-for="lv in levels" :key="lv.id" :value="lv.level_code">
                      {{ lv.level_name }}
                    </a-option>
                  </a-select>
                </a-form-item>
              </a-col>
              <a-col :span="8">
                <a-form-item label="状态">
                  <a-select v-model="searchForm.status" placeholder="全部" allow-clear>
                    <a-option :value="1">正常</a-option>
                    <a-option :value="0">禁用</a-option>
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
      <template #title>
        <span class="card-title">会员列表</span>
      </template>
      <template #extra>
        <a-space>
          <a-button type="primary" @click="showAddModal">
            <template #icon><icon-plus /></template>
            新增会员
          </a-button>
          <a-button @click="handleExport">
            <template #icon><icon-download /></template>
            导出
          </a-button>
        </a-space>
      </template>

      <a-table
        :columns="columns"
        :data="members"
        :loading="loading"
        :pagination="paginationConfig"
        @page-change="onPageChange"
        @page-size-change="onPageSizeChange"
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
    </a-card>

    <!-- 新增/编辑弹窗 -->
    <a-modal
      v-model:visible="modalVisible"
      :title="isEdit ? '编辑会员' : '新增会员'"
      @ok="handleSubmit"
      :width="520"
    >
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
import { ref, reactive, computed, onMounted } from 'vue'
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
  level: undefined,
  status: undefined
})

const form = reactive({
  member_name: '', member_code: '', phone: '', member_level: '', status: '1'
})

const stats = ref({ total: 0, active: 0, todayNew: 0, totalPoints: 0 })

const pagination = reactive({ current: 1, pageSize: 10, total: 0 })

const paginationConfig = computed(() => ({
  current: pagination.current,
  pageSize: pagination.pageSize,
  total: pagination.total,
  showTotal: true,
  showSizeChanger: true
}))

const columns = [
  { title: '头像', slotName: 'avatar', width: 70 },
  { title: '会员姓名', dataIndex: 'member_name', width: 120 },
  { title: '会员编号', dataIndex: 'member_code', width: 140 },
  { title: '手机号', dataIndex: 'phone', width: 130 },
  { title: '等级', slotName: 'member_level', width: 100 },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '积分', dataIndex: 'points', width: 80 },
  { title: '注册时间', dataIndex: 'created_at', width: 170 },
  { title: '操作', slotName: 'actions', width: 180, fixed: 'right' }
]

const getLevelColor = (lvl) => {
  const map = { '1': 'gold', '2': 'purple', '3': 'blue', '4': 'red' }
  return map[lvl] || 'gray'
}

const loadMembers = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(
      `${API_BASE}/members?page=${pagination.current}&page_size=${pagination.pageSize}&keyword=${searchForm.keyword}&level=${searchForm.level || ''}`,
      { headers: { 'Authorization': `Bearer ${token}` } }
    )
    const data = await res.json()
    if (data.code === 0) {
      members.value = data.data.list || []
      pagination.total = data.data.total || 0
    }
  } catch (e) {
    members.value = getMockMembers()
    pagination.total = members.value.length
  } finally {
    loading.value = false
  }
}

const getMockMembers = () => [
  { id: 1, member_name: '张三', member_code: 'M001', phone: '13800138000', member_level: '2', level_name: '黄金会员', status: 1, points: 5200, created_at: '2024-01-15 10:30:00' },
  { id: 2, member_name: '李四', member_code: 'M002', phone: '13800138001', member_level: '3', level_name: '铂金会员', status: 1, points: 12000, created_at: '2024-02-20 14:20:00' },
  { id: 3, member_name: '王五', member_code: 'M003', phone: '13800138002', member_level: '1', level_name: '普通会员', status: 0, points: 800, created_at: '2024-03-05 09:15:00' }
]

const loadLevels = async () => {
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`${API_BASE}/member/levels`, {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    const data = await res.json()
    if (data.code === 0) levels.value = data.data || []
  } catch (e) {
    levels.value = [
      { id: 1, level_code: '1', level_name: '普通会员' },
      { id: 2, level_code: '2', level_name: '黄金会员' },
      { id: 3, level_code: '3', level_name: '铂金会员' },
      { id: 4, level_code: '4', level_name: '钻石会员' }
    ]
  }
}

const loadStats = async () => {
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`${API_BASE}/dashboard/member-stats`, {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    const data = await res.json()
    if (data.code === 0) stats.value = data.data || { total: 0, active: 0, todayNew: 0, totalPoints: 0 }
  } catch (e) {
    stats.value = { total: 3, active: 2, todayNew: 0, totalPoints: 18000 }
  }
}

const handleSearch = () => {
  pagination.current = 1
  loadMembers()
}

const handleReset = () => {
  searchForm.keyword = ''
  searchForm.level = undefined
  searchForm.status = undefined
  pagination.current = 1
  loadMembers()
}

const onPageChange = (page) => {
  pagination.current = page
  loadMembers()
}

const onPageSizeChange = (size) => {
  pagination.pageSize = size
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

const showDetail = (record) => showEdit(record)

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
    } else Message.error(data.message || '操作失败')
  } catch (e) {
    if (isEdit.value) {
      const idx = members.value.findIndex(m => m.id === currentId.value)
      if (idx !== -1) members.value[idx] = { ...members.value[idx], ...form }
    } else {
      members.value.unshift({ ...form, id: Date.now(), points: 0, created_at: new Date().toLocaleString() })
    }
    Message.success(isEdit.value ? '更新成功（模拟）' : '创建成功（模拟）')
    modalVisible.value = false
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
    if (data.code === 0) { Message.success('删除成功'); loadMembers() }
    else Message.error(data.message || '删除失败')
  } catch (e) {
    members.value = members.value.filter(m => m.id !== record.id)
    Message.success('删除成功（模拟）')
  }
}

const handleExport = () => {
  Message.info('导出功能开发中')
}

onMounted(() => { loadMembers(); loadLevels(); loadStats() })
</script>

<style scoped>
.member-page-container {
  padding: 20px 24px;
  min-height: calc(100vh - 64px);
  background: #f5f7fa;
}

.general-card {
  border-radius: 8px;
}

.card-title {
  font-weight: 600;
  font-size: 15px;
}
</style>
