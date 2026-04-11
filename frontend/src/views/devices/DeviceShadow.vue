<template>
  <div class="pro-page-container">
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>设备管理</a-breadcrumb-item>
      <a-breadcrumb-item>设备影子</a-breadcrumb-item>
    </a-breadcrumb>
    <div class="pro-search-bar">
      <a-space>
        <a-input-search v-model="searchKeyword" placeholder="搜索设备ID" style="width: 280px" @search="loadShadows" search-button />
        <a-select v-model="filterStatus" placeholder="状态" allow-clear style="width: 120px" @change="loadShadows">
          <a-option value="online">在线</a-option>
          <a-option value="offline">离线</a-option>
        </a-select>
      </a-space>
    </div>
    <div class="pro-action-bar">
      <a-space>
        <a-button type="primary" @click="refreshAll">刷新全部</a-button>
        <a-button @click="loadShadows">刷新</a-button>
      </a-space>
    </div>
    <div class="pro-content-area">
      <a-table :columns="columns" :data="shadows" :loading="loading" :pagination="pagination" row-key="device_id" @page-change="handlePageChange" :scroll="{ x: 1200 }">
        <template #desired="{ record }">
          <a-tag :color="record.desired_updated ? 'green' : 'gray'">
            {{ record.desired_updated ? '已更新' : '未变化' }}
          </a-tag>
        </template>
        <template #reported="{ record }">
          <a-tag :color="record.reported_updated ? 'arcoblue' : 'gray'">
            {{ record.reported_updated ? '已上报' : '未上报' }}
          </a-tag>
        </template>
        <template #version="{ record }">
          <a-input-number v-model="record.version" size="small" style="width: 80px" readonly />
        </template>
        <template #updated_at="{ record }">
          {{ formatDate(record.updated_at) }}
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="viewDetail(record)">详情</a-button>
            <a-button type="text" size="small" @click="updateDesired(record)">更新Desired</a-button>
          </a-space>
        </template>
      </a-table>
    </div>
    <a-modal v-model:visible="detailVisible" title="设备影子详情" :width="800" :footer="null">
      <a-descriptions :column="2" bordered style="margin-bottom: 16px">
        <a-descriptions-item label="设备ID">{{ currentShadow?.device_id }}</a-descriptions-item>
        <a-descriptions-item label="设备名称">{{ currentShadow?.device_name || '-' }}</a-descriptions-item>
        <a-descriptions-item label="版本">{{ currentShadow?.version }}</a-descriptions-item>
        <a-descriptions-item label="最后更新">{{ formatDate(currentShadow?.updated_at) }}</a-descriptions-item>
      </a-descriptions>
      <a-tabs>
        <a-tab-pane key="desired" title="Desired状态">
          <a-alert type="info" style="margin-bottom: 12px">
            <template #title>说明</template>
            Desired状态是由云端下发的期望状态，设备端会同步此状态。
          </a-alert>
          <pre class="json-viewer">{{ JSON.stringify(currentShadow?.desired || {}, null, 2) }}</pre>
        </a-tab-pane>
        <a-tab-pane key="reported" title="Reported状态">
          <a-alert type="info" style="margin-bottom: 12px">
            <template #title>说明</template>
            Reported状态是设备端实际上报的当前状态。
          </a-alert>
          <pre class="json-viewer">{{ JSON.stringify(currentShadow?.reported || {}, null, 2) }}</pre>
        </a-tab-pane>
        <a-tab-pane key="delta" title="差异(Delta)">
          <a-alert type="warning" style="margin-bottom: 12px">
            <template #title>说明</template>
            Delta是desired与reported之间的差异，设备需要同步这些差异。
          </a-alert>
          <pre class="json-viewer">{{ JSON.stringify(currentShadow?.delta || {}, null, 2) }}</pre>
        </a-tab-pane>
      </a-tabs>
    </a-modal>
    <a-modal v-model:visible="updateVisible" title="更新Desired状态" @ok="submitDesired" :width="600" :loading="submitting">
      <a-form :model="desiredForm" layout="vertical">
        <a-form-item label="设备ID">
          <a-input v-model="desiredForm.device_id" readonly />
        </a-form-item>
        <a-form-item label="Desired状态(JSON)">
          <a-textarea v-model="desiredForm.json_str" :rows="10" placeholder='{"temperature": 25, "mode": "auto"}' />
        </a-form-item>
        <a-alert type="warning" message="请确保JSON格式正确，这将直接覆盖设备的Desired状态。" />
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { Message } from '@arco-design/web-vue'

