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
      <template #status="{ record }">
        <a-tag :color="record.status === 1 ? 'green' : 'gray'">{{ record.status === 1 ? '启用' : '禁用' }}</a-tag>
      </template>
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
const couponOptions = ref([])
const memberOptions = ref([])
const levelOptions = ref([])
const loading = ref(false)
const formLoading = ref(false)
const grantVisible = ref(false)

const filters = reactive({ keyword: '', mode: undefined, dateRange: [] })
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })

const paginationConfig = computed(() => ({
  current: pagination.current, pageSize: pagination.pageSize, total: pagination.total,
  showTotal: true, showPageSize: true, pageSizeOptions: [10, 20, 50, 100]
}))

const grantForm = reactive({
  couponId: undefined, mode: 'all', memberIds: [], levelId: undefined, count: 1
})

const columns = [
  { title: '优惠券名称', dataIndex: 'couponName', width: 180 },
  { title: '发放方式', slotName: 'mode', width: 120 },
  { title: '发放对象', slotName: 'target', width: 160 },
  { title: '发放数量', dataIndex: 'count', width: 100 },
  { title: '发放时间', dataIndex: 'createTime', width: 170 },
  { title: '状态', slotName: 'status', width: 90 }
]

const getModeColor = (m) => ({ member: 'blue', level: 'purple', all: 'green' }[m] || 'gray')
const getModeText = (m) => ({ member: '指定会员', level: '指定等级', all: '全部会员' }[m] || m)
const getStatusColor = (s) => ({ success: 'green', pending: 'blue', failed: 'red' }[s] || 'gray')
const getStatusText = (s) => ({ success: '成功', pending: '进行中', failed: '失败' }[s] || s)

const loadData = async () => {
  loading.value = true
  try {
    const params = { page: pagination.current, pageSize: pagination.pageSize }
    if (filters.keyword) params.keyword = filters.keyword
    if (filters.mode) params.mode = filters.mode
    if (filters.dateRange && filters.dateRange.length === 2) {
      params.startDate = filters.dateRange[0]
      params.endDate = filters.dateRange[1]
    }
    const res = await api.getCouponGrantList(params)
    const d = res.data || {}
    dataList.value = d.list || []
    pagination.total = d.total || 0
  } catch (err) {
    Message.error('加载发放记录失败: ' + err.message)
  } finally {
    loading.value = false
  }
}

const loadCouponOptions = async () => {
  try {
    const res = await api.getCouponList({ page: 1, pageSize: 100 })
    couponOptions.value = res.data?.list || []
  } catch (err) { /* ignore */ }
}

const loadMemberOptions = async () => {
  try {
    const res = await api.getMemberList({ page: 1, pageSize: 100 })
    memberOptions.value = res.data?.list || []
  } catch (err) { /* ignore */ }
}

const loadLevelOptions = async () => {
  try {
    const res = await api.getLevelList()
    levelOptions.value = res.data || []
  } catch (err) { /* ignore */ }
}

const showGrantDrawer = () => {
  loadCouponOptions()
  loadMemberOptions()
  loadLevelOptions()
  Object.assign(grantForm, { couponId: undefined, mode: 'all', memberIds: [], levelId: undefined, count: 1 })
  grantVisible.value = true
}

const handleGrant = async () => {
  if (!grantForm.couponId) { Message.warning('请选择优惠券'); return }
  formLoading.value = true
  try {
    const payload = { couponId: grantForm.couponId, count: grantForm.count }
    if (grantForm.mode === 'member') payload.memberIds = grantForm.memberIds
    if (grantForm.mode === 'level') payload.levelId = grantForm.levelId
    if (grantForm.mode === 'all') payload.grantToAll = true
    await api.grantCoupon(payload)
    Message.success('发放成功')
    grantVisible.value = false
    loadData()
  } catch (err) {
    Message.error(err.message || '发放失败')
  } finally {
    formLoading.value = false
  }
}

const handleExport = () => {
  Message.info('导出功能开发中')
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
