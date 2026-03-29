<template>
  <div class="templates-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>系统管理</a-breadcrumb-item>
      <a-breadcrumb-item>邮件模板</a-breadcrumb-item>
    </a-breadcrumb>

    <a-row :gutter="16">
      <!-- 左侧模板列表 -->
      <a-col :span="14">
        <a-card class="list-card">
          <template #title>
            <div class="card-title">
              <icon-email />
              <span>邮件模板列表</span>
            </div>
          </template>

          <template #extra>
            <a-button type="primary" @click="handleCreate">
              <template #icon><icon-plus /></template>
              新建模板
            </a-button>
          </template>

          <a-form :model="query" layout="inline" style="margin-bottom: 16px">
            <a-form-item field="name" label="模板名称">
              <a-input v-model="query.name" placeholder="请输入模板名称" allow-clear style="width: 160px" />
            </a-form-item>
            <a-form-item field="type" label="模板类型">
              <a-select v-model="query.type" placeholder="选择类型" allow-clear style="width: 140px">
                <a-option value="system">系统通知</a-option>
                <a-option value="device">设备告警</a-option>
                <a-option value="user">用户相关</a-option>
                <a-option value="ota">OTA升级</a-option>
              </a-select>
            </a-form-item>
            <a-form-item>
              <a-button type="primary" @click="handleSearch">查询</a-button>
            </a-form-item>
          </a-form>

          <a-table
            :columns="columns"
            :data="templates"
            :loading="loading"
            :pagination="pagination"
            @change="handleTableChange"
            row-key="id"
            :row-selection="{ type: 'checkbox', showCheckedAll: true }"
            @selection-change="handleSelectionChange"
          >
            <template #status="{ record }">
              <a-tag :color="record.status === 1 ? 'green' : 'gray'">
                {{ record.status === 1 ? '启用' : '禁用' }}
              </a-tag>
            </template>
            <template #actions="{ record }">
              <a-button type="text" size="small" @click="handleEdit(record)">
                <template #icon><icon-edit /></template>
                编辑
              </a-button>
              <a-popconfirm content="确定删除该模板?" @ok="handleDelete(record.id)">
                <a-button type="text" size="small" status="danger">
                  <template #icon><icon-delete /></template>
                  删除
                </a-button>
              </a-popconfirm>
            </template>
          </a-table>

          <div v-if="selectedIds.length > 0" class="batch-actions">
            <a-space>
              <span>已选择 {{ selectedIds.length }} 项</span>
              <a-button type="text" size="small" status="danger" @click="handleBatchDelete">
                <template #icon><icon-delete /></template>
                批量删除
              </a-button>
            </a-space>
          </div>
        </a-card>
      </a-col>

      <!-- 右侧变量说明 -->
      <a-col :span="10">
        <a-card class="variables-card">
          <template #title>
            <div class="card-title">
              <icon-info-circle />
              <span>变量说明</span>
            </div>
          </template>

          <a-tabs default-active-key="system">
            <a-tab-pane key="system" title="系统变量">
              <a-descriptions :column="1" size="small">
                <a-descriptions-item label="{{system_name}}">
                  系统名称
                </a-descriptions-item>
                <a-descriptions-item label="{{system_url}}">
                  系统访问地址
                </a-descriptions-item>
                <a-descriptions-item label="{{current_time}}">
                  当前时间
                </a-descriptions-item>
                <a-descriptions-item label="{{admin_email}}">
                  管理员邮箱
                </a-descriptions-item>
                <a-descriptions-item label="{{company_name}}">
                  公司名称
                </a-descriptions-item>
              </a-descriptions>
            </a-tab-pane>

            <a-tab-pane key="user" title="用户变量">
              <a-descriptions :column="1" size="small">
                <a-descriptions-item label="{{user_id}}">
                  用户ID
                </a-descriptions-item>
                <a-descriptions-item label="{{user_name}}">
                  用户名称
                </a-descriptions-item>
                <a-descriptions-item label="{{user_email}}">
                  用户邮箱
                </a-descriptions-item>
                <a-descriptions-item label="{{user_phone}}">
                  用户手机号
                </a-descriptions-item>
                <a-descriptions-item label="{{register_time}}">
                  注册时间
                </a-descriptions-item>
              </a-descriptions>
            </a-tab-pane>

            <a-tab-pane key="device" title="设备变量">
              <a-descriptions :column="1" size="small">
                <a-descriptions-item label="{{device_id}}">
                  设备ID
                </a-descriptions-item>
                <a-descriptions-item label="{{device_name}}">
                  设备名称
                </a-descriptions-item>
                <a-descriptions-item label="{{device_model}}">
                  设备型号
                </a-descriptions-item>
                <a-descriptions-item label="{{device_sn}}">
                  设备序列号
                </a-descriptions-item>
                <a-descriptions-item label="{{device_status}}">
                  设备状态
                </a-descriptions-item>
                <a-descriptions-item label="{{last_online_time}}">
                  最后在线时间
                </a-descriptions-item>
              </a-descriptions>
            </a-tab-pane>

            <a-tab-pane key="ota" title="OTA变量">
              <a-descriptions :column="1" size="small">
                <a-descriptions-item label="{{firmware_version}}">
                  固件版本
                </a-descriptions-item>
                <a-descriptions-item label="{{ota_version}}">
                  升级版本
                </a-descriptions-item>
                <a-descriptions-item label="{{download_url}}">
                  下载地址
                </a-descriptions-item>
                <a-descriptions-item label="{{release_notes}}">
                  更新说明
                </a-descriptions-item>
              </a-descriptions>
            </a-tab-pane>
          </a-tabs>

          <a-divider>使用示例</a-divider>
          <a-card size="small" class="example-card">
            <div class="example-content">
              <p><strong>主题:</strong> {{system_name}} - 设备告警通知</p>
              <p><strong>正文:</strong></p>
              <div style="background: #f5f7fa; padding: 12px; border-radius: 4px; margin-top: 8px">
                <p>尊敬的 {{user_name}}：</p>
                <p>您的设备 <strong>{{device_name}}</strong> ({{device_sn}}) 发生告警。</p>
                <p>告警时间：{{current_time}}</p>
                <p>请及时处理。</p>
              </div>
            </div>
          </a-card>
        </a-card>
      </a-col>
    </a-row>

    <!-- 编辑弹窗 -->
    <a-modal
      v-model:visible="modalVisible"
      :title="isEdit ? '编辑邮件模板' : '新建邮件模板'"
      :width="800"
      :mask-closable="false"
      @before-ok="handleSubmit"
      @cancel="handleCancel"
    >
      <a-form :model="form" layout="vertical">
        <a-form-item field="name" label="模板名称" required>
          <a-input v-model="form.name" placeholder="请输入模板名称" />
        </a-form-item>

        <a-form-item field="type" label="模板类型" required>
          <a-select v-model="form.type" placeholder="选择模板类型">
            <a-option value="system">系统通知</a-option>
            <a-option value="device">设备告警</a-option>
            <a-option value="user">用户相关</a-option>
            <a-option value="ota">OTA升级</a-option>
          </a-select>
        </a-form-item>

        <a-form-item field="subject" label="邮件主题" required>
          <a-input v-model="form.subject" placeholder="请输入邮件主题，支持变量替换" />
          <template #extra>
            <span class="form-hint">可用变量: {{systemVariables}}</span>
          </template>
        </a-form-item>

        <a-form-item field="content" label="邮件正文" required>
          <a-textarea
            v-model="form.content"
            placeholder="请输入邮件正文内容，支持HTML格式"
            :auto-size="{ minRows: 8, maxRows: 16 }"
          />
          <template #extra>
            <span class="form-hint">支持HTML标签，可用变量格式: {"{"+变量名+"}"}</span>
          </template>
        </a-form-item>

        <a-form-item field="status" label="状态">
          <a-switch v-model="form.status" />
          <span style="margin-left: 8px">{{ form.status ? '启用' : '禁用' }}</span>
        </a-form-item>
      </a-form>

      <template #footer>
        <a-button @click="handleCancel">取消</a-button>
        <a-button type="primary" :loading="submitting" @click="handleSubmit">提交</a-button>
      </template>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import * as api from '@/api/system'

