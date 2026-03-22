<template>
  <a-layout class="pet-conversations">
    <a-layout-sider v-model:collapsed="collapsed" :trigger="null" collapsible>
      <div class="logo">
        <span v-if="!collapsed">MDM 控制台</span>
      </div>
      <a-menu v-model:selectedKeys="selectedKeys" theme="dark" mode="inline" @click="handleMenuClick">
        <a-menu-item key="dashboard">
          <span>设备大盘</span>
        </a-menu-item>
        <a-menu-item key="ota">
          <span>OTA 固件</span>
        </a-menu-item>
        <a-menu-item key="pet">
          <span>宠物配置</span>
        </a-menu-item>
        <a-menu-item key="status">
          <span>设备状态</span>
        </a-menu-item>
      </a-menu>
    </a-layout-sider>

    <a-layout>
      <a-layout-header class="header">
        <div class="header-left">
          <a-button type="text" @click="collapsed = !collapsed">
            <span v-if="collapsed">☰</span>
            <span v-else>✕</span>
          </a-button>
        </div>
        <div class="header-title">
          <span>宠物对话</span>
        </div>
        <div class="header-right">
          <a-select v-model="selectedPetId" placeholder="选择宠物" style="width: 180px" allow-search @change="onPetChange">
            <a-option v-for="pet in pets" :key="pet.pet_id" :value="pet.pet_id">
              {{ pet.pet_name }}
            </a-option>
          </a-select>
        </div>
      </a-layout-header>

      <a-layout-content class="content">
        <a-spin :spinning="loading">
          <!-- 未选择宠物 -->
          <a-empty v-if="!selectedPetId" description="请先选择宠物" style="margin-top: 80px" />

          <div v-else class="chat-layout">
            <!-- 左侧：宠物信息和历史记录 -->
            <div class="sidebar-panel">
              <!-- 宠物信息卡片 -->
              <a-card class="pet-info-card" size="small">
                <a-descriptions :column="1" size="small">
                  <a-descriptions-item label="宠物名称">
                    <a-space>
                      <a-avatar :size="24" :style="{ backgroundColor: currentPet?.avatar_color || '#165dff' }">
                        {{ currentPet?.pet_name?.charAt(0) || '?' }}
                      </a-avatar>
                      {{ currentPet?.pet_name || '-' }}
                    </a-space>
                  </a-descriptions-item>
                  <a-descriptions-item label="宠物性格">
                    {{ getPersonalityText(currentPet?.personality) }}
                  </a-descriptions-item>
                  <a-descriptions-item label="绑定设备">
                    {{ currentPet?.device_id || '未绑定' }}
                  </a-descriptions-item>
                </a-descriptions>
              </a-card>

              <!-- 对话历史列表 -->
              <a-card class="history-card" size="small">
                <template #title>
                  <span>对话历史</span>
                </template>
                <a-space direction="vertical" style="width: 100%" :size="4">
                  <a-button
                    v-for="conv in conversations"
                    :key="conv.conversation_id"
                    :type="selectedConvId === conv.conversation_id ? 'primary' : 'outline'"
                    :size="size"
                    style="width: 100%; text-align: left; justify-content: flex-start"
                    @click="selectConversation(conv)"
                  >
                    <a-space>
                      <span>{{ conv.title || '对话 ' + conv.conversation_id }}</span>
                      <a-badge v-if="conv.unread_count > 0" :count="conv.unread_count" :max-count="99" />
                    </a-space>
                  </a-button>
                  <a-empty v-if="conversations.length === 0" description="暂无对话记录" />
                </a-space>
                <div style="margin-top: 12px">
                  <a-button type="primary" style="width: 100%" @click="startNewConversation">
                    新建对话
                  </a-button>
                </div>
              </a-card>

              <!-- 统计信息 -->
              <a-card size="small">
                <a-statistic title="对话总数" :value="conversations.length" />
              </a-card>
            </div>

            <!-- 右侧：对话内容 -->
            <div class="chat-panel">
              <!-- 无选中对话 -->
              <a-empty v-if="!selectedConvId" description="请选择或新建对话" />

              <div v-else class="chat-area">
                <!-- 对话标题栏 -->
                <div class="chat-header">
                  <span style="font-weight: 500; font-size: 14px">{{ currentConvTitle }}</span>
                  <a-space>
                    <a-button size="small" @click="clearConversation">清空</a-button>
                    <a-button size="small" type="primary" @click="exportConversation">导出</a-button>
                  </a-space>
                </div>

                <!-- 消息列表 -->
                <div class="message-list" ref="messageListRef">
                  <div
                    v-for="msg in messages"
                    :key="msg.id"
                    :class="['message-item', msg.direction]"
                  >
                    <div class="message-bubble">
                      <div class="message-text">{{ msg.content }}</div>
                      <div class="message-meta">
                        <span>{{ msg.time }}</span>
                        <span v-if="msg.status">{{ msg.status }}</span>
                      </div>
                    </div>
                  </div>
                  <div v-if="messages.length === 0" style="text-align: center; padding: 40px; color: #999">
                    开始和宠物对话吧
                  </div>
                </div>

                <!-- 发送区域 -->
                <div class="send-area">
                  <a-input-group compact style="display: flex; gap: 8px">
                    <a-select v-model="sendType" style="width: 100px">
                      <a-option value="text">文本</a-option>
                      <a-option value="command">指令</a-option>
                    </a-select>
                    <a-input
                      v-model="sendText"
                      :placeholder="sendType === 'command' ? '输入指令，如 feed, play, walk' : '输入对话内容'"
                      @press-enter="handleSend"
                      style="flex: 1"
                    />
                    <a-button type="primary" :loading="sending" @click="handleSend">发送</a-button>
                  </a-input-group>
                  <!-- 快捷指令 -->
                  <div class="quick-sends">
                    <a-space wrap style="gap: 6px">
                      <a-tag
                        v-for="cmd in quickCommands"
                        :key="cmd"
                        color="arcoblue"
                        style="cursor: pointer"
                        @click="quickSend(cmd)"
                      >
                        {{ cmd }}
                      </a-tag>
                    </a-space>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </a-spin>
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>

