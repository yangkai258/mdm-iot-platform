<template>
  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>具身智能</a-breadcrumb-item>
      <a-breadcrumb-item>动作库</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 操作区 -->
    <a-card class="action-header-card">
      <a-space wrap>
        <a-button type="primary" @click="openRecordModal">
          <template #icon><icon-video-camera /></template>
          录制动作
        </a-button>
        <a-button @click="openShareModal" :disabled="!selectedAction">
          <template #icon><icon-share /></template>
          分享
        </a-button>
        <a-button type="text" danger @click="handleDelete" :disabled="!selectedAction || selectedAction.action_type === 'built-in'">
          <template #icon><icon-delete /></template>
          删除
        </a-button>
        <a-divider type="vertical" />
        <a-input-search
          v-model="keyword"
          placeholder="搜索动作名称/标签"
          style="width: 240px"
          @search="loadActions"
          @change="debouncedLoad"
        />
        <a-select v-model="filterType" placeholder="动作类型" allow-clear style="width: 140px" @change="loadActions">
          <a-option value="built-in">内置</a-option>
          <a-option value="learned">学习</a-option>
          <a-option value="custom">自定义</a-option>
        </a-select>
      </a-space>
    </a-card>

    <!-- 动作分类统计 -->
    <a-row :gutter="12" style="margin-top: 16px">
      <a-col :xs="8" :sm="8" :md="6">
        <a-card class="stat-card" :class="{ active: filterType === 'built-in' }" @click="filterType = 'built-in'; loadActions()">
          <div class="stat-value">{{ stats.built_in || 0 }}</div>
          <div class="stat-label">内置动作</div>
        </a-card>
      </a-col>
      <a-col :xs="8" :sm="8" :md="6">
        <a-card class="stat-card" :class="{ active: filterType === 'learned' }" @click="filterType = 'learned'; loadActions()">
          <div class="stat-value">{{ stats.learned || 0 }}</div>
          <div class="stat-label">学习动作</div>
        </a-card>
      </a-col>
      <a-col :xs="8" :sm="8" :md="6">
        <a-card class="stat-card" :class="{ active: filterType === 'custom' }" @click="filterType = 'custom'; loadActions()">
          <div class="stat-value">{{ stats.custom || 0 }}</div>
          <div class="stat-label">自定义动作</div>
        </a-card>
      </a-col>
      <a-col :xs="24" :sm="24" :md="6">
        <a-card class="stat-card" @click="filterType = ''; loadActions()">
          <div class="stat-value">{{ stats.total || 0 }}</div>
          <div class="stat-label">全部动作</div>
        </a-card>
      </a-col>
    </a-row>

    <!-- 动作列表 -->
    <a-row :gutter="16" style="margin-top: 16px">
      <a-col :xs="24" :sm="24" :md="16">
        <a-card title="动作列表">
          <template #extra>
            <a-pagination
              :current="page"
              :page-size="pageSize"
              :total="total"
              :show-total="true"
              @change="loadActions"
              size="small"
            />
          </template>
          <a-table
            :columns="columns"
            :data="actions"
            :loading="loading"
            :pagination="false"
            row-key="id"
            :row-class="rowClass"
            @row-click="selectAction"
            size="small"
          >
            <template #action_type="{ record }">
              <a-tag :color="getTypeColor(record.action_type)">
                {{ getTypeText(record.action_type) }}
              </a-tag>
            </template>
            <template #difficulty="{ record }">
              <a-tag v-if="record.difficulty" :color="getDifficultyColor(record.difficulty)">
                {{ getDifficultyText(record.difficulty) }}
              </a-tag>
              <span v-else>-</span>
            </template>
            <template #score="{ record }">
              <a-rate :model-value="record.score || 0" readonly allow-half :count="5" size="small" />
            </template>
            <template #actions="{ record }">
              <a-button type="text" size="small" @click.stop="playAction(record)">
                <icon-play-circle />
              </a-button>
              <a-button type="text" size="small" @click.stop="openLearnModal(record)" v-if="record.action_type === 'built-in'">
                <icon-robot />
              </a-button>
            </template>
          </a-table>
        </a-card>
      </a-col>

      <!-- 动作详情 -->
      <a-col :xs="24" :sm="24" :md="8">
        <a-card title="动作详情" v-if="selectedAction">
          <template #extra>
            <a-tag :color="getTypeColor(selectedAction.action_type)">
              {{ getTypeText(selectedAction.action_type) }}
            </a-tag>
          </template>
          <a-descriptions :column="1" size="small">
            <a-descriptions-item label="动作名称">{{ selectedAction.action_name }}</a-descriptions-item>
            <a-descriptions-item label="难度">
              <a-tag v-if="selectedAction.difficulty" :color="getDifficultyColor(selectedAction.difficulty)">
                {{ getDifficultyText(selectedAction.difficulty) }}
              </a-tag>
              <span v-else>-</span>
            </a-descriptions-item>
            <a-descriptions-item label="时长">
              {{ selectedAction.duration_ms ? `${(selectedAction.duration_ms / 1000).toFixed(1)}s` : '-' }}
            </a-descriptions-item>
            <a-descriptions-item label="评分">
              <a-rate :model-value="selectedAction.score || 0" readonly allow-half :count="5" />
            </a-descriptions-item>
            <a-descriptions-item label="下载量">{{ selectedAction.downloads || 0 }}</a-descriptions-item>
            <a-descriptions-item label="创建时间">{{ formatTime(selectedAction.created_at) }}</a-descriptions-item>
          </a-descriptions>
          <a-divider>标签</a-divider>
          <a-tag v-for="tag in selectedAction.tags" :key="tag" class="action-tag">{{ tag }}</a-tag>
          <a-empty v-if="!selectedAction.tags?.length" description="暂无标签" />
          <a-divider>描述</a-divider>
          <p class="action-desc">{{ selectedAction.description || '暂无描述' }}</p>

          <!-- 动作执行 -->
          <a-divider>执行</a-divider>
          <a-select v-model="executeDeviceId" placeholder="选择设备" size="small" style="width: 100%; margin-bottom: 8px">
            <a-option v-for="d in devices" :key="d.device_id" :value="d.device_id">
              {{ d.device_name || d.device_id }}
            </a-option>
          </a-select>
          <a-button type="primary" long @click="executeSelectedAction" :loading="executing" :disabled="!executeDeviceId">
            <template #icon><icon-play-circle /></template>
            执行动作
          </a-button>

          <!-- 动作执行历史 -->
          <a-divider>执行历史</a-divider>
          <a-list :data-source="executions" size="small">
            <template #renderItem="{ item }">
              <a-list-item>
                <a-list-item-meta>
                  <template #title>
                    <a-space>
                      <a-tag :color="getExecStatusColor(item.status)" size="small">
                        {{ getExecStatusText(item.status) }}
                      </a-tag>
                      {{ item.execution_type }}
                    </a-space>
                  </template>
                  <template #description>
                    {{ formatTime(item.start_time) }}
                    <span v-if="item.interruption_reason" class="text-danger"> · {{ item.interruption_reason }}</span>
                  </template>
                </a-list-item-meta>
              </a-list-item>
            </template>
          </a-list>
        </a-card>
        <a-card v-else title="动作详情">
          <a-empty description="请选择一个动作" />
        </a-card>
      </a-col>
    </a-row>

    <!-- 录制动作弹窗 -->
    <a-modal
      v-model:visible="recordModalVisible"
      title="录制动作"
      @before-ok="submitRecord"
      @cancel="recordModalVisible = false"
      :width="500"
    >
      <a-form :model="recordForm" layout="vertical">
        <a-form-item label="动作名称" required>
          <a-input v-model="recordForm.action_name" placeholder="输入动作名称" />
        </a-form-item>
        <a-form-item label="动作描述">
          <a-textarea v-model="recordForm.description" :rows="3" placeholder="动作描述..." />
        </a-form-item>
        <a-form-item label="难度">
          <a-select v-model="recordForm.difficulty" placeholder="选择难度">
            <a-option value="easy">简单</a-option>
            <a-option value="medium">中等</a-option>
            <a-option value="hard">困难</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="标签">
          <a-select v-model="recordForm.tags" multiple placeholder="选择/输入标签" allow-create>
            <a-option value="dance">舞蹈</a-option>
            <a-option value="greeting">问候</a-option>
            <a-option value="exercise">运动</a-option>
            <a-option value="fun">娱乐</a-option>
            <a-option value="utility">工具</a-option>
          </a-select>
        </a-form-item>
        <a-alert type="info">
          录制将启动设备动作捕捉，录制完成后自动保存。
        </a-alert>
      </a-form>
      <template #footer>
        <a-button @click="recordModalVisible = false">取消</a-button>
        <a-button type="primary" :loading="recording" @click="startRecording">
          <template #icon><icon-video-camera /></template>
          开始录制
        </a-button>
      </template>
    </a-modal>

    <!-- 学习动作弹窗 -->
    <a-modal
      v-model:visible="learnModalVisible"
      title="学习动作"
      @before-ok="submitLearn"
      @cancel="learnModalVisible = false"
    >
      <a-result
        status="info"
        title="AI 学习"
        sub-title="系统将使用 AI 模型分析并学习此动作的要点"
      >
        <template #extra>
          <a-form :model="learnForm" layout="vertical" style="width: 400px">
            <a-form-item label="学习次数">
              <a-input-number v-model="learnForm.learn_count" :min="1" :max="100" style="width: 100%" />
            </a-form-item>
            <a-form-item label="备注">
              <a-textarea v-model="learnForm.note" :rows="2" placeholder="备注信息..." />
            </a-form-item>
          </a-form>
        </template>
      </a-result>
    </a-modal>

    <!-- 分享弹窗 -->
    <a-modal
      v-model:visible="shareModalVisible"
      title="分享动作"
      @before-ok="submitShare"
      @cancel="shareModalVisible = false"
    >
      <a-result
        status="success"
        title="即将分享动作"
        :sub-title="`分享「${selectedAction?.action_name}」到公共动作库`"
      />
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getActionLibrary, recordAction, learnAction, executeAction, shareAction, deleteAction, getActionExecutions } from '@/api/embodied'
import { Message, Modal } from '@arco-design/web-vue'

