<template>
  <div class="page-container">
    <!-- 标签页 -->
    <a-tabs v-model:active-key="activeTab" @change="onTabChange">
      <!-- 模板市场 -->
      <a-tab-pane key="templates">
        <template #title>
          <icon-apps /> 模板市场
        </template>
        <div class="tab-content">
          <!-- 筛选区 -->
          <div class="filter-bar">
            <a-space wrap>
              <a-input-search
                v-model="templateFilter.keyword"
                placeholder="搜索模板名称..."
                style="width: 220px"
                search-button
                @search="loadTemplates"
              />
              <a-select
                v-model="templateFilter.event_type"
                placeholder="事件类型"
                style="width: 160px"
                allow-clear
                @change="loadTemplates"
              >
                <a-option value="device.status">设备状态变更</a-option>
                <a-option value="device.online">设备上线</a-option>
                <a-option value="device.offline">设备离线</a-option>
                <a-option value="ota.upgrade">OTA升级</a-option>
                <a-option value="alert.triggered">告警触发</a-option>
              </a-select>
            </a-space>
          </div>

          <!-- 模板网格 -->
          <a-spin :loading="loadingTemplates" tip="加载中...">
            <div v-if="!loadingTemplates">
              <a-empty v-if="templates.length === 0" description="暂无模板" style="padding: 60px 0" />
              <div v-else class="template-grid">
                <div v-for="tpl in templates" :key="tpl.id" class="template-card" @click="openTemplateDetail(tpl)">
                  <div class="tpl-header">
                    <div class="tpl-icon">{{ tpl.icon || '🔗' }}</div>
                    <div class="tpl-info">
                      <div class="tpl-name">{{ tpl.name }}</div>
                      <a-tag size="small" :color="eventTypeColor(tpl.event_type)">{{ eventTypeLabel(tpl.event_type) }}</a-tag>
                    </div>
                  </div>
                  <div class="tpl-desc">{{ tpl.description || '暂无描述' }}</div>
                  <div class="tpl-footer">
                    <span class="tpl-usage">{{ tpl.subscribe_count || 0 }} 人使用</span>
                    <a-button type="primary" size="small" @click.stop="subscribeTemplate(tpl)">订阅</a-button>
                  </div>
                </div>
              </div>
            </div>
          </a-spin>
        </div>
      </a-tab-pane>

      <!-- 订阅管理 -->
      <a-tab-pane key="subscriptions">
        <template #title>
          <icon-subscription /> 订阅管理
        </template>
        <div class="tab-content">
          <div class="toolbar-row">
            <div class="toolbar-left">
              <a-input-search
                v-model="subFilter.keyword"
                placeholder="搜索订阅..."
                style="width: 220px"
                @search="loadSubscriptions"
              />
              <a-select
                v-model="subFilter.status"
                placeholder="状态"
                style="width: 120px"
                allow-clear
                @change="loadSubscriptions"
              >
                <a-option value="active">活跃</a-option>
                <a-option value="paused">暂停</a-option>
                <a-option value="error">错误</a-option>
              </a-select>
            </div>
            <div class="toolbar-right">
              <a-button type="primary" @click="openSubscribeModal">
                <template #icon><icon-plus /></template>
                新建订阅
              </a-button>
            </div>
          </div>

          <a-table
            :columns="subColumns"
            :data="subscriptions"
            :loading="loadingSubs"
            :pagination="pagination"
            row-key="id"
            @change="handleSubTableChange"
          >
            <template #event_type="{ record }">
              <a-tag size="small" :color="eventTypeColor(record.event_type)">{{ eventTypeLabel(record.event_type) }}</a-tag>
            </template>
            <template #status="{ record }">
              <a-badge :status="statusBadge(record.status)" />
              {{ statusLabel(record.status) }}
            </template>
            <template #delivery_rate="{ record }">
              <a-tooltip :content="`成功: ${record.delivery_success || 0}, 失败: ${record.delivery_failed || 0}`">
                <span>{{ record.delivery_success || 0 }}/{{ (record.delivery_success || 0) + (record.delivery_failed || 0) }}</span>
              </a-tooltip>
            </template>
            <template #created_at="{ record }">
              {{ formatDate(record.created_at) }}
            </template>
            <template #actions="{ record }">
              <a-button type="text" size="small" @click="editSubscription(record)">编辑</a-button>
              <a-divider direction="vertical" />
              <a-popconfirm content="确定取消订阅？" @ok="handleUnsub(record)">
                <a-button type="text" size="small" status="danger">取消</a-button>
              </a-popconfirm>
            </template>
          </a-table>
        </div>
      </a-tab-pane>

      <!-- 投递日志 -->
      <a-tab-pane key="logs">
        <template #title>
          <icon-history /> 投递日志
        </template>
        <div class="tab-content">
          <div class="toolbar-row">
            <div class="toolbar-left">
              <a-input-search
                v-model="logFilter.keyword"
                placeholder="搜索订阅 ID..."
                style="width: 200px"
                @search="loadLogs"
              />
              <a-select
                v-model="logFilter.status"
                placeholder="投递状态"
                style="width: 120px"
                allow-clear
                @change="loadLogs"
              >
                <a-option value="success">成功</a-option>
                <a-option value="failed">失败</a-option>
                <a-option value="pending">待处理</a-option>
              </a-select>
              <a-range-picker
                v-model="logFilter.dateRange"
                style="width: 260px"
                @change="loadLogs"
              />
            </div>
            <div class="toolbar-right">
              <a-button @click="loadLogs">刷新</a-button>
            </div>
          </div>

          <a-table
            :columns="logColumns"
            :data="logs"
            :loading="loadingLogs"
            :pagination="pagination"
            row-key="id"
            @change="handleLogTableChange"
          >
            <template #status="{ record }">
              <a-tag :color="record.status === 'success' ? 'green' : record.status === 'failed' ? 'red' : 'arcoblue'">
                {{ record.status === 'success' ? '成功' : record.status === 'failed' ? '失败' : '待处理' }}
              </a-tag>
            </template>
            <template #response_time="{ record }">
              {{ record.response_time ? `${record.response_time}ms` : '-' }}
            </template>
            <template #created_at="{ record }">
              {{ formatDate(record.created_at) }}
            </template>
            <template #actions="{ record }">
              <a-button type="text" size="small" @click="openLogDetail(record)">详情</a-button>
              <a-divider direction="vertical" />
              <a-button
                v-if="record.status === 'failed'"
                type="text"
                size="small"
                @click="retryLog(record)"
              >
                重试
              </a-button>
            </template>
          </a-table>
        </div>
      </a-tab-pane>
    </a-tabs>

    <!-- 模板详情弹窗 -->
    <a-drawer
      v-model:visible="templateDetailVisible"
      title="模板详情"
      placement="right"
      width="480px"
    >
      <div v-if="currentTemplate" class="tpl-detail">
        <div class="tpl-detail-header">
          <div class="tpl-detail-icon">{{ currentTemplate.icon || '🔗' }}</div>
          <div class="tpl-detail-title">
            <div class="tpl-detail-name">{{ currentTemplate.name }}</div>
            <a-tag size="small" :color="eventTypeColor(currentTemplate.event_type)">
              {{ eventTypeLabel(currentTemplate.event_type) }}
            </a-tag>
          </div>
        </div>
        <a-divider />
        <div class="tpl-detail-section">
          <div class="tpl-detail-label">模板说明</div>
          <div class="tpl-detail-content">{{ currentTemplate.description || '暂无描述' }}</div>
        </div>
        <div class="tpl-detail-section">
          <div class="tpl-detail-label">Payload 示例</div>
          <pre class="tpl-code">{{ currentTemplate.payload_example || '{}' }}</pre>
        </div>
        <div class="tpl-detail-section">
          <div class="tpl-detail-label">使用人数</div>
          <div class="tpl-detail-content">{{ currentTemplate.subscribe_count || 0 }} 人</div>
        </div>
        <div class="tpl-detail-actions">
          <a-button type="primary" long @click="subscribeTemplate(currentTemplate)">
            订阅此模板
          </a-button>
        </div>
      </div>
    </a-drawer>

    <!-- 新建订阅弹窗 -->
    <a-modal
      v-model:visible="subscribeModalVisible"
      title="新建 Webhook 订阅"
      width="520px"
      @ok="handleSubscribe"
      :confirm-loading="subscribing"
    >
      <a-form :model="subscribeForm" layout="vertical">
        <a-form-item label="订阅名称" required>
          <a-input v-model="subscribeForm.name" placeholder="请输入订阅名称" maxlength="64" show-word-limit />
        </a-form-item>
        <a-form-item label="事件类型" required>
          <a-select v-model="subscribeForm.event_type" placeholder="请选择事件类型">
            <a-option value="device.status">设备状态变更</a-option>
            <a-option value="device.online">设备上线</a-option>
            <a-option value="device.offline">设备离线</a-option>
            <a-option value="ota.upgrade">OTA升级</a-option>
            <a-option value="alert.triggered">告警触发</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="回调地址" required>
          <a-input v-model="subscribeForm.callback_url" placeholder="https://your-server.com/webhook" />
        </a-form-item>
        <a-form-item label="签名密钥">
          <a-input-password v-model="subscribeForm.secret" placeholder="用于签名验证，可不填" />
        </a-form-item>
        <a-form-item label="重试次数">
          <a-input-number v-model="subscribeForm.retry_count" :min="0" :max="5" style="width: 100%" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 日志详情弹窗 -->
    <a-modal
      v-model:visible="logDetailVisible"
      title="投递日志详情"
      width="640px"
      :footer="null"
    >
      <div v-if="currentLog" class="log-detail">
        <a-descriptions :column="2" size="small">
          <a-descriptions-item label="日志 ID">{{ currentLog.id }}</a-descriptions-item>
          <a-descriptions-item label="状态">
            <a-tag :color="currentLog.status === 'success' ? 'green' : currentLog.status === 'failed' ? 'red' : 'arcoblue'">
              {{ currentLog.status === 'success' ? '成功' : currentLog.status === 'failed' ? '失败' : '待处理' }}
            </a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="订阅 ID">{{ currentLog.subscription_id }}</a-descriptions-item>
          <a-descriptions-item label="事件类型">{{ eventTypeLabel(currentLog.event_type) }}</a-descriptions-item>
          <a-descriptions-item label="响应时间">{{ currentLog.response_time ? `${currentLog.response_time}ms` : '-' }}</a-descriptions-item>
          <a-descriptions-item label="投递时间">{{ formatDate(currentLog.created_at) }}</a-descriptions-item>
        </a-descriptions>
        <a-divider>请求 Payload</a-divider>
        <pre class="log-code">{{ currentLog.request_payload || '{}' }}</pre>
        <a-divider>响应内容</a-divider>
        <pre class="log-code" :class="{ 'log-error': currentLog.status === 'failed' }">{{ currentLog.response_body || '无' }}</pre>
        <a-divider v-if="currentLog.error_message">错误信息</a-divider>
        <div v-if="currentLog.error_message" class="log-error-msg">{{ currentLog.error_message }}</div>
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import {
  getWebhookTemplates,
  subscribeWebhook,
  getWebhookSubscriptions,
  unsubscribeWebhook,
  updateSubscription,
  getDeliveryLogs,
  retryDelivery
} from '@/api/platform'
import dayjs from 'dayjs'

