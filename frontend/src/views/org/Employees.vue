<template>
  <div class="page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item><a href="#/dashboard">首页</a></a-breadcrumb-item>
      <a-breadcrumb-item>组织管理</a-breadcrumb-item>
      <a-breadcrumb-item>员工管理</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 搜索框（左上角） -->
    <div class="search-bar">
      <a-input-search
        v-model="searchKey"
        placeholder="搜索..."
        style="width: 280px"
        @search="loadEmployees"
      />
    </div>

    <!-- 操作按钮组（靠左） -->
    <div class="action-bar">
      <a-space :size="12">
        <a-button type="primary" @click="openOnboardModal">「新建」</a-button>
        <a-button @click="showFilter = !showFilter">「筛选」</a-button>
        <a-button @click="loadEmployees">「刷新」</a-button>
      </a-space>
    </div>

    <!-- 筛选面板 -->
    <a-card v-if="showFilter" :bordered="false" style="margin-bottom: 12px">
      <a-space wrap>
        <a-select v-model="filterCompany" placeholder="所属公司" style="width: 160px" allow-clear>
          <a-option v-for="c in companies" :key="c.id" :value="c.id">{{ c.company_name }}</a-option>
        </a-select>
        <a-select v-model="filterStatus" placeholder="员工状态" style="width: 120px" allow-clear>
          <a-option :value="1">在职</a-option>
          <a-option :value="2">离职</a-option>
        </a-select>
        <a-button type="primary" @click="loadEmployees">查询</a-button>
        <a-button @click="resetFilter">重置</a-button>
      </a-space>
    </a-card>

    <!-- 员工列表 -->
    <a-card :bordered="false">
      <a-table
        :columns="columns"
        :data="filteredData"
        :loading="loading"
        :pagination="{ pageSize: 10 }"
        row-key="id"
        scroll="x"
      >
        <template #empCode="{ record }">
          <a-link @click="openDetailDrawer(record)">{{ record.emp_code }}</a-link>
        </template>
        <template #empName="{ record }">
          <a-space>
            <a-avatar v-if="record.photo" :src="record.photo" :size="24" />
            <span>{{ record.emp_name }}</span>
          </a-space>
        </template>
        <template #gender="{ record }">
          {{ record.gender === 'male' ? '男' : record.gender === 'female' ? '女' : '-' }}
        </template>
        <template #empStatus="{ record }">
          <a-tag :color="getStatusColor(record.emp_status)">
            {{ getStatusText(record.emp_status) }}
          </a-tag>
        </template>
        <template #entryDate="{ record }">
          {{ formatDate(record.entry_date) }}
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="openDetailDrawer(record)">「详情」</a-button>
            <a-button
              v-if="record.emp_status === 1"
              type="text"
              size="small"
              status="warning"
              @click="handleLeave(record)"
            >「办理离职」</a-button>
            <a-button type="text" size="small" @click="openUserLinkModal(record)">「关联用户」</a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 员工详情抽屉 -->
    <a-drawer
      v-model:visible="detailVisible"
      :title="`员工详情 — ${currentRecord?.emp_name}`"
      width="500"
      @cancel="detailVisible = false"
    >
      <a-descriptions :column="1" bordered size="small">
        <a-descriptions-item label="工号">{{ currentRecord?.emp_code }}</a-descriptions-item>
        <a-descriptions-item label="姓名">{{ currentRecord?.emp_name }}</a-descriptions-item>
        <a-descriptions-item label="性别">{{ currentRecord?.gender === 'male' ? '男' : currentRecord?.gender === 'female' ? '女' : '-' }}</a-descriptions-item>
        <a-descriptions-item label="手机号">{{ currentRecord?.phone || '-' }}</a-descriptions-item>
        <a-descriptions-item label="邮箱">{{ currentRecord?.email || '-' }}</a-descri

