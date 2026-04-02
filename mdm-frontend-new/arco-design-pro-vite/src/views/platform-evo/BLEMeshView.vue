<template>
    <Breadcrumb :items="['Home','Console','']" />


  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>平台演进</a-breadcrumb-item>
      <a-breadcrumb-item>BLE Mesh 网络</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 搜索栏 -->
    <div class="pro-search-bar">
      <a-space>
        <a-input v-model="searchKey" placeholder="搜索节点名称 / MAC" allow-clear style="width: 220px" @press-enter="loadNodes" />
        <a-select v-model="filterType" placeholder="节点类型" allow-clear style="width: 140px">
          <a-option value="relay">中继节点</a-option>
          <a-option value="leaf">叶节点</a-option>
          <a-option value="gateway">网关节点</a-option>
        </a-select>
        <a-select v-model="filterOnline" placeholder="在线状态" allow-clear style="width: 120px">
          <a-option :value="true">在线</a-option>
          <a-option :value="false">离线</a-option>
        </a-select>
        <a-button @click="loadNodes"><template #icon><icon-refresh /></template>刷新</a-button>
      </a-space>
    </div>

    <!-- 操作按钮 -->
    <div class="pro-action-bar">
      <a-space>
        <a-button type="primary" @click="showProvisionModal = true">
          <template #icon><icon-plus /></template>
          节点配网
        </a-button>
        <a-button @click="showMeshConfigModal = true">
          <template #icon><icon-settings /></template>
          网络配置
        </a-button>
        <a-button @click="initNetwork">
          <template #icon><icon-share-internal /></template>
          初始化网络
        </a-button>
        <a-button @click="reconnectNetwork">
          <template #icon><icon-refresh /></template>
          重连网络
        </a-button>
      </a-space>
    </div>

    <!-- 统计卡片 -->
    <div class="metric-cards">
      <a-row :gutter="16">
        <a-col :span="6">
          <a-statistic title="节点总数" :value="stats.total_nodes" :loading="statsLoading" animation />
        </a-col>
        <a-col :span="6">
          <a-statistic title="在线节点" :value="stats.online_nodes" :loading="statsLoading" animation>
            <template #extra><a-badge status="success" /></template>
          </a-statistic>
        </a-col>
        <a-col :span="6">
          <a-statistic title="网络深度" :value="stats.max_depth" :loading="statsLoading" animation />
        </a-col>
        <a-col :span="6">
          <a-statistic title="丢包率" :value="stats.packet_loss_rate" :loading="statsLoading" animation>
            <template #suffix>%</template>
          </a-statistic>
        </a-col>
      </a-row>
    </div>

    <!-- 网络拓扑图 + 节点列表 -->
    <a-row :gutter="16" style="margin-top: 16px">
      <a-col :span="14">
        <a-card title="网络拓扑" :header-extra-text="`共 ${topology.nodes?.length || 0} 个节点`">
          <template #extra>
            <a-space>
              <a-button size="small" @click="loadTopology"><icon-refresh /></a-button>
              <a-button size="small" @click="showTopologyHelp = true">图例</a-button>
            </a-space>
          </template>
          <div class="topology-container" ref="topologyContainer">
            <svg width="100%" height="360" viewBox="0 0 800 360">
              <!-- 边 -->
              <line v-for="(edge, i) in topology.edges" :key="'e' + i"
                :x1="edge.x1" :y1="edge.y1" :x2="edge.x2" :y2="edge.y2"
                :stroke="edge.active ? '#00B42A' : '#86909C'"
                :stroke-width="edge.active ? 2 : 1"
                :stroke-dasharray="edge.active ? '' : '4 2'"
              />
              <!-- 节点 -->
              <g v-for="(node, i) in topology.nodes" :key="'n' + i"
                :transform="`translate(${node.x}, ${node.y})`"
                @click="selectNode(node)"
                style="cursor: pointer"
              >
                <circle r="22" :fill="nodeColor(node)" stroke="#165DFF" stroke-width="1.5" />
                <text y="5" text-anchor="middle" fill="white" font-size="12" font-weight="bold">
                  {{ node.name?.substring(0, 3) || 'N' + i }}
                </text>
                <circle v-if="node.online" cx="16" cy="-16" r="5" fill="#00B42A" stroke="white" stroke-width="1" />
                <text y="-8" text-anchor="middle" fill="#4E5969" font-size="10">{{ node.rssi }} dBm</text>
              </g>
            </svg>
          </div>
        </a-card>
      </a-col>
      <a-col :span="10">
        <a-card title="选中节点详情" :bordered="false">
          <template #extra>
            <a-tag v-if="selectedNode?.online" color="green">在线</a-tag>
            <a-tag v-else color="gray">离线</a-tag>
          </template>
          <a-descriptions v-if="selectedNode" :column="1" size="small">
            <a-descriptions-item label="节点名称">{{ selectedNode.name }}</a-descriptions-item>
            <a-descriptions-item label="MAC 地址">{{ selectedNode.mac || '-' }}</a-descriptions-item>
            <a-descriptions-item label="节点类型">{{ nodeTypeText(selectedNode.type) }}</a-descriptions-item>
            <a-descriptions-item label="信号强度">{{ selectedNode.rssi }} dBm</a-descriptions-item>
            <a-descriptions-item label="父节点">{{ selectedNode.parent || '-' }}</a-descriptions-item>
            <a-descriptions-item label="子节点数">{{ selectedNode.children?.length || 0 }}</a-descriptions-item>
            <a-descriptions-item label="深度">{{ selectedNode.depth }}</a-descriptions-item>
            <a-descriptions-item label="最后活跃">{{ formatTime(selectedNode.last_seen) }}</a-descriptions-item>
          </a-descriptions>
          <a-empty v-else description="点击拓扑图中的节点查看详情" />
          <a-divider />
          <a-space v-if="selectedNode">
            <a-button size="small" type="primary" @click="configNode(selectedNode)">配置</a-button>
            <a-button size="small" status="warning" @click="reconnectNode(selectedNode)">重连</a-button>
            <a-button size="small" status="danger" @click="removeNode(selectedNode)">移除</a-button>
          </a-space>
        </a-card>
      </a-col>
    </a-row>

    <!-- 节点列表 -->
    <a-divider>节点列表</a-divider>
    <a-table :data="nodes" :loading="loading" :pagination="{ total: total, current: page, pageSize: pageSize, showTotal: true }" @page-change="onPageChange" row-key="id">
      <template #columns>
        <a-table-column title="节点名称" data-index="name" :width="160">
          <template #cell="{ record }">
            <a-link @click="selectNodeById(record)">{{ record.name }}</a-link>
          </template>
      </a-table>
        </a-table-column>
        <a-table-column title="MAC 地址" data-index="mac" :width="160" />
        <a-table-column title="类型" data-index="type" :width="100">
          <template #cell="{ record }">
            <a-tag :color="nodeTypeColor(record.type)">{{ nodeTypeText(record.type) }}</a-tag>
          </template>
        </a-table-column>
        <a-table-column title="在线" data-index="online" :width="80">
          <template #cell="{ record }">
            <a-badge :status="record.online ? 'success' : 'default'" :text="record.online ? '在线' : '离线'" />
          </template>
        </a-table-column>
        <a-table-column title="信号" data-index="rssi" :width="100">
          <template #cell="{ record }">{{ record.rssi }} dBm</template>
        </a-table-column>
        <a-table-column title="父节点" data-index="parent_name" :width="120">
          <template #cell="{ record }">{{ record.parent_name || '-' }}</template>
        </a-table-column>
        <a-table-column title="深度" data-index="depth" :width="80" />
        <a-table-column title="最后活跃" data-index="last_seen" :width="160">
          <template #cell="{ record }">{{ formatTime(record.last_seen) }}</template>
        </a-table-column>
        <a-table-column title="操作" :width="160" fixed="right">
          <template #cell="{ record }">
            <a-space>
              <a-button size="small" @click="selectNodeById(record)">查看</a-button>
              <a-button size="small" @click="configNode(record)">配置</a-button>
              <a-button size="small" status="danger" @click="removeNode(record)">移除</a-button>
            </a-space>
          </template>
        </a-table-column>
      </template>
    </a-table>

    <!-- 配网弹窗 -->
    <a-modal v-model:visible="showProvisionModal" title="节点配网" :width="480" @before-ok="handleProvision" @cancel="showProvisionModal = false">
      <a-form :model="provisionForm" layout="vertical">
        <a-form-item label="节点 MAC" required>
          <a-input v-model="provisionForm.mac" placeholder="请输入节点 MAC 地址" />
        </a-form-item>
        <a-form-item label="节点名称">
          <a-input v-model="provisionForm.name" placeholder="请输入节点名称" />
        </a-form-item>
        <a-form-item label="节点类型" required>
          <a-select v-model="provisionForm.type" placeholder="请选择类型">
            <a-option value="relay">中继节点</a-option>
            <a-option value="leaf">叶节点</a-option>
            <a-option value="gateway">网关节点</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="父节点">
          <a-select v-model="provisionForm.parent_id" placeholder="选择父节点（可选）" allow-clear>
            <a-option v-for="n in onlineNodes" :key="n.id" :value="n.id">{{ n.name }}</a-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 网络配置弹窗 -->
    <a-modal v-model:visible="showMeshConfigModal" title="BLE Mesh 网络配置" :width="500" @before-ok="handleMeshConfig" @cancel="showMeshConfigModal = false">
      <a-form :model="meshConfigForm" layout="vertical">
        <a-form-item label="网络名称">
          <a-input v-model="meshConfigForm.network_name" placeholder="网络名称" />
        </a-form-item>
        <a-form-item label="网络 ID">
          <a-input v-model="meshConfigForm.network_id" placeholder="网络 ID" />
        </a-form-item>
        <a-form-item label="最大深度">
          <a-input-number v-model="meshConfigForm.max_depth" :min="1" :max="10" style="width: 100%" />
        </a-form-item>
        <a-form-item label="广播间隔 (ms)">
          <a-input-number v-model="meshConfigForm.advertising_interval" :min="20" :step="10" style="width: 100%" />
        </a-form-item>
        <a-form-item label="扫描间隔 (ms)">
          <a-input-number v-model="meshConfigForm.scan_interval" :min="20" :step="10" style="width: 100%" />
        </a-form-item>
        <a-form-item label="重连策略">
          <a-select v-model="meshConfigForm.reconnect_policy">
            <a-option value="aggressive">激进重连</a-option>
            <a-option value="normal">普通</a-option>
            <a-option value="conservative">保守</a-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 节点配置弹窗 -->
    <a-modal v-model:visible="showNodeConfigModal" title="节点配置" :width="440" @before-ok="handleNodeConfig" @cancel="showNodeConfigModal = false">
      <a-form :model="nodeConfigForm" layout="vertical">
        <a-form-item label="节点名称">
          <a-input v-model="nodeConfigForm.name" />
        </a-form-item>
        <a-form-item label="节点类型">
          <a-select v-model="nodeConfigForm.type">
            <a-option value="relay">中继节点</a-option>
            <a-option value="leaf">叶节点</a-option>
            <a-option value="gateway">网关节点</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="父节点">
          <a-select v-model="nodeConfigForm.parent_id" placeholder="选择父节点" allow-clear>
            <a-option v-for="n in onlineNodes" :key="n.id" :value="n.id">{{ n.name }}</a-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 图例弹窗 -->
    <a-modal v-model:visible="showTopologyHelp" title="拓扑图图例" :width="400">
      <a-space direction="vertical" fill>
        <a-space>
          <svg width="30" height="30"><circle cx="15" cy="15" r="10" fill="#165DFF" stroke="#165DFF" /></svg>
          <span>网关节点</span>
        </a-space>
        <a-space>
          <svg width="30" height="30"><circle cx="15" cy="15" r="10" fill="#0FC6C2" stroke="#0FC6C2" /></svg>
          <span>中继节点</span>
        </a-space>
        <a-space>
          <svg width="30" height="30"><circle cx="15" cy="15" r="10" fill="#FF7D00" stroke="#FF7D00" /></svg>
          <span>叶节点</span>
        </a-space>
        <a-space>
          <svg width="30" height="30">
            <line x1="5" y1="15" x2="25" y2="15" stroke="#00B42A" stroke-width="2" />
          </svg>
          <span>活跃连接</span>
        </a-space>
        <a-space>
          <svg width="30" height="30">
            <line x1="5" y1="15" x2="25" y2="15" stroke="#86909C" stroke-width="1" stroke-dasharray="4 2" />
          </svg>
          <span>非活跃连接</span>
        </a-space>
      </a-space>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const API_BASE = '/api'

