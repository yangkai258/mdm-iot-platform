<template>
  <div class="page-container">
    <div class="search-form">
      <a-form :model="form" layout="inline">
        <a-form-item label="群组名称"><a-input v-model="form.name" placeholder="请输入" /></a-form-item>
        <a-form-item>
          <a-button type="primary" @click="handleSearch">搜索</a-button>
          <a-button @click="handleReset">重置</a-button>
        </a-form-item>
      </a-form>
    </div>
    <div class="toolbar">
      <a-button type="primary" @click="handleCreate">新建群组</a-button>
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
import * as analytics from '@/api/analytics'

const loading = ref(false)
const data = ref([])
const modalVisible = ref(false)
const modalTitle = ref('新建群组')

const form = reactive({
  name: '',
  type: 'registration',
  description: ''
})

const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0
})

const columns = [
  { title: '群组名称', dataIndex: 'name', width: 200 },
  { title: '群组类型', dataIndex: 'type', width: 120 },
  { title: '成员数量', dataIndex: 'size', width: 100 },
  { title: '创建时间', dataIndex: 'created_at', width: 170 }
]

const handleSearch = () => {
  loadData()
}

const handleReset = () => {
  form.name = ''
  form.type = 'registration'
  form.description = ''
  loadData()
}

const handleCreate = () => {
  modalTitle.value = '新建群组'
  modalVisible.value = true
}

const handleSubmit = async () => {
  modalVisible.value = false
  try {
    await analytics.createCohort(form)
    Message.success('保存成功')
    loadData()
  } catch (e) {
    Message.error('保存失败')
  }
}

const loadData = async () => {
  loading.value = true
  try {
    const res = await analytics.getCohortList({ keyword: form.name })
    data.value = res.data?.list || []
    pagination.total = data.value.length
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
