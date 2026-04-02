<template>
    <Breadcrumb :items="['Home','Console','']" />


  <div class="container">
    <a-card class="general-card" title="宠物控制台">
      <div class="chat-container">
        <div v-for="msg in messages" :key="msg.id" :class="['message', msg.role]">
          <div class="message-content">{{ msg.content }}</div>
        </div>
      </div>
      <a-input-search v-model="inputText" placeholder="输入消息..." @search="sendMessage" search-button class="input-area" />
    </a-card>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { Message } from '@arco-design/web-vue'
import Breadcrumb from '@/components/Breadcrumb.vue'

const inputText = ref('')
const messages = ref([])

const sendMessage = async () => {
  if (!inputText.value) return
  messages.value.push({ id: Date.now(), role: 'user', content: inputText.value })
  const text = inputText.value
  inputText.value = ''
  try {
    const res = await fetch('/api/v1/pet/conversations', {
      method: 'POST',
      headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token'), 'Content-Type': 'application/json' },
      body: JSON.stringify({ message: text })
    }).then(r => r.json())
    if (res.code === 0 && res.data) {
      messages.value.push({ id: Date.now() + 1, role: 'assistant', content: res.data.content || res.data.response || '收到消息' })
    }
  } catch { messages.value.push({ id: Date.now() + 1, role: 'assistant', content: '抱歉，服务器错误' }) }
}
</script>

<style scoped>
.chat-container { height: 400px; overflow-y: auto; margin-bottom: 16px; }
.message { margin: 8px 0; padding: 12px; border-radius: 8px; max-width: 80%; }
.message.user { background: #1890ff; color: white; margin-left: auto; }
.message.assistant { background: #f0f0f0; color: #333; margin-right: auto; }
.input-area { margin-top: 16px; }
</style>
