<template>
  <div class="promotion-page">
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>会员管理</a-breadcrumb-item>
      <a-breadcrumb-item>促销活动</a-breadcrumb-item>
      <a-breadcrumb-item>买赠促销</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 搜索筛选 -->
    <a-card class="action-card">
      <a-space wrap>
        <a-input-search v-model="filters.keyword" placeholder="搜索活动名称" style="width: 220px" search-button @search="loadData" />
        <a-select v-model="filters.status" placeholder="状态" allow-clear style="width: 120px" @change="loadData">
          <a-option value="active">进行中</a-option>
          <a-option value="pending">未开始</a-option>
          <a-option value="ended">已结束</a-option>
        </a-select>
        <a-button type="primary" @click="showCreateDrawer">新建活动</a-button>
        <a-button @click="loadData">刷新</a-button>
      </a-space>
    </a-card>

    <!-- 活动列表 -->
    <a-card class="table-card">
      <a-table
        :columns="columns"
        :data="dataList"
        :loading="loading"
        :pagination="paginationConfig"
        @page-change="onPageChange"
        @page-size-change="onPageSizeChange"
        row-key="id"
        :scroll="{ x: 1200 }"
      >
        <template #rule="{ record }">
          <span>买 <strong>{{ record.buyCount }}</strong> 件，送 <strong>{{ record.giftCount }}</strong> 件 {{ record.giftName }}</span>
        </template>
        <template #status="{ record }">
          <a-tag :color="getStatusColor(record.status)">{{ getStatusText(record.status) }}</a-tag>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="showEdit(record)">编辑</a-button>
            <a-button type="text" size="small" status="danger" @click="handleDelete(record)">删除</a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 新建/编辑活动抽屉 -->
    <a-drawer v-model:visible="formVisible" :title="isEdit ? '编辑买赠活动' : '新建买赠活动'" :width="560">
      <a-form :model="form" layout="vertical">
        <a-form-item label="活动名称" required>
          <a-input v-model="form.name" placeholder="请输入活动名称" />
        </a-form-item>
        <a-row :gutter="12">
          <a-col :span="12">
            <a-form-item label="购买数量">
              <a-input-number v-model="form.buyCount" :min="1" style="width: 100%" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="赠送数量">
              <a-input-number v-model="form.giftCount" :min="1" style="width: 100%" />
            </a-form-item>
          </a-col>
        </a-row>
        <a-form-item label="赠送商品">
          <a-input v-model="form.giftName" placeholder="请输入赠送商品名称" />
        </a-form-item>
        <a-form-item label="适用商品">
          <a-select v-model="form.productIds" placeholder="选择适用商品" multiple>
            <a-option v-for="p in productOptions" :key="p.id" :value="p.id">{{ p.name }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="活动时间">
          <a-range-picker v-model="form.dateRange" style="width: 100%" />
        </a-form-item>
        <a-form-item label="活动说明">
          <a-textarea v-model="form.description" :rows="3" placeholder="描述活动规则" />
        </a-form-item>
      </a-form>
      <template #footer>
        <a-button @click="formVisible = false">取消</a-button>
        <a-button type="primary" :loading="formLoading" @click="handleFormSubmit">{{ isEdit ? '保存' : '创建' }}</a-button>
      </template>
    </a-drawer>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const dataList = ref([])
const productOptions = ref([])
const loading = ref(false)
const formLoading = ref(false)
const formVisible = ref(false)
const isEdit = ref(false)
const currentRecord = ref(null)

const filters = reactive({ keyword: '', status: undefined })
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })

const paginationConfig = computed(() => ({
  current: pagination.current, pageSize: pagination.pageSize, total: pagination.total,
  showTotal: true, showPageSize: true, pageSizeOptions: [10, 20, 50, 100]
}))

const form = reactive({
  name: '', buyCount: 1, giftCount: 1, giftName: '',
  productIds: [], dateRange: [], description: ''
})

const columns = [
  { title: '活动名称', dataIndex: 'name', width: 200 },
  { title: '买赠规则', slotName: 'rule', width: 280 },
  { title: '适用商品', dataIndex: 'productName', width: 160, ellipsis: true },
  { title: '时间范围', dataIndex: 'dateRange', width: 220 },
  { title: '状态', slotName: 'status', width: 90 },
  { title: '操作', slotName: 'actions', width: 150 }
]

const getStatusColor = (s) => ({ active: 'green', pending: 'blue', ended: 'gray' }[s] || 'gray')
const getStatusText = (s) => ({ active: '进行中', pending: '未开始', ended: '已结束' }[s] || s)

const loadData = async () => {
  loading.value = true
  try {
    const d = {
      list: [
        { id: 1, name: '买一送一', buyCount: 1, giftCount: 1, giftName: '同款商品', productName: '全场商品', dateRange: '2026-03-01 至 2026-06-30', status: 'active' },
        { id: 2, name: '买三送一', buyCount: 3, giftCount: 1, giftName: '指定配件', productName: '主粮系列', dateRange: '2026-04-01 至 2026-05-31', status: 'active' }
      ],
      total: 2
    }
    dataList.value = d.list
    pagination.total = d.total
  } catch (err) {
    dataList.value = []
  } finally {
    loading.value = false
  }
}

const showCreateDrawer = () => {
  isEdit.value = false
  Object.assign(form, { name: '', buyCount: 1, giftCount: 1, giftName: '', productIds: [], dateRange: [], description: '' })
  formVisible.value = true
}

const showEdit = (record) => {
  isEdit.value = true
  currentRecord.value = record
  Object.assign(form, { name: record.name, buyCount: record.buyCount, giftCount: record.giftCount, giftName: record.giftName, productIds: [], dateRange: [], description: record.description || '' })
  formVisible.value = true
}

const handleFormSubmit = async () => {
  if (!form.name) { Message.warning('请填写活动名称'); return }
  formLoading.value = true
  try {
    await new Promise(r => setTimeout(r, 500))
    Message.success(isEdit.value ? '更新成功' : '创建成功')
    formVisible.value = false
    loadData()
  } catch (err) {
    Message.error(err.message || '操作失败')
  } finally {
    formLoading.value = false
  }
}

const handleDelete = async (record) => {
  try {
    await new Promise(r => setTimeout(r, 300))
    Message.success('删除成功')
    loadData()
  } catch (err) {
    Message.error(err.message || '删除失败')
  }
}

const onPageChange = (page) => { pagination.current = page; loadData() }
const onPageSizeChange = (pageSize) => { pagination.pageSize = pageSize; pagination.current = 1; loadData() }

onMounted(() => loadData())
</script>

<style scoped>
.promotion-page { padding: 20px 24px; min-height: calc(100vh - 64px); background: #f5f7fa; }
.breadcrumb { margin-bottom: 16px; }
.action-card { margin-bottom: 16px; }
.table-card { border-radius: 8px; }
</style>
