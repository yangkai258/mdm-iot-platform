<template>
  <div class="member-page">
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>会员管理</a-breadcrumb-item>
      <a-breadcrumb-item>会员权益管理</a-breadcrumb-item>
    </a-breadcrumb>

    <a-row :gutter="16" class="stats-row">
      <a-col :span="8">
        <a-card class="stat-card">
          <a-statistic title="权益总数" :value="stats.total || 0" />
        </a-card>
      </a-col>
      <a-col :span="8">
        <a-card class="stat-card">
          <a-statistic title="已激活" :value="stats.active || 0" />
        </a-card>
      </a-col>
      <a-col :span="8">
        <a-card class="stat-card">
          <a-statistic title="覆盖会员" :value="stats.coverMembers || 0" />
        </a-card>
      </a-col>
    </a-row>

    <a-card class="action-card">
      <a-space wrap>
        <a-input-search v-model="filters.keyword" placeholder="搜索权益名称" style="width: 220px" search-button @search="loadData" />
        <a-select v-model="filters.levelId" placeholder="适用等级" allow-clear style="width: 140px" @change="loadData">
          <a-option v-for="lv in levelOptions" :key="lv.id" :value="lv.id">{{ lv.name }}</a-option>
        </a-select>
        <a-button type="primary" @click="showCreate">新建</a-button>
        <a-button @click="loadData">刷新</a-button>
      </a-space>
    </a-card>

    <a-card class="table-card">
      <a-table :columns="columns" :data="dataList" :loading="loading" :pagination="paginationConfig" @page-change="onPageChange" row-key="id" :scroll="{ x: 1100 }">
        <template #status="{ record }">
          <a-tag :color="record.status === 1 ? 'green' : 'gray'">{{ record.status === 1 ? '启用' : '禁用' }}</a-tag>
        </template>
        <template #actions="{ record }">
          <a-button type="text" size="small" @click="showEdit(record)">编辑</a-button>
          <a-button type="text" size="small" status="danger" @click="handleDelete(record)">删除</a-button>
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:visible="formVisible" :title="isEdit ? '编辑权益' : '新建权益'" :width="520px" @before-ok="handleSubmit" @cancel="formVisible = false">
      <a-form :model="form" layout="vertical">
        <a-form-item label="权益名称" required>
          <a-input v-model="form.name" placeholder="请输入权益名称" />
        </a-form-item>
        <a-form-item label="适用等级" required>
          <a-select v-model="form.levelId" placeholder="请选择适用等级" style="width: 100%;">
            <a-option v-for="lv in levelOptions" :key="lv.id" :value="lv.id">{{ lv.name }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="权益内容">
          <a-textarea v-model="form.content" placeholder="请输入权益内容" :rows="3" />
        </a-form-item>
        <a-form-item label="权益图标">
          <a-input v-model="form.icon" placeholder="请输入图标名称或URL" />
        </a-form-item>
        <a-form-item label="排序">
          <a-input-number v-model="form.sort" :min="0" style="width: 100%;" />
        </a-form-item>
        <a-form-item label="状态">
          <a-switch v-model="form.status" checked-value="1" unchecked-value="0" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { Message } from '@arco-design/web-vue'

const loading = ref(false)
const formVisible = ref(false)
const isEdit = ref(false)
const dataList = ref([])
const levelOptions = ref([
  { id: 1, name: '普通会员' },
  { id: 2, name: '银卡会员' },
  { id: 3, name: '金卡会员' },
  { id: 4, name: '钻石会员' }
])
const stats = ref({ total: 0, active: 0, coverMembers: 0 })
const filters = reactive({ keyword: '', levelId: '' })
const paginationConfig = reactive({ current: 1, pageSize: 10, total: 0 })
const form = reactive({ id: null, name: '', levelId: '', content: '', icon: '', sort: 0, status: '1' })

const columns = [
  { title: '权益名称', dataIndex: 'name', width: 180 },
  { title: '适用等级', dataIndex: 'levelName', width: 150 },
  { title: '权益内容', dataIndex: 'content', ellipsis: true },
  { title: '排序', dataIndex: 'sort', width: 80 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '操作', slotName: 'actions', width: 150, fixed: 'right' }
]

const mockData = () => [
  { id: 1, name: '生日双倍积分', levelId: 2, levelName: '银卡会员', content: '生日当月消费享受双倍积分', sort: 1, status: 1 },
  { id: 2, name: '专属折扣', levelId: 3, levelName: '金卡会员', content: '全场商品9折优惠', sort: 2, status: 1 },
  { id: 3, name: '免费配送', levelId: 4, levelName: '钻石会员', content: '全年无限次免费配送', sort: 3, status: 1 },
  { id: 4, name: '优先客服', levelId: 1, levelName: '普通会员', content: '享受优先客服通道', sort: 10, status: 1 }
]

const loadData = () => {
  loading.value = true
  setTimeout(() => {
    dataList.value = mockData()
    stats.value = { total: mockData().length, active: mockData().filter(d => d.status === 1).length, coverMembers: 5678 }
    paginationConfig.total = dataList.value.length
    loading.value = false
  }, 400)
}

const onPageChange = (page) => { paginationConfig.current = page; loadData() }
const showCreate = () => { isEdit.value = false; Object.assign(form, { id: null, name: '', levelId: '', content: '', icon: '', sort: 0, status: '1' }); formVisible.value = true }
const showEdit = (record) => { isEdit.value = true; Object.assign(form, { ...record, status: String(record.status) }); formVisible.value = true }
const handleSubmit = (done) => {
  if (!form.name || !form.levelId) { Message.error('请填写必填项'); done(false); return }
  setTimeout(() => { Message.success(isEdit.value ? '更新成功' : '创建成功'); formVisible.value = false; loadData(); done(true) }, 400)
}
const handleDelete = () => { Message.success('删除成功'); loadData() }

loadData()
</script>

<style scoped>
.member-page { padding: 20px; }
.breadcrumb { margin-bottom: 16px; }
.stats-row { margin-bottom: 16px; }
.stat-card { text-align: center; }
.action-card { margin-bottom: 16px; }
</style>
