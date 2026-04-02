<template>
    <Breadcrumb :items="['Home','Console','']" />


  <div class="page-container">
    <!-- 搜索筛选区 -->
    <div class="search-bar">
      <a-space>
        <a-select
          v-model="selectedPetId"
          placeholder="选择宠物"
          style="width: 200px"
          allow-search
          @change="handlePetChange"
        >
          <a-option v-for="pet in petList" :key="pet.device_id" :value="pet.device_id">
            {{ pet.pet_name }} ({{ pet.device_id }})
          </a-option>
        </a-select>
        <a-button @click="loadVitals">
          <template #icon><icon-refresh /></template>
          刷新
        </a-button>
      </a-space>
    </div>

    <!-- 操作按钮区 -->
    <div class="action-bar">
      <a-space>
        <a-button type="primary" @click="goToChart">查看实时曲线</a-button>
        <a-button @click="goToHistory">历史回放</a-button>
        <a-button @click="goToPrediction">行为预测</a-button>
      </a-space>
    </div>

    <!-- 四大指标卡片 -->
    <div class="vitals-cards">
      <a-row :gutter="16">
        <!-- 心率卡片 -->
        <a-col :span="6">
          <a-card class="vital-card" :class="{ 'abnormal': vitals.is_abnormal && vitals.abnormal_items.includes('heart_rate') }">
            <div class="vital-icon heart">
              <icon-heart-fill :size="32" />
            </div>
            <div class="vital-content">
              <div class="vital-label">心率</div>
              <div class="vital-value" :class="{ 'abnormal-value': vitals.is_abnormal && vitals.abnormal_items.includes('heart_rate') }">
                {{ vitals.heart_rate || '--' }}
                <span class="vital-unit">bpm</span>
              </div>
              <div class="vital-range">正常范围: 60-140</div>
            </div>
          </a-card>
        </a-col>

        <!-- 呼吸频率卡片 -->
        <a-col :span="6">
          <a-card class="vital-card" :class="{ 'abnormal': vitals.is_abnormal && vitals.abnormal_items.includes('respiratory_rate') }">
            <div class="vital-icon respiratory">
              <icon-wind :size="32" />
            </div>
            <div class="vital-content">
              <div class="vital-label">呼吸频率</div>
              <div class="vital-value" :class="{ 'abnormal-value': vitals.is_abnormal && vitals.abnormal_items.includes('respiratory_rate') }">
                {{ vitals.respiratory_rate || '--' }}
                <span class="vital-unit">次/分</span>
              </div>
              <div class="vital-range">正常范围: 10-30</div>
            </div>
          </a-card>
        </a-col>

        <!-- 体温卡片 -->
        <a-col :span="6">
          <a-card class="vital-card" :class="{ 'abnormal': vitals.is_abnormal && vitals.abnormal_items.includes('body_temp') }">
            <div class="vital-icon temp">
              <icon-temperature :size="32" />
            </div>
            <div class="vital-content">
              <div class="vital-label">体温</div>
              <div class="vital-value" :class="{ 'abnormal-value': vitals.is_abnormal && vitals.abnormal_items.includes('body_temp') }">
                {{ vitals.body_temp != null ? vitals.body_temp.toFixed(1) : '--' }}
                <span class="vital-unit">℃</span>
              </div>
              <div class="vital-range">正常范围: 38-39.5</div>
            </div>
          </a-card>
        </a-col>

        <!-- 活动量卡片 -->
        <a-col :span="6">
          <a-card class="vital-card">
            <div class="vital-icon activity">
              <icon-activity :size="32" />
            </div>
            <div class="vital-content">
              <div class="vital-label">活动量</div>
              <div class="vital-value">
                {{ vitals.activity_level != null ? vitals.activity_level : '--' }}
                <span class="vital-unit">分</span>
              </div>
              <div class="vital-progress">
                <a-progress
                  :percent="vitals.activity_level || 0"
                  :stroke-width="8"
                  :show-text="false"
                  :color="getActivityColor(vitals.activity_level)"
                />
              </div>
            </div>
          </a-card>
        </a-col>
      </a-row>
    </div>

    <!-- 异常告警 -->
    <div v-if="vitals.is_abnormal" class="abnormal-alert">
      <a-alert type="error" :show-icon="true">
        <template #title>检测到异常指标</template>
        <template #content>
          <span>异常项目: {{ vitals.abnormal_items.join(', ') }}</span>
          <span style="margin-left: 16px">更新时间: {{ vitals.timestamp }}</span>
        </template>
      </a-alert>
    </div>

    <!-- 详细数据表格 -->
    <div class="content-area">
      <a-card title="详细数据">
        <a-table
          :columns="columns"
          :data="[vitals]"
          :loading="loading"
          :pagination="false"
          row-key="device_id"
        >
          <template #is_abnormal="{ record }">
            <a-badge
              :status="record.is_abnormal ? 'error' : 'success'"
              :text="record.is_abnormal ? '异常' : '正常'"
            />
          </template>
          <template #timestamp="{ record }">
            {{ record.timestamp || '--' }}
          </template>
        </a-table>
      </a-card>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { Message } from '@arco-design/web-vue'
import axios from 'axios'

const router = useRouter()

const API_BASE = '/api/digital-twin'

const petList = ref([])
const selectedPetId = ref('')
const loading = ref(false)
const refreshTimer = ref(null)

const vitals = reactive({
  device_id: '',
  heart_rate: null,
  respiratory_rate: null,
  body_temp: null,
  activity_level: null,
  is_abnormal: false,
  abnormal_items: [],
  timestamp: ''
})

