<template>
  <div class="page-container">
    <div class="search-form">
      <a-form :model="form" layout="inline">
        <a-form-item label="模型名称"><a-input v-model="form.model_name" placeholder="请输入" /></a-form-item>
        <a-form-item>
          <a-button type="primary" @click="handleSearch">搜索</a-button>
          <a-button @click="handleReset">重置</a-button>
        </a-form-item>
      </a-form>
    </div>
    <div class="toolbar">
      <a-button type="primary" @click="handleCreate">新建</a-button>
    </div>
    <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" />
    <a-modal v-model:visible="modalVisible" :title="modalTitle">
      <a-form :model="form" label-col-flex="100px">
        <a-form-item label="模型名称"><a-input v-model="form.model_name" /></a-form-item>
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
import { getAiModels, postAiModel } from '@/api/ai'

const loading = ref(false)
const data = ref([])
const modalVisible = ref(false)
const modalTitle = ref('新建')

const form = reactive({
  model_name: '',
  model_id: '',
  model_type: ''
})

const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0
})

const columns = [
  { title: '模型名称', dataIndex: 'model_name', width: 200 },
  { title: '当前版本', dataIndex: 'current_version', width: 160 },
  { title: '状态', dataIndex: 'status', width: 120 },
  { title: '发布时间', dataIndex: 'published_at', width: 160 }
]

const handleSearch = () => {
  pagination.current = 1
  loadData()
}

const handleReset = () => {
  form.model_name = ''
  form.model_id = ''
  form.model_type = ''
  pagination.current = 1
  loadData()
}

const handleCreate = () => {
  modalTitle.value = '新建'
  modalVisible.value = true
}

const handleSubmit = async () => {
  modalVisible.value = false
  Message.success('保存成功')
}

const loadData = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.current,
      page_size: pagination.pageSize,
      model_name: form.model_name || undefined
    }
    const res = await getAiModels(params)
    if (res.code === 0) {
      data.value = res.data?.list || []
      pagination.total = res.data?.total || 0
    }
  } catch (e) {
    Message.error('加载失败')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.page-container { background: #fff; border-radius: 4px; padding: 20px; }
.search-form { margin-bottom: 16px; padding: 16px; background: #f7f8fa; border-radius: 4px; }
.toolbar { margin-bottom: 16px; }
</style>
