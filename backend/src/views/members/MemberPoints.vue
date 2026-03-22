<template>
  <div class="member-points-page">
    <!-- 面包屑 -->
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>会员管理</a-breadcrumb-item>
      <a-breadcrumb-item>积分规则</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 统计卡片 -->
    <a-row :gutter="16" class="stats-row">
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="会员总数" :value="stats.total || 0">
            <template #prefix><icon-user /></template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="总积分池" :value="stats.totalPoints || 0" :value-style="{ color: '#ff6b00' }">
            <template #prefix><icon-star /></template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="本月消耗" :value="stats.monthUsed || 0" :value-style="{ color: '#52c41a' }" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="积分规则" :value="rulesSet ? '已配置' : '未配置'" :value-style="{ color: rulesSet ? '#52c41a' : '#ff4d4f' }">
            <template #prefix><icon-settings /></template>
          </a-statistic>
        </a-card>
      </a-col>
    </a-row>

    <!-- 积分规则配置 -->
    <a-card class="rule-card">
      <template #title>
        <a-space>
          <icon-settings style="font-size: 16px;" />
          <span style="font-weight: 600; font-size: 15px;">积分规则配置</span>
        </a-space>
      </template>
      <template #extra>
        <a-button type="primary" @click="showRulesModal = true">{{ rulesSet ? '编辑规则' : '设置规则' }}</a-button>
      </template>

      <a-descriptions :column="3" v-if="rulesSet">
        <a-descriptions-item label="基础积分倍率">
          <span style="font-weight: 600; color: #ff6b00;">{{ pointsRules.baseRate || 1 }} 元 = 1 积分</span>
        </a-descriptions-item>
        <a-descriptions-item label="积分抵现比例">
          <span>{{ pointsRules.pointsToMoney || 100 }} 积分 = 1 元</span>
        </a-descriptions-item>
        <a-descriptions-item label="最高抵扣比例">
          <span>{{ ((pointsRules.maxDeductRate || 0.5) * 100).toFixed(0) }}%</span>
        </a-descriptions-item>
        <a-descriptions-item label="最低抵扣积分">
          <span>{{ pointsRules.minPointsToUse || 100 }} 积分</span>
        </a-descriptions-item>
        <a-descriptions-item label="生日双倍积分">
          <a-tag :color="pointsRules.birthdayDouble ? 'green' : 'gray'">
            {{ pointsRules.birthdayDouble ? '开启' : '关闭' }}
          </a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="签到积分">
          <span>{{ pointsRules.signInEnabled ? pointsRules.signInPoints + ' 积分/次' : '未开启' }}</span>
        </a-descriptions-item>
      </a-descriptions>
      <a-empty v-else description="暂未配置积分规则，请点击右侧按钮设置" />
    </a-card>

    <!-- 操作区域 -->
    <a-card class="action-card">
      <a-space wrap>
        <a-select v-model="filters.level" placeholder="会员等级" allow-clear style="width: 120px" @change="handleSearch">
          <a-option v-for="lv in levelList" :key="lv.id" :value="lv.id">{{ lv.name }}</a-option>
        </a-select>
        <a-input-search v-model="filters.keyword" placeholder="搜索会员名称/手机号" style="width: 200px" search-button @search="handleSearch" />
        <a-button type="primary" @click="showAdjustDrawer = true">积分调整</a-button>
        <a-button @click="loadMembers">刷新</a-button>
      </a-space>
    </a-card>

    <!-- 会员积分列表 -->
    <a-card>
      <a-table
        :columns="columns"
        :data="memberList"
        :loading="loading"
        :pagination="paginationConfig"
        @page-change="onPageChange"
        @page-size-change="onPageSizeChange"
        row-key="id"
        :scroll="{ x: 1100 }"
      >
        <template #level="{ record }">
          <a-tag :color="getLevelColor(record.levelId)">{{ record.levelName || '普通' }}</a-tag>
        </template>
        <template #points="{ record }">
          <span style="color: #ff6b00; font-weight: 600;">{{ record.totalPoints || 0 }}</span>
        </template>
        <template #availablePoints="{ record }">
          <span style="color: #1890ff;">{{ record.availablePoints || 0 }}</span>
        </template>
        <template #growthValue="{ record }">
          <span style="color: #52c41a;">{{ record.growthValue || record.growth_value || 0 }}</span>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="openDetail(record)">详情</a-button>
            <a-button type="text" size="small" @click="adjustPoints(record)">调整</a-button>
            <a-button type="text" size="small" @click="viewHistory(record)">记录</a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 积分调整抽屉 -->
    <a-drawer v-model:visible="showAdjustDrawer" title="积分调整" :width="440" @before-ok="handleAdjustSubmit" :loading="formLoading">
      <a-form :model="adjustForm" layout="vertical">
        <a-form-item label="会员" required>
          <a-select
            v-model="adjustForm.memberId"
            placeholder="搜索选择会员"
            searchable
            :filter-option="false"
            @search="searchMembers"
          >
            <a-option v-for="m in searchResults" :key="m.id" :value="m.id">
              {{ m.name }} ({{ m.mobile }})
            </a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="调整类型" required>
          <a-radio-group v-model="adjustForm.type">
            <a-radio value="add">增加</a-radio>
            <a-radio value="deduct">扣除</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item label="积分数量" required>
          <a-input-number v-model="adjustForm.points" :min="1" :max="100000" placeholder="请输入积分数量" style="width: 100%" />
        </a-form-item>
        <a-form-item label="调整原因" required>
          <a-select v-model="adjustForm.reason" placeholder="选择调整原因">
            <a-option value="complaint_compensation">投诉补偿</a-option>
            <a-option value="activity_reward">活动奖励</a-option>
            <a-option value="birthday_gift">生日礼物</a-option>
            <a-option value="error_correction">纠错调整</a-option>
            <a-option value="violation_punish">违规处罚</a-option>
            <a-option value="other">其他</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="备注">
          <a-textarea v-model="adjustForm.remark" :rows="2" placeholder="补充说明" />
        </a-form-item>
      </a-form>
    </a-drawer>

    <!-- 积分详情抽屉 -->
    <a-drawer v-model:visible="showDetailDrawer" title="积分详情" :width="520">
      <template v-if="currentMember">
        <a-descriptions :column="1" bordered size="small">
          <a-descriptions-item label="会员名称">{{ currentMember.name }}</a-descriptions-item>
          <a-descriptions-item label="手机号">{{ currentMember.mobile }}</a-descriptions-item>
          <a-descriptions-item label="会员等级">
            <a-tag :color="getLevelColor(currentMember.levelId)">{{ currentMember.levelName || '普通' }}</a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="总积分">
            <span style="color: #ff6b00; font-weight: 600; font-size: 18px;">{{ currentMember.totalPoints || 0 }}</span>
          </a-descriptions-item>
          <a-descriptions-item label="可用积分">
            <span style="color: #1890ff; font-weight: 600;">{{ currentMember.availablePoints || 0 }}</span>
          </a-descriptions-item>
          <a-descriptions-item label="冻结积分">{{ currentMember.frozenPoints || 0 }}</a-descriptions-item>
          <a-descriptions-item label="成长值">{{ currentMember.growthValue || 0 }}</a-descriptions-item>
          <a-descriptions-item label="注册时间">{{ formatTime(currentMember.createdAt) }}</a-descriptions-item>
        </a-descriptions>
        <a-divider>积分记录</a-divider>
        <a-spin :loading="historyLoading">
          <a-list :data="pointHistory" size="small" :virtual-list-props="{ height: 300 }">
            <template #item="{ item }">
              <a-list-item>
                <a-list-item-meta
                  :title="getFlowTypeText(item.type) + ' - ' + (item.description || '-')"
                  :description="formatTime(item.createdAt)"
                />
                <template #extra>
                  <span :style="{ color: (item.type === 1 || item.type === 'add') ? '#52c41a' : '#ff4d4f', fontWeight: 600 }">
                    {{ (item.type === 1 || item.type === 'add') ? '+' : '-' }}{{ item.points }}
                  </span>
                </template>
              </a-list-item>
            </template>
          </a-list>
          <a-empty v-if="!pointHistory.length && !historyLoading" description="暂无积分记录" />
        </a-spin>
      </template>
    </a-drawer>

    <!-- 积分规则弹窗 -->
    <a-modal v-model:visible="showRulesModal" :title="rulesSet ? '编辑积分规则' : '设置积分规则'" :width="520" @before-ok="handleRulesSubmit" :loading="formLoading">
      <a-form :model="rulesForm" layout="vertical">
        <a-form-item label="基础积分倍率" required>
          <a-input-number v-model="rulesForm.baseRate" :min="0.1" :step="0.1" style="width: 100%">
            <template #suffix>元 = 1 积分</template>
          </a-input-number>
        </a-form-item>
        <a-form-item label="积分抵现比例" required>
          <a-input-number v-model="rulesForm.pointsToMoney" :min="1" style="width: 100%">
            <template #suffix>积分 = 1 元</template>
          </a-input-number>
        </a-form-item>
        <a-form-item label="单笔最高抵扣比例">
          <a-input-number v-model="rulesForm.maxDeductRate" :min="0" :max="1" :step="0.05" style="width: 100%">
            <template #suffix>%</template>
          </a-input-number>
        </a-form-item>
        <a-form-item label="最低抵扣积分">
          <a-input-number v-model="rulesForm.minPointsToUse" :min="1" style="width: 100%" />
        </a-form-item>
        <a-form-item label="生日双倍积分">
          <a-switch v-model="rulesForm.birthdayDouble" />
        </a-form-item>
        <a-form-item label="开启签到积分">
          <a-switch v-model="rulesForm.signInEnabled" />
        </a-form-item>
        <a-form-item v-if="rulesForm.signInEnabled" label="签到积分">
          <a-input-number v-model="rulesForm.signInPoints" :min="1" style="width: 100%">
            <template #suffix>积分/次</template>
          </a-input-number>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import * as api from '@/api/member'

