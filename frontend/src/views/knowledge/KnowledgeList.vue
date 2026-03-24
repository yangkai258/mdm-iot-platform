<template>
  <div class="page-container">
    <div class="search-form">
      <a-form :model="form" layout="inline">
        <a-form-item label="名称"><a-input v-model="form.name" placeholder="请输入" /></a-form-item>
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
        <a-form-item label="分类"><a-input v-model="form.category" /></a-form-item>
        <a-form-item label="问题"><a-input v-model="form.question" /></a-form-item>
        <a-form-item label="答案"><a-input v-model="form.answer" /></a-form-item>
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
const modalVisible = ref(false)
const modalTitle = ref('新建')
const isEdit = ref(false)

const form = reactive({ id: '', name: '', category: '', question: '', answer: '' })

const columns = [
  { title: '分类', dataIndex: 'category' },
  { title: '问题', dataIndex: 'question' },
  { title: '答案', dataIndex: 'answer' }
]

const pagination = reactive({ total: 0, current: 1, pageSize: 10 })
const data = ref([])

const loadData = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const res = await fetch('/api/v1/knowledge/list', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${token}` },
      body: JSON.stringify({ keyword: form.name })
    })
    const resData = await res.json()
    if (resData.code === 0) {
      data.value = resData.data || []
    } else {
      loadMockData()
    }
  } catch {
    loadMockData()
  } finally {
    pagination.total = data.value.length
    loading.value = false
  }
}

const loadMockData = () => {
  data.value = [
    { id: '1', category: '宠物喂养', question: '如何选择适合的猫粮？', answer: '选择猫粮时应关注成分表，优先选择肉类为主要成分的猫粮。' },
    { id: '2', category: '健康护理', question: '宠物疫苗接种时间？', answer: '幼猫/犬在6-8周开始接种疫苗，共需3-4次。' },
    { id: '3', category: '行为训练', question: '如何纠正狗狗乱叫？', answer: '通过正向训练，用奖励的方式鼓励安静行为。' }
  ]
}

const handleSearch = () => loadData()
const handleReset = () => { form.name = ''; loadData() }

const handleCreate = () => {
  isEdit.value = false
  modalTitle.value = '新建'
  Object.assign(form, { id: '', category: '', question: '', answer: '' })
  modalVisible.value = true
}

const handleSubmit = () => {
  if (!form.category || !form.question) {
    Message.warning('请填写完整信息')
    return
  }
  if (isEdit.value) {
    const idx = data.value.findIndex(item => item.id === form.id)
    if (idx !== -1) data.value[idx] = { ...form }
    Message.success('编辑成功')
  } else {
    data.value.unshift({ ...form, id: Date.now().toString() })
    pagination.total++
    Message.success('添加成功')
  }
  modalVisible.value = false
}

onMounted(() => { loadData() })
</script>

<style scoped>
.page-container { background: #fff; border-radius: 4px; padding: 20px; }
.search-form { margin-bottom: 16px; padding: 16px; background: #f7f8fa; border-radius: 4px; }
.toolbar { margin-bottom: 16px; }
</style>