const columns = [
  { title: '设备ID', dataIndex: 'device_id', ellipsis: true },
  { title: '心率 (bpm)', dataIndex: 'heart_rate' },
  { title: '呼吸 (次/分)', dataIndex: 'respiratory_rate' },
  { title: '体温 (℃)', dataIndex: 'body_temp' },
  { title: '活动量', dataIndex: 'activity_level' },
  { title: '状态', slotName: 'is_abnormal' },
  { title: '更新时间', slotName: 'timestamp' }
]

// 获取宠物列表
const loadPets = async () => {
  try {
    const res = await axios.get(`${API_BASE}/pets`, {
      headers: { Authorization: `Bearer ${localStorage.getItem('token')}` }
    })
    if (res.data.code === 0 || res.data.code === 200) {
      petList.value = res.data.data?.list || res.data.data || []
    }
  } catch {
    // 模拟数据
    petList.value = [
      { device_id: 'PET001', pet_name: '小白' },
      { device_id: 'PET002', pet_name: '旺财' }
    ]
  }
  if (petList.value.length > 0 && !selectedPetId.value) {
    selectedPetId.value = petList.value[0].device_id
    loadVitals()
  }
}

// 加载生命体征数据
const loadVitals = async () => {
  if (!selectedPetId.value) {
    Message.warning('请先选择宠物')
    return
  }
  loading.value = true
  try {
    const res = await axios.get(`${API_BASE}/vitals/current/${selectedPetId.value}`, {
      headers: { Authorization: `Bearer ${localStorage.getItem('token')}` }
    })
    if (res.data.code === 0 || res.data.code === 200) {
      const data = res.data.data
      Object.assign(vitals, data)
    }
  } catch {
    // 模拟数据
    const mockVitals = {
      device_id: selectedPetId.value,
      heart_rate: Math.floor(Math.random() * 60) + 70,
      respiratory_rate: Math.floor(Math.random() * 15) + 15,
      body_temp: (Math.random() * 1.5 + 38).toFixed(1),
      activity_level: Math.floor(Math.random() * 100),
      is_abnormal: Math.random() > 0.8,
      abnormal_items: Math.random() > 0.8 ? ['heart_rate'] : [],
      timestamp: new Date().toLocaleString('zh-CN')
    }
    Object.assign(vitals, mockVitals)
    Message.warning('使用模拟数据')
  } finally {
    loading.value = false
  }
}

// 自动刷新
const startAutoRefresh = () => {
  refreshTimer.value = setInterval(() => {
    loadVitals()
  }, 30000) // 每30秒刷新
}

// 宠物切换
const handlePetChange = () => {
  loadVitals()
}

// 跳转实时曲线
const goToChart = () => {
  if (selectedPetId.value) {
    router.push(`/digital-twin/vitals-chart?device=${selectedPetId.value}`)
  } else {
    router.push('/digital-twin/vitals-chart')
  }
}

// 跳转历史回放
const goToHistory = () => {
  if (selectedPetId.value) {
    router.push(`/digital-twin/history?device=${selectedPetId.value}`)
  } else {
    router.push('/digital-twin/history')
  }
}

// 跳转行为预测
const goToPrediction = () => {
  if (selectedPetId.value) {
    router.push(`/digital-twin/behavior?device=${selectedPetId.value}`)
  } else {
    router.push('/digital-twin/behavior')
  }
}

// 获取活动量颜色
const getActivityColor = (level) => {
  if (level == null) return '#1650ff'
  if (level < 30) return '#00b42a'
  if (level < 70) return '#1650ff'
  return '#ff7d00'
}

onMounted(() => {
  loadPets()
  startAutoRefresh()
})

onUnmounted(() => {
  if (refreshTimer.value) {
    clearInterval(refreshTimer.value)
  }
})
</script>

<style scoped>
.page-container {
  padding: 16px;
}

.search-bar {
  margin-bottom: 12px;
}

.action-bar {
  margin-bottom: 16px;
}

.vitals-cards {
  margin-bottom: 16px;
}

.vital-card {
  border-radius: 8px;
  transition: all 0.3s;
}

.vital-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.vital-card.abnormal {
  border: 2px solid #ff7d00;
  background: linear-gradient(135deg, #fff 0%, #fff5f0 100%);
}

.vital-card :deep(.arco-card-body) {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 20px;
}

.vital-icon {
  width: 56px;
  height: 56px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
}

.vital-icon.heart {
  background: linear-gradient(135deg, #ff7d00 0%, #f53f3f 100%);
}

.vital-icon.respiratory {
  background: linear-gradient(135deg, #00b42a 0%, #1650ff 100%);
}

.vital-icon.temp {
  background: linear-gradient(135deg, #f53f3f 0%, #ff7d00 100%);
}

.vital-icon.activity {
  background: linear-gradient(135deg, #1650ff 0%, #722ed1 100%);
}

.vital-content {
  flex: 1;
}

.vital-label {
  font-size: 14px;
  color: #86909c;
  margin-bottom: 4px;
}

.vital-value {
  font-size: 28px;
  font-weight: 600;
  color: #1d2129;
  line-height: 1.2;
}

.vital-value.abnormal-value {
  color: #f53f3f;
}

.vital-unit {
  font-size: 14px;
  font-weight: 400;
  color: #86909c;
  margin-left: 4px;
}

.vital-range {
  font-size: 12px;
  color: #86909c;
  margin-top: 4px;
}

.vital-progress {
  margin-top: 8px;
  width: 100%;
}

.abnormal-alert {
  margin-bottom: 16px;
}

.content-area {
  background: #fff;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
}
</style>
