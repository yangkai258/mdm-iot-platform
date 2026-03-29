<template>
  <div class="page-container">
    <a-card class="general-card" title="会员积分">
      <template #extra>
        <a-button type="primary" @click="handleCreate"><icon-plus />新建</a-button>
      </template>
      <div class="search-form">
        <a-form :model="filters" layout="inline">
          <a-form-item label="关键词"><a-input v-model="filters.keyword" placeholder="请输入" /></a-form-item>
          <a-form-item>
            <a-button type="primary" @click="loadMembers">查询</a-button>
            <a-button @click="filters.keyword = ''; loadMembers()">重置</a-button>
          </a-form-item>
        </a-form>
      </div>
      <a-table :columns="columns" :data="memberList" :loading="loading" :pagination="pagination" row-key="id">
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

import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import { IconPlus } from '@arco-design/web-vue/es/icon'

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

const loadMembers = async () => {
  loading.value = true
  try {
    const params = { page: pagination.current, page_size: pagination.pageSize }
    if (filters.level) params.level = filters.level
    if (filters.keyword) params.keyword = filters.keyword
    memberList.value = []
    pagination.total = 0
  } catch (err) {
    Message.error('加载会员列表失败')
  } finally {
    loading.value = false
  }
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
  Message.success('积分调整成功')
  showAdjustDrawer.value = false
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
  Message.info('查看积分记录')
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
.page-container { padding: 16px; }
.search-form { margin-bottom: 16px; padding: 16px; background: var(--color-fill-lightest); border-radius: 4px; }
</style>
