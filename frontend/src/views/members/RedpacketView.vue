<template>
  <div class="redpacket-page">
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>会员管理</a-breadcrumb-item>
      <a-breadcrumb-item>会员红包</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 统计卡片 -->
    <a-row :gutter="16" class="stats-row">
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="红包总数" :value="stats.total || 0" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="已发放" :value="stats.issued || 0" :value-style="{ color: '#1890ff' }" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="已使用" :value="stats.used || 0" :value-style="{ color: '#52c41a' }" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="待使用" :value="stats.unused || 0" :value-style="{ color: '#faad14' }" />
        </a-card>
      </a-col>
    </a-row>

    <!-- 搜索筛选 -->
    <a-card class="action-card">
      <a-space wrap>
        <a-input-search v-model="filters.keyword" placeholder="搜索红包名称" style="width: 220px" search-button @search="loadData" />
        <a-select v-model="filters.status" placeholder="状态" allow-clear style="width: 120px" @change="loadData">
          <a-option value="active">有效</a-option>
          <a-option value="inactive">未激活</a-option>
          <a-option value="expired">已过期</a-option>
        </a-select>
        <a-button type="primary" @click="showCreateDrawer">新建红包</a-button>
        <a-button @click="loadData">刷新</a-button>
      </a-space>
    </a-card>

    <!-- 红包列表 -->
    <a-card class="table-card">
      <a-table
        :columns="columns"
        :data="dataList"
        :loading="loading"
        :pagination="paginationConfig"
        @page-change="onPageChange"
        @page-size-change="onPageSizeChange"
        row-key="id"
        :scroll="{ x: 1100 }"
      >
        <template #amount="{ record }">
          <span style="color: #ff6b00; font-weight: 600;">¥{{ record.amount || 0 }}</span>
        </template>
        <template #status="{ record }">
          <a-tag :color="getStatusColor(record.status)">{{ getStatusText(record.status) }}</a-tag>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="showGrant(record)">发放</a-button>
            <a-button type="text" size="small" @click="showEdit(record)">编辑</a-button>
            <a-button type="text" size="small" status="danger" @click="handleDelete(record)">删除</a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 新建/编辑红包抽屉 -->
    <a-drawer v-model:visible="formVisible" :title="isEdit ? '编辑红包' : '新建红包'" :width="520">
      <a-form :model="form" layout="vertical">
        <a-form-item label="红包名称" required>
          <a-input v-model="form.name" placeholder="请输入红包名称" />
        </a-form-item>
        <a-row :gutter="12">
          <a-col :span="12">
            <a-form-item label="红包金额">
              <a-input-number v-model="form.amount" :min="0" :precision="2" style="width: 100%" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="发行总量">
              <a-input-number v-model="form.totalCount" :min="0" style="width: 100%" placeholder="0=不限量" />
            </a-form-item>
          </a-col>
        </a-row>
        <a-form-item label="有效期">
          <a-range-picker v-model="form.dateRange" style="width: 100%" />
        </a-form-item>
        <a-form-item label="使用说明">
          <a-textarea v-model="form.description" :rows="3" placeholder="描述红包使用规则" />
        </a-form-item>
      </a-form>
      <template #footer>
        <a-button @click="formVisible = false">取消</a-button>
        <a-button type="primary" :loading="formLoading" @click="handleFormSubmit">{{ isEdit ? '保存' : '创建' }}</a-button>
      </template>
    </a-drawer>

    <!-- 发放红包抽屉 -->
    <a-drawer v-model:visible="grantVisible" title="发放红包" :width="520">
      <a-form layout="vertical">
        <a-form-item label="红包">
          <a-input :value="currentRedpacket?.name" disabled />
        </a-form-item>
        <a-form-item label="发放方式" required>
          <a-radio-group v-model="grantForm.mode">
            <a-radio value="member">指定会员</a-radio>
            <a-radio value="level">指定等级</a-radio>
            <a-radio value="all">全部会员</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item v-if="grantForm.mode === 'member'" label="选择会员">
          <a-select v-model="grantForm.memberIds" placeholder="选择会员" multiple searchable>
            <a-option v-for="m in memberOptions" :key="m.id" :value="m.id">{{ m.name }} ({{ m.mobile }})</a-option>
          </a-select>
        </a-form-item>
        <a-form-item v-if="grantForm.mode === 'level'" label="选择等级">
          <a-select v-model="grantForm.levelId" placeholder="选择会员等级">
            <a-option v-for="lv in levelOptions" :key="lv.id" :value="lv.id">{{ lv.name }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="发放数量">
          <a-input-number v-model="grantForm.count" :min="1" :max="10" style="width: 100%" />
        </a-form-item>
      </a-form>
      <template #footer>
        <a-button @click="grantVisible = false">取消</a-button>
        <a-button type="primary" @click="handleGrant">确认发放</a-button>
      </template>
    </a-drawer>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const dataList = ref([])
const memberOptions = ref([])
const levelOptions = ref([])
const loading = ref(false)
const formLoading = ref(false)
const formVisible = ref(false)
const grantVisible = ref(false)
const isEdit = ref(false)
const currentRedpacket = ref(null)

const filters = reactive({ keyword: '', status: undefined })
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })
const stats = reactive({ total: 0, issued: 0, used: 0, unused: 0 })

