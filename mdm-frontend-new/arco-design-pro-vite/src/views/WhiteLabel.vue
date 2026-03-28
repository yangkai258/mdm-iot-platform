<template>
  <div class="container">
    <a-card>
      <template #title><icon-palette /> 白标定制</template>
      <a-form :model="formData" layout="vertical" style="max-width: 500px">
        <a-form-item label="品牌名称"><a-input v-model="formData.brandName" /></a-form-item>
        <a-form-item label="Logo"><a-upload /></a-form-item>
        <a-form-item label="主题色"><a-color-picker v-model="formData.primaryColor" /></a-form-item>
        <a-form-item><a-button type="primary" @click="handleSave">保存</a-button></a-form-item>
      </a-form>
    </a-card>
  </div>
</template>
<script setup>
import { ref, reactive } from 'vue'
const formData = reactive({ brandName: '', primaryColor: '#165qff' })
const handleSave = () => {}
</script>
<style scoped>.container { padding: 16px; }</style>
