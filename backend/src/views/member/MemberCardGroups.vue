<template>
  <div class="member-card-groups-page">
    <!-- 面包屑 -->
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>会员管理</a-breadcrumb-item>
      <a-breadcrumb-item>会员卡分组</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 搜索筛选 -->
    <a-card class="action-card">
      <a-space wrap>
        <a-input-search v-model="filters.keyword" placeholder="搜索分组名称/编码" style="width: 260px" search-button @search="handleSearch" />
        <a-button @click="handleSearch">「搜索」</a-button>
        <a-button @click="resetFilters">「重置」</a-button>
      </a-space>
    </a-card>

    <!-- 操作栏 -->
    <a-card class="table-card">
      <template #extra>
        <a-space>
          <a-button type="primary" @click="showCreate">「新建」</a-button>
          <a-button @click="handleExport">「导出」</a-button>
          <a-button @click="loadData">🔄</a-button>
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
        <template #cardTypeCount="{ record }">
          <a-tag color="arcoblue">{{ record.cardTypeCount || 0 }} 个类型</a-tag>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="showDetail(record)">详情</a-button>
            <a-button type="text" size="small" @click="showEdit(record)">编辑</a-button>
            <a-button type="text" size="small" status="danger" @click="handleDelete(record)">删除</a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 详情抽屉 -->
    <a-drawer v-model:visible="detailVisible" title="卡分组详情" :width="480">
      <template v-if="currentRecord">
        <a-descriptions :column="1" bordered size="small">
          <a-descriptions-item label="分组名称">{{ currentRecord.name }}</a-descriptions-item>
          <a-descriptions-item label="分组编码">{{ currentRecord.code }}</a-descriptions-item>
          <a-descriptions-item label="卡类型数">{{ currentRecord.cardTypeCount || 0 }}</a-descriptions-item>
          <a-descriptions-item label="描述">{{ currentRecord.description || '-' }}</a-descriptions-item>
          <a-descriptions-item label="创建时间">{{ currentRecord.createdAt || '-' }}</a-descriptions-item>
        </a-descriptions>
      </template>
    </a-drawer>

    <!-- 新增/编辑弹窗 -->
    <a-modal
      v-model:visible="formVisible"
      :title="isEdit ? '编辑分组' : '新建分组'"
      @before-ok="handleFormSubmit"
      @cancel="formVisible = false"
      :width="520"
      :loading="formLoading"
    >
      <a-form :model="form" layout="vertical">
        <a-form-item label="分组名称" field="name" :rules="[{ required: true, message: '请输入分组名称' }]">
          <a-input v-model="form.name" placeholder="如：储值卡组" />
        </a-form-item>
        <a-form-item label="分组编码" field="code" :rules="[{ required: true, message: '请输入分组编码' }]">
          <a-input v-model="form.code" placeholder="如：STORAGE_CARD" :disabled="isEdit" />
        </a-form-item>
        <a-form-item label="排序值" field="sort">
          <a-input-number v-model="form.sort" :min="0" style="width: 100%" />
        </a-form-item>
        <a-form-item label="描述">
          <a-textarea v-model="form.description" :rows="3" placeholder="简要描述该分组" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { Message, Modal } from '@arco-design/web-vue'
import * as api from '@/api/member'

const dataList = ref([])
const loading = ref(false)
const formLoading = ref(false)
const formVisible = ref(false)
const detailVisible = ref(false)
const isEdit = ref(false)
const currentId = ref(null)
const currentRecord = ref(null)

const filters = reactive({ keyword: '' })
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })

const paginationConfig = computed(() => ({
  current: pagination.current, pageSize: pagination.pageSize, total: pagination.total,
  showTotal: true, showPageSize: true, pageSizeOptions: [10, 20, 50, 100]
}))

const form = reactive({ name: '', code: '', sort: 0, description: '' })

const columns = [
  { title: '分组名称', dataIndex: 'name', width: 180 },
  { title: '分组编码', dataIndex: 'code', width: 160 },
  { title: '卡类型数', slotName: 'cardTypeCount', width: 120 },
  { title: '排序', dataIndex: 'sort', width: 80 },
  { title: '描述', dataIndex: 'description', ellipsis: true },
  { title: '创建时间', dataIndex: 'createdAt', width: 170 },
  { title: '操作', slotName: 'actions', width: 150, fixed: 'right' }
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
  } catch (err) { Message.error('加载失败: ' + err.message) }
  finally { loading.value = false }
}

const handleSearch = () => { pagination.current = 1; loadData() }
const resetFilters = () => { filters.keyword = ''; pagination.current = 1; loadData() }
const onPageChange = (page) => { pagination.current = page; loadData() }
const onPageSizeChange = (pageSize) => { pagination.pageSize = pageSize; pagination.current = 1; loadData() }

const showCreate = () => {
  isEdit.value = false
  Object.assign(form, { name: '', code: '', sort: 0, description: '' })
  formVisible.value = true
}

const showEdit = (record) => {
  isEdit.value = true; currentId.value = record.id
  Object.assign(form, { name: record.name, code: record.code, sort: record.sort || 0, description: record.description || '' })
  formVisible.value = true
}

const showDetail = (record) => { currentRecord.value = record; detailVisible.value = true }

const handleFormSubmit = async (done) => {
  if (!form.name || !form.code) { Message.warning('请填写名称和编码'); done(false); return }
  formLoading.value = true
  try {
    if (isEdit.value) { await api.updateCardGroup(currentId.value, { ...form }); Message.success('更新成功') }
    else { await api.createCardGroup({ ...form }); Message.success('创建成功') }
    formVisible.value = false; loadData(); done(true)
  } catch (err) { Message.error(err.message || '操作失败'); done(false) }
  finally { formLoading.value = false }
}

const handleDelete = (record) => {
  Modal.warning({
    title: '确认删除', content: `确定要删除分组「${record.name}」吗？`,
    okText: '确认删除', onOk: async () => {
      try { await api.deleteCardGroup(record.id); Message.success('删除成功'); loadData() }
      catch (err) { Message.error(err.message || '删除失败') }
    }
  })
}

const handleExport = () => { Message.info('导出功能开发中') }

onMounted(() => loadData())
</script>

<style scoped>
.member-card-groups-page { padding: 20px 24px; min-height: calc(100vh - 64px); background: #f5f7fa; }
.breadcrumb { margin-bottom: 16px; }
.action-card { margin-bottom: 16px; }
.table-card { border-radius: 8px; }
</style>
