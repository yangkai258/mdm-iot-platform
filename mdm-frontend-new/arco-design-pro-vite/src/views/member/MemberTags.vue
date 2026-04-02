<template>
    <Breadcrumb :items="['Home','Console','']" />


  <div class="member-tags-page">
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>会员管理</a-breadcrumb-item>
      <a-breadcrumb-item>会员标签</a-breadcrumb-item>
    </a-breadcrumb>

    <a-card class="action-card">
      <a-space wrap>
        <a-input-search v-model="filters.keyword" placeholder="搜索标签名称" style="width: 240px" search-button @search="handleSearch" />
        <a-button @click="handleSearch">「搜索」</a-button>
        <a-button @click="resetFilters">「重置」</a-button>
      </a-space>
    </a-card>

    <a-card class="table-card">
      <template #extra>
        <a-space>
          <a-button type="primary" @click="showCreate">「新建」</a-button>
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
        :scroll="{ x: 800 }"
      >
        <template #tagType="{ record }">
          <a-tag :color="getTagTypeColor(record.tagType)">{{ getTagTypeText(record.tagType) }}</a-tag>
        </template>
        <template #memberCount="{ record }">
          <a-badge :value="record.memberCount || 0" :max-count="99999" />
        </template>
        <template #status="{ record }">
          <a-switch v-model="record.status" :checked-value="1" :unchecked-value="0" @change="handleStatusChange(record)" />
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="showEdit(record)">编辑</a-button>
            <a-button type="text" size="small" status="danger" @click="handleDelete(record)">删除</a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:visible="formVisible" :title="isEdit ? '编辑标签' : '新建标签'" @before-ok="handleFormSubmit" :width="520" :loading="formLoading">
      <a-form :model="form" layout="vertical">
        <a-form-item label="标签名称" field="name" :rules="[{ required: true, message: '请输入标签名称' }]">
          <a-input v-model="form.name" placeholder="如：优质客户" />
        </a-form-item>
        <a-form-item label="标签类型" field="tagType">
          <a-select v-model="form.tagType" placeholder="选择标签类型">
            <a-option value="manual">手动标签</a-option>
            <a-option value="auto">自动标签</a-option>
            <a-option value="system">系统标签</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="描述">
          <a-textarea v-model="form.description" :rows="2" placeholder="简要描述" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { Message, Modal } from '@arco-design/web-vue'

const dataList = ref([])
const loading = ref(false)
const formLoading = ref(false)
const formVisible = ref(false)
const isEdit = ref(false)
const currentId = ref(null)

const filters = reactive({ keyword: '' })
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })

const paginationConfig = computed(() => ({
  current: pagination.current, pageSize: pagination.pageSize, total: pagination.total,
  showTotal: true, showPageSize: true, pageSizeOptions: [10, 20, 50, 100]
}))

const form = reactive({ name: '', tagType: 'manual', description: '' })

const columns = [
  { title: '标签名称', dataIndex: 'name', width: 180 },
  { title: '标签类型', slotName: 'tagType', width: 120 },
  { title: '会员数', slotName: 'memberCount', width: 100 },
  { title: '描述', dataIndex: 'description', ellipsis: true },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '操作', slotName: 'actions', width: 150 }
]

const tagTypeMap = { manual: '手动标签', auto: '自动标签', system: '系统标签' }
const tagTypeColorMap = { manual: 'blue', auto: 'green', system: 'purple' }

const getTagTypeText = (t) => tagTypeMap[t] || t
const getTagTypeColor = (t) => tagTypeColorMap[t] || 'gray'

const loadData = async () => {
  loading.value = true
  // Mock data
  setTimeout(() => {
    dataList.value = [
      { id: 1, name: '优质客户', tagType: 'auto', memberCount: 128, description: '消费金额排名前10%', status: 1 },
      { id: 2, name: '新客户', tagType: 'system', memberCount: 56, description: '注册时间30天内', status: 1 },
      { id: 3, name: '沉睡客户', tagType: 'auto', memberCount: 89, description: '90天未消费', status: 1 }
    ]
    pagination.total = 3
    loading.value = false
  }, 300)
}

const handleSearch = () => { pagination.current = 1; loadData() }
const resetFilters = () => { filters.keyword = ''; pagination.current = 1; loadData() }
const onPageChange = (page) => { pagination.current = page; loadData() }
const onPageSizeChange = (pageSize) => { pagination.pageSize = pageSize; pagination.current = 1; loadData() }

const showCreate = () => { isEdit.value = false; Object.assign(form, { name: '', tagType: 'manual', description: '' }); formVisible.value = true }
const showEdit = (record) => { isEdit.value = true; currentId.value = record.id; Object.assign(form, { name: record.name, tagType: record.tagType, description: record.description || '' }); formVisible.value = true }

const handleFormSubmit = async (done) => {
  if (!form.name) { Message.warning('请填写标签名称'); done(false); return }
  formLoading.value = true
  setTimeout(() => { Message.success(isEdit.value ? '更新成功' : '创建成功'); formVisible.value = false; loadData(); formLoading.value = false; done(true) }, 300)
}

const handleDelete = (record) => {
  Modal.warning({ title: '确认删除', content: `确定要删除标签「${record.name}」吗？`, okText: '确认删除', onOk: () => { Message.success('删除成功'); loadData() }})
}

const handleStatusChange = (record) => { Message.success(`${record.name}已${record.status === 1 ? '启用' : '禁用'}`) }

onMounted(() => loadData())
</script>

<style scoped>
.member-tags-page { padding: 20px 24px; min-height: calc(100vh - 64px); background: #f5f7fa; }
.breadcrumb { margin-bottom: 16px; }
.action-card { margin-bottom: 16px; }
.table-card { border-radius: 8px; }
</style>
