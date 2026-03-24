<template>
  <div class="member-card-page">
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>会员管理</a-breadcrumb-item>
      <a-breadcrumb-item>会员卡管理</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 统计卡片 -->
    <a-row :gutter="16" class="stats-row">
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="卡类型总数" :value="stats.totalTypes || 0" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="已发行卡数" :value="stats.issuedCards || 0" :value-style="{ color: '#1890ff' }" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="活跃卡片" :value="stats.activeCards || 0" :value-style="{ color: '#52c41a' }" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="今日发行" :value="stats.todayIssued || 0" :value-style="{ color: '#faad14' }" />
        </a-card>
      </a-col>
    </a-row>

    <!-- 搜索筛选 -->
    <a-card class="action-card">
      <a-space wrap>
        <a-input-search v-model="filters.keyword" placeholder="搜索卡类型名称/编码" style="width: 240px" search-button @search="loadCardTypes" />
        <a-select v-model="filters.status" placeholder="状态" allow-clear style="width: 120px" @change="loadCardTypes">
          <a-option :value="1">启用</a-option>
          <a-option :value="0">禁用</a-option>
        </a-select>
        <a-button type="primary" @click="showCreateType">新建卡类型</a-button>
        <a-button @click="loadCardTypes">刷新</a-button>
      </a-space>
    </a-card>

    <!-- 卡类型列表 -->
    <a-card class="table-card">
      <template #title><span style="font-weight: 600;">卡类型列表</span></template>
      <a-table
        :columns="typeColumns"
        :data="cardTypeList"
        :loading="loading"
        :pagination="paginationConfig"
        @page-change="onPageChange"
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
        <template #cardFee="{ record }">
          {{ record.cardFee > 0 ? '¥' + record.cardFee.toFixed(2) : '免费' }}
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="showEditType(record)">编辑</a-button>
            <a-button type="text" size="small" @click="showIssueCard(record)">发行</a-button>
            <a-button type="text" size="small" status="danger" @click="handleDeleteType(record)">删除</a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 发行记录 -->
    <a-card class="table-card" style="margin-top: 16px;">
      <template #title><span style="font-weight: 600;">发行记录</span></template>
      <template #extra>
        <a-button type="primary" size="small" @click="showBatchIssue">批量发行</a-button>
      </template>
      <a-table
        :columns="issueColumns"
        :data="issueRecords"
        :loading="issueLoading"
        :pagination="issuePaginationConfig"
        @page-change="onIssuePageChange"
        row-key="id"
      >
        <template #cardType="{ record }">
          <a-tag>{{ record.cardTypeName || '-' }}</a-tag>
        </template>
        <template #status="{ record }">
          <a-tag :color="record.status === 1 ? 'green' : 'gray'">{{ record.status === 1 ? '正常' : '失效' }}</a-tag>
        </template>
      </a-table>
    </a-card>

    <!-- 新建/编辑卡类型弹窗 -->
    <a-modal
      v-model:visible="typeModalVisible"
      :title="isEditType ? '编辑卡类型' : '新建卡类型'"
      @before-ok="handleTypeSubmit"
      @cancel="typeModalVisible = false"
      :width="560"
      :loading="formLoading"
    >
      <a-form :model="typeForm" layout="vertical">
        <a-form-item label="卡类型名称" required>
          <a-input v-model="typeForm.name" placeholder="如：黄金会员卡" />
        </a-form-item>
        <a-form-item label="卡类型编码" required>
          <a-input v-model="typeForm.code" placeholder="如：GOLD_CARD" :disabled="isEditType" />
        </a-form-item>
        <a-row :gutter="12">
          <a-col :span="12">
            <a-form-item label="折扣率">
              <a-input-number v-model="typeForm.discountRate" :min="0" :max="1" :step="0.05" style="width: 100%">
                <template #suffix>如 0.95</template>
              </a-input-number>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="积分倍率">
              <a-input-number v-model="typeForm.pointsRate" :min="0" :step="0.5" style="width: 100%">
                <template #suffix>x</template>
              </a-input-number>
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="12">
          <a-col :span="12">
            <a-form-item label="有效期(天)">
              <a-input-number v-model="typeForm.validDays" :min="0" style="width: 100%" placeholder="0=永久" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="工本费(元)">
              <a-input-number v-model="typeForm.cardFee" :min="0" :precision="2" style="width: 100%" />
            </a-form-item>
          </a-col>
        </a-row>
        <a-form-item label="卡面颜色">
          <a-input v-model="typeForm.color" placeholder="#1890ff" />
        </a-form-item>
        <a-form-item label="权益描述">
          <a-textarea v-model="typeForm.benefits" :rows="3" placeholder="描述该卡类型的权益" />
        </a-form-item>
        <a-form-item label="状态">
          <a-switch v-model="typeFormStatus" checked-value="1" unchecked-value="0" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 发行卡片抽屉 -->
    <a-drawer v-model:visible="issueModalVisible" title="发行会员卡" :width="420">
      <a-form :model="issueForm" layout="vertical">
        <a-form-item label="卡类型">
          <a-input :value="currentCardType?.name" disabled />
        </a-form-item>
        <a-form-item label="会员" required>
          <a-select v-model="issueForm.memberId" placeholder="选择会员" searchable>
            <a-option v-for="m in memberOptions" :key="m.id" :value="m.id">{{ m.name }} ({{ m.mobile }})</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="备注">
          <a-textarea v-model="issueForm.remark" :rows="2" />
        </a-form-item>
      </a-form>
      <template #footer>
        <a-button @click="issueModalVisible = false">取消</a-button>
        <a-button type="primary" @click="handleIssueCard">确认发行</a-button>
      </template>
    </a-drawer>

    <!-- 批量发行弹窗 -->
    <a-modal v-model:visible="batchIssueVisible" title="批量发行会员卡" :width="480">
      <a-form :model="batchForm" layout="vertical">
        <a-form-item label="卡类型" required>
          <a-select v-model="batchForm.cardTypeId" placeholder="选择卡类型">
            <a-option v-for="t in cardTypeList" :key="t.id" :value="t.id">{{ t.name }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="发行数量" required>
          <a-input-number v-model="batchForm.count" :min="1" :max="1000" style="width: 100%" />
        </a-form-item>
        <a-form-item label="备注">
          <a-textarea v-model="batchForm.remark" :rows="2" />
        </a-form-item>
      </a-form>
      <template #footer>
        <a-button @click="batchIssueVisible = false">取消</a-button>
        <a-button type="primary" @click="handleBatchIssue">确认发行</a-button>
      </template>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { Message, Modal } from '@arco-design/web-vue'
import * as api from '@/api/member'

const cardTypeList = ref([])
const issueRecords = ref([])
const memberOptions = ref([])
const loading = ref(false)
const issueLoading = ref(false)
const formLoading = ref(false)
const typeModalVisible = ref(false)
const issueModalVisible = ref(false)
const batchIssueVisible = ref(false)
const isEditType = ref(false)
const currentId = ref(null)
const currentCardType = ref(null)

const filters = reactive({ keyword: '', status: undefined })
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })
const issuePagination = reactive({ current: 1, pageSize: 20, total: 0 })
const stats = reactive({ totalTypes: 0, issuedCards: 0, activeCards: 0, todayIssued: 0 })

