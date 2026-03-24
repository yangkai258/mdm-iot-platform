<template>
  <div class="page-container">
    <div class="page-header">
      <div class="page-title">全球化设置</div>
      <div class="page-desc">统一管理多区域部署、时区配置和数据驻留策略</div>
    </div>
    <div class="quick-entries">
      <div class="entry-card" v-for="entry in entries" :key="entry.path" @click="goTo(entry.path)">
        <div class="entry-icon">
          <component :is="entry.icon" />
        </div>
        <div class="entry-info">
          <div class="entry-name">{{ entry.name }}</div>
          <div class="entry-desc">{{ entry.desc }}</div>
        </div>
        <div class="entry-status">
          <a-badge :color="entry.statusColor" :text="entry.statusText" />
        </div>
        <icon-right class="entry-arrow" />
      </div>
    </div>
    <div class="config-status">
      <div class="section-title">当前配置状态</div>
      <div class="config-grid">
        <div class="config-item">
          <div class="config-label">当前时区</div>
          <div class="config-value">{{ currentTz || '未设置' }}</div>
        </div>
        <div class="config-item">
          <div class="config-label">区域数量</div>
          <div class="config-value">{{ regionCount }} 个</div>
        </div>
        <div class="config-item">
          <div class="config-label">活跃 AI 节点</div>
          <div class="config-value">{{ activeNodes }} 个</div>
        </div>
        <div class="config-item">
          <div class="config-label">数据驻留规则</div>
          <div class="config-value">{{ ruleCount }} 条</div>
        </div>
        <div class="config-item">
          <div class="config-label">默认区域</div>
          <div class="config-value">{{ defaultRegion || '未设置' }}</div>
        </div>
        <div class="config-item">
          <div class="config-label">同步状态</div>
          <div class="config-value">
            <a-badge :color="syncStatusColor" :text="syncStatusText" />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getRegions, getTimezone, getAINodes, getDataResidencyRules, getSyncStatus } from '@/api/globalization'

const router = useRouter()
const currentTz = ref('')
const regionCount = ref(0)
const activeNodes = ref(0)
const ruleCount = ref(0)
const defaultRegion = ref('')
const syncStatusColor = ref('gray')
const syncStatusText = ref('未知')

const entries = [
  {
    path: '/globalization/region',
    name: '区域管理',
    desc: '配置和管理多区域数据库架构',
    icon: 'IconApps',
    statusColor: 'green',
    statusText: '已配置'
  },
  {
    path: '/globalization/timezone',
    name: '时区设置',
    desc: '系统、租户、部门的时区配置',
    icon: 'IconClockCircle',
    statusColor: 'green',
    statusText: '已配置'
  },
  {
    path: '/globalization/data-residency',
    name: '数据驻留规则',
    desc: '配置数据类型与存储区域的映射',
    icon: 'IconSafe',
    statusColor: 'green',
    statusText: '已配置'
  },
  {
    path: '/globalization/ai-node',
    name: '区域 AI 节点',
    desc: '监控各区域的 AI 推理节点状态',
    icon: 'IconRobot',
    statusColor: 'green',
    statusText: '已配置'
  },
  {
    path: '/globalization/sync-status',
    name: '跨区域同步状态',
    desc: '查看跨区域数据同步状态和手动触发',
    icon: 'IconSync',
    statusColor: 'green',
    statusText: '已配置'
  }
]

function goTo(path) {
  router.push(path)
}

onMounted(async () => {
  try {
    const [tzRes, regionsRes, nodesRes, rulesRes, syncRes] = await Promise.all([
      getTimezone(),
      getRegions(),
      getAINodes(),
      getDataResidencyRules(),
      getSyncStatus()
    ])

    const tzData = tzRes.data || tzRes
    currentTz.value = tzData?.timezone || ''

    const regions = regionsRes.data || regionsRes || []
    regionCount.value = regions.length
    const defaultReg = regions.find(r => r.is_default)
    defaultRegion.value = defaultReg ? defaultReg.region_name : ''

    const nodes = nodesRes.data || nodesRes || []
    activeNodes.value = nodes.filter(n => n.health_status === 'online').length

    const rules = rulesRes.data || rulesRes || []
    ruleCount.value = rules.length

    const syncData = syncRes.data || syncRes
    if (syncData) {
      syncStatusColor.value = syncData.status === 'healthy' ? 'green' : 'orange'
      syncStatusText.value = syncData.status === 'healthy' ? '正常' : '异常'
    }
  } catch (e) {
    console.error('加载配置状态失败', e)
  }
})
</script>

<style scoped>
.page-container { background: #fff; border-radius: 4px; padding: 20px; }
.page-header { margin-bottom: 20px; }
.page-title { font-size: 18px; font-weight: 600; margin-bottom: 4px; }
.page-desc { font-size: 13px; color: var(--color-text-3); }
.quick-entries { display: grid; grid-template-columns: repeat(2, 1fr); gap: 12px; margin-bottom: 20px; }
.entry-card {
  display: flex; align-items: center; gap: 12px; padding: 16px;
  background: var(--color-bg-white); border: 1px solid var(--color-border);
  border-radius: 8px; cursor: pointer; transition: all 0.2s;
}
.entry-card:hover { border-color: rgb(var(--primary-6)); background: var(--color-fill-1); }
.entry-icon {
  width: 44px; height: 44px; border-radius: 10px; background: var(--color-fill-2);
  display: flex; align-items: center; justify-content: center;
  font-size: 22px; color: var(--color-text-3); flex-shrink: 0;
}
.entry-info { flex: 1; }
.entry-name { font-size: 14px; font-weight: 600; margin-bottom: 2px; }
.entry-desc { font-size: 12px; color: var(--color-text-3); }
.entry-status { margin-right: 8px; }
.entry-arrow { color: var(--color-text-3); }
.config-status { padding: 16px; background: #f7f8fa; border-radius: 4px; }
.section-title { font-size: 14px; font-weight: 600; margin-bottom: 16px; }
.config-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 16px; }
.config-item { display: flex; flex-direction: column; gap: 4px; }
.config-label { font-size: 12px; color: var(--color-text-3); }
.config-value { font-size: 14px; font-weight: 500; }
</style>