<script setup>
import { ref, reactive, computed, onMounted, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { Message } from '@arco-design/web-vue'
import { getPets, getConversations, sendCommand } from '../../api/pet.js'

const router = useRouter()

const collapsed = ref(false)
const selectedKeys = ref(['pet'])
const loading = ref(false)
const sending = ref(false)
const pets = ref([])
const selectedPetId = ref('')
const selectedConvId = ref('')
const sendText = ref('')
const sendType = ref('text')
const messageListRef = ref(null)
const size = ref('small')

const conversations = ref([])
const messages = ref([])
let msgIdCounter = 1

const quickCommands = ['feed', 'play', 'walk', 'rest', 'dance', 'speak', 'sit']

const currentPet = computed(() => pets.value.find(p => p.pet_id === selectedPetId.value) || null)
const currentConvTitle = computed(() => {
  const conv = conversations.value.find(c => c.conversation_id === selectedConvId.value)
  return conv?.title || `对话 ${selectedConvId.value}`
})

const handleMenuClick = ({ key }) => {
  if (key === 'dashboard') router.push('/dashboard')
  else if (key === 'ota') router.push('/ota')
  else if (key === 'pet') router.push('/pet')
  else if (key === 'status') router.push('/status')
}

const getPersonalityText = (personality) => {
  const texts = { 1: '活泼好动', 2: '温顺安静', 3: '好奇宝宝', 4: '独立自主', 5: '粘人依赖' }
  return texts[personality] || '未知'
}

const loadPets = async () => {
  loading.value = true
  try {
    const res = await getPets()
    if (res.code === 0) {
      pets.value = res.data.list || []
    }
  } catch (err) {
    pets.value = [
      { pet_id: 'PET001', pet_name: '小橘', personality: 1, device_id: 'DEV001', avatar_color: '#f6ad55' },
      { pet_id: 'PET002', pet_name: '布丁', personality: 2, device_id: '', avatar_color: '#fc8181' },
      { pet_id: 'PET003', pet_name: '豆豆', personality: 3, device_id: 'DEV003', avatar_color: '#68d391' }
    ]
  } finally {
    loading.value = false
  }
}

const loadConversations = async () => {
  if (!selectedPetId.value) return
  try {
    const res = await getConversations(selectedPetId.value)
    if (res.code === 0) {
      conversations.value = res.data?.list || []
    }
  } catch (err) {
    conversations.value = [
      { conversation_id: 'C001', title: '小橘的日常', unread_count: 2 },
      { conversation_id: 'C002', title: '喂食记录', unread_count: 0 },
      { conversation_id: 'C003', title: '健康检查', unread_count: 1 }
    ]
  }
}

const loadMessages = async (convId) => {
  if (!selectedPetId.value || !convId) return
  try {
    const res = await getConversations(selectedPetId.value, { conversation_id: convId })
    if (res.code === 0) {
      messages.value = res.data?.messages || []
    }
  } catch (err) {
    messages.value = [
      { id: 1, direction: 'in', content: '你好呀！我是小橘~', time: '10:00', status: '' },
      { id: 2, direction: 'out', content: '小橘，今天感觉怎么样？', time: '10:01', status: '已发送' },
      { id: 3, direction: 'in', content: '我很开心！主人今天陪我玩耍了 🎾', time: '10:02', status: '' }
    ]
  }
}

const selectConversation = (conv) => {
  selectedConvId.value = conv.conversation_id
  loadMessages(conv.conversation_id)
  conv.unread_count = 0
}

const startNewConversation = () => {
  const newId = `C${Date.now()}`
  conversations.value.unshift({
    conversation_id: newId,
    title: `新对话 ${newId}`,
    unread_count: 0
  })
  selectedConvId.value = newId
  messages.value = []
  const petName = currentPet.value?.pet_name || '宠物'
  addMessage('in', `你好！我是${petName}，有什么可以帮你的吗？`)
}

const handleSend = async () => {
  if (!sendText.value.trim()) {
    Message.warning('请输入内容')
    return
  }
  if (!selectedConvId.value) {
    Message.warning('请先选择或新建对话')
    return
  }

  const text = sendText.value.trim()
  const direction = 'out'
  const now = new Date().toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })

  // 添加用户消息
  addMessage(direction, text, now)
  sendText.value = ''

  sending.value = true
  try {
    await sendCommand(selectedPetId.value, {
      conversation_id: selectedConvId.value,
      content: text,
      type: sendType.value
    })
    addMessage('in', `收到指令: ${text}，正在处理...`, now)
  } catch (err) {
    // 模拟回复
    setTimeout(() => {
      const replies = [
        '好的，我收到了！🐾',
        '收到指令啦，正在执行~',
        '明白！马上处理~',
        '收到！(*^▽^*)',
        '指令已接收，等我一下哦~'
      ]
      const reply = replies[Math.floor(Math.random() * replies.length)]
      const replyTime = new Date().toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
      addMessage('in', reply, replyTime)
    }, 800)
  } finally {
    sending.value = false
  }
}

