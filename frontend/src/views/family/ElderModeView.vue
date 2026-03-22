<template>
  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>家庭管理</a-breadcrumb-item>
      <a-breadcrumb-item>老人陪伴模式</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 页面标题 -->
    <div class="pro-page-header">
      <h2 class="pro-page-title">老人陪伴模式</h2>
      <p class="pro-page-desc">为老年用户配置简化界面、大字体和高对比度显示</p>
    </div>

    <!-- 搜索筛选区 -->
    <div class="pro-search-bar">
      <a-space>
        <a-select v-model="elderFilter" placeholder="选择老人账号" allow-clear style="width: 200px" @change="loadElderModes">
          <a-option v-for="e in elders" :key="e.id" :value="e.id">{{ e.name }}</a-option>
        </a-select>
        <a-select v-model="statusFilter" placeholder="模式状态" allow-clear style="width: 140px" @change="loadElderModes">
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
        <a-button @click="loadElderModes">
          <template #icon><icon-refresh /></template>
          刷新
        </a-button>
      </a-space>
    </div>

    <!-- 数据内容区 -->
    <div class="pro-content-area">
      <a-table
        :columns="columns"
        :data="elderModes"
        :loading="loading"
        :pagination="pagination"
        @page-change="onPageChange"
        row-key="id"
      >
        <template #elder_name="{ record }">
          <a-avatar :style="{ backgroundColor: '#8e4e9c' }" :size="32">
            {{ record.elder_name?.charAt(0) || '?' }}
          </a-avatar>
          <span style="margin-left: 8px">{{ record.elder_name }}</span>
        </template>
        <template #enabled="{ record }">
          <a-switch v-model="record.enabled" @change="toggleMode(record)" :disabled="savingId === record.id" />
        </template>
        <template #simplified_ui="{ record }">
          <a-tag :color="record.simplified_ui ? 'green' : 'gray'">
            {{ record.simplified_ui ? '已启用' : '已禁用' }}
          </a-tag>
        </template>
        <template #font_size="{ record }">
          {{ getFontLabel(record.font_size) }}
        </template>
        <template #high_contrast="{ record }">
          <a-tag :color="record.high_contrast ? 'green' : 'gray'">
            {{ record.high_contrast ? '已启用' : '已禁用' }}
          </a-tag>
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
    <a-modal v-model:visible="configModalVisible" :title="isEdit ? '编辑老人模式配置' : '新增老人模式配置'" @ok="handleSave" :width="560" @close="resetForm">
      <a-form :model="form" layout="vertical">
        <a-form-item label="老人账号" required>
          <a-select v-model="form.elder_id" placeholder="请选择老人账号" :disabled="isEdit">
            <a-option v-for="e in elders" :key="e.id" :value="e.id">{{ e.name }} ({{ e.phone }})</a-option>
          </a-select>
        </a-form-item>

        <a-divider>界面显示</a-divider>

        <a-form-item label="启用简化界面">
          <a-switch v-model="form.simplified_ui" />
        </a-form-item>
        <a-form-item label="隐藏复杂功能" v-if="form.simplified_ui">
          <a-checkbox-group v-model="form.hidden_features">
            <a-checkbox value="ota">OTA升级</a-checkbox>
            <a-checkbox value="analytics">数据分析</a-checkbox>
            <a-checkbox value="notifications">通知管理</a-checkbox>
            <a-checkbox value="advanced">高级设置</a-checkbox>
          </a-checkbox-group>
        </a-form-item>

        <a-divider>字体与视觉</a-divider>

        <a-form-item label="字体大小">
          <a-radio-group v-model="form.font_size">
            <a-radio value="large">大（1.2倍）</a-radio>
            <a-radio value="xlarge">超大（1.5倍）</a-radio>
            <a-radio value="xxlarge">极大（2倍）</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item label="高对比度模式">
          <a-switch v-model="form.high_contrast" />
        </a-form-item>
        <a-form-item label="深色模式">
          <a-switch v-model="form.dark_mode" />
        </a-form-item>

        <a-divider>语音与交互</a-divider>

        <a-form-item label="语音播报">
          <a-switch v-model="form.voice_announce" />
        </a-form-item>
        <a-form-item label="触摸反馈音效" v-if="form.voice_announce">
          <a-switch v-model="form.touch_sound" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const elderModes = ref<any[]>([])