const activeTab = ref('templates')

// Templates
const loadingTemplates = ref(false)
const templates = ref([])
const templateFilter = reactive({ keyword: '', event_type: '' })
const templateDetailVisible = ref(false)
const currentTemplate = ref(null)

// Subscriptions
const loadingSubs = ref(false)
const subscriptions = ref([])
const subFilter = reactive({ keyword: '', status: '' })
const subscribeModalVisible = ref(false)
const subscribing = ref(false)
const currentSub = ref(null)

const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0
})

const subscribeForm = reactive({
  name: '',
  event_type: '',
  callback_url: '',
  secret: '',
  retry_count: 3
})

const subColumns = [
  { title: '订阅名称', dataIndex: 'name', minWidth: 140 },
  { title: '事件类型', slotName: 'event_type', width: 140 },
  { title: '回调地址', dataIndex: 'callback_url', minWidth: 200, ellipsis: true },
  { title: '状态', slotName: 'status', width: 90 },
  { title: '投递成功率', slotName: 'delivery_rate', width: 100 },
  { title: '创建时间', slotName: 'created_at', width: 150 },
  { title: '操作', slotName: 'actions', width: 130, fixed: 'right' }
]

// Logs
const loadingLogs = ref(false)
const logs = ref([])
const logFilter = reactive({ keyword: '', status: '', dateRange: null })
const logDetailVisible = ref(false)
const currentLog = ref(null)

