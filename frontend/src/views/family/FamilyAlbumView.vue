<template>
  <div class="page-container">
    <div class="search-form">
      <a-form :model="form" layout="inline">
        <a-form-item label="照片名称">
          <a-input v-model="form.keyword" placeholder="请输入" style="width: 200px" />
        </a-form-item>
        <a-form-item label="上传者">
          <a-select v-model="form.uploader_id" placeholder="请选择" allow-clear style="width: 160px">
            <a-option v-for="m in members" :key="m.id" :value="m.id">{{ m.name }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item>
          <a-button type="primary" @click="handleSearch">搜索</a-button>
          <a-button @click="handleReset">重置</a-button>
        </a-form-item>
      </a-form>
    </div>
    <div class="toolbar">
      <a-button type="primary" @click="handleUpload">上传照片</a-button>
    </div>
    <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="onPageChange" />
    <a-modal v-model:visible="modalVisible" title="照片详情" :width="600">
      <a-form :model="form" label-col-flex="100px">
        <a-form-item label="照片名称">
          <a-input v-model="form.name" />
        </a-form-item>
        <a-form-item label="上传者">
          <a-input v-model="form.uploader_name" readonly />
        </a-form-item>
        <a-form-item label="上传时间">
          <a-input v-model="form.created_at" readonly />
        </a-form-item>
      </a-form>
      <template #footer>
        <a-button @click="modalVisible = false">关闭</a-button>
      </template>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const loading = ref(false)
const data = ref<any[]>([])
const members = ref<any[]>([])
const modalVisible = ref(false)

const form = reactive({
  name: '',
  uploader_name: '',
  created_at: ''
})

const pagination = reactive({
  current: 1,
  pageSize: 12,
  total: 0
})

const columns = [
  { title: '照片名称', dataIndex: 'name', width: 200, ellipsis: true },
  { title: '上传者', dataIndex: 'uploader_name', width: 120 },
  { title: '文件大小', dataIndex: 'size', width: 120 },
  { title: '上传时间', dataIndex: 'created_at', width: 180 },
  { title: '操作', slotName: 'actions', width: 120 }
]

async function loadMembers() {
  try {
    const res = await fetch('/api/v1/family/members')
    const json = await res.json()
    members.value = json.data?.list || []
  } catch {}
}

async function loadData() {
  loading.value = true
  try {
    const params = new URLSearchParams()
    if (form.keyword) params.append('keyword', form.keyword)
    if (form.uploader_id) params.append('uploader_id', String(form.uploader_id))
    params.append('page', String(pagination.current))
    params.append('page_size', String(pagination.pageSize))

    const res = await fetch(`/api/v1/family/album?${params}`)
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
  form.keyword = ''
  form.uploader_id = ''
  pagination.current = 1
  loadData()
}

function handleUpload() {
  Message.info('上传功能开发中')
}

function onPageChange(page: number) {
  pagination.current = page
  loadData()
}

onMounted(() => {
  loadMembers()
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
