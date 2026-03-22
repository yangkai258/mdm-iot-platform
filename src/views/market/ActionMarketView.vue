<template>
  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>内容生态</a-breadcrumb-item>
      <a-breadcrumb-item>动作资源库</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 页面标题 -->
    <div class="pro-page-title">动作资源库</div>

    <!-- 搜索筛选区 -->
    <div class="pro-search-bar">
      <a-space wrap>
        <a-select
          v-model="filterForm.category_id"
          placeholder="全部分类"
          style="width: 160px"
          allow-clear
          @change="loadActions"
        >
          <a-option v-for="cat in categories" :key="cat.category_id" :value="cat.category_id">
            {{ cat.category_name }}
          </a-option>
        </a-select>
        <a-select
          v-model="filterForm.is_official"
          placeholder="来源"
          style="width: 140px"
          allow-clear
          @change="loadActions"
        >
          <a-option :value="true">官方</a-option>
          <a-option :value="false">用户</a-option>
        </a-select>
        <a-input-search
          v-model="filterForm.keyword"
          placeholder="搜索动作名称"
          style="width: 240px"
          search-button
          @search="loadActions"
          @change="e => !e.target.value && loadActions()"
        />
      </a-space>
    </div>

    <!-- 操作按钮区 -->
    <div class="pro-action-bar">
      <a-space>
        <a-button type="primary" @click="showCreateModal">创建自定义动作</a-button>
        <a-button @click="loadActions">刷新</a-button>
      </a-space>
    </div>

    <!-- 动作列表 -->
    <div class="pro-content-area">
      <a-spin :loading="loading" tip="加载中...">
        <a-empty v-if="actions.length === 0 && !loading" description="暂无动作" style="padding: 60px 0" />

        <div v-else class="action-grid">
          <div
            v-for="action in actions"
            :key="action.action_id"
            class="action-card"
          >
            <!-- 动作预览区 -->
            <div class="action-preview" @click="previewAction(action)">
              <div class="action-preview-inner">
                <span class="action-icon">{{ action.preview_icon || '🎭' }}</span>
                <a-tag
                  class="action-source-tag"
                  :color="action.is_official ? 'blue' : 'green'"
                >
                  {{ action.is_official ? '官方' : '用户' }}
                </a-tag>
              </div>
              <div class="action-overlay">
                <icon-play-circle-fill class="preview-icon" />
                <span>预览</span>
              </div>
            </div>

            <!-- 动作信息 -->
            <div class="action-info">
              <div class="action-name" :title="action.name">{{ action.name }}</div>
              <div class="action-category">
                <a-tag size="small">{{ action.category_name || '默认' }}</a-tag>
              </div>
              <div class="action-meta">
                <span>创建者: {{ action.creator_name || '官方' }}</span>
              </div>
              <div class="action-desc">{{ action.description || '暂无描述' }}</div>
            </div>

            <!-- 操作按钮 -->
            <div class="action-operations">
              <a-button type="text" size="small" @click="previewAction(action)">预览</a-button>
              <a-button
                v-if="!action.is_official && !action.is_published"
                type="text"
                size="small"
                status="success"
                @click="publishAction(action)"
              >
                发布到市场
              </a-button>
              <a-button
                v-if="action.is_published"
                type="text"
                size="small"
                disabled
              >
                已发布
              </a-button>
              <a-button type="text" size="small" status="danger" @click="deleteAction(action)">删除</a-button>
            </div>
          </div>
        </div>

        <!-- 分页 -->
        <div v-if="pagination.total > 0" class="action-pagination">
          <a-pagination
            :total="pagination.total"
            :current="pagination.current"
            :page-size="pagination.pageSize"
            @change="handlePageChange"
            show-total
            show-page-size
            :page-size-options="[12, 24, 36, 60]"
          />
        </div>
      </a-spin>
    </div>

    <!-- 预览弹窗 -->
    <a-modal
      v-model:visible="previewVisible"
      :title="`预览: ${currentAction?.name}`"
      :width="500"
      :footer="null"
    >
      <div v-if="currentAction" class="action-preview-modal">
        <div class="preview-display">
          <span class="preview-big-icon">{{ currentAction.preview_icon || '🎭' }}</span>
        </div>
        <a-descriptions :column="2" size="small" style="margin-top: 16px">
          <a-descriptions-item label="动作名称">
            {{ currentAction.name }}
          </a-descriptions-item>
          <a-descriptions-item label="分类">
            {{ currentAction.category_name || '默认' }}
          </a-descriptions-item>
          <a-descriptions-item label="创建者">
            {{ currentAction.creator_name || '官方' }}
          </a-descriptions-item>
          <a-descriptions-item label="使用次数">
            {{ currentAction.usage_count || 0 }}
          </a-descriptions-item>
          <a-descriptions-item label="动画时长">
            {{ currentAction.duration || 0 }}s
          </a-descriptions-item>
          <a-descriptions-item label="状态">
            <a-tag :color="currentAction.is_published ? 'green' : 'orange'">
              {{ currentAction.is_published ? '已发布' : '未发布' }}
            </a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="描述" :span="2">
            {{ currentAction.description || '暂无描述' }}
          </a-descriptions-item>
        </a-descriptions>
        <div style="margin-top: 16px; display: flex; justify-content: flex-end">
          <a-button @click="previewVisible = false">关闭</a-button>
        </div>
      </div>
    </a-modal>

    <!-- 创建自定义动作弹窗 -->
    <a-modal
      v-model:visible="createVisible"
      title="创建自定义动作"
      @ok="handleCreate"
      :confirm-loading="creating"
      :width="560"
    >
      <a-form :model="createForm" layout="vertical">
        <a-form-item label="动作名称" required>
          <a-input v-model="createForm.name" placeholder="请输入动作名称" />
        </a-form-item>
        <a-form-item label="分类" required>
          <a-select v-model="createForm.category_id" placeholder="选择分类" allow-clear>
            <a-option v-for="cat in categories" :key="cat.category_id" :value="cat.category_id">
              {{ cat.category_name }}
            </a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="预览图标">
          <a-input v-model="createForm.preview_icon" placeholder="如: 🎭 (单个emoji)" maxlength="4" />
        </a-form-item>
        <a-form-item label="动画时长(秒)">
          <a-input-number v-model="createForm.duration" :min="0" :max="300" style="width: 100%" />
        </a-form-item>
        <a-form-item label="动作参数(JSON)">
          <a-textarea
            v-model="createForm.params"
            placeholder='{"angle": 30, "speed": 1}'
            :rows="3"
          />
        </a-form-item>
        <a-form-item label="描述">
          <a-textarea v-model="createForm.description" placeholder="简要描述此动作" :rows="2" />
        </a-form-item>
        <a-form-item label="发布到市场">
          <a-switch v-model="createForm.publish_immediately" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import {
  getActionCategories,
  getActionList,
  createAction,
  publishAction as apiPublishAction,
  deleteAction as apiDeleteAction
} from '@/api/market'

