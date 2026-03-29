<template>
  <div class="settings-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>系统管理</a-breadcrumb-item>
      <a-breadcrumb-item>系统参数</a-breadcrumb-item>
    </a-breadcrumb>

    <a-row :gutter="16">
      <!-- 左侧分类导航 -->
      <a-col :span="5">
        <a-card class="category-card">
          <template #title>
            <div class="card-title">
              <icon-settings />
              <span>参数分类</span>
            </div>
          </template>
          <a-menu
            :selected-keys="[activeCategory]"
            @menu-item-click="handleCategoryChange"
          >
            <a-menu-item v-for="cat in categories" :key="cat.key">
              <template #icon><component :is="cat.icon" /></template>
              {{ cat.label }}
            </a-menu-item>
          </a-menu>
        </a-card>
      </a-col>

      <!-- 右侧参数列表 -->
      <a-col :span="19">
        <a-card class="params-card">
          <template #title>
            <div class="card-title">
              <icon-edit />
              <span>{{ currentCategoryLabel }} - 参数配置</span>
            </div>
          </template>

          <template #extra>
            <a-button v-if="hasChanges" type="primary" :loading="saving" @click="handleSave">
              <template #icon><icon-save /></template>
              保存修改
            </a-button>
          </template>

          <a-form :model="form" layout="vertical" :loading="loading">
            <a-table
              :columns="columns"
              :data="params"
              :loading="loading"
              :pagination="false"
              row-key="id"
              size="small"
            >
              <template #type="{ record }">
                <a-tag :color="typeColors[record.type] || 'arcoblue'">
                  {{ typeLabels[record.type] || record.type }}
                </a-tag>
              </template>
              <template #value="{ record }">
                <!-- 文本类型 -->
                <a-input
                  v-if="record.type === 'string'"
                  v-model="record.value"
                  :placeholder="record.default_value"
                  style="width: 280px"
                  @change="markChanged(record)"
                />
                <!-- 数字类型 -->
                <a-input-number
                  v-else-if="record.type === 'number'"
                  v-model="record.value"
                  :min="record.min_value || 0"
                  :max="record.max_value || 999999"
                  style="width: 120px"
                  @change="markChanged(record)"
                />
                <!-- 布尔类型 -->
                <a-switch
                  v-else-if="record.type === 'boolean'"
                  v-model="record.value"
                  @change="markChanged(record)"
                />
                <!-- 枚举类型 -->
                <a-select
                  v-else-if="record.type === 'enum'"
                  v-model="record.value"
                  :options="parseOptions(record.options)"
                  style="width: 200px"
                  @change="markChanged(record)"
                />
                <!-- JSON类型 -->
                <a-textarea
                  v-else-if="record.type === 'json'"
                  v-model="record.value"
                  :placeholder="record.default_value"
                  style="width: 300px; min-height: 60px"
                  @change="markChanged(record)"
                />
                <span v-else>{{ record.value }}</span>
              </template>
              <template #actions="{ record }">
                <a-button
                  v-if="isChanged(record)"
                  type="text"
                  status="warning"
                  size="mini"
                  @click="resetParam(record)"
                >
                  <template #icon><icon-refresh /></template>
                  重置
                </a-button>
              </template>
            </a-table>

            <a-empty v-if="!loading && params.length === 0" style="margin-top: 48px">
              <template #image>
                <icon-empty />
              </template>
              该分类暂无参数配置
            </a-empty>
          </a-form>
        </a-card>
      </a-col>
    </a-row>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import * as api from '@/api/system'

const loading = ref(false)
const saving = ref(false)
const activeCategory = ref('general')
const params = ref([])
const changedParams = ref(new Set())

// 分类定义
const categories = [
  { key: 'general', label: '通用设置', icon: 'icon-settings' },
  { key: 'security', label: '安全设置', icon: 'icon-safe' },
  { key: 'notification', label: '通知设置', icon: 'icon-notification' },
  { key: 'device', label: '设备管理', icon: 'icon-device' },
  { key: 'integration', label: '集成配置', icon: 'icon-connection' },
  { key: 'storage', label: '存储配置', icon: 'icon-storage' },
  { key: 'advanced', label: '高级设置', icon: 'icon-tool' }
]

const currentCategoryLabel = computed(() => {
  const cat = categories.find(c => c.key === activeCategory.value)
  return cat ? cat.label : '通用设置'
})

