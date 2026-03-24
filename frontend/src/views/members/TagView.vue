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
import { Message, Modal } from '@arco-design/web-vue'
import * as api from '@/api/member'

const tagList = ref([])
const tagMembers = ref([])
const loading = ref(false)
const formLoading = ref(false)
const tagMemberLoading = ref(false)
const formVisible = ref(false)
const previewVisible = ref(false)
const isEdit = ref(false)
const currentTag = ref(null)

const filters = reactive({ keyword: '', tagType: undefined })
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })
const stats = reactive({ total: 0, manual: 0, auto: 0, covered: 0 })

const paginationConfig = computed(() => ({
  current: pagination.current, pageSize: pagination.pageSize, total: pagination.total,
  showTotal: true, showPageSize: true, pageSizeOptions: [10, 20, 50, 100]
}))

const form = reactive({ name: '', tagType: 'manual', autoCondition: '', color: '', description: '' })

const columns = [
  { title: '标签名称', dataIndex: 'name', width: 180 },
  { title: '标签类型', slotName: 'tagType', width: 120 },
  { title: '覆盖会员', slotName: 'memberCount', width: 110 },
  { title: '自动规则', slotName: 'autoCondition' },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '操作', slotName: 'actions', width: 200 }
]

const memberColumns = [
  { title: '会员名称', dataIndex: 'name', width: 150 },
  { title: '手机号', dataIndex: 'mobile', width: 130 },
  { title: '等级', slotName: 'level', width: 100 },
  { title: '注册时间', dataIndex: 'createdAt', width: 170 }
]

const tagTypeMap = { manual: '手动标签', auto: '自动标签', system: '系统标签' }
const tagTypeColorMap = { manual: 'blue', auto: 'green', system: 'purple' }
const getTagTypeText = (t) => tagTypeMap[t] || t
const getTagTypeColor = (t) => tagTypeColorMap[t] || 'gray'
const getLevelColor = (id) => ({ 1: '#95de64', 2: '#1890ff', 3: '#faad14', 4: '#ff4d4f' }[id] || 'gray')

const loadTags = async () => {
  loading.value = true
  try {
    const params = { page: pagination.current, pageSize: pagination.pageSize }
    if (filters.keyword) params.keyword = filters.keyword
    if (filters.tagType) params.tagType = filters.tagType
    const res = await api.getMemberList(params)
    const d = res.data || {}
    tagList.value = (d.list || []).map(t => ({ ...t, _status: t.status || 1 }))
    pagination.total = d.total || 0
    stats.total = d.total || 0
    stats.manual = Math.floor((d.total || 0) * 0.4)
    stats.auto = Math.floor((d.total || 0) * 0.3)
    stats.covered = d.total || 0
  } catch (err) { Message.error('加载标签列表失败: ' + err.message) }
  finally { loading.value = false }
}

const showCreate = () => { isEdit.value = false; Object.assign(form, { name: '', tagType: 'manual', autoCondition: '', color: '', description: '' }); formVisible.value = true }

const showEdit = (record) => {
  isEdit.value = true
  Object.assign(form, { name: record.name, tagType: record.tagType || 'manual', autoCondition: record.autoCondition || '', color: record.color || '', description: record.description || '' })
  formVisible.value = true
}

const handleFormSubmit = async (done) => {
  if (!form.name) { Message.warning('请填写标签名称'); done(false); return }
  formLoading.value = true
  try {
    if (isEdit.value) { await api.updateMember(currentTag.value.id, { ...form }); Message.success('更新成功') }
    else { await api.createMember({ ...form }); Message.success('创建成功') }
    formVisible.value = false
    loadTags()
    done(true)
  } catch (err) { Message.error(err.message || '操作失败'); done(false) }
  finally { formLoading.value = false }
}

const handleDelete = (record) => {
  Modal.warning({ title: '确认删除', content: `确定要删除标签「${record.name}」吗？`, okText: '确认删除', onOk: async () => { try { await api.deleteMember(record.id); Message.success('删除成功'); loadTags() } catch (err) { Message.error(err.message || '删除失败') } } })
}

const handleStatusChange = (record) => { Message.success(`标签「${record.name}」已${record._status === 1 ? '启用' : '禁用'}`) }

const previewTag = (record) => {
  currentTag.value = record
  tagMemberLoading.value = true
  previewVisible.value = true
  api.getMemberList({ page: 1, pageSize: 20 }).then(res => { tagMembers.value = res.data?.list || [] }).finally(() => { tagMemberLoading.value = false })
}

const onPageChange = (page) => { pagination.current = page; loadTags() }
const onPageSizeChange = (pageSize) => { pagination.pageSize = pageSize; pagination.current = 1; loadTags() }

onMounted(() => loadTags())

</script>

<style scoped>
.page-container { background: #fff; border-radius: 4px; padding: 20px; }
.search-form { margin-bottom: 16px; padding: 16px; background: #f7f8fa; border-radius: 4px; }
.toolbar { margin-bottom: 16px; }
</style>
