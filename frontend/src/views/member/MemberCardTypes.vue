<template>
  <div class="member-card-types-page">
    <!-- 面包屑 -->
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>会员管理</a-breadcrumb-item>
      <a-breadcrumb-item>会员卡类型</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 搜索筛选 -->
    <a-card class="action-card">
      <a-space wrap>
        <a-input-search
          v-model="filters.keyword"
          placeholder="搜索卡类型名称/编码"
          style="width: 260px"
          search-button
          @search="handleSearch"
        />
        <a-select v-model="filters.status" placeholder="状态" allow-clear style="width: 120px" @change="handleSearch">
          <a-option :value="1">启用</a-option>
          <a-option :value="0">禁用</a-option>
        </a-select>
        <a-button @click="handleSearch">「搜索」</a-button>
        <a-button @click="resetFilters">「重置」</a-button>
      </a-space>
    </a-card>

    <!-- 操作栏 -->
    <a-card class="table-card">
      <template #extra>
        <a-space>
          <a-button type="primary" @click="showCreate">「新建」</a-button>
          <a-button @click="handleExport">「导出」</a-button>
          <a-button @click="loadData">🔄</a-button>
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
        :scroll="{ x: 1100 }"
      >
        <template #status="{ record }">
          <a-tag :color="record.status === 1 ? 'green' : 'gray'">{{ record.status === 1 ? '启用' : '禁用' }}</a-tag>
        </template>
        <template #discountRate="{ record }">
          <a-tag color="orange">{{ ((record.discountRate || 1) * 100).toFixed(0) }}%</a-tag>
        </template>
        <template #pointsRate="{ record }">
          <span style="color: #ff6b00; font-weight: 600;">{{ (record.pointsRate || 1).toFixed(1) }}x</span>
        </template>
        <template #validDays="{ record }">
          {{ record.validDays ? record.validDays + '天' : '永久' }}
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="showDetail(record)">详情</a-button>
            <a-button type="text" size="small" @click="showEdit(record)">编辑</a-button>
            <a-button type="text" size="small" status="danger" @click="handleDelete(record)">删除</a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 详情抽屉 -->
    <a-drawer v-model:visible="detailVisible" title="卡类型详情" :width="480">
      <template v-if="currentRecord">
        <a-descriptions :column="1" bordered size="small">
          <a-descriptions-item label="卡类型名称">{{ currentRecord.name }}</a-descriptions-item>
          <a-descriptions-item label="卡类型编码">{{ currentRecord.code }}</a-descriptions-item>
          <a-descriptions-item label="所属分组">{{ currentRecord.groupName || '-' }}</a-descriptions-item>
          <a-descriptions-item label="折扣率">
            <a-tag color="orange">{{ ((currentRecord.discountRate || 1) * 100).toFixed(0) }}%</a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="积分倍率">{{ (currentRecord.pointsRate || 1).toFixed(1) }}x</a-descriptions-item>
          <a-descriptions-item label="有效期">{{ currentRecord.validDays ? currentRecord.validDays + '天' : '永久' }}</a-descriptions-item>
          <a-descriptions-item label="状态">
            <a-tag :color="currentRecord.status === 1 ? 'green' : 'gray'">{{ currentRecord.status === 1 ? '启用' : '禁用' }}</a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="描述">{{ currentRecord.description || '-' }}</a-descriptions-item>
          <a-descriptions-item label="创建时间">{{ currentRecord.createdAt || '-' }}</a-descriptions-item>
        </a-descriptions>
      </template>
    </a-drawer>

    <!-- 新增/编辑弹窗 -->
    <a-modal
      v-model:visible="formVisible"
      :title="isEdit ? '编辑卡类型' : '新建卡类型'"
      @before-ok="handleFormSubmit"
      @cancel="formVisible = false"
      :width="560"
      :loading="formLoading"
    >
      <a-form :model="form" layout="vertical" ref="formRef">
        <a-form-item label="卡类型名称" field="name" :rules="[{ required: true, message: '请输入卡类型名称' }]">
          <a-input v-model="form.name" placeholder="如：黄金会员卡" />
        </a-form-item>
        <a-form-item label="卡类型编码" field="code" :rules="[{ required: true, message: '请输入卡类型编码' }]">
          <a-input v-model="form.code" placeholder="如：GOLD_CARD" :disabled="isEdit" />
        </a-form-item>
        <a-form-item label="所属分组" field="groupId">
          <a-select v-model="form.groupId" placeholder="选择所属分组" allow-clear>
            <a-option v-for="g in groupList" :key="g.id" :value="g.id">{{ g.name }}</a-option>
          </a-select>
        </a-form-item>
        <a-row :gutter="12">
          <a-col :span="12">
            <a-form-item label="折扣率" field="discountRate" :rules="[{ required: true, message: '请输入折扣率' }]">
              <a-input-number v-model="form.discountRate" :min="0" :max="1" :step="0.05" style="width: 100%">
                <template #suffix>如 0.95</template>
              </a-input-number>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="积分倍率" field="pointsRate" :rules="[{ required: true, message: '请输入积分倍率' }]">
              <a-input-number v-model="form.pointsRate" :min="0" :step="0.5" style="width: 100%">
                <template #suffix>x</template>
              </a-input-number>
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="12">
          <a-col :span="12">
            <a-form-item label="有效期(天)" field="validDays">
              <a-input-number v-model="form.validDays" :min="0" style="width: 100%" placeholder="0表示永久有效" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="状态" field="status">
              <a-radio-group v-model="form.status">
                <a-radio :value="1">启用</a-radio>
                <a-radio :value="0">禁用</a-radio>
              </a-radio-group>
            </a-form-item>
          </a-col>
        </a-row>
        <a-form-item label="描述">
          <a-textarea v-model="form.description" :rows="2" placeholder="简要描述该卡类型" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { Message, Modal } from '@arco-design/web-vue'
