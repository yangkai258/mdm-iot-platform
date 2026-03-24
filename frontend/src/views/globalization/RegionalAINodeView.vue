<template>
  <div class="page-container">
    <div class="search-form">
      <a-form :model="form" layout="inline">
        <a-form-item label="区域">
          <a-select v-model="filters.regionCode" placeholder="请选择" allow-clear style="width: 160px">
            <a-option v-for="r in regions" :key="r.region_code" :value="r.region_code">{{ r.region_name }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="状态">
          <a-select v-model="filters.status" placeholder="请选择" allow-clear style="width: 120px">
            <a-option value="online">在线</a-option>
            <a-option value="offline">离线</a-option>
            <a-option value="standby">备用</a-option>
          </a-select>
        </a-form-item>
        <a-form-item>
          <a-button type="primary" @click="loadNodes">搜索</a-button>
          <a-button @click="handleReset">重置</a-button>
        </a-form-item>
      </a-form>
    </div>
    <div class="toolbar">
      <a-button type="primary" @click="openRegisterModal">注册新节点</a-button>
    </div>
    <a-table :columns="columns" :data="filteredNodes" :loading="loading" :pagination="pagination" row-key="id">
      <template #nodeId="{ record }">
        <span class="mono">{{ record.node_id || record.id }}</span>
      </template>
      <template #model="{ record }">
        <span>{{ record.model || 'mini-claw' }}</span>
      </template>
      <template #status="{ record }">
        <a-badge :color="nodeStatusColor(record.health_status)" :text="nodeStatusLabel(record.health_status)" />
      </template>
      <template #qps="{ record }">
        <span v-if="record.health_status === 'online'">{{ record.qps || 0 }}/{{ record.qps_limit || 100 }}</span>
        <span v-else class="na">--</span>
      </template>
      <template #load="{ record }">
        <span v-if="record.health_status === 'online'">
          <a-progress :percent="record.load_factor || 0" size="small" :show-text="true" />
        </span>
        <span v-else class="na">--</span>
      </template>
      <template #actions="{ record }">
        <a-button type="text" size="small" @click="openDetailModal(record)">详情</a-button>
        <a-button type="text" size="small" status="danger" @click="handleDeregister(record)">注销</a-button>
      </template>
    </a-table>
    <a-modal v-model:visible="registerVisible" title="注册新节点" :width="480" @before-ok="handleRegister">
      <a-form :model="form" label-col-flex="100px">
        <a-form-item label="节点 ID" required>
          <a-input v-model="form.node_id" placeholder="如：node-001" />
        </a-form-item>
        <a-form-item label="所属区域" required>
          <a-select v-model="form.region_code" placeholder="选择区域">
            <a-option v-for="r in regions" :key="r.region_code" :value="r.region_code">{{ r.region_name }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="模型">
          <a-input v-model="form.model" placeholder="如：mini-claw" />
        </a-form-item>
        <a-form-item label="端点地址" required>
          <a-input v-model="form.endpoint" placeholder="如：https://ai-node-001.example.com" />
        </a-form-item>
        <a-form-item label="QPS 限制">
          <a-input-number v-model="form.qps_limit" placeholder="100" :min="1" style="width: 100%" />
        </a-form-item>
      </a-form>
      <template #footer>
        <a-button @click="registerVisible = false">取消</a-button>
        <a-button type="primary" @click="handleRegister">确定</a-button>
      </template>
    </a-modal>
    <a-modal v-model:visible="detailVisible" title="节点详情" :width="520">
      <div class="detail-grid" v-if="detailNode">
        <div class="detail-item">
          <span class="detail-label">节点 ID</span>
          <span class="detail-value mono">{{ detailNode.node_id || detailNode.id }}</span>
        </div>
        <div class="detail-item">
          <span class="detail-label">模型</span>
          <span class="detail-value">{{ detailNode.model || 'mini-claw' }}</span>
        </div>
        <div class="detail-item">
          <span class="detail-label">区域</span>
          <span class="detail-value">{{ detailNode.region_code }}</span>
        </div>
        <div class="detail-item">
          <span class="detail-label">状态</span>
          <span class="detail-value">
            <a-badge :color="nodeStatusColor(detailNode.health_status)" :text="nodeStatusLabel(detailNode.health_status)" />
          </span>
        </div>
        <div class="detail-item">
          <span class="detail-label">QPS</span>
          <span class="detail-value">{{ detailNode.qps || 0 }} / {{ detailNode.qps_limit || 100 }}</span>
        </div>
        <div class="detail-item">
          <span class="detail-label">负载</span>
          <span class="detail-value">{{ detailNode.load_factor || 0 }}%</span>
        </div>
        <div class="detail-item">
          <span class="detail-label">端点</span>
          <span class="detail-value">{{ detailNode.endpoint || '-' }}</span>
        </div>
        <div class="detail-item">
          <span class="detail-label">最后心跳</span>
          <span class="detail-value">{{ formatDate(detailNode.last_heartbeat) }}</span>
        </div>
      </div>
      <template #footer></template>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { Message, Modal } from '@arco-design/web-vue'
import { getRegions, getAINodes, registerAINode, deregisterAINode } from '@/api/globalization'
import dayjs from 'dayjs'

const loading = ref(false)
const nodes = ref([])
const regions = ref([])
const filters = reactive({ regionCode: '', status: '' })
const pagination = reactive({ pageSize: 20, current: 1, total: 0 })

const columns = [
  { title: '节点ID', slotName: 'nodeId', width: 140 },
  { title: '模型', slotName: 'model', width: 120 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: 'QPS', slotName: 'qps', width: 110 },
  { title: '负载', slotName: 'load', width: 140 },
  { title: '操作', slotName: 'actions', width: 120 }
]

const filteredNodes = computed(() => {
  return nodes.value.filter(n => {
    if (filters.regionCode && n.region_code !== filters.regionCode) return false
    if (filters.status && n.health_status !== filters.status) return false
    return true
  })
})

const registerVisible = ref(false)
const detailVisible = ref(false)
const detailNode = ref(null)
const form = reactive({
  node_id: '',
  region_code: '',
  model: 'mini-claw',
  endpoint: '',
  qps_limit: 100
})

function nodeStatusColor(status) {
  const map = { online: 'green', offline: 'red', standby: 'orange' }
  return map[status] || 'default'
}

function nodeStatusLabel(status) {
  const map = { online: '在线', offline: '离线', standby: '备用' }
  return map[status] || status
}

function formatDate(date) {
  if (!date) return '-'
  return dayjs(date).format('YYYY-MM-DD HH:mm:ss')
}

async function loadNodes() {
  loading.value = true
  try {
    const res = await getAINodes(filters)
    nodes.value = res.data || res || []
    pagination.total = filteredNodes.value.length
  } catch (e) {
    console.error('加载AI节点失败', e)
  } finally {
    loading.value = false
  }
}

function handleReset() {
  filters.regionCode = ''
  filters.status = ''
  pagination.current = 1
  loadNodes()
}

onMounted(async () => {
  try {
    const [nodesRes, regionsRes] = await Promise.all([
      getAINodes(),
      getRegions()
    ])
    nodes.value = nodesRes.data || nodesRes || []
    regions.value = regionsRes.data || regionsRes || []
    pagination.total = nodes.value.length
  } catch (e) {
    console.error('加载数据失败', e)
  }
})

function openRegisterModal() {
  Object.assign(form, { node_id: '', region_code: '', model: 'mini-claw', endpoint: '', qps_limit: 100 })
  registerVisible.value = true
}

function openDetailModal(record) {
  detailNode.value = record
  detailVisible.value = true
}

async function handleRegister() {
  if (!form.node_id || !form.region_code || !form.endpoint) {
    Message.warning('请填写必填项')
    return
  }
  try {
    await registerAINode(form)
    Message.success('注册成功')
    registerVisible.value = false
    await loadNodes()
  } catch (e) {
    Message.error('注册失败')
  }
}

function handleDeregister(record) {
  Modal.warning({
    title: '确认注销',
    content: `确定要注销节点「${record.node_id || record.id}」吗？`,
    okText: '注销',
    onOk: async () => {
      try {
        await deregisterAINode(record.id)
        Message.success('注销成功')
        await loadNodes()
      } catch (e) {
        Message.error('注销失败')
      }
    }
  })
}
</script>

<style scoped>
.page-container { background: #fff; border-radius: 4px; padding: 20px; }
.search-form { margin-bottom: 16px; padding: 16px; background: #f7f8fa; border-radius: 4px; }
.toolbar { margin-bottom: 16px; }
.mono { font-family: monospace; font-size: 12px; }
.na { color: var(--color-text-3); }
.detail-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 16px; }
.detail-item { display: flex; flex-direction: column; gap: 4px; }
.detail-label { font-size: 12px; color: var(--color-text-3); }
.detail-value { font-size: 14px; word-break: break-all; }
</style>
