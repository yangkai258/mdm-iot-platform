<template>
  <div class="container">
    <a-card>
      <template #title>
        <a-space><icon-safe /> 应用权限管理</a-space>
      </template>

      <a-tabs default-active-key="permissions">
        <a-tab-pane key="permissions" title="权限列表">
          <a-space style="margin-bottom: 16px">
            <a-button type="primary" @click="handleApply">
              <template #icon><icon-plus /></template>
              申请权限
            </a-button>
            <a-button @click="handleBatchApprove">
              <template #icon><icon-check /></template>
              批量审批
            </a-button>
          </a-space>

          <a-table :columns="columns" :data="tableData" :row-selection="{ type: 'checkbox' }">
            <template #type="{ record }">
              <a-tag :color="record.type === 'sensitive' ? 'red' : 'blue'">{{ record.type === 'sensitive' ? '敏感' : '普通' }}</a-tag>
            </template>
            <template #status="{ record }">
              <a-tag :color="getStatusColor(record.status)">{{ getStatusLabel(record.status) }}</a-tag>
            </template>
            <template #actions="{ record }">
              <a-link @click="handleApprove(record)">审批</a-link>
              <a-link @click="handleView(record)">详情</a-link>
            </template>
          </a-table>
        </a-tab-pane>

        <a-tab-pane key="apply" title="申请记录">
          <a-table :columns="applyColumns" :data="applyData" />
        </a-tab-pane>
      </a-tabs>
    </a-card>

    <a-modal v-model:visible="modalVisible" title="权限申请" @ok="handleSubmit">
      <a-form :model="formData" layout="vertical">
        <a-form-item label="应用" required>
          <a-select v-model="formData.appId" placeholder="选择应用">
            <a-option value="app1">智能宠物App</a-option>
            <a-option value="app2">设备管理器</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="权限" required>
          <a-select v-model="formData.permissions" multiple placeholder="选择权限">
            <a-option value="p1">设备控制权限</a-option>
            <a-option value="p2">位置信息权限</a-option>
            <a-option value="p3">相册访问权限</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="申请理由">
          <a-textarea v-model="formData.reason" placeholder="请输入申请理由" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'

const modalVisible = ref(false)
const formData = reactive({ appId: '', permissions: [], reason: '' })

const columns = [
  { title: '应用名称', dataIndex: 'appName' },
  { title: '用户', dataIndex: 'userName' },
  { title: '权限名称', dataIndex: 'permissionName' },
  { title: '权限类型', slotName: 'type' },
  { title: '申请时间', dataIndex: 'applyTime' },
  { title: '状态', slotName: 'status' },
  { title: '操作', slotName: 'actions', width: 120 }
]

const tableData = ref([
  { id: 1, appName: '智能宠物App', userName: '张三', permissionName: '位置信息', type: 'sensitive', applyTime: '2026-03-28 10:00', status: 'pending' },
  { id: 2, appName: '设备管理器', userName: '李四', permissionName: '设备控制', type: 'normal', applyTime: '2026-03-28 09:30', status: 'approved' }
])

const applyColumns = [
  { title: '应用', dataIndex: 'appName' },
  { title: '权限', dataIndex: 'permissionName' },
  { title: '申请时间', dataIndex: 'applyTime' },
  { title: '状态', dataIndex: 'status' }
]
const applyData = ref([])

const getStatusColor = (s) => ({ pending: 'orange', approved: 'green', rejected: 'red' }[s] || 'gray')
const getStatusLabel = (s) => ({ pending: '待审批', approved: '已通过', rejected: '已拒绝' }[s] || s)

const handleApply = () => { modalVisible.value = true }
const handleBatchApprove = () => { }
const handleApprove = (r) => { }
const handleView = (r) => { }
const handleSubmit = () => { modalVisible.value = false }
</script>

<style scoped>
.container { padding: 16px; }
</style>