const paginationConfig = computed(() => ({
  current: pagination.current, pageSize: pagination.pageSize, total: pagination.total,
  showTotal: true, showPageSize: true, pageSizeOptions: [10, 20, 50, 100]
}))
const issuePaginationConfig = computed(() => ({
  current: issuePagination.current, pageSize: issuePagination.pageSize, total: issuePagination.total,
  showTotal: true, showPageSize: true
}))

const typeForm = reactive({
  name: '', code: '', discountRate: 1, pointsRate: 1, validDays: 0,
  cardFee: 0, color: '#1890ff', benefits: ''
})
const typeFormStatus = ref('1')
const issueForm = reactive({ memberId: undefined, remark: '' })
const batchForm = reactive({ cardTypeId: undefined, count: 1, remark: '' })

const typeColumns = [
  { title: '卡类型名称', dataIndex: 'name', width: 180 },
  { title: '编码', dataIndex: 'code', width: 120 },
  { title: '折扣率', slotName: 'discountRate', width: 100 },
  { title: '积分倍率', slotName: 'pointsRate', width: 110 },
  { title: '有效期', slotName: 'validDays', width: 100 },
  { title: '工本费', slotName: 'cardFee', width: 100 },
  { title: '卡面颜色', dataIndex: 'color', width: 100 },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '操作', slotName: 'actions', width: 180 }
]

const issueColumns = [
  { title: '卡号', dataIndex: 'cardNo', width: 200 },
  { title: '卡类型', slotName: 'cardType', width: 150 },
  { title: '会员', dataIndex: 'memberName', width: 120 },
  { title: '手机号', dataIndex: 'mobile', width: 130 },
  { title: '状态', slotName: 'status', width: 90 },
  { title: '发行时间', dataIndex: 'createdAt', width: 170 }
]