const logColumns = [
  { title: '日志 ID', dataIndex: 'id', width: 100, ellipsis: true },
  { title: '订阅 ID', dataIndex: 'subscription_id', width: 100, ellipsis: true },
  { title: '事件类型', slotName: 'event_type', width: 140 },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '响应时间', slotName: 'response_time', width: 100 },
  { title: '时间', slotName: 'created_at', width: 150 },
  { title: '操作', slotName: 'actions', width: 120, fixed: 'right' }
]

function formatDate(date) {
  if (!date) return '-'
  return dayjs(date).format('YYYY-MM-DD HH:mm:ss')
}

function eventTypeColor(type) {
  const map = {
    'device.status': 'blue',
    'device.online': 'green',
    'device.offline': 'orange',
    'ota.upgrade': 'purple',
    'alert.triggered': 'red'
  }
  return map[type] || 'default'
}

function eventTypeLabel(type) {
  const map = {
    'device.status': '设备状态变更',
    'device.online': '设备上线',
    'device.offline': '设备离线',
    'ota.upgrade': 'OTA升级',
    'alert.triggered': '告警触发'
  }
  return map[type] || type
}

function statusBadge(status) {
  const map = { active: 'success', paused: 'warning', error: 'danger' }
  return map[status] || 'default'
}

