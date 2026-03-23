<template>
  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>内容生态</a-breadcrumb-item>
      <a-breadcrumb-item>内容审核</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 页面标题 -->
    <div class="pro-page-title">内容审核</div>

    <!-- 待审核数量提示 -->
    <div v-if="pendingCount > 0" class="review-tip">
      <icon-exclamation-circle-fill style="color: #ff7d00; font-size: 18px;" />
      <span>当前有 <a-badge :count="pendingCount" :max-count="99" :style="{ backgroundColor: '#ff7d00' }" /> 项内容待审核</span>
    </div>

    <!-- Tab 页签 -->
    <a-tabs v-model:active-key="activeTab" class="review-tabs" @change="handleTabChange">
      <a-tab-pane key="all" title="全部审核" />
      <a-tab-pane key="plugin" title="插件审核" />
      <a-tab-pane key="emoticon" title="表情审核" />
      <a-tab-pane key="action" title="动作审核" />
      <a-tab-pane key="voice" title="声音审核" />
    </a-tabs>

    <!-- 搜索筛选区 -->
    <div class="pro-search-bar">
      <a-space wrap>
        <a-select
          v-model="filterForm.status"
          placeholder="审核状态"
          style="width: 140px"
          allow-clear
          @change="loadReviews"
        >
          <a-option value="pending">待审核</a-option>
          <a-option value="approved">已通过</a-option>
          <a-option value="rejected">已拒绝</a-option>
        </a-select>
        <a-input-search
          v-model="filterForm.keyword"
          placeholder="搜索内容名称"
          style="width: 240px"
          search-button
          @search="loadReviews"
          @change="e => !e.target.value && loadReviews()"
        />
      </a-space>
    </div>

    <!-- 操作按钮区 -->
    <div class="pro-action-bar">
      <a-space>
        <a-button @click="loadReviews">
          <template #icon><icon-refresh /></template>
          刷新
        </a-button>
        <a-button v-if="activeTab === 'all'" type="primary" @click="showHistoryModal">
          <template #icon><icon-history /></template>
          审核历史
        </a-button>
      </a-space>
    </div>

    <!-- 审核列表 -->
    <div class="pro-content-area">
      <a-spin :loading="loading" tip="加载中...">
        <div v-if="!loading">
          <!-- 空状态 -->
          <a-empty v-if="reviews.length === 0" description="暂无待审核内容" style="padding: 60px 0" />

          <!-- 审核列表表格 -->
          <a-table
            v-else
            :columns="columns"
            :data="reviews"
            :pagination="false"
            :row-key="record => record.id"
            row-class="review-row"
          >
            <template #columns>
              <a-table-column title="序号" data-index="index" :width="60">
                <template #cell="{ rowIndex }">{{ rowIndex + 1 + (pagination.current - 1) * pagination.pageSize }}</template>
              </a-table-column>
              <a-table-column title="内容名称" data-index="name">
                <template #cell="{ record }">
                  <div class="review-name-cell">
                    <span class="review-type-icon">
                      <icon-plugin v-if="record.content_type === 'plugin'" />
                      <icon-face-smile v-else-if="record.content_type === 'emoticon'" />
                      <icon-swap v-else-if="record.content_type === 'action'" />
                      <icon-music v-else-if="record.content_type === 'voice'" />
                      <icon-file v-else />
                    </span>
                    <span class="review-name">{{ record.name }}</span>
                  </div>
                </template>
              </a-table-column>
              <a-table-column title="内容类型" data-index="content_type" :width="100">
                <template #cell="{ record }">
                  <a-tag :color="getTypeColor(record.content_type)">
                    {{ getTypeName(record.content_type) }}
                  </a-tag>
                </template>
              </a-table-column>
              <a-table-column title="开发者" data-index="developer_name" :width="120" />
              <a-table-column title="提交时间" data-index="submitted_at" :width="160">
                <template #cell="{ record }">{{ record.submitted_at || '-' }}</template>
              </a-table-column>
              <a-table-column title="状态" data-index="status" :width="100">
                <template #cell="{ record }">
                  <a-tag :color="getStatusColor(record.status)">
                    {{ getStatusName(record.status) }}
                  </a-tag>
                </template>
              </a-table-column>
              <a-table-column title="操作" :width="140" align="center">
                <template #cell="{ record }">
                  <a-space v-if="record.status === 'pending'">
                    <a-button type="primary" size="small" @click="handleApprove(record)">
                      通过
                    </a-button>
                    <a-button status="danger" size="small" @click="handleReject(record)">
                      拒绝
                    </a-button>
                  </a-space>
                  <a-button v-else size="small" @click="showDetail(record)">
                    详情
                  </a-button>
                </template>
              </a-table-column>
            </template>
          </a-table>

          <!-- 分页 -->
          <div v-if="pagination.total > 0" class="review-pagination">
            <a-pagination
              :total="pagination.total"
              :current="pagination.current"
              :page-size="pagination.pageSize"
              @change="handlePageChange"
              show-total
              show-page-size
              :page-size-options="[10, 20, 50, 100]"
            />
          </div>
        </div>
      </a-spin>
    </div>

    <!-- 审核详情弹窗 -->
    <a-modal
      v-model:visible="detailVisible"
      :title="`内容详情 - ${currentReview?.name}`"
      :width="640"
      :footer="currentReview?.status === 'pending' ? modalFooter : null"
    >
      <div v-if="currentReview" class="review-detail">
        <a-descriptions :column="2" size="small" bordered>
          <a-descriptions-item label="内容名称" :span="2">
            {{ currentReview.name }}
          </a-descriptions-item>
          <a-descriptions-item label="内容类型">
            <a-tag :color="getTypeColor(currentReview.content_type)">
              {{ getTypeName(currentReview.content_type) }}
            </a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="审核状态">
            <a-tag :color="getStatusColor(currentReview.status)">
              {{ getStatusName(currentReview.status) }}
            </a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="开发者">
            {{ currentReview.developer_name }}
          </a-descriptions-item>
          <a-descriptions-item label="提交时间">
            {{ currentReview.submitted_at || '-' }}
          </a-descriptions-item>
          <a-descriptions-item label="价格">
            {{ currentReview.is_free ? '免费' : `¥${currentReview.price}` }}
          </a-descriptions-item>
          <a-descriptions-item label="下载量">
            {{ currentReview.download_count || 0 }}
          </a-descriptions-item>
          <a-descriptions-item label="描述" :span="2">
            {{ currentReview.description || '暂无描述' }}
          </a-descriptions-item>
          <a-descriptions-item label="拒绝原因" :span="2" v-if="currentReview.rejected_reason">
            <span style="color: #f53f3f;">{{ currentReview.rejected_reason }}</span>
          </a-descriptions-item>
        </a-descriptions>

        <!-- 预览区 -->
        <div v-if="currentReview.preview_url" class="detail-preview">
          <div class="preview-label">内容预览</div>
          <img
            v-if="currentReview.content_type === 'emoticon'"
            :src="currentReview.preview_url"
            class="preview-img"
          />
          <video
            v-else-if="currentReview.content_type === 'action'"
            :src="currentReview.preview_url"
            controls
            class="preview-video"
          />
          <audio
            v-else-if="currentReview.content_type === 'voice'"
            :src="currentReview.preview_url"
            controls
            class="preview-audio"
          />
          <a-image
            v-else
            :src="currentReview.preview_url"
            :preview="true"
          />
        </div>
      </div>

      <template #footer v-if="currentReview?.status === 'pending'">
        <a-space>
          <a-button @click="detailVisible = false">关闭</a-button>
          <a-button status="danger" @click="handleReject(currentReview)">拒绝</a-button>
          <a-button type="primary" @click="handleApprove(currentReview)">通过</a-button>
        </a-space>
      </template>
    </a-modal>

    <!-- 审核历史弹窗 -->
    <a-modal
      v-model:visible="historyVisible"
      title="审核历史"
      :width="800"
      :footer="null"
    >
      <a-spin :loading="historyLoading" tip="加载中...">
        <a-table
          :columns="historyColumns"
          :data="historyRecords"
          :pagination="historyPagination"
          :row-key="record => record.id"
          @page-change="handleHistoryPageChange"
        >
          <template #columns>
            <a-table-column title="内容名称" data-index="name" />
            <a-table-column title="内容类型" data-index="content_type" :width="100">
              <template #cell="{ record }">
                <a-tag :color="getTypeColor(record.content_type)">
                  {{ getTypeName(record.content_type) }}
                </a-tag>
              </template>
            </a-table-column>
            <a-table-column title="审核结果" data-index="status" :width="100">
              <template #cell="{ record }">
                <a-tag :color="getStatusColor(record.status)">
                  {{ getStatusName(record.status) }}
                </a-tag>
              </template>
            </a-table-column>
            <a-table-column title="审核时间" data-index="reviewed_at" :width="160" />
            <a-table-column title="审核人" data-index="reviewer_name" :width="100" />
          </template>
        </a-table>
      </a-spin>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import {
  getReviewList,
  getReviewPendingCount,
  getReviewDetail,
  reviewContent,
  getReviewHistory
} from '@/api/market'

