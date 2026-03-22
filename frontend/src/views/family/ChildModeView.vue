<template>
  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>家庭管理</a-breadcrumb-item>
      <a-breadcrumb-item>儿童模式</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 页面标题 -->
    <div class="pro-page-header">
      <h2 class="pro-page-title">儿童模式</h2>
      <p class="pro-page-desc">为儿童配置安全的内容过滤和使用时间限制</p>
    </div>

    <!-- 搜索筛选区 -->
    <div class="pro-search-bar">
      <a-space>
        <a-select v-model="childFilter" placeholder="选择儿童账号" allow-clear style="width: 200px" @change="loadChildModes">
          <a-option v-for="c in children" :key="c.id" :value="c.id">{{ c.name }}</a-option>
        </a-select>
        <a-select v-model="statusFilter" placeholder="模式状态" allow-clear style="width: 140px" @change="loadChildModes">
          <a-option value="enabled">已启用</a-option>
          <a-option value="disabled">已禁用</a-option>
        </a-select>
      </a-space>
    </div>

    <!-- 操作按钮区 -->
    <div class="pro-action-bar">
      <a-space>
        <a-button type="primary" @click="showConfigModal(null)">
          <template #icon><icon-settings /></template>
          新增配置
        </a-button>
        <a-button @click="loadChildModes">
          <template #icon><icon-refresh /></template>
          刷新
        </a-button>
      </a-space>
    </div>

    <!-- 数据内容区 -->
    <div class="pro-content-area">
      <a-table
        :columns="columns"
        :data="childModes"
        :loading="loading"
        :pagination="pagination"
        @page-change="onPageChange"
        row-key="id"
      >
        <template #child_name="{ record }">
          <a-avatar :style="{ backgroundColor: '#ff7d00' }" :size="32">
            {{ record.child_name?.charAt(0) || '?' }}
          </a-avatar>
          <span style="margin-left: 8px">{{ record.child_name }}</span>
        </template>
        <template #enabled="{ record }">
          <a-switch v-model="record.enabled" @change="toggleMode(record)" :disabled="savingId === record.id" />
        </template>
        <template #content_filter="{ record }">
          <a-tag :color="record.content_filter_enabled ? 'green' : 'gray'">
            {{ record.content_filter_enabled ? '已启用' : '已禁用' }}
          </a-tag>
        </template>
        <template #time_limit="{ record }">
          <span v-if="record.time_limit_enabled">
            {{ record.daily_time_limit }}分钟/天
          </span>
          <a-tag v-else color="gray">未限制</a-tag>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="showConfigModal(record)">配置</a-button>
            <a-button type="text" size="small" status="danger" @click="handleDelete(record)">删除</a-button>
          </a-space>
        </template>
      </a-table>
    </div>

    <!-- 配置弹窗 -->
    <a-modal v-model:visible="configModalVisible" :title="isEdit ? '编辑儿童模式配置' : '新增儿童模式配置'" @ok="handleSave" :width="560" @close="resetForm">
      <a-form :model="form" layout="vertical">
        <a-form-item label="儿童账号" required>
          <a-select v-model="form.child_id" placeholder="请选择儿童账号" :disabled="isEdit">
            <a-option v-for="c in children" :key="c.id" :value="c.id">{{ c.name }} ({{ c.phone }})</a-option>
          </a-select>
        </a-form-item>

        <a-divider>内容过滤</a-divider>

        <a-form-item label="启用内容过滤">
          <a-switch v-model="form.content_filter_enabled" />
        </a-form-item>
        <a-form-item label="过滤级别" v-if="form.content_filter_enabled">
          <a-radio-group v-model="form.filter_level">
            <a-radio value="strict">严格（适合低龄儿童）</a-radio>
            <a-radio value="moderate">中等（适合学龄儿童）</a-radio>
            <a-radio value="light">轻度（适合青少年）</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item label="禁止访问的分类" v-if="form.content_filter_enabled">
          <a-checkbox-group v-model="form.blocked_categories">
            <a-checkbox value="violence">暴力内容</a-checkbox>
            <a-checkbox value="adult">成人内容</a-checkbox>
            <a-checkbox value="gambling">赌博</a-checkbox>
            <a-checkbox value="social">社交媒体</a-checkbox>
            <a-checkbox value="games">游戏</a-checkbox>
          </a-checkbox-group>
        </a-form-item>

        <a-divider>使用时间限制</a-divider>

        <a-form-item label="启用时间限制">
          <a-switch v-model="form.time_limit_enabled" />
        </a-form-item>
        <a-form-item label="每日时长限制（分钟）" v-if="form.time_limit_enabled">
          <a-input-number v-model="form.daily_time_limit" :min="15" :max="480" :step="15" style="width: 200px" />
        </a-form-item>
        <a-form-item label="允许使用的时间段" v-if="form.time_limit_enabled">
          <a-space direction="vertical">
            <a-time-picker-range v-model="form.allowed_time_range" format="HH:mm" style="width: 300px" />
          </a-space>
        </a-form-item>
        <a-form-item label="禁用日期" v-if="form.time_limit_enabled">
          <a-checkbox-group v-model="form.disabled_days">
            <a-checkbox :value="1">周一</a-checkbox>
            <a-checkbox :value="2">周二</a-checkbox>
            <a-checkbox :value="3">周三</a-checkbox>
            <a-checkbox :value="4">周四</a-checkbox>
            <a-checkbox :value="5">周五</a-checkbox>
            <a-checkbox :value="6">周六</a-checkbox>
            <a-checkbox :value="0">周日</a-checkbox>
          </a-checkbox-group>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const childModes = ref<any[]>([])