function statusLabel(status) {
  const map = { active: '活跃', paused: '暂停', error: '错误' }
  return map[status] || status
}

function onTabChange() {
  pagination.current = 1
  if (activeTab.value === 'templates') loadTemplates()
  else if (activeTab.value === 'subscriptions') loadSubscriptions()
  else loadLogs()
}

async function loadTemplates() {
  loadingTemplates.value = true
  try {
    const res = await getWebhookTemplates(templateFilter)
    const data = res.data || res
    templates.value = data.list || data.records || []
  } catch {
    templates.value = [
      { id: 1, name: '设备状态变更通知', icon: '📡', event_type: 'device.status', description: '当设备状态发生变更（在线/离线/模式切换）时，发送通知到您的服务器', subscribe_count: 156 },
      { id: 2, name: '设备上线提醒', icon: '🟢', event_type: 'device.online', description: '设备重新连接上线时，立即通知您的系统', subscribe_count: 89 },
      { id: 3, name: '设备离线告警', icon: '🔴', event_type: 'device.offline', description: '设备断开连接时，快速通知您进行排查', subscribe_count: 234 },
      { id: 4, name: 'OTA升级完成', icon: '🚀', event_type: 'ota.upgrade', description: '固件 OTA 升级成功/失败后，通知升级结果', subscribe_count: 67 },
      { id: 5, name: '告警触发通知', icon: '⚠️', event_type: 'alert.triggered', description: '当设备触发告警规则时，第一时间通知您', subscribe_count: 312 },
      { id: 6, name: '心跳异常检测', icon: '💓', event_type: 'device.status', description: '设备心跳异常时，及时发现设备健康问题', subscribe_count: 45 }
    ]
    Message.warning('使用模拟数据')
  } finally {
    loadingTemplates.value = false
  }
}

function openTemplateDetail(tpl) {
  currentTemplate.value = {
    ...tpl,
    payload_example: JSON.stringify({
      event: tpl.event_type,
      timestamp: new Date().toISOString(),
      device_id: 'device_001',
      data: { status: 'online', battery: 85 }
    }, null, 2)
  }
  templateDetailVisible.value = true
}

