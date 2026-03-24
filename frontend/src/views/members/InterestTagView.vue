<template>
  <div class="member-page">
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>会员管理</a-breadcrumb-item>
      <a-breadcrumb-item>兴趣分类标签</a-breadcrumb-item>
    </a-breadcrumb>

    <a-row :gutter="16" class="stats-row">
      <a-col :span="8">
        <a-card class="stat-card">
          <a-statistic title="标签总数" :value="stats.total || 0" />
        </a-card>
      </a-col>
      <a-col :span="8">
        <a-card class="stat-card">
          <a-statistic title="兴趣分类数" :value="stats.categoryCount || 0" />
        </a-card>
      </a-col>
      <a-col :span="8">
        <a-card class="stat-card">
          <a-statistic title="已打标会员" :value="stats.taggedMembers || 0" />
        </a-card>
      </a-col>
    </a-row>

    <a-card class="action-card">
      <a-space wrap>
        <a-input-search v-model="filters.keyword" placeholder="搜索标签名称" style="width: 220px" search-button @search="loadData" />
        <a-select v-model="filters.category" placeholder="兴趣分类" allow-clear style="width: 160px" @change="loadData">
          <a-option value="food">美食</a-option>
          <a-option value="travel">旅游</a-option>
          <a-option value="digital">数码</a-option>
          <a-option value="fashion">时尚</a-option>
          <a-option value="sports">运动</a-option>
        </a-select>
        <a-button type="primary" @click="showCreate">新建</a-button>
        <a-button @click="loadData">刷新</a-button>
      </a-space>
    </a-card>

    <a-card class="table-card">
      <a-table :columns="columns" :data="dataList" :loading="loading" :pagination="paginationConfig" @page-change="onPageChange" row-key="id" :scroll="{ x: 900 }">
        <template #category="{ record }">
          <a-tag>{{ categoryMap[record.category] || record.category }}</a-tag>
        </template>
        <template #status="{ record }">
          <a-tag :color="record.status === 1 ? 'green' : 'gray'">{{ record.status === 1 ? '启用' : '禁用' }}</a-tag>
        </template>
        <template #actions="{ record }">
          <a-button type="text" size="small" @click="showEdit(record)">编辑</a-button>
          <a-button type="text" size="small" status="danger" @click="handleDelete(record)">删除</a-button>
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:visible="formVisible" :title="isEdit ? '编辑兴趣标签' : '新建兴趣标签'" :width="480px" @before-ok="handleSubmit" @cancel="formVisible = false">
      <a-form :model="form" layout="vertical">
        <a-form-item label="标签名称" required>
          <a-input v-model="form.name" placeholder="请输入标签名称" />
        </a-form-item>
        <a-form-item label="兴趣分类" required>
          <a-select v-model="form.category" placeholder="请选择兴趣分类" style="width: 100%;">
            <a-option value="food">美食</a-option>
            <a-option value="travel">旅游</a-option>
            <a-option value="digital">数码</a-option>
            <a-option value="fashion">时尚</a-option>
            <a-option value="sports">运动</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="描述">
          <a-textarea v-model="form.description" placeholder="标签描述" :rows="3" />
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
const stats = ref({ total: 0, categoryCount: 5, taggedMembers: 0 })
const filters = reactive({ keyword: '', category: '' })
const paginationConfig = reactive({ current: 1, pageSize: 10, total: 0 })
const categoryMap = { food: '美食', travel: '旅游', digital: '数码', fashion: '时尚', sports: '运动' }
const form = reactive({ id: null, name: '', category: '', description: '', status: '1' })

const columns = [
  { title: '标签名称', dataIndex: 'name', width: 200 },
  { title: '兴趣分类', slotName: 'category', width: 150 },
  { title: '包含会员数', dataIndex: 'memberCount', width: 150 },
  { title: '描述', dataIndex: 'description', ellipsis: true },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '操作', slotName: 'actions', width: 150, fixed: 'right' }
]

const mockData = () => [
  { id: 1, name: '美食爱好者', category: 'food', memberCount: 2345, description: '喜欢品尝各类美食', status: 1 },
  { id: 2, name: '旅行达人', category: 'travel', memberCount: 1234, description: '热爱旅游出行', status: 1 },
  { id: 3, name: '数码玩家', category: 'digital', memberCount: 876, description: '数码产品爱好者', status: 1 }
]

const loadData = () => {
  loading.value = true
  setTimeout(() => {
    dataList.value = mockData()
    stats.value.taggedMembers = mockData().reduce((s, d) => s + d.memberCount, 0)
    paginationConfig.total = dataList.value.length
    loading.value = false
  }, 400)
}

const onPageChange = (page) => { paginationConfig.current = page; loadData() }
const showCreate = () => { isEdit.value = false; Object.assign(form, { id: null, name: '', category: '', description: '', status: '1' }); formVisible.value = true }
const showEdit = (record) => { isEdit.value = true; Object.assign(form, { ...record, status: String(record.status) }); formVisible.value = true }
const handleSubmit = (done) => {
  if (!form.name || !form.category) { Message.error('请填写必填项'); done(false); return }
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
