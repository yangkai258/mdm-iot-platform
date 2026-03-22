<template>
  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>技术架构</a-breadcrumb-item>
      <a-breadcrumb-item>BLE Mesh 网络</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 统计卡片 -->
    <a-row :gutter="16" class="stats-row">
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="Mesh 网络数" :value="stats.networks" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="在线设备" :value="stats.onlineDevices" :value-style="{ color: '#52c41a' }" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="离线设备" :value="stats.offlineDevices" :value-style="{ color: '#ff4d4f' }" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="消息速率" :value="stats.msgRate" suffix="msg/s" />
        </a-card>
      </a-col>
    </a-row>

    <!-- 标签页 -->
    <a-tabs v-model:active-key="activeTab" class="pro-tabs">
      <a-tab-pane key="devices" title="设备列表">
        <div class="pro-search-bar">
          <a-space>
            <a-input-search v-model="searchKeyword" placeholder="搜索设备ID或MAC" style="width: 280px" @search="loadDevices" search-button />
            <a-select v-model="filterOnline" placeholder="在线状态" style="width: 120px" allow-clear>
              <a-option :value="true">在线</a-option>
              <a-option :value="false">离线</a-option>
            </a-select>
          </a-space>
        </div>
        <div class="pro-action-bar">
          <a-space>
            <a-button type="primary" @click="showNetworkModal">创建网络</a-button>
            <a-button @click="loadDevices">刷新</a-button>
          </a-space>
        </div>
        <div class="pro-content-area">
          <a-table :columns="deviceColumns" :data="filteredDevices" :loading="loading" :pagination="pagination" @change="handleTableChange" row-key="id">
            <template #deviceId="{ record }">
              <a-space>
                <a-avatar :size="24" :style="{ backgroundColor: record.is_online ? '#52c41a' : '#8c8c8c' }">
                  {{ record.device_id.charAt(0) }}
                </a-avatar>
                <span>{{ record.device_id }}</span>
              </a-space>
            </template>
            <template #isOnline="{ record }">
              <a-badge :status="record.is_online ? 'success' : 'default'" :text="record.is_online ? '在线' : '离线'" />
            </template>
            <template #rssi="{ record }">
              <span :style="{ color: record.rssi > -60 ? '#52c41a' : record.rssi > -80 ? '#faad14' : '#ff4d4f' }">
                {{ record.rssi }} dBm
              </span>
            </template>
            <template #role="{ record }">
              <a-tag :color="getRoleColor(record.role)">{{ getRoleText(record.role) }}</a-tag>
            </template>
            <template #actions="{ record }">
              <a-space>
                <a-button v-if="record.role !== 'relay'" type="text" size="small" @click="handleUpgradeRole(record)">升级Relay</a-button>
                <a-button v-if="record.is_online" type="text" size="small" @click="showConnectionModal(record)">连接控制</a-button>
                <a-button type="text" size="small" status="danger" @click="handleRemoveDevice(record)">移除</a-button>
              </a-space>
            </template>
          </a-table>
        </div>
      </a-tab-pane>

      <a-tab-pane key="topology" title="网络拓扑">
        <div class="pro-search-bar">
          <a-space>
            <a-select v-model="selectedNetwork" placeholder="选择Mesh网络" style="width: 240px" @change="loadTopology">
              <a-option v-for="n in networks" :key="n.id" :value="n.id">{{ n.name }}</a-option>
            </a-select>
            <a-button @click="loadTopology">刷新拓扑</a-button>
          </a-space>
        </div>
        <div class="topology-area">
          <a-empty description="网络拓扑图" v-if="topologyNodes.length === 0">
            <template #image><icon-wifi :size="48" /></template>
          </a-empty>
          <div v-else class="topology-canvas">
            <svg class="topology-lines" :width="canvasWidth" :height="canvasHeight">
              <line v-for="line in topologyLines" :key="line.key"
                :x1="line.x1" :y1="line.y1" :x2="line.x2" :y2="line.y2"
                :stroke="line.strength > -60 ? '#52c41a' : line.strength > -80 ? '#faad14' : '#ff4d4f'"
                stroke-width="2" stroke-opacity="0.5" />
            </svg>
            <div v-for="node in topologyNodes" :key="node.id" class="topology-node"
              :style="{ left: node.x + 'px', top: node.y + 'px' }">
              <div class="node-icon" :class="{ online: node.is_online, offline: !node.is_online, relay: node.role === 'relay' }">
                <icon-wifi :size="20" />
              </div>
              <div class="node-label">{{ node.device_id }}</div>
              <div class="node-rssi">{{ node.rssi }}dBm</div>
            </div>
          </div>
          <div class="topology-legend">
            <span class="legend-item"><span class="dot online"></span> 在线</span>
            <span class="legend-item"><span class="dot offline"></span> 离线</span>
            <span class="legend-item"><span class="dot relay"></span> Relay节点</span>
          </div>
        </div>
        <a-card class="connections-card" v-if="selectedNetwork && connections.length">
          <template #title><span>连接详情</span></template>
          <a-table :columns="connColumns" :data="connections" :pagination="false" row-key="id" size="small">
            <template #strength="{ record }">
              <a-progress :percent="Math.round((record.rssi + 100) * -1 * -1)" :stroke-width="6" size="small" />
            </template>
          </a-table>
        </a-card>
      </a-tab-pane>
    </a-tabs>

    <!-- 创建网络弹窗 -->
    <a-modal v-model:visible="networkModalVisible" title="创建 Mesh 网络" @ok="handleCreateNetwork" :confirm-loading="networkCreating" :width="480">
      <a-form :model="networkForm" layout="vertical">
        <a-form-item label="网络名称" required>
          <a-input v-model="networkForm.name" placeholder="例如: Home-Mesh-01" />
        </a-form-item>
        <a-form-item label="网络ID" required>
          <a-input v-model="networkForm.mesh_id" placeholder="6位十六进制, 例如: 0x0001" />
        </a-form-item>
        <a-form-item label="安全模式">
          <a-select v-model="networkForm.security" placeholder="选择安全模式">
            <a-option value="secure">安全模式</a-option>
            <a-option value="unsecure">非安全模式</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="描述">
          <a-textarea v-model="networkForm.description" placeholder="网络描述" :rows="2" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 连接控制弹窗 -->
    <a-modal v-model:visible="connectionModalVisible" title="连接控制" @ok="handleConnectionControl" :confirm-loading="connectionLoading" :width="480">
      <a-descriptions :column="1" bordered size="small">
        <a-descriptions-item label="设备ID">{{ currentDevice?.device_id }}</a-descriptions-item>
        <a-descriptions-item label="MAC地址">{{ currentDevice?.mac_address }}</a-descriptions-item>
        <a-descriptions-item label="当前角色"><a-tag :color="getRoleColor(currentDevice?.role)">{{ getRoleText(currentDevice?.role) }}</a-tag></a-descriptions-item>
        <a-descriptions-item label="信号强度">{{ currentDevice?.rssi }} dBm</a-descriptions-item>
      </a-descriptions>
      <a-divider>连接操作</a-divider>
      <a-form layout="vertical">
        <a-form-item label="操作类型">
          <a-radio-group v-model="connectionForm.action">
            <a-radio value="connect">连接</a-radio>
            <a-radio value="disconnect">断开</a-radio>
            <a-radio value="reconnect">重连</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item label="目标节点" v-if="connectionForm.action !== 'disconnect'">
          <a-select v-model="connectionForm.target_id" placeholder="选择目标节点" allow-clear>
            <a-option v-for="d in onlineNeighborDevices" :key="d.device_id" :value="d.device_id">{{ d.device_id }} ({{ d.rssi }}dBm)</a-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import axios from 'axios'

