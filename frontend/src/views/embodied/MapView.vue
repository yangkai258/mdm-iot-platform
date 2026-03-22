<template>
  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>具身智能</a-breadcrumb-item>
      <a-breadcrumb-item>
        <router-link :to="`/embodied/${deviceId}/perception`">设备 {{ deviceId }}</router-link>
      </a-breadcrumb-item>
      <a-breadcrumb-item>地图管理</a-breadcrumb-item>
    </a-breadcrumb>

    <a-row :gutter="16" class="map-main">
      <!-- 地图列表 -->
      <a-col :xs="24" :sm="24" :md="8">
        <a-card title="地图列表" class="map-list-card">
          <template #extra>
            <a-button type="primary" size="small" @click="openCreateModal">
              <template #icon><icon-plus /></template>
              新建
            </a-button>
          </template>
          <a-list :data-source="maps" :loading="loading" size="small">
            <template #renderItem="{ item }">
              <a-list-item
                :class="{ 'active-map-item': item.id === activeMapId }"
                @click="selectMap(item)"
              >
                <a-list-item-meta>
                  <template #title>
                    <a-space>
                      {{ item.map_name || `地图 ${item.id}` }}
                      <a-tag v-if="item.is_active" color="green" size="small">激活</a-tag>
                      <a-tag v-if="item.map_type === 'grid'" size="small">栅格</a-tag>
                      <a-tag v-if="item.map_type === 'semantic'" size="small">语义</a-tag>
                    </a-space>
                  </template>
                  <template #description>
                    <span class="text-muted">分辨率: {{ item.resolution }}m · 版本: v{{ item.version }}</span>
                  </template>
                </a-list-item-meta>
                <template #actions>
                  <a-button type="text" size="small" @click.stop="editMap(item)">
                    <icon-edit />
                  </a-button>
                  <a-button type="text" size="small" @click.stop="setActiveMap(item)">
                    <icon-check-circle />
                  </a-button>
                </template>
              </a-list-item>
            </template>
          </a-list>
          <a-empty v-if="!loading && !maps.length" description="暂无地图" />
        </a-card>

        <!-- 当前位置 -->
        <a-card title="当前位置" class="localization-card" style="margin-top: 16px">
          <a-descriptions :column="1" size="small">
            <a-descriptions-item label="X 坐标">{{ localization?.x?.toFixed(3) || '-' }}</a-descriptions-item>
            <a-descriptions-item label="Y 坐标">{{ localization?.y?.toFixed(3) || '-' }}</a-descriptions-item>
            <a-descriptions-item label="朝向">{{ localization?.theta !== undefined ? `${(localization.theta * 180 / Math.PI).toFixed(1)}°` : '-' }}</a-descriptions-item>
            <a-descriptions-item label="置信度">
              <a-progress
                v-if="localization?.confidence"
                :percent="Math.round(localization.confidence * 100)"
                :color="getConfidenceColor(localization.confidence)"
                size="small"
              />
              <span v-else>-</span>
            </a-descriptions-item>
          </a-descriptions>
          <a-button long type="primary" @click="calibrate" style="margin-top: 12px">
            <template #icon><icon-send /></template>
            重新标定
          </a-button>
        </a-card>
      </a-col>

      <!-- 地图可视化 -->
      <a-col :xs="24" :sm="24" :md="16">
        <a-card title="地图可视化" class="map-canvas-card">
          <template #extra>
            <a-space>
              <a-select v-model="displayType" size="small" style="width: 120px">
                <a-option value="grid">栅格图</a-option>
                <a-option value="semantic">语义图</a-option>
                <a-option value="fusion">融合图</a-option>
              </a-select>
              <a-button size="small" @click="refreshMap">刷新</a-button>
            </a-space>
          </template>
          <div class="map-canvas-container" ref="mapCanvasRef">
            <canvas ref="canvasRef" class="map-canvas" />
            <div v-if="!activeMap" class="map-empty-overlay">
              <a-empty description="请从左侧选择地图" />
            </div>
          </div>
        </a-card>

        <!-- 位置标注 -->
        <a-card title="位置标注" class="annotation-card" style="margin-top: 16px">
          <template #extra>
            <a-button type="primary" size="small" @click="openAnnotationModal">
              <template #icon><icon-plus /></template>
              添加标注
            </a-button>
          </template>
          <a-table
            :columns="annotationColumns"
            :data="annotations"
            :loading="loadingAnnotations"
            :pagination="{ pageSize: 5, showTotal: true }"
            row-key="id"
            size="small"
          >
            <template #annotation_type="{ record }">
              <a-tag :color="getAnnotationColor(record.annotation_type)">
                {{ getAnnotationText(record.annotation_type) }}
              </a-tag>
            </template>
            <template #position="{ record }">
              ({{ record.position?.x?.toFixed(2) }}, {{ record.position?.y?.toFixed(2) }})
            </template>
            <template #actions="{ record }">
              <a-button type="text" size="small" @click="navigateToAnnotation(record)">
                <icon-location />
              </a-button>
              <a-button type="text" size="small" @click="deleteAnnotation(record)">
                <icon-delete />
              </a-button>
            </template>
          </a-table>
        </a-card>
      </a-col>
    </a-row>

    <!-- 创建/编辑地图弹窗 -->
    <a-modal
      v-model:visible="mapModalVisible"
      :title="editingMap ? '编辑地图' : '创建地图'"
      @before-ok="submitMap"
      @cancel="mapModalVisible = false"
    >
      <a-form :model="mapForm" layout="vertical">
        <a-form-item label="地图名称" required>
          <a-input v-model="mapForm.map_name" placeholder="请输入地图名称" />
        </a-form-item>
        <a-form-item label="地图类型" required>
          <a-select v-model="mapForm.map_type" placeholder="选择类型">
            <a-option value="grid">栅格地图</a-option>
            <a-option value="semantic">语义地图</a-option>
            <a-option value="topological">拓扑地图</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="分辨率 (m)">
          <a-input-number v-model="mapForm.resolution" :min="0.01" :max="1" :step="0.01" style="width: 100%" />
        </a-form-item>
        <a-form-item label="描述">
          <a-textarea v-model="mapForm.description" :rows="3" placeholder="地图描述..." />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 添加标注弹窗 -->
    <a-modal
      v-model:visible="annotationModalVisible"
      title="添加位置标注"
      @before-ok="submitAnnotation"
      @cancel="annotationModalVisible = false"
    >
      <a-form :model="annotationForm" layout="vertical">
        <a-form-item label="标注类型" required>
          <a-select v-model="annotationForm.annotation_type" placeholder="选择类型">
            <a-option value="charging_station">充电站</a-option>
            <a-option value="home_base">Home Base</a-option>
            <a-option value="waypoint">路径点</a-option>
            <a-option value="landmark">地标</a-option>
            <a-option value="forbidden_zone">禁区</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="标注名称" required>
          <a-input v-model="annotationForm.label" placeholder="标注名称" />
        </a-form-item>
        <a-form-item label="X 坐标">
          <a-input-number v-model="annotationForm.position_x" :min="0" style="width: 100%" />
        </a-form-item>
        <a-form-item label="Y 坐标">
          <a-input-number v-model="annotationForm.position_y" :min="0" style="width: 100%" />
        </a-form-item>
        <a-form-item label="描述">
          <a-textarea v-model="annotationForm.description" :rows="2" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, nextTick } from 'vue'
