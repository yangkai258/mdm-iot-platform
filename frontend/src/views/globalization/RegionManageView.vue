<template>
  <div class="page-container">
    <div class="search-form">
      <a-form :model="form" layout="inline">
        <a-form-item label="区域类型">
          <a-select v-model="filters.regionType" placeholder="请选择" allow-clear style="width: 140px">
            <a-option v-for="t in regionTypes" :key="t.value" :value="t.value">{{ t.label }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="状态">
          <a-select v-model="filters.status" placeholder="请选择" allow-clear style="width: 120px">
            <a-option value="active">活跃</a-option>
            <a-option value="inactive">停用</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="区域名称">
          <a-input v-model="filters.keyword" placeholder="搜索区域名称" style="width: 160px" allow-clear />
        </a-form-item>
        <a-form-item>
          <a-button type="primary" @click="handleFilterChange">搜索</a-button>
          <a-button @click="handleReset">重置</a-button>
        </a-form-item>
      </a-form>
    </div>
    <div class="toolbar">
      <a-button type="primary" @click="openAddModal">添加区域</a-button>
    </div>
    <a-table :columns="columns" :data="filteredRegions" :loading="loading" :pagination="pagination" row-key="id">
      <template #regionName="{ record }">
        <span class="region-name">{{ record.region_name }}</span>
      </template>
      <template #regionCode="{ record }">
        <span class="mono">{{ record.region_code }}</span>
      </template>
      <template #regionType="{ record }">
        <a-tag :color="typeColor(record.region_type)">{{ typeLabel(record.region_type) }}</a-tag>
      </template>
      <template #status="{ record }">
        <span class="status-dot">
          <a-badge :color="statusColor(record.is_active ? 'active' : 'inactive')" />
          {{ record.is_active ? '活跃' : '停用' }}
        </span>
      </template>
      <template #dbStatus="{ record }">
        <a-badge :color="record.db_online ? 'green' : 'gray'" :text="record.db_online ? '在线' : '离线'" />
      </template>
      <template #aiStatus="{ record }">
        <a-badge :color="record.ai_online ? 'green' : 'gray'" :text="record.ai_online ? '在线' : '离线'" />
      </template>
      <template #actions="{ record }">
        <a-button type="text" size="small" @click="openDetailModal(record)">详情</a-button>
        <a-button type="text" size="small" @click="openEditModal(record)">编辑</a-button>
        <a-button type="text" size="small" status="danger" @click="handleDelete(record)">删除</a-button>
      </template>
    </a-table>
    <a-modal v-model:visible="formVisible" :title="editingRegion ? '编辑区域' : '添加区域'" :width="520" @before-ok="handleSubmit">
      <a-form :model="form" label-col-flex="100px">
        <a-form-item label="区域名称" required>
          <a-input v-model="form.region_name" placeholder="如：中国东部" />
        </a-form-item>
        <a-form-item label="区域代码" required>
          <a-input v-model="form.region_code" placeholder="如：cn-east" />
        </a-form-item>
        <a-form-item label="区域类型" required>
          <a-select v-model="form.region_type" placeholder="选择区域类型">
            <a-option v-for="t in regionTypes" :key="t.value" :value="t.value">{{ t.label }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="数据库 Schema">
          <a-input v-model="form.db_schema" placeholder="如：public_cn_east" />
        </a-form-item>
        <a-form-item label="AI 端点">
          <a-input v-model="form.ai_endpoint" placeholder="https://ai.cn-east.example.com" />
        </a-form-item>
        <a-form-item label="设为默认">
          <a-switch v-model="form.is_default" />
        </a-form-item>
        <a-form-item label="启用状态">
          <a-switch v-model="form.is_active" />
        </a-form-item>
      </a-form>
      <template #footer>
        <a-button @click="formVisible = false">取消</a-button>
        <a-button type="primary" @click="handleSubmit">确定</a-button>
      </template>
    </a-modal>
    <a-modal v-model:visible="detailVisible" title="区域详情" :width="560">
      <div class="detail-grid" v-if="detailRegion">
        <div class="detail-item">
          <span class="detail-label">区域名称</span>
          <span class="detail-value">{{ detailRegion.region_name }}</span>
        </div>
        <div class="detail-item">
          <span class="detail-label">区域代码</span>
          <span class="detail-value mono">{{ detailRegion.region_code }}</span>
        </div>
        <div class="detail-item">
          <span class="detail-label">类型</span>
          <span class="detail-value">
            <a-tag :color="typeColor(detailRegion.region_type)">{{ typeLabel(detailRegion.region_type) }}</a-tag>
          </span>
        </div>
        <div class="detail-item">
          <span class="detail-label">状态</span>
          <span class="detail-value">
            <a-badge :color="statusColor(detailRegion.is_active ? 'active' : 'inactive')" :text="detailRegion.is_active ? '活跃' : '停用'" />
          </span>
        </div>
        <div class="detail-item">
          <span class="detail-label">数据库</span>
          <span class="detail-value">{{ detailRegion.db_schema || '-' }}</span>
        </div>
        <div class="detail-item">
          <span class="detail-label">AI 端点</span>
          <span class="detail-value">{{ detailRegion.ai_endpoint || '-' }}</span>
        </div>
        <div class="detail-item">
          <span class="detail-label">创建时间</span>
          <span class="detail-value">{{ formatDate(detailRegion.created_at) }}</span>
        </div>
        <div class="detail-item">
          <span class="detail-label">更新时间</span>
          <span class="detail-value">{{ formatDate(detailRegion.updated_at) }}</span>
        </div>
      </div>
      <template #footer></template>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { Message, Modal } from '@arco-design/web-vue'
import { getRegions, createRegion, updateRegion, deleteRegion } from '@/api/globalization'
import { REGION_TYPES } from '@/composables/useGlobalization'
import dayjs from 'dayjs'

const regionTypes = REGION_TYPES
const loading = ref(false)
const regions = ref([])
const filters = reactive({ regionType: '', status: '', keyword: '' })
const pagination = reactive({ pageSize: 20, current: 1, total: 0 })

const columns = [
  { title: '区域名称', slotName: 'regionName', width: 160 },
  { title: '区域代码', slotName: 'regionCode', width: 120 },
  { title: '类型', slotName: 'regionType', width: 110 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '数据库', slotName: 'dbStatus', width: 100 },
  { title: 'AI节点', slotName: 'aiStatus', width: 100 },
  { title: '操作', slotName: 'actions', width: 160 }
]

const filteredRegions = computed(() => {
  return regions.value.filter(r => {
    if (filters.regionType && r.region_type !== filters.regionType) return false
    if (filters.status) {
      const isActive = filters.status === 'active'
      if (r.is_active !== isActive) return false
    }
    if (filters.keyword && !r.region_name.toLowerCase().includes(filters.keyword.toLowerCase())) return false
    return true
  })
})

const formVisible = ref(false)
const editingRegion = ref(null)
const detailVisible = ref(false)
const detailRegion = ref(null)
const form = reactive({
  region_name: '',
  region_code: '',
  region_type: 'primary',
  db_schema: '',
  ai_endpoint: '',
  is_default: false,
  is_active: true
})

function typeColor(type) {
  const map = { primary: 'blue', backup: 'orange', ai_node: 'purple', storage: 'cyan' }
  return map[type] || 'default'
}

function typeLabel(type) {
  const map = { primary: '主区域', backup: '备份区域', ai_node: '推理节点', storage: '存储节点' }
  return map[type] || type
}

function statusColor(status) {
  return status === 'active' ? 'green' : 'gray'
}

function formatDate(date) {
  if (!date) return '-'
  return dayjs(date).format('YYYY-MM-DD HH:mm:ss')
}

function handleFilterChange() {
  pagination.current = 1
}

function handleReset() {
  filters.regionType = ''
  filters.status = ''
  filters.keyword = ''
  pagination.current = 1
}

onMounted(async () => {
  loading.value = true
  try {
    const res = await getRegions()
    regions.value = res.data || res || []
    pagination.total = regions.value.length
  } catch (e) {
    console.error('加载区域列表失败', e)
  } finally {
    loading.value = false
  }
})

function openAddModal() {
  editingRegion.value = null
  Object.assign(form, { region_name: '', region_code: '', region_type: 'primary', db_schema: '', ai_endpoint: '', is_default: false, is_active: true })
  formVisible.value = true
}

function openEditModal(record) {
  editingRegion.value = record
  Object.assign(form, {
    region_name: record.region_name,
    region_code: record.region_code,
    region_type: record.region_type,
    db_schema: record.db_schema || '',
    ai_endpoint: record.ai_endpoint || '',
    is_default: record.is_default || false,
    is_active: record.is_active
  })
  formVisible.value = true
}

function openDetailModal(record) {
  detailRegion.value = record
  detailVisible.value = true
}

async function handleSubmit() {
  try {
    if (editingRegion.value) {
      await updateRegion(editingRegion.value.id, form)
      Message.success('更新成功')
    } else {
      await createRegion(form)
      Message.success('创建成功')
    }
    formVisible.value = false
    const res = await getRegions()
    regions.value = res.data || res || []
  } catch (e) {
    Message.error('操作失败')
  }
}

function handleDelete(record) {
  Modal.warning({
    title: '确认删除',
    content: `确定要删除区域「${record.region_name}」吗？此操作不可恢复。`,
    okText: '删除',
    onOk: async () => {
      try {
        await deleteRegion(record.id)
        Message.success('删除成功')
        const res = await getRegions()
        regions.value = res.data || res || []
      } catch (e) {
        Message.error('删除失败')
      }
    }
  })
}
</script>

<style scoped>
.page-container { background: #fff; border-radius: 4px; padding: 20px; }
.search-form { margin-bottom: 16px; padding: 16px; background: #f7f8fa; border-radius: 4px; }
.toolbar { margin-bottom: 16px; }
.region-name { font-weight: 500; }
.mono { font-family: monospace; font-size: 12px; }
.status-dot { display: flex; align-items: center; gap: 6px; }
.detail-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 16px; }
.detail-item { display: flex; flex-direction: column; gap: 4px; }
.detail-label { font-size: 12px; color: var(--color-text-3); }
.detail-value { font-size: 14px; word-break: break-all; }
</style>
