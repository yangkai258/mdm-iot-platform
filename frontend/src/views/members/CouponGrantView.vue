<template>
  <div class="coupon-grant-page">
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>会员管理</a-breadcrumb-item>
      <a-breadcrumb-item>优惠券发放</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 搜索筛选 -->
    <a-card class="action-card">
      <a-space wrap>
        <a-input-search v-model="filters.keyword" placeholder="搜索优惠券名称" style="width: 220px" search-button @search="loadData" />
        <a-select v-model="filters.mode" placeholder="发放方式" allow-clear style="width: 140px" @change="loadData">
          <a-option value="member">指定会员</a-option>
          <a-option value="level">指定等级</a-option>
          <a-option value="all">全部会员</a-option>
        </a-select>
        <a-range-picker v-model="filters.dateRange" style="width: 260px" @change="loadData" />
        <a-button type="primary" @click="showGrantDrawer">发放优惠券</a-button>
        <a-button @click="handleExport">导出</a-button>
        <a-button @click="loadData">刷新</a-button>
      </a-space>
    </a-card>

    <!-- 发放记录列表 -->
    <a-card class="table-card">
      <a-table
        :columns="columns"
        :data="dataList"
        :loading="loading"
        :pagination="paginationConfig"
        @page-change="onPageChange"
        @page-size-change="onPageSizeChange"
        row-key="id"
        :scroll="{ x: 1100 }"
      >
        <template #mode="{ record }">
          <a-tag :color="getModeColor(record.mode)">{{ getModeText(record.mode) }}</a-tag>
        </template>
        <template #target="{ record }">
          <span v-if="record.mode === 'member'">{{ record.memberName || '-' }}</span>
          <span v-else-if="record.mode === 'level'">{{ record.levelName || '-' }}</span>
          <span v-else>全部会员</span>
        </template>
        <template #status="{ record }">
          <a-tag :color="getStatusColor(record.status)">{{ getStatusText(record.status) }}</a-tag>
        </template>
      </a-table>
    </a-card>

    <!-- 发放优惠券抽屉 -->
    <a-drawer v-model:visible="grantVisible" title="发放优惠券" :width="520">
      <a-form layout="vertical">
        <a-form-item label="优惠券" required>
          <a-select v-model="grantForm.couponId" placeholder="选择优惠券" searchable>
            <a-option v-for="c in couponOptions" :key="c.id" :value="c.id">{{ c.name }}</a-option>
          </a-select>
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
        <a-button type="primary" :loading="formLoading" @click="handleGrant">确认发放</a-button>
      </template>
    </a-drawer>
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
.coupon-grant-page { padding: 20px 24px; min-height: calc(100vh - 64px); background: #f5f7fa; }
.breadcrumb { margin-bottom: 16px; }
.action-card { margin-bottom: 16px; }
.table-card { border-radius: 8px; }
</style>
