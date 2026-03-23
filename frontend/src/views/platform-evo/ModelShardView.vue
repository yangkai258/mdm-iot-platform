<template>
  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>平台演进</a-breadcrumb-item>
      <a-breadcrumb-item>模型分片</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 搜索栏 -->
    <div class="pro-search-bar">
      <a-space>
        <a-input v-model="searchKey" placeholder="搜索分片名称 / 版本" allow-clear style="width: 240px" @press-enter="loadShards" />
        <a-select v-model="filterVersion" placeholder="版本" allow-clear style="width: 140px">
          <a-option v-for="v in versionOptions" :key="v" :value="v">{{ v }}</a-option>
        </a-select>
        <a-button @click="loadShards"><template #icon><icon-refresh /></template>刷新</a-button>
      </a-space>
    </div>

    <!-- 操作按钮 -->
    <div class="pro-action-bar">
      <a-space>
        <a-button type="primary" @click="showCreateModal = true">
          <template #icon><icon-plus /></template>
          创建分片
        </a-button>
        <a-button @click="showVersionPanel = true">
          <template #icon><icon-history /></template>
          版本管理
        </a-button>
      </a-space>
    </div>

    <!-- 统计卡片 -->
    <div class="metric-cards">
      <a-row :gutter="16">
        <a-col :span="6">
          <a-statistic title="总分片数" :value="stats.total_shards" :loading="statsLoading" animation />
        </a-col>
        <a-col :span="6">
          <a-statistic title="已加载" :value="stats.loaded_shards" :loading="statsLoading" animation />
        </a-col>
        <a-col :span="6">
          <a-statistic title="版本数" :value="stats.total_versions" :loading="statsLoading" animation />
        </a-col>
        <a-col :span="6">
          <a-statistic title="设备覆盖" :value="stats.devices_covered" :loading="statsLoading" animation>
            <template #suffix>台</template>
          </a-statistic>
        </a-col>
      </a-row>
    </div>

    <!-- 分片列表 -->
    <a-table :data="shards" :loading="loading" :pagination="{ total: total, current: page, pageSize: pageSize, showTotal: true, showPageSize: true }" @page-change="onPageChange" row-key="id">
      <template #columns>
        <a-table-column title="分片名称" data-index="name" :width="180">
          <template #cell="{ record }">
            <a-link @click="viewDetail(record)">{{ record.name }}</a-link>
          </template>
        </a-table-column>
        <a-table-column title="版本" data-index="version" :width="100">
          <template #cell="{ record }">
            <a-tag color="arcoblue" size="small">{{ record.version }}</a-tag>
          </template>
        </a-table-column>
        <a-table-column title="大小" data-index="size_mb" :width="100">
          <template #cell="{ record }">{{ record.size_mb }} MB</template>
        </a-table-column>
        <a-table-column title="层范围" data-index="layer_range" :width="120">
          <template #cell="{ record }">{{ record.layer_start }} - {{ record.layer_end }}</template>
        </a-table-column>
        <a-table-column title="加载设备" data-index="loaded_device" :width="140">
          <template #cell="{ record }">{{ record.loaded_device || '-' }}</template>
        </a-table-column>
        <a-table-column title="加载状态" data-index="load_status" :width="100">
          <template #cell="{ record }">
            <a-badge :status="loadStatusBadge(record.load_status)" :text="loadStatusText(record.load_status)" />
          </template>
        </a-table-column>
        <a-table-column title="创建时间" data-index="created_at" :width="160">
          <template #cell="{ record }">{{ formatTime(record.created_at) }}</template>
        </a-table-column>
        <a-table-column title="操作" :width="200" fixed="right">
          <template #cell="{ record }">
            <a-space>
              <a-button size="small" @click="viewDetail(record)">详情</a-button>
              <a-button v-if="record.load_status === 'idle'" size="small" type="primary" @click="showLoadModal(record)">加载</a-button>
              <a-button v-if="record.load_status === 'loaded'" size="small" status="warning" @click="unloadShard(record)">卸载</a-button>
              <a-button size="small" @click="editShard(record)">编辑</a-button>
              <a-button size="small" status="danger" @click="removeShard(record)">删除</a-button>
            </a-space>
          </template>
        </a-table-column>
      </template>
    </a-table>

    <!-- 版本管理面板 -->
    <a-drawer v-model:visible="showVersionPanel" title="版本管理" :width="680" @close="showVersionPanel = false">
      <template #footer>
        <a-space>
          <a-button @click="showVersionPanel = false">关闭</a-button>
          <a-button type="primary" @click="showCreateVersionModal = true">
            <template #icon><icon-plus /></template>
            创建版本
          </a-button>
        </a-space>
      </template>
      <a-table :data="versions" :loading="versionsLoading" row-key="id" size="small">
        <template #columns>
          <a-table-column title="版本号" data-index="version" :width="100" />
          <a-table-column title="分片数" data-index="shard_count" :width="80" />
          <a-table-column title="总大小" data-index="total_size_mb" :width="100">
            <template #cell="{ record }">{{ record.total_size_mb }} MB</template>
          </a-table-column>
          <a-table-column title="状态" data-index="status" :width="100">
            <template #cell="{ record }">
              <a-tag :color="versionStatusColor(record.status)">{{ versionStatusText(record.status) }}</a-tag>
            </template>
          </a-table-column>
          <a-table-column title="创建时间" data-index="created_at" :width="160">
            <template #cell="{ record }">{{ formatTime(record.created_at) }}</template>
          </a-table-column>
          <a-table-column title="操作" :width="120">
            <template #cell="{ record }">
              <a-space>
                <a-button v-if="record.status === 'draft'" size="small" type="primary" @click="publishVersion(record)">发布</a-button>
                <a-button size="small" @click="viewVersionDetail(record)">详情</a-button>
              </a-space>
            </template>
          </a-table-column>
        </template>
      </a-table>
    </a-drawer>

    <!-- 创建分片弹窗 -->
    <a-modal v-model:visible="showCreateModal" title="创建分片" :width="500" @before-ok="handleCreate" @cancel="showCreateModal = false">
      <a-form :model="createForm" layout="vertical">
        <a-form-item label="分片名称" required>
          <a-input v-model="createForm.name" placeholder="请输入分片名称" />
        </a-form-item>
        <a-form-item label="版本" required>
          <a-select v-model="createForm.version" placeholder="请选择版本">
            <a-option v-for="v in versionOptions" :key="v" :value="v">{{ v }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="起始层">
          <a-input-number v-model="createForm.layer_start" :min="0" placeholder="起始层索引" style="width: 100%" />
        </a-form-item>
        <a-form-item label="结束层">
          <a-input-number v-model="createForm.layer_end" :min="0" placeholder="结束层索引" style="width: 100%" />
        </a-form-item>
        <a-form-item label="描述">
          <a-textarea v-model="createForm.description" placeholder="分片描述" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 创建版本弹窗 -->
    <a-modal v-model:visible="showCreateVersionModal" title="创建版本" :width="460" @before-ok="handleCreateVersion" @cancel="showCreateVersionModal = false">
      <a-form :model="createVersionForm" layout="vertical">
        <a-form-item label="版本号" required>
          <a-input v-model="createVersionForm.version" placeholder="如 v1.0.0" />
        </a-form-item>
        <a-form-item label="描述">
          <a-textarea v-model="createVersionForm.description" placeholder="版本描述" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 加载分片弹窗 -->
    <a-modal v-model:visible="showLoadModalFlag" title="加载分片到设备" :width="440" @before-ok="handleLoad" @cancel="showLoadModalFlag = false">
      <a-form :model="loadForm" layout="vertical">
        <a-form-item label="分片">
          <a-input :model-value="loadForm.shard_name" disabled />
        </a-form-item>
        <a-form-item label="目标设备" required>
          <a-select v-model="loadForm.device_id" placeholder="请选择设备">
            <a-option v-for="d in devices" :key="d.id" :value="d.id">{{ d.name }}</a-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 详情抽屉 -->
    <a-drawer v-model:visible="showDetailDrawer" :title="`分片详情: ${currentShard?.name}`" :width="520" @close="showDetailDrawer = false">
      <a-descriptions :column="2" bordered>
        <a-descriptions-item label="分片名称">{{ currentShard?.name }}</a-descriptions-item>
        <a-descriptions-item label="版本">{{ currentShard?.version }}</a-descriptions-item>
        <a-descriptions-item label="大小">{{ currentShard?.size_mb }} MB</a-descriptions-item>
        <a-descriptions-item label="层范围">{{ currentShard?.layer_start }} - {{ currentShard?.layer_end }}</a-descriptions-item>
        <a-descriptions-item label="加载状态">
          <a-badge :status="loadStatusBadge(currentShard?.load_status)" :text="loadStatusText(currentShard?.load_status)" />
        </a-descriptions-item>
        <a-descriptions-item label="加载设备">{{ currentShard?.loaded_device || '-' }}</a-descriptions-item>
        <a-descriptions-item label="创建时间" :span="2">{{ formatTime(currentShard?.created_at) }}</a-descriptions-item>
        <a-descriptions-item label="描述" :span="2">{{ currentShard?.description || '-' }}</a-descriptions-item>
      </a-descriptions>
    </a-drawer>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const API_BASE = '/api/v1'

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

const shards = ref([])
const devices = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)
const loading = ref(false)
const statsLoading = ref(false)
const searchKey = ref('')
const filterVersion = ref('')
const versionOptions = ref([])
const stats = ref({ total_shards: 0, loaded_shards: 0, total_versions: 0, devices_covered: 0 })

