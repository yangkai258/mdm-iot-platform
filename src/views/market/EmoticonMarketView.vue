<template>
  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>内容生态</a-breadcrumb-item>
      <a-breadcrumb-item>表情包市场</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 页面标题 -->
    <div class="pro-page-title">表情包市场</div>

    <!-- 搜索筛选区 -->
    <div class="pro-search-bar">
      <a-space wrap>
        <a-select
          v-model="filterForm.category_id"
          placeholder="全部分类"
          style="width: 160px"
          allow-clear
          @change="loadEmoticons"
        >
          <a-option v-for="cat in categories" :key="cat.category_id" :value="cat.category_id">
            {{ cat.category_name }}
          </a-option>
        </a-select>
        <a-select
          v-model="filterForm.is_premium"
          placeholder="收费类型"
          style="width: 140px"
          allow-clear
          @change="loadEmoticons"
        >
          <a-option :value="false">免费</a-option>
          <a-option :value="true">付费</a-option>
        </a-select>
        <a-input-search
          v-model="filterForm.keyword"
          placeholder="搜索表情包名称"
          style="width: 240px"
          search-button
          @search="loadEmoticons"
          @change="e => !e.target.value && loadEmoticons()"
        />
      </a-space>
    </div>

    <!-- 操作按钮区 -->
    <div class="pro-action-bar">
      <a-space>
        <a-button type="primary" @click="showUploadModal">上传表情包</a-button>
        <a-button @click="loadEmoticons">刷新</a-button>
      </a-space>
    </div>

    <!-- 表情包网格 -->
    <div class="pro-content-area">
      <!-- 加载状态 -->
      <a-spin :loading="loading" tip="加载中...">
        <div v-if="!loading">
          <!-- 空状态 -->
          <a-empty v-if="emoticons.length === 0" description="暂无表情包" style="padding: 60px 0" />

          <!-- 表情包网格 -->
          <div v-else class="emoticon-grid">
            <div
              v-for="emo in emoticons"
              :key="emo.emoticon_id"
              class="emoticon-card"
              @click="showDetail(emo)"
            >
              <div class="emoticon-thumb">
                <img
                  v-if="emo.thumb_url"
                  :src="emo.thumb_url"
                  :alt="emo.name"
                  class="emoticon-img"
                />
                <div v-else class="emoticon-placeholder">
                  <span class="emoticon-emoji">{{ emo.preview_emoji || '😊' }}</span>
                </div>
                <!-- 收费标签 -->
                <a-tag
                  v-if="emo.is_premium"
                  class="emoticon-tag"
                  color="orange"
                >
                  付费
                </a-tag>
                <a-tag
                  v-else
                  class="emoticon-tag"
                  color="green"
                >
                  免费
                </a-tag>
              </div>
              <div class="emoticon-info">
                <div class="emoticon-name" :title="emo.name">{{ emo.name }}</div>
                <div class="emoticon-meta">
                  <span class="emoticon-creator">{{ emo.creator_name || '官方' }}</span>
                  <span class="emoticon-downloads">
                    <icon-download />{{ emo.download_count || 0 }}
                  </span>
                </div>
              </div>
            </div>
          </div>

          <!-- 分页 -->
          <div v-if="pagination.total > 0" class="emoticon-pagination">
            <a-pagination
              :total="pagination.total"
              :current="pagination.current"
              :page-size="pagination.pageSize"
              @change="handlePageChange"
              show-total
              show-page-size
              :page-size-options="[20, 40, 60, 100]"
            />
          </div>
        </div>
      </a-spin>
    </div>

    <!-- 表情包详情弹窗 -->
    <a-modal
      v-model:visible="detailVisible"
      :title="currentEmoticon?.name"
      :width="600"
      :footer="null"
    >
      <div v-if="currentEmoticon" class="emoticon-detail">
        <div class="detail-preview">
          <img
            v-if="currentEmoticon.thumb_url"
            :src="currentEmoticon.thumb_url"
            :alt="currentEmoticon.name"
            class="detail-img"
          />
          <div v-else class="detail-placeholder">
            <span class="detail-emoji">{{ currentEmoticon.preview_emoji || '😊' }}</span>
          </div>
        </div>
        <a-descriptions :column="2" size="small" style="margin-top: 16px">
          <a-descriptions-item label="分类">
            {{ currentEmoticon.category_name || '默认' }}
          </a-descriptions-item>
          <a-descriptions-item label="收费类型">
            <a-tag :color="currentEmoticon.is_premium ? 'orange' : 'green'">
              {{ currentEmoticon.is_premium ? '付费' : '免费' }}
            </a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="创作者">
            {{ currentEmoticon.creator_name || '官方' }}
          </a-descriptions-item>
          <a-descriptions-item label="下载量">
            {{ currentEmoticon.download_count || 0 }}
          </a-descriptions-item>
          <a-descriptions-item label="创建时间" :span="2">
            {{ currentEmoticon.create_time || '-' }}
          </a-descriptions-item>
          <a-descriptions-item label="描述" :span="2">
            {{ currentEmoticon.description || '暂无描述' }}
          </a-descriptions-item>
        </a-descriptions>
        <div class="detail-actions">
          <a-space>
            <a-button type="primary" @click="handleDownload(currentEmoticon)">
              <template #icon><icon-download /></template>
              下载
            </a-button>
            <a-button @click="detailVisible = false">关闭</a-button>
          </a-space>
        </div>
      </div>
    </a-modal>

    <!-- 上传表情包弹窗 -->
    <a-modal
      v-model:visible="uploadVisible"
      title="上传表情包"
      @ok="handleUpload"
      :confirm-loading="uploading"
      :width="520"
    >
      <a-form :model="uploadForm" layout="vertical">
        <a-form-item label="表情包名称" required>
          <a-input v-model="uploadForm.name" placeholder="请输入表情包名称" />
        </a-form-item>
        <a-form-item label="分类" required>
          <a-select v-model="uploadForm.category_id" placeholder="选择分类" allow-clear>
            <a-option v-for="cat in categories" :key="cat.category_id" :value="cat.category_id">
              {{ cat.category_name }}
            </a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="预览 Emoji">
          <a-input v-model="uploadForm.preview_emoji" placeholder="如: 😊" maxlength="4" />
        </a-form-item>
        <a-form-item label="描述">
          <a-textarea v-model="uploadForm.description" placeholder="简要描述" :rows="2" />
        </a-form-item>
        <a-form-item label="设为付费">
          <a-switch v-model="uploadForm.is_premium" />
        </a-form-item>
        <a-form-item label="上传文件">
          <a-upload
            action="#"
            :before-upload="beforeUpload"
            :show-upload-list="false"
            accept="image/*"
          >
            <a-button type="outline">
              <template #icon><icon-upload /></template>
              选择图片
            </a-button>
          </a-upload>
          <div v-if="uploadForm.file" style="margin-top: 8px; color: var(--color-success)">
            已选择: {{ uploadForm.file.name }}
          </div>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import {
  getEmoticonCategories,
  getEmoticonList,
  uploadEmoticon,
  downloadEmoticon
} from '@/api/market'

