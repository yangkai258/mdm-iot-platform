<template>
  <div class="container">
    <a-row :gutter="16">
      <a-col :span="6">
        <a-card>
          <a-statistic title="综合健康评分" :value="healthScore" :precision="1" suffix="分">
            <template #prefix>
              <a-progress :percent="healthScore" :show-text="false" :stroke-width="8" />
            </template>
          </a-statistic>
          <a-space direction="vertical" fill style="margin-top: 16px">
            <a-tag :color="scoreColor"> {{ scoreLabel }} </a-tag>
            <span class="update-time">更新于 {{ updateTime }}</span>
          </a-space>
        </a-card>
      </a-col>
      <a-col :span="18">
        <a-card title="评分因素分析">
          <a-row :gutter="16">
            <a-col :span="8" v-for="factor in factors" :key="factor.name">
              <a-card size="small">
                <a-statistic :title="factor.name" :value="factor.score" suffix="分" :value-style="{ color: factor.color }" />
                <a-progress :percent="factor.score" :show-text="false" :stroke-color="factor.color" style="margin-top: 8px" />
              </a-card>
            </a-col>
          </a-row>
        </a-card>
      </a-col>
    </a-row>

    <a-card title="健康评分趋势 (7日)" style="margin-top: 16px">
      <template #extra>
        <a-space>
          <a-radio-group v-model="timeRange" type="button">
            <a-radio value="7d">近7天</a-radio>
            <a-radio value="30d">近30天</a-radio>
          </a-radio-group>
          <a-button @click="handleExport">
            <template #icon><icon-download /></template>
            导出
          </a-button>
        </a-space>
      </template>
      <div style="height:250px;background:#f5f5f5;border-radius:4px;display:flex;align-items:center;justify-content:center;color:#999;font-size:14px">Chart placeholder</div>
    </a-card>

    <a-card title="设备详情" style="margin-top: 16px">
      <a-descriptions :columns="3" bordered>
        <a-descriptions-item label="设备ID">{{ deviceInfo.deviceId }}</a-descriptions-item>
        <a-descriptions-item label="设备名称">{{ deviceInfo.deviceName }}</a-descriptions-item>
        <a-descriptions-item label="设备型号">{{ deviceInfo.model }}</a-descriptions-item>
        <a-descriptions-item label="固件版本">{{ deviceInfo.firmware }}</a-descriptions-item>
        <a-descriptions-item label="在线状态">
          <a-tag :color="deviceInfo.online ? 'green' : 'gray'">{{ deviceInfo.online ? '在线' : '离线' }}</a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="最后活跃">{{ deviceInfo.lastActive }}</a-descriptions-item>
      </a-descriptions>
    </a-card>
  </div>
</template>

<script setup>
import { ref, computed, reactive } from 'vue'

const deviceInfo = reactive({
  deviceId: 'D001',
  deviceName: 'M5Stack-001',
  model: 'M5Stack-Grey',
  firmware: 'v2.1.0',
  online: true,
  lastActive: '2026-03-28 10:00:00'
})

const healthScore = ref(85.5)
const updateTime = ref('2026-03-28 10:05:00')
const timeRange = ref('7d')

const scoreColor = computed(() => {
  if (healthScore.value >= 80) return 'green'
  if (healthScore.value >= 60) return 'orange'
  return 'red'
})

const scoreLabel = computed(() => {
  if (healthScore.value >= 80) return '优秀'
  if (healthScore.value >= 60) return '良好'
  return '需关注'
})

const factors = ref([
  { name: '网络稳定性', score: 92, color: '#67C23A' },
  { name: '固件状态', score: 88, color: '#409EFF' },
  { name: '电池健康', score: 75, color: '#E6A23C' },
  { name: '温度状态', score: 95, color: '#67C23A' },
  { name: '存储使用', score: 60, color: '#E6A23C' },
  { name: '连接质量', score: 82, color: '#409EFF' }
])

const chartOption = reactive({
  tooltip: { trigger: 'axis' },
  xAxis: { type: 'category', data: ['03-22', '03-23', '03-24', '03-25', '03-26', '03-27', '03-28'] },
  yAxis: { type: 'value', min: 0, max: 100 },
  series: [{ name: '健康评分', type: 'line', smooth: true, data: [82, 85, 83, 87, 84, 86, 85.5], areaStyle: {} }]
})

const handleExport = () => { }
</script>

<style scoped>
.container { padding: 16px; }
.update-time { font-size: 12px; color: #909399; }
</style>
