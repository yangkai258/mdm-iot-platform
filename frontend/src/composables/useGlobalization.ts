import { ref, computed } from 'vue'
import {
  getRegions,
  getRegion,
  createRegion,
  updateRegion,
  deleteRegion,
  getRegionNodes,
  createRegionNode,
  deleteRegionNode,
  getRegionHealth,
  getTimezone,
  updateTimezone,
  getTimezoneList,
  getDataResidencyRules,
  getDataResidencyRule,
  createDataResidencyRule,
  updateDataResidencyRule,
  deleteDataResidencyRule,
  getAINodes,
  getAINode,
  registerAINode,
  deregisterAINode,
  getSyncRecords,
  triggerSync,
  getSyncStatus
} from '@/api/globalization'

// 时区列表（IANA 标准）
export const TIMEZONE_LIST = [
  { value: 'Pacific/Honolulu', label: 'Pacific/Honolulu (UTC-10)', offset: -10 },
  { value: 'America/Anchorage', label: 'America/Anchorage (UTC-9)', offset: -9 },
  { value: 'America/Los_Angeles', label: 'America/Los_Angeles (UTC-8)', offset: -8 },
  { value: 'America/Denver', label: 'America/Denver (UTC-7)', offset: -7 },
  { value: 'America/Chicago', label: 'America/Chicago (UTC-6)', offset: -6 },
  { value: 'America/New_York', label: 'America/New_York (UTC-5)', offset: -5 },
  { value: 'America/Sao_Paulo', label: 'America/Sao_Paulo (UTC-3)', offset: -3 },
  { value: 'UTC', label: 'UTC (UTC+0)', offset: 0 },
  { value: 'Europe/London', label: 'Europe/London (UTC+0)', offset: 0 },
  { value: 'Europe/Paris', label: 'Europe/Paris (UTC+1)', offset: 1 },
  { value: 'Europe/Berlin', label: 'Europe/Berlin (UTC+1)', offset: 1 },
  { value: 'Europe/Moscow', label: 'Europe/Moscow (UTC+3)', offset: 3 },
  { value: 'Asia/Dubai', label: 'Asia/Dubai (UTC+4)', offset: 4 },
  { value: 'Asia/Karachi', label: 'Asia/Karachi (UTC+5)', offset: 5 },
  { value: 'Asia/Kolkata', label: 'Asia/Kolkata (UTC+5:30)', offset: 5.5 },
  { value: 'Asia/Dhaka', label: 'Asia/Dhaka (UTC+6)', offset: 6 },
  { value: 'Asia/Bangkok', label: 'Asia/Bangkok (UTC+7)', offset: 7 },
  { value: 'Asia/Shanghai', label: 'Asia/Shanghai (UTC+8)', offset: 8 },
  { value: 'Asia/Singapore', label: 'Asia/Singapore (UTC+8)', offset: 8 },
  { value: 'Asia/Tokyo', label: 'Asia/Tokyo (UTC+9)', offset: 9 },
  { value: 'Asia/Seoul', label: 'Asia/Seoul (UTC+9)', offset: 9 },
  { value: 'Australia/Sydney', label: 'Australia/Sydney (UTC+10)', offset: 10 },
  { value: 'Pacific/Auckland', label: 'Pacific/Auckland (UTC+12)', offset: 12 }
]

// 区域类型
export const REGION_TYPES = [
  { value: 'primary', label: '主区域' },
  { value: 'backup', label: '备份区域' },
  { value: 'ai_node', label: '推理节点' },
  { value: 'storage', label: '存储节点' }
]

// 数据类型
export const DATA_TYPES = [
  { value: 'user_data', label: '用户数据' },
  { value: 'device_data', label: '设备数据' },
  { value: 'ai_training_data', label: 'AI训练数据' },
  { value: 'ai_inference_data', label: 'AI推理数据' },
  { value: 'log_data', label: '日志数据' }
]

// 状态映射
export function statusColor(status: string) {
  const map: Record<string, string> = {
    active: 'green',
    inactive: 'gray',
    pending: 'yellow',
    online: 'green',
    offline: 'red',
    standby: 'orange'
  }
  return map[status] || 'default'
}

export function statusLabel(status: string) {
  const map: Record<string, string> = {
    active: '活跃',
    inactive: '停用',
    pending: '待审批',
    online: '在线',
    offline: '离线',
    standby: '备用'
  }
  return map[status] || status
}

// 节点类型
export const NODE_TYPES = [
  { value: 'ai', label: 'AI推理节点' },
  { value: 'storage', label: '存储节点' },
  { value: 'api', label: 'API节点' },
  { value: 'mqtt', label: 'MQTT节点' }
]

export function useGlobalization() {
  // 区域管理
  const regions = ref<any[]>([])
  const regionsLoading = ref(false)

  async function loadRegions(params?: any) {
    regionsLoading.value = true
    try {
      const res = await getRegions(params)
      regions.value = res.data || res || []
    } catch (e) {
      console.error('加载区域列表失败', e)
    } finally {
      regionsLoading.value = false
    }
  }

  // 时区
  const currentTimezone = ref<any>(null)
  const timezoneList = ref(TIMEZONE_LIST)

  async function loadTimezone() {
    try {
      const res = await getTimezone()
      currentTimezone.value = res.data || res
    } catch (e) {
      console.error('加载时区失败', e)
    }
  }

  async function saveTimezone(data: any) {
    await updateTimezone(data)
  }

  // 数据驻留
  const residencyRules = ref<any[]>([])
  const residencyLoading = ref(false)

  async function loadResidencyRules(params?: any) {
    residencyLoading.value = true
    try {
      const res = await getDataResidencyRules(params)
      residencyRules.value = res.data || res || []
    } catch (e) {
      console.error('加载数据驻留规则失败', e)
    } finally {
      residencyLoading.value = false
    }
  }

  // AI 节点
  const aiNodes = ref<any[]>([])
  const aiNodesLoading = ref(false)

  async function loadAINodes(params?: any) {
    aiNodesLoading.value = true
    try {
      const res = await getAINodes(params)
      aiNodes.value = res.data || res || []
    } catch (e) {
      console.error('加载AI节点失败', e)
    } finally {
      aiNodesLoading.value = false
    }
  }

  // 同步状态
  const syncRecords = ref<any[]>([])
  const syncLoading = ref(false)

  async function loadSyncRecords(params?: any) {
    syncLoading.value = true
    try {
      const res = await getSyncRecords(params)
      syncRecords.value = res.data || res || []
    } catch (e) {
      console.error('加载同步记录失败', e)
    } finally {
      syncLoading.value = false
    }
  }

  return {
    // 区域
    regions,
    regionsLoading,
    loadRegions,
    // 时区
    currentTimezone,
    timezoneList,
    loadTimezone,
    saveTimezone,
    // 数据驻留
    residencyRules,
    residencyLoading,
    loadResidencyRules,
    // AI 节点
    aiNodes,
    aiNodesLoading,
    loadAINodes,
    // 同步
    syncRecords,
    syncLoading,
    loadSyncRecords
  }
}
