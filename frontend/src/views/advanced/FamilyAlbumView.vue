<template>
  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>高级功能</a-breadcrumb-item>
      <a-breadcrumb-item>家庭相册</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 页面标题 -->
    <div class="pro-page-header">
      <h2 class="pro-page-title">家庭相册</h2>
      <p class="pro-page-desc">管理家庭共享照片，支持上传、下载、分享与相册管理</p>
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
        <a-select v-model="albumFilter" placeholder="所属相册" allow-clear style="width: 160px" @change="loadPhotos">
          <a-option v-for="a in albums" :key="a.id" :value="a.id">{{ a.name }}</a-option>
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
        <a-button @click="showAlbumModal(null)">
          <template #icon><icon-folder-add /></template>
          新建相册
        </a-button>
        <a-button @click="loadPhotos">
          <template #icon><icon-refresh /></template>
          刷新
        </a-button>
      </a-space>
      <input ref="fileInputRef" type="file" accept="image/*" multiple style="display:none" @change="handleFileChange" />
    </div>

    <!-- 相册标签页 -->
    <div class="album-tabs" style="margin-bottom: 12px">
      <a-radio-group v-model="viewMode" type="button">
        <a-radio value="grid">网格视图</a-radio>
        <a-radio value="list">列表视图</a-radio>
      </a-radio-group>
      <a-select v-model="albumFilter" placeholder="按相册筛选" allow-clear style="width: 160px; margin-left: 12px" @change="loadPhotos">
        <a-option value="">全部相册</a-option>
        <a-option v-for="a in albums" :key="a.id" :value="a.id">{{ a.name }}</a-option>
      </a-select>
    </div>

    <!-- 数据内容区 -->
    <div class="pro-content-area">
      <!-- 空状态 -->
      <a-empty v-if="!loading && photos.length === 0" description="暂无照片，上传一张开始记录家庭回忆吧" />

      <!-- 照片网格视图 -->
      <div v-else-if="viewMode === 'grid'" class="photo-grid">
        <div v-for="photo in photos" :key="photo.id" class="photo-card">
          <div class="photo-thumb" @click="showPreview(photo)">
            <img :src="photo.url" :alt="photo.name" @error="handleImgError" />
            <div class="photo-overlay">
              <icon-eye @click.stop="showPreview(photo)" />
              <icon-download @click.stop="downloadPhoto(photo)" />
              <icon-share @click.stop="sharePhoto(photo)" />
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
              <a-button type="text" size="mini" @click="downloadPhoto(photo)">下载</a-button>
              <a-button type="text" size="mini" status="danger" @click="handleDelete(photo)">删除</a-button>
            </div>
          </div>
        </div>
      </div>

      <!-- 照片列表视图 -->
      <a-table
        v-else
        :columns="listColumns"
        :data="photos"
        :loading="loading"
        :pagination="pagination"
        @page-change="onPageChange"
        row-key="id"
      >
        <template #thumb="{ record }">
          <img :src="record.url" :alt="record.name" style="width:48px;height:48px;object-fit:cover;border-radius:4px" @error="handleImgError" />
        </template>
        <template #name="{ record }">
          <span style="font-weight:500">{{ record.name }}</span>
        </template>
        <template #uploader_name="{ record }">
          <a-avatar :size="20" :style="{ backgroundColor: '#00b42a' }">
            {{ record.uploader_name?.charAt(0) || '?' }}
          </a-avatar>
          <span style="margin-left:6px">{{ record.uploader_name }}</span>
        </template>
        <template #album_name="{ record }">
          <a-tag>{{ record.album_name || '默认相册' }}</a-tag>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="showPreview(record)">预览</a-button>
            <a-button type="text" size="small" @click="downloadPhoto(record)">下载</a-button>
            <a-button type="text" size="small" @click="sharePhoto(record)">分享</a-button>
            <a-button type="text" size="small" status="danger" @click="handleDelete(record)">删除</a-button>
          </a-space>
        </template>
      </a-table>

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
    <a-modal v-model:visible="uploadProgressVisible" title="上传进度" :closable="false" :footer="null" :width="400">
      <a-progress :percent="uploadPercent" :status="uploadPercent === 100 ? 'success' : 'normal'" />
      <p style="text-align:center;margin-top:8px">{{ uploadStatusText }}</p>
    </a-modal>

    <!-- 照片预览 -->
    <a-modal v-model:visible="previewVisible" title="照片预览" :width="800" @close="previewPhoto = null">
      <div v-if="previewPhoto" class="preview-container">
        <img :src="previewPhoto.url" :alt="previewPhoto.name" class="preview-image" />
        <div class="preview-info">
          <a-descriptions :column="1" size="small">
            <a-descriptions-item label="照片名称">{{ previewPhoto.name }}</a-descriptions-item>
            <a-descriptions-item label="相册">{{ previewPhoto.album_name || '默认相册' }}</a-descriptions-item>
            <a-descriptions-item label="上传者">{{ previewPhoto.uploader_name }}</a-descriptions-item>
            <a-descriptions-item label="上传时间">{{ formatDate(previewPhoto.created_at) }}</a-descriptions-item>
            <a-descriptions-item label="文件大小">{{ formatSize(previewPhoto.size) }}</a-descriptions-item>
          </a-descriptions>
          <div style="margin-top: 12px">
            <a-button type="primary" @click="downloadPhoto(previewPhoto)">
              <template #icon><icon-download /></template>
              下载
            </a-button>
            <a-button style="margin-left: 8px" @click="sharePhoto(previewPhoto)">
              <template #icon><icon-share /></template>
              分享
            </a-button>
          </div>
        </div>
      </div>
    </a-modal>

    <!-- 新建相册弹窗 -->
    <a-modal v-model:visible="albumModalVisible" title="新建相册" @ok="handleAlbumSave" :width="420">
      <a-form :model="albumForm" layout="vertical">
        <a-form-item label="相册名称" required>
          <a-input v-model="albumForm.name" placeholder="请输入相册名称" />
        </a-form-item>
        <a-form-item label="相册描述">
          <a-textarea v-model="albumForm.description" placeholder="请输入相册描述（可选）" :rows="3" />
        </a-form-item>
        <a-form-item label="相册可见性">
          <a-radio-group v-model="albumForm.visibility">
            <a-radio value="family">仅家庭成员</a-radio>
            <a-radio value="public">公开</a-radio>
          </a-radio-group>
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 分享弹窗 -->
    <a-modal v-model:visible="shareModalVisible" title="分享照片" @ok="handleShare" :width="420">
      <a-form :model="shareForm" layout="vertical">
        <a-form-item label="分享方式">
          <a-radio-group v-model="shareForm.type">
            <a-radio value="link">生成分享链接</a-radio>
            <a-radio value="member">分享给家庭成员</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item label="链接有效期" v-if="shareForm.type === 'link'">
          <a-select v-model="shareForm.expire_hours" placeholder="选择有效期">
            <a-option :value="24">24小时</a-option>
            <a-option :value="72">3天</a-option>
            <a-option :value="168">7天</a-option>
            <a-option :value="0">永久有效</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="选择成员" v-if="shareForm.type === 'member'">
          <a-select v-model="shareForm.member_ids" multiple placeholder="选择要分享的家庭成员">
            <a-option v-for="m in members" :key="m.id" :value="m.id">{{ m.name }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="分享链接" v-if="shareForm.type === 'link' && shareLink">
          <a-input :value="shareLink" readonly />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 相册管理抽屉 -->
    <a-drawer v-model:visible="manageDrawerVisible" :width="520" title="相册管理" @close="manageDrawerVisible = false">
      <div class="album-list">
        <div v-for="album in albums" :key="album.id" class="album-item">
          <div class="album-item-info">
            <div class="album-item-name">{{ album.name }}</div>
            <div class="album-item-meta">{{ album.photo_count || 0 }} 张照片 · {{ formatDate(album.created_at) }}</div>
          </div>
          <a-space>
            <a-button type="text" size="small" @click="showAlbumModal(album)">编辑</a-button>
            <a-button type="text" size="small" status="danger" @click="handleDeleteAlbum(album)">删除</a-button>
          </a-space>
        </div>
      </div>
    </a-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const photos = ref<any[]>([])
const albums = ref<any[]>([])
const members = ref<any[]>([])
const loading = ref(false)
const keyword = ref('')
const uploaderFilter = ref('')
const albumFilter = ref('')
const dateRange = ref<any>(null)
const viewMode = ref('grid')
const previewVisible = ref(false)
const previewPhoto = ref<any>(null)
const uploadProgressVisible = ref(false)
const uploadPercent = ref(0)
const uploadStatusText = ref('')
const albumModalVisible = ref(false)
const shareModalVisible = ref(false)
const shareLink = ref('')
const manageDrawerVisible = ref(false)
const fileInputRef = ref<any>(null)

const listColumns = [
  { title: '缩略图', slotName: 'thumb', width: 80 },
  { title: '名称', slotName: 'name', width: 200 },
  { title: '上传者', slotName: 'uploader_name', width: 160 },
  { title: '相册', slotName: 'album_name', width: 140 },
  { title: '上传时间', dataIndex: 'created_at', width: 160 },
  { title: '操作', slotName: 'actions', width: 280 }
]

const pagination = reactive({
  current: 1,
  pageSize: 24,
  total: 0
})

const albumForm = reactive({
  id: null as number | null,
  name: '',
  description: '',
  visibility: 'family'
})

const shareForm = reactive({
  type: 'link',
  expire_hours: 72,
  member_ids: [] as number[]
})

function formatDate(date: string) {
  if (!date) return ''
  return new Date(date).toLocaleDateString('zh-CN')
}

function formatSize(bytes: number) {
  if (!bytes) return '未知'
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
  return (bytes / (1024 * 1024)).toFixed(1) + ' MB'
}

function handleImgError(e: any) {
  e.target.src = 'data:image/svg+xml,<svg xmlns="http://www.w3.org/2000/svg" width="100" height="100"><rect fill="%23f2f3f5" width="100" height="100"/><text x="50" y="55" text-anchor="middle" fill="%23888" font-size="12">图片加载失败</text></svg>'
}

async function loadPhotos() {
  loading.value = true
  try {
    const params = new URLSearchParams()
    if (keyword.value) params.append('keyword', keyword.value)
    if (uploaderFilter.value) params.append('uploader_id', uploaderFilter.value)
    if (albumFilter.value) params.append('album_id', albumFilter.value)
    if (dateRange.value && dateRange.value.length === 2) {
      params.append('start_date', dateRange.value[0])
      params.append('end_date', dateRange.value[1])
    }
    params.append('page', String(pagination.current))
    params.append('page_size', String(pagination.pageSize))

    const res = await fetch(`/api/v1/advanced/album?${params}`, { credentials: 'include' })
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

async function loadAlbums() {
  try {
    const res = await fetch('/api/v1/advanced/album/albums', { credentials: 'include' })
    const data = await res.json()
    albums.value = data.data?.list || data.data || []
  } catch { /* ignore */ }
}

async function loadMembers() {
  try {
    const res = await fetch('/api/v1/family/members', { credentials: 'include' })
    const data = await res.json()
    members.value = data.data?.list || data.data?.members || []
  } catch { /* ignore */ }
}

function onPageChange(page: number) {
  pagination.current = page
  loadPhotos()
}

function triggerUpload() {
  fileInputRef.value?.click()
}

async function handleFileChange(e: any) {
  const files = Array.from(e.target.files) as File[]
  if (!files.length) return

  uploadProgressVisible.value = true
  uploadPercent.value = 0
  uploadStatusText.value = `正在上传 0/${files.length}`

  let completed = 0
  for (const file of files) {
    try {
      const formData = new FormData()
      formData.append('file', file)
      if (albumFilter.value) formData.append('album_id', String(albumFilter.value))

      const res = await fetch('/api/v1/advanced/album/upload', {
        method: 'POST',
        credentials: 'include',
        body: formData
      })
      const data = await res.json()
      if (data.code === 0 || data.code === 200) {
        completed++
        uploadPercent.value = Math.round((completed / files.length) * 100)
        uploadStatusText.value = `正在上传 ${completed}/${files.length}`
      }
    } catch { /* skip failed uploads */ }
  }

  uploadStatusText.value = '上传完成'
  setTimeout(() => {
    uploadProgressVisible.value = false
    loadPhotos()
    loadAlbums()
  }, 1000)
}

function showPreview(photo: any) {
  previewPhoto.value = photo
  previewVisible.value = true
}

async function downloadPhoto(photo: any) {
  try {
    const res = await fetch(photo.url)
    const blob = await res.blob()
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = photo.name
    a.click()
    URL.revokeObjectURL(url)
    Message.success('开始下载')
  } catch {
    Message.error('下载失败')
  }
}

function sharePhoto(photo: any) {
  shareForm.type = 'link'
  shareForm.expire_hours = 72
  shareForm.member_ids = []
  shareLink.value = ''
  previewPhoto.value = photo
  shareModalVisible.value = true
}

async function handleShare() {
  try {
    const payload: any = {
      photo_id: previewPhoto.value?.id,
      type: shareForm.type,
      expire_hours: shareForm.expire_hours
    }
    if (shareForm.type === 'member') {
      payload.member_ids = shareForm.member_ids
    }

    const res = await fetch('/api/v1/advanced/album/share', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify(payload)
    })
    const data = await res.json()
    if (data.code === 0 || data.code === 200) {
      if (shareForm.type === 'link') {
        shareLink.value = data.data?.share_url || ''
        Message.success('分享链接已生成')
      } else {
        Message.success('已分享给选定成员')
        shareModalVisible.value = false
      }
    } else {
      Message.error(data.message || '分享失败')
    }
  } catch {
    Message.error('网络错误')
  }
}

function showAlbumModal(album: any) {
  if (album) {
    albumForm.id = album.id
    albumForm.name = album.name
    albumForm.description = album.description || ''
    albumForm.visibility = album.visibility || 'family'
  } else {
    albumForm.id = null
    albumForm.name = ''
    albumForm.description = ''
    albumForm.visibility = 'family'
  }
  albumModalVisible.value = true
}

async function handleAlbumSave() {
  if (!albumForm.name.trim()) {
    Message.warning('请输入相册名称')
    return
  }
  try {
    const url = albumForm.id
      ? `/api/v1/advanced/album/albums/${albumForm.id}`
      : '/api/v1/advanced/album/albums'
    const method = albumForm.id ? 'PUT' : 'POST'
    const res = await fetch(url, {
      method,
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify(albumForm)
    })
    const data = await res.json()
    if (data.code === 0 || data.code === 200) {
      Message.success(albumForm.id ? '相册已更新' : '相册已创建')
      albumModalVisible.value = false
      loadAlbums()
    } else {
      Message.error(data.message || '操作失败')
    }
  } catch {
    Message.error('网络错误')
  }
}

async function handleDeleteAlbum(album: any) {
  try {
    const res = await fetch(`/api/v1/advanced/album/albums/${album.id}`, {
      method: 'DELETE',
      credentials: 'include'
    })
    const data = await res.json()
    if (data.code === 0 || data.code === 200) {
      Message.success('相册已删除')
      loadAlbums()
    } else {
      Message.error(data.message || '删除失败')
    }
  } catch {
    Message.error('网络错误')
  }
}

async function handleDelete(photo: any) {
  try {
    const res = await fetch(`/api/v1/advanced/album/${photo.id}`, {
      method: 'DELETE',
      credentials: 'include'
    })
    const data = await res.json()
    if (data.code === 0 || data.code === 200) {
      Message.success('已删除')
      loadPhotos()
      loadAlbums()
    } else {
      Message.error(data.message || '删除失败')
    }
  } catch {
    Message.error('网络错误')
  }
}

onMounted(() => {
  loadPhotos()
  loadAlbums()
  loadMembers()
})
</script>

<style scoped>
.photo-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 16px;
  margin-bottom: 16px;
}

.photo-card {
  background: var(--color-bg-2);
  border-radius: 8px;
  overflow: hidden;
  border: 1px solid var(--color-border);
  transition: box-shadow 0.2s;
}
.photo-card:hover {
  box-shadow: 0 4px 16px rgba(0,0,0,0.1);
}

.photo-thumb {
  position: relative;
  padding-top: 75%;
  background: #f2f3f5;
  cursor: pointer;
  overflow: hidden;
}
.photo-thumb img {
  position: absolute;
  top: 0; left: 0;
  width: 100%; height: 100%;
  object-fit: cover;
}
.photo-overlay {
  position: absolute;
  top: 0; left: 0; right: 0; bottom: 0;
  background: rgba(0,0,0,0.4);
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  opacity: 0;
  transition: opacity 0.2s;
}
.photo-thumb:hover .photo-overlay { opacity: 1; }
.photo-overlay .arco-icon {
  color: #fff;
  font-size: 20px;
  cursor: pointer;
}

.photo-info {
  padding: 10px 12px;
}
.photo-name {
  font-size: 13px;
  font-weight: 500;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  margin-bottom: 6px;
}
.photo-meta {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: #666;
  margin-bottom: 6px;
}
.photo-actions {
  display: flex;
  gap: 4px;
}

.photo-pagination {
  display: flex;
  justify-content: center;
  margin-top: 16px;
}

.preview-container {
  display: flex;
  flex-direction: column;
  gap: 16px;
}
.preview-image {
  max-width: 100%;
  max-height: 60vh;
  object-fit: contain;
  border-radius: 8px;
}

.album-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
.album-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px;
  border: 1px solid var(--color-border);
  border-radius: 8px;
}
.album-item-name {
  font-weight: 500;
  margin-bottom: 4px;
}
.album-item-meta {
  font-size: 12px;
  color: #666;
}
</style>
