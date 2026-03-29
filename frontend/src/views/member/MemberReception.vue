<template>
  <div class="member-reception-page">
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>会员管理</a-breadcrumb-item>
      <a-breadcrumb-item>会员服务</a-breadcrumb-item>
      <a-breadcrumb-item>会员接待</a-breadcrumb-item>
    </a-breadcrumb>

    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="8"><a-card hoverable><a-statistic title="今日接待" :value="stats.todayCount || 0" /></a-card></a-col>
      <a-col :span="8"><a-card hoverable><a-statistic title="本周接待" :value="stats.weekCount || 0" /></a-card></a-col>
      <a-col :span="8"><a-card hoverable><a-statistic title="本月接待" :value="stats.monthCount || 0" /></a-card></a-col>
    </a-row>

    <a-card class="action-card">
      <a-space>
        <a-button type="primary" @click="showCreate">新建接待</a-button>
        <a-button @click="loadData">刷新</a-button>
      </a-space>
    </a-card>

    <a-card>
      <a-table :columns="columns" :data="recordList" :loading="loading" row-key="id" :pagination="{ pageSize: 10 }">
        <template #member="{ record }">
          <a-space>
            <a-avatar :size="24">{{ record.memberName?.charAt(0) || '?' }}</a-avatar>
            <span>{{ record.memberName || record.memberId }}</span>
          </a-space>
        </template>
      </a-table>
        <template #channel="{ record }"><a-tag>{{ record.channel || '-' }}</a-tag></template>
        <template #type="{ record }"><a-tag :color="getTypeColor(record.type)">{{ record.type }}</a-tag></template>
        <template #createTime="{ record }">{{ record.createTime?.slice(0, 19) }}</template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="showDetail(record)">详情</a-button>
            <a-button type="text" size="small" @click="showEdit(record)">编辑</a-button>
            <a-button type="text" size="small" status="danger" @click="handleDelete(record)">删除</a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 新增/编辑接待记录 -->
    <a-modal v-model:visible="formVisible" :title="isEdit ? '编辑接待' : '新建接待'"
      @before-ok="handleFormSubmit" @cancel="formVisible = false" :width="560" :loading="formLoading">
      <a-form :model="form" layout="vertical">
        <a-form-item label="会员" :rules="[{ required: true, message: '请输入会员ID或手机号' }]">
          <a-input v-model="form.memberId" placeholder="会员ID或手机号" />
        </a-form-item>
        <a-row :gutter="12">
          <a-col :span="12">
            <a-form-item label="接待方式">
              <a-select v-model="form.type" placeholder="选择方式">
                <a-option value="到店">到店</a-option>
                <a-option value="电话">电话</a-option>
                <a-option value="线上">线上</a-option>
                <a-option value="上门">上门</a-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="接待渠道">
              <a-input v-model="form.channel" placeholder="如：门店A" />
            </a-form-item>
          </a-col>
        </a-row>
        <a-form-item label="接待内容">
          <a-textarea v-model="form.content" :rows="4" placeholder="记录接待内容和会员反馈" />
        </a-form-item>
        <a-form-item label="跟进状态">
          <a-select v-model="form.followStatus" placeholder="选择状态">
            <a-option value="pending">待跟进</a-option>
            <a-option value="following">跟进中</a-option>
            <a-option value="resolved">已解决</a-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 详情弹窗 -->
    <a-modal v-model:visible="detailVisible" title="接待详情" :width="560" footer="">
      <a-descriptions :column="2" bordered>
        <a-descriptions-item label="会员">{{ detailRecord.memberName || detailRecord.memberId }}</a-descriptions-item>
        <a-descriptions-item label="接待方式">{{ detailRecord.type }}</a-descriptions-item>
        <a-descriptions-item label="渠道">{{ detailRecord.channel || '-' }}</a-descriptions-item>
        <a-descriptions-item label="状态">{{ detailRecord.followStatus }}</a-descriptions-item>
        <a-descriptions-item label="接待时间" :span="2">{{ detailRecord.createTime?.slice(0, 19) }}</a-descriptions-item>
        <a-descriptions-item label="接待内容" :span="2">{{ detailRecord.content || '-' }}</a-descriptions-item>
      </a-descriptions>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message, Modal } from '@arco-design/web-vue'
import * as api from '@/api/marketing'

const recordList = ref([])
const stats = ref({})
const loading = ref(false)
const formLoading = ref(false)
const formVisible = ref(false)
const detailVisible = ref(false)
const isEdit = ref(false)
const currentId = ref(null)
const detailRecord = ref({})

const form = reactive({ memberId: '', type: '到店', channel: '', content: '', followStatus: 'pending' })

const columns = [
  { title: '会员', slotName: 'member', width: 180 },
  { title: '接待方式', slotName: 'type', width: 100 },
  { title: '渠道', slotName: 'channel', width: 120 },
  { title: '跟进状态', dataIndex: 'followStatus', width: 120 },
  { title: '接待时间', slotName: 'createTime', width: 180 },
  { title: '操作', slotName: 'actions', width: 180 }
]

const getTypeColor = (t) => ({ 到店: 'blue', 电话: 'green', 线上: 'purple', 上门: 'orange' }[t] || 'gray')

const loadData = async () => {
  loading.value = true
  try {
    const res = await api.getMemberReceptionList()
    recordList.value = res.data?.list || []
    stats.value = res.data?.stats || {}
  } catch (err) { Message.error('加载数据失败: ' + err.message) } finally { loading.value = false }
}

const showCreate = () => { isEdit.value = false; Object.assign(form, { memberId: '', type: '到店', channel: '', content: '', followStatus: 'pending' }); formVisible.value = true }
const showEdit = (record) => { isEdit.value = true; currentId.value = record.id; Object.assign(form, { ...record }); formVisible.value = true }
const showDetail = (record) => { detailRecord.value = record; detailVisible.value = true }

const handleFormSubmit = async (done) => {
  formLoading.value = true
  try {
    if (isEdit.value) await api.updateMemberReception(currentId.value, { ...form })
    else await api.createMemberReception({ ...form })
    Message.success('保存成功'); formVisible.value = false; loadData(); done(true)
  } catch (err) { Message.error(err.message || '操作失败'); done(false) } finally { formLoading.value = false }
}

const handleDelete = (record) => {
  Modal.warning({ title: '确认删除', content: '确定删除该接待记录吗？', okText: '确认删除',
    onOk: async () => { try { await api.deleteMemberReception(record.id); Message.success('删除成功'); loadData() } catch (err) { Message.error(err.message || '删除失败') } }
  })
}

onMounted(() => loadData())
</script>

<style scoped>
.member-reception-page { padding: 20px 24px; min-height: calc(100vh - 64px); background: #f5f7fa; }
.breadcrumb { margin-bottom: 16px; }
.action-card { margin-bottom: 16px; }
</style>
