<template>
  <div class="pet-config-container">

    <a-card class="general-card">
      <template #title><span class="card-title">宠物查询</span></template>
      <a-row :gutter="16">
        <a-col :flex="1">
          <a-form :model="searchForm" layout="vertical" size="small">
            <a-row :gutter="16">
              <a-col :span="8">
                <a-form-item label="关键词">
                  <a-input v-model="searchForm.keyword" placeholder="宠物名称" allow-clear />
                </a-form-item>
              </a-col>
              <a-col :span="8">
                <a-form-item label="宠物性格">
                  <a-select v-model="searchForm.personality" placeholder="全部" allow-clear>
                    <a-option :value="1">活泼好动</a-option>
                    <a-option :value="2">温顺安静</a-option>
                    <a-option :value="3">好奇宝宝</a-option>
                    <a-option :value="4">独立自主</a-option>
                    <a-option :value="5">粘人依赖</a-option>
                  </a-select>
                </a-form-item>
              </a-col>
            </a-row>
          </a-form>
        </a-col>
        <a-divider style="height: 84px" direction="vertical" />
        <a-col :flex="'86px'" style="text-align: right">
          <a-space direction="vertical" :size="18">
            <a-button type="primary" @click="handleSearch">
              <template #icon><icon-search /></template>
              查询
            </a-button>
            <a-button @click="handleReset">
              <template #icon><icon-refresh /></template>
              重置
            </a-button>
          </a-space>
        </a-col>
      </a-row>
    </a-card>

    <a-card class="general-card" style="margin-top: 16px">
      <template #title><span class="card-title">宠物列表</span></template>
      <template #extra>
        <a-space>
          <a-button type="primary" @click="showAddModal">
            <template #icon><icon-plus /></template>
            添加宠物
          </a-button>
          <a-button @click="loadPets">
            <template #icon><icon-refresh /></template>
            刷新
          </a-button>
        </a-space>
      </template>

      <a-table
        :columns="columns"
        :data="filteredPets"
        :loading="loading"
        :pagination="paginationConfig"
        @page-change="onPageChange"
        @page-size-change="onPageSizeChange"
        row-key="pet_id"
      >
        <template #avatar="{ record }">
          <a-avatar :size="40" :style="{ backgroundColor: record.avatar_color }">
            {{ record.pet_name?.charAt(0) || '?' }}
          </a-avatar>
        </template>
        <template #personality="{ record }">
          <a-tag :color="getPersonalityColor(record.personality)">
            {{ getPersonalityText(record.personality) }}
          </a-tag>
        </template>
        <template #dndEnabled="{ record }">
          <a-tag :color="record.dnd_enabled ? 'green' : 'gray'">
            {{ record.dnd_enabled ? '已启用' : '已禁用' }}
          </a-tag>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="editPet(record)">配置</a-button>
            <a-button type="text" size="small" status="danger" @click="deletePet(record)">删除</a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 添加/编辑宠物弹窗 -->
    <a-modal
      v-model:visible="modalVisible"
      :title="isEdit ? '编辑宠物配置' : '添加宠物'"
      @ok="handleSave"
      :confirm-loading="saving"
      :width="520"
    >
      <a-form :model="petForm" layout="vertical">
        <a-form-item label="宠物名称" required>
          <a-input v-model="petForm.pet_name" placeholder="请输入宠物名称" />
        </a-form-item>
        <a-form-item label="宠物性格" required>
          <a-select v-model="petForm.personality" placeholder="选择宠物性格">
            <a-option :value="1">活泼好动</a-option>
            <a-option :value="2">温顺安静</a-option>
            <a-option :value="3">好奇宝宝</a-option>
            <a-option :value="4">独立自主</a-option>
            <a-option :value="5">粘人依赖</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="绑定设备">
          <a-select v-model="petForm.device_id" placeholder="选择绑定的设备" allow-clear>
            <a-option v-for="dev in devices" :key="dev.device_id" :value="dev.device_id">
              {{ dev.device_id }} ({{ dev.hardware_model }})
            </a-option>
          </a-select>
        </a-form-item>
        <a-divider>免打扰设置</a-divider>
        <a-form-item label="启用免打扰">
          <a-switch v-model="petForm.dnd_enabled" />
        </a-form-item>
        <a-form-item v-if="petForm.dnd_enabled" label="免打扰时间">
          <a-space>
            <a-time-picker v-model="petForm.dnd_start_time" format="HH:mm" placeholder="开始时间" style="width: 140px" />
            <span>至</span>
            <a-time-picker v-model="petForm.dnd_end_time" format="HH:mm" placeholder="结束时间" style="width: 140px" />
          </a-space>
        </a-form-item>
        <a-divider>提醒设置</a-divider>
        <a-form-item label="运动提醒">
          <a-switch v-model="petForm.settings.exercise_reminder" />
        </a-form-item>
        <a-form-item label="喂食提醒">
          <a-switch v-model="petForm.settings.feeding_reminder" />
        </a-form-item>
        <a-form-item label="健康监测">
          <a-switch v-model="petForm.settings.health_monitoring" />
        </a-form-item>
        <a-form-item label="备注">
          <a-textarea v-model="petForm.remark" placeholder="备注信息" :rows="2" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import dayjs from 'dayjs'

