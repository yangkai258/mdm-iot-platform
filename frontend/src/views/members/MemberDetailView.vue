<template>
  <div class="container">
    <a-card class="general-card" title="会员详情">
      <template #extra>
        <a-space :size="12">
          <a-button type="primary" @click="showEdit"><icon-edit />编辑</a-button>
          <a-button @click="loadData"><icon-refresh />刷新</a-button>
        </a-space>
      </template>
      <a-tabs v-model:active-key="activeTab">
        <a-tab-pane key="info" title="基本信息">
          <a-descriptions :column="2" bordered>
            <a-descriptions-item label="姓名">{{ member.name }}</a-descriptions-item>
            <a-descriptions-item label="手机号">{{ member.mobile }}</a-descriptions-item>
            <a-descriptions-item label="会员等级">{{ member.levelName }}</a-descriptions-item>
            <a-descriptions-item label="状态">{{ member.status }}</a-descriptions-item>
            <a-descriptions-item label="积分">{{ member.points }}</a-descriptions-item>
            <a-descriptions-item label="累计消费">{{ member.totalConsume }}</a-descriptions-item>
          </a-descriptions>
        </a-tab-pane>
        <a-tab-pane key="orders" title="订单记录">
          <a-table :columns="orderColumns" :data="orders" :loading="loading" :pagination="paginationConfig" row-key="id" />
        </a-tab-pane>
        <a-tab-pane key="points" title="积分记录">
          <a-table :columns="pointsColumns" :data="pointsRecords" :loading="loading" :pagination="paginationConfig" row-key="id" />
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>
      </a-table>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { Message } from '@arco-design/web-vue'
import Breadcrumb from '@/components/Breadcrumb.vue'

const route = useRoute()
const memberId = computed(() => route.params.id)
const member = ref({})
const orders = ref([])
const pointsRecords = ref([])
const loading = ref(false)
const activeTab = ref('info')
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })
const paginationConfig = computed(() => ({ current: pagination.current, pageSize: pagination.pageSize, total: pagination.total, showTotal: true }))
const orderColumns = [
  { title: '订单号', dataIndex: 'orderNo', width: 200 },
  { title: '金额', dataIndex: 'amount', width: 120 },
  { title: '状态', dataIndex: 'status', width: 100 },
  { title: '时间', dataIndex: 'createdAt', width: 170 }
]
const pointsColumns = [
  { title: '类型', dataIndex: 'type', width: 90 },
  { title: '积分', dataIndex: 'points', width: 100 },
  { title: '原因', dataIndex: 'reason', ellipsis: true },
  { title: '时间', dataIndex: 'createdAt', width: 170 }
]

const loadData = async () => {
  if (!memberId.value) return
  loading.value = true
  try {
    const res = await fetch(`/api/v1/members/${memberId.value}`, {
      headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token') }
    }).then(r => r.json())
    member.value = res.data || {}
  } catch { Message.error('加载失败') } finally { loading.value = false }
}

const showEdit = () => { Message.info('编辑功能开发中') }

onMounted(() => loadData())
</script>
</a-card>
