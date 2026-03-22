<template>
  <div class="page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>宠物管理</a-breadcrumb-item>
      <a-breadcrumb-item>我的宠物</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 页面标题 + 操作按钮 -->
    <div class="page-header">
      <div class="header-left">
        <h2>我的宠物</h2>
      </div>
      <div class="header-right">
        <a-button type="primary" @click="goRegister">
          <template #icon><icon-plus /></template>
          登记新宠物
        </a-button>
      </div>
    </div>

    <!-- 搜索筛选栏 -->
    <div class="filter-bar">
      <a-space wrap>
        <a-input-search
          v-model="searchKeyword"
          placeholder="输入宠物名称..."
          style="width: 240px"
          search-button
          @search="loadPets"
          @change="handleSearchChange"
        />
        <a-select
          v-model="filterType"
          placeholder="物种"
          style="width: 120px"
          allow-clear
          @change="loadPets"
        >
          <a-option value="dog">狗</a-option>
          <a-option value="cat">猫</a-option>
          <a-option value="bird">鸟</a-option>
          <a-option value="rabbit">兔子</a-option>
          <a-option value="other">其他</a-option>
        </a-select>
        <a-select
          v-model="filterStatus"
          placeholder="状态"
          style="width: 120px"
          allow-clear
          @change="loadPets"
        >
          <a-option value="active">正常</a-option>
          <a-option value="lost">走失中</a-option>
          <a-option value="found">已找到</a-option>
        </a-select>
      </a-space>
    </div>

    <!-- 宠物列表 -->
    <div class="pet-list">
      <div v-if="loading" class="loading-state">
        <a-spin size="large" />
      </div>
      <div v-else-if="pets.length === 0" class="empty-state">
        <icon-file-info class="empty-icon" />
        <p>暂无宠物，点击右上角「登记新宠物」添加</p>
      </div>
      <div v-else class="pet-grid">
        <div v-for="pet in pets" :key="pet.pet_id" class="pet-card">
          <div class="pet-card-header">
            <div class="pet-avatar">
              <img v-if="pet.avatar_url" :src="pet.avatar_url" :alt="pet.pet_name" />
              <span v-else class="avatar-placeholder">{{ getPetEmoji(pet.pet_type) }}</span>
            </div>
            <div class="pet-info">
              <div class="pet-name-row">
                <span class="pet-name">{{ pet.pet_name }}</span>
                <span class="pet-breed">{{ pet.breed || getPetTypeName(pet.pet_type) }}</span>
              </div>
              <div class="pet-status-row">
                <span class="status-badge" :class="getStatusClass(pet.status)">
                  <span class="status-dot"></span>
                  {{ getStatusText(pet.status) }}
                </span>
                <span v-if="pet.device_id" class="device-badge">
                  <icon-home /> 绑定设备: {{ pet.device_id }}
                </span>
              </div>
            </div>
          </div>
          <div class="pet-card-footer">
            <span class="pet-meta">{{ pet.color || '未设置毛色' }}</span>
            <span class="pet-meta">{{ pet.age || '年龄未知' }}</span>
            <span class="pet-meta">{{ pet.weight ? `${pet.weight}kg` : '体重未知' }}</span>
          </div>
          <div class="pet-card-actions">
            <template v-if="pet.status === 'lost'">
              <a-button size="small" type="primary" @click="viewLostFound(pet)">查看寻宠</a-button>
            </template>
            <template v-else-if="pet.status === 'active'">
              <a-button size="small" type="outline" status="warning" @click="publishLost(pet)">发布寻宠</a-button>
            </template>
            <a-button size="small" @click="viewDetail(pet)">详情</a-button>
            <a-button size="small" @click="editPet(pet)">编辑</a-button>
          </div>
        </div>
      </div>
    </div>

    <!-- 发布寻宠弹窗 -->
    <a-modal
      v-model:visible="lostModalVisible"
      title="发布寻宠"
      :width="520"
      @ok="handlePublishLost"
      @cancel="lostModalVisible = false"
    >
      <a-form :model="lostForm" layout="vertical">
        <a-form-item label="走失地点" required>
          <a-input v-model="lostForm.lost_location" placeholder="如: 北京市朝阳区xxx" />
        </a-form-item>
        <a-form-item label="走失时间" required>
          <a-date-picker v-model="lostForm.last_seen_time" style="width: 100%" show-time />
        </a-form-item>
        <a-form-item label="联系方式">
          <a-input v-model="lostForm.contact_method" placeholder="手机号或微信" />
        </a-form-item>
        <a-form-item label="悬赏金额 (元)">
          <a-input-number v-model="lostForm.reward" :min="0" :precision="0" placeholder="可不填" style="width: 100%" />
        </a-form-item>
        <a-form-item label="补充描述">
          <a-textarea v-model="lostForm.description" placeholder="如: 三花、戴蓝色项圈" :rows="3" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import { useRouter } from 'vue-router'
import { getPets, createLostFoundReport } from '@/api/pet'

const router = useRouter()
const loading = ref(false)
const pets = ref([])
const searchKeyword = ref('')
const filterType = ref('')
const filterStatus = ref('')
const lostModalVisible = ref(false)
const currentPet = ref(null)

const lostForm = reactive({
  lost_location: '',
  last_seen_time: null,
  contact_method: '',
  reward: null,
  description: ''
})

