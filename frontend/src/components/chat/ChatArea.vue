<template>
  <div class="chat-area">
    <div class="chat-messages" ref="messagesContainer">
      <div v-if="!messages.length && !loading" class="empty-state">
        <icon-message :size="48" />
        <p>开始和宠物聊天吧</p>
      </div>
      
      <TransitionGroup name="message" tag="div" class="messages-list">
        <ChatMessage
          v-for="msg in messages"
          :key="msg.id"
          :message="msg"
        />
      </TransitionGroup>
      
      <div v-if="loading" class="loading-indicator">
        <a-spin size="small" />
        <span>宠物正在思考...</span>
      </div>
    </div>
    
    <ChatInput
      @send="handleSend"
      :disabled="loading"
    />
  </div>
</template>

<script setup>
import { ref, watch, nextTick } from 'vue'
import ChatMessage from './ChatMessage.vue'
import ChatInput from './ChatInput.vue'

const props = defineProps({
  messages: {
    type: Array,
    default: () => []
  },
  loading: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['send'])

const messagesContainer = ref(null)

const handleSend = (content) => {
  emit('send', content)
}

// 自动滚动到底部
watch(() => props.messages.length, async () => {
  await nextTick()
  scrollToBottom()
})

watch(() => props.loading, async (newVal, oldVal) => {
  if (newVal === false && oldVal === true) {
    await nextTick()
    scrollToBottom()
  }
})

function scrollToBottom() {
  if (messagesContainer.value) {
    messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
  }
}
</script>

<style scoped>
.chat-area {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: #fff;
  border-radius: 8px;
  overflow: hidden;
}

.chat-messages {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
  display: flex;
  flex-direction: column;
}

.empty-state {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: var(--color-text-3);
  gap: 12px;
}

.empty-state p {
  margin: 0;
  font-size: 14px;
}

.messages-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.loading-indicator {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px;
  color: var(--color-text-3);
  font-size: 14px;
}

/* 消息动画 */
.message-enter-active {
  transition: all 0.3s ease;
}

.message-enter-from {
  opacity: 0;
  transform: translateY(20px);
}
</style>
