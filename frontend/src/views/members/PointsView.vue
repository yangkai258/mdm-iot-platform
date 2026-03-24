<template>
  <div class="points-page">
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>会员管理</a-breadcrumb-item>
      <a-breadcrumb-item>积分规则</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 统计卡片 -->
    <a-row :gutter="16" class="stats-row">
      <a-col :span="6"><a-card class="stat-card"><a-statistic title="会员总数" :value="stats.total || 0" /></a-card></a-col>
      <a-col :span="6"><a-card class="stat-card"><a-statistic title="今日新增" :value="stats.todayNew || 0" :value-style="{ color: '#52c41a' }" /></a-card></a-col>
      <a-col :span="6"><a-card class="stat-card"><a-statistic title="总积分池" :value="stats.totalPoints || 0" /></a-card></a-col>
      <a-col :span="6"><a-card class="stat-card"><a-statistic title="本月消耗" :value="stats.monthUsed || 0" :value-style="{ color: '#ff4d4f' }" /></a-card></a-col>
    </a-row>

    <!-- 搜索筛选 -->
    <a-card class="action-card">
      <a-space wrap>
        <a-input-search v-model="filters.keyword" placeholder="搜索会员名称/手机号" style="width: 220px" search-button @search="loadMembers" />
        <a-select v-model="filters.levelId" placeholder="会员等级" allow-clear style="width: 140px" @change="loadMembers">
          <a-option v-for="lv in levelList" :key="lv.id" :value="lv.id">{{ lv.name }}</a-option>
        </a-select>
        <a-button type="primary" @click="showAdjustDrawer = true">积分调整</a-button>
        <a-button @click="loadMembers">刷新</a-button>
      </a-space>
    </a-card>

    <!-- Tab: 会员积分 / 积分规则 -->
    <a-card class="table-card">
      <a-tabs v-model:active-key="activeTab">
        <a-tab-pane key="members" title="会员积分">
          <a-table :columns="memberColumns" :data="memberList" :loading="loading" :pagination="paginationConfig"
            @page-change="onPageChange" @page-size-change="onPageSizeChange" row-key="id" :scroll="{ x: 1000 }">
            <template #level="{ record }"><a-tag :color="getLevelColor(record.levelId)">{{ record.levelName || '普通' }}</a-tag></template>
            <template #points="{ record }"><span style="color: #ff6b00; font-weight: 600;">{{ record.totalPoints || 0 }}</span></template>
            <template #growthValue="{ record }"><span style="color: #52c41a;">{{ record.growthValue || 0 }}</span></template>
            <template #actions="{ record }">
              <a-space>
                <a-button type="text" size="small" @click="adjustPoints(record)">调整</a-button>
                <a-button type="text" size="small" @click="viewHistory(record)">记录</a-button>
              </a-space>
            </template>
          </a-table>
        </a-tab-pane>
        <a-tab-pane key="rules" title="积分规则">
          <a-form :model="rulesForm" layout="vertical" style="max-width: 600px;">
            <a-form-item label="消费1元获得积分">
              <a-input-number v-model="rulesForm.consumeRate" :min="0" :step="1" style="width: 200px" />
            </a-form-item>
            <a-form-item label="每日登录获得积分">
              <a-input-number v-model="rulesForm.loginPoints" :min="0" :step="1" style="width: 200px" />
            </a-form-item>
            <a-form-item label="分享获得积分">
              <a-input-number v-model="rulesForm.sharePoints" :min="0" :step="1" style="width: 200px" />
            </a-form-item>
            <a-form-item label="积分抵现比例">
              <a-input-number v-model="rulesForm.exchangeRate" :min="0" :step="0.01" style="width: 200px">
                <template #suffix>积分=1元</template>
              </a-input-number>
            </a-form-item>
            <a-form-item label="积分过期规则">
              <a-select v-model="rulesForm.expireType" placeholder="选择过期规则">
                <a-option value="never">永不过期</a-option>
                <a-option value="year">每年底过期</a-option>
                <a-option value="month12">12个月后过期</a-option>
              </a-select>
            </a-form-item>
            <a-form-item>
              <a-button type="primary" @click="handleSaveRules">保存规则</a-button>
            </a-form-item>
          </a-form>
        </a-tab-pane>
        <a-tab-pane key="flow" title="积分流水">
          <a-table :columns="flowColumns" :data="flowList" :loading="flowLoading" :pagination="flowPaginationConfig"
            @page-change="onFlowPageChange" row-key="id" :scroll="{ x: 900 }">
            <template #type="{ record }"><a-tag :color="record.type === 'add' ? 'green' : 'red'">{{ record.type === 'add' ? '获得' : '消耗' }}</a-tag></template>
            <template #points="{ record }"><span :style="{ color: record.type === 'add' ? '#52c41a' : '#ff4d4f', fontWeight: 600 }">{{ record.type === 'add' ? '+' : '-' }}{{ record.points }}</span></template>
          </a-table>
        </a-tab-pane>
      </a-tabs>
    </a-card>

    <!-- 积分调整抽屉 -->
    <a-drawer v-model:visible="showAdjustDrawer" title="积分调整" :width="420">
      <a-form :model="adjustForm" layout="vertical">
        <a-form-item label="会员" required>
          <a-select v-model="adjustForm.memberId" placeholder="选择会员" searchable>
            <a-option v-for="m in memberOptions" :key="m.id" :value="m.id">{{ m.name }} ({{ m.mobile }})</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="调整类型" required>
          <a-radio-group v-model="adjustForm.type">
            <a-radio value="add">增加</a-radio>
            <a-radio value="deduct">扣除</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item label="积分数量" required>
          <a-input-number v-model="adjustForm.points" :min="1" :max="1000000" style="width: 100%" />
        </a-form-item>
        <a-form-item label="调整原因" required>
          <a-textarea v-model="adjustForm.reason" :rows="3" placeholder="请输入调整原因" />
        </a-form-item>
      </a-form>
      <template #footer>
        <a-button @click="showAdjustDrawer = false">取消</a-button>
        <a-button type="primary" :loading="formLoading" @click="handleAdjust">确认调整</a-button>
      </template>
    </a-drawer>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import * as api from '@/api/member'

