<template>
  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>具身智能</a-breadcrumb-item>
      <a-breadcrumb-item>
        <router-link :to="`/embodied/${deviceId}/perception`">设备 {{ deviceId }}</router-link>
      </a-breadcrumb-item>
      <a-breadcrumb-item>导航控制</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 导航控制区 -->
    <a-row :gutter="16">
      <!-- 地图与路径显示 -->
      <a-col :xs="24" :sm="24" :md="16">
        <a-card title="导航地图" class="nav-map-card">
          <template #extra>
            <a-tag :color="navStatusColor">{{ navStatusText }}</a-tag>
          </template>
          <div class="nav-canvas-wrapper">
            <canvas ref="navCanvasRef" class="nav-canvas" />
            <!-- 点击下发目标 -->
            <div class="canvas-overlay" @click="handleCanvasClick" />
          </div>
          <div class="nav-path-info" v-if="currentPath.length">
            <a-descriptions :column="3" size="small">
              <a-descriptions-item label="路径点数">{{ currentPath.length }}</a-descriptions-item>
              <a-descriptions-item label="总长度">
                {{ pathLength.toFixed(2) }} m
              </a-descriptions-item>
              <a-descriptions-item label="预计时间">
                {{ Math.round(pathLength / 0.5) }} s
              </a-descriptions-item>
            </a-descriptions>
          </div>
        </a-card>

        <!-- 路径历史 -->
        <a-card title="导航历史" style="margin-top: 16px">
          <a-table
            :columns="historyColumns"
            :data="navHistory"
            :loading="loadingHistory"
            :pagination="{ pageSize: 5, showTotal: true }"
            row-key="id"
            size="small"
          >
            <template #status="{ record }">
              <a-tag :color="getNavStatusColor(record.status)">
                {{ getNavStatusText(record.status) }}
              </a-tag>
            </template>
            <template #start_time="{ record }">
              {{ formatTime(record.started_at || record.start_time) }}
            </template>
          </a-table>
        </a-card>
      </a-col>

      <!-- 控制面板 -->
      <a-col :xs="24" :sm="24" :md="8">
        <!-- 紧急停止 -->
        <a-card class="emergency-card">
          <a-button
            type="primary"
            danger
            long
            size="large"
            :loading="stopping"
            @click="handleEmergencyStop"
          >
            <template #icon><icon-stop /></template>
            紧急停止 (E-STOP)
          </a-button>
        </a-card>

        <!-- 目标点设置 -->
        <a-card title="目标点设置" style="margin-top: 16px">
          <a-form :model="targetForm" layout="vertical">
            <a-form-item label="目标 X (m)">
              <a-input-number v-model="targetForm.x" :min="0" :precision="3" style="width: 100%" />
            </a-form-item>
            <a-form-item label="目标 Y (m)">
              <a-input-number v-model="targetForm.y" :min="0" :precision="3" style="width: 100%" />
            </a-form-item>
            <a-form-item label="预设位置">
              <a-select v-model="presetSelect" placeholder="选择预设位置" allow-clear @change="applyPreset">
                <a-option value="charging">充电桩</a-option>
                <a-option value="home">Home Base</a-option>
                <a-option value="kitchen">厨房</a-option>
                <a-option value="living">客厅</a-option>
                <a-option value="bedroom">卧室</a-option>
              </a-select>
            </a-form-item>
            <a-button type="primary" long @click="startNavigation" :loading="navigating" :disabled="!targetForm.x || !targetForm.y">
              <template #icon><icon-send /></template>
              开始导航
            </a-button>
          </a-form>
        </a-card>

        <!-- 跟随模式 -->
        <a-card title="跟随模式" style="margin-top: 16px">
          <a-form :model="followForm" layout="vertical">
            <a-form-item label="跟随目标">
              <a-select v-model="followForm.target_id" placeholder="选择跟随目标">
                <a-option value="user-1">用户 1</a-option>
                <a-option value="user-2">用户 2</a-option>
                <a-option value="user-3">用户 3</a-option>
              </a-select>
            </a-form-item>
            <a-form-item label="跟随距离 (m)">
              <a-slider v-model="followForm.distance" :min="0.5" :max="3" :step="0.1" :show-ticks="true" />
              <span class="text-muted">{{ followForm.distance }} m</span>
            </a-form-item>
            <a-button type="primary" long @click="startFollow" :disabled="!followForm.target_id">
              <template #icon><icon-user /></template>
              启动跟随
            </a-button>
            <a-button type="secondary" long style="margin-top: 8px" @click="stopFollow" :disabled="!isFollowing">
              <template #icon><icon-close /></template>
              停止跟随
            </a-button>
          </a-form>
        </a-card>

        <!-- 自主探索 -->
        <a-card title="自主探索" style="margin-top: 16px">
          <a-space direction="vertical" style="width: 100%">
            <a-descriptions :column="1" size="small">
              <a-descriptions-item label="探索状态">
                <a-tag :color="exploreStatus === 'exploring' ? 'green' : 'gray'">
                  {{ exploreStatus === 'exploring' ? '探索中' : '空闲' }}
                </a-tag>
              </a-descriptions-item>
              <a-descriptions-item label="已探索区域">
                <a-progress
                  :percent="exploreProgress"
                  :color="exploreProgress > 80 ? 'green' : 'arcoblue'"
                  size="small"
                />
              </a-descriptions-item>
            </a-descriptions>
            <a-button
              v-if="exploreStatus !== 'exploring'"
              type="primary"
              long
              @click="startExplore"
            >
              <template #icon><icon-compass /></template>
              启动探索
            </a-button>
            <a-button v-else type="warning" long @click="stopExplore">
              <template #icon><icon-stop /></template>
              停止探索
            </a-button>
          </a-space>
        </a-card>
      </a-col>
    </a-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed, nextTick } from 'vue'
