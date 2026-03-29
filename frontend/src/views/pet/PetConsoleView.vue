<template>
  <div class="page-container pet-console-view">
    <!-- 面包屑 -->
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>宠物管理</a-breadcrumb-item>
      <a-breadcrumb-item>宠物控制台</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 页面标题 -->
    <div class="page-header">
      <h2>宠物控制台</h2>
    </div>

    <!-- 顶部工具栏 -->
    <div class="console-toolbar">
      <div class="toolbar-left">
        <a-select
          v-model="selectedDeviceId"
          placeholder="选择设备"
          style="width: 200px"
          @change="handleDeviceChange"
        >
          <a-option v-for="device in deviceList" :key="device.device_id" :value="device.device_id">
            {{ device.name || device.device_id }}
          </a-option>
        </a-select>
      </div>
      <div class="toolbar-right">
        <a-space>
          <a-button @click="showHistoryDrawer = true">
            <template #icon>
              <icon-history />
            </template>
            历史
          </a-button>
          <a-button @click="showSettingsModal = true">
            <template #icon>
              <icon-settings />
            </template>
            设置
          </a-button>
        </a-space>
      </div>
    </div>

    <!-- 主控制台区域 -->
    <div class="console-main">
      <!-- 左侧：聊天区域 (70%) -->
      <div class="chat-section">
        <ChatArea
          :messages="messages"
          :loading="messageLoading"
          @send="handleSendMessage"
        />
      </div>

      <!-- 右侧：宠物状态卡片 (30%) -->
      <div class="status-section">
        <PetStatusCard
          :pet-status="petStatus"
          :loading="statusLoading"
        />
        <QuickActions
          @action="handleQuickAction"
        />
      </div>
    </div>

    <!-- 历史会话抽屉 -->
    <a-drawer
      v-model:visible="showHistoryDrawer"
      title="历史会话"
      placement="left"
      :width="360"
    >
      <ConversationList
        :device-id="selectedDeviceId"
        :active-conversation-id="activeConversationId"
        @select="handleConversationSelect"
        @create="handleCreateConversation"
      />
    </a-drawer>

    <!-- 设置弹窗 -->
    <PetSettingsModal
      v-model:visible="showSettingsModal"
      :device-id="selectedDeviceId"
      :settings="petSettings"
      @save="handleSaveSettings"
    />
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import ChatArea from '@/components/chat/ChatArea.vue'
import ChatMessage from '@/components/chat/ChatMessage.vue'
import ChatInput from '@/components/chat/ChatInput.vue'
import PetStatusCard from '@/components/pet/PetStatusCard.vue'
import QuickActions from '@/components/pet/QuickActions.vue'
import ConversationList from '@/components/chat/ConversationList.vue'
import PetSettingsModal from '@/components/pet/PetSettingsModal.vue'
import { getPetStatus, sendMessage, executeAction, updatePetSettings, getConversations, createConversation } from '@/api/pet'
import { useWebSocket } from '@/composables/useWebSocket'

// 状态
const selectedDeviceId = ref('')
const deviceList = ref([])
const messages = ref([])
const messageLoading = ref(false)
const statusLoading = ref(false)
const petStatus = ref({
  name: '未知宠物',
  type: 'cat',
  mood: 0,
  energy: 0,
  hunger: 0,
  is_online: false,
  last_seen: null
})
const petSettings = ref({
  name: '',
  type: 'cat',
  mood: 50,
  energy: 50,
  hunger: 50
})

// UI 状态
const showHistoryDrawer = ref(false)
const showSettingsModal = ref(false)
const activeConversationId = ref(null)

// WebSocket
const wsUrl = ref('')
const { connect: connectWs, disconnect: disconnectWs, send: sendWs, isConnected } = useWebSocket()

// 初始化
onMounted(async () => {
  await loadDevices()
  if (deviceList.value.length > 0) {
    selectedDeviceId.value = deviceList.value[0].device_id
    await loadPetStatus()
    connectWebSocket()
  }
})

onUnmounted(() => {
  disconnectWs()
})

// 加载设备列表
async function loadDevices() {
  try {
    const res = await getPetStatus()
    if (res.data) {
      deviceList.value = res.data.devices || []
    }
  } catch (error) {
    console.error('加载设备列表失败:', error)
  }
}

// 加载宠物状态
async function loadPetStatus() {
  if (!selectedDeviceId.value) return
  
  statusLoading.value = true
  try {
    const res = await getPetStatus(selectedDeviceId.value)
    if (res.data) {
      petStatus.value = res.data
    }
  } catch (error) {
    console.error('加载宠物状态失败:', error)
  } finally {
    statusLoading.value = false
  }
}

