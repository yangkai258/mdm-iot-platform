<template>
  <div class="page-container">
    <!-- 筛选栏 -->
    <div class="search-bar">
      <a-space wrap>
        <a-select
          v-model="selectedPetId"
          placeholder="选择宠物"
          style="width: 200px"
          allow-search
          @change="loadMoments"
        >
          <a-option v-for="pet in petList" :key="pet.device_id" :value="pet.device_id">
            {{ pet.pet_name }} ({{ pet.device_id }})
          </a-option>
        </a-select>
        <a-select v-model="filterCategory" placeholder="分类" style="width: 140px" allow-clear @change="loadMoments">
          <a-option value="playing">玩耍</a-option>
          <a-option value="eating">进食</a-option>
          <a-option value="sleeping">睡眠</a-option>
          <a-option value="social">社交</a-option>
          <a-option value="cute">可爱瞬间</a-option>
          <a-option value="milestone">里程碑</a-option>
        </a-select>
        <a-range-picker v-model="dateRange" style="width: 260px" @change="loadMoments" />
        <a-button type="primary" @click="loadMoments">
          <template #icon><icon-search /></template>
          搜索
        </a-button>
        <a-button @click="resetFilters">重置</a-button>
      </a-space>
    </div>

    <!-- 操作栏 -->
    <div class="action-bar">
      <a-space>
        <a-button type="primary" @click="openDownloadModal">
          <template #icon><icon-download /></template>
          批量下载
        </a-button>
        <a-button @click="shareSelected">
          <template #icon><icon-share-external /></template>
          分享选中
        </a-button>
        <a-button @click="loadMoments">
          <template #icon><icon-refresh /></template>
          刷新
        </a-button>
      </a-space>
      <a-space>
        <span class="stat-text">共 {{ total }} 个精彩瞬间</span>
      </a-space>
    </div>

    <!-- 统计卡片 -->
    <a-row :gutter="16" class="stats-row">
      <a-col :xs="24" :sm="8">
        <a-card class="stat-card">
          <a-statistic title="本月精彩瞬间" :value="monthStats.total" :value-from="0" :duration="1000" />
        </a-card>
      </a-col>
      <a-col :xs="24" :sm="8">
        <a-card class="stat-card">
          <a-statistic title="已下载" :value="monthStats.downloaded" color="green" />
        </a-card>
      </a-col>
      <a-col :xs="24" :sm="8">
        <a-card class="stat-card">
          <a-statistic title="已分享" :value="monthStats.shared" color="arcoblue" />
        </a-card>
      </a-col>
    </a-row>

    <!-- 精彩瞬间网格 -->
    <div class="moments-grid">
      <a-row :gutter="[16, 16]">
        <a-col :xs="24" :sm="12" :md="8" :lg="6" v-for="moment in moments" :key="moment.id">
          <a-card class="moment-card" :class="{ selected: selectedIds.includes(moment.id) }">
            <template #title>
              <div class="moment-header">
                <a-tag :color="getCategoryColor(moment.category)" size="small">
                  {{ getCategoryName(moment.category) }}
                </a-tag>
                <a-checkbox :model-value="selectedIds.includes(moment.id)" @change="toggleSelect(moment.id)" />
              </div>
            </template>
            <template #extra>
              <a-tag :color="moment.is_downloaded ? 'green' : 'gray'" size="small">
                {{ moment.is_downloaded ? '已下载' : '未下载' }}
              </a-tag>
            </template>
            <!-- 媒体预览 -->
            <div class="moment-media" @click="openPreview(moment)">
              <div v-if="moment.media_type === 'video'" class="media-placeholder video-placeholder">
                <icon-video-camera />
                <span class="duration">{{ moment.duration }}</span>
              </div>
              <div v-else class="media-placeholder image-placeholder">
                <icon-image />
              </div>
            </div>
            <!-- 描述 -->
            <div class="moment-desc">
              <div class="moment-title">{{ moment.title }}</div>
              <div class="moment-text">{{ moment.description }}</div>
            </div>
            <!-- 元信息 -->
            <div class="moment-meta">
              <span><icon-clock-circle /> {{ moment.captured_at }}</span>
              <span><icon-thumb-up /> {{ moment.ai_score }}</span>
            </div>
            <template #actions>
              <a-button type="text" size="small" @click="downloadMoment(moment)">
                <template #icon><icon-download /></template>
                下载
              </a-button>
              <a-button type="text" size="small" @click="shareMoment(moment)">
                <template #icon><icon-share-external /></template>
                分享
              </a-button>
              <a-button type="text" size="small" @click="openPreview(moment)">
                <template #icon><icon-full-screen /></template>
                预览
              </a-button>
            </template>
          </a-card>
        </a-col>
      </a-row>
      <a-empty v-if="!moments.length && !loading" description="暂无精彩瞬间" />
    </div>

    <!-- 分页 -->
    <div class="pagination-wrapper" v-if="total > 0">
      <a-pagination
        :total="total"
        :current="page"
        :page-size="pageSize"
        show-total
        show-page-size
        :page-size-options="[12, 24, 48]"
        @page-size-change="onPageSizeChange"
        @change="onPageChange"
      />
    </div>

    <!-- 预览弹窗 -->
    <a-modal v-model:visible="previewVisible" :width="720" footer="null" @cancel="previewVisible = false">
      <template #title>{{ previewMoment?.title }}</template>
      <div class="preview-content">
        <div v-if="previewMoment?.media_type === 'video'" class="video-player">
          <icon-video-camera :size="64" />
          <span>视频播放区域</span>
        </div>
        <div v-else class="image-preview">
          <icon-image :size="64" />
          <span>图片预览区域</span>
        </div>
        <div class="preview-info">
          <a-descriptions :column="1" size="small">
            <a-descriptions-item label="分类">
              <a-tag :color="getCategoryColor(previewMoment?.category)">
                {{ getCategoryName(previewMoment?.category) }}
              </a-tag>
            </a-descriptions-item>
            <a-descriptions-item label="描述">{{ previewMoment?.description }}</a-descriptions-item>
            <a-descriptions-item label="AI评分">
              <a-rate :model-value="previewMoment?.ai_score / 20" readonly allow-half />
            </a-descriptions-item>
            <a-descriptions-item label="捕获时间">{{ previewMoment?.captured_at }}</a-descriptions-item>
          </a-descriptions>
        </div>
      </div>
    </a-modal>

    <!-- 分享弹窗 -->
    <a-modal v-model:visible="shareVisible" title="分享" @before-ok="handleShare">
      <a-form :model="shareForm" layout="vertical">
        <a-form-item label="分享方式">
          <a-radio-group v-model="shareForm.type">
            <a-radio value="link">链接</a-radio>
            <a-radio value="image">图片</a-radio>
            <a-radio value="video">视频</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item label="有效期">
          <a-select v-model="shareForm.expire">
            <a-option value="1h">1小时</a-option>
            <a-option value="24h">24小时</a-option>
            <a-option value="7d">7天</a-option>
            <a-option value="never">永久</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="访问密码">
          <a-input v-model="shareForm.password" placeholder="可选，设置访问密码" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const petList = ref<any[]>([])
