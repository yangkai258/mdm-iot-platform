<template>
  <div class="container">
    <a-row :gutter="16" style="margin-bottom: 16px">
      <a-col :span="6">
        <a-card>
          <a-statistic title="告警总数" :value="stats.total" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card>
          <a-statistic title="今日新增" :value="stats.today" suffix="条" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card>
          <a-statistic title="未处理" :value="stats.unprocessed" suffix="条" :value-style="{ color: '#F56C6C' }" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card>
          <a-statistic title="处理率" :value="stats.processRate" suffix="%" />
        </a-card>
      </a-col>
    </a-row>

    <a-row :gutter="16">
      <a-col :span="16">
        <a-card title="告警趋势">
          <template #extra>
            <a-space>
              <a-radio-group v-model="timeRange" type="button">
                <a-radio value="today">今日</a-radio>
                <a-radio value="week">本周</a-radio>
                <a-radio value="month">本月</a-radio>
              </a-radio-group>
            </a-space>
          </template>
          <a-chart :option="trendChart" style="height: 250px" />
        </a-card>
      </a-col>
      <a-col :span="8">
        <a-card title="告警类型分布">
          <a-chart :option="pieChart" style="height: 250px" />
        </a-card>
      </a-col>
    </a-row>

    <a-card title="Top 5 告警源" style="margin-top: 16px">
      <a-table :columns="topColumns" :data="topDevices" size="small">
        <template #count="{ record }">
          <a-progress :percent="record.percent" :show-text="true" size="small" />
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'

const stats = reactive({ total: 1256, today: 45, unprocessed: 23, processRate: 98.2 })
const timeRange = ref('week')

const trendChart = reactive({
  tooltip: { trigger: 'axis' },
  legend: { data: ['告警数', '处理数'] },
  xAxis: { type: 'category', data: ['周一', '周二', '周三', '周四', '周五', '周六', '周日'] },
  yAxis: { type: 'value' },
  series: [
    { name: '告警数', type: 'bar', data: [42, 55, 38, 62, 48, 35, 45] },
    { name: '处理数', type: 'bar', data: [40, 53, 38, 60, 46, 35, 43] }
  ]
})

const pieChart = reactive({
  tooltip: { trigger: 'item' },
  series: [{ type: 'pie', radius: '60%', data: [
    { name: '设备离线', value: 35 },
    { name: '温度异常', value: 25 },
    { name: '电量不足', value: 20 },
    { name: '固件升级', value: 15 },
    { name: '其他', value: 5 }
  ]}]
})

const topColumns = [
  { title: '设备ID', dataIndex: 'deviceId' },
  { title: '设备名称', dataIndex: 'deviceName' },
  { title: '告警次数', slotName: 'count' },
  { title: '最近告警', dataIndex: 'lastAlert' }
]
const topDevices = ref([
  { deviceId: 'D001', deviceName: 'M5Stack-001', count: 45, percent: 90, lastAlert: '2026-03-28 10:00' },
  { deviceId: 'D002', deviceName: 'M5Stack-002', count: 38, percent: 76, lastAlert: '2026-03-28 09:30' }
])
</script>

<style scoped>
.container { padding: 16px; }
</style>