const API_BASE = '/api/v1/mesh'
const activeTab = ref('devices')
const loading = ref(false)
const searchKeyword = ref('')
const filterOnline = ref(null)
const selectedNetwork = ref(null)
const networkModalVisible = ref(false)
const connectionModalVisible = ref(false)
const networkCreating = ref(false)
const connectionLoading = ref(false)
const networkForm = reactive({ name: '', mesh_id: '', security: 'secure', description: '' })
const connectionForm = reactive({ action: 'connect', target_id: '' })
const currentDevice = ref(null)
const stats = reactive({ networks: 0, onlineDevices: 0, offlineDevices: 0, msgRate: 0 })
const devices = ref([])
const networks = ref([])
const topologyNodes = ref([])
const connections = ref([])
const canvasWidth = 800
const canvasHeight = 400
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })

const deviceColumns = [
  { title: '设备ID', slotName: 'deviceId', width: 160 },
  { title: 'MAC地址', dataIndex: 'mac_address', width: 160 },
  { title: '在线状态', slotName: 'isOnline', width: 100 },
  { title: 'RSSI', slotName: 'rssi', width: 100 },
  { title: '角色', slotName: 'role', width: 100 },
  { title: '所属网络', dataIndex: 'network_name', width: 140 },
  { title: '最后通信', dataIndex: 'last_comm_time', width: 180 },
  { title: '操作', slotName: 'actions', width: 200 }
]

const connColumns = [
  { title: '源设备', dataIndex: 'source', width: 120 },
  { title: '目标设备', dataIndex: 'target', width: 120 },
  { title: '信号强度', slotName: 'strength', width: 180 },
  { title: '状态', dataIndex: 'status', width: 100 }
]

