<template>
  <div class="container">
    <a-card>
      <template #title>
        <a-space><icon-mind-mapping /> 设备状态预测</a-space>
      </template>
      <template #extra>
        <a-button @click="handleRefresh">
          <template #icon><icon-refresh /></template>
          刷新预测
        </a-button>
      </template>

      <a-row :gutter="16" style="margin-bottom: 16px">
        <a-col :span="4">
          <a-card size="small">
            <a-statistic title="预测置信度" :value="confidence" suffix="%" />
          </a-card>
        </a-col>
        <a-col :span="20">
          <a-card size="small">
            <a-space>
              <span>预测时间范围:</span>
              <a-radio-group v-model="predictRange" type="button">
                <a-radio value="1h">1小时</a-radio>
                <a-radio value="6h">6小时</a-radio>
                <a-radio value="24h">24小时</a-radio>
              </a-radio-group>
            </a-space>
          </a-card>
        </a-col>
      </a-row>

      <a-card title="预测趋势">
        <a-chart :option="predictChart" style="height: 350px" />
      </a-card>

      <a-row :gutter="16" style="margin-top: 16px">
        <a-col :span="12">
          <a-card title="异常标记">
            <a-table :columns="anomalyColumns" :data="anomalies" size="small">
              <template #severity="{ record }">
                <a-tag :color="getSeverityColor(record.severity)">{{ record.severity }}</a-tag>
              </template>
            </a-table>
          </a-card>
        </a-col>
        <a-col :span="12">
          <a-card title="预测建议">
            <a-list>
              <a-list-item v-for="suggestion in suggestions" :key="suggestion.id">
                <a-list-item-meta :title="suggestion.title" :description="suggestion.desc" />
              </a-list-item>
            </a-list>
          </a-card>
        </a-col>
      </a-row>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'

const confidence = ref(87.5)
const predictRange = ref('6h')

const predictChart = reactive({
  tooltip: { trigger: 'axis' },
  legend: { data: ['实际值', '预测值', '置信区间'] },
  xAxis: { type: 'category', data: ['10:00', '11:00', '12:00', '13:00', '14:00', '15:00', '16:00'] },
  yAxis: { type: 'value', name: '状态值' },
  series: [
    { name: '实际值', type: 'line', data: [85, 82, 88, null, null, null, null], smooth: true },
    { name: '预测值', type: 'line', data: [null, null, 88, 86, 84, 82, 80], smooth: true, linestyle: { type: 'dashed' } },
    { name: '置信区间', type: 'area', data: [null, null, 88, 89, 87, 85, 83], itemStyle: { opacity: 0.2 } }
  ]
})

const anomalyColumns = [
  { title: '时间', dataIndex: 'time' },
  { title: '类型', dataIndex: 'type' },
  { title: '严重程度', slotName: 'severity' },
  { title: '预测值', dataIndex: 'predicted' },
  { title: '实际值', dataIndex: 'actual' }
]
const anomalies = ref([
  { time: '2026-03-28 12:00', type: '温度异常', severity: 'warning', predicted: 45, actual: 52 },
  { time: '2026-03-28 13:00', type: '电量消耗过快', severity: 'info', predicted: 80, actual: 75 }
])

const suggestions = ref([
  { id: 1, title: '建议充电', desc: '电量低于20%，建议在2小时内充电' },
  { id: 2, title: '温度关注', desc: '预测显示温度将升高，建议减少高负载操作' }
])

const getSeverityColor = (s) => ({ warning: 'orange', error: 'red', info: 'blue' }[s] || 'gray')
const handleRefresh = () => { }
</script>

<style scoped>
.container { padding: 16px; }
</style>
