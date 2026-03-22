<template>
  <a-card class="quick-actions-card" :bordered="false">
    <template #title>
      <div class="card-title">
        <icon-thunderbolt />
        <span>快捷指令</span>
      </div>
    </template>
    
    <div class="actions-grid">
      <a-button
        v-for="action in actions"
        :key="action.name"
        :class="['action-btn', action.category]"
        @click="handleAction(action.name)"
      >
        <span class="action-icon">{{ action.icon }}</span>
        <span class="action-name">{{ action.label }}</span>
      </a-button>
    </div>
  </a-card>
</template>

<script setup>
import { Message } from '@arco-design/web-vue'

const emit = defineEmits(['action'])

const actions = [
  { name: 'come', label: '过来', icon: '👋', category: 'movement' },
  { name: 'sit', label: '坐下', icon: '🧘', category: 'movement' },
  { name: 'spin', label: '转圈', icon: '🔄', category: 'movement' },
  { name: 'sleep', label: '睡觉', icon: '😴', category: 'life' },
  { name: 'play', label: '玩耍', icon: '🎾', category: 'life' },
  { name: 'eat', label: '投喂', icon: '🍖', category: 'life' }
]

const handleAction = (actionName) => {
  emit('action', actionName)
}
</script>

<style scoped>
.quick-actions-card {
  background: #fff;
  border-radius: 8px;
}

.card-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  font-weight: 500;
}

.actions-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
}

.action-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 16px 8px;
  border-radius: 12px;
  border: 1px solid var(--color-fill-3);
  background: var(--color-fill-1);
  cursor: pointer;
  transition: all 0.2s ease;
}

.action-btn:hover {
  border-color: var(--color-primary);
  background: var(--color-primary-light-1);
}

.action-btn:active {
  transform: scale(0.95);
}

.action-icon {
  font-size: 24px;
}

.action-name {
  font-size: 12px;
  color: var(--color-text-2);
}

.action-btn.movement:hover .action-name {
  color: var(--color-primary);
}

.action-btn.life:hover .action-name {
  color: #52c41a;
}
</style>
