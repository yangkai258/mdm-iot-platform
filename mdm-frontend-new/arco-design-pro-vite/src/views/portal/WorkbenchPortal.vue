<template>
    <Breadcrumb :items="['Home','Console','']" />


  <div class="container">
    <a-card class="general-card" title="工作台门户">
      <template #extra>
        <a-button @click="loadData"><icon-refresh />刷新</a-button>
      </template>
      <a-row :gutter="16" style="margin-bottom: 16px">
        <a-col :span="6">
          <a-statistic title="设备总数" :value="stats.totalDevices" />
        </a-col>
        <a-col :span="6">
          <a-statistic title="在线设备" :value="stats.onlineDevices" color="green" />
        </a-col>
        <a-col :span="6">
          <a-statistic title="会员总数" :value="stats.totalMembers" />
        </a-col>
        <a-col :span="6">
          <a-statistic title="今日新增" :value="stats.todayNew" color="blue" />
        </a-col>
      </a-row>
      <a-divider style="margin: 0 0 16px 0" />
      <a-row :gutter="16">
        <a-col :span="12">
          <a-card title="快捷入口" size="small">
            <a-space wrap>
              <a-button v-for="shortcut in shortcuts" :key="shortcut.label" @click="$router.push(shortcut.path)">{{ shortcut.label }}</a-button>
            </a-space>
          </a-card>
        </a-col>
        <a-col :span="12">
          <a-card title="最近活动" size="small">
            <a-list :data="recentActivities" size="small">
              <template #item="{ item }">
                <a-list-item>{{ item }}</a-list-item>
              </template>
            </a-list>
          </a-card>
        </a-col>
      </a-row>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import Breadcrumb from '@/components/Breadcrumb.vue'

const stats = reactive({ totalDevices: 0, onlineDevices: 0, totalMembers: 0, todayNew: 0 })
const shortcuts = ref([
  { label: '设备列表', path: '/devices' },
  { label: '会员管理', path: '/members' },
  { label: '告警规则', path: '/alerts/rules' },
  { label: 'OTA升级', path: '/ota/packages' }
])
const recentActivities = ref(['暂无最近活动'])

const loadData = async () => {
  try {
    const res = await fetch('/api/v1/dashboard/stats', {
      headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token') }
    }).then(r => r.json())
    Object.assign(stats, res.data || {})
  } catch {}
}

onMounted(() => loadData())
</script>
