<template>
    <Breadcrumb :items="['Home','Console','']" />


  <div class="container">
    <a-card class="general-card" title="个人桌面">
      <a-row :gutter="16" style="margin-bottom: 16px">
        <a-col :span="8">
          <a-card title="快捷操作" size="small">
            <a-space wrap>
              <a-button v-for="item in shortcuts" :key="item.label" @click="$router.push(item.path)">{{ item.label }}</a-button>
            </a-space>
          </a-card>
        </a-col>
        <a-col :span="8">
          <a-card title="我的任务" size="small">
            <a-statistic title="待完成" :value="myTasks.pending" />
            <a-statistic title="已完成" :value="myTasks.done" />
          </a-card>
        </a-col>
        <a-col :span="8">
          <a-card title="通知" size="small">
            <a-list :data="notifications" size="small" />
          </a-card>
        </a-col>
      </a-row>
    </a-card>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import Breadcrumb from '@/components/Breadcrumb.vue'

const shortcuts = ref([
  { label: '设备列表', path: '/devices' },
  { label: '会员管理', path: '/members' },
  { label: '告警管理', path: '/alerts' }
])
const myTasks = ref({ pending: 0, done: 0 })
const notifications = ref(['暂无新通知'])
</script>
