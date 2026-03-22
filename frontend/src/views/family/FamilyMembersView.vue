<template>
  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>家庭管理</a-breadcrumb-item>
      <a-breadcrumb-item>家庭成员</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 页面标题 -->
    <div class="pro-page-header">
      <h2 class="pro-page-title">家庭成员</h2>
    </div>

    <!-- 搜索筛选区 -->
    <div class="pro-search-bar">
      <a-space>
        <a-input-search
          v-model="keyword"
          placeholder="搜索成员姓名/手机号"
          style="width: 280px"
          @search="loadMembers"
          search-button
        />
        <a-select v-model="roleFilter" placeholder="成员角色" allow-clear style="width: 160px" @change="loadMembers">
          <a-option value="owner">户主</a-option>
          <a-option value="adult">成人</a-option>
          <a-option value="child">儿童</a-option>
          <a-option value="elder">老人</a-option>
          <a-option value="guest">访客</a-option>
        </a-select>
        <a-select v-model="statusFilter" placeholder="状态" allow-clear style="width: 120px" @change="loadMembers">
          <a-option value="active">正常</a-option>
          <a-option value="pending">待激活</a-option>
          <a-option value="disabled">已禁用</a-option>
        </a-select>
      </a-space>
    </div>

    <!-- 操作按钮区 -->
    <div class="pro-action-bar">
      <a-space>
        <a-button type="primary" @click="showInviteModal">
          <template #icon><icon-user-add /></template>
          邀请成员
        </a-button>
        <a-button @click="loadMembers">
          <template #icon><icon-refresh /></template>
          刷新
        </a-button>
      </a-space>
    </div>

    <!-- 数据内容区 -->
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
          <a-avatar :style="{ backgroundColor: getRoleColor(record.role) }" :size="40">
            {{ record.name?.charAt(0) || '?' }}
          </a-avatar>
        </template>
        <template #role="{ record }">
          <a-tag :color="getRoleColor(record.role)">{{ getRoleLabel(record.role) }}</a-tag>
        </template>
        <template #status="{ record }">
          <a-tag :color="getStatusColor(record.status)">{{ getStatusLabel(record.status) }}</a-tag>
        </template>
        <template #joined_at="{ record }">
          {{ formatDate(record.joined_at) }}
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="showEditRoleModal(record)">编辑角色</a-button>
            <a-button type="text" size="small" status="danger" @click="handleRemove(record)" :disabled="record.role === 'owner'">
              移除
            </a-button>
          </a-space>
        </template>
      </a-table>
    </div>

    <!-- 邀请成员弹窗 -->
    <a-modal v-model:visible="inviteModalVisible" title="邀请成员" @ok="handleInvite" :width="480" @close="inviteForm = { phone: '', role: 'adult' }">
      <a-form :model="inviteForm" layout="vertical">
        <a-form-item label="手机号码" required>
          <a-input v-model="inviteForm.phone" placeholder="请输入被邀请人手机号" />
        </a-form-item>
        <a-form-item label="成员角色" required>
          <a-select v-model="inviteForm.role" placeholder="请选择角色">
            <a-option value="adult">成人</a-option>
            <a-option value="child">儿童</a-option>
            <a-option value="elder">老人</a-option>
            <a-option value="guest">访客</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="邀请留言">
          <a-textarea v-model="inviteForm.message" placeholder="可选填写邀请留言" :rows="3" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 编辑角色弹窗 -->
    <a-modal v-model:visible="editRoleModalVisible" title="编辑成员角色" @ok="handleUpdateRole" :width="400" @close="editRoleForm = { id: null, role: '' }">
      <a-form :model="editRoleForm" layout="vertical">
        <a-form-item label="成员角色" required>
          <a-select v-model="editRoleForm.role" placeholder="请选择角色">
            <a-option value="adult">成人</a-option>
            <a-option value="child">儿童</a-option>
            <a-option value="elder">老人</a-option>
            <a-option value="guest">访客</a-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const members = ref<any[]>([])
const loading = ref(false)
const keyword = ref('')
const roleFilter = ref('')
const statusFilter = ref('')
const inviteModalVisible = ref(false)
const editRoleModalVisible = ref(false)

const inviteForm = reactive({
  phone: '',
  role: 'adult',
  message: ''
})

const editRoleForm = reactive({
  id: null as number | null,
  role: ''
})

