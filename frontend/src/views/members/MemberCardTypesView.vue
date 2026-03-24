<template>
  <div class="card-types-page">
    <!-- 面包屑 -->
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>会员管理</a-breadcrumb-item>
      <a-breadcrumb-item>会员卡类型</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 搜索筛选区 -->
    <a-card class="action-card">
      <a-space wrap>
        <a-input-search
          v-model="filters.keyword"
          placeholder="搜索卡类型名称"
          style="width: 220px"
          search-button
          @search="handleSearch"
        />
        <a-select v-model="filters.status" placeholder="状态" allow-clear style="width: 120px" @change="handleSearch">
          <a-option :value="1">启用</a-option>
          <a-option :value="0">禁用</a-option>
        </a-select>
        <a-button @click="handleSearch">筛选</a-button>
        <a-button @click="resetFilters">重置</a-button>
      </a-space>
    </a-card>

    <!-- 操作+表格 -->
    <a-card class="table-card">
      <template #title>
        <a-space>
          <span style="font-weight: 600; font-size: 15px;">会员卡类型</span>
          <a-badge :count="pagination.total" :max-count="99999" />
        </a-space>
      </template>
      <template #extra>
        <a-space>
          <a-button @click="handleExport">导出</a-button>
          <a-button type="primary" @click="showCreateModal">新建</a-button>
        </a-space>
      </template>

      <a-table
        :columns="columns"
        :data="dataList"
        :loading="loading"
        :pagination="paginationConfig"
        @page-change="onPageChange"
        @page-size-change="onPageSizeChange"
        row-key="id"
        :scroll="{ x: 900 }"
      >
        <template #status="{ record }">
          <a-tag :color="record.status === 1 ? 'green' : 'gray'">{{ record.status === 1 ? '启用' : '禁用' }}</a-tag>
        </template>
        <template #discountRate="{ record }">
          <span>{{ (record.discountRate * 100).toFixed(0) }}%</span>
        </template>
        <template #pointsRate="{ record }">
          <span>{{ (record.pointsRate || 1).toFixed(1) }}x</span>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="showEdit(record)">编辑</a-button>
            <a-button type="text" size="small" status="danger" @click="handleDelete(record)">删除</a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 新增/编辑弹窗 -->
    <a-modal
      v-model:visible="formVisible"
      :title="isEdit ? '编辑卡类型' : '新建卡类型'"
      @before-ok="handleFormSubmit"
      @cancel="formVisible = false"
      :width="480"
      :loading="formLoading"
    >
      <a-form :model="form" layout="vertical" ref="formRef">
        <a-form-item label="卡类型名称" field="name" :rules="[{ required: true, message: '请输入卡类型名称' }]">
          <a-input v-model="form.name" placeholder="请输入卡类型名称" />
        </a-form-item>
        <a-form-item label="折扣率" field="discountRate" :rules="[{ required: true, message: '请输入折扣率' }]">
          <a-input-number v-model="form.discountRate" :min="0" :max="1" :step="0.01" placeholder="如 0.9 表示9折" style="width: 100%" />
          <template #extra>输入0~1之间的小数，如 0.9 表示9折，1 表示不打折</template>
        </a-form-item>
        <a-form-item label="积分倍率" field="pointsRate">
          <a-input-number v-model="form.pointsRate" :min="0" :max="100" :step="0.1" placeholder="如 1.5 表示1.5倍积分" style="width: 100%" />
          <template #extra>积分倍率，默认为 1.0</template>
        </a-form-item>
        <a-form-item label="状态">
          <a-radio-group v-model="form.status">
            <a-radio :value="1">启用</a-radio>
            <a-radio :value="0">禁用</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item label="描述">
          <a-textarea v-model="form.description" :rows="2" placeholder="描述信息" />
        </a-form-item>
      </a-form>
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

const filters = reactive({ keyword: '', status: undefined })
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
  discountRate: 1,
  pointsRate: 1,
  status: 1,
  description: ''
})

const columns = [
  { title: 'ID', dataIndex: 'id', width: 70 },
  { title: '卡类型名称', dataIndex: 'name', width: 160 },
  { title: '折扣率', slotName: 'discountRate', width: 100 },
  { title: '积分倍率', slotName: 'pointsRate', width: 100 },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '描述', dataIndex: 'description', ellipsis: true },
  { title: '操作', slotName: 'actions', width: 120, fixed: 'right' }
]

const loadData = async () => {
  loading.value = true
  try {
    const params = { page: pagination.current, pageSize: pagination.pageSize }
    if (filters.keyword) params.keyword = filters.keyword
    if (filters.status !== undefined) params.status = filters.status

    const res = await api.getCardTypeList(params)
    const d = res.data || {}
    dataList.value = d.list || []
    pagination.total = d.total || 0
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
  filters.status = undefined
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
  Object.assign(form, { name: '', discountRate: 1, pointsRate: 1, status: 1, description: '' })
  formVisible.value = true
}

const showEdit = (record) => {
  isEdit.value = true
  currentId.value = record.id
  Object.assign(form, {
    name: record.name,
    discountRate: record.discountRate ?? 1,
    pointsRate: record.pointsRate ?? 1,
    status: record.status ?? 1,
    description: record.description || ''
  })
  formVisible.value = true
}

const handleFormSubmit = async (done) => {
  formLoading.value = true
  try {
    if (isEdit.value) {
      await api.updateCardType(currentId.value, { ...form })
      Message.success('更新成功')
    } else {
      await api.createCardType({ ...form })
      Message.success('创建成功')
    }
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
    await api.deleteCardType(record.id)
    Message.success('删除成功')
    loadData()
  } catch (err) {
    Message.error(err.message || '删除失败')
  }
}

const handleExport = () => {
  Message.info('导出功能开发中')
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.card-types-page {
  padding: 20px 24px;
  min-height: calc(100vh - 64px);
  background: #f5f7fa;
}
.breadcrumb { margin-bottom: 16px; }
.action-card { margin-bottom: 16px; }
.table-card { border-radius: 8px; }
</style>