const memberList = ref([])
const levelList = ref([])
const searchResults = ref([])
const pointHistory = ref([])

const loading = ref(false)
const formLoading = ref(false)
const historyLoading = ref(false)
const rulesSet = ref(false)
const showAdjustDrawer = ref(false)
const showDetailDrawer = ref(false)
const showRulesModal = ref(false)
const currentMember = ref(null)

const filters = reactive({ level: undefined, keyword: '' })

const pagination = reactive({ current: 1, pageSize: 20, total: 0 })
const paginationConfig = computed(() => ({
  current: pagination.current, pageSize: pagination.pageSize, total: pagination.total,
  showTotal: true, showPageSize: true
}))

const stats = reactive({ total: 0, totalPoints: 0, monthUsed: 0 })

const pointsRules = reactive({
  baseRate: 1, pointsToMoney: 100, maxDeductRate: 0.5,
  minPointsToUse: 100, birthdayDouble: false, signInEnabled: false, signInPoints: 5
})

const rulesForm = reactive({ ...pointsRules })

const adjustForm = reactive({
  memberId: undefined,
  type: 'add',
  points: 0,
  reason: '',
  remark: ''
})

const columns = [
  { title: '会员名称', dataIndex: 'name', width: 120 },
  { title: '手机号', dataIndex: 'mobile', width: 130 },
  { title: '等级', slotName: 'level', width: 100 },
  { title: '总积分', slotName: 'points', width: 100 },
  { title: '可用积分', slotName: 'availablePoints', width: 100 },
  { title: '成长值', slotName: 'growthValue', width: 100 },
  { title: '注册时间', dataIndex: 'createdAt', width: 170 },
  { title: '操作', slotName: 'actions', width: 160, fixed: 'right' }
]

