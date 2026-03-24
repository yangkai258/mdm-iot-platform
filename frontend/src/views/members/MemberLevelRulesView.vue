<template>
  <div class="level-rules-page">
    <!-- 面包屑 -->
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>会员管理</a-breadcrumb-item>
      <a-breadcrumb-item>会员升级规则</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 搜索筛选区 -->
    <a-card class="action-card">
      <a-space wrap>
        <a-input-search
          v-model="filters.keyword"
          placeholder="搜索等级名称"
          style="width: 220px"
          search-button
          @search="handleSearch"
        />
        <a-button @click="handleSearch">筛选</a-button>
        <a-button @click="resetFilters">重置</a-button>
      </a-space>
    </a-card>

    <!-- 操作+表格 -->
    <a-card class="table-card">
      <template #title>
        <a-space>
          <span style="font-weight: 600; font-size: 15px;">会员升级规则</span>
          <a-badge :count="pagination.total" :max-count="99999" />
        </a-space>
      </template>
      <template #extra>
        <a-space>
          <a-button type="primary" @click="showCreateModal">新建规则</a-button>
        </a-space>
      </template>

      <a-table
        :columns="columns"
        :data="dataList"
        :loading="loading"
        :pagination="paginationConfig"
        @page-change="onPageChange"
        @page-size-change="onPageSizeChange"
        row-key="id"
        :scroll="{ x: 1000 }"
      >
        <template #conditionType="{ record }">
          <a-tag :color="conditionColor(record.conditionType)">{{ getConditionLabel(record.conditionType) }}</a-tag>
        </template>
        <template #conditionValue="{ record }">
          <span v-if="record.conditionType === 'consume_amount'">消费满 ¥{{ record.conditionValue }}</span>
          <span v-else-if="record.conditionType === 'order_count'">消费 {{ record.conditionValue }} 次</span>
          <span v-else-if="record.conditionType === 'points'">累计 {{ record.conditionValue }} 积分</span>
          <span v-else>{{ record.conditionValue }}</span>
        </template>
        <template #requiredPoints="{ record }">
          <span style="color: #ff6b00;">{{ record.requiredPoints || 0 }}</span>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="showEdit(record)">编辑</a-button>
            <a-button type="text" size="small" status="danger" @click="handleDelete(record)">删除</a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 新增/编辑弹窗 -->
    <a-modal
      v-model:visible="formVisible"
      :title="isEdit ? '编辑升级规则' : '新建升级规则'"
      @before-ok="handleFormSubmit"
      @cancel="formVisible = false"
      :width="500"
      :loading="formLoading"
    >
      <a-form :model="form" layout="vertical">
        <a-form-item label="等级名称" field="name" :rules="[{ required: true, message: '请输入等级名称' }]">
          <a-input v-model="form.name" placeholder="如：黄金会员" />
        </a-form-item>
        <a-form-item label="升级条件类型" field="conditionType" :rules="[{ required: true, message: '请选择条件类型' }]">
          <a-select v-model="form.conditionType" placeholder="请选择条件类型">
            <a-option value="consume_amount">累计消费金额</a-option>
            <a-option value="order_count">累计消费次数</a-option>
            <a-option value="points">累计积分</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="条件阈值" field="conditionValue" :rules="[{ required: true, message: '请输入条件阈值' }]">
          <a-input-number v-model="form.conditionValue" :min="0" :max="99999999" style="width: 100%" placeholder="达到该值触发升级" />
        </a-form-item>
        <a-form-item label="所需积分">
          <a-input-number v-model="form.requiredPoints" :min="0" :max="99999999" style="width: 100%" placeholder="升级所需积分，为0表示无需积分" />
        </a-form-item>
        <a-form-item label="描述">
          <a-textarea v-model="form.description" :rows="2" placeholder="规则描述" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import * as api from '@/api/member'

const dataList = ref([])
const loading = ref(false)
const formLoading = ref(false)
const formVisible = ref(false)
const isEdit = ref(false)
const currentId = ref(null)

const filters = reactive({ keyword: '' })
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })

const paginationConfig = computed(() => ({
  current: pagination.current,
  pageSize: pagination.pageSize,
  total: pagination.total,
  showTotal: true,
  showPageSize: true,
  pageSizeOptions: [10, 20, 50, 100]
}))

const form = reactive({
  name: '',
  conditionType: 'consume_amount',
  conditionValue: 0,
  requiredPoints: 0,
  description: ''
})

const columns = [
  { title: 'ID', dataIndex: 'id', width: 70 },
  { title: '等级名称', dataIndex: 'name', width: 160 },
  { title: '升级条件类型', slotName: 'conditionType', width: 150 },
  { title: '升级条件值', slotName: 'conditionValue', width: 160 },
  { title: '所需积分', slotName: 'requiredPoints', width: 120 },
  { title: '描述', dataIndex: 'description', ellipsis: true },
  { title: '操作', slotName: 'actions', width: 120, fixed: 'right' }
]

const getConditionLabel = (type) => ({
  consume_amount: '累计消费金额',
  order_count: '累计消费次数',
  points: '累计积分'
}[type] || type)

const conditionColor = (type) => ({
  consume_amount: 'blue',
  order_count: 'green',
  points: 'orange'
}[type] || 'gray')

const loadData = async () => {
  loading.value = true
  try {
    const params = { page: pagination.current, pageSize: pagination.pageSize }
    if (filters.keyword) params.keyword = filters.keyword

    const res = await api.getUpgradeRules()
    const d = res.data || {}
    dataList.value = Array.isArray(d) ? d : (d.list || [])
    pagination.total = dataList.value.length
  } catch (err) {
    Message.error('加载失败: ' + err.message)
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.current = 1
  loadData()
}

const resetFilters = () => {
  filters.keyword = ''
  pagination.current = 1
  loadData()
}

const onPageChange = (page) => {
  pagination.current = page
  loadData()
}

const onPageSizeChange = (pageSize) => {
  pagination.pageSize = pageSize
  pagination.current = 1
  loadData()
}

const showCreateModal = () => {
  isEdit.value = false
  Object.assign(form, { name: '', conditionType: 'consume_amount', conditionValue: 0, requiredPoints: 0, description: '' })
  formVisible.value = true
}

const showEdit = (record) => {
  isEdit.value = true
  currentId.value = record.id
  Object.assign(form, {
    name: record.name,
    conditionType: record.conditionType || 'consume_amount',
    conditionValue: record.conditionValue ?? 0,
    requiredPoints: record.requiredPoints ?? 0,
    description: record.description || ''
  })
  formVisible.value = true
}

const handleFormSubmit = async (done) => {
  formLoading.value = true
  try {
    const payload = { ...form }
    await api.updateUpgradeRules(payload)
    Message.success('保存成功')
    formVisible.value = false
    loadData()
    done(true)
  } catch (err) {
    Message.error(err.message || '操作失败')
    done(false)
  } finally {
    formLoading.value = false
  }
}

const handleDelete = async (record) => {
  try {
    // 降级规则没有单独的 delete，先用 update 模拟
    Message.info('请通过编辑调整规则')
  } catch (err) {
    Message.error(err.message || '操作失败')
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.level-rules-page {
  padding: 20px 24px;
  min-height: calc(100vh - 64px);
  background: #f5f7fa;
}
.breadcrumb { margin-bottom: 16px; }
.action-card { margin-bottom: 16px; }
.table-card { border-radius: 8px; }
</style>
