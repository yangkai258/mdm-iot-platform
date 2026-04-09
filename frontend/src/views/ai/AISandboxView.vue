<template>
  <div class="pro-page-container">
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>AI 功能</a-breadcrumb-item>
      <a-breadcrumb-item>AI 沙箱测试</a-breadcrumb-item>
    </a-breadcrumb>

    <a-row :gutter="16">
      <a-col :span="12">
        <a-card title="测试输入" class="test-input-card">
          <a-form :model="testForm" layout="vertical">
            <a-form-item label="模型">
              <a-select v-model="testForm.model" placeholder="选择模型">
                <a-option value="behavior">行为识别模型 v2.1.0</a-option>
                <a-option value="emotion">情感分析模型 v2.0.5</a-option>
                <a-option value="voice">语音合成模型 v2.2.0</a-option>
              </a-select>
            </a-form-item>
            <a-form-item label="输入类型">
              <a-radio-group v-model="testForm.inputType">
                <a-radio value="text">文本</a-radio>
                <a-radio value="voice">语音</a-radio>
              </a-radio-group>
            </a-form-item>
            <a-form-item v-if="testForm.inputType === 'text'" label="文本输入">
              <a-textarea v-model="testForm.textInput" :rows="8" placeholder="输入要测试的文本内容..." show-word-limit />
            </a-form-item>
            <a-form-item v-else label="语音输入">
              <a-space direction="vertical" style="width: 100%">
                <a-upload action="/" accept="audio/*" :limit="1" />
                <a-button type="outline" status="warning">开始录音</a-button>
              </a-space>
            </a-form-item>
            <a-form-item label="参数配置">
              <a-collapse :default-active-key="['basic']">
                <a-collapse-item key="basic" title="基础参数">
                  <a-form :model="testForm.params" layout="vertical">
                    <a-form-item label="温度"><a-slider v-model="testForm.params.temperature" :min="0" :max="1" :step="0.1" show-input /></a-form-item>
                    <a-form-item label="最大令牌数"><a-input-number v-model="testForm.params.max_tokens" :min="1" :max="2048" style="width:100%" /></a-form-item>
                    <a-form-item label="Top-P"><a-slider v-model="testForm.params.top_p" :min="0" :max="1" :step="0.05" show-input /></a-form-item>
                  </a-form>
                </a-collapse-item>
              </a-collapse>
            </a-form-item>
            <a-form-item>
              <a-space>
                <a-button type="primary" @click="runTest" :loading="testing">运行测试</a-button>
                <a-button @click="handleReset">重置</a-button>
                <a-button type="outline" @click="showBatchModal = true">批量测试</a-button>
              </a-space>
            </a-form-item>
          </a-form>
        </a-card>
      </a-col>

      <a-col :span="12">
        <a-card title="测试结果" class="test-result-card">
          <template #extra>
            <a-space>
              <a-button type="text" size="small" @click="copyResult" :disabled="!testResult">复制</a-button>
              <a-button type="text" size="small" @click="exportResult" :disabled="!testResult">导出</a-button>
            </a-space>
          </template>
          <div v-if="testing" class="result-loading"><a-spin size="large" /><p>模型推理中，请稍候...</p></div>
          <div v-else-if="!testResult" class="result-empty"><p>点击「运行测试」查看输出结果</p></div>
          <div v-else class="result-content">
            <a-descriptions :column="1" bordered size="small">
              <a-descriptions-item label="模型">{{ testResult.model }}</a-descriptions-item>
              <a-descriptions-item label="版本">{{ testResult.version }}</a-descriptions-item>
              <a-descriptions-item label="推理延迟">{{ testResult.latency_ms }}ms</a-descriptions-item>
              <a-descriptions-item label="置信度">{{ (testResult.confidence * 100).toFixed(1) }}%</a-descriptions-item>
              <a-descriptions-item label="Token使用">{{ testResult.tokens_used }} / {{ testResult.tokens_max }}</a-descriptions-item>
            </a-descriptions>
            <a-divider>输出结果</a-divider>
            <div class="output-display"><pre>{{ testResult.output }}</pre></div>
          </div>
        </a-card>
        <a-card title="历史测试记录" class="history-card" style="margin-top: 16px">
          <a-table :columns="historyColumns" :data="historyList" :pagination="{ pageSize: 5 }" size="small" row-key="id">
            <template #input_type="{ record }"><a-tag>{{ record.input_type === 'text' ? '文本' : '语音' }}</a-tag></template>
            <template #actions="{ record }"><a-button type="text" size="small" @click="loadHistory(record)">复用</a-button></template>
          </a-table>
        </a-card>
      </a-col>
    </a-row>

    <a-modal v-model:visible="showBatchModal" title="批量测试" :width="700">
      <a-form layout="vertical">
        <a-form-item label="模型">
          <a-select v-model="batchForm.model" placeholder="选择模型">
            <a-option value="behavior">行为识别模型 v2.1.0</a-option>
            <a-option value="emotion">情感分析模型 v2.0.5</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="测试用例 (每行一个)">
          <a-textarea v-model="batchForm.cases" :rows="10" placeholder="输入测试用例，每行一个..." />
        </a-form-item>
        <a-form-item label="执行方式">
          <a-radio-group v-model="batchForm.mode">
            <a-radio value="sequential">顺序执行</a-radio>
            <a-radio value="parallel">并行执行</a-radio>
          </a-radio-group>
        </a-form-item>
      </a-form>
      <template #footer>
        <a-button @click="showBatchModal = false">取消</a-button>
        <a-button type="primary" @click="runBatch" :loading="batchRunning">开始批量测试</a-button>
      </template>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { Message } from '@arco-design/web-vue'

