<template>
  <div class="page-container">
    <a-card class="general-card" title="跨境设备管控">
      <template #extra>
        <a-button type="primary" @click="openPolicyModal"><icon-plus />新建策略</a-button>
      </template>
      <a-tabs v-model:active-key="activeTab">
        <a-tab-pane key="policies" title="区域策略">
          <div class="search-form">
            <a-form :model="policyForm" layout="inline">
              <a-form-item label="策略名称"><a-input v-model="policyForm.name" placeholder="请输�? /></a-form-item>
              <a-form-item label="状�?>
                <a-select v-model="policyForm.status" placeholder="选择状�? allow-clear style="width: 120px">
                  <a-option value="active">启用</a-option>
                  <a-option value="inactive">停用</a-option>
                </a-select>
              </a-form-item>
              <a-form-item><a-button type="primary" @click="loadPolicies">查询</a-button><a-button @click="Object.keys(policyForm).forEach(k => policyForm[k] = ''); loadPolicies()">重置</a-button></a-form-item>
            </a-form>
          </div>
          <a-table :columns="policyColumns" :data="policies" :loading="loading" :pagination="pagination" @page-change="onPolicyPageChange" row-key="id">
            <template #status="{ record }">
              <a-tag :color="record.status === 'active' ? 'green' : 'gray'">{{ record.status === 'active' ? '启用' : '停用' }}</a-tag>
            </template>
            <template #actions="{ record }">
              <a-button type="text" size="small" @click="editPolicy(record)">编辑</a-button>
              <a-button type="text" size="small" status="danger" @click="deletePolicy(record)">删除</a-button>
            </template>
          </a-table>
        </a-tab-pane>
        <a-tab-pane key="alerts" title="跨境告警">
          <div class="search-form">
            <a-form :model="alertForm" layout="inline">
              <a-form-item label="设备ID"><a-input v-model="alertForm.device_id" placeholder="请输�? /></a-form-item>
              <a-form-item label="告警类型">
                <a-select v-model="alertForm.alert_type" placeholder="选择类型" allow-clear style="width: 140px">
                  <a-option value="cross_in">跨境进入</a-option>
                  <a-option value="cross_out">跨境离开</a-option>
                  <a-option value="region_change">区域变更</a-option>
                </a-select>
              </a-form-item>
              <a-form-item label="时间范围"><a-range-picker v-model="alertForm.time_range" style="width: 240px" /></a-form-item>
              <a-form-item><a-button type="primary" @click="loadAlerts">查询</a-button><a-button @click="Object.assign(alertForm, { device_id: '', alert_type: '', time_range: [] }); loadAlerts()">重置</a-button></a-form-item>
            </a-form>
          </div>
          <a-table :columns="alertColumns" :data="alerts" :loading="alertsLoading" :pagination="alertPagination" @page-change="onAlertPageChange" row-key="id">
            <template #alert_type="{ record }">
              <a-tag :color="record.alert_type === 'cross_in' ? 'green' : 'orange'">{{ record.alert_type === 'cross_in' ? '跨境进入' : record.alert_type === 'cross_out' ? '跨境离开' : '区域变更' }}</a-tag>
            </template>
          </a-table>
        </a-tab-pane>
      </a-tabs>
    </a-card>
    <a-modal v-model:visible="policyModalVisible" :title="isEditPolicy ? '编辑策略' : '新建策略'" @before-ok="handlePolicySubmit" :loading="submitting" :width="560">
      <a-form :model="policyFormData" layout="vertical">
        <a-form-item label="策略名称" required><a-input v-model="policyFormData.name" placeholder="请输入策略名�? /></a-form-item>
        <a-form-item label="源区�?>
          <a-select v-model="policyFormData.source_region" placeholder="选择源区�?>
            <a-option value="cn">中国大陆</a-option>
            <a-option value="us">美国</a-option>
            <a-option value="eu">欧盟</a-option>
            <a-option value="other">其他</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="目标区域">
          <a-select v-model="policyFormData.target_region" placeholder="选择目标区域">
            <a-option value="cn">中国大陆</a-option>
            <a-option value="us">美国</a-option>
            <a-option value="eu">欧盟</a-option>
            <a-option value="other">其他</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="管控动作">
          <a-select v-model="policyFormData.action" placeholder="选择管控动作">
            <a-option value="allow">允许</a-option>
            <a-option value="warn">告警</a-option>
            <a-option value="block">阻止</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="状�?><a-switch v-model="policyFormData.is_active" /></a-form-item>
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
  { title: '策略名称', dataIndex: 'name', width: 180 },
  { title: '源区�?, dataIndex: 'source_region', width: 120 },
  { title: '目标区域', dataIndex: 'target_region', width: 120 },
  { title: '管控动作', dataIndex: 'action', width: 100 },
  { title: '状�?, slotName: 'status', width: 80 },
  { title: '创建时间', dataIndex: 'created_at', width: 170 },
  { title: '操作', slotName: 'actions', width: 120 }
]

const alertColumns = [
  { title: '告警ID', dataIndex: 'id', width: 80 },
  { title: '设备ID', dataIndex: 'device_id', width: 120 },
  { title: '设备名称', dataIndex: 'device_name', width: 140 },
  { title: '告警类型', slotName: 'alert_type', width: 100 },
  { title: '源区�?, dataIndex: 'source_region', width: 120 },
  { title: '目标区域', dataIndex: 'target_region', width: 120 },
  { title: '触发时间', dataIndex: 'created_at', width: 170 }
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
  } catch (e) { Message.error('加载策略失败') } finally { loading.value = false }
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
  } catch (e) { Message.error('加载告警失败') } finally { alertsLoading.value = false }
}

const openPolicyModal = () => { isEditPolicy.value = false; Object.assign(policyFormData, { id: null, name: '', source_region: '', target_region: '', action: 'warn', is_active: true }); policyModalVisible.value = true }
const editPolicy = (record) => { isEditPolicy.value = true; Object.assign(policyFormData, record); policyModalVisible.value = true }

const handlePolicySubmit = async (done) => {
  if (!policyFormData.name) { Message.warning('请输入策略名�?); done(false); return }
  submitting.value = true
  try {
    const token = localStorage.getItem('token')
    const url = isEditPolicy.value ? `/api/v1/device/cross-border/policies/${policyFormData.id}` : '/api/v1/device/cross-border/policies'
    const res = await fetch(url, { method: isEditPolicy.value ? 'PUT' : 'POST', headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' }, body: JSON.stringify(policyFormData) }).then(r => r.json())
    if (res.code === 0) { Message.success(isEditPolicy.value ? '更新成功' : '创建成功'); policyModalVisible.value = false; loadPolicies() }
    else { Message.error(res.message || '操作失败') }
    done(true)
  } catch (e) { Message.error('操作失败'); done(false) } finally { submitting.value = false }
}

const deletePolicy = async (record) => {
  try {
    const token = localStorage.getItem('token')
    await fetch(`/api/v1/device/cross-border/policies/${record.id}`, { method: 'DELETE', headers: { 'Authorization': `Bearer ${token}` } })
    Message.success('删除成功'); loadPolicies()
  } catch (e) { Message.error('删除失败') }
}

const onPolicyPageChange = (page) => { pagination.current = page; loadPolicies() }
const onAlertPageChange = (page) => { alertPagination.current = page; loadAlerts() }

onMounted(() => { loadPolicies(); loadAlerts() })
</script>

<style scoped>
.page-container { padding: 16px; }
.search-form { margin-bottom: 16px; padding: 16px; background: var(--color-fill-lightest); border-radius: 4px; }
</style>