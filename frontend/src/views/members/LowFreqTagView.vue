<template>
  <div class="member-page">
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>会员管理</a-breadcrumb-item>
      <a-breadcrumb-item>低频购买标签</a-breadcrumb-item>
    </a-breadcrumb>

    <a-row :gutter="16" class="stats-row">
      <a-col :span="8">
        <a-card class="stat-card">
          <a-statistic title="标签总数" :value="stats.total || 0" />
        </a-card>
      </a-col>
      <a-col :span="8">
        <a-card class="stat-card">
          <a-statistic title="覆盖会员数" :value="stats.coverMembers || 0" />
        </a-card>
      </a-col>
      <a-col :span="8">
        <a-card class="stat-card">
          <a-statistic title="沉默会员数" :value="stats.silentMembers || 0" />
        </a-card>
      </a-col>
    </a-row>

    <a-card class="action-card">
      <a-space wrap>
        <a-input-search v-model="filters.keyword" placeholder="搜索标签名称" style="width: 220px" search-button @search="loadData" />
        <a-button type="primary" @click="showCreate">新建</a-button>
        <a-button @click="loadData">刷新</a-button>
      </a-space>
    </a-card>

    <a-card class="table-card">
      <a-table :columns="columns" :data="dataList" :loading="loading" :pagination="paginationConfig" @page-change="onPageChange" row-key="id" :scroll="{ x: 900 }">
        <template #status="{ record }">
          <a-tag :color="record.status === 1 ? 'green' : 'gray'">{{ record.status === 1 ? '启用' : '禁用' }}</a-tag>
        </template>
        <template #actions="{ record }">
          <a-button type="text" size="small" @click="showEdit(record)">编辑</a-button>
          <a-button type="text" size="small" status="danger" @click="handleDelete(record)">删除</a-button>
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:visible="formVisible" :title="isEdit ? '编辑低频标签' : '新建低频标签'" :width="480px" @before-ok="handleSubmit" @cancel="formVisible = false">
      <a-form :model="form" layout="vertical">
        <a-form-item label="标签名称" required>
          <a-input v-model="form.name" placeholder="请输入标签名称" />
        </a-form-item>
        <a-form-item label="消费间隔天数" required>
          <a-input-number v-model="form.intervalDays" :min="1" placeholder="多少天未消费" style="width: 100%;" />
          <div style="color:#999;font-size:12px;margin-top:4px;">超过此天数未消费的会员将被打上此标签</div>
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
const stats = ref({ total: 0, coverMembers: 0, silentMembers: 0 })
const filters = reactive({ keyword: '' })
const paginationConfig = reactive({ current: 1, pageSize: 10, total: 0 })
const form = reactive({ id: null, name: '', intervalDays: 90, description: '', status: '1' })

const columns = [
  { title: '标签名称', dataIndex: 'name', width: 200 },
  { title: '定义条件（消费间隔天数）', dataIndex: 'intervalDays', width: 220 },
  { title: '包含会员数', dataIndex: 'memberCount', width: 150 },
  { title: '描述', dataIndex: 'description', ellipsis: true },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '操作', slotName: 'actions', width: 150, fixed: 'right' }
]

const mockData = () => [
  { id: 1, name: '沉睡会员', intervalDays: 90, memberCount: 3456, description: '90天未消费', status: 1 },
  { id: 2, name: '流失风险', intervalDays: 60, memberCount: 1234, description: '60天未消费', status: 1 },
  { id: 3, name: '边缘客户', intervalDays: 30, memberCount: 567, description: '30天未消费', status: 1 }
]

const loadData = () => {
  loading.value = true
  setTimeout(() => {
    dataList.value = mockData()
    stats.value = { total: mockData().length, coverMembers: mockData().reduce((s, d) => s + d.memberCount, 0), silentMembers: mockData()[0].memberCount }
    paginationConfig.total = dataList.value.length
    loading.value = false
  }, 400)
}

const onPageChange = (page) => { paginationConfig.current = page; loadData() }
const showCreate = () => { isEdit.value = false; Object.assign(form, { id: null, name: '', intervalDays: 90, description: '', status: '1' }); formVisible.value = true }
const showEdit = (record) => { isEdit.value = true; Object.assign(form, { ...record, status: String(record.status) }); formVisible.value = true }
const handleSubmit = (done) => {
  if (!form.name || !form.intervalDays) { Message.error('请填写必填项'); done(false); return }
  setTimeout(() => { Message.success(isEdit.value ? '更新成功' : '创建成功'); formVisible.value = false; loadData(); done(true) }, 400)
}
const handleDelete = (record) => { Message.success('删除成功'); loadData() }

loadData()
</script>

<style scoped>
.member-page { padding: 20px; }
.breadcrumb { margin-bottom: 16px; }
.stats-row { margin-bottom: 16px; }
.stat-card { text-align: center; }
.action-card { margin-bottom: 16px; }
</style>
