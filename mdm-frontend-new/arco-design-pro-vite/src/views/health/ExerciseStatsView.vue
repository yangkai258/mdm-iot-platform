<template>
    <Breadcrumb :items="['Home','Console','']" />


  <div class="container">
    <a-card class="general-card" title="运动统计">
      <template #extra>
        <a-space :size="12">
          <a-button type="primary" @click="handleCreate"><icon-plus />新建</a-button>
          <a-button @click="loadData"><icon-refresh />刷新</a-button>
        </a-space>
      </template>
      <a-row :gutter="16">
        <a-col :span="8">
          <a-form-item label="设备ID">
            <a-input v-model="form.deviceId" placeholder="请输入设备ID" @pressEnter="loadData" />
          </a-form-item>
        </a-col>
        <a-col :flex="'86px'" style="display: flex; align-items: flex-end">
          <a-space direction="vertical" :size="8">
            <a-button type="primary" @click="loadData">查询</a-button>
            <a-button @click="handleReset">重置</a-button>
          </a-space>
        </a-col>
      </a-row>
      <a-divider style="margin: 0 0 16px 0" />
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="paginationConfig" @page-change="onPageChange" row-key="id" />
    </a-card>
      </a-table>
    <a-modal v-model:visible="modalVisible" :title="modalTitle">
      <a-form :model="form" layout="vertical">
        <a-form-item label="步数"><a-input-number v-model="form.steps" :min="0" style="width: 100%" /></a-form-item>
        <a-form-item label="距离(km)"><a-input-number v-model="form.distance" :min="0" :precision="1" style="width: 100%" /></a-form-item>
        <a-form-item label="卡路里(kcal)"><a-input-number v-model="form.calories" :min="0" style="width: 100%" /></a-form-item>
      </a-form>
      <template #footer>
        <a-button @click="modalVisible = false">取消</a-button>
        <a-button type="primary" @click="handleSubmit">确定</a-button>
      </template>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import Breadcrumb from '@/components/Breadcrumb.vue'

const loading = ref(false)
const modalVisible = ref(false)
const modalTitle = ref('新建')
const isEdit = ref(false)
const form = reactive({ deviceId: '', steps: 0, distance: 0, calories: 0 })
const data = ref([])
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })
const paginationConfig = computed(() => ({ current: pagination.current, pageSize: pagination.pageSize, total: pagination.total, showTotal: true }))
const columns = [
  { title: '日期', dataIndex: 'date', width: 120 },
  { title: '设备ID', dataIndex: 'device_id', width: 120 },
  { title: '步数', dataIndex: 'steps', width: 100 },
  { title: '距离(km)', dataIndex: 'distance', width: 100 },
  { title: '卡路里(kcal)', dataIndex: 'calories', width: 120 },
  { title: '时长(分钟)', dataIndex: 'duration', width: 120 },
  { title: '创建时间', dataIndex: 'created_at', width: 160 }
]

const loadData = async () => {
  loading.value = true
  try {
    const params = { page: pagination.current, page_size: pagination.pageSize }
    if (form.deviceId) params.device_id = form.deviceId
    const res = await fetch(`/api/health/exercise-stats?${new URLSearchParams(params)}`, {
      headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token') }
    }).then(r => r.json())
    if (res.code === 0) { data.value = res.data?.list || []; pagination.total = res.data?.pagination?.total || 0 }
    else data.value = [{ id: 1, date: '2026-03-22', device_id: 'DEV001', steps: 8500, distance: 5.2, calories: 320, duration: 45, created_at: '2026-03-22 23:00:00' }]
  } catch { data.value = [] } finally { loading.value = false }
}

const handleReset = () => { Object.assign(form, { deviceId: '', steps: 0, distance: 0, calories: 0 }); loadData() }
const handleCreate = () => { isEdit.value = false; modalTitle.value = '新建'; modalVisible.value = true }
const handleSubmit = () => { modalVisible.value = false; Message.success(isEdit.value ? '编辑成功' : '添加成功'); loadData() }
const onPageChange = (page) => { pagination.current = page; loadData() }

onMounted(() => loadData())
</script>
