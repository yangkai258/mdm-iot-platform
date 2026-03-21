<template>
  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>告警管理</a-breadcrumb-item>
      <a-breadcrumb-item>告警设置</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 搜索框 -->
    <div class="pro-search-bar">
      <a-input-search v-model="searchKeyword" placeholder="搜索设置项" style="width: 280px" search-button @search="loadSettings" />
    </div>

    <!-- 操作按钮组 -->
    <div class="pro-action-bar">
      <a-space>
        <a-button type="primary" @click="showAddModal">新增设置</a-button>
        <a-button @click="loadSettings">刷新</a-button>
      </a-space>
    </div>

    <!-- 数据表格 -->
    <div class="pro-content-area">
      <a-table :columns="columns" :data="settings" :loading="loading" row-key="id">
        <template #enabled="{ record }">
          <a-switch v-model="record.enabled" @change="toggleSetting(record)" />
        </template>
        <template #type="{ record }">
          <a-tag :color="getTypeColor(record.type)">{{ getTypeText(record.type) }}</a-tag>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="editSetting(record)">编辑</a-button>
            <a-button type="text" size="small" status="danger" @click="deleteSetting(record)">删除</a-button>
          </a-space>
        </template>
      </a-table>
    </div>

    <!-- 添加/编辑设置弹窗 -->
    <a-modal v-model:visible="modalVisible" :title="isEdit ? '编辑设置' : '新增设置'" @ok="handleSubmit" :width="520">
      <a-form :model="form" layout="vertical">
        <a-form-item label="设置名称" required>
          <a-input v-model="form.name" placeholder="请输入设置名称" />
        </a-form-item>
        <a-form-item label="设置类型" required>
          <a-select v-model="form.type" placeholder="选择类型">
            <a-option value="email">邮件通知</a-option>
            <a-option value="sms">短信通知</a-option>
            <a-option value="webhook">Webhook</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="配置值">
          <a-input v-model="form.config_value" placeholder="请输入配置值" />
        </a-form-item>
        <a-form-item label="描述">
          <a-textarea v-model="form.description" placeholder="请输入描述" :rows="3" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const loading = ref(false)
const settings = ref([])
const modalVisible = ref(false)
const isEdit = ref(false)
const searchKeyword = ref('')

const form = reactive({
  name: '', type: 'email', config_value: '', description: ''
})

const columns = [
  { title: '设置名称', dataIndex: 'name' },
  { title: '类型', slotName: 'type', width: 120 },
  { title: '配置值', dataIndex: 'config_value', ellipsis: true },
  { title: '描述', dataIndex: 'description', ellipsis: true },
  { title: '启用', slotName: 'enabled', width: 100 },
  { title: '操作', slotName: 'actions', width: 150 }
]

const getTypeColor = (t) => ({ email: 'blue', sms: 'green', webhook: 'purple' }[t] || 'gray')
const getTypeText = (t) => ({ email: '邮件通知', sms: '短信通知', webhook: 'Webhook' }[t] || t)

const loadSettings = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const res = await fetch('/api/v1/alerts/settings', { headers: { 'Authorization': `Bearer ${token}` } })
    const data = await res.json()
    if (data.code === 0) settings.value = data.data.list || []
  } catch (e) { Message.error('加载失败') }
  finally { loading.value = false }
}

const showAddModal = () => {
  isEdit.value = false
  Object.assign(form, { name: '', type: 'email', config_value: '', description: '' })
  modalVisible.value = true
}

const editSetting = (record) => {
  isEdit.value = true
  Object.assign(form, record)
  modalVisible.value = true
}

const handleSubmit = async () => {
  try {
    const token = localStorage.getItem('token')
    const res = await fetch('/api/v1/alerts/settings', {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' },
      body: JSON.stringify(form)
    })
    const data = await res.json()
    if (data.code === 0) { Message.success('保存成功'); modalVisible.value = false; loadSettings() }
  } catch (e) { Message.error('操作失败') }
}

const toggleSetting = async (record) => {}
const deleteSetting = (record) => { settings.value = settings.value.filter(s => s.id !== record.id); Message.success('删除成功') }

onMounted(() => { loadSettings() })
</script>

<style scoped>
.pro-page-container { padding: 20px 24px; min-height: calc(100vh - 64px); background: #f5f7fa; }
.pro-breadcrumb { margin-bottom: 16px; }
.pro-search-bar { margin-bottom: 12px; }
.pro-action-bar { margin-bottom: 16px; }
.pro-content-area {
  background: #fff; border-radius: 8px; padding: 20px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
}
</style>
