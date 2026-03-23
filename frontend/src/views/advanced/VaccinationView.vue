<template>
  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>高级功能</a-breadcrumb-item>
      <a-breadcrumb-item>疫苗接种</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 页面标题 -->
    <div class="pro-page-header">
      <h2 class="pro-page-title">疫苗接种</h2>
      <p class="pro-page-desc">管理宠物疫苗接种记录，设置接种提醒</p>
    </div>

    <!-- 标签页切换 -->
    <a-tabs v-model="activeTab" class="pro-tabs">
      <a-tab-pane key="records" title="接种记录">
        <!-- 搜索筛选区 -->
        <div class="pro-search-bar">
          <a-space>
            <a-select v-model="petFilter" placeholder="选择宠物" allow-clear style="width: 160px" @change="loadRecords">
              <a-option v-for="p in pets" :key="p.id" :value="p.id">{{ p.name }}</a-option>
            </a-select>
            <a-select v-model="vaccineTypeFilter" placeholder="疫苗类型" allow-clear style="width: 160px" @change="loadRecords">
              <a-option value="rabies">狂犬疫苗</a-option>
              <a-option value="distemper">犬瘟热</a-option>
              <a-option value="parvovirus">犬细小病毒</a-option>
              <a-option value="feline">猫瘟</a-option>
              <a-option value="other">其他</a-option>
            </a-select>
            <a-input-search v-model="keyword" placeholder="搜索疫苗名称" style="width: 200px" @search="loadRecords" search-button />
          </a-space>
        </div>

        <!-- 操作按钮区 -->
        <div class="pro-action-bar">
          <a-space>
            <a-button type="primary" @click="showRecordModal(null)">
              <template #icon><icon-add /></template>
              添加记录
            </a-button>
            <a-button @click="loadRecords">
              <template #icon><icon-refresh /></template>
              刷新
            </a-button>
          </a-space>
        </div>

        <!-- 接种记录列表 -->
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
              <a-avatar :style="{ backgroundColor: '#1659f5' }" :size="28">
                {{ record.pet_name?.charAt(0) || '?' }}
              </a-avatar>
              <span style="margin-left: 8px">{{ record.pet_name }}</span>
            </template>
            <template #vaccine_name="{ record }">
              <span class="vaccine-name">{{ record.vaccine_name }}</span>
            </template>
            <template #vaccine_type="{ record }">
              <a-tag :color="getVaccineTypeColor(record.vaccine_type)">
                {{ getVaccineTypeLabel(record.vaccine_type) }}
              </a-tag>
            </template>
            <template #dose_number="{ record }">
              {{ record.dose_number }} / {{ record.total_doses }}
            </template>
            <template #inoculation_date="{ record }">
              {{ formatDate(record.inoculation_date) }}
            </template>
            <template #next_due_date="{ record }">
              <span v-if="record.next_due_date" :class="isOverdue(record.next_due_date) ? 'overdue-text' : ''">
                {{ formatDate(record.next_due_date) }}
                <a-tag v-if="isOverdue(record.next_due_date)" color="red" size="small">已过期</a-tag>
                <a-tag v-else-if="isDueSoon(record.next_due_date)" color="orange" size="small">即将到期</a-tag>
              </span>
              <span v-else class="text-muted">—</span>
            </template>
            <template #vet_name="{ record }">
              {{ record.vet_name || '—' }}
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

      <a-tab-pane key="reminders" title="提醒设置">
        <!-- 操作按钮区 -->
        <div class="pro-action-bar">
          <a-space>
            <a-button type="primary" @click="showReminderModal(null)">
              <template #icon><icon-add /></template>
              添加提醒
            </a-button>
          </a-space>
        </div>

        <!-- 提醒列表 -->
        <div class="pro-content-area">
          <a-table
            :columns="reminderColumns"
            :data="reminders"
            :loading="reminderLoading"
            :pagination="reminderPagination"
            @page-change="onReminderPageChange"
            row-key="id"
          >
            <template #pet_name="{ record }">
              <a-avatar :style="{ backgroundColor: '#0fc6c2' }" :size="28">
                {{ record.pet_name?.charAt(0) || '?' }}
              </a-avatar>
              <span style="margin-left: 8px">{{ record.pet_name }}</span>
            </template>
            <template #vaccine_name="{ record }">
              {{ record.vaccine_name }}
            </template>
            <template #reminder_date="{ record }">
              {{ formatDate(record.reminder_date) }}
            </template>
            <template #advance_days="{ record }">
              提前 {{ record.advance_days }} 天
            </template>
            <template #enabled="{ record }">
              <a-switch v-model="record.enabled" @change="toggleReminder(record)" :disabled="savingId === record.id" />
            </template>
            <template #actions="{ record }">
              <a-space>
                <a-button type="text" size="small" @click="showReminderModal(record)">编辑</a-button>
                <a-button type="text" size="small" status="danger" @click="handleDeleteReminder(record)">删除</a-button>
              </a-space>
            </template>
          </a-table>
        </div>
      </a-tab-pane>
    </a-tabs>

    <!-- 接种记录弹窗 -->
    <a-modal
      v-model:visible="recordModalVisible"
      :title="isEditRecord ? '编辑接种记录' : '添加接种记录'"
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
        <a-form-item label="疫苗名称" required>
          <a-input v-model="recordForm.vaccine_name" placeholder="请输入疫苗名称" />
        </a-form-item>
        <a-form-item label="疫苗类型" required>
          <a-select v-model="recordForm.vaccine_type" placeholder="请选择疫苗类型">
            <a-option value="rabies">狂犬疫苗</a-option>
            <a-option value="distemper">犬瘟热</a-option>
            <a-option value="parvovirus">犬细小病毒</a-option>
            <a-option value="feline">猫瘟</a-option>
            <a-option value="other">其他</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="接种日期" required>
          <a-date-picker v-model="recordForm.inoculation_date" style="width: 100%" />
        </a-form-item>
        <a-row :gutter="12">
          <a-col :span="12">
            <a-form-item label="本次剂次">
              <a-input-number v-model="recordForm.dose_number" :min="1" :max="10" style="width: 100%" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="总剂次">
              <a-input-number v-model="recordForm.total_doses" :min="1" :max="10" style="width: 100%" />
            </a-form-item>
          </a-col>
        </a-row>
        <a-form-item label="下次接种日期">
          <a-date-picker v-model="recordForm.next_due_date" style="width: 100%" />
        </a-form-item>
        <a-form-item label="接种医院">
          <a-input v-model="recordForm.hospital" placeholder="请输入接种医院" />
        </a-form-item>
        <a-form-item label="兽医姓名">
          <a-input v-model="recordForm.vet_name" placeholder="请输入兽医姓名" />
        </a-form-item>
        <a-form-item label="备注">
          <a-textarea v-model="recordForm.notes" placeholder="请输入备注" :rows="3" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 提醒设置弹窗 -->
    <a-modal
      v-model:visible="reminderModalVisible"
      :title="isEditReminder ? '编辑提醒' : '添加提醒'"
      @ok="handleSaveReminder"
      :width="520"
      @close="resetReminderForm"
    >
      <a-form :model="reminderForm" layout="vertical">
        <a-form-item label="宠物" required>
          <a-select v-model="reminderForm.pet_id" placeholder="请选择宠物">
            <a-option v-for="p in pets" :key="p.id" :value="p.id">{{ p.name }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="疫苗名称" required>
          <a-input v-model="reminderForm.vaccine_name" placeholder="请输入疫苗名称" />
        </a-form-item>
        <a-form-item label="提醒日期" required>
          <a-date-picker v-model="reminderForm.reminder_date" style="width: 100%" />
        </a-form-item>
        <a-form-item label="提前提醒天数">
          <a-input-number v-model="reminderForm.advance_days" :min="1" :max="30" :step="1" style="width: 100%" />
          <span class="form-help-text">将在接种日期前 N 天发送提醒</span>
        </a-form-item>
        <a-form-item label="提醒方式">
          <a-checkbox-group v-model="reminderForm.channels">
            <a-checkbox value="app">App 推送</a-checkbox>
            <a-checkbox value="sms">短信</a-checkbox>
            <a-checkbox value="email">邮件</a-checkbox>
          </a-checkbox-group>
        </a-form-item>
        <a-form-item label="备注">
          <a-textarea v-model="reminderForm.notes" placeholder="请输入备注" :rows="2" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import {
  getVaccinationRecords,
  createVaccinationRecord,
  updateVaccinationRecord,
  deleteVaccinationRecord,
  getVaccinationReminders,
  setVaccinationReminder,
  updateVaccinationReminder,
  deleteVaccinationReminder
} from '@/api/advanced'
import { Message } from '@arco-design/web-vue'

const activeTab = ref('records')

// Mock pets data
const pets = ref([
  { id: 1, name: '小橘' },
  { id: 2, name: '豆豆' }
])

// ===== 接种记录 =====
const records = ref([])
const loading = ref(false)
const recordModalVisible = ref(false)
const isEditRecord = ref(false)
const savingId = ref(null)

const recordForm = reactive({
  id: null,
  pet_id: null,
  vaccine_name: '',
  vaccine_type: '',
  inoculation_date: null,
  dose_number: 1,
  total_doses: 1,
  next_due_date: null,
  hospital: '',
  vet_name: '',
  notes: ''
})

const recordColumns = [
  { title: '宠物', slotName: 'pet_name', width: 150 },
  { title: '疫苗名称', slotName: 'vaccine_name', width: 160 },
  { title: '疫苗类型', slotName: 'vaccine_type', width: 120 },
  { title: '剂次', slotName: 'dose_number', width: 80 },
  { title: '接种日期', slotName: 'inoculation_date', width: 130 },
  { title: '下次接种', slotName: 'next_due_date', width: 180 },
  { title: '兽医', slotName: 'vet_name', width: 120 },
  { title: '操作', slotName: 'actions', width: 120 }
]

const pagination = reactive({ current: 1, pageSize: 10, total: 0 })
const petFilter = ref(null)
const vaccineTypeFilter = ref(null)
const keyword = ref('')

async function loadRecords() {
  loading.value = true
  try {
    const params = {
      page: pagination.current,
      page_size: pagination.pageSize,
      pet_id: petFilter.value,
      vaccine_type: vaccineTypeFilter.value,
      keyword: keyword.value
    }
    const res = await getVaccinationRecords(params)
    records.value = res.data?.items || res.data || []
    pagination.total = res.data?.pagination?.total || records.value.length
  } catch (e) {
    Message.error('加载接种记录失败')
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
      vaccine_name: record.vaccine_name,
      vaccine_type: record.vaccine_type,
      inoculation_date: record.inoculation_date ? new Date(record.inoculation_date) : null,
      dose_number: record.dose_number,
      total_doses: record.total_doses,
      next_due_date: record.next_due_date ? new Date(record.next_due_date) : null,
      hospital: record.hospital,
      vet_name: record.vet_name,
      notes: record.notes
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
    vaccine_name: '',
    vaccine_type: '',
    inoculation_date: null,
    dose_number: 1,
    total_doses: 1,
    next_due_date: null,
    hospital: '',
    vet_name: '',
    notes: ''
  })
}