import { useRoute } from 'vue-router'
import { getMaps, updateMap, getLocalization, calibrateLocalization } from '@/api/embodied'
import { Message, Modal } from '@arco-design/web-vue'

const route = useRoute()
const deviceId = ref(route.params.device_id as string)

const loading = ref(false)
const loadingAnnotations = ref(false)
const maps = ref<any[]>([])
const activeMapId = ref<number | null>(null)
const activeMap = ref<any>(null)
const localization = ref<any>(null)
const annotations = ref<any[]>([])
const displayType = ref('grid')

const mapCanvasRef = ref<HTMLElement>()
const canvasRef = ref<HTMLCanvasElement>()

const mapModalVisible = ref(false)
const annotationModalVisible = ref(false)
const editingMap = ref<any>(null)

const mapForm = ref({ map_name: '', map_type: 'grid', resolution: 0.05, description: '' })
const annotationForm = ref({ annotation_type: '', label: '', position_x: 0, position_y: 0, description: '' })

const annotationColumns = [
  { title: '类型', dataIndex: 'annotation_type', slotName: 'annotation_type', width: 120 },
  { title: '名称', dataIndex: 'label', width: 140 },
  { title: '位置', dataIndex: 'position', slotName: 'position', width: 160 },
  { title: '描述', dataIndex: 'description' },
  { title: '操作', slotName: 'actions', width: 80 }
]

