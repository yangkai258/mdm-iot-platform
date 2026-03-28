<template>
  <div class="container">
    <a-card>
      <template #title><icon-image /> 家庭相册</template>
      <a-button type="primary" @click="handleUpload"><template #icon><icon-upload /></template>上传照片</a-button>
      <a-image-preview-group>
        <a-row :gutter="[8, 8]" style="margin-top: 16px">
          <a-col :span="6" v-for="photo in photos" :key="photo.id">
            <a-image :src="photo.url" width="100%" />
          </a-col>
        </a-row>
      </a-image-preview-group>
    </a-card>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const photos = ref([
  { id: 1, url: 'https://picsum.photos/200?random=1' },
  { id: 2, url: 'https://picsum.photos/200?random=2' }
])
const handleUpload = () => {}
</script>
<style scoped>.container { padding: 16px; }</style>