const loading = ref(false)
const actions = ref<any[]>([])
const selectedAction = ref<any>(null)
const stats = ref<any>({})
const keyword = ref('')
const filterType = ref('')
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)

const recordModalVisible = ref(false)
const learnModalVisible = ref(false)
const shareModalVisible = ref(false)
const recording = ref(false)
const executing = ref(false)
const executeDeviceId = ref('')
const executions = ref<any[]>([])
const devices = ref<any[]>([])

const recordForm = ref({ action_name: '', description: '', difficulty: 'medium', tags: [] as string[] })
const learnForm = ref({ learn_count: 10, note: '' })

const columns = [
  { title: 'ID', dataIndex: 'id', width: 70 },
  { title: '动作名称', dataIndex: 'action_name', width: 160 },
  { title: '类型', dataIndex: 'action_type', slotName: 'action_type', width: 100 },
  { title: '难度', dataIndex: 'difficulty', slotName: 'difficulty', width: 90 },
  { title: '时长', dataIndex: 'duration_ms', width: 80 },
  { title: '评分', dataIndex: 'score', slotName: 'score', width: 120 },
  { title: '操作', slotName: 'actions', width: 80 }
]

let loadTimer: number | null = null

async function loadActions() {
  try {
    loading.value = true
    const params: any = { page: page.value, page_size: pageSize.value }
    if (keyword.value) params.keyword = keyword.value
    if (filterType.value) params.action_type = filterType.value
    const res = await getActionLibrary(params)
    const data = res.data
    actions.value = data?.actions || data?.list || data || []
    total.value = data?.total || actions.value.length

    // 计算统计
    const all: any[] = actions.value
    stats.value = {
      built_in: all.filter((a: any) => a.action_type === 'built-in').length,
      learned: all.filter((a: any) => a.action_type === 'learned').length,
      custom: all.filter((a: any) => a.action_type === 'custom').length,
      total: total.value
    }
  } catch (err: any) {
    Message.error('加载动作库失败: ' + err.message)
  } finally {
    loading.value = false
  }
}

