<template>
  <div class="page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>第三方集成</a-breadcrumb-item>
      <a-breadcrumb-item>宠物医疗</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 页面标题 + 操作按钮 -->
    <div class="page-header">
      <div class="header-left">
        <h2>宠物医疗</h2>
      </div>
      <div class="header-right">
        <a-button type="primary" @click="showCreateModal">
          <template #icon><icon-plus /></template>
          新建预约
        </a-button>
      </div>
    </div>

    <!-- 搜索筛选栏 -->
    <div class="filter-bar">
      <a-space wrap>
        <a-input-search
          v-model="searchKeyword"
          placeholder="搜索宠物名称..."
          style="width: 240px"
          search-button
          @search="loadAppointments"
        />
        <a-select
          v-model="filterStatus"
          placeholder="预约状态"
          style="width: 140px"
          allow-clear
          @change="loadAppointments"
        >
          <a-option value="pending">待确认</a-option>
          <a-option value="confirmed">已确认</a-option>
          <a-option value="completed">已完成</a-option>
          <a-option value="cancelled">已取消</a-option>
        </a-select>
        <a-select
          v-model="filterDept"
          placeholder="科室"
          style="width: 140px"
          allow-clear
          @change="loadAppointments"
        >
          <a-option value="general">全科</a-option>
          <a-option value="surgery">外科</a-option>
          <a-option value="dermatology">皮肤科</a-option>
          <a-option value="dentistry">牙科</a-option>
          <a-option value="ophthalmology">眼科</a-option>
        </a-select>
        <a-range-picker
          v-model="dateRange"
          style="width: 260px"
          @change="loadAppointments"
        />
      </a-space>
    </div>

    <!-- 预约列表 -->
    <div class="appointment-list">
      <div v-if="loading" class="loading-state">
        <a-spin size="large" />
      </div>
      <div v-else-if="filteredAppointments.length === 0" class="empty-state">
        <icon-calendar class="empty-icon" />
        <p>暂无预约记录</p>
        <a-button type="primary" @click="showCreateModal">立即预约</a-button>
      </div>
      <a-table
        v-else
        :columns="columns"
        :data="filteredAppointments"
        :pagination="{ pageSize: 10 }"
        row-key="id"
        @page-change="(page) => currentPage = page"
      >
        <template #status="{ record }">
          <a-tag :color="getStatusColor(record.status)" class="status-tag">
            {{ getStatusText(record.status) }}
          </a-tag>
        </template>
        <template #pet="{ record }">
          <div class="pet-cell">
            <span class="pet-avatar-small">{{ getPetEmoji(record.pet_type) }}</span>
            <span>{{ record.pet_name }}</span>
          </div>
        </template>
        <template #time="{ record }">
          <div class="time-cell">
            <div>{{ formatDate(record.appointment_date) }}</div>
            <div class="time-slot">{{ record.time_slot }}</div>
          </div>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button size="small" @click="viewDetail(record)">详情</a-button>
            <a-button
              v-if="record.status === 'pending'"
              size="small"
              type="primary"
              @click="confirmAppointment(record)"
            >
              确认
            </a-button>
            <a-button
              v-if="record.status === 'pending' || record.status === 'confirmed'"
              size="small"
              status="warning"
              @click="cancelAppointment(record)"
            >
              取消
            </a-button>
          </a-space>
        </template>
      </a-table>
    </div>

    <!-- 创建预约弹窗 -->
    <a-modal
      v-model:visible="createModalVisible"
      title="新建预约"
      :width="600"
      @ok="handleCreate"
      @cancel="createModalVisible = false"
    >
      <a-form :model="appointmentForm" layout="vertical">
        <a-form-item label="选择宠物" required>
          <a-select
            v-model="appointmentForm.pet_id"
            placeholder="请选择宠物"
            style="width: 100%"
          >
            <a-option v-for="pet in petList" :key="pet.pet_id" :value="pet.pet_id">
              <a-space>
                <span>{{ getPetEmoji(pet.pet_type) }}</span>
                <span>{{ pet.pet_name }}</span>
              </a-space>
            </a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="选择医院" required>
          <a-select
            v-model="appointmentForm.hospital_id"
            placeholder="请选择医院"
            style="width: 100%"
            @change="loadDepartments"
          >
            <a-option v-for="h in hospitalList" :key="h.id" :value="h.id">{{ h.name }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="选择科室" required>
          <a-select
            v-model="appointmentForm.department"
            placeholder="请选择科室"
            style="width: 100%"
          >
            <a-option v-for="d in departmentList" :key="d.value" :value="d.value">{{ d.label }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="预约日期" required>
          <a-date-picker
            v-model="appointmentForm.appointment_date"
            style="width: 100%"
            :disabled-date="disabledDate"
          />
        </a-form-item>
        <a-form-item label="预约时间段" required>
          <a-select
            v-model="appointmentForm.time_slot"
            placeholder="请选择时间段"
            style="width: 100%"
          >
            <a-option value="09:00-10:00">09:00-10:00</a-option>
            <a-option value="10:00-11:00">10:00-11:00</a-option>
            <a-option value="11:00-12:00">11:00-12:00</a-option>
            <a-option value="14:00-15:00">14:00-15:00</a-option>
            <a-option value="15:00-16:00">15:00-16:00</a-option>
            <a-option value="16:00-17:00">16:00-17:00</a-option>
            <a-option value="17:00-18:00">17:00-18:00</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="症状描述" required>
          <a-textarea
            v-model="appointmentForm.symptoms"
            placeholder="请描述宠物的主要症状"
            :rows="3"
          />
        </a-form-item>
        <a-form-item label="备注">
          <a-textarea
            v-model="appointmentForm.notes"
            placeholder="其他补充说明（可选）"
            :rows="2"
          />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 预约详情弹窗 -->
    <a-drawer
      v-model:visible="detailDrawerVisible"
      title="预约详情"
      :width="480"
      @cancel="detailDrawerVisible = false"
    >
      <a-descriptions :column="1" bordered>
        <a-descriptions-item label="预约编号">{{ currentAppointment?.id }}</a-descriptions-item>
        <a-descriptions-item label="宠物名称">{{ currentAppointment?.pet_name }}</a-descriptions-item>
        <a-descriptions-item label="预约医院">{{ currentAppointment?.hospital_name }}</a-descriptions-item>
        <a-descriptions-item label="科室">{{ getDeptName(currentAppointment?.department) }}</a-descriptions-item>
        <a-descriptions-item label="预约日期">{{ formatDate(currentAppointment?.appointment_date) }}</a-descriptions-item>
        <a-descriptions-item label="时间段">{{ currentAppointment?.time_slot }}</a-descriptions-item>
        <a-descriptions-item label="预约状态">
          <a-tag :color="getStatusColor(currentAppointment?.status)">
            {{ getStatusText(currentAppointment?.status) }}
          </a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="症状描述">{{ currentAppointment?.symptoms }}</a-descriptions-item>
        <a-descriptions-item label="备注">{{ currentAppointment?.notes || '无' }}</a-descriptions-item>
        <a-descriptions-item label="创建时间">{{ currentAppointment?.created_at }}</a-descriptions-item>
      </a-descriptions>
      <template #footer>
        <a-space>
          <a-button
            v-if="currentAppointment?.status === 'pending'"
            type="primary"
            @click="confirmAppointment(currentAppointment); detailDrawerVisible = false"
          >
            确认预约
          </a-button>
          <a-button
            v-if="currentAppointment?.status === 'pending' || currentAppointment?.status === 'confirmed'"
            status="warning"
            @click="cancelAppointment(currentAppointment); detailDrawerVisible = false"
          >
            取消预约
          </a-button>
          <a-button @click="detailDrawerVisible = false">关闭</a-button>
        </a-space>
      </template>
    </a-drawer>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import { IconPlus, IconCalendar } from '@arco-design/web-vue/es/icon'

const loading = ref(false)
const appointments = ref([])
const searchKeyword = ref('')
const filterStatus = ref('')
const filterDept = ref('')
const dateRange = ref([])
const currentPage = ref(1)
const createModalVisible = ref(false)
const detailDrawerVisible = ref(false)
const currentAppointment = ref(null)
const petList = ref([])
const hospitalList = ref([])

const departmentList = ref([
  { value: 'general', label: '全科' },
  { value: 'surgery', label: '外科' },
  { value: 'dermatology', label: '皮肤科' },
  { value: 'dentistry', label: '牙科' },
  { value: 'ophthalmology', label: '眼科' }
])

const appointmentForm = reactive({
  pet_id: null,
  hospital_id: null,
  department: '',
  appointment_date: null,
  time_slot: '',
  symptoms: '',
  notes: ''
})

const columns = [
  { title: '预约编号', dataIndex: 'id', width: 120 },
  { title: '宠物', slotName: 'pet', width: 140 },
  { title: '医院', dataIndex: 'hospital_name', width: 160 },
  { title: '科室', dataIndex: 'department', width: 100 },
  { title: '预约时间', slotName: 'time', width: 140 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '操作', slotName: 'actions', width: 180 }
]

const filteredAppointments = computed(() => {
  let result = appointments.value
  if (searchKeyword.value) {
    const kw = searchKeyword.value.toLowerCase()
    result = result.filter(a => a.pet_name.toLowerCase().includes(kw))
  }
  if (filterStatus.value) {
    result = result.filter(a => a.status === filterStatus.value)
  }
  if (filterDept.value) {
    result = result.filter(a => a.department === filterDept.value)
  }
  return result
})

function getPetEmoji(type) {
  const map = { dog: '🐕', cat: '🐈', bird: '🐦', rabbit: '🐰', other: '🐾' }
  return map[type] || '🐾'
}

function getStatusColor(status) {
  const map = {
    pending: 'orange',
    confirmed: 'arcoblue',
    completed: 'green',
    cancelled: 'gray'
  }
  return map[status] || 'gray'
}

function getStatusText(status) {
  const map = {
    pending: '待确认',
    confirmed: '已确认',
    completed: '已完成',
    cancelled: '已取消'
  }
  return map[status] || status
}

function getDeptName(dept) {
  const d = departmentList.value.find(item => item.value === dept)
  return d ? d.label : dept
}

function formatDate(dateStr) {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
}

function disabledDate(date) {
  return date && date < new Date(new Date().setHours(0, 0, 0, 0))
}

async function loadAppointments() {
  loading.value = true
  try {
    const res = await fetch('/api/v1/pet-hospital/appointments', {
      headers: { Authorization: `Bearer ${localStorage.getItem('token') || ''}` }
    })
    const data = await res.json()
    if (data.data) {
      appointments.value = data.data
    } else {
      loadMockAppointments()
    }
  } catch {
    loadMockAppointments()
  } finally {
    loading.value = false
  }
}

function loadMockAppointments() {
  appointments.value = [
    {
      id: 'APT2026030101',
      pet_id: 'p1',
      pet_name: '小旺',
      pet_type: 'dog',
      hospital_name: '阳光宠物医院',
      department: 'general',
      appointment_date: '2026-03-25',
      time_slot: '10:00-11:00',
      status: 'pending',
      symptoms: '食欲不振，精神萎靡',
      notes: '',
      created_at: '2026-03-20 14:30:00'
    },
    {
      id: 'APT2026030102',
      pet_id: 'p2',
      pet_name: '小美',
      pet_type: 'cat',
      hospital_name: '爱宠宠物诊所',
      department: 'dermatology',
      appointment_date: '2026-03-24',
      time_slot: '14:00-15:00',
      status: 'confirmed',
      symptoms: '皮肤红肿，脱毛',
      notes: '之前有过皮肤病史',
      created_at: '2026-03-18 09:15:00'
    },
    {
      id: 'APT2026030103',
      pet_id: 'p1',
      pet_name: '小旺',
      pet_type: 'dog',
      hospital_name: '阳光宠物医院',
      department: 'surgery',
      appointment_date: '2026-03-15',
      time_slot: '09:00-10:00',
      status: 'completed',
      symptoms: '常规体检',
      notes: '',
      created_at: '2026-03-10 16:00:00'
    }
  ]
}

function showCreateModal() {
  appointmentForm.pet_id = null
  appointmentForm.hospital_id = null
  appointmentForm.department = ''
  appointmentForm.appointment_date = null
  appointmentForm.time_slot = ''
  appointmentForm.symptoms = ''
  appointmentForm.notes = ''
  createModalVisible.value = true
}

async function handleCreate() {
  if (!appointmentForm.pet_id || !appointmentForm.hospital_id || !appointmentForm.department ||
      !appointmentForm.appointment_date || !appointmentForm.time_slot || !appointmentForm.symptoms) {
    Message.warning('请填写完整的预约信息')
    return
  }
  try {
    const res = await fetch('/api/v1/pet-hospital/appointments', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${localStorage.getItem('token') || ''}`
      },
      body: JSON.stringify(appointmentForm)
    })
    const data = await res.json()
    if (data.id) {
      Message.success('预约创建成功')
      createModalVisible.value = false
      loadAppointments()
    }
  } catch {
    // mock
    const newAppt = {
      id: `APT${Date.now()}`,
      ...appointmentForm,
      pet_name: petList.value.find(p => p.pet_id === appointmentForm.pet_id)?.pet_name || '',
      hospital_name: hospitalList.value.find(h => h.id === appointmentForm.hospital_id)?.name || '',
      status: 'pending',
      created_at: new Date().toLocaleString()
    }
    appointments.value.unshift(newAppt)
    Message.success('预约创建成功')
    createModalVisible.value = false
  }
}

function viewDetail(record) {
  currentAppointment.value = record
  detailDrawerVisible.value = true
}

async function confirmAppointment(record) {
  try {
    await fetch(`/api/v1/pet-hospital/appointments/${record.id}/confirm`, {
      method: 'PUT',
      headers: { Authorization: `Bearer ${localStorage.getItem('token') || ''}` }
    })
    record.status = 'confirmed'
    Message.success('预约已确认')
  } catch {
    record.status = 'confirmed'
    Message.success('预约已确认')
  }
}

async function cancelAppointment(record) {
  try {
    await fetch(`/api/v1/pet-hospital/appointments/${record.id}/cancel`, {
      method: 'PUT',
      headers: { Authorization: `Bearer ${localStorage.getItem('token') || ''}` }
    })
    record.status = 'cancelled'
    Message.success('预约已取消')
  } catch {
    record.status = 'cancelled'
    Message.success('预约已取消')
  }
}

async function loadDepartments() {
  // 动态加载科室列表，根据医院ID获取
}

onMounted(() => {
  loadAppointments()
  petList.value = [
    { pet_id: 'p1', pet_name: '小旺', pet_type: 'dog' },
    { pet_id: 'p2', pet_name: '小美', pet_type: 'cat' }
  ]
  hospitalList.value = [
    { id: 'h1', name: '阳光宠物医院' },
    { id: 'h2', name: '爱宠宠物诊所' },
    { id: 'h3', name: '康乐宠物医疗中心' }
  ]
})
</script>

<style scoped>
.page-container {
  padding: 24px;
  min-height: 100vh;
  background: #f5f6f7;
}

.breadcrumb {
  margin-bottom: 16px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.page-header h2 {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
  color: var(--color-text-1);
}

.filter-bar {
  background: #fff;
  padding: 16px;
  border-radius: 8px;
  margin-bottom: 16px;
}

.appointment-list {
  background: #fff;
  border-radius: 8px;
  padding: 16px;
}

.loading-state,
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 0;
  gap: 16px;
}

.empty-icon {
  font-size: 48px;
  color: #ccc;
}

.empty-state p {
  color: #999;
  font-size: 14px;
  margin: 0;
}

.status-tag {
  border-radius: 4px;
}

.pet-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.pet-avatar-small {
  font-size: 18px;
}

.time-cell {
  font-size: 13px;
}

.time-slot {
  color: #86909c;
  font-size: 12px;
}
</style>