const activeTab = ref('members')
const memberList = ref([])
const flowList = ref([])
const memberOptions = ref([])
const levelList = ref([])
const loading = ref(false)
const flowLoading = ref(false)
const formLoading = ref(false)
const showAdjustDrawer = ref(false)
const currentMember = ref(null)

const filters = reactive({ keyword: '', levelId: undefined })
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })
const flowPagination = reactive({ current: 1, pageSize: 20, total: 0 })
const stats = reactive({ total: 0, todayNew: 0, totalPoints: 0, monthUsed: 0 })
const adjustForm = reactive({ memberId: undefined, type: 'add', points: 0, reason: '' })
const rulesForm = reactive({ consumeRate: 1, loginPoints: 5, sharePoints: 10, exchangeRate: 100, expireType: 'never' })

const paginationConfig = computed(() => ({ current: pagination.current, pageSize: pagination.pageSize, total: pagination.total, showTotal: true, showPageSize: true, pageSizeOptions: [10, 20, 50, 100] }))
const flowPaginationConfig = computed(() => ({ current: flowPagination.current, pageSize: flowPagination.pageSize, total: flowPagination.total, showTotal: true, showPageSize: true }))

const memberColumns = [
  { title: '会员名称', dataIndex: 'name', width: 150 },
  { title: '手机号', dataIndex: 'mobile', width: 130 },
  { title: '会员等级', slotName: 'level', width: 100 },
  { title: '当前积分', slotName: 'points', width: 120 },
  { title: '成长值', slotName: 'growthValue', width: 100 },
  { title: '注册时间', dataIndex: 'createdAt', width: 160 },
  { title: '操作', slotName: 'actions', width: 160, fixed: 'right' }
]

const flowColumns = [
  { title: '类型', slotName: 'type', width: 90 },
  { title: '积分', slotName: 'points', width: 100 },
  { title: '会员', dataIndex: 'memberName', width: 120 },
  { title: '来源/原因', dataIndex: 'reason', ellipsis: true },
  { title: '时间', dataIndex: 'createdAt', width: 170 }
]

const getLevelColor = (id) => ({ 1: '#95de64', 2: '#1890ff', 3: '#faad14', 4: '#ff4d4f' }[id] || 'gray')

const loadMembers = async () => {
  loading.value = true
  try {
    const params = { page: pagination.current, pageSize: pagination.pageSize }
    if (filters.keyword) params.keyword = filters.keyword
    if (filters.levelId) params.levelId = filters.levelId
    const res = await api.getMemberList(params)
    const d = res.data || {}
    memberList.value = d.list || []
    pagination.total = d.total || 0
    stats.total = d.total || 0
  } catch (err) { Message.error('加载会员列表失败') }
  finally { loading.value = false }
}

const loadFlow = async () => {
  flowLoading.value = true
  try {
    const res = await api.getPointsFlow({ page: flowPagination.current, pageSize: flowPagination.pageSize })
    const d = res.data || {}
    flowList.value = d.list || []
    flowPagination.total = d.total || 0
  } catch (err) { /* ignore */ }
  finally { flowLoading.value = false }
}

const loadLevels = async () => {
  try { const res = await api.getLevelList(); levelList.value = res.data || [] } catch (err) { /* ignore */ }
}

const loadMemberOptions = async () => {
  try { const res = await api.getMemberList({ page: 1, pageSize: 100 }); memberOptions.value = res.data?.list || [] } catch (err) { /* ignore */ }
}

const loadRules = async () => {
  try {
    const res = await api.getPointsRules()
    if (res.data) Object.assign(rulesForm, res.data)
  } catch (err) { /* ignore */ }
}

const adjustPoints = (record) => {
  currentMember.value = record
  adjustForm.memberId = record.id
  showAdjustDrawer.value = true
}

const handleAdjust = async () => {
  if (!adjustForm.memberId || !adjustForm.points || !adjustForm.reason) { Message.warning('请填写完整信息'); return }
  formLoading.value = true
  try {
    await api.adjustPoints({ ...adjustForm })
    Message.success('积分调整成功')
    showAdjustDrawer.value = false
    loadMembers()
  } catch (err) { Message.error(err.message || '调整失败') }
  finally { formLoading.value = false }
}

const handleSaveRules = async () => {
  try { await api.updatePointsRules(rulesForm); Message.success('规则保存成功') }
  catch (err) { Message.error(err.message || '保存失败') }
}

const viewHistory = (record) => {
  activeTab.value = 'flow'
  flowPagination.current = 1
  loadFlow()
}

const onPageChange = (page) => { pagination.current = page; loadMembers() }
const onPageSizeChange = (pageSize) => { pagination.pageSize = pageSize; pagination.current = 1; loadMembers() }
const onFlowPageChange = (page) => { flowPagination.current = page; loadFlow() }

onMounted(() => { loadMembers(); loadLevels(); loadRules() })
</script>

<style scoped>
.points-page { padding: 20px 24px; min-height: calc(100vh - 64px); background: #f5f7fa; }
.breadcrumb { margin-bottom: 16px; }
.stats-row { margin-bottom: 16px; }
.stat-card { border-radius: 8px; text-align: center; }
.action-card { margin-bottom: 16px; }
.table-card { border-radius: 8px; }
</style>
