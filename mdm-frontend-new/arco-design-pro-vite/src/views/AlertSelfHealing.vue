<template>
  <div class="container">
    <a-card>
      <template #title>
        <a-space><icon-bot /> 自愈建议</a-space>
      </template>

      <a-row :gutter="16">
        <a-col :span="16">
          <a-card title="告警详情">
            <a-descriptions :column="2" bordered>
              <a-descriptions-item label="告警ID">{{ alert.id }}</a-descriptions-item>
              <a-descriptions-item label="告警类型">
                <a-tag color="orange">{{ alert.type }}</a-tag>
              </a-descriptions-item>
              <a-descriptions-item label="设备ID">{{ alert.deviceId }}</a-descriptions-item>
              <a-descriptions-item label="设备名称">{{ alert.deviceName }}</a-descriptions-item>
              <a-descriptions-item label="触发时间">{{ alert.triggerTime }}</a-descriptions-item>
              <a-descriptions-item label="告警级别">
                <a-tag :color="getLevelColor(alert.level)">{{ alert.level }}</a-tag>
              </a-descriptions-item>
              <a-descriptions-item label="告警内容" :span="2">{{ alert.message }}</a-descriptions-item>
            </a-descriptions>
          </a-card>

          <a-card title="自愈建议" style="margin-top: 16px">
            <a-list>
              <a-list-item v-for="suggestion in suggestions" :key="suggestion.id">
                <a-card size="small" :class="{ selected: selectedSuggestion?.id === suggestion.id }">
                  <template #title>
                    <a-space>
                      <icon-check-circle :style="{ color: suggestion.successRate > 80 ? '#67C23A' : '#E6A23C' }" />
                      {{ suggestion.title }}
                      <a-tag size="small">{{ suggestion.successRate }}% 成功率</a-tag>
                    </a-space>
                  </template>
                  <a-space direction="vertical" fill>
                    <p>{{ suggestion.description }}</p>
                    <a-space>
                      <a-button type="primary" size="small" @click="handleExecute(suggestion)">
                        <template #icon><icon-play /></template>
                        执行建议
                      </a-button>
                      <a-button size="small" @click="handleViewHistory(suggestion)">查看历史</a-button>
                    </a-space>
                  </a-space>
                </a-card>
              </a-list-item>
            </a-list>
          </a-card>
        </a-col>

        <a-col :span="8">
          <a-card title="执行结果">
            <a-result v-if="executeResult" :type="executeResult.success ? 'success' : 'error'" :title="executeResult.message">
              <template #subtitle>
                执行时间: {{ executeResult.time }}
              </template>
            </a-result>
            <a-empty v-else description="暂无执行结果" />
          </a-card>

          <a-card title="历史执行记录" style="margin-top: 16px">
            <a-timeline>
              <a-timeline-item v-for="h in history" :key="h.time" :color="h.success ? 'green' : 'red'">
                <p>{{ h.action }}</p>
                <span class="time">{{ h.time }}</span>
              </a-timeline-item>
            </a-timeline>
          </a-card>
        </a-col>
      </a-row>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'

const alert = reactive({
  id: 'A001',
  type: '设备离线',
  deviceId: 'D001',
  deviceName: 'M5Stack-001',
  triggerTime: '2026-03-28 10:00:00',
  level: 'warning',
  message: '设备已离线超过30分钟'
})

const suggestions = ref([
  { id: 1, title: '重启设备', description: '发送重启指令到设备', successRate: 95, history: [] },
  { id: 2, title: '检查网络', description: '验证设备网络连接状态', successRate: 88, history: [] },
  { id: 3, title: '重置MQTT连接', description: '断开并重新建立MQTT连接', successRate: 82, history: [] }
])

const selectedSuggestion = ref(null)
const executeResult = ref(null)
const history = ref([
  { action: '执行: 重启设备', time: '2026-03-28 09:30', success: true },
  { action: '执行: 检查网络', time: '2026-03-28 09:15', success: false }
])

const getLevelColor = (level) => ({ warning: 'orange', major: 'orange', critical: 'red' }[level] || 'gray')

const handleExecute = (suggestion) => {
  selectedSuggestion.value = suggestion
  executeResult.value = { success: true, message: '执行成功', time: new Date().toLocaleString() }
}
const handleViewHistory = (suggestion) => { }
</script>

<style scoped>
.container { padding: 16px; }
.selected { border: 2px solid #409EFF; }
.time { font-size: 12px; color: #909399; }
</style>
