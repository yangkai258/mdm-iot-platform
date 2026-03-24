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

const dataList = ref([])
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

const form = reactive({ name: '', code: '', description: '', enabled: true })

const columns = [
  { title: '类型名称', dataIndex: 'name', width: 200 },
  { title: '类型编码', dataIndex: 'code', width: 160 },
  { title: '描述', dataIndex: 'description', width: 300, ellipsis: true },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '操作', slotName: 'actions', width: 150 }
]

const getEnabledText = (enabled) => enabled ? '启用' : '禁用'
const getEnabledColor = (enabled) => enabled ? 'green' : 'gray'

const loadData = async () => {
  loading.value = true
  try {
    const d = {
      list: [
        { id: 1, name: '满额减', code: 'amount_reduce', description: '消费满指定金额减免一定金额', enabled: true },
        { id: 2, name: '满额折', code: 'amount_discount', description: '消费满指定金额享受折扣优惠', enabled: true },
        { id: 3, name: '直减', code: 'direct_reduce', description: '直接减免指定金额', enabled: true },
        { id: 4, name: '买赠', code: 'buy_gift', description: '购买指定商品赠送赠品', enabled: true },
        { id: 5, name: '红包', code: 'redpacket', description: '会员红包优惠', enabled: true }
      ],
      total: 5
    }
    dataList.value = d.list
    pagination.total = d.total
  } catch (err) {
    dataList.value = []
  } finally {
    loading.value = false
  }
}

const showCreateDrawer = () => {
  isEdit.value = false
  Object.assign(form, { name: '', code: '', description: '', enabled: true })
  formVisible.value = true
}

const showEdit = (record) => {
  isEdit.value = true
  Object.assign(form, { name: record.name, code: record.code, description: record.description || '', enabled: record.enabled })
  formVisible.value = true
}

const handleFormSubmit = async () => {
  if (!form.name) { Message.warning('请填写类型名称'); return }
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