const testing = ref(false)
const testResult = ref(null)
const showBatchModal = ref(false)
const batchRunning = ref(false)
const historyList = ref([])

const testForm = reactive({
  model: 'behavior', inputType: 'text', textInput: '',
  params: { temperature: 0.7, max_tokens: 1024, top_p: 0.9 }
})

const batchForm = reactive({ model: 'behavior', cases: '', mode: 'sequential' })

const historyColumns = [
  { title: '时间', dataIndex: 'created_at', width: 160 },
  { title: '模型', dataIndex: 'model_name', width: 120 },
  { title: '输入类型', slotName: 'input_type', width: 80 },
  { title: '延迟', dataIndex: 'latency_ms', width: 80 },
  { title: '操作', width: 80, slotName: 'actions' }
]

const runTest = async () => {
  if (!testForm.textInput && testForm.inputType === 'text') { Message.warning('请输入测试文本'); return }
  testing.value = true; testResult.value = null
  setTimeout(() => {
    testing.value = false
    testResult.value = {
      model: testForm.model === 'behavior' ? '行为识别模型' : '情感分析模型',
      version: '2.1.0',
      latency_ms: Math.floor(80 + Math.random() * 200),
      confidence: 0.7 + Math.random() * 0.29,
      tokens_used: Math.floor(50 + Math.random() * 500),
      tokens_max: 1024,
      output: '识别结果：行为类型 = move, 置信度 = 85.3%, 建议动作 = 移动到位置 A'
    }
    Message.success('测试完成')
  }, 1500)
}

const handleReset = () => { testForm.textInput = ''; testForm.params = { temperature: 0.7, max_tokens: 1024, top_p: 0.9 }; testResult.value = null }
const copyResult = () => Message.success('已复制')
const exportResult = () => Message.info('导出功能开发中')
const loadHistory = (record) => { testForm.model = record.model; testForm.textInput = record.input }
const runBatch = async () => { batchRunning.value = true; setTimeout(() => { batchRunning.value = false; showBatchModal.value = false; Message.success('批量测试完成') }, 2000) }
</script>

<style scoped>
.pro-page-container { padding: 16px; }
.result-loading { text-align: center; padding: 40px; }
.result-empty { text-align: center; padding: 40px; color: var(--color-text-4); }
.output-display { background: #f7f8fa; padding: 12px; border-radius: 4px; }
.output-display pre { margin: 0; white-space: pre-wrap; font-family: 'Courier New', monospace; }
</style>
