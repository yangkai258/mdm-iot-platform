<template>
    <Breadcrumb :items="['Home','Console','']" />



  <div class="pro-page-container">

    <!-- 面包屑 -->

    <a-breadcrumb class="pro-breadcrumb">

      <a-breadcrumb-item>首页</a-breadcrumb-item>

      <a-breadcrumb-item>设备管理</a-breadcrumb-item>

      <a-breadcrumb-item>配对管理</a-breadcrumb-item>

    </a-breadcrumb>



    <!-- 搜索栏 -->

    <div class="pro-search-bar">

      <a-space>

        <a-input-search v-model="searchKeyword" placeholder="搜索设备ID/用户" style="width: 280px" @search="loadRequests" search-button />

        <a-select v-model="filterStatus" placeholder="请求状态" allow-clear style="width: 130px" @change="loadRequests">

          <a-option value="pending">待审批</a-option>

          <a-option value="approved">已批准</a-option>

          <a-option value="rejected">已拒绝</a-option>

          <a-option value="expired">已过期</a-option>

        </a-select>

      </a-space>

    </div>



    <!-- 操作按钮 -->

    <div class="pro-action-bar">

      <a-space>

        <a-button type="primary" @click="loadRequests">刷新</a-button>

      </a-space>

    </div>



    <!-- 配对请求列表 -->

    <div class="pro-content-area">

      <a-table :columns="columns" :data="requests" :loading="loading" :pagination="pagination" row-key="id" @page-change="handlePageChange">

        <template #status="{ record }">

          <a-tag :color="getStatusColor(record.status)">{{ getStatusText(record.status) }}</a-tag>

        </template>

      </a-table>

        <template #actions="{ record }">

          <a-space>

            <a-button type="primary" size="small" :disabled="record.status !== 'pending'" @click="approveRequest(record)">批准</a-button>

            <a-button type="primary" status="danger" size="small" :disabled="record.status !== 'pending'" @click="rejectRequest(record)">拒绝</a-button>

          </a-space>

        </template>

      </a-table>

    </div>



    <!-- 配对历史 -->

    <a-card class="history-card" title="配对历史记录">

      <a-table :columns="historyColumns" :data="history" :loading="historyLoading" :pagination="historyPagination" row-key="id" @page-change="handleHistoryPageChange">

        <template #action="{ record }">

          <a-tag :color="record.action === 'approved' ? 'green' : 'red'">{{ record.action === 'approved' ? '批准' : '拒绝' }}</a-tag>

        </template>

      </a-table>

      </a-table>

    </a-card>



    <!-- 拒绝原因弹窗 -->

    <a-modal v-model:visible="rejectModalVisible" title="拒绝配对请求" @ok="handleReject" :width="480" :loading="submitting">

      <a-form layout="vertical">

        <a-form-item label="目标设备">

          <a-input :value="selectedRequest?.device_name + ' (' + selectedRequest?.device_id + ')'" disabled />

        </a-form-item>

        <a-form-item label="拒绝原因" required>

          <a-select v-model="rejectReason" placeholder="选择或输入原因">

            <a-option value="unauthorized">未经授权的设备</a-option>

            <a-option value="duplicate">重复配对</a-option>

            <a-option value="policy">违反设备策略</a-option>

            <a-option value="user_cancel">用户取消</a-option>

            <a-option value="other">其他原因</a-option>

          </a-select>

        </a-form-item>

        <a-form-item label="备注">

          <a-textarea v-model="rejectNote" placeholder="补充说明" :rows="2" />

        </a-form-item>

      </a-form>

    </a-modal>

  </div>

</template>



<script setup>

import { ref, reactive, onMounted } from 'vue'

import { Message } from '@arco-design/web-vue'



const loading = ref(false)

const historyLoading = ref(false)

const submitting = ref(false)

const requests = ref([])

const history = ref([])

const searchKeyword = ref('')

const filterStatus = ref('')

const selectedRequest = ref(null)

const rejectModalVisible = ref(false)

const rejectReason = ref('')

const rejectNote = ref('')



const pagination = reactive({ current: 1, pageSize: 10, total: 0 })

const historyPagination = reactive({ current: 1, pageSize: 10, total: 0 })



const columns = [

  { title: '请求ID', dataIndex: 'id', width: 80 },

  { title: '设备ID', dataIndex: 'device_id', ellipsis: true },

  { title: '设备名称', dataIndex: 'device_name' },

  { title: '申请用户', dataIndex: 'user_name', width: 120 },

  { title: '设备型号', dataIndex: 'device_model', width: 120 },

  { title: '请求时间', dataIndex: 'requested_at', width: 170 },

  { title: '状态', slotName: 'status', width: 90 },

  { title: '操作', slotName: 'actions', width: 160, fixed: 'right' }

]



