<template>
  <div class="pro-page-container">
    <!-- 筛选栏 -->
    <div class="pro-filter-bar">
      <a-card>
        <a-space wrap>
          <a-input v-model="filterDeviceId" placeholder="设备ID" style="width: 160px" allow-clear @change="loadReplays" />
          <a-select v-model="filterType" placeholder="回放类型" allow-clear style="width: 140px" @change="loadReplays">
            <a-option value="device">设备事件</a-option>
            <a-option value="ai">AI决策</a-option>
            <a-option value="system">系统日志</a-option>
          </a-select>
          <a-select v-model="filterStatus" placeholder="状态" allow-clear style="width: 120px" @change="loadReplays">
            <a-option value="ready">就绪</a-option>
            <a-option value="playing">播放中</a-option>
            <a-option value="paused">已暂停</a-option>
            <a-option value="finished">已结束</a-option>
          </a-select>
          <a-range-picker v-model="dateRange" style="width: 260px" @change="loadReplays" />
          <a-button type="primary" @click="loadReplays">
            <template #icon><icon-search /></template>
            搜索
          </a-button>
        </a-space>
      </a-card>
    </div>

    <!-- 操作栏 -->
    <div class="pro-action-bar">
      <a-space>
        <a-button type="primary" @click="startNewReplay">
          <template #icon><icon-video-camera /></template>
          新建回放
        </a-button>
        <a-button @click="loadReplays">
          <template #icon><icon-refresh /></template>
          刷新
        </a-button>
      </a-space>
      <a-space v-if="currentReplay">
        <a-tag color="arcoblue">正在播放: {{ currentReplay.name }}</a-tag>
      </a-space>
    </div>

    <!-- 回放列表 -->
    <div class="pro-content-area">
      <a-row :gutter="[16, 16]">
        <a-col :xs="24" :sm="12" :md="8" v-for="replay in replays" :key="replay.id">
          <a-card class="replay-card" :class="{ active: currentReplay?.id === replay.id }">
            <template #title>
              <div class="replay-title">
                <icon-video-camera />
                <span>{{ replay.name }}</span>
              </div>
            </template>
            <template #extra>
              <a-tag :color="getStatusColor(replay.status)" size="small">
                {{ getStatusName(replay.status) }}
              </a-tag>
            </template>
            <a-descriptions :column="1" size="small">
              <a-descriptions-item label="设备ID">{{ replay.device_id }}</a-descriptions-item>
              <a-descriptions-item label="回放类型">{{ getTypeName(replay.type) }}</a-descriptions-item>
              <a-descriptions-item label="时间范围">
                {{ replay.start_time }} ~ {{ replay.end_time }}
              </a-descriptions-item>
              <a-descriptions-item label="事件数量">
                <a-badge :text="`${replay.event_count} 条事件`" />
              </a-descriptions-item>
              <a-descriptions-item label="时长">{{ replay.duration }}</a-descriptions-item>
            </a-descriptions>
            <template #actions>
              <a-button type="text" size="small" :disabled="currentReplay?.id === replay.id && replay.status === 'playing'" @click="playReplay(replay)">
                <template #icon><icon-play-circle /></template>
                {{ currentReplay?.id === replay.id && replay.status === 'playing' ? '播放中' : '播放' }}
              </a-button>
              <a-button type="text" size="small" :disabled="currentReplay?.id !== replay.id" @click="pauseReplay">
                <template #icon><icon-pause /></template>
                暂停
              </a-button>
              <a-button type="text" size="small" @click="viewReplayDetail(replay)">
                <template #icon><icon-eye /></template>
                详情
              </a-button>
            </template>
          </a-card>
        </a-col>
      </a-row>
      <a-empty v-if="!replays.length && !loading" description="暂无回放记录" />
    </div>

    <!-- 分页 -->
    <div class="pagination-wrapper" v-if="total > 0">
      <a-pagination
        :total="total"
        :current="page"
        :page-size="pageSize"
        show-total
        @change="onPageChange"
      />
    </div>

    <!-- 回放播放器弹窗 -->
    <a-modal v-model:visible="playerVisible" :title="`回放: ${playerReplay?.name}`" :width="900" footer="null" @cancel="closePlayer">
      <div class="player-wrapper">
        <!-- 时间轴 -->
        <div class="timeline-bar">
          <div class="timeline-track" @click="seekTimeline">
            <div class="timeline-progress" :style="{ width: playProgress + '%' }" />
            <div class="timeline-marker" :style="{ left: playProgress + '%' }" />
            <!-- 事件标注点 -->
            <div
              v-for="evt in timelineEvents"
              :key="evt.id"
              class="timeline-event"
              :class="evt.type"
              :style="{ left: evt.position + '%' }"
              :title="`${evt.time} - ${evt.type}`"
              @click.stop="jumpToEvent(evt)"
            />
          </div>
          <div class="timeline-time">
            <span>{{ currentTimeStr }}</span>
            <span>{{ totalTimeStr }}</span>
          </div>
        </div>

        <!-- 回放控制栏 -->
        <div class="player-controls">
          <a-space>
            <a-button-group>
              <a-button @click="skipBackward">
                <template #icon><icon-rewind-line /></template>
              </a-button>
              <a-button type="primary" @click="togglePlay">
                <template #icon><icon-play-circle v-if="!isPlaying" /><icon-pause v-else /></template>
                {{ isPlaying ? '暂停' : '播放' }}
              </a-button>
              <a-button @click="skipForward">
                <template #icon><icon-fast-forward-line /></template>
              </a-button>
            </a-button-group>
            <a-button @click="setPlaybackSpeed(0.5)">0.5x</a-button>
            <a-button @click="setPlaybackSpeed(1)">1x</a-button>
            <a-button @click="setPlaybackSpeed(2)">2x</a-button>
          </a-space>
          <a-space>
            <span class="speed-label">播放速度: {{ playbackSpeed }}x</span>
          </a-space>
        </div>

        <!-- 回放内容区 -->
        <div class="player-content">
          <div class="player-main">
            <div class="replay-visual">
              <icon-video-camera :size="48" />
              <span>回放可视化区域</span>
              <span class="replay-device">{{ playerReplay?.device_id }}</span>
            </div>
          </div>
          <!-- 事件日志 -->
          <div class="event-log">
            <div class="log-title">事件日志</div>
            <div class="log-list">
              <div
                v-for="evt in visibleEvents"
                :key="evt.id"
                class="log-item"
                :class="evt.type"
              >
                <span class="log-time">{{ evt.time }}</span>
                <a-tag :color="getEventColor(evt.type)" size="small">{{ evt.type }}</a-tag>
                <span class="log-msg">{{ evt.message }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </a-modal>

    <!-- 详情弹窗 -->
    <a-modal v-model:visible="detailVisible" title="回放详情" :width="640" footer="null">
      <a-descriptions :column="1" size="small" bordered>
        <a-descriptions-item label="回放名称">{{ detailReplay?.name }}</a-descriptions-item>
        <a-descriptions-item label="设备ID">{{ detailReplay?.device_id }}</a-descriptions-item>
        <a-descriptions-item label="类型">{{ getTypeName(detailReplay?.type) }}</a-descriptions-item>
        <a-descriptions-item label="状态">
          <a-tag :color="getStatusColor(detailReplay?.status)">{{ getStatusName(detailReplay?.status) }}</a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="时间范围">{{ detailReplay?.start_time }} ~ {{ detailReplay?.end_time }}</a-descriptions-item>
        <a-descriptions-item label="总时长">{{ detailReplay?.duration }}</a-descriptions-item>
        <a-descriptions-item label="事件数量">{{ detailReplay?.event_count }} 条</a-descriptions-item>
        <a-descriptions-item label="创建时间">{{ detailReplay?.created_at }}</a-descriptions-item>
      </a-descriptions>
    </a-modal>

    <!-- 新建回放弹窗 -->
    <a-modal v-model:visible="createVisible" title="新建回放" @before-ok="handleCreate">
      <a-form :model="createForm" layout="vertical">
        <a-form-item label="回放名称" required>
          <a-input v-model="createForm.name" placeholder="请输入回放名称" />
        </a-form-item>
        <a-form-item label="设备ID" required>
          <a-select v-model="createForm.device_id" placeholder="请选择设备" allow-search>
            <a-option value="pet-001">pet-001 (小白)</a-option>
            <a-option value="pet-002">pet-002 (旺财)</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="回放类型" required>
          <a-select v-model="createForm.type" placeholder="请选择类型">
            <a-option value="device">设备事件</a-option>
            <a-option value="ai">AI决策</a-option>
            <a-option value="system">系统日志</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="时间范围" required>
          <a-range-picker v-model="createForm.dateRange" style="width: 100%" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const replays = ref<any[]>([])
const loading = ref(false)
const total = ref(0)
const page = ref(1)
const pageSize = ref(9)
const filterDeviceId = ref('')
const filterType = ref('')
const filterStatus = ref('')
const dateRange = ref<any[]>([])
const currentReplay = ref<any>(null)
const playerVisible = ref(false)
const playerReplay = ref<any>(null)
const detailVisible = ref(false)
const detailReplay = ref<any>(null)
const createVisible = ref(false)

const isPlaying = ref(false)
const playProgress = ref(0)
const playbackSpeed = ref(1)
const currentTime = ref(0)
const totalTime = ref(3600)
let playTimer: any = null

const visibleEvents = ref<any[]>([])
const timelineEvents = ref<any[]>([])

const createForm = ref({
  name: '',
  device_id: '',
  type: 'device',
  dateRange: [],
})

const statusMap: Record<string, string> = { ready: '就绪', playing: '播放中', paused: '已暂停', finished: '已结束' }
const typeMap: Record<string, string> = { device: '设备事件', ai: 'AI决策', system: '系统日志' }
const eventColorMap: Record<string, string> = { alert: 'red', action: 'blue', sensor: 'green', ai: 'purple', system: 'gray' }

const getStatusName = (s: string) => statusMap[s] || s
const getTypeName = (t: string) => typeMap[t] || t
const getStatusColor = (s: string) => ({ ready: 'green', playing: 'arcoblue', paused: 'orange', finished: 'gray' }[s] || 'gray')
const getEventColor = (t: string) => eventColorMap[t] || 'gray'

const currentTimeStr = computed(() => formatTime(currentTime.value))
const totalTimeStr = computed(() => formatTime(totalTime.value))

const formatTime = (seconds: number) => {
  const h = Math.floor(seconds / 3600)
  const m = Math.floor((seconds % 3600) / 60)
  const s = Math.floor(seconds % 60)
  return `${String(h).padStart(2, '0')}:${String(m).padStart(2, '0')}:${String(s).padStart(2, '0')}`
}

const loadReplays = () => {
  loading.value = true
  setTimeout(() => {
    replays.value = [
      { id: 'R1', name: '2026-03-23 设备事件', device_id: 'pet-001', type: 'device', status: 'ready', start_time: '2026-03-23 08:00', end_time: '2026-03-23 18:00', event_count: 156, duration: '10小时', created_at: '2026-03-23 18:30' },
      { id: 'R2', name: 'AI决策回放', device_id: 'pet-001', type: 'ai', status: 'ready', start_time: '2026-03-22 09:00', end_time: '2026-03-22 21:00', event_count: 89, duration: '12小时', created_at: '2026-03-22 21:30' },
      { id: 'R3', name: '系统日志回放', device_id: 'pet-002', type: 'system', status: 'ready', start_time: '2026-03-21 00:00', end_time: '2026-03-21 23:59', event_count: 234, duration: '24小时', created_at: '2026-03-21 23:59' },
    ]
    total.value = 3
    loading.value = false
  }, 500)
}

const startNewReplay = () => {
  createForm.value = { name: '', device_id: '', type: 'device', dateRange: [] }
  createVisible.value = true
}

const handleCreate = () => {
  const newReplay = {
    id: `R${Date.now()}`,
    name: createForm.value.name || `回放-${Date.now()}`,
    device_id: createForm.value.device_id,
    type: createForm.value.type,
    status: 'ready',
    start_time: createForm.value.dateRange[0] || new Date().toISOString().slice(0, 10),
    end_time: createForm.value.dateRange[1] || new Date().toISOString().slice(0, 10),
    event_count: 0,
    duration: '未知',
    created_at: new Date().toISOString().slice(0, 19).replace('T', ' '),
  }
  replays.value.unshift(newReplay)
  total.value++
  Message.success('回放任务已创建')
  createVisible.value = false
}

const playReplay = (replay: any) => {
  currentReplay.value = replay
  playerReplay.value = replay
  replay.status = 'playing'
  isPlaying.value = true
  playProgress.value = 0
  currentTime.value = 0
  totalTime.value = 3600
  playerVisible.value = true

  visibleEvents.value = [
    { id: 'E1', time: '00:00:10', type: 'sensor', message: '心率数据: 72 bpm' },
    { id: 'E2', time: '00:00:25', type: 'alert', message: '体温异常告警触发' },
    { id: 'E3', time: '00:01:05', type: 'ai', message: 'AI决策: 建议检查' },
    { id: 'E4', time: '00:01:30', type: 'action', message: '发送通知给主人' },
    { id: 'E5', time: '00:02:00', type: 'system', message: '告警已确认' },
  ]

  timelineEvents.value = [
    { id: 'E1', type: 'sensor', position: 10, time: '00:00:10' },
    { id: 'E2', type: 'alert', position: 25, time: '00:00:25' },
    { id: 'E3', type: 'ai', position: 65, time: '00:01:05' },
    { id: 'E4', type: 'action', position: 90, time: '00:01:30' },
    { id: 'E5', type: 'system', position: 100, time: '00:02:00' },
  ]

  startPlayback()
}

const startPlayback = () => {
  stopPlayback()
  playTimer = setInterval(() => {
    currentTime.value += playbackSpeed.value
    playProgress.value = (currentTime.value / totalTime.value) * 100
    if (currentTime.value >= totalTime.value) {
      isPlaying.value = false
      if (currentReplay.value) currentReplay.value.status = 'finished'
      stopPlayback()
    }
  }, 1000)
}

const stopPlayback = () => {
  if (playTimer) {
    clearInterval(playTimer)
    playTimer = null
  }
}

const togglePlay = () => {
  if (isPlaying.value) {
    isPlaying.value = false
    if (currentReplay.value) currentReplay.value.status = 'paused'
    stopPlayback()
  } else {
    isPlaying.value = true
    if (currentReplay.value) currentReplay.value.status = 'playing'
    startPlayback()
  }
}

const pauseReplay = () => {
  isPlaying.value = false
  if (currentReplay.value) currentReplay.value.status = 'paused'
  stopPlayback()
}

const skipForward = () => {
  currentTime.value = Math.min(currentTime.value + 30, totalTime.value)
  playProgress.value = (currentTime.value / totalTime.value) * 100
}

const skipBackward = () => {
  currentTime.value = Math.max(currentTime.value - 30, 0)
  playProgress.value = (currentTime.value / totalTime.value) * 100
}

const setPlaybackSpeed = (speed: number) => {
  playbackSpeed.value = speed
}

const seekTimeline = (e: MouseEvent) => {
  const track = e.currentTarget as HTMLElement
  const rect = track.getBoundingClientRect()
  const ratio = (e.clientX - rect.left) / rect.width
  currentTime.value = Math.floor(ratio * totalTime.value)
  playProgress.value = ratio * 100
}

const jumpToEvent = (evt: any) => {
  currentTime.value = Math.floor((evt.position / 100) * totalTime.value)
  playProgress.value = evt.position
}

const closePlayer = () => {
  stopPlayback()
  isPlaying.value = false
  if (currentReplay.value) currentReplay.value.status = 'finished'
  playerVisible.value = false
}

const viewReplayDetail = (replay: any) => {
  detailReplay.value = replay
  detailVisible.value = true
}

const onPageChange = (p: number) => {
  page.value = p
  loadReplays()
}

onMounted(() => {
  loadReplays()
})

onUnmounted(() => {
  stopPlayback()
})
</script>

<style scoped lang="less">
.pro-page-container {
  padding: 16px;
}
.pro-breadcrumb {
  margin-bottom: 12px;
}
.pro-filter-bar {
  margin-bottom: 12px;
}
.pro-action-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}
.pro-content-area {
  margin-bottom: 16px;
}
.replay-card {
  transition: all 0.3s;
  &.active {
    border-color: rgb(var(--primary-6));
    box-shadow: 0 0 0 2px rgb(var(--primary-1));
  }
  .replay-title {
    display: flex;
    align-items: center;
    gap: 8px;
    font-weight: 600;
  }
}
.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
}
.player-wrapper {
  .timeline-bar {
    margin-bottom: 16px;
    .timeline-track {
      position: relative;
      height: 8px;
      background: var(--color-fill-2);
      border-radius: 4px;
      cursor: pointer;
      margin-bottom: 8px;
    }
    .timeline-progress {
      height: 100%;
      background: linear-gradient(90deg, rgb(var(--primary-5)), rgb(var(--primary-6)));
      border-radius: 4px;
      transition: width 0.3s;
    }
    .timeline-marker {
      position: absolute;
      top: 50%;
      transform: translate(-50%, -50%);
      width: 14px;
      height: 14px;
      background: rgb(var(--primary-6));
      border-radius: 50%;
      border: 2px solid white;
      box-shadow: 0 2px 4px rgba(0,0,0,0.2);
    }
    .timeline-event {
      position: absolute;
      top: -4px;
      width: 8px;
      height: 16px;
      border-radius: 2px;
      transform: translateX(-50%);
      cursor: pointer;
      &.alert { background: rgb(var(--danger-6)); }
      &.action { background: rgb(var(--primary-6)); }
      &.sensor { background: rgb(var(--success-6)); }
      &.ai { background: rgb(var(--arcoblue-6)); }
      &.system { background: var(--color-fill-3); }
    }
    .timeline-time {
      display: flex;
      justify-content: space-between;
      font-size: 12px;
      color: var(--color-text-3);
    }
  }
  .player-controls {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;
    padding: 8px;
    background: var(--color-fill-1);
    border-radius: 8px;
    .speed-label {
      font-size: 13px;
      color: var(--color-text-3);
    }
  }
  .player-content {
    display: flex;
    gap: 16px;
    height: 360px;
    .player-main {
      flex: 1;
      .replay-visual {
        height: 100%;
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        gap: 12px;
        background: var(--color-fill-1);
        border-radius: 8px;
        color: var(--color-text-3);
        .replay-device {
          font-size: 12px;
          color: var(--color-text-3);
        }
      }
    }
    .event-log {
      width: 320px;
      border: 1px solid var(--color-border);
      border-radius: 8px;
      padding: 12px;
      overflow: hidden;
      .log-title {
        font-weight: 600;
        font-size: 14px;
        margin-bottom: 12px;
        padding-bottom: 8px;
        border-bottom: 1px solid var(--color-border);
      }
      .log-list {
        overflow-y: auto;
        max-height: 280px;
        .log-item {
          display: flex;
          align-items: center;
          gap: 8px;
          padding: 6px 0;
          font-size: 12px;
          border-bottom: 1px solid var(--color-fill-1);
          &:last-child { border-bottom: none; }
          .log-time { color: var(--color-text-3); min-width: 60px; }
          .log-msg { flex: 1; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
        }
      }
    }
  }
}
</style>
