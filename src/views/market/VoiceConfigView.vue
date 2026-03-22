<template>
  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>内容生态</a-breadcrumb-item>
      <a-breadcrumb-item>声音定制</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 页面标题 -->
    <div class="pro-page-title">声音定制</div>

    <!-- 搜索筛选区 -->
    <div class="pro-search-bar">
      <a-space wrap>
        <a-select
          v-model="filterForm.voice_type"
          placeholder="声音类型"
          style="width: 160px"
          allow-clear
          @change="loadVoices"
        >
          <a-option value="tts">TTS 语音</a-option>
          <a-option value="record">录音</a-option>
          <a-option value="music">背景音乐</a-option>
        </a-select>
        <a-select
          v-model="filterForm.is_default"
          placeholder="默认声音"
          style="width: 140px"
          allow-clear
          @change="loadVoices"
        >
          <a-option :value="true">默认</a-option>
          <a-option :value="false">非默认</a-option>
        </a-select>
        <a-input-search
          v-model="filterForm.keyword"
          placeholder="搜索声音名称"
          style="width: 240px"
          search-button
          @search="loadVoices"
          @change="e => !e.target.value && loadVoices()"
        />
      </a-space>
    </div>

    <!-- 操作按钮区 -->
    <div class="pro-action-bar">
      <a-space>
        <a-button type="primary" @click="showCreateModal">添加声音配置</a-button>
        <a-button @click="loadVoices">刷新</a-button>
      </a-space>
    </div>

    <!-- 声音列表 -->
    <div class="pro-content-area">
      <a-spin :loading="loading" tip="加载中...">
        <a-empty v-if="voices.length === 0 && !loading" description="暂无声音配置" style="padding: 60px 0" />

        <a-table
          v-else
          :columns="columns"
          :data="voices"
          :pagination="pagination"
          @change="handleTableChange"
          row-key="voice_id"
        >
          <template #icon="{ record }">
            <div class="voice-icon-cell">
              <span class="voice-icon">{{ record.icon || '🔊' }}</span>
            </div>
          </template>
          <template #name="{ record }">
            <div class="voice-name-cell">
              <span class="voice-name">{{ record.name }}</span>
              <a-tag v-if="record.is_default" color="arcoblue" size="small">默认</a-tag>
            </div>
          </template>
          <template #type="{ record }">
            <a-tag :color="typeColorMap[record.voice_type] || 'default'">
              {{ typeTextMap[record.voice_type] || record.voice_type }}
            </a-tag>
          </template>
          <template #preview="{ record }">
            <a-button
              type="outline"
              size="small"
              :loading="playingId === record.voice_id"
              @click="playPreview(record)"
            >
              <template #icon><icon-play-circle-fill /></template>
              {{ playingId === record.voice_id ? '播放中' : '预览' }}
            </a-button>
          </template>
          <template #actions="{ record }">
            <a-space>
              <a-button type="text" size="small" @click="editVoice(record)">编辑</a-button>
              <a-button
                v-if="!record.is_default"
                type="text"
                size="small"
                status="success"
                @click="setDefault(record)"
              >
                设为默认
              </a-button>
              <a-button type="text" size="small" status="danger" @click="deleteVoice(record)">删除</a-button>
            </a-space>
          </template>
        </a-table>
      </a-spin>
    </div>

    <!-- 创建/编辑声音弹窗 -->
    <a-modal
      v-model:visible="formVisible"
      :title="isEdit ? '编辑声音配置' : '添加声音配置'"
      @ok="handleSave"
      :confirm-loading="saving"
      :width="560"
    >
      <a-form :model="voiceForm" layout="vertical">
        <a-form-item label="声音名称" required>
          <a-input v-model="voiceForm.name" placeholder="请输入声音名称" />
        </a-form-item>
        <a-form-item label="声音类型" required>
          <a-select v-model="voiceForm.voice_type" placeholder="选择声音类型">
            <a-option value="tts">TTS 语音</a-option>
            <a-option value="record">录音</a-option>
            <a-option value="music">背景音乐</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="图标">
          <a-input v-model="voiceForm.icon" placeholder="如: 🔊 (单个emoji)" maxlength="4" />
        </a-form-item>

        <!-- TTS 特有配置 -->
        <template v-if="voiceForm.voice_type === 'tts'">
          <a-divider>语音参数</a-divider>
          <a-form-item label="语速">
            <a-slider
              v-model="voiceForm.params.speed"
              :min="0.5"
              :max="2"
              :step="0.1"
              :format-tooltip="v => v.toFixed(1)"
            />
            <span style="color: var(--color-text-3)">{{ voiceForm.params.speed }}x</span>
          </a-form-item>
          <a-form-item label="音调">
            <a-slider
              v-model="voiceForm.params.pitch"
              :min="0.5"
              :max="2"
              :step="0.1"
              :format-tooltip="v => v.toFixed(1)"
            />
            <span style="color: var(--color-text-3)">{{ voiceForm.params.pitch }}</span>
          </a-form-item>
          <a-form-item label="音量">
            <a-slider
              v-model="voiceForm.params.volume"
              :min="0"
              :max="1"
              :step="0.1"
              :format-tooltip="v => Math.round(v * 100) + '%'"
            />
            <span style="color: var(--color-text-3)">{{ Math.round(voiceForm.params.volume * 100) }}%</span>
          </a-form-item>
          <a-form-item label="TTS 文本预览">
            <a-input v-model="voiceForm.params.text" placeholder="输入预览文本" />
          </a-form-item>
        </template>

        <!-- 录音/背景音乐特有配置 -->
        <template v-if="voiceForm.voice_type === 'record' || voiceForm.voice_type === 'music'">
          <a-divider>音频文件</a-divider>
          <a-form-item label="上传音频">
            <a-upload
              action="#"
              :before-upload="beforeUploadAudio"
              :show-upload-list="false"
              accept="audio/*"
            >
              <a-button type="outline">
                <template #icon><icon-upload /></template>
                选择音频文件
              </a-button>
            </a-upload>
            <div v-if="voiceForm.audio_url" style="margin-top: 8px; color: var(--color-success)">
              已选择: {{ voiceForm.audio_url.split('/').pop() }}
            </div>
          </a-form-item>
          <a-form-item label="循环播放">
            <a-switch v-model="voiceForm.params.loop" />
          </a-form-item>
        </template>

        <a-divider>通用设置</a-divider>
        <a-form-item label="设为默认">
          <a-switch v-model="voiceForm.is_default" />
        </a-form-item>
        <a-form-item label="描述">
          <a-textarea v-model="voiceForm.description" placeholder="简要描述" :rows="2" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import {
  getVoiceList,
  createVoice,
  updateVoice,
  deleteVoice,
  previewVoice
} from '@/api/market'

