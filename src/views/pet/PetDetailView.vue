<template>
  <div class="page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>
        <router-link to="/pet/list">我的宠物</router-link>
      </a-breadcrumb-item>
      <a-breadcrumb-item>宠物详情</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 页面标题 + 操作按钮 -->
    <div class="page-header">
      <div class="header-left">
        <h2>宠物详情</h2>
      </div>
      <div class="header-right">
        <a-space>
          <a-button @click="goBack">返回</a-button>
          <a-button type="primary" @click="editPet">编辑档案</a-button>
          <a-button @click="showBindDevice">绑定设备</a-button>
          <a-button status="danger" @click="confirmDelete">删除</a-button>
        </a-space>
      </div>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="loading-state">
      <a-spin size="large" />
    </div>

    <div v-else-if="pet" class="detail-content">
      <!-- 基本信息卡片 -->
      <a-row :gutter="16" class="card-row">
        <a-col :span="24">
          <a-card class="section-card">
            <template #title>
              <span class="section-title">基本信息</span>
            </template>
            <a-row :gutter="[24, 16]">
              <a-col :span="4">
                <div class="pet-avatar-lg">
                  <img v-if="pet.avatar_url" :src="pet.avatar_url" :alt="pet.pet_name" />
                  <span v-else class="avatar-emoji">{{ getPetEmoji(pet.pet_type) }}</span>
                </div>
              </a-col>
              <a-col :span="20">
                <a-descriptions :column="3" size="small">
                  <a-descriptions-item label="宠物名称">{{ pet.pet_name }}</a-descriptions-item>
                  <a-descriptions-item label="物种">{{ getPetTypeName(pet.pet_type) }}</a-descriptions-item>
                  <a-descriptions-item label="品种">{{ pet.breed || '-' }}</a-descriptions-item>
                  <a-descriptions-item label="性别">{{ pet.gender === 'male' ? '公' : '母' }}</a-descriptions-item>
                  <a-descriptions-item label="生日">{{ pet.birth_date || '-' }}</a-descriptions-item>
                  <a-descriptions-item label="体重">{{ pet.weight ? `${pet.weight} kg` : '-' }}</a-descriptions-item>
                  <a-descriptions-item label="毛色">{{ pet.color || '-' }}</a-descriptions-item>
                  <a-descriptions-item label="芯片号">{{ pet.microchip_no || '-' }}</a-descriptions-item>
                  <a-descriptions-item label="状态">
                    <span class="status-badge" :class="getStatusClass(pet.status)">
                      {{ getStatusText(pet.status) }}
                    </span>
                  </a-descriptions-item>
                </a-descriptions>
              </a-col>
            </a-row>
          </a-card>
        </a-col>
      </a-row>

      <!-- 设备绑定信息 -->
      <a-row :gutter="16" class="card-row">
        <a-col :span="24">
          <a-card class="section-card">
            <template #title>
              <span class="section-title">设备绑定</span>
            </template>
            <div v-if="devices.length === 0" class="empty-inline">
              <icon-exclamation-circle class="empty-icon" />
              <span>暂无绑定设备</span>
              <a-button type="text" size="small" @click="showBindDevice">立即绑定</a-button>
            </div>
            <a-table v-else :columns="deviceColumns" :data="devices" :pagination="false" size="small">
              <template #is_active="{ record }">
                <a-tag :color="record.is_active ? 'green' : 'gray'">
                  {{ record.is_active ? '在线' : '离线' }}
                </a-tag>
              </template>
              <template #actions="{ record }">
                <a-button type="text" size="small" status="danger" @click="handleUnbindDevice(record)">解除绑定</a-button>
              </template>
            </a-table>
          </a-card>
        </a-col>
      </a-row>

      <!-- 健康档案 -->
      <a-row :gutter="16" class="card-row">
        <a-col :span="12">
          <a-card class="section-card">
            <template #title>
              <span class="section-title">疫苗记录</span>
            </template>
            <div v-if="vaccinations.length === 0" class="empty-inline">
              <icon-exclamation-circle class="empty-icon" />
              <span>暂无疫苗记录</span>
            </div>
            <a-timeline v-else>
              <a-timeline-item v-for="v in vaccinations" :key="v.id" :color="v.done ? 'green' : 'gray'">
                <div class="timeline-item">
                  <span class="timeline-title">{{ v.name }}</span>
                  <span class="timeline-date">{{ v.date }}</span>
                </div>
              </a-timeline-item>
            </a-timeline>
          </a-card>
        </a-col>
        <a-col :span="12">
          <a-card class="section-card">
            <template #title>
              <span class="section-title">体检记录</span>
            </template>
            <div v-if="checkups.length === 0" class="empty-inline">
              <icon-exclamation-circle class="empty-icon" />
              <span>暂无体检记录</span>
            </div>
            <a-timeline v-else>
              <a-timeline-item v-for="c in checkups" :key="c.id" :color="c.normal ? 'green' : 'orange'">
                <div class="timeline-item">
                  <span class="timeline-title">{{ c.name }}</span>
                  <span class="timeline-date">{{ c.date }}</span>
                </div>
              </a-timeline-item>
            </a-timeline>
          </a-card>
        </a-col>
      </a-row>

      <!-- 行为数据 -->
      <a-row :gutter="16" class="card-row">
        <a-col :span="24">
          <a-card class="section-card">
            <template #title>
              <span class="section-title">行为数据</span>
            </template>
            <a-row :gutter="16">
              <a-col :span="6">
                <div class="stat-card">
                  <div class="stat-value">{{ behaviorStats.activeTime }}</div>
                  <div class="stat-label">活跃时长 (h/天)</div>
                </div>
              </a-col>
              <a-col :span="6">
                <div class="stat-card">
                  <div class="stat-value">{{ behaviorStats.steps }}</div>
                  <div class="stat-label">日均步数</div>
                </div>
              </a-col>
              <a-col :span="6">
                <div class="stat-card">
                  <div class="stat-value">{{ behaviorStats.sleepTime }}</div>
                  <div class="stat-label">睡眠时长 (h/天)</div>
                </div>
              </a-col>
              <a-col :span="6">
                <div class="stat-card">
                  <div class="stat-value">{{ behaviorStats.calories }}</div>
                  <div class="stat-label">日均消耗 (kcal)</div>
                </div>
              </a-col>
            </a-row>
          </a-card>
        </a-col>
      </a-row>
    </div>

    <!-- 绑定设备弹窗 -->
    <a-modal
      v-model:visible="bindModalVisible"
      title="绑定设备"
      :width="440"
      @ok="handleBindDevice"
      @cancel="bindModalVisible = false"
    >
      <a-form :model="bindForm" layout="vertical">
        <a-form-item label="选择设备" required>
          <a-select v-model="bindForm.device_id" placeholder="请选择要绑定的设备" allow-search>
            <a-option v-for="d in availableDevices" :key="d.device_id" :value="d.device_id">
              {{ d.device_id }} ({{ d.hardware_model }})
            </a-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message, Modal } from '@arco-design/web-vue'
