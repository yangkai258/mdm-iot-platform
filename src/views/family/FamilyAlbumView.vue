<template>
  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>家庭管理</a-breadcrumb-item>
      <a-breadcrumb-item>家庭相册</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 页面标题 -->
    <div class="pro-page-header">
      <h2 class="pro-page-title">家庭相册</h2>
      <p class="pro-page-desc">管理家庭共享照片，珍藏每一刻美好回忆</p>
    </div>

    <!-- 搜索筛选区 -->
    <div class="pro-search-bar">
      <a-space>
        <a-input-search
          v-model="keyword"
          placeholder="搜索照片名称"
          style="width: 260px"
          @search="loadPhotos"
          search-button
        />
        <a-select v-model="uploaderFilter" placeholder="上传者" allow-clear style="width: 160px" @change="loadPhotos">
          <a-option v-for="m in members" :key="m.id" :value="m.id">{{ m.name }}</a-option>
        </a-select>
        <a-range-picker v-model="dateRange" style="width: 260px" @change="loadPhotos" />
      </a-space>
    </div>

    <!-- 操作按钮区 -->
    <div class="pro-action-bar">
      <a-space>
        <a-button type="primary" @click="triggerUpload">
          <template #icon><icon-upload /></template>
          上传照片
        </a-button>
        <a-button @click="loadPhotos">
          <template #icon><icon-refresh /></template>
          刷新
        </a-button>
      </a-space>
      <input ref="fileInputRef" type="file" accept="image/*" multiple style="display:none" @change="handleFileChange" />
    </div>

    <!-- 数据内容区 - 照片网格 -->
    <div class="pro-content-area">
      <!-- 空状态 -->
      <a-empty v-if="!loading && photos.length === 0" description="暂无照片，上传一张开始记录家庭回忆吧" />

      <!-- 照片网格 -->
      <div v-else class="photo-grid">
        <div v-for="photo in photos" :key="photo.id" class="photo-card">
          <div class="photo-thumb" @click="showPreview(photo)">
            <img :src="photo.url" :alt="photo.name" @error="handleImgError" />
            <div class="photo-overlay">
              <icon-eye @click.stop="showPreview(photo)" />
            </div>
          </div>
          <div class="photo-info">
            <div class="photo-name" :title="photo.name">{{ photo.name }}</div>
            <div class="photo-meta">
              <a-avatar :size="20" :style="{ backgroundColor: '#00b42a' }">
                {{ photo.uploader_name?.charAt(0) || '?' }}
              </a-avatar>
              <span class="photo-uploader">{{ photo.uploader_name }}</span>
              <span class="photo-date">{{ formatDate(photo.created_at) }}</span>
            </div>
            <div class="photo-actions">
              <a-button type="text" size="mini" @click="showPreview(photo)">预览</a-button>
              <a-button type="text" size="mini" status="danger" @click="handleDelete(photo)">删除</a-button>
            </div>
          </div>
        </div>
      </div>

      <!-- 分页 -->
      <div class="photo-pagination" v-if="pagination.total > 0">
        <a-pagination
          :current="pagination.current"
          :page-size="pagination.pageSize"
          :total="pagination.total"
          @change="onPageChange"
          show-total
        />
      </div>
    </div>

    <!-- 上传进度 -->
    <a-modal v-model:visible="uploadProgressVisible" title="上传中" :closable="false" :footer="null" :width="400">
      <a-progress :percent="uploadPercent" :status="'normal'" />
      <p style="text-align:center;margin-top:8px">{{ uploadStatusText }}</p>
    </a-modal>

    <!-- 照片预览 -->
    <a-modal v-model:visible="previewVisible" title="照片预览" :width="800" @close="previewPhoto = null">
      <div v-if="previewPhoto" class="preview-container">
        <img :src="previewPhoto.url" :alt="previewPhoto.name" class="preview-image" />
        <div class="preview-info">
          <a-descriptions :column="1" size="small">
            <a-descriptions-item label="照片名称">{{ previewPhoto.name }}</a-descriptions-item>
            <a-descriptions-item label="上传者">{{ previewPhoto.uploader_name }}</a-descriptions-item>
            <a-descriptions-item label="上传时间">{{ formatDate(previewPhoto.created_at) }}</a-descriptions-item>
            <a-descriptions-item label="文件大小">{{ formatSize(previewPhoto.size) }}</a-descriptions-item>
          </a-descriptions>
        </div>
      </div>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const photos = ref<any[]>([])
const members = ref<any[]>([])
const loading = ref(false)
const keyword = ref('')
const uploaderFilter = ref('')
const dateRange = ref<any[]>([])
const fileInputRef = ref<HTMLInputElement | null>(null)
const uploadProgressVisible = ref(false)
const uploadPercent = ref(0)
const uploadStatusText = ref('')
const previewVisible = ref(false)
const previewPhoto = ref<any | null>(null)

const pagination = reactive({
  current: 1,
  pageSize: 12,
  total: 0
})

async function loadMembers() {
  try {
    const res = await fetch('/api/v1/family/members', { credentials: 'include' })
    const data = await res.json()
    members.value = data.data?.list || data.data?.members || []
  } catch { /* ignore */ }
}

