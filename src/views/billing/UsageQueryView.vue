<template>
  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>订阅管理</a-breadcrumb-item>
      <a-breadcrumb-item>用量查询</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 页面标题 -->
    <div class="pro-page-header">
      <h2 class="pro-page-title">当前用量</h2>
    </div>

    <!-- 操作栏 -->
    <div class="pro-action-bar">
      <a-space>
        <a-button @click="loadUsage">
          <template #icon><icon-refresh /></template>
          刷新
        </a-button>
      </a-space>
    </div>

    <!-- 三大指标卡片 -->
    <div class="metric-cards" v-loading="loading">
      <a-row :gutter="16">
        <a-col :span="8">
          <a-card class="usage-card">
            <div class="usage-label">API调用</div>
            <div class="usage-value">{{ usage.api_calls?.used?.toLocaleString() || 0 }}</div>
            <div class="usage-quota">{{ (usage.api_calls?.used || 0).toLocaleString() }}/{{ quotaLabel(usage.api_calls?.quota) }}</div>
            <a-progress
              :percent="usage.api_calls?.percent || 0"
              :color="progressColor(usage.api_calls?.percent)"
              :show-text="true"
            />
          </a-card>
        </a-col>
        <a-col :span="8">
          <a-card class="usage-card">
            <div class="usage-label">设备数</div>
            <div class="usage-value">{{ usage.devices?.used || 0 }}</div>
            <div class="usage-quota">{{ usage.devices?.used || 0 }}/{{ quotaLabel(usage.devices?.quota) }}</div>
            <a-progress
              :percent="usage.devices?.percent || 0"
              :color="progressColor(usage.devices?.percent)"
              :show-text="true"
            />
          </a-card>
        </a-col>
        <a-col :span="8">
          <a-card class="usage-card">
            <div class="usage-label">存储空间</div>
            <div class="usage-value">{{ usage.storage?.used || 0 }}GB</div>
            <div class="usage-quota">{{ usage.storage?.used || 0 }}GB/{{ quotaLabel(usage.storage?.quota) }}GB</div>
            <a-progress
              :percent="usage.storage?.percent || 0"
              :color="progressColor(usage.storage?.percent)"
              :show-text="true"
            />
          </a-card>
        </a-col>
      </a-row>

      <!-- 配额明细表格 -->
      <a-card class="quota-table-card" style="margin-top: 16px;">
        <template #title>配额明细</template>
        <a-table :columns="columns" :data="quotaList" :loading="loading" :pagination="false" />
      </a-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { billingApi } from '@/api/billing'

const loading = ref(false)
const usage = ref<any>({
  api_calls: { used: 12456, quota: 50000, percent: 35 },
  devices: { used: 8, quota: 50, percent: 16 },
  storage: { used: 45, quota: 100, percent: 45 }
})

const columns = [
  { title: '类型', dataIndex: 'type' },
  { title: '已用', dataIndex: 'used' },
  { title: '配额', dataIndex: 'quota' },
  { title: '状态', dataIndex: 'status' }
]

const quotaList = computed(() => [
  {
    type: 'API调用',
    used: (usage.value.api_calls?.used || 0).toLocaleString(),
    quota: quotaLabel(usage.value.api_calls?.quota),
    status: statusTag(usage.value.api_calls?.percent)
  },
  {
    type: '设备数',
    used: usage.value.devices?.used || 0,
    quota: quotaLabel(usage.value.devices?.quota),
    status: statusTag(usage.value.devices?.percent)
  },
  {
    type: '存储空间',
    used: `${usage.value.storage?.used || 0}GB`,
    quota: `${quotaLabel(usage.value.storage?.quota)}GB`,
    status: statusTag(usage.value.storage?.percent)
  }
])

const quotaLabel = (val?: number) => {
  if (val === -1 || val === undefined || val === null) return '不限'
  return val?.toLocaleString()
}

const progressColor = (percent?: number) => {
  if (!percent) return 'rgb(var(--primary-6))'
  if (percent >= 90) return 'rgb(var(--danger-6))'
  if (percent >= 70) return 'rgb(var(--warning-6))'
  return 'rgb(var(--primary-6))'
}

const statusTag = (percent?: number) => {
  if (!percent) return '🟢 正常'
  if (percent >= 90) return '🔴 告警'
  if (percent >= 70) return '🟡 预警'
  return '🟢 正常'
}

const loadUsage = async () => {
  loading.value = true
  try {
    const res = await billingApi.getCurrentUsage()
    if (res.code === 0 || res.code === 200) {
      usage.value = res.data || usage.value
    }
  } catch (e) {
    // use mock
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadUsage()
})
</script>

<style scoped>
.metric-cards {
  padding: 8px 0;
}
.usage-card {
  text-align: center;
}
.usage-label {
  font-size: 14px;
  color: var(--color-text-3, #86909c);
  margin-bottom: 8px;
}
.usage-value {
  font-size: 28px;
  font-weight: 700;
  color: var(--color-text-1, #1f2329);
}
.usage-quota {
  font-size: 12px;
  color: var(--color-text-3, #86909c);
  margin-bottom: 12px;
}
</style>
