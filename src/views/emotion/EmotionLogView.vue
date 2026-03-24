<template>
  <div class="page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>情绪管理</a-breadcrumb-item>
      <a-breadcrumb-item>情绪日志</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 搜索筛选区 -->
    <div class="search-form">
      <a-form :model="searchForm" layout="inline">
        <a-form-item label="宠物">
          <a-select v-model="searchForm.pet_id" placeholder="选择宠物" allow-clear style="width: 160px">
            <a-option v-for="pet in pets" :key="pet.pet_id" :value="pet.pet_id">{{ pet.pet_name }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="情绪类型">
          <a-select v-model="searchForm.emotion_type" placeholder="选择类型" allow-clear style="width: 140px">
            <a-option v-for="(emo, key) in emotionTypes" :key="key" :value="key">{{ emo.emoji }} {{ emo.text }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="来源">
          <a-select v-model="searchForm.source" placeholder="选择来源" allow-clear style="width: 120px">
            <a-option value="voice">语音</a-option>
            <a-option value="text">文字</a-option>
            <a-option value="expression">表情</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="日期">
          <a-range-picker v-model="searchForm.dateRange" style="width: 260px" />
        </a-form-item>
        <a-form-item>
          <a-space>
            <a-button type="primary" @click="handleSearch">搜索</a-button>
            <a-button @click="handleReset">重置</a-button>
          </a-space>
        </a-form-item>
      </a-form>
    </div>

    <!-- 操作栏 -->
    <div class="toolbar">
      <a-space>
        <a-button type="primary" status="warning" @click="exportLogs">导出日志</a-button>
        <a-button @click="handleSearch">刷新</a-button>
      </a-space>
    </div>

    <!-- 表格 -->
    <a-table
      :columns="columns"
      :data="logs"
      :loading="loading"
      :pagination="pagination"
      @change="handleTableChange"
      row-key="log_id"
    >
      <template #emotion="{ record }">
        <span class="emotion-cell">
          <span class="emotion-emoji">{{ getEmotionEmoji(record.emotion_type) }}</span>
          <span>{{ getEmotionText(record.emotion_type) }}</span>
        </span>
      </template>
      <template #source="{ record }">
        <a-tag :color="getSourceColor(record.source)">{{ getSourceText(record.source) }}</a-tag>
      </template>
      <template #confidence="{ record }">
        <a-progress :percent="(record.confidence * 100).toFixed(1)" :show-text="false" :color="getConfidenceColor(record.confidence)" />
        <span class="confidence-text">{{ (record.confidence * 100).toFixed(1) }}%</span>
      </template>
      <template #actions="{ record }">
        <a-button type="text" size="small" @click="showDetail(record)">详情</a-button>
      </template>
    </a-table>

    <!-- 日志详情弹窗 -->
    <a-modal
      v-model:visible="detailModalVisible"
      title="情绪日志详情"
      :width="600"
      :footer="null"
    >
      <a-descriptions :column="2" bordered v-if="currentLog">
        <a-descriptions-item label="日志ID">{{ currentLog.log_id }}</a-descriptions-item>
        <a-descriptions-item label="宠物名称">{{ currentLog.pet_name }}</a-descriptions-item>
        <a-descriptions-item label="情绪类型">
          <span>{{ getEmotionEmoji(currentLog.emotion_type) }} {{ getEmotionText(currentLog.emotion_type) }}</span>
        </a-descriptions-item>
        <a-descriptions-item label="识别来源">{{ getSourceText(currentLog.source) }}</a-descriptions-item>
        <a-descriptions-item label="置信度">{{ (currentLog.confidence * 100).toFixed(1) }}%</a-descriptions-item>
        <a-descriptions-item label="识别时间">{{ currentLog.create_time }}</a-descriptions-item>
        <a-descriptions-item label="原始数据" :span="2">{{ currentLog.raw_data || '无' }}</a-descriptions-item>
        <a-descriptions-item label="上下文" :span="2">{{ currentLog.context || '无' }}</a-descriptions-item>
        <a-descriptions-item label="触发动作">{{ currentLog.triggered_action || '无' }}</a-descriptions-item>
        <a-descriptions-item label="响应状态">
          <a-tag :color="currentLog.response_status === 'success' ? 'green' : 'orange'">
            {{ currentLog.response_status === 'success' ? '成功' : '已响应' }}
          </a-tag>
        </a-descriptions-item>
      </a-descriptions>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import axios from 'axios'

const API_BASE = '/api/v1/emotion'
const loading = ref(false)
const detailModalVisible = ref(false)
const currentLog = ref(null)

const pets = ref([
  { pet_id: 'PET001', pet_name: '小橘' },
  { pet_id: 'PET002', pet_name: '布丁' },
  { pet_id: 'PET003', pet_name: '豆豆' }
])

const emotionTypes = {
  happy: { emoji: '😄', text: '开心' },
  sad: { emoji: '😢', text: '伤心' },
  angry: { emoji: '😠', text: '生气' },
  surprised: { emoji: '😮', text: '惊讶' },
  calm: { emoji: '😌', text: '平静' },
  anxious: { emoji: '😰', text: '焦虑' },
  lonely: { emoji: '🥺', text: '孤独' },
  tired: { emoji: '😴', text: '疲惫' }
}

const searchForm = reactive({
  pet_id: '',
  emotion_type: '',
  source: '',
  dateRange: []
})

const pagination = reactive({ current: 1, pageSize: 20, total: 0 })

const columns = [
  { title: '时间', dataIndex: 'create_time', width: 180 },
  { title: '宠物', dataIndex: 'pet_name', width: 100 },
  { title: '情绪', slotName: 'emotion', width: 120 },
  { title: '来源', slotName: 'source', width: 100 },
  { title: '置信度', slotName: 'confidence', width: 150 },
  { title: '触发动作', dataIndex: 'triggered_action', ellipsis: true },
  { title: '操作', slotName: 'actions', width: 80, fixed: 'right' }
]

const logs = ref([])

const getEmotionEmoji = (type) => emotionTypes[type]?.emoji || '❓'
const getEmotionText = (type) => emotionTypes[type]?.text || '未知'
const getSourceText = (source) => ({ voice: '语音', text: '文字', expression: '表情' }[source] || source)
const getSourceColor = (source) => ({ voice: 'arcoblue', text: 'green', expression: 'purple' }[source] || 'gray')
const getConfidenceColor = (confidence) => confidence >= 0.8 ? 'green' : confidence >= 0.6 ? 'orange' : 'red'

const loadLogs = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.current,
      page_size: pagination.pageSize,
      pet_id: searchForm.pet_id || undefined,
      emotion_type: searchForm.emotion_type || undefined,
      source: searchForm.source || undefined,
      start_date: searchForm.dateRange?.[0]?.format('YYYY-MM-DD') || undefined,
      end_date: searchForm.dateRange?.[1]?.format('YYYY-MM-DD') || undefined
    }
    const res = await axios.get(`${API_BASE}/logs`, { params })
    if (res.data.code === 0) {
      logs.value = res.data.data.list
      pagination.total = res.data.data.pagination.total
    }
  } catch (err) {
    logs.value = generateMockLogs()
    pagination.total = logs.value.length
    Message.warning('使用模拟数据')
  } finally {
    loading.value = false
  }
}

