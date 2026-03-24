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
        <a-form-item label="表情包名称"><a-input v-model="form.name" /></a-form-item>
        <a-form-item label="预览"><a-input v-model="form.preview_emoji" placeholder="如: 😊" /></a-form-item>
        <a-form-item label="描述"><a-textarea v-model="form.description" :rows="2" /></a-form-item>
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

const form = reactive({ id: '', name: '', preview_emoji: '', description: '' })

const columns = [
  { title: '表情包名称', dataIndex: 'name' },
  { title: '分类', dataIndex: 'category_name', width: 100 },
  { title: '预览', dataIndex: 'preview_emoji', width: 80 },
  { title: '创作者', dataIndex: 'creator_name', width: 120 },
  { title: '收费类型', dataIndex: 'price_type', width: 100 },
  { title: '下载量', dataIndex: 'download_count', width: 100 }
]

const pagination = reactive({ total: 0, current: 1, pageSize: 20 })
const data = ref([])
const categories = ref([
  { category_id: 'cat1', category_name: '可爱' },
  { category_id: 'cat2', category_name: '搞笑' },
  { category_id: 'cat3', category_name: '日常' },
  { category_id: 'cat4', category_name: '节日' },
  { category_id: 'cat5', category_name: '宠物' }
])

const loadEmoticons = async () => {
  loading.value = true
  try {
    const res = await fetch(`/api/v1/market/emoticons?keyword=${form.name || ''}`)
    const resData = await res.json()
    if (resData.code === 0) {
      data.value = (resData.data?.list || []).map(e => ({
        ...e,
        price_type: e.is_premium ? '付费' : '免费'
      }))
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
    { emoticon_id: 'emo1', name: '开心喵', category_id: 'cat5', category_name: '宠物', preview_emoji: '😺', creator_name: '官方', is_premium: false, price_type: '免费', download_count: 1234 },
    { emoticon_id: 'emo2', name: '搞怪汪', category_id: 'cat5', category_name: '宠物', preview_emoji: '🐶', creator_name: '官方', is_premium: false, price_type: '免费', download_count: 856 },
    { emoticon_id: 'emo3', name: '快乐每一天', category_id: 'cat3', category_name: '日常', preview_emoji: '😊', creator_name: '用户小明', is_premium: true, price_type: '付费', download_count: 2341 },
    { emoticon_id: 'emo4', name: '节日祝福', category_id: 'cat4', category_name: '节日', preview_emoji: '🎉', creator_name: '官方', is_premium: false, price_type: '免费', download_count: 5678 },
    { emoticon_id: 'emo5', name: '卖萌专用', category_id: 'cat1', category_name: '可爱', preview_emoji: '🥰', creator_name: '用户小红', is_premium: true, price_type: '付费', download_count: 4321 }
  ]
}

const handleSearch = () => loadEmoticons()
const handleReset = () => { form.name = ''; loadEmoticons() }

const handleCreate = () => {
  isEdit.value = false
  modalTitle.value = '新建'
  Object.assign(form, { id: '', name: '', preview_emoji: '', description: '' })
  modalVisible.value = true
}

const handleSubmit = () => {
  if (!form.name) { Message.warning('请填写表情包名称'); return }
  if (isEdit.value) {
    const idx = data.value.findIndex(e => e.emoticon_id === form.id)
    if (idx !== -1) data.value[idx] = { ...data.value[idx], name: form.name }
    Message.success('编辑成功')
  } else {
    data.value.unshift({
      emoticon_id: `emo${Date.now()}`,
      name: form.name,
      category_id: null,
      category_name: '默认',
      preview_emoji: form.preview_emoji || '😊',
      creator_name: '当前用户',
      is_premium: false,
      price_type: '免费',
      download_count: 0
    })
    pagination.total++
    Message.success('添加成功')
  }
  modalVisible.value = false
}

onMounted(() => { loadEmoticons() })
</script>

<style scoped>
.page-container { background: #fff; border-radius: 4px; padding: 20px; }
.search-form { margin-bottom: 16px; padding: 16px; background: #f7f8fa; border-radius: 4px; }
.toolbar { margin-bottom: 16px; }
</style>
