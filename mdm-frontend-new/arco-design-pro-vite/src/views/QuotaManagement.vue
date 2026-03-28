<template>
  <div class="container">
    <a-card>
      <template #title><icon-check-circle /> 配额管理</template>
      <a-form layout="vertical" style="max-width: 500px">
        <a-form-item label="API配额(每天)"><a-input-number v-model="form.apiQuota" /></a-form-item>
        <a-form-item label="存储配额"><a-input-number v-model="form.storageQuota" suffix="GB" /></a-form-item>
        <a-form-item><a-button type="primary" @click="handleSave">保存</a-button></a-form-item>
      </a-form>
    </a-card>
  </div>
</template>
<script setup>
import { ref, reactive } from 'vue'
const form = reactive({ apiQuota: 10000, storageQuota: 10 })
const handleSave = () => {}
</script>
<style scoped>.container { padding: 16px; }</style>