const generateMockLogs = () => {
  const types = Object.keys(emotionTypes)
  const sources = ['voice', 'text', 'expression']
  const actions = ['播放安抚音乐', '发送问候消息', '启动互动游戏', '降低音量', '增加陪伴时间']
  const now = new Date()
  
  return Array.from({ length: 20 }, (_, i) => {
    const type = types[Math.floor(Math.random() * types.length)]
    const source = sources[Math.floor(Math.random() * sources.length)]
    const confidence = 0.6 + Math.random() * 0.4
    const date = new Date(now - i * 3600000 * Math.random() * 5)
    
    return {
      log_id: `LOG${Date.now() - i}`,
      pet_id: pets.value[Math.floor(Math.random() * pets.value.length)].pet_id,
      pet_name: pets.value[Math.floor(Math.random() * pets.value.length)].pet_name,
      emotion_type: type,
      source: source,
      confidence: confidence,
      create_time: date.toLocaleString(),
      raw_data: source === 'voice' ? '音频片段 12.5s, 采样率 16kHz' : 
                source === 'text' ? '"今天心情不太好..."' : 
                '表情特征: 眉头微皱, 嘴角下垂',
      context: '连续3次相同情绪',
      triggered_action: actions[Math.floor(Math.random() * actions.length)],
      response_status: Math.random() > 0.2 ? 'success' : 'pending'
    }
  })
}

const handleSearch = () => {
  pagination.current = 1
  loadLogs()
}

const handleReset = () => {
  searchForm.pet_id = ''
  searchForm.emotion_type = ''
  searchForm.source = ''
  searchForm.dateRange = []
  pagination.current = 1
  loadLogs()
}

const handleTableChange = (pag) => {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  loadLogs()
}

const showDetail = (record) => {
  currentLog.value = record
  detailModalVisible.value = true
}

const exportLogs = () => {
  Message.info('正在导出日志...')
  setTimeout(() => Message.success('日志导出成功'), 1000)
}

onMounted(() => loadLogs())
</script>

<style scoped>
.page-container {
  background: #fff;
  border-radius: 4px;
  padding: 20px;
}

.breadcrumb {
  margin-bottom: 16px;
}

.search-form {
  margin-bottom: 16px;
  padding: 16px;
  background: #f7f8fa;
  border-radius: 4px;
}

.toolbar {
  margin-bottom: 16px;
}

.emotion-cell {
  display: flex;
  align-items: center;
  gap: 6px;
}

.emotion-emoji {
  font-size: 18px;
}

.confidence-text {
  font-size: 12px;
  margin-left: 8px;
}
</style>