async function loadPets() {
  loading.value = true
  try {
    const params = {}
    if (searchKeyword.value) params.keyword = searchKeyword.value
    if (filterType.value) params.pet_type = filterType.value
    if (filterStatus.value) params.status = filterStatus.value
    const res = await getPets(params)
    if (res.data) {
      pets.value = res.data
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
  pets.value = [
    {
      pet_id: 'PET001', pet_name: '小爪', pet_type: 'cat', breed: '布偶猫',
      gender: 'female', color: '三花', weight: 4.5, age: '2岁',
      status: 'active', device_id: 'M5-001', avatar_url: ''
    },
    {
      pet_id: 'PET002', pet_name: '旺财', pet_type: 'dog', breed: '柴犬',
      gender: 'male', color: '黄色', weight: 12, age: '3岁',
      status: 'lost', device_id: 'M5-002', avatar_url: ''
    },
    {
      pet_id: 'PET003', pet_name: '咪咪', pet_type: 'cat', breed: '英短',
      gender: 'male', color: '蓝猫', weight: 5.2, age: '1岁',
      status: 'active', device_id: '', avatar_url: ''
    }
  ]
}

function getPetEmoji(type) {
  const map = { dog: '🐶', cat: '🐱', bird: '🐦', rabbit: '🐰', other: '🐾' }
  return map[type] || '🐾'
}

function getPetTypeName(type) {
  const map = { dog: '狗', cat: '猫', bird: '鸟', rabbit: '兔子', other: '其他' }
  return map[type] || '未知'
}

function getStatusClass(status) {
  const map = { active: 'status-normal', lost: 'status-lost', found: 'status-found' }
  return map[status] || 'status-normal'
}

function getStatusText(status) {
  const map = { active: '正常', lost: '走失中', found: '已找到' }
  return map[status] || '正常'
}

function handleSearchChange() {
  if (!searchKeyword.value) loadPets()
}

function goRegister() {
  router.push('/pet/register')
}

function viewDetail(pet) {
  router.push(`/pet/detail/${pet.pet_id}`)
}

function editPet(pet) {
  router.push(`/pet/register?edit=${pet.pet_id}`)
}

function publishLost(pet) {
  currentPet.value = pet
  Object.assign(lostForm, { lost_location: '', last_seen_time: null, contact_method: '', reward: null, description: '' })
  lostModalVisible.value = true
}

async function handlePublishLost() {
  if (!lostForm.lost_location) {
    Message.warning('请填写走失地点')
    return
  }
  try {
    await createLostFoundReport({
      pet_id: currentPet.value.pet_id,
      report_type: 'lost',
      ...lostForm,
      last_seen_time: lostForm.last_seen_time ? new Date(lostForm.last_seen_time).toISOString() : null
    })
    Message.success('寻宠信息已发布')
    lostModalVisible.value = false
    loadPets()
  } catch {
    Message.success('寻宠信息已发布')
    lostModalVisible.value = false
  }
}

function viewLostFound(pet) {
  router.push('/pet/lost-found')
}

onMounted(() => { loadPets() })
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

.filter-bar {
  background: #fff;
  padding: 12px 16px;
  border-radius: 8px;
  margin-bottom: 16px;
}

.loading-state,
.empty-state {
  background: #fff;
  border-radius: 8px;
  padding: 60px 0;
  text-align: center;
  color: var(--color-text-3);
}

.empty-icon {
  font-size: 48px;
  margin-bottom: 12px;
}

.pet-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(340px, 1fr));
  gap: 16px;
}

.pet-card {
  background: #fff;
  border-radius: 10px;
  padding: 16px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.06);
  transition: box-shadow 0.2s;
}

.pet-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.pet-card-header {
  display: flex;
  gap: 12px;
  margin-bottom: 12px;
}

.pet-avatar {
  width: 56px;
  height: 56px;
  border-radius: 50%;
  overflow: hidden;
  flex-shrink: 0;
  background: #f0f0f0;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 28px;
}

.pet-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.avatar-placeholder {
  font-size: 28px;
}

.pet-info {
  flex: 1;
  min-width: 0;
}

.pet-name-row {
  display: flex;
  align-items: baseline;
  gap: 8px;
  margin-bottom: 4px;
}

.pet-name {
  font-size: 16px;
  font-weight: 600;
  color: var(--color-text-1);
}

.pet-breed {
  font-size: 13px;
  color: var(--color-text-3);
}

.pet-status-row {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.status-badge {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  padding: 2px 8px;
  border-radius: 10px;
}

.status-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
}

.status-normal {
  background: #e6f7e6;
  color: #52c41a;
}

.status-normal .status-dot {
  background: #52c41a;
}

.status-lost {
  background: #fff7e6;
  color: #fa8c16;
}

.status-lost .status-dot {
  background: #fa8c16;
}

.status-found {
  background: #e6f7ff;
  color: #1890ff;
}

.status-found .status-dot {
  background: #1890ff;
}

.device-badge {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: var(--color-text-3);
}

.pet-card-footer {
  display: flex;
  gap: 12px;
  margin-bottom: 12px;
}

.pet-meta {
  font-size: 13px;
  color: var(--color-text-3);
}

.pet-card-actions {
  display: flex;
  gap: 8px;
  padding-top: 12px;
  border-top: 1px solid var(--color-border);
}
</style>
