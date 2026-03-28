<template>
  <div class="container">
    <a-card>
      <template #title>
        <a-space><icon-lock /> 远程锁定</a-space>
      </template>

      <a-form :model="formData" layout="vertical" style="max-width: 600px">
        <a-form-item label="设备ID" required>
          <a-select v-model="formData.deviceId" placeholder="选择要锁定的设备">
            <a-option value="D001">设备-001</a-option>
            <a-option value="D002">设备-002</a-option>
          </a-select>
        </a-form-item>

        <a-form-item label="锁定类型">
          <a-radio-group v-model="formData.lockType">
            <a-radio value="temporary">临时锁定</a-radio>
            <a-radio value="permanent">永久锁定</a-radio>
          </a-radio-group>
        </a-form-item>

        <a-form-item v-if="formData.lockType === 'temporary'" label="锁定时长">
          <a-input-number v-model="formData.duration" :min="1" :max="1440" /> 分钟
        </a-form-item>

        <a-form-item label="锁定原因">
          <a-textarea v-model="formData.reason" placeholder="请输入锁定原因" />
        </a-form-item>

        <a-form-item>
          <a-result status="warning" title="警告：此操作将导致设备无法使用">
            <template #sub-title>
              确认要远程锁定此设备吗？
            </template>
          </a-result>
        </a-form-item>

        <a-form-item>
          <a-button type="primary" status="danger" @click="handleLock">确认锁定</a-button>
          <a-button @click="handleCancel">取消</a-button>
        </a-form-item>
      </a-form>

      <a-divider>锁定记录</a-divider>

      <a-table :columns="columns" :data="records" />
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'

const formData = reactive({ deviceId: '', lockType: 'temporary', duration: 60, reason: '' })

const columns = [
  { title: '时间', dataIndex: 'time' },
  { title: '设备ID', dataIndex: 'deviceId' },
  { title: '锁定类型', dataIndex: 'lockType' },
  { title: '状态', dataIndex: 'status' },
  { title: '操作人', dataIndex: 'operator' }
]
const records = ref([])

const handleLock = () => { }
const handleCancel = () => { }
</script>

<style scoped>
.container { padding: 16px; }
</style>
