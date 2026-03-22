<template>
  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>具身智能</a-breadcrumb-item>
      <a-breadcrumb-item>
        <router-link :to="`/embodied/${deviceId}/perception`">设备 {{ deviceId }}</router-link>
      </a-breadcrumb-item>
      <a-breadcrumb-item>安全禁区</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 紧急停止与安全概览 -->
    <a-row :gutter="16" style="margin-top: 16px">
      <a-col :xs="24" :sm="24" :md="8">
        <a-card class="emergency-card" :bordered="false">
          <div class="emergency-title">
            <icon-warning size="20" style="color: #ff4d4f" />
            紧急停止
          </div>
          <a-button
            type="primary"
            danger
            long
            size="large"
            :loading="stopping"
            @click="handleEmergencyStop"
          >
            <template #icon><icon-stop /></template>
            一键紧急停止
          </a-button>
          <div class="emergency-hint">立即停止设备所有运动</div>
        </a-card>

        <!-- 安全统计 -->
        <a-card title="安全统计" style="margin-top: 16px">
          <a-descriptions :column="1" size="small">
            <a-descriptions-item label="禁区数">
              <a-badge :count="safetyStats.forbidden" :number-style="{ backgroundColor: '#ff4d4f' }" :max-count="99" />
            </a-descriptions-item>
            <a-descriptions-item label="警戒区数">
              <a-badge :count="safetyStats.caution" :number-style="{ backgroundColor: '#ff7800' }" :max-count="99" />
            </a-descriptions-item>
            <a-descriptions-item label="安全区数">
              <a-badge :count="safetyStats.safe" :number-style="{ backgroundColor: '#00b42a' }" :max-count="99" />
            </a-descriptions-item>
            <a-descriptions-item label="今日告警">
              <a-badge :count="safetyStats.todayAlerts" :number-style="{ backgroundColor: '#1650d8' }" :max-count="99" />
            </a-descriptions-item>
            <a-descriptions-item label="紧急停止历史">
              <a-badge :count="safetyStats.emergencyStops" :number-style="{ backgroundColor: '#ff4d4f' }" :max-count="99" />
            </a-descriptions-item>
          </a-descriptions>
        </a-card>
      </a-col>

      <a-col :xs="24" :sm="24" :md="16">
        <!-- 禁区列表 -->
        <a-card title="安全区域配置">
          <template #extra>
            <a-button type="primary" size="small" @click="openZoneModal()">
              <template #icon><icon-plus /></template>
              添加区域
            </a-button>
          </template>
          <a-table
            :columns="zoneColumns"
            :data="zones"
            :loading="loadingZones"
            :pagination="{ pageSize: 8, showTotal: true }"
            row-key="id"
            size="small"
          >
            <template #zone_type="{ record }">
              <a-tag :color="getZoneColor(record.zone_type)">
                {{ getZoneText(record.zone_type) }}
              </a-tag>
            </template>
            <template #zone_shape="{ record }">
              {{ getShapeText(record.zone_shape) }}
            </template>
            <template #zone_data="{ record }">
              <span class="text-muted" v-if="record.zone_shape === 'rectangle'">
                {{ record.zone_data?.x1?.toFixed(1) }},{{ record.zone_data?.y1?.toFixed(1) }}
                → {{ record.zone_data?.x2?.toFixed(1) }},{{ record.zone_data?.y2?.toFixed(1) }}
              </span>
              <span class="text-muted" v-else-if="record.zone_shape === 'circle'">
                圆心({{ record.zone_data?.cx?.toFixed(1) }},{{ record.zone_data?.cy?.toFixed(1) }})
                r={{ record.zone_data?.r?.toFixed(1) }}m
              </span>
              <span v-else>-</span>
            </template>
            <template #is_enabled="{ record }">
              <a-switch v-model="record.is_enabled" @change="toggleZone(record)" size="small" />
            </template>
            <template #actions="{ record }">
              <a-button type="text" size="small" @click="openZoneModal(record)">
                <icon-edit />
              </a-button>
              <a-button type="text" size="small" @click="deleteZone(record)">
                <icon-delete />
              </a-button>
            </template>
          </a-table>
        </a-card>
      </a-col>
    </a-row>

    <!-- 安全日志 & 紧急停止历史 -->
    <a-row :gutter="16" style="margin-top: 16px">
      <a-col :xs="24" :sm="24" :md="12">
        <a-card title="安全日志">
          <template #extra>
            <a-button type="text" size="small" @click="loadSafetyLogs">
              <icon-refresh />
            </a-button>
          </template>
          <a-table
            :columns="logColumns"
            :data="safetyLogs"
            :loading="loadingLogs"
            :pagination="{ pageSize: 8, showTotal: true }"
            row-key="id"
            size="small"
          >
            <template #severity="{ record }">
              <a-tag :color="getSeverityColor(record.severity)">
                {{ getSeverityText(record.severity) }}
              </a-tag>
            </template>
            <template #event_type="{ record }">
              {{ getEventTypeText(record.event_type) }}
            </template>
            <template #resolved="{ record }">
              <a-tag :color="record.resolved ? 'green' : 'orange'" size="small">
                {{ record.resolved ? '已解决' : '未解决' }}
              </a-tag>
            </template>
            <template #created_at="{ record }">
              {{ formatTime(record.created_at) }}
            </template>
          </a-table>
        </a-card>
      </a-col>

      <a-col :xs="24" :sm="24" :md="12">
        <a-card title="紧急停止历史">
          <template #extra>
            <a-button type="text" size="small" @click="loadEmergencyHistory">
              <icon-refresh />
            </a-button>
          </template>
          <a-table
            :columns="emergencyColumns"
            :data="emergencyHistory"
            :loading="loadingEmergency"
            :pagination="{ pageSize: 8, showTotal: true }"
            row-key="id"
            size="small"
          >
            <template #stop_type="{ record }">
              <a-tag color="red">{{ record.stop_type || 'emergency' }}</a-tag>
            </template>
            <template #resolved="{ record }">
              <a-tag :color="record.resolved ? 'green' : 'orange'" size="small">
                {{ record.resolved ? '已恢复' : '待处理' }}
              </a-tag>
            </template>
            <template #actions="{ record }">
              <a-button
                type="text"
                size="small"
                :disabled="record.resolved"
                @click="resolveEmergency(record)"
              >
                <icon-check-circle />
              </a-button>
            </template>
          </a-table>
        </a-card>
      </a-col>
    </a-row>

    <!-- 区域编辑弹窗 -->
    <a-modal
      v-model:visible="zoneModalVisible"
      :title="editingZone ? '编辑安全区域' : '添加安全区域'"
      @before-ok="submitZone"
      @cancel="zoneModalVisible = false"
    >
      <a-form :model="zoneForm" layout="vertical">
        <a-form-item label="区域名称" required>
          <a-input v-model="zoneForm.zone_name" placeholder="输入区域名称" />
        </a-form-item>
        <a-form-item label="区域类型" required>
          <a-select v-model="zoneForm.zone_type" placeholder="选择类型">
            <a-option value="forbidden">禁区</a-option>
            <a-option value="caution">警戒区</a-option>
            <a-option value="safe">安全区</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="形状" required>
          <a-select v-model="zoneForm.zone_shape" placeholder="选择形状">
            <a-option value="rectangle">矩形</a-option>
            <a-option value="circle">圆形</a-option>
            <a-option value="polygon">多边形</a-option>
          </a-select>
        </a-form-item>
        <!-- 矩形参数 -->
        <template v-if="zoneForm.zone_shape === 'rectangle'">
          <a-form-item label="左下角 X">
            <a-input-number v-model="zoneForm.zone_data.x1" :min="0" style="width: 100%" />
          </a-form-item>
          <a-form-item label="左下角 Y">
            <a-input-number v-model="zoneForm.zone_data.y1" :min="0" style="width: 100%" />
          </a-form-item>
          <a-form-item label="右上角 X">
            <a-input-number v-model="zoneForm.zone_data.x2" :min="0" style="width: 100%" />
          </a-form-item>
          <a-form-item label="右上角 Y">
            <a-input-number v-model="zoneForm.zone_data.y2" :min="0" style="width: 100%" />
          </a-form-item>
        </template>
        <!-- 圆形参数 -->
        <template v-if="zoneForm.zone_shape === 'circle'">
          <a-form-item label="圆心 X">
            <a-input-number v-model="zoneForm.zone_data.cx" :min="0" style="width: 100%" />
          </a-form-item>
          <a-form-item label="圆心 Y">
            <a-input-number v-model="zoneForm.zone_data.cy" :min="0" style="width: 100%" />
          </a-form-item>
          <a-form-item label="半径 (m)">
            <a-input-number v-model="zoneForm.zone_data.r" :min="0.1" :max="10" :step="0.1" style="width: 100%" />
          </a-form-item>
        </template>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { getSafetyZones, createSafetyZone, updateSafetyZone, deleteSafetyZone, emergencyStop, getSafetyLogs } from '@/api/embodied'
