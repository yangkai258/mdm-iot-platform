<template>
  <div class="member-list-page">
    <!-- 面包屑 -->
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>会员管理</a-breadcrumb-item>
      <a-breadcrumb-item>会员列表</a-breadcrumb-item>
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
          <a-statistic title="活跃会员" :value="stats.active || 0" :value-style="{ color: '#52c41a' }">
            <template #prefix><icon-check-circle /></template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="今日新增" :value="stats.todayNew || 0">
            <template #prefix><icon-user-add /></template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="总积分" :value="stats.totalPoints || 0" :value-style="{ color: '#faad14' }">
            <template #prefix><icon-star /></template>
          </a-statistic>
        </a-card>
      </a-col>
    </a-row>

    <!-- 搜索筛选区 -->
    <a-card class="action-card">
      <a-space wrap>
        <a-input-search
          v-model="filters.keyword"
          placeholder="搜索姓名/手机号/会员编号"
          style="width: 240px"
          search-button
          @search="handleSearch"
        />
        <a-select v-model="filters.levelId" placeholder="会员等级" allow-clear style="width: 140px" @change="handleSearch">
          <a-option v-for="lv in levelList" :key="lv.id" :value="lv.id">{{ lv.name }}</a-option>
        </a-select>
        <a-select v-model="filters.status" placeholder="会员状态" allow-clear style="width: 120px" @change="handleSearch">
          <a-option :value="1">正常</a-option>
          <a-option :value="2">冻结</a-option>
          <a-option :value="3">禁用</a-option>
        </a-select>
        <a-button @click="handleSearch">筛选</a-button>
        <a-button @click="resetFilters">重置</a-button>
      </a-space>
    </a-card>

    <!-- 操作+表格 -->
    <a-card class="table-card">
      <template #title>
        <a-space>
          <span style="font-weight: 600; font-size: 15px;">会员列表</span>
          <a-badge :count="pagination.total" :max-count="99999" />
        </a-space>
      </template>
      <template #extra>
        <a-space>
          <a-button type="primary" @click="showCreateModal">新增会员</a-button>
          <a-button @click="loadMembers">刷新</a-button>
        </a-space>
      </template>

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
        <template #avatar="{ record }">
          <a-avatar :style="{ backgroundColor: getLevelColor(record.levelId) }" :size="32">
            {{ (record.name || record.mobile || '?').charAt(0) }}
          </a-avatar>
        </template>
        <template #level="{ record }">
          <a-tag :color="getLevelColor(record.levelId)">{{ record.levelName || '普通' }}</a-tag>
        </template>
        <template #status="{ record }">
          <a-tag :color="getStatusColor(record.status)">{{ getStatusText(record.status) }}</a-tag>
        </template>
        <template #totalPoints="{ record }">
          <span style="color: #ff6b00; font-weight: 600;">{{ record.totalPoints || 0 }}</span>
        </template>
        <template #totalConsume="{ record }">
          <span>¥{{ (record.totalConsume || 0).toFixed(2) }}</span>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="showDetail(record)">详情</a-button>
            <a-button type="text" size="small" @click="showEdit(record)">编辑</a-button>
            <a-button type="text" size="small" @click="showPointsAdjust(record)">积分</a-button>
            <a-dropdown trigger="click">
              <a-button type="text" size="small">更多</a-button>
              <template #content>
                <a-doption @click="showAdjustLevel(record)">调整等级</a-doption>
                <a-doption @click="showStatusModal(record)" v-if="record.status !== 3">禁用</a-doption>
                <a-doption status="danger" @click="handleDelete(record)">删除</a-doption>
              </template>
            </a-dropdown>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 新增/编辑弹窗 -->
    <a-modal
      v-model:visible="formVisible"
      :title="isEdit ? '编辑会员' : '新增会员'"
      @before-ok="handleFormSubmit"
      @cancel="formVisible = false"
      :width="520"
      :loading="formLoading"
    >
      <a-form :model="form" layout="vertical" ref="formRef">
        <a-form-item label="手机号" field="mobile" :rules="[{ required: !isEdit, message: '请输入手机号' }]">
          <a-input v-model="form.mobile" placeholder="请输入手机号" :disabled="isEdit" />
        </a-form-item>
        <a-form-item label="姓名" field="name" :rules="[{ required: true, message: '请输入姓名' }]">
          <a-input v-model="form.name" placeholder="请输入姓名" />
        </a-form-item>
        <a-form-item label="性别">
          <a-radio-group v-model="form.gender">
            <a-radio :value="0">未知</a-radio>
            <a-radio :value="1">男</a-radio>
            <a-radio :value="2">女</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item label="生日">
          <a-date-picker v-model="form.birthday" format="YYYY-MM-DD" style="width: 100%" />
        </a-form-item>
        <a-form-item label="邮箱">
          <a-input v-model="form.email" placeholder="请输入邮箱" />
        </a-form-item>
        <a-form-item label="备注">
          <a-textarea v-model="form.remark" :rows="2" placeholder="备注信息" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 积分调整抽屉 -->
    <a-drawer v-model:visible="pointsVisible" title="积分调整" :width="420">
      <a-form :model="pointsForm" layout="vertical">
        <a-form-item label="会员">
          <a-input :value="currentMember?.name + ' (' + currentMember?.mobile + ')'" disabled />
        </a-form-item>
        <a-form-item label="调整类型" required>
          <a-radio-group v-model="pointsForm.type">
            <a-radio value="add">增加</a-radio>
            <a-radio value="deduct">扣除</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item label="积分数量" required>
          <a-input-number v-model="pointsForm.points" :min="1" :max="1000000" style="width: 100%" />
        </a-form-item>
        <a-form-item label="调整原因" required>
          <a-textarea v-model="pointsForm.reason" :rows="3" placeholder="请输入调整原因" />
        </a-form-item>
      </a-form>
      <template #footer>
        <a-button @click="pointsVisible = false">取消</a-button>
        <a-button type="primary" :loading="formLoading" @click="handlePointsSubmit">确认调整</a-button>
      </template>
    </a-drawer>

    <!-- 会员详情 -->
    <a-drawer v-model:visible="detailVisible" title="会员详情" :width="520">
      <template v-if="currentMember">
        <a-descriptions :column="1" bordered size="small">
          <a-descriptions-item label="会员编号">{{ currentMember.memberNo }}</a-descriptions-item>
          <a-descriptions-item label="姓名">{{ currentMember.name }}</a-descriptions-item>
          <a-descriptions-item label="手机号">{{ currentMember.mobile }}</a-descriptions-item>
          <a-descriptions-item label="性别">{{ getGenderText(currentMember.gender) }}</a-descriptions-item>
          <a-descriptions-item label="生日">{{ currentMember.birthday || '-' }}</a-descriptions-item>
          <a-descriptions-item label="会员等级">
            <a-tag :color="getLevelColor(currentMember.levelId)">{{ currentMember.levelName || '普通' }}</a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="状态">
            <a-tag :color="getStatusColor(currentMember.status)">{{ getStatusText(currentMember.status) }}</a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="可用积分">
            <span style="color: #ff6b00; font-weight: 600;">{{ currentMember.totalPoints || 0 }}</span>
          </a-descriptions-item>
          <a-descriptions-item label="累计消费">
            <span>¥{{ (currentMember.totalConsume || 0).toFixed(2) }}</span>
          </a-descriptions-item>
          <a-descriptions-item label="订单数">{{ currentMember.totalOrderCount || 0 }}</a-descriptions-item>
          <a-descriptions-item label="所属门店">{{ currentMember.storeName || '-' }}</a-descriptions-item>
          <a-descriptions-item label="注册时间">{{ currentMember.createdAt || '-' }}</a-descriptions-item>
        </a-descriptions>
        <div style="margin-top: 16px">
          <a-space>
            <a-button type="primary" size="small" @click="goDetail(currentMember)">查看完整详情</a-button>
            <a-button size="small" @click="showEdit(currentMember)">编辑</a-button>
          </a-space>
        </div>
      </template>
    </a-drawer>

    <!-- 调整等级弹窗 -->
    <a-modal v-model:visible="levelModalVisible" title="调整会员等级" @before-ok="handleLevelSubmit" :width="400" :loading="formLoading">
      <a-form :model="levelForm" layout="vertical">
        <a-form-item label="会员">
          <a-input :value="currentMember?.name + ' (' + currentMember?.mobile + ')'" disabled />
        </a-form-item>
        <a-form-item label="目标等级" required>
          <a-select v-model="levelForm.levelId" placeholder="请选择目标等级">
            <a-option v-for="lv in levelList" :key="lv.id" :value="lv.id">{{ lv.name }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="调整原因" required>
          <a-select v-model="levelForm.reason" placeholder="请选择调整原因">
            <a-option value="upgrade">升级</a-option>
            <a-option value="downgrade">降级</a-option>
            <a-option value="vip">VIP专属</a-option>
            <a-option value="other">其他</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="备注说明">
          <a-textarea v-model="levelForm.remark" :rows="2" placeholder="备注说明" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Message } from '@arco-design/web-vue'