function debouncedLoad() {
  if (loadTimer) clearTimeout(loadTimer)
  loadTimer = window.setTimeout(() => loadActions(), 400)
}

function selectAction(record: any) {
  selectedAction.value = record
  loadExecutions(record.id)
}

function rowClass(record: any) {
  return selectedAction.value?.id === record.id ? 'selected-row' : ''
}

async function loadExecutions(actionId: number) {
  try {
    const res = await getActionExecutions(deviceId.value || '', { action_id: actionId, page_size: 10 })
    executions.value = res.data?.executions || res.data || []
  } catch {}
}

async function executeSelectedAction() {
  if (!selectedAction.value || !executeDeviceId.value) return
  try {
    executing.value = true
    await executeAction(executeDeviceId.value, { action_id: selectedAction.value.id })
    Message.success('动作已下发执行')
    await loadExecutions(selectedAction.value.id)
  } catch (err: any) {
    Message.error('执行失败: ' + err.message)
  } finally {
    executing.value = false
  }
}

function openRecordModal() {
  recordForm.value = { action_name: '', description: '', difficulty: 'medium', tags: [] }
  recordModalVisible.value = true
}

async function startRecording() {
  if (!recordForm.value.action_name) {
    Message.warning('请输入动作名称')
    return
  }
  try {
    recording.value = true
    await recordAction({ device_id: executeDeviceId.value, ...recordForm.value })
    Message.success('录制已启动')
    recordModalVisible.value = false
    await loadActions()
  } catch (err: any) {
    Message.error('录制失败: ' + err.message)
  } finally {
    recording.value = false
  }
}

