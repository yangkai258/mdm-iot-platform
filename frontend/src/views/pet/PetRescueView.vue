<template>
  <div class="page-container">
    <a-card class="general-card" title="瀵诲洖缃戠粶">
      <template #extra>
        <a-button type="primary" @click="openReportModal"><icon-plus />鐧昏涓㈠け</a-button>
      </template>
      <a-tabs v-model:active-key="activeTab">
        <a-tab-pane key="lost" title="涓㈠け瀹犵墿">
          <div class="search-form">
            <a-form :model="form" layout="inline">
              <a-form-item label="瀹犵墿鍚嶇О"><a-input v-model="form.pet_name" placeholder="璇疯緭锟" /></a-form-item>
              <a-form-item label="鐘讹拷?>
                <a-select v-model="form.status" placeholder="閫夋嫨鐘讹拷" allow-clear style="width: 140px">
                  <a-option value="lost">瀵绘壘锟?/a-option>
                  <a-option value="found">宸叉壘锟?/a-option>
                  <a-option value="closed">宸插叧锟?/a-option>
                </a-select>
              </a-form-item>
              <a-form-item><a-button type="primary" @click="loadLost">鏌ヨ</a-button><a-button @click="Object.keys(form).forEach(k => form[k] = ''); loadLost()">閲嶇疆</a-button></a-form-item>
            </a-form>
          </div>
          <a-table :columns="lostColumns" :data="lostPets" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id">
            <template #status="{ record }">
              <a-tag :color="record.status === 'lost' ? 'red' : record.status === 'found' ? 'green' : 'gray'">{{ record.status === 'lost' ? '瀵绘壘锟? : record.status === 'found' ? '宸叉壘锟? : '宸插叧锟? }}</a-tag>
            </template>
            <template #reward="{ record }">
              <span v-if="record.reward_amount">楼{{ record.reward_amount }}</span>
              <span v-else>-</span>
            </template>
            <template #actions="{ record }">
              <a-button type="text" size="small" @click="viewLostDetail(record)">璇︽儏</a-button>
              <a-button v-if="record.status === 'lost'" type="text" size="small" @click="markFound(record)">鏍囪鎵惧埌</a-button>
              <a-button v-if="record.status === 'lost'" type="text" size="small" status="warning" @click="closeCase(record)">鍏抽棴</a-button>
            </template>
          </a-table>
        </a-tab-pane>
        <a-tab-pane key="volunteers" title="蹇楁効鑰呯綉锟?>
          <template #extra><a-button type="primary" size="small" @click="openVolunteerModal"><icon-plus />娣诲姞蹇楁効锟?/a-button></template>
          <a-table :columns="volunteerColumns" :data="volunteers" :loading="volLoading" :pagination="volPagination" @page-change="onVolPageChange" row-key="id">
            <template #status="{ record }">
              <a-tag :color="record.status === 'active' ? 'green' : 'gray'">{{ record.status === 'active' ? '娲昏穬' : '鍋滅敤' }}</a-tag>
            </template>
            <template #actions="{ record }">
              <a-button type="text" size="small" @click="toggleVolunteer(record)">{{ record.status === 'active' ? '鍋滅敤' : '鍚敤' }}</a-button>
            </template>
          </a-table>
        </a-tab-pane>
      </a-tabs>
    </a-card>
    <a-modal v-model:visible="reportVisible" title="鐧昏涓㈠け瀹犵墿" @before-ok="handleReport" :loading="submitting" :width="600">
      <a-form :model="reportForm" layout="vertical">
        <a-form-item label="瀹犵墿鍚嶇О" required><a-input v-model="reportForm.pet_name" placeholder="璇疯緭鍏ュ疇鐗╁悕锟" /></a-form-item>
        <a-row :gutter="16">
          <a-col :span="12"><a-form-item label="瀹犵墿绫诲瀷"><a-input v-model="reportForm.pet_type" placeholder="锟? 锟" /></a-form-item></a-col>
          <a-col :span="12"><a-form-item label="鍝佺"><a-input v-model="reportForm.breed" placeholder="锟? 閲戞瘺" /></a-form-item></a-col>
        </a-row>
        <a-form-item label="涓㈠け鍦扮偣"><a-input v-model="reportForm.last_location" placeholder="璇疯緭鍏ヤ涪澶卞湴锟" /></a-form-item>
        <a-form-item label="涓㈠け鏃堕棿"><a-date-picker v-model="reportForm.lost_time" style="width: 100%" /></a-form-item>
        <a-form-item label="鎮祻閲戦(锟?"><a-input-number v-model="reportForm.reward_amount" :min="0" placeholder="鍙拷" style="width: 200px" /></a-form-item>
        <a-form-item label="瀹犵墿鐗瑰緛鎻忚堪"><a-textarea v-model="reportForm.description" :rows="3" placeholder="璇锋弿杩板疇鐗╃壒锟" /></a-form-item>
        <a-form-item label="鑱旂郴鏂瑰紡" required><a-input v-model="reportForm.contact" placeholder="璇疯緭鍏ヨ仈绯绘柟锟" /></a-form-item>
      </a-form>
    </a-modal>
    <a-modal v-model:visible="detailVisible" title="涓㈠け璇︽儏" :footer="null" :width="560">
      <a-descriptions :column="1" bordered>
        <a-descriptions-item label="瀹犵墿鍚嶇О">{{ currentLost?.pet_name }}</a-descriptions-item>
        <a-descriptions-item label="瀹犵墿绫诲瀷">{{ currentLost?.pet_type }}</a-descriptions-item>
        <a-descriptions-item label="涓㈠け鍦扮偣">{{ currentLost?.last_location }}</a-descriptions-item>
        <a-descriptions-item label="涓㈠け鏃堕棿">{{ currentLost?.lost_time }}</a-descriptions-item>
        <a-descriptions-item label="鎮祻閲戦">{{ currentLost?.reward_amount ? `楼${currentLost.reward_amount}` : '-' }}</a-descriptions-item>
        <a-descriptions-item label="鐘讹拷?><a-tag :color="currentLost?.status === 'lost' ? 'red' : 'green'">{{ currentLost?.status }}</a-tag></a-descriptions-item>
        <a-descriptions-item label="鎻忚堪">{{ currentLost?.description }}</a-descriptions-item>
        <a-descriptions-item label="鐧昏鏃堕棿">{{ currentLost?.created_at }}</a-descriptions-item>
      </a-descriptions>
    </a-modal>
    <a-modal v-model:visible="volunteerVisible" title="娣诲姞蹇楁効锟? @before-ok="handleAddVolunteer" :loading="submitting" :width="480">
      <a-form :model="volForm" layout="vertical">
        <a-form-item label="蹇楁効鑰呭锟? required><a-input v-model="volForm.name" placeholder="璇疯緭锟" /></a-form-item>
        <a-form-item label="鑱旂郴鐢佃瘽" required><a-input v-model="volForm.phone" placeholder="璇疯緭锟" /></a-form-item>
        <a-form-item label="鎵€鍦ㄥ尯锟?><a-input v-model="volForm.area" placeholder="锟? 鍖椾含鏈濋槼" /></a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const activeTab = ref('lost')
const loading = ref(false)
const volLoading = ref(false)
const submitting = ref(false)
const lostPets = ref([])
const volunteers = ref([])
const reportVisible = ref(false)
const detailVisible = ref(false)
const volunteerVisible = ref(false)
const currentLost = ref(null)
const form = reactive({ pet_name: '', status: '' })
const reportForm = reactive({ pet_name: '', pet_type: '', breed: '', last_location: '', lost_time: null, reward_amount: null, description: '', contact: '' })
const volForm = reactive({ name: '', phone: '', area: '' })
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })
const volPagination = reactive({ current: 1, pageSize: 20, total: 0 })

