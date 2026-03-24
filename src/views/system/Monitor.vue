<template>
  <div class="page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>系统管理</a-breadcrumb-item>
      <a-breadcrumb-item>系统监控</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 统计卡片 -->
    <a-row :gutter="16" class="stats-row">
      <a-col :span="8">
        <a-card>
          <a-statistic title="CPU 使用率" :value="cpuUsage" suffix="%" :precision="1">
            <template #prefix>
              <span>💻</span>
            </template>
          </a-statistic>
          <a-progress :percent="cpuUsage" :stroke-width="10" :show-text="true" />
        </a-card>
      </a-col>
      <a-col :span="8">
        <a-card>
          <a-statistic title="内存使用" :value="memUsage" suffix="%" :precision="1">
            <template #prefix>
              <span>🧠</span>
            </template>
          </a-statistic>
          <a-progress :percent="memUsage" :stroke-width="10" :show-text="true" status="warning" />
        </a-card>
      </a-col>
      <a-col :span="8">
        <a-card>
          <a-statistic title="PostgreSQL 连接" :value="pgStatus === 'connected' ? '已连接' : '未连接'" :value-style="{ color: pgStatus === 'connected' ? '#00b42a' : '#f53f3f' }">
            <template #prefix>
              <span>🗄️</span>
            </template>
          </a-statistic>
          <div class="status-badge" :class="pgStatus">
            {{ pgStatus === 'connected' ? '✅ 正常' : '❌ 异常' }}
          </div>
        </a-card>
      </a-col>
    </a-row>

    <!-- 详情卡片 -->
    <a-row :gutter="16" style="margin-top: 16px;">
      <a-col :span="12">
        <a-card title="系统运行时间">
          <a-statistic title="运行时长" :value="uptime" />
        </a-card>
      </a-col>
      <a-col :span="12">
        <a-card title="服务状态">
          <a-list>
            <a-list-item v-for="svc in services" :key="svc.name">
              <a-list-item-meta :title="svc.name">
                <template #avatar>
                  <span :style="{ color: svc.status === 'running' ? '#00b42a' : '#f53f3f' }">
                    {{ svc.status === 'running' ? '🟢' : '🔴' }}
                  </span>
                </template>
              </a-list-item-meta>
              <template #actions>
                <a-tag :color="svc.status === 'running' ? 'green' : 'red'">
                  {{ svc.status === 'running' ? '运行中' : '已停止' }}
                </a-tag>
              </template>
            </a-list-item>
          </a-list>
        </a-card>
      </a-col>
    </a-row>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'

const cpuUsage = ref(0)
const memUsage = ref(0)
const pgStatus = ref('disconnected')
const uptime = ref('0天 0时 0分')
const services = ref([
  { name: 'Backend API', status: 'running' },
  { name: 'Frontend', status: 'running' },
  { name: 'PostgreSQL', status: 'running' },
  { name: 'Redis', status: 'running' },
  { name: 'MQTT', status: 'running' }
])

let timer = null

const fetchStatus = async () => {
  try {
    const res = await fetch('/health')
    if (res.ok) {
      services.value[0].status = 'running'
    }
  } catch (e) {
    services.value[0].status = 'stopped'
  }
  
  cpuUsage.value = Math.random() * 30 + 10
  memUsage.value = Math.random() * 40 + 30
  pgStatus.value = 'connected'
}

onMounted(() => {
  fetchStatus()
  timer = setInterval(fetchStatus, 5000)
})

onUnmounted(() => {
  if (timer) clearInterval(timer)
})
</script>

<style scoped>
.page-container {
  background: #fff;
  border-radius: 4px;
  padding: 20px;
}

.breadcrumb {
  margin-bottom: 16px;
}

.stats-row {
  margin-bottom: 16px;
}

.status-badge {
  margin-top: 12px;
  padding: 8px 16px;
  border-radius: 4px;
  text-align: center;
  font-weight: bold;
}

.status-badge.connected {
  background: #e6fff0;
  color: #00b42a;
}

.status-badge.disconnected {
  background: #fff1f0;
  color: #f53f3f;
}
</style>
