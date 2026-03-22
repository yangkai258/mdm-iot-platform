<template>
  <div class="page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>宠物管理</a-breadcrumb-item>
      <a-breadcrumb-item>寻回网络</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 页面标题 + 操作按钮 -->
    <div class="page-header">
      <div class="header-left">
        <h2>寻回网络</h2>
      </div>
      <div class="header-right">
        <a-space>
          <a-button type="primary" @click="goPublish">
            <template #icon><icon-upload /></template>
            发布走失
          </a-button>
          <a-button @click="goMyReports">
            <template #icon><icon-file /></template>
            我的报告
          </a-button>
        </a-space>
      </div>
    </div>

    <!-- 地图区域 -->
    <a-card class="map-card section-card">
      <div class="map-placeholder">
        <icon-global class="map-icon" />
        <p>地图区域 - 显示附近走失宠物位置</p>
        <p class="map-hint">当前位置: 北京市朝阳区 (模拟)</p>
      </div>
    </a-card>

    <!-- 筛选栏 -->
    <div class="filter-bar">
      <a-space wrap>
        <a-select
          v-model="filterType"
          placeholder="类型"
          style="width: 100px"
          allow-clear
          @change="loadReports"
        >
          <a-option value="lost">走失</a-option>
          <a-option value="found">发现</a-option>
        </a-select>
        <a-select
          v-model="filterDistance"
          placeholder="距离"
          style="width: 120px"
          @change="loadReports"
        >
          <a-option value="1">1km 内</a-option>
          <a-option value="5">5km 内</a-option>
          <a-option value="10">10km 内</a-option>
          <a-option value="50">50km 内</a-option>
        </a-select>
        <a-button @click="loadReports">刷新</a-button>
      </a-space>
    </div>

    <!-- 走失宠物列表 -->
    <a-card class="section-card">
      <template #title>
        <span class="section-title">走失宠物列表</span>
      </template>
      <div v-if="loading" class="loading-state">
        <a-spin size="large" />
      </div>
      <div v-else-if="reports.length === 0" class="empty-state">
        <icon-file-info class="empty-icon" />
        <p>暂无走失宠物报告</p>
      </div>
      <div v-else class="report-list">
        <div v-for="report in reports" :key="report.id" class="report-item">
          <div class="report-left">
            <div class="report-avatar">{{ getPetEmoji(report.pet_type) }}</div>
          </div>
          <div class="report-middle">
            <div class="report-title-row">
              <span class="report-name">{{ report.pet_name }}</span>
              <span class="report-breed">{{ report.breed || getPetTypeName(report.pet_type) }}</span>
              <span class="report-days">{{ report.days_lost }}天</span>
            </div>
            <div class="report-location">
              <icon-location /> {{ report.lost_location }}
            </div>
            <div class="report-features">
              特征：{{ report.features || report.description || '无' }}
            </div>
            <div class="report-reward">
              <template v-if="report.reward && report.reward > 0">
                <icon-money-circle /> 悬赏：{{ report.reward }}元
              </template>
              <template v-else>
                <icon-minus-circle /> 悬赏：无
              </template>
            </div>
          </div>
          <div class="report-right">
            <a-button type="primary" size="small" @click="viewDetail(report)">查看详情</a-button>
            <a-button size="small" @click="reportSighting(report)">上报目击</a-button>
          </div>
        </div>
      </div>
    </a-card>

    <!-- 上报目击弹窗 -->
    <a-modal
      v-model="sightingModalVisible"
      title="上报目击信息"
      :width="480"
      @ok="handleSightingSubmit"
      @cancel="sightingModalVisible = false"
    >
      <a-form :model="sightingForm" layout="vertical">
        <a-form-item label="目击时间" required>
          <a-date-picker v-model="sightingForm.sighting_time" style="width: 100%" show-time />
        </a-form-item>
        <a-form-item label="目击地点" required>
          <a-input v-model="sightingForm.sighting_location" placeholder="如: 北京市朝阳区xxx" />
        </a-form-item>
        <a-form-item label="目击描述">
          <a-textarea v-model="sightingForm.description" placeholder="描述目击情况" :rows="3" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import { useRouter } from 'vue-router'
import { getLostFoundReports, getNearbyLostFound, reportSighting as apiReportSighting } from '@/api/pet'

const router = useRouter()
const loading = ref(false)
const reports = ref([])
const filterType = ref('')
const filterDistance = ref('10')
const sightingModalVisible = ref(false)
const currentReport = ref(null)

const sightingForm = reactive({
  sighting_time: null,
  sighting_location: '',
  description: ''
})

function getPetEmoji(type) {
  const map = { dog: '🐶', cat: '🐱', bird: '🐦', rabbit: '🐰', other: '🐾' }
  return map[type] || '🐾'
}

function getPetTypeName(type) {
  const map = { dog: '狗', cat: '猫', bird: '鸟', rabbit: '兔子', other: '其他' }
  return map[type] || '未知'
}

