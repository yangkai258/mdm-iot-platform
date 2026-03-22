<template>
  <div class="member-coupons">
    <!-- 统计卡片 -->
    <a-row :gutter="16" class="stats-row">
      <a-col :span="6">
        <a-card>
          <a-statistic title="优惠券总数" :value="stats.total" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card>
          <a-statistic title="已发放" :value="stats.issued" :value-style="{ color: '#1890ff' }" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card>
          <a-statistic title="已使用" :value="stats.used" :value-style="{ color: '#52c41a' }" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card>
          <a-statistic title="未使用" :value="stats.unused" :value-style="{ color: '#faad14' }" />
        </a-card>
      </a-col>
    </a-row>

    <!-- 操作区域 -->
    <a-card class="action-card">
      <a-space wrap>
        <a-select v-model="filters.status" placeholder="优惠券状态" allow-clear style="width: 120px" @change="loadCoupons">
          <a-option value="active">有效</a-option>
          <a-option value="expired">已过期</a-option>
          <a-option value="used">已使用</a-option>
        </a-select>
        <a-input-search v-model="filters.keyword" placeholder="搜索优惠券名称" style="width: 200px" search-button @search="loadCoupons" />
        <a-button type="primary" @click="showCreateDrawer = true">创建优惠券</a-button>
        <a-button @click="loadCoupons">刷新</a-button>
      </a-space>
    </a-card>

    <!-- 优惠券列表 -->
    <a-card>
      <a-table :columns="columns" :data="couponList" :loading="loading" :pagination="pagination" row-key="id">
        <template #type="{ record }">
          <a-tag>{{ getTypeText(record.type) }}</a-tag>
        </template>
        <template #status="{ record }">
          <a-tag :color="getStatusColor(record.status)">{{ getStatusText(record.status) }}</a-tag>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="viewDetail(record)">详情</a-button>
            <a-button type="text" size="small" @click="distributeCoupon(record)">发放</a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 创建优惠券抽屉 -->
    <a-drawer v-model:visible="showCreateDrawer" title="创建优惠券" :width="480">
      <a-form :model="form" layout="vertical">
        <a-form-item label="优惠券名称" required>
          <a-input v-model="form.name" placeholder="请输入名称" />
        </a-form-item>
        <a-form-item label="优惠券类型" required>
          <a-select v-model="form.type" placeholder="选择类型">
            <a-option value="discount">折扣券</a-option>
            <a-option value="cash">现金券</a-option>
            <a-option value="gift">礼品券</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="面值/折扣">
          <a-input-number v-model="form.value" :min="0" />
        </a-form-item>
        <a-form-item>
          <a-button type="primary" @click="handleCreate">创建</a-button>
        </a-form-item>
      </a-form>
    </a-drawer>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const loading = ref(false)
const couponList = ref([])
const showCreateDrawer = ref(false)

const filters = reactive({ status: undefined, keyword: '' })
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })
const stats = reactive({ total: 0, issued: 0, used: 0, unused: 0 })
const form = reactive({ name: '', type: '', value: 0 })

const columns = [
  { title: 'ID', dataIndex: 'id', width: 80 },
  { title: '优惠券名称', dataIndex: 'name' },
  { title: '类型', slotName: 'type', width: 100 },
  { title: '面值', dataIndex: 'value', width: 100 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '有效期', dataIndex: 'expire_time', width: 160 },
  { title: '操作', slotName: 'actions', width: 160 }
]

const getTypeText = (type) => ({ discount: '折扣券', cash: '现金券', gift: '礼品券' }[type] || type)
const getStatusText = (s) => ({ active: '有效', expired: '已过期', used: '已使用' }[s] || s)
const getStatusColor = (s) => ({ active: 'green', expired: 'gray', used: 'blue' }[s] || 'gray')

const loadCoupons = () => { loading.value = true; setTimeout(() => { loading.value = false }, 300) }
const viewDetail = (r) => Message.info('查看详情')
const distributeCoupon = (r) => Message.info('发放优惠券')
const handleCreate = () => { Message.success('创建成功'); showCreateDrawer.value = false }

onMounted(() => loadCoupons())
</script>

<style scoped>
.member-coupons { padding: 16px; }
.stats-row { margin-bottom: 16px; }
.action-card { margin-bottom: 16px; }
</style>