const loading = ref(false)
const submitting = ref(false)
const templates = ref([])
const selectedIds = ref([])
const modalVisible = ref(false)
const isEdit = ref(false)
const currentId = ref(null)

const query = reactive({
  name: '',
  type: ''
})

const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0
})

const form = reactive({
  name: '',
  type: 'system',
  subject: '',
  content: '',
  status: true
})

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '模板名称', dataIndex: 'name', ellipsis: true },
  { title: '类型', dataIndex: 'type', width: 100 },
  { title: '主题', dataIndex: 'subject', ellipsis: true },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '更新时间', dataIndex: 'updated_at', width: 180 },
  { title: '操作', slotName: 'actions', width: 160, fixed: 'right' }
]

const systemVariables = computed(() => {
  return Object.keys(getSystemVars()).join(', ')
})

function getSystemVars() {
  return {
    system_name: '系统名称',
    system_url: '系统地址',
    current_time: '当前时间',
    admin_email: '管理员邮箱',
    company_name: '公司名称'
  }
}

async function loadTemplates() {
  loading.value = true
  try {
    const res = await api.getEmailTemplates({
      page: pagination.current,
      page_size: pagination.pageSize,
      name: query.name,
      type: query.type
    })
    if (res.code === 0) {
      templates.value = res.data?.list || []
      pagination.total = res.data?.total || 0
    } else {
      // Mock data for demo
      templates.value = getMockTemplates()
      pagination.total = templates.value.length
    }
  } catch (e) {
    console.error('加载模板失败:', e)
    templates.value = getMockTemplates()
    pagination.total = templates.value.length
  } finally {
    loading.value = false
  }
}

