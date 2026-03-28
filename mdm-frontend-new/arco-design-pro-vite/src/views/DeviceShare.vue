<template>
  <div class="container">
    <a-card>
      <template #title>
        <a-space><icon-share /> 设备分享</a-space>
      </template>

      <a-tabs default-active-key="shares">
        <a-tab-pane key="shares" title="分享记录">
          <a-space style="margin-bottom: 16px">
            <a-button type="primary" @click="handleAddShare">
              <template #icon><icon-plus /></template>
              添加分享对象
            </a-button>
          </a-space>

          <a-table :columns="columns" :data="tableData">
            <template #permission="{ record }">
              <a-tag :color="getPermissionColor(record.permission)">{{ getPermissionLabel(record.permission) }}</a-tag>
            </template>
            <template #status="{ record }">
              <a-tag :color="record.status === 'active' ? 'green' : 'gray'">{{ record.status === 'active' ? '生效中' : '已失效' }}</a-tag>
            </template>
            <template #expires="{ record }">
              <span :class="{ 'expired': isExpired(record.expiresAt) }">{{ record.expiresAt || '永久' }}</span>
            </template>
            <template #actions="{ record }">
              <a-link @click="handleModify(record)">修改权限</a-link>
              <a-link @click="handleRevoke(record)">撤销</a-link>
            </template>
          </a-table>
        </a-tab-pane>

        <a-tab-pane key="devices" title="我收到的分享">
          <a-table :columns="receivedColumns" :data="receivedData">
            <template #permission="{ record }">
              <a-tag>{{ getPermissionLabel(record.permission) }}</a-tag>
            </template>
            <template #actions="{ record }">
              <a-link @click="handleAccept(record)">接受</a-link>
              <a-link @click="handleReject(record)">拒绝</a-link>
            </template>
          </a-table>
        </a-tab-pane>
      </a-tabs>
    </a-card>

    <a-drawer v-model:visible="drawerVisible" :title="drawerTitle" :width="400" @ok="handleSubmit">
      <a-form :model="formData" layout="vertical">
        <a-form-item label="分享对象" required>
          <a-select v-model="formData.targetUserId" placeholder="选择用户或输入用户ID">
            <a-option value="U001">张三</a-option>
            <a-option value="U002">李四</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="权限级别" required>
          <a-radio-group v-model="formData.permission">
            <a-radio value="read">只读</a-radio>
            <a-radio value="control">控制</a-radio>
            <a-radio value="manage">管理</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item label="有效期">
          <a-select v-model="formData.expiresType">
            <a-option value="permanent">永久</a-option>
            <a-option value="7d">7天</a-option>
            <a-option value="30d">30天</a-option>
            <a-option value="custom">自定义</a-option>
          </a-select>
        </a-form-item>
        <a-form-item v-if="formData.expiresType === 'custom'" label="到期时间">
          <a-date-picker v-model="formData.expiresAt" show-time />
        </a-form-item>
      </a-form>
    </a-drawer>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'

const drawerVisible = ref(false)
const drawerTitle = ref('添加分享')
const formData = reactive({ targetUserId: '', permission: 'read', expiresType: 'permanent', expiresAt: '' })
const columns = [
  { title: '分享对象', dataIndex: 'targetUserName' },
  { title: '用户ID', dataIndex: 'targetUserId' },
  { title: '权限级别', slotName: 'permission' },
  { title: '状态', slotName: 'status' },
  { title: '到期时间', slotName: 'expires' },
  { title: '分享时间', dataIndex: 'sharedAt' },
  { title: '操作', slotName: 'actions', width: 150 }
]
const receivedColumns = [
  { title: '分享人', dataIndex: 'fromUserName' },
  { title: '设备名称', dataIndex: 'deviceName' },
  { title: '权限级别', slotName: 'permission' },
  { title: '操作', slotName: 'actions', width: 150 }
]
const tableData = ref([
  { targetUserName: '张三', targetUserId: 'U001', permission: 'control', status: 'active', expiresAt: '2026-04-01', sharedAt: '2026-03-20' },
  { targetUserName: '李四', targetUserId: 'U002', permission: 'read', status: 'expired', expiresAt: '2026-03-25', sharedAt: '2026-03-15' }
])
const receivedData = ref([
  { fromUserName: '王五', deviceName: 'M5Stack-001', permission: 'control' }
])

const getPermissionColor = (p) => ({ read: 'blue', control: 'orange', manage: 'purple' }[p] || 'gray')
const getPermissionLabel = (p) => ({ read: '只读', control: '控制', manage: '管理' }[p] || p)
const isExpired = (date) => date ? new Date(date) < new Date() : false
const handleAddShare = () => { drawerVisible.value = true; drawerTitle.value = '添加分享' }
const handleModify = (record) => { drawerVisible.value = true; drawerTitle.value = '修改权限' }
const handleRevoke = (record) => { }
const handleAccept = (record) => { }
const handleReject = (record) => { }
const handleSubmit = () => { drawerVisible.value = false }
</script>

<style scoped>
.container { padding: 16px; }
.expired { color: #F56C6C; }
</style>
