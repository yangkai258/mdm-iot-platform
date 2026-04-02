<template>
    <Breadcrumb :items="['Home','Console','']" />


  <div class="pro-page-container">
    <!-- 操作栏 -->
    <div class="pro-action-bar">
      <a-space>
        <a-button type="primary" @click="openCreateModal"><icon-plus />创建虚拟宠物</a-button>
        <a-button @click="loadPets"><icon-refresh />刷新</a-button>
      </a-space>
    </div>

    <!-- 虚拟宠物列表 -->
    <div class="pro-content-area">
      <a-row :gutter="[16, 16]" v-if="pets.length">
        <a-col :xs="24" :sm="12" :md="8" :lg="6" v-for="pet in pets" :key="pet.id">
          <a-card class="pet-card">
            <template #title>
              <div class="pet-header">
                <span>{{ pet.pet_name }}</span>
                <a-tag :color="getPetTypeColor(pet.pet_type)" size="small">{{ getPetTypeName(pet.pet_type) }}</a-tag>
              </div>
            </template>
            <template #extra>
              <a-tag :color="getStatusColor(pet.status)" size="small">
                {{ getStatusName(pet.status) }}
              </a-tag>
            </template>
            <div class="pet-avatar">
              <div class="pet-avatar-placeholder">
                <span>{{ getPetEmoji(pet.pet_type) }}</span>
              </div>
            </div>
            <a-descriptions :column="1" size="small">
              <a-descriptions-item label="当前情绪">
                <a-tag v-if="pet.current_emotion">{{ getEmotionName(pet.current_emotion) }}</a-tag>
                <span v-else class="text-muted">无</span>
              </a-descriptions-item>
              <a-descriptions-item label="电量">
                <a-progress :percent="(pet.battery_level || 0) * 100" size="small" :color="getBatteryColor(pet.battery_level)" />
              </a-descriptions-item>
            </a-descriptions>
            <template #actions>
              <a-button type="text" size="small" :disabled="pet.status === 'running'" @click="startPet(pet)">启动</a-button>
              <a-button type="text" size="small" :disabled="pet.status !== 'running'" @click="stopPet(pet)">停止</a-button>
              <a-button type="text" size="small" @click="openInteractModal(pet)">交互</a-button>
            </template>
          </a-card>
        </a-col>
      </a-row>
      <a-empty v-else description="暂无虚拟宠物" />

      <div class="pro-pagination" v-if="total > 0">
        <a-pagination :total="total" :current="page" :page-size="pageSize" show-total @page-size-change="onPageSizeChange" @change="onPageChange" />
      </div>
    </div>

    <!-- 创建虚拟宠物弹窗 -->
    <a-modal v-model:visible="createModalVisible" title="创建虚拟宠物" @before-ok="handleCreate">
      <a-form :model="createForm" layout="vertical" ref="formRef">
        <a-form-item label="宠物名称" required>
          <a-input v-model="createForm.pet_name" placeholder="请输入宠物名称" />
        </a-form-item>
        <a-form-item label="宠物类型" required>
          <a-select v-model="createForm.pet_type" placeholder="请选择宠物类型">
            <a-option value="cat">猫</a-option>
            <a-option value="dog">狗</a-option>
            <a-option value="rabbit">兔子</a-option>
          </a-select>
        </a-form-item>
        <a-divider>性格配置</a-divider>
        <a-form-item label="友善度">
          <a-slider v-model="createForm.personality.friendliness" :min="0" :max="1" :step="0.1" show-input />
        </a-form-item>
        <a-form-item label="活泼度">
          <a-slider v-model="createForm.personality.playfulness" :min="0" :max="1" :step="0.1" show-input />
        </a-form-item>
        <a-form-item label="精力水平">
          <a-slider v-model="createForm.personality.energy_level" :min="0" :max="1" :step="0.1" show-input />
        </a-form-item>
        <a-divider>能力配置</a-divider>
        <a-form-item label="语音识别">
          <a-switch v-model="createForm.capabilities.speech_recognition" />
        </a-form-item>
        <a-form-item label="情绪识别">
          <a-switch v-model="createForm.capabilities.emotion_recognition" />
        </a-form-item>
        <a-form-item label="路径规划">
          <a-switch v-model="createForm.capabilities.path_planning" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 交互弹窗 -->
    <a-modal v-model:visible="interactModalVisible" title="虚拟交互" :width="480" @before-ok="handleInteract">
      <div v-if="currentPet" class="interact-content">
        <div class="interact-pet-info">
          <span class="pet-emoji">{{ getPetEmoji(currentPet.pet_type) }}</span>
          <span>{{ currentPet.pet_name }}</span>
          <a-tag v-if="currentPet.current_emotion">{{ getEmotionName(currentPet.current_emotion) }}</a-tag>
        </div>
        <a-divider>交互类型</a-divider>
        <a-form :model="interactForm" layout="vertical">
          <a-form-item label="交互方式" required>
            <a-select v-model="interactForm.interaction_type" placeholder="请选择交互方式">
              <a-option value="touch">触摸</a-option>
              <a-option value="voice">语音指令</a-option>
              <a-option value="gesture">手势指令</a-option>
              <a-option value="environmental">环境交互</a-option>
            </a-select>
          </a-form-item>
          <a-form-item label="位置 X">
            <a-input-number v-model="interactForm.position.x" :min="0" />
          </a-form-item>
          <a-form-item label="位置 Y">
            <a-input-number v-model="interactForm.position.y" :min="0" />
          </a-form-item>
          <template v-if="interactForm.interaction_type === 'touch'">
            <a-form-item label="力度">
              <a-slider v-model="interactForm.parameters.force" :min="0" :max="1" :step="0.1" show-input />
            </a-form-item>
            <a-form-item label="持续时间(ms)">
              <a-input-number v-model="interactForm.parameters.duration_ms" :min="0" :max="5000" />
            </a-form-item>
          </template>
        </a-form>
        <!-- 交互结果 -->
        <a-divider v-if="interactResult">交互结果</a-divider>
        <a-card v-if="interactResult" size="small" class="interact-result-card">
          <a-descriptions :column="2" size="small">
            <a-descriptions-item label="情绪">{{ getEmotionName(interactResult.emotion) }}</a-descriptions-item>
            <a-descriptions-item label="动作">{{ interactResult.action }}</a-descriptions-item>
            <a-descriptions-item label="动画">{{ interactResult.animation }}</a-descriptions-item>
            <a-descriptions-item label="声音">{{ interactResult.sound }}</a-descriptions-item>
            <a-descriptions-item label="响应时间">{{ interactResult.response_time_ms }}ms</a-descriptions-item>
          </a-descriptions>
        </a-card>
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import { IconPlus, IconRefresh } from '@arco-design/web-vue/es/icon'
import { getSimulationPets, createSimulationPet, updateSimulationPet, interactWithPet } from '@/api/simulation'