// 版本
const showVersionPanel = ref(false)
const versions = ref([])
const versionsLoading = ref(false)
const showCreateVersionModal = ref(false)
const createVersionForm = reactive({ version: '', description: '' })

// 创建
const showCreateModal = ref(false)
const createForm = reactive({ name: '', version: '', layer_start: 0, layer_end: 0, description: '' })

// 加载
const showLoadModalFlag = ref(false)
const loadForm = reactive({ shard_id: '', shard_name: '', device_id: '' })

// 详情
const showDetailDrawer = ref(false)
const currentShard = ref(null)

function formatTime(t) {
  if (!t) return '-'
  return new Date(t).toLocaleString('zh-CN')
}

function loadStatusBadge(s) {
  return { loaded: 'success', loading: 'processing', idle: 'default', error: 'danger' }[s] || 'default'
}
function loadStatusText(s) {
  return { loaded: '已加载', loading: '加载中', idle: '空闲', error: '异常' }[s] || s
}
function versionStatusColor(s) {
  return { published: 'green', draft: 'gray', deprecated: 'red' }[s] || 'gray'
}
function versionStatusText(s) {
  return { published: '已发布', draft: '草稿', deprecated: '已废弃' }[s] || s
}

async function loadShards() {
  loading.value = true
  try {
    const params = { page: page.value, page_size: pageSize.value }
    if (searchKey.value) params.keyword = searchKey.value
    if (filterVersion.value) params.version = filterVersion.value
    const qs = new URLSearchParams(params).toString()
    const res = await apiRequest(`${API_BASE}/model-shards?${qs}`)
    shards.value = res.data?.list || []
    total.value = res.data?.total || 0
  } catch (e) {
    Message.error('加载分片列表失败')
  } finally {
    loading.value = false
  }
}