const getLevelColor = (id) => ({ 1: '#95de64', 2: '#1890ff', 3: '#faad14', 4: '#ff4d4f' }[id] || 'gray')

const getFlowTypeText = (t) => ({ 1: '获取', 2: '抵扣', 3: '调整' }[t] || { add: '获取', deduct: '扣除' }[t] || t)

const loadMembers = async () => {
  loading.value = true
  try {
    const params = { page: pagination.current, pageSize: pagination.pageSize }
    if (filters.level) params.levelId = filters.level
    if (filters.keyword) params.keyword = filters.keyword

    const res = await api.getMemberList(params)
    const d = res.data || {}
    memberList.value = d.list || []
    pagination.total = d.total || 0

    stats.total = d.total || 0
    stats.totalPoints = (d.list || []).reduce((s, m) => s + (m.totalPoints || 0), 0)
  } catch (err) {
    Message.error('加载会员列表失败: ' + err.message)
  } finally {
    loading.value = false
  }
}

const loadLevels = async () => {
  try {
    const res = await api.getLevelList()
    levelList.value = res.data || []
  } catch (err) {}
}

const loadRules = async () => {
  try {
    const res = await api.getPointsRules()
    if (res.data) {
      Object.assign(pointsRules, res.data)
      rulesSet.value = true
    }
  } catch (err) {
    rulesSet.value = false
  }
}

