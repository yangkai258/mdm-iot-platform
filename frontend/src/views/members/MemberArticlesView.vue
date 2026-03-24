<template>
  <div class="member-page">
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>会员管理</a-breadcrumb-item>
      <a-breadcrumb-item>会员推文流水</a-breadcrumb-item>
    </a-breadcrumb>

    <a-row :gutter="16" class="stats-row">
      <a-col :span="8">
        <a-card class="stat-card">
          <a-statistic title="文章总数" :value="stats.total || 0" />
        </a-card>
      </a-col>
      <a-col :span="8">
        <a-card class="stat-card">
          <a-statistic title="今日发布" :value="stats.today || 0" />
        </a-card>
      </a-col>
      <a-col :span="8">
        <a-card class="stat-card">
          <a-statistic title="总阅读量" :value="stats.totalViews || 0" />
        </a-card>
      </a-col>
    </a-row>

    <a-card class="action-card">
      <a-space wrap>
        <a-input-search v-model="filters.keyword" placeholder="搜索文章标题" style="width: 240px" search-button @search="loadData" />
        <a-select v-model="filters.status" placeholder="状态" allow-clear style="width: 120px" @change="loadData">
          <a-option :value="1">已发布</a-option>
          <a-option :value="0">草稿</a-option>
        </a-select>
        <a-button type="primary" @click="showCreate">新建</a-button>
        <a-button @click="handleExport">导出</a-button>
        <a-button @click="loadData">刷新</a-button>
      </a-space>
    </a-card>

    <a-card class="table-card">
      <a-table :columns="columns" :data="dataList" :loading="loading" :pagination="paginationConfig" @page-change="onPageChange" row-key="id" :scroll="{ x: 1200 }">
        <template #status="{ record }">
          <a-tag :color="record.status === 1 ? 'green' : 'gray'">{{ record.status === 1 ? '已发布' : '草稿' }}</a-tag>
        </template>
        <template #actions="{ record }">
          <a-button type="text" size="small" @click="showEdit(record)">编辑</a-button>
          <a-button type="text" size="small" status="danger" @click="handleDelete(record)">删除</a-button>
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:visible="formVisible" :title="isEdit ? '编辑推文' : '新建推文'" :width="560px" @before-ok="handleSubmit" @cancel="formVisible = false">
      <a-form :model="form" layout="vertical">
        <a-form-item label="文章标题" required>
          <a-input v-model="form.title" placeholder="请输入文章标题" />
        </a-form-item>
        <a-form-item label="内容摘要">
          <a-textarea v-model="form.summary" placeholder="请输入内容摘要" :rows="3" />
        </a-form-item>
        <a-form-item label="文章内容">
          <a-textarea v-model="form.content" placeholder="请输入文章内容" :rows="5" />
        </a-form-item>
        <a-form-item label="推送渠道">
          <a-checkbox-group v-model="form.channels">
            <a-checkbox value="sms">短信</a-checkbox>
            <a-checkbox value="wechat">微信公众号</a-checkbox>
            <a-checkbox value="app">APP推送</a-checkbox>
          </a-checkbox-group>
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
const stats = ref({ total: 0, today: 0, totalViews: 0 })
const filters = reactive({ keyword: '', status: '' })
const paginationConfig = reactive({ current: 1, pageSize: 10, total: 0 })
const form = reactive({ id: null, title: '', summary: '', content: '', channels: ['wechat'], status: '1' })

const columns = [
  { title: '文章标题', dataIndex: 'title', width: 220 },
  { title: '内容摘要', dataIndex: 'summary', ellipsis: true },
  { title: '发布时间', dataIndex: 'publishTime', width: 180 },
  { title: '阅读量', dataIndex: 'views', width: 100 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '操作', slotName: 'actions', width: 150, fixed: 'right' }
]

const mockData = () => [
  { id: 1, title: '新会员专享福利来袭', summary: '新注册会员可领取专属优惠券，满100减20', publishTime: '2026-03-20 10:00:00', views: 2345, status: 1 },
  { id: 2, title: '生日月double福利', summary: '3月生日会员享受双倍积分和专属折扣', publishTime: '2026-03-15 10:00:00', views: 1890, status: 1 },
  { id: 3, title: '新品上市抢先看', summary: '春季新品预告，会员优先购买权', publishTime: '2026-03-10 10:00:00', views: 3200, status: 1 }
]

const loadData = () => {
  loading.value = true
  setTimeout(() => {
    dataList.value = mockData()
    stats.value = { total: mockData().length, today: 0, totalViews: mockData().reduce((s, d) => s + d.views, 0) }
    paginationConfig.total = dataList.value.length
    loading.value = false
  }, 400)
}

const onPageChange = (page) => { paginationConfig.current = page; loadData() }
const showCreate = () => { isEdit.value = false; Object.assign(form, { id: null, title: '', summary: '', content: '', channels: ['wechat'], status: '1' }); formVisible.value = true }
const showEdit = (record) => { isEdit.value = true; Object.assign(form, { ...record, status: String(record.status), channels: ['wechat'] }); formVisible.value = true }
const handleSubmit = (done) => {
  if (!form.title) { Message.error('请填写必填项'); done(false); return }
  setTimeout(() => { Message.success(isEdit.value ? '更新成功' : '创建成功'); formVisible.value = false; loadData(); done(true) }, 400)
}
const handleDelete = () => { Message.success('删除成功'); loadData() }
const handleExport = () => { Message.success('导出成功') }

loadData()
</script>

<style scoped>
.member-page { padding: 20px; }
.breadcrumb { margin-bottom: 16px; }
.stats-row { margin-bottom: 16px; }
.stat-card { text-align: center; }
.action-card { margin-bottom: 16px; }
</style>
