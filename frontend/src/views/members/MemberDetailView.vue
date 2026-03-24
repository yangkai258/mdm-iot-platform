<template>
  <div class="member-detail-page">
    <!-- 面包屑 -->
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item><a @click="$router.push('/members')">会员管理</a></a-breadcrumb-item>
      <a-breadcrumb-item>会员详情</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 基本信息卡片 -->
    <a-row :gutter="16" class="info-row">
      <a-col :span="16">
        <a-card class="profile-card">
          <a-row :gutter="24">
            <a-col :span="4" style="text-align:center">
              <a-avatar :style="{ backgroundColor: levelColor, width: '72px', height: '72px', fontSize: '28px' }">
                {{ (member.name || member.mobile || '?').charAt(0) }}
              </a-avatar>
            </a-col>
            <a-col :span="20">
              <div class="profile-header">
                <div class="profile-name">
                  <span style="font-size: 20px; font-weight: 600;">{{ member.name || '-' }}</span>
                  <a-tag :color="levelColor" style="margin-left: 8px">{{ member.levelName || '普通会员' }}</a-tag>
                  <a-tag :color="statusColor">{{ statusText }}</a-tag>
                </div>
                <div class="profile-meta">
                  <span><icon-mobile /> {{ member.mobile || '-' }}</span>
                  <span style="margin-left: 16px"><icon-user /> 会员编号: {{ member.memberNo || '-' }}</span>
                  <span style="margin-left: 16px"><icon-clock /> 注册: {{ member.createdAt || '-' }}</span>
                </div>
              </div>
            </a-col>
          </a-row>
        </a-card>
      </a-col>
      <a-col :span="8">
        <a-card class="quick-action-card">
          <a-space direction="vertical" style="width:100%">
            <a-button type="primary" long @click="showPointsAdjust">积分调整</a-button>
            <a-button long @click="showRechargeModal">充值</a-button>
            <a-button type="outline" long @click="showCouponModal">发放优惠券</a-button>
          </a-space>
        </a-card>
      </a-col>
    </a-row>

    <!-- 标签 + 统计 -->
    <a-row :gutter="16" class="stats-row">
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="可用积分" :value="member.totalPoints || 0" :value-style="{ color: '#ff6b00' }">
            <template #prefix><icon-star /></template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="累计消费" :value="member.totalConsume || 0" prefix="¥" :precision="2" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="订单数" :value="member.totalOrderCount || 0" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="优惠券" :value="member.couponCount || 0" />
        </a-card>
      </a-col>
    </a-row>

    <!-- Tab 区域 -->
    <a-card class="tab-card">
      <a-tabs v-model:active-key="activeTab">
        <!-- 基本信息 -->
        <a-tab-pane key="info" title="基本信息">
          <a-descriptions :column="2" bordered size="small">
            <a-descriptions-item label="姓名">{{ member.name || '-' }}</a-descriptions-item>
            <a-descriptions-item label="手机号">{{ member.mobile || '-' }}</a-descriptions-item>
            <a-descriptions-item label="性别">{{ getGenderText(member.gender) }}</a-descriptions-item>
            <a-descriptions-item label="生日">{{ member.birthday || '-' }}</a-descriptions-item>
            <a-descriptions-item label="邮箱">{{ member.email || '-' }}</a-descriptions-item>
            <a-descriptions-item label="会员等级">
              <a-tag :color="levelColor">{{ member.levelName || '普通' }}</a-tag>
            </a-descriptions-item>
            <a-descriptions-item label="所属门店">{{ member.storeName || '-' }}</a-descriptions-item>
            <a-descriptions-item label="状态">
              <a-tag :color="statusColor">{{ statusText }}</a-tag>
            </a-descriptions-item>
            <a-descriptions-item label="注册时间" :span="2">{{ member.createdAt || '-' }}</a-descriptions-item>
            <a-descriptions-item label="备注" :span="2">{{ member.remark || '-' }}</a-descriptions-item>
          </a-descriptions>
          <div style="margin-top: 16px">
            <a-button type="primary" @click="showEditInfo">编辑信息</a-button>
          </div>
        </a-tab-pane>

        <!-- 会员卡信息 -->
        <a-tab-pane key="card" title="会员卡信息">
          <a-space style="margin-bottom: 12px">
            <a-button type="primary" size="small" @click="showBindCard">绑定会员卡</a-button>
          </a-space>
          <a-table :columns="cardColumns" :data="cards" :loading="cardsLoading" :pagination="false" row-key="id">
            <template #cardType="{ record }">
              <a-tag>{{ record.cardTypeName || '-' }}</a-tag>
            </template>
            <template #status="{ record }">
              <a-tag :color="record.status === 1 ? 'green' : 'gray'">{{ record.status === 1 ? '正常' : '失效' }}</a-tag>
            </template>
            <template #expireDate="{ record }">
              {{ record.expireDate || '永久' }}
            </template>
          </a-table>
          <a-empty v-if="!cards.length && !cardsLoading" description="暂无会员卡" style="margin-top: 32px" />
        </a-tab-pane>

        <!-- 积分记录 -->
        <a-tab-pane key="points" title="积分记录">
          <a-space style="margin-bottom: 12px">
            <a-button type="primary" size="small" @click="showPointsAdjust">调整积分</a-button>
          </a-space>
          <a-table :columns="pointsColumns" :data="pointsRecords" :loading="pointsLoading"
            :pagination="pointsPaginationConfig" @page-change="onPointsPageChange" row-key="id">
            <template #type="{ record }">
              <a-tag :color="record.type === 'add' ? 'green' : 'red'">
                {{ record.type === 'add' ? '获得' : '消耗' }}
              </a-tag>
            </template>
          </a-table>
          <a-empty v-if="!pointsRecords.length && !pointsLoading" description="暂无积分记录" style="margin-top: 32px" />
        </a-tab-pane>

        <!-- 消费记录 -->
        <a-tab-pane key="orders" title="消费记录">
          <a-table :columns="orderColumns" :data="orders" :loading="ordersLoading"
            :pagination="ordersPaginationConfig" @page-change="onOrdersPageChange" row-key="id">
            <template #status="{ record }">
              <a-tag :color="getOrderStatusColor(record.status)">{{ getOrderStatusText(record.status) }}</a-tag>
            </template>
            <template #amount="{ record }">
              <span style="color: #ff6b00; font-weight: 600;">¥{{ (record.amount || 0).toFixed(2) }}</span>
            </template>
          </a-table>
          <a-empty v-if="!orders.length && !ordersLoading" description="暂无消费记录" style="margin-top: 32px" />
        </a-tab-pane>

        <!-- 标签 -->
        <a-tab-pane key="tags" title="会员标签">
          <a-space style="margin-bottom: 12px">
            <a-button type="primary" size="small" @click="showAddTag">添加标签</a-button>
          </a-space>
          <a-space wrap>
            <a-tag v-for="tag in memberTags" :key="tag.id" closable @close="handleRemoveTag(tag)">
              {{ tag.name }}
            </a-tag>
            <a-tag v-if="!memberTags.length" style="color: #999">暂无标签</a-tag>
          </a-space>
        </a-tab-pane>

        <!-- 优惠券 -->
        <a-tab-pane key="coupons" title="优惠券">
          <a-space style="margin-bottom: 12px">
            <a-button type="primary" size="small" @click="showCouponModal">发放优惠券</a-button>
          </a-space>
          <a-table :columns="couponColumns" :data="memberCoupons" :loading="couponsLoading"
            :pagination="couponPaginationConfig" @page-change="onCouponsPageChange" row-key="id">
            <template #status="{ record }">
              <a-tag :color="getCouponStatusColor(record.status)">{{ getCouponStatusText(record.status) }}</a-tag>
            </template>
          </a-table>
          <a-empty v-if="!memberCoupons.length && !couponsLoading" description="暂无优惠券" style="margin-top: 32px" />
        </a-tab-pane>
      </a-tabs>
    </a-card>

    <!-- 编辑信息弹窗 -->
    <a-modal v-model:visible="editInfoVisible" title="编辑会员信息" @before-ok="handleEditInfo" :width="520">
      <a-form :model="editForm" layout="vertical">
        <a-form-item label="姓名" required><a-input v-model="editForm.name" placeholder="请输入姓名" /></a-form-item>
        <a-form-item label="性别">
          <a-radio-group v-model="editForm.gender">
            <a-radio :value="0">未知</a-radio><a-radio :value="1">男</a-radio><a-radio :value="2">女</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item label="生日"><a-date-picker v-model="editForm.birthday" format="YYYY-MM-DD" style="width:100%" /></a-form-item>
        <a-form-item label="邮箱"><a-input v-model="editForm.email" placeholder="请输入邮箱" /></a-form-item>
        <a-form-item label="备注"><a-textarea v-model="editForm.remark" :rows="2" /></a-form-item>
      </a-form>
    </a-modal>

    <!-- 积分调整抽屉 -->
    <a-drawer v-model:visible="pointsAdjustVisible" title="积分调整" :width="420">
      <a-form :model="pointsForm" layout="vertical">
        <a-form-item label="会员"><a-input :value="member.name + ' (' + member.mobile + ')'" disabled /></a-form-item>
        <a-form-item label="调整类型" required>
          <a-radio-group v-model="pointsForm.type">
            <a-radio value="add">增加</a-radio><a-radio value="deduct">扣除</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item label="积分数量" required>
          <a-input-number v-model="pointsForm.points" :min="1" :max="1000000" style="width:100%" />
        </a-form-item>
        <a-form-item label="原因" required>
          <a-textarea v-model="pointsForm.reason" :rows="3" placeholder="请输入调整原因" />
        </a-form-item>
      </a-form>
      <template #footer>
        <a-button @click="pointsAdjustVisible = false">取消</a-button>
        <a-button type="primary" :loading="pointsLoading" @click="handlePointsAdjust">确认调整</a-button>
      </template>
    </a-drawer>

    <!-- 发放优惠券抽屉 -->
    <a-drawer v-model:visible="grantCouponVisible" title="发放优惠券" :width="420">
      <a-form :model="couponForm" layout="vertical">
        <a-form-item label="会员"><a-input :value="member.name + ' (' + member.mobile + ')'" disabled /></a-form-item>
        <a-form-item label="选择优惠券" required>
          <a-select v-model="couponForm.couponId" placeholder="请选择优惠券">
            <a-option v-for="c in couponOptions" :key="c.id" :value="c.id">{{ c.name }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="发放数量">
          <a-input-number v-model="couponForm.count" :min="1" :max="10" style="width:100%" />
        </a-form-item>
      </a-form>
      <template #footer>
        <a-button @click="grantCouponVisible = false">取消</a-button>
        <a-button type="primary" @click="handleGrantCoupon">确认发放</a-button>
      </template>
    </a-drawer>

    <!-- 添加标签弹窗 -->
    <a-modal v-model:visible="addTagVisible" title="添加标签" :width="400">
      <a-form :model="tagForm" layout="vertical">
        <a-form-item label="选择标签" required>
          <a-select v-model="tagForm.tagId" placeholder="请选择标签" multiple>
            <a-option v-for="t in availableTags" :key="t.id" :value="t.id">{{ t.name }}</a-option>
          </a-select>
        </a-form-item>
      </a-form>
      <template #footer>
        <a-button @click="addTagVisible = false">取消</a-button>
        <a-button type="primary" @click="handleAddTag">确认添加</a-button>
      </template>
    </a-modal>

    <!-- 充值弹窗 -->
    <a-modal v-model:visible="rechargeVisible" title="账户充值" :width="400">
      <a-form :model="rechargeForm" layout="vertical">
        <a-form-item label="会员"><a-input :value="member.name + ' (' + member.mobile + ')'" disabled /></a-form-item>
        <a-form-item label="充值类型" required>
          <a-select v-model="rechargeForm.type" placeholder="选择充值类型">
            <a-option value="balance">余额</a-option>
            <a-option value="points">积分</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="金额" required>
          <a-input-number v-model="rechargeForm.amount" :min="0.01" style="width:100%" />
        </a-form-item>
        <a-form-item label="备注"><a-textarea v-model="rechargeForm.remark" :rows="2" /></a-form-item>
      </a-form>
      <template #footer>
        <a-button @click="rechargeVisible = false">取消</a-button>
        <a-button type="primary" @click="handleRecharge">确认充值</a-button>
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
.member-detail-page { padding: 20px 24px; min-height: calc(100vh - 64px); background: #f5f7fa; }
.breadcrumb { margin-bottom: 16px; }
.info-row { margin-bottom: 16px; }
.profile-card { border-radius: 8px; }
.profile-header { padding: 8px 0; }
.profile-name { margin-bottom: 8px; }
.profile-meta { color: #666; font-size: 14px; }
.profile-meta span { display: inline-flex; align-items: center; gap: 4px; }
.quick-action-card { border-radius: 8px; height: 100%; }
.stats-row { margin-bottom: 16px; }
.stat-card { border-radius: 8px; text-align: center; }
.tab-card { border-radius: 8px; }
</style>