const paginationConfig = computed(() => ({
  current: pagination.current, pageSize: pagination.pageSize, total: pagination.total,
  showTotal: true, showPageSize: true, pageSizeOptions: [10, 20, 50, 100]
}))

const form = reactive({ name: '', amount: 0, totalCount: 0, dateRange: [], description: '' })

const grantForm = reactive({ mode: 'all', memberIds: [], levelId: undefined, count: 1 })

const columns = [
  { title: '红包名称', dataIndex: 'name', width: 180 },
  { title: '金额', slotName: 'amount', width: 110 },
  { title: '发行总量', dataIndex: 'totalCount', width: 110 },
  { title: '有效期至', dataIndex: 'endTime', width: 170 },
  { title: '状态', slotName: 'status', width: 90 },
  { title: '操作', slotName: 'actions', width: 200 }
]

const getStatusColor = (s) => ({ active: 'green', inactive: 'gray', expired: 'orange' }[s] || 'gray')
const getStatusText = (s) => ({ active: '有效', inactive: '未激活', expired: '已过期' }[s] || s)

const loadData = async () => {
  loading.value = true
  try {
    const params = { page: pagination.current, pageSize: pagination.pageSize }
    if (filters.keyword) params.keyword = filters.keyword
    if (filters.status) params.status = filters.status
    const res = await mockApi()
    const d = res.data || {}
    dataList.value = d.list || []
    pagination.total = d.total || 0
    stats.total = d.total || 0
  } catch (err) {
    Message.error('加载红包列表失败: ' + err.message)
  } finally {
    loading.value = false
  }
}

const mockApi = () => Promise.resolve({
  data: {
    list: [
      { id: 1, name: '新人红包', amount: 10, totalCount: 1000, usedCount: 300, endTime: '2026-12-31', status: 'active' },
      { id: 2, name: '节日红包', amount: 20, totalCount: 500, usedCount: 150, endTime: '2026-06-30', status: 'active' }
    ],
    total: 2
  }
})

const showCreateDrawer = () => {
  isEdit.value = false
  Object.assign(form, { name: '', amount: 0, totalCount: 0, dateRange: [], description: '' })
  formVisible.value = true
}

const showEdit = (record) => {
  isEdit.value = true
  currentRedpacket.value = record
  Object.assign(form, { name: record.name, amount: record.amount, totalCount: record.totalCount, dateRange: [], description: record.description || '' })
  formVisible.value = true
}

const showGrant = (record) => {
  currentRedpacket.value = record
  loadMemberOptions()
  loadLevelOptions()
  Object.assign(grantForm, { mode: 'all', memberIds: [], levelId: undefined, count: 1 })
  grantVisible.value = true
}

const loadMemberOptions = async () => {
  try {
    const res = await import('@/api/member').then(m => m.getMemberList({ page: 1, pageSize: 100 }))
    memberOptions.value = res.data?.list || []
  } catch (err) { /* ignore */ }
}

const loadLevelOptions = async () => {
  try {
    const res = await import('@/api/member').then(m => m.getLevelList())
    levelOptions.value = res.data || []
  } catch (err) { /* ignore */ }
}

const handleFormSubmit = async () => {
  if (!form.name) { Message.warning('请填写名称'); return }
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

const handleGrant = async () => {
  try {
    await new Promise(r => setTimeout(r, 500))
    Message.success('发放成功')
    grantVisible.value = false
  } catch (err) {
    Message.error(err.message || '发放失败')
  }
}

const onPageChange = (page) => { pagination.current = page; loadData() }
const onPageSizeChange = (pageSize) => { pagination.pageSize