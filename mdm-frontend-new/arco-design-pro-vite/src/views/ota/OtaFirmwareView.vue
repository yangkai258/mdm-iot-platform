<template>
  <div class="ota-firmware-view-container">
    <Breadcrumb :items="[{title: '首页', path: '/'},{title: 'OTA升级'},{title: '固件管理'}]" />

    <!-- 统计卡片 -->
    <a-row :gutter="16" style="margin-bottom: 16px">
      <a-col :span="6">
        <a-card class="general-card">
          <a-statistic title="固件版本" :value="stats.versions" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="general-card">
          <a-statistic title="设备总数" :value="stats.totalDevices" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="general-card">
          <a-statistic title="已升级" :value="stats.upgraded" :value-style="{ color: '#52c41a' }" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="general-card">
          <a-statistic title="待升级" :value="stats.pending" :value-style="{ color: '#faad14' }" />
        </a-card>
      </a-col>
    </a-row>

    <!-- 搜索筛选区 -->
    <a-card class="general-card" style="margin-bottom: 16px">
      <a-row :gutter="16">
        <a-col :flex="1">
          <a-form :model="searchForm" layout="vertical">
            <a-row :gutter="16">
              <a-col :span="8">
                <a-form-item label="关键词">
                  <a-input v-model="searchForm.keyword" placeholder="搜索固件版本/硬件型号" allow-clear />
                </a-form-item>
              </a-col>
              <a-col :span="8">
                <a-form-item label="硬件型号">
                  <a-select v-model="searchForm.hardware_model" placeholder="全部" allow-clear>
                    <a-option value="MDM-Pro-200">MDM-Pro-200</a-option>
                    <a-option value="MDM-Mini-100">MDM-Mini-100</a-option>
                    <a-option value="MDM-Lite-50">MDM-Lite-50</a-option>
                  </a-select>
                </a-form-item>
              </a-col>
              <a-col :span="8">
                <a-form-item label="状态">
                  <a-select v-model="searchForm.status" placeholder="全部" allow-clear>
                    <a-option :value="1">已发布</a-option>
                    <a-option :value="2">测试中</a-option>
                    <a-option :value="3">已废弃</a-option>
                  </a-select>
                </a-form-item>
              </a-col>
            </a-row>
          </a-form>
        </a-col>
        <a-divider style="height: 84px" direction="vertical" />
        <a-col :flex="'86px'" style="text-align: right">
          <a-space direction="vertical" :size="18">
            <a-button type="primary" @click="handleSearch">
              <template #icon><icon-search /></template>
              查询
            </a-button>
            <a-button @click="handleReset">
              <template #icon><icon-refresh /></template>
              重置
            </a-button>
          </a-space>
        </a-col>
      </a-row>
    </a-card>

    <!-- 数据表格 -->
    <a-card class="general-card" style="margin-bottom: 16px">
      <template #title>固件列表</template>
      <template #extra>
        <a-space>
          <a-button type="primary" @click="showUploadModal">
            <template #icon><icon-plus /></template>
            上传固件
          </a-button>
          <a-button @click="loadFirmwares">
            <template #icon><icon-refresh /></template>
            刷新
          </a-button>
        </a-space>
      </template>

      <a-table
        :columns="columns"
        :data="filteredFirmwares"
        :loading="loading"
        :pagination="paginationConfig"
        @page-change="onPageChange"
        @page-size-change="onPageSizeChange"
        row-key="id"
      >
        <template #version="{ record }">
          <a-tag color="blue">v{{ record.version }}</a-tag>
        </template>
        <template #status="{ record }">
          <a-badge :status="getStatusBadge(record.status)" :text="getStatusText(record.status)" />
        </template>
        <template #fileSize="{ record }">
          {{ formatFileSize(record.file_size) }}
        </template>
        <template #uploadTime="{ record }">
          {{ formatTime(record.upload_time) }}
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="handleDeploy(record)">下发升级</a-button>
            <a-button type="text" size="small" status="danger" @click="handleDelete(record)">删除</a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 升级任务 -->
    <a-card class="general-card">
      <template #title>升级任务</template>
      <a-table
        :columns="taskColumns"
        :data="tasks"
        :loading="taskLoading"
        :pagination="false"
        row-key="id"
      >
        <template #status="{ record }">
          <a-tag :color="getTaskStatusColor(record.status)">
            {{ getTaskStatusText(record.status) }}
          </a-tag>
        </template>
        <template #progress="{ record }">
          <a-progress :percent="record.progress" :stroke-width="6" :show-text="true" />
        </template>
        <template #createTime="{ record }">
          {{ formatTime(record.create_time) }}
        </template>
      </a-table>
    </a-card>

    <!-- 上传固件弹窗 -->
    <a-modal v-model:visible="uploadModalVisible" title="上传固件" @ok="handleUpload" :confirm-loading="uploading" :width="520">
      <a-form :model="uploadForm" layout="vertical">
        <a-form-item label="固件版本" required>
          <a-input v-model="uploadForm.version" placeholder="例如: 1.2.1" />
        </a-form-item>
        <a-form-item label="硬件型号" required>
          <a-select v-model="uploadForm.hardware_model" placeholder="选择硬件型号">
            <a-option value="MDM-Pro-200">MDM-Pro-200</a-option>
            <a-option value="MDM-Mini-100">MDM-Mini-100</a-option>
            <a-option value="MDM-Lite-50">MDM-Lite-50</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="固件描述">
          <a-textarea v-model="uploadForm.description" placeholder="固件更新说明" :rows="3" />
        </a-form-item>
        <a-form-item label="固件文件" required>
          <a-upload :auto-upload="false" :show-upload-list="true" @change="handleFileChange">
            <a-button>选择文件</a-button>
          </a-upload>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const API_BASE = '/api/v1'