const children = ref<any[]>([])
const loading = ref(false)
const savingId = ref<number | null>(null)
const childFilter = ref('')
const statusFilter = ref('')
const configModalVisible = ref(false)
const isEdit = ref(false)

const form = reactive({
  id: null as number | null,
  child_id: null as number | null,
  content_filter_enabled: false,
  filter_level: 'moderate',
  blocked_categories: [] as string[],
  time_limit_enabled: false,
  daily_time_limit: 60,
  allowed_time_range: null as any,
  disabled_days: [] as number[]
})

const columns = [
  { title: '儿童账号', dataIndex: 'child_name', slotName: 'child_name', width: 200 },
  { title: '模式开关', dataIndex: 'enabled', slotName: 'enabled', width: 120 },
  { title: '内容过滤', dataIndex: 'content_filter', slotName: 'content_filter', width: 120 },
  { title: '时间限制', dataIndex: 'time_limit', slotName: 'time_limit', width: 160 },
  { title: '操作', slotName: 'actions', width: 160 }
]

const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0
})

async function loadChildren() {
  try {
    const res = await fetch('/api/v1/family/members?role=child', { credentials: 'include' })
    const data = await res.json()
    children.value = data.data?.list || data.data?.members || []
  } catch { /* ignore */ }
}

async function loadChildModes() {
  loading.value = true
  try {
    const params = new URLSearchParams()
    if (childFilter.value) params.append('child_id', childFilter.value)
    if (statusFilter.value) params.append('status', statusFilter.value)
    params.append('page', String(pagination.current))
    params.append('page_size', String(pagination.pageSize))

    const res = await fetch(`/api/v1/family/child-mode?${params}`, { credentials: 'include' })
    const data = await res.json()
    if (data.code === 0 || data.code === 200) {
      childModes.value = data.data?.list || data.data || []
      pagination.total = data.data?.total || 0
    } else {
      Message.error(data.message || '加载失败')
    }
  } catch {
    Message.error('网络错误')
  } finally {
    loading.value = false
  }
}

function onPageChange(page: number) {
  pagination.current = page
  loadChildModes()
}

function showConfigModal(record: any) {
  if (record) {
    isEdit.value = true
    form.id = record.id
    form.child_id = record.child_id
    form.content_filter_enabled = record.content_filter_enabled
    form.filter_level = record.filter_level || 'moderate'
    form.blocked_categories = record.blocked_categories || []
    form.time_limit_enabled = record.time_limit_enabled
    form.daily_time_limit = record.daily_time_limit || 60
    form.allowed_time_range = record.allowed_time_range
    form.disabled_days = record.disabled_days || []
  } else {
    isEdit.value = false
    resetForm()
  }
  configModalVisible.value = true
}

function resetForm() {
  form.id = null
  form.child_id = null
  form.content_filter_enabled = false
  form.filter_level = 'moderate'
  form.blocked_categories = []
  form.time_limit_enabled = false
  form.daily_time_limit = 60
  form.allowed_time_range = null
  form.disabled_days = []
}

async function handleSave() {
  if (!isEdit.value && !form.child_id) {
    Message.warning('请选择儿童账号')
    return
  }
  try {
    const payload = { ...form }
    if (payload.allowed_time_range && payload.allowed_time_range.length === 2) {
      payload.allowed_start_time = payload.allowed_time_range[0]
      payload.allowed_end_time = payload.allowed_time_range[1]
    }
    delete payload.allowed_time_range

    const url = isEdit.value ? `/api/v1/family/child-mode/${form.id}` : '/api/v1/family/child-mode'
    const method = isEdit.value ? 'PUT' : 'POST'
    const res = await fetch(url, {
      method,
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify(payload)
    })
    const data = await res.json()
    if (data.code === 0 || data.code === 200) {
      Message.success(isEdit.value ? '配置已更新' : '配置已保存')
      configModalVisible.value = false
      loadChildModes()
    } else {
      Message.error(data.message || '保存失败')
    }
  } catch {
    Message.error('网络错误')
  }
}

async function toggleMode(record: any) {
  savingId.value = record.id
  try {
    const res = await fetch(`/api/v1/family/child-mode/${record.id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({ enabled: record.enabled })
    })
    const data = await res.json()
    if (data.code === 0 || data.code === 200) {
      Message.success(record.enabled ? '已启用' : '已禁用')
    } else {
      record.enabled = !record.enabled
      Message.error(data.message || '操作失败')
    }
  } catch {
    record.enabled = !record.enabled
    Message.error('网络错误')
  } finally {
    savingId.value = null
  }
}

async function handleDelete(record: any) {
  try {
    const res = await fetch(`/api/v1/family/child-mode/${record.id}`, {
      method: 'DELETE',
      credentials: 'include'
    })
    const data = await res.json()
    if (data.code === 0 || data.code === 200) {
      Message.success('已删除')
      loadChildModes()
    } else {
      Message.error(data.message || '删除失败')
    }
  } catch {
    Message.error('网络错误')
  }
}

onMounted(() => {
  loadChildren()
  loadChildModes()
})
</script>
