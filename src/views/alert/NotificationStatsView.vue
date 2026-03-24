<template>
  <div class="notification-stats-view">
    <!-- 时间范围选择 -->
    <div class="filter-section">
      <a-form layout="inline">
        <a-form-item label="时间范围">
          <a-range-picker
            v-model="dateRange"
            style="width: 280px"
            @change="onDateChange"
          />
        </a-form-item>
        <a-form-item>
          <a-space>
            <a-button type="primary" @click="handleSearch">
              <template #icon><icon-search /></template>
              「查询」
            </a-button>
            <a-button @click="handleRefresh">
              <template #icon><icon-refresh /></template>
              「刷新」
            </a-button>
          </a-space>
        </a-form-item>
      </a-form>
    </div>

    <!-- 统计卡片 -->
    <div class="stats-cards">
      <a-row :gutter="16">
        <a-col :span="6">
          <a-card class="stat-card">
            <div class="stat-content">
              <div class="stat-icon sent">
                <icon-send />
              </div>
              <div class="stat-info">
                <div class="stat-value">{{ stats?.total_sent || 0 }}</div>
                <div class="stat-label">总发送量</div>
              </div>
            </div>
          </a-card>
        </a-col>
        <a-col :span="6">
          <a-card class="stat-card">
            <div class="stat-content">
              <div class="stat-icon failed">
                <icon-close />
              </div>
              <div class="stat-info">
                <div class="stat-value">{{ stats?.total_failed || 0 }}</div>
                <div class="stat-label">失败数量</div>
              </div>
            </div>
          </a-card>
        </a-col>
        <a-col :span="6">
          <a-card class="stat-card">
            <div class="stat-content">
              <div class="stat-icon rate">
                <icon-check-circle />
              </div>
              <div class="stat-info">
                <div class="stat-value">{{ stats?.success_rate || 0 }}%</div>
                <div class="stat-label">成功率</div>
              </div>
            </div>
          </a-card>
        </a-col>
        <a-col :span="6">
          <a-card class="stat-card">
            <div class="stat-content">
              <div class="stat-icon channels">
                <icon-all-application />
              </div>
              <div class="stat-info">
                <div class="stat-value">{{ stats?.by_channel?.length || 0 }}</div>
                <div class="stat-label">活跃渠道</div>
              </div>
            </div>
          </a-card>
        </a-col>
      </a-row>
    </div>

    <!-- 渠道统计 -->
    <div class="channel-stats">
      <div class="section-title">各渠道发送统计</div>
      <a-table
        :data="stats?.by_channel || []"
        :loading="loading"
        :columns="channelColumns"
        :pagination="false"
        stripe
      >
        <template #columns>
          <a-table-column title="渠道类型" data-index="channel_type">
            <template #cell="{ record }">
              <a-tag :color="channelTypeColor(record.channel_type)" size="small">
                {{ channelTypeLabel(record.channel_type) }}
              </a-tag>
            </template>
          </a-table-column>
          <a-table-column title="发送成功" data-index="sent">
            <template #cell="{ record }">
              <span class="text-success">{{ record.sent }}</span>
            </template>
          </a-table-column>
          <a-table-column title="发送失败" data-index="failed">
            <template #cell="{ record }">
              <span class="text-danger">{{ record.failed }}</span>
            </template>
          </a-table-column>
          <a-table-column title="成功率" data-index="success_rate">
            <template #cell="{ record }">
              <a-progress
                :percent="record.success_rate"
                :color="record.success_rate >= 90 ? 'green' : record.success_rate >= 70 ? 'orange' : 'red'"
                size="small"
              />
            </template>
          </a-table-column>
        </template>
      </a-table>
    </div>

    <!-- 趋势图表 -->
    <div class="trend-chart">
      <div class="section-title">发送趋势</div>
      <div class="chart-placeholder">
        <a-empty v-if="!stats?.by_day?.length" description="暂无数据" />
        <div v-else class="chart-content">
          <!-- 简化图表展示 -->
          <div class="simple-chart">
            <div
              v-for="(day, index) in stats?.by_day || []"
              :key="index"
              class="chart-bar"
            >
              <div
                class="bar-sent"
                :style="{ height: getBarHeight(day.sent) + 'px' }"
                :title="'发送: ' + day.sent"
              ></div>
              <div
                class="bar-failed"
                :style="{ height: getBarHeight(day.failed) + 'px' }"
                :title="'失败: ' + day.failed"
              ></div>
            </div>
          </div>
          <div class="chart-labels">
            <span
              v-for="(day, index) in (stats?.by_day || []).slice(-7)"
              :key="index"
              class="chart-label"
            >
              {{ formatDate(day.date) }}
            </span>
          </div>
          <div class="chart-legend">
            <span class="legend-item"><span class="dot sent"></span>发送</span>
            <span class="legend-item"><span class="dot failed"></span>失败</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import {
  IconSearch,
  IconRefresh,
  IconSend,
  IconClose,
  IconCheckCircle,
  IconApps
} from '@arco-design/web-vue/es/icon'
import { useNotificationStats } from '@/composables/useNotification'
import { CHANNEL_TYPE_MAP } from '@/composables/useNotification'

