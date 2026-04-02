<template>
    <Breadcrumb :items="['Home','Console','']" />


  <div class="container">
    <a-card class="general-card" title="优惠券发放">
      <template #extra>
        <a-space :size="12">
          <a-button type="primary" @click="showGrantDrawer"><icon-plus />发放优惠券</a-button>
          <a-button @click="loadData"><icon-refresh />刷新</a-button>
        </a-space>
      </template>
      <a-row :gutter="16">
        <a-col :span="8">
          <a-form-item label="优惠券名称">
            <a-input v-model="filters.keyword" placeholder="请输入" @pressEnter="loadData" />
          </a-form-item>
        </a-col>
        <a-col :span="8">
          <a-form-item label="发放方式">
            <a-select v-model="filters.mode" placeholder="请选择" allow-clear style="width: 100%">
              <a-option value="member">指定会员</a-option>
              <a-option value="level">指定等级</a-option>
              <a-option value="all">全部会员</a-option>
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
      <a-table :columns="columns" :data="dataList" :loading="loading" :pagination="paginationConfig" @page-change="onPageChange" row-key="id">
        <template #status="{ record }">
          <a-tag :color="record.status === 1 ? 'green' : 'gray'">{{ record.status === 1 ? '启用' : '禁用' }}</a-tag>
        </template>
      </a-table>
        <template #actions="{ record }">
          <a-button type="text" size="small" @click="showEdit(record)">编辑</a-button>
          <a-button type="text" size="small" @click="handleDelete(record)">删除</a-button>
        </template>
      </a-table>
    </a-card>
    <a-modal v-model:visible="modalVisible" :title="modalTitle">
      <a-form :model="form" label-col-flex="100px">
        <a-form-item label="名称"><a-input v-model="form.name" /></a-form-item>
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
const couponOptions = ref([])
const loading = ref(false)
const formLoading = ref(false)
const modalVisible = ref(false)
const isEdit = ref(false)
const filters = reactive({ keyword: '', mode: undefined })
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })
const paginationConfig = computed(() => ({ current: pagination.current, pageSize: pagination.pageSize, total: pagination.total, showTotal: true }))
const form = reactive({ couponId: undefined, mode: 'all', count: 1 })
const modalTitle = computed(() => isEdit.value ? '编辑' : '新建')
const columns = [
  { title: '优惠券名称', dataIndex: 'couponName', width: 180 },
  { title: '发放方式', slotName: 'mode', width: 120 },
  { title: '发放对象', slotName: 'target', width: 160 },
  { title: '发放数量', dataIndex: 'count', width: 100 },
  { title: '发放时间', dataIndex: 'createTime', width: 170 },
  { title: '状态', slotName: 'status', width: 90 }
]

const loadData = async () => {
  loading.value = true
  try {
    const params = { page: pagination.current, pageSize: pagination.pageSize }
    if (filters.keyword) params.keyword = filters.keyword
    if (filters.mode) params.mode = filters.mode
    const res = await fetch(`/api/v1/members/coupons/grant?${new URLSearchParams(params)}`, {
      headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token') }
    }).then(r => r.json())
    dataList.value = res.data?.list || []
    pagination.total = res.data?.total || 0
  } catch { dataList.value = [] } finally { loading.value = false }
}

const showGrantDrawer = () => { Object.assign(form, { couponId: undefined, mode: 'all', count: 1 }); modalVisible.value = true }
const showEdit = (record) => { isEdit.value = true; Object.assign(form, record); modalVisible.value = true }
const handleFormSubmit = () => { if (!form.name) { Message.warning('请填写名称'); return }; modalVisible.value = false; Message.success(isEdit.value ? '更新成功' : '创建成功'); loadData() }
const handleDelete = () => { Message.success('删除成功'); loadData() }
const onPageChange = (page) => { pagination.current = page; loadData() }

onMounted(() => loadData())
</script>