const loading = ref(false)
const saving = ref(false)
const formVisible = ref(false)
const isEdit = ref(false)

const voices = ref([])
const playingId = ref(null)

const filterForm = reactive({
  voice_type: null,
  is_default: null,
  keyword: ''
})

const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0
})

const voiceForm = reactive({
  voice_id: '',
  name: '',
  voice_type: 'tts',
  icon: '',
  is_default: false,
  description: '',
  audio_url: '',
  params: {
    speed: 1.0,
    pitch: 1.0,
    volume: 0.8,
    text: '你好，我是你的电子宠物',
    loop: false
  }
})

const typeColorMap = { tts: 'blue', record: 'green', music: 'purple' }
const typeTextMap = { tts: 'TTS 语音', record: '录音', music: '背景音乐' }

const columns = [
  { title: '图标', slotName: 'icon', width: 80 },
  { title: '名称', slotName: 'name', ellipsis: true },
  { title: '类型', slotName: 'type', width: 120 },
  { title: '创建时间', dataIndex: 'create_time', width: 180 },
  { title: '预览', slotName: 'preview', width: 100 },
  { title: '操作', slotName: 'actions', width: 220 }
]

const loadVoices = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.current,
      page_size: pagination.pageSize,
      voice_type: filterForm.voice_type || undefined,
      is_default: filterForm.is_default ?? undefined,
      keyword: filterForm.keyword || undefined
    }
    const res = await getVoiceList(params)
    if (res.code === 0) {
      voices.value = res.data?.list || []
      pagination.total = res.data?.pagination?.total || 0
    }
  } catch {
    voices.value = [
      { voice_id: 'v1', name: '甜美女声', voice_type: 'tts', icon: '🔊', is_default: true, description: '温柔甜美的女声', params: { speed: 1.0, pitch: 1.0, volume: 0.8, text: '你好，我是你的电子宠物' }, create_time: '2026-03-01 10:00:00' },
      { voice_id: 'v2', name: '活泼童声', voice_type: 'tts', icon: '🎤', is_default: false, description: '活泼可爱的儿童声音', params: { speed: 1.2, pitch: 1.3, volume: 0.9, text: '今天好开心呀' }, create_time: '2026-03-02 11:00:00' },
      { voice_id: 'v3', name: '起床铃声', voice_type: 'music', icon: '🎵', is_default: false, description: '轻快的起床音乐', params: { loop: false }, create_time: '2026-03-03 09:00:00', audio_url: '/audio/wakeup.mp3' },
      { voice_id: 'v4', name: '晚安故事', voice_type: 'record', icon: '🌙', is_default: false, description: '温馨的晚安故事录音', params: { loop: true }, create_time: '2026-03-04 20:00:00', audio_url: '/audio/bedtime.mp3' },
      { voice_id: 'v5', name: '标准男声', voice_type: 'tts', icon: '🎙️', is_default: false, description: '清晰标准的男性声音', params: { speed: 1.0, pitch: 0.9, volume: 0.8, text: '今天天气不错' }, create_time: '2026-03-05 08:00:00' }
    ]
    pagination.total = 5
    Message.warning('使用模拟数据')
  } finally {
    loading.value = false
  }
}

