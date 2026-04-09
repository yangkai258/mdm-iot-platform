<template>
  <div class="page-container">
    <a-card class="general-card" title="体重追踪">
      <template #extra>
        <a-button type="primary" @click="showAddModal = true"><icon-plus />记录体重</a-button>
      </template>
      <a-row :gutter="16">
        <a-col :span="16">
          <div class="search-form">
            <a-form :model="form" layout="inline">
              <a-form-item label="设备ID">
                <a-input v-model="form.deviceId" placeholder="请输入设备ID" style="width: 160px" />
              </a-form-item>
              <a-form-item label="日期范围">
                <a-range-picker v-model="form.dateRange" style="width: 240px" />
              </a-form-item>
              <a-form-item>
                <a-button type="primary" @click="loadData">查询</a-button>
                <a-button @click="handleReset">重置</a-button>
              </a-form-item>
            </a-form>
          </div>
          <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id">
            <template #weight="{ record }">
              <span :class="{ 'weight-warning': isWeightAbnormal(record.weight) }">{{ record.weight }}kg</span>
            </template>
            <template #actions="{ record }">
              <a-button type="primary" size="small" @click="handleEdit(record)">编辑</a-button>
              <a-button size="small" status="danger" @click="handleDelete(record)">删除</a-button>
            </template>
          </a-table>
        </a-col>
        <a-col :span="8">
          <a-card title="体重趋势" size="small">
            <div style="height: 300px; display: flex; align-items: center; justify-content: center;">
              <a-empty description="趋势图表" />
            </div>
          </a-card>
          <a-card title="目标设置" size="small" style="margin-top: 16px">
            <a-form :model="targetForm" layout="vertical">
              <a-form-item label="目标体重(kg)">
                <a-input-number v-model="targetForm.targetWeight" :min="0" :precision="1" style="width: 100%" />
              </a-form-item>
              <a-form-item>
                <a-button type="primary" @click="handleSaveTarget">保存目标</a-button>
              </a-form-item>
            </a-form>
          </a-card>
        </a-col>
      </a-row>
    </a-card>

    <a-modal v-model:visible="showAddModal" :title="isEditing ? '编辑体重' : '记录体重'" @ok="handleSubmit">
      <a-form :model="weightForm" layout="vertical">
        <a-form-item label="设备">
          <a-select v-model="weightForm.deviceId" placeholder="请选择设备">
            <a-option value="DEV001">DEV001 - 豆豆</a-option>
            <a-option value="DEV002">DEV002 - 旺财</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="体重(kg)">
          <a-input-number v-model="weightForm.weight" :min="0" :precision="1" style="width: 100%" />
        </a-form-item>
        <a-form-item label="记录日期">
          <a-date-picker v-model="weightForm.date" style="width: 100%" />
        </a-form-item>
        <a-form-item label="备注">
          <a-textarea v-model="weightForm.remark" :rows="2" placeholder="请输入备�? />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import { IconPlus } from '@arco-design/web-vue/es/icon'

const loading = ref(false)
const showAddModal = ref(false)
const isEditing = ref(false)
const form = reactive({ deviceId: '', dateRange: [] })
const data = ref([])
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })
const targetForm = reactive({ targetWeight: 5.0 })
const weightForm = reactive({ id: '', deviceId: '', weight: 0, date: null, remark: '' })

const columns = [
  { title: '设备ID', dataIndex: 'device_id', width: 120 },
  { title: '宠物名称', dataIndex: 'pet_name', width: 100 },
  { title: '体重', slotName: 'weight', width: 100 },
  { title: '变化', dataIndex: 'change', width: 80 },
  { title: '记录日期', dataIndex: 'record_date', width: 120 },
  { title: '备注', dataIndex: 'remark', ellipsis: true },
  { title: '操作', slotName: 'actions', width: 140 }
]

const isWeightAbnormal = (weight) => weight < 3 || weight > 8

const loadData = async () => {
  loading.value = true
  try {
    const res = await fetch('/api/v1/health/weight').then(r => r.json())
    if (res.code === 0) {
      data.value = res.data?.list || []
    } else {
      loadMockData()
    }
  } catch { loadMockData() } finally { loading.value = false }
}

const loadMockData = () => {
  data.value = [
    { id: 1, device_id: 'DEV001', pet_name: '豆豆', weight: 5.2, change: '+0.1', record_date: '2026-04-09', remark: '正常' },
    { id: 2, device_id: 'DEV001', pet_name: '豆豆', weight: 5.1, change: '0', record_date: '2026-04-08', remark: '' },
    { id: 3, device_id: 'DEV001', pet_name: '豆豆', weight: 5.1, change: '-0.1', record_date: '2026-04-07', remark: '节食�? },
    { id: 4, device_id: 'DEV002', pet_name: '旺财', weight: 6.8, change: '+0.2', record_date: '2026-04-09', remark: '正常' },
    { id: 5, device_id: 'DEV002', pet_name: '旺财', weight: 6.6, change: '0', record_date: '2026-04-08', remark: '' }
  ]
}

const handleReset = () => {
  form.deviceId = ''
  form.dateRange = []
  loadData()
}

const handleEdit = (record) => {
  isEditing.value = true
  Object.assign(weightForm, record)
  showAddModal.value = true
}

const handleDelete = (record) => {
  data.value = data.value.filter(d => d.id !== record.id)
  Message.success('删除成功')
}

const handleSubmit = () => {
  if (isEditing.value) {
    const idx = data.value.findIndex(d => d.id === weightForm.id)
    if (idx !== -1) data.value[idx] = { ...weightForm }
    Message.success('编辑成功')
  } else {
    data.value.unshift({
      id: Date.now(),
      device_id: weightForm.deviceId,
      pet_name: weightForm.deviceId === 'DEV001' ? '豆豆' : '旺财',
      weight: weightForm.weight,
      change: 'new',
      record_date: weightForm.date?.format('YYYY-MM-DD') || new Date().toLocaleDateString(),
      remark: weightForm.remark
    })
    Message.success('添加成功')
  }
  showAddModal.value = false
}

const handleSaveTarget = () => {
  Message.success('目标体重已保�?)
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
.weight-warning { color: #f53f3f; font-weight: 500; }
</style>