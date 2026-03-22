/**
 * 性能管理 API
 * 系统性能指标、缓存管理、数据库状态
 */

const BASE_URL = '/api/v1/performance'

export interface SystemMetrics {
  cpu_usage: number       // CPU 使用率 0-100
  memory_usage: number    // 内存使用率 0-100
  memory_total: number     // 内存总量 MB
  memory_used: number      // 已用内存 MB
  disk_usage: number       // 磁盘使用率 0-100
  disk_total: number       // 磁盘总量 GB
  disk_used: number        // 已用磁盘 GB
  network_rx: number       // 网络下行速率 B/s
  network_tx: number        // 网络上行速率 B/s
  uptime: number           // 运行时间（秒）
  load_avg: number[]       // 系统负载 [1min, 5min, 15min]
  goroutines: number       // Go 协程数
  timestamp: string
}

export interface CacheStats {
  hits: number             // 缓存命中次数
  misses: number           // 缓存未命中次数
  hit_rate: number          // 命中率 0-100
  keys: number              // 总 key 数
  memory_used: number       // 内存使用量 bytes
  memory_total: number      // 内存上限 bytes
  cmd_get: number           // GET 命令次数
  cmd_set: number           // SET 命令次数
  cmd_del: number           // DEL 命令次数
  evictions: number         // 驱逐次数
  expired: number           // 过期 key 数
}

export interface DbStats {
  status: string            // connected | disconnected | error
  version: string           // PostgreSQL 版本
  connections: number       // 当前连接数
  connections_max: number   // 最大连接数
  queries_total: number     // 总查询数
  queries_per_sec: number   // QPS
  transaction_per_sec: number // TPS
  cache_hit_rate: number    // 缓存命中率
  slow_queries: number      // 慢查询数
  latency_ms: number        // 平均延迟 ms
  databases: Array<{
    name: string
    size_mb: number
    tables: number
  }>
}

const performanceApi = {
  // 系统性能指标
  getSystemMetrics: () => {
    const token = localStorage.getItem('token')
    return fetch(`${BASE_URL}/system`, {
      headers: { 'Authorization': `Bearer ${token}` }
    }).then(r => r.json())
  },

  // 缓存统计
  getCacheStats: () => {
    const token = localStorage.getItem('token')
    return fetch(`${BASE_URL}/cache`, {
      headers: { 'Authorization': `Bearer ${token}` }
    }).then(r => r.json())
  },

  // 数据库状态
  getDbStats: () => {
    const token = localStorage.getItem('token')
    return fetch(`${BASE_URL}/database`, {
      headers: { 'Authorization': `Bearer ${token}` }
    }).then(r => r.json())
  },

  // 清空缓存
  clearCache: (type?: 'all' | 'memory' | 'expired') => {
    const token = localStorage.getItem('token')
    return fetch(`${BASE_URL}/cache/clear`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ type: type || 'all' })
    }).then(r => r.json())
  },

  // 历史性能数据
  getMetricsHistory: (params?: { start?: string; end?: string; interval?: string }) => {
    const token = localStorage.getItem('token')
    const qs = params ? '?' + new URLSearchParams(params as any).toString() : ''
    return fetch(`${BASE_URL}/history${qs}`, {
      headers: { 'Authorization': `Bearer ${token}` }
    }).then(r => r.json())
  }
}

export default performanceApi
