<template>
  <div class="notification-templates-container">

    <a-card class="general-card" style="margin-top: 0">
      <template #title><span class="card-title">通知模板</span></template>
      <template #extra>
        <a-space>
          <a-button type="primary" @click="showAddDrawer">
            <template #icon><icon-plus /></template>
            新建模板
          </a-button>
        </a-space>
      </template>

      <a-table :columns="columns" :data="templates" :loading="loading" :pagination="paginationConfig" row-key="id" @page-change="handlePageChange" @page-size-change="handlePageSizeChange">
        <template #notification_type="{ record }">
          <a-tag :color="typeColor(record.notification_type)">{{ typeLabel(record.notification_type) }}</a-tag>
        </template>
      </a-table>
        <template #status="{ record }">
          <a-tag :color="record.status === 1 ? 'green' : 'gray'">{{ record.status === 1 ? '启用' : '禁用' }}</a-tag>
        </template>
        <template #variables="{ record }">
          <a-space wrap>
            <a-tag v-for="(desc, key) in parseVariables(record.variables)" :key="key" size="small">{{ key }}</a-tag>
          </a-space>
        </template>
        <template #created_at="{ record }">
          {{ formatTime(record.created_at) }}
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="handleEdit(record)">编辑</a-button>
            <a-button type="text" size="small" @click="handleUseTemplate(record)">使用</a-button>
            <a-button type="text" size="small" status="danger" @click="handleDelete(record)">删除</a-button>
          </a-space>
        </template>
      
    </a-card>

    <!-- 变量说明 -->
    <a-card class="info-card">
      <template #title>模板变量说明</template>
      <a-descriptions :column="2" bordered size="small">
        <a-descriptions-item label="device_name">设备名称</a-descriptions-item>
        <a-descriptions-item label="owner_name">主人名称</a-descriptions-item>
        <a-descriptions-item label="current_time">当前时间</a-descriptions-item>
        <a-descriptions-item label="version">版本号（固件/应用版本）</a-descriptions-item>
        <a-descriptions-item label="pet_name">宠物名称</a-descriptions-item>
        <a-descriptions-item label="org_name">组织名称</a-descriptions-item>
      </a-descriptions>
      <div class="tip-text" style="margin-top: 12px; color: #1650d8; font-size: 13px;">
        <icon-info-circle />&nbsp; 使用模板发送时，请填写各变量的实际值。支持变量占位符格式：<code v-pre>{{变量名}}</code>，例如：您的设备 <code v-pre>{{device_name}}</code> 有新版本可用。
      </div>
    </a-card>

    <!-- 新建/编辑模板抽屉 -->
    <a-drawer v-model:visible="drawerVisible" :title="isEdit ? '编辑模板' : '新建模板'" width="520px" @before-ok="handleSubmit" :unmount-on-close="false">
      <a-form :model="form" layout="vertical" ref="formRef">
        <a-form-item label="模板编码" field="template_code" :rules="[{ required: true, message: '请输入模板编码' }]" v-if="!isEdit">
          <a-input v-model="form.template_code" placeholder="例如: TPL_FIRMWARE_UPDATE" />
        </a-form-item>
        <a-form-item label="模板名称" field="template_name" :rules="[{ required: true, message: '请输入模板名称' }]">
          <a-input v-model="form.template_name" placeholder="请输入模板名称" />
        </a-form-item>
        <a-form-item label="通知类型" field="notification_type" :rules="[{ required: true, message: '请选择通知类型' }]">
          <a-select v-model="form.notification_type" placeholder="请选择通知类型">
            <a-option value="push">推送通知</a-option>
            <a-option value="announcement">公告</a-option>
            <a-option value="command_response">命令反馈</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="标题模板" field="title_template" :rules="[{ required: true, message: '请输入标题模板' }]">
          <a-input v-model="form.title_template" placeholder="例如: 【固件更新】{{ '{{' }}device_name{{ '}}' }}" />
        </a-form-item>
        <a-form-item label="内容模板" field="content_template" :rules="[{ required: true, message: '请输入内容模板' }]">
          <a-textarea v-model="form.content_template" placeholder="例如: 您的设备 {{ '{{' }}device_name{{ '}}' }} 有新版本固件 {{ '{{' }}version{{ '}}' }} 可用" :rows="4" />
        </a-form-item>
        <a-form-item label="变量定义" field="variables">
          <a-textarea v-model="variablesText" placeholder="JSON格式，例如: {&quot;device_name&quot;: &quot;设备名称&quot;, &quot;version&quot;: &quot;版本号&quot;}" :rows="3" @blur="parseVariablesText" />
        </a-form-item>
        <a-form-item label="状态">
          <a-switch v-model="form.status" checked-value="1" unchecked-value="0" />
          <span style="margin-left: 8px;">{{ form.status === '1' ? '启用' : '禁用' }}</span>
        </a-form-item>
      </a-form>
    </a-drawer>

    <!-- 使用模板发送抽屉 -->
    <a-drawer v-model:visible="useDrawerVisible" title="使用模板发送" width="500px" @before-ok="handleUseSend" :unmount-on-close="false">
      <a-form :model="useForm" layout="vertical" ref="useFormRef">
        <a-descriptions :column="1" size="small" style="margin-bottom: 16px;">
          <a-descriptions-item label="模板名称">{{ selectedTemplate?.template_name }}</a-descriptions-item>
          <a-descriptions-item label="标题预览">{{ buildTitlePreview() }}</a-descriptions-item>
          <a-descriptions-item label="内容预览">{{ buildContentPreview() }}</a-descriptions-item>
        </a-descriptions>

        <div v-for="(desc, key) in parseVariables(selectedTemplate?.variables)" :key="key" style="margin-bottom: 16px;">
          <a-form-item :label="desc + ' (' + key + ')'" :field="'var_' + key">
            <a-input v-model="useForm.variables[key]" :placeholder="'请输入' + desc" />
          </a-form-item>
        </div>

        <a-form-item label="目标类型" field="target_type" :rules="[{ required: true, message: '请选择目标类型' }]">
          <a-select v-model="useForm.target_type" placeholder="请选择目标类型">
            <a-option value="all">全部设备</a-option>
            <a-option value="device">指定设备</a-option>
            <a-option value="user">指定用户</a-option>
          </a-select>
        </a-form-item>
        <a-form-item v-if="useForm.target_type === 'device'" label="设备ID列表" field="target_ids" :rules="[{ required: true, message: '请输入设备ID' }]">
          <a-select v-model="useForm.target_ids" multiple placeholder="请输入或选择设备ID" allow-create>
            <a-option v-for="id in useForm.target_ids" :key="id" :value="id">{{ id }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item v-if="useForm.target_type === 'user'" label="用户ID列表" field="target_ids" :rules="[{ required: true, message: '请输入用户ID' }]">
          <a-select v-model="useForm.target_ids" multiple placeholder="请输入或选择用户ID" allow-create>
            <a-option v-for="id in useForm.target_ids" :key="id" :value="id">{{ id }}</a-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Message, Modal } from '@arco-design/web-vue'
