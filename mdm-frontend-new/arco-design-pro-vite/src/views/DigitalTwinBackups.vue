<template>
  <div class="container">
    <a-card>
      <template #title>
        <a-space><icon-cloud /> 数字孪生备份管理</a-space>
      </template>
      <template #extra>
        <a-button type="primary" @click="handleBackup">
          <template #icon><icon-plus /></template>
          立即备份
        </a-button>
      </template>

      <a-row :gutter="16" style="margin-bottom: 16px">
        <a-col :span="4">
          <a-card size="small">
            <a-statistic title="备份次数" :value="stats.totalBackups" />
          </a-card>
        </a-col>
        <a-col :span="4">
          <a-card size="small">
            <a-statistic title="总大小" :value="stats.totalSize" suffix="MB" />
          </a-card>
        </a-col>
        <a-col :span="4">
          <a-card size="small">
            <a-statistic title="最近备份" :value="stats.lastBackup" />
          </a-card>
        </a-col>
      </a-row>

      <a-card title="备份列表">
        <template #extra>
          <a-button @click="handleAutoBackup">
            <template #icon><icon-settings /></template>
            设置自动备份
          </a-button>
        </template>

        <a-table :columns="columns" :data="backups">
          <template #type="{ record }">
            <a-tag :color="record.type === 'full' ? 'blue' : 'green'">{{ record.type === 'full' ? '全量' : '增量' }}</a-tag>
          </template>
          <template #status="{ record }">
            <a-badge :status="record.status === 'ready' ? 'success' : 'processing'" :text="record.status === 'ready' ? '就绪' : '创建中'" />
          </template>
          <template #actions="{ record }">
            <a-link @click="handleRestore(record)">恢复</a-link>
            <a-link @click="handleView(record)">详情</a-link>
            <a-link @click="handleDelete(record)">删除</a-link>
          </template>
        </a-table>
      </a-card>
    </a-card>

    <a-modal v-model:visible="restoreVisible" title="恢复备份" @ok="handleRestoreConfirm">
      <a-result status="warning" title="恢复后将覆盖当前数据">
        <template #sub-title>
          备份时间: {{ selectedBackup?.createdAt }}
        </template>
      </a-result>
    </a-modal>

    <a-modal v-model:visible="autoBackupVisible" title="自动备份设置" @ok="handleSaveAutoBackup">
      <a-form :model="autoBackupForm" layout="vertical">
        <a-form-item label="自动备份">
          <a-switch v-model="autoBackupForm.enabled" />
        </a-form-item>
        <a-form-item label="备份频率">
          <a-select v-model="autoBackupForm.frequency">
            <a-option value="daily">每天</a-option>
            <a-option value="weekly">每周</a-option>
            <a-option value="monthly">每月</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="保留数量">
          <a-input-number v-model="autoBackupForm.keepCount" :min="1" :max="10" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'

const stats = reactive({ totalBackups: 12, totalSize: 256, lastBackup: '2小时前' })
const restoreVisible = ref(false)
const autoBackupVisible = ref(false)
const selectedBackup = ref(null)
const autoBackupForm = reactive({ enabled: true, frequency: 'daily', keepCount: 7 })

const columns = [
  { title: '备份ID', dataIndex: 'id', width: 100 },
  { title: '类型', slotName: 'type', width: 80 },
  { title: '大小', dataIndex: 'size', width: 100 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '创建时间', dataIndex: 'createdAt' },
  { title: '操作', slotName: 'actions', width: 180 }
]

const backups = ref([
  { id: 'B001', type: 'full', size: '120MB', status: 'ready', createdAt: '2026-03-28 08:00' },
  { id: 'B002', type: 'incremental', size: '25MB', status: 'ready', createdAt: '2026-03-27 08:00' }
])

const handleBackup = () => { }
const handleRestore = (r) => { selectedBackup.value = r; restoreVisible.value = true }
const handleView = (r) => { }
const handleDelete = (r) => { }
const handleRestoreConfirm = () => { restoreVisible.value = false }
const handleAutoBackup = () => { autoBackupVisible.value = true }
const handleSaveAutoBackup = () => { autoBackupVisible.value = false }
</script>

<style scoped>
.container { padding: 16px; }
</style>
