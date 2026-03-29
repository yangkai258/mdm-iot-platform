<template>
  <div class="page-container">
    <Breadcrumb :items="['menu.apps', 'menu.apps.versions']" />
    <!-- 应用信息卡片 -->
    <a-card v-if="appInfo" class="app-info-card">
      <a-row :gutter="16" align="center">
        <a-col :span="4">
          <a-avatar :size="48" shape="square">
            <img v-if="appInfo.icon_url" :src="appInfo.icon_url" />
            <span v-else>{{ appInfo.app_name?.charAt(0) }}</span>
          </a-avatar>
        </a-col>
        <a-col :span="16">
          <div style="font-size: 18px; font-weight: 600">{{ appInfo.app_name }}</div>
          <div style="color: #666; margin-top: 4px">
            {{ appInfo.bundle_id }} · {{ appInfo.app_type?.toUpperCase() }}
          </div>
        </a-col>
        <a-col :span="4" style="text-align: right">
          <a-button type="primary" @click="showUploadDrawer = true">上传新版本</a-button>
          <a-button style="margin-top: 8px" @click="router.push('/apps/list')">返回列表</a-button>
        </a-col>
      </a-row>
    </a-card>

    <!-- 版本列表 -->
    <a-card class="version-card">
      <template #title><span>版本列表</span></template>
      <a-table :columns="columns" :data="versionList" :loading="loading" :pagination="false" row-key="id">
        <template #version="{ record }">
          <a-tag :color="record.is_latest ? 'blue' : 'gray'">{{ record.version }}<span v-if="record.is_latest"> (最新)</span></a-tag>
        </template>
        <template #mandatory="{ record }">
          <a-tag :color="record.is_mandatory ? 'orange' : 'default'">{{ record.is_mandatory ? '强制更新' : '可选更新' }}</a-tag>
        </template>
        <template #fileSize="{ record }">{{ formatFileSize(record.file_size) }}</template>
        <template #installCount="{ record }">{{ record.install_count || 0 }}</template>
        <template #createdAt="{ record }">{{ formatTime(record.created_at) }}</template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="openDetail(record)">详情</a-button>
            <a-button type="text" size="small" status="danger" @click="deleteVersion(record)">删除</a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 上传版本抽屉 -->
    <a-drawer v-model:visible="showUploadDrawer" title="上传新版本" :width="480">
      <a-form :model="uploadForm" layout="vertical" @submit-success="handleUploadSubmit">
        <a-form-item label="版本号" required><a-input v-model="uploadForm.version" placeholder="如 v2.1.0" /></a-form-item>
        <a-form-item label="版本数字" required><a-input-number v-model="uploadForm.version_code" :min="1" placeholder="如 21" style="width: 100%" /></a-form-item>
        <a-form-item label="安装包 URL" required><a-input v-model="uploadForm.file_url" placeholder="https://cdn.example.com/app-v2.1.0.apk" /></a-form-item>
        <a-form-item label="文件大小 (字节)"><a-input-number v-model="uploadForm.file_size" :min="0" placeholder="如 52428800" style="width: 100%" /></a-form-item>
        <a-form-item label="MD5 校验码"><a-input v-model="uploadForm.file_md5" placeholder="文件 MD5 校验码" /></a-form-item>
        <a-form-item label="最低系统版本"><a-input v-model="uploadForm.min_os_version" placeholder="如 14.0" /></a-form-item>
        <a-form-item label="更新日志"><a-textarea v-model="uploadForm.release_notes" :rows="4" placeholder="本次更新内容..." /></a-form-item>
        <a-form-item label="强制更新"><a-switch v-model="uploadForm.is_mandatory" /></a-form-item>
        <a-form-item>
          <a-space>
            <a-button type="primary" html-type="submit">上传</a-button>
            <a-button @click="showUploadDrawer = false">取消</a-button>
          </a-space>
        </a-form-item>
      </a-form>
    </a-drawer>

    <!-- 版本详情抽屉 -->
    <a-drawer v-model:visible="showDetailDrawer" title="版本详情" :width="480">
      <a-descriptions v-if="currentVersion" :column="1" bordered size="small">
        <a-descriptions-item label="版本号">{{ currentVersion.version }}</a-descriptions-item>
        <a-descriptions-item label="版本数字">{{ currentVersion.version_code }}</a-descriptions-item>
        <a-descriptions-item label="最低系统版本">{{ currentVersion.min_os_version || '-' }}</a-descriptions-item>
        <a-descriptions-item label="文件大小">{{ formatFileSize(currentVersion.file_size) }}</a-descriptions-item>
        <a-descriptions-item label="MD5">{{ currentVersion.file_md5 || '-' }}</a-descriptions-item>
        <a-descriptions-item label="安装包 URL"><a :href="currentVersion.file_url" target="_blank">{{ currentVersion.file_url }}</a></a-descriptions-item>
        <a-descriptions-item label="强制更新">
          <a-tag :color="currentVersion.is_mandatory ? 'orange' : 'default'">{{ currentVersion.is_mandatory ? '是' : '否' }}</a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="最新版本">
          <a-tag :color="currentVersion.is_latest ? 'blue' : 'default'">{{ currentVersion.is_latest ? '是' : '否' }}</a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="安装次数">{{ currentVersion.install_count || 0 }}</a-descriptions-item>
        <a-descriptions-item label="更新日志">{{ currentVersion.release_notes || '-' }}</a-descriptions-item>
        <a-descriptions-item label="上传时间">{{ formatTime(currentVersion.created_at) }}</a-descriptions-item>
      </a-descriptions>
    </a-drawer>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import axios from 'axios'
