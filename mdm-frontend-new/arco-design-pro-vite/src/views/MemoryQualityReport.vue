<template>
  <div class="container">
    <a-card>
      <template #title>
        <a-space><icon-file-text /> 记忆质量报告</a-space>
      </template>

      <a-row :gutter="16" style="margin-bottom: 16px">
        <a-col :span="6">
          <a-card>
            <a-statistic title="质量评分" :value="quality.score" suffix="分">
              <template #prefix>
                <a-progress :percent="quality.score" :show-text="false" :stroke-width="10" />
              </template>
            </a-statistic>
          </a-card>
        </a-col>
        <a-col :span="6">
          <a-card>
            <a-statistic title="完整性" :value="quality.completeness" suffix="%" />
          </a-card>
        </a-col>
        <a-col :span="6">
          <a-card>
            <a-statistic title="准确性" :value="quality.accuracy" suffix="%" />
          </a-card>
        </a-col>
        <a-col :span="6">
          <a-card>
            <a-statistic title="时效性" :value="quality.timeliness" suffix="%" />
          </a-card>
        </a-col>
      </a-row>

      <a-row :gutter="16">
        <a-col :span="12">
          <a-card title="质量雷达图">
            <div style="height:250px;background:#f5f5f5;border-radius:4px;display:flex;align-items:center;justify-content:center;color:#999;font-size:14px">Chart placeholder</div>
          </a-card>
        </a-col>
        <a-col :span="12">
          <a-card title="问题记忆列表">
            <a-list>
              <a-list-item v-for="issue in issues" :key="issue.id">
                <a-list-item-meta :title="issue.title" :description="issue.description" />
                <template #actions>
                  <a-link @click="handleFix(issue)">修复</a-link>
                </template>
              </a-list-item>
            </a-list>
          </a-card>
        </a-col>
      </a-row>

      <a-space style="margin-top: 16px">
        <a-button type="primary" @click="handleGenerate">生成报告</a-button>
        <a-button @click="handleExport">
          <template #icon><icon-download /></template>
          导出报告
        </a-button>
      </a-space>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'

const quality = reactive({ score: 85, completeness: 92, accuracy: 88, timeliness: 78 })

const radarChart = reactive({
  radar: {
    indicator: [
      { name: '完整性', max: 100 }, { name: '准确性', max: 100 },
      { name: '时效性', max: 100 }, { name: '一致性', max: 100 }
    ]
  },
  series: [{ type: 'radar', data: [{ value: [92, 88, 78, 85], name: '质量评分' }] }]
})

const issues = ref([
  { id: 1, title: '记忆片段丢失', description: '3月份部分记忆数据不完整' },
  { id: 2, title: '时间戳错误', description: '部分记忆的时间戳存在偏差' }
])

const handleGenerate = () => { }
const handleExport = () => { }
const handleFix = (issue) => { }
</script>

<style scoped>
.container { padding: 16px; }
</style>
