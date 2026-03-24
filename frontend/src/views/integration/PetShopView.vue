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
        <a-form-item label="商品名称"><a-input v-model="form.name" /></a-form-item>
        <a-form-item label="分类"><a-input v-model="form.category" /></a-form-item>
        <a-form-item label="价格"><a-input-number v-model="form.price" :min="0" style="width:100%" /></a-form-item>
        <a-form-item label="库存"><a-input-number v-model="form.stock" :min="0" style="width:100%" /></a-form-item>
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

const form = reactive({ id: '', name: '', category: '', price: 0, stock: 0 })

const columns = [
  { title: '名称', dataIndex: 'name' },
  { title: '分类', dataIndex: 'category' },
  { title: '价格', dataIndex: 'price' },
  { title: '库存', dataIndex: 'stock' }
]

const pagination = reactive({ total: 0, current: 1, pageSize: 10 })
const data = ref([])

const loadProducts = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const res = await fetch('/api/v1/integration/products', {
      headers: { 'Authorization': `Bearer ${token}` }
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
    { id: '1', name: '天然狗粮 5kg', category: '食品', price: 299, stock: 50 },
    { id: '2', name: '智能逗猫棒', category: '玩具', price: 129, stock: 30 },
    { id: '3', name: '宠物冬季外套', category: '服装', price: 199, stock: 25 },
    { id: '4', name: '宠物益生菌', category: '保健', price: 89, stock: 100 }
  ]
}

const handleSearch = () => loadProducts()
const handleReset = () => { form.name = ''; loadProducts() }

const handleCreate = () => {
  isEdit.value = false
  modalTitle.value = '新建'
  Object.assign(form, { id: '', name: '', category: '', price: 0, stock: 0 })
  modalVisible.value = true
}

const handleSubmit = () => {
  if (!form.name) { Message.warning('请填写商品名称'); return }
  if (isEdit.value) {
    const idx = data.value.findIndex(p => p.id === form.id)
    if (idx !== -1) data.value[idx] = { ...form }
    Message.success('编辑成功')
  } else {
    data.value.unshift({ ...form, id: Date.now().toString() })
    pagination.total++
    Message.success('添加成功')
  }
  modalVisible.value = false
}

onMounted(() => loadProducts())
</script>

<style scoped>
.page-container { background: #fff; border-radius: 4px; padding: 20px; }
.search-form { margin-bottom: 16px; padding: 16px; background: #f7f8fa; border-radius: 4px; }
.toolbar { margin-bottom: 16px; }
</style>
