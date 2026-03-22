import { ref, reactive, computed, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import * as api from '@/api/notification'
import type {
  NotificationChannel,
  NotificationLog,
  AlertHistory,
  NotificationStats,
  NotificationPeriod,
  NotificationLogParams,
  AlertHistoryParams
} from '@/api/notification'

// 告警级别映射
export const SEVERITY_MAP = {
  1: { label: '提示', color: 'arcoblue' },
  2: { label: '警告', color: 'orange' },
  3: { label: '严重', color: 'red' },
  4: { label: '紧急', color: 'red' }
}

// 告警状态映射
export const STATUS_MAP = {
  0: { label: '未处理', color: 'red' },
  1: { label: '已确认', color: 'orange' },
  2: { label: '已解决', color: 'green' },
  3: { label: '已忽略', color: 'gray' }
}

// 通知状态映射
export const LOG_STATUS_MAP = {
  pending: { label: '等待中', color: 'arcoblue' },
  sent: { label: '已发送', color: 'green' },
  failed: { label: '失败', color: 'red' },
  retrying: { label: '重试中', color: 'orange' }
}

// 渠道类型映射
export const CHANNEL_TYPE_MAP = {
  email: { label: '邮件', icon: 'email' },
  sms: { label: '短信', icon: 'message' },
  webhook: { label: 'Webhook', icon: 'link' }
}

export function useNotificationChannels() {
  const loading = ref(false)
  const channels = ref<NotificationChannel[]>([])
  const currentChannel = ref<NotificationChannel | null>(null)

  // 获取各类型渠道
  const emailChannels = computed(() => channels.value.filter(ch => ch.channel_type === 'email'))
  const smsChannels = computed(() => channels.value.filter(ch => ch.channel_type === 'sms'))
  const webhookChannels = computed(() => channels.value.filter(ch => ch.channel_type === 'webhook'))

  // 加载所有渠道
  const loadChannels = async () => {
    loading.value = true
    try {
      const res = await api.getNotificationChannels()
      if (res.code === 0) {
        channels.value = res.data || []
      }
    } catch (e) {
      Message.error('加载通知渠道失败')
    } finally {
      loading.value = false
    }
  }

  // 加载单个渠道
  const loadChannel = async (id: number) => {
    try {
      const res = await api.getNotificationChannel(id)
      if (res.code === 0) {
        currentChannel.value = res.data
        return res.data
      }
    } catch (e) {
      Message.error('加载渠道详情失败')
    }
    return null
  }

  // 创建渠道
  const createChannel = async (data: Partial<NotificationChannel>) => {
    try {
      const res = await api.createNotificationChannel(data)
      if (res.code === 0) {
        Message.success('创建成功')
        await loadChannels()
        return res.data
      }
      Message.error(res.code + '')
    } catch (e) {
      Message.error('创建失败')
    }
    return null
  }

  // 更新渠道
  const updateChannel = async (id: number, data: Partial<NotificationChannel>) => {
    try {
      const res = await api.updateNotificationChannel(id, data)
      if (res.code === 0) {
        Message.success('更新成功')
        await loadChannels()
        return res.data
      }
      Message.error(res.code + '')
    } catch (e) {
      Message.error('更新失败')
    }
    return null
  }

  // 删除渠道
  const deleteChannel = async (id: number) => {
    try {
      const res = await api.deleteNotificationChannel(id)
      if (res.code === 0) {
        Message.success('删除成功')
        await loadChannels()
        return true
      }
      Message.error('删除失败')
    } catch (e) {
      Message.error('删除失败')
    }
    return false
  }

  // 测试渠道
  const testChannel = async (id: number, testData?: { recipient?: string }) => {
    try {
      const res = await api.testNotificationChannel(id, testData)
      if (res.code === 0) {
        if (res.data?.success) {
          Message.success('测试成功：' + (res.data.message || '通知发送成功'))
        } else {
          Message.error('测试失败：' + (res.data?.error || '未知错误'))
        }
        return res.data
      }
    } catch (e) {
      Message.error('测试请求失败')
    }
    return null
  }

  // 启用/禁用渠道
  const toggleChannel = async (id: number, enabled: boolean) => {
    try {
      const res = await api.toggleNotificationChannel(id, enabled)
      if (res.code === 0) {
        Message.success(enabled ? '已启用' : '已停用')
        await loadChannels()
        return true
      }
    } catch (e) {
      Message.error('操作失败')
    }
    return false
  }

  return {
    loading,
    channels,
    currentChannel,
    emailChannels,
    smsChannels,
    webhookChannels,
    loadChannels,
    loadChannel,
    createChannel,
    updateChannel,
    deleteChannel,
    testChannel,
    toggleChannel
  }
}

export function useNotificationLogs() {
  const loading = ref(false)
  const logs = ref<NotificationLog[]>([])
  const pagination = reactive({
    current: 1,
    pageSize: 20,
    total: 0
  })
  const filters = reactive<NotificationLogParams>({
    channel_type: undefined,
    status: undefined,
    start_time: undefined,
    end_time: undefined
  })

  // 加载日志
  const loadLogs = async () => {
    loading.value = true
    try {
      const res = await api.getNotificationLogs({
        ...filters,
        page: pagination.current,
        page_size: pagination.pageSize
      })
      if (res.code === 0) {
        logs.value = res.data?.list || []
        pagination.total = res.data?.total || 0
      }
    } catch (e) {
      Message.error('加载通知日志失败')
    } finally {
      loading.value = false
    }
  }

  // 查看详情
  const getLogDetail = async (id: number) => {
    try {
      const res = await api.getNotificationLogDetail(id)
      if (res.code === 0) {
        return res.data
      }
    } catch (e) {
      Message.error('加载日志详情失败')
    }
    return null
  }

  // 重置筛选
  const resetFilters = () => {
    pagination.current = 1
    filters.channel_type = undefined
    filters.status = undefined
    filters.start_time = undefined
    filters.end_time = undefined
  }

  return {
    loading,
    logs,
    pagination,
    filters,
    loadLogs,
    getLogDetail,
    resetFilters
  }
}

export function useAlertHistory() {
  const loading = ref(false)
  const alerts = ref<AlertHistory[]>([])
  const pagination = reactive({
    current: 1,
    pageSize: 20,
    total: 0
  })
  const filters = reactive<AlertHistoryParams>({
    device_id: undefined,
    alert_type: undefined,
    severity: undefined,
    status: undefined,
    start_time: undefined,
    end_time: undefined
  })

  // 加载告警历史
  const loadAlerts = async () => {
    loading.value = true
    try {
      const res = await api.getAlertHistory({
        ...filters,
        page: pagination.current,
        page_size: pagination.pageSize
      })
      if (res.code === 0) {
        alerts.value = res.data?.list || []
        pagination.total = res.data?.total || 0
      }
    } catch (e) {
      Message.error('加载告警历史失败')
    } finally {
      loading.value = false
    }
  }

  // 确认处理
  const confirmAlert = async (id: number) => {
    try {
      const res = await api.confirmAlertHistory(id)
      if (res.code === 0) {
        Message.success('已确认')
        await loadAlerts()
        return true
      }
    } catch (e) {
      Message.error('操作失败')
    }
    return false
  }

  // 解决告警
  const resolveAlert = async (id: number, remark?: string) => {
    try {
      const res = await api.resolveAlertHistory(id, remark)
      if (res.code === 0) {
        Message.success('已解决')
        await loadAlerts()
        return true
      }
    } catch (e) {
      Message.error('操作失败')
    }
    return false
  }

  // 导出
  const exportAlerts = async () => {
    try {
      const blob = await api.batchExportAlertHistory(filters)
      const url = URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = `alert-history-${Date.now()}.xlsx`
      a.click()
      URL.revokeObjectURL(url)
      Message.success('导出成功')
    } catch (e) {
      Message.error('导出失败')
    }
  }

  // 重置筛选
  const resetFilters = () => {
    pagination.current = 1
    filters.device_id = undefined
    filters.alert_type = undefined
    filters.severity = undefined
    filters.status = undefined
    filters.start_time = undefined
    filters.end_time = undefined
  }

  return {
    loading,
    alerts,
    pagination,
    filters,
    loadAlerts,
    confirmAlert,
    resolveAlert,
    exportAlerts,
    resetFilters
  }
}

export function useNotificationStats() {
  const loading = ref(false)
  const stats = ref<NotificationStats | null>(null)
  const dateRange = ref<[string, string]>(['', ''])

  // 加载统计
  const loadStats = async () => {
    loading.value = true
    try {
      const res = await api.getNotificationStats(
        dateRange.value[0] || undefined,
        dateRange.value[1] || undefined
      )
      if (res.code === 0) {
        stats.value = res.data
      }
    } catch (e) {
      Message.error('加载统计数据失败')
    } finally {
      loading.value = false
    }
  }

  return {
    loading,
    stats,
    dateRange,
    loadStats
  }
}

export function useNotificationPeriods() {
  const loading = ref(false)
  const periods = ref<NotificationPeriod[]>([])

  // 默认全天
  const defaultPeriod: NotificationPeriod = {
    name: '全天',
    start_time: '00:00',
    end_time: '23:59',
    enabled: true
  }

  // 加载时段
  const loadPeriods = async () => {
    loading.value = true
    try {
      const res = await api.getNotificationPeriods()
      if (res.code === 0) {
        periods.value = res.data || [defaultPeriod]
      }
    } catch (e) {
      periods.value = [defaultPeriod]
    } finally {
      loading.value = false
    }
  }

  // 保存时段
  const savePeriods = async (newPeriods: NotificationPeriod[]) => {
    try {
      const res = await api.updateNotificationPeriods(newPeriods)
      if (res.code === 0) {
        Message.success('保存成功')
        periods.value = newPeriods
        return true
      }
    } catch (e) {
      Message.error('保存失败')
    }
    return false
  }

  return {
    loading,
    periods,
    defaultPeriod,
    loadPeriods,
    savePeriods
  }
}
