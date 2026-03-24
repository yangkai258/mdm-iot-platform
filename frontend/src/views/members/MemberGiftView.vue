<template>
  <div class="member-page">
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>会员管理</a-breadcrumb-item>
      <a-breadcrumb-item>会员礼包</a-breadcrumb-item>
    </a-breadcrumb>

    <a-row :gutter="16" class="stats-row">
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="礼包总数" :value="stats.total || 0" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="发放总数" :value="stats.issued || 0" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="领取数" :value="stats.claimed || 0" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="总价值" :value="stats.totalValue || 0" suffix="元" />
        </a-card>
      </a-col>
    </a-row>

    <a-card class="action-card">
      <a-space wrap>
        <a-input-search v-model="filters.keyword" placeholder="搜索礼包名称" style="width: 220px" search-button @search="loadData" />
        <a-select v-model="filters.status" placeholder="状态" allow-clear style="width: 120px" @change="loadData">
          <a-option :value="1">启用</a-option>
          <a-option :value="0">禁用</a-option>
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
          <a-button type="text" size="small" @click="handleGrant(record)">发放</a-button>
          <a-button type="text" size="small" status="danger" @click="handleDelete(record)">删除</a-button>
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:visible="formVisible" :title="isEdit ? '编辑礼包' : '新建礼包'" :width="520px" @before-ok="handleSubmit" @cancel="formVisible = false">
      <a-form :model="form" layout="vertical">
        <a-form-item label="礼包名称" required>
          <a-input v-model="form.name" placeholder="请输入礼包名称" />
        </a-form-item>
        <a-form-item label="包含内容">
          <a-textarea v-model="form.content" placeholder="礼包包含内容，如：优惠券、积分、实物等" :rows="3" />
        </a-form-item>
        <a-form-item label="礼包价值（元）">
          <a-input-number v-model="form.value" :min="0" :precision="2" style="width: 100%;" />
        </a-form-item>
        <a-form-item label="有效期（天）">
          <a-input-number v-model="form.validDays" :min="1" style="width: 100%;" />
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
const stats = ref({ total: 0, issued: 0, claimed: 0, totalValue: 0 })
const filters = reactive({ keyword: '', status: '' })
const paginationConfig = reactive({ current: 1, pageSize: 10, total: 0 })
const form = reactive({ id: null, name: '', content: '', value: 0, validDays: 30, status: '1' })

const columns = [
  { title: '礼包名称', dataIndex: 'name', width: 200 },
  { title: '包含内容', dataIndex: 'content', ellipsis: true },
  { title: '价值（元）', dataIndex: 'value', width: 120 },
  { title: '已发放数', dataIndex: 'issuedCount', width: 120 },
  { title: '已领取数', dataIndex: 'claimedCount', width: 120 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '操作', slotName: 'actions', width: 200, fixed: 'right' }
]

const mockData = () => [
  { id: 1, name: '新会员礼包', content: '新人优惠券满100减20、100积分', value: 50, issuedCount: 234, claimedCount: 210, status: 1 },
  { id: 2, name: '生日专属礼包', content: '生日优惠券满200减50、双倍积分卡', value: 100, issuedCount: 567, claimedCount: 520, status: 1 },
  { id: 3, name: '节日礼包', content: '节日专属优惠券包', value: 80, issuedCount: 1234, claimedCount: 1100, status: 1 }
]

const loadData = () => {
  loading.value = true
  setTimeout(() => {
    dataList.value = mockData()
    const total = mockData()
    stats.value = {
      total: total.length,
      issued: total.reduce((s, d) => s + d.issuedCount, 0),
      claimed: total.reduce((s, d) => s + d.claimedCount, 0),
      totalValue: total.reduce((s, d) => s + d.value * d.issuedCount, 0)
    }
    paginationConfig.total = dataList.value.length
    loading.value = false
  }, 400)
}

const onPageChange = (page) => { paginationConfig.current = page; loadData() }
const showCreate = () => { isEdit.value = false; Object.assign(form, { id: null, name: '', content: '', value: 0, validDays: 30, status: '1' }); formVisible.value = true }
const showEdit = (record) => { isEdit.value = true; Object.assign(form, { ...record, status: String(record.status) }); formVisible.value = true }
const handleSubmit = (done) => {
  if (!form.name) { Message.error('请填写必填项'); done(false); return }
  setTimeout(() => { Message.success(isEdit.value ? '更新成功' : '创建成功'); formVisible.value = false; loadData(); done(true) }, 400)
}
const handleGrant = (record) => { Message.success(`已向符合条件的会员发放「${record.name}」`) }
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