async function loadStats() {
  statsLoading.value = true
  try {
    const res = await apiRequest(`${API_BASE}/model-shards/stats`)
    stats.value = res.data || {}
  } catch (e) { /* ignore */ }
  finally { statsLoading.value = false }
}

async function loadVersions() {
  versionsLoading.value = true
  try {
    const res = await apiRequest(`${API_BASE}/model-shards/versions?page_size=50`)
    versions.value = res.data?.list || []
  } catch (e) {
    Message.error('加载版本列表失败')
  } finally {
    versionsLoading.value = false
  }
}

async function loadDevices() {
  try {
    const res = await apiRequest(`${API_BASE}/devices?page_size=100`)
    devices.value = res.data?.list || []
  } catch (e) { /* ignore */ }
}

async function handleCreate(done) {
  if (!createForm.name || !createForm.version) {
    Message.warning('请填写必填项')
    done(false)
    return
  }
  try {
    await apiRequest(`${API_BASE}/model-shards`, { method: 'POST', body: JSON.stringify(createForm) })
    Message.success('创建成功')
    showCreateModal.value = false
    loadShards()
    loadStats()
    done(true)
  } catch (e) {
    Message.error('创建失败: ' + e.message)
    done(false)
  }
}

async function handleCreateVersion(done) {
  if (!createVersionForm.version) {
    Message.warning('请填写版本号')
    done(false)
    return
  }
  try {
    await apiRequest(`${API_BASE}/model-shards/versions`, { method: 'POST', body: JSON.stringify(createVersionForm) })
    Message.success('创建成功')
    showCreateVersionModal.value = false
    createVersionForm.version = ''
    createVersionForm.description = ''
    loadVersions()
    done(true)
  } catch (e) {
    Message.error('创建失败: ' + e.message)
    done(false)
  }
}