const lostColumns = [
  { title: 'ID', dataIndex: 'id', width: 70 },
  { title: '瀹犵墿鍚嶇О', dataIndex: 'pet_name', width: 120 },
  { title: '绫诲瀷', dataIndex: 'pet_type', width: 80 },
  { title: '涓㈠け鍦扮偣', dataIndex: 'last_location', ellipsis: true },
  { title: '涓㈠け鏃堕棿', dataIndex: 'lost_time', width: 170 },
  { title: '鎮祻', slotName: 'reward', width: 100 },
  { title: '鐘讹拷?, slotName: 'status', width: 90 },
  { title: '鐧昏鏃堕棿', dataIndex: 'created_at', width: 170 },
  { title: '鎿嶄綔', slotName: 'actions', width: 160 }
]

const volunteerColumns = [
  { title: 'ID', dataIndex: 'id', width: 70 },
  { title: '濮撳悕', dataIndex: 'name', width: 120 },
  { title: '鐢佃瘽', dataIndex: 'phone', width: 130 },
  { title: '鍖哄煙', dataIndex: 'area', width: 140 },
  { title: '鐘讹拷?, slotName: 'status', width: 80 },
  { title: '鍔犲叆鏃堕棿', dataIndex: 'created_at', width: 170 },
  { title: '鎿嶄綔', slotName: 'actions', width: 80 }
]

const loadLost = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const params = new URLSearchParams({ page: pagination.current, page_size: pagination.pageSize })
    if (form.pet_name) params.append('pet_name', form.pet_name)
    if (form.status) params.append('status', form.status)
    const res = await fetch(`/api/v1/pet/rescue/lost?${params}`, { headers: { 'Authorization': `Bearer ${token}` } }).then(r => r.json())
    if (res.code === 0) { lostPets.value = res.data?.list || []; pagination.total = res.data?.total || 0 }
    else { lostPets.value = [] }
  } catch (e) { Message.error('鍔犺浇澶辫触') } finally { loading.value = false }
}

