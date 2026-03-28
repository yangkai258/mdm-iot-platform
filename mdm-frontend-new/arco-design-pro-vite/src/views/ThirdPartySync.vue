<template>
  <div class="container">
    <a-card>
      <template #title><icon-sync /> 第三方同步</template>
      <a-form layout="inline">
        <a-form-item label="平台"><a-select v-model="platform" style="width: 150px"><a-option value="wechat">微信</a-option></a-select></a-form-item>
        <a-form-item><a-button type="primary" @click="handleSync">同步</a-button></a-form-item>
      </a-form>
      <a-descriptions :column="2" bordered style="margin-top: 16px">
        <a-descriptions-item label="最后同步">2026-03-28 10:00</a-descriptions-item>
        <a-descriptions-item label="同步用户">256</a-descriptions-item>
      </a-descriptions>
    </a-card>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const platform = ref('wechat')
const handleSync = () => {}
</script>
<style scoped>.container { padding: 16px; }</style>