async function publishVersion(record) {
  try {
    await apiRequest(`${API_BASE}/model-shards/versions/${record.id}/publish`, { method: 'POST' })
    Message.success('发布成功')
    loadVersions()
  } catch (e) {
    Message.error('发布失败: ' + e.message)
  }
}

function showLoadModal(record) {
  loadForm.shard_id = record.id
  loadForm.shard_name = record.name
  loadForm.device_id = ''
  showLoadModalFlag.value = true
}

async function handleLoad(done) {
  if (!loadForm.device_id) {
    Message.warning('请选择设备')
    done(false)
    return
  }
  try {
    await apiRequest(`${API_BASE}/model-shards/${loadForm.shard_id}/load`, {
      method: 'POST',
      body: JSON.stringify({ device_id: loadForm.device_id })
    })
    Message.success('加载成功')
    showLoadModalFlag.value = false
    loadShards()
    loadStats()
    done(true)
  } catch (e) {
    Message.error('加载失败: ' + e.message)
    done(false)
  }
}

async function unloadShard(record) {
  if (!record.loaded_device_id) {
    Message.warning('该分片未加载到设备')
    return
  }
  try {
    await apiRequest(`${API_BASE}/model-shards/${record.id}/unload`, {
      method: 'POST',
      body: JSON.stringify({ device_id: record.loaded_device_id })
    })
    Message.success('卸载成功')
    loadShards()
    loadStats()
  } catch (e) {
    Message.error('卸载失败: ' + e.message)
  }
}

function editShard(record) {
  Object.assign(createForm, { id: record.id, name: record.name, version: record.version, layer_start: record.layer_start, layer_end: record.layer_end, description: record.description })
  showCreateModal.value = true
}

async function removeShard(record) {
  try {
    await apiRequest(`${API_BASE}/model-shards/${record.id}`, { method: 'DELETE' })
    Message.success('删除成功')
    loadShards()
    loadStats()
  } catch (e) {
    Message.error('删除失败: ' + e.message)
  }
}

function viewDetail(record) {
  currentShard.value = record
  showDetailDrawer.value = true
}

function viewVersionDetail(record) {
  // 简化为弹窗显示详情
  Message.info(`版本 ${record.version}: ${record.shard_count} 个分片, ${record.total_size_mb} MB`)
}

function onPageChange(p) {
  page.value = p
  loadShards()
}

onMounted(() => {
  loadShards()
  loadStats()
  loadVersions()
  loadDevices()
})
</script>
