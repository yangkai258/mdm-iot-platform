<template>
  <div class="container">
    <a-card>
      <template #title><icon-history /> 宠物记忆</template>
      <a-timeline>
        <a-timeline-item v-for="memory in memories" :key="memory.id">
          <p>{{ memory.content }}</p>
          <span class="time">{{ memory.time }}</span>
        </a-timeline-item>
      </a-timeline>
    </a-card>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const memories = ref([
  { id: 1, content: '今天学会了新动作', time: '2026-03-28' },
  { id: 2, content: '和主人一起去公园', time: '2026-03-27' }
])
</script>
<style scoped>
.time { font-size: 12px; color: #909399; }
</style>