const pets = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)

const createModalVisible = ref(false)
const createForm = reactive({
  pet_name: '',
  pet_type: 'cat',
  personality: { friendliness: 0.8, playfulness: 0.6, energy_level: 0.7 },
  capabilities: { speech_recognition: true, emotion_recognition: true, path_planning: true },
  environment_id: 1
})

const interactModalVisible = ref(false)
const currentPet = ref(null)
const interactResult = ref(null)
const interactForm = reactive({
  interaction_type: 'touch',
  position: { x: 100, y: 200 },
  parameters: { force: 0.5, duration_ms: 500 }
})

async function loadPets() {
  try {
    const res = await getSimulationPets({ page: page.value, page_size: pageSize.value })
    pets.value = res.data?.items || res.data || []
    total.value = res.data?.total || 0
  } catch {
    Message.error('加载虚拟宠物列表失败')
  }
}

function openCreateModal() {
  Object.assign(createForm, {
    pet_name: '',
    pet_type: 'cat',
    personality: { friendliness: 0.8, playfulness: 0.6, energy_level: 0.7 },
    capabilities: { speech_recognition: true, emotion_recognition: true, path_planning: true },
    environment_id: 1
  })
  createModalVisible.value = true
}

async function handleCreate() {
  try {
    await createSimulationPet(createForm)
    Message.success('创建成功')
    createModalVisible.value = false
    loadPets()
  } catch {
    Message.error('创建失败')
    return false
  }
}

async function startPet(pet) {
  try {
    await updateSimulationPet(pet.id, { status: 'running' })
    Message.success('启动成功')
    loadPets()
  } catch {
    Message.error('启动失败')
  }
}

async function stopPet(pet) {
  try {
    await updateSimulationPet(pet.id, { status: 'idle' })
    Message.success('已停止')
    loadPets()
  } catch {
    Message.error('操作失败')
  }
}

function openInteractModal(pet) {
  currentPet.value = pet
  interactResult.value = null
  Object.assign(interactForm, {
    interaction_type: 'touch',
    position: { x: 100, y: 200 },
    parameters: { force: 0.5, duration_ms: 500 }
  })
  interactModalVisible.value = true
}

async function handleInteract() {
  if (!currentPet.value) return
  try {
    const res = await interactWithPet(currentPet.value.id, interactForm)
    interactResult.value = res.data?.pet_response || res.data
    Message.success('交互完成')
  } catch {
    Message.error('交互失败')
    return false
  }
}

function getPetTypeColor(type) {
  const colors = { cat: 'orange', dog: 'blue', rabbit: 'pink' }
  return colors[type] || 'gray'
}

function getPetTypeName(type) {
  const names = { cat: '猫', dog: '狗', rabbit: '兔子' }
  return names[type] || type
}

function getPetEmoji(type) {
  const emojis = { cat: '🐱', dog: '🐶', rabbit: '🐰' }
  return emojis[type] || '🐾'
}

function getStatusColor(status) {
  const colors = { idle: 'gray', running: 'green', paused: 'arcoblue' }
  return colors[status] || 'gray'
}

function getStatusName(status) {
  const names = { idle: '空闲', running: '运行中', paused: '已暂停' }
  return names[status] || status
}

function getEmotionName(emotion) {
  const names = { happy: '开心', sad: '悲伤', angry: '生气', fearful: '害怕', surprised: '惊讶', neutral: '平静' }
  return names[emotion] || emotion
}

function getBatteryColor(level) {
  if (!level) return 'gray'
  if (level > 0.5) return 'green'
  if (level > 0.2) return 'orange'
  return 'red'
}

function onPageChange(p) {
  page.value = p
  loadPets()
}

function onPageSizeChange(s) {
  pageSize.value = s
  page.value = 1
  loadPets()
}

onMounted(() => {
  loadPets()
})
</script>

<style scoped>
.pet-header {
  display: flex;
  align-items: center;
  gap: 8px;
}
.pet-avatar {
  display: flex;
  justify-content: center;
  margin-bottom: 12px;
}
.pet-avatar-placeholder {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  background: #F2F3F5;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 40px;
}
.interact-content {
  min-height: 200px;
}
.interact-pet-info {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 0;
}
.pet-emoji {
  font-size: 32px;
}
.interact-result-card {
  margin-top: 12px;
  background: #F7F8FA;
}
</style>