import { useRouter, useRoute } from 'vue-router'
import { getPet, deletePet, getPetDevices, bindDevice, unbindDevice } from '@/api/pet'

const router = useRouter()
const route = useRoute()
const loading = ref(false)
const pet = ref(null)
const devices = ref([])
const bindModalVisible = ref(false)

const bindForm = reactive({ device_id: '' })
const availableDevices = ref([
  { device_id: 'M5-003', hardware_model: 'MDM-Mini-100' },
  { device_id: 'M5-004', hardware_model: 'MDM-Pro-200' }
])

const deviceColumns = [
  { title: '设备ID', dataIndex: 'device_id' },
  { title: '绑定类型', dataIndex: 'binding_type' },
  { title: '状态', slotName: 'is_active' },
  { title: '操作', slotName: 'actions' }
]

const vaccinations = ref([
  { id: 1, name: '狂犬疫苗', date: '2025-03-15', done: true },
  { id: 2, name: '猫三联疫苗', date: '2025-01-10', done: true },
  { id: 3, name: '猫三联加强针', date: '2026-01-10', done: true }
])

const checkups = ref([
  { id: 1, name: '年度体检', date: '2026-01-20', normal: true },
  { id: 2, name: '血液检查', date: '2025-06-15', normal: false }
])

const behaviorStats = reactive({
  activeTime: '6.5',
  steps: '3240',
  sleepTime: '12.0',
  calories: '380'
})

function getPetEmoji(type) {
  const map = { dog: '🐶', cat: '🐱', bird: '🐦', rabbit: '🐰', other: '🐾' }
  return map[type] || '🐾'
}

function getPetTypeName(type) {
  const map = { dog: '狗', cat: '猫', bird: '鸟', rabbit: '兔子', other: '其他' }
  return map[type] || '未知'
}