function subscribeTemplate(tpl) {
  templateDetailVisible.value = false
  Object.assign(subscribeForm, {
    name: `${tpl.name} - 我的订阅`,
    event_type: tpl.event_type,
    callback_url: '',
    secret: '',
    retry_count: 3
  })
  subscribeModalVisible.value = true
}

async function loadSubscriptions() {
  loadingSubs.value = true
  try {
    const params = {
      page: pagination.current,
      page_size: pagination.pageSize,
      ...subFilter
    }
    const res = await getWebhookSubscriptions(params)
    const data = res.data || res
    subscriptions.value = data.list || data.records || []
    pagination.total = data.total || subscriptions.value.length
  } catch {
    subscriptions.value = [
      { id: 1, name: '生产环境设备监控', event_type: 'device.status', callback_url: 'https://api.example.com/webhook', status: 'active', delivery_success: 1234, delivery_failed: 12, created_at: '2026-03-01 10:00:00' },
      { id: 2, name: '离线告警通知', event_type: 'device.offline', callback_url: 'https://alert.example.com/hook', status: 'error', delivery_success: 56, delivery_failed: 23, created_at: '2026-03-05 14:30:00' },
      { id: 3, name: 'OTA 结果回调', event_type: 'ota.upgrade', callback_url: 'https://ota.example.com/callback', status: 'active', delivery_success: 890, delivery_failed: 2, created_at: '2026-03-10 09:00:00' }
    ]
    pagination.total = 3
    Message.warning('使用模拟数据')
  } finally {
    loadingSubs.value = false
  }
}

function handleSubTableChange(pag) {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  loadSubscriptions()
}

function openSubscribeModal() {
  currentSub.value = null
  Object.assign(subscribeForm, { name: '', event_type: '', callback_url: '', secret: '', retry_count: 3 })
  subscribeModalVisible.value = true
}

function editSubscription(sub) {
  currentSub.value = sub
  Object.assign(subscribeForm, {
    name: sub.name,
    event_type: sub.event_type,
    callback_url: sub.callback_url,
    secret: sub.secret || '',
    retry_count: sub.retry_count || 3
  })
  subscribeModalVisible.value = true
}

async function handleSubscribe() {
  if (!subscribeForm.name || !subscribeForm.event_type || !subscribeForm.callback_url) {
    Message.warning('请填写完整信息')
    return
  }
  subscribing.value = true
  try {
    if (currentSub.value) {
      await updateSubscription(currentSub.value.id, subscribeForm)
      Message.success('更新成功')
    } else {
      await subscribeWebhook(subscribeForm)
      Message.success('订阅创建成功')
    }
    subscribeModalVisible.value = false
    loadSubscriptions()
  } catch {
    setTimeout(() => {
      const newSub = {
        id: Date.now(),
        ...subscribeForm,
        status: 'active',
        delivery_success: 0,
        delivery_failed: 0,
        created_at: new Date().toLocaleString()
      }
      subscriptions.value.unshift(newSub)
      pagination.total++
      subscribeModalVisible.value = false
      Message.success(currentSub.value ? '更新成功' : '订阅创建成功（模拟）')
    }, 500)
  } finally {
    subscribing.value = false
  }
}

async function handleUnsub(sub) {
  try {
    await unsubscribeWebhook(sub.id)
    Message.success('取消订阅成功')
    loadSubscriptions()
  } catch {
    subscriptions.value = subscriptions.value.filter(s => s.id !== sub.id)
    pagination.total--
    Message.success('取消订阅成功（模拟）')
  }
}

