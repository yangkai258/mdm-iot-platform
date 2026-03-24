<template>
  <div class="card-groups-page">
    <!-- 面包屑 -->
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>会员管理</a-breadcrumb-item>
      <a-breadcrumb-item>会员卡分组</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 搜索筛选区 -->
    <a-card class="action-card">
      <a-space wrap>
        <a-input-search
          v-model="filters.keyword"
          placeholder="搜索分组名称"
          style="width: 220px"
          search-button
          @search="handleSearch"
        />
        <a-button @click="handleSearch">筛选</a-button>
        <a-button @click="resetFilters">重置</a-button>
      </a-space>
    </a-card>

    <!-- 操作+表格 -->
    <a-card class="table-card">
      <template #title>
        <a-space>
          <span style="font-weight: 600; font-size: 15px;">会员卡分组</span>
          <a-badge :count="pagination.total" :max-count="99999" />
        </a-space>
      </template>
      <template #extra>
        <a-space>
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
        :scroll="{ x: 800 }"
      >
        <template #cardTypeCount="{ record }">
          <span>{{ record.cardTypeCount || 0 }}</span>
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
      :title="isEdit ? '编辑分组' : '新建分组'"
      @before-ok="handleFormSubmit"
      @cancel="formVisible = false"
      :width="460"
      :loading="formLoading"
    >
      <a-form :model="form" layout="vertical">
        <a-form-item label="分组名称" field="name" :rules="[{ required: true, message: '请输入分组名称' }]">
          <a-input v-model="form.name" placeholder="请输入分组名称" />
        </a-form-item>
        <a-form-item label="描述">
          <a-textarea v-model="form.description" :rows="3" placeholder="描述信息" />
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

const form = reactive({ name: '', description: '' })

const columns = [
  { title: 'ID', dataIndex: 'id', width: 70 },
  { title: '分组名称', dataIndex: 'name', width: 200 },
  { title: '描述', dataIndex: 'description', ellipsis: true },
  { title: '包含卡类型数量', slotName: 'cardTypeCount', width: 150 },
  { title: '操作', slotName: 'actions', width: 120, fixed: 'right' }
]

const loadData = async () => {
  loading.value = true
  try {
    const params = { page: pagination.current, pageSize: pagination.pageSize }
    if (filters.keyword) params.keyword = filters.keyword

    const res = await api.getCardGroupList(params)
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
  Object.assign(form, { name: '', description: '' })
  formVisible.value = true
}

const showEdit = (record) => {
  isEdit.value = true
  currentId.value = record.id
  Object.assign(form, { name: record.name, description: record.description || '' })
  formVisible.value = true
}

const handleFormSubmit = async (done) => {
  formLoading.value = true
  try {
    if (isEdit.value) {
      await api.updateCardGroup(currentId.value, { ...form })
      Message.success('更新成功')
    } else {
      await api.createCardGroup({ ...form })
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
    await api.deleteCardGroup(record.id)
    Message.success('删除成功')
    loadData()
  } catch (err) {
    Message.error(err.message || '删除失败')
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.card-groups-page {
  padding: 20px 24px;
  min-height: calc(100vh - 64px);
  background: #f5f7fa;
}
.breadcrumb { margin-bottom: 16px; }
.action-card { margin-bottom: 16px; }
.table-card { border-radius: 8px; }
</style>