import * as api from '@/api/member'

const router = useRouter()
const memberList = ref([])
const levelList = ref([])
const loading = ref(false)
const formLoading = ref(false)
const formVisible = ref(false)
const detailVisible = ref(false)
const levelModalVisible = ref(false)
const pointsVisible = ref(false)
const isEdit = ref(false)
const currentId = ref(null)
const currentMember = ref(null)

const filters = reactive({
  keyword: '',
  levelId: undefined,
  status: undefined
})

const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0
})

const paginationConfig = computed(() => ({
  current: pagination.current,
  pageSize: pagination.pageSize,
  total: pagination.total,
  showTotal: true,
  showPageSize: true,
  pageSizeOptions: [10, 20, 50, 100]
}))

const stats = reactive({ total: 0, active: 0, todayNew: 0, totalPoints: 0 })

const form = reactive({
  mobile: '',
  name: '',
  gender: 0,
  birthday: '',
  email: '',
  remark: ''
})

const pointsForm = reactive({
  type: 'add',
  points: 0,
  reason: ''
})

const levelForm = reactive({
  levelId: undefined,
  reason: '',
  remark: ''
})

const columns = [
  { title: '头像', slotName: 'avatar', width: 70, fixed: 'left' },
  { title: '会员编号', dataIndex: 'memberNo', width: 170 },
  { title: '姓名', dataIndex: 'name', width: 100 },
  { title: '手机号', dataIndex: 'mobile', width: 130 },
  { title: '等级', slotName: 'level', width: 100 },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '积分', slotName: 'totalPoints', width: 100 },
  { title: '累计消费', slotName: 'totalConsume', width: 110 },
  { title: '订单数', dataIndex: 'totalOrderCount', width: 80 },
  { title: '注册时间', dataIndex: 'createdAt', width: 170 },
  { title: '操作', slotName: 'actions', width: 200, fixed: 'right' }
]

