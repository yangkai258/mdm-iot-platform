<template>
  <div class="coupon-inventory-page">
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>会员管理</a-breadcrumb-item>
      <a-breadcrumb-item>优惠券库存</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 搜索筛选 -->
    <a-card class="action-card">
      <a-space wrap>
        <a-input-search v-model="filters.keyword" placeholder="搜索优惠券名称" style="width: 220px" search-button @search="loadData" />
        <a-button type="primary" @click="showRechargeDrawer">充值库存</a-button>
        <a-button @click="loadData">刷新</a-button>
      </a-space>
    </a-card>

    <!-- 库存列表 -->
    <a-card class="table-card">
      <a-table
        :columns="columns"
        :data="dataList"
        :loading="loading"
        :pagination="paginationConfig"
        @page-change="onPageChange"
        @page-size-change="onPageSizeChange"
        row-key="id"
        :scroll="{ x: 900 }"
      >
        <template #usage="{ record }">
          <span style="color: #52c41a; font-weight: 600;">{{ record.usedCount || 0 }}</span>
          <span style="color: #999; margin: 0 4px;">/</span>
          <span>{{ record.totalCount || 0 }}</span>
        </template>
        <template #remaining="{ record }">
          <span :style="{ color: (record.totalCount - (record.usedCount || 0)) > 0 ? '#1890ff' : '#ff4d4f', fontWeight: 600 }">
            {{ (record.totalCount || 0) - (record.usedCount || 0) }}
          </span>
        </template>
        <template #actions="{ record }">
          <a-button type="text" size="small" @click="showRechargeDrawer(record)">充值</a-button>
        </template>
      </a-table>
    </a-card>

    <!-- 充值库存抽屉 -->
    <a-drawer v-model:visible="rechargeVisible" title="充值库存" :width="480">
      <a-form layout="vertical">
        <a-form-item label="优惠券">
          <a-select v-model="rechargeForm.couponId" placeholder="选择优惠券" searchable>
            <a-option v-for="c in couponOptions" :key="c.id" :value="c.id">{{ c.name }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="充值数量">
          <a-input-number v-model="rechargeForm.addCount" :min="1" style="width: 100%" />
        </a-form-item>
      </a-form>
      <template #footer>
        <a-button @click="rechargeVisible = false">取消</a-button>
        <a-button type="primary" :loading="formLoading" @click="handleRecharge">确认充值</a-button>
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
const loading = ref(false)
const formLoading = ref(false)
const rechargeVisible = ref(false)

const filters = reactive({ keyword: '' })
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })

const paginationConfig = computed(() => ({
  current: pagination.current, pageSize: pagination.pageSize, total: pagination.total,
  showTotal: true, showPageSize: true, pageSizeOptions: [10, 20, 50, 100]
}))

const rechargeForm = reactive({ couponId: undefined, addCount: 1 })

const columns = [
  { title: '优惠券名称', dataIndex: 'couponName', width: 200 },
  { title: '类型', dataIndex: 'type', width: 100 },
  { title: '总数量', dataIndex: 'totalCount', width: 110 },
  { title: '已使用', slotName: 'usage', width: 110 },
  { title: '剩余库存', slotName: 'remaining', width: 110 },
  { title: '操作', slotName: 'actions', width: 100 }
]

const loadData = async () => {
  loading.value = true
  try {
    const params = { page: pagination.current, pageSize: pagination.pageSize }
    if (filters.keyword) params.keyword = filters.keyword
    const res = await api.getCouponInventoryList(params)
    const d = res.data || {}
    dataList.value = d.list || []
    pagination.total = d.total || 0
  } catch (err) {
    Message.error('加载库存列表失败: ' + err.message)
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

const showRechargeDrawer = (record) => {
  rechargeForm.couponId = record?.id
  rechargeForm.addCount = 1
  loadCouponOptions()
  rechargeVisible.value = true
}

const handleRecharge = async () => {
  if (!rechargeForm.couponId) { Message.warning('请选择优惠券'); return }
  formLoading.value = true
  try {
    await api.rechargeCouponInventory({ couponId: rechargeForm.couponId, addCount: rechargeForm.addCount })
    Message.success('充值成功')
    rechargeVisible.value = false
    loadData()
  } catch (err) {
    Message.error(err.message || '充值失败')
  } finally {
    formLoading.value = false
  }
}

const onPageChange = (page) => { pagination.current = page; loadData() }
const onPageSizeChange = (pageSize) => { pagination.pageSize = pageSize; pagination.current = 1; loadData() }

onMounted(() => loadData())
</script>

<style scoped>
.coupon-inventory-page { padding: 20px 24px; min-height: calc(100vh - 64px); background: #f5f7fa; }
.breadcrumb { margin-bottom: 16px; }
.action-card { margin-bottom: 16px; }
.table-card { border-radius: 8px; }
</style>
