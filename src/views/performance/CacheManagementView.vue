<template>
  <div class="page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>性能管理</a-breadcrumb-item>
      <a-breadcrumb-item>缓存管理</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 缓存概览 -->
    <a-row :gutter="16">
      <a-col :xs="24" :sm="12" :md="6">
        <a-card class="stat-card">
          <a-statistic title="总 Key 数" :value="stats.keys">
            <template #prefix>🔑</template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :xs="24" :sm="12" :md="6">
        <a-card class="stat-card">
          <a-statistic title="命中率" :value="stats.hit_rate" suffix="%" :precision="1"
            :value-style="{ color: getHitRateColor(stats.hit_rate) }">
            <template #prefix>🎯</template>
          </a-statistic>
          <a-progress
            :percent="stats.hit_rate"
            :stroke-color="getHitRateColor(stats.hit_rate)"
            style="margin-top: 8px;"
          />
        </a-card>
      </a-col>
      <a-col :xs="24" :sm="12" :md="6">
        <a-card class="stat-card">
          <a-statistic title="命中次数" :value="stats.hits" />
          <a-statistic title="未命中次数" :value="stats.misses" :value-style="{ fontSize: '14px', color: '#86909c' }" />
        </a-card>
      </a-col>
      <a-col :xs="24" :sm="12" :md="6">
        <a-card class="stat-card">
          <a-statistic title="内存使用">
            <template #formatter>
              <span>{{ formatBytes(stats.memory_used) }}</span>
            </template>
          </a-statistic>
          <a-progress
            :percent="(stats.memory_used / stats.memory_total * 100)"
            :stroke-color="getMemColor(stats.memory_used / stats.memory_total * 100)"
            style="margin-top: 8px;"
          />
          <div class="mem-bar-label">{{ formatBytes(stats.memory_used) }} / {{ formatBytes(stats.memory_total) }}</div>
        </a-card>
      </a-col>
    </a-row>

    <!-- 操作区 -->
    <a-card title="🔧 缓存控制" class="section-card" style="margin-top: 16px;">
      <a-row :gutter="16" align="middle">
        <a-col :xs="24" :md="16">
          <a-space wrap>
            <a-button type="primary" status="danger" @click="handleClearCache('all')" :loading="clearing">
              清空所有缓存
            </a-button>
            <a-button status="warning" @click="handleClearCache('expired')" :loading="clearing">
              清除过期 Key
            </a-button>
            <a-button status="normal" @click="handleClearCache('memory')" :loading="clearing">
              释放内存
            </a-button>
          </a-space>
          <div class="clear-tip">
            <icon-info-circle :size="14" />
            <span>清空操作不可恢复，请谨慎操作</span>
          </div>
        </a-col>
        <a-col :xs="24" :md="8" style="text-align: right;">
          <a-button type="primary" @click="refresh">
            <template #icon><icon-refresh :spin="refreshing" /></template>
            刷新数据
          </a-button>
        </a-col>
      </a-row>
    </a-card>

    <!-- 命令统计 -->
    <a-card title="📈 命令统计" class="section-card" style="margin-top: 16px;">
      <a-row :gutter="16">
        <a-col :xs="12" :sm="6">
          <div class="cmd-stat">
            <div class="cmd-icon">📖</div>
            <div class="cmd-info">
              <div class="cmd-label">GET</div>
              <div class="cmd-value">{{ stats.cmd_get.toLocaleString() }}</div>
            </div>
          </div>
        </a-col>
        <a-col :xs="12" :sm="6">
          <div class="cmd-stat">
            <div class="cmd-icon">✏️</div>
            <div class="cmd-info">
              <div class="cmd-label">SET</div>
              <div class="cmd-value">{{ stats.cmd_set.toLocaleString() }}</div>
            </div>
          </div>
        </a-col>
        <a-col :xs="12" :sm="6">
          <div class="cmd-stat">
            <div class="cmd-icon">🗑️</div>
            <div class="cmd-info">
              <div class="cmd-label">DEL</div>
              <div class="cmd-value">{{ stats.cmd_del.toLocaleString() }}</div>
            </div>
          </div>
        </a-col>
        <a-col :xs="12" :sm="6">
          <div class="cmd-stat">
            <div class="cmd-icon">⚡</div>
            <div class="cmd-info">
              <div class="cmd-label">驱逐</div>
              <div class="cmd-value">{{ stats.evictions.toLocaleString() }}</div>
            </div>
          </div>
        </a-col>
      </a-row>
    </a-card>

    <!-- 缓存 Key 列表 -->
    <a-card title="🔑 缓存 Key 列表" class="section-card" style="margin-top: 16px;">
      <template #extra>
        <a-space>
          <a-input-search
            v-model="searchKey"
            placeholder="搜索 Key..."
            style="width: 200px;"
            @search="loadCacheKeys"
            search-button
          />
          <a-button @click="loadCacheKeys">
            <template #icon><icon-refresh /></template>
          </a-button>
        </a-space>
      </template>

      <a-table
        :columns="columns"
        :data="cacheKeys"
        :loading="loadingKeys"
        :pagination="{ current: page, pageSize: pageSize, total: totalKeys, showTotal: true, showPageSize: true }"
        @page-change="handlePageChange"
        @page-size-change="handlePageSizeChange"
        row-key="key"
        stripe
      >
        <template #empty>
          <a-empty description="暂无缓存数据" />
        </template>

        <template #type="{ record }">
          <a-tag :color="getKeyTypeColor(record.type)">{{ record.type }}</a-tag>
        </template>

        <template #ttl="{ record }">
          <a-tag :color="record.ttl < 0 ? 'gray' : record.ttl < 300 ? 'red' : 'arcoblue'">
            {{ record.ttl < 0 ? '永不过期' : record.ttl + 's' }}
          </a-tag>
        </template>

        <template #size="{ record }">
          {{ formatBytes(record.size) }}
        </template>

        <template #actions="{ record }">
          <a-space>
            <a-button size="small" @click="handleViewValue(record)">查看</a-button>
            <a-popconfirm content="确定删除此 Key？" @ok="handleDeleteKey(record.key)">
              <a-button size="small" status="danger">删除</a-button>
            </a-popconfirm>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- Key 值查看弹窗 -->
    <a-modal
      v-model:visible="valueModalVisible"
      title="Key 详情"
      :width="600"
      :footer="null"
    >
      <a-descriptions :column="1" bordered size="small">
        <a-descriptions-item label="Key">{{ currentKey }}</a-descriptions-item>
        <a-descriptions-item label="类型">{{ currentKeyInfo.type }}</a-descriptions-item>
        <a-descriptions-item label="TTL">
          <a-tag :color="currentKeyInfo.ttl < 0 ? 'gray' : 'arcoblue'">
            {{ currentKeyInfo.ttl < 0 ? '永不过期' : currentKeyInfo.ttl + 's' }}
          </a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="大小">{{ formatBytes(currentKeyInfo.size) }}</a-descriptions-item>
        <a-descriptions-item label="值">
          <div class="value-preview">{{ currentKeyValue }}</div>
        </a-descriptions-item>
      </a-descriptions>
    </a-modal>

    <!-- 清空确认弹窗 -->
    <a-modal
      v-model:visible="clearModalVisible"
      :title="clearAction === 'all' ? '确认清空所有缓存' : clearAction === 'expired' ? '确认清除过期 Key' : '确认释放内存'"
      @before-ok="confirmClear"
      @cancel="clearModalVisible = false"
    >
      <div class="clear-confirm-content">
        <a-result
          v-if="clearAction === 'all'"
          status="warning"
          title="清空所有缓存"
          subtitle="此操作将删除所有缓存数据，包括未过期的 Key，操作不可恢复！"
        />
        <a-result
          v-else-if="clearAction === 'expired'"
          status="warning"
          title="清除过期 Key"
          subtitle="此操作将删除所有已过期的缓存 Key，请确认。"
        />
        <a-result
          v-else
          status="warning"
          title="释放内存"
          subtitle="此操作将触发 Redis 内存整理，释放未使用的内存碎片。"
        />
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import performanceApi from '@/api/performance'
import { Message, Modal } from '@arco-design/web-vue'