const quickSend = (cmd) => {
  if (!selectedConvId.value) {
    startNewConversation()
  }
  sendType.value = 'command'
  sendText.value = cmd
  handleSend()
}

const addMessage = (direction, content, time) => {
  if (!time) {
    time = new Date().toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
  }
  messages.value.push({
    id: msgIdCounter++,
    direction,
    content,
    time,
    status: direction === 'out' ? '已发送' : ''
  })
  nextTick(() => {
    if (messageListRef.value) {
      messageListRef.value.scrollTop = messageListRef.value.scrollHeight
    }
  })
}

const clearConversation = () => {
  messages.value = []
  Message.success('对话已清空')
}

const exportConversation = () => {
  const lines = messages.value.map(m => {
    const prefix = m.direction === 'out' ? '[我]' : '[宠物]'
    return `${prefix} ${m.time}: ${m.content}`
  }).join('\n')
  const blob = new Blob([lines], { type: 'text/plain' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `conversation_${selectedConvId.value}_${Date.now()}.txt`
  a.click()
  URL.revokeObjectURL(url)
  Message.success('对话已导出')
}

const onPetChange = (petId) => {
  selectedConvId.value = ''
  messages.value = []
  loadConversations()
}

onMounted(() => {
  loadPets()
  loadConversations()
})
</script>

<style scoped>
.pet-conversations {
  min-height: 100vh;
}

.logo {
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  font-weight: bold;
  color: #fff;
}

.header {
  background: #fff;
  padding: 0 16px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
}

.header-left, .header-right {
  display: flex;
  align-items: center;
}

.header-title {
  font-size: 16px;
  font-weight: 500;
}

.content {
  margin: 16px;
}

.chat-layout {
  display: flex;
  gap: 16px;
  height: calc(100vh - 140px);
}

.sidebar-panel {
  width: 280px;
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.chat-panel {
  flex: 1;
  min-width: 0;
}

.chat-area {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: #fff;
  border-radius: 4px;
  overflow: hidden;
}

.chat-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  border-bottom: 1px solid #e5e7eb;
  background: #fafafa;
}

.message-list {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.message-item {
  display: flex;
}

.message-item.out {
  justify-content: flex-end;
}

.message-item.in {
  justify-content: flex-start;
}

.message-bubble {
  max-width: 70%;
  padding: 8px 12px;
  border-radius: 8px;
  word-break: break-word;
}

.message-item.out .message-bubble {
  background: #165dff;
  color: #fff;
}

.message-item.in .message-bubble {
  background: #f2f3f5;
  color: #333;
}

.message-text {
  font-size: 14px;
  line-height: 1.5;
}

.message-meta {
  display: flex;
  gap: 8px;
  font-size: 11px;
  opacity: 0.7;
  margin-top: 4px;
  justify-content: flex-end;
}

.send-area {
  padding: 12px 16px;
  border-top: 1px solid #e5e7eb;
}

.quick-sends {
  margin-top: 8px;
}

.pet-info-card {
  margin-bottom: 0;
}

.history-card {
  flex: 1;
  overflow-y: auto;
}
</style>
