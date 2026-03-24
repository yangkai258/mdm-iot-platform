<template>
  <div class="promotion-page">
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>会员管理</a-breadcrumb-item>
      <a-breadcrumb-item>促销活动</a-breadcrumb-item>
      <a-breadcrumb-item>最高等级促销</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 搜索筛选 -->
    <a-card class="action-card">
      <a-space wrap>
        <a-input-search v-model="filters.keyword" placeholder="搜索活动名称" style="width: 220px" search-button @search="loadData" />
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
        :scroll="{ x: 900 }"
      >
        <template #level="{ record }">
          <a-tag color="gold">{{ record.levelName }}</a-tag>
        </template>
        <template #discount="{ record }">
          <span v-if="record.discountType === 'percent'" style="color: #1890ff; font-weight: 600;">{{ (record.discountValue * 100).toFixed(0) }}折</span>
          <span v-else style="color: #ff6b00; font-weight: 600;">¥{{ record.discountValue }}</span>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="showEdit(record)">编辑</a-button>
            <a-button type="text" size="small" status="danger" @click="handleDelete(record)">删除</a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 新建/编辑抽屉 -->
    <a-drawer v-model:visible="formVisible" :title="isEdit ? '编辑活动' : '新建最高等级促销'" :width="560">
      <a-form :model="form" layout="vertical">
        <a-form-item label="活动名称" required>
          <a-input v-model="form.name" placeholder="请输入活动名称" />
        </a-form-item>
        <a-form-item label="适用等级" required>
          <a-select v-model="form.levelId" placeholder="选择会员等级">
            <a-option v-for="lv in levelOptions" :key="lv.id" :value="lv.id">{{ lv.name }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="优惠类型">
          <a-radio-group v-model="form.discountType">
            <a-radio value="percent">折扣</a-radio>
            <a-radio value="amount">金额</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item label="优惠内容">
          <a-input-number v-if="form.discountType === 'percent'" v-model="form.discountValue" :min="0" :max="1" :step="0.1" :precision="2" style="width: 100%" />
          <a-input-number v-else v-model="form.discountValue" :min="0" :precision="2" style="width: 100%" />
        </a-form-item>
        <a-form-item label="适用商品">
          <a-select v-model="form.productIds" placeholder="选择适用商品" multiple>
            <a-option v-for="p in productOptions" :key="p.id" :value="p.id">{{ p.name }}</a-option>
          </a-select>
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
import * as api from '@/api/member'

const dataList = ref([])
const levelOptions = ref([])
const productOptions = ref([])
const loading = ref(false)
const formLoading = ref(false)
const formVisible = ref(false)
const isEdit = ref(false)

const filters = reactive({ keyword: '' })
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })

const paginationConfig = computed(() => ({
  current: pagination.current, pageSize: pagination.pageSize, total: pagination.total,
  showTotal: true, showPageSize: true, pageSizeOptions: [10, 20, 50, 100]
}))

const form = reactive({
  name: '', levelId: undefined, discountType: 'percent',
  discountValue: 1, productIds: [], description: ''
})

const columns = [
  { title: '活动名称', dataIndex: 'name', width: 200 },
  { title: '适用等级', slotName: 'level', width: 140 },
  { title: '优惠内容', slotName: 'discount', width: 140 },
  { title: '适用商品', dataIndex: 'productName', width: 180, ellipsis: true },
  { title: '操作', slotName: 'actions', width: 150 }
]

const loadData = async () => {
  loading.value = true
  try {
    const d = {
      list: [
        { id: 1, name: 'VIP专属折扣', levelName: '黄金会员', discountType: 'percent', discountValue: 0.85, productName: '全场商品' },
        { id: 2, name: '钻石立减', levelName: '钻石会员', discountType: 'amount', discountValue: 50, productName: '指定商品' }
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

const loadLevelOptions = async () => {
  try {
    const res = await api.getLevelList()
    levelOptions.value = res.data || []
  } catch (err) { /* ignore */ }
}

const showCreateDrawer = () => {
  isEdit.value = false
  loadLevelOptions()
  Object.assign(form, { name: '', levelId: undefined, discountType: 'percent', discountValue: 1, productIds: [], description: '' })
  formVisible.value = true
}

const showEdit = (record) => {
  isEdit.value = true
  loadLevelOptions()
  Object.assign(form, { name: record.name, levelId: record.levelId, discountType: record.discountType, discountValue: record.discountValue, productIds: [], description: record.description || '' })
  formVisible.value = true
}

const handleFormSubmit = async () => {
  if (!form.name || !form.levelId) { Message.warning('请填写名称和选择等级'); return }
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
