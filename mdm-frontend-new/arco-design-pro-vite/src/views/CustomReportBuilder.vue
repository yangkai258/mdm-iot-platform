<template>
  <div class="container">
    <a-card>
      <template #title>
        <a-space><icon-bar-chart /> 自定义报表生成器</a-space>
      </template>

      <a-row :gutter="16">
        <a-col :span="6">
          <a-card title="维度配置" size="small">
            <a-list>
              <a-list-item v-for="dim in dimensions" :key="dim.key">
                <a-checkbox v-model="dim.selected" @change="handleDimChange">{{ dim.label }}</a-checkbox>
              </a-list-item>
            </a-list>
          </a-card>
          <a-card title="指标配置" size="small" style="margin-top: 16px">
            <a-list>
              <a-list-item v-for="metric in metrics" :key="metric.key">
                <a-checkbox v-model="metric.selected" @change="handleMetricChange">{{ metric.label }}</a-checkbox>
              </a-list-item>
            </a-list>
          </a-card>
        </a-col>

        <a-col :span="12">
          <a-card title="报表预览">
            <template #extra>
              <a-space>
                <a-select v-model="chartType" placeholder="图表类型">
                  <a-option value="line">折线图</a-option>
                  <a-option value="bar">柱状图</a-option>
                  <a-option value="pie">饼图</a-option>
                  <a-option value="table">表格</a-option>
                </a-select>
              </a-space>
            </template>
            <div style="height:250px;background:#f5f5f5;border-radius:4px;display:flex;align-items:center;justify-content:center;color:#999;font-size:14px">Chart placeholder</div>
          </a-card>
        </a-col>

        <a-col :span="6">
          <a-card title="图表配置" size="small">
            <a-form :model="chartConfig" layout="vertical">
              <a-form-item label="标题">
                <a-input v-model="chartConfig.title" />
              </a-form-item>
              <a-form-item label="X轴字段">
                <a-select v-model="chartConfig.xField" placeholder="选择X轴">
                  <a-option v-for="d in selectedDims" :key="d.key" :value="d.key">{{ d.label }}</a-option>
                </a-select>
              </a-form-item>
              <a-form-item label="Y轴字段">
                <a-select v-model="chartConfig.yField" placeholder="选择Y轴">
                  <a-option v-for="m in selectedMetrics" :key="m.key" :value="m.key">{{ m.label }}</a-option>
                </a-select>
              </a-form-item>
            </a-form>
          </a-card>
        </a-col>
      </a-row>

      <a-space style="margin-top: 16px">
        <a-button type="primary" @click="handleGenerate">生成报表</a-button>
        <a-button @click="handleSaveTemplate">保存模板</a-button>
        <a-button @click="handleSubscribe">定时推送</a-button>
        <a-button @click="handleExport">
          <template #icon><icon-download /></template>
          导出
        </a-button>
      </a-space>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'

const chartType = ref('line')
const chartConfig = reactive({ title: '', xField: '', yField: '' })

const dimensions = ref([
  { key: 'date', label: '日期', selected: true },
  { key: 'device', label: '设备', selected: false },
  { key: 'region', label: '地区', selected: false }
])

const metrics = ref([
  { key: 'count', label: '数量', selected: true },
  { key: 'amount', label: '金额', selected: false },
  { key: 'rate', label: '比率', selected: false }
])

const selectedDims = computed(() => dimensions.value.filter(d => d.selected))
const selectedMetrics = computed(() => metrics.value.filter(m => m.selected))

const previewChart = reactive({
  tooltip: { trigger: 'axis' },
  xAxis: { type: 'category', data: ['周一', '周二', '周三', '周四', '周五'] },
  yAxis: { type: 'value' },
  series: [{ type: chartType.value, data: [120, 200, 150, 80, 70] }]
})

const handleDimChange = () => { }
const handleMetricChange = () => { }
const handleGenerate = () => { }
const handleSaveTemplate = () => { }
const handleSubscribe = () => { }
const handleExport = () => { }
</script>

<style scoped>
.container { padding: 16px; }
</style>