const loadVolunteers = async () => {
  volLoading.value = true
  try {
    const token = localStorage.getItem('token')
    const params = new URLSearchParams({ page: volPagination.current, page_size: volPagination.pageSize })
    const res = await fetch(`/api/v1/pet/rescue/volunteers?${params}`, { headers: { 'Authorization': `Bearer ${token}` } }).then(r => r.json())
    if (res.code === 0) { volunteers.value = res.data?.list || []; volPagination.total = res.data?.total || 0 }
    else { volunteers.value = [] }
  } catch (e) { Message.error('鍔犺浇蹇楁効鑰呭け锟?) } finally { volLoading.value = false }
}

const openReportModal = () => { Object.assign(reportForm, { pet_name: '', pet_type: '', breed: '', last_location: '', lost_time: null, reward_amount: null, description: '', contact: '' }); reportVisible.value = true }

const handleReport = async (done) => {
  if (!reportForm.pet_name || !reportForm.contact) { Message.warning('璇峰～鍐欏繀濉」'); done(false); return }
  submitting.value = true
  try {
    const token = localStorage.getItem('token')
    const res = await fetch('/api/v1/pet/rescue/lost', { method: 'POST', headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' }, body: JSON.stringify(reportForm) }).then(r => r.json())
    if (res.code === 0) { Message.success('鐧昏鎴愬姛'); reportVisible.value = false; loadLost() }
    else { Message.error(res.message || '鐧昏澶辫触') }
    done(true)
  } catch (e) { Message.error('鐧昏澶辫触'); done(false) } finally { submitting.value = false }
}

const viewLostDetail = (record) => { currentLost.value = record; detailVisible.value = true }
const markFound = async (record) => {
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`/api/v1/pet/rescue/lost/${record.id}/found`, { method: 'POST', headers: { 'Authorization': `Bearer ${token}` } }).then(r => r.json())
    if (res.code === 0) { Message.success('宸叉爣璁颁负鎵惧埌'); loadLost() }
    else Message.error('鎿嶄綔澶辫触')
  } catch (e) { Message.error('鎿嶄綔澶辫触') }
}
const closeCase = async (record) => {
  try {
    const token = localStorage.getItem('token')
    await fetch(`/api/v1/pet/rescue/lost/${record.id}/close`, { method: 'POST', headers: { 'Authorization': `Bearer ${token}` } })
    Message.success('宸插叧锟?); loadLost()
  } catch (e) { Message.error('鎿嶄綔澶辫触') }
}

const openVolunteerModal = () => { Object.assign(volForm, { name: '', phone: '', area: '' }); volunteerVisible.value = true }
const handleAddVolunteer = async (done) => {
  if (!volForm.name || !volForm.phone) { Message.warning('璇峰～鍐欏繀濉」'); done(false); return }
  submitting.value = true
  try {
    const token = localStorage.getItem('token')
    const res = await fetch('/api/v1/pet/rescue/volunteers', { method: 'POST', headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' }, body: JSON.stringify(volForm) }).then(r => r.json())
    if (res.code === 0) { Message.success('娣诲姞鎴愬姛'); volunteerVisible.value = false; loadVolunteers() }
    else { Message.error(res.message || '娣诲姞澶辫触') }
    done(true)
  } catch (e) { Message.error('娣诲姞澶辫触'); done(false) } finally { submitting.value = false }
}

const toggleVolunteer = async (record) => {
  try {
    const token = localStorage.getItem('token')
    const newStatus = record.status === 'active' ? 'inactive' : 'active'
    await fetch(`/api/v1/pet/rescue/volunteers/${record.id}`, { method: 'PUT', headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' }, body: JSON.stringify({ status: newStatus }) })
    Message.success('鏇存柊鎴愬姛'); loadVolunteers()
  } catch (e) { Message.error('鎿嶄綔澶辫触') }
}

const onPageChange = (page) => { pagination.current = page; loadLost() }
const onVolPageChange = (page) => { volPagination.current = page; loadVolunteers() }

onMounted(() => { loadLost(); loadVolunteers() })
</script>

<style scoped>
.page-container { padding: 16px; }
.search-form { margin-bottom: 16px; padding: 16px; background: var(--color-fill-lightest); border-radius: 4px; }
</style>