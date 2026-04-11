<template>
  <div class="page-container">
    <a-card class="general-card" title="璺ㄥ璁惧绠℃帶">
      <template #extra>
        <a-button type="primary" @click="openPolicyModal"><icon-plus />鏂板缓绛栫暐</a-button>
      </template>
      <a-tabs v-model:active-key="activeTab">
        <a-tab-pane key="policies" title="鍖哄煙绛栫暐">
          <div class="search-form">
            <a-form :model="policyForm" layout="inline">
              <a-form-item label="绛栫暐鍚嶇О"><a-input v-model="policyForm.name" placeholder="璇疯緭锟" /></a-form-item>
              <a-form-item label="鐘讹拷?>
                <a-select v-model="policyForm.status" placeholder="閫夋嫨鐘讹拷" allow-clear style="width: 120px">
                  <a-option value="active">鍚敤</a-option>
                  <a-option value="inactive">鍋滅敤</a-option>
                </a-select>
              </a-form-item>
              <a-form-item><a-button type="primary" @click="loadPolicies">鏌ヨ</a-button><a-button @click="Object.keys(policyForm).forEach(k => policyForm[k] = ''); loadPolicies()">閲嶇疆</a-button></a-form-item>
            </a-form>
          </div>
          <a-table :columns="policyColumns" :data="policies" :loading="loading" :pagination="pagination" @page-change="onPolicyPageChange" row-key="id">
            <template #status="{ record }">
              <a-tag :color="record.status === 'active' ? 'green' : 'gray'">{{ record.status === 'active' ? '鍚敤' : '鍋滅敤' }}</a-tag>
            </template>
            <template #actions="{ record }">
              <a-button type="text" size="small" @click="editPolicy(record)">缂栬緫</a-button>
              <a-button type="text" size="small" status="danger" @click="deletePolicy(record)">鍒犻櫎</a-button>
            </template>
          </a-table>
        </a-tab-pane>
        <a-tab-pane key="alerts" title="璺ㄥ鍛婅">
          <div class="search-form">
            <a-form :model="alertForm" layout="inline">
              <a-form-item label="璁惧ID"><a-input v-model="alertForm.device_id" placeholder="璇疯緭锟" /></a-form-item>
              <a-form-item label="鍛婅绫诲瀷">
                <a-select v-model="alertForm.alert_type" placeholder="閫夋嫨绫诲瀷" allow-clear style="width: 140px">
                  <a-option value="cross_in">璺ㄥ杩涘叆</a-option>
                  <a-option value="cross_out">璺ㄥ绂诲紑</a-option>
                  <a-option value="region_change">鍖哄煙鍙樻洿</a-option>
                </a-select>
              </a-form-item>
              <a-form-item label="鏃堕棿鑼冨洿"><a-range-picker v-model="alertForm.time_range" style="width: 240px" /></a-form-item>
              <a-form-item><a-button type="primary" @click="loadAlerts">鏌ヨ</a-button><a-button @click="Object.assign(alertForm, { device_id: '', alert_type: '', time_range: [] }); loadAlerts()">閲嶇疆</a-button></a-form-item>
            </a-form>
          </div>
          <a-table :columns="alertColumns" :data="alerts" :loading="alertsLoading" :pagination="alertPagination" @page-change="onAlertPageChange" row-key="id">
            <template #alert_type="{ record }">
              <a-tag :color="record.alert_type === 'cross_in' ? 'green' : 'orange'">{{ record.alert_type === 'cross_in' ? '璺ㄥ杩涘叆' : record.alert_type === 'cross_out' ? '璺ㄥ绂诲紑' : '鍖哄煙鍙樻洿' }}</a-tag>
            </template>
          </a-table>
        </a-tab-pane>
      </a-tabs>
    </a-card>
    <a-modal v-model:visible="policyModalVisible" :title="isEditPolicy ? '缂栬緫绛栫暐' : '鏂板缓绛栫暐'" @before-ok="handlePolicySubmit" :loading="submitting" :width="560">
      <a-form :model="policyFormData" layout="vertical">
        <a-form-item label="绛栫暐鍚嶇О" required><a-input v-model="policyFormData.name" placeholder="璇疯緭鍏ョ瓥鐣ュ悕锟" /></a-form-item>
        <a-form-item label="婧愬尯锟?>
          <a-select v-model="policyFormData.source_region" placeholder="閫夋嫨婧愬尯锟">
            <a-option value="cn">涓浗澶ч檰</a-option>
            <a-option value="us">缇庡浗</a-option>
            <a-option value="eu">娆х洘</a-option>
            <a-option value="other">鍏朵粬</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="鐩爣鍖哄煙">
          <a-select v-model="policyFormData.target_region" placeholder="閫夋嫨鐩爣鍖哄煙">
            <a-option value="cn">涓浗澶ч檰</a-option>
            <a-option value="us">缇庡浗</a-option>
            <a-option value="eu">娆х洘</a-option>
            <a-option value="other">鍏朵粬</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="绠℃帶鍔ㄤ綔">
          <a-select v-model="policyFormData.action" placeholder="閫夋嫨绠℃帶鍔ㄤ綔">
            <a-option value="allow">鍏佽</a-option>
            <a-option value="warn">鍛婅</a-option>
            <a-option value="block">闃绘</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="鐘讹拷?><a-switch v-model="policyFormData.is_active" /></a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const activeTab = ref('policies')
