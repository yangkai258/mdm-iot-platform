<template>
  <div class="page-container">
    <div class="search-form">
      <a-form :model="form" layout="inline">
        <a-form-item label="行为ID"><a-input v-model="form.id" placeholder="请输入行为ID" /></a-form-item>
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
        <a-form-item label="名称"><a-input v-model="form.name" /></a-form-item>
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
import { useRoute, useRouter } from 'vue-router'
import { Message } from '@arco-design/web-vue'
import { getAiMonitorEventById } from '@/api/ai'

const route = useRoute()
const router = useRouter()

const loading = ref(false)
const data = ref([])
const modalVisible = ref(false)
const modalTitle = ref('新建')

const form = reactive({
  id: ''
})

const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0
})

const columns = [
  { title: '时间', dataIndex: 'created_at', width: 170 },
  { title: '设备', dataIndex: 'device_id', width: 140 },
  { title: '用户', dataIndex: 'user_id', width: 100 },
  { title: '模型', dataIndex: 'model_version', width: 160 },
  { title: '事件类型', dataIndex: 'behavior_type', width: 120 },
  { title: '延迟', dataIndex: 'latency_ms', width: 90 },
  { title: '置信度', dataIndex: 'confidence', width: 90 },
  { title: '状态', dataIndex: 'status', width: 90 }
]

const handleSearch = () => {
  loadData()
}

const handleReset = () => {
  form.id = ''
  loadData()
}

const handleCreate = () => {
  modalTitle.value = '新建'
  modalVisible.value = true
}

const handleSubmit = () => {
  modalVisible.value = false
  Message.success('保存成功')
}

const loadData = async () => {
  loading.value = true
  try {
    if (route.params.id) {
      const res = await getAiMonitorEventById(route.params.id)
      if (res.code === 0) {
        data.value = res.data ? [res.data] : []
      }
    } else {
      data.value = []
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
