<template>
  <div class="container">
    <a-card>
      <template #title>
        <a-space><icon-delete /> 远程擦除</a-space>
      </template>

      <a-result status="error" title="危险操作">
        <template #sub-title>
          此操作将永久删除设备上的所有数据，且不可恢复！
        </template>
      </a-result>

      <a-form :model="formData" layout="vertical" style="max-width: 600px; margin-top: 24px">
        <a-form-item label="设备ID" required>
          <a-select v-model="formData.deviceId" placeholder="选择要擦除的设备">
            <a-option value="D001">设备-001</a-option>
          </a-select>
        </a-form-item>

        <a-form-item label="擦除确认码" required>
          <a-input v-model="formData.confirmCode" placeholder="输入 CONFIRM-WIPE 以确认" />
        </a-form-item>

        <a-form-item label="擦除原因">
          <a-textarea v-model="formData.reason" placeholder="请输入擦除原因" />
        </a-form-item>

        <a-form-item>
          <a-button type="primary" status="danger" @click="handleWipe" :disabled="formData.confirmCode !== 'CONFIRM-WIPE'">
            确认擦除
          </a-button>
          <a-button @click="handleCancel">取消</a-button>
        </a-form-item>
      </a-form>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'

const formData = reactive({ deviceId: '', confirmCode: '', reason: '' })

const handleWipe = () => { }
const handleCancel = () => { }
</script>

<style scoped>
.container { padding: 16px; }
</style>