async function submitRecord(done: (val: boolean) => void) {
  done(true)
}

async function submitLearn(done: (val: boolean) => void) {
  if (!selectedAction.value) return
  try {
    await learnAction(selectedAction.value.id, learnForm.value)
    Message.success('学习任务已创建')
    learnModalVisible.value = false
    await loadActions()
    done(true)
  } catch (err: any) {
    Message.error('学习失败: ' + err.message)
    done(false)
  }
}

async function submitShare(done: (val: boolean) => void) {
  if (!selectedAction.value) return
  try {
    await shareAction(selectedAction.value.id)
    Message.success('已分享到公共动作库')
    shareModalVisible.value = false
    done(true)
  } catch (err: any) {
    Message.error('分享失败: ' + err.message)
    done(false)
  }
}

async function handleDelete() {
  if (!selectedAction.value) return
  Modal.warning({
    title: '确认删除',
    content: `确定删除动作「${selectedAction.value.action_name}」？`,
    onOk: async () => {
      try {
        await deleteAction(selectedAction.value.id)
        Message.success('已删除')
        selectedAction.value = null
        await loadActions()
      } catch (err: any) {
        Message.error('删除失败: ' + err.message)
      }
    }
  })
}

function openLearnModal(record: any) {
  selectedAction.value = record
  learnForm.value = { learn_count: 10, note: '' }
  learnModalVisible.value = true
}

async function playAction(record: any) {
  if (!executeDeviceId.value) {
    Message.warning('请先在详情区选择设备')
    return
  }
  selectedAction.value = record
  await executeSelectedAction()
}

function getTypeColor(type: string) {
  const map: Record<string, string> = { 'built-in': 'green', learned: 'arcoblue', custom: 'orange' }
  return map[type] || 'default'
}
function getTypeText(type: string) {
  const map: Record<string, string> = { 'built-in': '内置', learned: '学习', custom: '自定义' }
  return map[type] || type
}
function getDifficultyColor(d: string) {
  const map: Record<string, string> = { easy: 'green', medium: 'orange', hard: 'red' }
  return map[d] || 'default'
}
function getDifficultyText(d: string) {
  const map: Record<string, string> = { easy: '简单', medium: '中等', hard: '困难' }
  return map[d] || d
}
function getExecStatusColor(s: string) {
  const map: Record<string, string> = { running: 'arcoblue', completed: 'green', interrupted: 'orange', failed: 'red' }
  return map[s] || 'default'
}
function getExecStatusText(s: string) {
  const map: Record<string, string> = { running: '执行中', completed: '已完成', interrupted: '中断', failed: '失败' }
  return map[s] || s
}
function formatTime(t: string) {
  if (!t) return '-'
  return new Date(t).toLocaleString('zh-CN')
}

onMounted(() => {
  loadActions()
  // 模拟设备列表
  devices.value = [{ device_id: 'pet-001', device_name: '电子宠物 1号' }]
})
</script>

<style scoped>
.stat-card {
  text-align: center;
  cursor: pointer;
  transition: all 0.2s;
}
.stat-card.active {
  border-color: var(--color-primary);
}
.stat-value {
  font-size: 28px;
  font-weight: 600;
  color: var(--color-primary);
}
.stat-label {
  font-size: 12px;
  color: var(--color-text-3);
  margin-top: 4px;
}
.action-tag {
  margin: 4px 4px 4px 0;
}
.action-desc {
  color: var(--color-text-2);
  font-size: 13px;
  line-height: 1.6;
}
.selected-row {
  background: var(--color-primary-light-1) !important;
}
.text-danger {
  color: #ff4d4f;
}
</style>
