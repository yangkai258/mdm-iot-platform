<template>
  <div class="page-container">
    <a-card class="general-card" title="寻回网络">
      <template #extra>
        <a-button type="primary"><icon-plus />登记丢失</a-button>
      </template>
      <a-tabs v-model:active-tab="activeTab" @change="loadData">
        <a-tab-pane key="online" tab="在线定位" />
        <a-tab-pane key="lost" tab="丢失登记" />
        <a-tab-pane key="found" tab="已寻�? />
      </a-tabs>
      <div class="search-form">
        <a-form :model="form" layout="inline">
          <a-form-item label="设备ID">
            <a-input v-model="form.deviceId" placeholder="请输入设备ID" style="width: 160px" />
          </a-form-item>
          <a-form-item>
            <a-button type="primary" @click="loadData">查询</a-button>
            <a-button @click="handleReset">重置</a-button>
          </a-form-item>
        </a-form>
      </div>
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id">
        <template #status="{ record }">
          <a-tag :color="getStatusColor(record.status)">{{ getStatusText(record.status) }}</a-tag>
        </template>
        <template #actions="{ record }">
          <a-button type="primary" size="small" @click="handleLocate(record)">定位</a-button>
          <a-button v-if="record.status === 'lost'" type="primary" status="success" size="small" @click="handleFound(record)">已寻�?/a-button>
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import { IconPlus } from '@arco-design/web-vue/es/icon'

const loading = ref(false)
const activeTab = ref('online')
const form = reactive({ deviceId: '' })
const data = ref([])
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })

const columns = [
  { title: '设备ID', dataIndex: 'device_id', width: 140 },
  { title: '宠物名称', dataIndex: 'pet_name', width: 120 },
  { title: '最后位�?, dataIndex: 'last_location', width: 200 },
  { title: '最后在�?, dataIndex: 'last_online', width: 160 },
  { title: '电池', dataIndex: 'battery', width: 80 },
  { title: '状�?, slotName: 'status', width: 100 },
  { title: '操作', slotName: 'actions', width: 160 }
]

const getStatusColor = (status) => {
  const colors = { online: 'green', offline: 'gray', lost: 'red', found: 'blue' }
  return colors[status] || 'gray'
}

const getStatusText = (status) => {
  const texts = { online: '在线', offline: '离线', lost: '丢失�?, found: '已寻�? }
  return texts[status] || status
}

const loadData = async () => {
  loading.value = true
  try {
    const res = await fetch(`/api/v1/device/locate?status=${activeTab.value}`).then(r => r.json())
    if (res.code === 0) {
      data.value = res.data?.list || []
      pagination.total = res.data?.total || 0
    } else {
      loadMockData()
    }
  } catch { loadMockData() } finally { loading.value = false }
}

const loadMockData = () => {
  const mockData = {
    online: [
      { id: 1, device_id: 'DEV001', pet_name: '豆豆', last_location: '北京市朝阳区', last_online: '2026-04-09 19:30:00', battery: '85%', status: 'online' },
      { id: 2, device_id: 'DEV002', pet_name: '旺财', last_location: '上海市浦东新�?, last_online: '2026-04-09 19:25:00', battery: '72%', status: 'online' }
    ],
    lost: [
      { id: 3, device_id: 'DEV003', pet_name: '小白', last_location: '广州市天河区', last_online: '2026-04-08 14:20:00', battery: '5%', status: 'lost' }
    ],
    found: [
      { id: 4, device_id: 'DEV004', pet_name: '花花', last_location: '深圳市南山区', last_online: '2026-04-05 10:00:00', battery: '0%', status: 'found' }
    ]
  }
  data.value = mockData[activeTab.value] || []
  pagination.total = data.value.length
}

const handleReset = () => {
  form.deviceId = ''
  loadData()
}

const handleLocate = (record) => {
  Message.success(`正在定位: ${record.device_id}`)
}

const handleFound = (record) => {
  record.status = 'found'
  Message.success(`已更新状�? ${record.pet_name} 已寻回`)
}

const onPageChange = (page) => {
  pagination.current = page
  loadData()
}

onMounted(() => loadData())
</script>

<style scoped>
.page-container { padding: 16px; }
.search-form { margin-bottom: 16px; padding: 16px; background: var(--color-fill-lightest); border-radius: 4px; }
</style>