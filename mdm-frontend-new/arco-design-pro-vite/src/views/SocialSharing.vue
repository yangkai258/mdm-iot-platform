<template>
  <div class="container">
    <a-card>
      <template #title><icon-share /> 社交分享</template>
      <a-list>
        <a-list-item v-for="share in shares" :key="share.id">
          <a-list-item-meta :title="share.title" :description="share.time" />
        </a-list-item>
      </a-list>
    </a-card>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const shares = ref([{ id: 1, title: '分享到微信', time: '2026-03-28' }])
</script>
<style scoped>.container { padding: 16px; }</style>