const loading = ref(false)
const saving = ref(false)
const modalVisible = ref(false)
const isEdit = ref(false)

const searchForm = reactive({ keyword: '', personality: undefined })
const pets = ref([])
const devices = ref([])

const pagination = reactive({ current: 1, pageSize: 20, total: 0 })

const paginationConfig = computed(() => ({
  current: pagination.current,
  pageSize: pagination.pageSize,
  total: pagination.total,
  showTotal: true,
  showSizeChanger: true
}))

const petForm = reactive({
  pet_id: '', pet_name: '', personality: 1, device_id: '',
  dnd_enabled: false, dnd_start_time: null, dnd_end_time: null,
  settings: { exercise_reminder: true, feeding_reminder: true, health_monitoring: true },
  remark: ''
})

const columns = [
  { title: '头像', slotName: 'avatar', width: 80 },
  { title: '宠物名称', dataIndex: 'pet_name', ellipsis: true },
  { title: '宠物性格', slotName: 'personality', width: 120 },
  { title: '绑定设备', dataIndex: 'device_id', ellipsis: true },
  { title: '免打扰', slotName: 'dndEnabled', width: 100 },
  { title: '创建时间', dataIndex: 'create_time', width: 170 },
  { title: '操作', slotName: 'actions', width: 150, fixed: 'right' }
]

const API_BASE = '/api/v1'

const filteredPets = computed(() => {
  let result = pets.value
  if (searchForm.keyword) {
    result = result.filter(p => p.pet_name.includes(searchForm.keyword))
  }
  if (searchForm.personality !== undefined) {
    result = result.filter(p => p.personality === searchForm.personality)
  }
  return result
})

const loadPets = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`${API_BASE}/pets?page=${pagination.current}&page_size=${pagination.pageSize}`, {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    const data = await res.json()
    if (data.code === 0) {
      pets.value = data.data.list || []
      pagination.total = data.data.pagination?.total || 0
    }
  } catch (e) {
    pets.value = getMockPets()
    pagination.total = pets.value.length
  } finally { loading.value = false }
}

const getMockPets = () => [
  { pet_id: 'PET001', pet_name: '小橘', personality: 1, device_id: 'DEV001', dnd_enabled: true, dnd_start_time: '22:00', dnd_end_time: '07:00', avatar_color: '#f6ad55', create_time: '2026-03-15 10:00:00', settings: { exercise_reminder: true, feeding_reminder: true, health_monitoring: true }, remark: '' },
  { pet_id: 'PET002', pet_name: '布丁', personality: 2, device_id: 'DEV002', dnd_enabled: false, dnd_start_time: '', dnd_end_time: '', avatar_color: '#fc8181', create_time: '2026-03-10 14:30:00', settings: { exercise_reminder: true, feeding_reminder: false, health_monitoring: true }, remark: '需要特别关注' },
  { pet_id: 'PET003', pet_name: '豆豆', personality: 3, device_id: 'DEV003', dnd_enabled: true, dnd_start_time: '12:00', dnd_end_time: '14:00', avatar_color: '#68d391', create_time: '2026-03-05 09:15:00', settings: { exercise_reminder: false, feeding_reminder: true, health_monitoring: false }, remark: '' }
]