const getLevelColor = (id) => {
  const colors = { 1: '#95de64', 2: '#1890ff', 3: '#faad14', 4: '#ff4d4f' }
  return colors[id] || 'gray'
}

const getStatusColor = (s) => ({ 1: 'green', 2: 'orange', 3: 'red' }[s] || 'gray')
const getStatusText = (s) => ({ 1: '正常', 2: '冻结', 3: '禁用' }[s] || '未知')
const getGenderText = (g) => ({ 0: '未知', 1: '男', 2: '女' }[g] || '未知')

const loadMembers = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.current,
      pageSize: pagination.pageSize
    }
    if (filters.keyword) params.keyword = filters.keyword
    if (filters.levelId) params.levelId = filters.levelId
    if (filters.status) params.status = filters.status

    const res = await api.getMemberList(params)
    const d = res.data || {}
    memberList.value = d.list || []
    pagination.total = d.total || 0
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
  } catch (err) { /* ignore */ }
}

const loadStats = async () => {
  try {
    const res = await api.getMemberList({ page: 1, pageSize: 1 })
    stats.total = res.data?.total || 0
    stats.active = stats.total
    stats.todayNew = 0
    stats.totalPoints = 0
  } catch (err) { /* ignore */ }
}

const handleSearch = () => {
  pagination.current = 1
  loadMembers()
}

