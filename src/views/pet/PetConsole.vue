<template>
  <div class="pet-console">
    
    <div class="page-content"><a-card>
        <div class="chat-container">
          <div v-for="msg in messages" :key="msg.id" :class="['message', msg.role]">
            <div class="message-content">{{ msg.content }}</div>
          </div>
        </div>
        <a-input-search
          v-model="inputText"
          placeholder="杈撳叆娑堟伅..."
          @search="sendMessage"
          search-button
          class="input-area"
        />
      </a-card></div>
  </div>
</template>

<script>
export default {
  name: 'PetConsole',
  data() {
    return {
      inputText: '',
      messages: []
    }
  },
  methods: {
    sendMessage() {
      if (!this.inputText) return
      this.messages.push({
        id: Date.now(),
        role: 'user',
        content: this.inputText
      })
      this.inputText = ''
    }
  }
}
</script>

<style scoped>
.chat-container {
  height: 400px;
  overflow-y: auto;
  margin-bottom: 16px;
}
.message {
  margin: 8px 0;
  padding: 12px;
  border-radius: 8px;
  max-width: 80%;
}
.message.user {
  background: #1890ff;
  color: white;
  margin-left: auto;
}
.message.assistant {
  background: #f0f0f0;
  color: #333;
  margin-right: auto;
}
.input-area {
  margin-top: 16px;
}
</style>
