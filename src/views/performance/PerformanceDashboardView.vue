<template>
  <div class="page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>性能管理</a-breadcrumb-item>
      <a-breadcrumb-item>性能仪表盘</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 系统性能指标 -->
    <a-card title="📊 系统性能" class="section-card">
      <template #extra>
        <a-tag :color="systemStatusColor">{{ systemStatusText }}</a-tag>
      </template>
      <a-row :gutter="16">
        <!-- CPU -->
        <a-col :xs="24" :sm="12" :md="6">
          <div class="metric-item">
            <div class="metric-header">
              <span class="metric-icon">💻</span>
              <span class="metric-label">CPU 使用率</span>
            </div>
            <a-progress
              type="circle"
              :percent="Number(systemMetrics.cpu_usage.toFixed(1))"
              :width="100"
              :stroke-color="getCpuColor(systemMetrics.cpu_usage)"
            >
              <template #formatter>
                <span class="progress-text">{{ systemMetrics.cpu_usage.toFixed(1) }}%</span>
              </template>
            </a-progress>
          </div>
        </a-col>

        <!-- 内存 -->
        <a-col :xs="24" :sm="12" :md="6">
          <div class="metric-item">
            <div class="metric-header">
              <span class="metric-icon">🧠</span>
              <span class="metric-label">内存使用</span>
            </div>
            <a-progress
              type="circle"
              :percent="Number(systemMetrics.memory_usage.toFixed(1))"
              :width="100"
              :stroke-color="getMemColor(systemMetrics.memory_usage)"
            >
              <template #formatter>
                <span class="progress-text">{{ systemMetrics.memory_usage.toFixed(1) }}%</span>
              </template>
            </a-progress>
            <div class="metric-detail">
              {{ systemMetrics.memory_used }} / {{ systemMetrics.memory_total }} MB
            </div>
          </div>
        </a-col>

        <!-- 磁盘 -->
        <a-col :xs="24" :sm="12" :md="6">
          <div class="metric-item">
            <div class="metric-header">
              <span class="metric-icon">💾</span>
              <span class="metric-label">磁盘使用</span>
            </div>
            <a-progress
              type="circle"
              :percent="Number(systemMetrics.disk_usage.toFixed(1))"
              :width="100"
              :stroke-color="getDiskColor(systemMetrics.disk_usage)"
            >
              <template #formatter>
                <span class="progress-text">{{ systemMetrics.disk_usage.toFixed(1) }}%</span>
              </template>
            </a-progress>
            <div class="metric-detail">
              {{ systemMetrics.disk_used }} / {{ systemMetrics.disk_total }} GB
            </div>
          </div>
        </a-col>

        <!-- 系统负载 -->
        <a-col :xs="24" :sm="12" :md="6">
          <div class="metric-item">
            <div class="metric-header">
              <span class="metric-icon">⚖️</span>
              <span class="metric-label">系统负载</span>
            </div>
            <div class="load-avg">
              <div class="load-item">
                <span class="load-label">1min</span>
                <a-progress
                  :percent="Math.min(systemMetrics.load_avg[0] * 20, 100)"
                  :stroke-color="getLoadColor(systemMetrics.load_avg[0])"
                  :show-text="false"
                  :stroke-width="8"
                />
                <span class="load-value">{{ systemMetrics.load_avg[0].toFixed(2) }}</span>
              </div>
              <div class="load-item">
                <span class="load-label">5min</span>
                <a-progress
                  :percent="Math.min(systemMetrics.load_avg[1] * 20, 100)"
                  :stroke-color="getLoadColor(systemMetrics.load_avg[1])"
                  :show-text="false"
                  :stroke-width="8"
                />
                <span class="load-value">{{ systemMetrics.load_avg[1].toFixed(2) }}</span>
              </div>
              <div class="load-item">
                <span class="load-label">15min</span>
                <a-progress
                  :percent="Math.min(systemMetrics.load_avg[2] * 20, 100)"
                  :stroke-color="getLoadColor(systemMetrics.load_avg[2])"
                  :show-text="false"
                  :stroke-width="8"
                />
                <span class="load-value">{{ systemMetrics.load_avg[2].toFixed(2) }}</span>
              </div>
            </div>
            <div class="metric-detail">
              运行时间: {{ formatUptime(systemMetrics.uptime) }}
            </div>
          </div>
        </a-col>
      </a-row>

      <!-- 网络 & Go 协程 -->
      <a-row :gutter="16" style="margin-top: 24px;">
        <a-col :xs="24" :sm="12">
          <div class="sub-metric">
            <span class="sub-metric-icon">🌐</span>
            <span class="sub-metric-label">网络速率</span>
            <span class="sub-metric-value">
              ↓ {{ formatBytes(systemMetrics.network_rx) }}/s
              &nbsp;&nbsp;
              ↑ {{ formatBytes(systemMetrics.network_tx) }}/s
            </span>
          </div>
        </a-col>
        <a-col :xs="24" :sm="12">
          <div class="sub-metric">
            <span class="sub-metric-icon">⚙️</span>
            <span class="sub-metric-label">Go 协程</span>
            <span class="sub-metric-value">{{ systemMetrics.goroutines }}</span>
          </div>
        </a-col>
      </a-row>
    </a-card>

    <!-- 缓存统计 & 数据库状态 -->
    <a-row :gutter="16" style="margin-top: 16px;">
      <!-- 缓存统计 -->
      <a-col :xs="24" :lg="12">
        <a-card title="🚀 缓存统计" class="section-card">
          <template #extra>
            <a-tag color="arcoblue">Redis</a-tag>
          </template>
          <a-row :gutter="12">
            <a-col :span="8">
              <a-statistic title="命中率" :value="cacheStats.hit_rate" suffix="%" :precision="1"
                :value-style="{ color: getHitRateColor(cacheStats.hit_rate) }" />
              <a-progress
                :percent="cacheStats.hit_rate"
                :stroke-color="getHitRateColor(cacheStats.hit_rate)"
                style="margin-top: 8px;"
              />
            </a-col>
            <a-col :span="8">
              <a-statistic title="Key 总数" :value="cacheStats.keys" />
            </a-col>
            <a-col :span="8">
              <a-statistic title="内存使用" :value="formatBytes(cacheStats.memory_used)" />
              <div class="cache-bar">
                <a-progress
                  :percent="(cacheStats.memory_used / cacheStats.memory_total * 100)"
                  :stroke-color="getMemColor(cacheStats.memory_used / cacheStats.memory_total * 100)"
                  :show-text="false"
                  style="margin-top: 8px;"
                />
                <span class="cache-bar-label">{{ formatBytes(cacheStats.memory_used) }} / {{ formatBytes(cacheStats.memory_total) }}</span>
              </div>
            </a-col>
          </a-row>
          <a-divider />
          <a-descriptions :column="2" size="small">
            <a-descriptions-item label="GET 命令">{{ cacheStats.cmd_get.toLocaleString() }}</a-descriptions-item>
            <a-descriptions-item label="SET 命令">{{ cacheStats.cmd_set.toLocaleString() }}</a-descriptions-item>
            <a-descriptions-item label="DEL 命令">{{ cacheStats.cmd_del.toLocaleString() }}</a-descriptions-item>
            <a-descriptions-item label="驱逐次数">{{ cacheStats.evictions.toLocaleString() }}</a-descriptions-item>
            <a-descriptions-item label="命中">{{ cacheStats.hits.toLocaleString() }}</a-descriptions-item>
            <a-descriptions-item label="未命中">{{ cacheStats.misses.toLocaleString() }}</a-descriptions-item>
          </a-descriptions>
        </a-card>
      </a-col>

      <!-- 数据库状态 -->
      <a-col :xs="24" :lg="12">
        <a-card title="🗄️ 数据库状态" class="section-card">
          <template #extra>
            <a-tag :color="dbStats.status === 'connected' ? 'green' : 'red'">
              {{ dbStats.status === 'connected' ? '已连接' : '未连接' }}
            </a-tag>
          </template>
          <a-row :gutter="12">
            <a-col :span="8">
              <a-statistic title="版本" :value-style="{ fontSize: '14px' }">
                <template #formatter>
                  <span style="font-size: 13px;">{{ dbStats.version || 'N/A' }}</span>
                </template>
              </a-statistic>
            </a-col>
            <a-col :span="8">
              <a-statistic title="连接数" :value="`${dbStats.connections} / ${dbStats.connections_max}`" />
            </a-col>
            <a-col :span="8">
              <a-statistic title="平均延迟" :value="dbStats.latency_ms" suffix="ms" :precision="2" />
            </a-col>
          </a-row>
          <a-divider />
          <a-row :gutter="12">
            <a-col :span="8">
              <a-statistic title="QPS" :value="dbStats.queries_per_sec" :precision="1" />
            </a-col>
            <a-col :span="8">
              <a-statistic title="TPS" :value="dbStats.transaction_per_sec" :precision="1" />
            </a-col>
            <a-col :span="8">
              <a-statistic title="缓存命中率" :value="dbStats.cache_hit_rate" suffix="%" :precision="1" />
            </a-col>
          </a-row>
          <a-divider />
          <a-descriptions :column="2" size="small">
            <a-descriptions-item label="总查询数">{{ dbStats.queries_total.toLocaleString() }}</a-descriptions-item>
            <a-descriptions-item label="慢查询">
              <a-tag :color="dbStats.slow_queries > 0 ? 'orange' : 'green'" size="small">
                {{ dbStats.slow_queries }}
              </a-tag>
            </a-descriptions-item>
            <a-descriptions-item v-for="db in dbStats.databases" :key="db.name" :label="db.name">
              {{ db.size_mb.toFixed(1) }} MB ({{ db.tables }} 表)
            </a-descriptions-item>
          </a-descriptions>
        </a-card>
      </a-col>
    </a-row>

    <!-- 快捷操作 -->
    <a-card title="🔧 快速操作" class="section-card" style="margin-top: 16px;">
      <a-space>
        <a-button type="primary" @click="$router.push('/performance/cache')">
          <template #icon><IconDashboard /></template>
          缓存管理
        </a-button>
        <a-button @click="refreshAll">
          <template #icon><IconRefresh /></template>
          刷新数据
        </a-button>
      </a-space>
    </a-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import performanceApi from '@/api/performance'
