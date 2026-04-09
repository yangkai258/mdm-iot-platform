<template>
  <div class="page-container">
    <a-card class="general-card" title="多宠物管�?>
      <template #extra>
        <a-button type="primary" @click="openAddPet"><icon-plus />添加宠物</a-button>
      </template>
      <div class="search-form">
        <a-form :model="form" layout="inline">
          <a-form-item label="宠物名称"><a-input v-model="form.pet_name" placeholder="请输�? /></a-form-item>
          <a-form-item><a-button type="primary" @click="loadData">查询</a-button><a-button @click="Object.keys(form).forEach(k => form[k] = ''); loadData()">重置</a-button></a-form-item>
        </a-form>
      </div>
      <!-- 当前活跃宠物 -->
      <div v-if="activePet" class="active-pet-section">
        <a-alert type="info" class="active-pet-alert">
          <template #title>
            <span>当前管理宠物: <strong>{{ activePet.pet_name }}</strong></span>
          </template>
          <template #content>
            <a-space>
              <span>类型: {{ activePet.pet_type }}</span>
              <span>设备: {{ activePet.device_id }}</span>
            </a-space>
          </template>
        </a-alert>
      </div>
      <!-- 宠物卡片列表 -->
      <a-row :gutter="16" class="pet-card-list">
        <a-col :span="6" v-for="pet in data" :key="pet.id">
          <a-card class="pet-card" :class="{ 'active': pet.id === activePetId }" @click="switchPet(pet)">
            <template #extra>
              <a-tag v-if="pet.id === activePetId" color="arcoblue">当前</a-tag>
            </template>
            <div class="pet-avatar">
              <a-avatar v-if="pet.pet_photo" :src="pet.pet_photo" :size="72" />
              <a-avatar v-else :size="72" style="background-color: #165dff; font-size: 28px">{{ pet.pet_name?.[0] || '?' }}</a-avatar>
            </div>
            <div class="pet-info">
              <div class="pet-name">{{ pet.pet_name }}</div>
              <div class="pet-type">{{ pet.pet_type }} / {{ pet.breed || '-' }}</div>
              <div class="pet-device">设备: {{ pet.device_id || '未绑�? }}</div>
              <a-space style="margin-top: 8px">
                <a-button type="text" size="small" @click.stop="viewProfile(pet)">档案</a-button>
                <a-button type="text" size="small" @click.stop="bindDevice(pet)">绑定设备</a-button>
              </a-space>
            </div>
          </a-card>
        </a-col>
      </a-row>
      <div style="margin-top: 16px">
        <a-pagination v-model:current="pagination.current" :total="pagination.total" :page-size="pagination.pageSize" show-total @change="onPageChange" />
      </div>
    </a-card>
    <a-modal v-model:visible="addVisible" title="添加宠物" @before-ok="handleAdd" :loading="submitting" :width="560">
      <a-form :model="addForm" layout="vertical">
        <a-form-item label="宠物名称" required><a-input v-model="addForm.pet_name" placeholder="请输�? /></a-form-item>
        <a-form-item label="宠物类型" required>
          <a-select v-model="addForm.pet_type" placeholder="选择类型">
            <a-option value="dog">�?/a-option>
            <a-option value="cat">�?/a-option>
            <a-option value="bird">�?/a-option>
            <a-option value="other">其他</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="品种"><a-input v-model="addForm.breed" placeholder="请输入品�? /></a-form-item>
        <a-form-item label="性别">
          <a-select v-model="addForm.gender" placeholder="选择性别">
            <a-option value="male">�?/a-option>
            <a-option value="female">�?/a-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>
    <a-modal v-model:visible="deviceBindVisible" title="绑定设备" @before-ok="handleBindDevice" :loading="submitting" :width="480">
      <a-form layout="vertical">
        <a-form-item label="当前宠物">{{ selectedPet?.pet_name }}</a-form-item>
        <a-form-item label="选择设备">
          <a-select v-model="selectedDeviceId" placeholder="选择要绑定的设备" style="width: 100%">
            <a-option v-for="d in devices" :key="d.id" :value="d.id">{{ d.name }} ({{ d.id }})</a-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>
    <a-modal v-model:visible="profileVisible" title="宠物档案" :footer="null" :width="560">
      <a-descriptions :column="2" bordered>
        <a-descriptions-item label="宠物名称" :span="2">{{ selectedPet?.pet_name }}</a-descriptions-item>
        <a-descriptions-item label="类型">{{ selectedPet?.pet_type }}</a-descriptions-item>
        <a-descriptions-item label="品种">{{ selectedPet?.breed }}</a-descriptions-item>
        <a-descriptions-item label="性别">{{ selectedPet?.gender }}</a-descriptions-item>
        <a-descriptions-item label="绑定设备">{{ selectedPet?.device_id || '未绑�? }}</a-descriptions-item>
        <a-descriptions-item label="添加时间" :span="2">{{ selectedPet?.created_at }}</a-descriptions-item>
      </a-descriptions>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const loading = ref(false)
const submitting = ref(false)
const data = ref([])
const devices = ref([])
const activePetId = ref(null)
const addVisible = ref(false)
const deviceBindVisible = ref(false)
const profileVisible = ref(false)
const selectedPet = ref(null)
const selectedDeviceId = ref('')
const form = reactive({ pet_name: '' })
const addForm = reactive({ pet_name: '', pet_type: '', breed: '', gender: '' })
const pagination = reactive({ current: 1, pageSize: 12, total: 0 })

const activePet = ref(null)

const loadData = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const params = new URLSearchParams({ page: pagination.current, page_size: pagination.pageSize })
    if (form.pet_name) params.append('pet_name', form.pet_name)
    const res = await fetch(`/api/v1/pet/multi?${params}`, { headers: { 'Authorization': `Bearer ${token}` } }).then(r => r.json())
    if (res.code === 0) { data.value = res.data?.list || []; pagination.total = res.data?.total || 0 }
    else { data.value = [] }
    if (data.value.length > 0 && !activePetId.value) {
      activePetId.value = data.value[0].id
      activePet.value = data.value[0]
    }
  } catch (e) { Message.error('加载失败') } finally { loading.value = false }
}

const loadDevices = async () => {
  try {
    const token = localStorage.getItem('token')
    const res = await fetch('/api/v1/devices?page_size=200', { headers: { 'Authorization': `Bearer ${token}` } }).then(r => r.json())
    devices.value = res.data?.list || []
  } catch (e) { console.error('加载设备失败', e) }
}

const switchPet = (pet) => {
  activePetId.value = pet.id
  activePet.value = pet
  localStorage.setItem('active_pet_id', pet.id)
  Message.success(`已切换到宠物: ${pet.pet_name}`)
}

const viewProfile = (pet) => { selectedPet.value = pet; profileVisible.value = true }
const bindDevice = (pet) => { selectedPet.value = pet; selectedDeviceId.value = pet.device_id || ''; deviceBindVisible.value = true }

const handleAdd = async (done) => {
  if (!addForm.pet_name || !addForm.pet_type) { Message.warning('请填写必填项'); done(false); return }
  submitting.value = true
  try {
    const token = localStorage.getItem('token')
    const res = await fetch('/api/v1/pet/multi', { method: 'POST', headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' }, body: JSON.stringify(addForm) }).then(r => r.json())
    if (res.code === 0) { Message.success('添加成功'); addVisible.value = false; loadData() }
    else { Message.error(res.message || '添加失败') }
    done(true)
  } catch (e) { Message.error('添加失败'); done(false) } finally { submitting.value = false }
}

const handleBindDevice = async (done) => {
  if (!selectedDeviceId.value) { Message.warning('请选择设备'); done(false); return }
  submitting.value = true
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`/api/v1/pet/multi/${selectedPet.value.id}/bind`, { method: 'POST', headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' }, body: JSON.stringify({ device_id: selectedDeviceId.value }) }).then(r => r.json())
    if (res.code === 0) { Message.success('绑定成功'); deviceBindVisible.value = false; loadData() }
    else { Message.error(res.message || '绑定失败') }
    done(true)
  } catch (e) { Message.error('绑定失败'); done(false) } finally { submitting.value = false }
}

const openAddPet = () => { Object.assign(addForm, { pet_name: '', pet_type: '', breed: '', gender: '' }); addVisible.value = true }
const onPageChange = (page) => { pagination.current = page; loadData() }

onMounted(() => { activePetId.value = localStorage.getItem('active_pet_id'); loadData(); loadDevices() })
</script>

<style scoped>
.page-container { padding: 16px; }
.search-form { margin-bottom: 16px; padding: 16px; background: var(--color-fill-lightest); border-radius: 4px; }
.active-pet-section { margin-bottom: 16px; }
.pet-card-list { margin-top: 8px; }
.pet-card { cursor: pointer; transition: all 0.3s; text-align: center; }
.pet-card:hover { box-shadow: 0 4px 12px rgba(0,0,0,0.1); }
.pet-card.active { border: 2px solid #165dff; }
.pet-avatar { display: flex; justify-content: center; margin-bottom: 12px; }
.pet-info { text-align: center; }
.pet-name { font-size: 16px; font-weight: 600; margin-bottom: 4px; }
.pet-type { font-size: 12px; color: #666; margin-bottom: 4px; }
.pet-device { font-size: 12px; color: #999; }
</style>