import { useRoute } from 'vue-router'
import {
  navigateTo, stopNavigation, startFollow, stopFollow,
  getExploreStatus, startExplore, getNavigationHistory
} from '@/api/embodied'
import { Message, Modal } from '@arco-design/web-vue'

const route = useRoute()
const deviceId = ref(route.params.device_id as string)

const navCanvasRef = ref<HTMLCanvasElement>()
const navigating = ref(false)
const stopping = ref(false)
const isFollowing = ref(false)
const exploreStatus = ref<'idle' | 'exploring'>('idle')
const exploreProgress = ref(0)
const loadingHistory = ref(false)
const navHistory = ref<any[]>([])

const navStatus = ref<'idle' | 'navigating' | 'arrived' | 'failed'>('idle')
const currentPath = ref<Array<{ x: number; y: number }>>([])

const targetForm = ref({ x: 0, y: 0 })
const followForm = ref({ target_id: '', distance: 1.0 })
const presetSelect = ref('')

const navStatusColor = computed(() => {
  const map: Record<string, string> = { idle: 'gray', navigating: 'arcoblue', arrived: 'green', failed: 'red' }
  return map[navStatus.value] || 'gray'
})
const navStatusText = computed(() => {
  const map: Record<string, string> = { idle: '空闲', navigating: '导航中', arrived: '已到达', failed: '失败' }
  return map[navStatus.value] || '未知'
})
const pathLength = computed(() => {
  if (currentPath.value.length < 2) return 0
  let len = 0
  for (let i = 1; i < currentPath.value.length; i++) {
    const dx = currentPath.value[i].x - currentPath.value[i - 1].x
    const dy = currentPath.value[i].y - currentPath.value[i - 1].y
    len += Math.sqrt(dx * dx + dy * dy)
  }
  return len
})

const historyColumns = [
  { title: 'ID', dataIndex: 'id', width: 70 },
  { title: '状态', dataIndex: 'status', slotName: 'status', width: 100 },
  { title: '起点', dataIndex: 'start_pos', width: 120 },
  { title: '终点', dataIndex: 'target_pos', width: 120 },
  { title: '时间', dataIndex: 'start_time', slotName: 'start_time', width: 160 }
]

