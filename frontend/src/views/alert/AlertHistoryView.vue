<template>
  <div class="page-container">
    <div class="search-form">
      <a-form :model="form" layout="inline">
        <a-form-item label="设备"><a-input v-model="form.device_id" placeholder="设备ID" /></a-form-item>
        <a-form-item label="告警类型"><a-select v-model="form.alert_type" placeholder="选择类型" style="width: 140px" /></a-form-item>
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
import { useAlertHistory, SEVERITY_MAP, STATUS_MAP } from '@/composables/useNotification'

const { loading, alerts, pagination, filters, loadAlerts, confirmAlert, resolveAlert } = useAlertHistory()

const data = ref([])
const modalVisible = ref(false)
const modalTitle = ref('新建')

const form = reactive({
  device_id: '',
  alert_type: '',
  severity: '',
  status: ''
})

const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0
})

const columns = [
  { title: '时间', dataIndex: 'created_at', width: 170 },
  { title: '设备', dataIndex: 'device_id', width: 140 },
  { title: '告警类型', dataIndex: 'alert_type', width: 120 },
  { title: '级别', dataIndex: 'severity', width: 80 },
  { title: '告警消息', dataIndex: 'message', ellipsis: true },
  { title: '触发值', dataIndex: 'trigger_value', width: 100 },
  { title: '处理状态', dataIndex: 'status', width: 90 }
]

const handleSearch = () => {
  loadData()
}

const handleReset = () => {
  form.device_id = ''
  form.alert_type = ''
  form.severity = ''
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
    await loadAlerts()
    data.value = alerts.value
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
