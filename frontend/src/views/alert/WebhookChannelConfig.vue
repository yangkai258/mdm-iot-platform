<template>
  <div class="page-container">
    <div class="search-form">
      <a-form :model="form" layout="inline">
        <a-form-item label="渠道名称"><a-input v-model="form.channel_name" placeholder="请输入" /></a-form-item>
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
        <a-form-item label="渠道名称"><a-input v-model="form.channel_name" /></a-form-item>
        <a-form-item label="Webhook URL"><a-input v-model="form.webhook_url" placeholder="https://example.com/webhook" /></a-form-item>
        <a-form-item label="请求方式"><a-select v-model="form.webhook_method" style="width: 120px">
          <a-option value="POST">POST</a-option>
          <a-option value="PUT">PUT</a-option>
        </a-select></a-form-item>
        <a-form-item label="密钥"><a-input-password v-model="form.webhook_secret" /></a-form-item>
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
import { useNotificationChannels } from '@/composables/useNotification'

const { loading, channels, loadChannels, createChannel } = useNotificationChannels()

const data = ref([])
const modalVisible = ref(false)
const modalTitle = ref('新建')

const form = reactive({
  channel_name: '',
  webhook_url: '',
  webhook_method: 'POST',
  webhook_secret: ''
})

const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0
})

const columns = [
  { title: '渠道名称', dataIndex: 'channel_name', width: 200 },
  { title: 'Webhook URL', dataIndex: 'config', customRender: ({ record }) => record.config?.webhook_url || '-', ellipsis: true },
  { title: '请求方式', dataIndex: 'config', customRender: ({ record }) => record.config?.webhook_method || '-', width: 100 }
]

const handleSearch = () => {
  loadData()
}

const handleReset = () => {
  form.channel_name = ''
  loadData()
}

const handleCreate = () => {
  modalTitle.value = '新建'
  modalVisible.value = true
}

const handleSubmit = async () => {
  modalVisible.value = false
  try {
    await createChannel({
      channel_type: 'webhook',
      channel_name: form.channel_name,
      enabled: true,
      config: {
        webhook_url: form.webhook_url,
        webhook_method: form.webhook_method,
        webhook_secret: form.webhook_secret
      }
    })
    Message.success('保存成功')
    loadData()
  } catch (e) {
    Message.error('保存失败')
  }
}

const loadData = async () => {
  try {
    await loadChannels()
    data.value = channels.value.filter(ch => ch.channel_type === 'webhook')
    pagination.total = data.value.length
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
