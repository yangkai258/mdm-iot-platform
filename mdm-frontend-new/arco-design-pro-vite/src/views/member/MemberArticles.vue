<template>
    <Breadcrumb :items="['Home','Console','']" />


  <div class="member-articles-page">
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>会员管理</a-breadcrumb-item>
      <a-breadcrumb-item>会员服务</a-breadcrumb-item>
      <a-breadcrumb-item>会员推文流水</a-breadcrumb-item>
    </a-breadcrumb>

    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="8"><a-card hoverable><a-statistic title="本月推送" :value="stats.monthlyCount || 0" /></a-card></a-col>
      <a-col :span="8"><a-card hoverable><a-statistic title="已阅读" :value="stats.readCount || 0" /></a-card></a-col>
      <a-col :span="8"><a-card hoverable><a-statistic title="阅读率" :value="(stats.readRate || 0).toFixed(1) + '%'" /></a-card></a-col>
    </a-row>

    <a-card class="action-card">
      <a-space>
        <a-button type="primary" @click="showCreate">新建推文</a-button>
        <a-button @click="loadData">刷新</a-button>
      </a-space>
    </a-card>

    <a-card>
      <a-table :columns="columns" :data="articleList" :loading="loading" row-key="id" :pagination="{ pageSize: 10 }">
        <template #title="{ record }">
          <a-space direction="vertical" :size="0">
            <span style="font-weight: 500;">{{ record.title }}</span>
            <span style="color: #999; font-size: 12px;">{{ record.summary }}</span>
          </a-space>
        </template>
      </a-table>
        <template #type="{ record }"><a-tag :color="getTypeColor(record.type)">{{ record.type }}</a-tag></template>
        <template #sendTime="{ record }">{{ record.sendTime?.slice(0, 19) || '待发送' }}</template>
        <template #status="{ record }"><a-tag :color="getStatusColor(record.status)">{{ getStatusName(record.status) }}</a-tag></template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="showEdit(record)">编辑</a-button>
            <a-button type="text" size="small" status="danger" @click="handleDelete(record)">删除</a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:visible="formVisible" :title="isEdit ? '编辑推文' : '新建推文'"
      @before-ok="handleFormSubmit" @cancel="formVisible = false" :width="600" :loading="formLoading">
      <a-form :model="form" layout="vertical">
        <a-form-item label="标题" :rules="[{ required: true, message: '请输入标题' }]">
          <a-input v-model="form.title" placeholder="推文标题" />
        </a-form-item>
        <a-form-item label="摘要">
          <a-input v-model="form.summary" placeholder="推文摘要" />
        </a-form-item>
        <a-form-item label="推文类型">
          <a-select v-model="form.type" placeholder="选择类型">
            <a-option value="活动">活动</a-option>
            <a-option value="通知">通知</a-option>
            <a-option value="福利">福利</a-option>
            <a-option value="资讯">资讯</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="推送内容">
          <a-textarea v-model="form.content" :rows="4" placeholder="推文内容" />
        </a-form-item>
        <a-form-item label="封面图URL">
          <a-input v-model="form.coverUrl" placeholder="图片URL" />
        </a-form-item>
        <a-form-item label="发送时间">
          <a-date-picker v-model="form.sendTime" show-time style="width: 100%;" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message, Modal } from '@arco-design/web-vue'
import * as api from '@/api/marketing'

const articleList = ref([])
const stats = ref({})
const loading = ref(false)
const formLoading = ref(false)
const formVisible = ref(false)
const isEdit = ref(false)
const currentId = ref(null)

const form = reactive({ title: '', summary: '', type: '活动', content: '', coverUrl: '', sendTime: null })

const columns = [
  { title: '标题', slotName: 'title' },
  { title: '类型', slotName: 'type', width: 100 },
  { title: '发送时间', slotName: 'sendTime', width: 180 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '操作', slotName: 'actions', width: 140 }
]

const getTypeColor = (t) => ({ 活动: 'blue', 通知: 'orange', 福利: 'green', 资讯: 'purple' }[t] || 'gray')
const getStatusColor = (s) => ({ sent: 'green', pending: 'orange', draft: 'gray' }[s] || 'gray')
const getStatusName = (s) => ({ sent: '已发送', pending: '待发送', draft: '草稿' }[s] || s)

const loadData = async () => {
  loading.value = true
  try {
    const res = await api.getMemberArticleList()
    articleList.value = res.data?.list || []
    stats.value = res.data?.stats || {}
  } catch (err) { Message.error('加载数据失败: ' + err.message) } finally { loading.value = false }
}

const showCreate = () => { isEdit.value = false; Object.assign(form, { title: '', summary: '', type: '活动', content: '', coverUrl: '', sendTime: null }); formVisible.value = true }
const showEdit = (record) => { isEdit.value = true; currentId.value = record.id; Object.assign(form, { ...record }); formVisible.value = true }

const handleFormSubmit = async (done) => {
  if (!form.title) { Message.warning('请输入标题'); done(false); return }
  formLoading.value = true
  try {
    if (isEdit.value) await api.createMemberArticle({ id: currentId.value, ...form })
    else await api.createMemberArticle({ ...form })
    Message.success('保存成功'); formVisible.value = false; loadData(); done(true)
  } catch (err) { Message.error(err.message || '操作失败'); done(false) } finally { formLoading.value = false }
}

const handleDelete = (record) => {
  Modal.warning({ title: '确认删除', content: `确定删除推文「${record.title}」吗？`, okText: '确认删除',
    onOk: async () => { try { await api.deleteMemberArticle(record.id); Message.success('删除成功'); loadData() } catch (err) { Message.error(err.message || '删除失败') } }
  })
}

onMounted(() => loadData())
</script>

<style scoped>
.member-articles-page { padding: 20px 24px; min-height: calc(100vh - 64px); background: #f5f7fa; }
.breadcrumb { margin-bottom: 16px; }
.action-card { margin-bottom: 16px; }
</style>
