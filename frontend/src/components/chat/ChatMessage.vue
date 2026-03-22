<template>
  <div :class="['chat-message', message.role]">
    <!-- AI 头像 -->
    <div v-if="message.role === 'assistant'" class="avatar assistant-avatar">
      <span class="avatar-icon">🐱</span>
    </div>
    
    <!-- 消息内容 -->
    <div class="message-content-wrapper">
      <div class="message-bubble">
        {{ message.content }}
      </div>
      <div v-if="message.timestamp" class="message-time">
        {{ formatTime(message.timestamp) }}
      </div>
    </div>
    
    <!-- 用户头像 -->
    <div v-if="message.role === 'user'" class="avatar user-avatar">
      <span class="avatar-icon">👤</span>
    </div>
  </div>
</template>

<script setup>
defineProps({
  message: {
    type: Object,
    required: true,
    default: () => ({
      id: '',
      role: 'user',
      content: '',
      timestamp: null
    })
  }
})

function formatTime(timestamp) {
  if (!timestamp) return ''
  const date = new Date(timestamp)
  const hours = date.getHours().toString().padStart(2, '0')
  const minutes = date.getMinutes().toString().padStart(2, '0')
  return `${hours}:${minutes}`
}
</script>

<style scoped>
.chat-message {
  display: flex;
  gap: 12px;
  max-width: 80%;
}

.chat-message.user {
  align-self: flex-end;
  flex-direction: row-reverse;
}

.chat-message.assistant {
  align-self: flex-start;
}

.avatar {
  flex-shrink: 0;
  width: 36px;
  height: 36px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
}

.assistant-avatar {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.user-avatar {
  background: linear-gradient(135deg, #1890ff 0%, #0050b3 100%);
}

.avatar-icon {
  font-size: 18px;
}

.message-content-wrapper {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.message-bubble {
  padding: 12px 16px;
  border-radius: 16px;
  font-size: 14px;
  line-height: 1.5;
  word-break: break-word;
}

.chat-message.user .message-bubble {
  background: linear-gradient(135deg, #1890ff 0%, #0050b3 100%);
  color: #fff;
  border-bottom-right-radius: 4px;
}

.chat-message.assistant .message-bubble {
  background: #f5f5f5;
  color: var(--color-text-1);
  border-bottom-left-radius: 4px;
}

.message-time {
  font-size: 11px;
  color: var(--color-text-3);
  padding: 0 4px;
}

.chat-message.user .message-time {
  text-align: right;
}
</style>
