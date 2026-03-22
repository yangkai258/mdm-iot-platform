<template>
  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>具身智能</a-breadcrumb-item>
      <a-breadcrumb-item>
        <router-link :to="`/embodied/${deviceId}/perception`">设备 {{ deviceId }}</router-link>
      </a-breadcrumb-item>
      <a-breadcrumb-item>决策日志</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 决策上下文 -->
    <a-row :gutter="16" style="margin-top: 16px">
      <a-col :xs="24" :sm="24" :md="12">
        <a-card title="用户状态">
          <a-descriptions :column="2" size="small">
            <a-descriptions-item label="情绪">
              <a-tag :color="getEmotionColor(context.user_state?.emotion)">
                {{ getEmotionText(context.user_state?.emotion) }}
              </a-tag>
            </a-descriptions-item>
            <a-descriptions-item label="姿态">
              <a-tag>{{ context.user_state?.pose || '未知' }}</a-tag>
            </a-descriptions-item>
            <a-descriptions-item label="位置">
              <span v-if="context.user_state?.position">
                ({{ context.user_state.position.x?.toFixed(1) }}, {{ context.user_state.position.y?.toFixed(1) }})
              </span>
              <span v-else>-</span>
            </a-descriptions-item>
          </a-descriptions>
        </a-card>
      </a-col>

      <a-col :xs="24" :sm="24" :md="12">
        <a-card title="环境状态">
          <a-descriptions :column="2" size="small">
            <a-descriptions-item label="当前位置">
              {{ context.environment_state?.location || '未知' }}
            </a-descriptions-item>
            <a-descriptions-item label="时间">
              {{ context.environment_state?.time || '-' }}
            </a-descriptions-item>
            <a-descriptions-item label="场景">
              <a-tag color="arcoblue">{{ context.environment_state?.scene || '未知' }}</a-tag>
            </a-descriptions-item>
          </a-descriptions>
        </a-card>
      </a-col>

      <a-col :xs="24" :sm="24" :md="12" style="margin-top: 16px">
        <a-card title="设备状态">
          <a-descriptions :column="2" size="small">
            <a-descriptions-item label="情绪">
              <a-tag :color="getEmotionColor(context.pet_state?.emotion)">
                {{ getEmotionText(context.pet_state?.emotion) }}
              </a-tag>
            </a-descriptions-item>
            <a-descriptions-item label="电量">
              <a-progress
                v-if="context.pet_state?.battery !== undefined"
                :percent="Math.round(context.pet_state.battery * 100)"
                :color="getBatteryColor(context.pet_state.battery)"
                size="small"
                style="width: 120px"
              />
              <span v-else>-</span>
            </a-descriptions-item>
            <a-descriptions-item label="模式">
              <a-tag color="purple">{{ context.pet_state?.mode || '未知' }}</a-tag>
            </a-descriptions-item>
          </a-descriptions>
        </a-card>
      </a-col>

      <a-col :xs="24" :sm="24" :md="12" style="margin-top: 16px">
        <a-card title="任务进度">
          <a-descriptions :column="1" size="small">
            <a-descriptions-item label="当前任务">
              {{ context.task_state?.current_task || '无任务' }}
            </a-descriptions-item>
            <a-descriptions-item label="进度">
              <a-progress
                v-if="context.task_state?.progress !== undefined"
                :percent="Math.round(context.task_state.progress * 100)"
                size="small"
              />
              <span v-else>-</span>
            </a-descriptions-item>
          </a-descriptions>
          <a-button type="primary" long style="margin-top: 12px" @click="loadContext">
            <template #icon><icon-refresh /></template>
            刷新上下文
          </a-button>
        </a-card>
      </a-col>
    </a-row>

    <!-- 策略配置 -->
    <a-card title="决策策略配置" style="margin-top: 16px">
      <a-space wrap>
        <a-select v-model="activeStrategy" placeholder="选择策略" style="width: 180px" @change="applyStrategy">
          <a-option value="safety_first">安全优先</a-option>
          <a-option value="task_oriented">任务导向</a-option>
          <a-option value="interaction_first">交互优先</a-option>
          <a-option value="exploration">探索模式</a-option>
        </a-select>
        <a-button type="primary" @click="applyStrategy">
          <template #icon><icon-check-circle /></template>
          应用策略
        </a-button>
      </a-space>
      <a-divider />
      <a-row :gutter="12">
        <a-col :xs="24" :sm="12" :md="6">
          <div class="strategy-card" :class="{ active: activeStrategy === 'safety_first' }" @click="activeStrategy = 'safety_first'">
            <icon-warning size="24" style="color: #ff4d4f" />
            <div class="strategy-name">安全优先</div>
            <div class="strategy-desc">始终避开危险区域和障碍物</div>
          </div>
        </a-col>
        <a-col :xs="24" :sm="12" :md="6">
          <div class="strategy-card" :class="{ active: activeStrategy === 'task_oriented' }" @click="activeStrategy = 'task_oriented'">
            <icon-check-circle size="24" style="color: #00b42a" />
            <div class="strategy-name">任务导向</div>
            <div class="strategy-desc">专注于完成任务目标</div>
          </div>
        </a-col>
        <a-col :xs="24" :sm="12" :md="6">
          <div class="strategy-card" :class="{ active: activeStrategy === 'interaction_first' }" @click="activeStrategy = 'interaction_first'">
            <icon-message size="24" style="color: #1650d8" />
            <div class="strategy-name">交互优先</div>
            <div class="strategy-desc">优先响应用户交互</div>
          </div>
        </a-col>
        <a-col :xs="24" :sm="12" :md="6">
          <div class="strategy-card" :class="{ active: activeStrategy === 'exploration' }" @click="activeStrategy = 'exploration'">
            <icon-compass size="24" style="color: #a0a0a0" />
            <div class="strategy-name">探索模式</div>
            <div class="strategy-desc">主动探索新环境区域</div>
          </div>
        </a-col>
      </a-row>
    </a-card>

    <!-- 决策历史 -->
    <a-card title="决策历史" style="margin-top: 16px">
      <template #extra>
        <a-space>
          <a-select v-model="filterType" placeholder="决策类型" allow-clear style="width: 140px" @change="loadLogs">
            <a-option value="navigation">导航决策</a-option>
            <a-option value="interaction">交互决策</a-option>
            <a-option value="safety">安全决策</a-option>
            <a-option value="task">任务决策</a-option>
          </a-select>
          <a-button type="text" size="small" @click="loadLogs">
            <icon-refresh />
          </a-button>
        </a-space>
      </template>
      <a-table
        :columns="logColumns"
        :data="logs"
        :loading="loading"
        :pagination="{ pageSize: 10, showTotal: true, showPageSize: true, pageSizeOptions: [10, 20, 50] }"
        row-key="id"
        size="small"
        @expand="handleExpand"
        :expanded-keys="expandedKeys"
        @update:expanded-keys="expandedKeys = $event"
      >
        <template #decision_type="{ record }">
          <a-tag :color="getDecisionColor(record.decision_type)">
            {{ getDecisionText(record.decision_type) }}
          </a-tag>
        </template>
        <template #confidence="{ record }">
          <a-progress
            :percent="Math.round((record.confidence || 0) * 100)"
            :color="getConfidenceColor(record.confidence)"
            size="small"
            style="width: 100px"
          />
        </template>
        <template #chosen_action="{ record }">
          <a-tag>{{ record.chosen_action || '-' }}</a-tag>
        </template>
        <template #latency_ms="{ record }">
          <span class="text-muted">
            {{ record.latency_ms ? `${record.latency_ms}ms` : '-' }}
          </span>
        </template>
        <template #expanded-row="{ record }">
          <div class="expand-content">
            <a-descriptions :column="2" size="small" bordered>
              <a-descriptions-item label="决策ID">{{ record.id }}</a-descriptions-item>
              <a-descriptions-item label="设备ID">{{ record.device_id }}</a-descriptions-item>
              <a-descriptions-item label="决策时间">
                {{ formatTime(record.decided_at) }}
              </a-descriptions-item>
              <a-descriptions-item label="执行结果">
                <a-tag :color="getResultColor(record.execution_result)" size="small">
                  {{ getResultText(record.execution_result) }}
                </a-tag>
              </a-descriptions-item>
            </a-descriptions>

            <a-divider orientation="left">推理过程</a-divider>
            <div class="reasoning-box">
              <pre class="reasoning-text">{{ record.reasoning || '无推理记录' }}</pre>
            </div>

            <a-divider orientation="left">决策上下文</a-divider>
            <div class="context-box">
              <pre class="context-text">{{ JSON.stringify(record.context || {}, null, 2) }}</pre>
            </div>

            <a-divider orientation="left">动作参数</a-divider>
            <div class="params-box">
              <pre class="params-text">{{ record.action_params ? JSON.stringify(record.action_params, null, 2) : '无参数' }}</pre>
            </div>
          </div>
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { getDecisionContext, setDecisionStrategy, getDecisionLogs } from '@/api/embodied'
import { Message } from '@arco-design/web-vue'

