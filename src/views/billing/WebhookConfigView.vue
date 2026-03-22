<template>
  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>订阅管理</a-breadcrumb-item>
      <a-breadcrumb-item>Webhook配置</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 页面标题 -->
    <div class="pro-page-header">
      <h2 class="pro-page-title">Webhook 配置</h2>
    </div>

    <!-- 操作栏 -->
    <div class="pro-action-bar">
      <a-space>
        <a-button type="primary" @click="handleCreate">
          <template #icon><icon-plus /></template>
          创建Webhook
        </a-button>
        <a-button @click="loadWebhooks">
          <template #icon><icon-refresh /></template>
          刷新
        </a-button>
      </a-space>
    </div>

    <!-- Webhook 列表 -->
    <a-table
      :columns="columns"
      :data="webhooks"
      :loading="loading"
      :pagination="{ pageSize: 10 }"
      v-loading="loading"
    >
      <template #is_active="{ record }">
        <a-tag :color="record.is_active ? 'green' : 'gray'">
          {{ record.is_active ? '启用' : '停用' }}
        </a-tag>
      </template>
      <template #events="{ record }">
        <a-tag v-for="e in (record.events || [])" :key="e" size="small">{{ e }}</a-tag>
      </template>
      <template #operations="{ record }">
        <a-space>
          <a-button size="small" @click="handleTest(record)">测试</a-button>
          <a-button size="small" @click="handleEdit(record)">编辑</a-button>
          <a-button size="small" status="danger" @click="handleDelete(record)">删除</a-button>
        </a-space>
      </template>
    </a-table>

    <!-- 创建/编辑弹窗 -->
    <a-modal
      v-model:visible="modalVisible"
      :title="isEdit ? '编辑Webhook' : '创建Webhook'"
      @before-ok="confirmSave"
      width="560px"
    >
      <a-form :model="form" layout="vertical">
        <a-form-item label="Webhook名称" required>
          <a-input v-model="form.webhook_name" placeholder="请输入名称" />
        </a-form-item>
        <a-form-item label="回调地址" required>
          <a-input v-model="form.endpoint_url" placeholder="https://example.com/webhook" />
        </a-form-item>
        <a-form-item label="事件类型" required>
          <a-checkbox-group v-model="form.events">
            <a-checkbox value="subscription.created">订阅创建</a-checkbox>
            <a-checkbox value="subscription.renewed">订阅续费</a-checkbox>
            <a-checkbox value="subscription.cancelled">订阅取消</a-checkbox>
            <a-checkbox value="subscription.upgraded">订阅升级</a-checkbox>
            <a-checkbox value="usage.exceeded">用量超限</a-checkbox>
            <a-checkbox value="invoice.created">账单创建</a-checkbox>
          </a-checkbox-group>
        </a-form-item>
        <a-form-item label="密钥">
          <a-input-password v-model="form.secret_key" placeholder="用于签名验证（可选）" />
          <div class="form-hint">用于生成 HMAC-SHA256 签名的密钥</div>
        </a-form-item>
        <a-form-item label="启用状态">
          <a-switch v-model="form.is_active" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 事件历史 -->
    <a-card class="log-card" style="margin-top: 16px;" v-if="selectedWebhook">
      <template #title>
        <a-space>
          事件历史
          <a-tag>{{ selectedWebhook.webhook_name }}</a-tag>
        </a-space>
      </template>
      <a-table
        :columns="logColumns"
        :data="webhookLogs"
        :loading="logsLoading"
        :pagination="{ pageSize: 5 }"
      />
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import { billingApi } from '@/api/billing'

const loading = ref(false)
const logsLoading = ref(false)
const webhooks = ref<any[]>([])
const webhookLogs = ref<any[]>([])
const selectedWebhook = ref<any>(null)
const modalVisible = ref(false)
const isEdit = ref(false)
const editingId = ref<number | null>(null)

const form = ref({
  webhook_name: '',
  endpoint_url: '',
  events: [] as string[],
  secret_key: '',
  is_active: true
})

const columns = [
  { title: '名称', dataIndex: 'webhook_name' },
  { title: '回调地址', dataIndex: 'endpoint_url' },
  { title: '事件', slotName: 'events' },
  { title: '状态', slotName: 'is_active' },
  { title: '操作', slotName: 'operations', width: 220 }
]

const logColumns = [
  { title: '事件类型', dataIndex: 'event_type' },
  { title: '状态', dataIndex: 'status' },
  { title: 'HTTP状态', dataIndex: 'response_code' },
  { title: '时间', dataIndex: 'sent_at' }
]

const loadWebhooks = async () => {
  loading.value = true
  try {
    const res = await billingApi.getWebhooks()
    if (res.code === 0 || res.code === 200) {
      webhooks.value = res.data || []
    } else {
      webhooks.value = [
        {
          id: 1,
          webhook_name: '订阅通知',
          endpoint_url: 'https://example.com/webhook',
          events: ['subscription.created', 'subscription.renewed'],
          is_active: true
        }
      ]
    }
  } catch (e) {
    webhooks.value = []
  } finally {
    loading.value = false
  }
}

const loadWebhookLogs = async (id: number) => {
  logsLoading.value = true
  try {
    const res = await billingApi.getWebhookLogs(id)
    if (res.code === 0 || res.code === 200) {
      webhookLogs.value = res.data || []
    }
  } catch (e) {
    webhookLogs.value = []
  } finally {
    logsLoading.value = false
  }
}

const handleCreate = () => {
  isEdit.value = false
  editingId.value = null
  form.value = { webhook_name: '', endpoint_url: '', events: [], secret_key: '', is_active: true }
  modalVisible.value = true
}

const handleEdit = (record: any) => {
  isEdit.value = true
  editingId.value = record.id
  form.value = {
    webhook_name: record.webhook_name,
    endpoint_url: record.endpoint_url,
    events: record.events || [],
    secret_key: record.secret_key || '',
    is_active: record.is_active
  }
  modalVisible.value = true
}

const confirmSave = async (done: (val: boolean) => void) => {
  if (!form.value.webhook_name || !form.value.endpoint_url || form.value.events.length === 0) {
    Message.error('请填写完整信息')
    done(false)
    return
  }
  try {
    let res: any
    if (isEdit.value && editingId.value) {
      res = await billingApi.updateWebhook(editingId.value, form.value)
    } else {
      res = await billingApi.createWebhook(form.value)
    }
    if (res.code === 0 || res.code === 200) {
      Message.success(isEdit.value ? '更新成功' : '创建成功')
      modalVisible.value = false
      loadWebhooks()
      done(true)
    } else {
      Message.error('操作失败')
      done(false)
    }
  } catch (e) {
    Message.error('操作失败')
    done(false)
  }
}

const handleTest = async (record: any) => {
  selectedWebhook.value = record
  try {
    const res = await billingApi.testWebhook(record.id)
    if (res.code === 0 || res.code === 200) {
      Message.success('测试请求已发送')
    } else {
      Message.error('测试失败')
    }
    loadWebhookLogs(record.id)
  } catch (e) {
    Message.error('测试失败')
  }
}

const handleDelete = async (record: any) => {
  try {
    const res = await billingApi.deleteWebhook(record.id)
    if (res.code === 0 || res.code === 200) {
      Message.success('删除成功')
      loadWebhooks()
    } else {
      Message.error('删除失败')
    }
  } catch (e) {
    Message.error('删除失败')
  }
}

onMounted(() => {
  loadWebhooks()
})
</script>

<style scoped>
.form-hint {
  font-size: 12px;
  color: var(--color-text-3, #86909c);
  margin-top: 4px;
}
</style>
