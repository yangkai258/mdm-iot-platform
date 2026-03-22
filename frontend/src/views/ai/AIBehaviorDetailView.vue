<template>
  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>AI 管理</a-breadcrumb-item>
      <a-breadcrumb-item>
        <a @click="goBack">AI 行为日志</a>
      </a-breadcrumb-item>
      <a-breadcrumb-item>行为详情</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 加载状态 -->
    <a-spin v-if="loading" />
    <div v-else-if="record" class="detail-content">

      <!-- 基本信息卡片 -->
      <a-card class="detail-card" title="基本信息">
        <a-descriptions :column="3" bordered size="small">
          <a-descriptions-item label="行为ID">{{ record.id }}</a-descriptions-item>
          <a-descriptions-item label="设备ID">{{ record.device_id || '--' }}</a-descriptions-item>
          <a-descriptions-item label="用户ID">{{ record.user_id || '--' }}</a-descriptions-item>
          <a-descriptions-item label="模型版本">{{ record.model_version || '--' }}</a-descriptions-item>
          <a-descriptions-item label="事件类型">
            <a-tag :color="getBehaviorTypeColor(record.behavior_type)">
              {{ getBehaviorTypeText(record.behavior_type) }}
            </a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="时间">
            {{ formatTime(record.created_at) }}
          </a-descriptions-item>
        </a-descriptions>
      </a-card>

      <!-- 性能指标卡片 -->
      <a-card class="detail-card" title="性能指标">
        <a-row :gutter="16">
          <a-col :span="6">
            <a-statistic title="推理延迟" :value="record.latency_ms" suffix="ms" />
          </a-col>
          <a-col :span="6">
            <a-statistic title="置信度" :value="record.confidence ? (record.confidence * 100).toFixed(1) : '--'" suffix="%" />
          </a-col>
          <a-col :span="6">
            <a-statistic title="异常分数" :value="record.anomaly_score ? (record.anomaly_score * 100).toFixed(1) : '--'" suffix="%" />
          </a-col>
          <a-col :span="6">
            <a-statistic title="状态">
              <template #extra>
                <a-tag v-if="record.is_anomaly" color="red">异常</a-tag>
                <a-tag v-else-if="record.error_code" color="orange">错误</a-tag>
                <a-tag v-else color="green">正常</a-tag>
              </template>
            </a-statistic>
          </a-col>
        </a-row>

        <!-- 异常标记 -->
        <a-alert
          v-if="record.is_anomaly"
          type="error"
          title="检测到异常行为"
          style="margin-top: 16px"
        >
          <template #icon><icon-warning /></template>
          <div>异常分数: {{ (record.anomaly_score * 100).toFixed(2) }}%</div>
          <div v-if="record.error_message">错误信息: {{ record.error_message }}</div>
        </a-alert>

        <!-- 错误信息 -->
        <a-alert
          v-if="record.error_code"
          type="warning"
          :title="`错误代码: ${record.error_code}`"
          style="margin-top: 16px"
        >
          <div>{{ record.error_message }}</div>
        </a-alert>
      </a-card>

      <!-- 输入数据 -->
      <a-card class="detail-card" title="输入数据">
        <div class="data-block">
          <pre>{{ formatJson(record.input_summary) }}</pre>
        </div>
      </a-card>

      <!-- 输出数据 -->
      <a-card class="detail-card" title="输出数据">
        <div class="data-block">
          <pre>{{ formatJson(record.output_summary) }}</pre>
        </div>
      </a-card>

    </div>
    <a-empty v-else description="未找到该行为记录" />

    <!-- 返回按钮 -->
    <div class="bottom-action">
      <a-button @click="goBack">
        <template #icon><icon-left /></template>
        返回列表
      </a-button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Message } from '@arco-design/web-vue'
import { getAiMonitorEventById } from '@/api/ai'

const route = useRoute()
const router = useRouter()

const loading = ref(false)
const record = ref(null)

const getBehaviorTypeColor = (type) => ({
  intent_recognition: 'arcoblue',
  response_generation: 'green',
  action_selection: 'purple'
}[type] || 'gray')

const getBehaviorTypeText = (type) => ({
  intent_recognition: '意图识别',
  response_generation: '响应生成',
  action_selection: '动作选择'
}[type] || type)

const formatTime = (ts) => {
  if (!ts) return '--'
  return new Date(ts).toLocaleString('zh-CN', { hour12: false })
}

const formatJson = (obj) => {
  if (!obj) return '(无数据)'
  if (typeof obj === 'string') {
    try { return JSON.stringify(JSON.parse(obj), null, 2) } catch { return obj }
  }
  return JSON.stringify(obj, null, 2)
}

const goBack = () => {
  router.push('/ai/behavior-log')
}

const loadDetail = async () => {
  loading.value = true
  try {
    const res = await getAiMonitorEventById(route.params.id)
    if (res.code === 0) {
      record.value = res.data
    } else {
      Message.error('加载失败')
    }
  } catch (e) {
    Message.error('加载失败')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadDetail()
})
</script>

<style scoped>
.pro-page-container { padding: 20px 24px; min-height: calc(100vh - 64px); background: #f5f7fa; }
.pro-breadcrumb { margin-bottom: 16px; }

.detail-card {
  margin-bottom: 16px;
  border-radius: 8px;
}

.data-block {
  background: #f5f7fa;
  border-radius: 4px;
  padding: 12px 16px;
  max-height: 300px;
  overflow: auto;
}

.data-block pre {
  margin: 0;
  font-family: 'Courier New', monospace;
  font-size: 13px;
  white-space: pre-wrap;
  word-break: break-all;
}

.bottom-action {
  display: flex;
  justify-content: flex-start;
  margin-top: 8px;
}
</style>
