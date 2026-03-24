<template>
  <div class="page-container">
    <div class="search-form">
      <a-form :model="form" layout="inline">
        <a-form-item label="告警内容"><a-input v-model="form.keyword" placeholder="搜索告警内容" /></a-form-item>
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

const loading = ref(false)
const data = ref([])
const modalVisible = ref(false)
const modalTitle = ref('新建')

const form = reactive({
  keyword: '',
  severity: '',
  category: ''
})

const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0
})

const columns = [
  { title: '告警内容', dataIndex: 'message', width: 280, ellipsis: true },
  { title: '规则名称', dataIndex: 'rule_name', width: 150 },
  { title: '设备', dataIndex: 'device_name', width: 120 },
  { title: '严重级别', dataIndex: 'severity', width: 100 },
  { title: '状态', dataIndex: 'status', width: 100 },
  { title: '触发时间', dataIndex: 'created_at', width: 160 }
]

const handleSearch = () => {
  loadData()
}

const handleReset = () => {
  form.keyword = ''
  form.severity = ''
  form.category = ''
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
  setTimeout(() => {
    loading.value = false
  }, 300)
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
