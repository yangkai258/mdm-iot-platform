<template>
  <a-layout class="member-points">
    <a-layout-sider v-model:collapsed="collapsed" :trigger="null" collapsible>
      <div class="logo">
        <span v-if="!collapsed">MDM 控制台</span>
      </div>
      <a-menu v-model:selectedKeys="selectedKeys" theme="dark" mode="inline" @click="handleMenuClick">
        <a-menu-item key="dashboard">
          <span>设备大盘</span>
        </a-menu-item>
        <a-menu-item key="members">
          <span>会员管理</span>
        </a-menu-item>
      </a-menu>
    </a-layout-sider>

    <a-layout>
      <a-layout-header class="header">
        <div class="header-left">
          <a-button type="text" @click="collapsed = !collapsed">
            <span v-if="collapsed">☰</span>
            <span v-else>✕</span>
          </a-button>
        </div>
        <div class="header-title">
          <a-breadcrumb>
            <a-breadcrumb-item>会员管理</a-breadcrumb-item>
            <a-breadcrumb-item>积分管理</a-breadcrumb-item>
          </a-breadcrumb>
        </div>
        <div class="header-right"></div>
      </a-layout-header>

      <a-layout-content class="content">
        <!-- 统计卡片 -->
        <a-row :gutter="16" class="stats-row">
          <a-col :span="6">
            <a-card>
              <a-statistic title="会员总数" :value="stats.total" />
            </a-card>
          </a-col>
          <a-col :span="6">
            <a-card>
              <a-statistic title="今日新增" :value="stats.todayNew" :value-style="{ color: '#52c41a' }" />
            </a-card>
          </a-col>
          <a-col :span="6">
            <a-card>
              <a-statistic title="总积分池" :value="stats.totalPoints" />
            </a-card>
          </a-col>
          <a-col :span="6">
            <a-card>
              <a-statistic title="本月消耗" :value="stats.monthUsed" />
            </a-card>
          </a-col>
        </a-row>

        <!-- 操作栏 -->
        <a-card class="action-card">
          <a-space wrap>
            <a-select v-model="filters.level" placeholder="会员等级" allow-clear style="width: 120px" @change="loadMembers">
              <a-option value="gold">黄金</a-option>
              <a-option value="silver">白银</a-option>
              <a-option value="bronze">青铜</a-option>
            </a-select>
            <a-input-search v-model="filters.keyword" placeholder="搜索会员名称/手机号" style="width: 200px" search-button @search="loadMembers" />
            <a-button type="primary" @click="showAdjustDrawer = true">积分调整</a-button>
            <a-button @click="loadMembers">刷新</a-button>
          </a-space>
        </a-card>

        <!-- 会员积分列表 -->
        <a-card class="points-card">
          <a-table
            :columns="columns"
            :data="memberList"
            :loading="loading"
            :pagination="pagination"
            row-key="id"
            @page-change="handlePageChange"
            @page-size-change="handlePageSizeChange"
          >
            <template #level="{ record }">
              <a-tag :color="getLevelColor(record.level)">{{ getLevelText(record.level) }}</a-tag>
            </template>
            <template #points="{ record }">
              <span style="color: #ff6b00; font-weight: 600;">{{ record.points || 0 }}</span>
            </template>
            <template #growthValue="{ record }">
              <span style="color: #52c41a;">{{ record.growth_value || 0 }}</span>
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
      </a-layout-content>
    </a-layout>

    <!-- 积分调整抽屉 -->
    <a-drawer
      v-model:visible="showAdjustDrawer"
      title="积分调整"
      :width="400"
      @before-ok="handleAdjust"
    >
      <a-form :model="adjustForm" layout="vertical">
        <a-form-item label="会员" required>
          <a-select v-model="adjustForm.member_id" placeholder="选择会员" searchable>
            <a-option v-for="m in memberList" :key="m.id" :value="m.id">{{ m.member_name }} ({{ m.phone }})</a-option>
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
          <a-textarea v-model="adjustForm.reason" :rows="3" placeholder="请输入调整原因" />
        </a-form-item>
        <a-form-item>
          <a-space>
            <a-button type="primary" @click="handleAdjustSubmit">确认调整</a-button>
            <a-button @click="showAdjustDrawer = false">取消</a-button>
          </a-space>
        </a-form-item>
      </a-form>
    </a-drawer>

    <!-- 积分详情抽屉 -->
    <a-drawer
      v-model:visible="showDetailDrawer"
      title="积分详情"
      :width="520"
    >
      <template v-if="currentMember">
        <a-descriptions :column="1" bordered size="small">
          <a-descriptions-item label="会员名称">{{ currentMember.member_name }}</a-descriptions-item>
          <a-descriptions-item label="手机号">{{ currentMember.phone }}</a-descriptions-item>
          <a-descriptions-item label="会员等级">
            <a-tag :color="getLevelColor(currentMember.level)">{{ getLevelText(currentMember.level) }}</a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="当前积分">
            <span style="color: #ff6b00; font-weight: 600; font-size: 18px;">{{ currentMember.points || 0 }}</span>
          </a-descriptions-item>
          <a-descriptions-item label="成长值">{{ currentMember.growth_value || 0 }}</a-descriptions-item>
          <a-descriptions-item label="累计积分">{{ currentMember.total_earned || 0 }}</a-descriptions-item>
          <a-descriptions-item label="已消耗积分">{{ currentMember.total_used || 0 }}</a-descriptions-item>
          <a-descriptions-item label="注册时间">{{ formatTime(currentMember.created_at) }}</a-descriptions-item>
        </a-descriptions>
        <a-divider>积分记录</a-divider>
        <a-list :data="currentMember.point_history || []" size="small">
          <template #item="{ item }">
            <a-list-item>
              <a-list-item-meta :title="item.description" :description="formatTime(item.created_at)" />
              <template #extra>
                <span :style="{ color: item.type === 'add' ? '#52c41a' : '#ff4d4f' }">
                  {{ item.type === 'add' ? '+' : '-' }}{{ item.points }}
                </span>
              </template>
            </a-list-item>
          </template>
        </a-list>
      </template>
    </a-drawer>
  </a-layout>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'