const loading = ref(false)
const taskLoading = ref(false)
const uploading = ref(false)
const uploadModalVisible = ref(false)

const stats = reactive({ versions: 3, totalDevices: 156, upgraded: 142, pending: 14 })

const searchForm = reactive({
  keyword: '',
  hardware_model: undefined as string | undefined,
  status: undefined as number | undefined
})

const uploadForm = reactive({
  version: '', hardware_model: '', description: '', file: null as any
})

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '版本', slotName: 'version', width: 100 },
  { title: '硬件型号', dataIndex: 'hardware_model', width: 120 },
  { title: '文件大小', slotName: 'fileSize', width: 100 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '上传时间', slotName: 'uploadTime', width: 180 },
  { title: '操作', slotName: 'actions', width: 180 }
]

const taskColumns = [
  { title: '任务ID', dataIndex: 'id', width: 80 },
  { title: '固件版本', dataIndex: 'firmware_version', width: 100 },
  { title: '目标设备', dataIndex: 'target_devices', width: 100 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '进度', slotName: 'progress', width: 200 },
  { title: '创建时间', slotName: 'createTime', width: 180 }
]

const firmwares = ref<any[]>([
  { id: 1, version: '1.2.0', hardware_model: 'MDM-Pro-200', file_size: 5242880, status: 1, upload_time: '2026-03-19 10:00:00', description: '修复已知问题' },
  { id: 2, version: '1.1.5', hardware_model: 'MDM-Pro-200', file_size: 5242880, status: 2, upload_time: '2026-03-15 14:30:00', description: '优化性能' },
  { id: 3, version: '1.0.0', hardware_model: 'MDM-Pro-200', file_size: 5242880, status: 3, upload_time: '2026-03-01 09:00:00', description: '初始版本' }
])

const tasks = ref<any[]>([
  { id: 1, firmware_version: 'v1.2.0', target_devices: 50, status: 1, progress: 100, create_time: '2026-03-19 08:00:00' },
  { id: 2, firmware_version: 'v1.2.0', target_devices: 100, status: 2, progress: 45, create_time: '2026-03-19 10:00:00' }
])

const pagination = reactive({ current: 1, pageSize: 10, total: 3 })