const route = useRoute()
const deviceId = ref(route.params.device_id as string)

const loading = ref(false)
const context = ref<any>({
  user_state: {},
  environment_state: {},
  pet_state: {},
  task_state: {}
})
const logs = ref<any[]>([])
const filterType = ref('')
const activeStrategy = ref('safety_first')
const expandedKeys = ref<number[]>([])

const logColumns = [
  { title: 'ID', dataIndex: 'id', width: 70 },
  { title: '决策类型', dataIndex: 'decision_type', slotName: 'decision_type', width: 120 },
  { title: '执行动作', dataIndex: 'chosen_action', slotName: 'chosen_action', width: 140 },
  { title: '置信度', dataIndex: 'confidence', slotName: 'confidence', width: 130 },
  { title: '延迟', dataIndex: 'latency_ms', slotName: 'latency_ms', width: 90 },
  { title: '决策时间', dataIndex: 'decided_at', slotName: 'decided_at', width: 170 }
]

async function loadContext() {
  try {
    const res = await getDecisionContext(deviceId.value)
    context.value = res.data || {}
  } catch (err: any) {
    Message.error('加载上下文失败: ' + err.message)
  }
}

async function loadLogs() {
  try {
    loading.value = true
    const params: any = { page_size: 50 }
    if (filterType.value) params.decision_type = filterType.value
    const res = await getDecisionLogs(deviceId.value, params)
    logs.value = res.data?.logs || res.data || []
  } catch (err: any) {
    Message.error('加载决策日志失败: ' + err.message)
  } finally {
    loading.value = false
  }
}