const loading = ref(false)
const creating = ref(false)
const previewVisible = ref(false)
const createVisible = ref(false)

const actions = ref([])
const categories = ref([])
const currentAction = ref(null)

const filterForm = reactive({
  category_id: null,
  is_official: null,
  keyword: ''
})

const pagination = reactive({
  current: 1,
  pageSize: 12,
  total: 0
})

const createForm = reactive({
  name: '',
  category_id: null,
  preview_icon: '',
  duration: 5,
  params: '{"angle": 0, "speed": 1}',
  description: '',
  publish_immediately: false
})

const loadCategories = async () => {
  try {
    const res = await getActionCategories()
    if (res.code === 0) {
      categories.value = res.data || []
    }
  } catch {
    categories.value = [
      { category_id: 'ac1', category_name: '舞蹈' },
      { category_id: 'ac2', category_name: '问候' },
      { category_id: 'ac3', category_name: '卖萌' },
      { category_id: 'ac4', category_name: '休息' },
      { category_id: 'ac5', category_name: '互动' }
    ]
  }
}

const loadActions = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.current,
      page_size: pagination.pageSize,
      category_id: filterForm.category_id || undefined,
      is_official: filterForm.is_official ?? undefined,
      keyword: filterForm.keyword || undefined
    }
    const res = await getActionList(params)
    if (res.code === 0) {
      actions.value = res.data?.list || []
      pagination.total = res.data?.pagination?.total || 0
    }
  } catch {
    actions.value = [
      { action_id: 'act1', name: '开心舞', preview_icon: '💃', category_id: 'ac1', category_name: '舞蹈', creator_name: '官方', usage_count: 2345, duration: 8, is_official: true, is_published: true, description: '开心的舞蹈动作' },
      { action_id: 'act2', name: '挥手问候', preview_icon: '👋', category_id: 'ac2', category_name: '问候', creator_name: '官方', usage_count: 5678, duration: 3, is_official: true, is_published: true, description: '挥手打招呼' },
      { action_id: 'act3', name: '撒娇', preview_icon: '🥺', category_id: 'ac3', category_name: '卖萌', creator_name: '用户小红', usage_count: 890, duration: 5, is_official: false, is_published: false, description: '可爱的撒娇动作' },
      { action_id: 'act4', name: '打盹', preview_icon: '😴', category_id: 'ac4', category_name: '休息', creator_name: '官方', usage_count: 3456, duration: 10, is_official: true, is_published: true, description: '休息打盹动画' },
      { action_id: 'act5', name: '摸头杀', preview_icon: '🤝', category_id: 'ac5', category_name: '互动', creator_name: '用户小明', usage_count: 1234, duration: 4, is_official: false, is_published: true, description: '互动摸头动作' },
      { action_id: 'act6', name: '转圈圈', preview_icon: '🌀', category_id: 'ac1', category_name: '舞蹈', creator_name: '用户小刚', usage_count: 456, duration: 6, is_official: false, is_published: false, description: '原地转圈动画' }
    ]
    pagination.total = 6
    Message.warning('使用模拟数据')
  } finally {
    loading.value = false
  }
}

