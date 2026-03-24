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
.page-container { background: #fff; border-radius: 4px; padding: 20px; }
.search-form { margin-bottom: 16px; padding: 16px; background: #f7f8fa; border-radius: 4px; }
.toolbar { margin-bottom: 16px; }
</style>
