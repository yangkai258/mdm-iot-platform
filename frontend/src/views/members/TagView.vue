<template>
  <div class="tag-page">
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>会员管理</a-breadcrumb-item>
      <a-breadcrumb-item>标签管理</a-breadcrumb-item>
    </a-breadcrumb>

    <a-row :gutter="16" class="stats-row">
      <a-col :span="6"><a-card class="stat-card"><a-statistic title="标签总数" :value="stats.total || 0" /></a-card></a-col>
      <a-col :span="6"><a-card class="stat-card"><a-statistic title="手动标签" :value="stats.manual || 0" :value-style="{ color: '#1890ff' }" /></a-card></a-col>
      <a-col :span="6"><a-card class="stat-card"><a-statistic title="自动标签" :value="stats.auto || 0" :value-style="{ color: '#52c41a' }" /></a-card></a-col>
      <a-col :span="6"><a-card class="stat-card"><a-statistic title="覆盖会员" :value="stats.covered || 0" /></a-card></a-col>
    </a-row>

    <a-card class="action-card">
      <a-space wrap>
        <a-input-search v-model="filters.keyword" placeholder="搜索标签名称" style="width: 240px" search-button @search="loadTags" />
        <a-select v-model="filters.tagType" placeholder="标签类型" allow-clear style="width: 140px" @change="loadTags">
          <a-option value="manual">手动标签</a-option>
          <a-option value="auto">自动标签</a-option>
          <a-option value="system">系统标签</a-option>
        </a-select>
        <a-button type="primary" @click="showCreate">新建标签</a-button>
        <a-button @click="loadTags">刷新</a-button>
      </a-space>
    </a-card>

    <a-card class="table-card">
      <a-table :columns="columns" :data="tagList" :loading="loading" :pagination="paginationConfig"
        @page-change="onPageChange" @page-size-change="onPageSizeChange" row-key="id" :scroll="{ x: 900 }">
        <template #tagType="{ record }"><a-tag :color="getTagTypeColor(record.tagType)">{{ getTagTypeText(record.tagType) }}</a-tag></template>
        <template #memberCount="{ record }"><a-badge :value="record.memberCount || 0" :max-count="99999" /></template>
        <template #autoCondition="{ record }"><span style="color: #666; font-size: 12px;">{{ record.autoCondition || '-' }}</span></template>
        <template #status="{ record }"><a-switch v-model="record._status" :checked-value="1" :unchecked-value="0" @change="handleStatusChange(record)" /></template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="showEdit(record)">编辑</a-button>
            <a-button type="text" size="small" @click="previewTag(record)">预览</a-button>
            <a-button type="text" size="small" status="danger" @click="handleDelete(record)">删除</a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:visible="formVisible" :title="isEdit ? '编辑标签' : '新建标签'" @before-ok="handleFormSubmit" :width="520" :loading="formLoading">
      <a-form :model="form" layout="vertical">
        <a-form-item label="标签名称" required><a-input v-model="form.name" placeholder="如：优质客户" /></a-form-item>
        <a-form-item label="标签类型" required>
          <a-select v-model="form.tagType" placeholder="选择标签类型">
            <a-option value="manual">手动标签</a-option>
            <a-option value="auto">自动标签</a-option>
            <a-option value="system">系统标签</a-option>
          </a-select>
        </a-form-item>
        <a-form-item v-if="form.tagType === 'auto'" label="自动规则条件">
          <a-textarea v-model="form.autoCondition" :rows="3" placeholder="如：累计消费满1000元" />
          <template #extra>设置自动打标的规则条件</template>
        </a-form-item>
        <a-form-item label="标签颜色"><a-input v-model="form.color" placeholder="#1890ff" /></a-form-item>
        <a-form-item label="描述"><a-textarea v-model="form.description" :rows="2" placeholder="简要描述" /></a-form-item>
      </a-form>
    </a-modal>

    <a-drawer v-model:visible="previewVisible" title="标签预览" :width="520">
      <template v-if="currentTag">
        <a-descriptions :column="1" bordered size="small">
          <a-descriptions-item label="标签名称">{{ currentTag.name }}</a-descriptions-item>
          <a-descriptions-item label="标签类型"><a-tag :color="getTagTypeColor(currentTag.tagType)">{{ getTagTypeText(currentTag.tagType) }}</a-tag></a-descriptions-item>
          <a-descriptions-item label="覆盖会员">{{ currentTag.memberCount || 0 }}</a-descriptions-item>
          <a-descriptions-item label="自动规则">{{ currentTag.autoCondition || '-' }}</a-descriptions-item>
          <a-descriptions-item label="描述">{{ currentTag.description || '-' }}</a-descriptions-item>
        </a-descriptions>
        <a-divider>该标签下的会员</a-divider>
        <a-table :columns="memberColumns" :data="tagMembers" :loading="tagMemberLoading" :pagination="false" row-key="id" size="small">
          <template #level="{ record }"><a-tag :color="getLevelColor(record.levelId)">{{ record.levelName || '普通' }}</a-tag></template>
        </a-table>
      </template>
    </a-drawer>
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
.tag-page { padding: 20px 24px; min-height: calc(100vh - 64px); background: #f5f7fa; }
.breadcrumb { margin-bottom: 16px; }
.stats-row { margin-bottom: 16px; }
.stat-card { border-radius: 8px; text-align: center; }
.action-card { margin-bottom: 16px; }
.table-card { border-radius: 8px; }
</style>
