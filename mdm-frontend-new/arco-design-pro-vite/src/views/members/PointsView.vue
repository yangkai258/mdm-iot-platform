<template>
    <Breadcrumb :items="['Home','Console','']" />


  <div class="container">
    <a-card class="general-card" title="积分管理">
      <template #extra>
        <a-button @click="loadData"><icon-refresh />刷新</a-button>
      </template>
      <a-row :gutter="16">
        <a-col :span="8">
          <a-form-item label="关键词">
            <a-input v-model="filters.keyword" placeholder="姓名/手机号" @pressEnter="loadData" />
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
      <a-table :columns="memberColumns" :data="memberList" :loading="loading" :pagination="paginationConfig" @page-change="onPageChange" row-key="id" />
    </a-card>
  </div>
</template>
      </a-table>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import Breadcrumb from '@/components/Breadcrumb.vue'

const memberList = ref([])
const loading = ref(false)
const filters = reactive({ keyword: '', levelId: undefined })
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })
const paginationConfig = computed(() => ({ current: pagination.current, pageSize: pagination.pageSize, total: pagination.total, showTotal: true }))
const memberColumns = [
  { title: '会员名称', dataIndex: 'name', width: 150 },
  { title: '手机号', dataIndex: 'mobile', width: 130 },
  { title: '当前积分', slotName: 'points', width: 120 },
  { title: '成长值', slotName: 'growthValue', width: 100 },
  { title: '注册时间', dataIndex: 'createdAt', width: 160 }
]

const loadData = async () => {
  loading.value = true
  try {
    const params = { page: pagination.current, pageSize: pagination.pageSize }
    if (filters.keyword) params.keyword = filters.keyword
    const res = await fetch(`/api/v1/members/points?${new URLSearchParams(params)}`, {
      headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token') }
    }).then(r => r.json())
    memberList.value = res.data?.list || []
    pagination.total = res.data?.total || 0
  } catch { memberList.value = [] } finally { loading.value = false }
}

const onPageChange = (page) => { pagination.current = page; loadData() }

onMounted(() => loadData())
</script>
</a-card>
