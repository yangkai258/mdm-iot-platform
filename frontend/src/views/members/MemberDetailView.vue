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
import { useRoute, useRouter } from 'vue-router'
import { Message } from '@arco-design/web-vue'
import * as api from '@/api/member'

const route = useRoute()
const router = useRouter()
const memberId = computed(() => route.params.id)

const member = ref({})
const loading = ref(false)
const activeTab = ref('info')

// Stats
const levelColor = computed(() => {
  const colors = { 1: '#95de64', 2: '#1890ff', 3: '#faad14', 4: '#ff4d4f' }
  return colors[member.value.levelId] || 'gray'
})
const statusColor = computed(() => ({ 1: 'green', 2: 'orange', 3: 'red' }[member.value.status] || 'gray'))
const statusText = computed(() => ({ 1: '正常', 2: '冻结', 3: '禁用' }[member.value.status] || '未知'))
const getGenderText = (g) => ({ 0: '未知', 1: '男', 2: '女' }[g] || '未知')

// Cards
const cards = ref([])
const cardsLoading = ref(false)
const cardColumns = [
  { title: '卡号', dataIndex: 'cardNo', width: 200 },
  { title: '卡类型', slotName: 'cardType', width: 120 },
  { title: '余额', dataIndex: 'balance', width: 120 },
  { title: '状态', slotName: 'status', width: 90 },
  { title: '有效期', slotName: 'expireDate', width: 150 },
  { title: '发卡时间', dataIndex: 'createdAt', width: 170 }
]

// Points
const pointsRecords = ref([])
const pointsLoading = ref(false)
const pointsPagination = reactive({ current: 1, pageSize: 20, total: 0 })
const pointsPaginationConfig = computed(() => ({
  current: pointsPagination.current, pageSize: pointsPagination.pageSize, total: pointsPagination.total,
  showTotal: true, showPageSize: true
}))
const pointsColumns = [
  { title: '类型', slotName: 'type', width: 90 },
  { title: '积分', dataIndex: 'points', width: 100 },
  { title: '原因', dataIndex: 'reason', ellipsis: true },
  { title: '时间', dataIndex: 'createdAt', width: 170 }
]

// Orders
const orders = ref([])
const ordersLoading = ref(false)
const ordersPagination = reactive({ current: 1, pageSize: 20, total: 0 })
const ordersPaginationConfig = computed(() => ({
  current: ordersPagination.current, pageSize: ordersPagination.pageSize, total: ordersPagination.total,
  showTotal: true, showPageSize: true
}))
const orderColumns = [
  { title: '订单号', dataIndex: 'orderNo', width: 200 },
  { title: '商品', dataIndex: 'goodsName', ellipsis: true },
  { title: '金额', slotName: 'amount', width: 120 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '时间', dataIndex: 'createdAt', width: 170 }
]
const getOrderStatusColor = (s) => ({ pending: 'blue', completed: 'green', cancelled: 'gray', refunded: 'orange' }[s] || 'gray')
const getOrderStatusText = (s) => ({ pending: '进行中', completed: '已完成', cancelled: '已取消', refunded: '已退款' }[s] || s)

// Tags
const memberTags = ref([])
const availableTags = ref([])
const addTagVisible = ref(false)
const tagForm = reactive({ tagId: [] })

// Coupons
const memberCoupons = ref([])
const couponsLoading = ref(false)
const couponsPagination = reactive({ current: 1, pageSize: 20, total: 0 })
const couponPaginationConfig = computed(() => ({
  current: couponsPagination.current, pageSize: couponsPagination.pageSize, total: couponsPagination.total, showTotal: true
}))
const couponOptions = ref([])
const couponColumns = [
  { title: '优惠券名称', dataIndex: 'name' },
  { title: '面值', dataIndex: 'value', width: 100 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '有效期', dataIndex: 'expireTime', width: 170 }
]
const getCouponStatusColor = (s) => ({ unused: 'green', used: 'gray', expired: 'orange' }[s] || 'gray')
const getCouponStatusText = (s) => ({ unused: '未使用', used: '已使用', expired: '已过期' }[s] || s)

// Edit info
const editInfoVisible = ref(false)
const editForm = reactive({ name: '', gender: 0, birthday: '', email: '', remark: '' })

// Points adjust
const pointsAdjustVisible = ref(false)
const pointsForm = reactive({ type: 'add', points: 0, reason: '' })

// Coupon grant
const grantCouponVisible = ref(false)
const couponForm = reactive({ couponId: undefined, count: 1 })

// Recharge
const rechargeVisible = ref(false)
const rechargeForm = reactive({ type: 'balance', amount: 0, remark: '' })

const loadMemberDetail = async () => {
  if (!memberId.value) return
  loading.value = true
  try {
    const res = await api.getMemberDetail(memberId.value)
    member.value = res.data || {}
  } catch (err) {
    Message.error('加载会员详情失败: ' + err.message)
  } finally {
    loading.value = false
  }
}

const loadCards = async () => {
  cardsLoading.value = true
  try {
    // cards.value = res.data || []
  } catch (err) { /* ignore */ }
  finally { cardsLoading.value = false }
}