import { Message } from '@arco-design/web-vue'

// 系统指标
const systemMetrics = ref({
  cpu_usage: 0,
  memory_usage: 0,
  memory_total: 0,
  memory_used: 0,
  disk_usage: 0,
  disk_total: 0,
  disk_used: 0,
  network_rx: 0,
  network_tx: 0,
  uptime: 0,
  load_avg: [0, 0, 0],
  goroutines: 0,
  timestamp: ''
})

// 缓存统计
const cacheStats = ref({
  hits: 0,
  misses: 0,
  hit_rate: 0,
  keys: 0,
  memory_used: 0,
  memory_total: 0,
  cmd_get: 0,
  cmd_set: 0,
  cmd_del: 0,
  evictions: 0,
  expired: 0
})

// 数据库状态
const dbStats = ref({
  status: 'disconnected',
  version: '',
  connections: 0,
  connections_max: 0,
  queries_total: 0,
  queries_per_sec: 0,
  transaction_per_sec: 0,
  cache_hit_rate: 0,
  slow_queries: 0,
  latency_ms: 0,
  databases: []
})

let timer = null

// 系统状态
const systemStatusColor = computed(() => {
  if (systemMetrics.value.cpu_usage > 80 || systemMetrics.value.memory_usage > 80) return 'red'
  if (systemMetrics.value.cpu_usage > 60 || systemMetrics.value.memory_usage > 60) return 'orange'
  return 'green'
})

