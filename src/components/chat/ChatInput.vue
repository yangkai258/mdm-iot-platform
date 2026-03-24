<template>
  <div class="chat-input">
    <div class="input-container">
      <a-input
        v-model="inputText"
        :placeholder="placeholder"
        :disabled="disabled"
        @press-enter="handleSend"
        @keydown.enter.meta="handleSend"
        @keydown.enter.ctrl="handleSend"
        :max-length="500"
        show-word-limit
      >
        <template #suffix>
          <a-button
            type="text"
            :disabled="!inputText.trim() || disabled"
            @click="handleSend"
            class="send-btn"
          >
            <template #icon>
              <icon-send />
            </template>
          </a-button>
        </template>
      </a-input>
      
      <a-button
        type="primary"
        :disabled="!inputText.trim() || disabled"
        @click="handleSend"
        class="send-btn-text"
      >
        发送
      </a-button>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { IconSend } from '@arco-design/web-vue/es/icon'

const props = defineProps({
  disabled: {
    type: Boolean,
    default: false
  },
  placeholder: {
    type: String,
    default: '输入消息...'
  }
})

const emit = defineEmits(['send'])

const inputText = ref('')

const handleSend = () => {
  if (!inputText.value.trim() || props.disabled) return
  
  emit('send', inputText.value.trim())
  inputText.value = ''
}
</script>

<style scoped>
.chat-input {
  padding: 12px 16px;
  background: #fff;
  border-top: 1px solid var(--color-fill-3);
}

.input-container {
  display: flex;
  gap: 12px;
  align-items: center;
}

.input-container :deep(.arco-input-wrapper) {
  flex: 1;
  border-radius: 20px;
  background: var(--color-fill-1);
}

.input-container :deep(.arco-input) {
  padding: 8px 12px;
}

.send-btn {
  padding: 4px;
  color: var(--color-primary);
}

.send-btn:disabled {
  color: var(--color-text-4);
}

.send-btn-text {
  border-radius: 20px;
  padding: 0 20px;
  height: 36px;
}
</style>
