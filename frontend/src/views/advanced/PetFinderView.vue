<template>
  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>高级功能</a-breadcrumb-item>
      <a-breadcrumb-item>寻回网络</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 页面标题 -->
    <div class="pro-page-header">
      <h2 class="pro-page-title">寻回网络</h2>
      <p class="pro-page-desc">发布宠物走失报告，全社区协助寻回</p>
    </div>

    <!-- 搜索筛选区 -->
    <div class="pro-search-bar">
      <a-space>
        <a-select v-model="statusFilter" placeholder="报告状态" allow-clear style="width: 140px" @change="loadReports">
          <a-option value="active">寻找中</a-option>
          <a-option value="found">已找到</a-option>
          <a-option value="closed">已关闭</a-option>
        </a-select>
        <a-select v-model="petTypeFilter" placeholder="宠物类型" allow-clear style="width: 140px" @change="loadReports">
          <a-option value="dog">狗</a-option>
          <a-option value="cat">猫</a-option>
          <a-option value="other">其他</a-option>
        </a-select>
        <a-input-search v-model="keyword" placeholder="搜索宠物名称" style="width: 200px" @search="loadReports" search-button />
      </a-space>
    </div>

    <!-- 操作按钮区 -->
    <div class="pro-action-bar">
      <a-space>
        <a-button type="primary" @click="showCreateModal">
          <template #icon><icon-add /></template>
          发布寻宠报告
        </a-button>
        <a-button @click="loadReports">
          <template #icon><icon-refresh /></template>
          刷新
        </a-button>
      </a-space>
    </div>

    <!-- 寻宠报告列表 -->
    <div class="pro-content-area">
      <a-table
        :columns="columns"
        :data="reports"
        :loading="loading"
        :pagination="pagination"
        @page-change="onPageChange"
        row-key="id"
      >
        <template #pet_photo="{ record }">
          <img
            v-if="record.pet_photo"
            :src="record.pet_photo"
            :alt="record.pet_name"
            style="width:48px;height:48px;border-radius:50%;object-fit:cover"
            @error="e => (e.target as any).src = 'data:image/svg+xml,<svg xmlns=%22http://www.w3.org/2000/svg%22 width=%2248%22 height=%2248%22><circle fill=%22%23f2f3f5%22 cx=%2224%22 cy=%2224%22 r=%2224%22/><text x=%2224%22 y=%2230%22 text-anchor=%22middle%22 fill=%22%23888%22 font-size=%2212%22>?</text></svg>'"
          />
          <a-avatar v-else :size="48" :style="{ backgroundColor: '#1659f5' }">
            {{ record.pet_name?.charAt(0) || '?' }}
          </a-avatar>
        </template>
        <template #pet_name="{ record }">
          <div>
            <div style="font-weight:600">{{ record.pet_name }}</div>
            <div style="font-size:12px;color:#666">{{ record.pet_breed || record.pet_type }}</div>
          </div>
        </template>
        <template #status="{ record }">
          <a-tag :color="getStatusColor(record.status)">
            {{ getStatusLabel(record.status) }}
          </a-tag>
        </template>
        <template #last_seen_location="{ record }">
          <span>{{ record.last_seen_location || '-' }}</span>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="showDetailModal(record)">详情</a-button>
            <a-button type="text" size="small" @click="showOnMap(record)">地图</a-button>
            <a-button type="text" size="small" status="success" @click="markAsFound(record)" v-if="record.status === 'active'">已找到</a-button>
            <a-button type="text" size="small" @click="shareReport(record)">分享</a-button>
          </a-space>
        </template>
      </a-table>
    </div>

    <!-- 地图抽屉 -->
    <a-drawer v-model:visible="mapDrawerVisible" :width="800" :title="`寻宠地图 - ${mapPetName}`" @close="mapDrawerVisible = false">
      <div class="map-container" ref="mapContainerRef">
        <div class="map-placeholder">
          <div class="map-placeholder-content">
            <icon-location :size="48" style="color:#1659f5;margin-bottom:12px" />
            <p style="font-size:14px;color:#666;margin-bottom:8px">走失地点：{{ mapReport?.last_seen_location || '未知' }}</p>
            <p style="font-size:13px;color:#999;margin-bottom:16px">最后出现时间：{{ formatDate(mapReport?.last_seen_time) }}</p>
            <a-tag v-if="mapReport" :color="getStatusColor(mapReport.status)">
              {{ getStatusLabel(mapReport.status) }}
            </a-tag>
          </div>
          <!-- 地图标记点（模拟） -->
          <div class="map-markers" v-if="mapReport">
            <div
              class="map-marker map-marker-main"
              :style="mapMarkerStyle"
              @click="mapInfoVisible = !mapInfoVisible"
            >
              <icon-location :size="32" style="color:#f53f3f" />
            </div>
            <!-- 扩散圈 -->
            <div class="map-search-circle"></div>
          </div>
        </div>
      </div>

      <!-- 地图信息面板 -->
      <div v-if="mapInfoVisible && mapReport" class="map-info-panel">
        <a-descriptions :column="1" size="small">
          <a-descriptions-item label="宠物名称">{{ mapReport.pet_name }}</a-descriptions-item>
          <a-descriptions-item label="宠物类型">{{ mapReport.pet_type }}</a-descriptions-item>
          <a-descriptions-item label="走失时间">{{ formatDate(mapReport.last_seen_time) }}</a-descriptions-item>
          <a-descriptions-item label="走失地点">{{ mapReport.last_seen_location }}</a-descriptions-item>
          <a-descriptions-item label="联系方式">{{ mapReport.contact_phone }}</a-descriptions-item>
          <a-descriptions-item label="悬赏金额">{{ mapReport.reward || '暂无' }}</a-descriptions-item>
        </a-descriptions>
        <div style="margin-top:12px">
          <a-button type="primary" size="small" @click="showDetailModal(mapReport)">查看详情</a-button>
        </div>
      </div>
    </a-drawer>

    <!-- 创建/编辑报告弹窗 -->
    <a-modal
      v-model:visible="createModalVisible"
      :title="isEdit ? '编辑寻宠报告' : '发布寻宠报告'"
      @ok="handleSave"
      :width="600"
      @close="resetForm"
    >
      <a-form :model="form" layout="vertical">
        <a-form-item label="宠物名称" required>
          <a-input v-model="form.pet_name" placeholder="请输入宠物名称" />
        </a-form-item>

        <a-form-item label="宠物类型" required>
          <a-radio-group v-model="form.pet_type">
            <a-radio value="dog">狗</a-radio>
            <a-radio value="cat">猫</a-radio>
            <a-radio value="other">其他</a-radio>
          </a-radio-group>
        </a-form-item>

        <a-form-item label="宠物品种">
          <a-input v-model="form.pet_breed" placeholder="如：金毛、柯基、布偶猫等" />
        </a-form-item>

        <a-form-item label="宠物照片">
          <div class="upload-photo-area" @click="triggerPhotoUpload">
            <img v-if="form.pet_photo_url" :src="form.pet_photo_url" class="photo-preview" />
            <div v-else class="upload-placeholder">
              <icon-upload :size="32" style="color:#999" />
              <span style="color:#999;font-size:13px;margin-top:8px">点击上传照片</span>
            </div>
          </div>
          <input ref="photoInputRef" type="file" accept="image/*" style="display:none" @change="handlePhotoChange" />
        </a-form-item>

        <a-form-item label="走失时间" required>
          <a-date-picker
            v-model="form.last_seen_time"
            show-time
            format="YYYY-MM-DD HH:mm"
            style="width: 100%"
          />
        </a-form-item>

        <a-form-item label="走失地点" required>
          <a-input v-model="form.last_seen_location" placeholder="请输入走失地点" />
        </a-form-item>

        <a-form-item label="详细描述">
          <a-textarea v-model="form.description" placeholder="请描述宠物特征、颜色、是否有项圈等" :rows="4" />
        </a-form-item>

        <a-form-item label="联系方式" required>
          <a-input v-model="form.contact_phone" placeholder="请输入联系电话" />
        </a-form-item>

        <a-form-item label="悬赏金额（可选）">
          <a-input-number v-model="form.reward" :min="0" :precision="0" placeholder="输入悬赏金额" style="width: 200px" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 报告详情弹窗 -->
    <a-modal v-model:visible="detailModalVisible" :title="`寻宠报告 - ${detailReport?.pet_name}`" :width="640">
      <div v-if="detailReport" class="report-detail">
        <div class="report-detail-header">
          <img
            v-if="detailReport.pet_photo"
            :src="detailReport.pet_photo"
            :alt="detailReport.pet_name"
            class="report-detail-photo"
          />
          <div class="report-detail-info">
            <h3>{{ detailReport.pet_name }}</h3>
            <p>{{ detailReport.pet_breed || detailReport.pet_type }}</p>
            <a-tag :color="getStatusColor(detailReport.status)">{{ getStatusLabel(detailReport.status) }}</a-tag>
          </div>
        </div>

        <a-descriptions :column="1" size="small" style="margin-top:16px">
          <a-descriptions-item label="走失时间">{{ formatDate(detailReport.last_seen_time) }}</a-descriptions-item>
          <a-descriptions-item label="走失地点">{{ detailReport.last_seen_location }}</a-descriptions-item>
          <a-descriptions-item label="联系方式">{{ detailReport.contact_phone }}</a-descriptions-item>
          <a-descriptions-item label="悬赏金额">{{ detailReport.reward ? `${detailReport.reward} 元` : '暂无' }}</a-descriptions-item>
          <a-descriptions-item label="详细描述">{{ detailReport.description || '暂无' }}</a-descriptions-item>
          <a-descriptions-item label="发布时间">{{ formatDate(detailReport.created_at) }}</a-descriptions-item>
        </a-descriptions>
      </div>
    </a-modal>

    <!-- 分享弹窗 -->
    <a-modal v-model:visible="shareModalVisible" title="分享寻宠报告" @ok="handleShareReport" :width="420">
      <p style="margin-bottom:12px;color:#666">生成分享链接，让更多人看到这条寻宠信息</p>
      <a-input v-if="shareUrl" :model-value="shareUrl" readonly />
      <a-spin v-else />
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const reports = ref<any[]>([])
const loading = ref(false)
const statusFilter = ref('')
const petTypeFilter = ref('')
const keyword = ref('')
const createModalVisible = ref(false)
const detailModalVisible = ref(false)
const mapDrawerVisible = ref(false)
const shareModalVisible = ref(false)
const isEdit = ref(false)
const mapReport = ref<any>(null)
const mapPetName = ref('')
const mapInfoVisible = ref(false)
const shareUrl = ref('')
const detailReport = ref<any>(null)
const photoInputRef = ref<any>(null)

