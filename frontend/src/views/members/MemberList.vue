<template>
  <div class="page-container">
    <div class="search-form">
      <a-form :model="form" layout="inline">
        <a-form-item label="名称"><a-input v-model="form.name" placeholder="请输入" /></a-form-item>
        <a-form-item>
          <a-button type="primary" @click="handleSearch">搜索</a-button>
          <a-button @click="handleReset">重置</a-button>
        </a-form-item>
      </a-form>
    </div>
    <div class="toolbar">
      <a-button type="primary" @click="handleCreate">新建</a-button>
    </div>
    <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id">
      <template #actions="{ record }">
        <a-button type="text" size="small" @click="handleEdit(record)">编辑</a-button>
        <a-button type="text" size="small" @click="handleDelete(record)">删除</a-button>
      </template>
    </a-table>
    <a-modal v-model:visible="modalVisible" :title="modalTitle" @before-ok="handleSubmit" @cancel="modalVisible = false">
      <a-form :model="form" label-col-flex="100px">
        <a-form-item label="名称"><a-input v-model="form.name" placeholder="请输入" /></a-form-item>
      </a-form>
      <template #footer>
        <a-button @click="modalVisible = false">取消</a-button>
        <a-button type="primary" @click="handleSubmit">确定</a-button>
      </template>
    </a-modal>
  </div>
</template>

<script setup>

import { ref, reactive, computed, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import * as api from '@/api/member'

const memberList = ref([])
const levelList = ref([])
const loading = ref(false)
const formLoading = ref(false)
const formVisible = ref(false)
const detailVisible = ref(false)
const levelModalVisible = ref(false)
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
  { title: '操作', slotName: 'actions', width: 180, fixed: 'right' }
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
  } catch (err) {
    // ignore
  }
}

const loadStats = async () => {
  try {
    const res = await api.getMemberList({ page: 1, pageSize: 1 })
    stats.total = res.data?.total || 0
    // 估算今日新增和活跃会员（后端未提供专门接口时用列表数据）
    stats.todayNew = 0
    stats.active = stats.total
    stats.totalPoints = 0
  } catch (err) {
    // ignore
  }
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

onMounted(() => {
  loadMembers()
  loadLevels()
  loadStats()
})

</script>

<style scoped>
.page-container { background: #fff; border-radius: 4px; padding: 20px; }
.search-form { margin-bottom: 16px; padding: 16px; background: #f7f8fa; border-radius: 4px; }
.toolbar { margin-bottom: 16px; }
</style>