const resetFilters = () => {
  filters.keyword = ''
  filters.levelId = undefined
  filters.status = undefined
  pagination.current = 1
  loadMembers()
}

const onPageChange = (page) => {
  pagination.current = page
  loadMembers()
}

const onPageSizeChange = (pageSize) => {
  pagination.pageSize = pageSize
  pagination.current = 1
  loadMembers()
}

const showCreateModal = () => {
  isEdit.value = false
  Object.assign(form, { mobile: '', name: '', gender: 0, birthday: '', email: '', remark: '' })
  formVisible.value = true
}

const showEdit = (record) => {
  isEdit.value = true
  currentId.value = record.id
  Object.assign(form, {
    mobile: record.mobile,
    name: record.name,
    gender: record.gender || 0,
    birthday: record.birthday || '',
    email: record.email || '',
    remark: record.remark || ''
  })
  formVisible.value = true
}

const showDetail = async (record) => {
  try {
    const res = await api.getMemberDetail(record.id)
    currentMember.value = res.data || record
    detailVisible.value = true
  } catch (err) {
    currentMember.value = record
    detailVisible.value = true
  }
}

const goDetail = (record) => {
  router.push(`/members/detail/${record.id}`)
}

const handleFormSubmit = async (done) => {
  formLoading.value = true
  try {
    const payload = { ...form }
    if (isEdit.value) {
      await api.updateMember(currentId.value, payload)
      Message.success('更新成功')
    } else {
      await api.createMember(payload)
      Message.success('创建成功')
    }
    formVisible.value = false
    loadMembers()
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
    await api.deleteMember(record.id)
    Message.success('删除成功')
    loadMembers()
  } catch (err) {
    Message.error(err.message || '删除失败')
  }
}

const showAdjustLevel = (record) => {
  currentMember.value = record
  Object.assign(levelForm, { levelId: undefined, reason: '', remark: '' })
  levelModalVisible.value = true
}

const handleLevelSubmit = async (done) => {
  if (!levelForm.levelId || !levelForm.reason) {
    Message.warning('请填写完整信息')
    done(false)
    return
  }
  try {
    await api.updateMemberLevel(currentMember.value.id, levelForm)
    Message.success('等级调整成功')
    levelModalVisible.value = false
    loadMembers()
    done(true)
  } catch (err) {
    Message.error(err.message || '调整失败')
    done(false)
  }
}

const showStatusModal = async (record) => {
  try {
    await api.updateMemberStatus(record.id, { status: 3, reason: '后台禁用' })
    Message.success('会员已禁用')
    loadMembers()
  } catch (err) {
    Message.error(err.message || '操作失败')
  }
}

const showPointsAdjust = (record) => {
  currentMember.value = record
  Object.assign(pointsForm, { type: 'add', points: 0, reason: '' })
  pointsVisible.value = true
}

const handlePointsSubmit = async () => {
  if (!pointsForm.points || !pointsForm.reason) {
    Message.warning('请填写完整信息')
    return
  }
  try {
    await api.adjustPoints({ memberId: currentMember.value.id, ...pointsForm })
    Message.success('积分调整成功')
    pointsVisible.value = false
    loadMembers()
  } catch (err) {
    Message.error(err.message || '调整失败')
  }
}

onMounted(() => {
  loadMembers()
  loadLevels()
  loadStats()
})
</script>

<style scoped>
.member-list-page {
  padding: 20px 24px;
  min-height: calc(100vh - 64px);
  background: #f5f7fa;
}
.breadcrumb { margin-bottom: 16px; }
.stats-row { margin-bottom: 16px; }
.stat-card { border-radius: 8px; }
.action-card { margin-bottom: 16px; }
.table-card { border-radius: 8px; }
</style>