// 设备选择变更
async function handleDeviceChange(deviceId) {
  selectedDeviceId.value = deviceId
  await loadPetStatus()
  connectWebSocket()
  messages.value = []
  activeConversationId.value = null
}

// 连接 WebSocket
function connectWebSocket() {
  if (!selectedDeviceId.value) return
  
  disconnectWs()
  
  // 根据当前域名动态构建 WebSocket URL
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  const host = window.location.host
  wsUrl.value = `${protocol}//${host}/api/v1/ws/pets/${selectedDeviceId.value}`
  
  connectWs(wsUrl.value, {
    onMessage: handleWsMessage,
    onConnect: () => {
      console.log('WebSocket connected')
    },
    onDisconnect: () => {
      console.log('WebSocket disconnected')
    }
  })
}

// 处理 WebSocket 消息
function handleWsMessage(data) {
  try {
    const msg = JSON.parse(data)
    if (msg.type === 'status_update') {
      petStatus.value = { ...petStatus.value, ...msg.data }
    } else if (msg.type === 'ai_response') {
      messages.value.push({
        id: Date.now(),
        role: 'assistant',
        content: msg.data.content,
        timestamp: new Date().toISOString()
      })
    }
  } catch (error) {
    console.error('解析 WebSocket 消息失败:', error)
  }
}

// 发送消息
async function handleSendMessage(content) {
  if (!content.trim() || !selectedDeviceId.value) return

  // 添加用户消息
  messages.value.push({
    id: Date.now(),
    role: 'user',
    content: content,
    timestamp: new Date().toISOString()
  })

  messageLoading.value = true
  try {
    const res = await sendMessage(selectedDeviceId.value, {
      content,
      conversation_id: activeConversationId.value
    })
    
    if (res.data) {
      // 添加 AI 回复
      if (res.data.ai_response) {
        messages.value.push({
          id: Date.now() + 1,
          role: 'assistant',
          content: res.data.ai_response,
          timestamp: new Date().toISOString()
        })
      }
      activeConversationId.value = res.data.conversation_id
    }
  } catch (error) {
    Message.error('发送消息失败')
    console.error('发送消息失败:', error)
  } finally {
    messageLoading.value = false
  }
}

// 快捷操作
async function handleQuickAction(action) {
  if (!selectedDeviceId.value) return

  try {
    await executeAction(selectedDeviceId.value, { action })
    Message.success(`执行成功: ${action}`)
    await loadPetStatus()
  } catch (error) {
    Message.error(`执行失败: ${action}`)
    console.error('执行动作失败:', error)
  }
}

// 保存设置
async function handleSaveSettings(settings) {
  if (!selectedDeviceId.value) return

  try {
    await updatePetSettings(selectedDeviceId.value, settings)
    Message.success('设置已保存')
    await loadPetStatus()
  } catch (error) {
    Message.error('保存设置失败')
    console.error('保存设置失败:', error)
  }
}

// 选择历史会话
async function handleConversationSelect(conversationId) {
  activeConversationId.value = conversationId
  showHistoryDrawer.value = false
  // 加载会话消息
  // messages.value = await getConversationMessages(conversationId)
}

// 创建新会话
async function handleCreateConversation() {
  if (!selectedDeviceId.value) return

  try {
    const res = await createConversation(selectedDeviceId.value)
    if (res.data) {
      activeConversationId.value = res.data.id
      messages.value = []
      showHistoryDrawer.value = false
      Message.success('已创建新会话')
    }
  } catch (error) {
    Message.error('创建会话失败')
    console.error('创建会话失败:', error)
  }
}
</script>

<style scoped>
.pet-console-view {
  display: flex;
  flex-direction: column;
  height: calc(100vh - 120px);
}

.breadcrumb {
  margin-bottom: 12px;
}

.page-header {
  margin-bottom: 16px;
}

.page-header h2 {
  font-size: 20px;
  font-weight: 600;
  color: var(--color-text-1);
  margin: 0;
}

.console-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background: #fff;
  border-radius: 8px;
  margin-bottom: 16px;
}

.console-main {
  display: flex;
  gap: 16px;
  flex: 1;
  min-height: 0;
}

.chat-section {
  flex: 7;
  min-width: 0;
}

.status-section {
  flex: 3;
  display: flex;
  flex-direction: column;
  gap: 16px;
  min-width: 280px;
  max-width: 360px;
}
</style>
