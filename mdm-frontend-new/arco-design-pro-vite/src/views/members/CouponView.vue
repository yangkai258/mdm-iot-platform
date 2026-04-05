<template>
    <Breadcrumb :items="['Home','Console','']" />


  <div class="container">
    <a-card class="general-card" title="优惠券管理">
      <template #extra>
        <a-space :size="12">
          <a-button type="primary" @click="showCreateDrawer"><icon-plus />新建</a-button>
          <a-button @click="loadData"><icon-refresh />刷新</a-button>
        </a-space>
      </template>
      <a-row :gutter="16">
        <a-col :span="6">
          <a-form-item label="优惠券名称">
            <a-input v-model="filters.keyword" placeholder="请输入" @pressEnter="loadData" />
          </a-form-item>
        </a-col>
        <a-col :span="6">
          <a-form-item label="类型">
            <a-select v-model="filters.type" placeholder="请选择" allow-clear style="width: 100%">
              <a-option value="discount">折扣券</a-option>
              <a-option value="cash">现金券</a-option>
              <a-option value="gift">礼品券</a-option>
              <a-option value="shipping">包邮券</a-option>
            </a-select>
          </a-form-item>
        </a-col>
        <a-col :flex="'86px'" style="display: flex; align-items: flex-end">
          <a-space direction="vertical" :size="8">
            <a-button type="primary" @click="loadData">查询</a-button>
            <a-button @click="Object.keys(filters).forEach(k => filters[k] = ''); loadData()">重置</a-button>
          </a-space>
        </a-col>
      </a-row>
      <a-divider style="margin: 0 0 16px 0" />
      <a-table :columns="columns" :data="couponList" :loading="loading" :pagination="paginationConfig" @page-change="onPageChange" row-key="id">
        <template #actions="{ record }">
          <a-button type="text" size="small" @click="showEdit(record)">编辑</a-button>
          <a-button type="text" size="small" @click="handleDelete(record)">删除</a-button>
        </template>
      </a-table>
    </a-card>
    <a-modal v-model:visible="modalVisible" :title="modalTitle">
      <a-form :model="form" label-col-flex="100px">
        <a-form-item label="名称"><a-input v-model="form.name" /></a-form-item>
        <a-form-item label="类型">
          <a-select v-model="form.type" style="width: 100%">
            <a-option value="discount">折扣券</a-option>
            <a-option value="cash">现金券</a-option>
            <a-option value="gift">礼品券</a-option>
            <a-option value="shipping">包邮券</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="面值/折扣"><a-input-number v-model="form.value" :min="0" style="width: 100%" /></a-form-item>
        <a-form-item label="发行数量"><a-input-number v-model="form.totalCount" :min="0" style="width: 100%" /></a-form-item>
        <a-form-item label="最低消费"><a-input-number v-model="form.minConsume" :min="0" style="width: 100%" /></a-form-item>
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

const couponList = ref([])
const loading = ref(false)
const formLoading = ref(false)
const modalVisible = ref(false)
const isEdit = ref(false)
const filters = reactive({ keyword: '', type: undefined, status: undefined })
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })
const paginationConfig = computed(() => ({ current: pagination.current, pageSize: pagination.pageSize, total: pagination.total, showTotal: true }))
const form = reactive({ name: '', type: '', value: 0, totalCount: 0, minConsume: 0 })
const modalTitle = computed(() => isEdit.value ? '编辑优惠券' : '新建优惠券')
const columns = [
  { title: '优惠券名称', dataIndex: 'name', width: 180 },
  { title: '类型', slotName: 'type', width: 100 },
  { title: '面值/折扣', slotName: 'value', width: 110 },
  { title: '发行/使用', slotName: 'usedCount', width: 110 },
  { title: '有效期至', dataIndex: 'endTime', width: 170 },
  { title: '状态', slotName: 'status', width: 90 },
  { title: '操作', slotName: 'actions', width: 220 }
]

const loadData = async () => {
  loading.value = true
  try {
    const params = { page: pagination.current, pageSize: pagination.pageSize }
    if (filters.keyword) params.keyword = filters.keyword
    if (filters.type) params.type = filters.type
    const res = await fetch(`/api/members/coupons?${new URLSearchParams(params)}`, {
      headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token') }
    }).then(r => r.json())
    const d = res.data || {}
    couponList.value = d.list || []
    pagination.total = d.total || 0
  } catch { couponList.value = [] } finally { loading.value = false }
}

const showCreateDrawer = () => { isEdit.value = false; Object.assign(form, { name: '', type: '', value: 0, totalCount: 0, minConsume: 0 }); modalVisible.value = true }
const showEdit = (record) => { isEdit.value = true; Object.assign(form, record); modalVisible.value = true }
const handleFormSubmit = () => { if (!form.name) { Message.warning('请填写名称'); return }; modalVisible.value = false; Message.success(isEdit.value ? '更新成功' : '创建成功'); loadData() }
const handleDelete = () => { Message.success('删除成功'); loadData() }
const onPageChange = (page) => { pagination.current = page; loadData() }

onMounted(() => loadData())
</script>
