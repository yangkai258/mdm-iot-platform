<template>
  <div class="container">
    <a-card>
      <template #title><icon-shopping-cart /> 应用市场</template>
      <template #extra><a-button type="primary" @click="handleUpload"><template #icon><icon-upload /></template>上传应用</a-button></template>
      <a-list>
        <a-list-item v-for="app in apps" :key="app.id">
          <a-list-item-meta :title="app.name" :description="app.description" />
          <template #actions><a-button size="small">安装</a-button></template>
        </a-list-item>
      </a-list>
    </a-card>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const apps = ref([{ id: 1, name: '智能提醒', description: '基于AI的智能提醒应用' }])
const handleUpload = () => {}
</script>
<style scoped>.container { padding: 16px; }</style>