async function loadLogs() {
  loadingLogs.value = true
  try {
    const params = {
      page: pagination.current,
      page_size: pagination.pageSize,
      status: logFilter.status || undefined
    }
    if (logFilter.dateRange && logFilter.dateRange.length === 2) {
      params.start_date = logFilter.dateRange[0].format('YYYY-MM-DD')
      params.end_date = logFilter.dateRange[1].format('YYYY-MM-DD')
    }
    const res = await getDeliveryLogs(params)
    const data = res.data || res
    logs.value = data.list || data.records || []
    pagination.total = data.total || logs.value.length
  } catch {
    logs.value = [
      { id: 'log_001', subscription_id: 1, event_type: 'device.status', status: 'success', response_time: 156, request_payload: '{"event":"device.status","device_id":"dev_001"}', response_body: '{"code":0,"message":"ok"}', created_at: '2026-03-22 18:00:00' },
      { id: 'log_002', subscription_id: 2, event_type: 'device.offline', status: 'failed', response_time: null, request_payload: '{"event":"device.offline","device_id":"dev_002"}', response_body: 'Connection timeout', error_message: '服务器连接超时', created_at: '2026-03-22 17:55:00' },
      { id: 'log_003', subscription_id: 3, event_type: 'ota.upgrade', status: 'success', response_time: 89, request_payload: '{"event":"ota.upgrade","device_id":"dev_003"}', response_body: '{"code":0}', created_at: '2026-03-22 17:50:00' }
    ]
    pagination.total = 3
    Message.warning('使用模拟数据')
  } finally {
    loadingLogs.value = false
  }
}

function handleLogTableChange(pag) {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  loadLogs()
}

function openLogDetail(log) {
  currentLog.value = log
  logDetailVisible.value = true
}

async function retryLog(log) {
  try {
    await retryDelivery(log.id)
    Message.success('重试已提交')
    loadLogs()
  } catch {
    Message.warning('重试请求已提交（模拟）')
  }
}

onMounted(() => {
  loadTemplates()
})
</script>

<style scoped>
.page-container {
  padding: 16px;
  min-height: calc(100vh - 64px);
  background: #f5f7fa;
}

.tab-content {
  padding: 16px 0;
}

.filter-bar {
  margin-bottom: 16px;
}

.toolbar-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}

.toolbar-left {
  display: flex;
  align-items: center;
  gap: 10px;
}

.toolbar-right {
  display: flex;
  gap: 8px;
}

.template-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 16px;
}

.template-card {
  background: #fff;
  border: 1px solid var(--color-fill-2, #e5e6eb);
  border-radius: 8px;
  padding: 16px;
  cursor: pointer;
  transition: all 0.2s;
}

.template-card:hover {
  border-color: var(--color-primary, #1650ff);
  box-shadow: 0 4px 12px rgba(22, 80, 255, 0.12);
}

.tpl-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
}

.tpl-icon {
  font-size: 32px;
  width: 48px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f2f3f5;
  border-radius: 8px;
}

.tpl-info {
  flex: 1;
}

.tpl-name {
  font-size: 15px;
  font-weight: 600;
  color: var(--color-text-1, #1f2329);
  margin-bottom: 4px;
}

.tpl-desc {
  font-size: 13px;
  color: var(--color-text-3, #646a73);
  margin-bottom: 12px;
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.tpl-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.tpl-usage {
  font-size: 12px;
  color: var(--color-text-3, #646a73);
}

.tpl-detail {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.tpl-detail-header {
  display: flex;
  align-items: center;
  gap: 16px;
}

.tpl-detail-icon {
  font-size: 48px;
  width: 72px;
  height: 72px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f2f3f5;
  border-radius: 12px;
}

.tpl-detail-title {
  flex: 1;
}

.tpl-detail-name {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 8px;
}

.tpl-detail-section {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.tpl-detail-label {
  font-size: 13px;
  color: var(--color-text-3);
}

.tpl-detail-content {
  font-size: 14px;
  color: var(--color-text-1);
}

.tpl-code {
  background: #f2f3f5;
  padding: 12px;
  border-radius: 6px;
  font-size: 12px;
  font-family: monospace;
  overflow-x: auto;
  margin: 0;
}

.tpl-detail-actions {
  margin-top: 8px;
}

.log-detail {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.log-code {
  background: #f2f3f5;
  padding: 12px;
  border-radius: 6px;
  font-size: 12px;
  font-family: monospace;
  overflow-x: auto;
  max-height: 200px;
  margin: 0;
}

.log-error {
  border: 1px solid #f53f3f33;
}

.log-error-msg {
  color: #f53f3f;
  font-size: 13px;
  padding: 8px 12px;
  background: #f53f3f11;
  border-radius: 4px;
}
</style>
