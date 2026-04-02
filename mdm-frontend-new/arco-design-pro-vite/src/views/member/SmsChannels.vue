<template>
    <Breadcrumb :items="['Home','Console','']" />


  <div class="sms-channels-page">
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>会员管理</a-breadcrumb-item>
      <a-breadcrumb-item>会员服务</a-breadcrumb-item>
      <a-breadcrumb-item>会员短信通道</a-breadcrumb-item>
    </a-breadcrumb>

    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="8"><a-card hoverable><a-statistic title="通道总数" :value="channelList.length" /></a-card></a-col>
      <a-col :span="8"><a-card hoverable><a-statistic title="本月发送" :value="stats.monthlySent || 0" /></a-card></a-col>
      <a-col :span="8"><a-card hoverable><a-statistic title="成功率" :value="(stats.successRate || 0).toFixed(1) + '%'" /></a-card></a-col>
    </a-row>

    <a-card class="action-card">
      <a-space>
        <a-button type="primary" @click="showCreate">新增通道</a-button>
        <a-button @click="loadData">刷新</a-button>
      </a-space>
    </a-card>

    <a-card>
      <a-table :columns="columns" :data="channelList" :loading="loading" row-key="id" :pagination="false">
        <template #name="{ record }">
          <a-space>
            <a-avatar :style="{ backgroundColor: record.status === 'active' ? '#00D084' : '#999' }" :size="24">✉</a-avatar>
            {{ record.name }}
          </a-space>
        </template>
      </a-table>
        <template #provider="{ record }"><a-tag>{{ record.provider }}</a-tag></template>
        <template #status="{ record }">
          <a-tag :color="record.status === 'active' ? 'green' : 'gray'">{{ record.status === 'active' ? '启用' : '禁用' }}</a-tag>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="showTest(record)">测试</a-button>
            <a-button type="text" size="small" @click="showEdit(record)">编辑</a-button>
            <a-button type="text" size="small" status="danger" @click="handleDelete(record)">删除</a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:visible="formVisible" :title="isEdit ? '编辑通道' : '新增通道'"
      @before-ok="handleFormSubmit" @cancel="formVisible = false" :width="520" :loading="formLoading">
      <a-form :model="form" layout="vertical">
        <a-form-item label="通道名称" :rules="[{ required: true, message: '请输入通道名称' }]">
          <a-input v-model="form.name" placeholder="如：阿里云短信" />
        </a-form-item>
        <a-form-item label="服务商">
          <a-select v-model="form.provider" placeholder="选择服务商">
            <a-option value="阿里云">阿里云</a-option>
            <a-option value="腾讯云">腾讯云</a-option>
            <a-option value="华为云">华为云</a-option>
            <a-option value="梦网">梦网</a-option>
            <a-option value="其他">其他</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="AccessKeyId">
          <a-input v-model="form.accessKeyId" placeholder="AccessKeyId" />
        </a-form-item>
        <a-form-item label="AccessKeySecret">
          <a-input v-model="form.accessKeySecret" placeholder="AccessKeySecret" type="password" />
        </a-form-item>
        <a-form-item label="签名">
          <a-input v-model="form.signName" placeholder="短信签名，如：【XXX公司】" />
        </a-form-item>
        <a-form-item label="模板ID">
          <a-input v-model="form.templateCode" placeholder="短信模板ID" />
        </a-form-item>
        <a-form-item label="状态">
          <a-switch v-model="form.status" checked-value="active" unchecked-value="inactive" />
        </a-form-item>
      </a-form>
    </a-modal>

    <a-modal v-model:visible="testVisible" title="测试短信通道" @ok="handleTest" :loading="testing" @cancel="testVisible = false">
      <a-form layout="vertical">
        <a-form-item label="接收手机号">
          <a-input v-model="testPhone" placeholder="请输入手机号" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message, Modal } from '@arco-design/web-vue'
import * as api from '@/api/marketing'

const channelList = ref([])
const stats = ref({})
const loading = ref(false)
const formLoading = ref(false)
const testing = ref(false)
const formVisible = ref(false)
const testVisible = ref(false)
const isEdit = ref(false)
const currentId = ref(null)
const testPhone = ref('')

const form = reactive({ name: '', provider: '阿里云', accessKeyId: '', accessKeySecret: '', signName: '', templateCode: '', status: 'inactive' })

const columns = [
  { title: '通道名称', slotName: 'name', width: 180 },
  { title: '服务商', slotName: 'provider', width: 120 },
  { title: '签名', dataIndex: 'signName', width: 160 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '操作', slotName: 'actions', width: 180 }
]

const loadData = async () => {
  loading.value = true
  try {
    const res = await api.getSmsChannelList()
    channelList.value = res.data || []
    stats.value = res.data?.stats || {}
  } catch (err) { Message.error('加载数据失败: ' + err.message) } finally { loading.value = false }
}

const showCreate = () => { isEdit.value = false; Object.assign(form, { name: '', provider: '阿里云', accessKeyId: '', accessKeySecret: '', signName: '', templateCode: '', status: 'inactive' }); formVisible.value = true }
const showEdit = (record) => { isEdit.value = true; currentId.value = record.id; Object.assign(form, { ...record }); formVisible.value = true }
const showTest = (record) => { currentId.value = record.id; testPhone.value = ''; testVisible.value = true }

const handleFormSubmit = async (done) => {
  if (!form.name) { Message.warning('请输入通道名称'); done(false); return }
  formLoading.value = true
  try {
    if (isEdit.value) await api.updateSmsChannel(currentId.value, { ...form })
    else await api.createSmsChannel({ ...form })
    Message.success('保存成功'); formVisible.value = false; loadData(); done(true)
  } catch (err) { Message.error(err.message || '操作失败'); done(false) } finally { formLoading.value = false }
}

const handleTest = async () => {
  if (!testPhone.value) { Message.warning('请输入手机号'); return }
  testing.value = true
  try {
    await api.testSmsChannel(currentId.value)
    Message.success('测试短信已发送')
    testVisible.value = false
  } catch (err) { Message.error('测试失败: ' + err.message) } finally { testing.value = false }
}

const handleDelete = (record) => {
  Modal.warning({ title: '确认删除', content: `确定删除通道「${record.name}」吗？`, okText: '确认删除',
    onOk: async () => { try { await api.deleteSmsChannel(record.id); Message.success('删除成功'); loadData() } catch (err) { Message.error(err.message || '删除失败') } }
  })
}

onMounted(() => loadData())
</script>

<style scoped>
.sms-channels-page { padding: 20px 24px; min-height: calc(100vh - 64px); background: #f5f7fa; }
.breadcrumb { margin-bottom: 16px; }
.action-card { margin-bottom: 16px; }
</style>
