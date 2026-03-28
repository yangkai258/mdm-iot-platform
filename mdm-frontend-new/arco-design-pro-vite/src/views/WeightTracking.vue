<template>
  <div class="container">
    <a-card>
      <template #title>
        <a-space><icon-scales /> 体重追踪</a-space>
      </template>

      <a-row :gutter="16" style="margin-bottom: 16px">
        <a-col :span="6">
          <a-card>
            <a-statistic title="当前体重" :value="stats.currentWeight" suffix="kg" />
          </a-card>
        </a-col>
        <a-col :span="6">
          <a-card>
            <a-statistic title="目标体重" :value="stats.targetWeight" suffix="kg" />
          </a-card>
        </a-col>
        <a-col :span="6">
          <a-card>
            <a-statistic title="变化" :value="stats.change" suffix="kg" :value-style="{ color: stats.change < 0 ? '#67C23A' : '#F56C6C' }" />
          </a-card>
        </a-col>
        <a-col :span="6">
          <a-card>
            <a-statistic title="BMI" :value="stats.bmi" />
          </a-card>
        </a-col>
      </a-row>

      <a-row :gutter="16">
        <a-col :span="12">
          <a-card title="体重趋势">
            <a-chart :option="weightChart" style="height: 250px" />
          </a-card>
        </a-col>
        <a-col :span="12">
          <a-card title="体重记录">
            <template #extra>
              <a-button type="primary" @click="handleAdd">
                <template #icon><icon-plus /></template>
                添加记录
              </a-button>
            </template>
            <a-table :columns="columns" :data="records" size="small" :pagination="false" />
          </a-card>
        </a-col>
      </a-row>

      <a-card title="营养建议" style="margin-top: 16px">
        <a-list>
          <a-list-item v-for="suggestion in suggestions" :key="suggestion.id">
            <a-list-item-meta :title="suggestion.title" :description="suggestion.content" />
          </a-list-item>
        </a-list>
      </a-card>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'

const stats = reactive({ currentWeight: 12.5, targetWeight: 11.0, change: -0.8, bmi: 22.5 })

const weightChart = reactive({
  tooltip: { trigger: 'axis' },
  xAxis: { type: 'category', data: ['周一', '周二', '周三', '周四', '周五', '周六', '周日'] },
  yAxis: { type: 'value', name: 'kg' },
  series: [{ type: 'line', smooth: true, data: [13.3, 13.2, 13.0, 12.8, 12.6, 12.5, 12.5] }]
})

const columns = [
  { title: '日期', dataIndex: 'date' },
  { title: '体重', dataIndex: 'weight' },
  { title: '变化', dataIndex: 'change' }
]
const records = ref([
  { date: '2026-03-28', weight: '12.5kg', change: '-0.1kg' }
])

const suggestions = ref([
  { id: 1, title: '控制食量', content: '建议每日减少50g狗粮摄入' },
  { id: 2, title: '增加运动', content: '建议每日增加10分钟散步时间' }
])

const handleAdd = () => { }
</script>

<style scoped>
.container { padding: 16px; }
</style>