import { Message, Modal } from '@arco-design/web-vue'

const route = useRoute()
const deviceId = ref(route.params.device_id as string)

const loadingZones = ref(false)
const loadingLogs = ref(false)
const loadingEmergency = ref(false)
const stopping = ref(false)
const zones = ref<any[]>([])
const safetyLogs = ref<any[]>([])
const emergencyHistory = ref<any[]>([])
const safetyStats = ref({ forbidden: 0, caution: 0, safe: 0, todayAlerts: 0, emergencyStops: 0 })

const zoneModalVisible = ref(false)
const editingZone = ref<any>(null)
const zoneForm = ref({
  zone_name: '',
  zone_type: 'forbidden',
  zone_shape: 'rectangle',
  zone_data: { x1: 0, y1: 0, x2: 1, y2: 1 }
})

const zoneColumns = [
  { title: '名称', dataIndex: 'zone_name', width: 120 },
  { title: '类型', dataIndex: 'zone_type', slotName: 'zone_type', width: 100 },
  { title: '形状', dataIndex: 'zone_shape', slotName: 'zone_shape', width: 90 },
  { title: '参数', dataIndex: 'zone_data', slotName: 'zone_data' },
  { title: '启用', dataIndex: 'is_enabled', slotName: 'is_enabled', width: 70 },
  { title: '操作', slotName: 'actions', width: 80 }
]

