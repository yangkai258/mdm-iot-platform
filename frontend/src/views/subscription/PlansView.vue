<template>
  <div class="page-container">
    <a-card class="general-card" title="璁㈤槄濂楅">
      <template #extra>
        <a-button type="primary" @click="openCreateModal"><icon-plus />鏂板缓濂楅</a-button>
      </template>
      <div class="search-form">
        <a-form :model="form" layout="inline">
          <a-form-item label="濂楅鍚嶇О"><a-input v-model="form.plan_name" placeholder="璇疯緭鍏" /></a-form-item>
          <a-form-item label="鐘舵€?>
            <a-select v-model="form.status" placeholder="閫夋嫨鐘舵€" allow-clear style="width: 120px">
              <a-option value="active">鐢熸晥涓?/a-option>
              <a-option value="inactive">鍋滅敤</a-option>
            </a-select>
          </a-form-item>
          <a-form-item><a-button type="primary" @click="loadData">鏌ヨ</a-button><a-button @click="Object.keys(form).forEach(k => form[k] = ''); loadData()">閲嶇疆</a-button></a-form-item>
        </a-form>
      </div>
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id">
        <template #price="{ record }">
          <span style="color: #f53f3f; font-weight: 600">楼{{ record.price }}</span>
          <span style="color: #999; font-size: 12px">/{{ record.billing_cycle === 'monthly' ? '鏈? : '骞? }}</span>
        </template>
        <template #features="{ record }">
          <a-tooltip :content="record.features?.join('\n') || '鏃?" placement="top">
            <span class="features-text">{{ record.features?.length || 0 }} 椤瑰姛鑳?/span>
          </a-tooltip>
        </template>
        <template #status="{ record }">
          <a-tag :color="record.status === 'active' ? 'green' : 'gray'">{{ record.status === 'active' ? '鐢熸晥涓? : '鍋滅敤' }}</a-tag>
        </template>
        <template #actions="{ record }">
          <a-button type="text" size="small" @click="editPlan(record)">缂栬緫</a-button>
          <a-button type="text" size="small" @click="viewSubscribers(record)">璁㈤槄鑰?/a-button>
          <a-button type="text" size="small" status="danger" @click="deletePlan(record)">鍒犻櫎</a-button>
        </template>
      </a-table>
    </a-card>
    <!-- 鏂板缓/缂栬緫濂楅寮圭獥 -->
    <a-modal v-model:visible="formVisible" :title="isEdit ? '缂栬緫濂楅' : '鏂板缓濂楅'" @before-ok="handleSubmit" :loading="submitting" :width="600">
      <a-form :model="planForm" layout="vertical">
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="濂楅鍚嶇О" required><a-input v-model="planForm.plan_name" placeholder="璇疯緭鍏ュ椁愬悕绉" /></a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="濂楅浠ｇ爜"><a-input v-model="planForm.plan_code" placeholder="濡" basic" /></a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="浠锋牸(鍏?" required><a-input-number v-model="planForm.price" :min="0" :precision="2" style="width: 100%" /></a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="璁¤垂鍛ㄦ湡" required>
              <a-select v-model="planForm.billing_cycle" placeholder="閫夋嫨鍛ㄦ湡">
                <a-option value="monthly">鏈堜粯</a-option>
                <a-option value="yearly">骞翠粯</a-option>
                <a-option value="one-time">涓€娆℃€?/a-option>
              </a-select>
            </a-form-item>
          </a-col>
        </a-row>
        <a-form-item label="濂楅鎻忚堪"><a-textarea v-model="planForm.description" :rows="2" placeholder="濂楅鎻忚堪" /></a-form-item>
        <a-form-item label="濂楅鍔熻兘">
          <a-select v-model="planForm.features" multiple placeholder="閫夋嫨鍖呭惈鐨勫姛鑳" allow-create style="width: 100%">
            <a-option value="device_management">璁惧绠＄悊</a-option>
            <a-option value="ai_features">AI鍔熻兘</a-option>
            <a-option value="ota_upgrade">OTA鍗囩骇</a-option>
            <a-option value="analytics">鏁版嵁鍒嗘瀽</a-option>
            <a-option value="api_access">API璁块棶</a-option>
            <a-option value="priority_support">浼樺厛鏀寔</a-option>
            <a-option value="custom_branding">鑷畾涔夊搧鐗?/a-option>
            <a-option value="multi_user">澶氱敤鎴?/a-option>
          </a-select>
        </a-form-item>
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="璁惧涓婇檺"><a-input-number v-model="planForm.max_devices" :min="1" placeholder="鏃犻檺鍒" style="width: 100%" /></a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="瀛樺偍涓婇檺(GB)"><a-input-number v-model="planForm.max_storage_gb" :min="0" placeholder="鏃犻檺鍒" style="width: 100%" /></a-form-item>
          </a-col>
        </a-row>
        <a-form-item label="鐘舵€?><a-switch v-model="planForm.is_active" /></a-form-item>
      </a-form>
    </a-modal>
    <!-- 璁㈤槄鑰呭垪琛?-->
    <a-modal v-model:visible="subscriberVisible" title="璁㈤槄鑰呭垪琛? :width="700" :footer="null">
      <a-table :columns="subColumns" :data="subscribers" :loading="subLoading" :pagination="subPagination" @page-change="onSubPageChange" row-key="id">
        <template #status="{ record }">
          <a-tag :color="record.subscription_status === 'active' ? 'green' : 'gray'">{{ record.subscription_status === 'active' ? '娲昏穬' : '宸茶繃鏈? }}</a-tag>
        </template>
        <template #expires_at="{ record }">{{ record.expires_at || '姘镐箙' }}</template>
      </a-table>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const loading = ref(false)
const submitting = ref(false)
const subLoading = ref(false)
const data = ref([])
const subscribers = ref([])
const formVisible = ref(false)
const subscriberVisible = ref(false)
const isEdit = ref(false)
const selectedPlan = ref(null)
const form = reactive({ plan_name: '', status: '' })
const planForm = reactive({ id: null, plan_name: '', plan_code: '', price: 0, billing_cycle: 'monthly', description: '', features: [], max_devices: null, max_storage_gb: null, is_active: true })
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })
const subPagination = reactive({ current: 1, pageSize: 10, total: 0 })

const columns = [
  { title: 'ID', dataIndex: 'id', width: 70 },
  { title: '濂楅鍚嶇О', dataIndex: 'plan_name', width: 160 },
  { title: '濂楅浠ｇ爜', dataIndex: 'plan_code', width: 100 },
  { title: '浠锋牸', slotName: 'price', width: 120 },
  { title: '璁¤垂鍛ㄦ湡', dataIndex: 'billing_cycle', width: 100 },
  { title: '鍔熻兘鏁?, slotName: 'features', width: 100 },
  { title: '璁惧涓婇檺', dataIndex: 'max_devices', width: 100 },
  { title: '鐘舵€?, slotName: 'status', width: 80 },
  { title: '鍒涘缓鏃堕棿', dataIndex: 'created_at', width: 170 },
  { title: '鎿嶄綔', slotName: 'actions', width: 160 }
]

const subColumns = [
  { title: '鐢ㄦ埛', dataIndex: 'user_name', width: 120 },
  { title: '閭', dataIndex: 'user_email', ellipsis: true },
  { title: '璁㈤槄鐘舵€?, slotName: 'status', width: 90 },
  { title: '寮€濮嬫椂闂?, dataIndex: 'started_at', width: 170 },
  { title: '鍒版湡鏃堕棿', slotName: 'expires_at', width: 170 }
]

const loadData = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const params = new URLSearchParams({ page: pagination.current, page_size: pagination.pageSize })
    if (form.plan_name) params.append('plan_name', form.plan_name)
    if (form.status) params.append('status', form.status)
    const res = await fetch(`/api/v1/subscription/plans?${params}`, { headers: { 'Authorization': `Bearer ${token}` } }).then(r => r.json())
    if (res.code === 0) { data.value = res.data?.list || []; pagination.total = res.data?.total || 0 }
    else { data.value = [] }
  } catch (e) { Message.error('鍔犺浇澶辫触') } finally { loading.value = false }
}

const openCreateModal = () => { isEdit.value = false; Object.assign(planForm, { id: null, plan_name: '', plan_code: '', price: 0, billing_cycle: 'monthly', description: '', features: [], max_devices: null, max_storage_gb: null, is_active: true }); formVisible.value = true }
const editPlan = (record) => { isEdit.value = true; Object.assign(planForm, record); planForm.features = record.features || []; formVisible.value = true }

const handleSubmit = async (done) => {
  if (!planForm.plan_name) { Message.warning('璇疯緭鍏ュ椁愬悕绉?); done(false); return }
  submitting.value = true
  try {
    const token = localStorage.getItem('token')
    const url = isEdit.value ? `/api/v1/subscription/plans/${planForm.id}` : '/api/v1/subscription/plans'
    const res = await fetch(url, { method: isEdit.value ? 'PUT' : 'POST', headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' }, body: JSON.stringify(planForm) }).then(r => r.json())
    if (res.code === 0) { Message.success(isEdit.value ? '鏇存柊鎴愬姛' : '鍒涘缓鎴愬姛'); formVisible.value = false; loadData() }
    else { Message.error(res.message || '鎿嶄綔澶辫触') }
    done(true)
  } catch (e) { Message.error('鎿嶄綔澶辫触'); done(false) } finally { submitting.value = false }
}

const deletePlan = async (record) => {
  try {
    const token = localStorage.getItem('token')
    await fetch(`/api/v1/subscription/plans/${record.id}`, { method: 'DELETE', headers: { 'Authorization': `Bearer ${token}` } })
    Message.success('鍒犻櫎鎴愬姛'); loadData()
  } catch (e) { Message.error('鍒犻櫎澶辫触') }
}

const viewSubscribers = async (record) => {
  selectedPlan.value = record
  subscriberVisible.value = true
  subLoading.value = true
  try {
    const token = localStorage.getItem('token')
    const params = new URLSearchParams({ plan_id: record.id, page: subPagination.current, page_size: subPagination.pageSize })
    const res = await fetch(`/api/v1/subscription/subscribers?${params}`, { headers: { 'Authorization': `Bearer ${token}` } }).then(r => r.json())
    if (res.code === 0) { subscribers.value = res.data?.list || []; subPagination.total = res.data?.total || 0 }
    else { subscribers.value = [] }
  } catch (e) { Message.error('鍔犺浇澶辫触') } finally { subLoading.value = false }
}

const onPageChange = (page) => { pagination.current = page; loadData() }
const onSubPageChange = (page) => { subPagination.current = page; viewSubscribers(selectedPlan.value) }

onMounted(() => loadData())
</script>

<style scoped>
.page-container { padding: 16px; }
.search-form { margin-bottom: 16px; padding: 16px; background: var(--color-fill-lightest); border-radius: 4px; }
.features-text { cursor: pointer; color: #165dff; }
</style>
