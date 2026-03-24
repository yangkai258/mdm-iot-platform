<template>
  <div class="page-container">
    <div class="search-form">
      <a-form :model="form" layout="inline">
        <a-form-item label="儿童账号">
          <a-select v-model="form.child_id" placeholder="请选择" allow-clear style="width: 160px">
            <a-option v-for="c in children" :key="c.id" :value="c.id">{{ c.name }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="状态">
          <a-select v-model="form.status" placeholder="请选择" allow-clear style="width: 140px">
            <a-option value="enabled">已启用</a-option>
            <a-option value="disabled">已禁用</a-option>
          </a-select>
        </a-form-item>
        <a-form-item>
          <a-button type="primary" @click="handleSearch">搜索</a-button>
          <a-button @click="handleReset">重置</a-button>
        </a-form-item>
      </a-form>
    </div>
    <div class="toolbar">
      <a-button type="primary" @click="handleCreate">新建</a-button>
    </div>
    <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id" />
    <a-modal v-model:visible="modalVisible" :title="modalTitle" :width="560">
      <a-form :model="form" label-col-flex="100px">
        <a-form-item label="儿童账号">
          <a-select v-model="form.child_id" placeholder="请选择">
            <a-option v-for="c in children" :key="c.id" :value="c.id">{{ c.name }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="内容过滤">
          <a-switch v-model="form.content_filter_enabled" />
        </a-form-item>
        <a-form-item label="时间限制">
          <a-switch v-model="form.time_limit_enabled" />
        </a-form-item>
        <a-form-item label="每日时长(分钟)">
          <a-input-number v-model="form.daily_time_limit" :min="15" :max="480" style="width: 100%" />
        </a-form-item>
      </a-form>
      <template #footer>
        <a-button @click="modalVisible = false">取消</a-button>
        <a-button type="primary" @click="handleSubmit">确定</a-button>
      </template>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const loading = ref(false)
const data = ref<any[]>([])
const children = ref<any[]>([])
const modalVisible = ref(false)
const editingId = ref<number | null>(null)

const form = reactive({
  child_id: null as number | null,
  content_filter_enabled: false,
  time_limit_enabled: false,
  daily_time_limit: 60
})

const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0
})

const modalTitle = computed(() => editingId.value ? '编辑配置' : '新建配置')

const columns = [
  { title: '儿童账号', dataIndex: 'child_name', width: 160 },
  { title: '模式开关', dataIndex: 'enabled', width: 120 },
  { title: '内容过滤', dataIndex: 'content_filter_enabled', width: 120 },
  { title: '时间限制', dataIndex: 'time_limit', width: 160 },
  { title: '操作', slotName: 'actions', width: 160 }
]

async function loadChildren() {
  try {
    const res = await fetch('/api/v1/family/members?role=child')
    const json = await res.json()
    children.value = json.data?.list || []
  } catch {}
}

async function loadData() {
  loading.value = true
  try {
    const params = new URLSearchParams()
    if (form.child_id) params.append('child_id', String(form.child_id))
    if (form.status) params.append('status', form.status)
    params.append('page', String(pagination.current))
    params.append('page_size', String(pagination.pageSize))

    const res = await fetch(`/api/v1/family/child-mode?${params}`)
    const json = await res.json()
    data.value = json.data?.list || []
    pagination.total = json.data?.total || 0
  } catch {
    Message.error('加载失败')
  } finally {
    loading.value = false
  }
}

function handleSearch() {
  pagination.current = 1
  loadData()
}

function handleReset() {
  form.child_id = null
  form.status = ''
  pagination.current = 1
  loadData()
}

function handleCreate() {
  editingId.value = null
  form.child_id = null
  form.content_filter_enabled = false
  form.time_limit_enabled = false
  form.daily_time_limit = 60
  modalVisible.value = true
}

async function handleSubmit() {
  try {
    const method = editingId.value ? 'PUT' : 'POST'
    const url = editingId.value ? `/api/v1/family/child-mode/${editingId.value}` : '/api/v1/family/child-mode'
    await fetch(url, {
      method,
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(form)
    })
    Message.success('保存成功')
    modalVisible.value = false
    loadData()
  } catch {
    Message.error('保存失败')
  }
}

function onPageChange(page: number) {
  pagination.current = page
  loadData()
}

onMounted(() => {
  loadChildren()
  loadData()
})
</script>

<style scoped>
.page-container {
  background: #fff;
  border-radius: 4px;
  padding: 20px;
}
.search-form {
  margin-bottom: 16px;
  padding: 16px;
  background: #f7f8fa;
  border-radius: 4px;
}
.toolbar {
  margin-bottom: 16px;
}
</style>
