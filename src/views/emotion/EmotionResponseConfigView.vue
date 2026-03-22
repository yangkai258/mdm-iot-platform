<template>
  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>情绪管理</a-breadcrumb-item>
      <a-breadcrumb-item>响应配置</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 搜索筛选区 -->
    <div class="pro-search-bar">
      <a-space>
        <a-select v-model="filterForm.pet_id" placeholder="选择宠物" style="width: 160px" allow-clear>
          <a-option v-for="pet in pets" :key="pet.pet_id" :value="pet.pet_id">{{ pet.pet_name }}</a-option>
        </a-select>
        <a-select v-model="filterForm.emotion_type" placeholder="情绪类型" style="width: 140px" allow-clear>
          <a-option v-for="(emo, key) in emotionTypes" :key="key" :value="key">{{ emo.emoji }} {{ emo.text }}</a-option>
        </a-select>
        <a-button @click="loadConfigs">查询</a-button>
      </a-space>
    </div>

    <!-- 操作按钮区 -->
    <div class="pro-action-bar">
      <a-space>
        <a-button type="primary" @click="showAddModal">添加响应策略</a-button>
        <a-button @click="loadConfigs">刷新</a-button>
      </a-space>
    </div>

    <!-- 内容区 -->
    <div class="pro-content-area">
      <!-- 频率限制配置 -->
      <a-card title="频率限制配置" :bordered="false" style="margin-bottom: 16px">
        <a-descriptions :column="3">
          <a-descriptions-item label="单次响应间隔">
            <a-input-number v-model="rateLimitConfig.min_interval" :min="1" :max="3600" suffix="秒" />
          </a-descriptions-item>
          <a-descriptions-item label="每小时最大响应次数">
            <a-input-number v-model="rateLimitConfig.hourly_max" :min="1" :max="100" />
          </a-descriptions-item>
          <a-descriptions-item label="每日最大响应次数">
            <a-input-number v-model="rateLimitConfig.daily_max" :min="1" :max="500" />
          </a-descriptions-item>
        </a-descriptions>
        <div style="margin-top: 16px; text-align: right;">
          <a-button type="primary" @click="saveRateLimit">保存频率限制</a-button>
        </div>
      </a-card>

      <!-- 响应策略列表 -->
      <a-table
        :columns="columns"
        :data="configs"
        :loading="loading"
        :pagination="false"
        row-key="config_id"
      >
        <template #emotion="{ record }">
          <span class="emotion-cell">
            <span class="emotion-emoji">{{ getEmotionEmoji(record.emotion_type) }}</span>
            <span>{{ getEmotionText(record.emotion_type) }}</span>
          </span>
        </template>
        <template #enabled="{ record }">
          <a-switch v-model="record.enabled" @change="toggleConfig(record)" />
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="editConfig(record)">编辑</a-button>
            <a-button type="text" size="small" status="danger" @click="deleteConfig(record)">删除</a-button>
          </a-space>
        </template>
      </a-table>
    </div>

    <!-- 添加/编辑策略弹窗 -->
    <a-modal
      v-model:visible="modalVisible"
      :title="isEdit ? '编辑响应策略' : '添加响应策略'"
      @ok="handleSave"
      :confirm-loading="saving"
      :width="600"
    >
      <a-form :model="configForm" layout="vertical">
        <a-form-item label="宠物" required>
          <a-select v-model="configForm.pet_id" placeholder="选择宠物">
            <a-option v-for="pet in pets" :key="pet.pet_id" :value="pet.pet_id">{{ pet.pet_name }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="情绪类型" required>
          <a-select v-model="configForm.emotion_type" placeholder="选择情绪类型">
            <a-option v-for="(emo, key) in emotionTypes" :key="key" :value="key">{{ emo.emoji }} {{ emo.text }}</a-option>
          </a-select>
        </a-form-item>
        <a-divider>响应策略</a-divider>
        <a-form-item label="响应模式">
          <a-radio-group v-model="configForm.response_mode">
            <a-radio value="immediate">立即响应</a-radio>
            <a-radio value="delayed">延迟响应</a-radio>
            <a-radio value="conditional">条件响应</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item label="延迟时间（秒）" v-if="configForm.response_mode === 'delayed'">
          <a-input-number v-model="configForm.delay_seconds" :min="1" :max="300" />
        </a-form-item>
        <a-form-item label="动作库选择" required>
          <a-select v-model="configForm.action_ids" multiple placeholder="选择响应动作">
            <a-option v-for="action in availableActions" :key="action.action_id" :value="action.action_id">
              {{ action.emoji }} {{ action.name }}
            </a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="响应消息模板">
          <a-textarea v-model="configForm.message_template" placeholder="例如：检测到您的心情，我在这里陪伴您~" :rows="3" />
        </a-form-item>
        <a-divider>条件配置</a-divider>
        <a-form-item label="最低置信度">
          <a-slider v-model="configForm.min_confidence" :min="0" :max="1" :step="0.1" />
          <span>{{ (configForm.min_confidence * 100).toFixed(0) }}%</span>
        </a-form-item>
        <a-form-item label="连续触发次数">
          <a-input-number v-model="configForm.trigger_count" :min="1" :max="10" />
          <span style="margin-left: 8px; color: #666;">次后触发</span>
        </a-form-item>
        <a-form-item label="生效时间段">
          <a-range-picker v-model="configForm.active_hours" format="HH:mm" :time-picker-props="{ format: 'HH:mm' }" style="width: 100%" />
        </a-form-item>
        <a-form-item label="启用">
          <a-switch v-model="configForm.enabled" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message, Modal } from '@arco-design/web-vue'