const filteredDevices = computed(() => {
  let result = devices.value
  if (filterOnline.value !== null) result = result.filter(d => d.is_online === filterOnline.value)
  if (searchKeyword.value) {
    const kw = searchKeyword.value.toLowerCase()
    result = result.filter(d => d.device_id.toLowerCase().includes(kw) || d.mac_address.toLowerCase().includes(kw))
  }
  return result
})

const onlineNeighborDevices = computed(() =>
  devices.value.filter(d => d.is_online && d.device_id !== currentDevice.value?.device_id)
)

const topologyLines = computed(() => {
  const lines = []
  const nodeMap = {}
  topologyNodes.value.forEach(n => { nodeMap[n.device_id] = n })
  topologyNodes.value.forEach(node => {
    node._connections?.forEach(conn => {
      const target = nodeMap[conn.target]
      if (target) {
        lines.push({
          key: `${node.device_id}-${conn.target}`,
          x1: node.x + 30, y1: node.y + 30, x2: target.x + 30, y2: target.y + 30, strength: conn.strength
        })
      }
    })
  })
  return lines
})

const loadDevices = async () => {
  loading.value = true
  try {
    const res = await axios.get(`${API_BASE}/devices`)
    if (res.data.code === 0) {
      devices.value = res.data.data.list
      pagination.total = res.data.data.pagination?.total || devices.value.length
    }
  } catch {
    devices.value = [
      { id: 1, device_id: 'MESH-001', mac_address: 'AA:BB:CC:DD:EE:01', is_online: true, rssi: -55, role: 'relay', network_name: 'Home-Mesh-01', last_comm_time: '2026-03-22 17:00:00', _connections: [{ target: 'MESH-002', strength: -55 }] },
      { id: 2, device_id: 'MESH-002', mac_address: 'AA:BB:CC:DD:EE:02', is_online: true, rssi: -68, role: 'node', network_name: 'Home-Mesh-01', last_comm_time: '2026-03-22 16:58:00', _connections: [{ target: 'MESH-003', strength: -68 }] },
      { id: 3, device_id: 'MESH-003', mac_address: 'AA:BB:CC:DD:EE:03', is_online: true, rssi: -75, role: 'node', network_name: 'Home-Mesh-01', last_comm_time: '2026-03-22 16:55:00', _connections: [] },
      { id: 4, device_id: 'MESH-004', mac_address: 'AA:BB:CC:DD:EE:04', is_online: false, rssi: -90, role: 'node', network_name: 'Home-Mesh-01', last_comm_time: '2026-03-21 10:00:00', _connections: [] },
      { id: 5, device_id: 'MESH-005', mac_address: 'AA:BB:CC:DD:EE:05', is_online: true, rssi: -62, role: 'relay', network_name: 'Home-Mesh-02', last_comm_time: '2026-03-22 17:01:00', _connections: [] }
    ]
    pagination.total = devices.value.length
    Message.warning('使用模拟数据')
  } finally { loading.value = false }
  updateStats()
}

const loadNetworks = async () => {
  try {
    const res = await axios.get(`${API_BASE}/networks`)
    if (res.data.code === 0) networks.value = res.data.data.list
  } catch {
    networks.value = [
      { id: 1, name: 'Home-Mesh-01', mesh_id: '0x0001', security: 'secure', description: '家庭主网络' },
      { id: 2, name: 'Home-Mesh-02', mesh_id: '0x0002', security: 'secure', description: '家庭副网络' }
    ]
  }
  stats.networks = networks.value.length
}

const loadTopology = () => {
  if (!selectedNetwork.value) return
  topologyNodes.value = [
    { id: 1, device_id: 'MESH-001', is_online: true, rssi: -55, role: 'relay', x: 120, y: 80, _connections: [{ target: 'MESH-002', strength: -55 }, { target: 'MESH-003', strength: -58 }] },
    { id: 2, device_id: 'MESH-002', is_online: true, rssi: -68, role: 'node', x: 300, y: 160, _connections: [{ target: 'MESH-001', strength: -68 }, { target: 'MESH-003', strength: -72 }] },
    { id: 3, device_id: 'MESH-003', is_online: true, rssi: -75, role: 'node', x: 480, y: 80, _connections: [{ target: 'MESH-001', strength: -75 }, { target: 'MESH-002', strength: -72 }] },
    { id: 4, device_id: 'MESH-004', is_online: false, rssi: -90, role: 'node', x: 300, y: 280, _connections: [] }
  ]
  connections.value = [
    { id: 1, source: 'MESH-001', target: 'MESH-002', rssi: -55, status: 'connected' },
    { id: 2, source: 'MESH-001', target: 'MESH-003', rssi: -58, status: 'connected' },
    { id: 3, source: 'MESH-002', target: 'MESH-003', rssi: -72, status: 'connected' }
  ]
}