// 缓存统计
const stats = ref({
  keys: 0,
  hits: 0,
  misses: 0,
  hit_rate: 0,
  memory_used: 0,
  memory_total: 0,
  cmd_get: 0,
  cmd_set: 0,
  cmd_del: 0,
  evictions: 0,
  expired: 0
})

// 缓存 Key 列表
const cacheKeys = ref([])
const loadingKeys = ref(false)
const searchKey = ref('')
const page = ref(1)
const pageSize = ref(10)
const totalKeys = ref(0)

const columns = [
  { title: 'Key', dataIndex: 'key', ellipsis: true },
  { title: '类型', dataIndex: 'type', slotName: 'type', width: 80 },
  { title: 'TTL', dataIndex: 'ttl', slotName: 'ttl', width: 120 },
  { title: '大小', dataIndex: 'size', slotName: 'size', width: 100 },
  { title: '操作', slotName: 'actions', width: 160, fixed: 'right' }
]

// Key 详情
const valueModalVisible = ref(false)
const currentKey = ref('')
const currentKeyInfo = ref({ type: '', ttl: 0, size: 0 })
const currentKeyValue = ref('')

// 清空缓存
const clearing = ref(false)
const refreshing = ref(false)
const clearModalVisible = ref(false)
const clearAction = ref('all')

// 格式化
const formatBytes = (bytes) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i]
}

// 颜色
const getHitRateColor = (v) => {
  if (v < 50) return '#f53f3f'
  if (v < 80) return '#ff7d00'
  return '#00b42a'
}