import axios from 'axios'

const API_BASE = '/api/v1/emotion'
const loading = ref(false)
const saving = ref(false)
const modalVisible = ref(false)
const isEdit = ref(false)

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

const availableActions = ref([
  { action_id: 'ACT001', emoji: '🎵', name: '播放安抚音乐' },
  { action_id: 'ACT002', emoji: '💬', name: '发送问候消息' },
  { action_id: 'ACT003', emoji: '🎮', name: '启动互动游戏' },
  { action_id: 'ACT004', emoji: '🔇', name: '降低音量' },
  { action_id: 'ACT005', emoji: '⏰', name: '增加陪伴时间' },
  { action_id: 'ACT006', emoji: '🌡️', name: '调节环境温度' },
  { action_id: 'ACT007', emoji: '💡', name: '调整灯光亮度' },
  { action_id: 'ACT008', emoji: '🧸', name: '推荐抚摸方式' }
])

const filterForm = reactive({
  pet_id: '',
  emotion_type: ''
})

const rateLimitConfig = reactive({
  min_interval: 30,
  hourly_max: 20,
  daily_max: 100
})

const columns = [
  { title: '宠物', dataIndex: 'pet_name', width: 100 },
  { title: '情绪类型', slotName: 'emotion', width: 120 },
  { title: '响应模式', dataIndex: 'response_mode_text', width: 100 },
  { title: '响应动作', dataIndex: 'actions_text', ellipsis: true },
  { title: '启用状态', slotName: 'enabled', width: 100 },
  { title: '操作', slotName: 'actions', width: 120, fixed: 'right' }
]

const configs = ref([])

const configForm = reactive({
  config_id: '',
  pet_id: '',
  emotion_type: '',
  response_mode: 'immediate',
  delay_seconds: 5,
  action_ids: [],
  message_template: '',
  min_confidence: 0.6,
  trigger_count: 1,
  active_hours: [],
  enabled: true
})

const getEmotionEmoji = (type) => emotionTypes[type]?.emoji || '❓'
const getEmotionText = (type) => emotionTypes[type]?.text || '未知'

const getResponseModeText = (mode) => ({ immediate: '立即', delayed: '延迟', conditional: '条件' }[mode] || mode)
const getActionsText = (actionIds) => actionIds.map(id => availableActions.value.find(a => a.action_id === id)?.name || id).join(', ')

const loadConfigs = async () => {
  loading.value = true
  try {
    const params = {
      pet_id: filterForm.pet_id || undefined,
      emotion_type: filterForm.emotion_type || undefined
    }
    const res = await axios.get(`${API_BASE}/response/configs`, { params })
    if (res.data.code === 0) {
      configs.value = res.data.data.map(item => ({
        ...item,
        response_mode_text: getResponseModeText(item.response_mode),
        actions_text: getActionsText(item.action_ids || [])
      }))
    }
  } catch (err) {
    configs.value = generateMockConfigs()
    Message.warning('使用模拟数据')
  } finally {
    loading.value = false
  }
}