const loading = ref(false)
const alertsLoading = ref(false)
const submitting = ref(false)
const policies = ref([])
const alerts = ref([])
const policyModalVisible = ref(false)
const isEditPolicy = ref(false)
const selectedPolicy = ref(null)

const policyForm = reactive({ name: '', status: '' })
const alertForm = reactive({ device_id: '', alert_type: '', time_range: [] })
const policyFormData = reactive({ id: null, name: '', source_region: '', target_region: '', action: 'warn', is_active: true })

const pagination = reactive({ current: 1, pageSize: 20, total: 0 })
const alertPagination = reactive({ current: 1, pageSize: 20, total: 0 })

const policyColumns = [
  { title: 'ID', dataIndex: 'id', width: 70 },
  { title: '绛栫暐鍚嶇О', dataIndex: 'name', width: 180 },
  { title: '婧愬尯锟?, dataIndex: 'source_region', width: 120 },
  { title: '鐩爣鍖哄煙', dataIndex: 'target_region', width: 120 },
  { title: '绠℃帶鍔ㄤ綔', dataIndex: 'action', width: 100 },
  { title: '鐘讹拷?, slotName: 'status', width: 80 },
  { title: '鍒涘缓鏃堕棿', dataIndex: 'created_at', width: 170 },
  { title: '鎿嶄綔', slotName: 'actions', width: 120 }
]

const alertColumns = [
  { title: '鍛婅ID', dataIndex: 'id', width: 80 },
  { title: '璁惧ID', dataIndex: 'device_id', width: 120 },
  { title: '璁惧鍚嶇О', dataIndex: 'device_name', width: 140 },
  { title: '鍛婅绫诲瀷', slotName: 'alert_type', width: 100 },
  { title: '婧愬尯锟?, dataIndex: 'source_region', width: 120 },
  { title: '鐩爣鍖哄煙', dataIndex: 'target_region', width: 120 },
  { title: '瑙﹀彂鏃堕棿', dataIndex: 'created_at', width: 170 }
]

const loadPolicies = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const params = new URLSearchParams({ page: pagination.current, page_size: pagination.pageSize })
    if (policyForm.name) params.append('name', policyForm.name)
    if (policyForm.status) params.append('status', policyForm.status)
    const res = await fetch(`/api/v1/device/cross-border/policies?${params}`, { headers: { 'Authorization': `Bearer ${token}` } }).then(r => r.json())
    if (res.code === 0) { policies.value = res.data?.list || []; pagination.total = res.data?.total || 0 }
    else { policies.value = [] }
  } catch (e) { Message.error('鍔犺浇绛栫暐澶辫触') } finally { loading.value = false }
}

const loadAlerts = async () => {
  alertsLoading.value = true
  try {
    const token = localStorage.getItem('token')
    const params = new URLSearchParams({ page: alertPagination.current, page_size: alertPagination.pageSize })
    if (alertForm.device_id) params.append('device_id', alertForm.device_id)
    if (alertForm.alert_type) params.append('alert_type', alertForm.alert_type)
    const res = await fetch(`/api/v1/device/cross-border/alerts?${params}`, { headers: { 'Authorization': `Bearer ${token}` } }).then(r => r.json())
    if (res.code === 0) { alerts.value = res.data?.list || []; alertPagination.total = res.data?.total || 0 }
    else { alerts.value = [] }
  } catch (e) { Message.error('鍔犺浇鍛婅澶辫触') } finally { alertsLoading.value = false }
}

const openPolicyModal = () => { isEditPolicy.value = false; Object.assign(policyFormData, { id: null, name: '', source_region: '', target_region: '', action: 'warn', is_active: true }); policyModalVisible.value = true }
const editPolicy = (record) => { isEditPolicy.value = true; Object.assign(policyFormData, record); policyModalVisible.value = true }

const handlePolicySubmit = async (done) => {
  if (!policyFormData.name) { Message.warning('璇疯緭鍏ョ瓥鐣ュ悕锟?); done(false); return }
  submitting.value = true
  try {
    const token = localStorage.getItem('token')
    const url = isEditPolicy.value ? `/api/v1/device/cross-border/policies/${policyFormData.id}` : '/api/v1/device/cross-border/policies'
    const res = await fetch(url, { method: isEditPolicy.value ? 'PUT' : 'POST', headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' }, body: JSON.stringify(policyFormData) }).then(r => r.json())
    if (res.code === 0) { Message.success(isEditPolicy.value ? '鏇存柊鎴愬姛' : '鍒涘缓鎴愬姛'); policyModalVisible.value = false; loadPolicies() }
    else { Message.error(res.message || '鎿嶄綔澶辫触') }
    done(true)
  } catch (e) { Message.error('鎿嶄綔澶辫触'); done(false) } finally { submitting.value = false }
}

const deletePolicy = async (record) => {
  try {
    const token = localStorage.getItem('token')
    await fetch(`/api/v1/device/cross-border/policies/${record.id}`, { method: 'DELETE', headers: { 'Authorization': `Bearer ${token}` } })
    Message.success('鍒犻櫎鎴愬姛'); loadPolicies()
  } catch (e) { Message.error('鍒犻櫎澶辫触') }
}

const onPolicyPageChange = (page) => { pagination.current = page; loadPolicies() }
const onAlertPageChange = (page) => { alertPagination.current = page; loadAlerts() }

onMounted(() => { loadPolicies(); loadAlerts() })
</script>

<style scoped>
.page-container { padding: 16px; }
.search-form { margin-bottom: 16px; padding: 16px; background: var(--color-fill-lightest); border-radius: 4px; }
</style>