function getToken() {
  return localStorage.getItem('token') || ''
}

async function apiRequest(url, options = {}) {
  const res = await fetch(url, {
    ...options,
    headers: {
      'Authorization': `Bearer ${getToken()}`,
      'Content-Type': 'application/json',
      ...(options.headers || {})
    }
  })
  const data = await res.json()
  if (data.code !== 0 && data.code !== 200) throw new Error(data.message || '请求失败')
  return data
}

const nodes = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)
const loading = ref(false)
const statsLoading = ref(false)
const searchKey = ref('')
const filterType = ref('')
const filterOnline = ref('')
const stats = ref({ total_nodes: 0, online_nodes: 0, max_depth: 0, packet_loss_rate: 0 })
const topology = ref({ nodes: [], edges: [] })
const selectedNode = ref(null)
const onlineNodes = ref([])

// 配网
const showProvisionModal = ref(false)
const provisionForm = reactive({ mac: '', name: '', type: '', parent_id: '' })

// 网络配置
const showMeshConfigModal = ref(false)
const meshConfigForm = reactive({ network_name: '', network_id: '', max_depth: 5, advertising_interval: 100, scan_interval: 100, reconnect_policy: 'normal' })

// 节点配置
const showNodeConfigModal = ref(false)
const nodeConfigForm = reactive({ id: '', name: '', type: '', parent_id: '' })