// 状态
const loading = ref(false)
const historyLoading = ref(false)
const activeTab = ref('all')
const reviews = ref([])
const pendingCount = ref(0)
const historyRecords = ref([])
const currentReview = ref(null)
const detailVisible = ref(false)
const historyVisible = ref(false)

// 筛选表单
const filterForm = reactive({
  status: 'pending',
  keyword: '',
  content_type: ''
})

// 分页
const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0
})

const historyPagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0
})

// 表格列
const columns = [
  { title: '序号', dataIndex: 'index', width: 60 },
  { title: '内容名称', dataIndex: 'name' },
  { title: '内容类型', dataIndex: 'content_type', width: 100 },
  { title: '开发者', dataIndex: 'developer_name', width: 120 },
  { title: '提交时间', dataIndex: 'submitted_at', width: 160 },
  { title: '状态', dataIndex: 'status', width: 100 },
  { title: '操作', width: 140, align: 'center' }
]

const historyColumns = [
  { title: '内容名称', dataIndex: 'name' },
  { title: '内容类型', dataIndex: 'content_type', width: 100 },
  { title: '审核结果', dataIndex: 'status', width: 100 },
  { title: '审核时间', dataIndex: 'reviewed_at', width: 160 },
  { title: '审核人', dataIndex: 'reviewer_name', width: 100 }
]