const mapMarkerStyle = reactive({
  top: '50%',
  left: '50%'
})

const form = reactive({
  id: null as number | null,
  pet_name: '',
  pet_type: 'dog',
  pet_breed: '',
  pet_photo: null as File | null,
  pet_photo_url: '',
  last_seen_time: null as any,
  last_seen_location: '',
  description: '',
  contact_phone: '',
  reward: null as number | null
})

const columns = [
  { title: '照片', slotName: 'pet_photo', width: 80 },
  { title: '宠物信息', slotName: 'pet_name', width: 160 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '走失地点', slotName: 'last_seen_location', width: 200 },
  { title: '发布时间', dataIndex: 'created_at', width: 160 },
  { title: '操作', slotName: 'actions', width: 260 }
]

const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0
})

function getStatusColor(status: string) {
  const map: Record<string, string> = {
    active: 'red',
    found: 'green',
    closed: 'gray'
  }
  return map[status] || 'default'
}

function getStatusLabel(status: string) {
  const map: Record<string, string> = {
    active: '寻找中',
    found: '已找到',
    closed: '已关闭'
  }
  return map[status] || status
}

function formatDate(date: string) {
  if (!date) return '-'
  return new Date(date).toLocaleString('zh-CN')
}

async function loadReports() {
  loading.value = true
  try {
    const params = new URLSearchParams()
    if (statusFilter.value) params.append('status', statusFilter.value)
    if (petTypeFilter.value) params.append('pet_type', petTypeFilter.value)
    if (keyword.value) params.append('keyword', keyword.value)
    params.append('page', String(pagination.current))
    params.append('page_size', String(pagination.pageSize))

    const res = await fetch(`/api/v1/advanced/pet-finder/reports?${params}`, { credentials: 'include' })
    const data = await res.json()
    if (data.code === 0 || data.code === 200) {
      reports.value = data.data?.list || data.data || []
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
  loadReports()
}

function showCreateModal() {
  isEdit.value = false
  resetForm()
  createModalVisible.value = true
}

function showDetailModal(report: any) {
  detailReport.value = report
  detailModalVisible.value = true
}

function showOnMap(report: any) {
  mapReport.value = report
  mapPetName.value = report.pet_name
  mapInfoVisible.value = false
  mapDrawerVisible.value = true
}

function showCreateModal(record?: any) {
  if (record) {
    isEdit.value = true
    form.id = record.id
    form.pet_name = record.pet_name
    form.pet_type = record.pet_type
    form.pet_breed = record.pet_breed || ''
    form.pet_photo_url = record.pet_photo || ''
    form.last_seen_time = record.last_seen_time
    form.last_seen_location = record.last_seen_location || ''
    form.description = record.description || ''
    form.contact_phone = record.contact_phone || ''
    form.reward = record.reward || null
  } else {
    isEdit.value = false
    resetForm()
  }
  createModalVisible.value = true
}

function resetForm() {
  form.id = null
  form.pet_name = ''
  form.pet_type = 'dog'
  form.pet_breed = ''
  form.pet_photo = null
  form.pet_photo_url = ''
  form.last_seen_time = null
  form.last_seen_location = ''
  form.description = ''
  form.contact_phone = ''
  form.reward = null
}

function triggerPhotoUpload() {
  photoInputRef.value?.click()
}

function handlePhotoChange(e: any) {
  const file = e.target.files?.[0]
  if (!file) return
  form.pet_photo = file
  form.pet_photo_url = URL.createObjectURL(file)
}

async function handleSave() {
  if (!form.pet_name || !form.last_seen_location || !form.contact_phone) {
    Message.warning('请填写必填项')
    return
  }
  try {
    const payload: any = { ...form }
    if (payload.last_seen_time instanceof Date) {
      payload.last_seen_time = payload.last_seen_time.toISOString()
    }
    delete payload.pet_photo

    const url = isEdit.value
      ? `/api/v1/advanced/pet-finder/reports/${form.id}`
      : '/api/v1/advanced/pet-finder/reports'
    const method = isEdit.value ? 'PUT' : 'POST'

    const res = await fetch(url, {
      method,
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify(payload)
    })
    const data = await res.json()
    if (data.code === 0 || data.code === 200) {
      Message.success(isEdit.value ? '报告已更新' : '报告已发布')
      createModalVisible.value = false
      loadReports()
    } else {
      Message.error(data.message || '操作失败')
    }
  } catch {
    Message.error('网络错误')
  }
}

async function markAsFound(report: any) {
  try {
    const res = await fetch(`/api/v1/advanced/pet-finder/reports/${report.id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({ status: 'found' })
    })
    const data = await res.json()
    if (data.code === 0 || data.code === 200) {
      Message.success('已标记为找到！')
      loadReports()
    } else {
      Message.error(data.message || '操作失败')
    }
  } catch {
    Message.error('网络错误')
  }
}

function shareReport(report: any) {
  shareUrl.value = ''
  detailReport.value = report
  shareModalVisible.value = true
  generateShareUrl(report)
}

async function generateShareUrl(report: any) {
  try {
    const res = await fetch('/api/v1/advanced/pet-finder/reports/share', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({ report_id: report.id })
    })
    const data = await res.json()
    if (data.code === 0 || data.code === 200) {
      shareUrl.value = data.data?.url || ''
    }
  } catch { /* ignore */ }
}

async function handleShareReport() {
  shareModalVisible.value = false
}

onMounted(() => {
  loadReports()
})
</script>

<style scoped>
.map-container {
  width: 100%;
  height: 400px;
  background: #f2f3f5;
  border-radius: 8px;
  overflow: hidden;
  position: relative;
}

.map-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  position: relative;
}

.map-placeholder-content {
  text-align: center;
  z-index: 2;
}

.map-markers {
  position: absolute;
  inset: 0;
}

.map-marker {
  position: absolute;
  transform: translate(-50%, -100%);
  cursor: pointer;
  z-index: 3;
  filter: drop-shadow(0 2px 4px rgba(0,0,0,0.2));
}

.map-marker-main {
  top: 40%;
  left: 45%;
}

.map-search-circle {
  position: absolute;
  top: 40%;
  left: 45%;
  width: 120px;
  height: 120px;
  border: 2px dashed rgba(245, 63, 63, 0.4);
  border-radius: 50%;
  transform: translate(-50%, -100%);
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 0.6; transform: translate(-50%, -100%) scale(1); }
  50% { opacity: 0.3; transform: translate(-50%, -100%) scale(1.2); }
}

.map-info-panel {
  position: absolute;
  bottom: 16px;
  left: 16px;
  background: #fff;
  border-radius: 8px;
  padding: 16px;
  box-shadow: 0 4px 16px rgba(0,0,0,0.15);
  z-index: 10;
  min-width: 240px;
}

.report-detail-header {
  display: flex;
  gap: 16px;
  align-items: center;
}

.report-detail-photo {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  object-fit: cover;
}

.report-detail-info h3 {
  margin: 0 0 4px 0;
}

.upload-photo-area {
  width: 120px;
  height: 120px;
  border: 2px dashed var(--color-border);
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  overflow: hidden;
  transition: border-color 0.2s;
}
.upload-photo-area:hover {
  border-color: var(--color-primary);
}

.photo-preview {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.upload-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
}
</style>
