<template>
  <div class="container">
    <a-card>
      <template #title>
        <a-space><icon-mind-mapping /> 行为预测</a-space>
      </template>

      <a-row :gutter="16" style="margin-bottom: 16px">
        <a-col :span="6">
          <a-card>
            <a-statistic title="预测准确率" :value="stats.accuracy" suffix="%" />
          </a-card>
        </a-col>
        <a-col :span="6">
          <a-card>
            <a-statistic title="预测次数" :value="stats.predictions" />
          </a-card>
        </a-col>
        <a-col :span="6">
          <a-card>
            <a-statistic title="模型版本" value="v2.1" />
          </a-card>
        </a-col>
      </a-row>

      <a-row :gutter="16">
        <a-col :span="12">
          <a-card title="行为预测趋势">
            <div style="height:250px;background:#f5f5f5;border-radius:4px;display:flex;align-items:center;justify-content:center;color:#999;font-size:14px">Chart placeholder</div>
          </a-card>
        </a-col>
        <a-col :span="12">
          <a-card title="预测列表">
            <a-table :columns="columns" :data="predictions" size="small">
              <template #confidence="{ record }">
                <a-progress :percent="record.confidence" :show-text="true" size="small" />
              </template>
            </a-table>
          </a-card>
        </a-col>
      </a-row>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'

const stats = reactive({ accuracy: 78, predictions: 1250 })

const trendChart = reactive({
  tooltip: { trigger: 'axis' },
  xAxis: { type: 'category', data: ['周一', '周二', '周三', '周四', '周五', '周六', '周日'] },
  yAxis: { type: 'value', min: 0, max: 100, name: '准确率%' },
  series: [{ type: 'line', smooth: true, data: [72, 75, 78, 76, 80, 82, 78], areaStyle: {} }]
})

const columns = [
  { title: '时间', dataIndex: 'time' },
  { title: '预测行为', dataIndex: 'behavior' },
  { title: '置信度', slotName: 'confidence' },
  { title: '实际行为', dataIndex: 'actual' }
]
const predictions = ref([
  { time: '2026-03-28 10:00', behavior: '进食', confidence: 85, actual: '进食' }
])
</script>

<style scoped>
.container { padding: 16px; }
</style>
