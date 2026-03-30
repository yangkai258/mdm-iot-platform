<template>
  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>设备管理</a-breadcrumb-item>
      <a-breadcrumb-item>远程控制</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 页面布局：左侧设备卡片 + 右侧操作面板 -->
    <a-row :gutter="16">
      <!-- 左侧：设备列表 -->
      <a-col :span="14">
        <a-card title="设备列表" class="device-list-card">
          <template #extra>
            <a-button size="small" @click="loadDevices">刷新</a-button>
          </template>
          <!-- 搜索栏 -->
          <div class="search-bar">
            <a-space>
              <a-input-search v-model="searchKeyword" placeholder="搜索设备ID/名称" style="width: 220px" @search="loadDevices" search-button />
              <a-select v-model="filterOnline" placeholder="在线状态" allow-clear style="width: 110px" @change="loadDevices">
                <a-option value="online">在线</a-option>
                <a-option value="offline">离线</a-option>
              </a-select>
            </a-space>
          </div>

          <!-- 设备表格 -->
          <a-table
            :columns="deviceColumns"
            :data="devices"
            :loading="loading"
            :pagination="pagination"
            row-key="device_id"
            :scroll="{ x: 600 }"
            @page-change="handlePageChange"
            size="small"
          >
            <template #is_online="{ record }">
              <span class="online-dot" :class="record.is_online ? 'online' : 'offline'"></span>
              {{ record.is_online ? '在线' : '离线' }}
            </template>
            <template #actions="{ record }">
              <a-space size="small">
                <a-button type="primary" size="mini" :disabled="!record.is_online" @click="showLockModal(record)">锁定</a-button>
                <a-button type="primary" size="mini" :disabled="!record.is_online" @click="showUnlockModal(record)" v-if="record.is_locked">解锁</a-button>
                <a-button type="primary" status="danger" size="mini" :disabled="!record.is_online" @click="showWipeModal(record)">擦除</a-button>
              </a-space>
            </template>
          </a-table>
        </a-card>
      </a-col>

      <!-- 右侧：设备详情 + 操作面板 -->
      <a-col :span="10">
        <!-- 选中设备信息 -->
        <a-card title="设备详情" class="device-info-card" v-if="selectedDevice">
          <a-descriptions :column="1" size="small">
            <a-descriptions-item label="设备ID">{{ selectedDevice.device_id }}</a-descriptions-item>
            <a-descriptions-item label="硬件型号">{{ selectedDevice.hardware_model || '-' }}</a-descriptions-item>
            <a-descriptions-item label="固件版本">{{ selectedDevice.firmware_version || '-' }}</a-descriptions-item>
            <a-descriptions-item label="在线状态">
              <span class="online-dot" :class="selectedDevice.is_online ? 'online' : 'offline'"></span>
              {{ selectedDevice.is_online ? '在线' : '离线' }}
            </a-descriptions-item>
            <a-descriptions-item label="绑定用户">{{ selectedDevice.bind_user_id || '-' }}</a-descriptions-item>
            <a-descriptions-item label="锁定状态">
              <a-tag :color="selectedDevice.is_locked ? 'red' : 'green'">
                {{ selectedDevice.is_locked ? '已锁定' : '未锁定' }}
              </a-tag>
            </a-descriptions-item>
          </a-descriptions>
        </a-card>

        <!-- 操作面板 -->
        <a-card title="远程操作" class="action-card">
          <template #title>
            <span>远程操作</span>
          </template>
          <div v-if="!selectedDevice" class="no-device-tip">
            <icon-info-circle style="font-size: 32px; color: #999; margin-bottom: 8px" />
            <p>请从左侧列表选择设备</p>
          </div>
          <div v-else class="action-grid">
            <a-button type="primary" long :disabled="!selectedDevice.is_online" @click="showLockModal(selectedDevice)">
              <template #icon><icon-lock /></template>
              锁定设备
            </a-button>
            <a-button type="primary" long :disabled="!selectedDevice.is_online || !selectedDevice.is_locked" @click="showUnlockModal(selectedDevice)">
              <template #icon><icon-unlock /></template>
              解锁设备
            </a-button>
            <a-button type="primary" status="warning" long :disabled="!selectedDevice.is_online" @click="showWipeModal(selectedDevice)">
              <template #icon><icon-delete /></template>
              远程擦除（数据）
            </a-button>
            <a-button type="primary" status="danger" long :disabled="!selectedDevice.is_online" @click="showFactoryResetModal(selectedDevice)">
              <template #icon><icon-refresh /></template>
              恢复出厂设置
            </a-button>
          </div>
        </a-card>

        <!-- 操作历史 -->
        <a-card title="擦除历史" class="history-card">
          <a-table
            :columns="historyColumns"
            :data="wipeHistory"
            :loading="historyLoading"
            :pagination="historyPagination"
            row-key="task_id"
            size="small"
            @page-change="handleHistoryPageChange"
          >
            <template #wipe_type="{ record }">
              <a-tag :color="record.wipe_type === 'factory_reset' ? 'red' : 'orange'">
                {{ record.wipe_type === 'factory_reset' ? '恢复出厂' : '数据擦除' }}
              </a-tag>
            </template>
            <template #status="{ record }">
              <a-tag :color="getWipeStatusColor(record.status)">{{ getWipeStatusText(record.status) }}</a-tag>
            </template>
          </a-table>
        </a-card>
      </a-col>
    </a-row>

    <!-- 锁定弹窗 -->
    <a-modal v-model:visible="lockModalVisible" title="锁定设备" :width="480" :loading="submitting" @before-ok="handleLock" @cancel="lockModalVisible = false">
      <a-form :model="lockForm" layout="vertical">
        <a-form-item label="目标设备">
          <a-input :value="selectedDevice?.device_id" disabled />
        </a-form-item>
        <a-form-item label="锁定原因">
          <a-select v-model="lockForm.reason" placeholder="选择原因">
            <a-option value="lost">设备丢失</a-option>
            <a-option value="stolen">设备被盗</a-option>
            <a-option value="maintenance">设备维护</a-option>
            <a-option value="security">安全风险</a-option>
            <a-option value="other">其他原因</a-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 解锁弹窗 -->
    <a-modal v-model:visible="unlockModalVisible" title="解锁设备" :width="480" :loading="submitting" @before-ok="handleUnlock" @cancel="unlockModalVisible = false">
      <a-form :model="unlockForm" layout="vertical">
        <a-form-item label="目标设备">
          <a-input :value="selectedDevice?.device_id" disabled />
        </a-form-item>
        <a-form-item label="解锁码" required>
          <a-input-password v-model="unlockForm.unlock_code" placeholder="请输入解锁码" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 擦除弹窗 -->
    <a-modal v-model:visible="wipeModalVisible" title="远程擦除" :width="520" :loading="submitting" @before-ok="handleWipe" @cancel="wipeModalVisible = false" :mask-closable="false">
      <a-result status="warning" title="危险操作确认" subtitle="远程擦除将清除设备上的所有数据，且不可恢复！">
        <template #extra>
          <a-form :model="wipeForm" layout="vertical">
            <a-form-item label="擦除类型">
              <a-radio-group v-model="wipeForm.wipe_type">
                <a-radio value="data_only">仅擦除数据</a-radio>
                <a-radio value="factory_reset">恢复出厂设置</a-radio>
              </a-radio-group>
            </a-form-item>
            <a-form-item label="确认码" required>
              <a-input v-model="wipeForm.confirm_code" placeholder="请输入 CONFIRM-WIPE 以确认" />
            </a-form-item>
          </a-form>
        </template>
      </a-result>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import { lockDevice, unlockDevice, wipeDevice, getDeviceSecurityStatus } from '@/api/security'

