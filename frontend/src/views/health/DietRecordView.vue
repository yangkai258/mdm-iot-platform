<template>
  <div class="page-container">
    <a-card class="general-card" title="楗璁板綍">
      <template #extra>
        <a-button type="primary" @click="showAddModal = true"><icon-plus />娣诲姞璁板綍</a-button>
      </template>
      <a-row :gutter="16">
        <a-col :span="16">
          <div class="search-form">
            <a-form :model="form" layout="inline">
              <a-form-item label="璁惧ID">
                <a-select v-model="form.deviceId" placeholder="璇烽€夋嫨璁惧" allow-clear style="width: 140px">
                  <a-option value="DEV001">DEV001 - 璞嗚眴</a-option>
                  <a-option value="DEV002">DEV002 - 鏃鸿储</a-option>
                </a-select>
              </a-form-item>
              <a-form-item label="鏃ユ湡">
                <a-date-picker v-model="form.date" style="width: 140px" />
              </a-form-item>
              <a-form-item>
                <a-button type="primary" @click="loadData">鏌ヨ</a-button>
                <a-button @click="handleReset">閲嶇疆</a-button>
              </a-form-item>
            </a-form>
          </div>
          <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id">
            <template #meal-type="{ record }">
              <a-tag :color="getMealTypeColor(record.meal_type)">{{ getMealTypeText(record.meal_type) }}</a-tag>
            </template>
            <template #actions="{ record }">
              <a-button type="primary" size="small" @click="handleEdit(record)">缂栬緫</a-button>
              <a-button size="small" status="danger" @click="handleDelete(record)">鍒犻櫎</a-button>
            </template>
          </a-table>
        </a-col>
        <a-col :span="8">
          <a-card title="钀ュ吇缁熻" size="small">
            <a-statistic title="浠婃棩鎽勫叆" :value="todayCalories" suffix="kcal" />
            <a-divider />
            <a-row :gutter="16">
              <a-col :span="8">
                <a-statistic title="铔嬬櫧锟? :value="nutrients.protein" suffix="g" />
              </a-col>
              <a-col :span="8">
                <a-statistic title="鑴傝偑" :value="nutrients.fat" suffix="g" />
              </a-col>
              <a-col :span="8">
                <a-statistic title="纰虫按" :value="nutrients.carbs" suffix="g" />
              </a-col>
            </a-row>
          </a-card>
        </a-col>
      </a-row>
    </a-card>

    <a-modal v-model:visible="showAddModal" :title="isEditing ? '缂栬緫璁板綍' : '娣诲姞璁板綍'" @ok="handleSubmit">
      <a-form :model="dietForm" layout="vertical">
        <a-form-item label="璁惧">
          <a-select v-model="dietForm.deviceId" placeholder="璇烽€夋嫨璁惧">
            <a-option value="DEV001">DEV001 - 璞嗚眴</a-option>
            <a-option value="DEV002">DEV002 - 鏃鸿储</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="椁愭">
          <a-select v-model="dietForm.mealType" placeholder="璇烽€夋嫨椁愭">
            <a-option value="breakfast">鏃╅</a-option>
            <a-option value="lunch">鍗堥</a-option>
            <a-option value="dinner">鏅氶</a-option>
            <a-option value="snack">闆堕</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="椋熺墿">
          <a-input v-model="dietForm.food" placeholder="璇疯緭鍏ラ鐗╁悕锟" />
        </a-form-item>
        <a-form-item label="鐑噺(kcal)">
          <a-input-number v-model="dietForm.calories" :min="0" style="width: 100%" />
        </a-form-item>
        <a-form-item label="鏃堕棿">
          <a-time-picker v-model="dietForm.time" format="HH:mm" style="width: 100%" />
        </a-form-item>
        <a-form-item label="澶囨敞">
          <a-textarea v-model="dietForm.remark" :rows="2" placeholder="璇疯緭鍏ュ锟" />
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
  { title: '鏃堕棿', dataIndex: 'time', width: 80 },
  { title: '璁惧ID', dataIndex: 'device_id', width: 100 },
  { title: '瀹犵墿', dataIndex: 'pet_name', width: 80 },
  { title: '椁愭', slotName: 'meal_type', width: 80 },
  { title: '椋熺墿', dataIndex: 'food', ellipsis: true },
  { title: '鐑噺', dataIndex: 'calories', width: 80 },
  { title: '澶囨敞', dataIndex: 'remark', ellipsis: true },
  { title: '鎿嶄綔', slotName: 'actions', width: 140 }
]

const getMealTypeColor = (type) => {
  const colors = { breakfast: 'orange', lunch: 'green', dinner: 'blue', snack: 'purple' }
  return colors[type] || 'gray'
}

const getMealTypeText = (type) => {
  const texts = { breakfast: '鏃╅', lunch: '鍗堥', dinner: '鏅氶', snack: '闆堕' }
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
    { id: 1, time: '08:00', device_id: 'DEV001', pet_name: '璞嗚眴', meal_type: 'breakfast', food: '鐙楃伯 50g', calories: 150, remark: '' },
    { id: 2, time: '12:00', device_id: 'DEV001', pet_name: '璞嗚眴', meal_type: 'lunch', food: '鐙楃伯 80g', calories: 240, remark: '' },
    { id: 3, time: '18:00', device_id: 'DEV001', pet_name: '璞嗚眴', meal_type: 'dinner', food: '鐙楃伯 60g + 楦¤兏锟?, calories: 320, remark: '鍔犻' },
    { id: 4, time: '10:00', device_id: 'DEV002', pet_name: '鏃鸿储', meal_type: 'snack', food: '闆堕', calories: 50, remark: '' },
    { id: 5, time: '08:00', device_id: 'DEV002', pet_name: '鏃鸿储', meal_type: 'breakfast', food: '鐙楃伯 80g', calories: 200, remark: '' }
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
  Message.success('鍒犻櫎鎴愬姛')
}

const handleSubmit = () => {
  if (isEditing.value) {
    const idx = data.value.findIndex(d => d.id === dietForm.id)
    if (idx !== -1) data.value[idx] = { ...dietForm }
    Message.success('缂栬緫鎴愬姛')
  } else {
    data.value.unshift({
      id: Date.now(),
      device_id: dietForm.deviceId,
      pet_name: dietForm.deviceId === 'DEV001' ? '璞嗚眴' : '鏃鸿储',
      ...dietForm
    })
    Message.success('娣诲姞鎴愬姛')
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