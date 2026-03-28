<template>
  <div class="container">
    <a-card>
      <template #title>
        <a-space><icon-check-circle /> 知识审核工作流</a-space>
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

      <a-table :columns="columns" :data="pendingList" :row-selection="{ type: 'checkbox' }">
        <template #knowledgeTitle="{ record }">
          <a-space direction="vertical">
            <span>{{ record.title }}</span>
            <span class="meta">类型: {{ record.type }} | 字数: {{ record.wordCount }}</span>
          </a-space>
        </template>
        <template #submitter="{ record }">
          <a-space direction="vertical">
            <span>{{ record.submitter }}</span>
            <span class="meta">{{ record.submitTime }}</span>
          </a-space>
        </template>
        <template #actions="{ record }">
          <a-button type="primary" size="small" @click="handleApprove(record)">通过</a-button>
          <a-button size="small" @click="handleView(record)">查看</a-button>
          <a-button status="danger" size="small" @click="handleReject(record)">拒绝</a-button>
        </template>
      </a-table>

      <a-card title="审核历史" style="margin-top: 16px">
        <a-table :columns="historyColumns" :data="auditHistory" size="small" :pagination="pagination" />
      </a-card>
    </a-card>

    <a-modal v-model:visible="rejectVisible" title="拒绝原因" @ok="handleRejectSubmit">
      <a-form :model="rejectForm" layout="vertical">
        <a-form-item label="拒绝原因">
          <a-textarea v-model="rejectForm.reason" placeholder="请输入拒绝原因" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'

const rejectVisible = ref(false)
const rejectForm = reactive({ reason: '' })
const selectedRecord = ref(null)
const pagination = reactive({ current: 1, pageSize: 10, total: 30 })

const columns = [
  { title: '知识标题', slotName: 'knowledgeTitle' },
  { title: '提交人', slotName: 'submitter' },
  { title: '提交时间', dataIndex: 'submitTime', width: 180 },
  { title: '操作', slotName: 'actions', width: 200 }
]

const pendingList = ref([
  { id: 1, title: '宠物饮食健康指南', type: '科普', wordCount: 1500, submitter: '张三', submitTime: '2026-03-28 10:00' },
  { id: 2, title: 'AI对话技巧大全', type: '教程', wordCount: 3000, submitter: '李四', submitTime: '2026-03-28 09:30' }
])

const historyColumns = [
  { title: '知识标题', dataIndex: 'title' },
  { title: '审核人', dataIndex: 'auditor' },
  { title: '审核时间', dataIndex: 'auditTime' },
  { title: '结果', dataIndex: 'result' }
]
const auditHistory = ref([
  { title: '宠物行为解读', auditor: 'admin', auditTime: '2026-03-27 15:00', result: '通过' }
])

const handleBatchApprove = () => { }
const handleBatchReject = () => { }
const handleApprove = (r) => { }
const handleReject = (r) => { selectedRecord.value = r; rejectVisible.value = true }
const handleView = (r) => { }
const handleRejectSubmit = () => { rejectVisible.value = false }
</script>

<style scoped>
.container { padding: 16px; }
.meta { font-size: 12px; color: #909399; }
</style>
