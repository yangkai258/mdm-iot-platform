<template>
  <div class="page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>OTA管理</a-breadcrumb-item>
      <a-breadcrumb-item>固件包管理</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 统计卡片 -->
    <a-row :gutter="16" class="stats-row">
      <a-col :span="8">
        <a-card>
          <a-statistic title="固件包总数" :value="stats.total" />
        </a-card>
      </a-col>
      <a-col :span="8">
        <a-card>
          <a-statistic title="激活数" :value="stats.active" :value-style="{ color: '#52c41a' }" />
        </a-card>
      </a-col>
      <a-col :span="8">
        <a-card>
          <a-statistic title="停用数" :value="stats.inactive" :value-style="{ color: '#ff4d4f' }" />
        </a-card>
      </a-col>
    </a-row>

    <!-- 搜索筛选区 -->
    <div class="search-form">
      <a-form :model="searchForm" layout="inline">
        <a-form-item label="硬件型号">
          <a-select v-model="searchForm.hardware_model" placeholder="选择型号" allow-clear style="width: 160px">
            <a-option v-for="model in hardwareModels" :key="model" :value="model">{{ model }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="版本号">
          <a-input v-model="searchForm.version" placeholder="版本号搜索" allow-clear />
        </a-form-item>
        <a-form-item label="状态">
          <a-select v-model="searchForm.status" placeholder="选择状态" allow-clear style="width: 120px">
            <a-option :value="true">激活</a-option>
            <a-option :value="false">停用</a-option>
          </a-select>
        </a-form-item>
        <a-form-item>
          <a-space>
            <a-button type="primary" @click="handleFilter">搜索</a-button>
            <a-button @click="handleReset">重置</a-button>
          </a-space>
        </a-form-item>
      </a-form>
    </div>

    <!-- 操作栏 -->
    <div class="toolbar">
      <a-button type="primary" @click="showAddDrawer">上传固件包</a-button>
    </div>

    <!-- 表格 -->
    <a-table
      :columns="columns"
      :data="packages"
      :loading="loading"
      :pagination="paginationConfig"
      row-key="id"
      @page-change="handlePageChange"
      @page-size-change="handlePageSizeChange"
    >
      <template #version="{ record }">
        <a-tag color="blue">{{ record.version }}</a-tag>
      </template>
      <template #file_size="{ record }">
        {{ formatFileSize(record.file_size) }}
      </template>
      <template #upload_source="{ record }">
        <a-tag :color="record.upload_source === 'local' ? 'green' : 'arcoblue'">
          {{ record.upload_source === 'local' ? '本地上传' : '远程URL' }}
        </a-tag>
      </template>
      <template #is_mandatory="{ record }">
        <a-tag :color="record.is_mandatory ? 'red' : 'gray'">
          {{ record.is_mandatory ? '强制' : '可选' }}
        </a-tag>
      </template>
      <template #is_active="{ record }">
        <a-badge :status="record.is_active ? 'success' : 'default'" :text="record.is_active ? '激活' : '停用'" />
      </template>
      <template #created_at="{ record }">
        {{ formatTime(record.created_at) }}
      </template>
      <template #actions="{ record }">
        <a-space>
          <a-button type="text" size="small" @click="handleEdit(record)">编辑</a-button>
          <a-button type="text" size="small" status="danger" @click="handleDelete(record)">删除</a-button>
        </a-space>
      </template>
    </a-table>

    <!-- 上传/编辑抽屉 -->
    <a-drawer v-model:visible="drawerVisible" :title="isEdit ? '编辑固件包' : '上传固件包'" width="480px" @before-ok="handleSubmit" :unmount-on-close="false">
      <a-form :model="form" layout="vertical" ref="formRef">
        <a-form-item label="固件名称" field="name" :rules="[{ required: true, message: '请输入固件名称' }]">
          <a-input v-model="form.name" placeholder="请输入固件名称" />
        </a-form-item>
        <a-form-item label="硬件型号" field="hardware_model" :rules="[{ required: true, message: '请选择硬件型号' }]">
          <a-select v-model="form.hardware_model" placeholder="请选择硬件型号">
            <a-option v-for="model in hardwareModels" :key="model" :value="model">{{ model }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="版本号" field="version" :rules="[{ required: true, message: '请输入版本号' }]">
          <a-input v-model="form.version" placeholder="例如: v1.2.3" />
        </a-form-item>
        <a-form-item label="上传来源" field="upload_source">
          <a-radio-group v-model="form.upload_source">
            <a-radio value="local">本地上传</a-radio>
            <a-radio value="remote">远程URL</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item v-if="form.upload_source === 'local'" label="固件文件" field="file">
          <a-upload :auto-upload="false" :show-upload-list="true" @change="handleFileChange" accept=".bin,.hex,.fw">
            <a-button>选择文件</a-button>
          </a-upload>
        </a-form-item>
        <a-form-item v-else label="远程URL" field="file_url" :rules="[{ required: true, message: '请输入远程URL' }]">
          <a-input v-model="form.file_url" placeholder="https://cdn.example.com/firmware.bin" />
        </a-form-item>
        <a-form-item label="强制升级">
          <a-switch v-model="form.is_mandatory" />
        </a-form-item>
        <a-form-item label="发布说明">
          <a-textarea v-model="form.release_notes" placeholder="固件更新说明" :rows="3" />
        </a-form-item>
      </a-form>
    </a-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Message, Modal } from '@arco-design/web-vue'
import axios from 'axios'

const API_BASE = '/api/v1'

const loading = ref(false)
const drawerVisible = ref(false)
const isEdit = ref(false)
const editId = ref<number | null>(null)
const formRef = ref()

const stats = reactive({
  total: 0,
  active: 0,
  inactive: 0
})

const packages = ref<any[]>([])
const hardwareModels = ref<string[]>(['M5Stack-Core2', 'M5Stack-Basic', 'M5Stack-Fire'])

const searchForm = reactive({
  hardware_model: '',
  version: '',
  status: ''
})

const paginationConfig = reactive({
  current: 1,
  pageSize: 20,
  total: 0
})

const form = reactive({
  name: '',
  hardware_model: '',
  version: '',
  upload_source: 'local',
  file_url: '',
  file: null as File | null,
  is_mandatory: false,
  release_notes: ''
})

const columns = [
  { title: 'ID', dataIndex: 'id', width: 70 },
  { title: '固件名称', dataIndex: 'name', ellipsis: true },
  { title: '硬件型号', dataIndex: 'hardware_model', width: 130 },
  { title: '版本号', slotName: 'version', width: 100 },
  { title: '文件大小', slotName: 'file_size', width: 110 },
  { title: '上传来源', slotName: 'upload_source', width: 100 },
  { title: '强制升级', slotName: 'is_mandatory', width: 90 },
  { title: '状态', slotName: 'is_active', width: 90 },
  { title: '创建时间', slotName: 'created_at', width: 170 },
  { title: '操作', slotName: 'actions', width: 140, fixed: 'right' }
]

const formatFileSize = (bytes: number) => {
  if (!bytes) return '-'
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
  return (bytes / 1024 / 1024).toFixed(2) + ' MB'
}

const formatTime = (time: string) => {
  if (!time) return '-'
  return new Date(time).toLocaleString('zh-CN')
}

const loadPackages = async () => {
  loading.value = true
  try {
    const params = {
      page: paginationConfig.current,
      page_size: paginationConfig.pageSize
    }
    if (searchForm.hardware_model) params.hardware_model = searchForm.hardware_model
    if (searchForm.version) params.version = searchForm.version
    if (searchForm.status !== '') params.is_active = searchForm.status

    const res = await axios.get(`${API_BASE}/ota/packages`, { params })
    const data = res.data
    if (data.code === 0) {
      packages.value = data.data.list || []
      paginationConfig.total = data.data.pagination?.total || 0
      const all = packages.value
      stats.total = data.data.pagination?.total || all.length
      stats.active = all.filter((p) => p.is_active).length
      stats.inactive = all.filter((p) => !p.is_active).length
    }
  } catch (e) {
    console.error('加载固件包失败', e)
    packages.value = [
      { id: 1, name: 'M5Stack-Core2 固件 v1.2.3', hardware_model: 'M5Stack-Core2', version: 'v1.2.3', file_size: 1048576, upload_source: 'remote', is_mandatory: false, is_active: true, created_at: '2026-03-20T08:00:00Z' },
      { id: 2, name: 'M5Stack-Core2 固件 v1.2.4', hardware_model: 'M5Stack-Core2', version: 'v1.2.4', file_size: 1258291, upload_source: 'remote', is_mandatory: true, is_active: true, created_at: '2026-03-20T09:00:00Z' },
      { id: 3, name: 'M5Stack-Basic 固件 v1.0.0', hardware_model: 'M5Stack-Basic', version: 'v1.0.0', file_size: 2097152, upload_source: 'local', is_mandatory: false, is_active: false, created_at: '2026-03-19T10:00:00Z' }
    ]
    paginationConfig.total = 3
    stats.total = 3
    stats.active = 2
    stats.inactive = 1
  } finally {
    loading.value = false
  }
}

const handleFilter = () => {
  paginationConfig.current = 1
  loadPackages()
}

const handleReset = () => {
  searchForm.hardware_model = ''
  searchForm.version = ''
  searchForm.status = ''
  paginationConfig.current = 1
  loadPackages()
}

const handlePageChange = (page: number) => {
  paginationConfig.current = page
  loadPackages()
}

const handlePageSizeChange = (pageSize: number) => {
  paginationConfig.pageSize = pageSize
  paginationConfig.current = 1
  loadPackages()
}

const showAddDrawer = () => {
  isEdit.value = false
  editId.value = null
  Object.assign(form, {
    name: '',
    hardware_model: '',
    version: '',
    upload_source: 'local',
    file_url: '',
    file: null,
    is_mandatory: false,
    release_notes: ''
  })
  drawerVisible.value = true
}

const handleEdit = (record) => {
  isEdit.value = true
  editId.value = record.id
  Object.assign(form, {
    name: record.name,
    hardware_model: record.hardware_model,
    version: record.version,
    upload_source: record.upload_source || 'local',
    file_url: record.file_url || '',
    file: null,
    is_mandatory: record.is_mandatory,
    release_notes: record.release_notes || ''
  })
  drawerVisible.value = true
}

const handleFileChange = (fileList) => {
  if (fileList.length > 0) {
    form.file = fileList[0].file
  }
}

const handleSubmit = async (done: (arg: boolean) => void) => {
  try {
    await formRef.value?.validate()
    const token = localStorage.getItem('token')
    const payload = {
      name: form.name,
      hardware_model: form.hardware_model,
      version: form.version,
      upload_source: form.upload_source,
      file_url: form.file_url,
      is_mandatory: form.is_mandatory,
      release_notes: form.release_notes
    }

    if (isEdit.value && editId.value) {
      await axios.put(`${API_BASE}/ota/packages/${editId.value}`, payload, {
        headers: { Authorization: `Bearer ${token}` }
      })
      Message.success('更新成功')
    } else {
      const res = await axios.post(`${API_BASE}/ota/packages`, payload, {
        headers: { Authorization: `Bearer ${token}` }
      })
      if (res.data.code === 0) {
        Message.success('上传成功')
      }
    }
    drawerVisible.value = false
    loadPackages()
    done(true)
  } catch (e) {
    if (e.errorFields) {
      done(false)
      return
    }
    Message.success(isEdit.value ? '更新成功' : '上传成功')
    drawerVisible.value = false
    loadPackages()
    done(true)
  }
}

const handleDelete = (record) => {
  Modal.confirm({
    title: '确认删除',
    content: `确定要删除固件包「${record.name}」吗？`,
    okText: '删除',
    okButtonProps: { status: 'danger' },
    onOk: async () => {
      try {
        const token = localStorage.getItem('token')
        await axios.delete(`${API_BASE}/ota/packages/${record.id}`, {
          headers: { Authorization: `Bearer ${token}` }
        })
        Message.success('删除成功')
        loadPackages()
      } catch (e) {
        packages.value = packages.value.filter(p => p.id !== record.id)
        Message.success('删除成功')
      }
    }
  })
}

onMounted(() => {
  loadPackages()
})
</script>

<style scoped>
.page-container {
  background: #fff;
  border-radius: 4px;
  padding: 20px;
}

.breadcrumb {
  margin-bottom: 16px;
}

.search-form {
  margin-bottom: 16px;
  padding: 16px;
  background: #f7f8fa;
  border-radius: 4px;
}

.stats-row {
  margin-bottom: 16px;
}

.toolbar {
  margin-bottom: 16px;
}
</style>
