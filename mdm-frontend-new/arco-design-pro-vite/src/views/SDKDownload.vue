<template>
  <div class="container">
    <a-card>
      <template #title><icon-download /> SDK下载</template>
      <a-list>
        <a-list-item v-for="sdk in sdks" :key="sdk.id">
          <a-list-item-meta :title="sdk.name" :description="sdk.version" />
          <template #actions>
            <a-button @click="handleDownload(sdk)">下载</a-button>
          </template>
        </a-list-item>
      </a-list>
    </a-card>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const sdks = ref([{ id: 1, name: 'iOS SDK', version: 'v2.1.0' }, { id: 2, name: 'Android SDK', version: 'v2.1.0' }])
const handleDownload = (sdk) => {}
</script>
<style scoped>.container { padding: 16px; }</style>
