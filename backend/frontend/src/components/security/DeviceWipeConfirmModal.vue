<template>
  <a-modal
    :visible="visible"
    title="数据擦除确认"
    width="520px"
    @cancel="handleCancel"
    @before-ok="handleConfirm"
    :ok-loading="confirming"
    :ok-button-props="{ status: 'danger', disabled: !canConfirm }"
  >
    <div class="wipe-content">
      <!-- 警告提示 -->
      <a-alert class="danger-alert">
        <template #title>
          <span class="danger-title">危险操作警告</span>
        </template>
        <template #content>
          此操作将永久清除设备上的所有数据，且 <strong>无法恢复</strong>。
          请确认您了解此操作的后果。
        </template>
      </a-alert>

      <!-- 设备信息 -->
      <div class="device-summary" v-if="device">
        <div class="summary-title">即将擦除的设备</div>
        <div class="summary-row">
          <span class="summary-label">设备名称</span>
          <span class="summary-value">{{ device.name }}</span>
        </div>
        <div class="summary-row">
          <span class="summary-label">设备ID</span>
          <span class="summary-value mono">{{ device.id }}</span>
        </div>
      </div>

      <!-- 数据范围说明 -->
      <div class="wipe-scope">
        <div class="scope-title">将被擦除的数据范围</div>
        <ul class="scope-list">
          <li>设备上所有用户数据</li>
          <li>设备配置信息</li>
          <li>本地缓存文件</li>
          <li>设备证书（需重新注册）</li>
        </ul>
      </div>

      <!-- 确认输入 -->
      <div class="confirm-section">
        <div class="confirm-label">请输入 <span class="confirm-code">CONFIRM</span> 以确认操作</div>
        <a-input
          v-model="confirmText"
          placeholder="请输入 CONFIRM"
          style="width: 200px"
          @press-enter="handleConfirm"
        />
      </div>

      <!-- 密码验证 -->
      <div class="password-section">
        <div class="password-label">操作人员密码验证</div>
        <a-input
          v-model="password"
          type="password"
          placeholder="请输入当前账号密码"
          style="width: 240px"
          allow-clear
        />
      </div>
    </div>
  </a-modal>
</template>

<script setup>
import { ref, computed, watch } from 'vue'

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  device: {
    type: Object,
    default: null
  }
})

const emit = defineEmits(['update:visible', 'confirm'])

const confirming = ref(false)
const confirmText = ref('')
const password = ref('')

const canConfirm = computed(() => {
  return confirmText.value.trim().toUpperCase() === 'CONFIRM' && password.value.length >= 6
})

watch(() => props.visible, (val) => {
  if (val) {
    confirmText.value = ''
    password.value = ''
  }
})

function handleCancel() {
  emit('update:visible', false)
}

async function handleConfirm(done) {
  if (!canConfirm.value) {
    done(false)
    return
  }
  confirming.value = true
  try {
    await emit('confirm', { password: password.value })
    done(true)
  } catch (e) {
    done(false)
  } finally {
    confirming.value = false
  }
}
</script>

<style scoped>
.wipe-content {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.danger-alert {
  border-color: #f53f3f;
  background: #fff1f0;
}

.danger-title {
  color: #f53f3f;
  font-weight: 600;
}

.device-summary {
  padding: 12px;
  background: var(--color-fill-lightest);
  border-radius: 4px;
}

.summary-title {
  font-size: 13px;
  font-weight: 600;
  margin-bottom: 8px;
  color: var(--color-text-1);
}

.summary-row {
  display: flex;
  gap: 12px;
  margin-bottom: 6px;
}

.summary-label {
  color: var(--color-text-3);
  font-size: 13px;
  min-width: 80px;
}

.summary-value {
  font-size: 13px;
  word-break: break-all;
}

.mono {
  font-family: monospace;
  font-size: 12px;
}

.wipe-scope {
  padding: 12px;
  border: 1px solid var(--color-border);
  border-radius: 4px;
}

.scope-title {
  font-size: 13px;
  font-weight: 600;
  margin-bottom: 8px;
  color: var(--color-text-1);
}

.scope-list {
  margin: 0;
  padding-left: 20px;
  color: var(--color-text-2);
  font-size: 13px;
  line-height: 1.8;
}

.scope-list li {
  color: #f53f3f;
}

.confirm-section {
  display: flex;
  align-items: center;
  gap: 12px;
}

.confirm-label {
  font-size: 13px;
  color: var(--color-text-2);
}

.confirm-code {
  font-family: monospace;
  font-weight: 600;
  color: var(--color-text-1);
  background: var(--color-fill-lightest);
  padding: 2px 6px;
  border-radius: 3px;
}

.password-section {
  display: flex;
  align-items: center;
  gap: 12px;
}

.password-label {
  font-size: 13px;
  color: var(--color-text-2);
  min-width: 120px;
}
</style>
