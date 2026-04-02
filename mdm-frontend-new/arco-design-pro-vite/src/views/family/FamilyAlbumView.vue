<template>
    <Breadcrumb :items="['Home','Console','']" />


  <div class="container">
    <a-card class="general-card" title="家庭相册">
      <template #extra>
        <a-space :size="12">
          <a-button type="primary" @click="handleUpload"><icon-plus />上传照片</a-button>
          <a-button @click="handleSearch"><icon-refresh />刷新</a-button>
        </a-space>
      </template>
      <a-row :gutter="16">
        <a-col :span="8">
          <a-form-item label="照片名称">
            <a-input v-model="form.keyword" placeholder="请输入" @pressEnter="handleSearch" />
          </a-form-item>
        </a-col>
        <a-col :span="8">
          <a-form-item label="上传者">
            <a-select v-model="form.uploader_id" placeholder="请选择" allow-clear style="width: 100%">
              <a-option v-for="m in members" :key="m.id" :value="m.id">{{ m.name }}</a-option>
            </a-select>
          </a-form-item>
        </a-col>
        <a-col :flex="'86px'" style="display: flex; align-items: flex-end">
          <a-space direction="vertical" :size="8">
            <a-button type="primary" @click="handleSearch">查询</a-button>
            <a-button @click="handleReset">重置</a-button>
          </a-space>
        </a-col>
      </a-row>
      <a-divider style="margin: 0 0 16px 0" />
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id" />
    </a-card>
      </a-table>
    <a-modal v-model:visible="modalVisible" title="照片详情" :width="600">
      <a-form :model="form" label-col-flex="100px">
        <a-form-item label="照片名称"><a-input v-model="form.name" readonly /></a-form-item>
        <a-form-item label="上传者"><a-input v-model="form.uploader_name" readonly /></a-form-item>
        <a-form-item label="上传时间"><a-input v-model="form.created_at" readonly /></a-form-item>
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
import Breadcrumb from '@/components/Breadcrumb.vue'

const loading = ref(false)
const data = ref<any[]>([])
const members = ref<any[]>([])
const modalVisible = ref(false)
const form = reactive({ name: '', uploader_name: '', created_at: '' })
const pagination = reactive({ current: 1, pageSize: 12, total: 0 })
const columns = [
  { title: '照片名称', dataIndex: 'name', width: 200, ellipsis: true },
  { title: '上传者', dataIndex: 'uploader_name', width: 120 },
  { title: '文件大小', dataIndex: 'size', width: 120 },
  { title: '上传时间', dataIndex: 'created_at', width: 180 },
  { title: '操作', slotName: 'actions', width: 120 }
]

async function loadMembers() {
  try {
    const res = await fetch('/api/v1/family/members', {
      headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token') }
    }).then(r => r.json())
    members.value = res.data?.list || []
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
    const res = await fetch(`/api/v1/family/album?${params}`, {
      headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token') }
    }).then(r => r.json())
    data.value = res.data?.list || []
    pagination.total = res.data?.total || 0
  } catch { Message.error('加载失败') } finally { loading.value = false }
}

const handleSearch = () => { pagination.current = 1; loadData() }
const handleReset = () => { form.keyword = ''; form.uploader_id = ''; pagination.current = 1; loadData() }
const handleUpload = () => { Message.info('上传功能开发中') }
const onPageChange = (page: number) => { pagination.current = page; loadData() }

onMounted(() => { loadMembers(); loadData() })
</script>
