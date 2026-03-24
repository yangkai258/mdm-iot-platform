<template>
  <div class="page-container">
    <div class="search-form">
      <a-form :model="form" layout="inline">
        <a-form-item label="设置类型">
          <a-select v-model="form.type" placeholder="请选择" style="width: 160px">
            <a-option value="basic">基本信息</a-option>
            <a-option value="notification">通知设置</a-option>
            <a-option value="privacy">隐私设置</a-option>
          </a-select>
        </a-form-item>
        <a-form-item>
          <a-button type="primary" @click="handleSearch">搜索</a-button>
          <a-button @click="handleReset">重置</a-button>
        </a-form-item>
      </a-form>
    </div>
    <div class="toolbar">
      <a-button type="primary" @click="handleSave">保存设置</a-button>
    </div>
    <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" />
    <a-modal v-model:visible="modalVisible" title="编辑设置" :width="520">
      <a-form :model="form" label-col-flex="100px">
        <a-form-item label="设置名称">
          <a-input v-model="form.name" />
        </a-form-item>
        <a-form-item label="设置值">
          <a-input v-model="form.value" />
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
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const loading = ref(false)
const data = ref<any[]>([])
const modalVisible = ref(false)

const form = reactive({
  name: '',
  value: '',
  description: ''
})

const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0
})

const columns = [
  { title: '设置项', dataIndex: 'name', width: 200 },
  { title: '当前值', dataIndex: 'value', width: 200 },
  { title: '描述', dataIndex: 'description', ellipsis: true },
  { title: '操作', slotName: 'actions', width: 120 }
]

async function loadData() {
  loading.value = true
  try {
    const res = await fetch('/api/v1/family/settings')
    const json = await res.json()
    data.value = json.data ? [json.data] : []
    pagination.total = data.value.length
  } catch {
    Message.error('加载失败')
  } finally {
    loading.value = false
  }
}

function handleSearch() {
  loadData()
}

function handleReset() {
  form.name = ''
  form.value = ''
  form.description = ''
  loadData()
}

function handleSave() {
  Message.success('保存成功')
}

async function handleSubmit() {
  try {
    await fetch('/api/v1/family/settings/basic', {
      method: 'PUT',
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
