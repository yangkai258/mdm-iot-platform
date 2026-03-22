<template>
  <div class="webhook-config-form">
    <a-form :model="form" layout="vertical" :rules="formRules" ref="formRef">
      <!-- 基本信息 -->
      <a-divider>基本信息</a-divider>
      <a-form-item label="渠道名称" field="channel_name">
        <a-input v-model="form.channel_name" placeholder="如：钉钉机器人通知" />
      </a-form-item>
      <a-form-item label="启用该渠道">
        <a-switch v-model="form.enabled" />
        <span class="form-hint">关闭后该渠道将不会发送任何通知</span>
      </a-form-item>

      <!-- Webhook 配置 -->
      <a-divider>Webhook 配置</a-divider>
      <a-form-item label="Webhook URL" field="webhook_url" required>
        <a-input v-model="form.config.webhook_url" placeholder="https://oapi.dingtalk.com/robot/send?access_token=xxx" />
        <span class="form-hint">支持钉钉、企业微信、飞书等 Webhook 地址</span>
      </a-form-item>
      <a-form-item label="请求方法">
        <a-radio-group v-model="form.config.webhook_method">
          <a-radio value="POST">POST</a-radio>
          <a-radio value="PUT">PUT</a-radio>
        </a-radio-group>
      </a-form-item>

      <!-- 安全设置 -->
      <a-divider>安全设置</a-divider>
      <a-form-item label="签名密钥（可选）">
        <a-input-password
          v-model="form.config.webhook_secret"
          placeholder="用于签名验证的密钥，如钉钉机器人 secret"
        />
        <span class="form-hint">部分 Webhook 需要提供密钥进行签名验证</span>
      </a-form-item>

      <!-- Headers 配置 -->
      <a-divider>自定义 Headers</a-divider>
      <div class="headers-list">
        <div class="header-item" v-for="(header, index) in form.config.webhook_headers" :key="index">
          <a-input v-model="header.key" placeholder="Header 名称" style="width: 200px" />
          <a-input v-model="header.value" placeholder="Header 值" style="width: 200px" />
          <a-button type="text" status="danger" @click="removeHeader(index)">
            <icon-delete />
          </a-button>
        </div>
        <a-button type="dashed" @click="addHeader">
          <template #icon><icon-plus /></template>
          「添加 Header」
        </a-button>
      </div>

      <!-- 发送设置 -->
      <a-divider>发送设置</a-divider>
      <a-form-item label="超时时间（秒）">
        <a-input-number v-model="form.config.timeout_seconds" :min="1" :max="60" :default-value="30" style="width: 200px" />
        <span class="form-hint">请求超时时间，建议 5-30 秒</span>
      </a-form-item>
      <a-form-item label="重试次数">
        <a-input-number v-model="form.config.retry_count" :min="0" :max="5" :default-value="3" style="width: 200px" />
        <span class="form-hint">发送失败时的最大重试次数</span>
      </a-form-item>

      <!-- 测试配置 -->
      <a-divider>测试配置</a-divider>
      <a-form-item label="测试数据">
        <a-textarea
          v-model="testPayload"
          placeholder='测试 payload，JSON 格式，如：{"msg": "测试消息"}'
          :rows="4"
        />
        <span class="form-hint">不填则使用默认测试数据</span>
      </a-form-item>
    </a-form>

    <!-- 操作按钮 -->
    <div class="form-actions">
      <a-space>
        <a-button @click="$emit('cancel')">「取消」</a-button>
        <a-button @click="handleTest" :loading="testing" :disabled="!form.config.webhook_url">
          「发送测试」
        </a-button>
        <a-button type="primary" @click="handleSubmit" :loading="saving">
          「保存配置」
        </a-button>
      </a-space>
    </div>

    <!-- 测试结果 -->
    <a-modal
      v-model:visible="testResultVisible"
      title="测试结果"
      :footer="null"
      :width="400"
    >
      <a-result
        v-if="testResult"
        :status="testResult.success ? 'success' : 'error'"
        :title="testResult.success ? '测试成功' : '测试失败'"
        :sub-title="testResult.message || testResult.error || ''"
      />
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import { useNotificationChannels } from '@/composables/useNotification'

const props = defineProps({
  channel: {
    type: Object,
    default: null
  }
})

const emit = defineEmits(['success', 'cancel'])

const formRef = ref(null)
const saving = ref(false)
const testing = ref(false)
const testPayload = ref('')
const testResultVisible = ref(false)
const testResult = ref(null)

const { createChannel, updateChannel, testChannel } = useNotificationChannels()

const form = reactive({
  channel_name: '',
  channel_type: 'webhook',
  enabled: true,
  config: {
    webhook_url: '',
    webhook_method: 'POST',
    webhook_secret: '',
    webhook_headers: [],
    timeout_seconds: 30,
    retry_count: 3
  }
})

const formRules = {
  channel_name: [{ required: true, message: '请输入渠道名称' }],
  webhook_url: [
    { required: true, message: '请输入 Webhook URL' },
    { type: 'url', message: '请输入有效的 URL' }
  ]
}

function initForm() {
  if (props.channel) {
    form.channel_name = props.channel.channel_name || ''
    form.enabled = props.channel.enabled ?? true
    form.config = {
      webhook_url: props.channel.config?.webhook_url || '',
      webhook_method: props.channel.config?.webhook_method || 'POST',
      webhook_secret: props.channel.config?.webhook_secret || '',
      webhook_headers: props.channel.config?.webhook_headers || [],
      timeout_seconds: props.channel.config?.timeout_seconds || 30,
      retry_count: props.channel.config?.retry_count || 3
    }
  }
}

function addHeader() {
  form.config.webhook_headers.push({ key: '', value: '' })
}

function removeHeader(index) {
  form.config.webhook_headers.splice(index, 1)
}

async function handleSubmit() {
  try {
    await formRef.value?.validate()
  } catch {
    Message.warning('请检查表单填写')
    return
  }

  saving.value = true
  try {
    const data = {
      channel_type: 'webhook',
      channel_name: form.channel_name,
      enabled: form.enabled,
      config: form.config
    }

    if (props.channel?.id) {
      await updateChannel(props.channel.id, data)
    } else {
      await createChannel(data)
    }
    emit('success')
  } finally {
    saving.value = false
  }
}

async function handleTest() {
  testing.value = true
  testResult.value = null
  testResultVisible.value = true

  try {
    const channelData = {
      channel_type: 'webhook',
      channel_name: form.channel_name || '测试渠道',
      enabled: true,
      config: form.config
    }

    if (!props.channel?.id) {
      const created = await createChannel(channelData)
      if (created?.id) {
        testResult.value = await testChannel(created.id)
        await updateChannel(created.id, { enabled: false })
      }
    } else {
      testResult.value = await testChannel(props.channel.id)
    }
  } catch (e) {
    testResult.value = { success: false, error: '测试失败：' + (e.message || '未知错误') }
  } finally {
    testing.value = false
  }
}

onMounted(() => {
  initForm()
})
</script>

<style scoped>
.webhook-config-form {
  padding: 8px 0;
}
.form-hint {
  margin-left: 8px;
  font-size: 12px;
  color: #86909c;
}
.headers-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
.header-item {
  display: flex;
  align-items: center;
  gap: 8px;
}
.form-actions {
  margin-top: 24px;
  padding-top: 16px;
  border-top: 1px solid #e5e6eb;
  display: flex;
  justify-content: flex-end;
}
</style>