const showTopologyHelp = ref(false)

function formatTime(t) {
  if (!t) return '-'
  return new Date(t).toLocaleString('zh-CN')
}

function nodeTypeText(t) {
  return { relay: '中继节点', leaf: '叶节点', gateway: '网关节点' }[t] || t
}
function nodeTypeColor(t) {
  return { relay: 'cyan', leaf: 'orange', gateway: 'blue' }[t] || 'gray'
}
function nodeColor(node) {
  if (node.type === 'gateway') return '#165DFF'
  if (node.type === 'relay') return '#0FC6C2'
  return '#FF7D00'
}

async function loadNodes() {
  loading.value = true
  try {
    const params = { page: page.value, page_size: pageSize.value }
    if (searchKey.value) params.keyword = searchKey.value
    if (filterType.value) params.type = filterType.value
    if (filterOnline.value !== '') params.online = filterOnline.value
    const qs = new URLSearchParams(params).toString()
    const res = await apiRequest(`${API_BASE}/ble-mesh/nodes?${qs}`)
    nodes.value = res.data?.list || []
    total.value = res.data?.total || 0
    onlineNodes.value = nodes.value.filter(n => n.online)
  } catch (e) {
    Message.error('加载节点列表失败')
  } finally {
    loading.value = false
  }
}

