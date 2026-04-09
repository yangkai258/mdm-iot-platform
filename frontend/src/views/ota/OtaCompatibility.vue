<template>
  <div class="pro-page-container">
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>OTA升级</a-breadcrumb-item>
      <a-breadcrumb-item>固件兼容性矩阵</a-breadcrumb-item>
    </a-breadcrumb>

    <div class="pro-action-bar">
      <a-space>
        <a-button type="primary" @click="batchTest">批量测试</a-button>
        <a-button @click="loadMatrix">刷新</a-button>
      </a-space>
    </div>

    <a-card title="固件版本 × 设备型号兼容性矩阵" style="margin-bottom: 16px">
      <div class="matrix-wrapper">
        <table class="compat-matrix">
          <thead>
            <tr>
              <th class="matrix-header-cell">固件版本 \ 设备型号</th>
              <th v-for="model in deviceModels" :key="model" class="matrix-header-cell">{{ model }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="fw in firmwareVersions" :key="fw.id">
              <td class="matrix-header-cell fw-version">{{ fw.version }} <span class="fw-size">({{ fw.size }})</span></td>
              <td v-for="model in deviceModels" :key="model" class="matrix-cell">
                <div class="compat-status" :class="getCompatClass(fw.version, model)" @click="showDetail(fw, model)">
                  <component :is="getCompatIcon(fw.version, model)" />
                  <span>{{ getCompatText(fw.version, model) }}</span>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </a-card>

    <a-row :gutter="16">
      <a-col :span="8">
        <a-card title="兼容性说明">
          <a-space direction="vertical" fill>
            <a-tag color="green" checkable checked><icon-check-circle /> 已验证兼容</a-tag>
            <a-tag color="orange" checkable checked><icon-exclamation-circle /> 待验证</a-tag>
            <a-tag color="red" checkable checked><icon-close-circle /> 不兼容</a-tag>
          </a-space>
        </a-card>
      </a-col>
      <a-col :span="8">
        <a-card title="统计">
          <a-statistic title="已验证兼容" :value="stats.compatible" :value-style="{ color: '#52c41a' }" />
          <a-statistic title="待验证" :value="stats.pending" :value-style="{ color: '#faad14' }" style="margin-top: 8px" />
          <a-statistic title="不兼容" :value="stats.incompatible" :value-style="{ color: '#ff4d4f' }" style="margin-top: 8px" />
        </a-card>
      </a-col>
      <a-col :span="8">
        <a-card title="最近测试记录">
          <a-timeline>
            <a-timeline-item v-for="log in recentTests" :key="log.id" :color="log.result === 'success' ? 'green' : log.result === 'failed' ? 'red' : 'orange'">
              <p style="margin:0">{{ log.firmware_version }} → {{ log.device_model }}</p>
              <p style="margin:0;color:var(--color-text-3);font-size:12px">{{ formatDate(log.test_time) }}</p>
            </a-timeline-item>
          </a-timeline>
        </a-card>
      </a-col>
    </a-row>

    <!-- 详情弹窗 -->
    <a-modal v-model:visible="detailVisible" :title="`兼容性详情: ${detailInfo?.firmware} × ${detailInfo?.model}`" :width="600" :footer="null">
      <a-descriptions :column="2" bordered style="margin-bottom: 16px">
        <a-descriptions-item label="固件版本">{{ detailInfo?.firmware }}</a-descriptions-item>
        <a-descriptions-item label="设备型号">{{ detailInfo?.model }}</a-descriptions-item>
        <a-descriptions-item label="兼容状态">
          <a-tag :color="getCompatTagColor(detailInfo?.status)">{{ getCompatTagText(detailInfo?.status) }}</a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="测试时间">{{ formatDate(detailInfo?.test_time) }}</a-descriptions-item>
        <a-descriptions-item label="测试耗时" :span="2">{{ detailInfo?.duration_ms }}ms</a-descriptions-item>
        <a-descriptions-item label="测试日志" :span="2">
          <pre class="json-viewer">{{ detailInfo?.test_log || '无' }}</pre>
        </a-descriptions-item>
      </a-descriptions>
      <a-space>
        <a-button type="primary" @click="runTest(detailInfo)">重新测试</a-button>
        <a-button @click="detailVisible = false">关闭</a-button>
      </a-space>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import axios from 'axios'

const loading = ref(false)
const detailVisible = ref(false)
const deviceModels = ref<string[]>(['M5Stack-Core', 'M5Stack-Core2', 'M5Stack-StickC', 'M5Stack-Fire'])
const firmwareVersions = ref<any[]>([])
const compatibilityMap = ref<Record<string, string>>({})
const recentTests = ref<any[]>([])
const stats = reactive({ compatible: 0, pending: 0, incompatible: 0 })

const detailInfo = ref<any>(null)

// Mock compatibility data
const mockCompatData: Record<string, string> = {
  'v2.1.0-M5Stack-Core': 'compatible', 'v2.1.0-M5Stack-Core2': 'compatible', 'v2.1.0-M5Stack-StickC': 'pending', 'v2.1.0-M5Stack-Fire': 'incompatible',
  'v2.0.5-M5Stack-Core': 'compatible', 'v2.0.5-M5Stack-Core2': 'compatible', 'v2.0.5-M5Stack-StickC': 'compatible', 'v2.0.5-M5Stack-Fire': 'compatible',
  'v2.0.3-M5Stack-Core': 'compatible', 'v2.0.3-M5Stack-Core2': 'pending', 'v2.0.3-M5Stack-StickC': 'compatible', 'v2.0.3-M5Stack-Fire': 'compatible',
}

const loadMatrix = async () => {
  loading.value = true
  try {
    const res = await axios.get('/api/v1/ota/compatibility-matrix')
    firmwareVersions.value = res.data.firmware_versions || []
    compatibilityMap.value = res.data.compatibility || {}
    stats.compatible = res.data.stats?.compatible || 0
    stats.pending = res.data.stats?.pending || 0
    stats.incompatible = res.data.stats?.incompatible || 0
  } catch {
    firmwareVersions.value = [
      { id: 1, version: 'v2.1.0', size: '2.5MB' },
      { id: 2, version: 'v2.0.5', size: '2.3MB' },
      { id: 3, version: 'v2.0.3', size: '2.1MB' },
    ]
    compatibilityMap.value = mockCompatData
    stats.compatible = 8; stats.pending = 2; stats.incompatible = 1
  } finally { loading.value = false }

  try {
    const logRes = await axios.get('/api/v1/ota/compatibility-tests?page_size=5')
    recentTests.value = logRes.data.items || []
  } catch {
    recentTests.value = [
      { id: 1, firmware_version: 'v2.1.0', device_model: 'M5Stack-Fire', result: 'failed', test_time: new Date().toISOString() },
      { id: 2, firmware_version: 'v2.1.0', device_model: 'M5Stack-StickC', result: 'pending', test_time: new Date(Date.now() - 600000).toISOString() },
      { id: 3, firmware_version: 'v2.1.0', device_model: 'M5Stack-Core2', result: 'success', test_time: new Date(Date.now() - 1200000).toISOString() },
    ]
  }
}

const batchTest = async () => {
  try {
    await axios.post('/api/v1/ota/compatibility/batch-test')
    Message.success('批量测试已启动')
  } catch { Message.error('启动失败') }
}

const showDetail = (fw: any, model: string) => {
  const key = `${fw.version}-${model}`
  detailInfo.value = {
    firmware: fw.version, model, size: fw.size,
    status: compatibilityMap.value[key] || 'unknown',
    test_time: new Date().toISOString(), duration_ms: 0, test_log: '',
  }
  detailVisible.value = true
}

const runTest = async (info: any) => {
  try {
    await axios.post('/api/v1/ota/compatibility/test', { firmware_version: info.firmware, device_model: info.model })
    Message.success('测试已启动')
    detailVisible.value = false
    loadMatrix()
  } catch { Message.error('测试启动失败') }
}

const getCompatClass = (fw: string, model: string) => compatibilityMap.value[`${fw}-${model}`] || 'pending'
const getCompatText = (fw: string, model: string) => ({ compatible: '✓', pending: '?', incompatible: '✗' }[compatibilityMap.value[`${fw}-${model}`]] || '?')
const getCompatTagColor = (s: string) => ({ compatible: 'green', pending: 'orange', incompatible: 'red' }[s] || 'gray')
const getCompatTagText = (s: string) => ({ compatible: '已验证兼容', pending: '待验证', incompatible: '不兼容' }[s] || '未知')
const getCompatIcon = (fw: string, model: string) => 'span'
const formatDate = (d: string) => d ? new Date(d).toLocaleString('zh-CN') : '-'

onMounted(() => loadMatrix())
</script>

<style scoped lang="less">
.matrix-wrapper { overflow-x: auto; }
.compat-matrix { border-collapse: collapse; width: 100%; min-width: 600px; }
.matrix-header-cell {
  border: 1px solid var(--color-border);
  padding: 10px 14px;
  background: var(--color-fill-1);
  font-weight: 600;
  text-align: center;
  white-space: nowrap;
}
.matrix-cell { border: 1px solid var(--color-border); padding: 4px; text-align: center; }
.fw-version { text-align: left !important; font-weight: 500; .fw-size { font-weight: 400; color: var(--color-text-3); font-size: 12px; } }
.compat-status {
  display: flex; align-items: center; justify-content: center; gap: 4px;
  padding: 6px 10px; border-radius: 4px; cursor: pointer; font-size: 14px; font-weight: 600;
  transition: all 0.2s;
  &.compatible { color: #52c41a; background: rgba(82, 196, 26, 0.1); &:hover { background: rgba(82, 196, 26, 0.2); } }
  &.pending { color: #faad14; background: rgba(250, 173, 20, 0.1); &:hover { background: rgba(250, 173, 20, 0.2); } }
  &.incompatible { color: #ff4d4f; background: rgba(255, 77, 79, 0.1); &:hover { background: rgba(255, 77, 79, 0.2); } }
}
.json-viewer {
  background: var(--color-fill-1); border: 1px solid var(--color-border);
  border-radius: 4px; padding: 8px; font-size: 12px; max-height: 200px; overflow: auto;
  white-space: pre-wrap; word-break: break-all; margin: 0;
}
</style>
