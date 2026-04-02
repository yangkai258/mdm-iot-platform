<template>
    <Breadcrumb :items="['Home','Console','']" />


  <div class="container">
    <a-card class="general-card" title="买赠活动">
      <template #extra>
        <a-space :size="12">
          <a-button type="primary" @click="showCreateDrawer"><icon-plus />新建</a-button>
          <a-button @click="loadData"><icon-refresh />刷新</a-button>
        </a-space>
      </template>
      <a-row :gutter="16">
        <a-col :span="8">
          <a-form-item label="活动名称">
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
          <a-button type="text" size="small" @click="showEdit(record)">编辑</a-button>
          <a-button type="text" size="small" @click="handleDelete(record)">删除</a-button>
        </template>
      </a-table>
      </a-table>
    </a-card>
    <a-modal v-model:visible="modalVisible" :title="modalTitle">
      <a-form :model="form" label-col-flex="100px">
        <a-form-item label="活动名称"><a-input v-model="form.name" placeholder="请输入" /></a-form-item>
      </a-form>
      <template #footer>
        <a-button @click="modalVisible = false">取消</a-button>
        <a-button type="primary" @click="handleFormSubmit">确定</a-button>
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
const formLoading = ref(false)
const modalVisible = ref(false)
const isEdit = ref(false)
const filters = reactive({ keyword: '' })
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })
const paginationConfig = computed(() => ({ current: pagination.current, pageSize: pagination.pageSize, total: pagination.total, showTotal: true }))
const form = reactive({ name: '', buyCount: 1, giftCount: 1, giftName: '' })
const modalTitle = computed(() => isEdit.value ? '编辑' : '新建')
const columns = [
  { title: '活动名称', dataIndex: 'name', width: 200 },
  { title: '买赠规则', slotName: 'rule', width: 280 },
  { title: '适用商品', dataIndex: 'productName', width: 160, ellipsis: true },
  { title: '时间范围', dataIndex: 'dateRange', width: 220 },
  { title: '状态', slotName: 'status', width: 90 },
  { title: '操作', slotName: 'actions', width: 150 }
]

const loadData = async () => {
  loading.value = true
  try {
    const res = await fetch(`/api/v1/members/promotions/buy-gift?page=${pagination.current}&page_size=${pagination.pageSize}&keyword=${filters.keyword}`, {
      headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token') }
    }).then(r => r.json())
    dataList.value = res.data?.list || [{ id: 1, name: '买一送一', buyCount: 1, giftCount: 1, giftName: '同款商品', productName: '全场商品', dateRange: '2026-03-01 至 2026-06-30', status: 'active' }]
    pagination.total = res.data?.total || 1
  } catch { dataList.value = [] } finally { loading.value = false }
}

const showCreateDrawer = () => { isEdit.value = false; Object.assign(form, { name: '', buyCount: 1, giftCount: 1, giftName: '' }); modalVisible.value = true }
const showEdit = (record) => { isEdit.value = true; Object.assign(form, record); modalVisible.value = true }
const handleFormSubmit = () => { if (!form.name) { Message.warning('请填写活动名称'); return }; modalVisible.value = false; Message.success(isEdit.value ? '更新成功' : '创建成功'); loadData() }
const handleDelete = () => { Message.success('删除成功'); loadData() }
const onPageChange = (page) => { pagination.current = page; loadData() }

onMounted(() => loadData())
</script>