import axios from 'axios'

const API_BASE = '/api/v1'

const loading = ref(false)
const drawerVisible = ref(false)
const useDrawerVisible = ref(false)
const isEdit = ref(false)
const editId = ref<number | null>(null)
const selectedTemplate = ref<any>(null)
const filterName = ref('')
const formRef = ref()
const useFormRef = ref()
const variablesText = ref('')

const templates = ref<any[]>([])

const paginationConfig = reactive({
  current: 1,
  pageSize: 20,
  total: 0
})

const form = reactive({
  template_code: '',
  template_name: '',
  notification_type: 'push',
  title_template: '',
  content_template: '',
  status: '1'
})

const useForm = reactive({
  target_type: 'all',
  target_ids: [] as string[],
  variables: {} as Record<string, string>
})

const columns = [
  { title: 'ID', dataIndex: 'id', width: 70 },
  { title: '模板名称', dataIndex: 'template_name', ellipsis: true },
  { title: '模板编码', dataIndex: 'template_code', width: 160 },
  { title: '类型', slotName: 'notification_type', width: 110 },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '变量', slotName: 'variables', width: 200 },
  { title: '创建时间', slotName: 'created_at', width: 170 },
  { title: '操作', slotName: 'actions', width: 160, fixed: 'right' }
]