// 表格列定义
const columns = [
  { title: '参数名称', dataIndex: 'name', width: 200 },
  { title: '参数键', dataIndex: 'key', width: 180 },
  { title: '类型', slotName: 'type', width: 100 },
  { title: '参数值', slotName: 'value' },
  { title: '说明', dataIndex: 'description', ellipsis: true },
  { title: '操作', slotName: 'actions', width: 80 }
]

const typeLabels = {
  string: '字符串',
  number: '数字',
  boolean: '开关',
  enum: '枚举',
  json: 'JSON'
}

const typeColors = {
  string: 'arcoblue',
  number: 'green',
  boolean: 'orange',
  enum: 'purple',
  json: 'cyan'
}

function parseOptions(optionsStr) {
  if (!optionsStr) return []
  try {
    return JSON.parse(optionsStr)
  } catch {
    return optionsStr.split(',').map(v => ({ label: v.trim(), value: v.trim() }))
  }
}

function markChanged(record) {
  changedParams.value.add(record.id)
}

function isChanged(record) {
  return changedParams.value.has(record.id)
}

function resetParam(record) {
  record.value = record.default_value
  changedParams.value.delete(record.id)
}

const hasChanges = computed(() => changedParams.value.size > 0)

async function loadParams() {
  loading.value = true
  changedParams.value.clear()
  try {
    const res = await api.getSystemParamsByCategory(activeCategory.value)
    if (res.code === 0) {
      params.value = (res.data || []).map(p => ({
        ...p,
        value: p.type === 'boolean' ? Boolean(p.value) : p.value
      }))
    } else {
      // Mock data for demo when API not available
      params.value = getMockParams(activeCategory.value)
    }
  } catch (e) {
    console.error('加载参数失败:', e)
    params.value = getMockParams(activeCategory.value)
  } finally {
    loading.value = false
  }
}

