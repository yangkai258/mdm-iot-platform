<template>
  <div class="member-coupons-page">
    <!-- 面包屑 -->
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
          <a-statistic title="已使用" :value="stats.used || 0" :value-style="{ color: '#52c41a' }" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="库存余量" :value="stats.remain || 0" :value-style="{ color: '#faad14' }" />
        </a-card>
      </a-col>
    </a-row>

    <!-- 筛选操作区 -->
    <a-card class="action-card">
      <a-space wrap>
        <a-select v-model="filters.type" placeholder="优惠券类型" allow-clear style="width: 120px" @change="handleSearch">
          <a-option :value="1">满减券</a-option>
          <a-option :value="2">折扣券</a-option>
          <a-option :value="3">兑换券</a-option>
        </a-select>
        <a-select v-model="filters.status" placeholder="活动状态" allow-clear style="width: 120px" @change="handleSearch">
          <a-option :value="1">未开始</a-option>
          <a-option :value="2">进行中</a-option>
          <a-option :value="3">已结束</a-option>
        </a-select>
        <a-input-search v-model="filters.keyword" placeholder="搜索优惠券名称" style="width: 200px" search-button @search="handleSearch" />
        <a-button type="primary" @click="showCreateDrawer = true">创建优惠券</a-button>
        <a-button @click="loadCoupons">刷新</a-button>
      </a-space>
    </a-card>

    <!-- 优惠券列表 -->
    <a-card>
      <a-table
        :columns="columns"
        :data="couponList"
        :loading="loading"
        :pagination="paginationConfig"
        @page-change="onPageChange"
        @page-size-change="onPageSizeChange"
        row-key="id"
      >
        <template #type="{ record }">
          <a-tag>{{ getTypeText(record.type) }}</a-tag>
        </template>
        <template #discountValue="{ record }">
          <span>{{ record.type === 2 ? (record.discountValue * 10).toFixed(1) + '折' : '¥' + record.discountValue }}</span>
        </template>
        <template #status="{ record }">
          <a-tag :color="getStatusColor(record.status)">{{ getStatusText(record.status) }}</a-tag>
        </template>
        <template #stock="{ record }">
          <span>{{ record.remainStock || 0 }} / {{ record.totalStock || 0 }}</span>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="viewDetail(record)">详情</a-button>
            <a-button type="text" size="small" @click="showDistribute(record)">发放</a-button>
            <a-button type="text" size="small" @click="showEdit(record)">编辑</a-button>
            <a-button type="text" size="small" status="danger" @click="handleDelete(record)">删除</a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 创建/编辑优惠券抽屉 -->
    <a-drawer v-model:visible="showCreateDrawer" :title="isEdit ? '编辑优惠券' : '创建优惠券'" :width="520" @before-ok="handleFormSubmit" :loading="formLoading">
      <a-form :model="form" layout="vertical">
        <a-form-item label="优惠券名称" required>
          <a-input v-model="form.name" placeholder="请输入优惠券名称" />
        </a-form-item>
        <a-form-item label="优惠券类型" required>
          <a-select v-model="form.type" placeholder="选择类型">
            <a-option :value="1">满减券</a-option>
            <a-option :value="2">折扣券</a-option>
            <a-option :value="3">兑换券</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="优惠方式" required>
          <a-radio-group v-model="form.discountType">
            <a-radio :value="1">满减</a-radio>
            <a-radio :value="2">折扣</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-row :gutter="12">
          <a-col :span="12">
            <a-form-item label="优惠值" required>
              <a-input-number v-model="form.discountValue" :min="0" style="width: 100%" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="使用门槛(元)">
              <a-input-number v-model="form.minAmount" :min="0" style="width: 100%" placeholder="无门槛" />
            </a-form-item>
          </a-col>
        </a-row>
        <a-form-item label="总库存" required>
          <a-input-number v-model="form.totalStock" :min="1" style="width: 100%" />
        </a-form-item>
        <a-form-item label="每人限领">
          <a-input-number v-model="form.perLimit" :min="1" style="width: 100%" placeholder="不限制" />
        </a-form-item>
        <a-form-item label="有效期类型" required>
          <a-radio-group v-model="form.validType">
            <a-radio :value="1">固定日期</a-radio>
            <a-radio :value="2">领取后N天</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item v-if="form.validType === 1" label="有效期范围">
          <a-range-picker v-model="form.dateRange" style="width: 100%" />
        </a-form-item>
        <a-form-item v-if="form.validType === 2" label="有效天数">
          <a-input-number v-model="form.validDays" :min="1" style="width: 100%" />
        </a-form-item>
        <a-form-item label="使用说明">
          <a-textarea v-model="form.description" :rows="3" placeholder="使用说明或适用场景" />
        </a-form-item>
      </a-form>
    </a-drawer>

    <!-- 发放优惠券弹窗 -->
    <a-modal v-model:visible="distributeVisible" title="发放优惠券" :width="520" @before-ok="handleDistribute" :loading="formLoading">
      <a-form :model="distributeForm" layout="vertical">
        <a-form-item label="优惠券">
          <a-input :value="currentCoupon?.name" disabled />
        </a-form-item>
        <a-form-item label="发放数量">
          <a-input-number v-model="distributeForm.count" :min="1" :max="100" style="width: 100%" />
        </a-form-item>
        <a-form-item label="发放说明">
          <a-textarea v-model="distributeForm.note" :rows="2" placeholder="发放备注" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 优惠券详情 -->
    <a-drawer v-model:visible="detailVisible" title="优惠券详情" :width="440">
      <template v-if="currentCoupon">
        <a-descriptions :column="1" bordered size="small">
          <a-descriptions-item label="名称">{{ currentCoupon.name }}</a-descriptions-item>
          <a-descriptions-item label="类型">{{ getTypeText(currentCoupon.type) }}</a-descriptions-item>
          <a-descriptions-item label="面值">{{ currentCoupon.type === 2 ? (currentCoupon.discountValue * 10).toFixed(1) + '折' : '¥' + currentCoupon.discountValue }}</a-descriptions-item>
          <a-descriptions-item label="使用门槛">{{ currentCoupon.minAmount ? '满' + currentCoupon.minAmount + '元' : '无门槛' }}</a-descriptions-item>
          <a-descriptions-item label="总库存">{{ currentCoupon.totalStock }}</a-descriptions-item>
          <a-descriptions-item label="剩余库存">{{ currentCoupon.remainStock }}</a-descriptions-item>
          <a-descriptions-item label="已使用">{{ currentCoupon.usedCount || 0 }}</a-descriptions-item>
          <a-descriptions-item label="状态">
            <a-tag :color="getStatusColor(currentCoupon.status)">{{ getStatusText(currentCoupon.status) }}</a-tag>
          </a-descriptions-item>
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