import { Message } from '@arco-design/web-vue'

const router = useRouter()
const collapsed = ref(false)
const selectedKeys = ref(['members'])
const loading = ref(false)
const memberList = ref([])
const showAdjustDrawer = ref(false)
const showDetailDrawer = ref(false)
const currentMember = ref(null)

const filters = reactive({
  level: undefined,
  keyword: ''
})

const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0
})

const stats = reactive({
  total: 0,
  todayNew: 0,
  totalPoints: 0,
  monthUsed: 0
})

const adjustForm = reactive({
  member_id: undefined,
  type: 'add',
  points: 0,
  reason: ''
})

const columns = [
  { title: '会员名称', dataIndex: 'member_name', width: 150 },
  { title: '手机号', dataIndex: 'phone', width: 130 },
  { title: '会员等级', slotName: 'level', width: 100 },
  { title: '当前积分', slotName: 'points', width: 120 },
  { title: '成长值', slotName: 'growthValue', width: 100 },
  { title: '注册时间', dataIndex: 'created_at', width: 160 },
  { title: '操作', slotName: 'actions', width: 160, fixed: 'right' }
]

const handleMenuClick = ({ key }) => {
  const routes = { dashboard: '/dashboard', members: '/members/points' }
  if (routes[key]) router.push(routes[key])
  selectedKeys.value = [key]
}

const loadMembers = async () => {
  loading.value = true
  try {
    const params = { page: pagination.current, page_size: pagination.pageSize }
    if (filters.level) params.level = filters.level
    if (filters.keyword) params.keyword = filters.keyword

    const res = await axios.get('/api/v1/members/points', { params })
    const data = res.data
    if (data.code === 0) {
      memberList.value = data.data?.list || []
      pagination.total = data.data?.pagination?.total || 0
      updateStats()
    }
  } catch (err) {
    Message.error('加载会员列表失败')
  } finally {
    loading.value = false
  }
}

const updateStats = () => {
  stats.total = memberList.value.length
  stats.todayNew = memberList.value.filter(m => isToday(m.created_at)).length
  stats.totalPoints = memberList.value.reduce((sum, m) => sum + (m.points || 0), 0)
  stats.monthUsed = memberList.value.reduce((sum, m) => sum + (m.month_used || 0), 0)
}

const isToday = (date) => {
  if (!date) return false
  const today = new Date()
  const d = new Date(date)
  return d.toDateString() === today.toDateString()
}

const handlePageChange = (page) => {
  pagination.current = page
  loadMembers()
}

const handlePageSizeChange = (pageSize) => {
  pagination.pageSize = pageSize
  pagination.current = 1
  loadMembers()
}

const handleAdjust = (done) => { done(true) }

const handleAdjustSubmit = async () => {
  if (!adjustForm.member_id || !adjustForm.points || !adjustForm.reason) {
    Message.warning('请填写完整信息')
    return
  }
  try {
    const res = await axios.post('/api/v1/members/points/adjust', adjustForm)
    if (res.data.code === 0) {
      Message.success('积分调整成功')
      showAdjustDrawer.value = false
      resetAdjustForm()
      loadMembers()
    } else {
      Message.error(res.data.message || '调整失败')
    }
  } catch (err) {
    Message.error('调整失败')
  }
}

const resetAdjustForm = () => {
  adjustForm.member_id = undefined
  adjustForm.type = 'add'
  adjustForm.points = 0
  adjustForm.reason = ''
}

const adjustPoints = (record) => {
  adjustForm.member_id = record.id
  showAdjustDrawer.value = true
}

const openDetail = (record) => {
  currentMember.value = record
  showDetailDrawer.value = true
}

const viewHistory = (record) => {
  router.push({ path: '/members/points/history', query: { memberId: record.id } })
}

const getLevelColor = (level) => {
  const map = { gold: '#FFD700', silver: '#C0C0C0', bronze: '#CD7F32' }
  return map[level] || 'gray'
}

const getLevelText = (level) => {
  const map = { gold: '黄金', silver: '白银', bronze: '青铜' }
  return map[level] || level
}

const formatTime = (time) => {
  if (!time) return '-'
  return new Date(time).toLocaleString('zh-CN')
}

onMounted(() => {
  loadMembers()
})
</script>

<style scoped>
.member-points { min-height: 100vh; }
.header { background: #fff; padding: 0 16px; display: flex; align-items: center; gap: 16px; box-shadow: 0 1px 4px rgba(0,0,0,0.1); }
.header-left { display: flex; align-items: center; }
.header-title { font-size: 16px; font-weight: 500; }
.content { padding: 16px; background: #f0f2f5; }
.stats-row { margin-bottom: 16px; }
.action-card { margin-bottom: 16px; }
</style>