const handleSearch = () => { pagination.current = 1; loadMembers() }
const onPageChange = (page) => { pagination.current = page; loadMembers() }
const onPageSizeChange = (pageSize) => { pagination.pageSize = pageSize; pagination.current = 1; loadMembers() }

const searchMembers = async (keyword) => {
  if (!keyword || keyword.length < 1) { searchResults.value = []; return }
  try {
    const res = await api.getMemberList({ keyword, page: 1, pageSize: 20 })
    searchResults.value = res.data?.list || []
  } catch (err) { searchResults.value = [] }
}

const adjustPoints = (record) => {
  adjustForm.memberId = record.id
  adjustForm.type = 'add'
  adjustForm.points = 0
  adjustForm.reason = ''
  adjustForm.remark = ''
  searchResults.value = [{ id: record.id, name: record.name, mobile: record.mobile }]
  showAdjustDrawer.value = true
}

const handleAdjustSubmit = async (done) => {
  if (!adjustForm.memberId || !adjustForm.points || !adjustForm.reason) {
    Message.warning('请填写完整信息')
    done(false)
    return
  }
  formLoading.value = true
  try {
    const payload = {
      memberId: adjustForm.memberId,
      points: adjustForm.type === 'deduct' ? -Math.abs(adjustForm.points) : Math.abs(adjustForm.points),
      reason: adjustForm.reason,
      remark: adjustForm.remark
    }
    await api.adjustPoints(payload)
    Message.success('积分调整成功')
    showAdjustDrawer.value = false
    loadMembers()
    done(true)
  } catch (err) {
    Message.error(err.message || '调整失败')
    done(false)
  } finally {
    formLoading.value = false
  }
}

const openDetail = async (record) => {
  currentMember.value = record
  showDetailDrawer.value = true
  historyLoading.value = true
  try {
    const [detailRes, flowRes] = await Promise.all([
      api.getMemberDetail(record.id),
      api.getPointsFlow({ memberId: record.id, page: 1, pageSize: 50 })
    ])
    currentMember.value = { ...record, ...(detailRes.data || {}) }
    pointHistory.value = flowRes.data?.list || []
  } catch (err) {
    pointHistory.value = []
  } finally {
    historyLoading.value = false
  }
}

const viewHistory = (record) => openDetail(record)

const handleRulesSubmit = async (done) => {
  formLoading.value = true
  try {
    await api.updatePointsRules(rulesForm)
    Object.assign(pointsRules, rulesForm)
    rulesSet.value = true
    showRulesModal.value = false
    Message.success('规则保存成功')
    done(true)
  } catch (err) {
    Message.error(err.message || '保存失败')
    done(false)
  } finally {
    formLoading.value = false
  }
}

const formatTime = (t) => {
  if (!t) return '-'
  return new Date(t).toLocaleString('zh-CN')
}

onMounted(() => {
  loadMembers()
  loadLevels()
  loadRules()
})
</script>

<style scoped>
.member-points-page {
  padding: 20px 24px;
  min-height: calc(100vh - 64px);
  background: #f5f7fa;
}
.breadcrumb { margin-bottom: 16px; }
.stats-row { margin-bottom: 16px; }
.stat-card { border-radius: 8px; }
.rule-card { margin-bottom: 16px; }
.action-card { margin-bottom: 16px; }
</style>
