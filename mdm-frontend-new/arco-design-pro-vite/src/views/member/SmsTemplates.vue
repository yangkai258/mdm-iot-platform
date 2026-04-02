<template>
    <Breadcrumb :items="['Home','Console','']" />


  <div class="sms-templates-page">
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>会员管理</a-breadcrumb-item>
      <a-breadcrumb-item>会员服务</a-breadcrumb-item>
      <a-breadcrumb-item>短信模板设置</a-breadcrumb-item>
    </a-breadcrumb>

    <a-card class="action-card">
      <a-space>
        <a-button type="primary" @click="showCreate">新建模板</a-button>
        <a-button @click="loadData">刷新</a-button>
      </a-space>
    </a-card>

    <a-card>
      <a-table :columns="columns" :data="templateList" :loading="loading" row-key="id" :pagination="{ pageSize: 10 }">
        <template #name="{ record }"><a-tag color="arcoblue">{{ record.name }}</a-tag></template>
        <template #type="{ record }"><a-tag>{{ record.type }}</a-tag></template>
      </a-table>
        <template #content="{ record }">
          <span style="color: #666; font-size: 13px;">{{ record.content?.slice(0, 50) }}{{ record.content?.length > 50 ? '...' : '' }}</span>
        </template>
        <template #status="{ record }"><a-tag :color="record.status === 'active' ? 'green' : 'gray'">{{ record.status === 'active' ? '已启用' : '已禁用' }}</a-tag></template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="showEdit(record)">编辑</a-button>
            <a-button type="text" size="small" status="danger" @click="handleDelete(record)">删除</a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:visible="formVisible" :title="isEdit ? '编辑模板' : '新建模板'"
      @before-ok="handleFormSubmit" @cancel="formVisible = false" :width="600" :loading="formLoading">
      <a-form :model="form" layout="vertical">
        <a-form-item label="模板名称" :rules="[{ required: true, message: '请输入模板名称' }]">
          <a-input v-model="form.name" placeholder="如：入会欢迎短信" />
        </a-form-item>
        <a-form-item label="模板类型">
          <a-select v-model="form.type" placeholder="选择类型">
            <a-option value="入会">入会</a-option>
            <a-option value="生日">生日</a-option>
            <a-option value="活动">活动</a-option>
            <a-option value="节日">节日</a-option>
            <a-option value="积分">积分</a-option>
            <a-option value="促销">促销</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="短信内容" :rules="[{ required: true, message: '请输入短信内容' }]">
          <a-textarea v-model="form.content" :rows="4" placeholder="请输入短信内容，支持变量：{name}会员姓名、{level}等级、{points}积分等" />
        </a-form-item>
        <a-form-item label="状态"><a-switch v-model="form.status" checked-value="active" unchecked-value="inactive" /></a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message, Modal } from '@arco-design/web-vue'
import * as api from '@/api/marketing'

const templateList = ref([])
const loading = ref(false)
const formLoading = ref(false)
const formVisible = ref(false)
const isEdit = ref(false)
const currentId = ref(null)

const form = reactive({ name: '', type: '入会', content: '', status: 'active' })

const columns = [
  { title: '模板名称', slotName: 'name', width: 160 },
  { title: '类型', slotName: 'type', width: 100 },
  { title: '短信内容', slotName: 'content' },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '操作', slotName: 'actions', width: 140 }
]

const loadData = async () => {
  loading.value = true
  try {
    const res = await api.getSmsTemplateList()
    templateList.value = res.data || []
  } catch (err) { Message.error('加载数据失败: ' + err.message) } finally { loading.value = false }
}

const showCreate = () => { isEdit.value = false; Object.assign(form, { name: '', type: '入会', content: '', status: 'active' }); formVisible.value = true }
const showEdit = (record) => { isEdit.value = true; currentId.value = record.id; Object.assign(form, { ...record }); formVisible.value = true }

const handleFormSubmit = async (done) => {
  if (!form.name || !form.content) { Message.warning('请填写完整信息'); done(false); return }
  formLoading.value = true
  try {
    if (isEdit.value) await api.updateSmsTemplate(currentId.value, { ...form })
    else await api.createSmsTemplate({ ...form })
    Message.success('保存成功'); formVisible.value = false; loadData(); done(true)
  } catch (err) { Message.error(err.message || '操作失败'); done(false) } finally { formLoading.value = false }
}

const handleDelete = (record) => {
  Modal.warning({ title: '确认删除', content: `确定删除模板「${record.name}」吗？`, okText: '确认删除',
    onOk: async () => { try { await api.deleteSmsTemplate(record.id); Message.success('删除成功'); loadData() } catch (err) { Message.error(err.message || '删除失败') } }
  })
}

onMounted(() => loadData())
</script>

<style scoped>
.sms-templates-page { padding: 20px 24px; min-height: calc(100vh - 64px); background: #f5f7fa; }
.breadcrumb { margin-bottom: 16px; }
.action-card { margin-bottom: 16px; }
</style>