const loadCardTypes = async () => {
  loading.value = true
  try {
    const params = { page: pagination.current, pageSize: pagination.pageSize }
    if (filters.keyword) params.keyword = filters.keyword
    if (filters.status !== undefined) params.status = filters.status
    const res = await api.getCardTypeList(params)
    const d = res.data || {}
    cardTypeList.value = d.list || []
    pagination.total = d.total || 0
    stats.totalTypes = d.total || 0
  } catch (err) {
    Message.error('加载卡类型失败: ' + err.message)
  } finally {
    loading.value = false
  }
}

const loadIssueRecords = async () => {
  issueLoading.value = true
  try {
    const params = { page: issuePagination.current, pageSize: issuePagination.pageSize }
    const res = await api.getCardGroupList(params)
    const d = res.data || {}
    issueRecords.value = d.list || []
    issuePagination.total = d.total || 0
  } catch (err) {
    issueRecords.value = []
  } finally {
    issueLoading.value = false
  }
}

const loadMemberOptions = async () => {
  try {
    const res = await api.getMemberList({ page: 1, pageSize: 100 })
    memberOptions.value = res.data?.list || []
  } catch (err) { /* ignore */ }
}

const showCreateType = () => {
  isEditType.value = false
  currentId.value = null
  Object.assign(typeForm, { name: '', code: '', discountRate: 1, pointsRate: 1, validDays: 0, cardFee: 0, color: '#1890ff', benefits: '' })
  typeFormStatus.value = '1'
  typeModalVisible.value = true
}

const showEditType = (record) => {
  isEditType.value = true
  currentId.value = record.id
  Object.assign(typeForm, {
    name: record.name,
    code: record.code,
    discountRate: record.discountRate || 1,
    pointsRate: record.pointsRate || 1,
    validDays: record.validDays || 0,
    cardFee: record.cardFee || 0,
    color: record.color || '#1890ff',
    benefits: record.benefits || ''
  })
  typeFormStatus.value = String(record.status || 1)
  typeModalVisible.value = true
}

const handleTypeSubmit = async (done) => {
  if (!typeForm.name || !typeForm.code) {
    Message.warning('请填写名称和编码')
    done(false)
    return
  }
  formLoading.value = true
  try {
    const payload = { ...typeForm, status: parseInt(typeFormStatus.value) }
    if (isEditType.value) {
      await api.updateCardType(currentId.value, payload)
      Message.success('更新成功')
    } else {
      await api.createCardType(payload)
      Message.success('创建成功')
    }
    typeModalVisible.value = false
    loadCardTypes()
    done(true)
  } catch (err) {
    Message.error(err.message || '操作失败')
    done(false)
  } finally {
    formLoading.value = false
  }
}

const handleDeleteType = (record) => {
  Modal.warning({
    title: '确认删除',
    content: `确定要删除卡类型「${record.name}」吗？`,
    okText: '确认删除',
    onOk: async () => {
      try {
        await api.deleteCardType(record.id)
        Message.success('删除成功')
        loadCardTypes()
      } catch (err) {
        Message.error(err.message || '删除失败')
      }
    }
  })
}

const showIssueCard = (record) => {
  currentCardType.value = record
  issueForm.memberId = undefined
  issueForm.remark = ''
  issueModalVisible.value = true
}

const handleIssueCard = async () => {
  if (!issueForm.memberId) {
    Message.warning('请选择会员')
    return
  }
  Message.success('发行成功')
  issueModalVisible.value = false
  loadIssueRecords()
}

const showBatchIssue = () => {
  batchForm.cardTypeId = undefined
  batchForm.count = 1
  batchForm.remark = ''
  batchIssueVisible.value = true
}

const handleBatchIssue = async () => {
  if (!batchForm.cardTypeId) {
    Message.warning('请选择卡类型')
    return
  }
  Message.success('批量发行成功')
  batchIssueVisible.value = false
  loadIssueRecords()
}

const onPageChange = (page) => { pagination.current = page; loadCardTypes() }
const onIssuePageChange = (page) => { issuePagination.current = page; loadIssueRecords() }

onMounted(() => {
  loadCardTypes()
  loadIssueRecords()
})
</script>

<style scoped>
.member-card-page { padding: 20px 24px; min-height: calc(100vh - 64px); background: #f5f7fa; }
.breadcrumb { margin-bottom: 16px; }
.stats-row { margin-bottom: 16px; }
.stat-card { border-radius: 8px; text-align: center; }
.action-card { margin-bottom: 16px; }
.table-card { border-radius: 8px; }
</style>
