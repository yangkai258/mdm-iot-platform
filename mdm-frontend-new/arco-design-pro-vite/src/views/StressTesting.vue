<template>
  <div class="container">
    <a-card>
      <template #title>
        <a-space><icon-launch /> 压力测试</a-space>
      </template>

      <a-row :gutter="16">
        <a-col :span="12">
          <a-card title="测试配置">
            <a-form :model="formData" layout="vertical">
              <a-form-item label="并发用户数">
                <a-input-number v-model="formData.concurrentUsers" :min="1" :max="10000" />
              </a-form-item>
              <a-form-item label="持续时间">
                <a-input-number v-model="formData.duration" :min="1" :max="3600" /> 秒
              </a-form-item>
              <a-form-item label="QPS目标">
                <a-input-number v-model="formData.targetQps" :min="1" />
              </a-form-item>
              <a-form-item>
                <a-button type="primary" @click="handleStart">开始测试</a-button>
                <a-button v-if="isRunning" status="danger" @click="handleStop">停止</a-button>
              </a-form-item>
            </a-form>
          </a-card>
        </a-col>

        <a-col :span="12">
          <a-card title="实时指标">
            <a-row :gutter="16">
              <a-col :span="12">
                <a-statistic title="当前QPS" :value="metrics.currentQps" />
              </a-col>
              <a-col :span="12">
                <a-statistic title="平均响应时间" :value="metrics.avgResponseTime" suffix="ms" />
              </a-col>
              <a-col :span="12">
                <a-statistic title="成功率" :value="metrics.successRate" suffix="%" />
              </a-col>
              <a-col :span="12">
                <a-statistic title="总请求数" :value="metrics.totalRequests" />
              </a-col>
            </a-row>
          </a-card>
        </a-col>
      </a-row>

      <a-card title="测试结果" style="margin-top: 16px">
        <a-chart :option="resultChart" style="height: 300px" />
      </a-card>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'

const formData = reactive({ concurrentUsers: 100, duration: 60, targetQps: 1000 })
const metrics = reactive({ currentQps: 0, avgResponseTime: 0, successRate: 100, totalRequests: 0 })
const isRunning = ref(false)

const resultChart = reactive({
  tooltip: { trigger: 'axis' },
  legend: { data: ['QPS', '响应时间'] },
  xAxis: { type: 'category', data: [] },
  yAxis: [{ type: 'value', name: 'QPS' }, { type: 'value', name: 'ms' }],
  series: []
})

const handleStart = () => { isRunning.value = true }
const handleStop = () => { isRunning.value = false }
</script>

<style scoped>
.container { padding: 16px; }
</style>