async function startNavigation() {
  try {
    navigating.value = true
    const res = await navigateTo(deviceId.value, { target_x: targetForm.value.x, target_y: targetForm.value.y })
    navStatus.value = 'navigating'
    if (res.data?.path) currentPath.value = res.data.path
    Message.success('导航已开始')
    drawNavigation()
  } catch (err: any) {
    Message.error('导航失败: ' + err.message)
    navStatus.value = 'failed'
  } finally {
    navigating.value = false
  }
}

async function handleEmergencyStop() {
  Modal.warning({
    title: '确认紧急停止',
    content: '确定要执行紧急停止吗？这将立即停止设备所有运动。',
    onOk: async () => {
      try {
        stopping.value = true
        await stopNavigation(deviceId.value)
        await stopFollow(deviceId.value)
        navStatus.value = 'idle'
        currentPath.value = []
        Message.warning('已执行紧急停止')
        drawNavigation()
      } catch (err: any) {
        Message.error('停止失败: ' + err.message)
      } finally {
        stopping.value = false
      }
    }
  })
}

async function stopFollow() {
  try {
    await stopFollow(deviceId.value)
    isFollowing.value = false
    Message.success('已停止跟随')
  } catch (err: any) {
    Message.error('停止跟随失败: ' + err.message)
  }
}

async function startFollowAction() {
  try {
    await startFollow(deviceId.value, { target_id: followForm.value.target_id, distance: followForm.value.distance })
    isFollowing.value = true
    Message.success('跟随模式已启动')
  } catch (err: any) {
    Message.error('跟随启动失败: ' + err.message)
  }
}

function applyPreset(val: string) {
  const presets: Record<string, { x: number; y: number }> = {
    charging: { x: 1.0, y: 1.0 },
    home: { x: 0.0, y: 0.0 },
    kitchen: { x: 3.0, y: 2.0 },
    living: { x: 5.0, y: 3.0 },
    bedroom: { x: 2.0, y: 5.0 }
  }
  if (presets[val]) {
    targetForm.value.x = presets[val].x
    targetForm.value.y = presets[val].y
  }
}

function handleCanvasClick(e: MouseEvent) {
  const canvas = navCanvasRef.value
  if (!canvas) return
  const rect = canvas.getBoundingClientRect()
  const scaleX = 10 / canvas.width
  const scaleY = 10 / canvas.height
  targetForm.value.x = parseFloat(((e.clientX - rect.left) * scaleX).toFixed(3))
  targetForm.value.y = parseFloat(((e.clientY - rect.top) * scaleY).toFixed(3))
  Message.info(`已设置目标点: (${targetForm.value.x}, ${targetForm.value.y})`)
}

async function startExplore() {
  try {
    await startExplore(deviceId.value)
    exploreStatus.value = 'exploring'
    Message.success('探索已启动')
  } catch (err: any) {
    Message.error('启动探索失败: ' + err.message)
  }
}

async function stopExplore() {
  try {
    await stopNavigation(deviceId.value)
    exploreStatus.value = 'idle'
    Message.success('探索已停止')
  } catch (err: any) {
    Message.error('停止探索失败: ' + err.message)
  }
}

async function loadExploreStatus() {
  try {
    const res = await getExploreStatus(deviceId.value)
    exploreStatus.value = res.data?.status || 'idle'
    exploreProgress.value = res.data?.explored_ratio ? Math.round(res.data.explored_ratio * 100) : 0
  } catch {}
}

async function loadHistory() {
  try {
    loadingHistory.value = true
    const res = await getNavigationHistory(deviceId.value)
    navHistory.value = (res.data?.history || res.data || []).map((h: any) => ({
      ...h,
      start_pos: h.start ? `(${h.start.x?.toFixed(1)}, ${h.start.y?.toFixed(1)})` : '-',
      target_pos: h.target ? `(${h.target.x?.toFixed(1)}, ${h.target.y?.toFixed(1)})` : '-'
    }))
  } catch {} finally {
    loadingHistory.value = false
  }
}