const elders = ref<any[]>([])
const loading = ref(false)
const savingId = ref<number | null>(null)
const elderFilter = ref('')
const statusFilter = ref('')
const configModalVisible = ref(false)
const isEdit = ref(false)

const form = reactive({
  id: null as number | null,
  elder_id: null as number | null,
  simplified_ui: true,
  hidden_features: [] as string[],
  font_size: 'large',
  high_contrast: false,
  dark_mode: false,
  voice_announce: false,
  touch_sound: true
})

const columns = [
  { title: '老人账号', dataIndex: 'elder_name', slotName: 'elder_name', width: 200 },
  { title: '模式开关', dataIndex: 'enabled', slotName: 'enabled', width: 120 },
  { title: '简化界面', dataIndex: 'simplified_ui', slotName: 'simplified_ui', width: 120 },
  { title: '字体大小', dataIndex: 'font_size', slotName: 'font_size', width: 140 },
  { title: '高对比度', dataIndex: 'high_contrast', slotName: 'high_contrast', width: 120 },
  { title: '操作', slotName: 'actions', width: 160 }
]

const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0
})

const fontLabels: Record<string, string> = {
  large: '大（1.2倍）',
  xlarge: '超大（1.5倍）',
  xxlarge: '极大（2倍）'
}

function getFontLabel(size: string) {
  return fontLabels[size] || size || '-'
}

async function loadElders() {
  try {
    const res = await fetch('/api/v1/family/members?role=elder', { credentials: 'include' })
    const data = await res.json()
    elders.value = data.data?.list || data.data?.members || []
  } catch { /* ignore */ }
}

async function loadElderModes() {
  loading.value = true
  try {
    const params = new URLSearchParams()
    if (elderFilter.value) params.append('elder_id', elderFilter.value)
    if (statusFilter.value) params.append('status', statusFilter.value)
    params.append('page', String(pagination.current))
    params.append('page_size', String(pagination.pageSize))

    const res = await fetch(`/api/v1/family/elder-mode?${params}`, { credentials: 'include' })
    const data = await res.json()
    if (data.code === 0 || data.code === 200) {
      elderModes.value = data.data?.list || data.data || []
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
  loadElderModes()
}

function showConfigModal(record: any) {
  if (record) {
    isEdit.value = true
    form.id = record.id
    form.elder_id = record.elder_id
    form.simplified_ui = record.simplified_ui ?? true
    form.hidden_features = record.hidden_features || []
    form.font_size = record.font_size || 'large'
    form.high_contrast = record.high_contrast ?? false
    form.dark_mode = record.dark_mode ?? false
    form.voice_announce = record.voice_announce ?? false
    form.touch_sound = record.touch_sound ?? true
  } else {
    isEdit.value = false
    resetForm()
  }
  configModalVisible.value = true
}

function resetForm() {
  form.id = null
  form.elder_id = null
  form.simplified_ui = true
  form.hidden_features = []
  form.font_size = 'large'
  form.high_contrast = false
  form.dark_mode = false
  form.voice_announce = false
  form.touch_sound = true
}

async function handleSave() {
  if (!isEdit.value && !form.elder_id) {
    Message.warning('请选择老人账号')
    return
  }
  try {
    const url = isEdit.value ? `/api/v1/family/elder-mode/${form.id}` : '/api/v1/family/elder-mode'
    const method = isEdit.value ? 'PUT' : 'POST'
    const res = await fetch(url, {
      method,
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({ ...form })
    })
    const data = await res.json()
    if (data.code === 0 || data.code === 200) {
      Message.success(isEdit.value ? '配置已更新' : '配置已保存')
      configModalVisible.value = false
      loadElderModes()
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
    const res = await fetch(`/api/v1/family/elder-mode/${record.id}`, {
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
    const res = await fetch(`/api/v1/family/elder-mode/${record.id}`, {
      method: 'DELETE',
      credentials: 'include'
    })
    const data = await res.json()
    if (data.code === 0 || data.code === 200) {
      Message.success('已删除')
      loadElderModes()
    } else {
      Message.error(data.message || '删除失败')
    }
  } catch {
    Message.error('网络错误')
  }
}

onMounted(() => {
  loadElders()
  loadElderModes()
})
</script>
