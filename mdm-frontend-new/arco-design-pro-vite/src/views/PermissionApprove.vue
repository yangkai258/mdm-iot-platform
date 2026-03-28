<template>
  <div class="container">
    <a-card>
      <template #title>
        <a-space><icon-check-circle /> 权限审批</a-space>
      </template>

      <a-space style="margin-bottom: 16px">
        <a-button type="primary" @click="handleBatchApprove">
          <template #icon><icon-check /></template>
          批量通过
        </a-button>
        <a-button status="danger" @click="handleBatchReject">
          <template #icon><icon-close /></template>
          批量拒绝
        </a-button>
      </a-space>

      <a-table :columns="columns" :data="requests" :row-selection="{ type: 'checkbox' }">
        <template #permissionType="{ record }">
          <a-tag>{{ record.permissionName }}</a-tag>
        </template>
        <template #expires="{ record }">
          <span v-if="record.expiresAt">{{ record.expiresAt }}</span>
          <a-tag v-else color="green">永久</a-tag>
        </span>
        </template>
        <template #actions="{ record }">
          <a-button type="primary" size="small" @click="handleApprove(record)">通过</a-button>
          <a-button size="small" @click="handleViewDetail(record)">详情</a-button>
          <a-button status="danger" size="small" @click="handleReject(record)">拒绝</a-button>
        </template>
      </a-table>
    </a-card>

    <a-drawer v-model:visible="detailVisible" title="申请详情" :width="500">
      <a-descriptions :column="1" bordered v-if="selectedRequest">
        <a-descriptions-item label="申请人">{{ selectedRequest.requester }}</a-descriptions-item>
        <a-descriptions-item label="权限">{{ selectedRequest.permissionName }}</a-descriptions-item>
        <a-descriptions-item label="申请时间">{{ selectedRequest.applyTime }}</a-descriptions-item>
        <a-descriptions-item label="有效期">{{ selectedRequest.expiresAt || '永久' }}</a-descriptions-item>
        <a-descriptions-item label="申请理由">{{ selectedRequest.reason }}</a-descriptions-item>
      </a-descriptions>
      <a-divider>审批</a-divider>
      <a-form :model="approveForm" layout="vertical">
        <a-form-item label="审批意见">
          <a-textarea v-model="approveForm.comment" placeholder="请输入审批意见" />
        </a-form-item>
        <a-space>
          <a-button type="primary" @click="handleApproveConfirm">通过</a-button>
          <a-button status="danger" @click="handleRejectConfirm">拒绝</a-button>
        </a-space>
      </a-form>
    </a-drawer>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'

const detailVisible = ref(false)
const selectedRequest = ref(null)
const approveForm = reactive({ comment: '' })

const columns = [
  { title: '申请人', dataIndex: 'requester' },
  { title: '部门', dataIndex: 'department' },
  { title: '权限', slotName: 'permissionType' },
  { title: '申请时间', dataIndex: 'applyTime' },
  { title: '有效期', slotName: 'expires' },
  { title: '状态', dataIndex: 'status' },
  { title: '操作', slotName: 'actions', width: 200 }
]

const requests = ref([
  { id: 1, requester: '张三', department: '研发部', permissionName: '设备控制', applyTime: '2026-03-28 10:00', expiresAt: '2026-04-28', status: '待审批' },
  { id: 2, requester: '李四', department: '运维部', permissionName: '数据导出', applyTime: '2026-03-28 09:00', expiresAt: '', status: '待审批' }
])

const handleBatchApprove = () => { }
const handleBatchReject = () => { }
const handleApprove = (r) => { selectedRequest.value = r; detailVisible.value = true }
const handleViewDetail = (r) => { selectedRequest.value = r; detailVisible.value = true }
const handleReject = (r) => { selectedRequest.value = r; detailVisible.value = true }
const handleApproveConfirm = () => { detailVisible.value = false }
const handleRejectConfirm = () => { detailVisible.value = false }
</script>

<style scoped>
.container { padding: 16px; }
</style>