const systemStatusText = computed(() => {
  if (systemMetrics.value.cpu_usage > 80 || systemMetrics.value.memory_usage > 80) return '负载过高'
  if (systemMetrics.value.cpu_usage > 60 || systemMetrics.value.memory_usage > 60) return '负载中等'
  return '运行正常'
})

// 格式化
const formatBytes = (bytes) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i]
}

const formatUptime = (seconds) => {
  const d = Math.floor(seconds / 86400)
  const h = Math.floor((seconds % 86400) / 3600)
  const m = Math.floor((seconds % 3600) / 60)
  const s = seconds % 60
  return `${d}天 ${h}时 ${m}分 ${s}秒`
}

// 颜色函数
const getCpuColor = (v) => {
  if (v > 80) return '#f53f3f'
  if (v > 60) return '#ff7d00'
  return '#00b42a'
}

const getMemColor = (v) => {
  if (v > 80) return '#f53f3f'
  if (v > 60) return '#ff7d00'
  return '#165dff'
}

const getDiskColor = (v) => {
  if (v > 90) return '#f53f3f'
  if (v > 70) return '#ff7d00'
  return '#14c9c9'
}

const getLoadColor = (v) => {
  if (v > 2) return '#f53f3f'
  if (v > 1) return '#ff7d00'
  return '#00b42a'
}