const loading = ref(false)
const formLoading = ref(false)
const couponList = ref([])
const showCreateDrawer = ref(false)
const distributeVisible = ref(false)
const detailVisible = ref(false)
const isEdit = ref(false)
const currentId = ref(null)
const currentCoupon = ref(null)

const filters = reactive({ type: undefined, status: undefined, keyword: '' })

const pagination = reactive({ current: 1, pageSize: 20, total: 0 })
const paginationConfig = computed(() => ({
  current: pagination.current, pageSize: pagination.pageSize, total: pagination.total,
  showTotal: true, showPageSize: true
}))

const stats = reactive({ total: 0, issued: 0, used: 0, remain: 0 })

const form = reactive({
  name: '', type: 1, discountType: 1, discountValue: 0,
  minAmount: 0, totalStock: 100, perLimit: 1,
  validType: 1, dateRange: [], validDays: 7, description: ''
})

const distributeForm = reactive({ count: 1, note: '' })

const columns = [
  { title: 'ID', dataIndex: 'id', width: 70 },
  { title: '名称', dataIndex: 'name', width: 160 },
  { title: '类型', slotName: 'type', width: 100 },
  { title: '面值', slotName: 'discountValue', width: 100 },
  { title: '使用门槛', dataIndex: 'minAmount', width: 100 },
  { title: '库存', slotName: 'stock', width: 120 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '操作', slotName: 'actions', width: 200 }
]