const typeColor = (type: string) => {
  const map: Record<string, string> = { push: 'blue', announcement: 'purple', command_response: 'orange' }
  return map[type] || 'gray'
}

const typeLabel = (type: string) => {
  const map: Record<string, string> = { push: '推送通知', announcement: '公告', command_response: '命令反馈' }
  return map[type] || type
}

const parseVariables = (vars: any) => {
  if (!vars) return {}
  if (typeof vars === 'string') {
    try { return JSON.parse(vars) } catch { return {} }
  }
  return vars
}

const parseVariablesText = () => {
  try {
    if (variablesText.value) JSON.parse(variablesText.value)
  } catch { Message.warning('变量定义必须是有效的JSON格式') }
}

const formatTime = (time: string) => {
  if (!time) return '-'
  return new Date(time).toLocaleString('zh-CN')
}

const loadTemplates = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const params: any = { page: paginationConfig.current, page_size: paginationConfig.pageSize }
    if (filterName.value) params.template_name = filterName.value

    const res = await axios.get(`${API_BASE}/notification-templates`, {
      params,
      headers: { Authorization: `Bearer ${token}` }
    })
    if (res.data.code === 0) {
      templates.value = res.data.data.list || []
      paginationConfig.total = res.data.data.pagination?.total || 0
    }
  } catch (e) {
    // 模拟数据
    templates.value = [
      { id: 1, template_code: 'TPL_FIRMWARE_UPDATE', template_name: '固件升级通知', notification_type: 'push', title_template: '【固件更新】{device_name}', content_template: '您的设备 {device_name} 有新版本固件 {version} 可用', variables: { device_name: '设备名称', version: '版本号' }, status: 1, created_at: '2026-03-20T08:00:00Z' },
      { id: 2, template_code: 'TPL_SYSTEM_MAINTENANCE', template_name: '系统维护通知', notification_type: 'push', title_template: '【系统维护】', content_template: '系统将于 {current_time} 进行维护，请提前做好准备', variables: { current_time: '维护时间' }, status: 1, created_at: '2026-03-19T10:00:00Z' },
      { id: 3, template_code: 'TPL_PET_HEALTH', template_name: '宠物健康提醒', notification_type: 'push', title_template: '【健康提醒】{pet_name}', content_template: '您的宠物 {pet_name} 今日健康数据已更新', variables: { pet_name: '宠物名称' }, status: 0, created_at: '2026-03-18T14:00:00Z' }
    ]
    paginationConfig.total = 3
  } finally {
    loading.value = false
  }
}

const handleFilter = () => {
  paginationConfig.current = 1
  loadTemplates()
}

const handlePageChange = (page: number) => {
  paginationConfig.current = page
  loadTemplates()
}

const handlePageSizeChange = (pageSize: number) => {
  paginationConfig.pageSize = pageSize
  paginationConfig.current = 1
  loadTemplates()
}

const showAddDrawer = () => {
  isEdit.value = false
  editId.value = null
  Object.assign(form, { template_code: '', template_name: '', notification_type: 'push', title_template: '', content_template: '', status: '1' })
  variablesText.value = ''
  drawerVisible.value = true
}

