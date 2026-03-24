<template>
  <div class="member-page">
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>会员管理</a-breadcrumb-item>
      <a-breadcrumb-item>短信模板设置</a-breadcrumb-item>
    </a-breadcrumb>

    <a-row :gutter="16" class="stats-row">
      <a-col :span="8">
        <a-card class="stat-card">
          <a-statistic title="模板总数" :value="stats.total || 0" />
        </a-card>
      </a-col>
      <a-col :span="8">
        <a-card class="stat-card">
          <a-statistic title="已启用" :value="stats.enabled || 0" />
        </a-card>
      </a-col>
      <a-col :span="8">
        <a-card class="stat-card">
          <a-statistic title="本月发送" :value="stats.monthSent || 0" />
        </a-card>
      </a-col>
    </a-row>

    <a-card class="action-card">
      <a-space wrap>
        <a-input-search v-model="filters.keyword" placeholder="搜索模板名称" style="width: 220px" search-button @search="loadData" />
        <a-select v-model="filters.type" placeholder="模板类型" allow-clear style="width: 140px" @change="loadData">
          <a-option value="marketing">营销类</a-option>
          <a-option value="notice">通知类</a-option>
          <a-option value="verify">验证类</a-option>
        </a-select>
        <a-button type="primary" @click="showCreate">新建</a-button>
        <a-button @click="loadData">刷新</a-button>
      </a-space>
    </a-card>

    <a-card class="table-card">
      <a-table :columns="columns" :data="dataList" :loading="loading" :pagination="paginationConfig" @page-change="onPageChange" row-key="id" :scroll="{ x: 1000 }">
        <template #status="{ record }">
          <a-tag :color="record.status === 1 ? 'green' : 'gray'">{{ record.status === 1 ? '启用' : '禁用' }}</a-tag>
        </template>
        <template #actions="{ record }">
          <a-button type="text" size="small" @click="showEdit(record)">编辑</a-button>
          <a-button type="text" size="small" status="danger" @click="handleDelete(record)">删除</a-button>
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:visible="formVisible" :title="isEdit ? '编辑短信模板' : '新建短信模板'" :width="520px" @before-ok="handleSubmit" @cancel="formVisible = false">
      <a-form :model="form" layout="vertical">
        <a-form-item label="模板名称" required>
          <a-input v-model="form.name" placeholder="请输入模板名称" />
        </a-form-item>
        <a-form-item label="模板类型">
          <a-select v-model="form.type" placeholder="请选择模板类型" style="width: 100%;">
            <a-option value="marketing">营销类</a-option>
            <a-option value="notice">通知类</a-option>
            <a-option value="verify">验证类</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="模板内容" required>
          <a-textarea v-model="form.content" placeholder="请输入短信内容，支持变量：{name}会员姓名、{points}积分、{coupon}优惠券" :rows="4" />
          <div style="color:#999;font-size:12px;margin-top:4px;">
            可用变量：{name}会员姓名、{points}积分余额、{coupon}优惠券名称、{date}日期
          </div>
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
const stats = ref({ total: 0, enabled: 0, monthSent: 0 })
const filters = reactive({ keyword: '', type: '' })
const paginationConfig = reactive({ current: 1, pageSize: 10, total: 0 })
const form = reactive({ id: null, name: '', type: 'notice', content: '', status: '1' })

const columns = [
  { title: '模板名称', dataIndex: 'name', width: 200 },
  { title: '模板类型', dataIndex: 'typeName', width: 120 },
  { title: '模板内容', dataIndex: 'content', ellipsis: true },
  { title: '变量', dataIndex: 'variables', width: 200 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '操作', slotName: 'actions', width: 150, fixed: 'right' }
]

const typeNameMap = { marketing: '营销类', notice: '通知类', verify: '验证类' }

const mockData = () => [
  { id: 1, name: '生日祝福', type: 'notice', typeName: '通知类', content: '亲爱的{name}，祝您生日快乐！本月可享受双倍积分，快来选购心仪商品吧~', variables: '{name}', status: 1 },
  { id: 2, name: '新会员欢迎', type: 'marketing', typeName: '营销类', content: '欢迎{name}加入会员大家庭！您已获得{coupon}，首单满100减20，立即使用！', variables: '{name},{coupon}', status: 1 },
  { id: 3, name: '积分到期提醒', type: 'notice', typeName: '通知类', content: '亲爱的{name}，您有{points}积分将于月底过期，请及时使用！', variables: '{name},{points}', status: 1 }
]

const loadData = () => {
  loading.value = true
  setTimeout(() => {
    dataList.value = mockData()
    stats.value = { total: mockData().length, enabled: mockData().filter(d => d.status === 1).length, monthSent: 1234 }
    paginationConfig.total = dataList.value.length
    loading.value = false
  }, 400)
}

const onPageChange = (page) => { paginationConfig.current = page; loadData() }
const showCreate = () => { isEdit.value = false; Object.assign(form, { id: null, name: '', type: 'notice', content: '', status: '1' }); formVisible.value = true }
const showEdit = (record) => { isEdit.value = true; Object.assign(form, { ...record, status: String(record.status) }); formVisible.value = true }
const handleSubmit = (done) => {
  if (!form.name || !form.content) { Message.error('请填写必填项'); done(false); return }
  setTimeout(() => { Message.success(isEdit.value ? '更新成功' : '创建成功'); formVisible.value = false; loadData(); done(true) }, 400)
}
const handleDelete = () => { Message.success('删除成功'); loadData() }

loadData()
</script>