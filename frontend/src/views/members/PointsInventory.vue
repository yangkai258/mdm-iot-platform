<template>
  <div class="container">
    <a-card class="general-card" title="积分库存">
      <template #extra>
        <a-space :size="12">
          <a-button type="primary" @click="showRechargeModal = true"><icon-plus />充值库存</a-button>
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
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id" />
    </a-card>
      </a-table>
    <a-modal v-model:visible="showRechargeModal" title="充值积分">
      <a-form :model="rechargeForm" label-col-flex="100px">
        <a-form-item label="充值积分"><a-input-number v-model="rechargeForm.amount" :min="1" style="width: 100%" /></a-form-item>
        <a-form-item label="备注"><a-input v-model="rechargeForm.remark" /></a-form-item>
      </a-form>
      <template #footer>
        <a-button @click="showRechargeModal = false">取消</a-button>
        <a-button type="primary" @click="handleRecharge">确定</a-button>
      </template>
    </a-modal>
  </div>
</template>

<script setup>

import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import Breadcrumb from '@/components/Breadcrumb.vue'

const API_BASE = '/api/v1'
const data = ref([])
const loading = ref(false)
const showRechargeModal = ref(false)
const detailVisible = ref(false)
const current = ref(null)

const filters = reactive({ keyword: '' })
const pagination = reactive({ current: 1, pageSize: 10, total: 0 })
const stats = reactive({ totalPool: 0, issued: 0, consumed: 0 })
const rechargeForm = reactive({ amount: 0, remark: '' })

const columns = [
  { title: 'ID', dataIndex: 'id', width: 80 },
  { title: '类型', slotName: 'type', width: 100 },
  { title: '积分数量', slotName: 'amount', width: 120 },
  { title: '说明', dataIndex: 'description' },
  { title: '操作人', dataIndex: 'operator', width: 120 },
  { title: '时间', dataIndex: 'created_at', width: 170 },
  { title: '操作', slotName: 'actions', width: 100 }
]

const formatTime = (t) => t ? new Date(t).toLocaleString('zh-CN') : '-'

const loadData = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const params = new URLSearchParams({ page: pagination.current, page_size: pagination.pageSize })
    if (filters.keyword) params.append('keyword', filters.keyword)
    const res = await fetch(`${API_BASE}/member/points/inventory?${params}`, {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    const resp = await res.json()
    if (resp.code === 0) {
      data.value = resp.data?.list || resp.data || []
      pagination.total = resp.data?.total || 0
      if (resp.data?.stats) Object.assign(stats, resp.data.stats)
    }
  } catch (e) {
    Message.error('加载库存信息失败')
  } finally {
    loading.value = false
  }
}

const handleRecharge = async () => {
  if (!rechargeForm.amount || rechargeForm.amount <= 0) { Message.warning('请输入正确的积分数量'); return }
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`${API_BASE}/member/points/inventory/recharge`, {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' },
      body: JSON.stringify(rechargeForm)
    })
    const data = await res.json()
    if (data.code === 0) { Message.success('充值成功'); showRechargeModal.value = false; loadData() }
    else Message.error(data.message || '充值失败')
  } catch (e) { Message.error('充值失败') }
}

const showDetail = (record) => { current.value = record; detailVisible.value = true }
const onPageChange = (page) => { pagination.current = page; loadData() }

onMounted(() => loadData())

</script>

<style scoped>
.page-container { background: #fff; border-radius: 4px; padding: 20px; }
.search-form { margin-bottom: 16px; padding: 16px; background: #f7f8fa; border-radius: 4px; }
.toolbar { margin-bottom: 16px; }
</style>
</a-card>
