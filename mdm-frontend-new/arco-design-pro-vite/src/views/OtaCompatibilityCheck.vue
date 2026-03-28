<template>
  <div class="container">
    <a-card>
      <template #title>
        <a-space><icon-check-circle /> OTA固件兼容性检测</a-space>
      </template>

      <a-form :model="formData" layout="inline" style="margin-bottom: 16px">
        <a-form-item label="固件版本">
          <a-select v-model="formData.firmwareVersion" placeholder="选择固件版本" style="width: 200px">
            <a-option value="v2.1.0">v2.1.0</a-option>
            <a-option value="v2.0.0">v2.0.0</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="设备型号">
          <a-select v-model="formData.deviceModel" placeholder="选择型号" style="width: 200px">
            <a-option value="M5Stack-1">M5Stack-1</a-option>
            <a-option value="M5Stack-2">M5Stack-2</a-option>
          </a-select>
        </a-form-item>
        <a-form-item>
          <a-button type="primary" @click="handleCheck">开始检测</a-button>
        </a-form-item>
      </a-form>

      <a-row :gutter="16" style="margin-bottom: 16px">
        <a-col :span="6">
          <a-card>
            <a-statistic title="兼容设备" :value="stats.compatible" suffix="台" />
          </a-card>
        </a-col>
        <a-col :span="6">
          <a-card>
            <a-statistic title="不兼容" :value="stats.incompatible" suffix="台" />
          </a-card>
        </a-col>
        <a-col :span="6">
          <a-card>
            <a-statistic title="未知" :value="stats.unknown" suffix="台" />
          </a-card>
        </a-col>
      </a-row>

      <a-card title="兼容矩阵">
        <a-table :columns="columns" :data="matrix">
          <template #compatible="{ record }">
            <a-tag :color="record.compatible ? 'green' : 'red'">{{ record.compatible ? '兼容' : '不兼容' }}</a-tag>
          </template>
        </a-table>
      </a-card>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'

const formData = reactive({ firmwareVersion: '', deviceModel: '' })
const stats = reactive({ compatible: 85, incompatible: 10, unknown: 5 })

const columns = [
  { title: '设备型号', dataIndex: 'model' },
  { title: '当前版本', dataIndex: 'currentVersion' },
  { title: '目标版本', dataIndex: 'targetVersion' },
  { title: '兼容性', slotName: 'compatible' },
  { title: '风险等级', dataIndex: 'riskLevel' }
]
const matrix = ref([
  { model: 'M5Stack-1', currentVersion: 'v2.0.0', targetVersion: 'v2.1.0', compatible: true, riskLevel: '低' }
])

const handleCheck = () => { }
</script>

<style scoped>
.container { padding: 16px; }
</style>