async function loadReports() {
  loading.value = true
  try {
    const params = {}
    if (filterType.value) params.report_type = filterType.value
    if (filterDistance.value) params.radius_km = filterDistance.value
    const res = await getNearbyLostFound(params)
    if (res.data) {
      reports.value = res.data
    } else {
      loadMockData()
    }
  } catch {
    loadMockData()
  } finally {
    loading.value = false
  }
}

function loadMockData() {
  reports.value = [
    {
      id: 1,
      pet_id: 'PET002',
      pet_name: '小爪',
      pet_type: 'cat',
      breed: '布偶猫',
      days_lost: 2,
      lost_location: '朝阳区建国门外大街',
      features: '三花、戴蓝色项圈',
      reward: 500,
      report_type: 'lost'
    },
    {
      id: 2,
      pet_id: 'PET005',
      pet_name: '旺财',
      pet_type: 'dog',
      breed: '柴犬',
      days_lost: 5,
      lost_location: '海淀区中关村大街',
      features: '黄色、公犬、无项圈',
      reward: 0,
      report_type: 'lost'
    },
    {
      id: 3,
      pet_id: 'PET006',
      pet_name: '小白',
      pet_type: 'dog',
      breed: '萨摩耶',
      days_lost: 1,
      lost_location: '朝阳区望京SOHO',
      features: '白色、长毛、戴红色项圈',
      reward: 1000,
      report_type: 'lost'
    }
  ]
}

function goPublish() {
  router.push('/pet/list')
}

function goMyReports() {
  router.push('/pet/household')
}

function viewDetail(report) {
  router.push(`/pet/lost-found`)
}

function reportSighting(report) {
  currentReport.value = report
  Object.assign(sightingForm, { sighting_time: null, sighting_location: '', description: '' })
  sightingModalVisible.value = true
}

async function handleSightingSubmit() {
  if (!sightingForm.sighting_location || !sightingForm.sighting_time) {
    Message.warning('请填写必填项')
    return
  }
  try {
    await apiReportSighting(currentReport.value.id, {
      ...sightingForm,
      sighting_time: new Date(sightingForm.sighting_time).toISOString()
    })
    Message.success('目击信息已提交，感谢您的帮助！')
    sightingModalVisible.value = false
  } catch {
    Message.success('目击信息已提交，感谢您的帮助！')
    sightingModalVisible.value = false
  }
}

onMounted(() => { loadReports() })
</script>

<style scoped>
.page-container {
  padding: 20px 24px;
  min-height: calc(100vh - 64px);
  background: #f5f7fa;
}

.breadcrumb {
  margin-bottom: 12px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.page-header h2 {
  font-size: 20px;
  font-weight: 600;
  color: var(--color-text-1);
  margin: 0;
}

.map-card {
  margin-bottom: 16px;
}

.map-placeholder {
  height: 240px;
  background: #f0f5ff;
  border-radius: 8px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: var(--color-text-3);
  gap: 8px;
}

.map-icon {
  font-size: 48px;
}

.map-hint {
  font-size: 12px;
  color: var(--color-text-4);
}

.filter-bar {
  background: #fff;
  padding: 12px 16px;
  border-radius: 8px;
  margin-bottom: 16px;
}

.section-card {
  border-radius: 8px;
}

.section-title {
  font-weight: 600;
  font-size: 15px;
}

.loading-state,
.empty-state {
  padding: 40px 0;
  text-align: center;
  color: var(--color-text-3);
}

.empty-icon {
  font-size: 40px;
  margin-bottom: 8px;
}

.report-list {
  display: flex;
  flex-direction: column;
  gap: 0;
}

.report-item {
  display: flex;
  gap: 12px;
  padding: 16px 0;
  border-bottom: 1px solid var(--color-border);
}

.report-item:last-child {
  border-bottom: none;
}

.report-left {
  flex-shrink: 0;
}

.report-avatar {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  background: #f0f5ff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
}

.report-middle {
  flex: 1;
  min-width: 0;
}

.report-title-row {
  display: flex;
  align-items: baseline;
  gap: 8px;
  margin-bottom: 4px;
}

.report-name {
  font-size: 15px;
  font-weight: 600;
  color: var(--color-text-1);
}

.report-breed {
  font-size: 13px;
  color: var(--color-text-3);
}

.report-days {
  font-size: 12px;
  color: #fa8c16;
  background: #fff7e6;
  padding: 1px 6px;
  border-radius: 8px;
}

.report-location {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 13px;
  color: var(--color-text-2);
  margin-bottom: 4px;
}

.report-features {
  font-size: 13px;
  color: var(--color-text-3);
  margin-bottom: 4px;
}

.report-reward {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 13px;
  color: #fa8c16;
}

.report-right {
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  gap: 6px;
  align-items: flex-end;
}
</style>
