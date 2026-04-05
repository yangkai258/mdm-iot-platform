<template>
    <Breadcrumb :items="['Home','Console','']" />


  <div class="container">
    <a-card class="general-card" title="优惠券库存">
      <template #extra>
        <a-space :size="12">
          <a-button type="primary" @click="showRechargeDrawer(null)"><icon-plus />充值库存</a-button>
          <a-button @click="loadData"><icon-refresh />刷新</a-button>
        </a-space>
      </template>
      <a-row :gutter="16">
        <a-col :span="8">
          <a-form-item label="关键词">
            <a-input v-model="filters.keyword" placeholder="请输入" @pressEnter="loadData" />
          </a-form-item>
        </a-col>
        <a-col :flex="'86px'" style="display: flex; align-items: flex-end">
          <a-space direction="vertical" :size="8">
            <a-button type="primary" @click="loadData">查询</a-button>
            <a-button @click="filters.keyword = ''; loadData()">重置</a-button>
          </a-space>
        </a-col>
      </a-row>
      <a-divider style="margin: 0 0 16px 0" />
      <a-table :columns="columns" :data="dataList" :loading="loading" :pagination="paginationConfig" @page-change="onPageChange" row-key="id">
        <template #actions="{ record }">
          <a-button type="text" size="small" @click="showRechargeDrawer(record)">充值</a-button>
        </template>
      </a-table>
    </a-card>
    <a-modal v-model:visible="rechargeVisible" title="充值库存">
      <a-form :model="rechargeForm" label-col-flex="100px">
        <a-form-item label="优惠券"><a-input :value="selectedCouponName" readonly /></a-form-item>
        <a-form-item label="充值数量"><a-input-number v-model="rechargeForm.addCount" :min="1" style="width: 100%" /></a-form-item>
      </a-form>
      <template #footer>
        <a-button @click="rechargeVisible = false">取消</a-button>
        <a-button type="primary" @click="handleRecharge">确定</a-button>
      </template>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import Breadcrumb from '@/components/Breadcrumb.vue'

const dataList = ref([])
const loading = ref(false)
const rechargeVisible = ref(false)
const selectedCouponName = ref('')
const filters = reactive({ keyword: '' })
const rechargeForm = reactive({ couponId: undefined, addCount: 1 })
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })
const paginationConfig = computed(() => ({ current: pagination.current, pageSize: pagination.pageSize, total: pagination.total, showTotal: true }))
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
    const res = await fetch(`/api/members/coupons/inventory?page=${pagination.current}&page_size=${pagination.pageSize}&keyword=${filters.keyword}`, {
      headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token') }
    }).then(r => r.json())
    dataList.value = res.data?.list || []
    pagination.total = res.data?.total || 0
  } catch { dataList.value = [] } finally { loading.value = false }
}

const showRechargeDrawer = (record) => { rechargeForm.couponId = record?.id; rechargeForm.addCount = 1; rechargeVisible.value = true }
const handleRecharge = () => { rechargeVisible.value = false; Message.success('充值成功'); loadData() }
const onPageChange = (page) => { pagination.current = page; loadData() }

onMounted(() => loadData())
</script>
