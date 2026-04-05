<template>
  <div>
    <Breadcrumb :items="['Home','Console','']" />
    <div class="page-container">
      <a-card class="general-card" title="安全区域">
        <template #extra>
          <a-button type="primary" @click="handleCreate"><icon-plus />新建</a-button>
        </template>
        <div class="search-form">
          <a-form :model="form" layout="inline">
            <a-form-item label="区域名称"><a-input v-model="form.zone_name" placeholder="请输入区域名称" /></a-form-item>
            <a-form-item label="区域类型">
              <a-select v-model="form.zone_type" placeholder="选择类型" allow-clear style="width: 120px">
                <a-option value="forbidden">禁区</a-option>
                <a-option value="caution">警戒区</a-option>
                <a-option value="safe">安全区</a-option>
              </a-select>
            </a-form-item>
            <a-form-item>
              <a-button type="primary" @click="handleSearch">查询</a-button>
              <a-button @click="handleReset">重置</a-button>
            </a-form-item>
          </a-form>
        </div>
        <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" />
        <a-modal v-model:visible="modalVisible" :title="modalTitle" @before-ok="handleSubmit">
          <a-form :model="form" label-col-flex="100px">
            <a-form-item label="区域名称" required><a-input v-model="form.zone_name" placeholder="请输入区域名称" /></a-form-item>
            <a-form-item label="区域类型">
              <a-select v-model="form.zone_type" placeholder="选择类型">
                <a-option value="forbidden">禁区</a-option>
                <a-option value="caution">警戒区</a-option>
                <a-option value="safe">安全区</a-option>
              </a-select>
            </a-form-item>
            <a-form-item label="中心纬度"><a-input-number v-model="form.center_lat" placeholder="中心纬度" style="width:100%" /></a-form-item>
            <a-form-item label="中心经度"><a-input-number v-model="form.center_lng" placeholder="中心经度" style="width:100%" /></a-form-item>
            <a-form-item label="半径(米)"><a-input-number v-model="form.radius" placeholder="半径" style="width:100%" /></a-form-item>
          </a-form>
        </a-modal>
      </a-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { getSafetyZones, createSafetyZone, updateSafetyZone, deleteSafetyZone } from '@/api/embodied'
import { Message, Modal } from '@arco-design/web-vue'
import { IconPlus } from '@arco-design/web-vue/es/icon'

const route = useRoute()
const deviceId = ref(route.params.device_id as string || '')

const loading = ref(false)
const data = ref<any[]>([])
const modalVisible = ref(false)
const modalTitle = ref('新建安全区域')
const form = ref<any>({
  zone_name: '',
  zone_type: 'forbidden',
  boundary: '{}',
  center_lat: 0,
  center_lng: 0,
  radius: 10
})

const columns = [
  { title: 'ID', dataIndex: 'id', width: 70 },
  { title: '区域名称', dataIndex: 'zone_name', width: 140 },
  { title: '设备ID', dataIndex: 'device_id', width: 120 },
  { title: '类型', dataIndex: 'zone_type', width: 100 },
  { title: '中心坐标', dataIndex: 'center_lat', width: 160, render: ({ record }: any) => `${record.center_lat}, ${record.center_lng}` },
  { title: '半径(m)', dataIndex: 'radius', width: 90 },
  { title: '启用', dataIndex: 'is_active', width: 70 }
]

const pagination = ref({ current: 1, pageSize: 20, total: 0, showTotal: true, showPageSize: true })

async function loadData() {
  try {
    loading.value = true
    const params: any = { page: pagination.value.current, page_size: pagination.value.pageSize }
    const res = await getSafetyZones(deviceId.value, params)
    data.value = res.data?.zones || res.data || []
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
  form.value = { zone_name: '', zone_type: '', boundary: '{}', center_lat: 0, center_lng: 0, radius: 10 }
  handleSearch()
}

function handleCreate() {
  modalTitle.value = '新建安全区域'
  form.value = { zone_name: '', zone_type: 'forbidden', boundary: '{}', center_lat: 0, center_lng: 0, radius: 10 }
  modalVisible.value = true
}

async function handleSubmit(done: (val: boolean) => void) {
  try {
    await createSafetyZone(deviceId.value, form.value)
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