function getMockParams(category) {
  const mockData = {
    general: [
      { id: 1, key: 'system_name', name: '系统名称', type: 'string', value: 'MDM管理平台', default_value: 'MDM管理平台', description: '显示在页面标题和首页的名称', options: '' },
      { id: 2, key: 'system_logo', name: '系统Logo', type: 'string', value: '/logo.png', default_value: '/logo.png', description: '系统Logo图片路径', options: '' },
      { id: 3, key: 'timezone', name: '系统时区', type: 'enum', value: 'Asia/Shanghai', default_value: 'UTC', description: '系统使用的时区设置', options: JSON.stringify([{ label: '北京时间 (Asia/Shanghai)', value: 'Asia/Shanghai' }, { label: 'UTC', value: 'UTC' }]) },
      { id: 4, key: 'language', name: '默认语言', type: 'enum', value: 'zh-CN', default_value: 'zh-CN', description: '系统默认显示语言', options: JSON.stringify([{ label: '简体中文', value: 'zh-CN' }, { label: 'English', value: 'en-US' }]) },
      { id: 5, key: 'items_per_page', name: '分页大小', type: 'number', value: 20, default_value: 20, description: '列表每页显示的记录数', min_value: 10, max_value: 100, options: '' },
      { id: 6, key: 'maintenance_mode', name: '维护模式', type: 'boolean', value: false, default_value: false, description: '开启后普通用户无法登录系统', options: '' }
    ],
    security: [
      { id: 10, key: 'password_min_length', name: '密码最小长度', type: 'number', value: 8, default_value: 8, description: '用户密码的最少字符数', min_value: 6, max_value: 32, options: '' },
      { id: 11, key: 'password_require_special', name: '密码特殊字符', type: 'boolean', value: true, default_value: true, description: '密码必须包含特殊字符', options: '' },
      { id: 12, key: 'session_timeout', name: '会话超时(分钟)', type: 'number', value: 30, default_value: 30, description: '无操作后会话超时时间', min_value: 5, max_value: 1440, options: '' },
      { id: 13, key: 'login_max_attempts', name: '登录最大尝试', type: 'number', value: 5, default_value: 5, description: '连续失败后锁定账户', min_value: 3, max_value: 10, options: '' },
      { id: 14, key: 'lockout_duration', name: '锁定时长(分钟)', type: 'number', value: 30, default_value: 30, description: '账户锁定持续时间', min_value: 5, max_value: 1440, options: '' },
      { id: 15, key: 'mfa_enabled', name: '双因素认证', type: 'boolean', value: false, default_value: false, description: '是否启用双因素认证', options: '' }
    ],
    notification: [
      { id: 20, key: 'email_enabled', name: '启用邮件通知', type: 'boolean', value: true, default_value: true, description: '是否发送邮件通知', options: '' },
      { id: 21, key: 'sms_enabled', name: '启用短信通知', type: 'boolean', value: false, default_value: false, description: '是否发送短信通知', options: '' },
      { id: 22, key: 'webhook_enabled', name: '启用Webhook', type: 'boolean', value: false, default_value: false, description: '是否启用Webhook回调', options: '' },
      { id: 23, key: 'notify_admins', name: '通知管理员', type: 'boolean', value: true, default_value: true, description: '重要事件通知所有管理员', options: '' }
    ],
    device: [
      { id: 30, key: 'device_offline_timeout', name: '设备离线超时(秒)', type: 'number', value: 300, default_value: 300, description: '超过此时间无心跳视为离线', min_value: 60, max_value: 3600, options: '' },
      { id: 31, key: 'max_devices_per_user', name: '用户最大设备数', type: 'number', value: 10, default_value: 10, description: '单个用户可绑定的最大设备数', min_value: 1, max_value: 100, options: '' },
      { id: 32, key: 'auto_bind_enabled', name: '自动绑定', type: 'boolean', value: true, default_value: true, description: '新设备自动绑定到用户', options: '' },
      { id: 33, key: 'ota_auto_check', name: '自动检查OTA', type: 'boolean', value: true, default_value: true, description: '定期检查设备固件更新', options: '' }
    ],
    integration: [
      { id: 40, key: 'mqtt_broker_url', name: 'MQTT Broker地址', type: 'string', value: 'tcp://localhost:1883', default_value: 'tcp://localhost:1883', description: 'MQTT服务地址', options: '' },
      { id: 41, key: 'redis_url', name: 'Redis连接地址', type: 'string', value: 'redis://localhost:6379', default_value: 'redis://localhost:6379', description: 'Redis缓存服务器地址', options: '' }
    ],
    storage: [
      { id: 50, key: 'max_upload_size', name: '最大上传大小(MB)', type: 'number', value: 10, default_value: 10, description: '单文件最大上传大小', min_value: 1, max_value: 100, options: '' },
      { id: 51, key: 'allowed_file_types', name: '允许的文件类型', type: 'string', value: 'jpg,png,pdf,zip', default_value: 'jpg,png,pdf,zip', description: '允许上传的文件扩展名', options: '' }
    ],
    advanced: [
      { id: 60, key: 'debug_mode', name: '调试模式', type: 'boolean', value: false, default_value: false, description: '开启后显示详细日志', options: '' },
      { id: 61, key: 'cache_ttl', name: '缓存TTL(秒)', type: 'number', value: 3600, default_value: 3600, description: '数据缓存生存时间', min_value: 60, max_value: 86400, options: '' },
      { id: 62, key: 'log_level', name: '日志级别', type: 'enum', value: 'info', default_value: 'info', description: '系统日志记录级别', options: JSON.stringify([{ label: 'DEBUG', value: 'debug' }, { label: 'INFO', value: 'info' }, { label: 'WARN', value: 'warn' }, { label: 'ERROR', value: 'error' }]) }
    ]
  }
  return mockData[category] || []
}

async function handleSave() {
  saving.value = true
  try {
    const changed = params.value.filter(p => changedParams.value.has(p.id))
    const payload = changed.map(p => ({
      id: p.id,
      key: p.key,
      value: p.value
    }))
    await api.updateSystemSettings({ params: payload })
    Message.success('保存成功')
    changedParams.value.clear()
    await loadParams()
  } catch (e) {
    console.error('保存失败:', e)
    Message.error('保存失败')
  } finally {
    saving.value = false
  }
}

function handleCategoryChange(key) {
  activeCategory.value = key
  loadParams()
}

onMounted(() => {
  loadParams()
})
</script>

<style scoped>
.settings-container {
  padding: 20px 24px;
  min-height: calc(100vh - 64px);
  background: #f5f7fa;
}

.category-card, .params-card {
  border-radius: 8px;
  min-height: 400px;
}

.card-title {
  font-weight: 600;
  font-size: 15px;
}

.category-card :deep(.arco-card-header) {
  padding: 12px 16px;
}

.params-card {
  border-radius: 8px;
}

.card-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
  font-size: 15px;
}

:deep(.arco-table-th) {
  background: #fafafa;
}
</style>
