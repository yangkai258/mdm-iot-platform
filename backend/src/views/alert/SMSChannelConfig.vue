<template>
  <div class="sms-config-form">
    <a-form :model="form" layout="vertical" :rules="formRules" ref="formRef">
      <!-- 基本信息 -->
      <a-divider>基本信息</a-divider>
      <a-form-item label="渠道名称" field="channel_name">
        <a-input v-model="form.channel_name" placeholder="如：生产环境短信通知" />
      </a-form-item>
      <a-form-item label="启用该渠道">
        <a-switch v-model="form.enabled" />
        <span class="form-hint">关闭后该渠道将不会发送任何通知</span>
      </a-form-item>

      <!-- 运营商选择 -->
      <a-divider>短信服务商配置</a-divider>
      <a-form-item label="短信服务商" field="sms_provider" required>
        <a-radio-group v-model="form.config.sms_provider">
          <a-radio value="aliyun">
            <span>阿里云</span>
            <img src="data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHdpZHRoPSIyNCIgaGVpZ2h0PSIyNCIgdmlld0JveD0iMCAwIDI0IDI0IiBmaWxsPSJub25lIiBzdHJva2U9IiMxNjUwZDgiIHN0cm9rZS13aWR0aD0iMiI+PHBhdGggZD0iTTEyIDEyLjVMMTQgMTRNMTEuNSA2LjVMMTQgMy41TDE2IDMuNUwxMyA3LjVMMTEuNSA2LjVNMTIgMy41TDEwLjUgNi41TDkuNSA2LjVMMTIgMy41Ii8+PC9zdmc+" alt="aliyun" style="width:16px;height:16px;vertical-align:middle;margin-left:4px;" />
          </a-radio>
          <a-radio value="tencent">
            <span>腾讯云</span>
            <img src="data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHdpZHRoPSIyNCIgaGVpZ2h0PSIyNCIgdmlld0JveD0iMCAwIDI0IDI0IiBmaWxsPSJub25lIiBzdHJva2U9IiMxNjUwZDgiIHN0cm9rZS13aWR0aD0iMiI+PGNpcmNsZSBjeD0iMTIiIGN5PSIxMiIgcj0iNSIvPjxwYXRoIGQ9Ik0xMiA3bDUgNS01IDUiLz48L3N2Zz4=" alt="tencent" style="width:16px;height:16px;vertical-align:middle;margin-left:4px;" />
          </a-radio>
          <a-radio value="huawei">华为云</a-radio>
        </a-radio-group>
      </a-form-item>

      <!-- 阿里云配置 -->
      <template v-if="form.config.sms_provider === 'aliyun'">
        <a-form-item label="AccessKey ID" field="access_key">
          <a-input v-model="form.config.access_key" placeholder="请输入阿里云 AccessKey ID" />
        </a-form-item>
        <a-form-item label="AccessKey Secret" field="secret_key">
          <a-input-password v-model="form.config.secret_key" placeholder="请输入阿里云 AccessKey Secret" />
        </a-form-item>
        <a-form-item label="短信签名名称" field="sign_name">
          <a-input v-model="form.config.sign_name" placeholder="如：MDM告警系统" />
        </a-form-item>
        <a-form-item label="短信模板 Code">
          <a-input v-model="form.config.template_code" placeholder="如：SMS_xxxxxxx" />
        </a-form-item>
      </template>

      <!-- 腾讯云配置 -->
      <template v-else-if="form.config.sms_provider === 'tencent'">
        <a-form-item label="SecretId" field="secret_id">
          <a-input v-model="form.config.secret_id" placeholder="请输入腾讯云 SecretId" />
        </a-form-item>
        <a-form-item label="SecretKey" field="secret_key">
          <a-input-password v-model="form.config.secret_key" placeholder="请输入腾讯云 SecretKey" />
        </a-form-item>
        <a-form-item label="短信签名名称" field="sign_name">
          <a-input v-model="form.config.sign_name" placeholder="如：MDM告警系统" />
        </a-form-item>
        <a-form-item label="短信模板 ID">
          <a-input v-model="form.config.template_id" placeholder="如：12xxx" />
        </a-form-item>
      </template>

      <!-- 华为云配置 -->
      <template v-else-if="form.config.sms_provider === 'huawei'">
        <a-form-item label="AppKey" field="app_key">
          <a-input v-model="form.config.app_key" placeholder="请输入华为云 AppKey" />
        </a-form-item>
        <a-form-item label="AppSecret" field="app_secret">
          <a-input-password v-model="form.config.app_secret" placeholder="请输入华为云 AppSecret" />
        </a-form-item>
        <a-form-item label="短信签名名称" field="sign_name">
          <a-input v-model="form.config.sign_name" placeholder="如：MDM告警系统" />
        </a-form-item>
        <a-form-item label="短信模板 ID">
          <a-input v-model="form.config.template_id" placeholder="请输入模板 ID" />
        </a-form-item>
      </template>

      <!-- 测试配置 -->
      <a-divider>测试配置</a-divider>
      <a-form-item label="测试手机号">
        <a-input v-model="testPhone" placeholder="用于测试短信的手机号，如：13800138000" />
      </a-form-item>
    </a-form>

    <!-- 操作按钮 -->
    <div class="form-actions">
      <a-space>
        <a-button @click="$emit('cancel')">「取消」</a-button>
        <a-button @click="handleTest" :loading="testing" :disabled="!form.config.access_key && !form.config.secret_id">
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
const testPhone = ref('')
const testResultVisible = ref(false)
const testResult = ref(null)

const { createChannel, updateChannel, testChannel } = useNotificationChannels()

const form = reactive({
  channel_name: '',
  channel_type: 'sms',
  enabled: true,
  config: {
    sms_provider: 'aliyun',
    access_key: '',
    secret_key: '',
    secret_id: '',
    app_key: '',
    app_secret: '',
    sign_name: '',
    template_code: '',
    template_id: ''
  }
})

const formRules = {
  channel_name: [{ required: true, message: '请输入渠道名称' }],
  sms_provider: [{ required: true, message: '请选择短信服务商' }]
}

function initForm() {
  if (props.channel) {
    form.channel_name = props.channel.channel_name || ''
    form.enabled = props.channel.enabled ?? true
    form.config = {
      ...form.config,
      ...props.channel.config
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
      channel_type: 'sms',
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
  if (!testPhone.value) {
    Message.warning('请输入测试手机号')
    return
  }

  testing.value = true
  testResult.value = null
  testResultVisible.value = true

  try {
    const channelData: Partial<NotificationChannel> = {
      channel_type: 'sms',
      channel_name: form.channel_name || '测试渠道',
      enabled: true,
      config: form.config
    }

    let tempId: number | null = null

    if (!props.channel?.id) {
      const created = await createChannel(channelData)
      if (created?.id) {
        tempId = created.id
        testResult.value = await testChannel(created.id, { recipient: testPhone.value })
        await updateChannel(created.id, { enabled: false } as any)
      }
    } else {
      testResult.value = await testChannel(props.channel.id, { recipient: testPhone.value })
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
.sms-config-form {
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