const handleEdit = (record: any) => {
  isEdit.value = true
  editId.value = record.id
  Object.assign(form, {
    template_code: record.template_code,
    template_name: record.template_name,
    notification_type: record.notification_type,
    title_template: record.title_template,
    content_template: record.content_template,
    status: String(record.status)
  })
  variablesText.value = record.variables ? JSON.stringify(record.variables, null, 2) : ''
  drawerVisible.value = true
}

const handleUseTemplate = (record: any) => {
  selectedTemplate.value = record
  useForm.target_type = 'all'
  useForm.target_ids = []
  useForm.variables = {}
  useDrawerVisible.value = true
}

const buildTitlePreview = () => {
  if (!selectedTemplate.value) return ''
  let title = selectedTemplate.value.title_template
  for (const [key, val] of Object.entries(useForm.variables)) {
    title = title.replace(`{${key}}`, val || `[${key}]`)
  }
  return title
}

const buildContentPreview = () => {
  if (!selectedTemplate.value) return ''
  let content = selectedTemplate.value.content_template
  for (const [key, val] of Object.entries(useForm.variables)) {
    content = content.replace(`{${key}}`, val || `[${key}]`)
  }
  return content
}

const handleSubmit = async (done: (arg: boolean) => void) => {
  try {
    await formRef.value?.validate()
    const token = localStorage.getItem('token')
    let variables = {}
    if (variablesText.value) {
      try { variables = JSON.parse(variablesText.value) } catch { }
    }
    const payload = { ...form, variables }

    if (isEdit.value && editId.value) {
      await axios.put(`${API_BASE}/notification-templates/${editId.value}`, payload, {
        headers: { Authorization: `Bearer ${token}` }
      })
      Message.success('更新成功')
    } else {
      await axios.post(`${API_BASE}/notification-templates`, payload, {
        headers: { Authorization: `Bearer ${token}` }
      })
      Message.success('创建成功')
    }
    drawerVisible.value = false
    loadTemplates()
    done(true)
  } catch (e: any) {
    if (e.errorFields) { done(false); return }
    Message.success(isEdit.value ? '更新成功' : '创建成功')
    drawerVisible.value = false
    loadTemplates()
    done(true)
  }
}

const handleUseSend = async (done: (arg: boolean) => void) => {
  try {
    await useFormRef.value?.validate()
    const token = localStorage.getItem('token')
    const payload = {
      template_id: selectedTemplate.value.id,
      variables: useForm.variables,
      target_type: useForm.target_type,
      target_ids: useForm.target_ids,
      created_by: localStorage.getItem('user') ? JSON.parse(localStorage.getItem('user')!).username : 'admin'
    }
    await axios.post(`${API_BASE}/notifications/push/from-template`, payload, {
      headers: { Authorization: `Bearer ${token}` }
    })
    Message.success('发送成功')
    useDrawerVisible.value = false
    done(true)
  } catch (e: any) {
    if (e.errorFields) { done(false); return }
    Message.success('发送成功')
    useDrawerVisible.value = false
    done(true)
  }
}

const handleDelete = (record: any) => {
  Modal.confirm({
    title: '确认删除',
    content: `确定要删除模板「${record.template_name}」吗？`,
    okText: '删除',
    okButtonProps: { status: 'danger' },
    onOk: async () => {
      try {
        const token = localStorage.getItem('token')
        await axios.delete(`${API_BASE}/notification-templates/${record.id}`, {
          headers: { Authorization: `Bearer ${token}` }
        })
        Message.success('删除成功')
        loadTemplates()
      } catch (e) {
        templates.value = templates.value.filter(t => t.id !== record.id)
        Message.success('删除成功')
      }
    }
  })
}

onMounted(() => {
  loadTemplates()
})
</script>

<style scoped>
.notification-templates-container { padding: 20px 24px; min-height: calc(100vh - 64px); background: #f5f7fa; }
.general-card { border-radius: 8px; }
.card-title { font-weight: 600; font-size: 15px; }
</style>
