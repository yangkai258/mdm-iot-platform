<template>
  <a-card class="pet-status-card" :bordered="false">
    <div class="pet-header">
      <div class="pet-avatar">
        <span class="pet-icon">{{ petIcon }}</span>
      </div>
      <div class="pet-info">
        <div class="pet-name">{{ petStatus.name || '鏈煡瀹犵墿' }}</div>
        <div class="pet-status">
          <a-badge
            :status="petStatus.is_online ? 'success' : 'default'"
            :text="petStatus.is_online ? '鍦ㄧ嚎' : '绂荤嚎'"
          />
        </div>
      </div>
    </div>

    <a-divider />

    <div class="status-list">
      <!-- 蹇冩儏 -->
      <div class="status-item">
        <div class="status-label">
          <span class="status-icon">馃槉</span>
          <span>蹇冩儏</span>
        </div>
        <div class="status-value">
          <a-progress
            :percent="petStatus.mood || 0"
            :stroke-width="8"
            :show-text="false"
            :color="getMoodColor(petStatus.mood)"
            size="small"
          />
          <span class="value-text">{{ petStatus.mood || 0 }}%</span>
        </div>
      </div>

      <!-- 鑳介噺 -->
      <div class="status-item">
        <div class="status-label">
          <span class="status-icon">鈿?/span>
          <span>鑳介噺</span>
        </div>
        <div class="status-value">
          <a-progress
            :percent="petStatus.energy || 0"
            :stroke-width="8"
            :show-text="false"
            color="#52c41a"
            size="small"
          />
          <span class="value-text">{{ petStatus.energy || 0 }}%</span>
        </div>
      </div>

      <!-- 楗ラタ -->
      <div class="status-item">
        <div class="status-label">
          <span class="status-icon">馃崠</span>
          <span>楗遍搴?/span>
        </div>
        <div class="status-value">
          <a-progress
            :percent="petStatus.hunger || 0"
            :stroke-width="8"
            :show-text="false"
            color="#faad14"
            size="small"
          />
          <span class="value-text">{{ petStatus.hunger || 0 }}%</span>
        </div>
      </div>
    </div>

    <a-divider />

    <div class="last-seen" v-if="petStatus.last_seen">
      <icon-clock-circle />
      <span>鏈€鍚庡湪绾? {{ formatLastSeen(petStatus.last_seen) }}</span>
    </div>

    <a-spin v-if="loading" class="loading-overlay" />
  </a-card>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  petStatus: {
    type: Object,
    default: () => ({
      name: '鏈煡瀹犵墿',
      type: 'cat',
      mood: 0,
      energy: 0,
      hunger: 0,
      is_online: false,
      last_seen: null
    })
  },
  loading: {
    type: Boolean,
    default: false
  }
})

const petIcon = computed(() => {
  const type = props.petStatus.type || 'cat'
  const icons = {
    cat: '馃惐',
    dog: '馃惗',
    bird: '馃惁',
    rabbit: '馃惏',
    hamster: '馃惞',
    fish: '馃悷',
    turtle: '馃悽',
    default: '馃惥'
  }
  return icons[type] || icons.default
})

function getMoodColor(mood) {
  if (mood >= 70) return '#52c41a'
  if (mood >= 40) return '#faad14'
  return '#ff4d4f'
}

function formatLastSeen(timestamp) {
  if (!timestamp) return '鏈煡'
  const date = new Date(timestamp)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  
  if (diff < 60000) return '鍒氬垰'
  if (diff < 3600000) return `${Math.floor(diff / 60000)} 鍒嗛挓鍓峘
  if (diff < 86400000) return `${Math.floor(diff / 3600000)} 灏忔椂鍓峘
  return date.toLocaleDateString('zh-CN')
}
</script>

<style scoped>
.pet-status-card {
  background: #fff;
  border-radius: 8px;
  position: relative;
}

.pet-header {
  display: flex;
  align-items: center;
  gap: 12px;
}

.pet-avatar {
  width: 56px;
  height: 56px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 28px;
  flex-shrink: 0;
}

.pet-info {
  flex: 1;
  min-width: 0;
}

.pet-name {
  font-size: 16px;
  font-weight: 600;
  color: var(--color-text-1);
  margin-bottom: 4px;
}

.pet-status {
  font-size: 12px;
}

.status-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.status-item {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.status-label {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: var(--color-text-2);
}

.status-icon {
  font-size: 14px;
}

.status-value {
  display: flex;
  align-items: center;
  gap: 8px;
}

.status-value :deep(.arco-progress) {
  flex: 1;
}

.value-text {
  font-size: 12px;
  color: var(--color-text-3);
  min-width: 36px;
  text-align: right;
}

.last-seen {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: var(--color-text-3);
}

.loading-overlay {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
}
</style>

