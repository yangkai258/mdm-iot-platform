<template>
  <div class="container">
    <a-card>
      <template #title>
        <a-space><icon-file-text /> 情绪报告</a-space>
      </template>

      <a-tabs default-active-key="weekly">
        <a-tab-pane key="weekly" title="周报">
          <a-card>
            <a-descriptions :column="3" bordered>
              <a-descriptions-item label="报告周期">2026-03-22 至 2026-03-28</a-descriptions-item>
              <a-descriptions-item label="生成时间">2026-03-28 10:00</a-descriptions-item>
              <a-descriptions-item label="宠物名称">小白</a-descriptions-item>
            </a-descriptions>

            <a-row :gutter="16" style="margin-top: 24px">
              <a-col :span="12">
                <a-card title="情绪分布">
                  <div style="height:250px;background:#f5f5f5;border-radius:4px;display:flex;align-items:center;justify-content:center;color:#999;font-size:14px">Chart placeholder</div>
                </a-card>
              </a-col>
              <a-col :span="12">
                <a-card title="情绪趋势">
                  <div style="height:250px;background:#f5f5f5;border-radius:4px;display:flex;align-items:center;justify-content:center;color:#999;font-size:14px">Chart placeholder</div>
                </a-card>
              </a-col>
            </a-row>

            <a-card title="高频情绪事件" style="margin-top: 16px">
              <a-list>
                <a-list-item v-for="event in events" :key="event.id">
                  <a-list-item-meta :title="event.title" :description="event.time + ' | ' + event.emotion" />
                  <template #actions>
                    <a-tag :color="getEmotionColor(event.emotionType)">{{ event.emotionType }}</a-tag>
                  </template>
                </a-list-item>
              </a-list>
            </a-card>
          </a-card>
        </a-tab-pane>

        <a-tab-pane key="monthly" title="月报">
          <a-empty description="月报功能开发中" />
        </a-tab-pane>
      </a-tabs>

      <a-space style="margin-top: 16px">
        <a-button type="primary" @click="handleShare">
          <template #icon><icon-share /></template>
          分享报告
        </a-button>
        <a-button @click="handleExportPDF">
          <template #icon><icon-download /></template>
          导出PDF
        </a-button>
      </a-space>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'

const events = ref([
  { id: 1, title: '与主人长时间互动', time: '2026-03-28 10:00', emotion: '开心', emotionType: 'happy' },
  { id: 2, title: '独自休息', time: '2026-03-27 14:00', emotion: '平静', emotionType: 'neutral' }
])

const emotionPie = reactive({
  tooltip: { trigger: 'item' },
  series: [{ type: 'pie', radius: '60%', data: [
    { name: '开心', value: 45 }, { name: '平静', value: 30 }, { name: '难过', value: 15 }, { name: '生气', value: 10 }
  ]}]
})

const emotionTrend = reactive({
  tooltip: { trigger: 'axis' },
  xAxis: { type: 'category', data: ['周一', '周二', '周三', '周四', '周五', '周六', '周日'] },
  yAxis: { type: 'value', min: 0, max: 100 },
  series: [{ type: 'line', smooth: true, data: [72, 75, 68, 80, 78, 85, 82], areaStyle: {} }]
})

const getEmotionColor = (type) => ({
  happy: '#67C23A', neutral: '#909399', sad: '#409EFF', angry: '#F56C6C'
}[type] || 'gray')

const handleShare = () => { }
const handleExportPDF = () => { }
</script>

<style scoped>
.container { padding: 16px; }
</style>
