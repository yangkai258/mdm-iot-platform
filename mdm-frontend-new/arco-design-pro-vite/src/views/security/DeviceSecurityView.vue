<template>
    <Breadcrumb :items="['Home','Console','']" />


  <div class="page-container">
    <!-- 设备选择栏 -->
    <a-card class="filter-card">
      <div class="filter-row">
        <a-select
          v-model="selectedDeviceId"
          placeholder="选择设备"
          style="width: 280px"
          allow-search
          @change="loadDeviceInfo"
        >
          <a-option v-for="d in devices" :key="d.id" :value="d.id">
            {{ d.name }} ({{ d.device_type || '未知类型' }})
          </a-option>
        </a-select>
      </div>
    </a-card>

    <!-- 设备信息 + 操作区 -->
    <div class="content-grid" v-if="selectedDeviceId && currentDevice">
      <!-- 设备信息卡片 -->
      <a-card class="device-info-card">
        <template #title>
          <div class="card-title">
            <icon-info-circle />
            <span>设备信息</span>
          </div>
        </template>
        <div class="device-info">
          <div class="info-row">
            <span class="info-label">设备名称</span>
            <span class="info-value">{{ currentDevice.name || '-' }}</span>
          </div>
          <div class="info-row">
            <span class="info-label">设备类型</span>
            <span class="info-value">{{ currentDevice.device_type || '-' }}</span>
          </div>
          <div class="info-row">
            <span class="info-label">在线状态</span>
            <span class="info-value">
              <a-badge
                :color="currentDevice.online ? 'green' : 'gray'"
                :text="currentDevice.online ? '在线' : '离线'"
              />
            </span>
          </div>
          <div class="info-row">
            <span class="info-label">设备ID</span>
            <span class="info-value mono">{{ currentDevice.id }}</span>
          </div>
          <div class="info-row">
            <span class="info-label">最后在线</span>
            <span class="info-value">{{ formatDate(currentDevice.last_seen_at) }}</span>
          </div>
        </div>
      </a-card>

      <!-- 安全操作卡片 -->
      <a-card class="actions-card">
        <template #title>
          <div class="card-title">
            <icon-safe />
            <span>安全操作</span>
          </div>
        </template>
        <div class="action-list">
          <div class="action-item">
            <div class="action-info">
              <div class="action-name">锁定设备</div>
              <div class="action-desc">远程锁定设备，用户无法正常使用</div>
            </div>
            <a-button
              type="primary"
              :loading="locking"
              :disabled="!currentDevice.online || currentDevice.locked"
              @click="handleLock"
            >
              {{ currentDevice.locked ? '已锁定' : '锁定设备' }}
            </a-button>
          </div>
          <a-divider />
          <div class="action-item">
            <div class="action-info">
              <div class="action-name">解锁设备</div>
              <div class="action-desc">解除设备远程锁定，恢复正常使用</div>
            </div>
            <a-button
              :loading="unlocking"
              :disabled="!currentDevice.online || !currentDevice.locked"
              @click="handleUnlock"
            >
              解锁设备
            </a-button>
          </div>
          <a-divider />
          <div class="action-item danger-zone">
            <div class="action-info">
              <div class="action-name">数据擦除</div>
              <div class="action-desc danger-desc">危险操作！将清除设备上的所有数据，且无法恢复</div>
            </div>
            <a-button
              status="danger"
              :disabled="!currentDevice.online"
              @click="openWipeConfirm"
            >
              数据擦除
            </a-button>
          </div>
        </div>
      </a-card>

      <!-- 操作历史卡片 -->
      <a-card class="history-card">
        <template #title>
          <div class="card-title">
            <icon-history />
            <span>操作历史</span>
          </div>
        </template>
        <a-table
          :columns="historyColumns"
          :data="wipeHistory"
          :loading="historyLoading"
          :pagination="false"
          size="small"
        >
          <template #type="{ record }">
            <a-tag :color="record.type === 'lock' ? 'blue' : 'red'">
              {{ record.type === 'lock' ? '锁定' : record.type === 'unlock' ? '解锁' : '擦除' }}
            </a-tag>
          </template>
          <template #status="{ record }">
            <a-tag :color="statusColor(record.status)">{{ statusLabel(record.status) }}</a-tag>
          </template>
          <template #time="{ record }">
            {{ formatDate(record.created_at) }}
          </template>
        </a-table>
      </a-card>
    </div>

    <!-- 空状态 -->
    <a-card class="empty-card" v-else>
      <a-empty description="请从上方选择设备以进行安全操作" />
    </a-card>

    <!-- 擦除确认弹窗 -->
    <DeviceWipeConfirmModal
      v-model:visible="wipeModalVisible"
      :device="currentDevice"
      @confirm="handleWipeConfirm"
    />
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import { getDevices, lockDevice, unlockDevice, wipeDevice, confirmWipeDevice, getWipeHistory } from '@/api/security'
import DeviceWipeConfirmModal from '@/components/security/DeviceWipeConfirmModal.vue'
import dayjs from 'dayjs'

