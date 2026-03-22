<template>
  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>具身智能</a-breadcrumb-item>
      <a-breadcrumb-item>环境感知</a-breadcrumb-item>
      <a-breadcrumb-item>{{ deviceId }}</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 实时感知数据展示 -->
    <a-row :gutter="16" class="perception-cards">
      <!-- 视觉感知 -->
      <a-col :xs="24" :sm="24" :md="8">
        <a-card class="perception-card" title="视觉感知">
          <template #extra>
            <a-tag :color="visualStatus === 'active' ? 'green' : 'gray'">
              {{ visualStatus === 'active' ? '实时' : '离线' }}
            </a-tag>
          </template>
          <a-descriptions :column="1" size="small">
            <a-descriptions-item label="场景分类">
              <a-tag>{{ perceptionData.visual?.scene || '未知' }}</a-tag>
            </a-descriptions-item>
            <a-descriptions-item label="手势识别">
              {{ perceptionData.visual?.gesture || '无' }}
            </a-descriptions-item>
            <a-descriptions-item label="人体姿态">
              <a-tag v-if="perceptionData.visual?.human_pose">检测到</a-tag>
              <span v-else class="text-muted">未检测到</span>
            </a-descriptions-item>
          </a-descriptions>
          <a-divider>识别物体</a-divider>
          <a-tag v-for="obj in perceptionData.visual?.objects" :key="obj.label" class="object-tag">
            {{ obj.label }} ({{ (obj.confidence * 100).toFixed(0) }}%)
          </a-tag>
          <a-empty v-if="!perceptionData.visual?.objects?.length" description="暂无数据" />
        </a-card>
      </a-col>

      <!-- 深度感知 -->
      <a-col :xs="24" :sm="24" :md="8">
        <a-card class="perception-card" title="深度感知">
          <template #extra>
            <a-tag :color="depthStatus === 'active' ? 'green' : 'gray'">
              {{ depthStatus === 'active' ? '实时' : '离线' }}
            </a-tag>
          </template>
          <a-descriptions :column="1" size="small">
            <a-descriptions-item label="置信度">
              <a-progress
                :percent="(perceptionData.depth?.confidence * 100) || 0"
                :color="getConfidenceColor(perceptionData.depth?.confidence)"
                size="small"
              />
            </a-descriptions-item>
          </a-descriptions>
          <a-divider>障碍物检测</a-divider>
          <a-list size="small" :data-source="perceptionData.depth?.obstacles || []">
            <template #renderItem="{ item }">
              <a-list-item>
                <a-list-item-meta>
                  <template #title>
                    {{ item.direction }}方向障碍物
                  </template>
                  <template #description>
                    距离: {{ item.distance.toFixed(2) }}m
                  </template>
                </a-list-item-meta>
              </a-list-item>
            </template>
          </a-list>
          <a-empty v-if="!perceptionData.depth?.obstacles?.length" description="暂无障碍物" />
        </a-card>
      </a-col>

      <!-- 触觉感知 -->
      <a-col :xs="24" :sm="24" :md="8">
        <a-card class="perception-card" title="触觉感知">
          <template #extra>
            <a-tag :color="perceptionData.touch?.touched ? 'green' : 'gray'">
              {{ perceptionData.touch?.touched ? '已触摸' : '未触摸' }}
            </a-tag>
          </template>
          <a-descriptions :column="1" size="small">
            <a-descriptions-item label="触摸位置">
              {{ perceptionData.touch?.position || '无' }}
            </a-descriptions-item>
            <a-descriptions-item label="触摸力度">
              {{ perceptionData.touch?.force?.toFixed(2) || '0' }} N
            </a-descriptions-item>
          </a-descriptions>
          <a-divider>触摸状态可视化</a-divider>
          <div class="touch-viz">
            <div
              class="touch-indicator"
              :class="{ active: perceptionData.touch?.touched }"
            >
              <icon-robot />
            </div>
          </div>
        </a-card>
      </a-col>
    </a-row>

    <!-- 感知事件记录 -->
    <div class="pro-content-area" style="margin-top: 16px">
      <a-card title="感知事件记录">
        <template #extra>
          <a-space>
            <a-button type="primary" @click="loadPerception">刷新</a-button>
          </a-space>
        </template>
        <a-table
          :columns="columns"
          :data="events"
          :loading="loading"
          :pagination="{ pageSize: 10, showTotal: true, showPageSize: true, pageSizeOptions: [10, 20, 50, 100] }"
          row-key="id"
        >
          <template #event_type="{ record }">
            <a-tag :color="getEventTypeColor(record.event_type)">
              {{ getEventTypeText(record.event_type) }}
            </a-tag>
          </template>
          <template #description="{ record }">
            <a-tooltip>
              <template #title>{{ record.description }}</template>
              <span class="text-ellipsis">{{ record.description }}</span>
            </a-tooltip>
          </template>
          <template #created_at="{ record }">
            {{ formatTime(record.created_at) }}
          </template>
        </a-table>
      </a-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import { getPerception } from '@/api/embodied'
