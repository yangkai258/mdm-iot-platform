<template>
  <div class="page-container">
    <a-card class="general-card" title="浣撻噸杩借釜">
      <template #extra>
        <a-button type="primary" @click="showAddModal = true"><icon-plus />璁板綍浣撻噸</a-button>
      </template>
      <a-row :gutter="16">
        <a-col :span="16">
          <div class="search-form">
            <a-form :model="form" layout="inline">
              <a-form-item label="璁惧ID">
                <a-input v-model="form.deviceId" placeholder="璇疯緭鍏ヨ澶嘔D" style="width: 160px" />
              </a-form-item>
              <a-form-item label="鏃ユ湡鑼冨洿">
                <a-range-picker v-model="form.dateRange" style="width: 240px" />
              </a-form-item>
              <a-form-item>
                <a-button type="primary" @click="loadData">鏌ヨ</a-button>
                <a-button @click="handleReset">閲嶇疆</a-button>
              </a-form-item>
            </a-form>
          </div>
          <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id">
            <template #weight="{ record }">
              <span :class="{ 'weight-warning': isWeightAbnormal(record.weight) }">{{ record.weight }}kg</span>
            </template>
            <template #actions="{ record }">
              <a-button type="primary" size="small" @click="handleEdit(record)">缂栬緫</a-button>
              <a-button size="small" status="danger" @click="handleDelete(record)">鍒犻櫎</a-button>
            </template>
          </a-table>
        </a-col>
        <a-col :span="8">
          <a-card title="浣撻噸瓒嬪娍" size="small">
            <div style="height: 300px; display: flex; align-items: center; justify-content: center;">
              <a-empty description="瓒嬪娍鍥捐〃" />
            </div>
          </a-card>
          <a-card title="鐩爣璁剧疆" size="small" style="margin-top: 16px">
            <a-form :model="targetForm" layout="vertical">
              <a-form-item label="鐩爣浣撻噸(kg)">
                <a-input-number v-model="targetForm.targetWeight" :min="0" :precision="1" style="width: 100%" />
              </a-form-item>
              <a-form-item>
                <a-button type="primary" @click="handleSaveTarget">淇濆瓨鐩爣</a-button>
              </a-form-item>
            </a-form>
          </a-card>
        </a-col>
      </a-row>
    </a-card>

    <a-modal v-model:visible="showAddModal" :title="isEditing ? '缂栬緫浣撻噸' : '璁板綍浣撻噸'" @ok="handleSubmit">
      <a-form :model="weightForm" layout="vertical">
        <a-form-item label="璁惧">
          <a-select v-model="weightForm.deviceId" placeholder="璇烽€夋嫨璁惧">
            <a-option value="DEV001">DEV001 - 璞嗚眴</a-option>
            <a-option value="DEV002">DEV002 - 鏃鸿储</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="浣撻噸(kg)">
          <a-input-number v-model="weightForm.weight" :min="0" :precision="1" style="width: 100%" />
        </a-form-item>
        <a-form-item label="璁板綍鏃ユ湡">
          <a-date-picker v-model="weightForm.date" style="width: 100%" />
        </a-form-item>
        <a-form-item label="澶囨敞">
          <a-textarea v-model="weightForm.remark" :rows="2" placeholder="璇疯緭鍏ュ锟" />
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
  { title: '璁惧ID', dataIndex: 'device_id', width: 120 },
  { title: '瀹犵墿鍚嶇О', dataIndex: 'pet_name', width: 100 },
  { title: '浣撻噸', slotName: 'weight', width: 100 },
  { title: '鍙樺寲', dataIndex: 'change', width: 80 },
  { title: '璁板綍鏃ユ湡', dataIndex: 'record_date', width: 120 },
  { title: '澶囨敞', dataIndex: 'remark', ellipsis: true },
  { title: '鎿嶄綔', slotName: 'actions', width: 140 }
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
    { id: 1, device_id: 'DEV001', pet_name: '璞嗚眴', weight: 5.2, change: '+0.1', record_date: '2026-04-09', remark: '姝ｅ父' },
    { id: 2, device_id: 'DEV001', pet_name: '璞嗚眴', weight: 5.1, change: '0', record_date: '2026-04-08', remark: '' },
    { id: 3, device_id: 'DEV001', pet_name: '璞嗚眴', weight: 5.1, change: '-0.1', record_date: '2026-04-07', remark: '鑺傞锟? },
    { id: 4, device_id: 'DEV002', pet_name: '鏃鸿储', weight: 6.8, change: '+0.2', record_date: '2026-04-09', remark: '姝ｅ父' },
    { id: 5, device_id: 'DEV002', pet_name: '鏃鸿储', weight: 6.6, change: '0', record_date: '2026-04-08', remark: '' }
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
  Message.success('鍒犻櫎鎴愬姛')
}

const handleSubmit = () => {
  if (isEditing.value) {
    const idx = data.value.findIndex(d => d.id === weightForm.id)
    if (idx !== -1) data.value[idx] = { ...weightForm }
    Message.success('缂栬緫鎴愬姛')
  } else {
    data.value.unshift({
      id: Date.now(),
      device_id: weightForm.deviceId,
      pet_name: weightForm.deviceId === 'DEV001' ? '璞嗚眴' : '鏃鸿储',
      weight: weightForm.weight,
      change: 'new',
      record_date: weightForm.date?.format('YYYY-MM-DD') || new Date().toLocaleDateString(),
      remark: weightForm.remark
    })
    Message.success('娣诲姞鎴愬姛')
  }
  showAddModal.value = false
}

const handleSaveTarget = () => {
  Message.success('鐩爣浣撻噸宸蹭繚锟?)
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