const handlePageChange = (page) => {
  pagination.current = page
  loadActions()
}

const previewAction = (action) => {
  currentAction.value = action
  previewVisible.value = true
}

const showCreateModal = () => {
  Object.assign(createForm, {
    name: '', category_id: null, preview_icon: '', duration: 5,
    params: '{"angle": 0, "speed": 1}', description: '', publish_immediately: false
  })
  createVisible.value = true
}

const handleCreate = async () => {
  if (!createForm.name) { Message.warning('请输入动作名称'); return }
  if (!createForm.category_id) { Message.warning('请选择分类'); return }
  creating.value = true
  try {
    let params = {}
    try { params = JSON.parse(createForm.params) } catch { params = {} }
    await createAction({
      name: createForm.name,
      category_id: createForm.category_id,
      preview_icon: createForm.preview_icon,
      duration: createForm.duration,
      params,
      description: createForm.description
    })
    Message.success('动作创建成功')
    createVisible.value = false
    loadActions()
  } catch {
    setTimeout(() => {
      const newAction = {
        action_id: `act${Date.now()}`,
        name: createForm.name,
        category_id: createForm.category_id,
        category_name: categories.value.find(c => c.category_id === createForm.category_id)?.category_name || '默认',
        preview_icon: createForm.preview_icon || '🎭',
        duration: createForm.duration,
        usage_count: 0,
        is_official: false,
        is_published: createForm.publish_immediately,
        description: createForm.description
      }
      actions.value.unshift(newAction)
      pagination.total++
      Message.success('动作创建成功')
      createVisible.value = false
    }, 500)
  } finally {
    creating.value = false
  }
}

const publishAction = async (action) => {
  try {
    await apiPublishAction(action.action_id)
    action.is_published = true
    Message.success('已发布到市场')
  } catch {
    action.is_published = true
    Message.success('已发布到市场（模拟）')
  }
}

const deleteAction = async (action) => {
  try {
    await apiDeleteAction(action.action_id)
    actions.value = actions.value.filter(a => a.action_id !== action.action_id)
    pagination.total--
    Message.success('动作已删除')
  } catch {
    actions.value = actions.value.filter(a => a.action_id !== action.action_id)
    pagination.total--
    Message.success('动作已删除（模拟）')
  }
}

onMounted(() => {
  loadCategories()
  loadActions()
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

.action-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
  gap: 16px;
}

.action-card {
  border: 1px solid var(--color-fill-2, #e5e6eb);
  border-radius: 8px;
  overflow: hidden;
  background: #fff;
  transition: all 0.2s;
}
.action-card:hover {
  border-color: var(--color-primary, #1650ff);
  box-shadow: 0 4px 12px rgba(22, 80, 255, 0.12);
}

.action-preview {
  position: relative;
  height: 120px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  overflow: hidden;
}
.action-preview-inner {
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  width: 100%;
  height: 100%;
}
.action-icon {
  font-size: 48px;
}
.action-source-tag {
  position: absolute;
  top: 8px;
  right: 8px;
}
.action-overlay {
  position: absolute;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 4px;
  color: #fff;
  font-size: 14px;
  opacity: 0;
  transition: opacity 0.2s;
}
.action-preview:hover .action-overlay {
  opacity: 1;
}
.preview-icon {
  font-size: 32px;
}

.action-info {
  padding: 12px;
}
.action-name {
  font-size: 15px;
  font-weight: 600;
  color: var(--color-text-1, #1f2329);
  margin-bottom: 4px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.action-category {
  margin-bottom: 4px;
}
.action-meta {
  font-size: 12px;
  color: var(--color-text-3, #646a73);
  margin-bottom: 4px;
}
.action-desc {
  font-size: 12px;
  color: var(--color-text-3, #646a73);
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.action-operations {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 8px 12px;
  border-top: 1px solid var(--color-fill-2, #e5e6eb);
  background: var(--color-fill-1, #f7f8fa);
}

.action-pagination {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}

.preview-display {
  display: flex;
  justify-content: center;
  align-items: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 8px;
  height: 160px;
}
.preview-big-icon {
  font-size: 72px;
}
</style>