import { Message, Modal } from '@arco-design/web-vue'

const router = useRouter()
const route = useRoute()
const loading = ref(false)
const appInfo = ref(null)
const versionList = ref([])
const showUploadDrawer = ref(false)
const showDetailDrawer = ref(false)
const currentVersion = ref(null)

const uploadForm = reactive({
  version: '',
  version_code: null,
  file_url: '',
  file_size: null,
  file_md5: '',
  min_os_version: '',
  release_notes: '',
  is_mandatory: false
})

const columns = [
  { title: '版本', slotName: 'version', width: 120 },
  { title: '版本数字', dataIndex: 'version_code', width: 100 },
  { title: '更新类型', slotName: 'mandatory', width: 100 },
  { title: '最低系统版本', dataIndex: 'min_os_version', width: 120 },
  { title: '文件大小', slotName: 'fileSize', width: 120 },
  { title: '安装次数', slotName: 'installCount', width: 100 },
  { title: '上传时间', slotName: 'createdAt', width: 160 },
  { title: '操作', slotName: 'actions', width: 120, fixed: 'right' }
]

const appId = parseInt(route.params.id)

const loadAppInfo = async () => {
  try {
    const res = await axios.get(`/api/v1/apps/${appId}`)
    if (res.data.code === 0) appInfo.value = res.data.data
  } catch { Message.error('加载应用信息失败') }
}

const loadVersions = async () => {
  loading.value = true
  try {
    const res = await axios.get(`/api/v1/apps/${appId}/versions`)
    if (res.data.code === 0) versionList.value = res.data.data.list || []
  } catch { Message.error('加载版本列表失败') }
  finally { loading.value = false }
}

const handleUploadSubmit = async () => {
  try {
    const res = await axios.post(`/api/v1/apps/${appId}/versions`, uploadForm)
    if (res.data.code === 0) {
      Message.success('上传成功')
      showUploadDrawer.value = false
      resetUploadForm()
      loadVersions()
    } else {
      Message.error(res.data.message || '上传失败')
    }
  } catch { Message.error('上传失败') }
}

const resetUploadForm = () => {
  Object.assign(uploadForm, { version: '', version_code: null, file_url: '', file_size: null, file_md5: '', min_os_version: '', release_notes: '', is_mandatory: false })
}

const openDetail = (record) => { currentVersion.value = record; showDetailDrawer.value = true }
const deleteVersion = (record) => {
  Modal.confirm({
    title: '确认删除',
    content: `确定要删除版本 ${record.version} 吗？`,
    onOk: async () => {
      try {
        const res = await axios.delete(`/api/v1/apps/${appId}/versions/${record.id}`)
        if (res.data.code === 0) { Message.success('删除成功'); loadVersions() }
        else Message.error(res.data.message || '删除失败')
      } catch { Message.error('删除失败') }
    }
  })
}

const formatFileSize = (bytes) => {
  if (!bytes) return '-'
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
  if (bytes < 1024 * 1024 * 1024) return (bytes / (1024 * 1024)).toFixed(1) + ' MB'
  return (bytes / (1024 * 1024 * 1024)).toFixed(1) + ' GB'
}

const formatTime = (time) => time ? new Date(time).toLocaleString('zh-CN') : '-'

onMounted(() => { loadAppInfo(); loadVersions() })
</script>

<style scoped>
.app-info-card { margin-bottom: 16px; }
</style>
