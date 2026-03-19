<template>
  <a-layout class="ota-firmware">
    <a-layout-sider v-model:collapsed="collapsed" :trigger="null" collapsible>
      <div class="logo">
        <span v-if="!collapsed">MDM 控制台</span>
      </div>
      <a-menu v-model:selectedKeys="selectedKeys" theme="dark" mode="inline" @click="handleMenuClick">
        <a-menu-item key="dashboard">
          <span>设备大盘</span>
        </a-menu-item>
        <a-menu-item key="status">
          <span>设备状态</span>
        </a-menu-item>
        <a-menu-item key="pet">
          <span>宠物配置</span>
        </a-menu-item>
        <a-menu-item key="ota">
          <span>OTA 固件</span>
        </a-menu-item>
      </a-menu>
    </a-layout-sider>

    <a-layout>
      <a-layout-header class="header">
        <div class="header-left">
          <a-button type="text" @click="collapsed = !collapsed">
            <span v-if="collapsed">☰</span>
            <span v-else>✕</span>
          </a-button>
        </div>
        <div class="header-title">
          <span>OTA 固件管理</span>
        </div>
        <div class="header-right">
        </div>
      </a-layout-header>

      <a-layout-content class="content">
        <!-- 统计卡片 -->
        <a-row :gutter="16" class="stats-row">
          <a-col :span="6">
            <a-card>
              <a-statistic title="固件版本" :value="stats.versions" />
            </a-card>
          </a-col>
          <a-col :span="6">
            <a-card>
              <a-statistic title="设备总数" :value="stats.totalDevices" />
            </a-card>
          </a-col>
          <a-col :span="6">
            <a-card>
              <a-statistic title="已升级" :value="stats.upgraded" :value-style="{ color: '#52c41a' }" />
            </a-card>
          </a-col>
          <a-col :span="6">
            <a-card>
              <a-statistic title="待升级" :value="stats.pending" :value-style="{ color: '#faad14' }" />
            </a-card>
          </a-col>
        </a-row>

        <!-- 操作栏 -->
        <a-card class="action-card">
          <a-space>
            <a-button type="primary" @click="showUploadModal">上传固件</a-button>
            <a-button @click="loadFirmwares">刷新</a-button>
          </a-space>
        </a-card>

        <!-- 固件列表 -->
        <a-card class="firmware-card">
          <template #title>
            <span>固件列表</span>
          </template>

          <a-table
            :columns="columns"
            :data="firmwares"
            :loading="loading"
            :pagination="false"
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
                <a-button type="primary" size="small" @click="handleDeploy(record)">下发升级</a-button>
                <a-button size="small" @click="handleDelete(record)">删除</a-button>
              </a-space>
            </template>
          </a-table>
        </a-card>

        <!-- 升级任务 -->
        <a-card class="task-card">
          <template #title>
            <span>升级任务</span>
          </template>

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
        <a-modal v-model:visible="uploadModalVisible" title="上传固件" @ok="handleUpload" :confirm-loading="uploading">
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
              <a-upload
                :auto-upload="false"
                :show-upload-list="true"
                @change="handleFileChange"
              >
                <a-button>选择文件</a-button>
              </a-upload>
            </a-form-item>
          </a-form>
        </a-modal>
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Message } from '@arco-design/web-vue'
import axios from 'axios'

const router = useRouter()

const collapsed = ref(false)
const selectedKeys = ref(['ota'])
const loading = ref(false)
const taskLoading = ref(false)
const uploading = ref(false)
const uploadModalVisible = ref(false)

const stats = reactive({
  versions: 3,
  totalDevices: 156,
  upgraded: 142,
  pending: 14
})