async function loadStats() {
  statsLoading.value = true
  try {
    const res = await apiRequest(`${API_BASE}/ble-mesh/stats`)
    stats.value = res.data || {}
  } catch (e) { /* ignore */ }
  finally { statsLoading.value = false }
}

async function loadTopology() {
  try {
    const res = await apiRequest(`${API_BASE}/ble-mesh/topology`)
    topology.value = res.data || { nodes: [], edges: [] }
    // 自动布局
    autoLayoutTopology()
  } catch (e) {
    Message.error('加载拓扑失败')
  }
}

function autoLayoutTopology() {
  const nodes = topology.value.nodes || []
  const w = 780, h = 320
  const cols = Math.ceil(Math.sqrt(nodes.length))
  nodes.forEach((node, i) => {
    const col = i % cols
    const row = Math.floor(i / cols)
    node.x = 50 + col * (w / cols)
    node.y = 50 + row * (h / Math.ceil(nodes.length / cols))
  })
  // 生成边
  topology.value.edges = []
  nodes.forEach(node => {
    if (node.parent_id) {
      const parent = nodes.find(n => n.id === node.parent_id)
      if (parent) {
        topology.value.edges.push({
          x1: parent.x, y1: parent.y,
          x2: node.x, y2: node.y,
          active: node.online && parent.online
        })
      }
    }
  })
}

