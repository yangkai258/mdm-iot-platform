<template>
  <div class="page-container">
    <div class="search-form">
      <a-form :model="form" layout="inline">
        <a-form-item label="情绪类型">
          <a-select v-model="form.emotion" placeholder="请选择" allow-clear style="width: 140px">
            <a-option value="happy">开心</a-option>
            <a-option value="sad">难过</a-option>
            <a-option value="angry">生气</a-option>
            <a-option value="fear">害怕</a-option>
          </a-select>
        </a-form-item>
        <a-form-item>
          <a-button type="primary" @click="handleSearch">搜索</a-button>
          <a-button @click="handleReset">重置</a-button>
        </a-form-item>
      </a-form>
    </div>
    <div class="toolbar">
      <a-button type="primary" @click="handleCreate">新建配置</a-button>
    </div>
    <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="onPageChange" />
    <a-modal v-model:visible="modalVisible" :title="modalTitle" :width="520">
      <a-form :model="form" label-col-flex="100px">
        <a-form-item label="情绪类型">
          <a-select v-model="form.emotion">
            <a-option value="happy">开心</a-option>
            <a-option value="sad">难过</a-option>
            <a-option value="angry">生气</a-option>
            <a-option value="fear">害怕</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="响应动作">
          <a-select v-model="form.actions" multiple placeholder="请选择">
            <a-option value="dance">跳舞</a-option>
            <a-option value="sing">唱歌</a-option>
            <a-option value="wave">挥手</a-option>
            <a-option value="comfort">安慰语音</a-option>
            <a-option value="hug">拥抱</a-option>
            <a-option value="calm">播放平复音乐</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="描述">
          <a-textarea v-model="form.description" :rows="3" />
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
const modalVisible = ref(false)
const editingId = ref<number | null>(null)

const form = reactive({
  emotion: '',
  actions: [] as string[],
  description: ''
})

const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0
})

const modalTitle = computed(() => editingId.value ? '编辑配置' : '新建配置')

const columns = [
  { title: '情绪类型', dataIndex: 'emotion', width: 120 },
  { title: '响应动作', dataIndex: 'actions', ellipsis: true, width: 300 },
  { title: '描述', dataIndex: 'description', ellipsis: true },
  { title: '操作', slotName: 'actions', width: 120 }
]

async function loadData() {
  loading.value = true
  try {
    const res = await fetch('/api/v1/emotion/response-config')
    const json = await res.json()
    data.value = json.data || []
    pagination.total = data.value.length
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
  form.emotion = ''
  form.actions = []
  form.description = ''
  pagination.current = 1
  loadData()
}

function handleCreate() {
  editingId.value = null
  form.emotion = ''
  form.actions = []
  form.description = ''
  modalVisible.value = true
}

function handleEdit(record: any) {
  editingId.value = record.id
  form.emotion = record.emotion
  form.actions = record.actions || []
  form.description = record.description || ''
  modalVisible.value = true
}

async function handleSubmit() {
  try {
    const method = editingId.value ? 'PUT' : 'POST'
    const url = editingId.value ? `/api/v1/emotion/response-config/${editingId.value}` : '/api/v1/emotion/response-config'
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

onMounted(() => loadData())
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