const getMemColor = (v) => {
  if (v > 80) return '#f53f3f'
  if (v > 60) return '#ff7d00'
  return '#165dff'
}

const getKeyTypeColor = (type) => {
  const map = { string: 'blue', list: 'green', hash: 'orange', set: 'purple', zset: 'red' }
  return map[type] || 'gray'
}

// 加载缓存统计
const loadCacheStats = async () => {
  try {
    const res = await performanceApi.getCacheStats()
    if (res.code === 0) {
      stats.value = res.data
    }
  } catch (e) {
    console.error('加载缓存统计失败:', e)
  }
}

// 加载缓存 Key 列表
const loadCacheKeys = async () => {
  loadingKeys.value = true
  try {
    const token = localStorage.getItem('token')
    const params = {
      page: page.value,
      page_size: pageSize.value,
      keyword: searchKey.value
    }
    const qs = '?' + new URLSearchParams(params).toString()
    const res = await fetch(`/api/v1/performance/cache/keys${qs}`, {
      headers: { 'Authorization': `Bearer ${token}` }
    }).then(r => r.json())
    if (res.code === 0) {
      cacheKeys.value = res.data.list || []
      totalKeys.value = res.data.total || 0
    }
  } catch (e) {
    console.error('加载缓存 Key 列表失败:', e)
  } finally {
    loadingKeys.value = false
  }
}

// 查看 Key 值
const handleViewValue = async (record) => {
  currentKey.value = record.key
  currentKeyInfo.value = { type: record.type, ttl: record.ttl, size: record.size }
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`/api/v1/performance/cache/keys/${encodeURIComponent(record.key)}`, {
      headers: { 'Authorization': `Bearer ${token}` }
    }).then(r => r.json())
    if (res.code === 0) {
      currentKeyValue.value = typeof res.data === 'object' ? JSON.stringify(res.data, null, 2) : String(res.data)
    } else {
      currentKeyValue.value = '获取失败: ' + (res.message || '未知错误')
    }
  } catch (e) {
    currentKeyValue.value = '获取失败: ' + e.message
  }
  valueModalVisible.value = true
}

// 删除 Key
const handleDeleteKey = async (key) => {
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`/api/v1/performance/cache/keys/${encodeURIComponent(key)}`, {
      method: 'DELETE',
      headers: { 'Authorization': `Bearer ${token}` }
    }).then(r => r.json())
    if (res.code === 0) {
      Message.success('删除成功')
      loadCacheKeys()
      loadCacheStats()
    } else {
      Message.error('删除失败: ' + (res.message || '未知错误'))
    }
  } catch (e) {
    Message.error('删除失败: ' + e.message)
  }
}

// 清空缓存
const handleClearCache = (type) => {
  clearAction.value = type
  clearModalVisible.value = true
}

const confirmClear = async (done) => {
  clearing.value = true
  try {
    const res = await performanceApi.clearCache(clearAction.value)
    if (res.code === 0) {
      Message.success('操作成功')
      clearModalVisible.value = false
      loadCacheStats()
      loadCacheKeys()
    } else {
      Message.error('操作失败: ' + (res.message || '未知错误'))
    }
  } catch (e) {
    Message.error('操作失败: ' + e.message)
  } finally {
    clearing.value = false
    done()
  }
}

// 分页
const handlePageChange = (p) => {
  page.value = p
  loadCacheKeys()
}

const handlePageSizeChange = (size) => {
  pageSize.value = size
  page.value = 1
  loadCacheKeys()
}

// 刷新
const refresh = async () => {
  refreshing.value = true
  await Promise.all([loadCacheStats(), loadCacheKeys()])
  refreshing.value = false
  Message.success('刷新成功')
}

onMounted(() => {
  loadCacheStats()
  loadCacheKeys()
})
</script>

<style scoped>
.page-container {
  padding: 20px 24px;
  min-height: calc(100vh - 64px);
  background: #f5f7fa;
}

.breadcrumb {
  margin-bottom: 16px;
}

.section-card {
  border-radius: 8px;
}

.stat-card {
  border-radius: 8px;
  text-align: center;
}

.mem-bar-label {
  font-size: 11px;
  color: #86909c;
  margin-top: 4px;
  text-align: center;
}

.clear-tip {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-top: 8px;
  font-size: 12px;
  color: #86909c;
}

.cmd-stat {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: #f7f8fa;
  border-radius: 8px;
}

.cmd-icon {
  font-size: 28px;
}

.cmd-label {
  font-size: 12px;
  color: #86909c;
}

.cmd-value {
  font-size: 18px;
  font-weight: 600;
  color: #1d2129;
}

.value-preview {
  max-height: 300px;
  overflow: auto;
  background: #f7f8fa;
  padding: 12px;
  border-radius: 4px;
  font-family: 'Courier New', monospace;
  font-size: 12px;
  white-space: pre-wrap;
  word-break: break-all;
}

.clear-confirm-content {
  padding: 16px 0;
}
</style>
