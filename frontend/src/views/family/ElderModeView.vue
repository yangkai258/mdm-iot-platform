<template>
  <div class="page-container">
    <div class="search-form">
      <a-form :model="form" layout="inline">
        <a-form-item label="老人账号">
          <a-select v-model="form.elder_id" placeholder="请选择" allow-clear style="width: 160px">
            <a-option v-for="e in elders" :key="e.id" :value="e.id">{{ e.name }}</a-option>
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
        <a-form-item label="老人账号">
          <a-select v-model="form.elder_id" placeholder="请选择">
            <a-option v-for="e in elders" :key="e.id" :value="e.id">{{ e.name }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="简化界面">
          <a-switch v-model="form.simplified_ui" />
        </a-form-item>
        <a-form-item label="字体大小">
          <a-select v-model="form.font_size" style="width: 100%">
            <a-option value="large">大（1.2倍）</a-option>
            <a-option value="xlarge">超大（1.5倍）</a-option>
            <a-option value="xxlarge">极大（2倍）</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="高对比度">
          <a-switch v-model="form.high_contrast" />
        </a-form-item>
        <a-form-item label="语音播报">
          <a-switch v-model="form.voice_announce" />
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
const elders = ref<any[]>([])
const modalVisible = ref(false)
const editingId = ref<number | null>(null)

const form = reactive({
  elder_id: null as number | null,
  simplified_ui: true,
  font_size: 'large',
  high_contrast: false,
  voice_announce: false
})

const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0
})

const modalTitle = computed(() => editingId.value ? '编辑配置' : '新建配置')

const columns = [
  { title: '老人账号', dataIndex: 'elder_name', width: 160 },
  { title: '模式开关', dataIndex: 'enabled', width: 120 },
  { title: '简化界面', dataIndex: 'simplified_ui', width: 120 },
  { title: '字体大小', dataIndex: 'font_size', width: 140 },
  { title: '高对比度', dataIndex: 'high_contrast', width: 120 },
  { title: '操作', slotName: 'actions', width: 160 }
]

async function loadElders() {
  try {
    const res = await fetch('/api/v1/family/members?role=elder')
    const json = await res.json()
    elders.value = json.data?.list || []
  } catch {}
}

async function loadData() {
  loading.value = true
  try {
    const params = new URLSearchParams()
    if (form.elder_id) params.append('elder_id', String(form.elder_id))
    if (form.status) params.append('status', form.status)
    params.append('page', String(pagination.current))
    params.append('page_size', String(pagination.pageSize))

    const res = await fetch(`/api/v1/family/elder-mode?${params}`)
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
  form.elder_id = null
  form.status = ''
  pagination.current = 1
  loadData()
}

function handleCreate() {
  editingId.value = null
  form.elder_id = null
  form.simplified_ui = true
  form.font_size = 'large'
  form.high_contrast = false
  form.voice_announce = false
  modalVisible.value = true
}

async function handleSubmit() {
  try {
    const method = editingId.value ? 'PUT' : 'POST'
    const url = editingId.value ? `/api/v1/family/elder-mode/${editingId.value}` : '/api/v1/family/elder-mode'
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
  loadElders()
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
