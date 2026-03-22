<template>
  <div class="pro-page-container">
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>会员管理</a-breadcrumb-item>
      <a-breadcrumb-item>积分库存</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 统计卡片 -->
    <a-row :gutter="16" class="stats-row">
      <a-col :span="8">
        <a-card class="stat-card">
          <a-statistic title="积分总池" :value="stats.totalPool" />
        </a-card>
      </a-col>
      <a-col :span="8">
        <a-card class="stat-card">
          <a-statistic title="已发放" :value="stats.issued" :value-style="{ color: '#1890ff' }" />
        </a-card>
      </a-col>
      <a-col :span="8">
        <a-card class="stat-card">
          <a-statistic title="本月消耗" :value="stats.consumed" :value-style="{ color: '#52c41a' }" />
        </a-card>
      </a-col>
    </a-row>

    <!-- 搜索栏 -->
    <div class="pro-search-bar">
      <a-space wrap>
        <a-input-search v-model="filters.keyword" placeholder="搜索说明" style="width: 240px" search-button @search="loadData" />
      </a-space>
    </div>

    <!-- 操作栏 -->
    <div class="pro-action-bar">
      <a-space>
        <a-button type="primary" @click="showRechargeModal = true">充值</a-button>
        <a-button @click="loadData">刷新</a-button>
      </a-space>
    </div>

    <!-- 数据表格 -->
    <div class="pro-content-area">
      <a-table
        :columns="columns"
        :data="data"
        :loading="loading"
        :pagination="pagination"
        @page-change="onPageChange"
        row-key="id"
      >
        <template #type="{ record }">
          <a-tag :color="record.type === 'add' ? 'green' : 'red'">{{ record.type === 'add' ? '充值' : '消耗' }}</a-tag>
        </template>
        <template #amount="{ record }">
          <span :style="{ color: record.type === 'add' ? '#52c41a' : '#ff4d4f', fontWeight: 600 }">
            {{ record.type === 'add' ? '+' : '-' }}{{ record.amount }}
          </span>
        </template>
        <template #actions="{ record }">
          <a-button type="text" size="small" @click="showDetail(record)">详情</a-button>
        </template>
      </a-table>
    </div>

    <!-- 充值弹窗 -->
    <a-modal v-model:visible="showRechargeModal" title="积分充值" @ok="handleRecharge" :width="400" :mask-closable="false">
      <a-form :model="rechargeForm" layout="vertical">
        <a-form-item label="充值积分数量" required>
          <a-input-number v-model="rechargeForm.amount" :min="1" placeholder="请输入充值积分" style="width: 100%" />
        </a-form-item>
        <a-form-item label="说明">
          <a-textarea v-model="rechargeForm.remark" :rows="3" placeholder="请输入说明" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 详情抽屉 -->
    <a-drawer v-model:visible="detailVisible" title="库存详情" :width="480">
      <template v-if="current">
        <a-descriptions :column="1" bordered size="small">
          <a-descriptions-item label="类型">
            <a-tag :color="current.type === 'add' ? 'green' : 'red'">{{ current.type === 'add' ? '充值' : '消耗' }}</a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="积分数量">{{ current.amount }}</a-descriptions-item>
          <a-descriptions-item label="说明">{{ current.description || '-' }}</a-descriptions-item>
          <a-descriptions-item label="操作人">{{ current.operator || '-' }}</a-descriptions-item>
          <a-descriptions-item label="时间">{{ formatTime(current.created_at) }}</a-descriptions-item>
        </a-descriptions>
      </template>
    </a-drawer>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

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
.pro-page-container { padding: 20px 24px; min-height: calc(100vh - 64px); background: #f5f7fa; }
.pro-breadcrumb { margin-bottom: 16px; }
.stats-row { margin-bottom: 16px; }
.stat-card { border-radius: 8px; text-align: center; }
.pro-search-bar { margin-bottom: 12px; }
.pro-action-bar { margin-bottom: 16px; }
.pro-content-area { background: #fff; border-radius: 8px; padding: 20px; box-shadow: 0 1px 3px rgba(0,0,0,0.04); }
</style>