const logColumns = [
  { title: '事件', dataIndex: 'event_type', slotName: 'event_type', width: 120 },
  { title: '级别', dataIndex: 'severity', slotName: 'severity', width: 80 },
  { title: '详情', dataIndex: 'details', ellipsis: true },
  { title: '状态', dataIndex: 'resolved', slotName: 'resolved', width: 80 },
  { title: '时间', dataIndex: 'created_at', slotName: 'created_at', width: 160 }
]

const emergencyColumns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '类型', dataIndex: 'stop_type', slotName: 'stop_type', width: 100 },
  { title: '触发原因', dataIndex: 'reason', ellipsis: true },
  { title: '状态', dataIndex: 'resolved', slotName: 'resolved', width: 80 },
  { title: '时间', dataIndex: 'created_at', slotName: 'created_at', width: 160 },
  { title: '操作', slotName: 'actions', width: 60 }
]

async function loadZones() {
  try {
    loadingZones.value = true
    const res = await getSafetyZones(deviceId.value)
    zones.value = res.data?.zones || res.data || []
    updateStats()
  } catch (err: any) {
    Message.error('加载安全区域失败: ' + err.message)
  } finally {
    loadingZones.value = false
  }
}

async function loadSafetyLogs() {
  try {
    loadingLogs.value = true
    const res = await getSafetyLogs(deviceId.value, { page_size: 20 })
    safetyLogs.value = res.data?.logs || res.data || []
  } catch (err: any) {
    Message.error('加载安全日志失败: ' + err.message)
  } finally {
    loadingLogs.value = false
  }
}

async function loadEmergencyHistory() {
  try {
    loadingEmergency.value = true
    const res = await getSafetyLogs(deviceId.value, { event_type: 'emergency_stop', page_size: 20 })
    emergencyHistory.value = res.data?.logs || res.data || []
  } catch (err: any) {
    Message.error('加载紧急停止历史失败: ' + err.message)
  } finally {
    loadingEmergency.value = false
  }
}

function updateStats() {
  safetyStats.value = {
    forbidden: zones.value.filter(z => z.zone_type === 'forbidden').length,
    caution: zones.value.filter(z => z.zone_type === 'caution').length,
    safe: zones.value.filter(z => z.zone_type === 'safe').length,
    todayAlerts: safetyLogs.value.filter(l => {
      const d = new Date(l.created_at)
      const now = new Date()
      return d.toDateString() === now.toDateString() && l.severity !== 'info'
    }).length,
    emergencyStops: emergencyHistory.value.length
  }
}