async function loadMaps() {
  try {
    loading.value = true
    const res = await getMaps(deviceId.value)
    maps.value = res.data?.maps || res.data || []
    if (maps.value.length && !activeMapId.value) {
      const active = maps.value.find((m: any) => m.is_active)
      if (active) selectMap(active)
    }
  } catch (err: any) {
    Message.error('加载地图失败: ' + err.message)
  } finally {
    loading.value = false
  }
}

async function loadLocalization() {
  try {
    const res = await getLocalization(deviceId.value)
    localization.value = res.data
  } catch {}
}

async function selectMap(item: any) {
  activeMapId.value = item.id
  activeMap.value = item
  annotations.value = item.annotations || []
  await drawMap()
}

async function setActiveMap(item: any) {
  try {
    await updateMap(deviceId.value, { map_id: item.id, action: 'activate' })
    Message.success('已激活地图')
    await loadMaps()
  } catch (err: any) {
    Message.error('激活失败: ' + err.message)
  }
}

function openCreateModal() {
  editingMap.value = null
  mapForm.value = { map_name: '', map_type: 'grid', resolution: 0.05, description: '' }
  mapModalVisible.value = true
}

function editMap(item: any) {
  editingMap.value = item
  mapForm.value = {
    map_name: item.map_name || item.name || '',
    map_type: item.map_type || 'grid',
    resolution: item.resolution || 0.05,
    description: item.description || ''
  }
  mapModalVisible.value = true
}

async function submitMap(done: (val: boolean) => void) {
  try {
    await updateMap(deviceId.value, {
      ...mapForm.value,
      action: editingMap.value ? 'update' : 'create'
    })
    Message.success(editingMap.value ? '更新成功' : '创建成功')
    mapModalVisible.value = false
    await loadMaps()
    done(true)
  } catch (err: any) {
    Message.error('操作失败: ' + err.message)
    done(false)
  }
}

async function calibrate() {
  try {
    await calibrateLocalization(deviceId.value, { action: 'calibrate' })
    Message.success('标定已启动')
    setTimeout(loadLocalization, 1000)
  } catch (err: any) {
    Message.error('标定失败: ' + err.message)
  }
}

async function refreshMap() {
  await loadMaps()
  await loadLocalization()
  await drawMap()
}

async function drawMap() {
  await nextTick()
  if (!canvasRef.value || !activeMap.value) return
  const canvas = canvasRef.value
  const ctx = canvas.getContext('2d')
  if (!ctx) return

  const container = mapCanvasRef.value
  if (container) {
    canvas.width = container.clientWidth
    canvas.height = Math.max(400, container.clientHeight * 0.6)
  }

  ctx.fillStyle = '#f0f2f5'
  ctx.fillRect(0, 0, canvas.width, canvas.height)

  if (!activeMap.value.map_data) {
    ctx.fillStyle = '#999'
    ctx.font = '16px sans-serif'
    ctx.textAlign = 'center'
    ctx.fillText('暂无地图数据', canvas.width / 2, canvas.height / 2)
    return
  }

  // 简单渲染地图数据
  const mapData = activeMap.value.map_data
  const size = mapData.size || { width: 100, height: 100 }
  const cellW = canvas.width / size.width
  const cellH = canvas.height / size.height

  const grid = mapData.grid || []
  for (let y = 0; y < Math.min(grid.length, size.height); y++) {
    for (let x = 0; x < Math.min((grid[y] || []).length, size.width); x++) {
      const val = grid[y][x]
      if (val === 1) {
        ctx.fillStyle = '#1a1a1a'
      } else if (val === 0.5) {
        ctx.fillStyle = '#f0e68c'
      } else {
        ctx.fillStyle = '#ffffff'
      }
      ctx.fillRect(x * cellW, y * cellH, cellW + 0.5, cellH + 0.5)
    }
  }

  // 绘制当前位置
  if (localization.value) {
    const px = (localization.value.x / (size.width * (activeMap.value.resolution || 0.05))) * canvas.width
    const py = (localization.value.y / (size.height * (activeMap.value.resolution || 0.05))) * canvas.height
    ctx.beginPath()
    ctx.arc(px, py, 8, 0, Math.PI * 2)
    ctx.fillStyle = '#1650d8'
    ctx.fill()
    ctx.strokeStyle = '#fff'
    ctx.lineWidth = 2
    ctx.stroke()

    // 朝向箭头
    const angle = localization.value.theta || 0
    ctx.beginPath()
    ctx.moveTo(px, py)
    ctx.lineTo(px + Math.cos(angle) * 20, py + Math.sin(angle) * 20)
    ctx.strokeStyle = '#1650d8'
    ctx.lineWidth = 3
    ctx.stroke()
  }

  // 绘制标注
  if (annotations.value.length) {
    annotations.value.forEach(ann => {
      const ax = (ann.position?.x / (size.width * (activeMap.value.resolution || 0.05))) * canvas.width
      const ay = (ann.position?.y / (size.height * (activeMap.value.resolution || 0.05))) * canvas.height
      ctx.beginPath()
      ctx.arc(ax, ay, 5, 0, Math.PI * 2)
      ctx.fillStyle = getAnnotationColor(ann.annotation_type)
      ctx.fill()
    })
  }
}

