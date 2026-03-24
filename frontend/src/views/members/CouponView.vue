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
    <a-table :columns="columns" :data="couponList" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id">
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

const couponList = ref([])
const memberOptions = ref([])
const levelOptions = ref([])
const loading = ref(false)
const formLoading = ref(false)
const formVisible = ref(false)
const grantVisible = ref(false)
const detailVisible = ref(false)
const isEdit = ref(false)
const currentCoupon = ref(null)

const filters = reactive({ keyword: '', type: undefined, status: undefined })
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })
const stats = reactive({ total: 0, issued: 0, used: 0, unused: 0 })

const paginationConfig = computed(() => ({
  current: pagination.current, pageSize: pagination.pageSize, total: pagination.total,
  showTotal: true, showPageSize: true, pageSizeOptions: [10, 20, 50, 100]
}))

const form = reactive({
  name: '', type: '', value: 0, totalCount: 0,
  minConsume: 0, perLimit: 0, dateRange: [], description: ''
})

const grantForm = reactive({
  mode: 'all', memberIds: [], levelId: undefined, count: 1
})

const columns = [
  { title: '优惠券名称', dataIndex: 'name', width: 180 },
  { title: '类型', slotName: 'type', width: 100 },
  { title: '面值/折扣', slotName: 'value', width: 110 },
  { title: '最低消费', slotName: 'minConsume', width: 110 },
  { title: '发行/使用', slotName: 'usedCount', width: 110 },
  { title: '有效期至', dataIndex: 'endTime', width: 170 },
  { title: '状态', slotName: 'status', width: 90 },
  { title: '操作', slotName: 'actions', width: 220 }
]

const getTypeColor = (t) => ({ discount: 'blue', cash: 'orange', gift: 'purple', shipping: 'green' }[t] || 'gray')
const getTypeText = (t) => ({ discount: '折扣券', cash: '现金券', gift: '礼品券', shipping: '包邮券' }[t] || t)
const getStatusColor = (s) => ({ active: 'green', inactive: 'gray', expired: 'orange' }[s] || 'gray')
const getStatusText = (s) => ({ active: '有效', inactive: '未激活', expired: '已过期' }[s] || s)

const loadCoupons = async () => {
  loading.value = true
  try {
    const params = { page: pagination.current, pageSize: pagination.pageSize }
    if (filters.keyword) params.keyword = filters.keyword
    if (filters.type) params.type = filters.type
    if (filters.status) params.status = filters.status
    const res = await api.getCouponList(params)
    const d = res.data || {}
    couponList.value = d.list || []
    pagination.total = d.total || 0
    stats.total = d.total || 0
  } catch (err) {
    Message.error('加载优惠券列表失败: ' + err.message)
  } finally {
    loading.value = false
  }
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

const showCreateDrawer = () => {
  isEdit.value = false
  Object.assign(form, { name: '', type: '', value: 0, totalCount: 0, minConsume: 0, perLimit: 0, dateRange: [], description: '' })
  formVisible.value = true
}

const showEdit = (record) => {
  isEdit.value = true
  Object.assign(form, {
    name: record.name, type: record.type, value: record.value || 0,
    totalCount: record.totalCount || 0, minConsume: record.minConsume || 0,
    perLimit: record.perLimit || 0, dateRange: [], description: record.description || ''
  })
  formVisible.value = true
}

const showDetail = (record) => {
  currentCoupon.value = record
  detailVisible.value = true
}

const showGrant = (record) => {
  currentCoupon.value = record
  loadMemberOptions()
  loadLevelOptions()
  Object.assign(grantForm, { mode: 'all', memberIds: [], levelId: undefined, count: 1 })
  grantVisible.value = true
}

const handleFormSubmit = async () => {
  if (!form.name || !form.type) { Message.warning('请填写名称和类型'); return }
  formLoading.value = true
  try {
    if (isEdit.value) {
      await api.updateCoupon(currentCoupon.value.id, { ...form })
      Message.success('更新成功')
    } else {
      await api.createCoupon({ ...form })
      Message.success('创建成功')
    }
    formVisible.value = false
    loadCoupons()
  } catch (err) {
    Message.error(err.message || '操作失败')
  } finally {
    formLoading.value = false
  }
}

const handleDelete = async (record) => {
  try {
    await api.deleteCoupon(record.id)
    Message.success('删除成功')
    loadCoupons()
  } catch (err) {
    Message.error(err.message || '删除失败')
  }
}

const handleGrant = async () => {
  try {
    const payload = { couponId: currentCoupon.value.id, count: grantForm.count }
    if (grantForm.mode === 'member') payload.memberIds = grantForm.memberIds
    if (grantForm.mode === 'level') payload.levelId = grantForm.levelId
    if (grantForm.mode === 'all') payload.grantToAll = true
    await api.grantCoupon(payload)
    Message.success('发放成功')
    grantVisible.value = false
  } catch (err) {
    Message.error(err.message || '发放失败')
  }
}

const onPageChange = (page) => { pagination.current = page; loadCoupons() }
const onPageSizeChange = (pageSize) => { pagination.pageSize = pageSize; pagination.current = 1; loadCoupons() }

onMounted(() => loadCoupons())

</script>

<style scoped>
.page-container { background: #fff; border-radius: 4px; padding: 20px; }
.search-form { margin-bottom: 16px; padding: 16px; background: #f7f8fa; border-radius: 4px; }
.toolbar { margin-bottom: 16px; }
</style>