const uploadForm = reactive({
  version: '',
  hardware_model: '',
  description: '',
  file: null
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

// 模拟数据
const firmwares = ref([
  { id: 1, version: '1.2.0', hardware_model: 'MDM-Pro-200', file_size: 5242880, status: 1, upload_time: '2026-03-19 10:00:00', description: '修复已知问题' },
  { id: 2, version: '1.1.5', hardware_model: 'MDM-Pro-200', file_size: 5242880, status: 2, upload_time: '2026-03-15 14:30:00', description: '优化性能' },
  { id: 3, version: '1.0.0', hardware_model: 'MDM-Pro-200', file_size: 5242880, status: 3, upload_time: '2026-03-01 09:00:00', description: '初始版本' }
])

const tasks = ref([
  { id: 1, firmware_version: 'v1.2.0', target_devices: 50, status: 1, progress: 100, create_time: '2026-03-19 08:00:00' },
  { id: 2, firmware_version: 'v1.2.0', target_devices: 100, status: 2, progress: 45, create_time: '2026-03-19 10:00:00' }
])

const API_BASE = 'http://localhost:8080/api/v1'

const loadFirmwares = async () => {
  loading.value = true
  try {
    const res = await axios.get(`${API_BASE}/firmwares`)
    if (res.data.code === 0) {
      firmwares.value = res.data.data.list
    }
  } catch (err) {
    // 使用模拟数据
    console.log('使用模拟固件数据')
  } finally {
    loading.value = false
  }
}

const loadTasks = async () => {
  taskLoading.value = true
  try {
    const res = await axios.get(`${API_BASE}/ota/tasks`)
    if (res.data.code === 0) {
      tasks.value = res.data.data.list
    }
  } catch (err) {
    // 使用模拟数据
    console.log('使用模拟任务数据')
  } finally {
    taskLoading.value = false
  }
}

const handleMenuClick = ({ key }) => {
  if (key === 'dashboard') {
    router.push('/dashboard')
  } else if (key === 'ota') {
    router.push('/ota')
  } else if (key === 'pet') {
    router.push('/pet')
  } else if (key === 'status') {
    router.push('/status')
  }
}

const showUploadModal = () => {
  uploadModalVisible.value = true
}

const handleFileChange = (file) => {
  uploadForm.file = file
}

const handleUpload = () => {
  if (!uploadForm.version || !uploadForm.hardware_model) {
    Message.warning('请填写必填项')
    return
  }
  
  uploading.value = true
  
  // 模拟上传
  setTimeout(() => {
    firmwares.value.unshift({
      id: Date.now(),
      version: uploadForm.version,
      hardware_model: uploadForm.hardware_model,
      file_size: 5242880,
      status: 1,
      upload_time: new Date().toLocaleString(),
      description: uploadForm.description
    })
    
    uploadModalVisible.value = false
    uploading.value = false
    Message.success('固件上传成功')
    
    // 重置表单
    uploadForm.version = ''
    uploadForm.hardware_model = ''
    uploadForm.description = ''
    uploadForm.file = null
  }, 1500)
}

const handleDeploy = (record) => {
  Message.info(`开始下发固件 ${record.version} 到设备...`)
  
  // 添加新任务
  tasks.value.unshift({
    id: Date.now(),
    firmware_version: `v${record.version}`,
    target_devices: stats.pending,
    status: 2,
    progress: 0,
    create_time: new Date().toLocaleString()
  })
  
  Message.success('升级任务已创建')
}

const handleDelete = (record) => {
  firmwares.value = firmwares.value.filter(f => f.id !== record.id)
  Message.success('固件已删除')
}

const getStatusBadge = (status) => {
  const badges = { 1: 'success', 2: 'normal', 3: 'default' }
  return badges[status] || 'default'
}

const getStatusText = (status) => {
  const texts = { 1: '已发布', 2: '测试中', 3: '已废弃' }
  return texts[status] || '未知'
}

const getTaskStatusColor = (status) => {
  const colors = { 1: 'green', 2: 'blue', 3: 'red', 4: 'gray' }
  return colors[status] || 'default'
}

const getTaskStatusText = (status) => {
  const texts = { 1: '已完成', 2: '进行中', 3: '失败', 4: '已取消' }
  return texts[status] || '未知'
}

const formatFileSize = (bytes) => {
  if (!bytes) return '-'
  const mb = bytes / 1024 / 1024
  return mb.toFixed(2) + ' MB'
}

const formatTime = (time) => {
  if (!time) return '-'
  return time
}

onMounted(() => {
  loadFirmwares()
  loadTasks()
})
</script>

<style scoped>
.ota-firmware {
  min-height: 100vh;
}

.logo {
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  font-weight: bold;
  color: #fff;
}

.header {
  background: #fff;
  padding: 0 16px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
}

.header-left, .header-right {
  display: flex;
  align-items: center;
}

.header-title {
  font-size: 16px;
  font-weight: 500;
}

.content {
  margin: 16px;
}

.stats-row {
  margin-bottom: 16px;
}

.action-card, .firmware-card, .task-card {
  margin-bottom: 16px;
}
</style>