const loadDevices = async () => {
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`${API_BASE}/devices?page_size=100`, {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    const data = await res.json()
    if (data.code === 0) devices.value = data.data.list || []
  } catch (e) {
    devices.value = [
      { device_id: 'DEV001', hardware_model: 'MDM-Pro-200' },
      { device_id: 'DEV002', hardware_model: 'MDM-Mini-100' },
      { device_id: 'DEV003', hardware_model: 'MDM-Lite-50' }
    ]
  }
}

const handleSearch = () => { pagination.current = 1; loadPets() }
const handleReset = () => { searchForm.keyword = ''; searchForm.personality = undefined; pagination.current = 1; loadPets() }
const onPageChange = (page) => { pagination.current = page; loadPets() }
const onPageSizeChange = (size) => { pagination.pageSize = size; pagination.current = 1; loadPets() }

const showAddModal = () => {
  isEdit.value = false
  resetForm()
  modalVisible.value = true
}

const editPet = (record) => {
  isEdit.value = true
  Object.assign(petForm, {
    ...record,
    dnd_start_time: record.dnd_start_time ? dayjs(record.dnd_start_time, 'HH:mm') : null,
    dnd_end_time: record.dnd_end_time ? dayjs(record.dnd_end_time, 'HH:mm') : null,
    settings: { ...record.settings }
  })
  modalVisible.value = true
}

const handleSave = async () => {
  if (!petForm.pet_name) { Message.warning('请输入宠物名称'); return }
  saving.value = true
  try {
    const token = localStorage.getItem('token')
    const submitData = {
      ...petForm,
      dnd_start_time: petForm.dnd_start_time?.format?.('HH:mm') || '',
      dnd_end_time: petForm.dnd_end_time?.format?.('HH:mm') || ''
    }
    const method = isEdit.value ? 'PUT' : 'POST'
    const url = isEdit.value ? `${API_BASE}/pets/${petForm.pet_id}` : `${API_BASE}/pets`
    const res = await fetch(url, {
      method,
      headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' },
      body: JSON.stringify(submitData)
    })
    const data = await res.json()
    if (data.code === 0) { Message.success(isEdit.value ? '更新成功' : '添加成功'); modalVisible.value = false; loadPets(); return }
  } catch (e) {}
  // mock
  if (isEdit.value) {
    const idx = pets.value.findIndex(p => p.pet_id === petForm.pet_id)
    if (idx !== -1) pets.value[idx] = { ...pets.value[idx], ...petForm, dnd_start_time: petForm.dnd_start_time?.format?.('HH:mm'), dnd_end_time: petForm.dnd_end_time?.format?.('HH:mm') }
    Message.success('宠物配置已更新（模拟）')
  } else {
    pets.value.unshift({ pet_id: `PET${Date.now()}`, ...petForm, dnd_start_time: petForm.dnd_start_time?.format?.('HH:mm'), dnd_end_time: petForm.dnd_end_time?.format?.('HH:mm'), create_time: new Date().toLocaleString(), avatar_color: ['#f6ad55', '#fc8181', '#68d391', '#63b3ed', '#f687b3'][Math.floor(Math.random() * 5)] })
    Message.success('宠物添加成功（模拟）')
  }
  modalVisible.value = false
  saving.value = false
}

const deletePet = (record) => { pets.value = pets.value.filter(p => p.pet_id !== record.pet_id); Message.success('宠物已删除') }

const resetForm = () => {
  Object.assign(petForm, {
    pet_id: '', pet_name: '', personality: 1, device_id: '',
    dnd_enabled: false, dnd_start_time: null, dnd_end_time: null,
    settings: { exercise_reminder: true, feeding_reminder: true, health_monitoring: true }, remark: ''
  })
}

const getPersonalityColor = (p) => ({ 1: 'green', 2: 'blue', 3: 'purple', 4: 'orange', 5: 'pink' }[p] || 'default')
const getPersonalityText = (p) => ({ 1: '活泼好动', 2: '温顺安静', 3: '好奇宝宝', 4: '独立自主', 5: '粘人依赖' }[p] || '未知')

onMounted(() => { loadPets(); loadDevices() })
</script>

<style scoped>
.pet-config-container { padding: 20px 24px; min-height: calc(100vh - 64px); background: #f5f7fa; }
.general-card { border-radius: 8px; }
.card-title { font-weight: 600; font-size: 15px; }
</style>