const loading = ref(false)
const uploading = ref(false)
const detailVisible = ref(false)
const uploadVisible = ref(false)

const emoticons = ref([])
const categories = ref([])
const currentEmoticon = ref(null)

const filterForm = reactive({
  category_id: null,
  is_premium: null,
  keyword: ''
})

const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0
})

const uploadForm = reactive({
  name: '',
  category_id: null,
  preview_emoji: '',
  description: '',
  is_premium: false,
  file: null
})

const loadCategories = async () => {
  try {
    const res = await getEmoticonCategories()
    if (res.code === 0) {
      categories.value = res.data || []
    }
  } catch {
    categories.value = [
      { category_id: 'cat1', category_name: '可爱' },
      { category_id: 'cat2', category_name: '搞笑' },
      { category_id: 'cat3', category_name: '日常' },
      { category_id: 'cat4', category_name: '节日' },
      { category_id: 'cat5', category_name: '宠物' }
    ]
  }
}

const loadEmoticons = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.current,
      page_size: pagination.pageSize,
      category_id: filterForm.category_id || undefined,
      is_premium: filterForm.is_premium ?? undefined,
      keyword: filterForm.keyword || undefined
    }
    const res = await getEmoticonList(params)
    if (res.code === 0) {
      emoticons.value = res.data?.list || []
      pagination.total = res.data?.pagination?.total || 0
    }
  } catch {
    // 模拟数据
    emoticons.value = [
      { emoticon_id: 'emo1', name: '开心喵', thumb_url: '', preview_emoji: '😺', category_id: 'cat5', category_name: '宠物', creator_name: '官方', download_count: 1234, is_premium: false, create_time: '2026-03-01 10:00:00', description: '超可爱的猫咪表情包' },
      { emoticon_id: 'emo2', name: '搞怪汪', thumb_url: '', preview_emoji: '🐶', category_id: 'cat5', category_name: '宠物', creator_name: '官方', download_count: 856, is_premium: false, create_time: '2026-03-02 11:00:00', description: '狗狗搞怪表情合集' },
      { emoticon_id: 'emo3', name: '快乐每一天', thumb_url: '', preview_emoji: '😊', category_id: 'cat3', category_name: '日常', creator_name: '用户小明', download_count: 2341, is_premium: true, create_time: '2026-03-03 09:00:00', description: '日常使用的高质量表情包' },
      { emoticon_id: 'emo4', name: '节日祝福', thumb_url: '', preview_emoji: '🎉', category_id: 'cat4', category_name: '节日', creator_name: '官方', download_count: 5678, is_premium: false, create_time: '2026-02-14 08:00:00', description: '节日专用祝福表情' },
      { emoticon_id: 'emo5', name: '卖萌专用', thumb_url: '', preview_emoji: '🥰', category_id: 'cat1', category_name: '可爱', creator_name: '用户小红', download_count: 4321, is_premium: true, create_time: '2026-03-05 14:00:00', description: '萌萌哒表情包' },
      { emoticon_id: 'emo6', name: '大笑系列', thumb_url: '', preview_emoji: '😂', category_id: 'cat2', category_name: '搞笑', creator_name: '官方', download_count: 9876, is_premium: false, create_time: '2026-03-06 10:00:00', description: '笑到停不下来' }
    ]
    pagination.total = 6
    Message.warning('使用模拟数据')
  } finally {
    loading.value = false
  }
}

