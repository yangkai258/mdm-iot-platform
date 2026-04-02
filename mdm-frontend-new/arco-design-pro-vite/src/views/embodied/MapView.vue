<template>
    <Breadcrumb :items="['Home','Console','']" />


  <div class="page-container">
    <a-card class="general-card" title="地图管理">
      <template #extra>
        <a-button type="primary" @click="handleCreate"><icon-plus />新建</a-button>
      </template>
      <div class="search-form">
        <a-form :model="form" layout="inline">
          <a-form-item label="地图名称"><a-input v-model="form.map_name" placeholder="请输入地图名称" /></a-form-item>
          <a-form-item label="地图类型">
            <a-select v-model="form.map_type" placeholder="选择类型" allow-clear style="width: 140px">
              <a-option value="grid">栅格地图</a-option>
              <a-option value="semantic">语义地图</a-option>
              <a-option value="topological">拓扑地图</a-option>
            </a-select>
          </a-form-item>
          <a-form-item>
            <a-button type="primary" @click="handleSearch">查询</a-button>
            <a-button @click="handleReset">重置</a-button>
          </a-form-item>
        </a-form>
      </div>
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" />
      </a-table>
    <a-modal v-model:visible="modalVisible" :title="modalTitle" @before-ok="handleSubmit">
      <a-form :model="form" label-col-flex="100px">
        <a-form-item label="地图名称" required><a-input v-model="form.map_name" placeholder="请输入地图名称" /></a-form-item>
        <a-form-item label="地图类型">
          <a-select v-model="form.map_type" placeholder="选择类型">
            <a-option value="grid">栅格地图</a-option>
            <a-option value="semantic">语义地图</a-option>
            <a-option value="topological">拓扑地图</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="分辨率 (m)">
          <a-input-number v-model="form.resolution" :min="0.01" :max="1" :step="0.01" style="width: 100%" />
        </a-form-item>
        <a-form-item label="描述">
          <a-textarea v-model="form.description" :rows="3" placeholder="地图描述" />
        </a-form-item>
      </a-form>
    </a-modal>
    </a-card>`n</div></template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { getMaps, updateMap } from '@/api/embodied'
import { Message } from '@arco-design/web-vue'
import { IconPlus } from '@arco-design/web-vue/es/icon'

const route = useRoute()
const deviceId = ref(route.params.device_id as string || '')

const loading = ref(false)
const data = ref<any[]>([])
const modalVisible = ref(false)
const modalTitle = ref('新建地图')
const form = ref<any>({
  map_name: '',
  map_type: 'grid',
  resolution: 0.05,
  description: ''
})

const columns = [
  { title: 'ID', dataIndex: 'id', width: 70 },
  { title: '地图名称', dataIndex: 'map_name', width: 160 },
  { title: '类型', dataIndex: 'map_type', width: 100 },
  { title: '分辨率', dataIndex: 'resolution', width: 90 },
  { title: '版本', dataIndex: 'version', width: 80 },
  { title: '状态', dataIndex: 'is_active', width: 80 },
  { title: '创建时间', dataIndex: 'created_at', width: 170 }
]

const pagination = ref({ current: 1, pageSize: 20, total: 0, showTotal: true, showPageSize: true })

async function loadData() {
  try {
    loading.value = true
    const params: any = { page: pagination.value.current, page_size: pagination.value.pageSize }
    if (form.value.map_name) params.map_name = form.value.map_name
    if (form.value.map_type) params.map_type = form.value.map_type
    const res = await getMaps(deviceId.value, params)
    data.value = res.data?.maps || res.data || []
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
  form.value = { map_name: '', map_type: '', resolution: 0.05, description: '' }
  handleSearch()
}

function handleCreate() {
  modalTitle.value = '新建地图'
  form.value = { map_name: '', map_type: 'grid', resolution: 0.05, description: '' }
  modalVisible.value = true
}

async function handleSubmit(done: (val: boolean) => void) {
  try {
    await updateMap(deviceId.value, { ...form.value, action: 'create' })
    Message.success('创建成功')
    modalVisible.value = false
    loadData()
    done(true)
  } catch (err: any) {
    Message.error('创建失败: ' + err.message)
    done(false)
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.page-container { padding: 16px; }
.search-form { margin-bottom: 16px; padding: 16px; background: var(--color-fill-lightest); border-radius: 4px; }
</style>

