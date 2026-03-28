<template>
  <div class="container">
    <a-card>
      <template #title><icon-share-node /> 网络拓扑</template>
      <div class="topology-view">
        <svg width="800" height="500" viewBox="0 0 800 500">
          <circle cx="400" cy="250" r="30" fill="#409EFF" />
          <text x="400" y="255" text-anchor="middle" fill="white">Cloud</text>
          <circle v-for="(node, i) in nodes" :key="i" :cx="node.x" :cy="node.y" r="20" fill="#67C23A" />
        </svg>
      </div>
    </a-card>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const nodes = ref([{ x: 200, y: 150 }, { x: 600, y: 150 }, { x: 200, y: 350 }, { x: 600, y: 350 }])
</script>
<style scoped>.container { padding: 16px; }
.topology-view { background: #f5f7fa; border-radius: 8px; display: flex; justify-content: center; padding: 16px; }
</style>