ptions-item>
        <a-descriptions-item label="身份证号">{{ currentRecord?.id_card || '-' }}</a-descriptions-item>
        <a-descriptions-item label="部门">{{ currentRecord?.department?.dept_name || '-' }}</a-descriptions-item>
        <a-descriptions-item label="岗位">{{ currentRecord?.position?.pos_name || '-' }}</a-descriptions-item>
        <a-descriptions-item label="入职日期">{{ formatDate(currentRecord?.entry_date) }}</a-descriptions-item>
        <a-descriptions-item label="员工状态">
          <a-tag :color="getStatusColor(currentRecord?.emp_status)">
            {{ getStatusText(currentRecord?.emp_status) }}
          </a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="地址">{{ currentRecord?.address || '-' }}</a-descriptions-item>
        <a-descriptions-item label="备注">{{ currentRecord?.remark || '-' }}</a-descriptions-item>
      </a-descriptions>
    </a-drawer>

    <!-- 办理入职 - 全屏模态 (风格D) -->
    <a-modal
      v-model:visible="onboardVisible"
      title="办理入职"
      :fullscreen="true"
      :mask-closable="false"
      @before-ok="submitOnboard"
      @cancel="onboardVisible = false"
    >
      <a-form ref="onboardFormRef" :model="onboardForm" :rules="onboardFormRules" layout="vertical" style="max-width: 800px">
        <a-divider>基本信息</a-divider>
        <a-row :gutter="24">
          <a-col :span="12">
            <a-form-item label="工号" field="emp_code">
              <a-input v-model="onboardForm.emp_code" placeholder="请输入工号" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="姓名" field="emp_name">
              <a-input v-model="onboardForm.emp_name" placeholder="请输入姓名" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="性别" field="gender">
              <a-radio-group v-model="onboardForm.gender">
                <a-radio value="male">男</a-radio>
                <a-radio value="female">女</a-radio>
              </a-radio-group>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="出生日期" field="birth_date">
              <a-date-picker v-model="onboardForm.birth_date" style="width: 100%" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="手机号" field="phone">
              <a-input v-model="onboardForm.phone" placeholder="请输入手机号" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="邮箱" field="email">
              <a-input v-model="onboardForm.email" placeholder="请输入邮箱" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="身份证号" field="id_card">
              <a-input v-model="onboardForm.id_card" placeholder="请输入身份证号" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="入职日期" field="entry_date">
              <a-date-picker v-model="onboardForm.entry_date" style="width: 100%" />
            </a-form-item>
          </a-col>
        </a-row>

        <a-divider>组织信息</a-divider>
        <a-row :gutter="24">
          <a-col :span="12">
            <a-form-item label="所属公司" field="company_id">
              <a-select v-model="onboardForm.company_id" placeholder="请选择所属公司">
                <a-option v-for="c in companies" :key="c.id" :value="c.id">{{ c.company_name }}</a-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="部门" field="dept_id">
              <a-tree-select
                v-model="onboardForm.dept_id"
                :data="deptTreeData"
                placeholder="请选择部门"
                allow-clear
                :field-names="{ key: 'id', title: 'dept_name', children: 'children' }"
              />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="岗位" field="position_id">
              <a-select v-model="onboardForm.position_id" placeholder="请选择岗位" allow-clear>
                <a-option v-for="p in posts" :key="p.id" :value="p.id">{{ p.pos_name }}</a-option>
              </a-select>
            </a-form-item>
          </a-col>
        </a-row>

        <a-divider>联系信息</a-divider>
        <a-row :gutter="24">
          <a-col :span="8">
            <a-form-item label="省" field="province">
              <a-input v-model="onboardForm.province" placeholder="省" />
            </a-form-item>
          </a-col>
          <a-col :span="8">
            <a-form-item label="市" field="city">
              <a-input v-model="onboardForm.city" placeholder="市" />
            </a-form-item>
          </a-col>
          <a-col :span="8">
            <a-form-item label="区" field="district">
              <a-input v-model="onboardForm.district" placeholder="区" />
            </a-form-item>
          </a-col>
          <a-col :span="24">
            <a-form-item label="详细地址" field="address">
              <a-input v-model="onboardForm.address" placeholder="请输入详细地址" />
            </a-form-item>
          </a-col>
          <a-col :span="24">
            <a-form-item label="备注" field="remark">
              <a-textarea v-model="onboardForm.remark" placeholder="请输入备注" :rows="3" />
            </a-form-item>
          </a-col>
        </a-row>
      </a-form>
    </a-modal>

    <!-- 办理离职确认 -->
    <a-modal
      v-model="leaveModalVisible"
      title="办理离职"
      @before-ok="submitLeave"
      @cancel="leaveModalVisible = false"
    >
      <a-result status="warning" :title="`确定要为「${currentLeaveRecord?.emp_name}」办理离职吗？`">
        <template #subtitle>
          离职后该员工的账号将被禁用。
        </template>
      </a-result>
      <a-form :model="leaveForm" layout="vertical" style="max-width: 400px; margin-top: 16px">
        <a-form-item label="离职原因（可选）">
          <a-textarea v-model="leaveForm.reason" placeholder="请输入离职原因" :rows="3" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 关联用户弹窗 -->
    <a-modal
      v-model="linkUserVisible"
      :title="`关联用户 — ${currentRecord?.emp_name}`"
      @before-ok="submitLinkUser"
      @cancel="linkUserVisible = false"
    >
      <a-form :model="linkUserForm" layout="vertical" style="max-width: 400px">
        <a-form-item label="选择系统用户">
          <a-select v-model="linkUserForm.user_id" placeholder="请选择要关联的系统用户" allow-search>
            <a-option v-for="u in systemUsers" :key="u.id" :value="u.id">{{ u.username }} ({{ u.email }})</a-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const loading = ref(false)