const paginationConfig = computed(() => ({
  current: pagination.current,
  pageSize: pagination.pageSize,
  total: pagination.total,
  showTotal: true,
  showSizeChanger: true
}))

const filteredFirmwares = computed(() => {
  let result = firmwares.value
  if (searchForm.keyword) {
    result = result.filter(f =>
      f.version.includes(searchForm.keyword) || f.hardware_model.includes(searchForm.keyword)
    )
  }
  if (searchForm.hardware_model) {
    result = result.filter(f => f.hardware_model === searchForm.hardware_model)
  }
  if (searchForm.status !== undefined) {
    result = result.filter(f => f.status === searchForm.status)
  }
  return result
})

const loadFirmwares = async () => {
  loading.value = true
  try {
    const res = await fetch(`${API_BASE}/firmwares`)
    const data = await res.json()
    if (data.code === 0) firmwares.value = data.data?.list || firmwares.value
  } catch {} finally { loading.value = false }
}

const loadTasks = async () => {
  taskLoading.value = true
  try {
    const res = await fetch(`${API_BASE}/ota/tasks`)
    const data = await res.json()
    if (data.code === 0) tasks.value = data.data?.list || tasks.value
  } catch {} finally { taskLoading.value = false }
}

const handleSearch = () => { pagination.current = 1; pagination.total = filteredFirmwares.value.length }
const handleReset = () => { searchForm.keyword = ''; searchForm.hardware_model = undefined; searchForm.status = undefined; pagination.total = firmwares.value.length }
const onPageChange = (page: number) => { pagination.current = page }
const onPageSizeChange = (size: number) => { pagination.pageSize = size; pagination.current = 1 }

const showUploadModal = () => { uploadModalVisible.value = true }
const handleFileChange = (file: any) => { uploadForm.file = file }

const handleUpload = () => {
  if (!uploadForm.version || !uploadForm.hardware_model) {
    Message.warning('请填写必填项'); return
  }
  uploading.value = true
  setTimeout(() => {
    firmwares.value.unshift({
      id: Date.now(), version: uploadForm.version, hardware_model: uploadForm.hardware_model,
      file_size: 5242880, status: 1, upload_time: new Date().toLocaleString(), description: uploadForm.description
    })
    uploadModalVisible.value = false
    uploading.value = false
    Message.success('固件上传成功')
    Object.assign(uploadForm, { version: '', hardware_model: '', description: '', file: null })
  }, 1500)
}

const handleDeploy = (record: any) => {
  Message.info(`开始下发固件 ${record.version} 到设备...`)
  tasks.value.unshift({
    id: Date.now(), firmware_version: `v${record.version}`, target_devices: stats.pending,
    status: 2, progress: 0, create_time: new Date().toLocaleString()
  })
  Message.success('升级任务已创建')
}

const handleDelete = (record: any) => {
  firmwares.value = firmwares.value.filter(f => f.id !== record.id)
  Message.success('固件已删除')
}

const getStatusBadge = (status: number) => ({ 1: 'success', 2: 'normal', 3: 'default' }[status] || 'default')
const getStatusText = (status: number) => ({ 1: '已发布', 2: '测试中', 3: '已废弃' }[status] || '未知')
const getTaskStatusColor = (status: number) => ({ 1: 'green', 2: 'blue', 3: 'red', 4: 'gray' }[status] || 'default')
const getTaskStatusText = (status: number) => ({ 1: '已完成', 2: '进行中', 3: '失败', 4: '已取消' }[status] || '未知')
const formatFileSize = (bytes: number) => bytes ? (bytes / 1024 / 1024).toFixed(2) + ' MB' : '-'
const formatTime = (time: string) => time || '-'

onMounted(() => { loadFirmwares(); loadTasks() })
</script>

<style scoped>
.ota-firmware-view-container {
  padding: 20px 24px;
  min-height: calc(100vh - 64px);
  background: #f5f7fa;
}
.general-card {
  border-radius: 8px;
}
</style>
