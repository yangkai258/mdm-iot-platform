<template>
  <Breadcrumb :items="['Home','Console','']" />
  <div class="page-container">
    <a-card class="general-card" title="感知日志">
      <template #extra>
        <a-button @click="handleRefresh"><icon-refresh />刷新</a-button>
      </template>
      <div class="search-form">
        <a-form :model="form" layout="inline">
          <a-form-item label="事件类型">
            <a-select v-model="form.event_type" placeholder="选择类型" allow-clear style="width: 160px">
              <a-option value="object_detected">物体检测</a-option>
              <a-option value="scene_changed">场景变化</a-option>
              <a-option value="human_detected">人体检测</a-option>
              <a-option value="obstacle_detected">障碍物检测</a-option>
              <a-option value="touch_detected">触摸检测</a-option>
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
import { getPerception } from '@/api/embodied'
import { Message } from '@arco-design/web-vue'
import { IconRefresh } from '@arco-design/web-vue/es/icon'

const route = useRoute()
const deviceId = ref(route.params.device_id as string || '')

const loading = ref(false)
const data = ref<any[]>([])
const form = ref<any>({ event_type: '' })

const columns = [
  { title: 'ID', dataIndex: 'id', width: 80 },
  { title: '事件类型', dataIndex: 'event_type', width: 140 },
  { title: '描述', dataIndex: 'description', ellipsis: true },
  { title: '时间', dataIndex: 'created_at', width: 180 }
]

const pagination = ref({ current: 1, pageSize: 20, total: 0, showTotal: true, showPageSize: true })

async function loadData() {
  try {
    loading.value = true
    const params: any = { page: pagination.value.current, page_size: pagination.value.pageSize }
    const res = await getPerception(deviceId.value, params)
    data.value = res.data?.events || []
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
  form.value = { event_type: '' }
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