const searchKey = ref('')
const showFilter = ref(false)
const filterCompany = ref<number | ''>('')
const filterStatus = ref<number | ''>('')
const detailVisible = ref(false)
const onboardVisible = ref(false)
const leaveModalVisible = ref(false)
const linkUserVisible = ref(false)
const currentRecord = ref<any>(null)
const currentLeaveRecord = ref<any>(null)
const onboardFormRef = ref()

const companies = ref<any[]>([
  { id: 1, company_name: '深圳市智联科技有限公司' },
  { id: 2, company_name: '广州云智数据服务公司' },
])

const deptTreeData = ref<any[]>([
  { id: 1, dept_name: '技术部', children: [{ id: 3, dept_name: '前端组' }, { id: 4, dept_name: '后端组' }] },
  { id: 2, dept_name: '市场部', children: [] },
])

const posts = ref<any[]>([
  { id: 1, pos_name: '前端工程师' },
  { id: 2, pos_name: '后端工程师' },
  { id: 3, pos_name: '产品经理' },
])

const systemUsers = ref<any[]>([
  { id: 1, username: 'admin', email: 'admin@example.com' },
  { id: 2, username: 'zhangsan', email: 'zhangsan@example.com' },
  { id: 3, username: 'lisi', email: 'lisi@example.com' },
])

const onboardForm = ref<any>({
  emp_code: '',
  emp_name: '',
  gender: '',
  birth_date: '',
  phone: '',
  email: '',
  id_card: '',
  entry_date: '',
  company_id: null,
  dept_id: null,
  position_id: null,
  province: '',
  city: '',
  district: '',
  address: '',
  remark: '',
})

const onboardFormRules = {
  emp_code: [{ required: true, message: '请输入工号' }],
  emp_name: [{ required: true, message: '请输入姓名' }],
  company_id: [{ required: true, message: '请选择所属公司' }],
}

const leaveForm = ref({ reason: '' })

const linkUserForm = ref({ user_id: null as number | null })

const columns = [
  { title: '工号', slotName: 'empCode', width: 120 },
  { title: '姓名', slotName: 'empName', width: 120 },
  { title: '性别', slotName: 'gender', width: 60 },
  { title: '手机号', dataIndex: 'phone', width: 130 },
  { title: '部门', dataIndex: 'department.dept_name', width: 100 },
  { title: '岗位', dataIndex: 'position.pos_name', width: 100 },
  { title: '入职日期', slotName: 'entryDate', width: 110 },
  { title: '状态', slotName: 'empStatus', width: 80 },
  { title: '操作', slotName: 'actions', width: 200, fixed: 'right' },
]

const mockData = ref<any[]>([
  {
    id: 1,
    emp_code: 'E001',
    emp_name: '张三',
    gender: 'male',
    phone: '13800138001',
    email: 'zhangsan@zhilian.com',
    id_card: '440100199001011234',
    province: '广东省',
    city: '深圳市',
    district: '南山区',
    address: '科技园路1号',
    dept_id: 3,
    position_id: 1,
    company_id: 1,
    entry_date: '2026-01-15',
    emp_status: 1,
    status: 1,
    remark: '',
    department: { dept_name: '前端组' },
    position: { pos_name: '前端工程师' },
  },
  {
    id: 2,
    emp_code: 'E002',
    emp_name: '李四',
    gender: 'female',
    phone: '13800138002',
    email: 'lisi@zhilian.com',
    id_card: '440100199002021234',
    province: '广东省',
    city: '深圳市',
    district: '南山区',
    address: '科技园路2号',
    dept_id: 4,
    position_id: 2,
    company_id: 1,
    entry_date: '2026-01-20',
    emp_status: 1,
    status: 1,
    remark: '',
    department: { dept_name: '后端组' },
    position: { pos_name: '后端工程师' },
  },
  {
    id: 3,
    emp_code: 'E003',
    emp_name: '王五',
    gender: 'male',
    phone: '13900139001',
    email: 'wangwu@yunzhi.com',
    id_card: '440100199003031234',
    province: '广东省',
    city: '广州市',
    district: '天河区',
    address: '天河路100号',
    dept_id: null,
    position_id: 3,
    company_id: 2,
    entry_date: '2026-02-01',
    emp_status: 2,
    status: 0,
    remark: '已离职',
    department: null,
    position: { pos_name: '产品经理' },
  },
])

