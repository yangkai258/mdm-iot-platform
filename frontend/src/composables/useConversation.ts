import { ref } from 'vue'
import { 
  getConversations, 
  getConversationMessages, 
  createConversation,
  deleteConversation 
} from '@/api/pet'

export function useConversation(deviceId) {
  const conversations = ref([])
  const currentMessages = ref([])
  const activeConversationId = ref(null)
  const loading = ref(false)
  const error = ref<Error | null>(null)

  async function loadConversations() {
    if (!deviceId.value) return
    
    loading.value = true
    error.value = null
    
    try {
      const res = await getConversations(deviceId.value)
      if (res.data) {
        conversations.value = res.data
      }
    } catch (e) {
      error.value = e as Error
      console.error('Failed to load conversations:', e)
    } finally {
      loading.value = false
    }
  }

  async function loadMessages(conversationId) {
    loading.value = true
    error.value = null
    activeConversationId.value = conversationId
    
    try {
      const res = await getConversationMessages(conversationId)
      if (res.data) {
        currentMessages.value = res.data
      }
    } catch (e) {
      error.value = e as Error
      console.error('Failed to load messages:', e)
    } finally {
      loading.value = false
    }
  }

  async function create() {
    if (!deviceId.value) return

    loading.value = true
    error.value = null
    
    try {
      const res = await createConversation(deviceId.value)
      if (res.data) {
        activeConversationId.value = res.data.id
        currentMessages.value = []
        await loadConversations()
        return res.data
      }
    } catch (e) {
      error.value = e as Error
      console.error('Failed to create conversation:', e)
      throw e
    } finally {
      loading.value = false
    }
  }

  async function remove(conversationId) {
    loading.value = true
    error.value = null
    
    try {
      await deleteConversation(deviceId.value, conversationId)
      
      if (activeConversationId.value === conversationId) {
        activeConversationId.value = null
        currentMessages.value = []
      }
      
      await loadConversations()
    } catch (e) {
      error.value = e as Error
      console.error('Failed to delete conversation:', e)
      throw e
    } finally {
      loading.value = false
    }
  }

  function addMessage(message) {
    currentMessages.value.push({
      id: Date.now(),
      ...message,
      timestamp: new Date().toISOString()
    })
  }

  function clearMessages() {
    currentMessages.value = []
    activeConversationId.value = null
  }

  return {
    conversations,
    currentMessages,
    activeConversationId,
    loading,
    error,
    loadConversations,
    loadMessages,
    create,
    remove,
    addMessage,
    clearMessages
  }
}
