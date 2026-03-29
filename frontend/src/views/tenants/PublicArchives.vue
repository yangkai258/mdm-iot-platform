<template>
  <div class="public-archives">
    <a-card class="general-card" title="公共档案">
      <a-alert type="info">
        公共档案用于管理跨租户共享的基础数据，如设备型号库、协议模板等。
      </a-alert>
      <a-divider />
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="{ pageSize: 10 }" row-key="id">
        <template #type="{ record }">
          <a-tag>{{ record.type }}</a-tag>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="view(record)">查看</a-button>
            <a-button type="text" size="small" @click="edit(record)">编辑</a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import axios from 'axios'
import { Message } from '@arco-design/web-vue'

const getToken = () => localStorage.getItem('token')

const loading = ref(false)
const columns = [
  { title: '档案编号', dataIndex: 'id' },
  { title: '档案名称', dataIndex: 'name' },
  { title: '类型', slotName: 'type' },
  { title: '创建时间', dataIndex: 'created_at' },
  { title: '操作', slotName: 'actions', width: 160 },
]
const data = ref([
  { id: 'ARCH-001', name: 'M5Stack Core2 设备型号', type: '设备型号', created_at: '2026-01-15' },
  { id: 'ARCH-002', name: 'MQTT协议配置模板', type: '协议模板', created_at: '2026-02-01' },
  { id: 'ARCH-003', name: 'OTA升级包命名规范', type: '规范文档', created_at: '2026-02-10' },
])

const loadData = async () => {
  try {
    const res = await axios.get('/api/v1/public-archives', { headers: { Authorization: `Bearer ${getToken()}` } })
    if (res.data.code === 0) {
      data.value = res.data.data?.list || []
    }
  } catch { /* use mock */ }
}

const view = (record) => { Message.info(`查看 ${record.name}`) }
const edit = async (record) => {
  try {
    const res = await axios.get(`/api/v1/public-archives/${record.id}`, { headers: { Authorization: `Bearer ${getToken()}` } })
    if (res.data.code === 0) {
      Message.success('已加载详情')
    }
  } catch { Message.info(`编辑 ${record.name}`) }
}
</script>

<style scoped>
.public-archives { padding: 16px; }
</style>