import * as api from '@/api/member'

const dataList = ref([])
const groupList = ref([])
const loading = ref(false)
const formLoading = ref(false)
const formVisible = ref(false)
const detailVisible = ref(false)
const isEdit = ref(false)
const currentId = ref(null)
const currentRecord = ref(null)

const filters = reactive({ keyword: '', status: undefined })
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })

const paginationConfig = computed(() => ({
  current: pagination.current, pageSize: pagination.pageSize, total: pagination.total,
  showTotal: true, showPageSize: true, pageSizeOptions: [10, 20, 50, 100]
}))

const form = reactive({
  name: '', code: '', groupId: undefined, discountRate: 1, pointsRate: 1,
  validDays: 0, status: 1, description: ''
})

const columns = [
  { title: '卡类型名称', dataIndex: 'name', width: 160 },
  { title: '卡类型编码', dataIndex: 'code', width: 140 },
  { title: '所属分组', dataIndex: 'groupName', width: 130 },
  { title: '折扣率', slotName: 'discountRate', width: 100 },
  { title: '积分倍率', slotName: 'pointsRate', width: 100 },
  { title: '有效期', slotName: 'validDays', width: 110 },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '创建时间', dataIndex: 'createdAt', width: 170 },
  { title: '操作', slotName: 'actions', width: 150, fixed: 'right' }
]

const loadData = async () => {
  loading.value = true
  try {
    const params = { page: pagination.current, pageSize: pagination.pageSize }
    if (filters.keyword) params.keyword = filters.keyword
    if (filters.status !== undefined) params.status = filters.status
    const res = await api.getCardTypeList(params)
    const d = res.data || {}
    dataList.value = d.list || []
    pagination.total = d.total || 0
  } catch (err) {
    Message.error('加载失败: ' + err.message)
  } finally {
    loading.value = false
  }
}

const loadGroups = async () => {
  try {
    const res = await api.getCardGroupList({ page: 1, pageSize: 100 })
    groupList.value = res.data?.list || []
  } catch (err) { /* ignore */ }
}

const handleSearch = () => { pagination.current = 1; loadData() }

const resetFilters = () => { filters.keyword = ''; filters.status = undefined; pagination.current = 1; loadData() }

const onPageChange = (page) => { pagination.current = page; loadData() }
const onPageSizeChange = (pageSize) => { pagination.pageSize = pageSize; pagination.current = 1; loadData() }

const showCreate = () => {
  isEdit.value = false
  Object.assign(form, { name: '', code: '', groupId: undefined, discountRate: 1, pointsRate: 1, validDays: 0, status: 1, description: '' })
  formVisible.value = true
}

const showEdit = (record) => {
  isEdit.value = true; currentId.value = record.id
  Object.assign(form, {
    name: record.name, code: record.code, groupId: record.groupId,
    discountRate: record.discountRate || 1, pointsRate: record.pointsRate || 1,
    validDays: record.validDays || 0, status: record.status ?? 1, description: record.description || ''
  })
  formVisible.value = true
}

const showDetail = (record) => { currentRecord.value = record; detailVisible.value = true }

const handleFormSubmit = async (done) => {
  if (!form.name || !form.code) { Message.warning('请填写名称和编码'); done(false); return }
  formLoading.value = true
  try {
    const payload = { ...form }
    if (isEdit.value) { await api.updateCardType(currentId.value, payload); Message.success('更新成功') }
    else { await api.createCardType(payload); Message.success('创建成功') }
    formVisible.value = false; loadData(); done(true)
  } catch (err) { Message.error(err.message || '操作失败'); done(false) }
  finally { formLoading.value = false }
}

const handleDelete = (record) => {
  Modal.warning({
    title: '确认删除', content: `确定要删除卡类型「${record.name}」吗？`,
    okText: '确认删除', onOk: async () => {
      try { await api.deleteCardType(record.id); Message.success('删除成功'); loadData() }
      catch (err) { Message.error(err.message || '删除失败') }
    }
  })
}

const handleExport = () => { Message.info('导出功能开发中') }

onMounted(() => { loadData(); loadGroups() })
</script>

<style scoped>
.member-card-types-page { padding: 20px 24px; min-height: calc(100vh - 64px); background: #f5f7fa; }
.breadcrumb { margin-bottom: 16px; }
.action-card { margin-bottom: 16px; }
.table-card { border-radius: 8px; }
</style>
