<template>
  <div class="container">
    <a-card class="general-card" title="API文档">
      <a-tabs>
        <a-tab-pane key="list" title="接口列表">
          <a-table :columns="columns" :data="apiList" row-key="path">
            <template #method="{ record }"><a-tag :color="methodColor(record.method)">{{ record.method }}</a-tag></template>
          </a-table>
        </a-tab-pane>
        <a-tab-pane key="detail" title="接口详情">
          <a-empty description="选择一个接口查看详情" />
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import Breadcrumb from '@/components/Breadcrumb.vue'

const apiList = ref([
  { path: '/api/v1/devices', method: 'GET', description: '设备列表' },
  { path: '/api/v1/devices/:id', method: 'GET', description: '设备详情' },
  { path: '/api/v1/devices', method: 'POST', description: '创建设备' },
  { path: '/api/v1/members', method: 'GET', description: '会员列表' },
  { path: '/api/v1/alerts', method: 'GET', description: '告警列表' },
  { path: '/api/v1/ota/packages', method: 'GET', description: 'OTA包列表' },
])

const columns = [
  { title: '方法', slotName: 'method', width: 100 },
  { title: '路径', dataIndex: 'path', width: 240 },
  { title: '描述', dataIndex: 'description' }
]

const methodColor = (m) => ({ GET: 'blue', POST: 'green', PUT: 'orange', DELETE: 'red' }[m] || 'gray')
</script>
