<template>
  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>高级功能</a-breadcrumb-item>
      <a-breadcrumb-item>饮食记录</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 页面标题 -->
    <div class="pro-page-header">
      <h2 class="pro-page-title">饮食记录</h2>
      <p class="pro-page-desc">记录宠物每日饮食，追踪营养摄入与饮食习惯</p>
    </div>

    <!-- 标签页切换 -->
    <a-tabs v-model="activeTab" class="pro-tabs">
      <a-tab-pane key="records" title="饮食记录">
        <!-- 搜索筛选区 -->
        <div class="pro-search-bar">
          <a-space>
            <a-select v-model="petFilter" placeholder="选择宠物" allow-clear style="width: 160px" @change="loadRecords">
              <a-option v-for="p in pets" :key="p.id" :value="p.id">{{ p.name }}</a-option>
            </a-select>
            <a-select v-model="mealTypeFilter" placeholder="餐食类型" allow-clear style="width: 140px" @change="loadRecords">
              <a-option value="breakfast">早餐</a-option>
              <a-option value="lunch">午餐</a-option>
              <a-option value="dinner">晚餐</a-option>
              <a-option value="snack">零食</a-option>
            </a-select>
            <a-range-picker v-model="dateRange" style="width: 240px" @change="loadRecords" />
            <a-input-search v-model="keyword" placeholder="搜索食物名称" style="width: 200px" @search="loadRecords" search-button />
          </a-space>
        </div>

        <!-- 操作按钮区 -->
        <div class="pro-action-bar">
          <a-space>
            <a-button type="primary" @click="showRecordModal(null)">
              <template #icon><icon-add /></template>
              记录饮食
            </a-button>
            <a-button @click="loadRecords">
              <template #icon><icon-refresh /></template>
              刷新
            </a-button>
          </a-space>
        </div>

        <!-- 饮食记录列表 -->
        <div class="pro-content-area">
          <a-table
            :columns="recordColumns"
            :data="records"
            :loading="loading"
            :pagination="pagination"
            @page-change="onPageChange"
            row-key="id"
          >
            <template #pet_name="{ record }">
              <a-avatar :style="{ backgroundColor: '#ff7d00' }" :size="28">
                {{ record.pet_name?.charAt(0) || '?' }}
              </a-avatar>
              <span style="margin-left: 8px">{{ record.pet_name }}</span>
            </template>
            <template #food_name="{ record }">
              <span class="food-name">{{ record.food_name }}</span>
              <div v-if="record.food_brand" class="food-brand">{{ record.food_brand }}</div>
            </template>
            <template #meal_type="{ record }">
              <a-tag :color="getMealTypeColor(record.meal_type)">
                {{ getMealTypeLabel(record.meal_type) }}
              </a-tag>
            </template>
            <template #amount="{ record }">
              {{ record.amount }}{{ record.unit }}
            </template>
            <template #calories="{ record }">
              <span class="calories-value">{{ record.calories || '—' }}</span>
              <span v-if="record.calories" class="calories-unit"> kcal</span>
            </template>
            <template #recorded_at="{ record }">
              {{ formatDateTime(record.recorded_at) }}
            </template>
            <template #food_photo="{ record }">
              <img
                v-if="record.food_photo"
                :src="record.food_photo"
                class="food-thumbnail"
                @click="previewPhoto(record.food_photo)"
              />
              <span v-else class="text-muted">—</span>
            </template>
            <template #actions="{ record }">
              <a-space>
                <a-button type="text" size="small" @click="showRecordModal(record)">编辑</a-button>
                <a-button type="text" size="small" status="danger" @click="handleDeleteRecord(record)">删除</a-button>
              </a-space>
            </template>
          </a-table>
        </div>
      </a-tab-pane>

      <a-tab-pane key="summary" title="饮食摘要">
        <!-- 摘要统计 -->
        <a-row :gutter="16" style="margin-bottom: 20px">
          <a-col :span="6">
            <a-card class="stat-card" hoverable>
              <a-statistic :value="summaryData.total_records || 0" suffix="条">
                <template #prefix><icon-book :size="18" style="color:#1659f5"/></template>
                <template #title>本周记录</template>
              </a-statistic>
            </a-card>
          </a-col>
          <a-col :span="6">
            <a-card class="stat-card" hoverable>
              <a-statistic :value="summaryData.avg_calories || 0" :precision="0" suffix="kcal">
                <template #prefix><icon-fire :size="18" style="color:#ff7d00"/></template>
                <template #title>日均热量</template>
              </a-statistic>
            </a-card>
          </a-col>
          <a-col :span="6">
            <a-card class="stat-card" hoverable>
              <a-statistic :value="summaryData.avg_meals_per_day || 0" :precision="1" suffix="次">
                <template #prefix><icon-clock :size="18" style="color:#0fc6c2"/></template>
                <template #title>日均餐次</template>
              </a-statistic>
            </a-card>
          </a-col>
          <a-col :span="6">
            <a-card class="stat-card" hoverable>
              <a-statistic :value="summaryData.food_diversity || 0" :precision="0">
                <template #prefix><icon-safe :size="18" style="color:#00b42a"/></template>
                <template #title>食物多样性</template>
              </a-statistic>
            </a-card>
          </a-col>
        </a-row>

        <!-- 宠物选择 -->
        <div class="pro-search-bar" style="margin-bottom: 16px">
          <a-space>
            <span style="font-weight: 500; color: var(--color-text-2)">选择宠物：</span>
            <a-select v-model="summaryPetId" placeholder="选择宠物" style="width: 160px" @change="loadSummary">
              <a-option v-for="p in pets" :key="p.id" :value="p.id">{{ p.name }}</a-option>
            </a-select>
            <a-select v-model="summaryPeriod" style="width: 120px" @change="loadSummary">
              <a-option value="week">本周</a-option>
              <a-option value="month">本月</a-option>
            </a-select>
          </a-space>
        </div>

        <!-- 热量趋势图 -->
        <a-row :gutter="16">
          <a-col :span="16">
            <a-card title="每日热量摄入趋势" style="margin-bottom: 16px">
              <div ref="caloriesChartRef" style="height: 240px"></div>
            </a-card>
          </a-col>
          <a-col :span="8">
            <a-card title="餐食类型分布">
              <div ref="mealTypeChartRef" style="height: 240px"></div>
            </a-card>
          </a-col>
        </a-row>

        <!-- 营养素分布 -->
        <a-card title="营养素摄入比例" style="margin-top: 16px">
          <a-row :gutter="16">
            <a-col :span="8">
              <div ref="proteinChartRef" style="height: 200px"></div>
            </a-col>
            <a-col :span="8">
              <div ref="fatChartRef" style="height: 200px"></div>
            </a-col>
            <a-col :span="8">
              <div ref="carbChartRef" style="height: 200px"></div>
            </a-col>
          </a-row>
        </a-card>
      </a-tab-pane>
    </a-tabs>

    <!-- 饮食记录弹窗 -->
    <a-modal
      v-model:visible="recordModalVisible"
      :title="isEditRecord ? '编辑饮食记录' : '记录饮食'"
      @ok="handleSaveRecord"
      :width="560"
      @close="resetRecordForm"
    >
      <a-form :model="recordForm" layout="vertical">
        <a-form-item label="宠物" required>
          <a-select v-model="recordForm.pet_id" placeholder="请选择宠物" :disabled="isEditRecord">
            <a-option v-for="p in pets" :key="p.id" :value="p.id">{{ p.name }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="食物名称" required>
          <a-input v-model="recordForm.food_name" placeholder="请输入食物名称" />
        </a-form-item>
        <a-form-item label="食物品牌">
          <a-input v-model="recordForm.food_brand" placeholder="请输入食物品牌（可选）" />
        </a-form-item>
        <a-row :gutter="12">
          <a-col :span="12">
            <a-form-item label="数量" required>
              <a-input-number v-model="recordForm.amount" :min="0.1" :precision="1" :step="0.1" style="width: 100%" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="单位" required>
              <a-select v-model="recordForm.unit" placeholder="请选择单位">
                <a-option value="g">克 (g)</a-option>
                <a-option value="ml">毫升 (ml)</a-option>
                <a-option value="cup">杯</a-option>
                <a-option value="piece">块</a-option>
                <a-option value="bowl">碗</a-option>
              </a-select>
            </a-form-item>
          </a-col>
        </a-row>
        <a-form-item label="餐食类型" required>
          <a-select v-model="recordForm.meal_type" placeholder="请选择餐食类型">
            <a-option value="breakfast">早餐</a-option>
            <a-option value="lunch">午餐</a-option>
            <a-option value="dinner">晚餐</a-option>
            <a-option value="snack">零食</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="记录时间">
          <a-date-picker
            v-model="recordForm.recorded_at"
            style="width: 100%"
            show-time
            format="YYYY-MM-DD HH:mm"
          />
        </a-form-item>
        <a-form-item label="热量（kcal）">
          <a-input-number v-model="recordForm.calories" :min="0" :precision="0" style="width: 100%" placeholder="请输入热量值（可选）" />
        </a-form-item>
        <a-form-item label="备注">
          <a-textarea v-model="recordForm.notes" placeholder="请输入备注" :rows="2" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 图片预览弹窗 -->
    <a-modal v-model:visible="previewVisible" :footer="null" :width="600" @close="previewVisible = false">
      <img :src="previewUrl" style="width: 100%; border-radius: 8px" />
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import {
  getDietRecords,
  createDietRecord,
  updateDietRecord,
  deleteDietRecord,
  getDietSummary
} from '@/api/advanced'
import { Message } from '@arco-design/web-vue'

const activeTab = ref('records')

// Mock pets data
const pets = ref([
  { id: 1, name: '小橘' },
  { id: 2, name: '豆豆' }
])

// ===== 饮食记录 =====
const records = ref([])
const loading = ref(false)
const recordModalVisible = ref(false)
const isEditRecord = ref(false)

const recordForm = reactive({
  id: null,
  pet_id: null,
  food_name: '',
  food_brand: '',
  amount: 100,
  unit: 'g',
  meal_type: '',
  recorded_at: null,
  calories: null,
  notes: ''
})

const recordColumns = [
  { title: '宠物', slotName: 'pet_name', width: 130 },
  { title: '食物', slotName: 'food_name', width: 200 },
  { title: '餐食类型', slotName: 'meal_type', width: 100 },
  { title: '数量', slotName: 'amount', width: 90 },
  { title: '热量', slotName: 'calories', width: 100 },
  { title: '记录时间', slotName: 'recorded_at', width: 170 },
  { title: '照片', slotName: 'food_photo', width: 80 },
  { title: '操作', slotName: 'actions', width: 120 }
]

const pagination = reactive({ current: 1, pageSize: 10, total: 0 })
const petFilter = ref(null)
const mealTypeFilter = ref(null)
const dateRange = ref([])
const keyword = ref('')

async function loadRecords() {
  loading.value = true
  try {
    const params = {
      page: pagination.current,
      page_size: pagination.pageSize,
      pet_id: petFilter.value,
      meal_type: mealTypeFilter.value,
      keyword: keyword.value,
      start_date: dateRange.value[0]?.format('YYYY-MM-DD') || undefined,
      end_date: dateRange.value[1]?.format('YYYY-MM-DD') || undefined
    }
    const res = await getDietRecords(params)
    records.value = res.data?.items || res.data || []
    pagination.total = res.data?.pagination?.total || records.value.length
  } catch (e) {
    Message.error('加载饮食记录失败')
  } finally {
    loading.value = false
  }
}

function onPageChange(page) {
  pagination.current = page
  loadRecords()
}

function showRecordModal(record) {
  if (record) {
    isEditRecord.value = true
    Object.assign(recordForm, {
      id: record.id,
      pet_id: record.pet_id,
      food_name: record.food_name,
      food_brand: record.food_brand || '',
      amount: record.amount,
      unit: record.unit,
      meal_type: record.meal_type,
      recorded_at: record.recorded_at ? new Date(record.recorded_at) : null,
      calories: record.calories,
      notes: record.notes || ''
    })
  } else {
    isEditRecord.value = false
    resetRecordForm()
  }
  recordModalVisible.value = true
}

function resetRecordForm() {
  Object.assign(recordForm, {
    id: null,
    pet_id: null,
    food_name: '',
    food_brand: '',
    amount: 100,
    unit: 'g',
    meal_type: '',
    recorded_at: null,
    calories: null,
    notes: ''
  })
}

async function handleSaveRecord() {
  try {
    if (isEditRecord.value) {
      await updateDietRecord(recordForm.id, recordForm)
      Message.success('更新成功')
    } else {
      await createDietRecord(recordForm)
      Message.success('添加成功')
    }
    recordModalVisible.value = false
    loadRecords()
  } catch (e) {
    Message.error('保存失败')
  }
}

async function handleDeleteRecord(record) {
  try {
    await deleteDietRecord(record.id)
    Message.success('删除成功')
    loadRecords()
  } catch (e) {
    Message.error('删除失败')
  }
}

// ===== 饮食摘要 =====
const summaryData = ref({})
const summaryPetId = ref(null)
const summaryPeriod = ref('week')
const caloriesChartRef = ref(null)
const mealTypeChartRef = ref(null)
const proteinChartRef = ref(null)
const fatChartRef = ref(null)
const carbChartRef = ref(null)
const previewVisible = ref(false)
const previewUrl = ref('')

async function loadSummary() {
  try {
    const params = {
      pet_id: summaryPetId.value,
      period: summaryPeriod.value
    }
    const res = await getDietSummary(params)
    summaryData.value = res.data || res
    renderCharts(summaryData.value)
  } catch (e) {
    Message.error('加载饮食摘要失败')
  }
}

function renderCharts(data) {
  // Simple text/div based charts since we don't have echarts loaded
  // In production, you would integrate with ECharts or similar
  if (caloriesChartRef.value) {
    caloriesChartRef.value.innerHTML = '<div style="display:flex;align-items:center;justify-content:center;height:100%;color:var(--color-text-3)">图表加载中...</div>'
  }
  if (mealTypeChartRef.value) {
    mealTypeChartRef.value.innerHTML = '<div style="display:flex;align-items:center;justify-content:center;height:100%;color:var(--color-text-3)">图表加载中...</div>'
  }
  if (proteinChartRef.value) {
    proteinChartRef.value.innerHTML = '<div style="display:flex;align-items:center;justify-content:center;height:100%;color:var(--color-text-3)">图表加载中...</div>'
  }
  if (fatChartRef.value) {
    fatChartRef.value.innerHTML = '<div style="display:flex;align-items:center;justify-content:center;height:100%;color:var(--color-text-3)">图表加载中...</div>'
  }
  if (carbChartRef.value) {
    carbChartRef.value.innerHTML = '<div style="display:flex;align-items:center;justify-content:center;height:100%;color:var(--color-text-3)">图表加载中...</div>'
  }
}

function previewPhoto(url) {
  previewUrl.value = url
  previewVisible.value = true
}

// ===== 工具函数 =====
function formatDateTime(dateStr) {
  if (!dateStr) return '—'
  const d = new Date(dateStr)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')} ${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`
}

function getMealTypeColor(type) {
  const colors = {
    breakfast: 'orange',
    lunch: 'blue',
    dinner: 'green',
    snack: 'purple'
  }
  return colors[type] || 'gray'
}

function getMealTypeLabel(type) {
  const labels = {
    breakfast: '早餐',
    lunch: '午餐',
    dinner: '晚餐',
    snack: '零食'
  }
  return labels[type] || type
}

onMounted(() => {
  loadRecords()
  if (pets.value.length > 0) {
    summaryPetId.value = pets.value[0].id
    loadSummary()
  }
})
</script>

<style scoped>
.food-name {
  font-weight: 500;
}
.food-brand {
  font-size: 12px;
  color: var(--color-text-3);
  margin-top: 2px;
}
.calories-value {
  font-weight: 600;
  color: #ff7d00;
}
.calories-unit {
  font-size: 12px;
  color: var(--color-text-3);
}
.food-thumbnail {
  width: 48px;
  height: 48px;
  object-fit: cover;
  border-radius: 6px;
  cursor: pointer;
  border: 1px solid var(--color-border);
}
.text-muted {
  color: var(--color-text-3);
}
</style>