function getMockTemplates() {
  return [
    { id: 1, name: '设备告警通知', type: 'device', subject: '{{system_name}} - 设备告警通知', content: '<p>尊敬的 {{user_name}}：</p><p>您的设备发生告警，请及时处理。</p>', status: 1, updated_at: '2026-03-20 10:30:00' },
    { id: 2, name: 'OTA升级提醒', type: 'ota', subject: '{{system_name}} - 新版本可用', content: '<p>尊敬的 {{user_name}}：</p><p>设备 {{device_name}} 有新版本 {{ota_version}} 可升级。</p>', status: 1, updated_at: '2026-03-19 15:20:00' },
    { id: 3, name: '账户激活', type: 'user', subject: '激活您的 {{system_name}} 账户', content: '<p>欢迎加入 {{system_name}}！</p><p>请点击以下链接激活账户。</p>', status: 1, updated_at: '2026-03-18 09:00:00' },
    { id: 4, name: '系统维护通知', type: 'system', subject: '{{system_name}} - 系统维护通知', content: '<p>尊敬的 {{user_name}}：</p><p>系统将于 {{current_time}} 进行维护。</p>', status: 1, updated_at: '2026-03-15 14:00:00' },
    { id: 5, name: '设备离线告警', type: 'device', subject: '{{system_name}} - 设备离线提醒', content: '<p>设备 {{device_name}} 已离线超过指定时间。</p>', status: 0, updated_at: '2026-03-10 11:00:00' }
  ]
}

function handleSearch() {
  pagination.current = 1
  loadTemplates()
}

function handleTableChange(pag) {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  loadTemplates()
}

function handleSelectionChange(ids) {
  selectedIds.value = ids
}

function handleCreate() {
  isEdit.value = false
  currentId.value = null
  Object.assign(form, {
    name: '',
    type: 'system',
    subject: '',
    content: '',
    status: true
  })
  modalVisible.value = true
}

function handleEdit(record) {
  isEdit.value = true
  currentId.value = record.id
  Object.assign(form, {
    name: record.name,
    type: record.type,
    subject: record.subject,
    content: record.content,
    status: record.status === 1
  })
  modalVisible.value = true
}

async function handleSubmit(done) {
  if (!form.name || !form.subject || !form.content) {
    Message.error('请填写必填项')
    done(false)
    return
  }

  submitting.value = true
  try {
    const payload = {
      ...form,
      status: form.status ? 1 : 0
    }

    if (isEdit.value) {
      await api.updateEmailTemplate(currentId.value, payload)
      Message.success('更新成功')
    } else {
      await api.createEmailTemplate(payload)
      Message.success('创建成功')
    }

    modalVisible.value = false
    loadTemplates()
    done(true)
  } catch (e) {
    console.error('保存失败:', e)
    Message.error('保存失败')
    done(false)
  } finally {
    submitting.value = false
  }
}

function handleCancel() {
  modalVisible.value = false
}

async function handleDelete(id) {
  try {
    await api.deleteEmailTemplate(id)
    Message.success('删除成功')
    loadTemplates()
  } catch (e) {
    console.error('删除失败:', e)
    Message.error('删除失败')
  }
}

async function handleBatchDelete() {
  if (selectedIds.value.length === 0) return

  try {
    for (const id of selectedIds.value) {
      await api.deleteEmailTemplate(id)
    }
    Message.success(`已删除 ${selectedIds.value.length} 个模板`)
    selectedIds.value = []
    loadTemplates()
  } catch (e) {
    console.error('批量删除失败:', e)
    Message.error('批量删除失败')
  }
}

onMounted(() => {
  loadTemplates()
})
</script>

<style scoped>
.templates-container {
  padding: 16px 20px;
  min-height: calc(100vh - 64px);
  background: #f5f7fa;
}

.pro-breadcrumb {
  margin-bottom: 16px;
}

.list-card {
  border-radius: 8px;
}

.variables-card {
  border-radius: 8px;
  height: fit-content;
}

.card-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
  font-size: 15px;
}

.batch-actions {
  margin-top: 16px;
  padding: 12px;
  background: #f2f3f5;
  border-radius: 4px;
}

.form-hint {
  font-size: 12px;
  color: #888;
}

.example-card {
  background: #fafafa;
}

.example-content {
  font-size: 13px;
  line-height: 1.8;
}
</style>
