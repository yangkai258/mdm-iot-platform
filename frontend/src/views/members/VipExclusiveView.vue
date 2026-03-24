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
    <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id">
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
const levelOptions = ref([])
const productOptions = ref([])
const loading = ref(false)
const formLoading = ref(false)
const formVisible = ref(false)
const isEdit = ref(false)

const filters = reactive({ keyword: '' })
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })

const paginationConfig = computed(() => ({
  current: pagination.current, pageSize: pagination.pageSize, total: pagination.total,
  showTotal: true, showPageSize: true, pageSizeOptions: [10, 20, 50, 100]
}))

const form = reactive({
  name: '', levelId: undefined, discountType: 'percent',
  discountValue: 1, productIds: [], description: ''
})

const columns = [
  { title: '活动名称', dataIndex: 'name', width: 200 },
  { title: '适用等级', slotName: 'level', width: 140 },
  { title: '优惠内容', slotName: 'discount', width: 140 },
  { title: '适用商品', dataIndex: 'productName', width: 180, ellipsis: true },
  { title: '操作', slotName: 'actions', width: 150 }
]

const loadData = async () => {
  loading.value = true
  try {
    const d = {
      list: [
        { id: 1, name: 'VIP专属折扣', levelName: '黄金会员', discountType: 'percent', discountValue: 0.85, productName: '全场商品' },
        { id: 2, name: '钻石立减', levelName: '钻石会员', discountType: 'amount', discountValue: 50, productName: '指定商品' }
      ],
      total: 2
    }
    dataList.value = d.list
    pagination.total = d.total
  } catch (err) {
    dataList.value = []
  } finally {
    loading.value = false
  }
}

const loadLevelOptions = async () => {
  try {
    const res = await api.getLevelList()
    levelOptions.value = res.data || []
  } catch (err) { /* ignore */ }
}

const showCreateDrawer = () => {
  isEdit.value = false
  loadLevelOptions()
  Object.assign(form, { name: '', levelId: undefined, discountType: 'percent', discountValue: 1, productIds: [], description: '' })
  formVisible.value = true
}

const showEdit = (record) => {
  isEdit.value = true
  loadLevelOptions()
  Object.assign(form, { name: record.name, levelId: record.levelId, discountType: record.discountType, discountValue: record.discountValue, productIds: [], description: record.description || '' })
  formVisible.value = true
}

const handleFormSubmit = async () => {
  if (!form.name || !form.levelId) { Message.warning('请填写名称和选择等级'); return }
  formLoading.value = true
  try {
    await new Promise(r => setTimeout(r, 500))
    Message.success(isEdit.value ? '更新成功' : '创建成功')
    formVisible.value = false
    loadData()
  } catch (err) {
    Message.error(err.message || '操作失败')
  } finally {
    formLoading.value = false
  }
}

const handleDelete = async (record) => {
  try {
    await new Promise(r => setTimeout(r, 300))
    Message.success('删除成功')
    loadData()
  } catch (err) {
    Message.error(err.message || '删除失败')
  }
}

const onPageChange = (page) => { pagination.current = page; loadData() }
const onPageSizeChange = (pageSize) => { pagination.pageSize = pageSize; pagination.current = 1; loadData() }

onMounted(() => loadData())

</script>

<style scoped>
.page-container { background: #fff; border-radius: 4px; padding: 20px; }
.search-form { margin-bottom: 16px; padding: 16px; background: #f7f8fa; border-radius: 4px; }
.toolbar { margin-bottom: 16px; }
</style>
