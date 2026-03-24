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
        <a-form-item label="SMTP服务器"><a-input v-model="form.smtp_host" placeholder="smtp.example.com" /></a-form-item>
        <a-form-item label="端口"><a-input-number v-model="form.smtp_port" :min="1" :max="65535" style="width: 120px" /></a-form-item>
        <a-form-item label="用户名"><a-input v-model="form.smtp_user" /></a-form-item>
        <a-form-item label="密码"><a-input-password v-model="form.smtp_password" /></a-form-item>
        <a-form-item label="发件人"><a-input v-model="form.smtp_from" placeholder="noreply@example.com" /></a-form-item>
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
  smtp_host: '',
  smtp_port: 587,
  smtp_user: '',
  smtp_password: '',
  smtp_from: ''
})

const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0
})

const columns = [
  { title: '渠道名称', dataIndex: 'channel_name', width: 200 },
  { title: 'SMTP服务器', dataIndex: 'config', customRender: ({ record }) => record.config?.smtp_host || '-', width: 200 },
  { title: '端口', dataIndex: 'config', customRender: ({ record }) => record.config?.smtp_port || '-', width: 80 },
  { title: '发件人', dataIndex: 'config', customRender: ({ record }) => record.config?.smtp_from || '-', width: 180 }
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
      channel_type: 'email',
      channel_name: form.channel_name,
      enabled: true,
      config: {
        smtp_host: form.smtp_host,
        smtp_port: form.smtp_port,
        smtp_user: form.smtp_user,
        smtp_password: form.smtp_password,
        smtp_from: form.smtp_from
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
    data.value = channels.value.filter(ch => ch.channel_type === 'email')
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