const loading = ref(false)
const historyLoading = ref(false)
const submitting = ref(false)
const devices = ref([])
const wipeHistory = ref([])
const selectedDevice = ref(null)
const searchKeyword = ref('')
const filterOnline = ref('')
const lockModalVisible = ref(false)
const unlockModalVisible = ref(false)
const wipeModalVisible = ref(false)

const lockForm = reactive({ reason: '' })
const unlockForm = reactive({ unlock_code: '' })
const wipeForm = reactive({ wipe_type: 'data_only', confirm_code: '' })

const pagination = reactive({ current: 1, pageSize: 10, total: 0 })
const historyPagination = reactive({ current: 1, pageSize: 5, total: 0 })

const deviceColumns = [
  { title: '设备ID', dataIndex: 'device_id', width: 140, ellipsis: true },
  { title: '硬件型号', dataIndex: 'hardware_model', width: 120 },
  { title: '在线状态', slotName: 'is_online', width: 90 },
  { title: '绑定用户', dataIndex: 'bind_user_id', width: 100, ellipsis: true },
  { title: '操作', slotName: 'actions', width: 180 }
]

const historyColumns = [
  { title: '任务ID', dataIndex: 'task_id', width: 120, ellipsis: true },
  { title: '类型', slotName: 'wipe_type', width: 110 },
  { title: '状态', slotName: 'status', width: 90 },
  { title: '操作人', dataIndex: 'initiated_by', width: 100 },
  { title: '发起时间', dataIndex: 'initiated_at', width: 160 }
]

const getWipeStatusColor = (s) => ({ initiated: 'blue', completed: 'green', failed: 'red' }[s] || 'gray')
const getWipeStatusText = (s) => ({ initiated: '进行中', completed: '已完成', failed: '失败' }[s] || s)

const getToken = () => localStorage.getItem('token')

const loadDevices = async () => {
  loading.value = true
  try {
    const params = { page: pagination.current, page_size: pagination.pageSize }
    if (searchKeyword.value) params.keyword = searchKeyword.value
    if (filterOnline.value) params.status = filterOnline.value

    const res = await fetch(`/api/v1/devices?${new URLSearchParams(params)}`, {
      headers: { 'Authorization': `Bearer ${getToken()}` }
    })
    const json = await res.json()
    if (json.code === 0) {
      devices.value = json.data.list || []
      pagination.total = json.data.total || 0
    }
  } catch (e) {
    Message.error('加载设备列表失败')
  } finally {
    loading.value = false
  }
}

