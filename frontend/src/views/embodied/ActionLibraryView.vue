<template>
  <div class="page-container">
    <div class="search-form">
      <a-form :model="form" layout="inline">
        <a-form-item label="动作名称"><a-input v-model="form.action_name" placeholder="请输入动作名称" /></a-form-item>
        <a-form-item label="类型">
          <a-select v-model="form.action_type" placeholder="选择类型" allow-clear style="width: 120px">
            <a-option value="built-in">内置</a-option>
            <a-option value="learned">学习</a-option>
            <a-option value="custom">自定义</a-option>
          </a-select>
        </a-form-item>
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
    <a-modal v-model:visible="modalVisible" :title="modalTitle" @before-ok="handleSubmit">
      <a-form :model="form" label-col-flex="100px">
        <a-form-item label="动作名称" required><a-input v-model="form.action_name" placeholder="请输入动作名称" /></a-form-item>
        <a-form-item label="描述"><a-textarea v-model="form.description" :rows="3" placeholder="动作描述" /></a-form-item>
        <a-form-item label="难度">
          <a-select v-model="form.difficulty" placeholder="选择难度">
            <a-option value="easy">简单</a-option>
            <a-option value="medium">中等</a-option>
            <a-option value="hard">困难</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="标签">
          <a-select v-model="form.tags" multiple placeholder="选择标签" allow-create>
            <a-option value="dance">舞蹈</a-option>
            <a-option value="greeting">问候</a-option>
            <a-option value="exercise">运动</a-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getActionLibrary, recordAction, deleteAction } from '@/api/embodied'
import { Message, Modal } from '@arco-design/web-vue'

const loading = ref(false)
const data = ref<any[]>([])
const modalVisible = ref(false)
const modalTitle = ref('新建动作')
const form = ref<any>({
  action_name: '',
  description: '',
  difficulty: 'medium',
  tags: []
})

const columns = [
  { title: 'ID', dataIndex: 'id', width: 70 },
  { title: '动作名称', dataIndex: 'action_name', width: 160 },
  { title: '类型', dataIndex: 'action_type', width: 100 },
  { title: '难度', dataIndex: 'difficulty', width: 90 },
  { title: '时长', dataIndex: 'duration_ms', width: 80 },
  { title: '评分', dataIndex: 'score', width: 80 },
  { title: '创建时间', dataIndex: 'created_at', width: 170 }
]

const pagination = ref({ current: 1, pageSize: 20, total: 0, showTotal: true, showPageSize: true })

async function loadData() {
  try {
    loading.value = true
    const params: any = { page: pagination.value.current, page_size: pagination.value.pageSize }
    if (form.value.action_name) params.keyword = form.value.action_name
    if (form.value.action_type) params.action_type = form.value.action_type
    const res = await getActionLibrary(params)
    const resData = res.data
    data.value = resData?.actions || resData?.list || resData || []
    pagination.value.total = resData?.total || data.value.length
  } catch (err: any) {
    Message.error('加载失败: ' + err.message)
  } finally {
    loading.value = false
  }
}

function handleSearch() {
  pagination.value.current = 1
  loadData()
}

function handleReset() {
  form.value = { action_name: '', action_type: '', difficulty: '', tags: [] }
  handleSearch()
}

function handleCreate() {
  modalTitle.value = '新建动作'
  form.value = { action_name: '', description: '', difficulty: 'medium', tags: [] }
  modalVisible.value = true
}

async function handleSubmit(done: (val: boolean) => void) {
  try {
    await recordAction({ device_id: '', ...form.value })
    Message.success('创建成功')
    modalVisible.value = false
    loadData()
    done(true)
  } catch (err: any) {
    Message.error('创建失败: ' + err.message)
    done(false)
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