const loading = ref(false)
const detailVisible = ref(false)
const updateVisible = ref(false)
const submitting = ref(false)
const searchKeyword = ref('')
const filterStatus = ref('')
const currentShadow = ref<any>(null)
const shadows = ref<any[]>([])
const pagination = reactive({ current: 1, pageSize: 10, total: 0 })

const columns = [
  { title: '设备ID', dataIndex: 'device_id', fixed: 'left', width: 180, ellipsis: true },
  { title: '设备名称', dataIndex: 'device_name', width: 150, ellipsis: true },
  { title: 'Desired', slotName: 'desired', width: 100 },
  { title: 'Reported', slotName: 'reported', width: 100 },
  { title: '版本', slotName: 'version', width: 100 },
  { title: '更新时间', dataIndex: 'updated_at', slotName: 'updated_at', width: 180 },
  { title: '操作', slotName: 'actions', fixed: 'right', width: 160 },
]

const desiredForm = reactive({ device_id: '', json_str: '' })

const loadShadows = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const params = new URLSearchParams({ page: pagination.current, page_size: pagination.pageSize })
    if (searchKeyword.value) params.append('keyword', searchKeyword.value)
    if (filterStatus.value) params.append('status', filterStatus.value)
    const res = await fetch(`/api/v1/devices/shadows?${params}`, {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    const data = await res.json()
    if (data.code === 0) {
      shadows.value = data.data?.list || []
      pagination.total = data.data?.total || 0
    }
  } catch (e) {
    shadows.value = [
      { device_id: 'dev-001', device_name: '宠物设备1号', desired: { temperature: 25 }, reported: { temperature: 24 }, version: 5, updated_at: new Date().toISOString(), desired_updated: true, reported_updated: true, delta: { temperature: 1 } },
      { device_id: 'dev-002', device_name: '宠物设备2号', desired: { temperature: 26 }, reported: { temperature: 26 }, version: 3, updated_at: new Date(Date.now() - 3600000).toISOString(), desired_updated: false, reported_updated: true, delta: {} },
    ]
    pagination.total = 2
  } finally {
    loading.value = false
  }
}

const refreshAll = async () => {
  try {
    const token = localStorage.getItem('token')
    await fetch('/api/v1/devices/shadows/refresh', {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${token}` }
    })
    Message.success('已触发全部刷新')
    loadShadows()
  } catch (e) {
    Message.error('刷新失败')
  }
}

const viewDetail = (record: any) => {
  currentShadow.value = record
  detailVisible.value = true
}

const updateDesired = (record: any) => {
  currentShadow.value = record
  desiredForm.device_id = record.device_id
  desiredForm.json_str = JSON.stringify(record.desired || {}, null, 2)
  updateVisible.value = true
}

const submitDesired = async () => {
  submitting.value = true
  try {
    JSON.parse(desiredForm.json_str)
    const token = localStorage.getItem('token')
    const res = await fetch(`/api/v1/devices/${desiredForm.device_id}/shadow/desired`, {
      method: 'PUT',
      headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' },
      body: JSON.stringify({ desired: JSON.parse(desiredForm.json_str) })
    })
    const data = await res.json()
    if (data.code === 0) {
      Message.success('更新成功')
      updateVisible.value = false
      loadShadows()
    } else {
      Message.error(data.message || '更新失败')
    }
  } catch (e: any) {
    Message.error('JSON格式错误或更新失败')
  } finally {
    submitting.value = false
  }
}

const handlePageChange = (page: number) => {
  pagination.current = page
  loadShadows()
}

const formatDate = (date: string) => {
  if (!date) return '-'
  return new Date(date).toLocaleString('zh-CN')
}

loadShadows()
</script>

<style scoped>
.json-viewer {
  background: var(--color-fill-1);
  border: 1px solid var(--color-border);
  border-radius: 4px;
  padding: 12px;
  font-family: 'Courier New', monospace;
  font-size: 13px;
  max-height: 400px;
  overflow: auto;
  white-space: pre-wrap;
  word-break: break-all;
  margin: 0;
}
</style>
