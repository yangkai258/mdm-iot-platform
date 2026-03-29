<template>
  <div class="page-container">
    <a-card class="general-card" title="情绪识别配置">
      <template #extra>
        <a-button type="primary" @click="handleSave"><icon-save />保存配置</a-button>
      </template>
      <div class="search-form">
        <a-form :model="form" layout="inline">
          <a-form-item label="识别模式">
            <a-select v-model="form.mode" placeholder="请选择" style="width: 140px">
              <a-option value="audio">语音情绪识别</a-option>
              <a-option value="visual">视觉情绪识别</a-option>
              <a-option value="both">综合识别</a-option>
            </a-select>
          </a-form-item>
          <a-form-item>
            <a-button type="primary" @click="handleSearch">查询</a-button>
            <a-button @click="handleReset">重置</a-button>
          </a-form-item>
        </a-form>
      </div>
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" />
    <a-modal v-model:visible="modalVisible" title="编辑配置" :width="480">
      <a-form :model="form" label-col-flex="100px">
        <a-form-item label="识别模式">
          <a-select v-model="form.mode">
            <a-option value="audio">语音情绪识别</a-option>
            <a-option value="visual">视觉情绪识别</a-option>
            <a-option value="both">综合识别</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="识别灵敏度">
          <a-input-number v-model="form.sensitivity" :min="1" :max="10" style="width: 100%" />
        </a-form-item>
        <a-form-item label="置信度阈值">
          <a-input-number v-model="form.threshold" :min="0" :max="100" suffix="%" style="width: 100%" />
        </a-form-item>
      </a-form>
      <template #footer>
        <a-button @click="modalVisible = false">取消</a-button>
        <a-button type="primary" @click="handleSubmit">确定</a-button>
      </template>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import { IconSave } from '@arco-design/web-vue/es/icon'
import { getEmotionRecognizeConfig, updateEmotionRecognizeConfig } from '@/api/emotion'

const loading = ref(false)
const data = ref<any[]>([])
const modalVisible = ref(false)

const form = reactive({
  mode: 'both',
  sensitivity: 5,
  threshold: 70
})

const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0
})

const columns = [
  { title: '配置项', dataIndex: 'name', width: 200 },
  { title: '当前值', dataIndex: 'value', width: 200 },
  { title: '描述', dataIndex: 'description', ellipsis: true },
  { title: '操作', slotName: 'actions', width: 120 }
]

async function loadData() {
  try {
    loading.value = true
    const res = await getEmotionRecognizeConfig()
    data.value = res.data || []
    pagination.total = data.value.length
    if (res.data) {
      form.mode = res.data.mode || 'both'
      form.sensitivity = res.data.sensitivity || 5
      form.threshold = res.data.threshold || 70
    }
  } catch (err: any) {
    Message.error('加载失败: ' + err.message)
  } finally {
    loading.value = false
  }
}

function handleSearch() {
  loadData()
}

function handleReset() {
  form.mode = 'both'
  form.sensitivity = 5
  form.threshold = 70
  loadData()
}

function handleSave() {
  modalVisible.value = true
}

async function handleSubmit() {
  try {
    await updateEmotionRecognizeConfig(form)
    Message.success('保存成功')
    modalVisible.value = false
    loadData()
  } catch (err: any) {
    Message.error('保存失败: ' + err.message)
  }
}

onMounted(() => loadData())
</script>

<style scoped>
.page-container { padding: 16px; }
.search-form { margin-bottom: 16px; padding: 16px; background: var(--color-fill-lightest); border-radius: 4px; }
</style>