function drawNavigation() {
  nextTick(() => {
    const canvas = navCanvasRef.value
    if (!canvas) return
    const ctx = canvas.getContext('2d')
    if (!ctx) return
    canvas.width = canvas.offsetWidth || 600
    canvas.height = 400

    ctx.fillStyle = '#f0f2f5'
    ctx.fillRect(0, 0, canvas.width, canvas.height)

    const gridSize = 10
    const cellW = canvas.width / gridSize
    const cellH = canvas.height / gridSize

    for (let i = 0; i <= gridSize; i++) {
      ctx.strokeStyle = '#e0e0e0'
      ctx.beginPath()
      ctx.moveTo(i * cellW, 0)
      ctx.lineTo(i * cellW, canvas.height)
      ctx.stroke()
      ctx.beginPath()
      ctx.moveTo(0, i * cellH)
      ctx.lineTo(canvas.width, i * cellH)
      ctx.stroke()
    }

    // 绘制路径
    if (currentPath.value.length > 1) {
      ctx.beginPath()
      ctx.moveTo(currentPath.value[0].x / gridSize * canvas.width, currentPath.value[0].y / gridSize * canvas.height)
      for (let i = 1; i < currentPath.value.length; i++) {
        ctx.lineTo(currentPath.value[i].x / gridSize * canvas.width, currentPath.value[i].y / gridSize * canvas.height)
      }
      ctx.strokeStyle = '#1650d8'
      ctx.lineWidth = 3
      ctx.stroke()
    }

    // 绘制目标点
    if (targetForm.value.x && targetForm.value.y) {
      ctx.beginPath()
      ctx.arc(targetForm.value.x / gridSize * canvas.width, targetForm.value.y / gridSize * canvas.height, 8, 0, Math.PI * 2)
      ctx.fillStyle = '#ff4d4f'
      ctx.fill()
      ctx.strokeStyle = '#fff'
      ctx.lineWidth = 2
      ctx.stroke()
    }

    // 当前位置
    ctx.beginPath()
    ctx.arc(canvas.width / 2, canvas.height / 2, 10, 0, Math.PI * 2)
    ctx.fillStyle = '#1650d8'
    ctx.fill()
    ctx.strokeStyle = '#fff'
    ctx.lineWidth = 2
    ctx.stroke()
  })
}

function getNavStatusColor(s: string) {
  const map: Record<string, string> = { arrived: 'green', navigating: 'arcoblue', pending: 'gray', failed: 'red' }
  return map[s] || 'default'
}
function getNavStatusText(s: string) {
  const map: Record<string, string> = { arrived: '已到达', navigating: '导航中', pending: '等待中', failed: '失败' }
  return map[s] || s
}
function formatTime(t: string) {
  if (!t) return '-'
  return new Date(t).toLocaleString('zh-CN')
}

onMounted(async () => {
  await loadExploreStatus()
  await loadHistory()
  drawNavigation()
})
</script>

<style scoped>
.nav-map-card {
  height: 100%;
}
.nav-canvas-wrapper {
  position: relative;
  width: 100%;
  height: 400px;
}
.nav-canvas {
  width: 100%;
  height: 100%;
  display: block;
  border: 1px solid var(--color-border);
  border-radius: 4px;
}
.canvas-overlay {
  position: absolute;
  inset: 0;
  cursor: crosshair;
  border-radius: 4px;
}
.nav-path-info {
  margin-top: 12px;
  padding: 8px;
  background: var(--color-fill-1);
  border-radius: 4px;
}
.emergency-card {
  border: 2px solid #ff4d4f;
  border-radius: 8px;
}
.text-muted {
  color: var(--color-text-3);
  font-size: 12px;
}
</style>
