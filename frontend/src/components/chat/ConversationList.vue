<template>
  <div class="conversation-list">
    <!-- 新建会话按钮 -->
    <div class="list-header">
      <a-button type="primary" long @click="handleCreate">
        <template #icon>
          <icon-plus />
        </template>
        新建会话
      </a-button>
    </div>

    <!-- 会话列表 -->
    <div class="list-content">
      <div v-if="loading" class="loading-state">
        <a-spin size="small" />
      </div>
      
      <div v-else-if="!conversations.length" class="empty-state">
        <icon-file />
        <p>暂无历史会话</p>
      </div>

      <div v-else class="conversations">
        <div
          v-for="conv in conversations"
          :key="conv.id"
          :class="['conversation-item', { active: conv.id === activeConversationId }]"
          @click="handleSelect(conv.id)"
        >
          <div class="conv-icon">
            <icon-message />
          </div>
          <div class="conv-info">
            <div class="conv-title">{{ conv.title || '新会话' }}</div>
            <div class="conv-meta">
              <span class="conv-time">{{ formatTime(conv.updated_at) }}</span>
              <a-button
                type="text"
                size="mini"
                class="delete-btn"
                @click.stop="handleDelete(conv.id)"
              >
                <template #icon>
                  <icon-delete />
                </template>
              </a-button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import { getConversations, deleteConversation } from '@/api/pet'

const props = defineProps({
  deviceId: {
    type: String,
    default: ''
  },
  activeConversationId: {
    type: String,
    default: null
  }
})

const emit = defineEmits(['select', 'create'])

const conversations = ref([])
const loading = ref(false)

// 加载会话列表
async function loadConversations() {
  if (!props.deviceId) return
  
  loading.value = true
  try {
    const res = await getConversations(props.deviceId)
    if (res.data) {
      conversations.value = res.data
    }
  } catch (error) {
    console.error('加载会话列表失败:', error)
  } finally {
    loading.value = false
  }
}

watch(() => props.deviceId, () => {
  loadConversations()
}, { immediate: true })

function handleSelect(conversationId) {
  emit('select', conversationId)
}

function handleCreate() {
  emit('create')
}

async function handleDelete(conversationId) {
  try {
    await deleteConversation(props.deviceId, conversationId)
    Message.success('删除成功')
    await loadConversations()
  } catch (error) {
    Message.error('删除失败')
    console.error('删除会话失败:', error)
  }
}

function formatTime(timestamp) {
  if (!timestamp) return ''
  const date = new Date(timestamp)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  
  if (diff < 60000) return '刚刚'
  if (diff < 3600000) return `${Math.floor(diff / 60000)} 分钟前`
  if (diff < 86400000) return `${Math.floor(diff / 3600000)} 小时前`
  if (diff < 604800000) return `${Math.floor(diff / 86400000)} 天前`
  return date.toLocaleDateString('zh-CN')
}
</script>

<style scoped>
.conversation-list {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.list-header {
  padding-bottom: 16px;
  border-bottom: 1px solid var(--color-fill-3);
}

.list-content {
  flex: 1;
  overflow-y: auto;
  padding: 12px 0;
}

.loading-state {
  display: flex;
  justify-content: center;
  padding: 24px;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 48px 24px;
  color: var(--color-text-3);
  gap: 12px;
}

.empty-state p {
  margin: 0;
  font-size: 13px;
}

.conversations {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.conversation-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.conversation-item:hover {
  background: var(--color-fill-2);
}

.conversation-item.active {
  background: var(--color-primary-light-1);
}

.conv-icon {
  width: 36px;
  height: 36px;
  border-radius: 8px;
  background: var(--color-fill-2);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--color-text-3);
  flex-shrink: 0;
}

.conversation-item.active .conv-icon {
  background: var(--color-primary);
  color: #fff;
}

.conv-info {
  flex: 1;
  min-width: 0;
}

.conv-title {
  font-size: 14px;
  color: var(--color-text-1);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  margin-bottom: 4px;
}

.conv-meta {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.conv-time {
  font-size: 12px;
  color: var(--color-text-3);
}

.delete-btn {
  opacity: 0;
  transition: opacity 0.2s;
}

.conversation-item:hover .delete-btn {
  opacity: 1;
}
</style>