const historyColumns = [

  { title: 'ID', dataIndex: 'id', width: 70 },

  { title: '设备ID', dataIndex: 'device_id', width: 100 },

  { title: '设备名称', dataIndex: 'device_name', ellipsis: true },

  { title: '操作人', dataIndex: 'operator' },

  { title: '操作', slotName: 'action', width: 80 },

  { title: '时间', dataIndex: 'operated_at', width: 170 }

]



const getStatusColor = (s) => ({ pending: 'blue', approved: 'green', rejected: 'red', expired: 'gray' }[s] || 'gray')

const getStatusText = (s) => ({ pending: '待审批', approved: '已批准', rejected: '已拒绝', expired: '已过期' }[s] || s)



const loadRequests = async () => {

  loading.value = true

  try {

    const token = localStorage.getItem('token')

    const params = new URLSearchParams({ page: pagination.current, page_size: pagination.pageSize })

    if (searchKeyword.value) params.append('keyword', searchKeyword.value)

    if (filterStatus.value) params.append('status', filterStatus.value)

    const res = await fetch(`/api/v1/device/pairing/requests?${params}`, { headers: { 'Authorization': `Bearer ${token}` } })

    const data = await res.json()

    if (data.code === 0) {

      requests.value = data.data.list || []

      pagination.total = data.data.total || 0

    }

  } catch (e) {

    requests.value = [

      { id: 1, device_id: 'DEV-001', device_name: '测试设备A', user_name: '张三', device_model: 'M5Stack', requested_at: '2026-03-24 10:00:00', status: 'pending' },

      { id: 2, device_id: 'DEV-002', device_name: '测试设备B', user_name: '李四', device_model: 'M5Stack', requested_at: '2026-03-24 09:30:00', status: 'approved' },

      { id: 3, device_id: 'DEV-003', device_name: '测试设备C', user_name: '王五', device_model: 'M5Stack', requested_at: '2026-03-23 16:00:00', status: 'rejected' }

    ]

    pagination.total = 3

  } finally {

    loading.value = false

  }

}



const loadHistory = async () => {

  historyLoading.value = true

  try {

    const token = localStorage.getItem('token')

    const params = new URLSearchParams({ page: historyPagination.current, page_size: historyPagination.pageSize })

    const res = await fetch(`/api/v1/device/pairing/history?${params}`, { headers: { 'Authorization': `Bearer ${token}` } })

    const data = await res.json()

    if (data.code === 0) {

      history.value = data.data.list || []

      historyPagination.total = data.data.total || 0

    }

  } catch (e) {

    history.value = [

      { id: 1, device_id: 'DEV-002', device_name: '测试设备B', operator: '管理员', action: 'approved', operated_at: '2026-03-24 09:35:00' },

      { id: 2, device_id: 'DEV-003', device_name: '测试设备C', operator: '管理员', action: 'rejected', operated_at: '2026-03-23 16:30:00' }

    ]

    historyPagination.total = 2

  } finally {

    historyLoading.value = false

  }

}



const approveRequest = async (record) => {

  try {

    const token = localStorage.getItem('token')

    const res = await fetch(`/api/v1/device/pairing/requests/${record.id}/approve`, {

      method: 'POST',

      headers: { 'Authorization': `Bearer ${token}` }

    })

    const data = await res.json()

    if (data.code === 0) {

      Message.success('配对请求已批准')

      loadRequests()

      loadHistory()

    } else { Message.error(data.message || '批准失败') }

  } catch (e) { Message.error('批准失败') }

}



const rejectRequest = (record) => {

  selectedRequest.value = record

  rejectReason.value = ''

  rejectNote.value = ''

  rejectModalVisible.value = true

}



const handleReject = async () => {

  if (!rejectReason.value) { Message.warning('请选择拒绝原因'); return }

  submitting.value = true

  try {

    const token = localStorage.getItem('token')

    const res = await fetch(`/api/v1/device/pairing/requests/${selectedRequest.value.id}/reject`, {

      method: 'POST',

      headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' },

      body: JSON.stringify({ reason: rejectReason.value, note: rejectNote.value })

    })

    const data = await res.json()

    if (data.code === 0) {

      Message.success('配对请求已拒绝')

      rejectModalVisible.value = false

      loadRequests()

      loadHistory()

    } else { Message.error(data.message || '拒绝失败') }

  } catch (e) { Message.error('拒绝失败') }

  finally { submitting.value = false }

}



const handlePageChange = (page) => { pagination.current = page; loadRequests() }

const handleHistoryPageChange = (page) => { historyPagination.current = page; loadHistory() }



onMounted(() => { loadRequests(); loadHistory() })

</script>



<style scoped>

.pro-page-container { padding: 20px 24px; min-height: calc(100vh - 64px); background: #f5f7fa; }

.pro-breadcrumb { margin-bottom: 16px; }

.pro-search-bar { margin-bottom: 12px; }

.pro-action-bar { margin-bottom: 16px; }

.pro-content-area { background: #fff; border-radius: 8px; padding: 20px; box-shadow: 0 1px 3px rgba(0,0,0,0.04); margin-bottom: 16px; }

.history-card { border-radius: 8px; }

</style>