async function handleSaveRecord() {
  try {
    if (isEditRecord.value) {
      await updateVaccinationRecord(recordForm.id, recordForm)
      Message.success('更新成功')
    } else {
      await createVaccinationRecord(recordForm)
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
    await deleteVaccinationRecord(record.id)
    Message.success('删除成功')
    loadRecords()
  } catch (e) {
    Message.error('删除失败')
  }
}

// ===== 提醒设置 =====
const reminders = ref([])
const reminderLoading = ref(false)
const reminderModalVisible = ref(false)
const isEditReminder = ref(false)

const reminderForm = reactive({
  id: null,
  pet_id: null,
  vaccine_name: '',
  reminder_date: null,
  advance_days: 7,
  channels: ['app'],
  notes: ''
})

const reminderColumns = [
  { title: '宠物', slotName: 'pet_name', width: 130 },
  { title: '疫苗名称', slotName: 'vaccine_name', width: 160 },
  { title: '提醒日期', slotName: 'reminder_date', width: 130 },
  { title: '提前天数', slotName: 'advance_days', width: 110 },
  { title: '启用', slotName: 'enabled', width: 80 },
  { title: '操作', slotName: 'actions', width: 120 }
]

const reminderPagination = reactive({ current: 1, pageSize: 10, total: 0 })

async function loadReminders() {
  reminderLoading.value = true
  try {
    const params = {
      page: reminderPagination.current,
      page_size: reminderPagination.pageSize
    }
    const res = await getVaccinationReminders(params)
    reminders.value = res.data?.items || res.data || []
    reminderPagination.total = res.data?.pagination?.total || reminders.value.length
  } catch (e) {
    Message.error('加载提醒列表失败')
  } finally {
    reminderLoading.value = false
  }
}

function onReminderPageChange(page) {
  reminderPagination.current = page
  loadReminders()
}

function showReminderModal(reminder) {
  if (reminder) {
    isEditReminder.value = true
    Object.assign(reminderForm, {
      id: reminder.id,
      pet_id: reminder.pet_id,
      vaccine_name: reminder.vaccine_name,
      reminder_date: reminder.reminder_date ? new Date(reminder.reminder_date) : null,
      advance_days: reminder.advance_days,
      channels: reminder.channels || ['app'],
      notes: reminder.notes
    })
  } else {
    isEditReminder.value = false
    resetReminderForm()
  }
  reminderModalVisible.value = true
}

function resetReminderForm() {
  Object.assign(reminderForm, {
    id: null,
    pet_id: null,
    vaccine_name: '',
    reminder_date: null,
    advance_days: 7,
    channels: ['app'],
    notes: ''
  })
}

async function handleSaveReminder() {
  try {
    if (isEditReminder.value) {
      await updateVaccinationReminder(reminderForm.id, reminderForm)
      Message.success('更新成功')
    } else {
      await setVaccinationReminder(reminderForm)
      Message.success('添加成功')
    }
    reminderModalVisible.value = false
    loadReminders()
  } catch (e) {
    Message.error('保存失败')
  }
}

async function handleDeleteReminder(reminder) {
  try {
    await deleteVaccinationReminder(reminder.id)
    Message.success('删除成功')
    loadReminders()
  } catch (e) {
    Message.error('删除失败')
  }
}

async function toggleReminder(reminder) {
  savingId.value = reminder.id
  try {
    await updateVaccinationReminder(reminder.id, { enabled: reminder.enabled })
    Message.success(reminder.enabled ? '已启用' : '已禁用')
  } catch (e) {
    reminder.enabled = !reminder.enabled
    Message.error('操作失败')
  } finally {
    savingId.value = null
  }
}

// ===== 工具函数 =====
function formatDate(dateStr) {
  if (!dateStr) return '—'
  const d = new Date(dateStr)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
}

function getVaccineTypeColor(type) {
  const colors = {
    rabies: 'red',
    distemper: 'orange',
    parvovirus: 'blue',
    feline: 'purple',
    other: 'gray'
  }
  return colors[type] || 'gray'
}

function getVaccineTypeLabel(type) {
  const labels = {
    rabies: '狂犬疫苗',
    distemper: '犬瘟热',
    parvovirus: '犬细小',
    feline: '猫瘟',
    other: '其他'
  }
  return labels[type] || type
}

function isOverdue(dateStr) {
  if (!dateStr) return false
  return new Date(dateStr) < new Date()
}

function isDueSoon(dateStr) {
  if (!dateStr) return false
  const due = new Date(dateStr)
  const now = new Date()
  const diff = (due - now) / (1000 * 60 * 60 * 24)
  return diff >= 0 && diff <= 7
}

onMounted(() => {
  loadRecords()
  loadReminders()
})
</script>

<style scoped>
.vaccine-name {
  font-weight: 500;
}
.text-muted {
  color: var(--color-text-3);
}
.overdue-text {
  color: #f53f3f;
}
.form-help-text {
  font-size: 12px;
  color: var(--color-text-3);
  margin-top: 4px;
}
</style>
