<template>
  <div class="container">
    <a-card>
      <template #title>
        <a-space><icon-launch /> 灰度发布配置</a-space>
      </template>

      <a-steps :current="currentStep">
        <a-step title="选择策略" />
        <a-step title="选择设备" />
        <a-step title="确认配置" />
      </a-steps>

      <div class="step-content">
        <a-card v-if="currentStep === 0" title="选择灰度策略">
          <a-form :model="strategyForm" layout="vertical">
            <a-form-item label="灰度类型">
              <a-radio-group v-model="strategyForm.type">
                <a-radio value="percentage">百分比</a-radio>
                <a-radio value="whitelist">白名单</a-radio>
                <a-radio value="random">随机</a-radio>
              </a-radio-group>
            </a-form-item>
            <a-form-item v-if="strategyForm.type === 'percentage'" label="升级比例">
              <a-slider v-model="strategyForm.percentage" :min="1" :max="100" :show-input="true" />
            </a-form-item>
            <a-form-item v-if="strategyForm.type === 'whitelist'" label="白名单设备">
              <a-select v-model="strategyForm.whitelist" multiple placeholder="选择设备">
                <a-option value="D001">设备-001</a-option>
                <a-option value="D002">设备-002</a-option>
              </a-select>
            </a-form-item>
          </a-form>
          <a-space style="margin-top: 16px">
            <a-button type="primary" @click="currentStep++">下一步</a-button>
          </a-space>
        </a-card>

        <a-card v-if="currentStep === 1" title="选择目标设备">
          <a-space style="margin-bottom: 16px">
            <a-input-search v-model="keyword" placeholder="搜索设备" style="width: 200px" />
            <a-button @click="handleSelectAll">全选</a-button>
            <a-button @click="handleInverse">反选</a-button>
          </a-space>
          <a-table :columns="deviceColumns" :data="devices" :row-selection="{ type: 'checkbox' }" />
          <a-space style="margin-top: 16px">
            <a-button @click="currentStep--">上一步</a-button>
            <a-button type="primary" @click="currentStep++">下一步</a-button>
          </a-space>
        </a-card>

        <a-card v-if="currentStep === 2" title="确认配置">
          <a-descriptions :column="2" bordered>
            <a-descriptions-item label="策略类型">{{ strategyForm.type }}</a-descriptions-item>
            <a-descriptions-item label="目标设备数">{{ selectedCount }} 台</a-descriptions-item>
            <a-descriptions-item label="固件版本">{{ firmwareVersion }}</a-descriptions-item>
            <a-descriptions-item label="升级策略">{{ strategyForm.type === 'percentage' ? strategyForm.percentage + '%' : '白名单' }}</a-descriptions-item>
          </a-descriptions>
          <a-space style="margin-top: 16px">
            <a-button @click="currentStep--">上一步</a-button>
            <a-button type="primary" @click="handleDeploy">开始灰度发布</a-button>
          </a-space>
        </a-card>
      </div>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'

const currentStep = ref(0)
const keyword = ref('')
const strategyForm = reactive({ type: 'percentage', percentage: 10, whitelist: [] })
const firmwareVersion = ref('v2.1.0')
const selectedCount = ref(0)

const deviceColumns = [
  { title: '设备ID', dataIndex: 'deviceId' },
  { title: '设备名称', dataIndex: 'deviceName' },
  { title: '当前版本', dataIndex: 'currentVersion' },
  { title: '状态', dataIndex: 'status' }
]
const devices = ref([
  { deviceId: 'D001', deviceName: '设备-001', currentVersion: 'v2.0.0', status: '在线' }
])

const handleSelectAll = () => { }
const handleInverse = () => { }
const handleDeploy = () => { }
</script>

<style scoped>
.container { padding: 16px; }
.step-content { margin-top: 24px; }
</style>