function openZoneModal(zone?: any) {
  editingZone.value = zone || null
  if (zone) {
    zoneForm.value = {
      zone_name: zone.zone_name,
      zone_type: zone.zone_type,
      zone_shape: zone.zone_shape,
      zone_data: zone.zone_data || { x1: 0, y1: 0, x2: 1, y2: 1 }
    }
  } else {
    zoneForm.value = {
      zone_name: '',
      zone_type: 'forbidden',
      zone_shape: 'rectangle',
      zone_data: { x1: 0, y1: 0, x2: 1, y2: 1 }
    }
  }
  zoneModalVisible.value = true
}

async function submitZone(done: (val: boolean) => void) {
  try {
    if (editingZone.value) {
      await updateSafetyZone(deviceId.value, editingZone.value.id, zoneForm.value)
      Message.success('区域已更新')
    } else {
      await createSafetyZone(deviceId.value, zoneForm.value)
      Message.success('区域已创建')
    }
    zoneModalVisible.value = false
    await loadZones()
    done(true)
  } catch (err: any) {
    Message.error('操作失败: ' + err.message)
    done(false)
  }
}

async function toggleZone(zone: any) {
  try {
    await updateSafetyZone(deviceId.value, zone.id, { is_enabled: zone.is_enabled })
    Message.success(zone.is_enabled ? '区域已启用' : '区域已禁用')
  } catch (err: any) {
    zone.is_enabled = !zone.is_enabled
    Message.error('操作失败: ' + err.message)
  }
}

async function deleteZone(zone: any) {
  Modal.warning({
    title: '确认删除',
    content: `确定删除区域「${zone.zone_name}」？`,
    onOk: async () => {
      try {
        await deleteSafetyZone(deviceId.value, zone.id)
        Message.success('已删除')
        await loadZones()
      } catch (err: any) {
        Message.error('删除失败: ' + err.message)
      }
    }
  })
}

async function handleEmergencyStop() {
  Modal.warning({
    title: '⚠️ 确认紧急停止',
    content: '这将立即停止设备所有运动，可能中断当前任务。',
    okText: '确认停止',
    onOk: async () => {
      try {
        stopping.value = true
        await emergencyStop(deviceId.value)
        Message.warning('已执行紧急停止')
        await loadEmergencyHistory()
        await loadSafetyLogs()
      } catch (err: any) {
        Message.error('紧急停止失败: ' + err.message)
      } finally {
        stopping.value = false
      }
    }
  })
}

async function resolveEmergency(record: any) {
  try {
    await updateSafetyZone(deviceId.value, record.id, { resolved: true })
    Message.success('已标记为已恢复')
    await loadEmergencyHistory()
  } catch (err: any) {
    Message.error('操作失败: ' + err.message)
  }
}

function getZoneColor(type: string) {
  const map: Record<string, string> = { forbidden: 'red', caution: 'orange', safe: 'green' }
  return map[type] || 'default'
}
function getZoneText(type: string) {
  const map: Record<string, string> = { forbidden: '禁区', caution: '警戒区', safe: '安全区' }
  return map[type] || type
}
function getShapeText(shape: string) {
  const map: Record<string, string> = { rectangle: '矩形', circle: '圆形', polygon: '多边形' }
  return map[shape] || shape
}
function getSeverityColor(s: string) {
  const map: Record<string, string> = { info: 'gray', warning: 'orange', critical: 'red' }
  return map[s] || 'default'
}
function getSeverityText(s: string) {
  const map: Record<string, string> = { info: '信息', warning: '警告', critical: '严重' }
  return map[s] || s
}
function getEventTypeText(t: string) {
  const map: Record<string, string> = {
    collision: '碰撞',
    emergency_stop: '紧急停止',
    zone_violation: '区域越界',
    fall_prevention: '防跌落'
  }
  return map[t] || t
}
function formatTime(t: string) {
  if (!t) return '-'
  return new Date(t).toLocaleString('zh-CN')
}

onMounted(async () => {
  await Promise.all([loadZones(), loadSafetyLogs(), loadEmergencyHistory()])
})
</script>

<style scoped>
.emergency-card {
  border: 2px solid #ff4d4f;
  border-radius: 8px;
  background: #fff5f5;
}
.emergency-title {
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 12px;
  display: flex;
  align-items: center;
  gap: 8px;
  color: #ff4d4f;
}
.emergency-hint {
  text-align: center;
  font-size: 12px;
  color: var(--color-text-3);
  margin-top: 8px;
}
.text-muted {
  color: var(--color-text-3);
  font-size: 12px;
}
</style>
