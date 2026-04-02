<template>
    <Breadcrumb :items="['Home','Console','']" />


  <div class="knowledge-page">

    <a-card class="general-card">
      <template #title><span class="card-title">知识库查询</span></template>
      <a-row :gutter="16">
        <a-col :flex="1">
          <a-form :model="searchForm" layout="vertical" size="small">
            <a-row :gutter="16">
              <a-col :span="8">
                <a-form-item label="名称">
                  <a-input v-model="searchForm.name" placeholder="请输入" allow-clear />
                </a-form-item>
              </a-col>
              <a-col :span="8">
                <a-form-item label="分类">
                  <a-input v-model="searchForm.category" placeholder="请输入" allow-clear />
                </a-form-item>
              </a-col>
            </a-row>
          </a-form>
        </a-col>
        <a-divider style="height: 84px" direction="vertical" />
        <a-col :flex="'86px'" style="text-align: right">
          <a-space direction="vertical" :size="18">
            <a-button type="primary" @click="handleSearch">
              <template #icon><icon-search /></template>
              查询
            </a-button>
            <a-button @click="handleReset">
              <template #icon><icon-refresh /></template>
              重置
            </a-button>
          </a-space>
        </a-col>
      </a-row>
    </a-card>

    <a-card class="general-card" style="margin-top: 16px">
      <template #title><span class="card-title">知识列表</span></template>
      <template #extra>
        <a-button type="primary" @click="handleCreate">
          <template #icon><icon-plus /></template>
          新建
        </a-button>
      </template>
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id" />
    </a-card>

      </a-table>
    <a-modal v-model:visible="modalVisible" :title="modalTitle" @ok="handleSubmit" :width="560">
      <a-form :model="form" layout="vertical">
        <a-form-item label="分类">
          <a-input v-model="form.category" placeholder="请输入分类" />
        </a-form-item>
        <a-form-item label="问题">
          <a-input v-model="form.question" placeholder="请输入问题" />
        </a-form-item>
        <a-form-item label="答案">
          <a-textarea v-model="form.answer" placeholder="请输入答案" :rows="4" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const loading = ref(false)
const modalVisible = ref(false)
const isEdit = ref(false)
const modalTitle = computed(() => isEdit.value ? '编辑知识' : '新建知识')

const searchForm = reactive({ name: '', category: '' })

const form = reactive({ id: null, category: '', question: '', answer: '' })

const data = ref([])
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })

const columns = [
  { title: 'ID', dataIndex: 'id', width: 80 },
  { title: '分类', dataIndex: 'category', width: 140 },
  { title: '问题', dataIndex: 'question', ellipsis: true },
  { title: '答案', dataIndex: 'answer', ellipsis: true },
  { title: '更新时间', dataIndex: 'updated_at', width: 170 },
  { title: '操作', width: 120 }
]

const filteredData = computed(() => {
  let result = data.value
  if (searchForm.name) result = result.filter(d => d.question?.includes(searchForm.name))
  if (searchForm.category) result = result.filter(d => d.category?.includes(searchForm.category))
  return result
})

const loadData = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`${'/api/v1'}/knowledge?page=${pagination.current}&page_size=${pagination.pageSize}`, {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    const resData = await res.json()
    if (resData.code === 0) {
      data.value = resData.data?.list || []
      pagination.total = resData.data?.total || 0
    }
  } catch (e) {
    data.value = getMockData()
    pagination.total = data.value.length
  } finally { loading.value = false }
}

const getMockData = () => [
  { id: 1, category: '设备使用', question: '设备无法开机怎么办？', answer: '请检查电源连接或长按电源键10秒强制重启', updated_at: '2026-03-20 10:00:00' },
  { id: 2, category: '固件升级', question: 'OTA升级失败如何处理？', answer: '检查网络连接，确保设备在线后重试', updated_at: '2026-03-19 15:30:00' }
]

const handleSearch = () => { pagination.current = 1; loadData() }
const handleReset = () => { searchForm.name = ''; searchForm.category = ''; pagination.current = 1; loadData() }
const onPageChange = (page) => { pagination.current = page; loadData() }

const handleCreate = () => {
  isEdit.value = false
  Object.assign(form, { id: null, category: '', question: '', answer: '' })
  modalVisible.value = true
}

const handleSubmit = async () => {
  if (!form.question || !form.answer) { Message.warning('请填写完整信息'); return }
  try {
    const token = localStorage.getItem('token')
    const method = isEdit.value ? 'PUT' : 'POST'
    const url = isEdit.value ? `${'/api/v1'}/knowledge/${form.id}` : `${'/api/v1'}/knowledge`
    const res = await fetch(url, {
      method, headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' },
      body: JSON.stringify(form)
    })
    const data = await res.json()
    if (data.code === 0) { Message.success('保存成功'); modalVisible.value = false; loadData() }
  } catch (e) {
    if (isEdit.value) {
      const idx = data.value.findIndex(d => d.id === form.id)
      if (idx !== -1) data.value[idx] = { ...data.value[idx], ...form, updated_at: new Date().toLocaleString() }
    } else {
      data.value.unshift({ ...form, id: Date.now(), updated_at: new Date().toLocaleString() })
    }
    Message.success('保存成功（模拟）')
    modalVisible.value = false
  }
}

onMounted(() => { loadData() })
</script>

<style scoped>
.knowledge-page { padding: 20px 24px; min-height: calc(100vh - 64px); background: #f5f7fa; }
.general-card { border-radius: 8px; }
.card-title { font-weight: 600; font-size: 15px; }
</style>
