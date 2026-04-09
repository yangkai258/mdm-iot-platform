<template>
  <div class="page-container">
    <a-card class="general-card" title="饮食记录">
      <template #extra>
        <a-button type="primary" @click="showAddModal = true"><icon-plus />添加记录</a-button>
      </template>
      <a-row :gutter="16">
        <a-col :span="16">
          <div class="search-form">
            <a-form :model="form" layout="inline">
              <a-form-item label="设备ID">
                <a-select v-model="form.deviceId" placeholder="请选择设备" allow-clear style="width: 140px">
                  <a-option value="DEV001">DEV001 - 豆豆</a-option>
                  <a-option value="DEV002">DEV002 - 旺财</a-option>
                </a-select>
              </a-form-item>
              <a-form-item label="日期">
                <a-date-picker v-model="form.date" style="width: 140px" />
              </a-form-item>
              <a-form-item>
                <a-button type="primary" @click="loadData">查询</a-button>
                <a-button @click="handleReset">重置</a-button>
              </a-form-item>
            </a-form>
          </div>
          <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id">
            <template #meal-type="{ record }">
              <a-tag :color="getMealTypeColor(record.meal_type)">{{ getMealTypeText(record.meal_type) }}</a-tag>
            </template>
            <template #actions="{ record }">
              <a-button type="primary" size="small" @click="handleEdit(record)">编辑</a-button>
              <a-button size="small" status="danger" @click="handleDelete(record)">删除</a-button>
            </template>
          </a-table>
        </a-col>
        <a-col :span="8">
          <a-card title="营养统计" size="small">
            <a-statistic title="今日摄入" :value="todayCalories" suffix="kcal" />
            <a-divider />
            <a-row :gutter="16">
              <a-col :span="8">
                <a-statistic title="蛋白�? :value="nutrients.protein" suffix="g" />
              </a-col>
              <a-col :span="8">
                <a-statistic title="脂肪" :value="nutrients.fat" suffix="g" />
              </a-col>
              <a-col :span="8">
                <a-statistic title="碳水" :value="nutrients.carbs" suffix="g" />
              </a-col>
            </a-row>
          </a-card>
        </a-col>
      </a-row>
    </a-card>

    <a-modal v-model:visible="showAddModal" :title="isEditing ? '编辑记录' : '添加记录'" @ok="handleSubmit">
      <a-form :model="dietForm" layout="vertical">
        <a-form-item label="设备">
          <a-select v-model="dietForm.deviceId" placeholder="请选择设备">
            <a-option value="DEV001">DEV001 - 豆豆</a-option>
            <a-option value="DEV002">DEV002 - 旺财</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="餐次">
          <a-select v-model="dietForm.mealType" placeholder="请选择餐次">
            <a-option value="breakfast">早餐</a-option>
            <a-option value="lunch">午餐</a-option>
            <a-option value="dinner">晚餐</a-option>
            <a-option value="snack">零食</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="食物">
          <a-input v-model="dietForm.food" placeholder="请输入食物名�? />
        </a-form-item>
        <a-form-item label="热量(kcal)">
          <a-input-number v-model="dietForm.calories" :min="0" style="width: 100%" />
        </a-form-item>
        <a-form-item label="时间">
          <a-time-picker v-model="dietForm.time" format="HH:mm" style="width: 100%" />
        </a-form-item>
        <a-form-item label="备注">
          <a-textarea v-model="dietForm.remark" :rows="2" placeholder="请输入备�? />
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
const form = reactive({ deviceId: '', date: null })
const data = ref([])
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })
const todayCalories = ref(0)
const nutrients = reactive({ protein: 0, fat: 0, carbs: 0 })
const dietForm = reactive({ id: '', deviceId: '', mealType: '', food: '', calories: 0, time: null, remark: '' })

const columns = [
  { title: '时间', dataIndex: 'time', width: 80 },
  { title: '设备ID', dataIndex: 'device_id', width: 100 },
  { title: '宠物', dataIndex: 'pet_name', width: 80 },
  { title: '餐次', slotName: 'meal_type', width: 80 },
  { title: '食物', dataIndex: 'food', ellipsis: true },
  { title: '热量', dataIndex: 'calories', width: 80 },
  { title: '备注', dataIndex: 'remark', ellipsis: true },
  { title: '操作', slotName: 'actions', width: 140 }
]

const getMealTypeColor = (type) => {
  const colors = { breakfast: 'orange', lunch: 'green', dinner: 'blue', snack: 'purple' }
  return colors[type] || 'gray'
}

const getMealTypeText = (type) => {
  const texts = { breakfast: '早餐', lunch: '午餐', dinner: '晚餐', snack: '零食' }
  return texts[type] || type
}

const loadData = async () => {
  loading.value = true
  try {
    const res = await fetch('/api/v1/health/diet').then(r => r.json())
    if (res.code === 0) {
      data.value = res.data?.list || []
    } else {
      loadMockData()
    }
  } catch { loadMockData() } finally { loading.value = false }
}

const loadMockData = () => {
  data.value = [
    { id: 1, time: '08:00', device_id: 'DEV001', pet_name: '豆豆', meal_type: 'breakfast', food: '狗粮 50g', calories: 150, remark: '' },
    { id: 2, time: '12:00', device_id: 'DEV001', pet_name: '豆豆', meal_type: 'lunch', food: '狗粮 80g', calories: 240, remark: '' },
    { id: 3, time: '18:00', device_id: 'DEV001', pet_name: '豆豆', meal_type: 'dinner', food: '狗粮 60g + 鸡胸�?, calories: 320, remark: '加餐' },
    { id: 4, time: '10:00', device_id: 'DEV002', pet_name: '旺财', meal_type: 'snack', food: '零食', calories: 50, remark: '' },
    { id: 5, time: '08:00', device_id: 'DEV002', pet_name: '旺财', meal_type: 'breakfast', food: '狗粮 80g', calories: 200, remark: '' }
  ]
  todayCalories.value = 960
  nutrients.protein = 45
  nutrients.fat = 25
  nutrients.carbs = 80
}

const handleReset = () => {
  form.deviceId = ''
  form.date = null
  loadData()
}

const handleEdit = (record) => {
  isEditing.value = true
  Object.assign(dietForm, record)
  showAddModal.value = true
}

const handleDelete = (record) => {
  data.value = data.value.filter(d => d.id !== record.id)
  Message.success('删除成功')
}

const handleSubmit = () => {
  if (isEditing.value) {
    const idx = data.value.findIndex(d => d.id === dietForm.id)
    if (idx !== -1) data.value[idx] = { ...dietForm }
    Message.success('编辑成功')
  } else {
    data.value.unshift({
      id: Date.now(),
      device_id: dietForm.deviceId,
      pet_name: dietForm.deviceId === 'DEV001' ? '豆豆' : '旺财',
      ...dietForm
    })
    Message.success('添加成功')
  }
  showAddModal.value = false
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