const updateStats = () => {
  stats.onlineDevices = devices.value.filter(d => d.is_online).length
  stats.offlineDevices = devices.value.filter(d => !d.is_online).length
  stats.msgRate = 128
}

const handleTableChange = (pag) => { pagination.current = pag.current }

const showNetworkModal = () => { networkModalVisible.value = true }

const handleCreateNetwork = () => {
  if (!networkForm.name || !networkForm.mesh_id) { Message.warning('请填写必填项'); return }
  networkCreating.value = true
  setTimeout(() => {
    networks.value.push({ id: Date.now(), name: networkForm.name, mesh_id: networkForm.mesh_id, security: networkForm.security, description: networkForm.description })
    networkModalVisible.value = false
    networkCreating.value = false
    stats.networks = networks.value.length
    Message.success('Mesh 网络创建成功')
    Object.assign(networkForm, { name: '', mesh_id: '', security: 'secure', description: '' })
  }, 1000)
}

const showConnectionModal = (record) => {
  currentDevice.value = { ...record }
  connectionForm.action = 'connect'
  connectionForm.target_id = ''
  connectionModalVisible.value = true
}

const handleConnectionControl = () => {
  connectionLoading.value = true
  setTimeout(() => {
    connectionLoading.value = false
    connectionModalVisible.value = false
    Message.success(`设备 ${currentDevice.value.device_id} ${connectionForm.action === 'disconnect' ? '已断开' : '连接请求已发送'}`)
  }, 800)
}

const handleUpgradeRole = (record) => {
  const idx = devices.value.findIndex(d => d.id === record.id)
  if (idx !== -1) devices.value[idx].role = 'relay'
  Message.success(`设备 ${record.device_id} 已升级为 Relay 节点`)
}

const handleRemoveDevice = (record) => {
  devices.value = devices.value.filter(d => d.id !== record.id)
  Message.success('设备已从网络移除')
  updateStats()
}

const getRoleColor = (role) => ({ relay: 'blue', node: 'green', edge: 'orange' }[role] || 'default')
const getRoleText = (role) => ({ relay: 'Relay', node: '节点', edge: '边缘' }[role] || '未知')

onMounted(() => { loadDevices(); loadNetworks() })
</script>

<style scoped>
.pro-page-container { padding: 20px 24px; min-height: calc(100vh - 64px); background: #f5f7fa; }
.pro-breadcrumb { margin-bottom: 16px; }
.stats-row { margin-bottom: 16px; }
.stat-card { border-radius: 8px; text-align: center; }
.pro-tabs { background: #fff; border-radius: 8px; padding: 16px; box-shadow: 0 1px 3px rgba(0,0,0,0.04); }
.pro-search-bar { margin-bottom: 12px; }
.pro-action-bar { margin-bottom: 16px; }
.pro-content-area { margin-bottom: 16px; }
.topology-area { position: relative; background: #fafafa; border: 1px solid #e8e8e8; border-radius: 8px; min-height: 420px; margin-bottom: 16px; overflow: hidden; }
.topology-canvas { position: relative; width: 100%; height: 420px; }
.topology-lines { position: absolute; top: 0; left: 0; width: 100%; height: 100%; pointer-events: none; }
.topology-node { position: absolute; transform: translate(-50%, -50%); text-align: center; z-index: 1; }
.node-icon { width: 40px; height: 40px; border-radius: 50%; display: flex; align-items: center; justify-content: center; margin: 0 auto; background: #8c8c8c; color: #fff; }
.node-icon.online { background: #52c41a; }
.node-icon.offline { background: #8c8c8c; }
.node-icon.relay { background: #165dff; box-shadow: 0 0 8px rgba(22, 93, 255, 0.5); }
.node-label { font-size: 12px; margin-top: 4px; font-weight: 500; }
.node-rssi { font-size: 11px; color: #999; }
.topology-legend { position: absolute; bottom: 12px; right: 12px; background: rgba(255,255,255,0.9); padding: 8px 12px; border-radius: 4px; border: 1px solid #e8e8e8; }
.legend-item { display: inline-flex; align-items: center; margin-right: 16px; font-size: 12px; }
.dot { display: inline-block; width: 10px; height: 10px; border-radius: 50%; margin-right: 4px; }
.dot.online { background: #52c41a; }
.dot.offline { background: #8c8c8c; }
.dot.relay { background: #165dff; }
.connections-card { border-radius: 8px; margin-top: 16px; }
</style>
