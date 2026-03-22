/**
 * AI 质量监控 Composable
 */
import { ref, computed } from 'vue'
import {
  getAiQualityMetrics,
  getAiQualityMetricsTrend,
  getAiMonitorStats,
  getAiMonitorEvents
} from '@/api/ai'

export function useAIQuality() {
  // ============ 质量指标 ============
  const metricsLoading = ref(false)
  const metrics = ref({
    total_inferences: 0,
    avg_latency_ms: 0,
    error_rate: 0,
    avg_confidence: 0,
    // 环比变化
    total_inferences_change: 0,
    avg_latency_change: 0,
    error_rate_change: 0,
    avg_confidence_change: 0
  })

  // ============ 趋势数据 ============
  const trendLoading = ref(false)
  const latencyTrend = ref([])   // [{ time, value }]
  const errorRateTrend = ref([]) // [{ time, value }]

  // ============ 行为日志 ============
  const logLoading = ref(false)
  const logList = ref([])
  const logTotal = ref(0)

  // ============ 辅助方法 ============
  const timeRangeOptions = [
    { label: '最近1小时', value: '1h' },
    { label: '最近24小时', value: '24h' },
    { label: '最近7天', value: '7d' },
    { label: '自定义', value: 'custom' }
  ]

  /**
   * 加载质量指标
   */
  const loadMetrics = async (params = {}) => {
    metricsLoading.value = true
    try {
      const res = await getAiQualityMetrics(params)
      if (res.code === 0) {
        metrics.value = res.data
      }
    } finally {
      metricsLoading.value = false
    }
  }

  /**
   * 加载趋势数据
   */
  const loadTrend = async (params = {}) => {
    trendLoading.value = true
    try {
      const res = await getAiQualityMetricsTrend(params)
      if (res.code === 0) {
        latencyTrend.value = res.data.latency_trend || []
        errorRateTrend.value = res.data.error_rate_trend || []
      }
    } finally {
      trendLoading.value = false
    }
  }

  /**
   * 加载行为日志列表
   */
  const loadLogs = async (params = {}) => {
    logLoading.value = true
    try {
      const res = await getAiMonitorEvents(params)
      if (res.code === 0) {
        logList.value = res.data.list || []
        logTotal.value = res.data.total || 0
      }
    } finally {
      logLoading.value = false
    }
  }

  /**
   * 加载统计数据
   */
  const loadStats = async (params = {}) => {
    const res = await getAiMonitorStats(params)
    return res
  }

  /**
   * 格式化趋势数据（用于图表）
   */
  const formatTrendData = (data, key = 'value') => {
    return data.map(item => ({
      time: item.time || item.created_at,
      value: item[key] ?? item
    }))
  }

  /**
   * 模型调用分布
   */
  const modelDistribution = computed(() => {
    // 从日志列表聚合模型分布
    const dist = {}
    logList.value.forEach(log => {
      const model = log.model_version || 'unknown'
      dist[model] = (dist[model] || 0) + 1
    })
    return Object.entries(dist).map(([name, value]) => ({ name, value }))
  })

  /**
   * 异常事件告警
   */
  const anomalyAlerts = computed(() => {
    return logList.value
      .filter(log => log.is_anomaly)
      .slice(0, 10)
      .map(log => ({
        id: log.id,
        device_id: log.device_id,
        model_version: log.model_version,
        behavior_type: log.behavior_type,
        anomaly_score: log.anomaly_score,
        created_at: log.created_at
      }))
  })

  return {
    // 状态
    metricsLoading,
    metrics,
    trendLoading,
    latencyTrend,
    errorRateTrend,
    logLoading,
    logList,
    logTotal,

    // 辅助
    timeRangeOptions,
    modelDistribution,
    anomalyAlerts,

    // 方法
    loadMetrics,
    loadTrend,
    loadLogs,
    loadStats,
    formatTrendData
  }
}
