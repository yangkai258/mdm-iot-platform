<template>
  <div class="page-container">
    <a-card class="general-card" title="寻回网络">
      <template #extra>
        <a-button type="primary" @click="openReportModal"><icon-plus />登记丢失</a-button>
      </template>
      <a-tabs v-model:active-key="activeTab">
        <a-tab-pane key="lost" title="丢失宠物">
          <div class="search-form">
            <a-form :model="form" layout="inline">
              <a-form-item label="宠物名称"><a-input v-model="form.pet_name" placeholder="请输�? /></a-form-item>
              <a-form-item label="状�?>
                <a-select v-model="form.status" placeholder="选择状�? allow-clear style="width: 140px">
                  <a-option value="lost">寻找�?/a-option>
                  <a-option value="found">已找�?/a-option>
                  <a-option value="closed">已关�?/a-option>
                </a-select>
              </a-form-item>
              <a-form-item><a-button type="primary" @click="loadLost">查询</a-button><a-button @click="Object.keys(form).forEach(k => form[k] = ''); loadLost()">重置</a-button></a-form-item>
            </a-form>
          </div>
          <a-table :columns="lostColumns" :data="lostPets" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id">
            <template #status="{ record }">
              <a-tag :color="record.status === 'lost' ? 'red' : record.status === 'found' ? 'green' : 'gray'">{{ record.status === 'lost' ? '寻找�? : record.status === 'found' ? '已找�? : '已关�? }}</a-tag>
            </template>
            <template #reward="{ record }">
              <span v-if="record.reward_amount">¥{{ record.reward_amount }}</span>
              <span v-else>-</span>
            </template>
            <template #actions="{ record }">
              <a-button type="text" size="small" @click="viewLostDetail(record)">详情</a-button>
              <a-button v-if="record.status === 'lost'" type="text" size="small" @click="markFound(record)">标记找到</a-button>
              <a-button v-if="record.status === 'lost'" type="text" size="small" status="warning" @click="closeCase(record)">关闭</a-button>
            </template>
          </a-table>
        </a-tab-pane>
        <a-tab-pane key="volunteers" title="志愿者网�?>
          <template #extra><a-button type="primary" size="small" @click="openVolunteerModal"><icon-plus />添加志愿�?/a-button></template>
          <a-table :columns="volunteerColumns" :data="volunteers" :loading="volLoading" :pagination="volPagination" @page-change="onVolPageChange" row-key="id">
            <template #status="{ record }">
              <a-tag :color="record.status === 'active' ? 'green' : 'gray'">{{ record.status === 'active' ? '活跃' : '停用' }}</a-tag>
            </template>
            <template #actions="{ record }">
              <a-button type="text" size="small" @click="toggleVolunteer(record)">{{ record.status === 'active' ? '停用' : '启用' }}</a-button>
            </template>
          </a-table>
        </a-tab-pane>
      </a-tabs>
    </a-card>
    <a-modal v-model:visible="reportVisible" title="登记丢失宠物" @before-ok="handleReport" :loading="submitting" :width="600">
      <a-form :model="reportForm" layout="vertical">
        <a-form-item label="宠物名称" required><a-input v-model="reportForm.pet_name" placeholder="请输入宠物名�? /></a-form-item>
        <a-row :gutter="16">
          <a-col :span="12"><a-form-item label="宠物类型"><a-input v-model="reportForm.pet_type" placeholder="�? �? /></a-form-item></a-col>
          <a-col :span="12"><a-form-item label="品种"><a-input v-model="reportForm.breed" placeholder="�? 金毛" /></a-form-item></a-col>
        </a-row>
        <a-form-item label="丢失地点"><a-input v-model="reportForm.last_location" placeholder="请输入丢失地�? /></a-form-item>
        <a-form-item label="丢失时间"><a-date-picker v-model="reportForm.lost_time" style="width: 100%" /></a-form-item>
        <a-form-item label="悬赏金额(�?"><a-input-number v-model="reportForm.reward_amount" :min="0" placeholder="可�? style="width: 200px" /></a-form-item>
        <a-form-item label="宠物特征描述"><a-textarea v-model="reportForm.description" :rows="3" placeholder="请描述宠物特�? /></a-form-item>
        <a-form-item label="联系方式" required><a-input v-model="reportForm.contact" placeholder="请输入联系方�? /></a-form-item>
      </a-form>
    </a-modal>
    <a-modal v-model:visible="detailVisible" title="丢失详情" :footer="null" :width="560">
      <a-descriptions :column="1" bordered>
        <a-descriptions-item label="宠物名称">{{ currentLost?.pet_name }}</a-descriptions-item>
        <a-descriptions-item label="宠物类型">{{ currentLost?.pet_type }}</a-descriptions-item>
        <a-descriptions-item label="丢失地点">{{ currentLost?.last_location }}</a-descriptions-item>
        <a-descriptions-item label="丢失时间">{{ currentLost?.lost_time }}</a-descriptions-item>
        <a-descriptions-item label="悬赏金额">{{ currentLost?.reward_amount ? `¥${currentLost.reward_amount}` : '-' }}</a-descriptions-item>
        <a-descriptions-item label="状�?><a-tag :color="currentLost?.status === 'lost' ? 'red' : 'green'">{{ currentLost?.status }}</a-tag></a-descriptions-item>
        <a-descriptions-item label="描述">{{ currentLost?.description }}</a-descriptions-item>
        <a-descriptions-item label="登记时间">{{ currentLost?.created_at }}</a-descriptions-item>
      </a-descriptions>
    </a-modal>
    <a-modal v-model:visible="volunteerVisible" title="添加志愿�? @before-ok="handleAddVolunteer" :loading="submitting" :width="480">
      <a-form :model="volForm" layout="vertical">
        <a-form-item label="志愿者姓�? required><a-input v-model="volForm.name" placeholder="请输�? /></a-form-item>
        <a-form-item label="联系电话" required><a-input v-model="volForm.phone" placeholder="请输�? /></a-form-item>
        <a-form-item label="所在区�?><a-input v-model="volForm.area" placeholder="�? 北京朝阳" /></a-form-item>
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
  { title: '宠物名称', dataIndex: 'pet_name', width: 120 },
  { title: '类型', dataIndex: 'pet_type', width: 80 },
  { title: '丢失地点', dataIndex: 'last_location', ellipsis: true },
  { title: '丢失时间', dataIndex: 'lost_time', width: 170 },
  { title: '悬赏', slotName: 'reward', width: 100 },
  { title: '状�?, slotName: 'status', width: 90 },
  { title: '登记时间', dataIndex: 'created_at', width: 170 },
  { title: '操作', slotName: 'actions', width: 160 }
]

const volunteerColumns = [
  { title: 'ID', dataIndex: 'id', width: 70 },
  { title: '姓名', dataIndex: 'name', width: 120 },
  { title: '电话', dataIndex: 'phone', width: 130 },
  { title: '区域', dataIndex: 'area', width: 140 },
  { title: '状�?, slotName: 'status', width: 80 },
  { title: '加入时间', dataIndex: 'created_at', width: 170 },
  { title: '操作', slotName: 'actions', width: 80 }
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
  } catch (e) { Message.error('加载失败') } finally { loading.value = false }
}

const loadVolunteers = async () => {
  volLoading.value = true
  try {
    const token = localStorage.getItem('token')
    const params = new URLSearchParams({ page: volPagination.current, page_size: volPagination.pageSize })
    const res = await fetch(`/api/v1/pet/rescue/volunteers?${params}`, { headers: { 'Authorization': `Bearer ${token}` } }).then(r => r.json())
    if (res.code === 0) { volunteers.value = res.data?.list || []; volPagination.total = res.data?.total || 0 }
    else { volunteers.value = [] }
  } catch (e) { Message.error('加载志愿者失�?) } finally { volLoading.value = false }
}

const openReportModal = () => { Object.assign(reportForm, { pet_name: '', pet_type: '', breed: '', last_location: '', lost_time: null, reward_amount: null, description: '', contact: '' }); reportVisible.value = true }

const handleReport = async (done) => {
  if (!reportForm.pet_name || !reportForm.contact) { Message.warning('请填写必填项'); done(false); return }
  submitting.value = true
  try {
    const token = localStorage.getItem('token')
    const res = await fetch('/api/v1/pet/rescue/lost', { method: 'POST', headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' }, body: JSON.stringify(reportForm) }).then(r => r.json())
    if (res.code === 0) { Message.success('登记成功'); reportVisible.value = false; loadLost() }
    else { Message.error(res.message || '登记失败') }
    done(true)
  } catch (e) { Message.error('登记失败'); done(false) } finally { submitting.value = false }
}

const viewLostDetail = (record) => { currentLost.value = record; detailVisible.value = true }
const markFound = async (record) => {
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`/api/v1/pet/rescue/lost/${record.id}/found`, { method: 'POST', headers: { 'Authorization': `Bearer ${token}` } }).then(r => r.json())
    if (res.code === 0) { Message.success('已标记为找到'); loadLost() }
    else Message.error('操作失败')
  } catch (e) { Message.error('操作失败') }
}
const closeCase = async (record) => {
  try {
    const token = localStorage.getItem('token')
    await fetch(`/api/v1/pet/rescue/lost/${record.id}/close`, { method: 'POST', headers: { 'Authorization': `Bearer ${token}` } })
    Message.success('已关�?); loadLost()
  } catch (e) { Message.error('操作失败') }
}

const openVolunteerModal = () => { Object.assign(volForm, { name: '', phone: '', area: '' }); volunteerVisible.value = true }
const handleAddVolunteer = async (done) => {
  if (!volForm.name || !volForm.phone) { Message.warning('请填写必填项'); done(false); return }
  submitting.value = true
  try {
    const token = localStorage.getItem('token')
    const res = await fetch('/api/v1/pet/rescue/volunteers', { method: 'POST', headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' }, body: JSON.stringify(volForm) }).then(r => r.json())
    if (res.code === 0) { Message.success('添加成功'); volunteerVisible.value = false; loadVolunteers() }
    else { Message.error(res.message || '添加失败') }
    done(true)
  } catch (e) { Message.error('添加失败'); done(false) } finally { submitting.value = false }
}

const toggleVolunteer = async (record) => {
  try {
    const token = localStorage.getItem('token')
    const newStatus = record.status === 'active' ? 'inactive' : 'active'
    await fetch(`/api/v1/pet/rescue/volunteers/${record.id}`, { method: 'PUT', headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' }, body: JSON.stringify({ status: newStatus }) })
    Message.success('更新成功'); loadVolunteers()
  } catch (e) { Message.error('操作失败') }
}

const onPageChange = (page) => { pagination.current = page; loadLost() }
const onVolPageChange = (page) => { volPagination.current = page; loadVolunteers() }

onMounted(() => { loadLost(); loadVolunteers() })
</script>

<style scoped>
.page-container { padding: 16px; }
.search-form { margin-bottom: 16px; padding: 16px; background: var(--color-fill-lightest); border-radius: 4px; }
</style>