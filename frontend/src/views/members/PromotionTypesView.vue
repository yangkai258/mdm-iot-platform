<template>
  <div class="promotion-page">
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>会员管理</a-breadcrumb-item>
      <a-breadcrumb-item>促销活动</a-breadcrumb-item>
      <a-breadcrumb-item>促销类型</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 搜索筛选 -->
    <a-card class="action-card">
      <a-space wrap>
        <a-input-search v-model="filters.keyword" placeholder="搜索类型名称" style="width: 220px" search-button @search="loadData" />
        <a-button type="primary" @click="showCreateDrawer">新建类型</a-button>
        <a-button @click="loadData">刷新</a-button>
      </a-space>
    </a-card>

    <!-- 类型列表 -->
    <a-card class="table-card">
      <a-table
        :columns="columns"
        :data="dataList"
        :loading="loading"
        :pagination="paginationConfig"
        @page-change="onPageChange"
        @page-size-change="onPageSizeChange"
        row-key="id"
        :scroll="{ x: 800 }"
      >
        <template #status="{ record }">
          <a-tag :color="record.enabled ? 'green' : 'gray'">{{ record.enabled ? '启用' : '禁用' }}</a-tag>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="showEdit(record)">编辑</a-button>
            <a-button type="text" size="small" status="danger" @click="handleDelete(record)">删除</a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 新建/编辑抽屉 -->
    <a-drawer v-model:visible="formVisible" :title="isEdit ? '编辑促销类型' : '新建促销类型'" :width="520">
      <a-form :model="form" layout="vertical">
        <a-form-item label="类型名称" required>
          <a-input v-model="form.name" placeholder="请输入类型名称" />
        </a-form-item>
        <a-form-item label="类型编码">
          <a-input v-model="form.code" placeholder="请输入类型编码" :disabled="isEdit" />
        </a-form-item>
        <a-form-item label="描述">
          <a-textarea v-model="form.description" :rows="3" placeholder="描述该促销类型" />
        </a-form-item>
        <a-form-item label="状态">
          <a-switch v-model="form.enabled" />
        </a-form-item>
      </a-form>
      <template #footer>
        <a-button @click="formVisible = false">取消</a-button>
        <a-button type="primary" :loading="formLoading" @click="handleFormSubmit">{{ isEdit ? '保存' : '创建' }}</a-button>
      </template>
    </a-drawer>
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
.promotion-page { padding: 20px 24px; min-height: calc(100vh - 64px); background: #f5f7fa; }
.breadcrumb { margin-bottom: 16px; }
.action-card { margin-bottom: 16px; }
.table-card { border-radius: 8px; }
</style>
