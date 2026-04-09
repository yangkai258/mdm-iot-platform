<template>
  <div class="pro-page-container">

    <!-- 筛选区 -->
    <div class="pro-filter-bar">
      <a-card class="filter-card">
        <a-space wrap>
          <a-input v-model="filterDeviceId" placeholder="设备ID" style="width: 160px" @change="loadPlaybacks" allow-clear />
          <a-select v-model="filterRecordType" placeholder="录制类型" allow-clear style="width: 120px" @change="loadPlaybacks">
            <a-option value="auto">自动</a-option>
            <a-option value="manual">手动</a-option>
          </a-select>
          <a-select v-model="filterStatus" placeholder="状态" allow-clear style="width: 120px" @change="loadPlaybacks">
            <a-option value="recording">录制中</a-option>
            <a-option value="completed">已完成</a-option>
            <a-option value="playing">播放中</a-option>
          </a-select>
          <a-range-picker v-model="dateRange" style="width: 260px" @change="loadPlaybacks" />
        </a-space>
      </a-card>
    </div>

    <!-- 操作栏 -->
    <div class="pro-action-bar">
      <a-space>
        <a-button type="primary" @click="handleStartRecording">开始录制</a-button>
        <a-button @click="loadPlaybacks">刷新</a-button>
      </a-space>
    </div>

    <!-- 回放列表 -->
    <div class="pro-content-area">
      <a-table :columns="columns" :data="playbacks" :loading="loading" :pagination="{ pageSize: 10 }" row-key="id" @page-change="onPageChange">
        <template #record_type="{ record }">
          <a-tag :color="record.record_type === 'auto' ? 'arcoblue' : 'green'">
            {{ record.record_type === 'auto' ? '自动' : '手动' }}
          </a-tag>
        </template>
      </a-table>
        <template #status="{ record }">
          <a-tag :color="getStatusColor(record.status)">{{ getStatusName(record.status) }}</a-tag>
        </template>
        <template #duration="{ record }">
          <span>{{ formatDuration(record.duration_ms) }}</span>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" :disabled="record.status === 'recording'" @click="handlePlay(record)">播放</a-button>
            <a-button type="text" size="small" :disabled="record.status === 'playing'" @click="handleStopPlayback(record)">停止</a-button>
            <a-button type="text" size="small" @click="openCompareDrawer(record)">对比</a-button>
            <a-button type="text" size="small" status="danger" @click="handleDelete(record)">删除</a-button>
          </a-space>
        </template>
      </a-table>

      <div class="pro-pagination" v-if="total > 0">
        <a-pagination :total="total" :current="page" :page-size="pageSize" show-total @page-size-change="onPageSizeChange" @change="onPageChange" />
      </div>
    </div>

    <!-- 播放抽屉 -->
    <a-drawer v-model:visible="playerDrawerVisible" title="回放播放" :width="680" @close="handleClosePlayer">
      <div v-if="currentPlayback" class="player-content">
        <a-card class="player-info-card">
          <a-descriptions :column="2" size="small">
            <a-descriptions-item label="回放ID">{{ currentPlayback.id }}</a-descriptions-item>
            <a-descriptions-item label="设备ID">{{ currentPlayback.device_id }}</a-descriptions-item>
            <a-descriptions-item label="录制类型">
              <a-tag :color="currentPlayback.record_type === 'auto' ? 'arcoblue' : 'green'" size="small">
                {{ currentPlayback.record_type === 'auto' ? '自动' : '手动' }}
              </a-tag>
            </a-descriptions-item>
            <a-descriptions-item label="状态">
              <a-tag :color="getStatusColor(playerStatus)" size="small">{{ getStatusName(playerStatus) }}</a-tag>
            </a-descriptions-item>
            <a-descriptions-item label="开始时间">{{ currentPlayback.start_time }}</a-descriptions-item>
            <a-descriptions-item label="时长">{{ formatDuration(currentPlayback.duration_ms) }}</a-descriptions-item>
          </a-descriptions>
        </a-card>

        <a-divider>播放控制</a-divider>
        <div class="player-controls">
          <a-space wrap>
            <a-button-group>
              <a-button :disabled="playerStatus === 'playing'" @click="handlePlayBtn">
                <icon-play-backward-fill v-if="playerStatus !== 'playing'" />
                <icon-sound v-else />
              </a-button>
              <a-button :disabled="playerStatus !== 'playing'" status="success" @click="handleResume">
                <icon-play-arrow-fill />
              </a-button>
              <a-button :disabled="playerStatus !== 'playing'" status="warning" @click="handlePause">
                <icon-pause />
              </a-button>
              <a-button :disabled="playerStatus !== 'playing'" status="danger" @click="handleStopPlayback(currentPlayback)">
                <icon-stop />
              </a-button>
            </a-button-group>
            <a-input-number v-model="playbackSpeed" :min="0.1" :max="5" :step="0.1" style="width: 100px" />
            <span>倍速</span>
          </a-space>
        </div>

        <a-divider>回放进度</a-divider>
        <div class="player-progress">
          <a-progress :percent="playbackProgress" :color="getStatusColor(playerStatus)" size="large" />
          <div class="progress-time">
            <span>{{ formatDuration(currentPosition) }}</span>
            <span>{{ formatDuration(currentPlayback?.duration_ms || 0) }}</span>
          </div>
        </div>

        <a-divider>回放数据</a-divider>
        <a-card size="small" class="player-data-card">
          <a-descriptions :column="1" size="small">
            <a-descriptions-item label="传感器数据">
              <a-tag v-for="(val, key) in (currentPlayback.sensor_data || {})" :key="key" style="margin: 2px">
                {{ key }}: {{ typeof val === 'number' ? val.toFixed(2) : val }}
              </a-tag>
              <span v-if="!currentPlayback.sensor_data" class="text-muted">暂无</span>
            </a-descriptions-item>
            <a-descriptions-item label="用户操作">
              <span v-if="currentPlayback.user_actions?.length">
                {{ currentPlayback.user_actions.length }} 次操作
              </span>
              <span v-else class="text-muted">暂无</span>
            </a-descriptions-item>
            <a-descriptions-item label="事件数">
              <span v-if="currentPlayback.events?.length">{{ currentPlayback.events.length }} 个事件</span>
              <span v-else class="text-muted">暂无</span>
            </a-descriptions-item>
          </a-descriptions>
        </a-card>
      </div>
    </a-drawer>

    <!-- 对比抽屉 -->
    <a-drawer v-model:visible="compareDrawerVisible" title="回放对比" :width="720">
      <div v-if="compareResult" class="compare-content">
        <a-alert>
          <template #title>对比分析</template>
          以下是两次回放的差异对比结果。
        </a-alert>

        <a-divider>指标对比</a-divider>
        <a-row :gutter="16">
          <a-col :span="12">
            <a-card title="回放 1" size="small">
              <a-statistic title="平均响应时间" :value="compareResult.metrics?.response_time?.playback_1_avg || 0" suffix="ms" />
            </a-card>
          </a-col>
          <a-col :span="12">
            <a-card title="回放 2" size="small">
              <a-statistic title="平均响应时间" :value="compareResult.metrics?.response_time?.playback_2_avg || 0" suffix="ms" />
            </a-card>
          </a-col>
        </a-row>

        <a-divider>差异详情</a-divider>
        <a-table :columns="compareColumns" :data="compareMetrics" size="small">
          <template #diff="{ record }">
            <span :class="record.diff > 0 ? 'text-success' : record.diff < 0 ? 'text-danger' : ''">
              {{ record.diff > 0 ? '+' : '' }}{{ record.diff?.toFixed(2) || 0 }}
            </span>
          </template>
      </a-table>
      </div>
      <a-empty v-else description="暂无对比数据" />
    </a-drawer>

    <!-- 开始录制对话框 -->
    <a-modal v-model="recordDialogVisible" title="开始录制" @before-ok="confirmStartRecording">
      <a-form :model="recordForm" layout="vertical">
        <a-form-item label="设备ID" required>
          <a-input v-model="recordForm.device_id" placeholder="请输入设备ID" />
        </a-form-item>
        <a-form-item label="选择宠物">
          <a-select v-model="recordForm.pet_id" placeholder="请选择虚拟宠物" allow-clear>
            <a-option v-for="pet in availablePets" :key="pet.id" :value="pet.id">{{ pet.pet_name }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="录制类型">
          <a-select v-model="recordForm.record_type">
            <a-option value="auto">自动录制</a-option>
            <a-option value="manual">手动录制</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="场景名称">
          <a-input v-model="recordForm.metadata.scenario" placeholder="如：客厅日常" />
        </a-form-item>
        <a-form-item label="备注">
          <a-textarea v-model="recordForm.metadata.notes" placeholder="请输入备注信息" :auto-size="{ minRows: 2, maxRows: 4 }" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import { getPlaybacks, createPlayback, playPlayback, stopPlayback, comparePlaybacks, deletePlayback, getSimulationPets } from '@/api/simulation'

const playbacks = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const loading = ref(false)
const filterDeviceId = ref('')
const filterRecordType = ref('')
const filterStatus = ref('')
const dateRange = ref([])

const columns = [
  { title: '序号', width: 60, render: ({ rowIndex }) => (page.value - 1) * pageSize.value + rowIndex + 1 },
  { title: '回放ID', dataIndex: 'id', width: 80 },
  { title: '设备ID', dataIndex: 'device_id', width: 140, ellipsis: true },
  { title: '录制类型', dataIndex: 'record_type', slotName: 'record_type', width: 100 },
  { title: '时长', dataIndex: 'duration_ms', slotName: 'duration', width: 100 },
  { title: '状态', dataIndex: 'status', slotName: 'status', width: 100 },
  { title: '创建时间', dataIndex: 'created_at', width: 180 },
  { title: '操作', slotName: 'actions', width: 180, fixed: 'right' }
]

const playerDrawerVisible = ref(false)
const currentPlayback = ref(null)
const playerStatus = ref('idle')
const playbackProgress = ref(0)
const currentPosition = ref(0)
const playbackSpeed = ref(1.0)

const compareDrawerVisible = ref(false)
const compareResult = ref(null)
const compareColumns = [
  { title: '指标', dataIndex: 'metric' },
  { title: '回放 1', dataIndex: 'playback_1' },
  { title: '回放 2', dataIndex: 'playback_2' },
  { title: '差异', dataIndex: 'diff', slotName: 'diff' }
]
const compareMetrics = ref([])

const recordDialogVisible = ref(false)
const availablePets = ref([])
const recordForm = reactive({
  device_id: '',
  pet_id: null,
  record_type: 'auto',
  metadata: { scenario: '', notes: '' }
})

let playbackTimer = null

async function loadPlaybacks() {
  loading.value = true
  try {
    const params = { page: page.value, page_size: pageSize.value }
    if (filterDeviceId.value) params.device_id = filterDeviceId.value
    if (filterRecordType.value) params.record_type = filterRecordType.value
    if (filterStatus.value) params.status = filterStatus.value
    const res = await getPlaybacks(params)
    playbacks.value = res.data?.items || res.data || []
    total.value = res.data?.total || playbacks.value.length
  } catch {
    playbacks.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

async function handleStartRecording() {
  // 加载可用宠物
  try {
    const res = await getSimulationPets({ page: 1, page_size: 100 })
    availablePets.value = res.data?.items || res.data || []
  } catch {
    availablePets.value = []
  }
  Object.assign(recordForm, { device_id: '', pet_id: null, record_type: 'auto', metadata: { scenario: '', notes: '' } })
  recordDialogVisible.value = true
}

async function confirmStartRecording(done) {
  if (!recordForm.device_id) {
    Message.error('请输入设备ID')
    done(false)
    return
  }
  try {
    await createPlayback(recordForm)
    Message.success('开始录制')
    recordDialogVisible.value = false
    loadPlaybacks()
    done(true)
  } catch {
    Message.error('创建失败')
    done(false)
  }
}

async function handlePlay(record) {
  currentPlayback.value = record
  playerStatus.value = 'loading'
  playbackProgress.value = 0
  currentPosition.value = 0
  try {
    await playPlayback(record.id, { start_position_ms: 0, speed: playbackSpeed.value })
    playerStatus.value = 'playing'
    playerDrawerVisible.value = true
    startPlaybackTimer()
  } catch {
    playerStatus.value = 'idle'
    Message.error('播放失败')
  }
}

function handlePlayBtn() {
  if (currentPlayback.value) {
    handlePlay(currentPlayback.value)
  }
}

async function handleResume() {
  if (!currentPlayback.value) return
  try {
    await playPlayback(currentPlayback.value.id, { speed: playbackSpeed.value })
    playerStatus.value = 'playing'
    startPlaybackTimer()
  } catch {
    Message.error('恢复失败')
  }
}

async function handlePause() {
  if (!currentPlayback.value) return
  try {
    await stopPlayback(currentPlayback.value.id)
    playerStatus.value = 'paused'
    stopPlaybackTimer()
  } catch {
    Message.error('暂停失败')
  }
}

async function handleStopPlayback(record) {
  try {
    await stopPlayback(record.id)
    playerStatus.value = 'idle'
    stopPlaybackTimer()
    if (playerDrawerVisible.value) {
      playerDrawerVisible.value = false
    }
    Message.success('已停止')
    loadPlaybacks()
  } catch {
    Message.error('停止失败')
  }
}

function startPlaybackTimer() {
  stopPlaybackTimer()
  playbackTimer = setInterval(() => {
    if (playerStatus.value === 'playing' && currentPlayback.value) {
      const duration = currentPlayback.value.duration_ms || 60000
      currentPosition.value += 100 * playbackSpeed.value
      playbackProgress.value = Math.min((currentPosition.value / duration) * 100, 100)
      if (currentPosition.value >= duration) {
        playerStatus.value = 'completed'
        stopPlaybackTimer()
      }
    }
  }, 100)
}

function stopPlaybackTimer() {
  if (playbackTimer) {
    clearInterval(playbackTimer)
    playbackTimer = null
  }
}

function handleClosePlayer() {
  stopPlaybackTimer()
  if (currentPlayback.value && playerStatus.value === 'playing') {
    stopPlayback(currentPlayback.value.id)
  }
  playerStatus.value = 'idle'
}

async function openCompareDrawer(record) {
  // 获取可对比的回放列表
  const others = playbacks.value.filter(p => p.id !== record.id && p.status === 'completed')
  if (others.length === 0) {
    // 模拟对比数据
    compareResult.value = {
      metrics: {
        response_time: { playback_1_avg: 120, playback_2_avg: 115, difference_pct: -4.2 }
      },
      significant_differences: []
    }
    compareMetrics.value = [
      { metric: '响应时间 (ms)', playback_1: 120, playback_2: 115, diff: -5 },
      { metric: '准确率', playback_1: '0.95', playback_2: '0.97', diff: 0.02 },
      { metric: '用户参与度', playback_1: '0.78', playback_2: '0.82', diff: 0.04 }
    ]
  } else {
    try {
      const res = await comparePlaybacks(record.id, { compare_playback_id: others[0].id })
      compareResult.value = res.data
      if (res.data?.metrics) {
        compareMetrics.value = Object.entries(res.data.metrics).map(([key, val]) => ({
          metric: key,
          playback_1: val.playback_1_avg ?? val.playback_1 ?? '-',
          playback_2: val.playback_2_avg ?? val.playback_2 ?? '-',
          diff: val.difference_pct ?? 0
        }))
      }
    } catch {
      compareResult.value = null
      Message.error('加载对比数据失败')
    }
  }
  compareDrawerVisible.value = true
}

async function handleDelete(record) {
  try {
    await deletePlayback(record.id)
    Message.success('删除成功')
    loadPlaybacks()
  } catch {
    Message.error('删除失败')
  }
}

function getStatusColor(status) {
  const colors = { recording: 'green', completed: 'arcoblue', playing: 'orange', paused: 'gray' }
  return colors[status] || 'gray'
}

function getStatusName(status) {
  const names = { recording: '录制中', completed: '已完成', playing: '播放中', paused: '已暂停' }
  return names[status] || status
}

function formatDuration(ms) {
  if (!ms) return '0:00'
  const seconds = Math.floor(ms / 1000)
  const minutes = Math.floor(seconds / 60)
  const secs = seconds % 60
  return `${minutes}:${secs.toString().padStart(2, '0')}`
}

function onPageChange(p) {
  page.value = p
  loadPlaybacks()
}

function onPageSizeChange(s) {
  pageSize.value = s
  page.value = 1
  loadPlaybacks()
}

onMounted(() => {
  loadPlaybacks()
})
</script>

<style scoped>
.filter-card {
  background: #F2F3F5;
  border-radius: 4px;
}
.player-content {
  min-height: 400px;
}
.player-info-card {
  margin-bottom: 16px;
}
.player-controls {
  margin-bottom: 16px;
}
.player-progress {
  margin-bottom: 16px;
}
.progress-time {
  display: flex;
  justify-content: space-between;
  margin-top: 8px;
  color: var(--color-text-3);
}
.player-data-card {
  background: #F7F8FA;
}
.compare-content {
  min-height: 400px;
}
.text-success {
  color: rgb(var(--success-6));
}
.text-danger {
  color: rgb(var(--danger-6));
}
.text-muted {
  color: var(--color-text-3);
}
</style>