const {
  loading,
  stats,
  dateRange,
  loadStats
} = useNotificationStats()

const channelColumns = [
  { title: '渠道类型', dataIndex: 'channel_type' },
  { title: '发送成功', dataIndex: 'sent' },
  { title: '发送失败', dataIndex: 'failed' },
  { title: '成功率', dataIndex: 'success_rate' }
]

function channelTypeLabel(type) {
  return CHANNEL_TYPE_MAP[type]?.label || type
}

function channelTypeColor(type) {
  const colorMap = { email: 'blue', sms: 'green', webhook: 'orange' }
  return colorMap[type] || 'gray'
}

function formatDate(date) {
  if (!date) return ''
  return new Date(date).toLocaleDateString('zh-CN', { month: '2-digit', day: '2-digit' })
}

function getBarHeight(value) {
  const maxSent = Math.max(...(stats.value?.by_day?.map(d => d.sent) || [1]))
  const maxFailed = Math.max(...(stats.value?.by_day?.map(d => d.failed) || [1]))
  const max = Math.max(maxSent, maxFailed, 1)
  return Math.max(5, (value / max) * 100)
}

function onDateChange(value) {
  // dateRange already updated via v-model
}

function handleSearch() {
  loadStats()
}

function handleRefresh() {
  loadStats()
  Message.success('已刷新')
}

onMounted(() => {
  loadStats()
})
</script>

<style scoped>
.notification-stats-view {
  padding: 20px;
}
.filter-section {
  background: #fff;
  padding: 16px;
  border-radius: 8px 8px 0 0;
}
.stats-cards {
  background: #fff;
  padding: 16px;
  border-bottom: 1px solid #e5e6eb;
}
.stat-card {
  border-radius: 8px;
}
.stat-card :deep(.arco-card-body) {
  padding: 16px;
}
.stat-content {
  display: flex;
  align-items: center;
  gap: 16px;
}
.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
}
.stat-icon.sent {
  background: #e6f4ff;
  color: #1650d8;
}
.stat-icon.failed {
  background: #fff2e8;
  color: #f53f3f;
}
.stat-icon.rate {
  background: #e6fffb;
  color: #00b42a;
}
.stat-icon.channels {
  background: #f9f0ff;
  color: #722ed1;
}
.stat-value {
  font-size: 24px;
  font-weight: 600;
  color: #1d2129;
}
.stat-label {
  font-size: 14px;
  color: #86909c;
}
.channel-stats {
  background: #fff;
  padding: 16px;
  margin-top: 0;
}
.section-title {
  font-size: 16px;
  font-weight: 600;
  color: #1d2129;
  margin-bottom: 16px;
}
.trend-chart {
  background: #fff;
  padding: 16px;
  border-radius: 0 0 8px 8px;
  margin-top: 16px;
}
.chart-placeholder {
  min-height: 250px;
}
.chart-content {
  padding: 16px 0;
}
.simple-chart {
  display: flex;
  align-items: flex-end;
  justify-content: space-around;
  height: 180px;
  padding: 0 16px;
  border-bottom: 1px solid #e5e6eb;
}
.chart-bar {
  display: flex;
  gap: 4px;
  align-items: flex-end;
}
.bar-sent {
  width: 20px;
  background: #1650d8;
  border-radius: 2px 2px 0 0;
  min-height: 5px;
  transition: height 0.3s;
}
.bar-failed {
  width: 20px;
  background: #f53f3f;
  border-radius: 2px 2px 0 0;
  min-height: 5px;
  transition: height 0.3s;
}
.chart-labels {
  display: flex;
  justify-content: space-around;
  padding: 8px 16px;
}
.chart-label {
  font-size: 12px;
  color: #86909c;
  width: 40px;
  text-align: center;
}
.chart-legend {
  display: flex;
  justify-content: center;
  gap: 24px;
  padding-top: 12px;
}
.legend-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: #86909c;
}
.legend-item .dot {
  width: 12px;
  height: 12px;
  border-radius: 2px;
}
.legend-item .dot.sent {
  background: #1650d8;
}
.legend-item .dot.failed {
  background: #f53f3f;
}
.text-success {
  color: #00b42a;
  font-weight: 600;
}
.text-danger {
  color: #f53f3f;
  font-weight: 600;
}
</style>