// 获取类型颜色
const getTypeColor = (type) => {
  const colors = {
    plugin: 'blue',
    emoticon: 'arcoblue',
    action: 'green',
    voice: 'orange'
  }
  return colors[type] || 'gray'
}

// 获取类型名称
const getTypeName = (type) => {
  const names = {
    plugin: '插件',
    emoticon: '表情包',
    action: '动作',
    voice: '声音'
  }
  return names[type] || type
}

// 获取状态颜色
const getStatusColor = (status) => {
  const colors = {
    pending: 'gold',
    approved: 'green',
    rejected: 'red'
  }
  return colors[status] || 'gray'
}

// 获取状态名称
const getStatusName = (status) => {
  const names = {
    pending: '待审核',
    approved: '已通过',
    rejected: '已拒绝'
  }
  return names[status] || status
}

// Tab 切换
const handleTabChange = (key) => {
  filterForm.content_type = key === 'all' ? '' : key
  pagination.current = 1
  loadReviews()
}

// 加载审核列表
const loadReviews = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.current,
      page_size: pagination.pageSize,
      status: filterForm.status || undefined,
      keyword: filterForm.keyword || undefined,
      content_type: filterForm.content_type || undefined
    }
    const res = await getReviewList(params)
    reviews.value = res.data?.list || []
    pagination.total = res.data?.total || 0
  } catch {
    // 模拟数据
    reviews.value = [
      { id: 'rev1', name: '可爱猫咪表情包', content_type: 'emoticon', developer_name: '用户小明', submitted_at: '2026-03-22 10:30:00', status: 'pending', price: 0, is_free: true, download_count: 0, description: '一套可爱的猫咪表情包' },
      { id: 'rev2', name: '夜灯控制插件', content_type: 'plugin', developer_name: '开发者张三', submitted_at: '2026-03-22 09:15:00', status: 'pending', price: 5, is_free: false, download_count: 0, description: '控制设备夜灯模式的插件' },
      { id: 'rev3', name: '跳舞动作', content_type: 'action', developer_name: '用户小红', submitted_at: '2026-03-21 16:45:00', status: 'approved', price: 10, is_free: false, download_count: 234, description: '可爱的跳舞动作' },
      { id: 'rev4', name: '童声朗读', content_type: 'voice', developer_name: '用户小刚', submitted_at: '2026-03-21 14:20:00', status: 'rejected', price: 20, is_free: false, download_count: 0, description: '童声TTS声音', rejected_reason: '内容不符合规范' },
      { id: 'rev5', name: '节日烟花动作', content_type: 'action', developer_name: '官方', submitted_at: '2026-03-20 11:00:00', status: 'approved', price: 0, is_free: true, download_count: 1234, description: '节日烟花特效动作' },
      { id: 'rev6', name: '温控插件', content_type: 'plugin', developer_name: '开发者李四', submitted_at: '2026-03-22 08:00:00', status: 'pending', price: 8, is_free: false, download_count: 0, description: '温度控制插件' }
    ]
    pagination.total = reviews.value.length
    Message.warning('使用模拟数据')
  } finally {
    loading.value = false
  }
}