async function applyStrategy() {
  try {
    await setDecisionStrategy(deviceId.value, {
      strategy_type: activeStrategy.value,
      is_active: true
    })
    Message.success(`已应用策略: ${getStrategyText(activeStrategy.value)}`)
  } catch (err: any) {
    Message.error('策略应用失败: ' + err.message)
  }
}

function handleExpand(show: boolean, record: any) {
  if (show) {
    if (!expandedKeys.value.includes(record.id)) {
      expandedKeys.value = [...expandedKeys.value, record.id]
    }
  }
}

function getEmotionColor(e: string) {
  const map: Record<string, string> = { happy: 'green', sad: 'blue', angry: 'red', neutral: 'gray', excited: 'orange', calm: 'cyan' }
  return map[e || ''] || 'default'
}
function getEmotionText(e: string) {
  const map: Record<string, string> = { happy: '开心', sad: '悲伤', angry: '生气', neutral: '平静', excited: '兴奋', calm: '冷静' }
  return map[e || ''] || e || '未知'
}
function getBatteryColor(b: number) {
  if (b >= 0.5) return 'green'
  if (b >= 0.2) return 'orange'
  return 'red'
}
function getDecisionColor(t: string) {
  const map: Record<string, string> = { navigation: 'arcoblue', interaction: 'green', safety: 'orange', task: 'purple' }
  return map[t] || 'default'
}
function getDecisionText(t: string) {
  const map: Record<string, string> = { navigation: '导航决策', interaction: '交互决策', safety: '安全决策', task: '任务决策' }
  return map[t] || t
}
function getConfidenceColor(c: number) {
  if (c >= 0.8) return 'green'
  if (c >= 0.5) return 'orange'
  return 'red'
}
function getResultColor(r: string) {
  const map: Record<string, string> = { success: 'green', failed: 'red', partial: 'orange' }
  return map[r || ''] || 'default'
}
function getResultText(r: string) {
  const map: Record<string, string> = { success: '成功', failed: '失败', partial: '部分成功' }
  return map[r || ''] || r || '-'
}
function getStrategyText(s: string) {
  const map: Record<string, string> = { safety_first: '安全优先', task_oriented: '任务导向', interaction_first: '交互优先', exploration: '探索模式' }
  return map[s] || s
}
function formatTime(t: string) {
  if (!t) return '-'
  return new Date(t).toLocaleString('zh-CN')
}

onMounted(async () => {
  await Promise.all([loadContext(), loadLogs()])
})
</script>

<style scoped>
.strategy-card {
  border: 1px solid var(--color-border);
  border-radius: 8px;
  padding: 16px;
  text-align: center;
  cursor: pointer;
  transition: all 0.2s;
  background: var(--color-fill-1);
}
.strategy-card:hover {
  border-color: var(--color-primary);
}
.strategy-card.active {
  border-color: var(--color-primary);
  background: var(--color-primary-light-1);
}
.strategy-name {
  font-size: 14px;
  font-weight: 600;
  margin: 8px 0 4px;
}
.strategy-desc {
  font-size: 12px;
  color: var(--color-text-3);
}
.expand-content {
  padding: 8px 0;
}
.reasoning-box, .context-box, .params-box {
  background: var(--color-fill-1);
  border-radius: 4px;
  padding: 12px;
  margin-bottom: 8px;
}
.reasoning-text, .context-text, .params-text {
  margin: 0;
  font-size: 12px;
  color: var(--color-text-2);
  white-space: pre-wrap;
  word-break: break-all;
}
.text-muted {
  color: var(--color-text-3);
}
</style>
