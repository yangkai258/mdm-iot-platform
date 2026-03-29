<template>
  <div class="page-container">
    <a-card class="general-card" title="导航记录">
      <template #extra>
        <a-button @click="loadData"><icon-refresh />刷新</a-button>
      </template>
      <div class="search-form">
        <a-form :model="form" layout="inline">
          <a-form-item label="状态">
            <a-select v-model="form.status" placeholder="选择状态" allow-clear style="width: 120px">
              <a-option value="navigating">导航中</a-option>
              <a-option value="arrived">已到达</a-option>
              <a-option value="failed">失败</a-option>
            </a-select>
          </a-form-item>
          <a-form-item>
            <a-button type="primary" @click="handleSearch">查询</a-button>
            <a-button @click="handleReset">重置</a-button>
          </a-form-item>
        </a-form>
      </div>
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" />
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { getNavigationHistory } from '@/api/embodied'
import { Message } from '@arco-design/web-vue'
import { IconRefresh } from '@arco-design/web-vue/es/icon'

const route = useRoute()
const deviceId = ref(route.params.device_id as string || '')

const loading = ref(false)
const data = ref<any[]>([])
const form = ref<any>({ status: '' })

const columns = [
  { title: 'ID', dataIndex: 'id', width: 70 },
  { title: '状态', dataIndex: 'status', width: 100 },
  { title: '起点', dataIndex: 'start_pos', width: 120 },
  { title: '终点', dataIndex: 'target_pos', width: 120 },
  { title: '距离', dataIndex: 'distance', width: 80 },
  { title: '耗时', dataIndex: 'duration', width: 80 },
  { title: '开始时间', dataIndex: 'started_at', width: 170 }
]

const pagination = ref({ current: 1, pageSize: 20, total: 0, showTotal: true, showPageSize: true })

async function loadData() {
  try {
    loading.value = true
    const params: any = { page: pagination.value.current, page_size: pagination.value.pageSize }
    const res = await getNavigationHistory(deviceId.value, params)
    data.value = (res.data?.history || res.data || []).map((h: any) => ({
      ...h,
      start_pos: h.start ? `(${h.start.x?.toFixed(1)}, ${h.start.y?.toFixed(1)})` : '-',
      target_pos: h.target ? `(${h.target.x?.toFixed(1)}, ${h.target.y?.toFixed(1)})` : '-',
      distance: h.distance ? `${h.distance.toFixed(2)}m` : '-',
      duration: h.duration ? `${h.duration}s` : '-'
    }))
    pagination.value.total = data.value.length
  } catch (err: any) {
    Message.error('加载失败: ' + err.message)
  } finally {
    loading.value = false
  }
}

function handleSearch() {
  pagination.value.current = 1
  loadData()
}

function handleReset() {
  form.value = { status: '' }
  handleSearch()
}

function handleRefresh() {
  loadData()
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.page-container { padding: 16px; }
.search-form { margin-bottom: 16px; padding: 16px; background: var(--color-fill-lightest); border-radius: 4px; }
</style>