// 加载待审核数量
const loadPendingCount = async () => {
  try {
    const res = await getReviewPendingCount()
    pendingCount.value = res.data?.count || 0
  } catch {
    pendingCount.value = 3
  }
}

// 分页切换
const handlePageChange = (page) => {
  pagination.current = page
  loadReviews()
}

// 显示详情
const showDetail = (record) => {
  currentReview.value = record
  detailVisible.value = true
}

// 通过审核
const handleApprove = async (record) => {
  try {
    await reviewContent(record.id, { status: 'approved' })
    Message.success('审核通过')
    detailVisible.value = false
    loadReviews()
    loadPendingCount()
  } catch {
    record.status = 'approved'
    Message.success('审核通过（模拟）')
    detailVisible.value = false
  }
}

// 拒绝审核
const handleReject = async (record) => {
  try {
    await reviewContent(record.id, { status: 'rejected', reason: '内容不符合平台规范' })
    Message.success('已拒绝')
    detailVisible.value = false
    loadReviews()
    loadPendingCount()
  } catch {
    record.status = 'rejected'
    Message.success('已拒绝（模拟）')
    detailVisible.value = false
  }
}

// 显示审核历史
const showHistoryModal = () => {
  historyVisible.value = true
  loadHistory()
}

// 加载审核历史
const loadHistory = async () => {
  historyLoading.value = true
  try {
    const params = {
      page: historyPagination.current,
      page_size: historyPagination.pageSize
    }
    const res = await getReviewHistory(params)
    historyRecords.value = res.data?.list || []
    historyPagination.total = res.data?.total || 0
  } catch {
    historyRecords.value = [
      { id: 'h1', name: '跳舞动作', content_type: 'action', status: 'approved', reviewed_at: '2026-03-21 17:00:00', reviewer_name: '管理员' },
      { id: 'h2', name: '童声朗读', content_type: 'voice', status: 'rejected', reviewed_at: '2026-03-21 15:00:00', reviewer_name: '管理员' },
      { id: 'h3', name: '节日烟花动作', content_type: 'action', status: 'approved', reviewed_at: '2026-03-20 12:00:00', reviewer_name: '管理员' }
    ]
    historyPagination.total = historyRecords.value.length
    Message.warning('使用模拟数据')
  } finally {
    historyLoading.value = false
  }
}

// 审核历史分页
const handleHistoryPageChange = (page) => {
  historyPagination.current = page
  loadHistory()
}

onMounted(() => {
  loadPendingCount()
  loadReviews()
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
.review-tip {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  background: #fffbe6;
  border: 1px solid #ffe58f;
  border-radius: 6px;
  margin-bottom: 16px;
  font-size: 14px;
  color: #ad6800;
}
.review-tabs {
  background: #fff;
  border-radius: 8px 8px 0 0;
  padding: 0 16px;
}
.pro-search-bar { margin-bottom: 12px; }
.pro-action-bar { margin-bottom: 16px; }
.pro-content-area {
  background: #fff;
  border-radius: 0 0 8px 8px;
  padding: 20px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
  min-height: 400px;
}

.review-name-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}
.review-type-icon {
  font-size: 18px;
  color: var(--color-text-3);
}
.review-name {
  font-weight: 500;
}

.review-pagination {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}

.review-detail {
  padding: 8px 0;
}
.detail-preview {
  margin-top: 16px;
}
.preview-label {
  font-size: 14px;
  font-weight: 500;
  color: var(--color-text-1);
  margin-bottom: 8px;
}
.preview-img {
  max-width: 100%;
  border-radius: 8px;
}
.preview-video {
  width: 100%;
  max-height: 300px;
  border-radius: 8px;
}
.preview-audio {
  width: 100%;
}

:deep(.review-row:hover) {
  background: var(--color-fill-1);
}
</style>