const columns = [
  { title: '头像', dataIndex: 'avatar', slotName: 'avatar', width: 80 },
  { title: '姓名', dataIndex: 'name', width: 140 },
  { title: '手机号', dataIndex: 'phone', width: 140 },
  { title: '角色', dataIndex: 'role', slotName: 'role', width: 120 },
  { title: '状态', dataIndex: 'status', slotName: 'status', width: 100 },
  { title: '加入时间', dataIndex: 'joined_at', slotName: 'joined_at', width: 180 },
  { title: '操作', slotName: 'actions', width: 180 }
]

const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0
})

const roleColors: Record<string, string> = {
  owner: '#1650ff',
  adult: '#00b42a',
  child: '#ff7d00',
  elder: '#8e4e9c',
  guest: '#86909c'
}

const roleLabels: Record<string, string> = {
  owner: '户主',
  adult: '成人',
  child: '儿童',
  elder: '老人',
  guest: '访客'
}

const statusLabels: Record<string, string> = {
  active: '正常',
  pending: '待激活',
  disabled: '已禁用'
}

function getRoleColor(role: string) {
  return roleColors[role] || '#86909c'
}

function getRoleLabel(role: string) {
  return roleLabels[role] || role
}

function getStatusColor(status: string) {
  const map: Record<string, string> = { active: 'green', pending: 'orange', disabled: 'red' }
  return map[status] || 'gray'
}

function getStatusLabel(status: string) {
  return statusLabels[status] || status
}

function formatDate(dateStr: string) {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleString('zh-CN')
}

async function loadMembers() {
  loading.value = true
  try {
    const params = new URLSearchParams()
    if (keyword.value) params.append('keyword', keyword.value)
    if (roleFilter.value) params.append('role', roleFilter.value)
    if (statusFilter.value) params.append('status', statusFilter.value)
    params.append('page', String(pagination.current))
    params.append('page_size', String(pagination.pageSize))

    const res = await fetch(`/api/v1/family/members?${params}`, {
      credentials: 'include'
    })
    const data = await res.json()
    if (data.code === 0 || data.code === 200) {
      members.value = data.data?.list || data.data?.members || []
      pagination.total = data.data?.total || 0
    } else {
      Message.error(data.message || '加载失败')
    }
  } catch {
    Message.error('网络错误，请稍后重试')
  } finally {
    loading.value = false
  }
}

function onPageChange(page: number) {
  pagination.current = page
  loadMembers()
}

function showInviteModal() {
  inviteForm.phone = ''
  inviteForm.role = 'adult'
  inviteForm.message = ''
  inviteModalVisible.value = true
}

async function handleInvite() {
  if (!inviteForm.phone) {
    Message.warning('请输入手机号码')
    return
  }
  try {
    const res = await fetch('/api/v1/family/members/invite', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify(inviteForm)
    })
    const data = await res.json()
    if (data.code === 0 || data.code === 200) {
      Message.success('邀请已发送')
      inviteModalVisible.value = false
      loadMembers()
    } else {
      Message.error(data.message || '邀请失败')
    }
  } catch {
    Message.error('网络错误')
  }
}

function showEditRoleModal(record: any) {
  editRoleForm.id = record.id
  editRoleForm.role = record.role
  editRoleModalVisible.value = true
}

async function handleUpdateRole() {
  if (!editRoleForm.id || !editRoleForm.role) {
    Message.warning('请选择角色')
    return
  }
  try {
    const res = await fetch(`/api/v1/family/members/${editRoleForm.id}/role`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({ role: editRoleForm.role })
    })
    const data = await res.json()
    if (data.code === 0 || data.code === 200) {
      Message.success('角色更新成功')
      editRoleModalVisible.value = false
      loadMembers()
    } else {
      Message.error(data.message || '更新失败')
    }
  } catch {
    Message.error('网络错误')
  }
}

async function handleRemove(record: any) {
  try {
    const res = await fetch(`/api/v1/family/members/${record.id}`, {
      method: 'DELETE',
      credentials: 'include'
    })
    const data = await res.json()
    if (data.code === 0 || data.code === 200) {
      Message.success('成员已移除')
      loadMembers()
    } else {
      Message.error(data.message || '移除失败')
    }
  } catch {
    Message.error('网络错误')
  }
}

onMounted(() => {
  loadMembers()
})
</script>
