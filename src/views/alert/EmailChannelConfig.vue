<template>
  <div class="email-config-form">
    <a-form :model="form" layout="vertical" :rules="formRules" ref="formRef">
      <!-- 基本信息 -->
      <a-divider>基本信息</a-divider>
      <a-form-item label="渠道名称" field="channel_name">
        <a-input v-model="form.channel_name" placeholder="如：生产环境邮件通知" />
      </a-form-item>
      <a-form-item label="启用该渠道">
        <a-switch v-model="form.enabled" />
        <span class="form-hint">关闭后该渠道将不会发送任何通知</span>
      </a-form-item>

      <!-- SMTP 服务器配置 -->
      <a-divider>SMTP 服务器配置</a-divider>
      <a-form-item label="SMTP 服务器地址" field="smtp_host" required>
        <a-input v-model="form.config.smtp_host" placeholder="smtp.example.com" />
      </a-form-item>
      <a-form-item label="端口" field="smtp_port">
        <a-select v-model="form.config.smtp_port" placeholder="选择端口">
          <a-option :value="25">25 (非加密)</a-option>
          <a-option :value="465">465 (SSL)</a-option>
          <a-option :value="587">587 (TLS)</a-option>
        </a-select>
      </a-form-item>
      <a-form-item label="用户名" field="smtp_user">
        <a-input v-model="form.config.smtp_user" placeholder="alert@example.com" />
      </a-form-item>
      <a-form-item label="密码/授权码" field="smtp_password">
        <a-input-password v-model="form.config.smtp_password" placeholder="请输入密码或授权码" />
      </a-form-item>

      <!-- 安全设置 -->
      <a-divider>安全设置</a-divider>
      <a-form-item label="加密方式">
        <a-radio-group v-model="form.config.encryption">
          <a-radio value="tls">TLS</a-radio>
          <a-radio value="ssl">SSL</a-radio>
          <a-radio value="none">不加密</a-radio>
        </a-radio-group>
      </a-form-item>

      <!-- 发件人配置 -->
      <a-divider>发件人配置</a-divider>
      <a-form-item label="发件人邮箱" field="smtp_from" required>
        <a-input v-model="form.config.smtp_from" placeholder="alert@example.com" />
      </a-form-item>
      <a-form-item label="发件人名称">
        <a-input v-model="form.config.smtp_from_name" placeholder="MDM 告警系统" />
      </a-form-item>
      <a-form-item label="默认收件人（多个用逗号分隔）" field="smtp_to">
        <a-input v-model="form.config.smtp_to" placeholder="admin@example.com,ops@example.com" />
      </a-form-item>

      <!-- 测试配置 -->
      <a-divider>测试配置</a-divider>
      <a-form-item label="测试收件人">
        <a-input v-model="testRecipient" placeholder="用于测试邮件的收件地址" />
        <span class="form-hint">不填则使用默认收件人</span>
      </a-form-item>
    </a-form>

    <!-- 操作按钮 -->
    <div class="form-actions">
      <a-space>
        <a-button @click="$emit('cancel')">「取消」</a-button>
        <a-button @click="handleTest" :loading="testing" :disabled="!form.config.smtp_host">
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

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import { useNotificationChannels } from '@/composables/useNotification'
import type { NotificationChannel } from '@/api/notification'

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
const testRecipient = ref('')
const testResultVisible = ref(false)
const testResult = ref(null)

const { createChannel, updateChannel, testChannel } = useNotificationChannels()

const form = reactive({
  channel_name: '',
  channel_type: 'email',
  enabled: true,
  config: {
    smtp_host: '',
    smtp_port: 587,
    smtp_user: '',
    smtp_password: '',
    smtp_from: '',
    smtp_from_name: '',
    smtp_to: '',
    encryption: 'tls'
  }
})

const formRules = {
  channel_name: [{ required: true, message: '请输入渠道名称' }],
  'config.smtp_host': [{ required: true, message: '请输入 SMTP 服务器' }],
  'config.smtp_from': [
    { required: true, message: '请输入发件人邮箱' },
    { type: 'email', message: '请输入有效的邮箱地址' }
  ]
}

function initForm() {
  if (props.channel) {
    form.channel_name = props.channel.channel_name || ''
    form.enabled = props.channel.enabled ?? true
    form.config = {
      smtp_host: props.channel.config?.smtp_host || '',
      smtp_port: props.channel.config?.smtp_port || 587,
      smtp_user: props.channel.config?.smtp_user || '',
      smtp_password: props.channel.config?.smtp_password || '',
      smtp_from: props.channel.config?.smtp_from || '',
      smtp_from_name: props.channel.config?.smtp_from_name || '',
      smtp_to: props.channel.config?.smtp_to || '',
      encryption: props.channel.config?.encryption || 'tls'
    }
  }
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
    const data: Partial<NotificationChannel> = {
      channel_type: 'email',
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
    const channelData: Partial<NotificationChannel> = {
      channel_type: 'email',
      channel_name: form.channel_name || '测试渠道',
      enabled: true,
      config: form.config
    }

    let tempId: number | null = null

    if (!props.channel?.id) {
      // 临时创建渠道进行测试
      const created = await createChannel(channelData)
      if (created?.id) {
        tempId = created.id
        testResult.value = await testChannel(created.id, {
          recipient: testRecipient.value || form.config.smtp_to
        })
        // 删除临时渠道
        await updateChannel(created.id, { enabled: false } as any)
      }
    } else {
      testResult.value = await testChannel(props.channel.id, {
        recipient: testRecipient.value || form.config.smtp_to
      })
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
.email-config-form {
  padding: 8px 0;
}
.form-hint {
  margin-left: 8px;
  font-size: 12px;
  color: #86909c;
}
.form-actions {
  margin-top: 24px;
  padding-top: 16px;
  border-top: 1px solid #e5e6eb;
  display: flex;
  justify-content: flex-end;
}
</style>
