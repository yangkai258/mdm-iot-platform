<template>
  <div class="coupon-page">
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>会员管理</a-breadcrumb-item>
      <a-breadcrumb-item>优惠券管理</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 统计卡片 -->
    <a-row :gutter="16" class="stats-row">
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="优惠券总数" :value="stats.total || 0" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="已发放" :value="stats.issued || 0" :value-style="{ color: '#1890ff' }" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="已核销" :value="stats.used || 0" :value-style="{ color: '#52c41a' }" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="未使用" :value="stats.unused || 0" :value-style="{ color: '#faad14' }" />
        </a-card>
      </a-col>
    </a-row>

    <!-- 搜索筛选 -->
    <a-card class="action-card">
      <a-space wrap>
        <a-input-search v-model="filters.keyword" placeholder="搜索优惠券名称" style="width: 220px" search-button @search="loadCoupons" />
        <a-select v-model="filters.type" placeholder="优惠券类型" allow-clear style="width: 140px" @change="loadCoupons">
          <a-option value="discount">折扣券</a-option>
          <a-option value="cash">现金券</a-option>
          <a-option value="gift">礼品券</a-option>
          <a-option value="shipping">包邮券</a-option>
        </a-select>
        <a-select v-model="filters.status" placeholder="状态" allow-clear style="width: 120px" @change="loadCoupons">
          <a-option value="active">有效</a-option>
          <a-option value="inactive">未激活</a-option>
          <a-option value="expired">已过期</a-option>
        </a-select>
        <a-button type="primary" @click="showCreateDrawer">创建优惠券</a-button>
        <a-button @click="loadCoupons">刷新</a-button>
      </a-space>
    </a-card>

    <!-- 优惠券列表 -->
    <a-card class="table-card">
      <a-table
        :columns="columns"
        :data="couponList"
        :loading="loading"
        :pagination="paginationConfig"
        @page-change="onPageChange"
        @page-size-change="onPageSizeChange"
        row-key="id"
        :scroll="{ x: 1100 }"
      >
        <template #type="{ record }">
          <a-tag :color="getTypeColor(record.type)">{{ getTypeText(record.type) }}</a-tag>
        </template>
        <template #value="{ record }">
          <span v-if="record.type === 'discount'" style="color: #1890ff; font-weight: 600;">{{ ((1 - (record.value || 0)) * 100).toFixed(0) }}折</span>
          <span v-else style="color: #ff6b00; font-weight: 600;">¥{{ record.value || 0 }}</span>
        </template>
        <template #status="{ record }">
          <a-tag :color="getStatusColor(record.status)">{{ getStatusText(record.status) }}</a-tag>
        </template>
        <template #minConsume="{ record }">
          {{ record.minConsume > 0 ? '满¥' + record.minConsume : '无门槛' }}
        </template>
        <template #usedCount="{ record }">
          {{ record.usedCount || 0 }} / {{ record.totalCount || '∞' }}
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="showDetail(record)">详情</a-button>
            <a-button type="text" size="small" @click="showGrant(record)">发放</a-button>
            <a-button type="text" size="small" @click="showEdit(record)">编辑</a-button>
            <a-button type="text" size="small" status="danger" @click="handleDelete(record)">删除</a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 创建/编辑优惠券抽屉 -->
    <a-drawer v-model:visible="formVisible" :title="isEdit ? '编辑优惠券' : '创建优惠券'" :width="520">
      <a-form :model="form" layout="vertical">
        <a-form-item label="优惠券名称" required>
          <a-input v-model="form.name" placeholder="请输入优惠券名称" />
        </a-form-item>
        <a-form-item label="优惠券类型" required>
          <a-select v-model="form.type" placeholder="选择类型">
            <a-option value="discount">折扣券</a-option>
            <a-option value="cash">现金券</a-option>
            <a-option value="gift">礼品券</a-option>
            <a-option value="shipping">包邮券</a-option>
          </a-select>
        </a-form-item>
        <a-row :gutter="12">
          <a-col :span="12">
            <a-form-item label="面值/折扣">
              <a-input-number v-model="form.value" :min="0" style="width: 100%" />
              <template #extra>折扣填0-1之间</template>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="发行总量">
              <a-input-number v-model="form.totalCount" :min="0" style="width: 100%" placeholder="0=不限量" />
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="12">
          <a-col :span="12">
            <a-form-item label="最低消费">
              <a-input-number v-model="form.minConsume" :min="0" :precision="2" style="width: 100%" placeholder="0=无门槛" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="每人限领">
              <a-input-number v-model="form.perLimit" :min="0" style="width: 100%" placeholder="0=不限" />
            </a-form-item>
          </a-col>
        </a-row>
        <a-form-item label="有效期">
          <a-range-picker v-model="form.dateRange" style="width: 100%" />
        </a-form-item>
        <a-form-item label="使用说明">
          <a-textarea v-model="form.description" :rows="3" placeholder="描述优惠券使用规则" />
        </a-form-item>
      </a-form>
      <template #footer>
        <a-button @click="formVisible = false">取消</a-button>
        <a-button type="primary" :loading="formLoading" @click="handleFormSubmit">{{ isEdit ? '保存' : '创建' }}</a-button>
      </template>
    </a-drawer>

    <!-- 发放优惠券抽屉 -->
    <a-drawer v-model:visible="grantVisible" title="发放优惠券" :width="520">
      <a-form layout="vertical">
        <a-form-item label="优惠券">
          <a-input :value="currentCoupon?.name" disabled />
        </a-form-item>
        <a-form-item label="发放方式" required>
          <a-radio-group v-model="grantForm.mode">
            <a-radio value="member">指定会员</a-radio>
            <a-radio value="level">指定等级</a-radio>
            <a-radio value="all">全部会员</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item v-if="grantForm.mode === 'member'" label="选择会员">
          <a-select v-model="grantForm.memberIds" placeholder="选择会员" multiple searchable>
            <a-option v-for="m in memberOptions" :key="m.id" :value="m.id">{{ m.name }} ({{ m.mobile }})</a-option>
          </a-select>
        </a-form-item>
        <a-form-item v-if="grantForm.mode === 'level'" label="选择等级">
          <a-select v-model="grantForm.levelId" placeholder="选择会员等级">
            <a-option v-for="lv in levelOptions" :key="lv.id" :value="lv.id">{{ lv.name }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="发放数量">
          <a-input-number v-model="grantForm.count" :min="1" :max="10" style="width: 100%" />
        </a-form-item>
      </a-form>
      <template #footer>
        <a-button @click="grantVisible = false">取消</a-button>
        <a-button type="primary" @click="handleGrant">确认发放</a-button>
      </template>
    </a-drawer>

    <!-- 优惠券详情抽屉 -->
    <a-drawer v-model:visible="detailVisible" title="优惠券详情" :width="480">
      <template v-if="currentCoupon">
        <a-descriptions :column="1" bordered size="small">
          <a-descriptions-item label="优惠券名称">{{ currentCoupon.name }}</a-descriptions-item>
          <a-descriptions-item label="类型">
            <a-tag :color="getTypeColor(currentCoupon.type)">{{ getTypeText(currentCoupon.type) }}</a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="面值/折扣">
            <span v-if="currentCoupon.type === 'discount'" style="color: #1890ff">{{ ((1 - (currentCoupon.value || 0)) * 100).toFixed(0) }}折</span>
            <span v-else style="color: #ff6b00">¥{{ currentCoupon.value || 0 }}</span>
          </a-descriptions-item>
          <a-descriptions-item label="发行总量">{{ currentCoupon.totalCount || '不限量' }}</a-descriptions-item>
          <a-descriptions-item label="已使用">{{ currentCoupon.usedCount || 0 }}</a-descriptions-item>
          <a-descriptions-item label="状态">
            <a-tag :color="getStatusColor(currentCoupon.status)">{{ getStatusText(currentCoupon.status) }}</a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="有效期">{{ currentCoupon.startTime || '-' }} 至 {{ currentCoupon.endTime || '-' }}</a-descriptions-item>
          <a-descriptions-item label="最低消费">{{ currentCoupon.minConsume > 0 ? '满¥' + currentCoupon.minConsume : '无门槛' }}</a-descriptions-item>
          <a-descriptions-item label="使用说明">{{ currentCoupon.description || '-' }}</a-descriptions-item>
        </a-descriptions>
      </template>
    </a-drawer>
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
.coupon-page { padding: 20px 24px; min-height: calc(100vh - 64px); background: #f5f7fa; }
.breadcrumb { margin-bottom: 16px; }
.stats-row { margin-bottom: 16px; }
.stat-card { border-radius: 8px; text-align: center; }
.action-card { margin-bottom: 16px; }
.table-card { border-radius: 8px; }
</style>