const selectedPetId = ref('')
const filterCategory = ref('')
const dateRange = ref<any[]>([])
const moments = ref<any[]>([])
const loading = ref(false)
const total = ref(0)
const page = ref(1)
const pageSize = ref(12)
const selectedIds = ref<string[]>([])
const previewVisible = ref(false)
const shareVisible = ref(false)
const previewMoment = ref<any>(null)

const monthStats = ref({ total: 0, downloaded: 0, shared: 0 })

const shareForm = ref({ type: 'link', expire: '24h', password: '' })

const categoryMap: Record<string, string> = {
  playing: '玩耍',
  eating: '进食',
  sleeping: '睡眠',
  social: '社交',
  cute: '可爱瞬间',
  milestone: '里程碑',
}

const categoryColorMap: Record<string, string> = {
  playing: 'orange',
  eating: 'green',
  sleeping: 'purple',
  social: 'blue',
  cute: 'pink',
  milestone: 'gold',
}

const getCategoryName = (cat: string) => categoryMap[cat] || cat
const getCategoryColor = (cat: string) => categoryColorMap[cat] || 'gray'

const isSelected = (id: string) => selectedIds.value.includes(id)

const toggleSelect = (id: string) => {
  const idx = selectedIds.value.indexOf(id)
  if (idx >= 0) selectedIds.value.splice(idx, 1)
  else selectedIds.value.push(id)
}

const loadPets = async () => {
  petList.value = [
    { device_id: 'pet-001', pet_name: '小白' },
    { device_id: 'pet-002', pet_name: '旺财' },
  ]
  if (petList.value.length) selectedPetId.value = petList.value[0].device_id
}

