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
    <a-table :columns="columns" :data="dataList" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id">
      <template #actions="{ record }">
        <a-button type="text" size="small" @click="handleEdit(record)">编辑</a-button>
        <a-button type="text" size="small" @click="handleDelete(record)">删除</a-button>
      </template>
    </a-table>
    <a-modal v-model:visible="modalVisible" :title="modalTitle" @before-ok="handleSubmit" @cancel="modalVisible = false">
      <a-form :model="form" label-col-flex="100px">
        <a-form-item label="名称"><a-input v-model="form.name" placeholder="请输入" /></a-form-item>
      </a-form>
      <template #footer>
        <a-button @click="modalVisible = false">取消</a-button>
        <a-button type="primary" @click="handleSubmit">确定</a-button>
      </template>
    </a-modal>
  </div>
</template>

<script setup>

import { ref, reactive, computed, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import * as api from '@/api/member'

const dataList = ref([])
const loading = ref(false)
const formLoading = ref(false)
const formVisible = ref(false)
const isEdit = ref(false)
const currentId = ref(null)

const filters = reactive({ keyword: '' })
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })

const paginationConfig = computed(() => ({
  current: pagination.current,
  pageSize: pagination.pageSize,
  total: pagination.total,
  showTotal: true,
  showPageSize: true,
  pageSizeOptions: [10, 20, 50, 100]
}))

const form = reactive({
  name: '',
  conditionType: 'consume_amount',
  conditionValue: 0,
  requiredPoints: 0,
  description: ''
})

const columns = [
  { title: 'ID', dataIndex: 'id', width: 70 },
  { title: '等级名称', dataIndex: 'name', width: 160 },
  { title: '升级条件类型', slotName: 'conditionType', width: 150 },
  { title: '升级条件值', slotName: 'conditionValue', width: 160 },
  { title: '所需积分', slotName: 'requiredPoints', width: 120 },
  { title: '描述', dataIndex: 'description', ellipsis: true },
  { title: '操作', slotName: 'actions', width: 120, fixed: 'right' }
]

const getConditionLabel = (type) => ({
  consume_amount: '累计消费金额',
  order_count: '累计消费次数',
  points: '累计积分'
}[type] || type)

const conditionColor = (type) => ({
  consume_amount: 'blue',
  order_count: 'green',
  points: 'orange'
}[type] || 'gray')

const loadData = async () => {
  loading.value = true
  try {
    const params = { page: pagination.current, pageSize: pagination.pageSize }
    if (filters.keyword) params.keyword = filters.keyword

    const res = await api.getUpgradeRules()
    const d = res.data || {}
    dataList.value = Array.isArray(d) ? d : (d.list || [])
    pagination.total = dataList.value.length
  } catch (err) {
    Message.error('加载失败: ' + err.message)
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.current = 1
  loadData()
}

const resetFilters = () => {
  filters.keyword = ''
  pagination.current = 1
  loadData()
}

const onPageChange = (page) => {
  pagination.current = page
  loadData()
}

const onPageSizeChange = (pageSize) => {
  pagination.pageSize = pageSize
  pagination.current = 1
  loadData()
}

const showCreateModal = () => {
  isEdit.value = false
  Object.assign(form, { name: '', conditionType: 'consume_amount', conditionValue: 0, requiredPoints: 0, description: '' })
  formVisible.value = true
}

const showEdit = (record) => {
  isEdit.value = true
  currentId.value = record.id
  Object.assign(form, {
    name: record.name,
    conditionType: record.conditionType || 'consume_amount',
    conditionValue: record.conditionValue ?? 0,
    requiredPoints: record.requiredPoints ?? 0,
    description: record.description || ''
  })
  formVisible.value = true
}

const handleFormSubmit = async (done) => {
  formLoading.value = true
  try {
    const payload = { ...form }
    await api.updateUpgradeRules(payload)
    Message.success('保存成功')
    formVisible.value = false
    loadData()
    done(true)
  } catch (err) {
    Message.error(err.message || '操作失败')
    done(false)
  } finally {
    formLoading.value = false
  }
}

const handleDelete = async (record) => {
  try {
    // 降级规则没有单独的 delete，先用 update 模拟
    Message.info('请通过编辑调整规则')
  } catch (err) {
    Message.error(err.message || '操作失败')
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