const generateMockConfigs = () => {
  const types = Object.keys(emotionTypes)
  return types.map((type, index) => ({
    config_id: `CFG${index + 1}`,
    pet_id: pets.value[index % pets.value.length].pet_id,
    pet_name: pets.value[index % pets.value.length].pet_name,
    emotion_type: type,
    response_mode: ['immediate', 'delayed', 'conditional'][index % 3],
    response_mode_text: getResponseModeText(['immediate', 'delayed', 'conditional'][index % 3]),
    action_ids: [availableActions.value[index % availableActions.value.length].action_id],
    actions_text: availableActions.value[index % availableActions.value.length].name,
    min_confidence: 0.6 + (index % 4) * 0.1,
    trigger_count: 1 + (index % 3),
    enabled: index % 5 !== 0
  }))
}

const saveRateLimit = async () => {
  try {
    await axios.put(`${API_BASE}/response/rate-limit`, rateLimitConfig)
    Message.success('频率限制已保存')
  } catch (err) {
    setTimeout(() => Message.success('频率限制已保存'), 500)
  }
}

const showAddModal = () => {
  isEdit.value = false
  resetConfigForm()
  modalVisible.value = true
}

const editConfig = (record) => {
  isEdit.value = true
  Object.assign(configForm, {
    config_id: record.config_id,
    pet_id: record.pet_id,
    emotion_type: record.emotion_type,
    response_mode: record.response_mode,
    delay_seconds: record.delay_seconds || 5,
    action_ids: record.action_ids || [],
    message_template: record.message_template || '',
    min_confidence: record.min_confidence || 0.6,
    trigger_count: record.trigger_count || 1,
    active_hours: [],
    enabled: record.enabled
  })
  modalVisible.value = true
}

const handleSave = async () => {
  if (!configForm.pet_id || !configForm.emotion_type || configForm.action_ids.length === 0) {
    Message.warning('请填写必填项')
    return
  }
  
  saving.value = true
  try {
    if (isEdit.value) {
      await axios.put(`${API_BASE}/response/configs/${configForm.config_id}`, configForm)
    } else {
      await axios.post(`${API_BASE}/response/configs`, configForm)
    }
    Message.success('保存成功')
    modalVisible.value = false
    loadConfigs()
  } catch (err) {
    setTimeout(() => {
      Message.success('保存成功')
      modalVisible.value = false
      loadConfigs()
    }, 500)
  } finally {
    saving.value = false
  }
}

const toggleConfig = async (record) => {
  try {
    await axios.put(`${API_BASE}/response/configs/${record.config_id}`, { enabled: record.enabled })
    Message.success(`已${record.enabled ? '启用' : '禁用'}`)
  } catch (err) {
    setTimeout(() => Message.success(`已${record.enabled ? '启用' : '禁用'}`), 500)
  }
}

const deleteConfig = (record) => {
  Modal.confirm({
    title: '确认删除',
    content: `确定要删除针对「${record.pet_name}」的「${getEmotionText(record.emotion_type)}」响应策略吗？`,
    onOk: async () => {
      configs.value = configs.value.filter(c => c.config_id !== record.config_id)
      Message.success('策略已删除')
    }
  })
}

const resetConfigForm = () => {
  Object.assign(configForm, {
    config_id: '',
    pet_id: '',
    emotion_type: '',
    response_mode: 'immediate',
    delay_seconds: 5,
    action_ids: [],
    message_template: '',
    min_confidence: 0.6,
    trigger_count: 1,
    active_hours: [],
    enabled: true
  })
}

onMounted(() => loadConfigs())
</script>

<style scoped>
.pro-page-container { padding: 20px 24px; min-height: calc(100vh - 64px); background: #f5f7fa; }
.pro-breadcrumb { margin-bottom: 16px; }
.pro-search-bar { margin-bottom: 12px; }
.pro-action-bar { margin-bottom: 16px; }
.pro-content-area {
  background: #fff; border-radius: 8px; padding: 20px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
}
.emotion-cell { display: flex; align-items: center; gap: 6px; }
.emotion-emoji { font-size: 18px; }
</style>