const loadMoments = () => {
  loading.value = true
  setTimeout(() => {
    moments.value = [
      {
        id: 'm1',
        title: '第一次跳跃',
        description: '小旺成功跳上了窗台',
        category: 'milestone',
        media_type: 'video',
        duration: '0:32',
        captured_at: '2026-03-23 14:30',
        ai_score: 92,
        is_downloaded: true,
      },
      {
        id: 'm2',
        title: '午睡时光',
        description: '蜷成一团睡觉的可爱模样',
        category: 'sleeping',
        media_type: 'image',
        captured_at: '2026-03-23 13:15',
        ai_score: 88,
        is_downloaded: false,
      },
      {
        id: 'm3',
        title: '追逐尾巴',
        description: '追着自己的尾巴转圈圈',
        category: 'playing',
        media_type: 'video',
        duration: '1:05',
        captured_at: '2026-03-22 16:20',
        ai_score: 85,
        is_downloaded: true,
      },
      {
        id: 'm4',
        title: '首次见面',
        description: '和邻居猫咪初次打招呼',
        category: 'social',
        media_type: 'video',
        duration: '2:10',
        captured_at: '2026-03-22 10:05',
        ai_score: 79,
        is_downloaded: false,
      },
    ]
    monthStats.value = { total: 42, downloaded: 28, shared: 15 }
    total.value = 4
    loading.value = false
  }, 600)
}

const resetFilters = () => {
  filterCategory.value = ''
  dateRange.value = []
  if (petList.value.length) selectedPetId.value = petList.value[0].device_id
  loadMoments()
}

const openPreview = (moment: any) => {
  previewMoment.value = moment
  previewVisible.value = true
}

const downloadMoment = (moment: any) => {
  Message.success(`开始下载: ${moment.title}`)
}

const shareMoment = (moment: any) => {
  shareForm.value = { type: 'link', expire: '24h', password: '' }
  shareVisible.value = true
}

const handleShare = () => {
  Message.success('分享链接已复制到剪贴板')
  shareVisible.value = false
}

const shareSelected = () => {
  if (!selectedIds.value.length) {
    Message.warning('请先选择要分享的精彩瞬间')
    return
  }
  shareForm.value = { type: 'link', expire: '24h', password: '' }
  shareVisible.value = true
}

const openDownloadModal = () => {
  if (!selectedIds.value.length) {
    Message.warning('请先选择要下载的精彩瞬间')
    return
  }
  Message.success(`开始批量下载 ${selectedIds.value.length} 个文件`)
}

const onPageChange = (p: number) => {
  page.value = p
  loadMoments()
}

const onPageSizeChange = (size: number) => {
  pageSize.value = size
  page.value = 1
  loadMoments()
}

onMounted(() => {
  loadPets()
  loadMoments()
})
</script>

<style scoped lang="less">
.page-container {
  padding: 16px;
}
.search-bar {
  margin-bottom: 12px;
  padding: 12px;
  background: var(--color-bg-2);
  border-radius: 6px;
}
.action-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}
.stat-text {
  color: var(--color-text-3);
  font-size: 13px;
}
.stats-row {
  margin-bottom: 16px;
}
.stat-card {
  text-align: center;
}
.moments-grid {
  margin-bottom: 16px;
}
.moment-card {
  cursor: pointer;
  transition: all 0.3s;
  &.selected {
    border-color: rgb(var(--primary-6));
    background: rgb(var(--primary-1));
  }
}
.moment-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.moment-media {
  margin-bottom: 8px;
  .media-placeholder {
    height: 140px;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 8px;
    border-radius: 6px;
    font-size: 13px;
    color: var(--color-text-3);
  }
    .video-placeholder {
      background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
      color: white;
    }
    .image-placeholder {
      background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
      color: white;
    }
    .duration {
      font-size: 12px;
      background: rgba(0,0,0,0.4);
      padding: 2px 6px;
      border-radius: 4px;
    }
}
.moment-desc {
  margin-bottom: 8px;
  .moment-title {
    font-weight: 600;
    font-size: 14px;
    margin-bottom: 4px;
  }
  .moment-text {
    font-size: 12px;
    color: var(--color-text-3);
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
}
.moment-meta {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  color: var(--color-text-3);
  span {
    display: flex;
    align-items: center;
    gap: 4px;
  }
}
.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
}
.preview-content {
  .video-player,
  .image-preview {
    height: 360px;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 12px;
    background: var(--color-fill-1);
    border-radius: 8px;
    margin-bottom: 16px;
    color: var(--color-text-3);
  }
  .preview-info {
    padding: 0 8px;
  }
}
</style>