const getHitRateColor = (v) => {
  if (v < 50) return '#f53f3f'
  if (v < 80) return '#ff7d00'
  return '#00b42a'
}

// 加载数据
const loadSystemMetrics = async () => {
  try {
    const res = await performanceApi.getSystemMetrics()
    if (res.code === 0) {
      systemMetrics.value = res.data
    }
  } catch (e) {
    console.error('加载系统指标失败:', e)
  }
}

const loadCacheStats = async () => {
  try {
    const res = await performanceApi.getCacheStats()
    if (res.code === 0) {
      cacheStats.value = res.data
    }
  } catch (e) {
    console.error('加载缓存统计失败:', e)
  }
}

const loadDbStats = async () => {
  try {
    const res = await performanceApi.getDbStats()
    if (res.code === 0) {
      dbStats.value = res.data
    }
  } catch (e) {
    console.error('加载数据库状态失败:', e)
  }
}

const refreshAll = async () => {
  Message.info('正在刷新...')
  await Promise.all([loadSystemMetrics(), loadCacheStats(), loadDbStats()])
  Message.success('刷新成功')
}

onMounted(() => {
  refreshAll()
  timer = setInterval(refreshAll, 30000)
})

onUnmounted(() => {
  if (timer) clearInterval(timer)
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

.metric-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 8px;
}

.metric-header {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-bottom: 12px;
}

.metric-icon {
  font-size: 18px;
}

.metric-label {
  font-size: 14px;
  color: #4e5969;
  font-weight: 500;
}

.progress-text {
  font-size: 16px;
  font-weight: 600;
}

.metric-detail {
  margin-top: 8px;
  font-size: 12px;
  color: #86909c;
  text-align: center;
}

.load-avg {
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.load-item {
  display: flex;
  align-items: center;
  gap: 8px;
}

.load-label {
  width: 36px;
  font-size: 12px;
  color: #86909c;
}

.load-value {
  width: 40px;
  font-size: 12px;
  font-weight: 600;
  text-align: right;
}

.sub-metric {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 0;
}

.sub-metric-icon {
  font-size: 16px;
}

.sub-metric-label {
  font-size: 14px;
  color: #4e5969;
  font-weight: 500;
}

.sub-metric-value {
  font-size: 14px;
  color: #1d2129;
  font-weight: 600;
  margin-left: auto;
}

.cache-bar {
  margin-top: 4px;
}

.cache-bar-label {
  font-size: 11px;
  color: #86909c;
}
</style>