const filteredData = computed(() => {
  let data = mockData.value
  if (searchKey.value) {
    data = data.filter(item =>
      item.emp_name.includes(searchKey.value) ||
      item.emp_code.includes(searchKey.value) ||
      (item.phone && item.phone.includes(searchKey.value))
    )
  }
  if (filterCompany.value !== '') {
    data = data.filter(item => item.company_id === filterCompany.value)
  }
  if (filterStatus.value !== '') {
    data = data.filter(item => item.emp_status === filterStatus.value)
  }
  return data
})

const getStatusColor = (status: number) => {
  return status === 1 ? 'green' : 'red'
}

const getStatusText = (status: number) => {
  return status === 1 ? '在职' : '离职'
}

const formatDate = (dateStr: string | null) => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleDateString('zh-CN')
}

const loadEmployees = async () => {
  loading.value = true
  try {
    // API: GET /api/v1/tenants/:tenant_id/employees
    // const res = await axios.get('/api/v1/tenants/1/employees', { params: { keyword: searchKey.value } })
  } catch {
    // use mock
  } finally {
    loading.value = false
  }
}

const resetFilter = () => {
  searchKey.value = ''
  filterCompany.value = ''
  filterStatus.value = ''
  loadEmployees()
}

const openDetailDrawer = (record: any) => {
  currentRecord.value = record
  detailVisible.value = true
}

const openOnboardModal = () => {
  onboardForm.value = {
    emp_code: '',
    emp_name: '',
    gender: '',
    birth_date: '',
    phone: '',
    email: '',
    id_card: '',
    entry_date: '',
    company_id: null,
    dept_id: null,
    position_id: null,
    province: '',
    city: '',
    district: '',
    address: '',
    remark: '',
  }
  onboardFormRef.value?.clearValidate()
  onboardVisible.value = true
}

const submitOnboard = async (done: (val: boolean) => void) => {
  try {
    // API: POST /api/v1/tenants/:tenant_id/employees/onboard
    // const res = await axios.post('/api/v1/tenants/1/employees/onboard', onboardForm.value)
    const newRecord = {
      ...onboardForm.value,
      id: Date.now(),
      emp_status: 1,
      status: 1,
      department: deptTreeData.value.find(d => d.id === onboardForm.value.dept_id) || null,
      position: posts.value.find(p => p.id === onboardForm.value.position_id) || null,
    }
    mockData.value.unshift(newRecord)
    Message.success('入职办理成功')
    onboardVisible.value = false
    done(true)
  } catch {
    Message.error('入职办理失败')
    done(false)
  }
}

const handleLeave = (record: any) => {
  currentLeaveRecord.value = record
  leaveForm.value.reason = ''
  leaveModalVisible.value = true
}

const submitLeave = async (done: (val: boolean) => void) => {
  try {
    // API: PUT /api/v1/tenants/:tenant_id/employees/:id/leave
    // await axios.put(`/api/v1/tenants/1/employees/${currentLeaveRecord.value.id}/leave`, leaveForm.value)
    currentLeaveRecord.value.emp_status = 2
    currentLeaveRecord.value.status = 0
    Message.success('离职办理成功')
    leaveModalVisible.value = false
    done(true)
  } catch {
    Message.error('离职办理失败')
    done(false)
  }
}

const openUserLinkModal = (record: any) => {
  currentRecord.value = record
  linkUserForm.value.user_id = null
  linkUserVisible.value = true
}

const submitLinkUser = async (done: (val: boolean) => void) => {
  if (!linkUserForm.value.user_id) {
    Message.warning('请选择要关联的系统用户')
    done(false)
    return
  }
  try {
    // API: POST /api/v1/tenants/:tenant_id/employees/:id/link-user
    // await axios.post(`/api/v1/tenants/1/employees/${currentRecord.value.id}/link-user`, linkUserForm.value)
    const user = systemUsers.value.find(u => u.id === linkUserForm.value.user_id)
    Message.success(`已关联用户 ${user?.username}`)
    linkUserVisible.value = false
    done(true)
  } catch {
    Message.error('关联失败')
    done(false)
  }
}

onMounted(() => {
  loadEmployees()
})
</script>

<style scoped>
.page-container {
  padding: 20px;
}
.breadcrumb {
  margin-bottom: 12px;
}
.search-bar {
  margin-bottom: 12px;
}
.action-bar {
  margin-bottom: 16px;
}
</style>