const selectedDeviceId = ref(null)
const devices = ref([])
const currentDevice = ref(null)
const wipeHistory = ref([])
const historyLoading = ref(false)
const locking = ref(false)
const unlocking = ref(false)
const wipeModalVisible = ref(false)

const historyColumns = [
  { title: '操作类型', slotName: 'type', width: 100 },
  { title: '操作人', dataIndex: 'operator', width: 120 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '时间', slotName: 'time', width: 160 },
  { title: '备注', dataIndex: 'reason', ellipsis: true }
]

onMounted(async () => {
  try {
    const res = await getDevices({ page_size: 200 })
    devices.value = res.data || res || []
  } catch (e) {
    console.error('加载设备列表失败', e)
  }
})

async function loadDeviceInfo(deviceId) {
  if (!deviceId) return
  try {
    const res = await getDevices({ id: deviceId })
    const list = res.data || res || []
    currentDevice.value = Array.isArray(list) ? list[0] : list
    await loadWipeHistory(deviceId)
  } catch (e) {
    console.error('加载设备信息失败', e)
  }
}

async function loadWipeHistory(deviceId) {
  historyLoading.value = true
  try {
    const res = await getWipeHistory(deviceId)
    wipeHistory.value = res.data || res || []
  } catch (e) {
    console.error('加载操作历史失败', e)
  } finally {
    historyLoading.value = false
  }
}

async function handleLock() {
  locking.value = true
  try {
    await lockDevice(selectedDeviceId.value)
    Message.success('锁定命令已下发')
    await loadDeviceInfo(selectedDeviceId.value)
  } catch (e) {
    Message.error('锁定失败')
  } finally {
    locking.value = false
  }
}

async function handleUnlock() {
  unlocking.value = true
  try {
    await unlockDevice(selectedDeviceId.value)
    Message.success('解锁命令已下发')
    await loadDeviceInfo(selectedDeviceId.value)
  } catch (e) {
    Message.error('解锁失败')
  } finally {
    unlocking.value = false
  }
}

function openWipeConfirm() {
  wipeModalVisible.value = true
}

async function handleWipeConfirm(formData) {
  try {
    await wipeDevice(selectedDeviceId.value)
    await confirmWipeDevice(selectedDeviceId.value, formData)
    Message.success('擦除命令已下发')
    wipeModalVisible.value = false
    await loadDeviceInfo(selectedDeviceId.value)
  } catch (e) {
    Message.error('擦除操作失败')
  }
}

function statusColor(status) {
  const map = { pending: 'yellow', sent: 'blue', confirmed: 'green', failed: 'red' }
  return map[status] || 'default'
}

function statusLabel(status) {
  const map = { pending: '待确认', sent: '已下发', confirmed: '已完成', failed: '失败' }
  return map[status] || status
}

function formatDate(date) {
  if (!date) return '-'
  return dayjs(date).format('YYYY-MM-DD HH:mm:ss')
}
</script>

<style scoped>
.page-container {
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
  height: 100%;
  box-sizing: border-box;
}

.filter-card {
  flex-shrink: 0;
}

.filter-row {
  display: flex;
  align-items: center;
  gap: 10px;
}

.content-grid {
  flex: 1;
  display: grid;
  grid-template-columns: 300px 1fr;
  grid-template-rows: auto 1fr;
  gap: 12px;
  overflow: hidden;
}

.device-info-card {
  grid-row: 1;
}

.actions-card {
  grid-row: 1;
}

.history-card {
  grid-column: 1 / -1;
  grid-row: 2;
  overflow: auto;
}

.empty-card {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
}

.card-title {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 14px;
}

.device-info {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.info-row {
  display: flex;
  gap: 10px;
}

.info-label {
  width: 80px;
  color: var(--color-text-3);
  font-size: 13px;
  flex-shrink: 0;
}

.info-value {
  font-size: 13px;
  word-break: break-all;
}

.mono {
  font-family: monospace;
  font-size: 12px;
}

.action-list {
  display: flex;
  flex-direction: column;
}

.action-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 0;
}

.action-info {
  flex: 1;
}

.action-name {
  font-size: 14px;
  font-weight: 500;
  margin-bottom: 4px;
}

.action-desc {
  font-size: 12px;
  color: var(--color-text-3);
}

.danger-desc {
  color: #f53f3f;
}
</style>
