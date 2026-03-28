<template>
  <div class="container">
    <a-card>
      <template #title>
        <a-space><icon-lock /> 权限申请</a-space>
      </template>

      <a-form :model="formData" layout="vertical" style="max-width: 600px">
        <a-form-item label="权限类型" required>
          <a-cascader v-model="formData.permissionType" :options="permissionTree" placeholder="选择权限类型" />
        </a-form-item>

        <a-form-item label="申请理由" required>
          <a-textarea v-model="formData.reason" placeholder="请详细说明申请此权限的原因" :max-length="500" show-word-limit />
        </a-form-item>

        <a-form-item label="有效期">
          <a-select v-model="formData.expiresType">
            <a-option value="temporary">临时权限</a-option>
            <a-option value="permanent">永久权限</a-option>
          </a-select>
        </a-form-item>

        <a-form-item v-if="formData.expiresType === 'temporary'" label="截止日期">
          <a-date-picker v-model="formData.expiresAt" style="width: 100%" />
        </a-form-item>

        <a-form-item label="审批流程预览">
          <a-steps :current="2" size="small">
            <a-step title="提交申请" />
            <a-step title="部门主管审批" />
            <a-step title="管理员审批" />
            <a-step title="完成" />
          </a-steps>
        </a-form-item>

        <a-form-item>
          <a-space>
            <a-button type="primary" @click="handleSubmit">提交申请</a-button>
            <a-button @click="handleCancel">取消</a-button>
          </a-space>
        </a-form-item>
      </a-form>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'

const formData = reactive({
  permissionType: [], reason: '', expiresType: 'temporary', expiresAt: ''
})

const permissionTree = [
  { value: 'device', label: '设备权限', children: [{ value: 'device:read', label: '设备读取' }, { value: 'device:control', label: '设备控制' }] },
  { value: 'data', label: '数据权限', children: [{ value: 'data:read', label: '数据读取' }, { value: 'data:export', label: '数据导出' }] }
]

const handleSubmit = () => { }
const handleCancel = () => { }
</script>

<style scoped>
.container { padding: 16px; }
</style>
