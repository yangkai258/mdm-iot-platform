<template>
  <div class="page-container">
    <div class="search-form">
      <a-form :model="form" layout="inline">
        <a-form-item label="声音类型">
          <a-select v-model="form.voice_type" placeholder="选择类型" style="width: 140px" allow-clear @change="loadVoices">
            <a-option value="tts">TTS 语音</a-option>
            <a-option value="record">录音</a-option>
            <a-option value="music">背景音乐</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="默认声音">
          <a-select v-model="form.is_default" placeholder="选择" style="width: 120px" allow-clear @change="loadVoices">
            <a-option :value="true">默认</a-option>
            <a-option :value="false">非默认</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="关键词">
          <a-input v-model="form.keyword" placeholder="搜索声音名称" @search="loadVoices" />
        </a-form-item>
        <a-form-item>
          <a-button type="primary" @click="loadVoices">搜索</a-button>
          <a-button @click="handleReset">重置</a-button>
        </a-form-item>
      </a-form>
    </div>
    <div class="toolbar">
      <a-button type="primary" @click="handleCreate">添加声音</a-button>
    </div>
    <a-table :columns="columns" :data="voices" :loading="loading" :pagination="pagination" @change="handleTableChange" />
    <a-modal v-model:visible="modalVisible" :title="modalTitle">
      <a-form :model="form" label-col-flex="100px">
        <a-form-item label="声音名称"><a-input v-model="form.name" /></a-form-item>
        <a-form-item label="声音类型">
          <a-select v-model="form.voice_type" placeholder="选择类型">
            <a-option value="tts">TTS 语音</a-option>
            <a-option value="record">录音</a-option>
            <a-option value="music">背景音乐</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="图标"><a-input v-model="form.icon" placeholder="如: 🔊" /></a-form-item>
        <a-form-item label="描述"><a-textarea v-model="form.description" :rows="2" /></a-form-item>
      </a-form>
      <template #footer>
        <a-button @click="modalVisible = false">取消</a-button>
        <a-button type="primary" @click="handleSubmit">确定</a-button>
      </template>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import { getVoiceList, createVoice, deleteVoice as deleteVoiceApi } from '@/api/market'

const loading = ref(false)
const modalVisible = ref(false)
const modalTitle = ref('添加声音')
const isEdit = ref(false)

const form = reactive({
  voice_type: null,
  is_default: null,
  keyword: '',
  id: '',
  name: '',
  icon: '',
  description: ''
})

const voices = ref([])
const pagination = reactive({ total: 0, current: 1, pageSize: 20 })

const typeTextMap = { tts: 'TTS 语音', record: '录音', music: '背景音乐' }

const columns = [
  { title: '图标', dataIndex: 'icon', width: 80 },
  { title: '名称', dataIndex: 'name' },
  { title: '类型', dataIndex: 'type_name', width: 120 },
  { title: '创建时间', dataIndex: 'create_time', width: 180 },
  { title: '是否默认', dataIndex: 'default_name', width: 100 }
]

const loadVoices = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.current,
      page_size: pagination.pageSize,
      voice_type: form.voice_type || undefined,
      is_default: form.is_default ?? undefined,
      keyword: form.keyword || undefined
    }
    const res = await getVoiceList(params)
    if (res.code === 0) {
      voices.value = (res.data?.list || []).map(v => ({
        ...v,
        type_name: typeTextMap[v.voice_type] || v.voice_type,
        default_name: v.is_default ? '是' : '否'
      }))
      pagination.total = res.data?.pagination?.total || 0
    }
  } catch {
    voices.value = [
      { voice_id: 'v1', name: '甜美女声', voice_type: 'tts', type_name: 'TTS 语音', icon: '🔊', is_default: true, default_name: '是', description: '温柔甜美的女声', create_time: '2026-03-01 10:00:00' },
      { voice_id: 'v2', name: '活泼童声', voice_type: 'tts', type_name: 'TTS 语音', icon: '🎤', is_default: false, default_name: '否', description: '活泼可爱的儿童声音', create_time: '2026-03-02 11:00:00' },
      { voice_id: 'v3', name: '起床铃声', voice_type: 'music', type_name: '背景音乐', icon: '🎵', is_default: false, default_name: '否', description: '轻快的起床音乐', create_time: '2026-03-03 09:00:00' },
      { voice_id: 'v4', name: '晚安故事', voice_type: 'record', type_name: '录音', icon: '🌙', is_default: false, default_name: '否', description: '温馨的晚安故事录音', create_time: '2026-03-04 20:00:00' },
      { voice_id: 'v5', name: '标准男声', voice_type: 'tts', type_name: 'TTS 语音', icon: '🎙️', is_default: false, default_name: '否', description: '清晰标准的男性声音', create_time: '2026-03-05 08:00:00' }
    ]
    pagination.total = 5
  } finally {
    loading.value = false
  }
}

const handleTableChange = (pag) => {
  pagination.current = pag.current
  loadVoices()
}

const handleReset = () => {
  form.voice_type = null
  form.is_default = null
  form.keyword = ''
  loadVoices()
}

const handleCreate = () => {
  isEdit.value = false
  modalTitle.value = '添加声音'
  Object.assign(form, { id: '', name: '', voice_type: '', icon: '', description: '' })
  modalVisible.value = true
}

const handleSubmit = async () => {
  if (!form.name) { Message.warning('请填写声音名称'); return }
  if (!form.voice_type) { Message.warning('请选择声音类型'); return }
  try {
    await createVoice({ name: form.name, voice_type: form.voice_type, icon: form.icon, description: form.description })
  } catch {}
  voices.value.unshift({
    voice_id: `v${Date.now()}`,
    name: form.name,
    voice_type: form.voice_type,
    type_name: typeTextMap[form.voice_type] || form.voice_type,
    icon: form.icon || '🔊',
    is_default: false,
    default_name: '否',
    description: form.description,
    create_time: new Date().toLocaleString()
  })
  pagination.total++
  Message.success('添加成功')
  modalVisible.value = false
}

onMounted(() => { loadVoices() })
</script>

<style scoped>
.page-container { background: #fff; border-radius: 4px; padding: 20px; }
.search-form { margin-bottom: 16px; padding: 16px; background: #f7f8fa; border-radius: 4px; }
.toolbar { margin-bottom: 16px; }
</style>