import { Message } from '@arco-design/web-vue'

const route = useRoute()
const deviceId = ref(route.params.device_id as string)

const loading = ref(false)
const events = ref<any[]>([])
const perceptionData = ref<any>({
  visual: { objects: [], scene: '', gesture: '', human_pose: null },
  depth: { obstacles: [], confidence: 0 },
  touch: { touched: false }
})
const visualStatus = ref<'active' | 'inactive'>('inactive')
const depthStatus = ref<'active' | 'inactive'>('inactive')

const columns = [
  { title: '序号', dataIndex: 'id', width: 80 },
  { title: '事件类型', dataIndex: 'event_type', slotName: 'event_type', width: 120 },
  { title: '描述', dataIndex: 'description', slotName: 'description' },
  { title: '时间', dataIndex: 'created_at', slotName: 'created_at', width: 180 }
]

let pollingTimer: number | null = null

async function loadPerception() {
  try {
    loading.value = true
    const res = await getPerception(deviceId.value, { page_size: 20 })
    if (res.data) {
      perceptionData.value = res.data.current || perceptionData.value
      events.value = res.data.events || []
      visualStatus.value = res.data.visual_status || 'inactive'
      depthStatus.value = res.data.depth_status || 'inactive'
    }
  } catch (err: any) {
    Message.error('加载感知数据失败: ' + err.message)
  } finally {
    loading.value = false
  }
}

function getConfidenceColor(confidence: number) {
  if (confidence >= 0.8) return 'green'
  if (confidence >= 0.5) return 'orange'
  return 'red'
}

function getEventTypeColor(type: string) {
  const map: Record<string, string> = {
    object_detected: 'blue',
    scene_changed: 'cyan',
    human_detected: 'green',
    obstacle_detected: 'orange',
    touch_detected: 'purple'
  }
  return map[type] || 'default'
}

function getEventTypeText(type: string) {
  const map: Record<string, string> = {
    object_detected: '物体检测',
    scene_changed: '场景变化',
    human_detected: '人体检测',
    obstacle_detected: '障碍物检测',
    touch_detected: '触摸检测'
  }
  return map[type] || type
}

function formatTime(time: string) {
  if (!time) return '-'
  return new Date(time).toLocaleString('zh-CN')
}

onMounted(() => {
  loadPerception()
  // 轮询刷新，每3秒
  pollingTimer = window.setInterval(loadPerception, 3000)
})

onUnmounted(() => {
  if (pollingTimer) {
    clearInterval(pollingTimer)
  }
})
</script>

<style scoped>
.perception-card {
  height: 100%;
}
.object-tag {
  margin: 4px 4px 4px 0;
}
.text-muted {
  color: var(--color-text-3);
}
.text-ellipsis {
  display: inline-block;
  max-width: 300px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.touch-viz {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 120px;
}
.touch-indicator {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  background: var(--color-fill-2);
  display: flex;
  justify-content: center;
  align-items: center;
  font-size: 36px;
  color: var(--color-text-3);
  transition: all 0.3s;
}
.touch-indicator.active {
  background: var(--color-primary-light-1);
  color: var(--color-primary);
  box-shadow: 0 0 16px var(--color-primary-light-2);
}
</style>
