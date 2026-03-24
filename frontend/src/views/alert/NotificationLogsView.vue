<template>
  <div class="page-container">
    <div class="search-form">
      <a-form :model="form" layout="inline">
        <a-form-item label="渠道类型"><a-select v-model="form.channel_type" placeholder="选择渠道" style="width: 120px" /></a-form-item>
        <a-form-item label="发送状态"><a-select v-model="form.status" placeholder="选择状态" style="width: 120px" /></a-form-item>
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
import { Message } from '@arco-design/web-vue'
import { useNotificationLogs } from '@/composables/useNotification'

const { loading, logs, pagination, filters, loadLogs, getLogDetail } = useNotificationLogs()

const data = ref([])
const modalVisible = ref(false)
const modalTitle = ref('新建')

const form = reactive({
  channel_type: '',
  status: ''
})

const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0
})

const columns = [
  { title: '时间', dataIndex: 'created_at', width: 170 },
  { title: '渠道', dataIndex: 'channel_type', width: 90 },
  { title: '渠道名称', dataIndex: 'channel_name', width: 140 },
  { title: '接收人', dataIndex: 'recipient', width: 160 },
  { title: '主题', dataIndex: 'subject', ellipsis: true },
  { title: '状态', dataIndex: 'status', width: 90 },
  { title: '重试次数', dataIndex: 'attempt_count', width: 90 }
]

const handleSearch = () => {
  loadData()
}

const handleReset = () => {
  form.channel_type = ''
  form.status = ''
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
  try {
    await loadLogs()
    data.value = logs.value
    pagination.total = pagination.total
  } catch (e) {
    Message.error('加载失败')
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