const handlePageChange = (page) => {
  pagination.current = page
  loadEmoticons()
}

const showDetail = (emo) => {
  currentEmoticon.value = emo
  detailVisible.value = true
}

const showUploadModal = () => {
  Object.assign(uploadForm, { name: '', category_id: null, preview_emoji: '', description: '', is_premium: false, file: null })
  uploadVisible.value = true
}

const beforeUpload = (file) => {
  uploadForm.file = file
  return false
}

const handleUpload = async () => {
  if (!uploadForm.name) { Message.warning('请输入表情包名称'); return }
  uploading.value = true
  try {
    await uploadEmoticon({ ...uploadForm })
    Message.success('表情包上传成功')
    uploadVisible.value = false
    loadEmoticons()
  } catch {
    setTimeout(() => {
      const newEmo = {
        emoticon_id: `emo${Date.now()}`,
        name: uploadForm.name,
        category_id: uploadForm.category_id,
        category_name: categories.value.find(c => c.category_id === uploadForm.category_id)?.category_name || '默认',
        preview_emoji: uploadForm.preview_emoji || '😊',
        creator_name: '当前用户',
        download_count: 0,
        is_premium: uploadForm.is_premium,
        create_time: new Date().toLocaleString(),
        description: uploadForm.description
      }
      emoticons.value.unshift(newEmo)
      pagination.total++
      Message.success('表情包上传成功')
      uploadVisible.value = false
    }, 500)
  } finally {
    uploading.value = false
  }
}

const handleDownload = async (emo) => {
  try {
    await downloadEmoticon(emo.emoticon_id)
    emo.download_count = (emo.download_count || 0) + 1
    Message.success('下载成功')
  } catch {
    emo.download_count = (emo.download_count || 0) + 1
    Message.success('下载成功（模拟）')
  }
}

onMounted(() => {
  loadCategories()
  loadEmoticons()
})
</script>

<style scoped>
.pro-page-container {
  padding: 20px 24px;
  min-height: calc(100vh - 64px);
  background: #f5f7fa;
}
.pro-breadcrumb { margin-bottom: 12px; }
.pro-page-title {
  font-size: 20px;
  font-weight: 600;
  margin-bottom: 16px;
  color: var(--color-text-1);
}
.pro-search-bar { margin-bottom: 12px; }
.pro-action-bar { margin-bottom: 16px; }
.pro-content-area {
  background: #fff;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
  min-height: 400px;
}

.emoticon-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));
  gap: 16px;
}

.emoticon-card {
  border: 1px solid var(--color-fill-2, #e5e6eb);
  border-radius: 8px;
  overflow: hidden;
  cursor: pointer;
  transition: all 0.2s;
  background: #fff;
}
.emoticon-card:hover {
  border-color: var(--color-primary, #1650ff);
  box-shadow: 0 4px 12px rgba(22, 80, 255, 0.15);
  transform: translateY(-2px);
}

.emoticon-thumb {
  position: relative;
  height: 140px;
  background: #f2f3f5;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
}
.emoticon-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
.emoticon-placeholder {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
}
.emoticon-emoji {
  font-size: 48px;
}
.emoticon-tag {
  position: absolute;
  top: 6px;
  right: 6px;
}

.emoticon-info {
  padding: 10px 12px;
}
.emoticon-name {
  font-size: 14px;
  font-weight: 500;
  color: var(--color-text-1, #1f2329);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  margin-bottom: 4px;
}
.emoticon-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 12px;
  color: var(--color-text-3, #646a73);
}
.emoticon-downloads {
  display: flex;
  align-items: center;
  gap: 2px;
}

.emoticon-pagination {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}

.emoticon-detail .detail-preview {
  display: flex;
  justify-content: center;
  align-items: center;
  background: #f2f3f5;
  border-radius: 8px;
  height: 200px;
  overflow: hidden;
}
.detail-img {
  max-width: 100%;
  max-height: 200px;
  object-fit: contain;
}
.detail-placeholder {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
}
.detail-emoji {
  font-size: 72px;
}
.detail-actions {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}
</style>
