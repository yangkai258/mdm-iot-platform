<template>
  <div class="system-info">
    <Breadcrumb :items="['menu.tenant', 'menu.tenant.systemInfo']" />
    <a-card class="general-card" title="系统信息">
      <a-descriptions :column="2" bordered size="large">
        <a-descriptions-item label="系统名称">MDM 设备管理平台</a-descriptions-item>
        <a-descriptions-item label="系统版本">v2.4.1</a-descriptions-item>
        <a-descriptions-item label="Go 版本">go1.22</a-descriptions-item>
        <a-descriptions-item label="Vue 版本">3.4.x</a-descriptions-item>
        <a-descriptions-item label="数据库">PostgreSQL 16</a-descriptions-item>
        <a-descriptions-item label="MQTT Broker">EMQX 5.0</a-descriptions-item>
        <a-descriptions-item label="Redis 版本">7.2</a-descriptions-item>
        <a-descriptions-item label="部署环境">Docker</a-descriptions-item>
        <a-descriptions-item label="运行时长" :span="2">{{ uptime }}</a-descriptions-item>
        <a-descriptions-item label="系统状态" :span="2">
          <a-tag color="green">运行正常</a-tag>
        </a-descriptions-item>
      </a-descriptions>

      <a-divider>组件状态</a-divider>
      <a-row :gutter="[16, 16]">
        <a-col v-for="comp in components" :key="comp.name" :span="8">
          <a-card size="small">
            <a-space direction="vertical" fill>
              <a-typography-text bold>{{ comp.name }}</a-typography-text>
              <a-tag :color="comp.status === 'up' ? 'green' : 'red'">
                {{ comp.status === 'up' ? '在线' : '离线' }}
              </a-tag>
              <a-typography-text type="secondary">{{ comp.version }}</a-typography-text>
            </a-space>
          </a-card>
        </a-col>
      </a-row>
    </a-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'

const uptime = ref('计算中...')
const components = ref([
  { name: 'API Server', status: 'up', version: 'v1.0.0' },
  { name: 'MQTT Gateway', status: 'up', version: 'EMQX 5.0' },
  { name: 'Redis Cache', status: 'up', version: '7.2' },
  { name: 'PostgreSQL', status: 'up', version: '16.0' },
  { name: 'OTA Service', status: 'up', version: 'v1.2.0' },
  { name: 'Notification Service', status: 'up', version: 'v1.0.5' },
])

onMounted(() => {
  const start = new Date('2026-03-15T00:00:00Z')
  const diff = Date.now() - start.getTime()
  const days = Math.floor(diff / 86400000)
  const hours = Math.floor((diff % 86400000) / 3600000)
  uptime.value = `${days} 天 ${hours} 小时`
})
</script>

<style scoped>
.system-info { padding: 16px; }
</style>