const getTypeText = (t) => ({ 1: '满减券', 2: '折扣券', 3: '兑换券' }[t] || t)
const getStatusText = (s) => ({ 1: '未开始', 2: '进行中', 3: '已结束' }[s] || s)
const getStatusColor = (s) => ({ 1: 'blue', 2: 'green', 3: 'gray' }[s] || 'gray')

const loadCoupons = async () => {
  loading.value = true
  try {
    const params = { page: pagination.current, pageSize: pagination.pageSize }
    if (filters.type) params.type = filters.type
    if (filters.status) params.status = filters.status
    if (filters.keyword) params.keyword = filters.keyword

    const res = await api.getCouponList(params)
    const d = res.data || {}
    couponList.value = d.list || []
    pagination.total = d.total || 0

    // 汇总统计
    stats.total = d.total || 0
    stats.issued = (d.list || []).reduce((sum, c) => sum + (c.usedCount || 0), 0)
    stats.used = (d.list || []).reduce((sum, c) => sum + (c.usedCount || 0), 0)
    stats.remain = (d.list || []).reduce((sum, c) => sum + (c.remainStock || 0), 0)
  } catch (err) {
    Message.error('加载优惠券列表失败: ' + err.message)
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.current = 1
  loadCoupons()
}

const onPageChange = (page) => { pagination.current = page; loadCoupons() }
const onPageSizeChange = (pageSize) => { pagination.pageSize = pageSize; pagination.current = 1; loadCoupons() }

const showEdit = (record) => {
  isEdit.value = true
  currentId.value = record.id
  Object.assign(form, {
    name: record.name,
    type: record.type || 1,
    discountType: record.discountType || 1,
    discountValue: record.discountValue || 0,
    minAmount: record.minAmount || 0,
    totalStock: record.totalStock || 100,
    perLimit: record.perLimit || 1,
    validType: record.validType || 1,
    validDays: record.validDays || 7,
    description: record.description || ''
  })
  showCreateDrawer.value = true
}

const viewDetail = (record) => {
  currentCoupon.value = record
  detailVisible.value = true
}

const showDistribute = (record) => {
  currentCoupon.value = record
  distributeForm.count = 1
  distributeForm.note = ''
  distributeVisible.value = true
}

const handleFormSubmit = async (done) => {
  if (!form.name) { Message.warning('请填写名称'); done(false); return }
  formLoading.value = true
  try {
    const payload = {
      ...form,
      startTime: form.dateRange?.[0] || '',
      endTime: form.dateRange?.[1] || ''
    }
    if (isEdit.value) {
      await api.updateCoupon(currentId.value, payload)
      Message.success('更新成功')
    } else {
      await api.createCoupon(payload)
      Message.success('创建成功')
    }
    showCreateDrawer.value = false
    loadCoupons()
    done(true)
  } catch (err) {
    Message.error(err.message || '操作失败')
    done(false)
  } finally {
    formLoading.value = false
  }
}

const handleDistribute = async (done) => {
  if (!distributeForm.count) { Message.warning('请填写发放数量'); done(false); return }
  formLoading.value = true
  try {
    await api.grantCoupon({ couponId: currentCoupon.value.id, count: distributeForm.count, note: distributeForm.note })
    Message.success('发放成功')
    distributeVisible.value = false
    done(true)
  } catch (err) {
    Message.error(err.message || '发放失败')
    done(false)
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

onMounted(() => loadCoupons())
</script>

<style scoped>
.member-coupons-page {
  padding: 20px 24px;
  min-height: calc(100vh - 64px);
  background: #f5f7fa;
}
.breadcrumb { margin-bottom: 16px; }
.stats-row { margin-bottom: 16px; }
.stat-card { border-radius: 8px; }
.action-card { margin-bottom: 16px; }
</style>
