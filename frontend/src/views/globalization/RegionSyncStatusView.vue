<template>
  <div class="page-container">
    <!-- 同步状态概览 -->
    <a-card class="overview-card">
      <div class="overview-grid">
        <div class="overview-item" v-for="item in overviewStats" :key="item.label">
          <div class="overview-icon">
            <component :is="item.icon" />
          </div>
          <div class="overview-info">
            <div class="overview-value">{{ item.value }}</div>
            <div class="overview-label">{{ item.label }}</div>
          </div>
        </div>
      </div>
    </a-card>

    <!-- 同步记录列表 -->
    <a-card class="table-card">
      <template #title>
        <div class="section-title">同步记录</div>
      </template>
      <div class="filter-row">
        <a-select
          v-model="filters.syncStatus"
          placeholder="同步状态"
          style="width: 140px"
          allow-clear
          @change="loadSyncRecords"
        >
          <a-option value="completed">已完成</a-option>
          <a-option value="syncing">同步中</a-option>
          <a-option value="pending">等待中</a-option>
          <a-option value="failed">失败</a-option>
        </a-select>
        <a-button type="primary" @click="openTriggerModal">手动触发同步</a-button>
      </div>

      <a-table
        :columns="columns"
        :data="filteredRecords"
        :loading="loading"
        :pagination="{ pageSize: 20 }"
        row-key="id"
        style="margin-top: 12px"
      >
        <template #syncPath="{ record }">
          <span class="sync-path">
            <span class="region-tag">{{ record.source_region }}</span>
            <icon-swap class="swap-icon" />
            <span class="region-tag">{{ record.target_region }}</span>
          </span>
        </template>
        <template #syncType="{ record }">
          <a-tag>{{ syncTypeLabel(record.sync_type) }}</a-tag>
        </template>
        <template #status="{ record }">
          <a-badge :color="syncStatusColor(record.sync_status)" :text="syncStatusLabel(record.sync_status)" />
        </template>
        <template #recordsCount="{ record }">
          <span>{{ record.records_synced || 0 }}</span>
        </template>
        <template #duration="{ record }">
          <span v-if="record.started_at && record.completed_at">
            {{ calcDuration(record.started_at, record.completed_at) }}
          </span>
          <span v-else-if="record.started_at">--</span>
          <span v-else>-</span>
        </template>
        <template #error="{ record }">
          <span v-if="record.error_message" class="error-text">{{ record.error_message }}</span>
          <span v-else>-</span>
        </template>
        <template #actions="{ record }">
          <div class="action-btns">
            <a-button
              type="text"
              size="small"
              :disabled="record.sync_status === 'syncing'"
              @click="handleRetry(record)"
            >
              重试
            </a-button>
          </div>
        </template>
      </a-table>
    </a-card>

    <!-- 手动触发弹窗 -->
    <a-modal
      v-model:visible="triggerVisible"
      title="手动触发同步"
      :width="480"
      @before-ok="handleTriggerSync"
      @cancel="triggerVisible = false"
    >
      <a-form :model="form" layout="vertical">
        <a-form-item label="源区域" required>
          <a-select v-model="form.source_region" placeholder="选择源区域">
            <a-option v-for="r in regions" :key="r.region_code" :value="r.region_code">{{ r.region_name }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="目标区域" required>
          <a-select v-model="form.target_region" placeholder="选择目标区域">
            <a-option v-for="r in regions" :key="r.region_code" :value="r.region_code">{{ r.region_name }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="同步类型" required>
          <a-select v-model="form.sync_type" placeholder="选择同步类型">
            <a-option value="full">全量同步</a-option>
            <a-option value="incremental">增量同步</a-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import { getRegions, getSyncRecords, triggerSync } from '@/api/globalization'
import dayjs from 'dayjs'

const loading = ref(false)
const records = ref([])
const regions = ref([])
const filters = reactive({ syncStatus: '' })
const triggerVisible = ref(false)

const form = reactive({
  source_region: '',
  target_region: '',
  sync_type: 'incremental'
})

const columns = [
  { title: '同步路径', slotName: 'syncPath', width: 200 },
  { title: '同步类型', slotName: 'syncType', width: 100 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '同步条数', slotName: 'recordsCount', width: 100 },
  { title: '耗时', slotName: 'duration', width: 100 },
  { title: '错误信息', slotName: 'error', ellipsis: true },
  { title: '操作', slotName: 'actions', width: 80 }
]

const filteredRecords = computed(() => {
  if (!filters.syncStatus) return records.value
  return records.value.filter(r => r.sync_status === filters.syncStatus)
})

const overviewStats = computed(() => [
  { label: '总同步数', value: records.value.length, icon: 'IconSync' },
  { label: '已完成', value: records.value.filter(r => r.sync_status === 'completed').length, icon: 'IconCheckCircle' },
  { label: '进行中', value: records.value.filter(r => r.sync_status === 'syncing').length, icon: 'IconRefresh' },
  { label: '失败', value: records.value.filter(r => r.sync_status === 'failed').length, icon: 'IconCloseCircle' }
])

function syncTypeLabel(type) {
  const map = { full: '全量同步', incremental: '增量同步' }
  return map[type] || type
}

function syncStatusColor(status) {
  const map = { completed: 'green', syncing: 'blue', pending: 'yellow', failed: 'red' }
  return map[status] || 'default'
}

function syncStatusLabel(status) {
  const map = { completed: '已完成', syncing: '同步中', pending: '等待中', failed: '失败' }
  return map[status] || status
}

function calcDuration(start, end) {
  const ms = dayjs(end).diff(dayjs(start))
  if (ms < 1000) return `${ms}ms`
  if (ms < 60000) return `${(ms / 1000).toFixed(1)}s`
  return `${(ms / 60000).toFixed(1)}min`
}

async function loadSyncRecords() {
  loading.value = true
  try {
    const res = await getSyncRecords()
    records.value = res.data || res || []
  } catch (e) {
    console.error('加载同步记录失败', e)
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  try {
    const [recordsRes, regionsRes] = await Promise.all([
      getSyncRecords(),
      getRegions()
    ])
    records.value = recordsRes.data || recordsRes || []
    regions.value = regionsRes.data || regionsRes || []
  } catch (e) {
    console.error('加载数据失败', e)
  }
})

function openTriggerModal() {
  Object.assign(form, { source_region: '', target_region: '', sync_type: 'incremental' })
  triggerVisible.value = true
}

async function handleTriggerSync(done) {
  if (!form.source_region || !form.target_region) {
    Message.warning('请选择源区域和目标区域')
    done(false)
    return
  }
  if (form.source_region === form.target_region) {
    Message.warning('源区域和目标区域不能相同')
    done(false)
    return
  }
  try {
    await triggerSync(form)
    Message.success('同步任务已触发')
    triggerVisible.value = false
    await loadSyncRecords()
  } catch (e) {
    Message.error('触发失败')
    done(false)
  }
}

async function handleRetry(record) {
  try {
    await triggerSync({
      source_region: record.source_region,
      target_region: record.target_region,
      sync_type: record.sync_type
    })
    Message.success('重试任务已触发')
    await loadSyncRecords()
  } catch (e) {
    Message.error('重试失败')
  }
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

.overview-card { flex-shrink: 0; }

.overview-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
}

.overview-item {
  display: flex;
  align-items: center;
  gap: 12px;
}

.overview-icon {
  width: 40px;
  height: 40px;
  border-radius: 8px;
  background: var(--color-fill-2);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  color: var(--color-text-3);
}

.overview-info { display: flex; flex-direction: column; }

.overview-value { font-size: 24px; font-weight: 600; }

.overview-label { font-size: 12px; color: var(--color-text-3); }

.table-card { flex: 1; overflow: auto; }

.section-title { font-size: 14px; font-weight: 600; }

.filter-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.sync-path {
  display: flex;
  align-items: center;
  gap: 6px;
}

.region-tag {
  background: var(--color-fill-2);
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-family: monospace;
}

.swap-icon { color: var(--color-text-3); }

.error-text { color: #f53f3f; font-size: 12px; }

.action-btns { display: flex; gap: 4px; }
</style>