async function loadPhotos() {
  loading.value = true
  try {
    const params = new URLSearchParams()
    if (keyword.value) params.append('keyword', keyword.value)
    if (uploaderFilter.value) params.append('uploader_id', uploaderFilter.value)
    if (dateRange.value && dateRange.value.length === 2) {
      params.append('start_date', dateRange.value[0])
      params.append('end_date', dateRange.value[1])
    }
    params.append('page', String(pagination.current))
    params.append('page_size', String(pagination.pageSize))

    const res = await fetch(`/api/v1/family/album?${params}`, { credentials: 'include' })
    const data = await res.json()
    if (data.code === 0 || data.code === 200) {
      photos.value = data.data?.list || data.data || []
      pagination.total = data.data?.total || 0
    } else {
      Message.error(data.message || '加载失败')
    }
  } catch {
    Message.error('网络错误')
  } finally {
    loading.value = false
  }
}

function onPageChange(page: number) {
  pagination.current = page
  loadPhotos()
}

function triggerUpload() {
  fileInputRef.value?.click()
}

async function handleFileChange(e: Event) {
  const input = e.target as HTMLInputElement
  const files = input.files
  if (!files || files.length === 0) return

  uploadProgressVisible.value = true
  uploadPercent.value = 0
  uploadStatusText.value = `准备上传 ${files.length} 张照片...`

  for (let i = 0; i < files.length; i++) {
    uploadStatusText.value = `上传中 ${i + 1}/${files.length}: ${files[i].name}`
    uploadPercent.value = Math.round(((i + 1) / files.length) * 100)

    const formData = new FormData()
    formData.append('file', files[i])

    try {
      const res = await fetch('/api/v1/family/album/upload', {
        method: 'POST',
        credentials: 'include',
        body: formData
      })
      const data = await res.json()
      if (data.code !== 0 && data.code !== 200) {
        Message.error(`${files[i].name} 上传失败: ${data.message}`)
      }
    } catch {
      Message.error(`${files[i].name} 上传失败`)
    }
  }

  uploadStatusText.value = '上传完成'
  setTimeout(() => {
    uploadProgressVisible.value = false
    uploadPercent.value = 0
    input.value = ''
    loadPhotos()
  }, 500)
}

function showPreview(photo: any) {
  previewPhoto.value = photo
  previewVisible.value = true
}

function handleImgError(e: Event) {
  const img = e.target as HTMLImageElement
  img.src = 'data:image/svg+xml,<svg xmlns="http://www.w3.org/2000/svg" width="200" height="150" viewBox="0 0 200 150"><rect fill="%23f0f0f0" width="200" height="150"/><text x="50%" y="50%" text-anchor="middle" fill="%23888" font-size="14">图片加载失败</text></svg>'
}

function formatDate(dateStr: string) {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleDateString('zh-CN')
}

function formatSize(bytes: number) {
  if (!bytes) return '-'
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
  return (bytes / (1024 * 1024)).toFixed(1) + ' MB'
}

async function handleDelete(photo: any) {
  try {
    const res = await fetch(`/api/v1/family/album/${photo.id}`, {
      method: 'DELETE',
      credentials: 'include'
    })
    const data = await res.json()
    if (data.code === 0 || data.code === 200) {
      Message.success('照片已删除')
      loadPhotos()
    } else {
      Message.error(data.message || '删除失败')
    }
  } catch {
    Message.error('网络错误')
  }
}

onMounted(() => {
  loadMembers()
  loadPhotos()
})
</script>

<style scoped>
.photo-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
  gap: 16px;
  padding: 4px 0;
}

.photo-card {
  border: 1px solid var(--color-neutral-3, #e5e6e8);
  border-radius: 8px;
  overflow: hidden;
  background: #fff;
  transition: box-shadow 0.2s;
}

.photo-card:hover {
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
}

.photo-thumb {
  position: relative;
  width: 100%;
  height: 180px;
  overflow: hidden;
  cursor: pointer;
  background: #f5f5f5;
}

.photo-thumb img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s;
}

.photo-card:hover .photo-thumb img {
  transform: scale(1.05);
}

.photo-overlay {
  position: absolute;
  inset: 0;
  background: rgba(0, 0, 0, 0.3);
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transition: opacity 0.2s;
  color: #fff;
  font-size: 28px;
}

.photo-thumb:hover .photo-overlay {
  opacity: 1;
}

.photo-info {
  padding: 12px;
}

.photo-name {
  font-size: 14px;
  font-weight: 500;
  color: var(--color-text-1, #1f1f1f);
  margin-bottom: 8px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.photo-meta {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-bottom: 8px;
}

.photo-uploader {
  font-size: 12px;
  color: var(--color-text-3, #86909c);
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.photo-date {
  font-size: 12px;
  color: var(--color-text-3, #86909c);
  flex-shrink: 0;
}

.photo-actions {
  display: flex;
  gap: 4px;
}

.photo-pagination {
  margin-top: 24px;
  display: flex;
  justify-content: center;
}

.preview-container {
  text-align: center;
}

.preview-image {
  max-width: 100%;
  max-height: 500px;
  border-radius: 8px;
  object-fit: contain;
}

.preview-info {
  margin-top: 16px;
  text-align: left;
}
</style>