function getConfidenceColor(conf: number) {
  if (conf >= 0.8) return 'green'
  if (conf >= 0.5) return 'orange'
  return 'red'
}

function getAnnotationColor(type: string) {
  const map: Record<string, string> = {
    charging_station: 'green',
    home_base: 'blue',
    waypoint: 'cyan',
    landmark: 'orange',
    forbidden_zone: 'red'
  }
  return map[type] || 'default'
}

function getAnnotationText(type: string) {
  const map: Record<string, string> = {
    charging_station: '充电站',
    home_base: 'Home Base',
    waypoint: '路径点',
    landmark: '地标',
    forbidden_zone: '禁区'
  }
  return map[type] || type
}

function openAnnotationModal() {
  annotationForm.value = { annotation_type: '', label: '', position_x: 0, position_y: 0, description: '' }
  annotationModalVisible.value = true
}

async function submitAnnotation(done: (val: boolean) => void) {
  try {
    await updateMap(deviceId.value, {
      action: 'add_annotation',
      annotation_type: annotationForm.value.annotation_type,
      label: annotationForm.value.label,
      position: { x: annotationForm.value.position_x, y: annotationForm.value.position_y },
      description: annotationForm.value.description
    })
    Message.success('标注已添加')
    annotationModalVisible.value = false
    if (activeMap.value) selectMap(activeMap.value)
    done(true)
  } catch (err: any) {
    Message.error('添加失败: ' + err.message)
    done(false)
  }
}

async function navigateToAnnotation(record: any) {
  try {
    const { navigateTo } = await import('@/api/embodied')
    await navigateTo(deviceId.value, { target_x: record.position.x, target_y: record.position.y })
    Message.success('已下发导航目标')
  } catch (err: any) {
    Message.error('导航失败: ' + err.message)
  }
}

async function deleteAnnotation(record: any) {
  Modal.warning({
    title: '确认删除',
    content: `确定删除标注「${record.label}」？`,
    onOk: async () => {
      try {
        await updateMap(deviceId.value, { action: 'delete_annotation', annotation_id: record.id })
        Message.success('已删除')
        if (activeMap.value) selectMap(activeMap.value)
      } catch (err: any) {
        Message.error('删除失败: ' + err.message)
      }
    }
  })
}

watch(displayType, drawMap)

onMounted(async () => {
  await loadMaps()
  await loadLocalization()
  await drawMap()
})
</script>

<style scoped>
.map-list-card, .localization-card {
  height: 100%;
}
.active-map-item {
  background: var(--color-primary-light-1);
  border-radius: 4px;
}
.map-canvas-container {
  position: relative;
  width: 100%;
  height: 420px;
  border: 1px solid var(--color-border);
  border-radius: 4px;
  overflow: hidden;
}
.map-canvas {
  width: 100%;
  height: 100%;
  display: block;
}
.map-empty-overlay {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255,255,255,0.8);
}
.text-muted {
  color: var(--color-text-3);
  font-size: 12px;
}
</style>