const loadDeviceSecurityStatus = async (device) => {
  try {
    const res = await getDeviceSecurityStatus(device.device_id)
    if (res.code === 0) {
      device.is_locked = res.data.is_locked
      device.is_wiped = res.data.is_wiped
    }
  } catch (e) { /* silent */ }
}

const loadWipeHistory = async () => {
  if (!selectedDevice.value) return
  historyLoading.value = true
  try {
    const res = await fetch(`/api/v1/devices/${selectedDevice.value.device_id}/wipe-history?${new URLSearchParams({ page: historyPagination.current, page_size: historyPagination.pageSize })}`, {
      headers: { 'Authorization': `Bearer ${getToken()}` }
    })
    const json = await res.json()
    if (json.code === 0) {
      wipeHistory.value = json.data.list || []
      historyPagination.total = json.data.total || 0
    }
  } catch (e) { /* silent */ } finally {
    historyLoading.value = false
  }
}

const handlePageChange = (page) => { pagination.current = page; loadDevices() }
const handleHistoryPageChange = (page) => { historyPagination.current = page; loadWipeHistory() }

const showLockModal = (device) => {
  selectedDevice.value = device
  lockForm.reason = ''
  lockModalVisible.value = true
}

const showUnlockModal = (device) => {
  selectedDevice.value = device
  unlockForm.unlock_code = ''
  unlockModalVisible.value = true
}

const showWipeModal = (device) => {
  selectedDevice.value = device
  wipeForm.wipe_type = 'data_only'
  wipeForm.confirm_code = ''
  wipeModalVisible.value = true
}

const showFactoryResetModal = (device) => {
  selectedDevice.value = device
  wipeForm.wipe_type = 'factory_reset'
  wipeForm.confirm_code = ''
  wipeModalVisible.value = true
}

const handleLock = async (done) => {
  if (!selectedDevice.value) { done(false); return }
  submitting.value = true
  try {
    const res = await lockDevice(selectedDevice.value.device_id, lockForm.reason)
    if (res.code === 0) {
      Message.success(res.message || '锁定命令已下发')
      lockModalVisible.value = false
      loadDevices()
      done(true)
    } else {
      Message.error(res.message || '锁定失败')
      done(false)
    }
  } catch (e) {
    Message.error('锁定失败')
    done(false)
  } finally {
    submitting.value = false
  }
}

const handleUnlock = async (done) => {
  if (!unlockForm.unlock_code) {
    Message.warning('请输入解锁码')
    done(false)
    return
  }
  submitting.value = true
  try {
    const res = await unlockDevice(selectedDevice.value.device_id, unlockForm.unlock_code)
    if (res.code === 0) {
      Message.success('解锁命令已下发')
      unlockModalVisible.value = false
      loadDevices()
      done(true)
    } else {
      Message.error(res.message || '解锁失败')
      done(false)
    }
  } catch (e) {
    Message.error('解锁失败')
    done(false)
  } finally {
    submitting.value = false
  }
}

const handleWipe = async (done) => {
  if (wipeForm.confirm_code !== 'CONFIRM-WIPE') {
    Message.warning('请输入确认码 CONFIRM-WIPE')
    done(false)
    return
  }
  submitting.value = true
  try {
    const res = await wipeDevice(selectedDevice.value.device_id, {
      wipe_type: wipeForm.wipe_type,
      confirm_code: wipeForm.confirm_code
    })
    if (res.code === 0) {
      Message.success('擦除命令已下发')
      wipeModalVisible.value = false
      loadWipeHistory()
      done(true)
    } else {
      Message.error(res.message || '擦除失败')
      done(false)
    }
  } catch (e) {
    Message.error('擦除失败')
    done(false)
  } finally {
    submitting.value = false
  }
}

onMounted(() => { loadDevices() })
</script>

<style scoped>
.pro-page-container { padding: 20px 24px; min-height: calc(100vh - 64px); background: #f5f7fa; }
.pro-breadcrumb { margin-bottom: 16px; }
.search-bar { margin-bottom: 12px; }
.device-list-card { border-radius: 8px; margin-bottom: 16px; }
.device-info-card { border-radius: 8px; margin-bottom: 16px; }
.action-card { border-radius: 8px; margin-bottom: 16px; }
.history-card { border-radius: 8px; }
.no-device-tip { text-align: center; padding: 32px 0; color: #999; }
.action-grid { display: flex; flex-direction: column; gap: 12px; }
.online-dot {
  display: inline-block;
  width: 8px; height: 8px;
  border-radius: 50%;
  margin-right: 4px;
}
.online-dot.online { background: #00b42a; }
.online-dot.offline { background: #8a8a8a; }
</style>

