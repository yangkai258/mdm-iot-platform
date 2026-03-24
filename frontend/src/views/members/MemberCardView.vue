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

const cardTypeList = ref([])
const issueRecords = ref([])
const memberOptions = ref([])
const loading = ref(false)
const issueLoading = ref(false)
const formLoading = ref(false)
const typeModalVisible = ref(false)
const issueModalVisible = ref(false)
const batchIssueVisible = ref(false)
const isEditType = ref(false)
const currentId = ref(null)
const currentCardType = ref(null)

const filters = reactive({ keyword: '', status: undefined })
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })
const issuePagination = reactive({ current: 1, pageSize: 20, total: 0 })
const stats = reactive({ totalTypes: 0, issuedCards: 0, activeCards: 0, todayIssued: 0 })

const paginationConfig = computed(() => ({
  current: pagination.current, pageSize: pagination.pageSize, total: pagination.total,
  showTotal: true, showPageSize: true, pageSizeOptions: [10, 20, 50, 100]
}))
const issuePaginationConfig = computed(() => ({
  current: issuePagination.current, pageSize: issuePagination.pageSize, total: issuePagination.total,
  showTotal: true, showPageSize: true
}))

const typeForm = reactive({
  name: '', code: '', discountRate: 1, pointsRate: 1, validDays: 0,
  cardFee: 0, color: '#1890ff', benefits: ''
})
const typeFormStatus = ref('1')
const issueForm = reactive({ memberId: undefined, remark: '' })
const batchForm = reactive({ cardTypeId: undefined, count: 1, remark: '' })

const typeColumns = [
  { title: '卡类型名称', dataIndex: 'name', width: 180 },
  { title: '编码', dataIndex: 'code', width: 120 },
  { title: '折扣率', slotName: 'discountRate', width: 100 },
  { title: '积分倍率', slotName: 'pointsRate', width: 110 },
  { title: '有效期', slotName: 'validDays', width: 100 },
  { title: '工本费', slotName: 'cardFee', width: 100 },
  { title: '卡面颜色', dataIndex: 'color', width: 100 },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '操作', slotName: 'actions', width: 180 }
]

const issueColumns = [
  { title: '卡号', dataIndex: 'cardNo', width: 200 },
  { title: '卡类型', slotName: 'cardType', width: 150 },
  { title: '会员', dataIndex: 'memberName', width: 120 },
  { title: '手机号', dataIndex: 'mobile', width: 130 },
  { title: '状态', slotName: 'status', width: 90 },
  { title: '发行时间', dataIndex: 'createdAt', width: 170 }
]

const loadCardTypes = async () => {
  loading.value = true
  try {
    const params = { page: pagination.current, pageSize: pagination.pageSize }
    if (filters.keyword) params.keyword = filters.keyword
    if (filters.status !== undefined) params.status = filters.status
    const res = await api.getCardTypeList(params)
    const d = res.data || {}
    cardTypeList.value = d.list || []
    pagination.total = d.total || 0
    stats.totalTypes = d.total || 0
  } catch (err) {
    Message.error('加载卡类型失败: ' + err.message)
  } finally {
    loading.value = false
  }
}

const loadIssueRecords = async () => {
  issueLoading.value = true
  try {
    const params = { page: issuePagination.current, pageSize: issuePagination.pageSize }
    const res = await api.getCardGroupList(params)
    const d = res.data || {}
    issueRecords.value = d.list || []
    issuePagination.total = d.total || 0
  } catch (err) {
    issueRecords.value = []
  } finally {
    issueLoading.value = false
  }
}

const loadMemberOptions = async () => {
  try {
    const res = await api.getMemberList({ page: 1, pageSize: 100 })
    memberOptions.value = res.data?.list || []
  } catch (err) { /* ignore */ }
}

const showCreateType = () => {
  isEditType.value = false
  currentId.value = null
  Object.assign(typeForm, { name: '', code: '', discountRate: 1, pointsRate: 1, validDays: 0, cardFee: 0, color: '#1890ff', benefits: '' })
  typeFormStatus.value = '1'
  typeModalVisible.value = true
}

const showEditType = (record) => {
  isEditType.value = true
  currentId.value = record.id
  Object.assign(typeForm, {
    name: record.name,
    code: record.code,
    discountRate: record.discountRate || 1,
    pointsRate: record.pointsRate || 1,
    validDays: record.validDays || 0,
    cardFee: record.cardFee || 0,
    color: record.color || '#1890ff',
    benefits: record.benefits || ''
  })
  typeFormStatus.value = String(record.status || 1)
  typeModalVisible.value = true
}

const handleTypeSubmit = async (done) => {
  if (!typeForm.name || !typeForm.code) {
    Message.warning('请填写名称和编码')
    done(false)
    return
  }
  formLoading.value = true
  try {
    const payload = { ...typeForm, status: parseInt(typeFormStatus.value) }
    if (isEditType.value) {
      await api.updateCardType(currentId.value, payload)
      Message.success('更新成功')
    } else {
      await api.createCardType(payload)
      Message.success('创建成功')
    }
    typeModalVisible.value = false
    loadCardTypes()
    done(true)
  } catch (err) {
    Message.error(err.message || '操作失败')
    done(false)
  } finally {
    formLoading.value = false
  }
}

const handleDeleteType = (record) => {
  Modal.warning({
    title: '确认删除',
    content: `确定要删除卡类型「${record.name}」吗？`,
    okText: '确认删除',
    onOk: async () => {
      try {
        await api.deleteCardType(record.id)
        Message.success('删除成功')
        loadCardTypes()
      } catch (err) {
        Message.error(err.message || '删除失败')
      }
    }
  })
}

const showIssueCard = (record) => {
  currentCardType.value = record
  issueForm.memberId = undefined
  issueForm.remark = ''
  issueModalVisible.value = true
}

const handleIssueCard = async () => {
  if (!issueForm.memberId) {
    Message.warning('请选择会员')
    return
  }
  Message.success('发行成功')
  issueModalVisible.value = false
  loadIssueRecords()
}

const showBatchIssue = () => {
  batchForm.cardTypeId = undefined
  batchForm.count = 1
  batchForm.remark = ''
  batchIssueVisible.value = true
}

const handleBatchIssue = async () => {
  if (!batchForm.cardTypeId) {
    Message.warning('请选择卡类型')
    return
  }
  Message.success('批量发行成功')
  batchIssueVisible.value = false
  loadIssueRecords()
}

const onPageChange = (page) => { pagination.current = page; loadCardTypes() }
const onIssuePageChange = (page) => { issuePagination.current = page; loadIssueRecords() }

onMounted(() => {
  loadCardTypes()
  loadIssueRecords()
})

</script>

<style scoped>
.page-container { background: #fff; border-radius: 4px; padding: 20px; }
.search-form { margin-bottom: 16px; padding: 16px; background: #f7f8fa; border-radius: 4px; }
.toolbar { margin-bottom: 16px; }
</style>