function selectNode(node) {
  selectedNode.value = node
}

function selectNodeById(record) {
  selectedNode.value = record
}

async function handleProvision(done) {
  if (!provisionForm.mac || !provisionForm.type) {
    Message.warning('请填写必填项')
    done(false)
    return
  }
  try {
    await apiRequest(`${API_BASE}/ble-mesh/nodes/${provisionForm.mac}/provision`, {
      method: 'POST',
      body: JSON.stringify(provisionForm)
    })
    Message.success('配网成功')
    showProvisionModal.value = false
    Object.assign(provisionForm, { mac: '', name: '', type: '', parent_id: '' })
    loadNodes()
    loadStats()
    loadTopology()
    done(true)
  } catch (e) {
    Message.error('配网失败: ' + e.message)
    done(false)
  }
}

async function loadMeshConfig() {
  try {
    const res = await apiRequest(`${API_BASE}/ble-mesh/config`)
    Object.assign(meshConfigForm, res.data || {})
  } catch (e) { /* ignore */ }
}

async function handleMeshConfig(done) {
  try {
    await apiRequest(`${API_BASE}/ble-mesh/config`, {
      method: 'PUT',
      body: JSON.stringify(meshConfigForm)
    })
    Message.success('配置成功')
    showMeshConfigModal.value = false
    done(true)
  } catch (e) {
    Message.error('配置失败: ' + e.message)
    done(false)
  }
}

function configNode(node) {
  Object.assign(nodeConfigForm, { id: node.id, name: node.name, type: node.type, parent_id: node.parent_id || '' })
  showNodeConfigModal.value = true
}

async function handleNodeConfig(done) {
  try {
    await apiRequest(`${API_BASE}/ble-mesh/nodes/${nodeConfigForm.id}`, {
      method: 'PUT',
      body: JSON.stringify(nodeConfigForm)
    })
    Message.success('配置成功')
    showNodeConfigModal.value = false
    loadNodes()
    loadTopology()
    done(true)
  } catch (e) {
    Message.error('配置失败: ' + e.message)
    done(false)
  }
}

async function removeNode(node) {
  try {
    await apiRequest(`${API_BASE}/ble-mesh/nodes/${node.id}`, { method: 'DELETE' })
    Message.success('移除成功')
    if (selectedNode.value?.id === node.id) selectedNode.value = null
    loadNodes()
    loadStats()
    loadTopology()
  } catch (e) {
    Message.error('移除失败: ' + e.message)
  }
}

async function initNetwork() {
  try {
    await apiRequest(`${API_BASE}/ble-mesh/network/init`, { method: 'POST' })
    Message.success('网络初始化成功')
    loadNodes()
    loadStats()
    loadTopology()
  } catch (e) {
    Message.error('初始化失败: ' + e.message)
  }
}

async function reconnectNetwork() {
  try {
    await apiRequest(`${API_BASE}/ble-mesh/network/reconnect`, { method: 'POST' })
    Message.success('重连已触发')
    loadNodes()
    loadTopology()
  } catch (e) {
    Message.error('重连失败: ' + e.message)
  }
}

async function reconnectNode(node) {
  try {
    await apiRequest(`${API_BASE}/ble-mesh/nodes/${node.id}/reconnect`, { method: 'POST' })
    Message.success('节点重连已触发')
    loadNodes()
    loadTopology()
  } catch (e) {
    Message.error('重连失败: ' + e.message)
  }
}

function onPageChange(p) {
  page.value = p
  loadNodes()
}

onMounted(() => {
  loadNodes()
  loadStats()
  loadTopology()
  loadMeshConfig()
})
</script>

<style scoped>
.topology-container {
  background: #F7F8FA;
  border-radius: 4px;
  overflow: hidden;
}
</style>