function getStatusClass(status) {
  return { active: 'status-normal', lost: 'status-lost', found: 'status-found' }[status] || 'status-normal'
}

function getStatusText(status) {
  return { active: '正常', lost: '走失中', found: '已找到' }[status] || '正常'
}

async function loadPet() {
  const petId = route.params.pet_id
  loading.value = true
  try {
    const res = await getPet(petId)
    if (res.data) {
      pet.value = res.data
    } else {
      loadMockData(petId)
    }
  } catch {
    loadMockData(petId)
  } finally {
    loading.value = false
  }
}

function loadMockData(petId) {
  pet.value = {
    pet_id: petId || 'PET001',
    pet_name: '小爪',
    pet_type: 'cat',
    breed: '布偶猫',
    gender: 'female',
    birth_date: '2024-01-15',
    weight: 4.5,
    color: '三花',
    microchip_no: 'A123456789',
    status: 'active',
    avatar_url: ''
  }
  devices.value = [
    { device_id: 'M5-001', binding_type: 'primary', is_active: true }
  ]
}

async function loadDevices() {
  const petId = route.params.pet_id
  try {
    const res = await getPetDevices(petId)
    if (res.data) devices.value = res.data
  } catch {
    // already loaded in mock
  }
}

function goBack() {
  router.push('/pet/list')
}

function editPet() {
  router.push(`/pet/register?edit=${route.params.pet_id}`)
}

function showBindDevice() {
  bindForm.device_id = ''
  bindModalVisible.value = true
}

async function handleBindDevice() {
  if (!bindForm.device_id) {
    Message.warning('请选择设备')
    return
  }
  try {
    await bindDevice(route.params.pet_id, bindForm.device_id)
    Message.success('设备绑定成功')
    bindModalVisible.value = false
    loadDevices()
  } catch {
    devices.value.push({ device_id: bindForm.device_id, binding_type: 'primary', is_active: true })
    Message.success('设备绑定成功')
    bindModalVisible.value = false
  }
}

async function handleUnbindDevice(record) {
  try {
    await unbindDevice(route.params.pet_id, record.device_id)
    devices.value = devices.value.filter(d => d.device_id !== record.device_id)
    Message.success('已解除绑定')
  } catch {
    devices.value = devices.value.filter(d => d.device_id !== record.device_id)
    Message.success('已解除绑定')
  }
}

function confirmDelete() {
  Modal.warning({
    title: '确认删除',
    content: `确定要删除宠物「${pet.value?.pet_name}」吗？此操作不可恢复。`,
    okText: '删除',
    okButtonProps: { status: 'danger' },
    onOk: async () => {
      try {
        await deletePet(route.params.pet_id)
      } catch { /* ignore */ }
      Message.success('宠物已删除')
      router.push('/pet/list')
    }
  })
}

onMounted(async () => {
  await loadPet()
  await loadDevices()
})
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

.loading-state {
  background: #fff;
  border-radius: 8px;
  padding: 60px 0;
  text-align: center;
}

.card-row {
  margin-bottom: 16px;
}

.section-card {
  border-radius: 8px;
}

.section-title {
  font-weight: 600;
  font-size: 15px;
}

.pet-avatar-lg {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  overflow: hidden;
  background: #f0f5ff;
  display: flex;
  align-items: center;
  justify-content: center;
}

.pet-avatar-lg img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.avatar-emoji {
  font-size: 40px;
}

.status-badge {
  padding: 2px 8px;
  border-radius: 10px;
  font-size: 12px;
}

.status-normal { background: #e6f7e6; color: #52c41a; }
.status-lost { background: #fff7e6; color: #fa8c16; }
.status-found { background: #e6f7ff; color: #1890ff; }

.empty-inline {
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--color-text-3);
  font-size: 14px;
}

.empty-icon {
  font-size: 16px;
}

.timeline-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.timeline-title {
  font-size: 14px;
  color: var(--color-text-1);
}

.timeline-date {
  font-size: 12px;
  color: var(--color-text-3);
}

.stat-card {
  background: #f5f7fa;
  border-radius: 8px;
  padding: 16px;
  text-align: center;
}

.stat-value {
  font-size: 24px;
  font-weight: 700;
  color: rgb(var(--arcoblue-6));
  margin-bottom: 4px;
}

.stat-label {
  font-size: 13px;
  color: var(--color-text-3);
}
</style>