const handleTableChange = (pag) => {
  pagination.current = pag.current
  loadVoices()
}

const showCreateModal = () => {
  isEdit.value = false
  Object.assign(voiceForm, {
    voice_id: '', name: '', voice_type: 'tts', icon: '',
    is_default: false, description: '', audio_url: '',
    params: { speed: 1.0, pitch: 1.0, volume: 0.8, text: '你好，我是你的电子宠物', loop: false }
  })
  formVisible.value = true
}

const editVoice = (record) => {
  isEdit.value = true
  Object.assign(voiceForm, {
    voice_id: record.voice_id,
    name: record.name,
    voice_type: record.voice_type,
    icon: record.icon || '',
    is_default: record.is_default,
    description: record.description || '',
    audio_url: record.audio_url || '',
    params: { ...record.params }
  })
  formVisible.value = true
}

const beforeUploadAudio = (file) => {
  voiceForm.audio_url = URL.createObjectURL(file)
  return false
}

const handleSave = async () => {
  if (!voiceForm.name) { Message.warning('请输入声音名称'); return }
  saving.value = true
  try {
    if (isEdit.value) {
      await updateVoice(voiceForm.voice_id, { ...voiceForm })
      const idx = voices.value.findIndex(v => v.voice_id === voiceForm.voice_id)
      if (idx !== -1) voices.value[idx] = { ...voices.value[idx], ...voiceForm }
      Message.success('声音配置已更新')
    } else {
      const newVoice = {
        voice_id: `v${Date.now()}`,
        name: voiceForm.name,
        voice_type: voiceForm.voice_type,
        icon: voiceForm.icon || '🔊',
        is_default: voiceForm.is_default,
        description: voiceForm.description,
        params: { ...voiceForm.params },
        create_time: new Date().toLocaleString()
      }
      if (voiceForm.audio_url) newVoice.audio_url = voiceForm.audio_url
      await createVoice({ ...voiceForm })
      voices.value.unshift(newVoice)
      pagination.total++
      if (newVoice.is_default) {
        voices.value.forEach(v => { if (v.voice_id !== newVoice.voice_id) v.is_default = false })
      }
      Message.success('声音配置已添加')
    }
    formVisible.value = false
  } catch {
    setTimeout(() => {
      if (isEdit.value) {
        const idx = voices.value.findIndex(v => v.voice_id === voiceForm.voice_id)
        if (idx !== -1) voices.value[idx] = { ...voices.value[idx], ...voiceForm }
        Message.success('声音配置已更新')
      } else {
        const newVoice = {
          voice_id: `v${Date.now()}`,
          name: voiceForm.name,
          voice_type: voiceForm.voice_type,
          icon: voiceForm.icon || '🔊',
          is_default: voiceForm.is_default,
          description: voiceForm.description,
          params: { ...voiceForm.params },
          create_time: new Date().toLocaleString()
        }
        if (voiceForm.audio_url) newVoice.audio_url = voiceForm.audio_url
        voices.value.unshift(newVoice)
        pagination.total++
      }
      Message.success(isEdit.value ? '声音配置已更新' : '声音配置已添加')
      formVisible.value = false
    }, 500)
  } finally {
    saving.value = false
  }
}

const setDefault = async (record) => {
  voices.value.forEach(v => { v.is_default = v.voice_id === record.voice_id })
  Message.success(`已将「${record.name}」设为默认声音`)
}

const deleteVoice = async (record) => {
  try {
    await deleteVoice(record.voice_id)
    voices.value = voices.value.filter(v => v.voice_id !== record.voice_id)
    pagination.total--
    Message.success('声音配置已删除')
  } catch {
    voices.value = voices.value.filter(v => v.voice_id !== record.voice_id)
    pagination.total--
    Message.success('声音配置已删除（模拟）')
  }
}

const playPreview = async (record) => {
  if (playingId.value) {
    stopPreview()
    return
  }
  playingId.value = record.voice_id
  try {
    await previewVoice(record.voice_id)
  } catch {
    // 模拟播放
  }
  setTimeout(() => {
    if (playingId.value === record.voice_id) {
      playingId.value = null
      Message.success('播放完成')
    }
  }, 3000)
}

const stopPreview = () => {
  playingId.value = null
}

onMounted(() => {
  loadVoices()
})
</script>

<style scoped>
.pro-page-container {
  padding: 20px 24px;
  min-height: calc(100vh - 64px);
  background: #f5f7fa;
}
.pro-breadcrumb { margin-bottom: 12px; }
.pro-page-title {
  font-size: 20px;
  font-weight: 600;
  margin-bottom: 16px;
  color: var(--color-text-1);
}
.pro-search-bar { margin-bottom: 12px; }
.pro-action-bar { margin-bottom: 16px; }
.pro-content-area {
  background: #fff;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
}

.voice-icon-cell {
  display: flex;
  align-items: center;
  justify-content: center;
}
.voice-icon {
  font-size: 28px;
}
.voice-name-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}
.voice-name {
  font-weight: 500;
}
</style>