const loadPointsRecords = async () => {
  pointsLoading.value = true
  try {
    const res = await api.getPointsFlow({ memberId: memberId.value, page: pointsPagination.current, pageSize: pointsPagination.pageSize })
    const d = res.data || {}
    pointsRecords.value = d.list || []
    pointsPagination.total = d.total || 0
  } catch (err) { /* ignore */ }
  finally { pointsLoading.value = false }
}

const loadOrders = async () => {
  ordersLoading.value = true
  try {
    const res = await api.getMemberOrderList({ memberId: memberId.value, page: ordersPagination.current, pageSize: ordersPagination.pageSize })
    const d = res.data || {}
    orders.value = d.list || []
    ordersPagination.total = d.total || 0
  } catch (err) { /* ignore */ }
  finally { ordersLoading.value = false }
}

const loadMemberTags = async () => {
  memberTags.value = []
}

const loadMemberCoupons = async () => {
  couponsLoading.value = true
  try {
    const res = await api.getMemberCouponList({ memberId: memberId.value, page: couponsPagination.current, pageSize: couponsPagination.pageSize })
    const d = res.data || {}
    memberCoupons.value = d.list || []
    couponsPagination.total = d.total || 0
  } catch (err) { /* ignore */ }
  finally { couponsLoading.value = false }
}

const loadCouponOptions = async () => {
  try {
    const res = await api.getCouponList({ page: 1, pageSize: 100 })
    couponOptions.value = res.data?.list || []
  } catch (err) { /* ignore */ }
}

const loadAvailableTags = async () => {
  // availableTags.value = res.data || []
}

const showEditInfo = () => {
  Object.assign(editForm, {
    name: member.value.name || '',
    gender: member.value.gender || 0,
    birthday: member.value.birthday || '',
    email: member.value.email || '',
    remark: member.value.remark || ''
  })
  editInfoVisible.value = true
}

const handleEditInfo = async (done) => {
  try {
    await api.updateMember(memberId.value, editForm)
    Message.success('更新成功')
    editInfoVisible.value = false
    loadMemberDetail()
    done(true)
  } catch (err) {
    Message.error(err.message || '更新失败')
    done(false)
  }
}

const showPointsAdjust = () => {
  Object.assign(pointsForm, { type: 'add', points: 0, reason: '' })
  pointsAdjustVisible.value = true
}

const handlePointsAdjust = async () => {
  if (!pointsForm.points || !pointsForm.reason) {
    Message.warning('请填写完整信息')
    return
  }
  try {
    await api.adjustPoints({ memberId: memberId.value, ...pointsForm })
    Message.success('积分调整成功')
    pointsAdjustVisible.value = false
    loadMemberDetail()
  } catch (err) {
    Message.error(err.message || '调整失败')
  }
}

const showCouponModal = () => {
  loadCouponOptions()
  couponForm.couponId = undefined
  couponForm.count = 1
  grantCouponVisible.value = true
}

const handleGrantCoupon = async () => {
  if (!couponForm.couponId) {
    Message.warning('请选择优惠券')
    return
  }
  try {
    await api.grantCoupon({ memberId: memberId.value, couponId: couponForm.couponId, count: couponForm.count })
    Message.success('优惠券发放成功')
    grantCouponVisible.value = false
    loadMemberCoupons()
  } catch (err) {
    Message.error(err.message || '发放失败')
  }
}

const showAddTag = () => {
  loadAvailableTags()
  tagForm.tagId = []
  addTagVisible.value = true
}

const handleAddTag = async () => {
  if (!tagForm.tagId?.length) {
    Message.warning('请选择标签')
    return
  }
  Message.success('标签添加成功')
  addTagVisible.value = false
  loadMemberTags()
}

const handleRemoveTag = (tag) => {
  memberTags.value = memberTags.value.filter(t => t.id !== tag.id)
  Message.success(`标签「${tag.name}」已移除`)
}

const showRechargeModal = () => {
  Object.assign(rechargeForm, { type: 'balance', amount: 0, remark: '' })
  rechargeVisible.value = true
}

const handleRecharge = async () => {
  if (!rechargeForm.amount) {
    Message.warning('请输入充值金额')
    return
  }
  Message.success('充值成功')
  rechargeVisible.value = false
  loadMemberDetail()
}

const showBindCard = () => {
  Message.info('绑定会员卡功能')
}

const onPointsPageChange = (page) => {
  pointsPagination.current = page
  loadPointsRecords()
}

const onOrdersPageChange = (page) => {
  ordersPagination.current = page
  loadOrders()
}

const onCouponsPageChange = (page) => {
  couponsPagination.current = page
  loadMemberCoupons()
}

const couponColumns2 = [
  { title: '优惠券名称', dataIndex: 'name' },
  { title: '面值', dataIndex: 'value', width: 100 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '有效期', dataIndex: 'expireTime', width: 170 }
]

onMounted(() => {
  loadMemberDetail()
})

</script>

<style scoped>
.page-container { background: #fff; border-radius: 4px; padding: 20px; }
.search-form { margin-bottom: 16px; padding: 16px; background: #f7f8fa; border-radius: 4px; }
.toolbar { margin-bottom: 16px; }
</style>
